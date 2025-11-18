# CAMI MCP Tools Audit & Normalization Plan

**Date**: 2025-11-17
**Purpose**: Audit existing MCP tools and plan normalization system integration

---

## Current MCP Tools (13 total)

### Agent Management (4 tools)

#### 1. `list_agents`
**Purpose**: List all available agents from sources
**Returns**: Agent names, versions, descriptions, categories
**Status**: ✅ Keep - Core functionality
**Integration**: Will need to check project state on each call (Phase 3 auto-detection)

#### 2. `deploy_agents`
**Purpose**: Deploy agents to project `.claude/agents/` directory
**Parameters**: `agent_names[]`, `target_path`, `overwrite`
**Status**: ⚠️ **NEEDS UPDATE** - Must write manifests
**Phase 2 Changes**:
- Write local manifest (`<project>/.claude/cami-manifest.yaml`)
- Write central manifest entry (`~/cami-workspace/deployments.yaml`)
- Calculate content hashes
- Record source information
- Track deployment timestamp

#### 3. `scan_deployed_agents`
**Purpose**: Scan project for deployed agents, compare with sources
**Returns**: Version status (up-to-date, outdated, not-deployed)
**Status**: ⚠️ **NEEDS UPDATE** - Must use manifests
**Phase 2 Changes**:
- Read from manifest instead of scanning directory
- Show source information
- Detect drift (not just version)
**Phase 3 Changes**:
- Add drift detection (content hash comparison)
- Show drift types

#### 4. `update_claude_md`
**Purpose**: Update CLAUDE.md with deployed agents section
**Status**: ✅ Keep - No changes needed
**Note**: Will benefit from manifest data but works independently

### Location Management (3 tools)

#### 5. `add_location`
**Purpose**: Register project directory for tracking
**Status**: ✅ Keep - Still useful
**Note**: Locations become more important with manifest system

#### 6. `list_locations`
**Purpose**: List all tracked project locations
**Status**: ✅ Keep - Core functionality
**Phase 4 Integration**: Show agent count per location from manifests

#### 7. `remove_location`
**Purpose**: Unregister project directory
**Status**: ✅ Keep - Core functionality
**Phase 2 Question**: Should this also remove manifest entry?

### Source Management (4 tools)

#### 8. `list_sources`
**Purpose**: List configured agent sources
**Returns**: Source names, paths, priorities, agent counts
**Status**: ✅ Keep - Core functionality
**Phase 4 Integration**: Show deployment coverage (how many projects use each source)

#### 9. `add_source`
**Purpose**: Clone git repo and add as agent source
**Status**: ✅ Keep - Core functionality
**Note**: Works perfectly for source management

#### 10. `update_source`
**Purpose**: Git pull on sources
**Status**: ✅ Keep - Core functionality
**Phase 4 Integration**: After update, show which projects need redeployment

#### 11. `source_status`
**Purpose**: Show git status of sources
**Status**: ✅ Keep - Core functionality
**Note**: Good for checking uncommitted changes

### Project Creation (1 tool)

#### 12. `create_project`
**Purpose**: Create new project with agents and vision doc
**Status**: ✅ Keep - Core functionality
**Phase 2 Changes**:
- Write manifests on project creation
- Set project state to `cami-native`

### Onboarding (1 tool)

#### 13. `onboard`
**Purpose**: Personalized setup guidance based on current state
**Status**: ✅ Keep - Core functionality
**Phase 4 Integration**: Check for unnormalized projects, offer to normalize

---

## New MCP Tools Needed

### Phase 1: Normalization Foundation

#### `detect_source_state`
**Purpose**: Analyze agent source for CAMI standards compliance
**Parameters**:
- `source_name`: Name of source to analyze (from config)
**Returns**:
```json
{
  "source_name": "my-custom-agents",
  "path": "/Users/lando/cami-workspace/sources/my-custom-agents",
  "is_compliant": false,
  "agent_count": 8,
  "issues": [
    {
      "agent": "frontend.md",
      "problems": ["missing version", "no description"]
    },
    {
      "agent": "backend.md",
      "problems": ["missing category"]
    }
  ],
  "recommendations": {
    "add_versions": ["frontend.md", "auth.md"],
    "add_descriptions": ["frontend.md"],
    "add_categories": ["backend.md", "payment.md"],
    "add_camiignore": true
  }
}
```

