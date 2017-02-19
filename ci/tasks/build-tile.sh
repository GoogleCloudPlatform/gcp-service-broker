#!/usr/bin/env bash

set -e

export service_broker_dir=src/gcp-service-broker

pushd "$service_broker_dir"
    zip /tmp/gcp-service-broker.zip -r . -x *.git* product/\* release/\*

    tile build
popd

mv "$service_broker_dir/product/*.pivotal" candidate/