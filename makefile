SHELL := /bin/bash

# ==============================================================================
# Building containers
DOCKER_IMAGE_NAME=signernode
DOCKER_FULL_IMAGE_NAME=ghcr.io/jffp113/${DOCKER_IMAGE_NAME}:latest

build-docker:
	docker build -t $(DOCKER_FULL_IMAGE_NAME) -f ./Docker/Dockerfile .

push-docker: build-docker
	docker push $(DOCKER_FULL_IMAGE_NAME)

# ==============================================================================
# Building go files
build:
	go build ./app/bootstrap/bootstrap.go
	go build ./app/signernode/signernode.go

clean:
	rm bootstrap
	rm signernode

# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor