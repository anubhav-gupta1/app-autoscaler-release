#!/bin/bash

set -euo pipefail
script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source "${script_dir}/common.sh"

function main(){
  step "cleaning up deployment ${DEPLOYMENT_NAME}"
  bosh_login
  cf_login

  cleanup_acceptance_run
  cleanup_service_broker
  cleanup_bosh_deployment
  delete_releases
  cleanup_credhub
}

[ "${BASH_SOURCE[0]}" == "${0}" ] && main "$@"
