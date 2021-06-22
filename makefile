GO_TEST=go test
COVERAGE_PACKAGES=github.com/davinash/yados/...

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

build:
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yados$(EXE)  cmd/main.go

test:
	GOFLAGS="-count=1" GO111MODULE=on $(GO_TEST) -timeout 50m github.com/davinash/yados/tests -v

test-single:
	GOFLAGS="-count=1" GO111MODULE=on go test -v github.com/davinash/yados/tests -testify.m $(TEST_NAME)

clean:
	rm -rf out/*
