//go:build tools

// collection of development tools, see
// https://github.com/golang/go/issues/25922
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
