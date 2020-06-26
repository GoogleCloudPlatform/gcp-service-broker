#!/usr/bin/env sh
set -e

echo "Installing zip"
apk update
apk add zip

export OUTPUT_DIR="/artifacts/servers"
export CURRENT_VERSION="$(cat /artifacts/metadata/version)"
mkdir -p $OUTPUT_DIR

echo "Creating source/CF app archive"
zip $OUTPUT_DIR/gcp-service-broker-$CURRENT_VERSION-cf-app.zip -r . -x *.git* product/\* release/\* examples/\* > /dev/null 2>&1
ls -la $OUTPUT_DIR
