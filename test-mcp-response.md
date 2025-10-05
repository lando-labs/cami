# MCP Response Format Test

## Testing the Schema Fix

### Before Fix (BROKEN)

The structured content was returning arrays directly:

```json
{
  "content": [
    {
      "type": "text",
      "text": "Available agents (12 total):\n\n..."
    }
  ],
  "structuredContent": [
    {
      "name": "architect",
      "version": "1.0.0"
    }
  ]  // ❌ ERROR: Array at top level!
}
```

**Error**: `Expected object, received array at path ["structuredContent"]`

### After Fix (WORKING)

The structured content now returns objects wrapping arrays:

```json
{
  "content": [
    {
      "type": "text",
      "text": "Available agents (12 total):\n\n..."
    }
  ],
  "structuredContent": {
    "agents": [
      {
        "name": "architect",
        "version": "1.0.0"
      }
    ]
  }  // ✅ SUCCESS: Object with array inside!
}
```

## Code Changes Summary

### Added Wrapper Types

```go
// Wraps []DeployResult
type DeployAgentsResponse struct {
    Results []DeployResult `json:"results"`
}

// Wraps []AgentInfo
type ListAgentsResponse struct {
    Agents []AgentInfo `json:"agents"`
}

// Wraps []AgentStatusInfo
type ScanDeployedAgentsResponse struct {
    Statuses []AgentStatusInfo `json:"statuses"`
}
```

### Updated Return Statements

**list_agents**:
```go
// Before: return ..., agentInfos, nil
// After:
return &mcp.CallToolResult{...}, &ListAgentsResponse{Agents: agentInfos}, nil
```

**deploy_agents**:
```go
// Before: return ..., deployResults, nil
// After:
return &mcp.CallToolResult{...}, &DeployAgentsResponse{Results: deployResults}, nil
```

**scan_deployed_agents**:
```go
// Before: return ..., statusInfos, nil
// After:
return &mcp.CallToolResult{...}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
```

## Why This Matters

The MCP protocol specification requires:

1. **Content** field: Array of content blocks (text, resources, etc.)
2. **StructuredContent** field: **MUST be a JSON object** (not array)

From the Go SDK docs:

> The Out type argument must also be a **map or struct**.

Arrays/slices marshal to JSON arrays `[]`, not objects `{}`, which violates the MCP spec.

## Testing Steps

1. **Rebuild**: `go build -o cami-mcp cmd/cami-mcp/main.go` ✅
2. **Restart Claude Code** to reload MCP server
3. **Test tools**:
   - `mcp__cami__list_agents` - Should list all agents
   - `mcp__cami__scan_deployed_agents` - Should scan project
   - `mcp__cami__deploy_agents` - Should deploy agents

All tools should now work without schema validation errors!
