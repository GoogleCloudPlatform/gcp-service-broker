#!/usr/bin/env bash

set -e

PROGNAME=gcp-service-broker
GODIR=github.com/GoogleCloudPlatform/gcp-service-broker

mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD $GOPATH/src/$GODIR


go get github.com/onsi/ginkgo/ginkgo

cd "${GOPATH}/src/${GODIR}"

ginkgo -r -race -skipPackage=integration,db_service .
