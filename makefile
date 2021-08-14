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

TEST_COUNT=1

build: getdeps format lint buildx

format:
	go mod tidy
	go fmt ./...
	${GOPATH}/bin/goimports -l -w .
	${GOPATH}/bin/errcheck -ignoretests -blank ./...

buildx:
	@echo "Building the product"
	GO111MODULE=on go build -ldflags "-X main.Version=$(VERSION)" -o out/yadosctl$(EXE)  cmd/cli/main.go

getdeps:
	@echo "Checking dependencies"
	@mkdir -p ${GOPATH}/bin 
	sh ./deps.sh

generate: getdeps
	@echo "Generating protobuf/grpc resources"
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@(PATH=$(GOPATH)/bin:$(PATH) && go generate internal/proto/generate.go)

lint:
	@echo "Running $@ check"
	GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean
	GO111MODULE=on ${GOPATH}/bin/golangci-lint run --build-tags kqueue --timeout=10m --skip-dirs internal/proto/gen --config ./.golangci.yml
	@echo "Running vet"
	go list ./... | grep -v gen | xargs go vet

test:
	GOFLAGS="-count=1" go test github.com/davinash/yados/tests -v -failfast

test-with-cover:
	GOFLAGS="-count=1" go test -coverprofile=tmp.cov --coverpkg $(COVERAGE_PACKAGES) github.com/davinash/yados/tests -v -failfast

test-single:
	go test -v github.com/davinash/yados/tests -testify.m $(TEST_NAME) -count $(TEST_COUNT)

clean:
	rm -rf out/*

