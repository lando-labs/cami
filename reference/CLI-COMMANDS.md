# CAMI CLI Commands

CAMI now supports both interactive TUI mode and CLI subcommands for programmatic agent management.

## Quick Start

```bash
# Build CAMI
make build

# Launch interactive TUI (no arguments)
./cami

# Use CLI commands (with subcommands)
./cami locations                                                # List deployment locations
./cami location add --name my-project --path ~/my-project      # Add a location
./cami deploy --agents frontend,backend --location ~/my-project # Deploy agents
./cami update-docs --location ~/my-project                      # Update CLAUDE.md
./cami list                                                     # List available agents
./cami scan --location ~/my-project                             # Scan deployed agents
```

## Commands

### locations

List all configured deployment locations.

**Flags:**
- `--output`: Output format: `text` or `json` (default: text)

**Examples:**

```bash
# List all deployment locations (text output)
./cami locations

# List with JSON output
./cami locations --output json
```

**JSON Output Format:**

```json
{
  "locations": [
    {
      "name": "my-project",
      "path": "/Users/lando/projects/my-app"
    },
    {
      "name": "client-app",
      "path": "/Users/lando/clients/app"
    }
  ],
  "count": 2
}
```

**When No Locations Configured:**

```
No deployment locations configured.

To add a location:
  cami location add --name <name> --path <path>
```

### location add

Add a new deployment location to the configuration.

**Flags:**
- `-n, --name` (required): Unique name for the location
- `-p, --path` (required): Absolute path to the project directory

**Examples:**

```bash
# Add a deployment location
./cami location add --name my-project --path /Users/lando/projects/my-app

# Add using short flags
./cami location add -n client-app -p ~/clients/app
```

**Validation:**
- Location name must be unique
- Location path must be unique
- Path must exist and be a directory
- Path can use `~` for home directory (will be expanded)

**Success Output:**

```
Successfully added location 'my-project' -> /Users/lando/projects/my-app
```

**Error Examples:**

```
Error: location with name "my-project" already exists
Error: location with path "/Users/lando/projects/my-app" already exists
Error: path does not exist: /Users/lando/projects/nonexistent
Error: path is not a directory: /Users/lando/projects/file.txt
```

### location remove

Remove a deployment location from the configuration.

**Flags:**
- `-n, --name` (required): Name of the location to remove

**Examples:**

```bash
# Remove a deployment location
./cami location remove --name my-project

# Remove using short flags
./cami location remove -n client-app
```

**Success Output:**

```
Successfully removed location 'my-project'
```

**Error Examples:**

```
Error: location with name "nonexistent" not found
Error: no deployment locations configured
```

### deploy

Deploy agents to a target project location.

**Flags:**
- `-a, --agents` (required): Comma-separated list of agent names
- `-l, --location` (required): Target project path
- `-o, --overwrite`: Overwrite existing files (default: false)
- `--output`: Output format: `text` or `json` (default: text)

**Examples:**

```bash
# Deploy two agents (text output)
./cami deploy --agents frontend,backend --location ~/projects/my-app

# Deploy with overwrite enabled
./cami deploy -a frontend,backend -l ~/projects/my-app --overwrite

# Deploy with JSON output for programmatic use
./cami deploy -a frontend,backend -l ~/projects/my-app --output json
```

**JSON Output Format:**

```json
{
  "success": true,
  "deployed": ["frontend", "backend"],
  "failed": [],
  "conflicts": [],
  "results": [
    {
      "agent": "frontend",
      "status": "success",
      "message": "Deployed successfully"
    },
    {
      "agent": "backend",
      "status": "success",
      "message": "Deployed successfully"
    }
  ]
}
```

**Status Values:**
- `success`: Agent deployed successfully
- `conflict`: File already exists (use --overwrite to replace)
- `failed`: Deployment failed (see message for details)

**Exit Codes:**
- `0`: All deployments successful
- `1`: One or more deployments failed or had conflicts

### update-docs

Update the CLAUDE.md file with deployed agent information.

**Flags:**
- `-l, --location` (required): Target project path
- `-s, --section`: Section name in CLAUDE.md (default: "Deployed Agents")
- `--dry-run`: Show changes without writing

**Examples:**

```bash
# Update CLAUDE.md with deployed agents
./cami update-docs --location ~/projects/my-app

# Use custom section name
./cami update-docs -l ~/projects/my-app --section "Available Agents"

# Preview changes without writing
./cami update-docs -l ~/projects/my-app --dry-run
```

**How It Works:**

