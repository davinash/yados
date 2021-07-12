GO_TEST=go test
COVERAGE_PACKAGES=github.com/davinash/yados/...
GOPATH := $(shell go env GOPATH)

FILE=VERSION
VERSION=$(shell cat ${FILE})
NAME=nvctl

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

ifeq ($(OS), Windows_NT)
	EXE=.exe
else
	EXE=
endif

build: getdeps
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yados$(EXE)     cmd/yados/main.go
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yadosctl$(EXE)  cmd/cli/main.go

getdeps:
	@mkdir -p ${GOPATH}/bin 
	@sh ./deps.sh

generate: getdeps
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	(PATH=$(GOPATH)/bin:$(PATH) && go generate internal/proto/generate.go)

lint:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint run --build-tags kqueue --timeout=10m --config ./.golangci.yml

test:
	go test github.com/davinash/yados/... -v -count=1

clean:
	rm -rf out/*
