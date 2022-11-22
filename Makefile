SHELL := /bin/bash

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
NAME := goboilerplate
MAINCMD := ./cmd/$(NAME)
IMAGE_TAG := latest

export GO111MODULE := on
export CGO_ENABLED := 0
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)

GO_EXEC ?= go
DOCKER_EXEC ?= docker

include build/*.mk

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
	@echo "Name  : $(NAME)"
	@echo "Cmd   : $(MAINCMD)"
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
