#!/usr/bin/env bash

set -e

export SERVICE_BROKER_DIR=src/gcp-service-broker
export CURRENT_VERSION="$(cat metadata/version)"

pushd "$SERVICE_BROKER_DIR"
    zip /tmp/gcp-service-broker.zip -r . -x *.git* product/\* release/\* examples/\*

    tile build "$CURRENT_VERSION"
    tile build "$CURRENT_VERSION-rc"

popd

mv "$SERVICE_BROKER_DIR/product/"*.pivotal tiles/
