#!/usr/bin/env sh
set -e

echo "Installing zip"
apk update
apk add zip

export CURRENT_VERSION="$(cat /artifacts/metadata/version)"
export REVISION="$(cat /artifacts/metadata/revision)"

# Bundle up the output
mkdir /artifacts/staging
mkdir /artifacts/bundle

echo "Staging files from the source"
cp CHANGELOG.md /artifacts/staging/
cp OSDF*.txt /artifacts/staging/
ls -la /artifacts/staging/

echo "Staging files from metadata"
cp /artifacts/metadata/version /artifacts/staging/
cp -r /artifacts/metadata/docs /artifacts/staging/
ls -la /artifacts/staging/

echo "Staging server binaries"
mkdir -p /artifacts/staging/servers
cp /artifacts/servers/* /artifacts/staging/servers
ls -la /artifacts/staging/

echo "Staging client binaries"
mkdir -p /artifacts/staging/clients
cp /artifacts/client-darwin/* /artifacts/staging/clients
cp /artifacts/client-linux/* /artifacts/staging/clients
cp /artifacts/client-windows/* /artifacts/staging/clients

echo "Creating release from /artifacts/staging"
ls -la /artifacts/staging/
zip /artifacts/bundle/gcp-service-broker-$CURRENT_VERSION-$REVISION.zip -r /artifacts/staging/*

echo Root
ls -la /

echo "/artifacts"
ls -la /artifacts

