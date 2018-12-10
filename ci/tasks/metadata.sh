#!/bin/sh

set -e

mkdir -p metadata/docs

./compiled-broker/gcp-service-broker version > metadata/version
./compiled-broker/gcp-service-broker generate tile > metadata/tile.yml
./compiled-broker/gcp-service-broker generate use > metadata/manifest.yml
./compiled-broker/gcp-service-broker generate customization > metadata/docs/customization.md
./compiled-broker/gcp-service-broker generate use > metadata/docs/use.md
