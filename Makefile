SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory

# Go parameters
GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint

# Project info
BINARY_NAME=deimos-client
MODULE_NAME=github.com/marsevilspirit/deimos-client

integration-test:
	@echo "Running integration tests..."
	@./integration-tests.sh

clean:
	@echo "Cleaning project..."
	@$(GOCLEAN)

fmt:
	@echo "Formatting code..."
	@$(GOFMT) ./...

lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint check"; \
		echo "To install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.55.2"; \
	fi

vet:
	@echo "Running go vet..."
	@$(GOCMD) vet ./...

install-lint:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	@echo "golangci-lint installed successfully!"

check: fmt vet lint
	@echo "All checks passed!"

help:
	@echo "Available commands:"
	@echo "  fmt              - Format code and tidy modules"
	@echo "  check            - Run all checks (fmt, vet, lint, test)"
	@echo "  clean            - Clean build artifacts"
	@echo "  integration-test - Run integration tests"
	@echo "  help             - Show this help message"

.PHONY: integration-test clean fmt lint vet check help
