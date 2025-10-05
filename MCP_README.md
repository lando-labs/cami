# CAMI MCP Server

Model Context Protocol (MCP) server for CAMI (Claude Agent Management Interface), enabling seamless integration with Claude Code and other MCP-compatible clients.

## Overview

The CAMI MCP server exposes CAMI's agent management capabilities through a standardized protocol interface. This allows Claude Code to:

- 📋 **List** available agents from the version-controlled repository
- 🚀 **Deploy** agents to project directories
- 📝 **Document** deployed agents in CLAUDE.md files
- 🔍 **Scan** projects to audit deployed agents and versions
- 📍 **Manage** deployment locations for quick access to projects

## Quick Start

### 1. Build the Server

```bash
cd /Users/lando/Development/cami
go build -o cami-mcp cmd/cami-mcp/main.go
```

### 2. Configure Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

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

### 3. Restart Claude Desktop

Restart Claude Desktop to load the MCP server.

### 4. Use in Conversations

```
You: "What agents are available in CAMI?"
Claude: [calls list_agents tool]

You: "Deploy the architect agent to /Users/username/my-project"
Claude: [calls deploy_agents tool]
```

## Available Tools

The CAMI MCP server provides **7 tools** for comprehensive agent management:

### Agent Deployment Tools

### 🚀 deploy_agents

Deploy selected agents to a target project's `.claude/agents/` directory.

**Parameters**:
- `agent_names` (string[]): Names of agents to deploy
- `target_path` (string): Absolute path to target project
- `overwrite` (boolean, optional): Overwrite existing files

**Example**:
```json
{
  "agent_names": ["architect", "backend", "frontend"],
  "target_path": "/Users/username/projects/my-app",
  "overwrite": false
}
```

### 📝 update_claude_md

Update a project's `CLAUDE.md` file with deployed agent documentation.

**Parameters**:
- `target_path` (string): Absolute path to target project

**Example**:
```json
{
  "target_path": "/Users/username/projects/my-app"
}
```

### 📋 list_agents

List all available agents from CAMI's version-controlled repository.

**Parameters**: None

**Example**:
```json
{}
```

### 🔍 scan_deployed_agents

Scan a project to find deployed agents and compare versions.

**Parameters**:
- `target_path` (string): Absolute path to target project

**Example**:
```json
{
  "target_path": "/Users/username/projects/my-app"
}
```

Returns status for each agent:
- ✓ **up-to-date**: Deployed version matches available version
- ⚠ **update-available**: Newer version available
- ○ **not-deployed**: Agent not present in project

### Location Management Tools

### 📍 add_location

Add a new deployment location to CAMI's configuration.

**Parameters**:
- `name` (string): Friendly name for the location
- `path` (string): Absolute path to project directory

**Example**:
```json
{
  "name": "my-project",
  "path": "/Users/username/projects/my-app"
}
```

### 📋 list_locations

List all configured deployment locations.

**Parameters**: None

**Example**:
```json
{}
```

Returns a list of locations with names and paths.

### 🗑️ remove_location

Remove a deployment location from CAMI's configuration.

**Parameters**:
- `name` (string): Name of location to remove

**Example**:
```json
{
  "name": "my-project"
}
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CAMI_VC_AGENTS_DIR` | Location of version-controlled agents | Auto-detected |

### Configuration File

CAMI stores deployment locations in `~/.cami.json`:

```json
{
  "deploy_locations": [
    {
      "name": "my-project",
      "path": "/Users/username/projects/my-app"
    },
    {
      "name": "another-project",
      "path": "/Users/username/projects/another-app"
    }
  ]
}
```

This file is automatically created and managed by the location management tools.

### Auto-Detection

If `CAMI_VC_AGENTS_DIR` is not set, the server attempts to find `vc-agents` in:
1. Current working directory: `./vc-agents`
2. Executable directory: `{exec_dir}/vc-agents`

## Architecture

```
┌─────────────────┐
│  Claude Code    │
└────────┬────────┘
         │ MCP Protocol (stdio)
         │
┌────────▼────────┐
│  CAMI MCP       │
│  Server         │
│  ┌──────────┐   │
│  │  Tools   │   │
│  └────┬─────┘   │
│       │         │
│  ┌────▼─────┐   │
│  │ Internal │   │
│  │ Packages │   │
│  └──────────┘   │
└────────┬────────┘
         │
┌────────▼────────┐
│  File System    │
│  • vc-agents/   │
│  • .claude/     │
│  • CLAUDE.md    │
└─────────────────┘
```

