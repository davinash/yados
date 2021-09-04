#!/usr/bin/env bash

set -o errexit
set -o nounset

# shellcheck disable=SC1068
GOPATH=$(go env GOPATH)

if test ! -f "${GOPATH}/bin/golangci-lint"; then
  echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.40.1
fi

if test ! -f "${GOPATH}/bin/protoc-gen-go"; then
  echo "Installing protoc-gen-go" && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if test ! -f "${GOPATH}/bin/protoc-gen-go-grpc"; then
  echo "Installing protoc-gen-go-grpc" && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

if test ! -f "${GOPATH}/bin/goimports"; then
  echo "Installing goimports" && go install golang.org/x/tools/cmd/goimports@latest
fi

if test ! -f "${GOPATH}/bin/errcheck"; then
  echo "Installing errcheck" && go install github.com/kisielk/errcheck@latest
fi

if test ! -f "${GOPATH}/bin/swag"; then
  echo "Installing swag" && go install github.com/swaggo/swag/cmd/swag@latest
fi

