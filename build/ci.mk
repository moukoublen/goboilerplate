.PHONY: checks
checks: vet staticcheck gofumpt goimports

.PHONY: image-ci
image-ci:
	$(DOCKER) build . -f $(CURDIR)/build/ci/Dockerfile -t ${NAME}-ci:latest

.PHONY: vet
vet:
	$(GO_EXEC) vet `${PACKAGES}`
	@echo ""

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

.PHONY: golangci-lint
golangci-lint:
	${TOOLSBIN}/golangci-lint run
	@echo ''

.PHONY: staticcheck
staticcheck:
	${TOOLSBIN}/staticcheck -f=stylish -checks=all,-ST1000 -tests ./...
	@echo ''

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="${TAGS}" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out
