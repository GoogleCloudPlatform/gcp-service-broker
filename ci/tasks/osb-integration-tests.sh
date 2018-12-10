#!/usr/bin/env bash

set -e

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
