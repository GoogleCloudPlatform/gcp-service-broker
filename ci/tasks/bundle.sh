#!/bin/sh
set -e

echo "Installing zip"
apk update
apk add zip

export CURRENT_VERSION="$(cat metadata/version)"

# Bundle up the output
mkdir staging

echo "Staging files from the source"
cp gcp-service-broker/CHANGELOG.md staging/
cp gcp-service-broker/OSDF*.txt staging/

echo "Staging files from metadata"
cp metadata/version staging/
cp -r metadata/docs staging/

echo "Staging server binaries"
mkdir -p staging/servers
cp tiles/* staging/servers

echo "Staging client binaries"
# TODO(josephlewis42) pack up cross-compiled binaries for windows/darwin/linux

echo "Creating release"
zip bundle/gcp-service-broker-$CURRENT_VERSION.zip -r staging/*
