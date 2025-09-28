# Makefile for the DigiSocialBlock Project
# This file automates common development and build tasks.

# Define Go package paths
PKG_DIR := ./pkg

# Define Go commands
GO := go

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
	@echo "  test         - Run all Go tests"
	@echo "  tidy         - Tidy Go modules"

.PHONY: test
test:
	@echo ">> Running Go tests..."
	cd $(PKG_DIR) && $(GO) test ./...

.PHONY: tidy
tidy:
	@echo ">> Tidying Go modules..."
	cd $(PKG_DIR) && $(GO) mod tidy

## --------------------------------------
## CI/CD Targets
## --------------------------------------

.PHONY: ci
ci: tidy test
	@echo ">> CI checks passed."