#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}" &> /dev/null || exit 1

make install
make install-deployer-gce

# currently equivalent to /home/prow/go/src/github.com/kubernetes/cloud-provider-gcp
K_REPO_ROOT="${REPO_ROOT}/../../kubernetes/cloud-provider-gcp"

# TODO(spiffxp): remove this when cloudprovider-gcp has a .bazelversion file
export USE_BAZEL_VERSION=5.3.0
# TODO(spiffxp): remove this when gce-build-up-down job updated to do this,
#                or when bazel 5.3.0 is preinstalled on kubekins image
if [ "${CI}" == "true" ]; then
  go install github.com/bazelbuild/bazelisk@latest
  mkdir -p /tmp/use-bazelisk
  ln -s "$(go env GOPATH)/bin/bazelisk" /tmp/use-bazelisk/bazel
  export PATH="/tmp/use-bazelisk:${PATH}"
fi

kubetest2 gce \
            -v 2 \
            --repo-root "$K_REPO_ROOT" \
            --build \
            --up \
            --down
