# CAMI Hooks Implementation Guide

This guide provides step-by-step instructions for implementing CAMI hooks in your Claude Code plugin.

## Overview

CAMI hooks provide automatic agent management across the Claude Code lifecycle:

1. **SessionStart** - Scan and display deployed agent status
2. **PreToolUse** - Validate agent deployments before they happen
3. **PostToolUse** - Auto-update CLAUDE.md after deployments
4. **SessionEnd** - Log session activity for audit trails

## Quick Start

### 1. Test Individual Hooks

Before integrating into your configuration, test each hook:

```bash
cd /Users/lando/Development/cami/reference/hook-examples

# Test SessionStart hook
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh

# Test PostToolUse hook
echo '{"tool_name":"mcp__cami__deploy_agents","cwd":"'$(pwd)'"}' | ./cami-post-tool.sh

# Test PreToolUse hook
echo '{"tool_name":"mcp__cami__deploy_agents","tool_input":{"agent_names":["architect","frontend"]}}' | ./cami-pre-tool-validate.sh

# Test SessionEnd hook
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-end.sh
```

### 2. Install Hooks (Project-Specific)

Add to `.claude/settings.json` in your project:

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

### 3. Install Hooks (Global)

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
    ],
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
    ],
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

## Plugin Integration

When creating a CAMI plugin, structure it like this:

```
cami-plugin/
├── .claude-plugin/
│   └── plugin.json
├── hooks/
│   └── hooks.json
├── scripts/
│   ├── cami-session-start.sh
│   ├── cami-pre-tool-validate.sh
│   ├── cami-post-tool.sh
│   └── cami-session-end.sh
└── README.md
```

### plugin.json

```json
{
  "name": "cami",
  "version": "1.0.0",
  "description": "Claude Agent Management Interface - Automated agent deployment and management",
  "author": {
    "name": "Your Name",
    "email": "[email protected]"
  },
  "repository": "https://github.com/yourusername/cami-plugin",
  "license": "MIT",
  "hooks": "./hooks/hooks.json",
  "mcpServers": {
    "cami": {
      "command": "node",
      "args": ["/path/to/cami/mcp/dist/index.js"],
      "env": {}
    }
  }
}
```

### hooks/hooks.json

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-session-start.sh"
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
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-pre-tool-validate.sh"
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
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-post-tool.sh"
          }
        ]
      }
    ],
    "SessionEnd": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-session-end.sh"
          }
        ]
      }
    ]
  }
}
```

## Customization Guide

### 1. Adjust VC Directory Path

If your CAMI installation is in a different location, update the scripts:

**In cami-session-start.sh:**
```bash
# Line 13 - Change this path
VC_DIR="$HOME/Development/cami/vc-agents"
```

**In cami-pre-tool-validate.sh:**
```bash
# Line 23 - Change this path
VC_DIR="$HOME/Development/cami/vc-agents"
```

### 2. Adjust CAMI CLI Path

**In cami-post-tool.sh:**
```bash
# Line 25 - Change this path
CAMI_CLI="$HOME/Development/cami/cami"
```

### 3. Customize Output Format

**Change SessionStart output:**

Edit `cami-session-start.sh` around lines 38-77 to customize the status display.

Example - Add emoji categories:
```bash
if [ "$deployed_version" = "$vc_version" ]; then
  echo "  ✅ $agent_name ($deployed_version)"
elif [ -n "$vc_version" ]; then
  echo "  🔄 $agent_name ($deployed_version → $vc_version)"
else
  echo "  ⚠️  $agent_name ($deployed_version)"
fi
```

### 4. Add Notifications

**macOS desktop notifications:**

Add to `cami-post-tool.sh` after line 36:
```bash
osascript -e 'display notification "CLAUDE.md updated" with title "CAMI"'
```

**Linux desktop notifications:**

```bash
notify-send "CAMI" "CLAUDE.md updated"
```

### 5. Add Logging

**Add detailed logging to any hook:**

```bash
# At the top of the script
LOG_FILE="$HOME/.cami/hook-debug.log"
echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $0 started" >> "$LOG_FILE"
echo "$TOOL_DATA" >> "$LOG_FILE"
```

### 6. Conditional Execution

**Only run in git repositories:**

```bash
if [ ! -d "$CWD/.git" ]; then
  exit 0
fi
```

**Only run for specific projects:**

```bash
# Check if project is in whitelist
ALLOWED_DIRS=("/path/to/project1" "/path/to/project2")

if [[ ! " ${ALLOWED_DIRS[@]} " =~ " ${CWD} " ]]; then
  exit 0
fi
```

**Only run during business hours:**

```bash
HOUR=$(date +%H)
if [ $HOUR -lt 9 ] || [ $HOUR -gt 17 ]; then
  exit 0
fi
```

## Advanced Features

### 1. Integration with Git

Auto-commit agent deployments:

```bash
# In cami-post-tool.sh after CLAUDE.md update
if [ -d "$CWD/.git" ]; then
  cd "$CWD"
  git add .claude/agents/ CLAUDE.md
  git commit -m "chore: update CAMI agents" --no-verify
fi
```

### 2. Slack Notifications

Send deployment notifications to Slack:

```bash
# In cami-post-tool.sh
WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"CAMI: Agents deployed in '"$CWD"'"}' \
  "$WEBHOOK_URL"
