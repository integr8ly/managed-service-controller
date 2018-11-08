TAG = 0.0.1
DOCKERORG = quay.io/integreatly
BROKER_IMAGE_NAME = managed-service-namespace
CERT_FILE = ./build/tmp/server-credentials/server.cert.pem
KEY_FILE =  ./build/tmp/server-credentials/server.key.pem
PORT =  8443

.phony: build_and_push
build_and_push: build_image push

.phony: push
push:
	docker push $(DOCKERORG)/$(BROKER_IMAGE_NAME):$(TAG)

.phony: build_image
build_image: build_binary
	docker build -t $(DOCKERORG)/$(BROKER_IMAGE_NAME):$(TAG) -f ./build/Dockerfile ./build

.phony: build_binary
build_binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/tmp/bin/webhook ./cmd/webhook

.phony: service_context
service_context:
	go run ./hack/credentials --service --service-name="managed-service-namespaces-webhook" --namespace=$(NAMESPACE)

.phony: deploy_template
deploy_template:
	./hack/deploy_template.sh $(NAMESPACE)

.phony: delete_template $(NAMESPACE)
delete_template:
	./hack/delete_template.sh $(NAMESPACE)

.phony: local_context
local_context:
	go run ./hack/credentials --local --ip-address=$(IP_ADDRESS)

.phony: deploy_local_config
deploy_local_config:
	./hack/deploy_local_config.sh $(IP_ADDRESS)

.phony: delete_local_config
delete_local_config:
	oc delete -f ./deploy/webhook.external.yaml

.phony: run
run:
	go run ./cmd/webhook/ --tls-cert-file=$(CERT_FILE) --tls-private-key-file=$(KEY_FILE) --port=$(PORT)