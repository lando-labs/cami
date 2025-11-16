.PHONY: help build install clean test lint release-all release-darwin release-linux release-windows

# Variables
BINARY_NAME=cami
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.3.0")
BUILD_DIR=build
INSTALL_DIR=$(HOME)/cami
BIN_DIR=/usr/local/bin

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

help: ## Show this help message
	@echo "CAMI Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build CAMI binary for current platform
	@echo "Building CAMI v$(VERSION) for current platform..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/cami
	@echo "✓ Built: $(BINARY_NAME)"

install: build ## Build and install CAMI locally
	@echo "Installing CAMI..."
	./install/install.sh
	@echo "✓ Installation complete!"

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

lint: ## Run linters
	@echo "Running linters..."
	gofmt -s -l .
	go vet ./...

# Cross-platform releases
release-all: release-darwin release-linux release-windows ## Build releases for all platforms
	@echo "✓ All releases built in $(BUILD_DIR)/"

release-darwin: ## Build macOS releases (amd64 and arm64)
	@echo "Building macOS releases..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/cami
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/cami
	@echo "✓ macOS releases built"

release-linux: ## Build Linux releases (amd64 and arm64)
	@echo "Building Linux releases..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cami
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/cami
	@echo "✓ Linux releases built"

release-windows: ## Build Windows releases (amd64 and arm64)
	@echo "Building Windows releases..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cami
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd/cami
	@echo "✓ Windows releases built"

# Package releases with install scripts
package: release-all ## Package releases with installation scripts
	@echo "Packaging releases..."
	@for platform in darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 windows-amd64 windows-arm64; do \
		echo "Packaging $$platform..."; \
		mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION); \
		if echo $$platform | grep -q windows; then \
			cp $(BUILD_DIR)/$(BINARY_NAME)-$$platform.exe $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION)/cami.exe; \
		else \
			cp $(BUILD_DIR)/$(BINARY_NAME)-$$platform $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION)/cami; \
		fi; \
		cp -r install $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION)/; \
		cp LICENSE $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION)/; \
		cp README.md $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION)/; \
		cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$$platform-$(VERSION).tar.gz $(BINARY_NAME)-$$platform-$(VERSION); \
		cd ..; \
		rm -rf $(BUILD_DIR)/$(BINARY_NAME)-$$platform-$(VERSION); \
	done
	@echo "✓ Packages created in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/*.tar.gz

# Development helpers
dev: ## Run in development mode (go run)
	@go run ./cmd/cami $(ARGS)

dev-mcp: ## Run MCP server in development mode
	@go run ./cmd/cami --mcp

# Version info
version: ## Show version information
	@echo "CAMI v$(VERSION)"
