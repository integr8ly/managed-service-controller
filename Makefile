TAG = 1.0.0
ORG = quay.io/integreatly
PROJECT_NAME = managed-service-controller
IMAGE = managed-service-controller

setup/dep:
	@echo Installing dep
	go get -u github.com/gobuffalo/packr/packr
	@echo setup complete

.PHONY: code/run
code/run:
	@./build/run.sh $(WATCH_NAMESPACE) $(PROJECT_NAME)

.PHONY: code/compile
code/compile:
	@./build/build.sh $(PROJECT_NAME)

.PHONY: code/gen
code/gen:
	@operator-sdk generate k8s

.PHONY: code/check
code/check:
	@diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: code/fix
code/fix:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: image/build
image/build: code/compile
	@./build/docker_build.sh $(ORG) $(IMAGE) $(TAG)

.PHONY: image/push
image/push:
	@docker push $(ORG)/$(IMAGE):$(TAG)

.PHONY: image/build/push
image/build/push: image/build image/push

.PHONY: cluster/prepare
cluster/prepare:
	@-oc create -f ./deploy/crds/integreatly_v1alpha1_managedservicenamespace_crd.yaml
	@-oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
	@-oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml
	@-oc create -f https://raw.githubusercontent.com/syndesisio/syndesis/master/install/operator/deploy/syndesis-crd.yml
	@-oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/crd.yaml
	@-oc create -f https://raw.githubusercontent.com/syndesisio/fuse-online-install/1.4.8/resources/fuse-online-image-streams.yml -n openshift
	@-oc create -f ./deploy/fuse-image-stream.yaml -n openshift
	@-oc create namespace enmasse

.PHONY: cluster/deploy
cluster/deploy:
	@-oc create -f ./deploy/role.yaml
	@-cat ./deploy/role_binding.yaml | sed -e 's/{NAMESPACE}/$(PROJECT_NAME)/g' | oc create -f -
	@-cat ./deploy/managed-service-rbacs/fuse-rbac.yaml | sed -e 's/{NAMESPACE}/$(PROJECT_NAME)/g' | oc create -f -
	@-cat ./deploy/managed-service-rbacs/integration-rbac.yaml | sed -e 's/{NAMESPACE}/$(PROJECT_NAME)/g' | oc create -f -
	@-oc create -f ./deploy/service_account.yaml
	@-oc create -f ./deploy/operator.yaml

.PHONY: cluster/deploy/remove
cluster/deploy/remove:
	@-oc delete -f ./deploy/role.yaml
	@-for i in $(oc get rolebindings -n $(PROJECT_NAME) | grep managed-service-controller | awk '{ print $1}');do oc delete $i; done
	@--oc delete -f ./deploy/managed-service-rbacs/
	@-oc delete -f ./deploy/service_account.yaml
	@-oc delete -f ./deploy/operator.yaml
	@-oc delete all -n $(PROJECT_NAME) -l for=managed-service-controller

.PHONY: cluster/clean
cluster/clean:
	@-oc delete -f ./deploy/crds/integreatly_v1alpha1_managedservicenamespace_crd.yaml
	@-oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
	@-oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml
	@-oc delete -f https://raw.githubusercontent.com/syndesisio/syndesis/master/install/operator/deploy/syndesis-crd.yml
	@-oc delete -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/crd.yaml
	@-oc delete -f https://raw.githubusercontent.com/syndesisio/fuse-online-install/1.4.8/resources/fuse-online-image-streams.yml -n openshift
	@-oc delete -f ./deploy/fuse-image-stream.yaml -n openshift
	@-oc delete namespace enmasse