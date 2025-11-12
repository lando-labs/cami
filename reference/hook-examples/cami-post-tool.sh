#!/bin/bash
#
# CAMI PostToolUse Hook
# Auto-updates CLAUDE.md after agent deployment
#
# This hook runs after any tool use.
# We filter for MCP deploy_agents calls.

# Read stdin (tool use info)
TOOL_DATA=$(cat)
TOOL_NAME=$(echo "$TOOL_DATA" | jq -r '.tool_name // empty')
CWD=$(echo "$TOOL_DATA" | jq -r '.cwd // empty')

# Exit if no tool name or CWD
if [ -z "$TOOL_NAME" ] || [ -z "$CWD" ]; then
  exit 0
fi

# Only run for CAMI deploy_agents MCP calls
if [[ "$TOOL_NAME" != *"mcp__cami__deploy_agents"* ]]; then
  exit 0
fi

# Check if CLAUDE.md exists
CLAUDE_MD="$CWD/CLAUDE.md"
if [ ! -f "$CLAUDE_MD" ]; then
  exit 0
fi

echo ""
echo "🔄 CAMI: Updating CLAUDE.md after agent deployment..."

# Call CAMI CLI to update docs
CAMI_CLI="$HOME/Development/cami/cami"

if [ ! -x "$CAMI_CLI" ]; then
  echo "⚠️  CAMI CLI not found at $CAMI_CLI"
  exit 0
fi

# Update documentation
if "$CAMI_CLI" update-docs "$CWD" 2>/dev/null; then
  echo "✅ CLAUDE.md updated successfully"
else
  echo "⚠️  Failed to update CLAUDE.md"
fi

echo ""

exit 0
