# Docker Compose Manager (DCM) Makefile

# Variables
BINARY_NAME=dcm
GO=go
GOFMT=gofmt
GOLINT=golint
GOVET=$(GO) vet
BUILD_DIR=.
CMD_DIR=./cmd/dcm
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: fmt lint vet build

# Build the application
.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Install the application
.PHONY: install
install:
	$(GO) install $(LDFLAGS) $(CMD_DIR)

# Run the application
.PHONY: run
run:
	$(GO) run $(CMD_DIR)

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	$(GO) clean

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

# Lint code
.PHONY: lint
lint:
	$(GOLINT) ./...

# Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Build for multiple platforms
.PHONY: release
release: clean
	# Linux
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	# macOS
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	# Windows
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all       - Format, lint, vet, and build"
	@echo "  build     - Build the binary"
	@echo "  install   - Install the binary"
	@echo "  run       - Run the application"
	@echo "  clean     - Remove build artifacts"
	@echo "  fmt       - Format code"
	@echo "  lint      - Run linter"
	@echo "  vet       - Run go vet"
	@echo "  test      - Run tests"
	@echo "  release   - Build for multiple platforms"
	@echo "  help      - Show this help"