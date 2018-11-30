TAG = 1.0.0
DOCKERORG = quay.io/integreatly
PROJECT_NAME = managed-service-controller

.phony: build_binary
build_binary:
	./build/build.sh $(PROJECT_NAME)

.phony: run
run:
	./build/run.sh $(WATCH_NAMESPACE) $(PROJECT_NAME)

.phony: run_local
run_local:
	./build/run_local.sh $(WATCH_NAMESPACE) $(PROJECT_NAME)

.phony: build_image
build_image: build_binary
	./build/docker_build.sh $(DOCKERORG) $(PROJECT_NAME) $(TAG)

.phony: push
push:
	docker push $(DOCKERORG)/$(PROJECT_NAME):$(TAG)

.phony: build_and_push
build_and_push: build_image push