# Phase 1 Implementation Plan: Normalization Foundation

**Goal**: Enable CAMI to normalize both sources (CAMI standards) and projects (deployment tracking)

**Timeline**: TBD
**Prerequisites**: Audit complete ✅

---

## Phase 1 Deliverables

### Core Packages
1. `internal/manifest/` - Manifest read/write utilities
2. `internal/normalize/` - Normalization logic for sources & projects
3. `internal/backup/` - Backup and restore functionality

### MCP Tools (5 new)
1. `detect_source_state` - Analyze source compliance
2. `normalize_source` - Fix source agents
3. `detect_project_state` - Analyze project state
4. `normalize_project` - Create manifests, link agents
5. `cleanup_backups` - Archive management

### Updated Tools (2 modified)
1. `add_source` - Auto-detect and offer normalization
2. `list_sources` - Show compliance status

### Documentation
1. Update workspace `CLAUDE.md` with normalization workflows
2. Add normalization examples
3. Migration guide for existing users

---

## Package Design

### 1. `internal/manifest/`

**Purpose**: Manifest file operations (central + local)

```go
package manifest

import (
    "time"
    "gopkg.in/yaml.v3"
)

// ProjectState represents the normalization state of a project
type ProjectState string

const (
    StateCAMINative  ProjectState = "cami-native"  // Fully normalized
    StateCAMILegacy  ProjectState = "cami-legacy"  // Old CAMI format
    StateCAMIAware   ProjectState = "cami-aware"   // Has agents, not tracked
    StateNonCAMI     ProjectState = "non-cami"     // No agents directory
)

// DeployedAgent represents an agent in a manifest
type DeployedAgent struct {
    Name          string    `yaml:"name"`
    Version       string    `yaml:"version"`
    Source        string    `yaml:"source"`          // Source name
    SourcePath    string    `yaml:"source_path"`     // Full path to source file
    Priority      int       `yaml:"priority"`
    DeployedAt    time.Time `yaml:"deployed_at"`
    ContentHash   string    `yaml:"content_hash"`    // SHA256 of normalized content
    MetadataHash  string    `yaml:"metadata_hash"`   // SHA256 of frontmatter only
    CustomOverride bool     `yaml:"custom_override"` // Intentionally customized
    NeedsUpgrade  bool      `yaml:"needs_upgrade,omitempty"` // Missing version, etc.
}

// ProjectManifest represents a project's deployment manifest (local)
type ProjectManifest struct {
    Version      string           `yaml:"version"`       // Schema version
    State        ProjectState     `yaml:"state"`
    NormalizedAt time.Time        `yaml:"normalized_at"`
    Agents       []DeployedAgent  `yaml:"agents"`
}

// ProjectDeployment represents a project in the central manifest
type ProjectDeployment struct {
    State        ProjectState     `yaml:"state"`
    NormalizedAt time.Time        `yaml:"normalized_at"`
    LastScanned  time.Time        `yaml:"last_scanned"`
    Agents       []DeployedAgent  `yaml:"agents"`
}

// CentralManifest represents the central deployments manifest
type CentralManifest struct {
    Version              string                        `yaml:"version"`
    LastUpdated          time.Time                     `yaml:"last_updated"`
    ManifestFormatVersion int                          `yaml:"manifest_format_version"`
    Deployments          map[string]ProjectDeployment `yaml:"deployments"` // Key: absolute project path
}

// Functions
func ReadProjectManifest(projectPath string) (*ProjectManifest, error)
func WriteProjectManifest(projectPath string, manifest *ProjectManifest) error
func ReadCentralManifest() (*CentralManifest, error)
func WriteCentralManifest(manifest *CentralManifest) error
func CalculateContentHash(filePath string) (string, error)
func CalculateMetadataHash(filePath string) (string, error)
func NormalizeContent(content []byte) []byte // Strip whitespace for hashing
```

### 2. `internal/normalize/`

**Purpose**: Source and project normalization logic

