# CAMI - Claude Agent Management Interface

**MCP-First Architecture for Claude Code Integration**

CAMI is a Model Context Protocol (MCP) server that enables Claude Code to dynamically manage specialized AI agents. It provides a single binary with dual modes: MCP server for Claude Code integration (primary) and CLI for scripting/automation (secondary).

## Architecture Overview

### Single Binary, Dual Modes

```bash
# MCP Server Mode (primary) - for Claude Code
$ cami --mcp
# Runs as MCP server on stdio for Claude Code integration

# CLI Mode (secondary) - for scripting and quick checks
$ cami list
$ cami deploy frontend backend ~/projects/my-app
$ cami scan ~/projects/my-app
```

### Global Agent Storage

CAMI uses a global agent repository at `~/.cami/sources/` instead of per-project storage:

```
~/.cami/
├── config.yaml           # Global configuration
├── sources/              # Global agent sources
│   ├── team-agents/     # Team/company agents (optional)
│   └── my-agents/       # Personal custom agents (optional)
└── cami                 # Single binary (MCP + CLI)
```

**Benefits of global storage:**
- Agents available across all projects without duplication
- Single source of truth for agent versions
- Easier to update agents globally
- Simpler mental model

### Configuration Format

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
    path: /Users/lando/projects/my-project
  - name: client-project
    path: /Users/lando/clients/acme-app
```

**Priority-based deduplication**: When the same agent exists in multiple sources, the highest priority source wins (my-agents: 200 > team-agents: 100).

## MCP Server Configuration

### Claude Code Setup

Add to your Claude Code MCP settings (`~/.claude/settings.json` or project `.claude/settings.local.json`):

```json
{
  "mcpServers": {
    "cami": {
      "command": "~/.cami/cami",
      "args": ["--mcp"]
    }
  }
}
```

### Installation

```bash
# Download binary (future: Homebrew available)
curl -L https://github.com/lando-labs/cami/releases/latest/download/cami-macos \
  -o ~/.cami/cami
chmod +x ~/.cami/cami

# Optional: Add to PATH for CLI convenience
ln -s ~/.cami/cami /usr/local/bin/cami
```

## MCP Tools Reference

CAMI provides 12 MCP tools for Claude Code to manage agents. These tools enable natural language workflows like "Add the frontend agent to this project" or "What agents do I have?".

### Core Agent Management

#### 1. `mcp__cami__list_agents`
**Purpose**: List all available agents from configured sources
**Use when**: User asks "what agents are available?" or wants to discover agents
**Returns**: Agent names, versions, descriptions, categories, and source information

**Example context**:
```
User: "What agents are available?"
Claude: *uses mcp__cami__list_agents*
"I found X agents across Y sources..."
```

#### 2. `mcp__cami__deploy_agents`
**Purpose**: Deploy selected agents to a project's `.claude/agents/` directory
**Use when**: User wants to add agents to a project
**Parameters**:
- `agent_names` (array of strings, required) - Names of agents to deploy
- `target_path` (string, required) - Absolute path to project directory
- `overwrite` (boolean, optional) - Overwrite existing agents if they already exist

**Behavior**:
- Creates `.claude/agents/` directory if it doesn't exist
- Detects conflicts and asks for confirmation if overwrite=false
- Copies agent files from sources to project
- Returns deployment status for each agent

**Example context**:
```
User: "Add frontend and backend agents"
Claude: *uses mcp__cami__deploy_agents with current working directory*
"✓ Deployed frontend (v1.1.0)
 ✓ Deployed backend (v1.1.0)"
```

#### 3. `mcp__cami__scan_deployed_agents`
**Purpose**: Scan a project to see what agents are deployed and their status
**Use when**: User asks "what agents are installed?" or wants to audit/update agents
**Parameters**:
- `target_path` (string, required) - Absolute path to project directory

**Returns**: For each deployed agent:
- Agent name and version
- Status: `up-to-date`, `update-available`, or `not-in-sources`
- Available version if update exists

**Example context**:
```
User: "What agents do I have?"
Claude: *uses mcp__cami__scan_deployed_agents*
"Found 3 deployed agents:
 - frontend (v1.0.0) → Update available to v1.1.0
 - backend (v1.1.0) → Up to date
 - custom-agent (v1.0.0) → Not found in sources"
