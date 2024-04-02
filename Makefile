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

.DEFAULT_GOAL=default
.PHONY: default
default: checks build

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.22
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags '$(TAGS)'
BUILD_FLAGS := -mod=vendor -a -ldflags '-s -w $(X_FLAGS)' -tags '$(TAGS)'
BUILD_FLAGS_DEBUG := -mod=vendor -ldflags '$(X_FLAGS)' -tags '$(TAGS)'

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

BUILD_OUTPUT ?= $(CURDIR)/output

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
		--progress='plain' \
		.

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


####################################################################################
## <ci & external tools> ###########################################################
####################################################################################
.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint-github-actions

.PHONY: vet
vet:
	$(GO_EXEC) vet `$(GO_PACKAGES)`
	@echo ""

TOOLS_DIR ?= $(shell pwd)/.tools
TOOLS_DB ?= $(TOOLS_DIR)/.db
TOOLS_BIN ?= $(TOOLS_DIR)/bin
export PATH := $(TOOLS_BIN):$(PATH)

.PHONY: tools
tools: \
	$(TOOLS_BIN)/goimports \
	$(TOOLS_BIN)/staticcheck \
	$(TOOLS_BIN)/golangci-lint \
	$(TOOLS_BIN)/gofumpt \
	$(TOOLS_BIN)/gojq

.PHONY: clean-tools
clean-tools:
	rm -rf $(TOOLS_DIR)

$(TOOLS_BIN):
	@mkdir -p $(TOOLS_BIN)

$(TOOLS_DB):
	@mkdir -p $(TOOLS_DB)

# In make >= 4.4. .NOTINTERMEDIATE will do the job.
.PRECIOUS: $(TOOLS_DB)/%.ver
$(TOOLS_DB)/%.ver: | $(TOOLS_DB)
	@rm -f $(TOOLS_DB)/$(word 1,$(subst ., ,$*)).*
	@touch $(TOOLS_DB)/$*.ver

define go_install
	@echo -e "Installing \e[1;36m$(1)\e[0m@\e[1;36m$(3)\e[0m using \e[1;36m$(GO_VER)\e[0m"
	GOBIN="$(TOOLS_BIN)" CGO_ENABLED=0 $(GO_EXEC) install -trimpath -ldflags '-s -w -extldflags "-static"' "$(2)@$(3)"
	@echo ""
endef

## <staticcheck>
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/c
STATICCHECK_CMD:=honnef.co/go/tools/cmd/staticcheck
STATICCHECK_VER:=2023.1.7
$(TOOLS_BIN)/staticcheck: $(TOOLS_DB)/staticcheck.$(STATICCHECK_VER).$(GO_VER).ver
	$(call go_install,staticcheck,$(STATICCHECK_CMD),$(STATICCHECK_VER))

.PHONY: staticcheck
staticcheck: $(TOOLS_BIN)/staticcheck
	staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
## </staticcheck>

## <golangci-lint>
# https://github.com/golangci/golangci-lint/releases
GOLANGCI-LINT_CMD:=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI-LINT_VER:=v1.57.2
$(TOOLS_BIN)/golangci-lint: $(TOOLS_DB)/golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver
	$(call go_install,golangci-lint,$(GOLANGCI-LINT_CMD),$(GOLANGCI-LINT_VER))

.PHONY: golangci-lint
golangci-lint: $(TOOLS_BIN)/golangci-lint
	golangci-lint run
	@echo ''

.PHONY: golangci-lint-github-actions
golangci-lint-github-actions: $(TOOLS_BIN)/golangci-lint
	golangci-lint run --out-format github-actions
	@echo ''
## </golangci-lint>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.19.0
$(TOOLS_BIN)/goimports: $(TOOLS_DB)/goimports.$(GOIMPORTS_VER).$(GO_VER).ver
	$(call go_install,goimports,$(GOIMPORTS_CMD),$(GOIMPORTS_VER))

