#!/bin/bash
json=${1:-${HOME}/Desktop/gcp.json}

script_dir=`dirname $0`
pushd $script_dir/../../.. > /dev/null
export GOPATH=`pwd`
popd > /dev/null

$script_dir/drop_tables.sh

export ROOT_SERVICE_ACCOUNT_JSON="$(< $json)"

m=$GOPATH/src/gcp-service-broker/manifest.yml

[[ -f $m ]] || echo "found it"
export SERVICES=$(ruby -ryaml -rjson -e 'puts JSON.pretty_generate(YAML.load(ARGF))' ${m}| jq -r '.applications[0].env.SERVICES')

echo $SERVICES

export PRECONFIGURED_PLANS=$(ruby -ryaml -rjson -e 'puts JSON.pretty_generate(YAML.load(ARGF))' ${m}| jq -r '.applications[0].env.PRECONFIGURED_PLANS')

echo $PRECONFIGURED_PLANS
export SECURITY_USER_NAME=admin
export SECURITY_USER_PASSWORD=admin
export DB_HOST=localhost
export DB_USERNAME=gcp-service-broker
export DB_PASSWORD=qwerty
echo $GOPATH
echo $ROOT_SERVICE_ACCOUNT_JSON
go run $GOPATH/src/gcp-service-broker/server.go