```go
package normalize

import (
    "github.com/lando/cami/internal/agent"
    "github.com/lando/cami/internal/manifest"
)

// SourceIssue represents a problem with an agent in a source
type SourceIssue struct {
    AgentFile string
    Problems  []string // "missing version", "no description", etc.
}

// SourceAnalysis represents the compliance state of a source
type SourceAnalysis struct {
    SourceName   string
    Path         string
    IsCompliant  bool
    AgentCount   int
    Issues       []SourceIssue
    MissingCAMIIgnore bool
}

// SourceNormalizationOptions specifies what to fix
type SourceNormalizationOptions struct {
    AddVersions     bool   // Add v1.0.0 to agents missing versions
    AddDescriptions bool   // Generate descriptions (use agent-architect)
    AddCategories   bool   // Auto-categorize agents
    CreateCAMIIgnore bool  // Add .camiignore template
}

// SourceNormalizationResult represents the outcome
type SourceNormalizationResult struct {
    Success       bool
    Changes       []string
    AgentsUpdated int
    BackupPath    string
}

// ProjectAnalysis represents the state of a project for normalization
type ProjectAnalysis struct {
    Path            string
    State           manifest.ProjectState
    HasAgentsDir    bool
    HasManifest     bool
    AgentCount      int
    Agents          []AgentAnalysis
    Recommendations ProjectRecommendations
}

// AgentAnalysis represents a single agent in a project
type AgentAnalysis struct {
    Name           string
    HasVersion     bool
    Version        string
    MatchesSource  string // Source name if matches, empty if no match
    IsTracked      bool   // In manifest?
    NeedsUpgrade   bool
}

// ProjectRecommendations suggests normalization actions
type ProjectRecommendations struct {
    MinimalRequired    bool // Create manifests
    StandardRecommended bool // Link sources, add versions
    FullOptional       bool // Rewrite with agent-architect
}

// ProjectNormalizationLevel specifies depth of normalization
type ProjectNormalizationLevel string

const (
    LevelMinimal  ProjectNormalizationLevel = "minimal"  // Just manifests
    LevelStandard ProjectNormalizationLevel = "standard" // Manifests + source links
    LevelFull     ProjectNormalizationLevel = "full"     // Complete rewrite
)

// ProjectNormalizationOptions specifies what to do
type ProjectNormalizationOptions struct {
    Level           ProjectNormalizationLevel
    UpgradeAgents   []string          // Agents to add versions to
    CopyToSource    map[string]string // agent name -> source name
    SkipAgents      []string          // Leave these alone
    CustomOverrides []string          // Mark as intentionally customized
}

// ProjectNormalizationResult represents the outcome
type ProjectNormalizationResult struct {
    Success      bool
    StateBefore  manifest.ProjectState
    StateAfter   manifest.ProjectState
    Changes      []string
    BackupPath   string
    UndoAvailable bool
}

// Functions
func AnalyzeSource(sourceName string, sourcePath string) (*SourceAnalysis, error)
func NormalizeSource(sourceName string, options SourceNormalizationOptions) (*SourceNormalizationResult, error)
func AnalyzeProject(projectPath string, availableSources []config.AgentSource) (*ProjectAnalysis, error)
func NormalizeProject(projectPath string, options ProjectNormalizationOptions) (*ProjectNormalizationResult, error)
```

### 3. `internal/backup/`

**Purpose**: Backup and archive management

```go
package backup

import (
    "time"
)

// BackupInfo represents a backup directory
type BackupInfo struct {
    Path      string
    Timestamp time.Time
    SizeBytes int64
}

// ArchiveAnalysis represents backup state
type ArchiveAnalysis struct {
    TotalBackups   int
    TotalSizeBytes int64
    OldestBackup   time.Time
    NewestBackup   time.Time
    Backups        []BackupInfo
}

// CleanupOptions specifies what to keep
type CleanupOptions struct {
    KeepRecent int // Number of recent backups to keep
}

// CleanupResult represents the outcome
type CleanupResult struct {
    RemovedCount  int
    FreedBytes    int64
    KeptBackups   []string
}

// Functions
func CreateBackup(targetPath string) (backupPath string, err error)
func ListBackups(targetPath string) ([]BackupInfo, error)
func AnalyzeArchive(targetPath string) (*ArchiveAnalysis, error)
func CleanupBackups(targetPath string, options CleanupOptions) (*CleanupResult, error)
func RestoreFromBackup(backupPath string, targetPath string) error
```

---

## Implementation Order

### Step 1: Manifest Package (Week 1)
- [ ] Design and implement manifest structs
- [ ] YAML read/write utilities
- [ ] Content hash calculation (normalized)
- [ ] Metadata hash calculation (frontmatter only)
- [ ] Unit tests for all functions

**Dependencies**: None
**Blockers**: None

### Step 2: Backup Package (Week 1)
- [ ] Backup creation (copy directory with timestamp)
- [ ] Backup listing and analysis
- [ ] Cleanup logic (keep N recent)
- [ ] Restore functionality
- [ ] Unit tests

**Dependencies**: None
**Blockers**: None

### Step 3: Normalize Package - Source (Week 2)
- [ ] Source analysis logic
  - [ ] Scan agents for missing frontmatter
  - [ ] Check for .camiignore
  - [ ] Categorize issues
- [ ] Source normalization
  - [ ] Add missing versions
  - [ ] Add missing categories
  - [ ] Create .camiignore from template
  - [ ] Integration with backup package
- [ ] Unit tests
- [ ] Integration tests

**Dependencies**: manifest, backup packages
**Blockers**: None

### Step 4: Normalize Package - Project (Week 2-3)
- [ ] Project analysis logic
  - [ ] Detect project state
  - [ ] Scan deployed agents
  - [ ] Match agents to sources
  - [ ] Generate recommendations
