.PHONY: build build-stable install test help

INSTALL_PATH := $(shell go env GOPATH)/bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build tm binary to ./tm
	go build -o tm ./cmd/tm

build-stable: ## Build stable tm binary to ./tm-stable
	go build -o tm-stable ./cmd/tm

install: ## Install tm binary to GOPATH/bin
	go install ./cmd/tm
	@echo "Installed to: $(INSTALL_PATH)"

test: ## Run all tests
	go test ./...
