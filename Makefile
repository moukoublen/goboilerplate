SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
VERSION ?= 0.0.0
X_FLAGS := \
		-X '$(MODULE)/build.Version=$(VERSION)' \
		-X '$(MODULE)/build.Branch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Commit=$(shell git rev-parse HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.CommitShort=$(shell git rev-parse --short HEAD 2>/dev/null || true)' \
		-X '$(MODULE)/build.Tag=$(shell git describe --tags 2>/dev/null || true)'
IMAGE_NAME ?= goboilerplate
IMAGE_TAG ?= latest

PACKAGES = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor -f '{{.Dir}}' ./...

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)

GO_EXEC ?= go
export GO_EXEC
DOCKER_EXEC ?= docker
export DOCKER_EXEC

.DEFAULT_GOAL=default
.PHONY: default
default: checks build

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.21
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags "$(TAGS)"
BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS)" -tags "$(TAGS)"
BUILD_FLAGS_DEBUG := -mod=vendor -ldflags "$(X_FLAGS)" -tags "$(TAGS)"

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

.PHONY: build
build: $(shell ls -d cmd/* | sed -e 's/\//./')

cmd.%: CMDNAME=$*
cmd.%:
	CGO_ENABLED=0 $(GO_EXEC) build $(BUILD_FLAGS) -o ./output/$(CMDNAME) ./cmd/$(CMDNAME)

dbg.%: BUILD_FLAGS=$(BUILD_FLAGS_DEBUG)
dbg.%: cmd.%
	@echo "debug binary done"

.PHONY: clean
clean:
	rm -rf ./output

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

.PHONY: image
image:
	$(DOCKER_EXEC) build . -f $(CURDIR)/build/docker/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG)

DOCKER_COMPOSE_EXEC ?= $(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml

.PHONY: compose-up
compose-up:
	$(DOCKER_COMPOSE_EXEC) up --force-recreate --build

.PHONY: compose-down
compose-down:
	$(DOCKER_COMPOSE_EXEC) down --volumes --rmi local --remove-orphans

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
	@echo ""
	@echo ">>> Path:"
	@echo $${PATH}


####################################################################################
## <ci & external tools> ###########################################################
####################################################################################
.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

.PHONY: vet
vet:
	$(GO_EXEC) vet `$(PACKAGES)`
	@echo ""

TOOLSDIR ?= $(shell pwd)/.ext
TOOLSBIN ?= $(TOOLSDIR)/bin
export TOOLSBIN
export PATH := $(TOOLSBIN):$(PATH)
TOOLSDB ?= $(TOOLSDIR)/.db

uppercase = $(shell echo '$(1)' | tr '[:lower:]' '[:upper:]')

.PHONY: tools
tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/staticcheck \
	$(TOOLSBIN)/golangci-lint \
	$(TOOLSBIN)/gofumpt

$(TOOLSBIN):
	@mkdir -p $(TOOLSBIN)

$(TOOLSDB):
	@mkdir -p $(TOOLSDB)

# In make >= 4.4. .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLSDB)/%.ver
$(TOOLSDB)/%.ver: | $(TOOLSDB)
	@rm -f $(TOOLSDB)/$(word 1,$(subst ., ,$*)).*
	@touch $(TOOLSDB)/$*.ver

# In make >= 4.4 .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLSBIN)/%
$(TOOLSBIN)/%: DSC=$*
$(TOOLSBIN)/%: VER=$($(call uppercase,$*)_VER)
$(TOOLSBIN)/%: CMD=$($(call uppercase,$*)_CMD)
$(TOOLSBIN)/%: $(TOOLSDB)/$$(DSC).$$(VER).$(GO_VER).ver
	@echo -e "Installing \e[1;36m$(DSC)\e[0m@\e[1;36m$(VER)\e[0m using \e[1;36m$(GO_VER)\e[0m"
	GOBIN="$(TOOLSBIN)" CGO_ENABLED=0 $(GO_EXEC) install -trimpath -ldflags '-s -w -extldflags "-static"' "$(CMD)@$(VER)"
	@echo ""

## <staticcheck>
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/c
STATICCHECK_CMD=honnef.co/go/tools/cmd/staticcheck
STATICCHECK_VER:=2023.1.6
$(TOOLSDB)/staticcheck.$(STATICCHECK_VER).$(GO_VER).ver:
$(TOOLSBIN)/staticcheck:

.PHONY: staticcheck
staticcheck: $(TOOLSBIN)/staticcheck
	staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
## </staticcheck>

## <golangci-lint>
# https://github.com/golangci/golangci-lint/releases
GOLANGCI-LINT_CMD:=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI-LINT_VER:=v1.55.2
$(TOOLSDB)/golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver:
$(TOOLSBIN)/golangci-lint:

.PHONY: golangci-lint
golangci-lint: $(TOOLSBIN)/golangci-lint
	golangci-lint run
	@echo ''
## </golangci-lint>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.17.0
$(TOOLSDB)/goimports.$(GOIMPORTS_VER).$(GO_VER).ver:
$(TOOLSBIN)/goimports:

.PHONY: goimports
goimports: $(TOOLSBIN)/goimports
	@echo '$(TOOLSBIN)/goimports -l `$(FOLDERS)`'
	@if [[ -n "$$(goimports -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'goimports errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make goimports.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make goimports.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: goimports.display
goimports.display: $(TOOLSBIN)/goimports
	goimports -d `$(FOLDERS)`

.PHONY: goimports.fix
goimports.fix: $(TOOLSBIN)/goimports
	goimports -w `$(FOLDERS)`
## </goimports>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.5.0
$(TOOLSDB)/gofumpt.$(GOFUMPT_VER).$(GO_VER).ver:
$(TOOLSBIN)/gofumpt:

.PHONY: gofumpt
gofumpt: $(TOOLSBIN)/gofumpt
	@echo '$(TOOLSBIN)/gofumpt -l `$(FOLDERS)`'
	@if [[ -n "$$(gofumpt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofumpt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make gofumpt.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make gofumpt.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofumpt.display
gofumpt.display:
	gofumpt -d `$(FOLDERS)`

.PHONY: gofumpt.fix
gofumpt.fix:
	gofumpt -w `$(FOLDERS)`
## </gofumpt>

## <gofmt>
.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `$(FOLDERS)`'
	@if [[ -n "$$(gofmt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofmt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    make gofmt.display'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    make gofmt.fix'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofmt.display
gofmt.display:
	gofmt -d `$(FOLDERS)`

.PHONY: gofmt.fix
gofmt.fix:
	gofmt -w `$(FOLDERS)`
## </gofmt>

## <gojq>
# https://github.com/itchyny/gojq/releases
GOJQ_CMD := github.com/itchyny/gojq/cmd/gojq
GOJQ_VER := v0.12.14
$(TOOLSBIN)/gojq:
$(TOOLSBIN)/.gojq.$(GOJQ_VER).$(GO_VER).ver:

.PHONY: gojq
gojq: $(TOOLSBIN)/gojq
## </gojq>

## <air>
# https://github.com/cosmtrek/air/releases
AIR_CMD:=github.com/cosmtrek/air
AIR_VER:=v1.49.0
$(TOOLSDB)/air.$(AIR_VER).$(GO_VER).ver:
$(TOOLSBIN)/air:

.PHONY: air
air: $(TOOLSBIN)/air
	$(TOOLSBIN)/air -c .air.toml
## </air>

## <protobuf>
# https://github.com/protocolbuffers/protobuf/releases
PROTOC_VER:=v25.2
$(TOOLSDB)/protoc.$(PROTOC_VER).ver:
$(TOOLSBIN)/protoc: $(TOOLSDB)/protoc.$(PROTOC_VER).ver
	./scripts/install-protoc --version $(PROTOC_VER) --destination $(TOOLSDIR)

# https://github.com/protocolbuffers/protobuf-go/releases
PROTOC-GEN-GO_CMD := google.golang.org/protobuf/cmd/protoc-gen-go
PROTOC-GEN-GO_VER := v1.32.0
$(TOOLSBIN)/protoc-gen-go:
$(TOOLSBIN)/.protoc-gen-go.$(PROTOBUFGO_VER).$(GO_VER).ver:

.PHONY: proto
proto: $(TOOLSBIN)/protoc $(TOOLSBIN)/protoc-gen-go
	$(TOOLSBIN)/protoc --version
	$(TOOLSBIN)/protoc-gen-go --version
## </protobuf>
####################################################################################
## </ci & external tools> ##########################################################
####################################################################################


# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types