- [ ] Project normalization
  - [ ] Minimal: Create manifests only
  - [ ] Standard: Add versions, link sources
  - [ ] Full: Integration with agent-architect (future)
  - [ ] Integration with backup package
- [ ] Unit tests
- [ ] Integration tests

**Dependencies**: manifest, backup, agent packages
**Blockers**: None

### Step 5: MCP Tools - Source Normalization (Week 3)
- [ ] `detect_source_state` tool
  - [ ] Parameters and response types
  - [ ] Integration with normalize.AnalyzeSource
  - [ ] Error handling
- [ ] `normalize_source` tool
  - [ ] Parameters and response types
  - [ ] Integration with normalize.NormalizeSource
  - [ ] Backup creation
  - [ ] Error handling and rollback
- [ ] Update `add_source` tool
  - [ ] Call detect_source_state after clone
  - [ ] Present findings to user
  - [ ] Offer normalization
- [ ] Update `list_sources` tool
  - [ ] Show compliance status
  - [ ] Indicator for non-compliant sources

**Dependencies**: normalize package (source)
**Blockers**: None

### Step 6: MCP Tools - Project Normalization (Week 4)
- [ ] `detect_project_state` tool
  - [ ] Parameters and response types
  - [ ] Integration with normalize.AnalyzeProject
  - [ ] Error handling
- [ ] `normalize_project` tool
  - [ ] Parameters and response types
  - [ ] Integration with normalize.NormalizeProject
  - [ ] Backup creation
  - [ ] Error handling and rollback
- [ ] `cleanup_backups` tool
  - [ ] Parameters and response types
  - [ ] Integration with backup.CleanupBackups
  - [ ] Threshold detection
  - [ ] Proactive suggestions

**Dependencies**: normalize package (project)
**Blockers**: None

### Step 7: Documentation (Week 4)
- [ ] Update workspace CLAUDE.md
  - [ ] Add source normalization workflows
  - [ ] Add project normalization workflows
  - [ ] Add examples for common scenarios
- [ ] Create migration guide
  - [ ] For existing CAMI users (legacy format)
  - [ ] For new users with existing agents
  - [ ] For teams sharing sources
- [ ] Update README
  - [ ] Normalization feature overview
  - [ ] Link to detailed docs

**Dependencies**: All tools implemented
**Blockers**: None

### Step 8: Integration Testing (Week 5)
- [ ] End-to-end source normalization tests
  - [ ] Add non-compliant source
  - [ ] Detect issues
  - [ ] Normalize
  - [ ] Verify changes
  - [ ] Test backup/restore
- [ ] End-to-end project normalization tests
  - [ ] Discover project with agents
  - [ ] Analyze state
  - [ ] Normalize (minimal, standard)
  - [ ] Verify manifests
  - [ ] Test backup/restore
- [ ] Archive management tests
  - [ ] Create multiple backups
  - [ ] Trigger cleanup threshold
  - [ ] Verify cleanup
- [ ] Real-world scenario testing
  - [ ] Test with actual projects
  - [ ] Test with team git sources
  - [ ] Test migration from old CAMI

**Dependencies**: All packages and tools
**Blockers**: None

---

## Testing Strategy

### Unit Tests
- All packages must have >80% code coverage
- Test both happy path and error cases
- Mock file system operations where needed

### Integration Tests
- Test package interactions (manifest + backup, etc.)
- Test MCP tool workflows
- Use temporary directories for testing

### Manual Testing Scenarios
1. Add custom source with missing frontmatter → normalize
2. Discover legacy project → normalize to standard
3. Create 15 backups → trigger cleanup
4. Normalize git source → verify override in my-agents/
5. Undo normalization → verify restore

---

## Success Criteria

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

---

## Open Questions

1. **Version inference**: When adding version to agents without one, always use v1.0.0 or try to infer from git history?
   - Proposal: Always v1.0.0 for simplicity

2. **Description generation**: Use agent-architect to generate descriptions or leave blank?
   - Proposal: Leave blank, mark as needing description, let user fill in

3. **Category auto-assignment**: How to auto-categorize agents without explicit category?
   - Proposal: Analyze agent name and description with simple heuristics (frontend → specialized, test → infrastructure, etc.)

4. **Manifest migration**: Should we auto-migrate old format, or require manual trigger?
   - Proposal: Auto-detect and offer migration, require user confirmation

5. **Backup location**: Should backups be in `.cami-backup-*/` or `.cami/backups/<timestamp>/`?
   - Proposal: `.cami-backup-<timestamp>/` in same directory as target (easier to find)

---

## Next Phase Preview

**Phase 2**: Deployment Tracking Integration
- Update `deploy_agents` to write manifests
- Update `scan_deployed_agents` to read manifests
- Update `create_project` to set state to cami-native
- Sync local ↔ central manifests on operations

This requires Phase 1 to be complete (manifest + normalization packages must exist).
