#!/bin/sh

set -e

mkdir metadata

./compiled-broker/gcp-service-broker version > metadata/version
./compiled-broker/gcp-service-broker generate tile > metadata/tile.yml
./compiled-broker/gcp-service-broker generate customization > metadata/customization.md
./compiled-broker/gcp-service-broker generate use > metadata/use.md
./compiled-broker/gcp-service-broker generate use > metadata/manifest.yml
