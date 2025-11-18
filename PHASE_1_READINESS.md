# Phase 1 Readiness Report

**Status**: ✅ Ready to Begin Implementation
**Date**: 2025-11-18
**Goal**: Normalization Foundation (Source + Project)

---

## Planning Complete

All planning documents are complete and comprehensive:

### 1. MCP Tools Audit ✅
**File**: `MCP_TOOLS_AUDIT.md`

- Audited all 13 existing MCP tools
- Confirmed no redundancy or overlap
- Designed 8 new tools for normalization
- Identified 4 tools requiring manifest integration
- **Result**: 21 total tools after Phase 1 (13 current + 8 new)

### 2. Implementation Plan ✅
**File**: `PHASE_1_IMPLEMENTATION.md`

- 3 core packages designed with full Go interfaces
- 8-step timeline (5 weeks)
- Testing strategy defined (>80% coverage target)
- Success criteria established (11 checkpoints)
- Open questions documented

### 3. Testing Patterns ✅
**Reference**: `internal/config/config_test.go`

Analyzed existing test patterns:
- testify/assert and testify/require
- Temporary directory isolation with `t.TempDir()`
- Environment variable override patterns
- Table-driven test structure
- Error message validation

---

## Design Decisions Confirmed

All 4 key decisions made by user:

| Decision | User Choice | Status |
|----------|-------------|--------|
| Auto-detection frequency | `add_source` only (conditional on `list_agents`) | ✅ Documented |
| Git source handling | Create overrides in my-agents/ | ✅ Documented |
| Backup retention | Archive process with threshold cleanup (10+ → keep 3) | ✅ Documented |
| Phase priority | Both normalizations together (whichever easier) | ✅ Documented |

---

## Package Architecture Ready

### 1. `internal/manifest/`
**Purpose**: Manifest read/write utilities

**Key Types**:
- `ProjectManifest` - Local project manifest
- `CentralManifest` - Central deployments manifest
- `DeployedAgent` - Agent tracking with hashes
- `ProjectState` - State enum (cami-native, cami-legacy, cami-aware, non-cami)

**Functions**:
- `ReadProjectManifest(projectPath string) (*ProjectManifest, error)`
- `WriteProjectManifest(projectPath string, manifest *ProjectManifest) error`
- `ReadCentralManifest() (*CentralManifest, error)`
- `WriteCentralManifest(manifest *CentralManifest) error`
- `CalculateContentHash(filePath string) (string, error)`
- `CalculateMetadataHash(filePath string) (string, error)`

**Dependencies**: None
**Blockers**: None

### 2. `internal/backup/`
**Purpose**: Backup and archive management

**Key Types**:
- `BackupInfo` - Backup directory metadata
- `ArchiveAnalysis` - Backup state analysis
- `CleanupOptions` - Cleanup configuration
- `CleanupResult` - Cleanup outcome

**Functions**:
- `CreateBackup(targetPath string) (backupPath string, err error)`
- `ListBackups(targetPath string) ([]BackupInfo, error)`
- `AnalyzeArchive(targetPath string) (*ArchiveAnalysis, error)`
- `CleanupBackups(targetPath string, options CleanupOptions) (*CleanupResult, error)`
- `RestoreFromBackup(backupPath string, targetPath string) error`

**Dependencies**: None
**Blockers**: None

### 3. `internal/normalize/`
**Purpose**: Source and project normalization logic

**Key Types**:
- `SourceAnalysis` - Source compliance state
- `SourceNormalizationOptions` - Source fix options
- `SourceNormalizationResult` - Source fix outcome
- `ProjectAnalysis` - Project normalization state
- `ProjectNormalizationOptions` - Project fix options
- `ProjectNormalizationResult` - Project fix outcome
- `ProjectNormalizationLevel` - Enum (minimal, standard, full)

**Functions**:
- `AnalyzeSource(sourceName string, sourcePath string) (*SourceAnalysis, error)`
- `NormalizeSource(sourceName string, options SourceNormalizationOptions) (*SourceNormalizationResult, error)`
- `AnalyzeProject(projectPath string, availableSources []config.AgentSource) (*ProjectAnalysis, error)`
- `NormalizeProject(projectPath string, options ProjectNormalizationOptions) (*ProjectNormalizationResult, error)`

**Dependencies**: manifest, backup, agent packages
**Blockers**: None

---

## MCP Tools Ready for Implementation

### New Tools (5)
1. `detect_source_state` - Analyze source compliance
2. `normalize_source` - Fix source agents
3. `detect_project_state` - Analyze project state
4. `normalize_project` - Create manifests, link agents
5. `cleanup_backups` - Archive management

### Updated Tools (2)
1. `add_source` - Auto-detect and offer normalization
2. `list_sources` - Show compliance status

---

## Implementation Timeline

**Week 1**:
- Manifest package (read/write, hashing)
- Backup package (create, restore, cleanup)

**Week 2**:
- Normalize package - source analysis and normalization
- Begin project normalization

