#!/bin/bash
#
# CAMI PreToolUse Hook
# Validates agent deployments before they happen
#
# This hook runs BEFORE tool execution.
# It can block the operation by returning exit code 2.

# Read stdin (tool use info)
TOOL_DATA=$(cat)
TOOL_NAME=$(echo "$TOOL_DATA" | jq -r '.tool_name // empty')
TOOL_INPUT=$(echo "$TOOL_DATA" | jq -r '.tool_input // empty')

# Exit if not a deploy_agents call
if [[ "$TOOL_NAME" != *"mcp__cami__deploy_agents"* ]]; then
  # Approve all other tools
  echo '{"decision":"approve"}'
  exit 0
fi

# Extract agent names from tool input
AGENT_NAMES=$(echo "$TOOL_INPUT" | jq -r '.agent_names[]?' 2>/dev/null)

if [ -z "$AGENT_NAMES" ]; then
  # No agents specified, approve
  echo '{"decision":"approve"}'
  exit 0
fi

# VC agents directory
VC_DIR="$HOME/Development/cami/vc-agents"

if [ ! -d "$VC_DIR" ]; then
  # VC directory not found, warn but approve
  echo '{"decision":"approve","systemMessage":"⚠️ CAMI VC directory not found - cannot validate versions"}'
  exit 0
fi

# Check each agent
MISSING=()
WARNINGS=()

while IFS= read -r agent; do
  # Skip empty lines
  [ -z "$agent" ] && continue

  agent_file="$VC_DIR/$agent.md"

  if [ ! -f "$agent_file" ]; then
    MISSING+=("$agent")
  else
    # Check if agent has version
    version=$(grep "^version:" "$agent_file" 2>/dev/null | head -1 | cut -d':' -f2 | tr -d ' ')
    if [ -z "$version" ]; then
      WARNINGS+=("$agent: no version found")
    fi
  fi
done <<< "$AGENT_NAMES"

# If there are missing agents, block with confirmation
if [ ${#MISSING[@]} -gt 0 ]; then
  missing_list=$(printf "\n  - %s" "${MISSING[@]}")

  # Return blocking JSON with reason
  cat <<EOF
{
  "decision": "ask",
  "reason": "⚠️  Some agents were not found in the VC repository:$missing_list\n\nThese agents may not exist or may be misspelled.\n\nContinue anyway?"
}
EOF
  exit 2
fi

# If there are warnings, show them but approve
if [ ${#WARNINGS[@]} -gt 0 ]; then
  warning_list=$(printf "\n  - %s" "${WARNINGS[@]}")

  cat <<EOF
{
  "decision": "approve",
  "systemMessage": "⚠️  Some warnings about agents:$warning_list"
}
EOF
  exit 0
fi

# All good, approve
echo '{"decision":"approve"}'
exit 0
