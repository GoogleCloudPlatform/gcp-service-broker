#!/usr/bin/env bash

set -e

export SERVICE_BROKER_DIR=src/gcp-service-broker
export CURRENT_VERSION="$(cat metadata/version)"

pushd "$SERVICE_BROKER_DIR"
    zip /tmp/gcp-service-broker.zip -r . -x *.git* product/\* release/\* examples/\*
    cp /tmp/gcp-service-broker.zip ../tiles/gcp-service-broker-$CURRENT_VERSION-cf-app.zip

    tile build "$CURRENT_VERSION"
    mv "product/"*.pivotal ../tiles/

    tile build "$CURRENT_VERSION-rc"
    mv "product/"*.pivotal ../tiles/
popd
