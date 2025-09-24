# Makefile for the DigiSocialBlock Project
# This file automates common development and build tasks.

# Define Go package paths
PKG_DIR := ./pkg
PROTO_DIR := ./proto

# Define Go commands
GO := go
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc

# Default target
.DEFAULT_GOAL := help

## --------------------------------------
## Development Targets
## --------------------------------------

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto-gen    - Generate Go code from Protocol Buffer definitions"
	@echo "  test         - Run all Go tests"
	@echo "  lint         - Lint the Go codebase"
	@echo "  tidy         - Tidy Go modules"

.PHONY: proto-gen
proto-gen:
	@echo ">> Generating Go code from .proto files..."
	$(PROTOC) --go_out=$(PKG_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo ">> Protocol Buffers generated successfully."

.PHONY: test
test:
	@echo ">> Running Go tests..."
	cd $(PKG_DIR) && $(GO) test ./...

.PHONY: lint
lint:
	@echo ">> Linting Go codebase..."
	cd $(PKG_DIR) && golangci-lint run

.PHONY: tidy
tidy:
	@echo ">> Tidying Go modules..."
	cd $(PKG_DIR) && $(GO) mod tidy

## --------------------------------------
## CI/CD Targets
## --------------------------------------

.PHONY: ci
ci: tidy lint test
	@echo ">> CI checks passed."