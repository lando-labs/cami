#!/bin/bash
#
# CAMI SessionStart Hook
# Scans deployed agents and shows status at session start
#
# This hook runs when Claude Code starts a session.
# stdout is injected into Claude's context.

# Read stdin (session info)
SESSION_DATA=$(cat)
CWD=$(echo "$SESSION_DATA" | jq -r '.cwd // empty')

# Exit if no CWD
if [ -z "$CWD" ]; then
  exit 0
fi

# Check if .claude/agents directory exists
if [ ! -d "$CWD/.claude/agents" ]; then
  exit 0
fi

# VC agents directory (adjust path as needed)
VC_DIR="$HOME/Development/cami/vc-agents"

if [ ! -d "$VC_DIR" ]; then
  echo "⚠️  CAMI: VC agents directory not found at $VC_DIR"
  exit 0
fi

# Output header
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 CAMI Agent Status"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Scan for deployed agents
AGENT_COUNT=0
UP_TO_DATE=0
UPDATES_AVAILABLE=0
UNKNOWN=0

for agent_file in "$CWD/.claude/agents"/*.md; do
  if [ -f "$agent_file" ]; then
    agent_name=$(basename "$agent_file" .md)
    AGENT_COUNT=$((AGENT_COUNT + 1))

    # Extract version from frontmatter (YAML-style)
    deployed_version=$(grep "^version:" "$agent_file" 2>/dev/null | head -1 | cut -d':' -f2 | tr -d ' ')

    # Check against VC repository
    vc_file="$VC_DIR/$agent_name.md"
    if [ -f "$vc_file" ]; then
      vc_version=$(grep "^version:" "$vc_file" 2>/dev/null | head -1 | cut -d':' -f2 | tr -d ' ')

      if [ "$deployed_version" = "$vc_version" ] && [ -n "$deployed_version" ]; then
        echo "  ✅ $agent_name ($deployed_version) - up to date"
        UP_TO_DATE=$((UP_TO_DATE + 1))
      elif [ -n "$vc_version" ]; then
        echo "  ⚠️  $agent_name ($deployed_version → $vc_version) - update available"
        UPDATES_AVAILABLE=$((UPDATES_AVAILABLE + 1))
      else
        echo "  ❓ $agent_name ($deployed_version) - version unknown"
        UNKNOWN=$((UNKNOWN + 1))
      fi
    else
      echo "  ❓ $agent_name ($deployed_version) - not found in VC"
      UNKNOWN=$((UNKNOWN + 1))
    fi
  fi
done

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Total: $AGENT_COUNT agents"
echo "  ✅ Up to date: $UP_TO_DATE"
if [ $UPDATES_AVAILABLE -gt 0 ]; then
  echo "  ⚠️  Updates available: $UPDATES_AVAILABLE"
fi
if [ $UNKNOWN -gt 0 ]; then
  echo "  ❓ Unknown: $UNKNOWN"
fi
echo ""

if [ $UPDATES_AVAILABLE -gt 0 ]; then
  echo "💡 Use CAMI MCP tools to update:"
  echo "   - scan_deployed_agents: detailed status"
  echo "   - deploy_agents: deploy updates"
  echo ""
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

exit 0
