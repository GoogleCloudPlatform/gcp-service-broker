// Copyright 2018 the Service Broker Project Authors.
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
	dbTypeProp     = "db.type"
	dbPathProp     = "db.path"

	DbTypeMysql   = "mysql"
	DbTypeSqlite3 = "sqlite3"
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

	viper.BindEnv(dbTypeProp, "DB_TYPE")
	viper.SetDefault(dbTypeProp, DbTypeMysql)

	viper.BindEnv(dbPathProp, "DB_PATH")
}

// pulls db credentials from the environment, connects to the db, and returns the db connection
func SetupDb(logger lager.Logger) *gorm.DB {
	dbType := viper.GetString(dbTypeProp)
	var db *gorm.DB
	var err error
	// if provided, use database injected by CF via VCAP_SERVICES environment variable
	if err := UseVcapServices(); err != nil {
		logger.Error("Invalid VCAP_SERVICES environment variable", err)
		os.Exit(1)
	}
	switch dbType {
	default:
		logger.Error("Database Setup", fmt.Errorf("Invalid database type %q, valid types are: sqlite3 and mysql", dbType))
		os.Exit(1)
	case DbTypeMysql:
		db, err = setupMysqlDb(logger)
	case DbTypeSqlite3:
		db, err = setupSqlite3Db(logger)
	}

	if err != nil {
		logger.Error("Database Setup", err)
		os.Exit(1)
	}

	return db
}

func setupSqlite3Db(logger lager.Logger) (*gorm.DB, error) {
	dbPath := viper.GetString(dbPathProp)
	if dbPath == "" {
		return nil, fmt.Errorf("You must set a database path when using SQLite3 databases")
	}

	logger.Info("WARNING: DO NOT USE SQLITE3 IN PRODUCTION!")

	return gorm.Open(DbTypeSqlite3, dbPath)
}

func setupMysqlDb(logger lager.Logger) (*gorm.DB, error) {
	// connect to database
	dbHost := viper.GetString(dbHostProp)
	dbUsername := viper.GetString(dbUserProp)
	dbPassword := viper.GetString(dbPassProp)

	if dbPassword == "" || dbHost == "" || dbUsername == "" {
		return nil, errors.New("DB_HOST, DB_USERNAME and DB_PASSWORD are required environment variables.")
	}

	dbPort := viper.GetString(dbPortProp)
	dbName := viper.GetString(dbNameProp)

	tlsStr, err := generateTlsStringFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error generating TLS string from env: %s", err)
	}

	logger.Info("Connecting to MySQL Database", lager.Data{
		"host": dbHost,
		"port": dbPort,
		"name": dbName,
	})

	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local%v", dbUsername, dbPassword, dbHost, dbPort, dbName, tlsStr)
	return gorm.Open(DbTypeMysql, connStr)
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
