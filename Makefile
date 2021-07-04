SHELL := /bin/bash

NAME := goboilerplate
MAINCMD := ./cmd/${NAME}
IMAGE := ${NAME}
IMAGE_TAG := latest

GO111MODULE := on
export GO111MODULE
CGO_ENABLED := 0
export CGO_ENABLED

GO_EXEC := go
DOCKER := docker
COMPOSE := docker-compose

PACKAGES = $(GO_EXEC) list -tags=${TAGS} -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=${TAGS} -mod=vendor -f '{{.Dir}}' ./...

.PHONY: build
build:
	@GO_EXEC=$(GO_EXEC) $(CURDIR)/scripts/build ${NAME} ${MAINCMD}

.PHONY: clean
clean:
	rm -f ${NAME}

.PHONY: env
env:
	$(GO_EXEC) env
	@echo ">>> Packages:"
	${PACKAGES}
	@echo ">>> Folders:"
	${FOLDERS}

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

#@go get -u golang.org/x/tools/cmd/goimports
.PHONY: goimports
goimports:
	@$(CURDIR)/scripts/goimports "${FOLDERS}"

.PHONY: gofmt
gofmt:
	@$(CURDIR)/scripts/gofmt "${FOLDERS}"

#@go get -u golang.org/x/lint/golint
.PHONY: golint
golint:
	@$(CURDIR)/scripts/golint "${PACKAGES}"

.PHONY: vet
vet:
	@echo ">>> go vet <<<"
	$(GO_EXEC) vet `${PACKAGES}`
	@echo ""

.PHONY: checks
checks: goimports gofmt golint vet

.PHONY: test
test:
	$(GO_EXEC) test -timeout 60s -tags="${TAGS}" -coverprofile cover.out -covermode atomic `${PACKAGES}`
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: up
up:
	$(COMPOSE) -f $(CURDIR)/deployments/compose/docker-compose.yml up --force-recreate --build

.PHONY: down
down:
	$(COMPOSE) -f $(CURDIR)/deployments/compose/docker-compose.yml down --volumes

.PHONY: image
image:
	$(DOCKER) build . -f $(CURDIR)/build/docker/Dockerfile -t ${IMAGE}:${IMAGE_TAG}

.PHONY: image-ci
image-ci:
	$(DOCKER) build $(CURDIR)/build/ci/ -t ${NAME}-ci:latest

.PHONY: default
default: checks build
