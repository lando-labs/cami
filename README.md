# CAMI - Claude Agent Management Interface

A comprehensive toolkit for managing and deploying Claude Code agents across projects. CAMI provides a rich terminal UI, powerful CLI commands, and native MCP server integration for seamless Claude Code workflows.

## Features

- **23 Version-Controlled Agents**: Curated collection of specialized agents for every development need
- **Interactive TUI**: Beautiful keyboard-driven interface built with Bubbletea
- **Full CLI Support**: Programmatic deployment and management commands
- **MCP Server Integration**: Native Claude Code/Desktop integration via Model Context Protocol
- **Agent Discovery**: Scan and update deployed agents across projects
- **Smart Documentation**: Automatic CLAUDE.md updates with deployed agent information
- **Multi-Select Deployment**: Deploy multiple agents at once
- **Location Management**: Configure and manage deployment locations
- **Conflict Detection**: Safe deployment with existing file detection
- **Version Tracking**: Compare deployed versions with available updates

## Installation

### CLI and TUI

Build from source:

```bash
go build -o cami cmd/cami/main.go
```

Or install to your PATH:

```bash
go install ./cmd/cami
```

### MCP Server (for Claude Code/Desktop)

Build the MCP server:

```bash
go build -o cami-mcp cmd/cami-mcp/main.go
```

Configure in Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "cami": {
      "command": "/absolute/path/to/cami-mcp",
      "env": {
        "CAMI_VC_AGENTS_DIR": "/absolute/path/to/cami/vc-agents"
      }
    }
  }
}
```

See [MCP_README.md](MCP_README.md) for detailed MCP server documentation.

## Usage

CAMI supports three modes of operation:

### 1. Interactive TUI Mode

Launch the interactive terminal UI (no arguments):

```bash
./cami
```

Perfect for browsing agents, managing locations, and interactive deployment.

### 2. CLI Commands

Use command-line interface for programmatic operations:

```bash
# Deploy agents
./cami deploy --agents frontend,backend --location ~/my-project

# Update documentation
./cami update-docs --location ~/my-project

# List available agents
./cami list

# Scan deployed agents
./cami scan --location ~/my-project
```

See [CLI-COMMANDS.md](CLI-COMMANDS.md) for complete CLI reference.

### 3. MCP Server (Claude Integration)

Once configured, Claude Code/Desktop can use CAMI tools directly:

```
You: "Deploy the architect and backend agents to my current project"
Claude: [uses CAMI MCP tools to deploy agents]

You: "What agents are available?"
Claude: [lists all 23 agents from CAMI]
```

## Available Agents

CAMI includes **23 specialized agents** covering the full development spectrum:

**Core Development:**
- `architect` - System architecture and design
- `frontend` - React, Vue, modern web frameworks
- `backend` - APIs, databases, server-side logic
- `mobile-native` - iOS/Android native development

**Specialized Domains:**
- `ai-ml-specialist` - AI/ML integration and deployment
- `blockchain-specialist` - Smart contracts and Web3
- `data-engineer` - Data pipelines and warehouses
- `embedded-systems` - IoT and firmware development
- `game-dev` - Game engines and mechanics

**Infrastructure & Operations:**
- `deploy` - Deployment infrastructure and CI/CD
- `devops` - Infrastructure as code and automation
- `gcp-firebase` - Google Cloud Platform and Firebase
- `performance-optimizer` - Performance analysis and optimization
- `security-specialist` - Security audits and compliance

**Integration & APIs:**
- `api-integrator` - Third-party API integration
- `mcp-specialist` - Model Context Protocol development

**Quality & Design:**
- `qa` - Testing and quality assurance
- `accessibility-expert` - WCAG compliance and inclusive design
- `designer` - Visual design and design systems
- `ux` - User experience and interaction design

**Documentation & Tools:**
- `docs-writer` - Technical documentation
- `terminal-specialist` - CLI tools and terminal UIs
- `agent-architect` - Claude agent design and optimization

Run `./cami list` for complete descriptions and version information.

## Keyboard Shortcuts (TUI Mode)

### Agent Selection View
- `↑/k` - Move up
- `↓/j` - Move down
- `space/x` - Select/deselect agent
- `enter/d` - Proceed to deployment
- `l` - Manage deployment locations
- `i` - Agent discovery and updates
- `q` - Quit

### Discovery View
- `↑/↓/j/k` - Navigate agents
- `←/→/h/l` - Switch between locations
- `u` - Update selected agent
- `U` - Update all agents at location
- `r` - Refresh scan
- `esc` - Back to agent selection

### Location Management View
- `↑/k` - Move up
- `↓/j` - Move down
- `a` - Add new location
- `d` - Delete selected location
- `esc` - Back to agent selection
- `q` - Quit

### Adding a Location
- `tab` - Switch between name and path fields
- Type to enter text
- `backspace` - Delete character
- `enter` - Save location
- `esc` - Cancel

### Deployment View
- `↑/k` - Move up
- `↓/j` - Move down
- `enter` - Deploy to selected location
- `esc` - Back to agent selection
- `q` - Quit

### Results View
- `enter` - Return to agent selection
- `q` - Quit

## Configuration

CAMI stores deployment locations in `~/.cami.json`. This file is created automatically when you add your first location.

Example configuration:

```json
{
  "deploy_locations": [
    {
      "name": "my-project",
      "path": "/Users/username/projects/my-project"
    }
  ]
}
```

## Agent Structure

Agents are stored in the `vc-agents/` directory with YAML frontmatter:

```markdown
---
name: architect
version: "1.0.0"
description: System architecture and design specialist
---

