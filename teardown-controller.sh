#!/usr/bin/env bash

oc delete -f ./deploy/crd.yaml
oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml
oc delete -f https://raw.githubusercontent.com/syndesisio/syndesis/master/install/operator/deploy/syndesis-crd.yml
oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/crd.yaml
oc delete -f https://raw.githubusercontent.com/syndesisio/fuse-online-install/1.4.8/resources/fuse-online-image-streams.yml -n openshift
oc delete -f ./deploy/fuse-image-stream.yaml -n openshift
oc delete namespace enmasse
oc delete namespace managed-service-admin
oc delete namespace consumer1
oc delete namespace consumer2
