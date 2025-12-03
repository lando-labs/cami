# CAMI Development - Claude Agent Management Interface

**CAMI Source Code Repository**

This is the CAMI development repository. You are working on the Go codebase that powers CAMI's MCP server and CLI.

---

## Claude Context: CAMI Developer Assistant

**You are assisting with CAMI development** - working on the Go codebase that implements the Model Context Protocol (MCP) server and CLI for Claude Code agent management.

**Your Focus Areas:**

1. **Go Development** - Write clean, idiomatic Go code following the project's patterns and conventions.

2. **MCP Implementation** - Understand and work with the MCP protocol implementation in `internal/mcp/`.

3. **Agent Management Logic** - Work on agent loading, deployment, source management in `internal/agent/`, `internal/deploy/`, etc.

4. **Cross-Platform Support** - Ensure code works on macOS, Linux, and Windows.

5. **Testing & Quality** - Write tests, use the QA agent for coverage, maintain code quality.

**Development Workflow:**

- Use `go run ./cmd/cami` for local testing
- The `.mcp.json` in this repo uses `go run` for zero-setup development
- Run `make build` to compile binary
- Run `make test` for tests
- Run `make lint` for code quality checks
- Use deployed agents (qa, agent-architect) for development tasks

**Remember:** This is the CAMI source code repo. Users install CAMI via releases and get a clean `~/cami-workspace/` workspace with templates from `install/templates/`.

---

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

### User Workspace Structure

CAMI creates a user workspace at `~/cami-workspace/` during installation:

```
~/cami-workspace/                          # User workspace (created by installer)
├── CLAUDE.md                    # User-facing CAMI documentation
├── README.md                    # User quick start guide
├── .mcp.json                    # Local MCP config
├── .gitignore                   # Git ignore rules
├── config.yaml                  # CAMI configuration
├── .claude/
│   └── agents/                  # CAMI's own agents
├── sources/                     # Agent sources
│   ├── my-agents/              # User's custom agents
│   ├── team-agents/            # (if added)
│   └── fullstack-guild/        # Example: guild added via add_source

/usr/local/bin/cami             # Binary (on PATH)
```

**Key Concepts:**
- User workspace is separate from source code repo
- Templates live in `install/templates/` in this repo
- Installer copies templates to `~/cami-workspace/` during setup
- Users can optionally track `~/cami-workspace/` with Git

### Configuration Format

User's `~/cami-workspace/config.yaml`:

```yaml
version: "1"
install_timestamp: 2025-01-18T10:30:00Z  # When CAMI was installed
setup_complete: true                      # Whether initial setup is complete
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
    path: ~/projects/my-project

default_projects_dir: ~/projects  # Where new projects are created by default
```

**New in v0.3.2**:
- `install_timestamp`: Automatically set during installation (used for robust fresh install detection)
- `setup_complete`: False on fresh install, true after adding real agent sources
- Existing configs are automatically migrated on first load

**Priority-based deduplication**: When the same agent exists in multiple sources, the lowest priority number wins (my-agents: 10 > team-agents: 50). Priority 1 = highest, 100 = lowest.

## Development Setup

### This Repository (Dev Mode)

This repo has `.mcp.json` configured for development:

```json
{
  "mcpServers": {
    "cami": {
      "command": "go",
      "args": ["run", "./cmd/cami", "--mcp"]
    }
  }
}
```

Open this directory in Claude Code and CAMI runs automatically via `go run`.

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make release-all

# Package releases with installer
make package

# Install locally (creates ~/cami-workspace/ workspace)
make install
```

### Testing

```bash
# Run tests
make test

# Run linters
make lint

# Test installation
make install
cd ~/cami-workspace
claude  # Test user experience
```

## Installation Templates

User workspace files are in `install/templates/`:

```
install/
├── templates/
│   ├── CLAUDE.md        # User-facing CAMI persona doc
│   ├── README.md        # User quick start guide
│   ├── .gitignore       # Git ignore for user workspace
│   └── .mcp.json        # Local MCP config (uses 'cami' on PATH)
└── install.sh           # Installation script
```

When modifying the user experience:
- Update templates in `install/templates/`
- Update installer in `install/install.sh`
- Test with `make install` and check `~/cami-workspace/`

## Project Structure

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
└── README.md              # Main documentation
```

