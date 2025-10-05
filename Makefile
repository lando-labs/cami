.PHONY: build install clean test run

# Binary name
BINARY=cami
VERSION=0.1.0

# Build the application
build:
	@echo "Building CAMI v$(VERSION)..."
	@go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY) cmd/cami/main.go
	@echo "Build complete: ./$(BINARY)"

# Install to $GOPATH/bin
install:
	@echo "Installing CAMI v$(VERSION)..."
	@go install -ldflags="-X main.version=$(VERSION)" ./cmd/cami
	@echo "Installed to $(shell go env GOPATH)/bin/cami"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY)
	@go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run: build
	@./$(BINARY)

# Show help
help:
	@echo "CAMI - Claude Agent Management Interface"
	@echo ""
	@echo "Available targets:"
	@echo "  make build   - Build the binary"
	@echo "  make install - Install to GOPATH/bin"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make test    - Run tests"
	@echo "  make run     - Build and run"
	@echo "  make help    - Show this help"
