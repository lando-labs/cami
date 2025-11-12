# Agent Workspace

This directory is your workspace for managing multiple agent sources.

## Structure

Each subdirectory can be:
- **Local workspace** (e.g., `my-agents/`) - Your personal agents
- **Cloned repository** (e.g., `lando-agents/`, `company-agents/`) - Remote sources

## Priority System

When multiple sources have the same agent name, the **highest priority wins**:
- **200**: Local experiments (my-agents/)
- **150**: Team/company sources
- **100**: Community sources

## Adding Sources

### Via Claude Code (Recommended)

```
@claude add this agent source: git@github.com:company/agents.git
```

### Manual Setup

```bash
cd vc-agents
git clone git@github.com:company/agents.git company
```

Then update `~/.cami/config.yaml` with the new source.

## Default Workspace

`my-agents/` is your default local workspace. Create agents here for experimentation,
then optionally contribute them to team or community sources later.

## Learn More

See [CAMI documentation](../README.md) for complete guide.
