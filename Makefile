# Copyright 2021 Adarga Limited
# SPDX-License-Identifier: Apache-2.0

SHELL := /bin/bash
.PHONY: generate generate-client generate-crd build image diff vet test compliance check
all: check build

# For compatibility with pachyderm, we need to use v1.17 code generator;
# Unfortunately, the docker image for this version doesn't include
# controller-gen so that needs to be installed separately
# (sigs.k8s.io/controller-tools/cmd/controller-gen)

KCG_IMAGE := quay.io/slok/kube-code-generator:v1.17.3
PROJECT_PACKAGE := github.com/adarga-ai/pachyderm-pipeline-controller

generate: generate-client generate-crd

generate-client:
	docker run -it --rm \
    -v $(PWD):/go/src/$(PROJECT_PACKAGE) \
    -e PROJECT_PACKAGE=$(PROJECT_PACKAGE) \
    -e CLIENT_GENERATOR_OUT=$(PROJECT_PACKAGE)/client \
    -e APIS_ROOT=$(PROJECT_PACKAGE)/apis \
    -e GROUPS_VERSION="pachyderm:v1alpha1" \
    -e GENERATION_TARGETS="deepcopy,client" \
    $(KCG_IMAGE)

generate-crd:
	controller-gen crd:trivialVersions=true paths="./apis/..." output:crd:artifacts:config=manifests

build:
	go build -o build/pachyderm-pipeline-controller .

# Docker image is built using Google Cloud Buildpacks
# Currently not published

VERSION   := 1.10
DKR_IMAGE := pachyderm-pipeline-controller

image:
	pack build \
		--builder gcr.io/buildpacks/builder:v1 \
		$(DKR_IMAGE):v$(VERSION)

diff:
	diff -u <(echo -n) <(gofmt -d ./)

vet:
	go vet ./...

test:
	go test ./...

# REUSE Spec compliance is tested using
# https://git.fsfe.org/reuse/tool

compliance:
	reuse lint

check: diff vet test compliance

