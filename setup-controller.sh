#!/usr/bin/env bash

oc create -f ./deploy/crd.yaml
oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml
oc create -f https://raw.githubusercontent.com/syndesisio/syndesis/master/install/operator/deploy/syndesis-crd.yml
oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/crd.yaml
oc create -f https://raw.githubusercontent.com/syndesisio/fuse-online-install/1.4.8/resources/fuse-online-image-streams.yml -n openshift
oc create -f ./deploy/fuse-image-stream.yaml -n openshift
oc create namespace enmasse
oc create namespace consumer1
oc create namespace consumer2
oc new-project managed-service-admin
