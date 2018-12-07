#!/usr/bin/env bash

set -e

# Setup Environment
GODIR=$GOPATH/src/github.com/GoogleCloudPlatform/gcp-service-broker
mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/src/gcp-service-broker $GODIR
cd $GODIR

# run unit tests
go test -cover ./...

# build the broker
go build

# test brokerpaks e2e
./gcp-service-broker pak test
