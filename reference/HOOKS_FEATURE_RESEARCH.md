# Claude Code Hooks Feature - Complete Research Report

**Date:** October 9, 2025
**Purpose:** Comprehensive research on Claude Code hooks for CAMI plugin development

---

## Table of Contents

1. [What are Hooks?](#what-are-hooks)
2. [Lifecycle Events](#lifecycle-events)
3. [Hook Actions & Capabilities](#hook-actions--capabilities)
4. [Configuration](#configuration)
5. [Technical Specification](#technical-specification)
6. [Use Cases for CAMI](#use-cases-for-cami)
7. [Implementation Recommendations](#implementation-recommendations)
8. [Working Examples](#working-examples)

---

## What are Hooks?

**Claude Code Hooks** are user-defined shell commands that execute automatically at specific points in Claude Code's lifecycle. They provide **deterministic control** over Claude Code's behavior, ensuring certain actions always happen rather than relying on the LLM to choose to run them.

### Key Characteristics

- **Automatic execution** - Hooks run automatically with the current environment's credentials
- **Event-driven** - Triggered by specific lifecycle events (tool use, prompts, sessions, etc.)
- **Shell-based** - Execute any shell command, script, or program
- **Input/Output** - Receive JSON via stdin, return JSON via stdout for control flow
- **Timeout** - Default 60-second timeout per hook execution
- **Parallel execution** - Multiple hooks for the same event can run concurrently

### Security Considerations

⚠️ **Important Security Notes:**
- Hooks run with current environment credentials
- Can execute arbitrary shell commands
- Potential for data exfiltration if not carefully implemented
- Always review hook implementations before registering
- Use absolute paths and validate all inputs

---

## Lifecycle Events

Claude Code provides **8 distinct lifecycle events** where hooks can be triggered:

### 1. UserPromptSubmit

**When it fires:** Immediately when user submits a prompt, before Claude processes it

**Special behavior:** stdout from this hook is added to Claude's context

**Use cases:**
- Logging user prompts
- Prompt validation/filtering
- Context injection (add dynamic information to prompts)
- Security filtering (block dangerous prompts)

**Can block:** Yes (exit code 2)

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "UserPromptSubmit",
  "user_prompt": "The actual prompt text"
}
```

---

### 2. PreToolUse

**When it fires:** Before any tool execution (Bash, Write, Edit, Read, etc.)

**Special behavior:** Can approve, block, or ask for permission

**Use cases:**
- Security validation (block dangerous commands)
- Permission management
- Logging tool usage
- File protection (prevent writes to sensitive files)
- Custom approval logic

**Can block:** Yes (exit code 2 or JSON decision: "block")

**Matchers:** Supports tool-specific patterns (e.g., "Write|Edit", "Bash", "Read")

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "PreToolUse",
  "tool_name": "Bash",
  "tool_input": {
    "command": "rm -rf /",
    "description": "Remove all files"
  }
}
```

**JSON Output (control flow):**
```json
{
  "decision": "block",
  "reason": "Dangerous command detected: rm -rf"
}
```

---

### 3. PostToolUse

**When it fires:** After successful tool execution

**Use cases:**
- Automatic code formatting (after Write/Edit)
- Logging successful operations
- Validation of results
- Triggering follow-up actions
- Git commits after file changes

**Can block:** Yes (but tool already executed, can only block continuation)

**Matchers:** Same tool-specific patterns as PreToolUse

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "PostToolUse",
  "tool_name": "Write",
  "tool_input": {
    "file_path": "/src/index.js",
    "content": "// Code here"
  },
  "tool_response": {
    "success": true,
    "message": "File written successfully"
  }
}
```

---

### 4. Stop

**When it fires:** When Claude Code finishes responding (main agent)

**Use cases:**
- Desktop notifications ("Claude is done")
- Continuation enforcement (prevent premature stopping)
- Session summaries
- Cost tracking

**Can block:** Yes (prevent Claude from stopping)

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "Stop",
  "stop_reason": "completed"
}
```

**JSON Output (force continuation):**
```json
{
  "continue": false,
  "decision": "block",
  "reason": "Please complete the remaining tasks: test the changes"
}
```

---

### 5. SubagentStop

**When it fires:** When a subagent (delegated agent) finishes responding

**Use cases:**
- Subagent task tracking
- Multi-agent workflow management
- Logging subagent completion

**Can block:** Yes

**JSON Input:** Similar to Stop event but for subagents

---

### 6. Notification

**When it fires:** During tool permission requests or input idle periods

**Use cases:**
- Custom notification systems
- Desktop alerts
- Mobile push notifications
- Slack/Discord integration

**Can block:** No

**JSON Input:**
```json
{
  "session_id": "abc123",
  "hook_event_name": "Notification",
  "notification_type": "permission_request",
  "message": "Claude is requesting permission to run command"
}
```

---

### 7. SessionStart

**When it fires:** When Claude Code starts a new session or resumes an existing session

**Special behavior:** stdout is added to Claude's context (like UserPromptSubmit)

**Use cases:**
- Environment setup
- Dependency installation (npm install, pip install)
- Loading project-specific context
- Memory/context injection
- Git worktree initialization
- **Auto-scan deployed agents** ✅

**Can block:** No

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "SessionStart",
  "timestamp": "2025-10-09T23:15:00Z"
}
```

---

### 8. SessionEnd

**When it fires:** On clean session termination (e.g., via /exit or Ctrl+D)

**Use cases:**
- Session cleanup
- Transcript archiving
- Cost/usage logging
- Backup operations
- Sending session summary to observability platform

**Can block:** No

**JSON Input:**
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/conversation.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "SessionEnd",
  "timestamp": "2025-10-09T23:45:00Z",
  "session_cost": 0.0234
}
```

---

### 9. PreCompact

**When it fires:** Before Claude Code compacts conversation history (to save tokens)

**Use cases:**
- Archiving full conversation before compaction
- Logging compaction events
- Custom compaction strategies

**Can block:** No

---

## Hook Actions & Capabilities

### What Hooks Can Do

1. **Execute Commands/Scripts**
   - Run any shell command
   - Execute scripts in any language (Python, Bash, Node.js, etc.)
   - Use environment variables

2. **Call MCP Tools**
   - Hooks can invoke MCP server tools
   - Integration with CAMI MCP server possible
   - Example: Auto-deploy agents via MCP

3. **Modify Behavior**
   - Block actions (exit code 2)
   - Approve actions automatically (bypass prompts)
   - Force continuation (prevent stopping)
   - Inject context into Claude's prompt

4. **Logging/Telemetry**
   - Log all tool usage
   - Track session costs
   - Monitor agent activity
   - Send data to observability platforms

5. **Validation**
   - Validate code quality
   - Check security constraints
   - Enforce coding standards
   - Prevent dangerous operations

6. **Automation**
   - Auto-format code
   - Auto-commit changes
   - Auto-deploy on completion
   - Auto-test after changes

---

## Configuration

### Configuration Locations

Hooks can be configured in multiple places (priority order):

1. **Project-local settings** (highest priority)
   - `.claude/settings.local.json` (gitignored)

2. **Project settings**
   - `.claude/settings.json` (checked into git)

3. **Plugin hooks**
   - `hooks/hooks.json` in plugin root
   - Referenced in `plugin.json`

4. **User global settings** (lowest priority)
   - `~/.claude/settings.json`

### File Format: hooks.json

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "ToolPattern",
        "hooks": [
          {
            "type": "command",
            "command": "script-to-execute",
            "description": "Optional description"
          }
        ]
      }
    ]
  }
}
```

### Configuration Examples

#### Basic Hook (SessionStart)

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/setup.sh"
          }
        ]
      }
    ]
  }
}
```

#### Hook with Matcher (PostToolUse)

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "npx prettier --write \"$file_path\""
          }
        ]
      }
    ]
  }
}
```

#### Multiple Hooks for Same Event

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/format.sh"
          },
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/lint.sh"
          }
        ]
      }
    ]
  }
}
```

