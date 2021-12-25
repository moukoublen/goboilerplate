SHELL := /bin/bash

NAME := goboilerplate
MAINCMD := ./cmd/${NAME}
IMAGE_TAG := latest

GO111MODULE := on
export GO111MODULE
CGO_ENABLED := 0
export CGO_ENABLED
GOFLAGS := -mod=vendor
export GOFLAGS

GO_EXEC := go
DOCKER := docker

PACKAGES = $(GO_EXEC) list -tags=${TAGS} -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=${TAGS} -mod=vendor -f '{{.Dir}}' ./...

.PHONY: build
build:
	@GO_EXEC=$(GO_EXEC) $(CURDIR)/scripts/build ${MAINCMD}

.PHONY: env
env:
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
	@if [[ -n "$$(goimports -l `${FOLDERS}` | tee /dev/stderr)" ]]; then \
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

.PHONY: goimports-w
goimports-w:
	goimports -w `${FOLDERS}`

.PHONY: gofmt
gofmt:
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

.PHONY: gofmt-w
gofmt-w:
	gofmt -w `${FOLDERS}`

.PHONY: gofumpt
gofumpt:
	@if [[ -n "$$(gofumpt -l `${FOLDERS}` | tee /dev/stderr)" ]]; then \
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

.PHONY: gofumpt-w
gofumpt-w:
	gofumpt -w `${FOLDERS}`

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: staticcheck
staticcheck:
	staticcheck -f=stylish -checks=all,-ST1000 -tests ./...

.PHONY: checks
checks: vet staticcheck gofumpt goimports

.PHONY: test
test:
	$(GO_EXEC) test -timeout 60s -tags="${TAGS}" -coverprofile cover.out -covermode atomic `${PACKAGES}`
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
