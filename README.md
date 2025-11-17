# CAMI - Claude Agent Management Interface

**Your AI Agent Guild Headquarters**

CAMI is a Model Context Protocol (MCP) server that enables Claude Code to dynamically manage specialized AI agents across all your projects. Single binary, clean workspace, conversation-first interface.

## Features

- **13 MCP Tools**: Native Claude Code integration for agent management
- **Global Agent Storage**: Single source of truth at `~/cami-workspace/sources/`
- **Priority-Based Deduplication**: Override agents with custom versions (lower priority number = higher precedence)
- **Smart Documentation**: Automatic CLAUDE.md updates with deployed agent information
- **Version Tracking**: Compare deployed versions with available updates
- **Multiple Sources**: Manage agents from Git repositories with priority-based loading
- **Git-Trackable Workspace**: Optionally version control your CAMI setup and custom agents

## Installation

### Download & Install

**Coming soon: Homebrew, direct downloads**

For now, build from source:

```bash
# Clone the repository
git clone https://github.com/lando-labs/cami.git
cd cami

# Build and install
make install
```

This creates:
- `~/cami-workspace/` - Your CAMI workspace
- `/usr/local/bin/cami` - Binary on your PATH

### First-Time Setup

```bash
# Open your CAMI workspace
cd ~/cami-workspace
claude

# Ask Claude to help you get started
```

```
You: "Help me get started with CAMI"
Claude: *uses mcp__cami__onboard to guide you*

You: "Add the official agent library"
Claude: *uses mcp__cami__add_source*

You: "What agents are available?"
Claude: *uses mcp__cami__list_agents*
```

That's it! CAMI will guide you through adding agent sources and deploying agents to your projects.

## Architecture

### Single Binary, Dual Modes

```bash
# MCP Server Mode (primary) - for Claude Code
$ cami --mcp
# Runs as MCP server on stdio

# CLI Mode (secondary) - for scripting
$ cami list
$ cami deploy frontend backend ~/projects/my-app
$ cami scan ~/projects/my-app
```

### Workspace Structure

```
~/cami-workspace/                          # Your CAMI workspace
├── CLAUDE.md                    # CAMI documentation and persona
├── README.md                    # Quick start guide
├── .mcp.json                    # MCP server configuration
├── .gitignore                   # Git ignore rules
├── config.yaml                  # CAMI configuration
├── .claude/
│   └── agents/                  # CAMI's own agents
├── sources/                     # Agent sources
│   ├── official-agents/        # (if added)
│   ├── team-agents/            # (if added)
│   └── my-agents/              # Your custom agents

/usr/local/bin/cami             # Binary on PATH
```

### Priority-Based Deduplication

When the same agent exists in multiple sources, **lower priority numbers win**:

```yaml
agent_sources:
  - name: my-agents
    priority: 10         # Highest priority (personal overrides)
  - name: team-agents
    priority: 50         # Medium priority (default)
  - name: official-agents
    priority: 100        # Lowest priority (public agents)
```

**Example**: If "frontend" agent exists in all three sources, the version from `my-agents` (priority 10) is used.

## MCP Tools

CAMI provides 13 MCP tools for Claude Code:

**Project Creation**
- `create_project` - Create new project with agents and documentation

**Agent Management**
- `list_agents` - List all available agents from configured sources
- `deploy_agents` - Deploy selected agents to `.claude/agents/`
- `scan_deployed_agents` - Check deployed agents and version status
- `update_claude_md` - Update CLAUDE.md with agent documentation

**Source Management**
- `list_sources` - List all configured agent sources
- `add_source` - Add new source by cloning Git repository
- `update_source` - Pull latest from Git sources
- `source_status` - Check Git status of sources

**Location Management**
- `add_location` - Register project directory for tracking
- `list_locations` - List all tracked project locations
- `remove_location` - Unregister project directory

**Onboarding**
- `onboard` - Get personalized setup guidance

See [CLAUDE.md](CLAUDE.md) for complete MCP tool documentation and workflows.

## CLI Commands

For scripting and automation:

```bash
# Agent management
cami list                        # List available agents
cami deploy <agents> <path>      # Deploy agents to project
cami scan <path>                 # Scan deployed agents
cami update-docs <path>          # Update CLAUDE.md

# Source management
cami source list                 # List agent sources
cami source add <git-url>        # Add new source
cami source update [name]        # Update sources (git pull)
cami source status               # Check git status

# Location management
cami locations list              # List tracked locations
cami locations add <name> <path> # Add location
cami locations remove <name>     # Remove location
```

## Agent Structure

Agents are markdown files with YAML frontmatter:

```markdown
---
name: frontend
version: "1.1.0"
description: Use this agent when building user interfaces...
---

# Frontend Agent

You are a specialized frontend development expert...
```