**Use Case**: User adds custom source with `add_source`, CAMI detects non-standard agents and offers to normalize the source.

#### `normalize_source`
**Purpose**: Normalize agent source to CAMI standards
**Parameters**:
- `source_name`: Name of source to normalize
- `fixes`: Array of fixes to apply
  - `add_versions`: boolean (add v1.0.0 to agents missing versions)
  - `add_descriptions`: boolean (use agent-architect to generate)
  - `add_categories`: boolean (auto-categorize or ask user)
  - `create_camiignore`: boolean (add standard .camiignore)

**Returns**:
```json
{
  "success": true,
  "changes": [
    "Added version 1.0.0 to frontend.md",
    "Added description to frontend.md",
    "Added category 'specialized' to backend.md",
    "Created .camiignore"
  ],
  "agents_updated": 6,
  "backup_created": "/Users/lando/cami-workspace/sources/my-custom-agents/.cami-backup-20251117"
}
```

#### `detect_project_state`
**Purpose**: Analyze project to determine normalization state
**Parameters**:
- `path`: Project directory to analyze
**Returns**:
```json
{
  "state": "cami-native | cami-legacy | cami-aware | non-cami",
  "has_agents": true,
  "agent_count": 5,
  "has_manifest": false,
  "agents": [
    {
      "name": "frontend",
      "has_version": false,
      "matches_source": "team-agents",
      "is_tracked": false
    }
  ],
  "recommendations": {
    "minimal_required": true,
    "standard_recommended": true,
    "full_optional": true
  }
}
```

#### `normalize_project`
**Purpose**: Normalize project to CAMI standards (AI-guided)
**Parameters**:
- `path`: Project directory
- `level`: "minimal" | "standard" | "full"
- `options`: Object with normalization choices
  - `upgrade_agents`: Array of agent names to upgrade
  - `copy_to_source`: Array of {agent: name, source: name}
  - `skip_agents`: Array of agents to leave alone
  - `custom_overrides`: Array of agents to mark as custom

**Returns**:
```json
{
  "success": true,
  "state_before": "cami-aware",
  "state_after": "cami-native",
  "changes": [
    "Created manifest",
    "Added version to frontend.md",
    "Linked backend.md to team-agents",
    "Copied custom-auth.md to my-agents/"
  ],
  "undo_available": true
}
```

**Implementation Notes**:
- Creates backup before modifying files
- Atomic operation (all or nothing)
- Returns detailed change summary for Claude to present

### Phase 3: Drift Detection

#### `detect_drift`
**Purpose**: Detailed drift analysis for a project
**Parameters**:
- `path`: Project directory
**Returns**:
```json
{
  "total_agents": 5,
  "drift_count": 2,
  "agents": [
    {
      "name": "frontend",
      "drift_type": "content",
      "deployed_version": "1.1.0",
      "source_version": "1.1.0",
      "deployed_hash": "abc123...",
      "source_hash": "def456...",
      "source_name": "team-agents",
      "message": "Local changes detected"
    },
    {
      "name": "backend",
      "drift_type": "version",
      "deployed_version": "1.0.0",
      "source_version": "1.1.0",
      "message": "Update available"
    }
  ]
}
```

### Phase 4: Overview System

#### `overview`
**Purpose**: Complete ecosystem overview
**Parameters**: None (uses current workspace)
**Returns**:
```json
{
  "sources": {
    "count": 2,
    "total_agents": 15,
    "by_source": [
      {
        "name": "team-agents",
        "agent_count": 12,
        "git_status": "clean",
        "deployed_instances": 18,
        "project_count": 3
      }
    ]
  },
  "projects": {
    "count": 3,
    "total_deployments": 22,
    "by_project": [
      {
        "name": "app-1",
        "path": "~/projects/app-1",
        "state": "cami-native",
        "agent_count": 8,
        "drift_count": 1
      }
    ]
  },
  "coverage": {
    "frontend": {"deployed": 3, "total_projects": 3},
    "backend": {"deployed": 3, "total_projects": 3},
    "payment": {"deployed": 1, "total_projects": 3}
  },
  "sync_status": {
    "sources_synced": true,
    "projects_need_update": ["app-1"]
  }
}
```