### Matcher Patterns

Matchers use **regular expression syntax** to filter which tool uses trigger the hook:

- `"Write"` - Matches only Write tool
- `"Write|Edit"` - Matches Write OR Edit
- `".*"` or `""` - Matches all tools
- `"Bash"` - Matches only Bash commands
- `"Read|Glob|Grep"` - Matches any read operations

**Events that support matchers:**
- PreToolUse
- PostToolUse

**Events that don't use matchers:**
- UserPromptSubmit
- Stop
- SubagentStop
- SessionStart
- SessionEnd
- PreCompact
- Notification

### Plugin Integration

In `plugin.json`:

```json
{
  "name": "cami-plugin",
  "version": "1.0.0",
  "hooks": "./hooks/hooks.json"
}
```

Or inline in `plugin.json`:

```json
{
  "name": "cami-plugin",
  "version": "1.0.0",
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/scan-agents.sh"
          }
        ]
      }
    ]
  }
}
```

### Conditional Hooks

Hooks can include conditional logic in the script itself:

```bash
#!/bin/bash
# Hook script with conditions

# Only run in git repositories
if [ -d .git ]; then
  echo "Running in git repo"
  git status
fi

# Only run for TypeScript files
if [[ "$file_path" == *.ts ]]; then
  npx prettier --write "$file_path"
fi
```

