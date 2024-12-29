GO_EXEC ?= go
export GO_EXEC

REPO_GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)
REPO_GIT_COMMIT = $(shell git rev-parse HEAD 2>/dev/null || true)
REPO_GIT_COMMIT_SHORT = $(shell git rev-parse --short HEAD 2>/dev/null || true)
REPO_GIT_TAG = $(shell git describe --tags 2>/dev/null || true)

MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")
X_FLAGS = \
		-X '$(MODULE)/build.Branch=$(REPO_GIT_BRANCH)' \
		-X '$(MODULE)/build.Commit=$(REPO_GIT_COMMIT)' \
		-X '$(MODULE)/build.CommitShort=$(REPO_GIT_COMMIT_SHORT)' \
		-X '$(MODULE)/build.Tag=$(REPO_GIT_TAG)'

GO_PACKAGES = $(GO_EXEC) list -tags='$(TAGS)' -mod=vendor ./...
GO_FOLDERS = $(GO_EXEC) list -tags='$(TAGS)' -mod=vendor -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go' -not -path './vendor/*'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)
BUILD_OUTPUT ?= $(CURDIR)/output


.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -go=1.23
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor

# https://go.dev/ref/mod#go-get
# -u flag tells go get to upgrade modules
# -t flag tells go get to consider modules needed to build tests of packages named on the command line.
# When -t and -u are used together, go get will update test dependencies as well.
.PHONY: go-deps-upgrade
go-deps-upgrade:
	$(GO_EXEC) get -u -t ./...
	$(GO_EXEC) mod tidy -go=1.23
	$(GO_EXEC) mod vendor

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags '$(TAGS)'
BUILD_FLAGS := -mod=vendor -a -ldflags '-s -w $(X_FLAGS)' -tags '$(TAGS)'
BUILD_FLAGS_DEBUG := -mod=vendor -ldflags '$(X_FLAGS)' -tags '$(TAGS)'

cmd.%: CMDNAME=$*
cmd.%:
	$(GO_EXEC) env
	@echo ''
	CGO_ENABLED=0 $(GO_EXEC) build $(BUILD_FLAGS) -o $(BUILD_OUTPUT)/$(CMDNAME) ./cmd/$(CMDNAME)

dbg.%: BUILD_FLAGS=$(BUILD_FLAGS_DEBUG)
dbg.%: cmd.%
	@echo "debug binary done"

.PHONY: build
build: $(shell ls -d cmd/* | sed -e 's/\//./')

.PHONY: clean
clean:
	rm -rf $(BUILD_OUTPUT)

.PHONY: test
test:
	CGO_ENABLED=1 $(GO_EXEC) test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@$(GO_EXEC) tool cover -func cover.out
	@rm cover.out

.PHONY: run
run:
	$(GO_EXEC) run -mod=vendor ./cmd/goboilerplate