**Key Components**:
- **MCP Protocol Layer**: Stdio transport with JSON-RPC
- **Tool Handlers**: Type-safe handlers with automatic schema generation
- **Internal Packages**: Direct integration with CAMI's agent, deploy, and docs packages
- **File System**: Read from vc-agents, write to .claude/agents/

## Development

### Prerequisites

- Go 1.24.4 or later
- CAMI project with vc-agents directory

### Building

```bash
# Standard build
go build -o cami-mcp cmd/cami-mcp/main.go

# Optimized release build
go build -ldflags="-s -w" -o cami-mcp cmd/cami-mcp/main.go
```

### Testing

```bash
# Build
go build -o cami-mcp cmd/cami-mcp/main.go

# Test with Claude Desktop (see Quick Start)

# View logs
tail -f ~/Library/Logs/Claude/mcp*.log
```

### Adding New Tools

See `reference/development-guide.md` for detailed instructions on adding new tools.

## Documentation

Comprehensive reference documentation:

- **[MCP Architecture](reference/mcp-architecture.md)**: Server architecture and design
- **[Tool Catalog](reference/tool-catalog.md)**: Complete tool reference with examples
- **[Development Guide](reference/development-guide.md)**: Contributing and extending the server

## Troubleshooting

### Server doesn't appear in Claude

1. Verify Claude Desktop config is valid JSON
2. Check binary path is absolute
3. Restart Claude Desktop
4. Check Claude logs: `~/Library/Logs/Claude/mcp*.log`

### "vc-agents directory not found"

Set `CAMI_VC_AGENTS_DIR` environment variable in Claude config:

```json
{
  "env": {
    "CAMI_VC_AGENTS_DIR": "/absolute/path/to/vc-agents"
  }
}
```

### Tool calls fail

1. Check stderr output from server
2. Verify parameters match tool schema
3. Ensure target paths are absolute and exist
4. Check file permissions

### Viewing Logs

```bash
# Server logs (stderr)
tail -f ~/Library/Logs/Claude/mcp-cami.log

# Claude Desktop logs
tail -f ~/Library/Logs/Claude/app.log
```

## Technical Details

- **Language**: Go 1.24.4
- **MCP SDK**: github.com/modelcontextprotocol/go-sdk v1.0.0
- **Transport**: stdio (standard input/output)
- **Protocol**: JSON-RPC 2.0
- **Binary Size**: ~7-8 MB (arm64 macOS)

## Security

- All target paths validated before operations
- Existing files not overwritten without explicit `overwrite: true`
- VC agents directory accessed read-only
- No network access (fully local)
- Error messages sanitized (no internal path leakage)

## Performance

- **Agent Loading**: O(n) where n = number of agent files
- **Deployment**: O(m) where m = number of agents to deploy
- **Scanning**: O(n) for version comparison
- **Memory**: <1 MB for typical agent metadata

## Future Enhancements

- [x] Location management (add/list/remove deployment locations)
- [ ] Resource providers (expose agents as MCP resources)
- [ ] Prompt templates (deployment workflow prompts)
- [ ] Progress notifications (long-running operations)
- [ ] Agent metadata caching
- [ ] Multi-project batch deployment
- [ ] Agent update command
- [ ] Version rollback support
- [ ] Deploy to location by name (use saved locations)

## License

Part of the CAMI project.

## Contributing

See `reference/development-guide.md` for development setup and contribution guidelines.

## Resources

- **MCP Specification**: https://modelcontextprotocol.io/
- **MCP Go SDK**: https://github.com/modelcontextprotocol/go-sdk
- **CAMI Project**: `/Users/lando/Development/cami`

## Support

For issues:
1. Check documentation in `reference/`
2. Review server logs (stderr)
3. Check Claude Desktop logs
4. Open GitHub issue with reproduction steps

---

**Version**: 0.2.0
**Status**: Production Ready
**Last Updated**: 2025-10-05
