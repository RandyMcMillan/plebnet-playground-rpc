SHELL := /bin/bash

PWD ?= pwd_unknown
TIME                                   := $(shell date +%s)
export TIME
PROJECT_NAME = $(notdir $(PWD))
export PROJECT_NAME
DOCKER:=$(shell which docker)
export DOCKER
DOCKER_COMPOSE:=$(shell which docker-compose)
export DOCKER_COMPOSE

ARCH                                   :=$(shell uname -m)
export ARCH
ifeq ($(ARCH),x86_64)
TRIPLET                                :=x86_64-linux-gnu
export TRIPLET
endif
ifeq ($(ARCH),arm64)
TRIPLET                                :=aarch64-linux-gnu
export TRIPLET
endif

GOPATH:=$(HOME)/go
export GOPATH
#export PATH=$(PATH):$(shell go env GOPATH)/bin

.PHONY: - install clean build release dry-release cov fmt help vet test
-: help
##    init:initialize some plebnet-playground (go and docker) dependancies
init:
	@[ "$(shell uname -s)" == "Darwin" ] && test brew &&  brew install -q golang docker docker-compose || echo $(shell uname -s)
	@[ "$(shell uname -s)" == "Linux"  ] && test apt  &&  apt  install    golang docker docker-compose || echo $(shell uname -s)
	git submodule update --init --recursive
	pushd cmd/docker/plebnet-playground-docker && make init
	$(DOCKER_COMPOSE) -f cmd/docker/plebnet-playground-docker/docker-compose.yaml pull
	test go && go mod download
	test go && go mod tidy
	test go && go get github.com/randymcmillan/plebnet-playground-rpc/internal/config
	test go && go get github.com/randymcmillan/plebnet-playground-rpc/internal/docker

##    install:installs dependencies
install:
	@[ "$(shell uname -s)" == "Darwin" ] && echo $(shell uname -s) & exit
	@[ "$(shell uname -s)" == "Linux"  ] && echo $(shell uname -s) & exit
	pushd cmd/docker/plebnet-playground-docker && make install

##    clean:cleans the binary
clean:
	@echo "Cleaning..."
	rm -rf dist
	go clean

##    build:build plebnet-playground-docker images and rpc binary
build:
	pushd cmd/docker/plebnet-playground-docker && make build para=true
	chmod u+x ./scripts/build
	./scripts/build

##    release:build and upload binaries to Github Releases
release:
	GITHUB_TOKEN=$(shell cat ~/GH_TOKEN.txt)  goreleaser --rm-dist

##    dry-release:build and test goreleaser
dry-release:
	goreleaser --snapshot --skip-publish --rm-dist

##    help:prints this help message
help:
	@echo ""
	@echo "Usage:"
	@echo ""
	@echo "make [arg]"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

##    fmt:Go Format
fmt:
	@echo "Gofmt..."
	@if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted"; exit 1; fi

##    vet:code analysis
vet:
	@echo "Vet..."
	@go vet ./...

##    test:runs go unit test with default values
test: clean install
	@echo "Testing..."
	go test -v -count=1 -race ./...

##    test-ci:runs travis tests
test-ci: clean
	@echo "Testing..."
	go test -short -v ./...

##    cov:generates coverage report
cov:
	@echo "Coverage..."
	go test -cover ./...

