# Makefile for bumert

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate
GOBUILD=$(GOCMD) build

# Variables
PKG := ./...
TAGS_DEBUG := -tags debug
GENERATOR_CMD := cmd/gen-release/main.go
GENERATOR_BIN := gen-release

.PHONY: all test test-debug fmt tidy generate build-generator help

all: test ## Run all tests (release mode)

build-generator: ## Build the release stub generator
	@echo "Building generator..."
	$(GOBUILD) -o $(GENERATOR_BIN) $(GENERATOR_CMD)

generate: build-generator ## Generate release stubs from debug code by running generator binary
	@echo "Generating release stubs..."
	./$(GENERATOR_BIN) -in bumert_debug.go

test: generate ## Run tests without debug assertions (generates stubs first)
	@echo "Running tests (release mode)..."
	$(GOTEST) $(PKG)

test-debug: generate ## Run tests with debug assertions enabled (generates stubs first)
	@echo "Running tests (debug mode)..."
	$(GOTEST) $(TAGS_DEBUG) $(PKG)

fmt: ## Format Go source code
	@echo "Formatting code..."
	$(GOFMT) $(PKG)

tidy: ## Tidy Go module dependencies
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'