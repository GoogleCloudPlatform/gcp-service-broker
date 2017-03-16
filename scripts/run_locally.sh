#!/bin/bash
json=$1
./drop_tables.sh
export GOPATH=`cd ../.. && pwd`
export ROOT_SERVICE_ACCOUNT_JSON=`cat $1`
export SECURITY_USER_NAME=admin
export SECURITY_USER_PASSWORD=admin
export DB_HOST=localhost
export DB_USERNAME=gcp-service-broker
export DB_PASSWORD=qwerty
echo $GOPATH
echo $ROOT_SERVICE_ACCOUNT_JSON
go run server.go
