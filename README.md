# CAMI - Claude Agent Management Interface

**MCP-First Architecture for Claude Code Integration**

CAMI is a Model Context Protocol (MCP) server that enables Claude Code to dynamically manage specialized AI agents. It provides a single binary with dual modes: MCP server for Claude Code integration (primary) and CLI for scripting/automation (secondary).

## Features

- **Version-Controlled Agents**: Manage specialized agents from Git repositories
- **MCP Server Integration**: 12 MCP tools for native Claude Code workflows
- **Global Agent Storage**: Single source of truth at `~/.cami/sources/`
- **Priority-Based Deduplication**: Override agents with custom versions from higher-priority sources
- **Smart Documentation**: Automatic CLAUDE.md updates with deployed agent information
- **Interactive TUI**: Beautiful keyboard-driven interface for quick deployments
- **Full CLI Support**: Programmatic deployment and management commands
- **Location Management**: Track and manage deployment across multiple projects
- **Version Tracking**: Compare deployed versions with available updates
- **Conflict Detection**: Safe deployment with existing file detection
- **.camiignore Support**: Flexible file filtering with glob patterns

## Quick Start

### 1. Build the Binary

```bash
cd /path/to/cami
go build -o cami ./cmd/cami
```

### 2. Initialize CAMI

```bash
# Add your agent repository
mkdir -p ~/.cami/sources
cd ~/.cami/sources
git clone <your-agent-repo-url> my-agents

# Create configuration
cat > ~/.cami/config.yaml << 'EOF'
version: "1"
agent_sources:
  - name: my-agents
    type: local
    path: ~/.cami/sources/my-agents
    priority: 100
    git:
      enabled: true
      remote: <your-agent-repo-url>

deploy_locations:
  - name: my-project
    path: /Users/yourname/projects/my-project
EOF
```

### 3. Configure for Claude Code

Add `.mcp.json` to your project:

```json
{
  "mcpServers": {
    "cami": {
      "command": "/absolute/path/to/cami",
      "args": ["--mcp"]
    }
  }
}
```

Or configure globally in `~/.claude.json`:

```json
{
  "projects": {
    "/your/project/path": {
      "mcpServers": {
        "cami": {
          "type": "stdio",
          "command": "/absolute/path/to/cami",
          "args": ["--mcp"],
          "env": {}
        }
      }
    }
  }
}
```

## Architecture

### Single Binary, Dual Modes

```bash
# MCP Server Mode (primary) - for Claude Code
$ cami --mcp
# Runs as MCP server on stdio for Claude Code integration

# CLI Mode (secondary) - for scripting and quick checks
$ cami list
$ cami deploy frontend backend ~/projects/my-app
$ cami scan ~/projects/my-app

# Interactive TUI (no args)
$ cami
# Launches terminal UI for browsing and deployment
```

### Global Agent Storage

CAMI uses a global agent repository at `~/.cami/sources/` instead of per-project storage:

```
~/.cami/
├── config.yaml           # Global configuration
├── sources/              # Global agent sources
│   ├── team-agents/     # Team/company agents (optional)
│   └── my-agents/       # Personal custom agents (optional)
└── cami                 # Single binary (MCP + CLI + TUI)
```

**Benefits of global storage:**
- Agents available across all projects without duplication
- Single source of truth for agent versions
- Easier to update agents globally
- Simpler mental model

### Priority-Based Deduplication

When the same agent exists in multiple sources, CAMI uses priority-based deduplication:

```yaml
agent_sources:
  - name: team-agents
    priority: 100        # Team/company agents (lower priority)

  - name: my-agents
    priority: 200        # Personal overrides (highest priority)
```

**Example**: If "frontend" agent exists in both sources, the version from `my-agents` (priority 200) is used.

## MCP Tools Reference

CAMI provides 13 MCP tools for Claude Code to manage agents. These tools enable natural language workflows like "Create a new project" or "What agents do I have?".

### Project Creation

1. **`mcp__cami__create_project`** - Create a new project with directory setup, agent deployment, and CLAUDE.md

### Core Agent Management

2. **`mcp__cami__list_agents`** - List all available agents from configured sources
3. **`mcp__cami__deploy_agents`** - Deploy selected agents to a project's `.claude/agents/` directory
4. **`mcp__cami__scan_deployed_agents`** - Scan a project to see what agents are deployed and their status
5. **`mcp__cami__update_claude_md`** - Update a project's CLAUDE.md with agent documentation

