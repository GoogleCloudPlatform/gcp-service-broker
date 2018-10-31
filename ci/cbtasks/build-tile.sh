#!/usr/bin/env bash
set -e

zip /tmp/gcp-service-broker.zip -r . -x *.git* product/\* release/\* examples/\*
tile build "$VERSION"
