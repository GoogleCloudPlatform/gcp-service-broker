#!/usr/bin/env sh
set -e

echo "Generating metadata"
mkdir -p /artifacts/metadata/docs

echo $COMMIT_SHA > /artifacts/metadata/revision
/workspace/compiled-broker/gcp-service-broker version > /artifacts/metadata/version
/workspace/compiled-broker/gcp-service-broker generate tile > /artifacts/metadata/tile.yml
/workspace/compiled-broker/gcp-service-broker generate use > /artifacts/metadata/manifest.yml
/workspace/compiled-broker/gcp-service-broker generate customization > /artifacts/metadata/docs/customization.md
/workspace/compiled-broker/gcp-service-broker generate use --destination-dir="/artifacts/metadata/docs/"
