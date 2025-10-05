# MCP Schema Validation Bug Fix

## Problem

The CAMI MCP server tools were failing with schema validation errors:

```
Error: Expected object, received array
Path: ["structuredContent"]
```

**Affected Tools**:
- `mcp__cami__scan_deployed_agents`
- `mcp__cami__list_agents`
- `mcp__cami__deploy_agents`

## Root Cause

According to the MCP protocol specification and the Go SDK documentation:

> **StructuredContent** is an optional value that represents the structured result of the tool call. **It must marshal to a JSON object.**

The tool handlers were returning **arrays/slices** (e.g., `[]AgentInfo`, `[]DeployResult`, `[]AgentStatusInfo`) as the structured output, but the MCP spec requires the `structuredContent` field to be a **JSON object**, not an array.

From the Go SDK `AddTool` documentation:

> The Out type argument must also be a **map or struct**.

## Solution

Wrapped all array responses in structs to ensure they marshal to JSON objects:

### Before (Broken)
```go
// Returns array directly - VIOLATES MCP SPEC
return &mcp.CallToolResult{
    Content: []mcp.Content{
        &mcp.TextContent{Text: responseText},
    },
}, agentInfos, nil  // []AgentInfo - array!
```

### After (Fixed)
```go
// Wrapper struct for MCP compliance
type ListAgentsResponse struct {
    Agents []AgentInfo `json:"agents"`
}

// Returns struct wrapping array - COMPLIES WITH MCP SPEC
return &mcp.CallToolResult{
    Content: []mcp.Content{
        &mcp.TextContent{Text: responseText},
    },
}, &ListAgentsResponse{Agents: agentInfos}, nil  // Object!
```

## Changes Made

### 1. Added Response Wrapper Structs

**File**: `cmd/cami-mcp/main.go`

```go
// DeployAgentsResponse wraps the deployment results (MCP requires object, not array)
type DeployAgentsResponse struct {
    Results []DeployResult `json:"results"`
}

// ListAgentsResponse wraps the agent list (MCP requires object, not array)
type ListAgentsResponse struct {
    Agents []AgentInfo `json:"agents"`
}

// ScanDeployedAgentsResponse wraps the status info (MCP requires object, not array)
type ScanDeployedAgentsResponse struct {
    Statuses []AgentStatusInfo `json:"statuses"`
}
```

### 2. Updated Tool Handlers

**deploy_agents** - Changed from:
```go
}, deployResults, nil
```
To:
```go
}, &DeployAgentsResponse{Results: deployResults}, nil
```

**list_agents** - Changed from:
```go
}, agentInfos, nil
```
To:
```go
}, &ListAgentsResponse{Agents: agentInfos}, nil
```

**scan_deployed_agents** - Changed from:
```go
}, statusInfos, nil
```
To:
```go
}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
```

## Expected Response Format

After the fix, tools now return properly structured MCP responses:

```json
{
  "content": [
    {
      "type": "text",
      "text": "Available agents (12 total):\n\n• architect (v1.0.0)\n  ..."
    }
  ],
  "structuredContent": {
    "agents": [
      {
        "name": "architect",
        "version": "1.0.0",
        "description": "...",
        "file_name": "architect.md"
      }
    ]
  },
  "isError": false
}
```

**Key Points**:
- ✅ `structuredContent` is now an **object** (with `agents` field)
- ✅ Array data is wrapped inside the object
- ✅ Complies with MCP protocol specification
- ✅ Claude Code can successfully parse the response

## Testing

1. **Rebuild the binary**:
   ```bash
   go build -o cami-mcp cmd/cami-mcp/main.go
   ```

2. **Restart Claude Code MCP server**:
   - The MCP server will automatically use the new binary
   - Or manually restart Claude Code

3. **Test the tools**:
   - Try `mcp__cami__list_agents` - should return list of available agents
   - Try `mcp__cami__scan_deployed_agents` - should scan project directory
   - Try `mcp__cami__deploy_agents` - should deploy selected agents

## Verification

After the fix, all tools should:
- ✅ Return successful responses without schema validation errors
- ✅ Provide both human-readable text (in `content`)
- ✅ Provide structured data (in `structuredContent` as object)
- ✅ Allow Claude Code to parse and use the structured data

## Files Modified

- **cmd/cami-mcp/main.go** - Added wrapper structs and updated tool handlers
- **cami-mcp** (binary) - Rebuilt with fixes

## Date

2025-10-05
