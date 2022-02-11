###############################################################################
# Requirements:
#    TOOLSBIN  must be defined to a directory path
###############################################################################

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="${TAGS}" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: image-ci
image-ci:
	$(DOCKER) build . -f $(CURDIR)/build/ci/Dockerfile -t ${NAME}-ci:latest

.PHONY: checks
checks: vet staticcheck gofumpt goimports

.PHONY: vet
vet:
	$(GO_EXEC) vet `${PACKAGES}`
	@echo ""

define install-log
	@echo -e "Installing \e[1;36m${1}\e[0m..."
endef

define go-install
	@$(GO_EXEC) install -a -trimpath -ldflags '-s -w -extldflags "-static"' "${1}"
endef

tools: \
	$(TOOLSBIN)/goimports \
	$(TOOLSBIN)/staticcheck \
	$(TOOLSBIN)/golangci-lint \
	$(TOOLSBIN)/gofumpt \
	$(TOOLSBIN)/migrate


$(TOOLSBIN):
	@mkdir -p $(TOOLSBIN)


########## goimports ##########################################################
# https://pkg.go.dev/golang.org/x/tools?tab=versions
VERSION_GOIMPORTS := 0.1.9

$(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._goimports_*
	@touch $(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS)

$(TOOLSBIN)/goimports: $(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS)
	$(call install-log,goimports)
	$(call go-install,golang.org/x/tools/cmd/goimports@v${VERSION_GOIMPORTS})
	@cp ${GOPATH}/bin/goimports $(TOOLSBIN)/goimports

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
###############################################################################


########## staticcheck ########################################################
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/
VERSION_STATICCHECK := 2021.1.2

$(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._staticcheck_*
	@touch $(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK)

$(TOOLSBIN)/staticcheck: $(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK)
	$(call install-log,staticcheck)
	$(call go-install,honnef.co/go/tools/cmd/staticcheck@$(VERSION_STATICCHECK))
	@cp ${GOPATH}/bin/staticcheck $(TOOLSBIN)/staticcheck

.PHONY: staticcheck
staticcheck:
	${TOOLSBIN}/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''
###############################################################################


########## golangci-lint ######################################################
# https://github.com/golangci/golangci-lint/releases
VERSION_GOLANGCILINT := 1.44.0

$(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._golangci-lint_*
	@touch $(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT)

$(TOOLSBIN)/golangci-lint: $(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT)
	$(call install-log,golangci-lint)
	$(call go-install,github.com/golangci/golangci-lint/cmd/golangci-lint@v$(VERSION_GOLANGCILINT))
	@cp ${GOPATH}/bin/golangci-lint $(TOOLSBIN)/golangci-lint

.PHONY: golangci-lint
golangci-lint:
	${TOOLSBIN}/golangci-lint run
	@echo ''
###############################################################################


########## gofumpt ############################################################
# https://github.com/mvdan/gofumpt/releases
VERSION_GOFUMPT := 0.2.1

$(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._gofumpt_*
	@touch $(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT)

$(TOOLSBIN)/gofumpt: $(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT)
	$(call install-log,gofumpt)
	$(call go-install,mvdan.cc/gofumpt@v$(VERSION_GOFUMPT))
	@cp ${GOPATH}/bin/gofumpt $(TOOLSBIN)/gofumpt

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
###############################################################################


########## golang-migrate #####################################################
# https://github.com/golang-migrate/migrate/releases
VERSION_MIGRATE := 4.15.1

$(TOOLSBIN)/._migrate_$(VERSION_MIGRATE): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._migrate_*
	@touch $(TOOLSBIN)/._migrate_$(VERSION_MIGRATE)

$(TOOLSBIN)/migrate: $(TOOLSBIN)/._migrate_$(VERSION_MIGRATE)
	$(call install-log,golang-migrate)
	@./scripts/install-migrate "$(VERSION_MIGRATE)" "$(TOOLSBIN)"
###############################################################################


########## gofmt ##############################################################
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
###############################################################################

