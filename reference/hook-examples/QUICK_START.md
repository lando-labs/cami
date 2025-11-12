# CAMI Hooks - Quick Start

**TL;DR:** Get CAMI agent scanning working in Claude Code in 5 minutes.

## What You Get

When you open a project with Claude Code:
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

## 5-Minute Setup

### Step 1: Test the Hook

```bash
cd /Users/lando/Development/cami/reference/hook-examples
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh
```

You should see agent status output.

### Step 2: Install Globally

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

**Note:** Replace `/Users/lando/Development/cami` with your actual CAMI path.

### Step 3: Restart Claude Code

Close and reopen Claude Code. You'll see the agent status when you start a session.

## Done!

That's it. Now every time you open a project, you'll see your CAMI agent status automatically.

## What's Next?

### Add Auto-Update (Optional)

Auto-update CLAUDE.md after deploying agents:

```json
{
  "hooks": {
    "SessionStart": [...],
    "PostToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-post-tool.sh"
          }
        ]
      }
    ]
  }
}
```

### Add Validation (Optional)

Validate agents before deployment:

```json
{
  "hooks": {
    "SessionStart": [...],
    "PostToolUse": [...],
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-pre-tool-validate.sh"
          }
        ]
      }
    ]
  }
}
```

### Add Logging (Optional)

Track session activity:

```json
{
  "hooks": {
    "SessionStart": [...],
    "PostToolUse": [...],
    "PreToolUse": [...],
    "SessionEnd": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-session-end.sh"
          }
        ]
      }
    ]
  }
}
```

## Troubleshooting

### Hook doesn't run

Check permissions:
```bash
ls -l /Users/lando/Development/cami/reference/hook-examples/*.sh
# Should show: -rwxr-xr-x
```

Make executable if needed:
```bash
chmod +x /Users/lando/Development/cami/reference/hook-examples/*.sh
```

### "jq: command not found"

Install jq:
```bash
brew install jq
```

### No output visible

SessionStart hooks inject stdout into Claude's context. You won't see it in the terminal, but Claude will have it in context.

To verify it's working, check debug logs:
```bash
tail -f ~/.claude/debug/$(ls -t ~/.claude/debug/ | head -1)
```

## More Information

- **[README.md](./README.md)** - Full feature documentation
- **[IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md)** - Customization guide
- **[HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md)** - Complete hooks research

## Need Help?

1. Test the script manually first
2. Check Claude Code debug logs
3. Review the implementation guide
4. Open a GitHub issue

## Configuration Location

Your settings file is at:
- **Global:** `~/.claude/settings.json`
- **Project:** `.claude/settings.json`
- **Local:** `.claude/settings.local.json` (gitignored)

## Example Complete Configuration

Copy this entire block to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-session-start.sh",
            "description": "CAMI: Show agent status"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-post-tool.sh",
            "description": "CAMI: Auto-update docs"
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-pre-tool-validate.sh",
            "description": "CAMI: Validate deployments"
          }
        ]
      }
    ],
    "SessionEnd": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-session-end.sh",
            "description": "CAMI: Log session"
          }
        ]
      }
    ]
  }
}
```

**Remember:** Change the paths to match your CAMI installation!

---

Happy CAMI-ing! 🎉
