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


### goimports ###
# https://pkg.go.dev/golang.org/x/tools?tab=versions
VERSION_GOIMPORTS := 0.1.8

$(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._goimports_*
	@touch $(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS)

$(TOOLSBIN)/goimports: $(TOOLSBIN)/._goimports_$(VERSION_GOIMPORTS)
	$(call install-log,goimports)
	$(call go-install,golang.org/x/tools/cmd/goimports@v${VERSION_GOIMPORTS})
	@cp ${GOPATH}/bin/goimports $(TOOLSBIN)/goimports
#############################################################################


### staticcheck ###
# https://github.com/dominikh/go-tools/releases    https://staticcheck.io/
VERSION_STATICCHECK := 2021.1.2

$(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._staticcheck_*
	@touch $(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK)

$(TOOLSBIN)/staticcheck: $(TOOLSBIN)/._staticcheck_$(VERSION_STATICCHECK)
	$(call install-log,staticcheck)
	$(call go-install,honnef.co/go/tools/cmd/staticcheck@$(VERSION_STATICCHECK))
	@cp ${GOPATH}/bin/staticcheck $(TOOLSBIN)/staticcheck
#############################################################################


### golangci-lint ###
# https://github.com/golangci/golangci-lint/releases
VERSION_GOLANGCILINT := 1.43.0

$(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._golangci-lint_*
	@touch $(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT)

$(TOOLSBIN)/golangci-lint: $(TOOLSBIN)/._golangci-lint_$(VERSION_GOLANGCILINT)
	$(call install-log,golangci-lint)
	$(call go-install,github.com/golangci/golangci-lint/cmd/golangci-lint@v$(VERSION_GOLANGCILINT))
	@cp ${GOPATH}/bin/golangci-lint $(TOOLSBIN)/golangci-lint
#############################################################################


### gofumpt ###
# https://github.com/mvdan/gofumpt/releases
VERSION_GOFUMPT := 0.2.1

$(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._gofumpt_*
	@touch $(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT)

$(TOOLSBIN)/gofumpt: $(TOOLSBIN)/._gofumpt_$(VERSION_GOFUMPT)
	$(call install-log,gofumpt)
	$(call go-install,mvdan.cc/gofumpt@v$(VERSION_GOFUMPT))
	@cp ${GOPATH}/bin/gofumpt $(TOOLSBIN)/gofumpt
#############################################################################


### golang-migrate ###
# https://github.com/golang-migrate/migrate/releases
VERSION_MIGRATE := 4.15.1

$(TOOLSBIN)/._migrate_$(VERSION_MIGRATE): | $(TOOLSBIN)
	@rm -f $(TOOLSBIN)/._migrate_*
	@touch $(TOOLSBIN)/._migrate_$(VERSION_MIGRATE)

$(TOOLSBIN)/migrate: $(TOOLSBIN)/._migrate_$(VERSION_MIGRATE)
	$(call install-log,golang-migrate)
	@./scripts/install-migrate "$(VERSION_MIGRATE)" "$(TOOLSBIN)"
#############################################################################

