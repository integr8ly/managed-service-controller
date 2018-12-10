#!/usr/bin/env bash

oc delete -f ./deploy/cr.yaml
oc delete -f ./deploy/operator-rbac.yaml
oc delete -f ./deploy/managed-service-rbacs/
oc delete -f ./deploy/operator.yaml
oc delete all -l for=managed-service-controller
oc process -f https://raw.githubusercontent.com/integr8ly/managed-service-broker/master/templates/broker.template.yaml  -p IMAGE_TAG=1.0.3 -p NAMESPACE=managed-service-admin  -p ROUTE_SUFFIX=127.0.0.1.nip.io  -p IMAGE_ORG=sedroche  -p CHE_DASHBOARD_URL=gfdgdf  -p LAUNCHER_DASHBOARD_URL=fgfdgfd  -p THREESCALE_DASHBOARD_URL=gfdgfdg  | oc delete -f -
oc delete serviceinstance --all --force --grace-period=0 --cascade=false -n consumer1
oc delete serviceinstance --all --force --grace-period=0 --cascade=false -n consumer2