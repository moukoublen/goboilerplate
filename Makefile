SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## this requires make >=4.4
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
	$(GO_EXEC) mod tidy -go=1.19
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags "$(TAGS)"
BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS)" -tags "$(TAGS)"

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

.PHONY: build
build: $(shell ls -d cmd/* | sed -e 's/\//./')
cmd.%: CMDNAME=$*
cmd.%:
	CGO_ENABLED=0 $(GO_EXEC) build $(BUILD_FLAGS) -o ./output/$(CMDNAME) ./cmd/$(CMDNAME)

.PHONY: clean
clean:
	rm -rf ./output

.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd
# man git-clean

.PHONY: image
image:
	$(DOCKER_EXEC) build . -f $(CURDIR)/build/docker/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: up
up:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml up --force-recreate --build

.PHONY: down
down:
	$(DOCKER_EXEC) compose -f $(CURDIR)/deployments/compose/docker-compose.yml down --volumes --rmi local --remove-orphans

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


###############################################################################
## ci & external tools #############################################################
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

TOOLSBIN ?= $(shell pwd)/.bin
export TOOLSBIN


uppercase = $(shell echo '$(1)' | tr '[:lower:]' '[:upper:]')


tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/staticcheck \
	$(TOOLSBIN)/golangci-lint \
	$(TOOLSBIN)/gofumpt

$(TOOLSBIN):
	@mkdir -p $(TOOLSBIN)

.PRECIOUS: $(TOOLSBIN)/.%.ver
$(TOOLSBIN)/.%.ver: | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/.$(word 1,$(subst ., ,$*)).*
	@touch $(TOOLSBIN)/.$*.ver

$(TOOLSBIN)/%: DSC=$*
$(TOOLSBIN)/%: VER=$($(call uppercase,$*)_VER)
$(TOOLSBIN)/%: CMD=$($(call uppercase,$*)_CMD)
$(TOOLSBIN)/%: $(TOOLSBIN)/.$$(DSC).$$(VER).$(GO_VER).ver
	@echo -e "Installing \e[1;36m$(DSC)\e[0m@\e[1;36m$(VER)\e[0m using \e[1;36m$(GO_VER)\e[0m"
	CGO_ENABLED=0 GOBIN="$(TOOLSBIN)" $(GO_EXEC) install -trimpath -ldflags '-s -w -extldflags "-static"' "$(CMD)@$(VER)"
	@echo ""

## <staticcheck>
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/c
STATICCHECK_CMD=honnef.co/go/tools/cmd/staticcheck
STATICCHECK_VER:=v0.3.3
$(TOOLSBIN)/staticcheck:
$(TOOLSBIN)/.staticcheck.$(STATICCHECK_VER).$(GO_VER).ver: # force not intermediate. In make >= 4.4. .NOTINTERMEDIATE will do the job.

.PHONY: staticcheck
staticcheck: $(TOOLSBIN)/staticcheck
	$(TOOLSBIN)/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
## </staticcheck>

## <golangci-lint>
# https://github.com/golangci/golangci-lint/releases
GOLANGCI-LINT_CMD:=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI-LINT_VER:=v1.51.0
$(TOOLSBIN)/golangci-lint:
$(TOOLSBIN)/.golangci-lint.$(GOLANGCI-LINT_VER).$(GO_VER).ver: # force not intermediate. In make >= 4.4. .NOTINTERMEDIATE will do the job.

.PHONY: golangci-lint
golangci-lint: $(TOOLSBIN)/golangci-lint
	$(TOOLSBIN)/golangci-lint run
	@echo ''
## </golangci-lint>

## <gofumpt>
# https://github.com/mvdan/gofumpt/releases
GOFUMPT_CMD:=mvdan.cc/gofumpt
GOFUMPT_VER:=v0.4.0
$(TOOLSBIN)/gofumpt:
$(TOOLSBIN)/.gofumpt.$(GOFUMPT_VER).$(GO_VER).ver: # force not intermediate. In make >= 4.4. .NOTINTERMEDIATE will do the job.

.PHONY: gofumpt
gofumpt: $(TOOLSBIN)/gofumpt
	@echo '$(TOOLSBIN)/gofumpt -l `$(FOLDERS)`'
	@if [[ -n "$$($(TOOLSBIN)/gofumpt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	$(TOOLSBIN)/gofumpt -d `$(FOLDERS)`

.PHONY: gofumpt.fix
gofumpt.fix:
	$(TOOLSBIN)/gofumpt -w `$(FOLDERS)`
## </gofumpt>

## <goimports>
# https://pkg.go.dev/golang.org/x/tools?tab=versions
GOIMPORTS_CMD := golang.org/x/tools/cmd/goimports
GOIMPORTS_VER := v0.5.0
$(TOOLSBIN)/goimports:
$(TOOLSBIN)/.goimports.$(GOIMPORTS_VER).$(GO_VER).ver: # force not intermediate. In make >= 4.4. .NOTINTERMEDIATE will do the job.

.PHONY: goimports
goimports: $(TOOLSBIN)/goimports
	@echo '$(TOOLSBIN)/goimports -l `$(FOLDERS)`'
	@if [[ -n "$$($(TOOLSBIN)/goimports -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
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
	$(TOOLSBIN)/goimports -d `$(FOLDERS)`

.PHONY: goimports.fix
goimports.fix: $(TOOLSBIN)/goimports
	$(TOOLSBIN)/goimports -w `$(FOLDERS)`
## </goimports>

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
###############################################################################
###############################################################################


# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types
