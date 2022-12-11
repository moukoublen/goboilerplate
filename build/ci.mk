# https://www.gnu.org/software/make/manual/make.html#Automatic-Variables
# https://www.gnu.org/software/make/manual/make.html#Prerequisite-Types

TOOLSBIN ?= $(shell pwd)/.bin
export TOOLSBIN

PACKAGES = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor ./...
FOLDERS = $(GO_EXEC) list -tags=$(TAGS) -mod=vendor -f '{{.Dir}}' ./...

GO_VER := $(shell go env GOVERSION)
GOPATH := $(shell go env GOPATH)

# https://pkg.go.dev/cmd/link
# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

define install-log
	@echo -e "Installing \e[1;36m$(1)\e[0m..."
endef

define go-install
	$(GO_EXEC) install -a -trimpath -ldflags '-s -w -extldflags "-static"' "$(1)"
	@echo ""
endef

tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/staticcheck \
	$(TOOLSBIN)/golangci-lint \
	$(TOOLSBIN)/gofumpt \
	$(TOOLSBIN)/migrate

$(TOOLSBIN):
	@mkdir -p $(TOOLSBIN)

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint

.PHONY: vet
vet:
	$(GO_EXEC) vet `$(PACKAGES)`
	@echo ""


########## goimports ##########################################################
# https://pkg.go.dev/golang.org/x/tools?tab=versions
VERSION_GOIMPORTS := v0.4.0
VERSION_FILE_GOIMPORTS := ._goimports_$(VERSION_GOIMPORTS)_$(GO_VER)

$(TOOLSBIN)/$(VERSION_FILE_GOIMPORTS): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._goimports_*
	@touch $(TOOLSBIN)/$(VERSION_FILE_GOIMPORTS)

$(TOOLSBIN)/goimports: $(TOOLSBIN)/$(VERSION_FILE_GOIMPORTS)
	$(call install-log,goimports)
	$(call go-install,golang.org/x/tools/cmd/goimports@$(VERSION_GOIMPORTS))
	@cp $(GOPATH)/bin/goimports $(TOOLSBIN)/goimports

.PHONY: goimports
goimports: $(TOOLSBIN)/goimports
	@echo '$(TOOLSBIN)/goimports -l `$(FOLDERS)`'
	@if [[ -n "$$($(TOOLSBIN)/goimports -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'goimports errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    goimports -d `$(FOLDERS)`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    goimports -w `$(FOLDERS)`'; \
		echo '  or'; \
		echo '    make goimports-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: goimports-w
goimports-w:
	$(TOOLSBIN)/goimports -w `$(FOLDERS)`
###############################################################################


########## staticcheck ########################################################
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/
VERSION_STATICCHECK := v0.3.3
VERSION_FILE_STATICCHECK := ._staticcheck_$(VERSION_STATICCHECK)_$(GO_VER)

$(TOOLSBIN)/$(VERSION_FILE_STATICCHECK): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._staticcheck_*
	@touch $(TOOLSBIN)/$(VERSION_FILE_STATICCHECK)

$(TOOLSBIN)/staticcheck: $(TOOLSBIN)/$(VERSION_FILE_STATICCHECK)
	$(call install-log,staticcheck)
	$(call go-install,honnef.co/go/tools/cmd/staticcheck@$(VERSION_STATICCHECK))
	@cp $(GOPATH)/bin/staticcheck $(TOOLSBIN)/staticcheck

.PHONY: staticcheck
staticcheck: $(TOOLSBIN)/staticcheck
	$(TOOLSBIN)/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
###############################################################################


########## golangci-lint ######################################################
# https://github.com/golangci/golangci-lint/releases
VERSION_GOLANGCILINT := v1.50.1
VERSION_FILE_GOLANGCILINT := ._golangci-lint_$(VERSION_GOLANGCILINT)_$(GO_VER)

$(TOOLSBIN)/$(VERSION_FILE_GOLANGCILINT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._golangci-lint_*
	@touch $(TOOLSBIN)/$(VERSION_FILE_GOLANGCILINT)

$(TOOLSBIN)/golangci-lint: $(TOOLSBIN)/$(VERSION_FILE_GOLANGCILINT)
	$(call install-log,golangci-lint)
	$(call go-install,github.com/golangci/golangci-lint/cmd/golangci-lint@$(VERSION_GOLANGCILINT))
	@cp $(GOPATH)/bin/golangci-lint $(TOOLSBIN)/golangci-lint

.PHONY: golangci-lint
golangci-lint: $(TOOLSBIN)/golangci-lint
	$(TOOLSBIN)/golangci-lint run
	@echo ''
###############################################################################


########## gofumpt ############################################################
# https://github.com/mvdan/gofumpt/releases
VERSION_GOFUMPT := v0.4.0
VERSION_FILE_GOFUMPT := ._gofumpt_$(VERSION_GOFUMPT)_$(GO_VER)

$(TOOLSBIN)/$(VERSION_FILE_GOFUMPT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._gofumpt_*
	@touch $(TOOLSBIN)/$(VERSION_FILE_GOFUMPT)

$(TOOLSBIN)/gofumpt: $(TOOLSBIN)/$(VERSION_FILE_GOFUMPT)
	$(call install-log,gofumpt)
	$(call go-install,mvdan.cc/gofumpt@$(VERSION_GOFUMPT))
	@cp $(GOPATH)/bin/gofumpt $(TOOLSBIN)/gofumpt

.PHONY: gofumpt
gofumpt: $(TOOLSBIN)/gofumpt
	@echo '$(TOOLSBIN)/gofumpt -l `$(FOLDERS)`'
	@if [[ -n "$$($(TOOLSBIN)/gofumpt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofumpt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    gofumpt -d `$(FOLDERS)`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    gofumpt -w `$(FOLDERS)`'; \
		echo '  or'; \
		echo '    make gofumpt-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofumpt-w
gofumpt-w:
	$(TOOLSBIN)/gofumpt -w `$(FOLDERS)`
###############################################################################


########## golang-migrate #####################################################
# https://github.com/golang-migrate/migrate/releases
VERSION_MIGRATE := 4.15.2
VERSION_FILE_MIGRATE := ._migrate_$(VERSION_MIGRATE)_$(GO_VER)

$(TOOLSBIN)/$(VERSION_FILE_MIGRATE): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._migrate_*
	@touch $(TOOLSBIN)/$(VERSION_FILE_MIGRATE)

$(TOOLSBIN)/migrate: $(TOOLSBIN)/$(VERSION_FILE_MIGRATE)
	$(call install-log,golang-migrate,$(VERSION_MIGRATE))
	@./scripts/install-migrate "$(VERSION_MIGRATE)" "$(TOOLSBIN)"
###############################################################################


########## gofmt ##############################################################
.PHONY: gofmt
gofmt:
	@echo 'gofmt -l `$(FOLDERS)`'
	@if [[ -n "$$(gofmt -l `$(FOLDERS)` | tee /dev/stderr)" ]]; then \
		echo 'gofmt errors'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To display the needed changes run:"; \
		echo '    gofmt -d `$(FOLDERS)`'; \
		echo ''; \
		echo -e "\e[0;34m→\e[0m To fix them run:"; \
		echo '    gofmt -w `$(FOLDERS)`'; \
		echo '  or'; \
		echo '    make gofmt-w'; \
		echo ''; \
		exit 1; \
	fi
	@echo ''

.PHONY: gofmt-w
gofmt-w:
	gofmt -w `$(FOLDERS)`
###############################################################################