1. Scans `.claude/agents/` in the target location
2. Extracts agent metadata (name, version, description)
3. Generates markdown section with agent documentation
4. Smart merge: Preserves existing CLAUDE.md content, only updates agent section
5. Creates CLAUDE.md if it doesn't exist

**Generated Section:**

The command creates a managed section marked with HTML comments:

```markdown
<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
## Deployed Agents

The following Claude Code agents are available in this project:

### frontend (v1.0.0)
Frontend development specialist for React, Vue, and modern web frameworks.

### backend (v1.0.0)
Backend development specialist for APIs, databases, and server-side logic.

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->
```

**Smart Merge:**

The command intelligently updates only the managed section:
- Content before the managed section is preserved
- Content after the managed section is preserved
- Only the agent list between markers is updated
- Safe to run multiple times

### list

List all available agents from the vc-agents directory.

**Flags:**
- `--output`: Output format: `text` or `json` (default: text)

**Examples:**

```bash
# List agents (text output)
./cami list

# List agents (JSON output)
./cami list --output json
```

**JSON Output Format:**

```json
{
  "count": 20,
  "agents": [
    {
      "name": "frontend",
      "version": "1.0.0",
      "description": "Frontend development specialist...",
      "file_path": "/path/to/vc-agents/frontend.md"
    },
    {
      "name": "backend",
      "version": "1.0.0",
      "description": "Backend development specialist...",
      "file_path": "/path/to/vc-agents/backend.md"
    }
  ]
}
```

### scan

Scan a project location and list deployed agents.

**Flags:**
- `-l, --location` (required): Target project path
- `--output`: Output format: `text` or `json` (default: text)

**Examples:**

```bash
# Scan deployed agents (text output)
./cami scan --location ~/projects/my-app

# Scan with JSON output
./cami scan -l ~/projects/my-app --output json
```

**JSON Output Format:**

```json
{
  "location": "/path/to/project",
  "count": 2,
  "agents": [
    {
      "name": "frontend",
      "version": "1.0.0",
      "description": "Frontend development specialist..."
    },
    {
      "name": "backend",
      "version": "1.0.0",
      "description": "Backend development specialist..."
    }
  ]
}
```

## Common Workflows

### Managing Deployment Locations

```bash
# Add a project location
./cami location add --name my-project --path ~/projects/my-app

# List all locations
./cami locations

# Deploy to a saved location (still requires full path)
./cami deploy --agents frontend,backend --location ~/projects/my-app

# Remove a location when project is archived
./cami location remove --name my-project
```

### Deploy and Document Agents

```bash
# 1. Deploy agents
./cami deploy --agents frontend,backend,qa --location ~/my-project

# 2. Update CLAUDE.md
./cami update-docs --location ~/my-project

# 3. Verify deployment
./cami scan --location ~/my-project
```

### Programmatic Deployment (with error handling)

```bash
#!/bin/bash

# Deploy agents and capture JSON output
RESULT=$(./cami deploy \
  --agents frontend,backend \
  --location ~/my-project \
  --output json)

# Check if deployment was successful
SUCCESS=$(echo "$RESULT" | jq -r '.success')

if [ "$SUCCESS" = "true" ]; then
  echo "Deployment successful!"
  ./cami update-docs --location ~/my-project
else
  echo "Deployment failed:"
  echo "$RESULT" | jq -r '.results[] | select(.status != "success") | "\(.agent): \(.message)"'
  exit 1
fi
```

### Preview Documentation Changes

```bash
# See what would be written to CLAUDE.md
./cami update-docs --location ~/my-project --dry-run
```

### List Available Agents for Deployment

```bash
# See all available agents
./cami list

# Get agent names in JSON for scripting
./cami list --output json | jq -r '.agents[].name'
```

## Integration with Claude Code

The CLI commands are designed to be used by Claude Code for programmatic agent management:

```bash
# Claude Code can deploy agents to the current project
cami deploy --agents frontend,backend --location $PWD --output json

# Then update the CLAUDE.md file
cami update-docs --location $PWD

# This creates a documented, versioned agent deployment
```

## Help

Get help for any command:

```bash
./cami --help
./cami deploy --help
./cami update-docs --help
./cami list --help
./cami scan --help
./cami locations --help
./cami location --help
./cami location add --help
./cami location remove --help
```

## Version

```bash
./cami --version
# Output: CAMI v0.2.0
#         Claude Agent Management Interface
```

## Backward Compatibility

The TUI mode is still fully functional - simply run `cami` without any subcommands to launch the interactive interface.
