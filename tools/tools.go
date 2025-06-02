//go:build tools

package main

import (
	_ "github.com/air-verse/air"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/itchyny/gojq/cmd/gojq"
	_ "github.com/vektra/mockery/v3"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "mvdan.cc/gofumpt"
	_ "mvdan.cc/sh/v3/cmd/shfmt"
)

func main() {
}
