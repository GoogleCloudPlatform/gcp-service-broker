#!/usr/bin/env bash

set -e

export GOPATH=${PWD}/gcp-service-broker-src
export PATH=${GOPATH}/bin:$PATH

cd ${PWD}/gcp-service-broker-src/brokerapi/brokers
go test