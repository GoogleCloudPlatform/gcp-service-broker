// Copyright the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////
//

package db_service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

const (
	caCertProp     = "db.ca.cert"
	clientCertProp = "db.client.cert"
	clientKeyProp  = "db.client.key"
	dbHostProp     = "db.host"
	dbUserProp     = "db.user"
	dbPassProp     = "db.password"
	dbPortProp     = "db.port"
	dbNameProp     = "db.name"
)

func init() {
	viper.BindEnv(caCertProp, "CA_CERT")
	viper.BindEnv(clientCertProp, "CLIENT_CERT")
	viper.BindEnv(clientKeyProp, "CLIENT_KEY")

	viper.BindEnv(dbHostProp, "DB_HOST")
	viper.BindEnv(dbUserProp, "DB_USERNAME")
	viper.BindEnv(dbPassProp, "DB_PASSWORD")

	viper.BindEnv(dbPortProp, "DB_PORT")
	viper.SetDefault(dbPortProp, "3306")
	viper.BindEnv(dbNameProp, "DB_NAME")
	viper.SetDefault(dbNameProp, "servicebroker")
}

// pulls db credentials from the environment, connects to the db, runs migrations, and returns the db connection
func SetupDb(logger lager.Logger) *gorm.DB {
	// connect to database
	dbHost := viper.GetString(dbHostProp)
	dbUsername := viper.GetString(dbUserProp)
	dbPassword := viper.GetString(dbPassProp)

	if dbPassword == "" || dbHost == "" || dbUsername == "" {
		panic("DB_HOST, DB_USERNAME and DB_PASSWORD are required environment variables.")
	}

	dbPort := viper.GetString(dbPortProp)
	dbName := viper.GetString(dbNameProp)

	tlsStr, err := generateTlsStringFromEnv()
	if err != nil {
		logger.Error("Error generating TLS string from env", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local%v", dbUsername, dbPassword, dbHost, dbPort, dbName, tlsStr)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		logger.Error("Error connecting to db", err)
		os.Exit(1)
	}

	return db
}

func generateTlsStringFromEnv() (string, error) {
	caCert := viper.GetString(caCertProp)
	clientCertStr := viper.GetString(clientCertProp)
	clientKeyStr := viper.GetString(clientKeyProp)
	tlsStr := "&tls=custom"

	// make sure ssl is set up for this connection
	if caCert != "" && clientCertStr != "" && clientKeyStr != "" {

		rootCertPool := x509.NewCertPool()

		if ok := rootCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			return "", fmt.Errorf("Error appending cert: %s", errors.New(""))
		}
		clientCert := make([]tls.Certificate, 0, 1)

		certs, err := tls.X509KeyPair([]byte(clientCertStr), []byte(clientKeyStr))
		if err != nil {
			return "", fmt.Errorf("Error parsing cert pair: %s", err)
		}
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
	} else {
		tlsStr = ""
	}

	return tlsStr, nil
}
