SHELL := /bin/bash

.PHONY: default
default: checks build

export TOOLSBIN := $(shell pwd)/.bin
include build/tools.Makefile

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
NAME := goboilerplate
MAINCMD := ./cmd/${NAME}
IMAGE_TAG := latest

export GO111MODULE := on
export CGO_ENABLED := 0
#export GOFLAGS := -mod=vendor

GO_EXEC := go
DOCKER := docker

PACKAGES = $(GO_EXEC) list -tags=${TAGS} -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=${TAGS} -mod=vendor -f '{{.Dir}}' ./...

VERSION = 0.0.0
VER_FLAGS = \
		-X '${MODULE}/build.Version=${VERSION}' \
		-X '${MODULE}/build.Branch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.Commit=$(shell git rev-parse HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.CommitShort=$(shell git rev-parse --short HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.Tag=$(shell git describe --tags 2>/dev/null || true)'

.PHONY: build
build:
	$(GO_EXEC) build -mod=vendor -ldflags "-extldflags -static ${VER_FLAGS}" ${MAINCMD}

.PHONY: env
env:
	@echo "Module: ${MODULE}"
	@echo "Name  : ${NAME}"
	@echo "Cmd   : ${MAINCMD}"
	$(GO_EXEC) env
	@echo ""
	@echo ">>> Packages:"
	${PACKAGES}
	@echo ""
	@echo ">>> Folders:"
	${FOLDERS}

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -compat=1.17
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

.PHONY: vet
vet:
	$(GO_EXEC) vet `${PACKAGES}`
	@echo ""

.PHONY: goimports
goimports:
	@echo '${TOOLSBIN}/goimports -l `${FOLDERS}`'
	@if [[ -n "$$(${TOOLSBIN}/goimports -l `${FOLDERS}` | tee /dev/stderr)" ]]; then \
		echo 'goimports errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    goimports -d `${FOLDERS}`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    goimports -w `${FOLDERS}`'; \
		echo '  or'; \
		echo '    make goimports-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: goimports-w
goimports-w:
	${TOOLSBIN}/goimports -w `${FOLDERS}`

.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `${FOLDERS}`'
	@if [[ -n "$$(gofmt -l `${FOLDERS}` | tee /dev/stderr)" ]]; then \
		echo 'gofmt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    gofmt -d `${FOLDERS}`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    gofmt -w `${FOLDERS}`'; \
		echo '  or'; \
		echo '    make gofmt-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofmt-w
gofmt-w:
	gofmt -w `${FOLDERS}`

.PHONY: gofumpt
gofumpt:
	@echo '${TOOLSBIN}/gofumpt -l `${FOLDERS}`'
	@if [[ -n "$$(${TOOLSBIN}/gofumpt -l `${FOLDERS}` | tee /dev/stderr)" ]]; then \
		echo 'gofumpt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    gofumpt -d `${FOLDERS}`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    gofumpt -w `${FOLDERS}`'; \
		echo '  or'; \
		echo '    make gofumpt-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofumpt-w
gofumpt-w:
	${TOOLSBIN}/gofumpt -w `${FOLDERS}`

.PHONY: golangci-lint
golangci-lint:
	${TOOLSBIN}/golangci-lint run
	@echo ''

.PHONY: staticcheck
staticcheck:
	${TOOLSBIN}/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''

.PHONY: checks
checks: vet staticcheck gofumpt goimports

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -tags="${TAGS}" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: up
up:
	$(DOCKER) compose -f $(CURDIR)/deployments/compose/docker-compose.yml up --force-recreate --build

.PHONY: down
down:
	$(DOCKER) compose -f $(CURDIR)/deployments/compose/docker-compose.yml down --volumes --rmi local --remove-orphans

.PHONY: image
image:
	$(DOCKER) build . -f $(CURDIR)/build/docker/Dockerfile -t ${NAME}:${IMAGE_TAG}

.PHONY: image-ci
image-ci:
	$(DOCKER) build . -f $(CURDIR)/build/ci/Dockerfile -t ${NAME}-ci:latest

.PHONY: default
default: checks build
