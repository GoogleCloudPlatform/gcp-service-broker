#!/usr/bin/env bash

set -e

# Setup Environment
GODIR=$GOPATH/src/github.com/GoogleCloudPlatform/gcp-service-broker
mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/src/gcp-service-broker $GODIR
cd $GODIR

# Run Tests
go get github.com/onsi/ginkgo/ginkgo
ginkgo -race integration_tests db_service