## MCP Tools Reference

CAMI provides 19 MCP tools for Claude Code to manage agents. These tools enable natural language workflows like "Add the frontend agent to this project" or "What agents do I have?".

**Tool Categories**:
- **Core Agent Management** (4 tools): Deploy, scan, list, and document agents
- **Source Management** (4 tools): Add, update, list, and check status of agent sources
- **Location Management** (3 tools): Track project directories for deployment
- **Normalization** (7 tools): Analyze and fix source/project compliance
- **Onboarding** (1 tool): Get personalized setup guidance

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

#### 5. `mcp__cami__list_sources` (Updated with Compliance Checking)
**Purpose**: List all configured agent sources with compliance status
**Use when**: User asks "where are my agents from?" or wants to see sources
**Returns**: For each source:
- Name and type
- Path
- Priority (for deduplication)
- Number of agents
- Git information (remote URL, status)
- **NEW**: Compliance status (✓ compliant or ⚠️ issues)
- **NEW**: Issue summary if non-compliant

**Example context**:
```
User: "Show me my agent sources"
Claude: *uses mcp__cami__list_sources*
"You have 2 agent sources:
 • ✓ fullstack-guild (priority 100) - 7 agents
   Path: ~/cami-workspace/sources/fullstack-guild
   Git: https://github.com/lando-labs/fullstack-guild.git (clean)
   Compliance: ✓ Compliant

 • ⚠️ team-agents (priority 50) - 15 agents
   Path: ~/cami-workspace/sources/team-agents
   Git: Not configured
   Compliance: ⚠️ Issues (no .camiignore, 3 agents with issues)"
```

#### 6. `mcp__cami__add_source` (Updated with Auto-Detection)
**Purpose**: Add a new agent source by cloning a Git repository with automatic compliance detection
**Use when**: User wants to add official agents, company sources, or team libraries
**Parameters**:
- `url` (string, required) - Git URL (SSH or HTTPS)
- `name` (string, optional) - Source name (defaults to repo name)
- `priority` (number, optional) - Priority for deduplication (lower = higher priority, default: 50)

**Behavior**:
- Clones repository to `~/cami-workspace/sources/<name>/`
- Updates `~/cami-workspace/config.yaml`
- Scans for agents in cloned repository
- **NEW**: Automatically checks compliance after cloning
- **NEW**: Reports issues if source is non-compliant
- **NEW**: Recommends using `normalize_source` to fix issues

**Example context**:
```
User: "Add the fullstack-guild"
Claude: *uses mcp__cami__add_source with url="https://github.com/lando-labs/fullstack-guild.git"*
"✓ Cloned fullstack-guild to ~/cami-workspace/sources/fullstack-guild
 ✓ Found 7 agents

 ## Source Compliance Check

 ✓ **Source is compliant**"
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
"✓ Updated fullstack-guild (3 new commits)
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
"fullstack-guild: Clean (on main, up to date with remote)
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
- Adds location to `~/cami-workspace/config.yaml`
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
 What kind of development do you primarily do?"
User: "Full stack web apps"
Claude: "I recommend the fullstack-guild - it has agents for React, Express, MongoDB, and more."
*uses mcp__cami__add_source*
"✓ Added fullstack-guild (7 agents available)
 Which agents would you like to add to your project?"
```

---

### Normalization (Phase 1 Complete - NEW!)

CAMI's normalization system ensures agent sources and projects follow CAMI standards for tracking, versioning, and deployment management.

#### 13. `mcp__cami__detect_source_state`
**Purpose**: Analyze an agent source for CAMI compliance
**Use when**: Need to check if a source follows CAMI standards before normalization
**Parameters**:
- `source_name` (string, required) - Name of the source to analyze

