#!/usr/bin/env bash

NAMESPACE=$1

SERVICE_CONTEXT=( $(go run ./hack/credentials --service --service-name="managed-service-namespaces-webhook" --namespace=$NAMESPACE) )
CA_BUNDLE=${SERVICE_CONTEXT[4]}
SERVER_KEY=${SERVICE_CONTEXT[9]}
SERVER_CERT=${SERVICE_CONTEXT[14]}

oc process -f ./deploy/webhook.template.yaml --param=NAMESPACE=$NAMESPACE --param=CA_BUNDLE=$CA_BUNDLE --param=SERVER_KEY=$SERVER_KEY --param=SERVER_CERT=$SERVER_CERT  | oc create -f - -n $NAMESPACE
