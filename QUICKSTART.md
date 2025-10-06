# CAMI Quick Start Guide

Get up and running with CAMI in 5 minutes.

## Step 1: Build CAMI

```bash
# Using Make
make build

# Or using Go directly
go build -o cami cmd/cami/main.go
```

## Step 2: Run CAMI

```bash
./cami
```

## Step 3: Add a Deployment Location

1. When CAMI launches, press `l` to open location management
2. Press `a` to add a new location
3. Enter a name (e.g., "my-project")
4. Press `tab` to switch to the path field
5. Enter the full path to your project (e.g., "/Users/username/projects/my-project")
6. Press `enter` to save
7. Press `esc` to return to agent selection

## Step 4: Deploy Agents

1. Use `↑/↓` or `j/k` to navigate the agent list
2. Press `space` to select agents you want to deploy
3. Press `enter` when ready to deploy
4. Select your target location with `↑/↓`
5. Press `enter` to deploy

## Step 5: Verify Deployment

Check your project directory:

```bash
ls -la /path/to/your/project/.claude/agents/
```

You should see your deployed agent files!

## Optional: MCP Server Setup

Enable CAMI in Claude Code for seamless agent management from any project.

### Build MCP Server

```bash
go build -o cami-mcp cmd/cami-mcp/main.go
```

### Configure Claude Code

Add to `~/.config/claude-code/mcp_settings.json`:

```json
{
  "mcpServers": {
    "cami": {
      "command": "/Users/lando/Development/cami/cami-mcp",
      "env": {
        "CAMI_VC_AGENTS_DIR": "/Users/lando/Development/cami/vc-agents"
      }
    }
  }
}
```

**Note**: Use absolute paths! Replace `/Users/lando/Development/cami` with your actual CAMI location.

### Restart Claude Code

Reload VSCode window: `Cmd+Shift+P` → "Developer: Reload Window"

### Use MCP Tools

Now you can manage agents directly in Claude Code conversations:

```
"What agents are available in CAMI?"
"Deploy the architect agent to /Users/username/my-project"
"Add my-project as a CAMI location at /Users/username/projects/my-app"
```

See [MCP_README.md](MCP_README.md) for complete MCP documentation.

## Common Commands

```bash
# Show version
./cami --version

# Show help
./cami --help

# Install to PATH (optional)
make install
```

## Tips

- **Keyboard-first**: Everything can be done without a mouse
- **Multi-select**: Select multiple agents before deploying
- **Conflict detection**: CAMI warns you if files already exist
- **Persistent config**: Your locations are saved in `~/.cami.json`

## Keyboard Reference

### Agent Selection
- `↑/k` - Move up
- `↓/j` - Move down
- `space` - Select agent
- `enter` - Deploy
- `l` - Locations
- `q` - Quit

### Location Management
- `a` - Add location
- `d` - Delete location
- `tab` - Switch fields (when adding)
- `esc` - Back

## What's Deployed?

When you deploy agents, CAMI:
1. Creates `.claude/agents/` directory in your project
2. Copies agent markdown files with full YAML frontmatter
3. Preserves versions and metadata
4. Makes agents available to Claude Code in that project

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check [reference/claude-agent-manager-interface.md](reference/claude-agent-manager-interface.md) for requirements
- Explore the agent files in `vc-agents/` to understand the structure

## Troubleshooting

**Problem**: "vc-agents directory not found"
- **Solution**: Run CAMI from the cami project directory

**Problem**: "path does not exist" when adding location
- **Solution**: Create the project directory first, or use an existing path

**Problem**: No agents showing up
- **Solution**: Ensure agent files in `vc-agents/` have proper YAML frontmatter

---

Happy deploying!
