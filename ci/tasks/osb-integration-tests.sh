#!/usr/bin/env bash

set -e

# Setup Environment
export SECURITY_USER_NAME=user
export SECURITY_USER_PASSWORD=password
export PORT=8080

echo "Running brokerpak tests"
./compiled-broker/gcp-service-broker pak test

echo "Starting server"
./compiled-broker/gcp-service-broker migrate
./compiled-broker/gcp-service-broker serve &

sleep 5
./compiled-broker/gcp-service-broker client run-examples
