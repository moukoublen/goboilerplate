SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

GO_EXEC ?= go
export GO_EXEC
DOCKER_EXEC ?= docker
export DOCKER_EXEC

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
VERSION ?= 0.0.0
X_FLAGS = \
		-X '$(MODULE)/build.Version=$(VERSION)' \
		-X '$(MODULE)/build.Branch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Commit=$(shell git rev-parse HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.CommitShort=$(shell git rev-parse --short HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Tag=$(shell git describe --tags 2>/dev/null || true)'
IMAGE_NAME ?= goboilerplate
IMAGE_TAG ?= latest

GO_PACKAGES = $(GO_EXEC) list -tags='$(TAGS)' -mod=vendor ./...
GO_FOLDERS = $(GO_EXEC) list -tags='$(TAGS)' -mod=vendor -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go' -not -path './vendor/*'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)
BUILD_OUTPUT ?= $(CURDIR)/output

include $(CURDIR)/scripts/tools.mk

.DEFAULT_GOAL=default
.PHONY: default
default: checks build

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.23
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

# https://go.dev/ref/mod#go-get
# -u flag tells go get to upgrade modules
# -t flag tells go get to consider modules needed to build tests of packages named on the command line.
# When -t and -u are used together, go get will update test dependencies as well.
.PHONY: go-deps-upgrade
go-deps-upgrade:
	$(GO_EXEC) get -u -t ./...
	$(GO_EXEC) mod tidy -go=1.23
	$(GO_EXEC) mod vendor

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags '$(TAGS)'
BUILD_FLAGS := -mod=vendor -a -ldflags '-s -w $(X_FLAGS)' -tags '$(TAGS)'
BUILD_FLAGS_DEBUG := -mod=vendor -ldflags '$(X_FLAGS)' -tags '$(TAGS)'

.PHONY: build
build: $(shell ls -d cmd/* | sed -e 's/\//./')

cmd.%: CMDNAME=$*
cmd.%:
	$(GO_EXEC) env
	@echo ''
	CGO_ENABLED=0 $(GO_EXEC) build $(BUILD_FLAGS) -o $(BUILD_OUTPUT)/$(CMDNAME) ./cmd/$(CMDNAME)

dbg.%: BUILD_FLAGS=$(BUILD_FLAGS_DEBUG)
dbg.%: cmd.%
	@echo "debug binary done"

.PHONY: clean
clean:
	rm -rf $(BUILD_OUTPUT)

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

## https://docs.docker.com/reference/cli/docker/buildx/build/
## --output='type=docker'
## --output='type=image,push=true'
## --platform=linux/arm64
## --platform=linux/amd64,linux/arm64,linux/arm/v7
## --platform=local
## --progress='plain'
## make DOCKER_BUILD_PLATFORM=linux/arm64 image
DOCKER_BUILD_PLATFORM ?= local
DOCKER_BUILD_OUTPUT ?= type=docker
.PHONY: image
image:
	$(DOCKER_EXEC) buildx build \
		--output='type=docker' \
		--file='$(CURDIR)/build/docker/Dockerfile' \
		--tag='$(IMAGE_NAME):$(IMAGE_TAG)' \
		--platform='$(DOCKER_BUILD_PLATFORM)' \
		--output='$(DOCKER_BUILD_OUTPUT)' \
		.

DOCKER_COMPOSE_EXEC ?= $(DOCKER_EXEC) compose -f $(CURDIR)/deployments/local/docker-compose.yml

.PHONY: compose-up
compose-up:
	$(DOCKER_COMPOSE_EXEC) up --force-recreate --build

.PHONY: compose-up-detach
compose-up-detach:
	$(DOCKER_COMPOSE_EXEC) up --force-recreate --build --detach

.PHONY: compose-down
compose-down:
	$(DOCKER_COMPOSE_EXEC) down --volumes --rmi local --remove-orphans

# If the first target is "compose-exec"
# remove the first argument 'compose-exec' and store the rest in DOCKER_COMPOSE_ARGS
# and ignore the subsequent arguments as make targets.
# (using spaces for indentation)
ifeq (compose-exec,$(firstword $(MAKECMDGOALS)))
    DOCKER_COMPOSE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    $(eval $(DOCKER_COMPOSE_ARGS):;@:)
endif

.PHONY: compose-exec
compose-exec:
	$(DOCKER_COMPOSE_EXEC) exec $(DOCKER_COMPOSE_ARGS)

.PHONY: env
env:
	@echo "Module: $(MODULE)"
	$(GO_EXEC) env
	@echo ""
	@echo ">>> Packages:"
	$(GO_PACKAGES)
	@echo ""
	@echo ">>> Folders:"
	$(GO_FOLDERS)
	@echo ""
	@echo ">>> Files:"
	$(GO_FILES)
	@echo ""
	@echo ">>> Tools:"
	@echo '$(TOOLS_BIN)'
	@echo ""
	@echo ">>> Path:"
	@echo "$${PATH}" | tr ':' '\n'

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint-github-actions

.PHONY: run
run:
	$(GO_EXEC) run -mod=vendor ./cmd/goboilerplate

.PHONY: ci-gen-n-format
ci-gen-n-format: goimports gofumpt
	./scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	./scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt
	@./scripts/sh-checks
	@./scripts/git-check-dirty

