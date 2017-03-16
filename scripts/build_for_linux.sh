#!/bin/bash
echo $0


cd `dirname $0`/..
rd=`pwd`
echo $rd

mkdir -p $GOPATH/bin

echo "building in $rd"

docker run --rm -v "$rd":/go/src/gcp-service-broker -v "$GOPATH/bin":/go/bin -w /go/src/gcp-service-broker -e GOOS=linux -e GOARCH=amd64 golang:1.8 go install -v
