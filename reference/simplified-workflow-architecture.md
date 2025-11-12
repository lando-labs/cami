# CAMI Simplified Workflow Architecture

**Version:** 1.0 (MCP-first, minimal CLI)
**Date:** 2025-11-10

---

## Core Philosophy

> **CLI does the minimum. MCP + Claude handles everything else.**

**Why?**
- Every team has different Git workflows (solo, team, fork, etc.)
- Claude is better at guiding users through context-specific decisions
- CAMI provides building blocks, Claude orchestrates

---

## What CAMI Provides

### Minimal CLI
- ✅ `cami init` - One question: where to store agents
- ✅ `cami build` - Build the CLI/MCP
- ✅ Help text pointing to MCP

### Core MCP Tools
- ✅ `onboard` - Interactive Claude-guided setup
- ✅ `source_add` - Clone remote sources to vc-agents/
- ✅ `source_list` - Show all sources
- ✅ `source_update` - Git pull all sources
- ✅ `source_status` - Git status for all sources
- ✅ `deploy_agents` - Deploy to projects
- ✅ `list_agents` - Browse available agents
- ✅ `scan_deployed_agents` - Check what's deployed

### Agent-Architect Integration
- ✅ Shipped in `.claude/agents/agent-architect.md`
- ✅ Available immediately after `cami init`
- ✅ Workspace-aware (knows about vc-agents structure)

---

## Repository Structure

### CAMI Repo (What Gets Cloned)

```
lando-labs/cami/
├── cmd/
│   ├── cami/          # Minimal CLI
│   └── cami-mcp/      # MCP server with onboarding
├── internal/          # Go packages
├── vc-agents/         # Agent workspace
│   ├── .gitignore     # **/  (ignore all subdirs)
│   ├── README.md      # Explains workspace concept
│   └── my-agents/     # Created by cami init
│       └── .gitkeep
├── .claude/
│   └── agents/
│       └── agent-architect.md   # Shipped with CAMI
├── go.mod
└── README.md          # Getting Started points to MCP onboarding
```

### After User Setup

```
cami/
├── ...
├── vc-agents/
│   ├── .gitignore            # **/
│   ├── my-agents/           # Default local (no git unless user adds)
│   │   └── my-agent.md
│   ├── company-agents/      # User added via source_add
│   │   ├── .git/
│   │   └── ...
│   └── community-agents/    # User added via source_add
│       ├── .git/
│       └── ...
└── ~/.cami/
    └── config.yaml
```

---

## User Journey

### Step 1: Install CAMI

```bash
git clone git@github.com:lando-labs/cami.git
cd cami
go build -o cami ./cmd/cami
go build -o cami-mcp ./cmd/cami-mcp
```

### Step 2: Initialize (Minimal CLI)

```bash
./cami init
```

**Single question:**
```
Welcome to CAMI! 🚀

CAMI manages Claude Code agents in a workspace directory.

? Where should we store your agents?
  > ./vc-agents (default, in this directory)
    Specify custom path

✓ Created vc-agents/my-agents (default local workspace)
✓ Configuration saved to ~/.cami/config.yaml

Next steps:
  1. Ensure MCP is configured in Claude Code settings
  2. Start Claude in this directory: claude
  3. Ask: @claude help me onboard with CAMI

For more info: cat README.md
```

**That's it for CLI!** Everything else is MCP + Claude.

### Step 3: MCP Configuration

**User's `~/.claude/claude_desktop_config.json`:**
```json
{
  "mcpServers": {
    "cami": {
      "command": "/Users/lando/lando-labs/cami/cami-mcp"
    }
  }
}
```

### Step 4: Claude-Guided Onboarding

```bash
cd ~/lando-labs/cami
claude
```

**User:**
```
@claude help me onboard with CAMI
```

