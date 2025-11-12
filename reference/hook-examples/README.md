# CAMI Hook Examples

This directory contains working examples of Claude Code hooks for CAMI.

## Files

- **hooks.json** - Hook configuration file
- **cami-session-start.sh** - SessionStart hook (scans agents)
- **cami-post-tool.sh** - PostToolUse hook (updates CLAUDE.md)
- **cami-session-end.sh** - SessionEnd hook (logs sessions)

## Installation

### Option 1: Global Configuration

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-session-start.sh"
          }
        ]
      }
    ]
  }
}
```

### Option 2: Project Configuration

Add to `.claude/settings.json` in your project:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${HOME}/Development/cami/reference/hook-examples/cami-session-start.sh"
          }
        ]
      }
    ]
  }
}
```

### Option 3: Plugin (Future)

When CAMI becomes a full plugin, hooks will be automatically installed.

## Testing

Test a hook manually:

```bash
# Create test session data
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh
```

## Hook Behavior

### SessionStart Hook

**When it runs:** Every time you start or resume a Claude Code session

**What it does:**
- Scans `.claude/agents/` directory
- Compares deployed versions with VC repository
- Shows status for each agent
- Displays summary with update count

**Output example:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📦 CAMI Agent Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✅ architect (1.1.0) - up to date
  ⚠️  frontend (1.0.0 → 1.1.0) - update available
  ✅ backend (1.1.0) - up to date

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total: 3 agents
  ✅ Up to date: 2
  ⚠️  Updates available: 1

💡 Use CAMI MCP tools to update:
   - scan_deployed_agents: detailed status
   - deploy_agents: deploy updates

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### PostToolUse Hook

**When it runs:** After any tool execution

**What it does:**
- Filters for `mcp__cami__deploy_agents` calls
- Calls CAMI CLI to update CLAUDE.md
- Shows success/failure message

**Output example:**
```
🔄 CAMI: Updating CLAUDE.md after agent deployment...
✅ CLAUDE.md updated successfully
```

### SessionEnd Hook

**When it runs:** When you exit Claude Code

**What it does:**
- Logs session end timestamp
- Records session ID and working directory
- Maintains last 100 session logs

**Log location:** `~/.cami/session-log.jsonl`

## Customization

### Change VC Directory

Edit `cami-session-start.sh`:

```bash
# Change this line:
VC_DIR="$HOME/Development/cami/vc-agents"
```

### Change CAMI CLI Path

Edit `cami-post-tool.sh`:

```bash
# Change this line:
CAMI_CLI="$HOME/Development/cami/cami"
```

### Change Log Location

Edit `cami-session-end.sh`:

```bash
# Change this line:
LOG_DIR="$HOME/.cami"
```

## Troubleshooting

### Hook not running

1. Check that scripts are executable: `ls -l *.sh`
2. Test script manually: `echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh`
3. Check Claude Code logs: `~/.claude/debug/`

### No output in context

- SessionStart and UserPromptSubmit hooks inject stdout into context
- Other hooks only log to stderr visible in transcripts

### jq not found

Install jq: `brew install jq`

### Permission denied

Make scripts executable: `chmod +x *.sh`

## Security Notes

- Hooks run with your user credentials
- Always review hook code before installing
- Don't run hooks from untrusted sources
- Use absolute paths to prevent PATH injection

## Future Enhancements

When packaged as a CAMI plugin:

1. **Auto-discovery** - Detect CAMI installation automatically
2. **Configuration UI** - Manage hooks via `/hooks` command
3. **Update notifications** - Desktop/mobile notifications
4. **Smart filtering** - Only show relevant agent updates
5. **Performance tracking** - Measure deployment times
6. **Integration** - Direct MCP tool calls from hooks

## Related Documentation

- [HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md) - Complete hooks documentation
- [Official Claude Code Hooks Guide](https://docs.claude.com/en/docs/claude-code/hooks-guide)
- [Official Hooks Reference](https://docs.claude.com/en/docs/claude-code/hooks)

## License

MIT License - Same as CAMI project
