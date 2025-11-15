# CAMI Quick Start

Get started with CAMI in 2 minutes.

## Step 1: Clone and Open

```bash
git clone <cami-repo-url>
cd cami
```

Open this directory in Claude Code.

## Step 2: Start Using CAMI

That's it! CAMI is automatically configured as an MCP server for this project.

Try it out:

```
You: "Help me get started with CAMI"
Claude: *uses mcp__cami__onboard*
```

## How It Works

The `.claude/settings.local.json` file in this repo configures CAMI to run via `go run` whenever you use Claude Code in this directory. No build step, no global installation needed.

## Next Steps

1. **First-time setup**: Ask Claude "Help me get started with CAMI"
2. **Add agent sources**: Claude will guide you to add your agent repositories
3. **Deploy agents**: Use natural language - "Add the frontend agent to this project"

See [README.md](README.md) for complete documentation.

## Requirements

- Go 1.21+ installed
- Claude Code

That's it!
