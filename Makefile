SHELL := /bin/bash

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
VERSION ?= 0.0.0
X_FLAGS := \
		-X '$(MODULE)/build.Version=$(VERSION)' \
		-X '$(MODULE)/build.Branch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Commit=$(shell git rev-parse HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.CommitShort=$(shell git rev-parse --short HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Tag=$(shell git describe --tags 2>/dev/null || true)'
IMAGE_NAME := goboilerplate
IMAGE_TAG := latest

export GO111MODULE := on
export CGO_ENABLED := 0
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)

GO_EXEC ?= go
export GO_EXEC
DOCKER_EXEC ?= docker
export DOCKER_EXEC

include build/ci.mk

.PHONY: default
default: checks build

.PHONY: up
up:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml up --force-recreate --build

.PHONY: down
down:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml down --volumes --rmi local --remove-orphans

.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd
# man git-clean

.PHONY: env
env:
	@echo "Module: $(MODULE)"
	$(GO_EXEC) env
	@echo ""
	@echo ">>> Packages:"
	$(PACKAGES)
	@echo ""
	@echo ">>> Folders:"
	$(FOLDERS)
	@echo ""
	@echo ">>> Tools:"
	@echo '$(TOOLSBIN)'

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags "$(TAGS)"
BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS)" -tags "$(TAGS)"

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

.PHONY: build
build: $(shell ls -d cmd/* | sed -e 's/\//./')
cmd.%: CMDNAME=$*
cmd.%:
	$(GO_EXEC) build $(BUILD_FLAGS) -o ./output/$(CMDNAME) ./cmd/$(CMDNAME)

.PHONY: image
image:
	$(DOCKER_EXEC) build . -f $(CURDIR)/build/docker/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.19
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

.PHONY: clean
clean:
	rm -rf ./output
