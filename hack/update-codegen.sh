#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go generate github.com/maiqueb/persistentips/pkg/crd/persistentip/v1alpha1/
