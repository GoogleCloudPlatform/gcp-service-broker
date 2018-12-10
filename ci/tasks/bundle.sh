#!/bin/sh
set -e

export CURRENT_VERSION="$(cat metadata/version)"

# Bundle up the output
mkdir staging

# Files from the source directory
cp gcp-service-broker/CHANGELOG.md staging/
cp gcp-service-broker/OSDF*.txt staging/

# Files from metadata
cp metadata/version staing/
cp -r metadata/docs staging/

# Server Side
mkdir -p staging/servers
cp tiles/* staging/servers

# Client Side
# TODO(josephlewis42) pack up cross-compiled binaries for windows/darwin/linux

zip bundle/gcp-service-broker-$CURRENT_VERSION.zip -r staging/*
