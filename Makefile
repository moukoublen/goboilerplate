SHELL := /bin/bash

.NOTPARALLEL:
.SECONDEXPANSION:
## NOTINTERMEDIATE requires make >=4.4
.NOTINTERMEDIATE:

include $(CURDIR)/scripts/go.mk
include $(CURDIR)/scripts/docker.mk
include $(CURDIR)/scripts/tools.mk

.DEFAULT_GOAL=default
.PHONY: default
default: checks build

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

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

.PHONY: checks
checks: vet staticcheck gofumpt goimports golangci-lint-github-actions

.PHONY: ci-gen-n-format
ci-gen-n-format: goimports gofumpt
	./scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	./scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt shellcheck
	@./scripts/git-check-dirty