---

## Technical Specification

### Exit Codes

- **Exit Code 0:** Success - hook ran without issues, execution continues normally
- **Exit Code 2:** Blocking error - tells Claude Code to halt the current action and process feedback
- **Exit Code 1 (or others):** General error - logs error but doesn't block

### Control Flow Priority

The layered system has clear priority:

1. **`continue: false`** - Overrides everything, completely stops Claude
2. **`decision: "block"`** - Blocks the specific action
3. **Exit code 2** - Simplest blocking mechanism

### JSON Input/Output

#### Common Input Fields (all events)

```json
{
  "session_id": "string",
  "transcript_path": "string",
  "cwd": "string",
  "hook_event_name": "string"
}
```

#### Event-Specific Input Fields

**PreToolUse / PostToolUse:**
```json
{
  "tool_name": "string",
  "tool_input": { /* varies by tool */ },
  "tool_response": { /* PostToolUse only */ }
}
```

**UserPromptSubmit:**
```json
{
  "user_prompt": "string"
}
```

**Stop / SubagentStop:**
```json
{
  "stop_reason": "string"
}
```

#### JSON Output Format (stdout)

```json
{
  "continue": true,          // false = stop all processing
  "decision": "approve",     // "approve" | "block" | "ask"
  "reason": "string",        // Required when decision = "block"
  "stopReason": "string",    // Message when continue = false
  "suppressOutput": false,   // Hide stdout from transcript
  "systemMessage": "string"  // Warning message to display
}
```

### Environment Variables

Available in hook scripts:

- `${CLAUDE_PLUGIN_ROOT}` - Plugin root directory
- `${CLAUDE_SESSION_ID}` - Current session ID
- `${CLAUDE_CWD}` - Current working directory
- `$file_path` - File path (for Write/Edit tools)
- Standard shell environment variables

### Timeouts

- Default: 60 seconds per hook
- Configurable per hook (future feature)
- Hook will be killed if it exceeds timeout

---

## Use Cases for CAMI

### 1. Auto-Scan Deployed Agents (SessionStart) ✅

**Purpose:** Automatically scan for deployed agents when a project is opened

**Implementation:**

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/scan-agents.sh",
            "description": "Scan deployed agents and inject status into context"
          }
        ]
      }
    ]
  }
}
```

**Script:** `scripts/scan-agents.sh`

```bash
#!/bin/bash

# Read stdin (session info)
SESSION_DATA=$(cat)
CWD=$(echo "$SESSION_DATA" | jq -r '.cwd')

# Check if .claude/agents directory exists
if [ ! -d "$CWD/.claude/agents" ]; then
  exit 0
fi

# Scan for deployed agents
echo "📦 CAMI Agent Status:"
echo ""
echo "Deployed agents in this project:"

