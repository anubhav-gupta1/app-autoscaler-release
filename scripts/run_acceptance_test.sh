#!/bin/bash

set -exuo pipefail
# shellcheck disable=SC2155
export script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# shellcheck disable=SC1091
source "${script_dir}/pr-vars.source.sh"
export SKIP_TEARDOWN=true

#export SUITES="broker"
export GINKGO_OPTS="--progress --fail-fast -v "
echo "Running acceptance tests for PR:${PR_NUMBER}"
export NODES=1
export SUITES="api"
"${CI_DIR}/autoscaler/scripts/run-acceptance-tests.sh"