.PHONY: goimports
goimports: $(TOOLS_BIN)/goimports
	@echo '$(TOOLS_BIN)/goimports -l `$(GO_FILES)`'
	@if [[ -n "$$(goimports -l `$(GO_FILES)` | tee /dev/stderr)" ]]; then \
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
goimports.display: $(TOOLS_BIN)/goimports
	goimports -d `$(GO_FOLDERS)`

.PHONY: goimports.fix
goimports.fix: $(TOOLS_BIN)/goimports
	goimports -w `$(GO_FOLDERS)`
## </goimports>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.6.0
$(TOOLS_BIN)/gofumpt: $(TOOLS_DB)/gofumpt.$(GOFUMPT_VER).$(GO_VER).ver
	$(call go_install,gofumpt,$(GOFUMPT_CMD),$(GOFUMPT_VER))

.PHONY: gofumpt
gofumpt: $(TOOLS_BIN)/gofumpt
	@echo '$(TOOLS_BIN)/gofumpt -l `$(GO_FOLDERS)`'
	@if [[ -n "$$(gofumpt -l `$(GO_FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	gofumpt -d `$(GO_FOLDERS)`

.PHONY: gofumpt.fix
gofumpt.fix:
	gofumpt -w `$(GO_FOLDERS)`
## </gofumpt>

## <gofmt>
.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `$(GO_FOLDERS)`'
	@if [[ -n "$$(gofmt -l `$(GO_FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	gofmt -d `$(GO_FOLDERS)`

.PHONY: gofmt.fix
gofmt.fix:
	gofmt -w `$(GO_FOLDERS)`
## </gofmt>

## <gojq>
# https://github.com/itchyny/gojq/releases
GOJQ_CMD := github.com/itchyny/gojq/cmd/gojq
GOJQ_VER := v0.12.14
$(TOOLS_BIN)/gojq: $(TOOLS_DB)/gojq.$(GOJQ_VER).$(GO_VER).ver
	$(call go_install,gojq,$(GOJQ_CMD),$(GOJQ_VER))

.PHONY: gojq
gojq: $(TOOLS_BIN)/gojq
## </gojq>

## <air>
# https://github.com/cosmtrek/air/releases
AIR_CMD:=github.com/cosmtrek/air
AIR_VER:=v1.51.0
$(TOOLS_BIN)/air: $(TOOLS_DB)/air.$(AIR_VER).$(GO_VER).ver
	$(call go_install,air,$(AIR_CMD),$(AIR_VER))

.PHONY: air
air: $(TOOLS_BIN)/air
	@exec $(TOOLS_BIN)/air -c .air.toml
## </air>

## <protobuf>
# https://github.com/protocolbuffers/protobuf/releases
PROTOC_VER:=v26.0
$(TOOLS_BIN)/protoc: $(TOOLS_DB)/protoc.$(PROTOC_VER).ver
	./scripts/install-protoc --version $(PROTOC_VER) --destination $(TOOLS_DIR)

# https://github.com/protocolbuffers/protobuf-go/releases
PROTOC-GEN-GO_CMD:=google.golang.org/protobuf/cmd/protoc-gen-go
PROTOC-GEN-GO_VER:=v1.33.0
$(TOOLS_BIN)/protoc-gen-go: $(TOOLS_DB)/protoc-gen-go.$(PROTOC-GEN-GO_VER).$(GO_VER).ver
	$(call go_install,protoc-gen-go,$(PROTOC-GEN-GO_CMD),$(PROTOC-GEN-GO_VER))

.PHONY: proto
proto: $(TOOLS_BIN)/protoc $(TOOLS_BIN)/protoc-gen-go
	$(TOOLS_BIN)/protoc --version
	$(TOOLS_BIN)/protoc-gen-go --version
## </protobuf>
####################################################################################
## </ci & external tools> ##########################################################
####################################################################################


# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types

.PHONY: run
run: $(TOOLS_BIN)/gojq
	$(GO_EXEC) run -mod=vendor ./cmd/goboilerplate | $(TOOLS_BIN)/gojq