for agent_file in "$CWD/.claude/agents"/*.md; do
  if [ -f "$agent_file" ]; then
    agent_name=$(basename "$agent_file" .md)

    # Extract version from frontmatter
    version=$(grep "^version:" "$agent_file" | cut -d':' -f2 | tr -d ' ')

    # Check against VC repository
    vc_version=$(grep "^version:" "$HOME/Development/cami/vc-agents/$agent_name.md" 2>/dev/null | cut -d':' -f2 | tr -d ' ')

    if [ "$version" = "$vc_version" ]; then
      echo "  ✅ $agent_name ($version) - up to date"
    else
      echo "  ⚠️  $agent_name ($version) - update available ($vc_version)"
    fi
  fi
done

echo ""
echo "Use CAMI MCP tools to update agents: deploy_agents, scan_deployed_agents"
```

**Output:** Injected into Claude's context at session start

---

### 2. Auto-Update CLAUDE.md (PostToolUse) ✅

**Purpose:** Automatically update CLAUDE.md after deploying agents

**Implementation:**

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/update-claude-md.sh",
            "description": "Auto-update CLAUDE.md after agent deployment"
          }
        ]
      }
    ]
  }
}
```

**Script:** `scripts/update-claude-md.sh`

```bash
#!/bin/bash

# Read stdin (tool use info)
TOOL_DATA=$(cat)
TOOL_NAME=$(echo "$TOOL_DATA" | jq -r '.tool_name')
CWD=$(echo "$TOOL_DATA" | jq -r '.cwd')

# Only run for MCP tool calls to deploy_agents
if [[ "$TOOL_NAME" == *"mcp__cami__deploy_agents"* ]]; then
  echo "Updating CLAUDE.md after agent deployment..."

  # Call CAMI MCP to update CLAUDE.md
  # This would need to be implemented via MCP protocol
  # For now, could call the CAMI CLI directly
  "$HOME/Development/cami/cami" update-docs "$CWD"

  echo "✅ CLAUDE.md updated"
fi
```

---

### 3. Log Agent Deployments (PostToolUse) ✅

**Purpose:** Track all agent deployments for audit/telemetry

**Implementation:**

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/log-deployment.sh",
            "description": "Log agent deployments"
          }
        ]
      }
    ]
  }
}
```

**Script:** `scripts/log-deployment.sh`

```bash
#!/bin/bash

# Read stdin
TOOL_DATA=$(cat)
TOOL_NAME=$(echo "$TOOL_DATA" | jq -r '.tool_name')
TOOL_INPUT=$(echo "$TOOL_DATA" | jq -r '.tool_input')

# Only log deploy_agents calls
if [[ "$TOOL_NAME" == *"deploy_agents"* ]]; then
  TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  SESSION_ID=$(echo "$TOOL_DATA" | jq -r '.session_id')

  # Extract agent names
  AGENTS=$(echo "$TOOL_INPUT" | jq -r '.agent_names[]' | tr '\n' ',' | sed 's/,$//')

  # Log to file
  LOG_FILE="$HOME/.cami/deployment-log.jsonl"
  echo "{\"timestamp\":\"$TIMESTAMP\",\"session\":\"$SESSION_ID\",\"agents\":\"$AGENTS\"}" >> "$LOG_FILE"

  echo "✅ Deployment logged"
fi
```

---

### 4. Validate Agent Versions (PreToolUse) ✅

**Purpose:** Warn before deploying outdated agent versions

**Implementation:**

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/validate-versions.sh",
            "description": "Validate agent versions before deployment"
          }
        ]
      }
    ]
  }
}
```

**Script:** `scripts/validate-versions.sh`

```bash
#!/bin/bash

# Read stdin
TOOL_DATA=$(cat)
TOOL_NAME=$(echo "$TOOL_DATA" | jq -r '.tool_name')

# Only validate deploy_agents calls
if [[ "$TOOL_NAME" == *"deploy_agents"* ]]; then
  TOOL_INPUT=$(echo "$TOOL_DATA" | jq -r '.tool_input')
  AGENTS=$(echo "$TOOL_INPUT" | jq -r '.agent_names[]')

  # Check each agent
  VC_DIR="$HOME/Development/cami/vc-agents"
  OUTDATED=""

  for agent in $AGENTS; do
    if [ ! -f "$VC_DIR/$agent.md" ]; then
      OUTDATED="$OUTDATED\n  - $agent: NOT FOUND in VC repository"
    fi
  done

  if [ -n "$OUTDATED" ]; then
    # Return blocking JSON
    echo "{\"decision\":\"ask\",\"reason\":\"Some agents may be outdated or missing:$OUTDATED\n\nContinue anyway?\"}"
    exit 2
  fi
fi

# Approve
echo '{"decision":"approve"}'
```

---

### 5. Notify on Outdated Agents (SessionStart) ✅

**Purpose:** Alert user if deployed agents are outdated

**Implementation:**

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/check-updates.sh",
            "description": "Check for agent updates"
          }
        ]
      }
    ]
  }
}
```

**Script:** `scripts/check-updates.sh`

```bash
#!/bin/bash

# Read stdin
SESSION_DATA=$(cat)
CWD=$(echo "$SESSION_DATA" | jq -r '.cwd')

# Check for updates
if [ -d "$CWD/.claude/agents" ]; then
  UPDATES_AVAILABLE=0

  for agent_file in "$CWD/.claude/agents"/*.md; do
    if [ -f "$agent_file" ]; then
      agent_name=$(basename "$agent_file" .md)
      deployed_version=$(grep "^version:" "$agent_file" | cut -d':' -f2 | tr -d ' ')
      vc_version=$(grep "^version:" "$HOME/Development/cami/vc-agents/$agent_name.md" 2>/dev/null | cut -d':' -f2 | tr -d ' ')

      if [ "$deployed_version" != "$vc_version" ] && [ -n "$vc_version" ]; then
        UPDATES_AVAILABLE=$((UPDATES_AVAILABLE + 1))
      fi
    fi
  done

  if [ $UPDATES_AVAILABLE -gt 0 ]; then
    echo ""
    echo "🔔 CAMI Update Available"
    echo "━━━━━━━━━━━━━━━━━━━━━━"
    echo "$UPDATES_AVAILABLE agent(s) have updates available."
    echo "Run 'scan_deployed_agents' to see details."
    echo ""
  fi
fi
```

---

### 6. Auto-Git Commit After Deployment (PostToolUse)

**Purpose:** Automatically commit agent changes to git

**Implementation:**

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/auto-commit.sh",
            "description": "Auto-commit agent deployments"
          }
        ]
      }
    ]
  }
}
```

---

## Implementation Recommendations

### Priority 1: Core Hooks

1. **SessionStart - Agent Scanner**
   - Scan deployed agents
   - Show update status
   - Inject into context

2. **PostToolUse - Auto-update CLAUDE.md**
   - Trigger after deploy_agents
   - Keep documentation in sync

### Priority 2: Enhanced Features

3. **SessionStart - Update Notifier**
   - Alert on outdated agents
   - Provide update instructions

4. **PostToolUse - Deployment Logger**
   - Track all deployments
   - Audit trail

### Priority 3: Advanced Features

5. **PreToolUse - Version Validator**
   - Warn before deploying old versions
   - Interactive approval

6. **SessionEnd - Summary Reporter**
   - Show deployment summary
   - Cost tracking

### Plugin Structure

```
cami-plugin/
├── .claude-plugin/
│   └── plugin.json
├── hooks/
│   └── hooks.json
├── scripts/
│   ├── scan-agents.sh
│   ├── update-claude-md.sh
│   ├── log-deployment.sh
│   ├── validate-versions.sh
│   ├── check-updates.sh
│   └── auto-commit.sh
└── README.md
```

### Best Practices

1. **Always validate input** from stdin
2. **Use absolute paths** (avoid relative paths)
3. **Handle errors gracefully** (exit codes)
4. **Keep hooks fast** (< 5 seconds ideal)
5. **Use environment variables** (`${CLAUDE_PLUGIN_ROOT}`)
6. **Test thoroughly** before deployment
7. **Document behavior** in hook descriptions
8. **Consider security** (never trust external input)

---

## Working Examples

### Example 1: Simple SessionStart Hook

**File:** `.claude/settings.json`

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo '📦 CAMI Plugin Active - Use deploy_agents to manage Claude Code agents'"
          }
        ]
      }
    ]
  }
}
```

---

### Example 2: Code Formatter (PostToolUse)

**File:** `hooks/hooks.json`

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "npx prettier --write \"$file_path\""
          }
        ]
      }
    ]
  }
}
```

---

### Example 3: Security Blocker (PreToolUse)

**File:** `scripts/security-check.py`

```python
#!/usr/bin/env python3
import json
import sys
import re

# Read stdin
data = json.load(sys.stdin)

if data.get('tool_name') == 'Bash':
    command = data.get('tool_input', {}).get('command', '')

    # Dangerous patterns
    dangerous = [
        r'rm\s+.*-[rf]',  # rm -rf
        r'sudo\s+rm',     # sudo rm
        r'chmod\s+777',   # chmod 777
    ]

    for pattern in dangerous:
        if re.search(pattern, command):
            # Block the command
            result = {
                "decision": "block",
                "reason": f"Dangerous command detected: {command}"
            }
            print(json.dumps(result))
            sys.exit(2)

# Approve
print(json.dumps({"decision": "approve"}))
sys.exit(0)
```

**Configuration:**

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/security-check.py"
          }
        ]
      }
    ]
  }
}
```

---

### Example 4: Complete CAMI Plugin Hook

**File:** `hooks/hooks.json`

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-session-start.sh",
            "description": "CAMI: Scan deployed agents and show status"
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
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-post-deploy.sh",
            "description": "CAMI: Auto-update CLAUDE.md after agent deployment"
          }
        ]
      }
    ],
    "SessionEnd": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/scripts/cami-session-end.sh",
            "description": "CAMI: Log session summary"
          }
        ]
      }
    ]
  }
}
```

---

## Resources

### Official Documentation
- [Claude Code Hooks Guide](https://docs.claude.com/en/docs/claude-code/hooks-guide)
- [Claude Code Hooks Reference](https://docs.claude.com/en/docs/claude-code/hooks)
- [Claude Code Plugins Reference](https://docs.claude.com/en/docs/claude-code/plugins-reference)
- [Anthropic Plugin Announcement](https://www.anthropic.com/news/claude-code-plugins)

### Community Examples
- [claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery) - Comprehensive examples
- [claude-hooks (johnlindquist)](https://github.com/johnlindquist/claude-hooks) - TypeScript hooks
- [claude-hooks (decider)](https://github.com/decider/claude-hooks) - Code quality validation
- [claude-git](https://github.com/listfold/claude-git) - Git worktree integration
- [GitButler Blog](https://blog.gitbutler.com/automate-your-ai-workflows-with-claude-code-hooks)

### Articles
- [Automate Your AI Workflows with Claude Code Hooks](https://blog.gitbutler.com/automate-your-ai-workflows-with-claude-code-hooks)
- [How I'm Using Claude Code Hooks](https://medium.com/@joe.njenga/use-claude-code-hooks-newest-feature-to-fully-automate-your-workflow-341b9400cfbe)

---

## Conclusion

Claude Code hooks provide a powerful, flexible mechanism for extending and automating Claude Code behavior. For CAMI, hooks offer perfect integration points to:

1. **Auto-scan agents** on project open (SessionStart)
2. **Auto-update documentation** after deployments (PostToolUse)
3. **Track deployments** for audit (PostToolUse)
4. **Validate versions** before deployment (PreToolUse)
5. **Notify on updates** at session start (SessionStart)

The hooks system is mature, well-documented, and production-ready. Implementation should be straightforward using shell scripts that integrate with the existing CAMI CLI and MCP server.

**Next Steps:**
1. Create `cami-plugin` directory structure
2. Implement core hooks (SessionStart scanner, PostToolUse updater)
3. Test hooks in development environment
4. Package as installable plugin
5. Document plugin usage

---

**Research completed:** October 9, 2025
**Compiled by:** Claude Code (Sonnet 4.5)
