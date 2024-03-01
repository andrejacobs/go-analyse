#!/bin/bash
set -eu
set -o pipefail

# Install golangci-lint
# https://golangci-lint.run/usage/install/#local-installation
# binary will be $(go env GOPATH)/bin/golangci-lint
GOLANGCI_LINT_VERSION="v1.56.2"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

golangci-lint --version