```

#### 4. `mcp__cami__update_claude_md`
**Purpose**: Update a project's CLAUDE.md with agent documentation
**Use when**: After deploying agents, to keep documentation in sync
**Parameters**:
- `target_path` (string, required) - Absolute path to project directory

**Behavior**:
- Scans `.claude/agents/` directory
- Reads frontmatter from each agent file
- Auto-generates "Deployed Agents" section
- Preserves other sections in CLAUDE.md
- Creates CLAUDE.md if it doesn't exist

**Example context**:
```
User: "Update my docs"
Claude: *uses mcp__cami__update_claude_md*
"✓ Updated CLAUDE.md with 3 deployed agents"
```

### Source Management

#### 5. `mcp__cami__list_sources`
**Purpose**: List all configured agent sources
**Use when**: User asks "where are my agents from?" or wants to see sources
**Returns**: For each source:
- Name and type
- Path
- Priority (for deduplication)
- Number of agents
- Git information (remote URL, status)

**Example context**:
```
User: "Show me my agent sources"
Claude: *uses mcp__cami__list_sources*
"You have 2 agent sources:
 1. lando-agents (100 priority) - 29 agents
    Path: ~/.cami/sources/lando-agents
    Git: git@github.com:lando-labs/lando-agents.git
 2. my-agents (200 priority) - 5 agents
    Path: ~/.cami/sources/my-agents
    Git: Not configured"
```

#### 6. `mcp__cami__add_source`
**Purpose**: Add a new agent source by cloning a Git repository
**Use when**: User wants to add official agents, company sources, or team libraries
**Parameters**:
- `url` (string, required) - Git URL (SSH or HTTPS)
- `name` (string, optional) - Source name (defaults to repo name)
- `priority` (number, optional) - Priority for deduplication (default: 100)

**Behavior**:
- Clones repository to `~/.cami/sources/<name>/`
- Updates `~/.cami/config.yaml`
- Scans for agents in cloned repository
- Returns agent count and source info

**Example context**:
```
User: "Add the official Lando agent library"
Claude: *uses mcp__cami__add_source with git@github.com:lando-labs/lando-agents.git*
"✓ Cloned lando-agents to ~/.cami/sources/lando-agents
 ✓ Found 29 agents"
```

#### 7. `mcp__cami__update_source`
**Purpose**: Update agent sources with git pull
**Use when**: User wants to get latest agents from sources
**Parameters**:
- `name` (string, optional) - Source name to update (updates all if not specified)

**Behavior**:
- Runs `git pull` on sources with git remotes
- Skips sources without git configured
- Returns update status for each source

**Example context**:
```
User: "Update my agents"
Claude: *uses mcp__cami__update_source*
"✓ Updated lando-agents (3 new commits)
 ⊘ Skipped my-agents (no git remote)"
```

#### 8. `mcp__cami__source_status`
**Purpose**: Show git status of agent sources
**Use when**: User wants to check if sources have uncommitted changes
**Returns**: For each source with git:
- Clean or has uncommitted changes
- Branch information
- Commit status

**Example context**:
```
User: "Check my agent sources"
Claude: *uses mcp__cami__source_status*
"lando-agents: Clean (on main, up to date with remote)
 my-agents: Modified (2 uncommitted changes)"
