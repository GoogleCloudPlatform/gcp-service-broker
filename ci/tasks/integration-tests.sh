#!/usr/bin/env bash

set -e

export GOPATH="${PWD}"
export PATH=$PATH:$GOPATH/bin

go get github.com/onsi/ginkgo/ginkgo

cd "${GOPATH}/src/gcp-service-broker"

ginkgo -race integration_tests db_service