**Returns**:
- Source path and agent count
- Compliance status (true/false)
- List of issues (missing versions, descriptions, names)
- Missing .camiignore indicator
- Recommended actions

**Example context**:
```
User: "Check if my team-agents source is compliant"
Claude: *uses mcp__cami__detect_source_state*
"# Source Analysis: team-agents

**Path:** ~/cami-workspace/sources/team-agents
**Agent Count:** 15
**Compliant:** false

⚠️ Missing .camiignore file

## Issues Found (3)

**frontend.md:**
  - missing version

**backend.md:**
  - missing version
  - missing description

## Recommended Action

Use `normalize_source` to fix these issues automatically."
```

---

#### 14. `mcp__cami__normalize_source`
**Purpose**: Fix source agents to meet CAMI standards
**Use when**: Source has compliance issues that need to be fixed
**Parameters**:
- `source_name` (string, required) - Source to normalize
- `add_versions` (boolean, optional) - Add v1.0.0 to agents missing versions
- `add_descriptions` (boolean, optional) - Add description placeholders
- `create_camiignore` (boolean, optional) - Create .camiignore file

**Features**:
- Creates backup before making changes
- Adds missing versions (default: v1.0.0)
- Adds description placeholders ("Description for X agent")
- Creates .camiignore with common patterns
- Reports all changes made

**Example context**:
```
User: "Fix the issues in my team-agents source"
Claude: *uses mcp__cami__normalize_source with add_versions=true, create_camiignore=true*
"# Source Normalization: team-agents

**Status:** ✓ Success
**Agents Updated:** 3
**Backup Created:** ~/cami-workspace/sources/.cami-backup-20251118-143000

## Changes Made

- Added version 1.0.0 to frontend.md
- Added version 1.0.0 to backend.md
- Added description placeholder to backend.md
- Created .camiignore file"
```

---

#### 15. `mcp__cami__detect_project_state`
**Purpose**: Analyze a project's normalization state
**Use when**: Need to understand project's CAMI integration level
**Parameters**:
- `project_path` (string, required) - Path to project directory

**Returns**:
- Project state (non-cami, cami-aware, cami-legacy, cami-native)
- Agents directory presence
- Manifest presence
- Agent count and details
- Version matching with sources
- Upgrade availability
- Normalization recommendations

**Project States**:
- `non-cami`: No `.claude/agents/` directory
- `cami-aware`: Has agents but no manifest
- `cami-legacy`: Has old CAMI format manifest
- `cami-native`: Fully normalized with current manifest format

**Example context**:
```
User: "What's the state of this project?"
Claude: *uses mcp__cami__detect_project_state with current directory*
"# Project Analysis

**Path:** /Users/lando/projects/my-app
**State:** cami-aware
**Has Agents Directory:** true
**Has Manifest:** false
**Agent Count:** 3

## Deployed Agents

**frontend** (v1.0.0) - matches available-sources
**backend** (v1.1.0) - matches available-sources (update available)
**custom-agent** (no version) - not in sources

## Recommendations

✓ **Minimal normalization required:** Create manifests for tracking
✓ **Standard normalization recommended:** Link agents to sources"
```

---

#### 16. `mcp__cami__normalize_project`
**Purpose**: Normalize project by creating manifests and linking agents
**Use when**: Project needs manifest creation or source linking
**Parameters**:
- `project_path` (string, required) - Project directory path
- `level` (string, required) - Normalization level: "minimal", "standard", or "full"

**Normalization Levels**:
- **minimal**: Create manifests only (basic tracking)
- **standard**: Manifests + source links (recommended)
- **full**: Complete rewrite with agent-architect (not yet implemented)

**Features**:
- Creates backup before normalization
- Creates local manifest (`.claude/cami-manifest.yaml`)
- Updates central manifest (`~/cami-workspace/deployments.yaml`)
- Links agents to sources with priority tracking
- Calculates content and metadata hashes
- Detects agents needing upgrades
- Provides undo capability via backup

