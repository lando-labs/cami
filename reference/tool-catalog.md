<!--
AI-Generated Documentation
Created by: mcp-specialist
Date: 2025-10-05
Purpose: Complete reference for all CAMI MCP tools
-->

# CAMI MCP Tool Catalog

Complete reference documentation for all tools exposed by the CAMI MCP server.

## Table of Contents

### Agent Management Tools
1. [deploy_agents](#deploy_agents)
2. [update_claude_md](#update_claude_md)
3. [list_agents](#list_agents)
4. [scan_deployed_agents](#scan_deployed_agents)

### Location Management Tools
5. [add_location](#add_location)
6. [list_locations](#list_locations)
7. [remove_location](#remove_location)

---

## deploy_agents

Deploy selected agents to a target project's `.claude/agents/` directory.

### Description

Use this tool when you want to add specific agents to a project. The tool handles conflict detection, creates necessary directories, and writes agent files with full frontmatter and content.

### Parameters

```typescript
{
  "agent_names": string[],    // Required: Array of agent names to deploy
  "target_path": string,      // Required: Absolute path to target project
  "overwrite": boolean        // Optional: Overwrite existing files (default: false)
}
```

#### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `agent_names` | string[] | Yes | Names of agents to deploy (e.g., `["architect", "backend"]`) |
| `target_path` | string | Yes | Absolute path to the target project directory |
| `overwrite` | boolean | No | Whether to overwrite existing agent files (default: `false`) |

### JSON Schema

```json
{
  "type": "object",
  "properties": {
    "agent_names": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Array of agent names to deploy (e.g. ['architect' 'backend'])"
    },
    "target_path": {
      "type": "string",
      "description": "Absolute path to target project directory"
    },
    "overwrite": {
      "type": "boolean",
      "description": "Whether to overwrite existing agent files (default: false)"
    }
  },
  "required": ["agent_names", "target_path"]
}
```

### Response Format

**Text Content**:
```
Deployed 2 agents to /path/to/project

✓ architect: Deployed successfully
✓ backend: Deployed successfully
```

**Structured Data**:
```json
[
  {
    "agent_name": "architect",
    "success": true,
    "message": "Deployed successfully"
  },
  {
    "agent_name": "backend",
    "success": true,
    "message": "Deployed successfully"
  }
]
```

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `invalid target path` | Path doesn't exist or isn't a directory | Provide valid directory path |
| `agents not found` | Agent names not in vc-agents | Check available agents with `list_agents` |
| `File already exists` (conflict) | Agent already deployed and overwrite=false | Set `overwrite: true` or manually remove |
| `deployment failed` | File system error during write | Check permissions |

### Example Invocations

#### Deploy single agent
```json
{
  "agent_names": ["architect"],
  "target_path": "/Users/username/projects/my-app"
}
```

#### Deploy multiple agents with overwrite
```json
{
  "agent_names": ["architect", "backend", "frontend"],
  "target_path": "/Users/username/projects/my-app",
  "overwrite": true
}
```

### Files Created

The tool creates:
- `/path/to/project/.claude/agents/architect.md`
- `/path/to/project/.claude/agents/backend.md`
- etc.

Each file contains full agent content including YAML frontmatter.

---

## update_claude_md

Update a project's `CLAUDE.md` file with documentation about deployed agents.

### Description

Use this tool after deploying agents to document them in the project's `CLAUDE.md` file. The tool adds or updates a managed section with agent information, preserving existing content.

### Parameters

```typescript
{
  "target_path": string      // Required: Absolute path to target project
}
```

#### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `target_path` | string | Yes | Absolute path to the target project directory |

### JSON Schema

```json
{
  "type": "object",
  "properties": {
    "target_path": {
      "type": "string",
      "description": "Absolute path to target project directory"
    }
  },
  "required": ["target_path"]
}
```

### Response Format

**Text Content**:
```
Updated CLAUDE.md at /path/to/project

Documented 3 agents:
  • architect (v1.2.0)
  • backend (v2.1.0)
  • frontend (v1.5.0)
```

**Structured Data**: `null`

### Behavior

1. Scans `.claude/agents/` directory for deployed agents
2. Reads existing `CLAUDE.md` (if exists)
3. Finds managed section markers or creates new section
4. Updates/inserts agent documentation
5. Preserves all other CLAUDE.md content
6. Writes updated file

### Managed Section Format

The tool manages a section between markers:

```markdown
<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
## Deployed Agents

The following Claude Code agents are available in this project:

### architect (v1.2.0)
System architect and planning specialist

### backend (v2.1.0)
Backend development expert

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->
```

**Important**: Content between markers is fully managed by CAMI. Manual edits will be overwritten.

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `invalid target path` | Path doesn't exist | Provide valid directory |
| `no agents directory found` | No `.claude/agents/` directory | Deploy agents first |
| `no agents found` | `.claude/agents/` is empty | Deploy at least one agent |
| `failed to write CLAUDE.md` | Permission error | Check file permissions |

### Example Invocation

```json
{
  "target_path": "/Users/username/projects/my-app"
}
```

---

## list_agents

List all available agents from CAMI's version-controlled agent repository.

### Description

Use this tool to discover what agents are available for deployment. Returns agent names, versions, and descriptions from the `vc-agents` directory.

### Parameters

```typescript
{}  // No parameters required
```

### JSON Schema

```json
{
  "type": "object",
  "properties": {},
  "required": []
}
```

### Response Format

**Text Content**:
```
Available agents (23 total):

• architect (v1.2.0)
  System architect and planning specialist

• backend (v2.1.0)
  Backend development expert specializing in APIs and databases

• frontend (v1.5.0)
  Frontend development with React, Vue, and modern frameworks

...
```

**Structured Data**:
```json
[
  {
    "name": "architect",
    "version": "1.2.0",
    "description": "System architect and planning specialist",
    "file_name": "architect.md"
  },
  {
    "name": "backend",
    "version": "2.1.0",
    "description": "Backend development expert specializing in APIs and databases",
    "file_name": "backend.md"
  },
  ...
]
```

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `failed to load agents` | Cannot read vc-agents directory | Check CAMI installation |
| `vc-agents directory not found` | Missing vc-agents | Set `CAMI_VC_AGENTS_DIR` env var |

### Example Invocation

```json
{}
```

### Use Cases

1. **Discovery**: Find available agents before deployment
2. **Validation**: Verify agent names before calling `deploy_agents`
3. **Documentation**: Generate agent inventory for reports
4. **Version Checking**: See latest agent versions

---

## scan_deployed_agents

Scan a project directory to find deployed agents and compare with available versions.

### Description

Use this tool to audit what agents are deployed in a project. Returns agent status (up-to-date, update-available, not-deployed) by comparing deployed versions with available versions.

### Parameters

```typescript
{
  "target_path": string      // Required: Absolute path to target project
}
```

#### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `target_path` | string | Yes | Absolute path to the target project directory |

### JSON Schema

```json
{
  "type": "object",
  "properties": {
    "target_path": {
      "type": "string",
      "description": "Absolute path to target project directory"
    }
  },
  "required": ["target_path"]
}
```

### Response Format

**Text Content**:
```
Scanning /Users/username/projects/my-app

Found 3 deployed agents

✓ architect: up-to-date (deployed: v1.2.0, available: v1.2.0)
⚠ backend: update-available (deployed: v2.0.0, available: v2.1.0)
○ frontend: not-deployed
✓ deploy: up-to-date (deployed: v1.0.0, available: v1.0.0)
...
```

**Structured Data**:
```json
[
  {
    "name": "architect",
    "deployed_version": "1.2.0",
    "available_version": "1.2.0",
    "status": "up-to-date"
  },
  {
    "name": "backend",
    "deployed_version": "2.0.0",
    "available_version": "2.1.0",
    "status": "update-available"
  },
  {
    "name": "frontend",
    "deployed_version": "",
    "available_version": "1.5.0",
    "status": "not-deployed"
  },
  ...
]
```

### Status Values

| Status | Symbol | Meaning |
|--------|--------|---------|
| `up-to-date` | ✓ | Deployed version matches available version |
| `update-available` | ⚠ | Newer version available in vc-agents |
| `not-deployed` | ○ | Agent not present in project |

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `invalid target path` | Path doesn't exist | Provide valid directory |
| `failed to load agents` | Cannot read vc-agents | Check CAMI installation |
| `failed to read agents directory` | Permission error | Check directory permissions |

### Example Invocations

#### Scan single project
```json
{
  "target_path": "/Users/username/projects/my-app"
}
```

### Use Cases

1. **Audit**: Check which agents are deployed
2. **Version Check**: Find outdated agents
3. **Update Planning**: Identify agents needing updates
4. **Compliance**: Verify required agents are deployed

### Notes

- Scans all available agents (not just deployed ones)
- Compares versions using exact string matching
- Shows all agents even if not deployed (status: not-deployed)
- Empty `.claude/agents/` directory shows all agents as not-deployed

---

## add_location

Add a new deployment location to CAMI's configuration.

### Description

Use this tool to register a project directory for agent deployment. Once added, locations can be listed and removed. This creates a persistent configuration for quick access to frequently used project directories.

### Parameters

```typescript
{
  "name": string,     // Required: Friendly name for the location
  "path": string      // Required: Absolute path to project directory
}
```

#### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Friendly name for the location (e.g., "my-project") |
| `path` | string | Yes | Absolute path to project directory |

### JSON Schema

```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Friendly name for the location (e.g. 'my-project')"
    },
    "path": {
      "type": "string",
      "description": "Absolute path to project directory"
    }
  },
  "required": ["name", "path"]
}
```

### Response Format

**Text Content**:
```
Added location 'my-project' at /Users/username/projects/my-app

Total locations: 3
```

**Structured Data**: `null`

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `location name is required` | Missing name parameter | Provide a name |
| `location path is required` | Missing path parameter | Provide a path |
| `path must be absolute` | Relative path provided | Use absolute path |
| `path does not exist` | Path doesn't exist | Verify path exists |
| `path is not a directory` | Path points to a file | Use directory path |
| `location with name already exists` | Duplicate name | Choose different name |
| `location with path already exists` | Duplicate path | Location already registered |

### Example Invocations

#### Add single location
```json
{
  "name": "my-project",
  "path": "/Users/username/projects/my-app"
}
```

#### Add another location
```json
{
  "name": "client-website",
  "path": "/Users/username/clients/acme-corp/website"
}
```

### Use Cases

1. **Quick Access**: Register frequently used projects
2. **Team Consistency**: Share location names across team
3. **Automation**: Script deployments using location names
4. **Organization**: Maintain list of active projects

---

## list_locations

List all configured deployment locations in CAMI.

### Description

Use this tool to see what project directories are registered for agent deployment. Returns all locations with their names and paths.

### Parameters

```typescript
{}  // No parameters required
```

### JSON Schema

```json
{
  "type": "object",
  "properties": {},
  "required": []
}
```

### Response Format

**Text Content**:
```
Configured locations (3 total):

• my-project
  /Users/username/projects/my-app

• client-website
  /Users/username/clients/acme-corp/website

• experimental
  /Users/username/experiments/new-app
```

**Structured Data**:
```json
{
  "locations": [
    {
      "name": "my-project",
      "path": "/Users/username/projects/my-app"
    },
    {
      "name": "client-website",
      "path": "/Users/username/clients/acme-corp/website"
    },
    {
      "name": "experimental",
      "path": "/Users/username/experiments/new-app"
    }
  ]
}
```

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `failed to load config` | Cannot read config file | Check file permissions |

### Example Invocation

```json
{}
```

### Use Cases

1. **Discovery**: See what locations are configured
2. **Verification**: Confirm location was added
3. **Audit**: Review all registered projects
4. **Selection**: Choose location for deployment

### Notes

- Returns empty list if no locations configured
- Shows helpful message when list is empty
- Locations stored in `~/.cami.json`

---

## remove_location

Remove a deployment location from CAMI's configuration.

### Description

Use this tool to unregister a project directory. This removes the location from the configuration but does not affect the project directory or its files.

### Parameters

```typescript
{
  "name": string      // Required: Name of location to remove
}
```

#### Parameter Details

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Name of location to remove |

### JSON Schema

```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Name of location to remove"
    }
  },
  "required": ["name"]
}
```

### Response Format

**Text Content**:
```
Removed location 'my-project'

Remaining locations: 2
```

**Structured Data**: `null`

### Error Scenarios

| Error | Cause | Solution |
|-------|-------|----------|
| `location name is required` | Missing name parameter | Provide a name |
| `location with name not found` | Location doesn't exist | Check name with list_locations |
| `failed to load config` | Cannot read config file | Check file permissions |
| `failed to save config` | Cannot write config file | Check file permissions |

### Example Invocations

#### Remove single location
```json
{
  "name": "my-project"
}
```

#### Remove by exact name
```json
{
  "name": "client-website"
}
```

### Use Cases

1. **Cleanup**: Remove unused locations
2. **Reorganization**: Clear old locations before adding new ones
3. **Error Correction**: Remove incorrectly added locations
4. **Project Completion**: Remove finished projects

### Notes

- Does not affect the project directory itself
- Does not delete any files
- Only removes the location from CAMI's configuration
- Use `list_locations` to verify removal

---

## Common Patterns

### Typical Workflow

#### Agent Deployment
```
1. list_agents          → Discover available agents
2. deploy_agents        → Deploy selected agents to project
3. update_claude_md     → Document deployed agents
4. scan_deployed_agents → Verify deployment
```

#### Location Management
```
1. add_location         → Register a project directory
2. list_locations       → View all registered locations
3. deploy_agents        → Deploy using saved location path
4. remove_location      → Clean up unused locations
```

#### Combined Workflow
```
1. add_location         → Register project as "my-app"
2. list_agents          → See available agents
3. deploy_agents        → Deploy to /path/from/location
4. update_claude_md     → Document deployed agents
5. scan_deployed_agents → Verify all agents current
```

### Error Handling

All tools return errors in JSON-RPC format:

```json
{
  "error": {
    "code": -32603,
    "message": "invalid target path: path does not exist"
  }
}
```

### Response Structure

Every successful tool call returns:
- **Content**: Array of content items (text, resources, etc.)
- **Data** (optional): Structured data for programmatic use

### Path Requirements

All tools requiring `target_path`:
- Must be **absolute paths** (not relative)
- Must exist and be accessible
- Must be directories (not files)

### Version Format

Agent versions follow semantic versioning: `major.minor.patch`
- Example: `1.2.0`, `2.1.5`, `0.3.0`
- Comparison is exact string match (not semver-aware)