```

### Location Management

#### 9. `mcp__cami__add_location`
**Purpose**: Register a project directory for agent deployment tracking
**Use when**: User wants to track a project in CAMI
**Parameters**:
- `name` (string, required) - Friendly name for the location
- `path` (string, required) - Absolute path to project directory

**Behavior**:
- Adds location to `~/.cami/config.yaml`
- Validates path exists
- Prevents duplicate locations

**Example context**:
```
User: "Track this project"
Claude: *uses mcp__cami__add_location with current directory*
"✓ Added location 'my-app' → /Users/lando/projects/my-app"
```

#### 10. `mcp__cami__list_locations`
**Purpose**: List all registered project locations
**Use when**: User asks "what projects am I tracking?"
**Returns**: Location names and paths

**Example context**:
```
User: "What projects am I tracking?"
Claude: *uses mcp__cami__list_locations*
"You have 3 tracked locations:
 - my-app → /Users/lando/projects/my-app
 - client-site → /Users/lando/clients/acme
 - cami → /Users/lando/dev/cami"
```

#### 11. `mcp__cami__remove_location`
**Purpose**: Unregister a project directory
**Use when**: User wants to stop tracking a project
**Parameters**:
- `name` (string, required) - Location name to remove

**Example context**:
```
User: "Stop tracking my-app"
Claude: *uses mcp__cami__remove_location*
"✓ Removed location 'my-app'"
```

### Onboarding

#### 12. `mcp__cami__onboard`
**Purpose**: Get personalized onboarding guidance based on current setup
**Use when**: User is new to CAMI or asks "what should I do next?" or "help me get started"

**Returns**: Structured analysis of:
- Has config file?
- Has agent sources?
- Has deployed agents?
- Has tracked locations?
- Recommended next steps

**Behavior**: Context-aware recommendations:
- No config → "Let me help you set up CAMI by adding the official agent library"
- Has sources but no deployed agents → "You have agents available. Which would you like to add to this project?"
- Has deployed agents → "Everything looks good! You have X agents deployed."

**Example context**:
```
User: "Help me get started with CAMI"
Claude: *uses mcp__cami__onboard*
"I see CAMI isn't configured yet. Let me help you set it up.
 I'll add the official Lando agent library with 29 professional agents."
*uses mcp__cami__add_source*
"✓ Added lando-agents (29 agents available)
 Which agents would you like to add to your project?"
```

## Common Workflows for Claude Code

### First-Time Setup

```
User: "I want to start using CAMI"

Claude workflow:
1. Use mcp__cami__onboard → Detect no config
2. Use mcp__cami__add_source → Clone lando-agents
3. Use mcp__cami__list_agents → Show available agents
4. Ask user which agents they want
5. Use mcp__cami__deploy_agents → Deploy selected agents
6. Use mcp__cami__update_claude_md → Document deployment
```

### Adding Agents to Current Project

```
User: "Add frontend and backend agents"

Claude workflow:
1. Use mcp__cami__list_agents → Verify agents exist
2. Use mcp__cami__deploy_agents → Deploy to current directory
3. Use mcp__cami__update_claude_md → Update documentation
```

### Updating Agents

```
User: "Update my agents"

Claude workflow:
1. Use mcp__cami__update_source → Pull latest from git
2. Use mcp__cami__scan_deployed_agents → Check deployed status
3. If updates available → Ask user if they want to redeploy
4. Use mcp__cami__deploy_agents with overwrite=true → Update agents
5. Use mcp__cami__update_claude_md → Update documentation
```

### Auditing Project Agents

```
User: "What agents do I have?"

Claude workflow:
1. Use mcp__cami__scan_deployed_agents → Scan current project
2. Show agent names, versions, and update status
3. If updates available → Suggest updating
```

### Creating Custom Agents

```
User: "Help me create a new agent"

Claude workflow (with agent-architect):
1. Invoke agent-architect to design agent
2. Save agent file to ~/.cami/sources/my-agents/
3. Use mcp__cami__list_agents → Verify agent discovered
4. Use mcp__cami__deploy_agents → Deploy to project
5. Use mcp__cami__update_claude_md → Document it
```

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

**Note**: The CLI is designed for power users and automation. For most users, interacting through Claude Code via MCP tools is recommended.

## Agent Discovery and Priority

### Multi-Source Deduplication

When the same agent exists in multiple sources, CAMI uses priority-based deduplication:

```yaml
agent_sources:
  - name: lando-agents
    priority: 100        # Official agents (lower priority)

  - name: company-agents
    priority: 150        # Company-specific (medium priority)

  - name: my-agents
    priority: 200        # Personal overrides (highest priority)