**Example context**:
```
User: "Normalize this project with standard tracking"
Claude: *uses mcp__cami__normalize_project with level="standard"*
"# Project Normalization

**Status:** ✓ Success
**State Before:** cami-aware
**State After:** cami-native
**Backup Created:** /Users/lando/projects/.cami-backup-20251118-143500

## Changes Made

- Created project manifest with source links

**Undo available:** Use backup.RestoreFromBackup to revert changes"
```

---

#### 17. `mcp__cami__cleanup_backups`
**Purpose**: Clean up old backup directories
**Use when**: Too many backups are consuming disk space
**Parameters**:
- `target_path` (string, required) - Path to directory with backups
- `keep_recent` (number, optional) - Number of recent backups to keep (default: 3)

**Features**:
- Analyzes backup state before cleanup
- Shows total backup count and size
- Keeps N most recent backups
- Removes older backups
- Reports freed disk space
- Lists remaining backups

**Example context**:
```
User: "I have too many backups, clean them up"
Claude: *uses mcp__cami__cleanup_backups with keep_recent=3*
"# Backup Cleanup

**Total Backups:** 15
**Total Size:** 45.30 MB

## Cleanup Results

**Removed:** 12 backups
**Freed:** 38.50 MB
**Kept:** 3 backups

**Remaining backups:**
- .cami-backup-20251118-143500
- .cami-backup-20251118-142000
- .cami-backup-20251118-140000"
```

---

#### 18. `mcp__cami__create_project`
**Purpose**: Create a new project with proper CAMI setup
**Use when**: User wants to start a new project with agents
**Parameters**:
- `name` (string, required) - Project name
- `description` (string, required) - Project description
- `agent_names` (array of strings, required) - Agents to deploy
- `path` (string, optional) - Project path (defaults to current directory)
- `vision_doc` (string, optional) - Vision document content

**Workflow**:
1. Gathers requirements from user
2. Recommends agents based on requirements
3. Creates project directory structure
4. Deploys selected agents
5. Creates initial manifest
6. Writes vision document

**Example context**:
```
User: "I want to create a new web app project"
Claude: "Let me help you create a new project. What tech stack are you using?"
User: "React frontend, Node.js backend"
Claude: *uses mcp__cami__list_agents to find relevant agents*
"I recommend these agents:
 - frontend (for React development)
 - backend (for Node.js APIs)
 - qa (for testing)

 Would you like to use these?"
User: "Yes"
Claude: *uses mcp__cami__create_project*
"✓ Created project 'my-web-app'
 ✓ Deployed 3 agents
 ✓ Created manifest
 ✓ Wrote vision document"
```

---

#### 19. `mcp__cami__update_claude_md` (Enhanced)
**Purpose**: Update project's CLAUDE.md with deployed agent documentation
**Use when**: After deploying agents or normalizing project
**Parameters**:
- `target_path` (string, required) - Path to project directory

**Behavior**:
- Scans `.claude/agents/` directory
- Reads frontmatter from each agent
- Auto-generates "Deployed Agents" section
- Preserves other sections
- Creates CLAUDE.md if missing
- **NEW**: Includes normalization status if manifest exists

**Example context**:
```
User: "Update my docs"
Claude: *uses mcp__cami__update_claude_md*
"✓ Updated CLAUDE.md with 3 deployed agents
 ✓ Added manifest tracking information"
```

---

## Common Workflows for Claude Code

### First-Time Setup

```
User: "I want to start using CAMI"

Claude workflow:
1. Use mcp__cami__onboard → Detect no config
2. Ask user about development focus
3. Use mcp__cami__add_source → Clone appropriate guild (fullstack-guild, content-guild, or game-dev-guild)
4. Use mcp__cami__list_agents → Show available agents
5. Ask user which agents they want
6. Use mcp__cami__deploy_agents → Deploy selected agents
7. Use mcp__cami__update_claude_md → Document deployment
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
2. Save agent file to ~/cami-workspace/sources/my-agents/
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
  - name: fullstack-guild
    priority: 100        # Public guilds (lowest priority)

  - name: team-agents
    priority: 50         # Team-specific (medium priority)

  - name: my-agents
    priority: 10         # Personal overrides (highest priority)
```

