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

GO_PACKAGES = go list -tags='$(TAGS)' ./...
GO_FOLDERS = go list -tags='$(TAGS)' -f '{{ .Dir }}' ./...
GO_FILES = find . -type f -name '*.go' -not -path './vendor/*'

export GO111MODULE := on
#export GOFLAGS := -mod=vendor
GOPATH := $(shell go env GOPATH)
GO_VER := $(shell go env GOVERSION)
BUILD_OUTPUT ?= $(CURDIR)/output


.PHONY: mod
mod:
	go mod tidy -go=1.23
	go mod verify

# https://go.dev/ref/mod#go-get
# -u flag tells go get to upgrade modules
# -t flag tells go get to consider modules needed to build tests of packages named on the command line.
# When -t and -u are used together, go get will update test dependencies as well.
.PHONY: go-deps-upgrade
go-deps-upgrade:
	go get -u -t ./...
	go mod tidy -go=1.23

# https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
# https://pkg.go.dev/cmd/compile
# https://pkg.go.dev/cmd/link

#BUILD_FLAGS := -mod=vendor -a -ldflags "-s -w $(X_FLAGS) -extldflags='-static'" -tags '$(TAGS)'
BUILD_FLAGS := -a -ldflags '-s -w $(X_FLAGS)' -tags '$(TAGS)'
BUILD_FLAGS_DEBUG := -ldflags '$(X_FLAGS)' -tags '$(TAGS)'

cmd.%: CMDNAME=$*
cmd.%:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT)/$(CMDNAME) ./cmd/$(CMDNAME)

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
	CGO_ENABLED=1 go test -timeout 60s -race -tags="$(TAGS)" -coverprofile cover.out -covermode atomic ./...
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: run
run:
	go run ./cmd/goboilerplate
