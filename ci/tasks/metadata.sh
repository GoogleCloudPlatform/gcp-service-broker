#!/bin/sh

set -e

echo "Installing dependencies"
apk update
apk add git

echo "Generating metadata"
mkdir -p metadata/docs

git --git-dir=gcp-service-broker/.git rev-parse HEAD > metadata/revision
./compiled-broker/gcp-service-broker version > metadata/version
./compiled-broker/gcp-service-broker generate tile > metadata/tile.yml
./compiled-broker/gcp-service-broker generate use > metadata/manifest.yml
./compiled-broker/gcp-service-broker generate customization > metadata/docs/customization.md
./compiled-broker/gcp-service-broker generate use --destination-dir="docs/"
