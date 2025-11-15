# CAMI - Claude Agent Management Interface

**MCP-First Architecture for Claude Code Integration**

CAMI is a Model Context Protocol (MCP) server that enables Claude Code to dynamically manage specialized AI agents. Single binary with dual modes: MCP server for Claude Code integration (primary) and CLI for scripting/automation (secondary).

## Features

- **MCP Server Integration**: 13 MCP tools for native Claude Code workflows
- **Global Agent Storage**: Single source of truth at `~/.cami/sources/`
- **Priority-Based Deduplication**: Override agents with custom versions (lower priority number = higher precedence)
- **Smart Documentation**: Automatic CLAUDE.md updates with deployed agent information
- **Version Tracking**: Compare deployed versions with available updates
- **Multiple Sources**: Manage agents from Git repositories with priority-based loading
- **.camiignore Support**: Flexible file filtering with glob patterns

## Quick Start

### Zero-Setup Mode (Recommended for Development)

```bash
# Clone and open
git clone <cami-repo-url>
cd cami
# Open in Claude Code
```

That's it! CAMI automatically runs via `go run` when you use Claude Code in this directory.

Try it: Ask Claude "Help me get started with CAMI"

See [QUICKSTART.md](QUICKSTART.md) for details.

### Production Installation

For using CAMI across multiple projects:

```bash
# Build and install
go build -o ~/.cami/cami ./cmd/cami

# Add to your project's .mcp.json
{
  "mcpServers": {
    "cami": {
      "command": "~/.cami/cami",
      "args": ["--mcp"]
    }
  }
}
```

### First-Time Setup

Open Claude Code and interact naturally:

```
You: "Help me get started with CAMI"
Claude: *uses mcp__cami__onboard*

You: "Add agent source from git@github.com:yourorg/agents.git"
Claude: *uses mcp__cami__add_source*

You: "Add frontend and backend agents to this project"
Claude: *uses mcp__cami__deploy_agents*
```

CAMI creates `~/.cami/config.yaml` automatically and deploys agents to `.claude/agents/` in your project.

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

### Global Agent Storage

```
~/.cami/
├── config.yaml           # Global configuration
├── sources/              # Global agent sources
│   ├── team-agents/     # Team/company agents
│   └── my-agents/       # Personal custom agents
└── cami                 # Single binary
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

## Common Workflows

**Via Claude Code (recommended):**

```
"Create a new meditation app project"
→ Claude gathers requirements, recommends agents, creates project

"What agents do I have?"
→ Claude scans deployed agents and shows versions

"Update my agents"
→ Claude pulls latest from sources and offers to redeploy updates

"Add the QA agent to this project"
→ Claude deploys agent and updates CLAUDE.md
```

**Via CLI:**

```bash
# Set up new project
cami source add git@github.com:yourorg/agents.git
cami deploy frontend backend ~/my-project
cami update-docs ~/my-project

# Update agents across projects
cami source update
cami scan ~/my-project  # Check for updates
```

## Development

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
│   └── cli/               # CLI commands
├── .claude/agents/        # Deployed agents for CAMI development
└── README.md              # This file
```

### Building

```bash
# Build binary
go build -o cami ./cmd/cami

# Run tests
go test ./...
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