### Phase 5: Drift Resolution

#### `resolve_drift`
**Purpose**: Apply drift resolution action
**Parameters**:
- `path`: Project directory
- `agent`: Agent name
- `action`: "pull_source" | "save_as" | "mark_custom" | "ignore"
- `options`: Action-specific options
  - `new_name`: For "save_as" action
  - `target_source`: For "save_as" action

**Returns**:
```json
{
  "success": true,
  "action_taken": "pull_source",
  "message": "Overwrote frontend.md with version from team-agents",
  "updated_hash": "xyz789..."
}
```

---

## Tool Integration Matrix

| Tool | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 |
|------|---------|---------|---------|---------|---------|
| `deploy_agents` | - | ✏️ Write manifests | - | - | - |
| `scan_deployed_agents` | - | ✏️ Read manifests | ✏️ Add drift | - | - |
| `list_agents` | ✏️ Check source state | - | ✏️ Auto-detect state | - | - |
| `list_locations` | - | - | - | ✏️ Show stats | - |
| `list_sources` | ✏️ Show compliance | - | - | ✏️ Show coverage | - |
| `add_source` | ✏️ Detect + offer normalize | - | - | - | - |
| `update_source` | - | - | - | ✏️ Show impact | - |
| `create_project` | - | ✏️ Write manifest | - | - | - |
| `onboard` | - | - | - | ✏️ Check normalize | - |

✏️ = Modification required

### Auto-Detection Trigger Points

**Source Normalization Detection**:
- ✅ `add_source` - Check immediately after clone (ALWAYS)
- ✅ `list_sources` - Show compliance status
- ⚠️ First deploy from non-compliant source - Block with helpful message

**Note**: `list_agents` will NOT auto-detect source issues. If a tool needs compliant sources (like drift detection), it will fail gracefully with a clear error message pointing to normalization.

**Project Normalization Detection**:
- ✅ `scan_deployed_agents` - Detect missing manifests (ALWAYS)
- ✅ `onboard` - Check all tracked locations
- ⚠️ First interaction with unnormalized project via tracking tool - Offer to normalize

---

## Manifest Schema Design

### Central Manifest
**Location**: `~/cami-workspace/deployments.yaml`

```yaml
version: "2"  # Schema version for migrations

# Metadata
last_updated: 2025-11-17T22:00:00Z
manifest_format_version: 2

# All deployments across all projects
deployments:
  /Users/lando/projects/my-app:
    state: cami-native  # cami-native | cami-legacy | cami-aware
    normalized_at: 2025-11-17T22:00:00Z
    last_scanned: 2025-11-17T22:00:00Z

    agents:
      - name: frontend
        version: 1.1.0
        source: team-agents
        source_path: /Users/lando/cami-workspace/sources/team-agents/frontend.md
        priority: 50
        deployed_at: 2025-11-17T21:00:00Z
        content_hash: sha256:abc123...
        metadata_hash: sha256:xyz789...
        custom_override: false

      - name: custom-auth
        version: 1.0.0
        source: my-agents
        source_path: /Users/lando/cami-workspace/sources/my-agents/custom-auth.md
        priority: 10
        deployed_at: 2025-11-17T21:00:00Z
        content_hash: sha256:def456...
        metadata_hash: sha256:abc123...
        custom_override: false

  /Users/lando/projects/legacy-app:
    state: cami-aware
    normalized_at: 2025-11-17T22:00:00Z
    last_scanned: 2025-11-17T22:00:00Z

    agents:
      - name: old-agent
        version: null  # No version in frontmatter
        source: untracked
        source_path: null
        deployed_at: 2025-11-17T22:00:00Z
        content_hash: sha256:ghi789...
        needs_upgrade: true
```

### Local Manifest
**Location**: `<project>/.claude/cami-manifest.yaml`

```yaml
version: "2"
state: cami-native
normalized_at: 2025-11-17T22:00:00Z

# Lightweight - just what's needed for project portability
agents:
  - name: frontend
    version: 1.1.0
    source: team-agents
    priority: 50
    deployed_at: 2025-11-17T21:00:00Z
    content_hash: sha256:abc123...

  - name: custom-auth
    version: 1.0.0
    source: my-agents
    priority: 10
    deployed_at: 2025-11-17T21:00:00Z
    content_hash: sha256:def456...
```

