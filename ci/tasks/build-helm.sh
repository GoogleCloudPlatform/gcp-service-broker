#!/bin/sh

VERSION=5.0.0

echo "dir"
ls

echo "parent"
ls ..

echo "$VERSION"

helm package --app-version=$VERSION --dependency-update --version=$VERSION --destination ../helm-chart "deployments/helm/gcp-service-broker"
