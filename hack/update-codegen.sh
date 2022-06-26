#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/casbin/k8s-gatekeeper/pkg/generated github.com/casbin/k8s-gatekeeper/pkg/apis \
  k8sauthz:v1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

# To use your own boilerplate text append:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt
mv github.com/casbin/k8s-gatekeeper/pkg/apis/k8sauthz/v1/zz_generated.deepcopy.go pkg/apis/k8sauthz/v1/zz_generated.deepcopy.go
mv github.com/casbin/k8s-gatekeeper/pkg/generated pkg/generated
rm -rf ./github.com