```

### 3. Metrics Collection

Track deployment metrics:

```bash
# In cami-post-tool.sh
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
METRICS_FILE="$HOME/.cami/metrics.jsonl"

echo "{\"timestamp\":\"$TIMESTAMP\",\"event\":\"deploy\",\"cwd\":\"$CWD\"}" >> "$METRICS_FILE"
```

### 4. Cost Tracking

Log session costs (requires SessionEnd enhancement):

```bash
# In cami-session-end.sh
SESSION_COST=$(echo "$SESSION_DATA" | jq -r '.session_cost // 0')

if [ "$SESSION_COST" != "0" ]; then
  echo "Session cost: \$$SESSION_COST" >> "$LOG_FILE"
fi
```

### 5. Multi-Language Support

Use Python for complex logic:

```python
#!/usr/bin/env python3
import json
import sys

# Read stdin
data = json.loads(sys.stdin.read())

# Complex validation logic
if validate_agents(data):
    print(json.dumps({"decision": "approve"}))
    sys.exit(0)
else:
    print(json.dumps({
        "decision": "block",
        "reason": "Validation failed"
    }))
    sys.exit(2)
```

## Troubleshooting

### Hook Not Running

**Check 1: Verify script permissions**
```bash
ls -l /Users/lando/Development/cami/reference/hook-examples/*.sh
# Should show: -rwxr-xr-x
```

**Check 2: Test script manually**
```bash
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh
```

**Check 3: Check Claude Code debug logs**
```bash
tail -f ~/.claude/debug/$(ls -t ~/.claude/debug/ | head -1)
```

### Hook Output Not Visible

- SessionStart and UserPromptSubmit: stdout → Claude's context
- Other hooks: stderr → visible in UI, stdout → logged but not shown

**To debug other hooks:**
```bash
# Redirect output to stderr
echo "Debug message" >&2
```

### JSON Parse Errors

**Validate JSON output:**
```bash
echo '{"decision":"approve"}' | jq .
```

**Common mistake - missing quotes:**
```bash
# Wrong
echo {decision:approve}

# Right
echo '{"decision":"approve"}'
```

### jq Command Not Found

**Install jq:**
```bash
# macOS
brew install jq

# Linux (Ubuntu/Debian)
apt-get install jq

# Linux (RHEL/CentOS)
yum install jq
```

### Timeout Issues

Hooks have a 60-second timeout. If your hook takes longer:

**Option 1: Run in background**
```bash
# Start long-running task in background
./long-task.sh &

# Return immediately
echo '{"decision":"approve"}'
exit 0
```

**Option 2: Optimize script**
```bash
# Use parallel processing
for file in *.txt; do
  process_file "$file" &
done
wait
```

### Permission Denied

**Make scripts executable:**
```bash
chmod +x /Users/lando/Development/cami/reference/hook-examples/*.sh
```

**Check directory permissions:**
```bash
ls -ld ~/.cami
# Should be writable: drwxr-xr-x
```

## Best Practices

### 1. Error Handling

Always include error handling:

```bash
#!/bin/bash
set -e  # Exit on error

# Trap errors
trap 'echo "Error on line $LINENO" >&2' ERR

# Validate inputs
if [ -z "$SESSION_DATA" ]; then
  echo "No session data provided" >&2
  exit 1
fi
```

### 2. Performance

Keep hooks fast:

```bash
# Bad - Sequential
for file in *.txt; do
  process_file "$file"
done

# Good - Parallel
for file in *.txt; do
  process_file "$file" &
done
wait
```

### 3. Security

Validate all inputs:

```bash
# Sanitize file paths
CWD=$(echo "$SESSION_DATA" | jq -r '.cwd // empty')
if [[ "$CWD" =~ [^a-zA-Z0-9/_-] ]]; then
  echo "Invalid CWD" >&2
  exit 1
fi
```

### 4. Logging

Log important events:

```bash
LOG_FILE="$HOME/.cami/hook.log"

log() {
  echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $*" >> "$LOG_FILE"
}

log "SessionStart hook executed in $CWD"
```

### 5. Testing

Test hooks in isolation:

```bash
# Test with mock data
test_session_start() {
  local TEST_DATA='{"session_id":"test","cwd":"'$(pwd)'"}'
  echo "$TEST_DATA" | ./cami-session-start.sh
}

test_session_start
```

## Next Steps

1. **Test hooks individually** - Verify each script works
2. **Install SessionStart hook** - Get immediate value
3. **Add PostToolUse hook** - Auto-update documentation
4. **Enable PreToolUse validation** - Prevent errors
5. **Add SessionEnd logging** - Track usage
6. **Create plugin package** - Distribute to others

## Resources

- [HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md) - Complete documentation
- [Official Hooks Guide](https://docs.claude.com/en/docs/claude-code/hooks-guide)
- [Official Hooks Reference](https://docs.claude.com/en/docs/claude-code/hooks)
- [Claude Code Plugins](https://docs.claude.com/en/docs/claude-code/plugins-reference)

## Support

For issues or questions:
- Open an issue on GitHub
- Check Claude Code documentation
- Review example repositories

## License

MIT License - Same as CAMI project