**Claude (uses `mcp__cami__onboard` tool):**
```
Welcome to CAMI! Let me help you get set up.

I see you've already run `cami init` - great!

✓ Found agent-architect in .claude/agents/
✓ Default workspace: vc-agents/my-agents

Let's set up your agent workflow.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CAMI's power comes from agent-architect, which researches
domains and creates sophisticated agents tailored to your needs.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Do you have any remote agent sources to add?
(e.g., company repository, community collections)

> No, just my local workspace for now

Great! Your setup:
  ✓ Local workspace: vc-agents/my-agents (priority 200)
  ✓ Agent-architect available: @agent-architect

Let's create your first agent to test the workflow.

What kind of agent would you like to create?
```

**User:**
```
A Terraform specialist
```

**Claude (uses agent-architect):**
```
@agent-architect create a Terraform specialist agent

[Agent-architect researches and creates agent]

✓ Created: vc-agents/my-agents/terraform-specialist.md

The agent is stored in your local workspace.

Would you like to:
  1. Deploy it to a project now?
  2. Test it here first?
  3. Just keep it in the workspace?

> Deploy it to a project

Which project? (provide path or I can list recent directories)

> ~/my-terraform-project

[Deploys terraform-specialist to ~/my-terraform-project/.claude/agents/]

✓ Deployed terraform-specialist to ~/my-terraform-project

To use it, you can either:
  1. Continue here: claude --continue
  2. Start new session: cd ~/my-terraform-project && claude

The agent will be available as @terraform-specialist.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Onboarding complete! 🎉

Quick reference:
  • Create agents: @agent-architect create a [description]
  • Deploy agents: Ask me to deploy agents to projects
  • Add remote sources: Ask me to add a GitHub repository
  • Check status: Ask me to show agent source status

Try creating more agents or ask me anything!
```

---

## MCP Tools (Full Spec)

### 1. `onboard`

**Description:** Interactive onboarding guide with Claude

**Parameters:** None (conversational)

**What it does:**
- Checks if `cami init` was run
- Explains agent-architect
- Offers to add remote sources
- Walks through creating first agent
- Explains deployment workflow
- Provides quick reference

**Implementation:**
```go
// Returns onboarding state and guidance text
type OnboardingState struct {
    InitComplete   bool
    HasAgentArch   bool
    SourceCount    int
    AgentCount     int
    NextStep       string
}
```

### 2. `source_add`

**Description:** Add remote agent source (clones to vc-agents/)

**Parameters:**
- `url` (required): Git URL (e.g., git@github.com:company/agents.git)
- `name` (optional): Source name (derived from URL if omitted)
- `priority` (optional): Priority (default: 150 for remote)

**What it does:**
```bash
# User asks Claude: "Add the company agent repository git@github.com:acme/agents.git"

# Claude calls: source_add(url: "git@github.com:acme/agents.git", name: "acme", priority: 150)

# CAMI executes:
cd vc-agents
git clone git@github.com:acme/agents.git acme

# Updates config.yaml:
agent_sources:
  - name: acme
    type: local
    path: /Users/lando/cami/vc-agents/acme
    priority: 150
    git:
      enabled: true
      remote: git@github.com:acme/agents.git
```

**Returns:**
```json
{
  "success": true,
  "source": "acme",
  "path": "vc-agents/acme",
  "agent_count": 15
}
```

### 3. `source_list`

**Description:** List all agent sources

**Parameters:** None

**Returns:**
```json
{
  "sources": [
    {
      "name": "my-agents",
      "type": "local",
      "path": "vc-agents/my-agents",
      "priority": 200,
      "git_enabled": false,
      "agent_count": 3
    },
    {
      "name": "acme",
      "type": "local",
      "path": "vc-agents/acme",
      "priority": 150,
      "git_enabled": true,
      "git_remote": "git@github.com:acme/agents.git",
      "git_status": "clean",
      "agent_count": 15
    }
  ]
}
```

### 4. `source_update`

**Description:** Update sources (git pull)

**Parameters:**
- `source` (optional): Specific source name, or all if omitted

**What it does:**
```bash
# User asks: "Update all agent sources"

# Claude calls: source_update()

# CAMI executes:
cd vc-agents/acme && git pull
cd vc-agents/other-source && git pull
# (skips my-agents - no git remote)
```

