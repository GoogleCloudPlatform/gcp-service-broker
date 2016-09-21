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
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// pulls db credentials from the environment, connects to the db, runs migrations, and returns the db connection
func SetupDb(logger lager.Logger) *gorm.DB {
	// connect to database
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	tlsStr, err := generateTlsStringFromEnv()
	if err != nil {
		logger.Error("Error generating TLS string from env", err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/servicebroker?charset=utf8&parseTime=True&loc=Local%v", dbUsername, dbPassword, dbHost, dbPort, tlsStr)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		logger.Error("Error connecting to db", err)
		os.Exit(1)
	}

	return db
}

func generateTlsStringFromEnv() (string, error) {
	caCert64 := os.Getenv("CA_CERT_B64")
	clientCert64 := os.Getenv("CLIENT_CERT_B64")
	clientKey64 := os.Getenv("CLIENT_KEY_B64")
	tlsStr := "&tls=custom"

	// make sure ssl is set up for this connection
	if caCert64 != "" && clientCert64 != "" && clientKey64 != "" {
		caCert, err := base64.StdEncoding.DecodeString(caCert64)
		if err != nil {
			return "", fmt.Errorf("Error decoding ca cert: %s", err)
		}

		clientCertStr, err := base64.StdEncoding.DecodeString(clientCert64)
		if err != nil {
			return "", fmt.Errorf("Error decoding client cert: %s", err)
		}

		clientKeyStr, err := base64.StdEncoding.DecodeString(clientKey64)
		if err != nil {
			return "", fmt.Errorf("Error decoding client key: %s", err)
		}

		rootCertPool := x509.NewCertPool()

		if ok := rootCertPool.AppendCertsFromPEM(caCert); !ok {
			return "", fmt.Errorf("Error appending cert: %s", errors.New(""))
		}
		clientCert := make([]tls.Certificate, 0, 1)

		certs, err := tls.X509KeyPair(clientCertStr, clientKeyStr)
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