**Example**: If "frontend" agent exists in all three sources, the version from `my-agents` (priority 10) is used because lower numbers = higher priority.

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

## Agent Classification System

### Overview

CAMI v0.4.0+ includes a three-class agent system that organizes agents by cognitive model and work style. Each class has different phase weights for the universal three-phase methodology (Research → Execute → Validate).

### The Three Classes

| Class | User-Friendly Name | Purpose | Phase Weights |
|-------|-------------------|---------|---------------|
| **workflow-specialist** | Task Automator | Execute specific, user-defined workflows | Research 15% → Execute 70% → Validate 15% |
| **technology-implementer** | Feature Builder | Build complete capabilities in specific domains | Research 30% → Execute 55% → Validate 15% |
| **strategic-planner** | System Architect | Architect systems, research, optimize at scale | Research 45% → Execute 30% → Validate 25% |

### Agent Schema

Agents include optional `class` and `specialty` fields in frontmatter:

```yaml
---
name: frontend
version: "1.1.0"
description: Use this agent PROACTIVELY when building user interfaces...
class: technology-implementer
specialty: react-development
---
```

**Fields**:
- `class` (optional): One of `workflow-specialist`, `technology-implementer`, or `strategic-planner`
- `specialty` (optional): Domain/specialty (e.g., `kubernetes-operations`, `react-development`, `system-architecture`)

### Code Implementation

**Agent Struct** (`internal/agent/agent.go`):
```go
type Agent struct {
    Name        string `yaml:"name"`
    Version     string `yaml:"version"`
    Description string `yaml:"description"`
    Class       string `yaml:"class,omitempty"`
    Specialty   string `yaml:"specialty,omitempty"`
    // ...
}
```

**Phase Weights** (`internal/agent/agent.go`):
```go
// GetPhaseWeights returns phase distribution for an agent
func (a *Agent) GetPhaseWeights() PhaseWeights {
    return GetPhaseWeightsByClass(a.Class)
}

// GetPhaseWeightsByClass maps class to phase percentages
func GetPhaseWeightsByClass(class string) PhaseWeights {
    // Returns default balanced weights if class not specified
}

// GetUserFriendlyClassName returns "Task Automator", "Feature Builder", etc.
func GetUserFriendlyClassName(class string) string {
    // Maps technical names to user-friendly names
}
```

### Auto-Classification

**Agent-architect v3.0.0+** automatically classifies agents based on request signals:

- **Workflow Specialist**: Requests mention checklists, procedures, workflows, repeatable processes
- **Technology Implementer**: Requests about building features, implementing capabilities, domain-specific work
- **Strategic Planner**: Requests about architecture, research, planning, optimization, high-level decisions

### Workflow Gathering (CAMI's Role)

CAMI (not agent-architect) gathers class-specific information before invoking agent-architect:

**For Workflow Specialists**:
- Offers 3 input methods: describe conversationally, provide file, or point to docs
- Gathers step-by-step workflow with success/failure criteria
- Confirms workflow before agent creation

**For Technology Implementers**:
- Asks about technology/framework versions
- Identifies integration points
- Notes patterns and conventions

**For Strategic Planners**:
- Understands constraints (timeline, budget, scale, team size)
- Identifies key tradeoffs
- Maps current state vs desired state

This separation ensures agent-architect focuses on agent generation while CAMI handles conversational requirements gathering.

### Backward Compatibility

- Fields are optional (`omitempty` in YAML)
- Existing agents without class/specialty continue to work
- Default balanced weights (30/50/20) used if class not specified
- Parser handles both old and new frontmatter formats

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
│   │   └── loader.go          # Load ~/cami-workspace/config.yaml
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
└── ~/cami-workspace/                   # User data directory
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
