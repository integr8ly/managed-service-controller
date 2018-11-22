#!/usr/bin/env bash

if ! which docker > /dev/null; then
	echo "docker needs to be installed"
	exit 1
fi

if (( $# != 3 )); then
	echo "DOCKERORG, PROJECT_NAME, and TAG must be set"
	exit 1
fi;

DOCKERORG=$1
PROJECT_NAME=$2
TAG=$3
IMAGE=$DOCKERORG/$PROJECT_NAME:$TAG

echo "building container ${IMAGE}..."
docker build -t "${IMAGE}" -f tmp/build/Dockerfile .