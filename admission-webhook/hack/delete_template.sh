#!/usr/bin/env bash

NAMESPACE=$1

oc process -f ./deploy/webhook.template.yaml --param=NAMESPACE=$NAMESPACE --param=CA_BUNDLE="PLACEHOLDER" --param=SERVER_KEY="PLACEHOLDER" --param=SERVER_CERT="PLACEHOLDER"  | oc delete -f - -n $NAMESPACE
oc delete clusterrolebindings -n testwebhook -l app=managed-service-namespaces-webhook