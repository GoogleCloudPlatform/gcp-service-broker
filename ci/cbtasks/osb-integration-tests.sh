#!/usr/bin/env bash
set -e

# Setup Environment
GODIR=$GOPATH/src/github.com/GoogleCloudPlatform/gcp-service-broker
mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD $GODIR
cd $GODIR

# Run Tests
go build

export SECURITY_USER_NAME=user
export SECURITY_USER_PASSWORD=password
export PORT=8080

export ROOT_SERVICE_ACCOUNT_JSON=$(cat ci/cbtasks/account.json)
export DB_PATH=test.sqlite3
export DB_TYPE=sqlite3

echo "Starting server"
./gcp-service-broker migrate
./gcp-service-broker serve &

sleep 5
./gcp-service-broker client run-examples