```

**Example**: If "frontend" agent exists in all three sources, the version from `my-agents` (priority 200) is used.

### Agent Versioning

Each agent has a version in its frontmatter:

```markdown
---
name: frontend
version: 1.1.0
description: Use this agent when building user interfaces...
---
```

CAMI tracks versions to detect when updates are available via `scan_deployed_agents`.

## Internal Architecture (for AI understanding)

### Code Structure

```
cami/
├── cmd/cami/main.go           # Single binary entry point
│   ├── main()                 # Mode detection: --mcp or CLI
│   ├── runMCPServer()         # MCP server mode
│   └── runCLI()               # CLI mode
├── internal/
│   ├── agent/                 # Agent loading and parsing
│   │   ├── agent.go           # Agent struct and frontmatter
│   │   └── loader.go          # LoadAgentsFromSources()
│   ├── config/                # Configuration management
│   │   ├── config.go          # Config struct
│   │   └── loader.go          # Load ~/.cami/config.yaml
│   ├── deploy/                # Agent deployment
│   │   └── deploy.go          # Deploy agents to projects
│   ├── docs/                  # CLAUDE.md management
│   │   └── claude.go          # Update deployed agents section
│   ├── discovery/             # Agent scanning
│   │   └── discovery.go       # Scan .claude/agents/
│   ├── cli/                   # CLI commands
│   │   └── commands.go        # CLI command implementations
│   └── tui/                   # Terminal UI
│       └── tui.go             # Interactive deployment interface
└── ~/.cami/                   # User data directory
    ├── config.yaml            # Global configuration
    ├── sources/               # Agent sources
    └── cami                   # Binary
```

### Key Functions (for code navigation)

**Agent Loading** (`internal/agent/loader.go`):
- `LoadAgentsFromSources(cfg *config.Config) ([]Agent, error)` - Load all agents from configured sources
- Priority-based deduplication happens here

**Deployment** (`internal/deploy/deploy.go`):
- `DeployAgents(agents []string, targetPath string, overwrite bool) error` - Deploy agents to `.claude/agents/`
- Handles conflict detection and directory creation

**Documentation** (`internal/docs/claude.go`):
- `UpdateClaudeMD(targetPath string, agents []Agent) error` - Update CLAUDE.md
- Preserves non-CAMI-managed sections
- Uses `<!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: 2025-11-14T12:27:04-06:00 -->
## Deployed Agents

The following Claude Code agents are available in this project:

### agent-architect (v1.1.0)
Use this agent PROACTIVELY when you need to create, refine, or optimize Claude Code agent configurations. This includes designing new agents from scratch, improving existing agent system prompts, establishing agent interaction patterns, defining agent responsibilities and boundaries, or architecting multi-agent systems with clear separation of concerns.

### qa (v1.1.0)
Use this agent PROACTIVELY when writing tests, analyzing test coverage, creating testing documentation, or maintaining testing standards. Invoke for unit tests, integration tests, E2E tests, test coverage analysis, test strategy planning, or quality assurance automation.

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->

---

## Additional Resources

### Architecture Planning Documents
- [MCP-First Architecture Plan](reference/mcp-first-architecture-plan.md) - Philosophy and migration strategy
- [Clean MCP-First Plan](reference/clean-mcp-first-plan.md) - Implementation details and phases
- [Open Source Strategy](reference/open-source-strategy.md) - Path to public release
- [Agent Classification System](reference/agent-classification-system-design.md) - Future agent categorization design

### Development
- See [README.md](README.md) for development setup and build instructions
- See [.claude/agents/](/.claude/agents/) for deployed agent configurations
