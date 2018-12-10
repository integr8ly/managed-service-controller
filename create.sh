#!/usr/bin/env bash

cat ./deploy/operator-rbac.yaml | sed -e 's/${NAMESPACE}/managed-service-admin/g' | oc create -f -
cat ./deploy/managed-service-rbacs/fuse-rbac.yaml | sed -e 's/${NAMESPACE}/managed-service-admin/g' | oc create -f -
cat ./deploy/managed-service-rbacs/integration-rbac.yaml | sed -e 's/${NAMESPACE}/managed-service-admin/g' | oc create -f -
oc create -f ./deploy/operator.yaml
oc create -f ./deploy/cr.yaml
oc process -f https://raw.githubusercontent.com/integr8ly/managed-service-broker/master/templates/broker.template.yaml -p IMAGE_TAG=1.0.3 -p NAMESPACE=managed-service-admin  -p ROUTE_SUFFIX=127.0.0.1.nip.io  -p IMAGE_ORG=sedroche  -p CHE_DASHBOARD_URL=gfdgdf  -p LAUNCHER_DASHBOARD_URL=fgfdgfd  -p THREESCALE_DASHBOARD_URL=gfdgfdg  | oc create -f -