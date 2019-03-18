#!/bin/sh

VERSION="$(cat metadata/version)"

echo "packaging version: $VERSION"
helm init --client-only
helm package --app-version=$VERSION --dependency-update --version=$VERSION --destination helm-chart "gcp-service-broker/deployments/helm/gcp-service-broker"
