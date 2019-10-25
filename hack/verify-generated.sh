#!/usr/bin/env bash
# Copyright 2019 Google LLC
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

#
# This is a helper script for validating if GCP Service Broker generators were executed and results were committed
#
set -o nounset
set -o errexit
set -o pipefail

readonly CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Prints first argument as header. Additionally prints current date.
shout() {
    echo -e "
#################################################################################################
# $(date)
# $1
#################################################################################################
"
}

shout "- Running GCP Service Broker generators..."
${CURRENT_DIR}/build.sh

shout "- Checking for modified files..."

# The porcelain format is used because it guarantees not to change in a backwards-incompatible
# way between Git versions or based on user configuration.
# source: https://git-scm.com/docs/git-status#_porcelain_format_version_1
if [[ -n "$(git status --porcelain)" ]]; then
    echo "Detected modified files:"
    git status --porcelain

    echo "
    Run:
        ./hack/build.sh
    in the root of the repository and commit changes.
    "
    exit 1
else
    echo "No issues detected. Have a nice day :-)"
fi

