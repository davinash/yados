#!/usr/bin/env bash

set -o errexit
set -o nounset

# shellcheck disable=SC1068
GOPATH=$(go env GOPATH)

if test ! -f "${GOPATH}/bin/golangci-lint"; then
  echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.40.1
fi

if test ! -f "${GOPATH}/bin/protoc-gen-go"; then
  echo "Installing protoc-gen-go" && go get -u google.golang.org/protobuf/cmd/protoc-gen-go
fi

if test ! -f "${GOPATH}/bin/protoc-gen-go-grpc"; then
  echo "Installing protoc-gen-go-grpc" && go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
fi
