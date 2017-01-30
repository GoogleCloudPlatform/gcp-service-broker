#!/usr/bin/env bash

set -e

export GOPATH=${PWD}

cd ${PWD}/src/gcp-service-broker-src/brokerapi/brokers
go test