---

## Data Flow

### Deployment Flow (Phase 2)
```
1. User: "Deploy frontend and backend to ~/projects/app"
2. CAMI uses deploy_agents tool
3. deploy_agents:
   a. Copies agent files to .claude/agents/
   b. Calculates content hashes
   c. Writes local manifest
   d. Updates central manifest
   e. Returns success
```

### Drift Detection Flow (Phase 3)
```
1. User: "Check status of ~/projects/app"
2. CAMI uses detect_drift tool
3. detect_drift:
   a. Reads local manifest
   b. Loads current agents from sources
   c. Compares versions and hashes
   d. Returns drift report
4. CAMI presents findings to user
```

### Overview Flow (Phase 4)
```
1. User: "Give me an overview"
2. CAMI uses overview tool
3. overview:
   a. Reads central manifest (all deployments)
   b. Reads source configs
   c. Calculates aggregations
   d. Checks git status for sources
   e. Returns comprehensive stats
4. CAMI formats and presents overview
```

---

## Implementation Questions & Decisions

### Q1: Tool Consolidation
**Question**: Should we merge `scan_deployed_agents` and `detect_drift`?

**Options**:
- A) Keep separate (scan = status, drift = detailed analysis)
- B) Merge into single tool with verbosity parameter
- C) Make scan call drift internally

**Decision**: **A - Keep separate**
- `scan_deployed_agents`: Quick status check (fast, used often)
- `detect_drift`: Detailed analysis when investigating issues
- Different use cases warrant separate tools

### Q2: Manifest Sync Strategy
**Question**: How do we keep central ↔ local manifests in sync?

**Options**:
- A) Central is source of truth, rebuild local on demand
- B) Local is source of truth, sync to central on operation
- C) Two-way sync with conflict resolution

**Decision**: **B - Local is source of truth**
- Local manifest lives with the project (portability)
- Operations update local first, then sync to central
- Central manifest is aggregation for overview
- On conflict (e.g., manual edit), local wins

### Q3: Backward Compatibility
**Question**: How do we handle old CAMI users without manifests?

**Options**:
- A) Auto-scan and create manifests on first operation
- B) Require explicit migration command
- C) Show warning, continue without manifests

**Decision**: **A - Auto-scan with confirmation**
- Detect legacy state in first tool call
- Offer to normalize/migrate with explanation
- User confirms, CAMI creates manifests
- Non-breaking: old deployments continue to work

### Q4: Undo System
**Question**: Should normalization support undo?

**Options**:
- A) Create backup directory before modifications
- B) Store changes in manifest for rollback
- C) No undo, just clear documentation

**Decision**: **A - Backup directory with archive management**
- Create `.cami-backup-<timestamp>/` in source or project directory before changes
- Store originals there
- Provide `undo_normalization` tool if needed
- **Archive Management**: Track backup count, propose cleanup when threshold reached (e.g., 10 backups)

#### Archive Cleanup Tool

**New Tool**: `cleanup_backups`
**Purpose**: Remove old normalization backups
**Trigger**: When backup count exceeds threshold (10), CAMI proactively offers cleanup
**Parameters**:
- `path`: Directory to clean (source or project)
- `keep_recent`: Number of recent backups to keep (default: 3)

**Example Workflow**:
```
CAMI detects 12 backups in my-agents/

CAMI: "I noticed you have 12 normalization backups in my-agents/.
      The 3 most recent are from this week, the rest are 30+ days old.

      Would you like me to clean up old backups? I'll keep the 3 most recent."

User: "Yes"

CAMI: *removes 9 old backups*
      ✓ Cleaned up 9 old backups
      ✓ Kept 3 recent backups
      ✓ Freed 2.3 MB
```

---

## Next Steps

1. ✅ Complete this audit
2. Create manifest utilities package (`internal/manifest/`)
3. Create normalization package (`internal/normalize/`)
4. Update `deploy_agents` to write manifests
5. Update `scan_deployed_agents` to use manifests
6. Implement `detect_project_state` tool
7. Implement `normalize_project` tool
8. Update CLAUDE.md with normalization workflows
9. Test with real projects
10. Document migration guide

---

## Source Normalization Workflows

### Scenario 1: User Adds Custom Source

