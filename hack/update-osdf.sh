#!/usr/bin/env bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the License);
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an AS IS BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eux
cd "${0%/*}"/..

if [ "$#" == "0" ]; then
	echo "Usage: $0 broker-version"
	exit 1
fi

temp_dir=$(mktemp -d)

go run ./tools/osdfgen/osdfgen.go -p . -o "$temp_dir/osdf.csv"

curl -X POST -F "name=gcp-service-broker" -F "version=$1" -F "file=@$temp_dir/osdf.csv" http://osdf-generator.cfapps.io/generate-disclosure > OSDF.txt
