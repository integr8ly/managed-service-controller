#!/usr/bin/env bash

IPADDRESS=$1

BASE64_CERT=$(go run ./hack/credentials --local --ip-address=$IPADDRESS | sed -n '/Base64 encoded CA cert/{n;p;}')

cat ./deploy/webhook.external.yaml | sed -e "s/<DOMAIN>/$IPADDRESS/g" -e "s/<CA_BUNDLE>/$BASE64_CERT/g" | oc create -f -

echo "Server key is at ./build/tmp/server-credentials/server.key.pem"
echo "Server cert is at ./build/tmp/server-credentials/server.cert.pem"