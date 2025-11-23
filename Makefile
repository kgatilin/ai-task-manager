.PHONY: build build-stable build-headless build-all size-comparison install test help

INSTALL_PATH := $(shell go env GOPATH)/bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build full tm binary with TUI support to ./tm
	go build -o tm ./cmd/tm

build-headless: ## Build lightweight headless binary without TUI to ./tm-headless
	go build -tags headless -o tm-headless ./cmd/tm

build-all: build build-headless ## Build both full and headless binaries

size-comparison: build-all ## Compare binary sizes
	@echo "Binary sizes:"
	@ls -lh tm tm-headless | awk '{print $$5, $$9}'

build-stable: ## Build stable tm binary to ./tm-stable
	go build -o tm-stable ./cmd/tm

install: ## Install tm binary to GOPATH/bin
	go install ./cmd/tm
	@echo "Installed to: $(INSTALL_PATH)"

test: ## Run all tests
	go test ./...