### Source Management

6. **`mcp__cami__list_sources`** - List all configured agent sources
7. **`mcp__cami__add_source`** - Add a new agent source by cloning a Git repository
8. **`mcp__cami__update_source`** - Update agent sources with git pull
9. **`mcp__cami__source_status`** - Show git status of agent sources

### Location Management

10. **`mcp__cami__add_location`** - Register a project directory for agent deployment tracking
11. **`mcp__cami__list_locations`** - List all registered project locations
12. **`mcp__cami__remove_location`** - Unregister a project directory

### Onboarding

13. **`mcp__cami__onboard`** - Get personalized onboarding guidance based on current setup

Example usage:
```
User: "Add frontend and backend agents to this project"
Claude: *uses mcp__cami__deploy_agents*
"✓ Deployed frontend (v1.1.0)
 ✓ Deployed backend (v1.1.0)"

User: "What agents are available?"
Claude: *uses mcp__cami__list_agents*
"I found X agents across Y sources..."
```

See [CLAUDE.md](CLAUDE.md) for complete MCP tool documentation and workflows.

## CLI Commands (Secondary Interface)

While MCP is the primary interface, CAMI provides CLI commands for scripting and quick checks:

```bash
# Agent management
cami list                           # List available agents
cami deploy <agents...> <path>      # Deploy agents to project
cami scan <path>                    # Scan deployed agents
cami update-docs <path>             # Update CLAUDE.md

# Source management
cami source list                    # List agent sources
cami source add <git-url>           # Add new source
cami source update [name]           # Update sources
cami source status                  # Check git status

# Location management
cami locations list                 # List tracked locations
cami locations add <name> <path>    # Add location
cami locations remove <name>        # Remove location

# Interactive TUI
cami                                # Launch TUI for deployment
cami --help                         # Show full help
```

## Agent Management

CAMI manages agents from Git repositories that you configure. Agents are markdown files with YAML frontmatter that define:

- **Name & Version**: Semantic versioning for tracking updates
- **Description**: What the agent specializes in
- **System Prompt**: Instructions for Claude Code when the agent is invoked

**Agent Structure:**
```markdown
---
name: example-agent
version: "1.0.0"
description: Use this agent when...
---

# Agent Instructions

Your specialized instructions here...
```

After adding agent sources, use `cami list` (CLI) or `mcp__cami__list_agents` (MCP) to see available agents and their descriptions.

## .camiignore Support

Each agent source can include a `.camiignore` file to exclude documentation and non-agent files:

```
# CAMI Ignore File
# Patterns for files that should not be loaded as agents

# Documentation files
README.md
CONTRIBUTING.md
LICENSE.md

# Patterns
*.txt
docs/

# Hidden files
.git/
.github/
```

Supports glob patterns (`*.md`, `docs/`, etc.) and comments (`#`).

## Configuration Format

`~/.cami/config.yaml`:

```yaml
version: "1"
agent_sources:
  - name: team-agents
    type: local
    path: ~/.cami/sources/team-agents
    priority: 100
    git:
      enabled: true
      remote: git@github.com:yourorg/team-agents.git

  - name: my-agents
    type: local
    path: ~/.cami/sources/my-agents
    priority: 200
    git:
      enabled: false

deploy_locations:
  - name: my-project
    path: /Users/username/projects/my-project
  - name: client-project
    path: /Users/username/clients/acme-app
```

## Agent Versioning

Each agent follows semantic versioning (MAJOR.MINOR.PATCH) in its frontmatter:

```markdown
---
name: frontend
version: "1.1.0"
description: Use this agent when building user interfaces...
---
```

CAMI tracks versions to detect when updates are available via `scan_deployed_agents`.

## Deployment Workflow

When you deploy agents, CAMI:

1. Creates `.claude/agents/` directory in the target location if needed
2. Copies selected agent files with YAML frontmatter and content
3. Detects conflicts with existing files (safe by default)
4. Shows deployment results with success/error/conflict status
5. Optionally updates `CLAUDE.md` with deployed agent documentation

### Example Workflows

**Via MCP (Claude Code):**
```
User: "Add frontend and backend agents"
Claude: Uses mcp__cami__deploy_agents → mcp__cami__update_claude_md

User: "What agents do I have?"
Claude: Uses mcp__cami__scan_deployed_agents

User: "Update my agents"
Claude: Uses mcp__cami__update_source → mcp__cami__scan_deployed_agents
```