**Returns:**
```json
{
  "updated": [
    {
      "source": "acme",
      "status": "up-to-date"
    },
    {
      "source": "community",
      "status": "updated",
      "new_commits": 3
    }
  ],
  "skipped": [
    {
      "source": "my-agents",
      "reason": "no git remote"
    }
  ]
}
```

### 5. `source_status`

**Description:** Show git status for all sources

**Parameters:** None

**Returns:**
```json
{
  "sources": [
    {
      "name": "my-agents",
      "git_enabled": false,
      "status": "not a git repo"
    },
    {
      "name": "acme",
      "git_enabled": true,
      "status": "clean",
      "branch": "main",
      "ahead": 0,
      "behind": 0
    },
    {
      "name": "community",
      "git_enabled": true,
      "status": "dirty",
      "branch": "main",
      "uncommitted_changes": 2,
      "modified_files": ["agent1.md", "agent2.md"]
    }
  ]
}
```

---

## Agent-Architect Workspace Awareness

**Enhanced system prompt for agent-architect:**

```markdown
## CAMI Workspace Context

You are running in a CAMI-managed workspace with multiple agent sources:

**Workspace structure:**
- `vc-agents/my-agents/` - User's local workspace (priority 200)
- `vc-agents/[other]/` - Additional sources (priorities vary)

**When creating agents:**
1. Ask user where to create the agent (if ambiguous)
2. Default to `my-agents` for experiments
3. Create in appropriate category subdirectory

**Available sources:**
[CAMI injects current sources list]

**Example:**
User: "Create a Terraform agent"

You:
I'll create a Terraform specialist agent. Where should I store it?
  - my-agents (your local workspace) ← default
  - acme (team repository)

[User selects my-agents]

Creating in: vc-agents/my-agents/infrastructure/terraform-specialist.md

[Creates agent with research]

✓ Created terraform-specialist.md

This agent is in your local workspace. To deploy it:
  - Use via CAMI MCP: Ask Claude to deploy it
  - Manual: Copy to project/.claude/agents/

To contribute to a team source later:
  cp vc-agents/my-agents/infrastructure/terraform-specialist.md vc-agents/acme/
  cd vc-agents/acme
  git add terraform-specialist.md
  git commit -m "Add Terraform specialist"
  git push
```

---

## Configuration File

**`~/.cami/config.yaml`:**

```yaml
# CAMI Configuration
version: 1

# Agent sources (priority: higher = higher precedence)
agent_sources:
  - name: my-agents
    type: local
    path: /Users/lando/cami/vc-agents/my-agents
    priority: 200
    git:
      enabled: false  # Not a git repo by default

  - name: acme-internal
    type: local
    path: /Users/lando/cami/vc-agents/acme-internal
    priority: 150
    git:
      enabled: true
      remote: git@github.com:acme-corp/agents.git

# Project deployment locations (optional, for convenience)
deploy_locations:
  - name: my-app
    path: /Users/lando/projects/my-app
```

---

## Documentation Updates

### README.md (in CAMI repo)