[Agent content here]
```

## Deployment Workflow

When you deploy agents, CAMI:

1. Creates `.claude/agents/` directory in the target location if needed
2. Copies selected agent files with YAML frontmatter and content
3. Detects conflicts with existing files (safe by default)
4. Shows deployment results with success/error/conflict status
5. Optionally updates `CLAUDE.md` with deployed agent documentation

### Example Workflow

```bash
# 1. Deploy agents to a project
./cami deploy --agents frontend,backend,qa --location ~/projects/my-app

# 2. Update project documentation
./cami update-docs --location ~/projects/my-app

# 3. Verify deployment
./cami scan --location ~/projects/my-app
```

Or use the TUI for interactive deployment:
1. Launch `./cami`
2. Select agents with `space`
3. Press `d` to deploy
4. Choose location
5. Press `i` to view and update deployed agents

## Quick Start

```bash
# Clone and build
cd /path/to/cami
go build -o cami cmd/cami/main.go

# Launch interactive mode
./cami

# Or use CLI commands
./cami deploy --agents architect,frontend --location ~/my-project
./cami update-docs --location ~/my-project

# For Claude Code integration, build MCP server
go build -o cami-mcp cmd/cami-mcp/main.go
# Then configure in Claude Desktop (see Installation section)
```

## Documentation

- **[README.md](README.md)** - This file, main entry point
- **[CLI-COMMANDS.md](CLI-COMMANDS.md)** - Complete CLI reference and examples
- **[MCP_README.md](MCP_README.md)** - MCP server setup and integration guide
- **[reference/](reference/)** - Detailed technical documentation
  - [mcp-architecture.md](reference/mcp-architecture.md) - MCP server architecture
  - [tool-catalog.md](reference/tool-catalog.md) - MCP tool reference
  - [development-guide.md](reference/development-guide.md) - Contributing guide
  - [cami-development-workflow.md](reference/cami-development-workflow.md) - Development workflow
  - [ai-generation-standard.md](reference/ai-generation-standard.md) - Agent generation standards

## Version

**CAMI v0.2.0** - Current Release

### What's New in v0.2.0
- Full CLI support with deploy, scan, list, and update-docs commands
- MCP server for native Claude Code/Desktop integration
- Agent discovery view in TUI (press `i`)
- Smart CLAUDE.md documentation updates
- JSON output format for programmatic use
- 23 specialized agents (expanded from initial set)

## Development

### Project Structure

```
cami/
├── cmd/
│   ├── cami/              # CLI/TUI entry point
│   └── cami-mcp/          # MCP server entry point
├── internal/
│   ├── agent/             # Agent parsing and metadata
│   ├── cli/               # CLI command implementations
│   ├── config/            # Configuration management
│   ├── deploy/            # Deployment engine
│   ├── discovery/         # Agent scanning and version comparison
│   ├── docs/              # CLAUDE.md update logic
│   ├── mcp/               # MCP server implementation
│   └── tui/               # Bubbletea TUI interface
├── vc-agents/             # 23 version-controlled agents
├── reference/             # Technical documentation
│   ├── mcp-architecture.md
│   ├── tool-catalog.md
│   ├── development-guide.md
│   └── ...
├── CLI-COMMANDS.md        # CLI reference
├── MCP_README.md          # MCP server guide
└── README.md              # This file
```

### Architecture

**CLI/TUI Application:**
- Entry: `cmd/cami/main.go`
- Mode detection: TUI (no args) vs CLI (with subcommands)
- Shared internal packages for core logic

**MCP Server:**
- Entry: `cmd/cami-mcp/main.go`
- Protocol: Model Context Protocol over stdio
- Tools: deploy_agents, list_agents, scan_deployed_agents, update_claude_md
- Integration: Direct usage of internal packages

**Shared Core:**
- `internal/agent` - YAML frontmatter parsing
- `internal/deploy` - File operations and conflict detection
- `internal/docs` - Smart CLAUDE.md merging
- `internal/discovery` - Version comparison logic

### Building

```bash
# Build CLI/TUI
go build -o cami cmd/cami/main.go

# Build MCP server
go build -o cami-mcp cmd/cami-mcp/main.go

# Build both
make build

# Run tests
go test ./...
```

### Contributing

See [reference/development-guide.md](reference/development-guide.md) for:
- Code organization patterns
- Adding new agents
- Extending MCP tools
- Testing guidelines

## Roadmap

### v0.2.0 (Current)
- ✅ Full CLI interface
- ✅ MCP server integration
- ✅ Agent discovery and updates
- ✅ Smart CLAUDE.md documentation
- ✅ JSON output for programmatic use

### v0.3.0 (Planned)
- Agent orchestration guidance
- Deployment history tracking
- Version rollback capability
- Agent dependency management
- Multi-project batch operations
- Agent usage analytics

### Future Considerations
- Web UI dashboard
- Agent marketplace/sharing
- Custom agent templates
- Integration with CI/CD pipelines
- Agent performance metrics
- Team collaboration features

## License

Proprietary - Product Team
