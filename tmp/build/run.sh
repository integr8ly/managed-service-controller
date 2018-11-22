#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go > /dev/null; then
	echo "golang needs to be installed"
	exit 1
fi

if (( $# != 2 )); then
	echo "WATCH_NAMESPACE and PROJECT_NAME must be set"
	exit 1
fi;

WATCH_NAMESPACE=$1
PROJECT_NAME=$2
BIN_DIR="$(pwd)/tmp/_output/bin"

echo "running "${PROJECT_NAME}"..."
KUBERNETES_CONFIG=$HOME/.kube/config OPERATOR_NAME=${PROJECT_NAME} WATCH_NAMESPACE=$WATCH_NAMESPACE ${BIN_DIR}/${PROJECT_NAME}