**Week 3**:
- Complete project normalization
- Implement source normalization MCP tools

**Week 4**:
- Implement project normalization MCP tools
- Update documentation (CLAUDE.md, README)

**Week 5**:
- Integration testing (E2E workflows)
- Real-world scenario validation

---

## Testing Strategy Confirmed

### Unit Tests
- Target: >80% code coverage
- Use testify/assert and testify/require
- Test happy path + error cases
- Mock file system operations

### Integration Tests
- Test package interactions
- Test MCP tool workflows
- Use temporary directories

### Manual Testing Scenarios
1. Add custom source with missing frontmatter → normalize
2. Discover legacy project → normalize to standard
3. Create 15 backups → trigger cleanup
4. Normalize git source → verify override in my-agents/
5. Undo normalization → verify restore

---

## Success Criteria (11 Checkpoints)

Phase 1 is complete when:
- ✅ All 5 new MCP tools are implemented and tested
- ✅ 2 existing tools are updated with normalization
- ✅ Manifest read/write works for central + local
- ✅ Source normalization fixes missing frontmatter
- ✅ Project normalization creates manifests and links sources
- ✅ Backup system creates and manages archives
- ✅ Archive cleanup triggers at threshold
- ✅ Documentation complete (CLAUDE.md, migration guide)
- ✅ All tests pass (unit + integration)
- ✅ Manual testing scenarios validated
- ✅ Code review complete (via qa agent)

---

## Open Questions for Implementation

These can be answered during implementation:

1. **Version inference**: When adding version to agents without one, always use v1.0.0 or try to infer from git history?
   - **Proposal**: Always v1.0.0 for simplicity

2. **Description generation**: Use agent-architect to generate descriptions or leave blank?
   - **Proposal**: Leave blank, mark as needing description, let user fill in

3. **Category auto-assignment**: How to auto-categorize agents without explicit category?
   - **Proposal**: Analyze agent name and description with simple heuristics

4. **Manifest migration**: Should we auto-migrate old format, or require manual trigger?
   - **Proposal**: Auto-detect and offer migration, require user confirmation

5. **Backup location**: Should backups be in `.cami-backup-*/` or `.cami/backups/<timestamp>/`?
   - **Proposal**: `.cami-backup-<timestamp>/` in same directory as target

---

## Next Steps (When Ready)

### Option 1: Start with Manifest Package
```bash
# Create package structure
mkdir -p internal/manifest
touch internal/manifest/manifest.go
touch internal/manifest/manifest_test.go

# Implement core types and functions
# Target: 1 week
```

### Option 2: Start with Backup Package
```bash
# Create package structure
mkdir -p internal/backup
touch internal/backup/backup.go
touch internal/backup/backup_test.go

# Implement backup operations
# Target: 1 week
```

### Option 3: Do Both in Parallel (Recommended)
Both packages are independent and can be developed simultaneously:
- Manifest package (Week 1)
- Backup package (Week 1)

This maximizes velocity and allows normalize package to start Week 2 with both dependencies ready.

---

## Files to Reference During Implementation

### Existing Code Patterns
- `internal/config/config.go` - Config struct patterns
- `internal/config/config_test.go` - Test patterns
- `internal/agent/agent.go` - Agent struct and frontmatter parsing
- `internal/deploy/deploy.go` - Deployment patterns
- `cmd/cami/main.go` - MCP tool implementations

### Documentation
- `MCP_TOOLS_AUDIT.md` - Tool specifications
- `PHASE_1_IMPLEMENTATION.md` - Package designs and timeline
- `CLAUDE.md` - Workspace-level guidance

---

## Risk Assessment

**Low Risk Items** ✅:
- Manifest package (straightforward YAML read/write)
- Backup package (standard file operations)
- Testing infrastructure (patterns established)

**Medium Risk Items** ⚠️:
- Content hashing with normalization (whitespace handling)
- Git source override workflow (user confusion potential)
- Archive cleanup threshold detection (UX design)

**Mitigation Strategies**:
- Implement hashing with comprehensive test coverage
- Clear documentation and examples for git source workflow
- Proactive user guidance in cleanup_backups tool

---

## Confidence Level

**Overall Readiness**: ✅ 95%

**Why 95% and not 100%**:
- 5 open questions that can be resolved during implementation
- Content hash normalization needs testing with real-world edge cases
- Archive cleanup UX needs refinement based on user feedback

**What makes us 95% ready**:
- ✅ All design decisions confirmed by user
- ✅ Package interfaces fully specified
- ✅ No architectural unknowns
- ✅ Testing patterns established
- ✅ Timeline realistic and achievable
- ✅ Dependencies clear and minimal
- ✅ Success criteria well-defined

---

## Recommendation

**Phase 1 is ready to begin implementation**.

The planning is comprehensive, design decisions are confirmed, and the implementation path is clear. We can start with Week 1 (manifest + backup packages) immediately.

**Awaiting user confirmation to proceed**.
