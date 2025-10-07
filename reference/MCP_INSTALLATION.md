# CAMI MCP Server Installation Guide

Quick installation guide for the CAMI MCP server.

## Prerequisites

- Go 1.24.4 or later
- CAMI project at `/Users/lando/Development/cami`
- Claude Desktop installed

## Installation Steps

### 1. Build the Server

```bash
cd /Users/lando/Development/cami
go build -o cami-mcp cmd/cami-mcp/main.go
```

**Expected output**:
- Binary created: `./cami-mcp`
- Size: ~7-8 MB
- Architecture: arm64 (macOS)

### 2. Verify Build

```bash
ls -lh cami-mcp
```

Should show executable permissions and correct size.

### 3. Configure Claude Desktop

Edit your Claude Desktop configuration:

**File**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Add** (or merge with existing):

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

**Important**:
- Use **absolute paths** for both `command` and `CAMI_VC_AGENTS_DIR`
- Ensure valid JSON syntax (use a JSON validator if needed)
- Don't include trailing commas

### 4. Restart Claude Desktop

1. Quit Claude Desktop completely
2. Relaunch Claude Desktop
3. Wait for initialization (~5-10 seconds)

### 5. Verify Installation

Open a new conversation in Claude and try:

```
What agents are available in CAMI?
```

Claude should use the `list_agents` tool and show available agents.

## Verification Checklist

- [ ] Binary built successfully (`./cami-mcp` exists)
- [ ] Claude Desktop config updated with absolute paths
- [ ] Claude Desktop restarted
- [ ] Claude responds to agent queries using CAMI tools

## Troubleshooting

### Server Not Appearing

**Symptoms**: Claude doesn't use CAMI tools

**Solutions**:
1. Check config file syntax (must be valid JSON)
2. Verify binary path is absolute
3. Restart Claude Desktop again
4. Check logs: `~/Library/Logs/Claude/mcp*.log`

### "vc-agents directory not found"

**Symptoms**: Server starts but tools fail

**Solutions**:
1. Verify `CAMI_VC_AGENTS_DIR` points to correct location
2. Check directory exists: `ls /Users/lando/Development/cami/vc-agents`
3. Ensure path is absolute (not relative)

### Permission Errors

**Symptoms**: Cannot execute binary or read files

**Solutions**:
```bash
# Make binary executable
chmod +x /Users/lando/Development/cami/cami-mcp

# Verify permissions
ls -l /Users/lando/Development/cami/cami-mcp
```

### Server Crashes

**Symptoms**: Server exits immediately

**Solutions**:
1. Check stderr logs
2. Verify Go version: `go version`
3. Rebuild: `go build -o cami-mcp cmd/cami-mcp/main.go`
4. Check Claude logs for error details

## Log Locations

- **MCP Server Logs**: `~/Library/Logs/Claude/mcp-cami.log`
- **Claude Desktop Logs**: `~/Library/Logs/Claude/app.log`
- **Server stderr**: Included in MCP logs

View logs in real-time:
```bash
tail -f ~/Library/Logs/Claude/mcp*.log
```

## Testing Tools

### Test 1: List Agents
```
User: "What agents are available?"
Expected: Claude lists all agents with versions
```

### Test 2: Deploy Agent
```
User: "Deploy the architect agent to /Users/username/test-project"
Expected: Claude confirms deployment and creates files
```

### Test 3: Scan Project
```
User: "What agents are in /Users/username/my-project?"
Expected: Claude shows deployed agents and their status
```

### Test 4: Update Docs
```
User: "Update the CLAUDE.md file in /Users/username/my-project"
Expected: Claude updates file with agent documentation
```

## Uninstallation

To remove the MCP server:

1. Remove from Claude Desktop config:
   ```bash
   # Edit and remove "cami" section
   nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. Restart Claude Desktop

3. Optionally delete binary:
   ```bash
   rm /Users/lando/Development/cami/cami-mcp
   ```

## Next Steps

After installation:

1. Read **[MCP_README.md](MCP_README.md)** for usage guide
2. See **[reference/tool-catalog.md](reference/tool-catalog.md)** for tool details
3. Check **[reference/development-guide.md](reference/development-guide.md)** for customization

## Support

If issues persist:

1. Check all documentation in `reference/`
2. Review logs (stderr + Claude Desktop)
3. Verify Go version and dependencies
4. Report issues with:
   - Go version (`go version`)
   - Error logs
   - Steps to reproduce

---

**Installation Time**: ~5 minutes
**Difficulty**: Easy
**Status**: Production Ready
