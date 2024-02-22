#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go generate github.com/k8snetworkplumbingwg/ipamclaims/pkg/crd/ipamclaims/v1alpha1/
