#!/bin/bash
set -e

# CAMI Installation Script
# This script installs CAMI and creates the user workspace

VERSION="0.3.0"
INSTALL_DIR="${CAMI_DIR:-$HOME/cami}"
BIN_DIR="${BIN_DIR:-/usr/local/bin}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR/templates"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print functions
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }
print_info() { echo -e "${YELLOW}ℹ${NC} $1"; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "unknown" ;;
    esac
}

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  CAMI Installation v$VERSION"
echo "  Claude Agent Management Interface"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Detect platform
OS=$(detect_os)
ARCH=$(detect_arch)

if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
    print_error "Unsupported platform: OS=$OS, ARCH=$ARCH"
    echo "Supported platforms: macOS (amd64/arm64), Linux (amd64/arm64), Windows (amd64/arm64)"
    exit 1
fi

print_info "Detected platform: $OS/$ARCH"

# Check if binary exists in script directory
BINARY_NAME="cami"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="cami.exe"
fi

# Look for binary (either in current dir during build or in releases)
BINARY_PATH=""
if [ -f "$SCRIPT_DIR/../cami" ]; then
    BINARY_PATH="$(cd "$SCRIPT_DIR/.." && pwd)/cami"
elif [ -f "$SCRIPT_DIR/cami" ]; then
    BINARY_PATH="$(cd "$SCRIPT_DIR" && pwd)/cami"
else
    print_error "CAMI binary not found. Please run 'make build' first or download a release."
    exit 1
fi

echo ""
print_info "Installation directory: $INSTALL_DIR"
print_info "Binary directory: $BIN_DIR"
echo ""

# Create CAMI workspace
if [ -d "$INSTALL_DIR" ]; then
    print_info "CAMI workspace already exists at $INSTALL_DIR"
    read -p "Continue and overwrite template files? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Installation cancelled"
        exit 1
    fi
else
    print_info "Creating CAMI workspace at $INSTALL_DIR"
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/sources/my-agents"
mkdir -p "$INSTALL_DIR/.claude/agents"

# Copy template files
print_info "Installing template files..."
cp "$TEMPLATE_DIR/CLAUDE.md" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/README.md" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/.gitignore" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/.mcp.json" "$INSTALL_DIR/"

# Deploy agent-architect (the only bundled agent)
print_info "Deploying agent-architect..."
cp "$TEMPLATE_DIR/agent-architect.md" "$INSTALL_DIR/.claude/agents/"

# Create initial config if it doesn't exist
if [ ! -f "$INSTALL_DIR/config.yaml" ]; then
    print_info "Creating initial config.yaml..."
    cat > "$INSTALL_DIR/config.yaml" <<EOF
version: "1"
agent_sources:
  - name: my-agents
    type: local
    path: $INSTALL_DIR/sources/my-agents
    priority: 10
    git:
      enabled: false

deploy_locations: []
EOF
fi

print_success "CAMI workspace created at $INSTALL_DIR"

# Install binary
print_info "Installing CAMI binary to $BIN_DIR..."

# Create bin directory if it doesn't exist
if [ ! -d "$BIN_DIR" ]; then
    print_info "Creating $BIN_DIR directory..."
    sudo mkdir -p "$BIN_DIR"
fi

if [ ! -w "$BIN_DIR" ]; then
    print_info "Need sudo permissions to install to $BIN_DIR"
    sudo cp "$BINARY_PATH" "$BIN_DIR/cami"
    sudo chmod +x "$BIN_DIR/cami"
else
    cp "$BINARY_PATH" "$BIN_DIR/cami"
    chmod +x "$BIN_DIR/cami"
fi

print_success "CAMI binary installed to $BIN_DIR/cami"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print_success "CAMI installation complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📂 CAMI workspace: $INSTALL_DIR"
echo "🔧 Binary location: $BIN_DIR/cami"
echo ""
echo "Next steps:"
echo ""
echo "  1. Open your CAMI workspace:"
echo "     $ cd $INSTALL_DIR"
echo "     $ claude"
echo ""
echo "  2. Ask Claude to help you get started:"
echo "     \"Help me get started with CAMI\""
echo ""
echo "  3. Optional - Set up global MCP access:"
echo "     Add to ~/.claude/settings.json:"
echo ""
echo "     {"
echo "       \"mcpServers\": {"
echo "         \"cami\": {"
echo "           \"command\": \"cami\","
echo "           \"args\": [\"--mcp\"]"
echo "         }"
echo "       }"
echo "     }"
echo ""
echo "  4. CLI commands work from anywhere:"
echo "     $ cami list"
echo "     $ cami source add <git-url>"
echo "     $ cami deploy <agents> <project-path>"
echo ""
echo "Documentation: https://github.com/lando-labs/cami"
echo ""
