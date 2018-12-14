#!/bin/sh

set -e

# set up alpine
apk update
apk add ca-certificates

# use the compiled broker
cd compiled-broker

# Setup Environment
export SECURITY_USER_NAME=user
export SECURITY_USER_PASSWORD=password
export PORT=8080

echo "Running brokerpak tests"
./gcp-service-broker pak test

echo "Starting server"
./gcp-service-broker migrate
./gcp-service-broker serve &

sleep 5
./gcp-service-broker client run-examples
