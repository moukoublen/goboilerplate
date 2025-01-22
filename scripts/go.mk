MODULE := $(shell cat go.mod | grep -e "^module" | sed "s/^module //")

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

# https://pkg.go.dev/cmd/go/internal/test
.PHONY: test
test:
	CGO_ENABLED=1 go test -timeout 30s -tags '$(TAGS)' -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-n-read
test-n-read: test
	@go tool cover -func coverage.txt

.PHONY: bench
bench: # runs all benchmarks
	CGO_ENABLED=1 go test -benchmem -run=^Benchmark$$ -mod=readonly -count=1 -v -race -bench=. ./...


.PHONY: run
run:
	go run ./cmd/goboilerplate