```markdown
# CAMI - Claude Agent Management Interface

Manage Claude Code agents with workspace-based version control.

## Quick Start

1. **Clone and build:**
   ```bash
   git clone git@github.com:lando-labs/cami.git
   cd cami
   go build -o cami ./cmd/cami
   go build -o cami-mcp ./cmd/cami-mcp
   ```

2. **Initialize:**
   ```bash
   ./cami init
   ```

3. **Configure MCP** (add to `~/.claude/claude_desktop_config.json`):
   ```json
   {
     "mcpServers": {
       "cami": {
         "command": "/absolute/path/to/cami/cami-mcp"
       }
     }
   }
   ```

4. **Start Claude and onboard:**
   ```bash
   claude
   ```

   Then ask:
   ```
   @claude help me onboard with CAMI
   ```

Claude will guide you through the rest!

## Core Concepts

**Workspace:** `vc-agents/` directory containing multiple agent sources
**Sources:** Git repos or local directories with agents
**Priorities:** Higher number = higher precedence (200 > 150 > 100)
**Agent-Architect:** Meta-agent that creates sophisticated agents via research

## What CAMI Does

✅ Clone remote sources to workspace
✅ Pull updates from sources
✅ Show git status
✅ Deploy agents to projects
✅ Track agent versions

## What CAMI Doesn't Do

❌ Git add/commit/push (use standard Git)
❌ Create PRs (use `gh` CLI)
❌ Manage auth (use Git credentials)

**Why?** Every team has different workflows. CAMI provides building blocks.

## Common Workflows

**Create an agent:**
```
@agent-architect create a [description]
```

**Add remote source:**
```
@claude add this agent repository: git@github.com:company/agents.git
```

**Deploy agents:**
```
@claude deploy the terraform agent to ~/my-project
```

**Update sources:**
```
@claude update all agent sources
```

**Check status:**
```
@claude show agent source status
```

## Documentation

- [Simplified Workflow](./reference/simplified-workflow-architecture.md)
- [Contribution Philosophy](./reference/contribution-philosophy.md)
- [Open Source Strategy](./reference/open-source-strategy.md)
```

---

## CLI Commands (Minimal)

### `cami init`

**Purpose:** One-time setup (storage location only)

**Usage:**
```bash
./cami init
```

**Interactive:**
```
Welcome to CAMI! 🚀

? Where should we store your agents?
  > ./vc-agents (default)
    Custom path: ___

✓ Created vc-agents/my-agents
✓ Configuration saved to ~/.cami/config.yaml

Next: Start Claude and ask for onboarding help
```

### `cami help`

**Purpose:** Show help text pointing to MCP

**Output:**
```
CAMI - Claude Agent Management Interface

Basic Commands:
  cami init          Initialize workspace (one-time setup)
  cami help          Show this help

Most operations are done via MCP + Claude Code:
  1. Configure CAMI MCP in Claude settings
  2. Start Claude: claude
  3. Ask: @claude help me onboard with CAMI

Documentation: https://github.com/lando-labs/cami
```

### `cami version`

**Purpose:** Show version

**Output:**
```
CAMI v1.0.0
```

---

## Open Source Implications

### What Changes

**Simplified onboarding:**
- No complex CLI wizard
- One question: where to store
- Everything else via MCP + Claude

**No "official" Lando agents (yet):**
- CAMI ships with agent-architect only
- Users create their own agents on demand
- Later: `lando-labs/lando-agents` as community collection
- Users can opt-in to add it as source

**MCP-first design:**
- CLI is minimal (init, help, version)
- All management via MCP tools
- Claude guides users through context-specific decisions

### Getting Started (Open Source Users)

```bash
# 1. Install
git clone git@github.com:lando-labs/cami.git
cd cami && go build -o cami ./cmd/cami && go build -o cami-mcp ./cmd/cami-mcp

# 2. Initialize
./cami init

# 3. Configure MCP (one time)
# Add to ~/.claude/claude_desktop_config.json

# 4. Start Claude
claude

# 5. Ask for help
@claude help me onboard with CAMI
```

Done! Claude handles the rest.

---

## Summary

**What we've simplified:**

1. ✅ **CLI:** Minimal (init, help, version)
2. ✅ **Onboarding:** One question (storage location)
3. ✅ **Management:** All via MCP + Claude
4. ✅ **Git:** Clone, pull, status only
5. ✅ **Workflow:** Claude guides based on user context
6. ✅ **No "official" repo:** Agent-architect creates on demand
7. ✅ **Workspace:** vc-agents/ with multiple sources

**Key insight:**
> Let Claude orchestrate. CAMI provides building blocks.

**User experience:**
- Install CAMI → `cami init` → Start Claude → Ask for help → Done
- Claude explains concepts as needed
- Workflows adapt to user context (solo dev, team, contributor)

**No more:**
- ❌ Complex CLI wizards
- ❌ Prescribed workflows
- ❌ One-size-fits-all assumptions