```
User: "Add agent source from git@github.com:myteam/custom-agents.git"

CAMI: *uses add_source tool*
      ✓ Cloned custom-agents to ~/cami-workspace/sources/custom-agents
      ✓ Found 8 agents

      *uses detect_source_state automatically*

      ⚠️ I noticed some agents don't follow CAMI standards:
      • 3 agents missing versions
      • 2 agents missing descriptions
      • No .camiignore file

      I can normalize this source to CAMI standards. This will:
      - Add version 1.0.0 to agents without versions
      - Generate descriptions for agents missing them
      - Create a .camiignore file

      All original files will be backed up first. Want me to normalize?

User: "Yes, do it"

CAMI: *uses normalize_source tool*
      ✓ Added versions to 3 agents
      ✓ Generated descriptions for 2 agents
      ✓ Created .camiignore
      ✓ Backup: ~/cami-workspace/sources/custom-agents/.cami-backup-20251117

      Your source is now CAMI-compliant! All 8 agents are ready to deploy.
```

### Scenario 2: User Creates Local Agents in Source

```
User: "I created some agents in my-agents/ but didn't add all the frontmatter"

CAMI: Let me check your my-agents source...
      *uses detect_source_state for my-agents*

      I found 5 agents with issues:
      • payment.md - missing version
      • auth.md - missing description and category
      • analytics.md - has version but missing category

      Want me to fix these?

User: "Yes, but I'll write my own descriptions"

CAMI: Got it! I'll:
      ✓ Add versions where missing
      ✓ Add categories where missing
      ✗ Skip description generation (you'll add those)

      *uses normalize_source with selective fixes*

      ✓ Fixed 5 agents

      You can now edit the descriptions in:
      - my-agents/payment.md
      - my-agents/auth.md
```

### Scenario 3: Git Source With Non-Standard Agents

```
User: "Add git@github.com:opensource/cool-agents.git"

CAMI: *adds source*
      ✓ Cloned cool-agents (12 agents found)

      *detects non-compliance*

      ⚠️ This source has some agents that don't follow CAMI standards.

      I can normalize them, but this is a git-managed source.

      Options:
      1. Normalize in place (you'll need to commit/push changes)
      2. Create a personal override in my-agents/ (keep original as-is)
      3. Use as-is (some features may not work)

      What would you like?

User: "Option 2 - create overrides"

CAMI: *creates normalized copies in my-agents/ with higher priority*
      ✓ Created 3 normalized overrides in my-agents/
      ✓ Set priority 10 (overrides cool-agents priority 50)

      Now when you deploy these agents, the normalized versions
      from my-agents/ will be used instead of the originals.
```

## Source Normalization vs Project Normalization

| Aspect | Source Normalization | Project Normalization |
|--------|---------------------|----------------------|
| **Target** | Agent files in sources/ | Agents deployed to projects |
| **Purpose** | Ensure agents meet CAMI standards | Track and sync deployments |
| **When** | After adding custom source | After discovering project |
| **Changes** | Add frontmatter, .camiignore | Create manifests, link sources |
| **Git Impact** | May require commit if git source | Only local manifest files |
| **Frequency** | One-time per source | Per project discovered |

## Summary

**Current Tools**: 13 (all have clear purpose)
**Tools to Modify**: 4 (deploy_agents, scan_deployed_agents, create_project, onboard)
**New Tools Needed**: 8
- Phase 1: detect_source_state, normalize_source, detect_project_state, normalize_project, cleanup_backups
- Phase 3: detect_drift
- Phase 4: overview
- Phase 5: resolve_drift

**Total After Implementation**: 21 tools

**Key Design Decisions**:
1. ✅ Source normalization triggers on `add_source` only (not list_agents)
2. ✅ Git sources with issues → offer to create normalized overrides in my-agents/
3. ✅ Archive management: Track backups, offer cleanup at threshold (10+)
4. ✅ Both normalizations in Phase 1 (logistically simpler, both needed for manifests)

**Key Insights**:
1. Existing tools are well-designed - we're adding tracking, not replacing
2. Need both **source normalization** (CAMI standards) and **project normalization** (deployment tracking)
3. Most changes are internal (writing/reading manifests), not interface changes
4. Source normalization particularly important for custom/team agents
5. Backup management prevents workspace clutter over time
