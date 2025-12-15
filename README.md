# CAMI - Claude Agent Management Interface

**Your AI Agent Guild Headquarters**

CAMI is a Model Context Protocol (MCP) server that enables Claude Code to dynamically manage specialized AI agents across all your projects. Single binary, clean workspace, conversation-first interface.

## Features

- **25 MCP Tools**: Native Claude Code integration for complete agent and skill lifecycle management
- **Skillset Architecture**: Separate agents (methodology) from skills (implementation) for flexible tech stack support
- **Global Agent Storage**: Single source of truth at `~/cami-workspace/sources/`
- **Priority-Based Deduplication**: Override agents/skills with custom versions (lower priority number = higher precedence)
- **Deployment Tracking**: Automatic manifest creation tracking agent/skill versions, sources, and hashes
- **Normalization System**: Analyze and standardize project agent deployments
- **Smart Documentation**: Automatic CLAUDE.md updates with deployed agent information
- **Multiple Sources**: Manage agents and skills from Git repositories with priority-based loading
- **Git-Trackable Workspace**: Optionally version control your CAMI setup and custom agents
- **No Sudo Required**: User-local installation to `~/.local/bin/`

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
- `~/.local/bin/cami` - Binary (user-local, no sudo required)

The MCP server works immediately after install. For CLI usage from any directory, add `~/.local/bin` to your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/.local/bin:$PATH"
```

### Platform Notes

**WSL (Windows Subsystem for Linux)**

CAMI works on WSL2. Requirements:
- Go 1.21+ installed in WSL (not Windows)
- `sudo` access configured in WSL
- Git configured in WSL

Then follow the standard installation instructions above.

**macOS & Linux**

Works on both Intel and Apple Silicon (arm64). No sudo required - binary installs to `~/.local/bin/`.

### First-Time Setup

```bash
# Open your CAMI workspace
cd ~/cami-workspace
claude

# Ask Claude to help you get started
```

```
You: "Help me get started with CAMI"
Claude: *uses mcp__cami__onboard to guide you through setup*

You: "Create a new agent for handling database operations"
Claude: *works with agent-architect to create a custom agent*

You: "Deploy the agent to my project"
Claude: *uses mcp__cami__deploy_agents*
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
│   └── agents/                  # Bundled agents (agent-architect, skill-architect, etc.)
├── sources/                     # Agent sources
│   ├── my-agents/              # Your custom agents
│   ├── team-agents/            # (if added)
│   └── fullstack-guild/        # Example: guild added via add_source

~/.local/bin/cami               # Binary (user-local)
```

### Priority-Based Deduplication

When the same agent exists in multiple sources, **lower priority numbers win**:

```yaml
agent_sources:
  - name: my-agents
    priority: 10         # Highest priority (personal overrides)
  - name: team-agents
    priority: 50         # Medium priority (default)
  - name: fullstack-guild
    priority: 100        # Lowest priority (public guilds)
```

**Example**: If "frontend" agent exists in all three sources, the version from `my-agents` (priority 10) is used.

### Example Agent Guilds

Ready-to-use agent collections you can add to CAMI:

| Guild | Focus | Agents | Repository |
|-------|-------|--------|------------|
| game-dev-guild | Phaser 3 games | 8 | [lando-labs/game-dev-guild](https://github.com/lando-labs/game-dev-guild) |
| content-guild | Writing & marketing | 6 | [lando-labs/content-guild](https://github.com/lando-labs/content-guild) |
| fullstack-guild | MERN stack | 7 | [lando-labs/fullstack-guild](https://github.com/lando-labs/fullstack-guild) |

Add a guild to CAMI:
```bash
# Via Claude Code
"Add the game-dev-guild agent source"

