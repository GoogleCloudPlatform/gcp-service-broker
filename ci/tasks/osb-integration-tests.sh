#!/usr/bin/env bash

set -e

# Setup Environment
GODIR=$GOPATH/src/github.com/GoogleCloudPlatform/gcp-service-broker
mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/src/gcp-service-broker $GODIR
cd $GODIR

# Run Tests
go build

export SECURITY_USER_NAME=user
export SECURITY_USER_PASSWORD=password
export PORT=8080

echo "Starting server"
./gcp-service-broker migrate
./gcp-service-broker serve &

sleep 5
./gcp-service-broker client run-examples
