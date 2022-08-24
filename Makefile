SHELL := /bin/bash

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
NAME := goboilerplate
MAINCMD := ./cmd/$(NAME)
IMAGE_TAG := latest

export GO111MODULE := on
export CGO_ENABLED := 0
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)

GO_EXEC := go
DOCKER_EXEC := docker

PACKAGES = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor -f '{{.Dir}}' ./...

export TOOLSBIN := $(shell pwd)/.bin

include build/*.mk

.PHONY: default
default: checks build

.PHONY: up
up:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml up --force-recreate --build

.PHONY: down
down:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml down --volumes --rmi local --remove-orphans


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
