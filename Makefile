.DEFAULT_GOAL := help

# ==============================================================================
# Build options

VERSION      = $(shell git describe --tags --always)
OUTPUT_DIR   = ./_output
PROTO_DIR    = ./api/apiserver/proto
PROTO_FILES  = $(shell find ${PROTO_DIR} -name *.proto)
GO_LDFLAGS  += -X main.Version=$(VERSION)
MAKEFLAGS   += --no-print-directory

# ==============================================================================
# Includes

include make-rules/tools.mk

# ==============================================================================
# Targets

## all: Build all.
.PHONY: all
all: lint api error test build

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: tools.verify.golangci-lint
	go mod tidy -compat=1.18
	golangci-lint run ./...

## api: Generate api proto.
.PHONY: api
api: tools.verify.protoc tools.verify.protoc-gen-go tools.verify.protoc-gen-go-grpc
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=paths=source_relative:$(PROTO_DIR) \
	       --go-grpc_out=paths=source_relative:$(PROTO_DIR) \
	       $(PROTO_FILES)

## error: Generate error code.
.PHONY: error
error: tools.verify.codegen
	codegen ./internal/pkg/code
	codegen -doc -output ./docs/error_code.md ./internal/pkg/code

## test: Run unit test.
.PHONY: test
test:
	@-mkdir -p $(OUTPUT_DIR)
	go test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out ./...

## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out
	go tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

## build: Build source code for host platform.
.PHONY: build
build:
	go build -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/ ./...

## update: Update all modules.
.PHONY: update
update:
	go get -u ./...
	go mod tidy -compat=1.18

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	-rm -vrf $(OUTPUT_DIR)

## docker: Docker build
.PHONY: docker
docker:
	docker build --build-arg VERSION=#{VERSION} -t chekwas/iam-apiserver:${VERSION} .
	docker push chekwas/iam-apiserver:${VERSION}

## help: Show help info.
.PHONY: help
help: Makefile
	@echo "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
