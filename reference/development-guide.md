<!--
AI-Generated Documentation
Created by: mcp-specialist
Date: 2025-10-05
Purpose: Development and contribution guide for CAMI MCP server
-->

# CAMI MCP Server Development Guide

Guide for developers working on or extending the CAMI MCP server.

## Table of Contents

1. [Development Setup](#development-setup)
2. [Building](#building)
3. [Testing](#testing)
4. [Adding New Tools](#adding-new-tools)
5. [Debugging](#debugging)
6. [Contributing](#contributing)
7. [Release Process](#release-process)

---

## Development Setup

### Prerequisites

- **Go**: 1.24.4 or later
- **Git**: For version control
- **Claude Desktop** (optional): For integration testing

### Clone and Setup

```bash
cd /Users/lando/Development/cami
git pull  # Ensure latest code

# Install dependencies
go mod download

# Verify setup
go version
go mod verify
```

### Dependencies

The MCP server uses:

```go
require (
    github.com/lando/cami/internal/agent
    github.com/lando/cami/internal/deploy
    github.com/lando/cami/internal/docs
    github.com/modelcontextprotocol/go-sdk v1.0.0
)
```

---

## Building

### Build MCP Server

```bash
# Build from project root
go build -o cami-mcp cmd/cami-mcp/main.go

# Verify binary
./cami-mcp --help  # Will start server (requires stdio input)
ls -lh cami-mcp
```

### Build Options

```bash
# Debug build (with symbols)
go build -o cami-mcp cmd/cami-mcp/main.go

# Optimized release build
go build -ldflags="-s -w" -o cami-mcp cmd/cami-mcp/main.go

# Cross-compilation (example: Linux)
GOOS=linux GOARCH=amd64 go build -o cami-mcp-linux cmd/cami-mcp/main.go
```

### Build Output

- **Binary Size**: ~7-8 MB (arm64 macOS)
- **Binary Location**: `./cami-mcp`
- **Architecture**: Matches build system (arm64, amd64, etc.)

---

## Testing

### Unit Testing

Currently, the MCP server doesn't have unit tests. To add them:

```bash
# Create test file
cat > cmd/cami-mcp/main_test.go << 'EOF'
package main

import "testing"

func TestGetVCAgentsDir(t *testing.T) {
    // Test implementation
}
EOF

# Run tests
go test ./cmd/cami-mcp/...
```

### Integration Testing

Test with MCP Inspector (if available):

```bash
# Run server manually
./cami-mcp

# Server expects JSON-RPC on stdin, outputs on stdout
# Logs appear on stderr
```

### Manual Testing with Claude Desktop

1. Update Claude Desktop config:

```json
{
  "mcpServers": {
    "cami": {
      "command": "/Users/lando/Development/cami/cami-mcp",
      "env": {
        "CAMI_VC_AGENTS_DIR": "/Users/lando/Development/cami/vc-agents"
      }
    }
  }
}
```

2. Restart Claude Desktop
3. Open a conversation
4. Type commands that should trigger CAMI tools
5. Check logs: `tail -f ~/Library/Logs/Claude/mcp*.log`

### Test Scenarios

#### Test 1: List Agents
```
User: "What agents are available in CAMI?"
Expected: Claude calls list_agents tool
```

#### Test 2: Deploy Agent
```
User: "Deploy the architect agent to /Users/username/test-project"
Expected: Claude calls deploy_agents with correct parameters
```

#### Test 3: Scan Project
```
User: "What agents are deployed in /Users/username/my-project?"
Expected: Claude calls scan_deployed_agents
```

---

## Adding New Tools

### Step 1: Define Argument Struct

```go
// Add to cmd/cami-mcp/main.go
type MyNewToolArgs struct {
    Param1 string `json:"param1" jsonschema:"required,description=First parameter"`
    Param2 int    `json:"param2" jsonschema:"description=Second parameter (optional)"`
}
```

**JSON Schema Tags**:
- `required`: Mark parameter as required
- `description`: Description for Claude (appears in tool schema)
- Use standard JSON tags for field names

### Step 2: Implement Handler

```go
mcp.AddTool(server, &mcp.Tool{
    Name: "my_new_tool",
    Description: "Clear description of what this tool does. " +
                 "Explain when to use it.",
}, func(ctx context.Context, req *mcp.CallToolRequest, args MyNewToolArgs) (*mcp.CallToolResult, any, error) {
    // 1. Validate inputs
    if args.Param1 == "" {
        return nil, nil, fmt.Errorf("param1 cannot be empty")
    }

    // 2. Perform operation using internal packages
    result, err := someInternalPackage.DoSomething(args.Param1, args.Param2)
    if err != nil {
        return nil, nil, fmt.Errorf("operation failed: %w", err)
    }

    // 3. Format response
    responseText := fmt.Sprintf("Operation completed: %v", result)

    // 4. Return result
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: responseText},
        },
    }, result, nil  // Optional structured data
})
```

### Step 3: Update Documentation

1. **Tool Catalog**: Add entry in `reference/tool-catalog.md`
2. **README**: Update README.md if user-facing
3. **CHANGELOG**: Document new feature

### Best Practices

1. **Naming**: Use snake_case (e.g., `deploy_agents`, not `deployAgents`)
2. **Descriptions**: Be specific about when to use the tool
3. **Error Messages**: Clear, actionable error messages
4. **Validation**: Validate all inputs before operations
5. **Logging**: Log important operations to stderr
6. **Response Format**: Provide both text (for Claude) and structured data (optional)

---

## Debugging

### Enable Debug Logging

The server logs to stderr by default. Increase verbosity:

```go
// In main.go
log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
```

### View MCP Communication

Since the server uses stdio transport, you can't easily see raw JSON-RPC messages. Options:

1. **Claude Desktop Logs**: Check `~/Library/Logs/Claude/mcp*.log`
2. **Wrapper Script**: Create a logging wrapper

```bash
#!/bin/bash
# cami-mcp-debug.sh
tee /tmp/cami-mcp-input.log | /Users/lando/Development/cami/cami-mcp 2>&1 | tee /tmp/cami-mcp-output.log
```

Update Claude config to use wrapper:
```json
{
  "command": "/Users/lando/Development/cami/cami-mcp-debug.sh"
}
```

### Common Issues

#### Issue: "vc-agents directory not found"
**Solution**: Set `CAMI_VC_AGENTS_DIR` environment variable in Claude config

#### Issue: "Server doesn't appear in Claude"
**Solution**:
1. Check Claude Desktop config syntax (valid JSON)
2. Verify binary path is absolute
3. Restart Claude Desktop
4. Check logs

#### Issue: "Tool calls fail silently"
**Solution**:
1. Check stderr logs (server-side)
2. Check Claude Desktop logs (client-side)
3. Verify parameter types match schema

### Debugging Tools

```bash
# Check binary info
file ./cami-mcp
otool -L ./cami-mcp  # macOS: Check dependencies

# Test binary runs
./cami-mcp  # Should wait for stdin (Ctrl+C to exit)

# Monitor logs
tail -f ~/Library/Logs/Claude/mcp*.log

# Check environment
env | grep CAMI
```

---

## Contributing

### Code Style

- **Format**: Use `gofmt` (enforced by `go fmt`)
- **Imports**: Group stdlib, external, internal
- **Comments**: Document exported functions and types
- **Error Handling**: Always check errors, wrap with context

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if available)
golangci-lint run
```

### Commit Messages

Follow conventional commits:

```
feat(mcp): add new tool for agent updates
fix(mcp): correct path validation in deploy_agents
docs(mcp): update tool catalog with examples
refactor(mcp): simplify agent scanning logic
```

### Pull Request Process

1. **Fork** repository (if external contributor)
2. **Branch** from main: `git checkout -b feature/my-new-tool`
3. **Implement** changes
4. **Test** thoroughly (manual + automated if available)
5. **Document** in reference docs
6. **Commit** with clear messages
7. **Push** and create PR
8. **Review** address feedback

---

## Release Process

### Version Bumping

Update version in `cmd/cami-mcp/main.go`:

```go
const (
    serverName    = "cami"
    serverVersion = "0.2.0"  // Increment here
)
```

### Build Release Binary

```bash
# Clean build
rm -f cami-mcp

# Build optimized binary
go build -ldflags="-s -w" -o cami-mcp cmd/cami-mcp/main.go

# Verify
./cami-mcp  # Test basic functionality
```

### Distribution

The MCP server binary is distributed as part of the CAMI project:

1. **Local Development**: Users build from source
2. **Release**: Include pre-built binary in releases (future)

### Versioning Strategy

Follow semantic versioning:
- **Major**: Breaking changes to tool interfaces
- **Minor**: New tools, backward-compatible features
- **Patch**: Bug fixes, documentation updates

---

## Architecture Notes

### Package Organization

```
cmd/cami-mcp/
└── main.go           # MCP server entry point, tool handlers

internal/
├── agent/            # Agent file parsing
├── deploy/           # Deployment logic
├── docs/             # CLAUDE.md management
└── discovery/        # Agent scanning (not used yet)
```

### Design Principles

1. **Thin Protocol Layer**: MCP server is a thin wrapper over internal packages
2. **Direct Integration**: Import and use internal packages directly (no CLI wrapping)
3. **Typed Handlers**: Leverage Go SDK's automatic schema generation
4. **Error Transparency**: Pass through internal package errors with context
5. **Logging to stderr**: Keep stdout clean for MCP protocol

### Future Enhancements

- [ ] Resource providers (expose agents as MCP resources)
- [ ] Prompt templates (deployment workflow prompts)
- [ ] Progress notifications (for long operations)
- [ ] Caching (agent metadata caching)
- [ ] Batch operations (multi-project deployment)
- [ ] Update command (update deployed agents)
- [ ] Rollback support (revert to previous versions)

---

## Resources

- **MCP Specification**: https://modelcontextprotocol.io/
- **Go SDK Docs**: https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk
- **CAMI Architecture**: See `reference/mcp-architecture.md`
- **Tool Reference**: See `reference/tool-catalog.md`

---

## Support

For issues or questions:
1. Check existing documentation
2. Review logs (stderr + Claude Desktop logs)
3. Open GitHub issue with:
   - Go version
   - CAMI version
   - Steps to reproduce
   - Error logs
