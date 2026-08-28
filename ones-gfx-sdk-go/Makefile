# ONES GFX SDK Go - Makefile

.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt

# Binary name (for future CLI)
BINARY_NAME=ones-gfx-sdk-mod
BINARY_DIR=bin

help: ## Show this help message
	@echo "ONES GFX SDK Go - Make targets"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

build-lib: ## Build the core library
	@echo "Building core library..."
	$(GOBUILD) ./ones_gfx
	$(GOBUILD) ./ones_gfx/resources
	@echo "✓ Library built successfully"

vet: ## Run go vet on all packages
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "✓ Vet passed"

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

check: vet ## Run static checks
	@echo "✓ All checks passed"

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	@echo "✓ Cleaned"

deps: ## Download dependencies
	$(GOGET) -v -d ./...

# Note: CLI binary build target will be added in v1.1
# build-cli: ## Build CLI binary (coming in v1.1)
# 	@echo "CLI build not yet available"

.PHONY: help build-lib vet fmt check clean deps
