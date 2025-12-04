#!/bin/bash
set -e

# CAMI Installation Script
# This script installs CAMI and creates the user workspace

VERSION="0.4.0"
INSTALL_DIR="${CAMI_DIR:-$HOME/cami-workspace}"
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
print_info "Workspace directory: $INSTALL_DIR"
print_info "Binary directory: $BIN_DIR (requires sudo)"
echo ""

# Ask user for custom workspace location
read -p "Customize workspace directory? (press Enter for default): " CUSTOM_DIR
if [ -n "$CUSTOM_DIR" ]; then
    # Expand ~ if present
    CUSTOM_DIR="${CUSTOM_DIR/#\~/$HOME}"
    INSTALL_DIR="$CUSTOM_DIR"
    print_info "Using custom workspace directory: $INSTALL_DIR"
fi

echo ""

# Create CAMI workspace
if [ -d "$INSTALL_DIR" ]; then
    print_info "CAMI workspace already exists at $INSTALL_DIR"
    echo ""
    echo "This will update:"
    echo "  ✓ Bundled agents (.claude/agents/agent-architect.md)"
    echo "  ✓ Template files (CLAUDE.md, README.md, .mcp.json, .gitignore)"
    echo "  ✓ Claude settings (if using default CAMI settings)"
    echo ""
    echo "This will NOT touch:"
    echo "  ✗ config.yaml"
    echo "  ✗ sources/my-agents/"
    echo "  ✗ Your agent sources"
    echo "  ✗ Custom Claude settings (if modified)"
    echo ""
    print_info "⚠️  If you've customized CLAUDE.md, back it up first!"
    echo ""
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

# Backup function for upgrades
backup_if_exists() {
    local file="$1"
    if [ -f "$file" ]; then
        local backup_dir="$INSTALL_DIR/.cami-backup-$(date +%Y%m%d-%H%M%S)"
        mkdir -p "$backup_dir"
        cp "$file" "$backup_dir/"
        echo "$backup_dir"
    fi
}

# Check if this is an upgrade
IS_UPGRADE=false
if [ -f "$INSTALL_DIR/CLAUDE.md" ]; then
    IS_UPGRADE=true
fi

# Copy template files (always update these)
print_info "Installing template files..."

if [ "$IS_UPGRADE" = true ]; then
    # Create a single backup directory for this upgrade
    BACKUP_DIR="$INSTALL_DIR/.cami-backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$BACKUP_DIR"

    # Backup user-modifiable files
    [ -f "$INSTALL_DIR/CLAUDE.md" ] && cp "$INSTALL_DIR/CLAUDE.md" "$BACKUP_DIR/"
    [ -f "$INSTALL_DIR/.claude/agents/agent-architect.md" ] && cp "$INSTALL_DIR/.claude/agents/agent-architect.md" "$BACKUP_DIR/"
    [ -f "$INSTALL_DIR/.claude/settings.json" ] && cp "$INSTALL_DIR/.claude/settings.json" "$BACKUP_DIR/"

    print_info "Backed up existing files to $BACKUP_DIR"
fi

cp "$TEMPLATE_DIR/CLAUDE.md" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/README.md" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/.gitignore" "$INSTALL_DIR/"
cp "$TEMPLATE_DIR/.mcp.json" "$INSTALL_DIR/"

# Deploy agent-architect (the only bundled agent)
print_info "Deploying agent-architect v4.0.0..."
cp "$TEMPLATE_DIR/agent-architect.md" "$INSTALL_DIR/.claude/agents/"

# Deploy settings.json with SessionStart hook for reconciliation
if [ ! -f "$INSTALL_DIR/.claude/settings.json" ]; then
    cp "$TEMPLATE_DIR/.claude/settings.json" "$INSTALL_DIR/.claude/"
    print_success "Installed .claude/settings.json with SessionStart hook"
else
    # Only update if it's our default template (check for reconcile hook)
    if grep -q "cami source reconcile" "$INSTALL_DIR/.claude/settings.json" 2>/dev/null; then
        if [ "$IS_UPGRADE" = true ]; then
            cp "$INSTALL_DIR/.claude/settings.json" "$BACKUP_DIR/" 2>/dev/null || true
        fi
        cp "$TEMPLATE_DIR/.claude/settings.json" "$INSTALL_DIR/.claude/"
        print_success "Updated .claude/settings.json"
    else
        print_info "Skipping .claude/settings.json (custom configuration detected)"
    fi
fi

# Copy templates to my-agents source (only if they don't exist)
if [ ! -f "$INSTALL_DIR/sources/my-agents/.camiignore" ]; then
    cp "$TEMPLATE_DIR/.camiignore" "$INSTALL_DIR/sources/my-agents/"
fi

# STRATEGIES.yaml - update if it's still the default template
if [ ! -f "$INSTALL_DIR/sources/my-agents/STRATEGIES.yaml" ]; then
    cp "$TEMPLATE_DIR/sources/my-agents/STRATEGIES.yaml" "$INSTALL_DIR/sources/my-agents/"
    print_success "Installed STRATEGIES.yaml template"
elif grep -q "^# CAMI Agent Strategies" "$INSTALL_DIR/sources/my-agents/STRATEGIES.yaml" 2>/dev/null; then
    # It's a CAMI template file, safe to update
    if [ "$IS_UPGRADE" = true ]; then
        cp "$INSTALL_DIR/sources/my-agents/STRATEGIES.yaml" "$BACKUP_DIR/" 2>/dev/null || true
    fi
    cp "$TEMPLATE_DIR/sources/my-agents/STRATEGIES.yaml" "$INSTALL_DIR/sources/my-agents/"
    print_success "Updated STRATEGIES.yaml template"
fi

# Create initial config if it doesn't exist
if [ ! -f "$INSTALL_DIR/config.yaml" ]; then
    print_info "Creating initial config.yaml..."

    # Prompt for default projects directory
    echo ""
    echo "Where do you want CAMI to create new projects by default?"
    echo "Examples: ~/projects, ~/dev, ~/workspace, ~/code"
    read -p "Default projects directory (press Enter for ~/projects): " PROJECTS_DIR

    if [ -z "$PROJECTS_DIR" ]; then
        PROJECTS_DIR="$HOME/projects"
    else
        # Expand ~ if present
        PROJECTS_DIR="${PROJECTS_DIR/#\~/$HOME}"
    fi

    # Create directory if it doesn't exist
    mkdir -p "$PROJECTS_DIR"
    print_success "Default projects directory: $PROJECTS_DIR"

    # Get current timestamp in RFC3339 format
    INSTALL_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    cat > "$INSTALL_DIR/config.yaml" <<EOF
version: "1"
install_timestamp: $INSTALL_TIMESTAMP
setup_complete: false
agent_sources:
  - name: my-agents
    type: local
    path: $INSTALL_DIR/sources/my-agents
    priority: 10
    git:
      enabled: false

deploy_locations: []
default_projects_dir: $PROJECTS_DIR
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

# If using custom location, inform user about CAMI_DIR
if [ "$INSTALL_DIR" != "$HOME/cami-workspace" ]; then
    echo "⚠️  Custom workspace location detected!"
    echo ""
    echo "  Add this to your ~/.zshrc or ~/.bashrc:"
    echo "     export CAMI_DIR=\"$INSTALL_DIR\""
    echo ""
    echo "  Then restart your terminal or run: source ~/.zshrc"
    echo ""
fi

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
