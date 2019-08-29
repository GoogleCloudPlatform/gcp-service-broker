#!/usr/bin/env bash

set -e

export OUTPUT_DIR=$PWD/tiles
export SERVICE_BROKER_DIR=src/gcp-service-broker
export CURRENT_VERSION="$(cat metadata/version)"

apt install -y zip

mkdir -p tiles

pushd "$SERVICE_BROKER_DIR"
    zip /tmp/gcp-service-broker.zip -r . -x *.git* product/\* release/\* examples/\*
    cp /tmp/gcp-service-broker.zip $OUTPUT_DIR/gcp-service-broker-$CURRENT_VERSION-cf-app.zip

    tile build "$CURRENT_VERSION"
    mv "product/"*.pivotal $OUTPUT_DIR

    tile build "$CURRENT_VERSION-rc"
    mv "product/"*.pivotal $OUTPUT_DIR
popd