# Via CLI
cami source add https://github.com/lando-labs/game-dev-guild.git
```

## MCP Tools

CAMI provides 25 MCP tools for Claude Code:

**Project Management**
- `create_project` - Create new project with agents and documentation
- `onboard` - Get personalized setup guidance

**Agent Management**
- `list_agents` - List all available agents from configured sources
- `deploy_agents` - Deploy agents to `.claude/agents/` (supports `with_skills` parameter)
- `scan_deployed_agents` - Check deployed agents and version status
- `update_claude_md` - Update CLAUDE.md with agent documentation

**Skill Management** (New in v0.5.0)
- `list_skills` - List available skills with filtering by tags/category
- `list_skillsets` - List skillset collections with tech stack info
- `deploy_skills` - Deploy skills to `.claude/skills/`
- `scan_deployed_skills` - Check deployed skills and version status
- `add_skill_source` - Add skill source repository
- `recommend_skills` - Suggest skills based on STRATEGIES.yaml tech stack

**Source Management**
- `list_sources` - List all configured agent sources with compliance status
- `add_source` - Add new source by cloning Git repository
- `update_source` - Pull latest from Git sources
- `source_status` - Check Git status of sources

**Location Management**
- `add_location` - Register project directory for tracking
- `list_locations` - List all tracked project locations
- `remove_location` - Unregister project directory

**Normalization**
- `detect_project_state` - Analyze project's CAMI integration level
- `normalize_project` - Create manifests and link agents to sources
- `detect_source_state` - Analyze source for CAMI compliance
- `normalize_source` - Fix source agents to meet CAMI standards
- `cleanup_backups` - Clean up old backup directories

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

`~/cami-workspace/config.yaml`:

```yaml
version: "1"
agent_sources:
  - name: team-agents
    type: local
    path: ~/cami-workspace/sources/team-agents
    priority: 50
    git:
      enabled: true
      remote: git@github.com:yourorg/team-agents.git

  - name: my-agents
    type: local
    path: ~/cami-workspace/sources/my-agents
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
"Help me get started with CAMI"
"Create a new agent for handling API integrations"
"Deploy the agent to ~/projects/my-app"
"What's the status of my deployed agents?"
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

**CAMI v0.5.0** - Current Release

### What's New in v0.5.0

#### Skillset Architecture (Major Feature)
CAMI now supports **separation of concerns** between agents (methodology/WHO) and skills (implementation/HOW):

- ✅ **Skills Support** - Deploy Claude Code skills alongside agents
  - Skills provide implementation patterns (code syntax, framework conventions)
  - Agents provide methodology (thinking patterns, quality standards, decision frameworks)
- ✅ **Skillset Collections** - Group related skills with optional SKILLSET.yaml metadata
- ✅ **6 New MCP Tools** for skill management:
  - `list_skills` - List available skills with filtering
  - `list_skillsets` - List skillset collections
  - `deploy_skills` - Deploy skills to `.claude/skills/`
  - `scan_deployed_skills` - Check deployed skill status
  - `add_skill_source` - Add skill source repositories
  - `recommend_skills` - Suggest skills based on STRATEGIES.yaml tech stack
- ✅ **Enhanced `deploy_agents`** - Deploy skills alongside agents with `with_skills` parameter
- ✅ **Tech Stack Matching** - Skillsets match STRATEGIES.yaml technologies via tags

#### New Bundled Agents
- ✅ **skill-architect v2.0.0** - Creates skills and skillsets with optimized descriptions
- ✅ **skilled-agent-architect v1.0.0** - Creates methodology-focused agents designed to work with skills

#### Installation Improvements
- ✅ **No sudo required** - Binary now installs to `~/.local/bin/` instead of `/usr/local/bin`
- ✅ **PATH detection** - Installer warns if `~/.local/bin` isn't in PATH

### Previous Release (v0.4.0)
- Agent Classification System (workflow-specialist, technology-implementer, strategic-planner)
- agent-architect v4.0.0 with auto-classification and model selection
- Enhanced onboarding with deployment counts
- Smart installer upgrades with automatic backups
- STRATEGIES.yaml documentation location

### Current Status

**Alpha Testing** - v0.5.0
- Skillset architecture complete and tested
- Agent-skill separation validated via POC
- Ready for early adopter testing
- Official agent guilds available (game-dev, content, fullstack)

### Roadmap

**v0.6.0 (Planned)**
- Remote agent sources (HTTP, direct Git URLs)
- Enhanced update detection with semantic versioning
- Team collaboration features
- Skill marketplace/discovery

## License

MIT License - See [LICENSE](LICENSE) file for details
