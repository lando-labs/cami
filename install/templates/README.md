# CAMI Workspace

Welcome to your CAMI workspace! This is your headquarters for managing Claude Code agents across all your projects.

## Quick Start

### First Time Setup

Open Claude Code in this directory and ask:

```
"Help me get started with CAMI"
```

CAMI will guide you through adding agent sources and configuring your setup.

### Common Tasks

**Deploy agents to a project:**
```
"Add frontend and backend agents to ~/projects/my-app"
```

**List available agents:**
```
"What agents are available?"
```

**Create a custom agent:**
```
"Help me create a new database agent"
```

**Update agent sources:**
```
"Update my agent sources"
```

## Directory Structure

```
~/cami-workspace/
├── CLAUDE.md                    # CAMI documentation and Claude context
├── .mcp.json                    # MCP server configuration (local to this directory)
├── config.yaml                  # CAMI configuration
├── .claude/
│   └── agents/                  # CAMI's own agents (qa, agent-architect, etc.)
├── sources/                     # Agent sources
│   ├── my-agents/              # Your custom agents
│   ├── team-agents/            # Team agents (if added)
│   └── fullstack-guild/        # Example: guild added via add_source
└── README.md                    # This file
```

## CLI Commands

CAMI commands work from anywhere:

```bash
# List available agents
cami list

# Deploy agents to a project
cami deploy frontend backend ~/projects/my-app

# Check deployed agents
cami scan ~/projects/my-app

# Manage sources
cami source list
cami source add git@github.com:yourorg/agents.git
cami source update

# Track project locations
cami locations add my-app ~/projects/my-app
cami locations list
```

## Agent Sources

Agents are organized in `sources/` directory:

- **my-agents/** - Your custom agents (tracked in git if you initialize this directory)
  - Includes a `.camiignore` file to exclude non-agent files from loading
- **[guild-name]/** - Agent guilds added via `add_source` (e.g., fullstack-guild, content-guild)
- **team-agents/** - Your team's shared agents (pulled from remote)

### Official Agent Guilds

Lando Labs maintains public agent guilds you can add:

| Guild | Focus |
|-------|-------|
| `fullstack-guild` | MERN stack web development |
| `content-guild` | Writing & marketing |
| `game-dev-guild` | Phaser 3 game development |

Add a guild: `"Add the fullstack-guild"`

### Adding Agent Sources

Use CAMI to add new sources:

```
"Add agent source from git@github.com:yourorg/agents.git"
```

Or via CLI:

```bash
cami source add git@github.com:yourorg/agents.git --priority 50
```

### Priority System

When the same agent exists in multiple sources, **lower priority number wins**:

- Priority 1 = Highest (overrides everything)
- Priority 50 = Default (standard sources)
- Priority 100 = Lowest (fallback defaults)

Example: If "frontend" agent exists in both `my-agents` (priority 10) and `fullstack-guild` (priority 100), the version from `my-agents` is used.

### Excluding Files with .camiignore

Each source directory can have a `.camiignore` file (like `.gitignore`) to exclude files from agent loading:

```
# sources/my-agents/.camiignore
README.md
*.txt
templates/
*-draft.md
```

The `my-agents/` directory includes a template `.camiignore` with common exclusions.

## Git Tracking (Optional)

Track this workspace with Git to share your CAMI setup:

```bash
git init
git add .
git commit -m "Initial CAMI setup"
git remote add origin <your-repo-url>
git push -u origin main
```

The `.gitignore` is configured to:
- ✅ Track your custom agents in `sources/my-agents/`
- ❌ Ignore pulled sources (managed by CAMI)
- ? Your choice on `config.yaml` (remove from .gitignore to track)

## Global MCP Setup (Optional)

To use CAMI from any Claude Code session (not just this directory):

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

This makes CAMI's MCP tools available in all projects while still using the agent sources configured here in `~/cami-workspace/`.

## Getting Help

- Open Claude Code in this directory and ask questions
- Run `cami --help` for CLI usage
- See `CLAUDE.md` for complete documentation
- Visit https://github.com/lando-labs/cami

## Next Steps

1. **Add agent sources**: Ask CAMI to add official or team agent libraries
2. **Deploy agents**: Add agents to your current projects
3. **Create custom agents**: Build specialized agents for your specific needs
4. **Track with Git**: Share your setup with your team (optional)

**Ready to build your agent guild!** 🎯