**Via CLI:**
```bash
# Deploy agents to a project
cami deploy --agents frontend,backend,qa --location ~/projects/my-app

# Update project documentation
cami update-docs --location ~/projects/my-app

# Verify deployment
cami scan --location ~/projects/my-app
```

**Via TUI:**
1. Launch `cami`
2. Select agents with `space`
3. Press `d` to deploy
4. Choose location
5. Press `i` to view and update deployed agents

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

## Development

### Project Structure

```
cami/
├── cmd/cami/main.go       # Single binary entry point
│   ├── main()             # Mode detection: --mcp, CLI, or TUI
│   ├── runMCPServer()     # MCP server mode
│   ├── runCLI()           # CLI mode
│   └── runTUI()           # TUI mode
├── internal/
│   ├── agent/             # Agent loading and parsing
│   │   ├── agent.go       # Agent struct and frontmatter
│   │   └── loader.go      # LoadAgentsFromSources(), .camiignore support
│   ├── config/            # Configuration management
│   │   ├── config.go      # Config struct
│   │   └── loader.go      # Load ~/.cami/config.yaml
│   ├── deploy/            # Agent deployment
│   │   └── deploy.go      # Deploy agents to projects
│   ├── docs/              # CLAUDE.md management
│   │   └── claude.go      # Update deployed agents section
│   ├── discovery/         # Agent scanning
│   │   └── discovery.go   # Scan .claude/agents/
│   ├── cli/               # CLI commands
│   │   └── commands.go    # CLI command implementations
│   └── tui/               # Terminal UI
│       └── tui.go         # Interactive deployment interface
├── .claude/agents/        # Deployed agents for CAMI development
├── reference/             # Technical documentation
└── README.md              # This file
```

### Building

```bash
# Build single binary
go build -o cami ./cmd/cami

# Install to PATH
go install ./cmd/cami

# Run tests
go test ./...
```

### Architecture

**Single Binary, Three Modes:**
- Entry: `cmd/cami/main.go`
- Mode detection: `--mcp` flag → MCP server, arguments → CLI, no args → TUI
- Shared internal packages for core logic

**MCP Server:**
- Protocol: Model Context Protocol over stdio
- Tools: 12 tools for agent management, source management, and locations
- Integration: Direct usage of internal packages

**Shared Core:**
- `internal/agent` - YAML frontmatter parsing, .camiignore support
- `internal/deploy` - File operations and conflict detection
- `internal/docs` - Smart CLAUDE.md merging
- `internal/discovery` - Version comparison logic
- `internal/config` - Global configuration management

## Documentation

- **[README.md](README.md)** - This file, main entry point
- **[CLAUDE.md](CLAUDE.md)** - Complete MCP tool documentation and workflows
- **[reference/](reference/)** - Detailed technical documentation
  - [mcp-first-architecture-plan.md](reference/mcp-first-architecture-plan.md) - Architecture philosophy
  - [clean-mcp-first-plan.md](reference/clean-mcp-first-plan.md) - Implementation details
  - [open-source-strategy.md](reference/open-source-strategy.md) - Path to public release
  - [agent-classification-system-design.md](reference/agent-classification-system-design.md) - Future categorization

## Version

**CAMI v0.3.0** - Current Release

### What's New in v0.3.0
- ✅ Single binary with dual modes (MCP + CLI/TUI)
- ✅ 13 MCP tools for complete Claude Code integration
- ✅ Project creation workflow with create_project tool
- ✅ Global agent storage at `~/.cami/sources/`
- ✅ Priority-based multi-source deduplication
- ✅ Source management tools (add, update, status)
- ✅ Location tracking across projects
- ✅ .camiignore support with glob patterns
- ✅ Comprehensive MCP-first documentation

### Roadmap

**v0.4.0 (Planned)**
- Agent classification system (3 tiers)
- Enhanced agent-architect with versioning strategy
- Remote agent sources (GitHub, Git, HTTP)
- Multi-workflow support (developers, consumers, teams)

**Future Considerations**
- Homebrew installation
- Agent marketplace/sharing
- Custom agent templates
- CI/CD pipeline integration
- Team collaboration features

## License

MIT License - See [LICENSE](LICENSE) file for details