## Configuration

`~/.cami/config.yaml`:

```yaml
version: "1"
agent_sources:
  - name: team-agents
    type: local
    path: ~/.cami/sources/team-agents
    priority: 50
    git:
      enabled: true
      remote: git@github.com:yourorg/team-agents.git

  - name: my-agents
    type: local
    path: ~/.cami/sources/my-agents
    priority: 10
    git:
      enabled: false

deploy_locations:
  - name: my-project
    path: /Users/username/projects/my-project
```

## .camiignore Support

Exclude files from agent loading with `.camiignore` in source directories:

```
# Documentation
README.md
LICENSE.md

# Patterns
*.txt
docs/

# Hidden files
.git/
.github/
```

## Using CAMI

### Working in Your CAMI Workspace

```bash
cd ~/cami-workspace
claude

# Natural language interface
"Add the official agent library"
"What agents are available?"
"Deploy frontend and backend agents to ~/projects/my-app"
"Create a custom database agent"
```

### Git Tracking (Optional)

Track your CAMI workspace to share setup with your team:

```bash
cd ~/cami-workspace
git init
git add .
git commit -m "Initial CAMI setup"
git remote add origin <your-repo-url>
git push -u origin main
```

The included `.gitignore` is configured to:
- ✅ Track your custom agents in `sources/my-agents/`
- ❌ Ignore pulled sources (managed by CAMI)
- ? Your choice on `config.yaml` (remove from .gitignore to track)

### CLI Commands

CAMI commands work from anywhere:

```bash
# Agent management
cami list                           # List available agents
cami deploy <agents> <path>         # Deploy agents to project
cami scan <path>                    # Scan deployed agents
cami update-docs <path>             # Update CLAUDE.md

# Source management
cami source list                    # List agent sources
cami source add <git-url>           # Add new source
cami source update [name]           # Update sources (git pull)
cami source status                  # Check git status

# Location management
cami locations list                 # List tracked locations
cami locations add <name> <path>    # Add location
cami locations remove <name>        # Remove location
```

### Global MCP Setup (Optional)

To use CAMI from any Claude Code session (not just ~/cami-workspace/):

Add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "cami": {
      "command": "cami",
      "args": ["--mcp"]
    }
  }
}
```

## Development

**Contributing to CAMI? Welcome!**

### Development Setup

```bash
# Clone the repository
git clone https://github.com/lando-labs/cami.git
cd cami

# Open in Claude Code (dev mode with go run)
claude
```

The `.mcp.json` in this repo uses `go run` for zero-setup development.

### Project Structure

```
cami/
├── cmd/cami/main.go       # Single binary entry point
├── internal/
│   ├── agent/             # Agent loading and parsing
│   ├── config/            # Configuration management
│   ├── deploy/            # Agent deployment
│   ├── docs/              # CLAUDE.md management
│   ├── discovery/         # Agent scanning
│   ├── cli/               # CLI commands
│   ├── mcp/               # MCP server implementation
│   └── tui/               # Terminal UI
├── install/
│   ├── templates/         # User workspace templates
│   └── install.sh         # Installation script
├── .claude/agents/        # Deployed agents for CAMI development
├── .mcp.json              # Dev mode: go run
├── Makefile               # Build, test, release targets
└── README.md              # This file
```

### Build & Test

```bash
# Build binary
make build

# Build for all platforms
make release-all

# Package releases with installer
make package

# Run tests
make test

# Run linters
make lint

# Install locally (creates ~/cami-workspace/ workspace)
make install
```

### Testing User Experience

```bash
# Install to ~/cami-workspace/
make install

# Test user workspace
cd ~/cami-workspace
claude

# Ask: "Help me get started with CAMI"
```

## Documentation

- **[README.md](README.md)** - This file (getting started)
- **[CLAUDE.md](CLAUDE.md)** - Complete MCP tool documentation and workflows

## Version

**CAMI v0.3.0** - Current Release

### What's New in v0.3.0
- ✅ Single binary with dual modes (MCP + CLI)
- ✅ 13 MCP tools for complete Claude Code integration
- ✅ Project creation workflow with `create_project` tool
- ✅ Global agent storage at `~/.cami/sources/`
- ✅ Inverted priority system (1 = highest, 100 = lowest)
- ✅ Priority-based multi-source deduplication
- ✅ Source management tools (add, update, status)
- ✅ .camiignore support with glob patterns

### Roadmap

**v0.4.0 (Planned)**
- Agent classification system
- Remote agent sources (HTTP, direct Git URLs)
- Enhanced agent-architect integration
- Team collaboration features

## License

MIT License - See [LICENSE](LICENSE) file for details
