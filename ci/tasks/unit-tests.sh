#!/usr/bin/env bash

set -e

ls -lah

PROGNAME=gcp-service-broker
GODIR=github.com/GoogleCloudPlatform/gcp-service-broker

mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/src/gcp-service-broker $GOPATH/src/$GODIR


go get github.com/onsi/ginkgo/ginkgo

cd "${GOPATH}/src/${GODIR}"

ls -lah

ginkgo -r -race -skipPackage=integration,db_service .
