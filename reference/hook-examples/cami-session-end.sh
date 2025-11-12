#!/bin/bash
#
# CAMI SessionEnd Hook
# Logs session summary
#
# This hook runs when Claude Code session ends.

# Read stdin (session info)
SESSION_DATA=$(cat)
SESSION_ID=$(echo "$SESSION_DATA" | jq -r '.session_id // "unknown"')
CWD=$(echo "$SESSION_DATA" | jq -r '.cwd // "unknown"')
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Log directory
LOG_DIR="$HOME/.cami"
LOG_FILE="$LOG_DIR/session-log.jsonl"

# Create log directory if it doesn't exist
mkdir -p "$LOG_DIR"

# Log session end
echo "{\"event\":\"session_end\",\"timestamp\":\"$TIMESTAMP\",\"session_id\":\"$SESSION_ID\",\"cwd\":\"$CWD\"}" >> "$LOG_FILE"

# Optional: Clean up old logs (keep last 100)
if [ -f "$LOG_FILE" ]; then
  line_count=$(wc -l < "$LOG_FILE")
  if [ "$line_count" -gt 100 ]; then
    tail -100 "$LOG_FILE" > "$LOG_FILE.tmp"
    mv "$LOG_FILE.tmp" "$LOG_FILE"
  fi
fi

exit 0
