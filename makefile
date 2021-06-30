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
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yados$(EXE)     cmd/yados/main.go
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yadosctl$(EXE)  cmd/cli/main.go

test:
	go test github.com/davinash/yados/... -v

clean:
	rm -rf out/*
