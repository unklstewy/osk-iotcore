# Makefile for osk-iotcore project
.PHONY: build run test lint fmt deps clean help install-tools

# Variables
BINARY_NAME := oskway
BINARY_PATH := ./cmd/oskway
BUILD_DIR := ./build
COVERAGE_DIR := ./coverage

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOIMPORTS := goimports
GOLANGCI_LINT := golangci-lint

# Build flags
LDFLAGS := -ldflags "-s -w"
BUILD_FLAGS := -trimpath $(LDFLAGS)

# Default target
all: deps fmt lint test build

# Build the application
build: protocols
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Generate Wayland protocol files
protocols:
	@echo "Generating Wayland protocol files..."
	@mkdir -p generated
	wayland-scanner client-header /usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml generated/xdg-shell-client-protocol.h
	wayland-scanner private-code /usr/share/wayland-protocols/stable/xdg-shell/xdg-shell.xml generated/xdg-shell-protocol.c
	gcc -c -fPIC generated/xdg-shell-protocol.c -o generated/xdg-shell-protocol.o `pkg-config --cflags wayland-client`
	ar rcs generated/libxdg-shell-protocol.a generated/xdg-shell-protocol.o

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests with race detection
test:
	@echo "Running tests with race detection..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Run linter
lint: install-tools
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run --config .golangci.yml

# Format code
fmt: install-tools
	@echo "Formatting code with goimports..."
	$(GOIMPORTS) -w -local github.com/iotcore/osk-iotcore .
	@echo "Formatting code with gofmt..."
	$(GOCMD) fmt ./...

# Install/update dependencies
deps:
	@echo "Installing/updating dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download
	$(GOMOD) verify

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@command -v $(GOIMPORTS) >/dev/null 2>&1 || { \
		echo "Installing goimports..."; \
		$(GOGET) golang.org/x/tools/cmd/goimports@latest; \
	}
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest; \
	}

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)

# Development build (with debug symbols)
dev-build:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -race -o $(BUILD_DIR)/$(BINARY_NAME)-dev $(BINARY_PATH)

# Run in development mode
dev-run: dev-build
	@echo "Running $(BINARY_NAME) in development mode..."
	$(BUILD_DIR)/$(BINARY_NAME)-dev

# Show coverage in terminal
coverage: test
	@echo "Test coverage summary:"
	$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Check for security issues
security: install-tools
	@echo "Running security checks..."
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	}
	gosec ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  test         - Run tests with race detection and coverage"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code with goimports and gofmt"
	@echo "  deps         - Install/update dependencies"
	@echo "  install-tools- Install development tools"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-build    - Build with debug symbols"
	@echo "  dev-run      - Run in development mode"
	@echo "  coverage     - Show test coverage summary"
	@echo "  bench        - Run benchmark tests"
	@echo "  security     - Run security checks"
	@echo "  help         - Show this help message"
