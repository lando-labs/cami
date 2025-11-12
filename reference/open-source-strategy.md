<!--
AI-Generated Documentation
Created by: architect
Date: 2025-11-10
Purpose: Comprehensive open source strategy for CAMI covering testing, versioning, workflows, and user experience
-->

# CAMI Open Source Strategy

**Version:** 1.0
**Status:** Planning Phase
**Target Release:** v1.0.0

This document provides a comprehensive architecture for preparing CAMI (Claude Agent Management Interface) for open source release. It addresses testing strategy, remote agent version control, workflow philosophy, repository structure, CLI UX, and development workflows.

---

## Table of Contents

1. [Testing Architecture](#1-testing-architecture)
2. [Remote Agent Sources Design](#2-remote-agent-sources-design)
3. [Multi-Workflow Support](#3-multi-workflow-support)
4. [Repository Strategy](#4-repository-strategy)
5. [CLI UX Improvements](#5-cli-ux-improvements)
6. [Development Workflow Documentation](#6-development-workflow-documentation)
7. [Implementation Roadmap](#implementation-roadmap)
8. [Open Source Readiness Checklist](#open-source-readiness-checklist)

---

## 1. Testing Architecture

### Current State
- **Zero test coverage** - No tests exist
- Manual testing only
- No CI/CD validation
- Difficult to refactor with confidence

### Goals
- **Production-ready test coverage** before open source release
- **Automated CI/CD testing** on every commit
- **Confidence in refactoring** and contributions
- **Documentation through tests** for contributors

### Testing Framework Selection

**Recommendation: Go standard `testing` package + `testify` for assertions**

**Rationale:**
- Go's testing package is well-designed, fast, and universally understood
- `testify/assert` provides readable assertions without reinventing the wheel
- `testify/mock` enables clean mocking patterns
- Minimal dependencies align with Go philosophy
- Easy onboarding for contributors familiar with Go

**Rejected alternatives:**
- `ginkgo/gomega`: Too heavy, different paradigm, steeper learning curve
- Pure `testing` only: Assertions too verbose, reduces readability

### Coverage Targets

**Overall Target: 80% coverage** with pragmatic exceptions:

| Package | Target | Priority | Rationale |
|---------|--------|----------|-----------|
| `internal/agent` | 95% | Critical | Core parsing logic, data integrity |
| `internal/deploy` | 90% | Critical | File operations, conflict detection |
| `internal/docs` | 90% | Critical | CLAUDE.md manipulation, data loss risk |
| `internal/config` | 85% | High | Configuration management |
| `internal/discovery` | 85% | High | Version comparison logic |
| `internal/cli` | 70% | Medium | Cobra integration, mostly glue code |
| `internal/tui` | 50% | Low | Bubbletea UI, high mocking cost |
| `cmd/*` | 40% | Low | Entry points, integration tests cover |

**Philosophy:** Focus on critical paths and business logic. Don't obsess over glue code and UI.

### Test Types and Structure

#### 1. Unit Tests
**Location:** `*_test.go` files alongside source
**Focus:** Individual functions and methods in isolation

**Example structure:**
```go
// internal/agent/agent_test.go
package agent_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/lando/cami/internal/agent"
)

func TestLoadAgent_ValidFrontmatter(t *testing.T) {
    // Arrange
    agentPath := "testdata/valid-agent.md"

    // Act
    ag, err := agent.LoadAgent(agentPath)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, "test-agent", ag.Name)
    assert.Equal(t, "1.0.0", ag.Version)
    assert.NotEmpty(t, ag.Content)
}

func TestLoadAgent_MissingFrontmatter(t *testing.T) {
    agentPath := "testdata/no-frontmatter.md"
    _, err := agent.LoadAgent(agentPath)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "missing frontmatter")
}
```

**Test data location:** `testdata/` subdirectories (Go convention)

#### 2. Integration Tests
**Location:** `internal/*_integration_test.go` or `test/integration/`
**Focus:** Multiple packages working together, file I/O

**Example:**
```go
// internal/deploy/integration_test.go
//go:build integration

package deploy_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/lando/cami/internal/agent"
    "github.com/lando/cami/internal/deploy"
)

func TestDeployAgent_FullWorkflow(t *testing.T) {
    // Create temp directory
    tmpDir := t.TempDir()

    // Load agent from real file
    ag, err := agent.LoadAgent("../../vc-agents/core/architect.md")
    require.NoError(t, err)

    // Deploy
    result, err := deploy.DeployAgent(ag, tmpDir, false)
    require.NoError(t, err)
    assert.True(t, result.Success)

    // Verify file exists
    deployedPath := filepath.Join(tmpDir, ".claude", "agents", "architect.md")
    assert.FileExists(t, deployedPath)

    // Verify content
    content, err := os.ReadFile(deployedPath)
    require.NoError(t, err)
    assert.Contains(t, string(content), "name: architect")
}
```

**Build tag:** Use `//go:build integration` to separate from unit tests

#### 3. E2E Tests (CLI/MCP)
**Location:** `test/e2e/`
**Focus:** Full command execution, MCP protocol compliance

**Example:**
```go
// test/e2e/cli_test.go
//go:build e2e

package e2e_test

import (
    "os/exec"
    "testing"
    "encoding/json"
)

func TestCLI_ListAgents_JSON(t *testing.T) {
    cmd := exec.Command("../../cami", "list", "--output", "json")
    output, err := cmd.CombinedOutput()
    require.NoError(t, err)

    var result struct {
        Count int `json:"count"`
        Agents []struct {
            Name string `json:"name"`
        } `json:"agents"`
    }

    err = json.Unmarshal(output, &result)
    require.NoError(t, err)
    assert.Greater(t, result.Count, 0)
}

func TestCLI_Deploy_ConflictDetection(t *testing.T) {
    tmpDir := t.TempDir()

    // First deployment should succeed
    cmd := exec.Command("../../cami", "deploy",
        "--agents", "architect",
        "--location", tmpDir)
    err := cmd.Run()
    require.NoError(t, err)

    // Second deployment without --overwrite should fail
    cmd = exec.Command("../../cami", "deploy",
        "--agents", "architect",
        "--location", tmpDir)
    err = cmd.Run()
    assert.Error(t, err) // Expect conflict
}
```

### Mocking Strategy

**File System Operations:**
Use `afero` for virtual filesystem testing (standard in Go community)

```go
import "github.com/spf13/afero"

func TestWithVirtualFS(t *testing.T) {
    // Create in-memory filesystem
    fs := afero.NewMemMapFs()

    // Create test files
    afero.WriteFile(fs, "/test/agent.md", []byte("content"), 0644)

    // Test functions that accept afero.Fs interface
    // (requires refactoring to dependency injection)
}
```

**Refactoring needed:** Inject filesystem dependency into functions:
```go
// Before
func LoadAgent(filePath string) (*Agent, error)

// After
func LoadAgent(fs afero.Fs, filePath string) (*Agent, error)

// Wrapper for backward compatibility
func LoadAgentFromDisk(filePath string) (*Agent, error) {
    return LoadAgent(afero.NewOsFs(), filePath)
}
```

**MCP Protocol Testing:**
Mock MCP server/client interactions using interfaces:

```go
type MCPServer interface {
    CallTool(name string, args map[string]interface{}) (*Result, error)
}

type mockMCPServer struct {
    mock.Mock
}

func (m *mockMCPServer) CallTool(name string, args map[string]interface{}) (*Result, error) {
    argsRet := m.Called(name, args)
    return argsRet.Get(0).(*Result), argsRet.Error(1)
}
```

### Test Data Fixtures

**Structure:**
```
internal/agent/testdata/
├── valid-agent.md              # Complete valid agent
├── invalid-frontmatter.md      # Malformed YAML
├── missing-version.md          # Missing required field
├── no-frontmatter.md           # No YAML block
└── empty-content.md            # Valid frontmatter, no content

internal/deploy/testdata/
├── existing-project/
│   └── .claude/agents/
│       └── architect.md        # Pre-deployed agent
└── empty-project/              # Clean slate

internal/docs/testdata/
├── claude-no-section.md        # CLAUDE.md without agent section
├── claude-with-section.md      # CLAUDE.md with existing section
└── claude-with-old-marker.md   # Old format markers (migration test)
```

### CI/CD Integration (GitHub Actions)

**Workflow file:** `.github/workflows/test.yml`

```yaml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests:
    name: Unit Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Download dependencies
        run: go mod download

      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt

  integration-tests:
    name: Integration Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run integration tests
        run: go test -v -tags=integration ./...

  e2e-tests:
    name: E2E Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build binaries
        run: |
          go build -o cami cmd/cami/main.go
          go build -o cami-mcp cmd/cami-mcp/main.go

      - name: Run E2E tests
        run: go test -v -tags=e2e ./test/e2e/...

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
```

### Coverage Reporting

**Tools:**
- `go test -coverprofile` for local coverage
- Codecov.io for public repos (free for open source)
- Coverage badge in README

**Enforcement:**
```bash
# Make target for local development
make test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out | grep total
    # Fail if below threshold
    @if [ $(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//') -lt 80 ]; then \
        echo "Coverage below 80%"; \
        exit 1; \
    fi
```

### Testing Principles

1. **Test behavior, not implementation** - Focus on inputs/outputs, not internal details
2. **Fast feedback** - Unit tests should run in <5 seconds total
3. **Isolated tests** - No shared state, each test independent
4. **Readable tests** - Tests are documentation
5. **Fail fast** - Clear error messages, easy debugging
6. **Pragmatic coverage** - 80% is good enough, don't chase 100%

### Sample Test Implementation Timeline

**Phase 1: Critical Paths (Week 1-2)**
- [ ] `internal/agent` - 95% coverage
- [ ] `internal/deploy` - 90% coverage
- [ ] Basic CI/CD setup

**Phase 2: High Priority (Week 3)**
- [ ] `internal/docs` - 90% coverage
- [ ] `internal/config` - 85% coverage
- [ ] Integration tests for deployment workflow

**Phase 3: Medium Priority (Week 4)**
- [ ] `internal/discovery` - 85% coverage
- [ ] `internal/cli` - 70% coverage
- [ ] E2E CLI tests

**Phase 4: Polish (Week 5)**
- [ ] MCP server protocol tests
- [ ] TUI basic tests
- [ ] Coverage reporting and badges

---

## 2. Remote Agent Sources Design

### Vision

Enable CAMI to pull agents from multiple sources:
1. **Official Lando Labs agents** (curated, high quality)
2. **Community agents** (third-party repositories)
3. **Private team agents** (company-internal libraries)
4. **Mixed environments** (all three simultaneously)

### Architecture

#### Configuration Format

**Location:** `~/.cami/config.yaml` (migrate from JSON to YAML for better comments/readability)

```yaml
# CAMI Configuration
version: "1.0"

# Deployment locations (existing)
deploy_locations:
  - name: my-project
    path: /Users/lando/projects/my-app
  - name: cami-dev
    path: /Users/lando/lando-labs/cami

# Agent sources (NEW)
agent_sources:
  # Official Lando Labs agents (default)
  - name: official
    type: github
    repo: lando-labs/lando-agents
    ref: main
    enabled: true
    priority: 100  # Higher priority = preferred for conflicts

  # Community agents
  - name: community
    type: github
    repo: awesome-claude/agents
    ref: v1.2.0  # Pin to specific version
    enabled: true
    priority: 50

  # Private company agents (SSH auth)
  - name: acme-internal
    type: git
    url: git@github.com:acme-corp/agents.git
    ref: main
    auth: ssh-key
    enabled: true
    priority: 75

  # Local development agents
  - name: my-agents
    type: local
    path: ~/my-custom-agents
    enabled: true
    priority: 200  # Highest priority for local dev

# Cache settings
cache:
  directory: ~/.cami/cache
  ttl: 3600  # Seconds (1 hour)
  auto_update: false  # Manual updates only
```

#### Source Types

**1. GitHub (Shorthand)**
```yaml
- name: official
  type: github
  repo: lando-labs/lando-agents  # Expands to github.com/lando-labs/lando-agents
  ref: main                      # Branch, tag, or commit SHA
```

**Implementation:**
- Use GitHub API for metadata
- Clone via HTTPS (no auth) or SSH (with auth)
- Support public and private repos

**2. Git (Full URL)**
```yaml
- name: team-agents
  type: git
  url: https://gitlab.com/myteam/agents.git
  ref: v2.0.0
  auth: token  # Use $CAMI_GIT_TOKEN env var
```

**Implementation:**
- Support any Git hosting (GitHub, GitLab, Bitbucket, self-hosted)
- Auth methods: none (public), ssh-key, token, username/password

**3. Local (Filesystem)**
```yaml
- name: dev-agents
  type: local
  path: /Users/lando/my-agents  # Absolute path
```

**Implementation:**
- Direct filesystem access
- No caching (always live)
- Highest priority by default (dev workflow)

**4. Registry (Future)**
```yaml
- name: registry
  type: registry
  url: https://agents.cami.dev/v1
  api_key: ${CAMI_REGISTRY_KEY}
```

**Implementation:** (Future phase)
- Centralized agent discovery
- Semantic versioning
- Dependency resolution
- Community ratings/reviews

### Caching Strategy

**Cache structure:**
```
~/.cami/cache/
├── sources/
│   ├── official/              # github:lando-labs/lando-agents
│   │   ├── .git/              # Git repository
│   │   ├── vc-agents/         # Agent files
│   │   └── .cami-source.yaml  # Metadata
│   ├── community/             # github:awesome-claude/agents
│   └── acme-internal/         # git@github.com:acme-corp/agents.git
└── index.json                 # Combined agent index
```

**Cache lifecycle:**
1. **First fetch:** Clone repository to cache
2. **Subsequent loads:** Read from cache
3. **Update check:** Compare local ref with remote
4. **Update:** `git pull` or re-clone
5. **TTL:** Automatic update if cache older than TTL (optional)

**Cache commands:**
```bash
cami sources update [source-name]  # Update specific or all sources
cami sources clear                  # Clear cache, force re-fetch
cami sources list                   # Show all sources and status
```

### Update Mechanism

**Manual update (default):**
```bash
cami sources update official  # Update single source
cami sources update           # Update all enabled sources
```

**Automatic update (opt-in):**
```yaml
cache:
  auto_update: true
  ttl: 3600  # Check for updates every hour
```

**Update behavior:**
1. Fetch remote metadata
2. Compare with cached version
3. If different, fetch updates
4. Rebuild agent index
5. Show update summary

**Example output:**
```
Updating agent sources...

✓ official (lando-labs/lando-agents)
  Updated: main (3 new agents, 5 updated)

✓ community (awesome-claude/agents)
  No changes

✗ acme-internal (git@github.com:acme-corp/agents.git)
  Failed: SSH authentication required

2 sources updated, 1 failed
Run 'cami list' to see updated agents
```

### Conflict Resolution

**Scenario:** Multiple sources have agents with the same name

**Resolution strategy: Priority-based**
```yaml
agent_sources:
  - name: my-agents
    priority: 200  # Highest
  - name: official
    priority: 100
  - name: community
    priority: 50   # Lowest
```

**Behavior:**
- Agent from higher priority source wins
- CLI shows source in output:
  ```
  Available Agents (45 total):

  ## Core (8 agents)

    architect (v1.1.0) [official]
      System architecture and design specialist

    architect (v0.9.0) [community]
      (Shadowed by higher priority source)
  ```

**Override command:**
```bash
# Deploy from specific source
cami deploy --agents architect --source community --location ~/project

# Show all versions of an agent
cami show architect --all-sources
```

### Authentication Handling

**SSH Keys:**
- Use system SSH agent automatically
- Respect `~/.ssh/config`
- No credentials stored by CAMI

**Personal Access Tokens:**
- Environment variable: `CAMI_GIT_TOKEN`
- Or per-source in config:
  ```yaml
  - name: private-repo
    type: github
    repo: company/agents
    auth: token
    token: ${GITHUB_TOKEN}  # Read from env
  ```

**Username/Password:**
- Discouraged (use tokens instead)
- Prompt if missing and terminal is interactive
- Never store in config file

### Implementation: internal/sources Package

**New package:** `internal/sources/`

```go
package sources

// Source represents an agent source
type Source struct {
    Name     string
    Type     SourceType
    Config   SourceConfig
    Priority int
    Enabled  bool
}

type SourceType string

const (
    SourceTypeGitHub   SourceType = "github"
    SourceTypeGit      SourceType = "git"
    SourceTypeLocal    SourceType = "local"
    SourceTypeRegistry SourceType = "registry"
)

type SourceConfig struct {
    // GitHub
    Repo string
    Ref  string

    // Git
    URL string
    Auth string

    // Local
    Path string

    // Registry
    RegistryURL string
    APIKey      string
}

// Manager handles multiple agent sources
type Manager struct {
    sources   []*Source
    cacheDir  string
    cacheTTL  time.Duration
}

func NewManager(config *config.Config) (*Manager, error)

func (m *Manager) LoadAllAgents() ([]*agent.Agent, error)

func (m *Manager) UpdateSource(name string) error

func (m *Manager) UpdateAllSources() error

func (m *Manager) ListSources() []*SourceStatus

// SourceStatus represents the current state of a source
type SourceStatus struct {
    Source      *Source
    LastUpdated time.Time
    AgentCount  int
    Status      string // "ok", "outdated", "error"
    Error       string
}
```

**Integration points:**
1. `cmd/cami/main.go` - Use `sources.Manager` instead of direct `agent.LoadAgents`
2. `cmd/cami-mcp/main.go` - Same for MCP server
3. `internal/cli/list.go` - Show source in output
4. New commands: `cami sources ...`

### Migration Path

**Phase 1: Backward compatibility (v0.3.0)**
- Introduce `agent_sources` config
- Default source: local `vc-agents/` directory
- Existing behavior unchanged

**Phase 2: Official repository (v0.4.0)**
- Launch `lando-labs/lando-agents` repository
- Default config includes official source
- Users can opt-in to remote sources

**Phase 3: Full multi-source (v0.5.0)**
- Git source types fully supported
- Conflict resolution active
- Documentation for community contributions

**Phase 4: Registry (v1.0.0)**
- Public registry at `agents.cami.dev`
- Discovery, ratings, dependencies
- Full ecosystem

---

## 3. Multi-Workflow Support

### Workflow Taxonomy

CAMI should support four primary workflows:

#### Workflow A: Agent Library Developer (Lando's Current Workflow)
**User:** Developer contributing to official agents
**Setup:** Works inside CAMI repository

```
~/lando-labs/cami/
├── cmd/cami/          # CLI/TUI
├── vc-agents/         # Edit agents here
└── .claude/           # MCP config points to ./vc-agents
```

**Workflow:**
1. Edit agent file: `vc-agents/core/architect.md`
2. Test with MCP (live reload)
3. Deploy to test project: `./cami deploy --agents architect --location ~/test-project`
4. Commit changes: `git commit -m "Update architect agent"`
5. Push to main branch

**Benefits:**
- Immediate feedback loop
- All tools available (CLI, TUI, MCP)
- Version control integrated

**Configuration:**
```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: local-dev
    type: local
    path: /Users/lando/lando-labs/cami/vc-agents
    priority: 200
```

#### Workflow B: Agent Consumer (Project Developer)
**User:** Developer using CAMI in their projects
**Setup:** CAMI installed globally, agents pulled from remote

```
~/.cami/
├── cache/sources/official/  # Cached from GitHub
└── config.yaml              # Sources configured

~/my-projects/
├── web-app/                 # Project 1
│   └── .claude/agents/      # Deployed agents
└── mobile-app/              # Project 2
    └── .claude/agents/      # Deployed agents
```

**Workflow:**
1. Install CAMI: `go install github.com/lando-labs/cami/cmd/cami@latest`
2. Configure sources: `cami sources add official github lando-labs/lando-agents`
3. Browse agents: `cami` (TUI) or `cami list`
4. Deploy to project: `cami deploy --agents architect,frontend --location ~/my-projects/web-app`
5. Update agents: `cami sources update && cami scan ~/my-projects/web-app`

**Benefits:**
- Clean separation (tool vs content)
- Always up-to-date agents
- Multiple projects easily managed

**Configuration:**
```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: official
    type: github
    repo: lando-labs/lando-agents
    ref: main
    priority: 100

deploy_locations:
  - name: web-app
    path: /Users/dev/my-projects/web-app
  - name: mobile-app
    path: /Users/dev/my-projects/mobile-app
```

#### Workflow C: Team Library (Company-Internal Agents)
**User:** Team with custom agents
**Setup:** Private Git repository + CAMI

```
github.com/acme-corp/claude-agents/  # Private repo
├── engineering/
│   ├── backend-specialist.md
│   └── frontend-react.md
└── data/
    ├── etl-engineer.md
    └── ml-ops.md
```

**Workflow:**
1. Create private repo: `acme-corp/claude-agents`
2. Add agents following CAMI format
3. Team members configure source:
   ```bash
   cami sources add acme git git@github.com:acme-corp/claude-agents.git --auth ssh-key
   ```
4. Use alongside official agents
5. Priority: company agents override official when names conflict

**Benefits:**
- Company-specific knowledge
- Private intellectual property
- Standardized across team
- Version controlled

**Configuration:**
```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: acme-internal
    type: git
    url: git@github.com:acme-corp/claude-agents.git
    ref: main
    auth: ssh-key
    priority: 150  # Higher than official

  - name: official
    type: github
    repo: lando-labs/lando-agents
    ref: main
    priority: 100
```

#### Workflow D: Polyglot (Mix of Everything)
**User:** Power user with multiple sources
**Setup:** Official + community + local + team

```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: local-experiments
    type: local
    path: ~/my-agents
    priority: 200

  - name: team
    type: git
    url: git@github.com:acme-corp/agents.git
    priority: 150

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100

  - name: community-awesome
    type: github
    repo: awesome-claude/agents
    priority: 75

  - name: community-specialized
    type: github
    repo: specialized-agents/library
    priority: 50
```

**Workflow:**
1. Local dev overrides everything
2. Team agents have priority over public
3. Official agents as stable baseline
4. Community agents for experimentation
5. Explicit source selection when needed

**Benefits:**
- Maximum flexibility
- Experimentation without risk
- Clear precedence
- Easy to disable sources

### Default Recommended Workflow

**For open source users: Workflow B (Agent Consumer)**

**Default config on first run:**
```yaml
# ~/.cami/config.yaml (auto-generated)
version: "1.0"

agent_sources:
  - name: official
    type: github
    repo: lando-labs/lando-agents
    ref: main
    enabled: true
    priority: 100

cache:
  directory: ~/.cami/cache
  ttl: 3600
  auto_update: false

deploy_locations: []
```

**First-run experience:**
```bash
$ cami

Welcome to CAMI v1.0.0!

No agent sources configured. Would you like to:
  1. Use official Lando Labs agents (recommended)
  2. Use local agents directory
  3. Configure custom source

Choice: 1

✓ Added official source (lando-labs/lando-agents)
✓ Fetching agent library...
✓ Loaded 25 agents

Press any key to continue...
```

### Migration Guides

**From v0.2.0 (local only) to v1.0.0 (multi-source):**

```markdown
## Migrating to CAMI v1.0.0

CAMI v1.0.0 introduces multi-source agent management. Your existing
setup will continue to work, but you can now access official agents.

### Automatic Migration

Run CAMI v1.0.0 for the first time:

```bash
./cami
```

You'll see:
- New config format (~/.cami/config.yaml instead of .cami.json)
- Official agent source added automatically
- Your local vc-agents/ directory as a source (if present)

### Manual Configuration

If you prefer manual setup:

1. Delete old config (backed up automatically):
   ```bash
   mv ~/.cami.json ~/.cami.json.bak
   ```

2. Create new config:
   ```bash
   cami sources add official github lando-labs/lando-agents
   ```

3. Restore deployment locations:
   ```bash
   cami location add my-project ~/path/to/project
   ```

### For Agent Developers

If you develop agents in vc-agents/ directory:

```bash
# Add local source with highest priority
cami sources add local-dev local /path/to/cami/vc-agents --priority 200
```

This preserves your development workflow.
```

---

## 4. Repository Strategy

### Recommendation: Separate Repositories (Multi-Repo)

**Structure:**
1. **`lando-labs/cami`** - CLI/TUI/MCP tool
2. **`lando-labs/lando-agents`** - Official curated agents
3. **`lando-labs/cami-docs`** - Documentation website (optional)

### Rationale

**Why separate?**

**Pros:**
- **Clear separation of concerns** - Tool vs content
- **Independent versioning** - Tool v1.0.0, Agents v2.3.0
- **Easier contributions** - Contributors can add agents without touching Go code
- **Smaller repos** - Faster clones, clearer history
- **Different audiences** - Tool users vs agent creators
- **CI/CD independence** - Agent updates don't trigger tool builds
- **License flexibility** - Could use different licenses if needed

**Cons:**
- **More repos to maintain** - 2 repos instead of 1
- **Version coordination** - Ensure compatibility
- **Initial complexity** - New users need to understand structure

**Mitigating cons:**
- Default config points to official agents (seamless)
- Semantic versioning ensures compatibility
- Clear documentation

### Repository: lando-labs/cami

**Purpose:** Core tool implementation

**Structure:**
```
lando-labs/cami/
├── cmd/
│   ├── cami/           # CLI/TUI entry point
│   └── cami-mcp/       # MCP server
├── internal/           # Core packages
├── test/               # Test suites
├── docs/               # Tool documentation
├── .github/
│   └── workflows/      # CI/CD
├── README.md           # Tool overview
├── LICENSE             # MIT or Apache 2.0
├── go.mod
└── Makefile
```

**README focus:**
- Installation instructions
- CLI/TUI usage
- MCP server setup
- Link to agent library
- Contributing to the tool

**Release cycle:** Semantic versioning
- v1.0.0 - Initial open source release
- v1.1.0 - New features (e.g., registry support)
- v1.0.1 - Bug fixes
- v2.0.0 - Breaking changes (e.g., config format change)

### Repository: lando-labs/lando-agents

**Purpose:** Official curated agent library

**Structure:**
```
lando-labs/lando-agents/
├── agents/
│   ├── core/
│   │   ├── architect.md
│   │   ├── backend.md
│   │   ├── frontend.md
│   │   ├── mobile-native.md
│   │   └── qa.md
│   ├── specialized/
│   │   ├── ai-ml-specialist.md
│   │   ├── blockchain-specialist.md
│   │   └── ...
│   ├── infrastructure/
│   ├── integration/
│   ├── design/
│   └── meta/
├── .github/
│   ├── workflows/
│   │   ├── validate.yml   # Validate agent format
│   │   └── release.yml    # Tag releases
│   └── ISSUE_TEMPLATE/
│       ├── new-agent.md   # Template for proposing new agents
│       └── update-agent.md
├── CONTRIBUTING.md         # How to contribute agents
├── README.md              # Agent library overview
├── TEMPLATE.md            # Agent creation template
└── LICENSE                # CC-BY-4.0 or MIT
```

**README focus:**
- Agent catalog (auto-generated?)
- How to use with CAMI
- Contributing new agents
- Agent quality standards
- Category explanations

**Contributing workflow:**
1. Fork repository
2. Create agent in appropriate category
3. Follow TEMPLATE.md format
4. Submit PR with agent description
5. Maintainers review quality
6. Merge and tag release

**Release cycle:** Calendar versioning (CalVer) or semantic
- v2025.01 - January 2025 agents
- v2025.02 - February 2025 agents
- Or v1.0.0, v1.1.0 (new agents), v2.0.0 (format changes)

**Quality gates (CI/CD):**
```yaml
# .github/workflows/validate.yml
name: Validate Agents

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Install CAMI validator
        run: go install github.com/lando-labs/cami/cmd/cami-validate@latest

      - name: Validate all agents
        run: |
          for agent in agents/**/*.md; do
            echo "Validating $agent"
            cami-validate "$agent"
          done

      - name: Check for duplicates
        run: |
          # Ensure no duplicate agent names
          cami-validate --check-duplicates agents/
```

### Agent Quality Standards

**Inclusion criteria for official repo:**

1. **Format compliance:**
   - Valid YAML frontmatter (name, version, description)
   - Semantic versioning
   - Markdown content

2. **Quality bars:**
   - Clear role definition
   - Specific domain expertise
   - Boundaries and limitations documented
   - Examples provided
   - No overlap with existing agents

3. **Content guidelines:**
   - Professional tone
   - Actionable instructions
   - Self-verification checklist
   - Tool usage guidance
   - No offensive content

4. **Maintenance:**
   - Regular updates for Claude API changes
   - Version bumps for improvements
   - Deprecation notices if replaced

**Review process:**
1. PR submitted with new agent
2. Automated validation passes
3. Maintainer reviews quality
4. Feedback and iteration
5. Approval and merge
6. Tag release (triggers version bump)

### Decision: agent-architect in Open Source?

**Question:** Should the `agent-architect` agent be open sourced?

**Recommendation: YES, but with caveats**

**Pros:**
- **Community empowerment** - Anyone can create quality agents
- **Ecosystem growth** - Lower barrier to entry
- **Dogfooding** - We use our own meta-agent
- **Transparency** - Shows our methodology
- **Education** - Teaches agent design principles

**Cons:**
- **Quality control** - Lower quality community agents
- **Support burden** - Questions about agent creation
- **Complexity** - Meta-agent is harder to understand

**Mitigation strategy:**

1. **Separate category:**
   ```
   agents/
   ├── core/            # Essential agents
   ├── specialized/     # Domain-specific
   └── meta/            # Advanced/meta agents
       └── agent-architect.md
   ```

2. **Clear documentation:**
   ```markdown
   # agent-architect

   **Advanced Meta-Agent - For experienced users**

   This agent helps create and optimize other agents. Use with caution.

   Prerequisites:
   - Deep understanding of Claude Code
   - Familiarity with agent architecture
   - Experience using multiple agents

   If you're new to CAMI, start with core agents instead.
   ```

3. **Supplement with template:**
   Create `AGENT_TEMPLATE.md` in the agents repo for simpler agent creation:
   ```markdown
   # Agent Template

   Use this template to create your own agents without agent-architect.

   ---
   name: your-agent-name
   version: "1.0.0"
   description: Brief description (one line)
   ---

   ## Role

   Define the agent's role and expertise...

   ## When to Use

   Describe situations where this agent should be invoked...

   ## Boundaries

   What this agent does NOT do...

   ## Examples

   Provide usage examples...
   ```

**Conclusion:** Open source agent-architect, but position as advanced tool with alternatives for beginners.

### Coordination Between Repositories

**Version compatibility:**
- CAMI tool declares minimum agent format version
- Agents repo declares compatible CAMI versions
- Breaking changes require major version bumps

**Cross-references:**
- CAMI README links to agents repo
- Agents README links to CAMI repo
- Release notes coordinate

**Example compatibility matrix:**
| CAMI Version | Compatible Agents | Notes |
|--------------|-------------------|-------|
| v1.0.x       | v1.x.x, v2.x.x    | All current agents |
| v0.2.x       | v1.x.x            | Local only, no sources |
| v2.0.x       | v3.x.x            | New agent format |

---

## 5. CLI UX Improvements

### Current State Analysis

**Existing commands:**
```bash
cami list [--output json]
cami deploy --agents <names> --location <path> [--overwrite]
cami update-docs --location <path> [--dry-run]
cami scan --location <path> [--output json]
cami locations
cami location add
cami location remove
```

**Strengths:**
- Clear command names
- Consistent `--location` flag
- JSON output option
- Good help text

**Weaknesses:**
- No short flags (`-o` for output, `-l` for location)
- `location` vs `locations` inconsistency
- Missing commands (init, search, upgrade)
- No interactive mode fallback
- Limited feedback during long operations
- No color output

### Improved Command Structure

#### Command Tree (Revised)

```
cami
├── init                       # Initialize CAMI config (new)
├── list [options]             # List agents
├── search <query> [options]   # Search agents (new)
├── show <agent> [options]     # Show agent details (new)
├── deploy [options]           # Deploy agents
├── scan <location> [options]  # Scan deployed agents
├── update <location> [options]# Update deployed agents (renamed from update-docs)
│
├── sources                    # Source management (new)
│   ├── list                   # List configured sources
│   ├── add <name> <type> <url># Add source
│   ├── remove <name>          # Remove source
│   ├── update [name]          # Update source(s)
│   └── info <name>            # Source details
│
├── locations                  # Location management (consistent plural)
│   ├── list                   # List locations
│   ├── add <name> <path>      # Add location
│   └── remove <name>          # Remove location
│
└── docs                       # Documentation commands (new)
    ├── update <location>      # Update CLAUDE.md
    └── generate <agent>       # Generate agent docs
```

### New Commands

#### 1. `cami init` - Interactive Setup

**Purpose:** First-time setup wizard

**Usage:**
```bash
cami init [--quick] [--no-interactive]
```

**Interactive flow:**
```
Welcome to CAMI!

Let's set up your agent management system.

? Where should agents be sourced from?
  > Official Lando Labs agents (recommended)
    Local directory
    Custom Git repository
    Skip for now

? Add a deployment location?
  > Yes
    No

? Location name: my-project
? Location path: /Users/lando/projects/my-app
  ✓ Path validated

✓ Configuration created at ~/.cami/config.yaml
✓ Fetched 25 agents from official source

Next steps:
  • Browse agents: cami list
  • Deploy agents: cami deploy --agents architect,frontend --location my-project
  • Launch TUI: cami
```

**Quick mode:**
```bash
cami init --quick
# Uses all defaults, non-interactive
```

#### 2. `cami search <query>` - Search Agents

**Purpose:** Find agents by name, category, or description

**Usage:**
```bash
cami search <query> [--category <cat>] [--source <source>] [-o json]
```

**Examples:**
```bash
# Search by keyword
$ cami search api
Found 3 agents matching "api":

  api-integrator (v1.1.0) [official]
  Integration / API Integration
  Third-party API integration, REST/GraphQL APIs, webhooks

  backend (v1.1.0) [official]
  Core / Backend Development
  APIs, databases, server-side logic (mentions: REST, GraphQL)

  mcp-specialist (v1.1.0) [official]
  Integration / MCP Development
  Model Context Protocol servers, API design

# Search within category
$ cami search --category infrastructure
Found 5 agents in Infrastructure:
  deploy, devops, gcp-firebase, performance-optimizer, security-specialist

# Search in specific source
$ cami search react --source community
```

#### 3. `cami show <agent>` - Agent Details

**Purpose:** Detailed agent information

**Usage:**
```bash
cami show <agent> [--source <source>] [--content]
```

**Example:**
```bash
$ cami show architect

architect (v1.1.0)
Source: official (lando-labs/lando-agents)
Category: Core / Architecture & Design

Description:
  System architecture specialist. Use for analyzing requirements,
  planning technical architecture, technology stack selection, and
  guiding system evolution.

Capabilities:
  • Requirements analysis
  • Architecture design and patterns
  • Technology stack recommendations
  • Refactoring guidance

Deployed to:
  ✓ my-project (/Users/lando/projects/my-app) - v1.1.0
  ⚠ old-project (/Users/lando/projects/old) - v1.0.0 (update available)

Usage:
  cami deploy --agents architect --location my-project
```

**With content:**
```bash
$ cami show architect --content
# Shows full agent markdown content
```

#### 4. `cami sources` - Source Management

See [Remote Agent Sources](#2-remote-agent-sources-design) for detailed design.

**Key commands:**
```bash
cami sources list                                     # List sources
cami sources add official github lando-labs/lando-agents  # Add source
cami sources update [name]                            # Update source(s)
cami sources remove <name>                            # Remove source
cami sources info <name>                              # Source details
```

### Improved Existing Commands

#### `cami list` - Enhanced

**Current:**
```bash
cami list [--output json]
```

**Enhanced:**
```bash
cami list [options]

Options:
  -c, --category <cat>    Filter by category
  -s, --source <source>   Filter by source
  -o, --output <format>   Output format: text, json, table (default: text)
  -v, --verbose           Show more details
  -q, --quiet             Show names only
```

**Examples:**
```bash
# List all agents (default)
cami list

# List core agents only
cami list --category core

# List from specific source
cami list --source official

# Quiet mode (names only, scriptable)
cami list --quiet
# Output: architect backend frontend qa ...

# Table format (new)
cami list --output table
┌─────────────────┬─────────┬──────────────┬─────────────────────────┐
│ Name            │ Version │ Category     │ Source                  │
├─────────────────┼─────────┼──────────────┼─────────────────────────┤
│ architect       │ 1.1.0   │ core         │ official                │
│ backend         │ 1.1.0   │ core         │ official                │
│ frontend        │ 1.1.0   │ core         │ official                │
│ frontend-react  │ 2.0.0   │ specialized  │ my-agents (local)       │
└─────────────────┴─────────┴──────────────┴─────────────────────────┘
```

#### `cami deploy` - Enhanced

**Current:**
```bash
cami deploy --agents <names> --location <path> [--overwrite]
```

**Enhanced:**
```bash
cami deploy [options]

Options:
  -a, --agents <names>     Comma-separated agent names (required*)
  -l, --location <path>    Target location (required*)
  -s, --source <source>    Deploy from specific source
  -f, --force              Overwrite existing files (alias for --overwrite)
      --overwrite          Overwrite existing files
  -d, --dry-run            Show what would be deployed
  -y, --yes                Skip confirmation prompts
  -u, --update-docs        Update CLAUDE.md after deployment
  -o, --output <format>    Output format: text, json (default: text)

* If omitted, interactive prompts will guide you
```

**Examples:**
```bash
# Interactive mode (no flags)
$ cami deploy
? Select agents to deploy: (space to select, enter to confirm)
  [x] architect
  [ ] backend
  [x] frontend
  [ ] qa

? Select location:
  > my-project (/Users/lando/projects/my-app)
    other-project (/Users/lando/projects/other)
    Enter path manually...

? Update CLAUDE.md after deployment? Yes

Deploying 2 agents to my-project...
✓ architect (v1.1.0) deployed
✓ frontend (v1.1.0) deployed
✓ CLAUDE.md updated

# Non-interactive (CI/CD)
$ cami deploy -a architect,frontend -l ~/project -y --update-docs
```

**Progress indicator for long operations:**
```bash
$ cami deploy -a architect,backend,frontend -l ~/project
Deploying agents...
  ✓ architect (1/3)
  ✓ backend (2/3)
  ⣽ frontend (3/3)  # Spinner animation
```

#### `cami scan` - Enhanced

**Current:**
```bash
cami scan --location <path> [--output json]
```

**Enhanced:**
```bash
cami scan <location> [options]

Options:
  -o, --output <format>   Output format: text, json, table
  -u, --updates-only      Show only agents with updates available
  -v, --verbose           Show detailed version information

Positional:
  location                Path to scan (or use --location flag)
```

**Examples:**
```bash
# Scan with positional arg (cleaner)
$ cami scan ~/project

# Show only outdated
$ cami scan ~/project --updates-only
Updates available (2):
  ⚠ backend: v1.0.0 → v1.1.0
  ⚠ frontend: v1.0.0 → v1.1.0

# Table format
$ cami scan ~/project -o table
┌───────────┬──────────────┬─────────────────┬──────────┐
│ Agent     │ Deployed     │ Available       │ Status   │
├───────────┼──────────────┼─────────────────┼──────────┤
│ architect │ v1.1.0       │ v1.1.0          │ ✓ OK     │
│ backend   │ v1.0.0       │ v1.1.0          │ ⚠ Update │
│ frontend  │ v1.0.0       │ v1.1.0          │ ⚠ Update │
└───────────┴──────────────┴─────────────────┴──────────┘
```

#### `cami update` - New Command (Replaces update-docs)

**Purpose:** Update deployed agents and/or documentation

**Usage:**
```bash
cami update <location> [options]

Options:
  -a, --agents <names>    Update specific agents (default: all outdated)
      --all               Update all agents, even if up-to-date
  -d, --docs-only         Only update CLAUDE.md
      --no-docs           Skip CLAUDE.md update
  -f, --force             Force update even if same version
  -y, --yes               Skip confirmation
      --dry-run           Show what would be updated
```

**Examples:**
```bash
# Update all outdated agents + docs
$ cami update ~/project
Found 2 outdated agents:
  backend: v1.0.0 → v1.1.0
  frontend: v1.0.0 → v1.1.0

? Update these agents? Yes

Updating agents...
✓ backend updated to v1.1.0
✓ frontend updated to v1.1.0
✓ CLAUDE.md updated

# Update docs only (old behavior)
$ cami update ~/project --docs-only
✓ CLAUDE.md updated with 4 deployed agents

# Update specific agent
$ cami update ~/project --agents backend
✓ backend updated to v1.1.0
```

### Universal Flags

**Available on all commands:**
```bash
-h, --help              Show help
-v, --version           Show version
-c, --config <path>     Use alternate config file
    --no-color          Disable colored output
    --debug             Enable debug logging
-q, --quiet             Minimal output
```

### Interactive Mode Fallback

**Philosophy:** Required flags can be omitted, CLI prompts interactively

**Example:**
```bash
# Missing required flags
$ cami deploy
? Select agents: (space to select)
  [ ] architect
  [ ] backend
  [x] frontend

? Select location:
  > my-project
    other-project
    [Enter path manually]

? Overwrite existing files? No

# Equivalent non-interactive
$ cami deploy -a frontend -l my-project
```

**Disable interactivity:**
```bash
export CAMI_NO_INTERACTIVE=1
# Or
cami deploy --no-interactive
# Fails if required flags missing
```

### Color and Formatting

**Use colors for clarity:**
- ✓ Green: Success
- ✗ Red: Error
- ⚠ Yellow: Warning
- ℹ Blue: Info
- ⣽ Cyan: Progress spinner

**Respect NO_COLOR environment variable:**
```bash
export NO_COLOR=1
# Or
cami --no-color list
```

**Libraries:**
- `github.com/fatih/color` - ANSI colors
- `github.com/briandowns/spinner` - Progress spinners
- `github.com/olekukonko/tablewriter` - Table formatting

### Error Messages

**Current (example):**
```
Error: failed to load agents: failed to walk vc-agents directory: open vc-agents: no such file or directory
```

**Improved:**
```
Error: Agent directory not found

  CAMI couldn't find the agent directory at:
  /Users/lando/cami/vc-agents

  Possible solutions:
    • Run 'cami init' to set up agent sources
    • Check your config at ~/.cami/config.yaml
    • Ensure CAMI is installed correctly

  Need help? Visit https://github.com/lando-labs/cami/issues
```

**Principles:**
1. **Clear problem statement** - What went wrong?
2. **Context** - What was CAMI trying to do?
3. **Actionable solutions** - How to fix it?
4. **Help pointer** - Where to get more help?

### Help Text Improvements

**Current:**
```bash
$ cami deploy --help
Deploy agents to a target location

Usage:
  cami deploy [flags]

Flags:
      --agents string     Comma-separated agent names
  -h, --help              help for deploy
      --location string   Target location path
      --overwrite         Overwrite existing files
```

**Enhanced:**
```bash
$ cami deploy --help
Deploy agents to a project

Usage:
  cami deploy [OPTIONS]
  cami deploy --agents <names> --location <path>

Description:
  Deploys selected agents to a project's .claude/agents/ directory.
  Creates the directory if it doesn't exist. By default, refuses to
  overwrite existing agent files (use --force to override).

Options:
  -a, --agents <names>     Comma-separated agent names (required*)
  -l, --location <path>    Target project path (required*)
  -f, --force              Overwrite existing files
  -u, --update-docs        Update CLAUDE.md after deployment
  -y, --yes                Skip confirmation prompts
      --dry-run            Show what would be deployed without doing it
  -o, --output <format>    Output format: text, json (default: text)
  -h, --help               Show this help message

  * If omitted, interactive prompts will guide you

Examples:
  # Interactive mode
  cami deploy

  # Deploy specific agents
  cami deploy --agents architect,frontend --location ~/my-project

  # Deploy and update docs
  cami deploy -a backend -l ~/project -u

  # Dry run to preview
  cami deploy -a architect -l ~/project --dry-run

  # CI/CD (non-interactive)
  cami deploy -a architect,backend -l ~/project -y -f

Learn more:
  cami list                 # See available agents
  cami scan <location>      # Check deployed agents
  cami update <location>    # Update deployed agents

Documentation: https://github.com/lando-labs/cami
```

### Command Aliases

**Enable shorter commands:**
```bash
cami ls          # Alias for: cami list
cami src         # Alias for: cami sources
cami loc         # Alias for: cami locations
```

**Implementation:**
```go
// internal/cli/root.go
listCmd := NewListCommand(vcAgentsDir)
listCmd.Aliases = []string{"ls"}

sourcesCmd := NewSourcesCommand()
sourcesCmd.Aliases = []string{"src"}
```

### Completion Scripts

**Generate shell completion:**
```bash
# Bash
cami completion bash > /usr/local/etc/bash_completion.d/cami

# Zsh
cami completion zsh > ~/.zsh/completion/_cami

# Fish
cami completion fish > ~/.config/fish/completions/cami.fish
```

**Cobra supports this out of the box:**
```go
rootCmd.AddCommand(completionCmd)
```

**Usage after setup:**
```bash
$ cami deploy --agents <TAB>
architect  backend  frontend  qa  devops  ...

$ cami deploy --location <TAB>
my-project  other-project  ~/
```

---

## 6. Development Workflow Documentation

### The "Agent Repo as Working Directory" Philosophy

**Current workflow (Lando's approach):**

```
Developer's Workspace:
~/lando-labs/cami/
├── cmd/                        # Tool source code
├── internal/                   # Core packages
├── vc-agents/                  # AGENTS LIVE HERE (content)
│   ├── core/
│   ├── specialized/
│   └── ...
├── .claude/                    # Claude Code config
│   └── settings.local.json     # MCP points to ./vc-agents
├── cami                        # Built CLI
└── cami-mcp                    # Built MCP server
```

**This is brilliant because:**
1. **Immediate feedback** - Edit agent → test with MCP → iterate
2. **All tools available** - CLI, TUI, MCP in one place
3. **Version control** - `git commit` agents alongside tool
4. **No context switching** - Everything in one workspace
5. **Dogfooding** - Use CAMI to manage CAMI agents

**Problem:** This only works for the CAMI project itself (monorepo style)

**Solution for open source:** Support this workflow while enabling others

### Recommended Workflows by User Type

#### 1. Open Source Agent Developer (Contribution Workflow)

**Goal:** Contribute agents to `lando-labs/lando-agents`

**Setup:**
```bash
# Fork and clone agents repo
git clone https://github.com/YOUR_USERNAME/lando-agents.git
cd lando-agents

# Install CAMI globally
go install github.com/lando-labs/cami/cmd/cami@latest

# Configure CAMI to use local agents
cami sources add local-dev local $(pwd)/agents --priority 200

# Add test project
cami locations add test-project ~/test-project
```

**Workflow:**
1. Create new agent: `agents/specialized/my-agent.md`
2. Test with CAMI: `cami show my-agent`
3. Deploy to test project: `cami deploy -a my-agent -l test-project`
4. Use Claude Code with MCP to test agent behavior
5. Iterate on agent content
6. Commit: `git commit -m "Add my-agent"`
7. Push and create PR

**Claude Code MCP config:**
```json
{
  "mcpServers": {
    "cami": {
      "command": "/Users/dev/go/bin/cami-mcp",
      "env": {
        "CAMI_CONFIG": "/Users/dev/.cami/config.yaml"
      }
    }
  }
}
```

**Benefits:**
- Agents repo is the working directory
- All agents visible and editable
- CAMI CLI available for testing
- PR-based contribution

#### 2. CAMI Tool Developer (Tool Contribution Workflow)

**Goal:** Contribute to `lando-labs/cami` tool

**Setup:**
```bash
# Clone CAMI repo
git clone https://github.com/lando-labs/cami.git
cd cami

# Build tool
go build -o cami cmd/cami/main.go
go build -o cami-mcp cmd/cami-mcp/main.go

# Configure to use official agents (not local)
./cami sources add official github lando-labs/lando-agents
```

**Workflow:**
1. Edit Go source code: `internal/cli/deploy.go`
2. Build: `go build -o cami cmd/cami/main.go`
3. Test: `./cami list`
4. Write tests: `internal/cli/deploy_test.go`
5. Run tests: `go test ./...`
6. Commit: `git commit -m "Add new CLI feature"`
7. Push and create PR

**Note:** This repo does NOT contain agents (they're pulled from remote)

#### 3. Hybrid Developer (Both Tool and Agents)

**Goal:** Develop both CAMI tool and agents (Lando's current workflow)

**Setup:**
```bash
# Clone both repos side by side
git clone https://github.com/lando-labs/cami.git
git clone https://github.com/lando-labs/lando-agents.git

cd cami

# Build tool
go build -o cami cmd/cami/main.go
go build -o cami-mcp cmd/cami-mcp/main.go

# Configure to use local agents
./cami sources add local-dev local ../lando-agents/agents --priority 200
```

**Alternative: Keep agents IN cami repo during dev**
```bash
cd cami

# Create local agents directory (gitignored)
mkdir -p vc-agents-dev

# Symlink or copy agents for testing
cp -r ../lando-agents/agents/* vc-agents-dev/

# Configure CAMI
./cami sources add local-dev local ./vc-agents-dev --priority 200
```

**Workflow:**
1. Edit tool: `internal/cli/deploy.go`
2. Edit agents: `../lando-agents/agents/core/architect.md`
3. Build tool: `go build -o cami cmd/cami/main.go`
4. Test integration: `./cami deploy -a architect -l ~/test`
5. Commit to both repos separately
6. Create PRs to both repos

**Benefits:**
- Full control over both tool and content
- Test new features with new agents
- Coordinated releases

#### 4. Team Internal Agent Developer

**Goal:** Develop private agents for company

**Setup:**
```bash
# Create company agents repo
git clone git@github.com:acme-corp/claude-agents.git
cd claude-agents

# Structure like official repo
mkdir -p agents/{core,specialized,custom}

# Install CAMI
go install github.com/lando-labs/cami/cmd/cami@latest

# Configure sources
cami sources add acme-internal local $(pwd)/agents --priority 150
cami sources add official github lando-labs/lando-agents --priority 100
```

**Workflow:**
1. Create agent: `agents/custom/acme-backend.md`
2. Test: `cami show acme-backend`
3. Deploy to team projects: `cami deploy -a acme-backend -l ~/team-project`
4. Commit to company repo: `git commit -m "Add acme-backend"`
5. Push to internal GitHub
6. Team members update: `cami sources update acme-internal`

**Benefits:**
- Private IP protected
- Company-specific knowledge
- Extends official agents
- Standard workflow across team

### Development Best Practices

#### Agent Development

**1. Use the template:**
```bash
# Copy template
cp AGENT_TEMPLATE.md agents/specialized/my-new-agent.md

# Fill in frontmatter
# Edit content
# Test
cami show my-new-agent
```

**2. Follow naming conventions:**
- `kebab-case` for file names
- `snake_case` or `kebab-case` for agent names
- Match file name to agent name
- Use `.md` extension

**3. Version semantically:**
- v1.0.0 - Initial release
- v1.1.0 - Added capabilities
- v1.0.1 - Bug fixes, clarifications
- v2.0.0 - Breaking changes (role redefinition)

**4. Test thoroughly:**
```bash
# Validate format
cami show my-agent

# Deploy to test project
cami deploy -a my-agent -l ~/test-project

# Use with Claude Code
# Verify behavior with real tasks
# Iterate on content

# Check for conflicts
cami list --quiet | grep my-agent
```

**5. Document well:**
- Clear role definition
- Specific when to use
- Examples
- Boundaries
- Self-verification checklist

#### Tool Development

**1. Write tests first:**
```go
// internal/agent/agent_test.go
func TestLoadAgent_ValidFrontmatter(t *testing.T) {
    // Test implementation
}
```

**2. Run tests before committing:**
```bash
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
```

**3. Update documentation:**
- Update README if user-facing change
- Add examples for new commands
- Update help text

**4. Follow Go conventions:**
- `gofmt` your code
- Run `golangci-lint`
- Document exported functions
- Keep packages focused

**5. Test CLI commands:**
```bash
# Build
go build -o cami cmd/cami/main.go

# Test
./cami list
./cami deploy -a architect -l ~/test
./cami --help
```

### Deployment Tracking (Future Feature)

**User request:** "Which projects use the `frontend` agent?"

**Design concept:**

**Database:** `~/.cami/deployments.db` (SQLite)

**Schema:**
```sql
CREATE TABLE deployments (
    id INTEGER PRIMARY KEY,
    agent_name TEXT NOT NULL,
    agent_version TEXT NOT NULL,
    location_name TEXT,
    location_path TEXT NOT NULL,
    deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source TEXT,  -- Which source it came from
    UNIQUE(agent_name, location_path)
);

CREATE INDEX idx_agent ON deployments(agent_name);
CREATE INDEX idx_location ON deployments(location_path);
```

**New command:**
```bash
cami where <agent>

# Example output:
$ cami where frontend

frontend (v1.1.0) is deployed to:

  ✓ my-project (/Users/lando/projects/my-app)
    Deployed: 2025-11-09 14:30:00
    Source: official

  ⚠ old-project (/Users/lando/projects/old-app)
    Deployed: 2025-10-15 09:00:00 (v1.0.0, outdated)
    Source: official

  ✓ test-project (/Users/lando/test)
    Deployed: 2025-11-10 08:45:00
    Source: local-dev

Total: 3 locations
```

**Other useful queries:**
```bash
# Show all deployments
cami deployments list

# Show deployments at a location
cami deployments list --location ~/project

# Show outdated deployments
cami deployments outdated

# Clean up stale entries (location no longer exists)
cami deployments clean
```

**Implementation notes:**
- Track deployments automatically during `cami deploy`
- Update during `cami update`
- Optional feature (can be disabled)
- Privacy: local database only
- Use for breaking change notifications

**Use cases:**
1. **Breaking changes:** "Which projects need updates?"
2. **Deprecation:** "Frontend agent split into react-specialist and vue-specialist"
3. **Team coordination:** "Who's using the experimental agent?"
4. **Cleanup:** "Remove agents from projects that no longer need them"

---

## Implementation Roadmap

### Phase 0: Pre-Release Preparation (Weeks 1-6)
**Goal:** Production-ready testing and code quality

**Week 1-2: Critical Path Testing**
- [ ] Set up test structure (`testdata/`, integration, e2e)
- [ ] Write unit tests for `internal/agent` (95% coverage target)
- [ ] Write unit tests for `internal/deploy` (90% coverage target)
- [ ] Set up GitHub Actions CI/CD

**Week 3: High Priority Testing**
- [ ] Unit tests for `internal/docs` (90% coverage)
- [ ] Unit tests for `internal/config` (85% coverage)
- [ ] Integration tests for deployment workflow
- [ ] Set up code coverage reporting (Codecov)

**Week 4: Medium Priority Testing**
- [ ] Unit tests for `internal/discovery` (85% coverage)
- [ ] Unit tests for `internal/cli` (70% coverage)
- [ ] E2E tests for CLI commands
- [ ] MCP protocol compliance tests

**Week 5: Polish and Documentation**
- [ ] Basic TUI tests (50% coverage)
- [ ] Test documentation
- [ ] Contributing guidelines
- [ ] Code review and cleanup

**Week 6: Pre-Release QA**
- [ ] Manual testing on macOS, Linux, Windows
- [ ] Performance testing (large agent libraries)
- [ ] Documentation review
- [ ] Create migration guide from v0.2.0

**Deliverables:**
- 80%+ test coverage
- Full CI/CD pipeline
- All tests passing
- Documentation complete

### Phase 1: Repository Split (Weeks 7-8)
**Goal:** Separate tool from agents

**Week 7: Agent Repository Creation**
- [ ] Create `lando-labs/lando-agents` repository
- [ ] Migrate `vc-agents/` content to new repo
- [ ] Set up agent validation CI/CD
- [ ] Create CONTRIBUTING.md for agents
- [ ] Create AGENT_TEMPLATE.md
- [ ] Add quality gates (format validation)

**Week 8: Tool Repository Updates**
- [ ] Update `lando-labs/cami` to remove bundled agents
- [ ] Default config points to official agents repo
- [ ] Update documentation (README, guides)
- [ ] Create migration tooling (`cami migrate`)
- [ ] Test backward compatibility

**Deliverables:**
- Two separate repositories
- Migration path documented
- Backward compatibility maintained

### Phase 2: Multi-Source Foundation (Weeks 9-11)
**Goal:** Enable remote agent sources

**Week 9: Source Infrastructure**
- [ ] Implement `internal/sources` package
- [ ] Support local source type
- [ ] Support GitHub source type
- [ ] Cache management (fetch, update, clear)
- [ ] Priority-based conflict resolution

**Week 10: CLI Integration**
- [ ] `cami sources` command group
- [ ] Integrate sources with `list`, `deploy`, `scan`
- [ ] Update configuration format (YAML)
- [ ] Backward compatibility for `.cami.json`
- [ ] Migration from JSON to YAML

**Week 11: Testing and Documentation**
- [ ] Integration tests for multi-source
- [ ] Test conflict resolution
- [ ] Documentation for sources feature
- [ ] User guide for workflows
- [ ] Release notes

**Deliverables:**
- Multi-source support (local, GitHub)
- CLI commands for source management
- Migration guide
- Beta release: v0.3.0-beta

### Phase 3: CLI UX Overhaul (Weeks 12-14)
**Goal:** Improved command structure and UX

**Week 12: New Commands**
- [ ] `cami init` - Interactive setup
- [ ] `cami search` - Agent search
- [ ] `cami show` - Agent details
- [ ] Restructure `cami locations` (consistent plural)
- [ ] Update `cami update` (replaces update-docs)

**Week 13: Enhanced Existing Commands**
- [ ] Short flags (`-a`, `-l`, `-o`)
- [ ] Interactive mode fallback
- [ ] Progress indicators and spinners
- [ ] Color output (respect NO_COLOR)
- [ ] Table output format
- [ ] Better error messages

**Week 14: Polish and Completion**
- [ ] Command aliases (`ls`, `src`)
- [ ] Shell completion scripts
- [ ] Improved help text
- [ ] Documentation updates
- [ ] E2E tests for new commands

**Deliverables:**
- Modern CLI experience
- Interactive and non-interactive modes
- Shell completion
- Release: v0.4.0

### Phase 4: Additional Source Types (Weeks 15-16)
**Goal:** Support Git URLs and authentication

**Week 15: Git Source Type**
- [ ] Implement generic Git source
- [ ] SSH key authentication
- [ ] Token authentication (HTTPS)
- [ ] Support GitLab, Bitbucket, self-hosted
- [ ] Auth configuration and env vars

**Week 16: Testing and Documentation**
- [ ] Integration tests for Git sources
- [ ] Test authentication methods
- [ ] Private repository guide
- [ ] Team workflow documentation
- [ ] Security best practices

**Deliverables:**
- Full Git source support
- Authentication methods
- Private repo workflows
- Release: v0.5.0

### Phase 5: Open Source Launch (Week 17)
**Goal:** Public release

**Week 17: Launch Preparation**
- [ ] Final QA pass
- [ ] Security audit
- [ ] License review (MIT/Apache 2.0)
- [ ] Website/landing page (optional)
- [ ] Blog post/announcement
- [ ] Social media assets
- [ ] GitHub issue templates
- [ ] Community guidelines (CODE_OF_CONDUCT.md)

**Launch:**
- [ ] Tag v1.0.0 release
- [ ] Publish to GitHub
- [ ] Announce on socials
- [ ] Submit to package managers (Homebrew, etc.)
- [ ] Monitor initial feedback

**Deliverables:**
- **CAMI v1.0.0** - Public open source release
- Complete documentation
- Community guidelines
- Launch announcement

### Phase 6: Post-Launch (Ongoing)
**Goal:** Community growth and feature additions

**Short-term (Weeks 18-22):**
- [ ] Community support and issue triage
- [ ] Accept first community contributions
- [ ] Bug fixes and quick wins
- [ ] Homebrew formula
- [ ] Performance optimizations

**Medium-term (Months 2-3):**
- [ ] Deployment tracking database
- [ ] `cami where` command
- [ ] Agent dependency management
- [ ] Multi-project batch operations
- [ ] Agent usage analytics

**Long-term (Months 4-6):**
- [ ] Registry server (`agents.cami.dev`)
- [ ] Web UI dashboard
- [ ] Agent marketplace
- [ ] CI/CD integrations
- [ ] Team collaboration features

### Release Schedule Summary

| Version | Target Date | Focus |
|---------|-------------|-------|
| v0.2.0 | Current | Manual testing, local agents |
| v0.3.0-beta | Week 11 | Multi-source foundation (beta) |
| v0.4.0 | Week 14 | CLI UX overhaul |
| v0.5.0 | Week 16 | Full Git support |
| v1.0.0 | Week 17 | **Open source launch** |
| v1.1.0 | Month 2 | Deployment tracking |
| v1.2.0 | Month 3 | Batch operations |
| v2.0.0 | Month 6 | Registry and marketplace |

---

## Open Source Readiness Checklist

### Code Quality
- [ ] 80%+ test coverage achieved
- [ ] All tests passing in CI/CD
- [ ] No critical security vulnerabilities
- [ ] Code reviewed and cleaned up
- [ ] No hardcoded secrets or credentials
- [ ] Proper error handling throughout
- [ ] Consistent code style (gofmt)
- [ ] All TODOs addressed or filed as issues

### Documentation
- [ ] README.md comprehensive and clear
- [ ] Installation instructions for all platforms
- [ ] Quick start guide (5-minute experience)
- [ ] Complete CLI reference
- [ ] MCP server setup guide
- [ ] Contributing guidelines
- [ ] Code of conduct
- [ ] License file (MIT or Apache 2.0)
- [ ] Architecture documentation
- [ ] API documentation for packages

### Repository Setup
- [ ] Repository public on GitHub
- [ ] Issue templates created
- [ ] PR template created
- [ ] Labels configured (bug, enhancement, help wanted)
- [ ] GitHub Actions CI/CD working
- [ ] Branch protection on main
- [ ] Release process documented
- [ ] Changelog maintained
- [ ] Git tags for releases
- [ ] GitHub Releases with binaries

### Community
- [ ] Code of conduct established
- [ ] Contribution guidelines clear
- [ ] Communication channels set up (Discussions, Discord, etc.)
- [ ] Maintainer team defined
- [ ] Response time expectations set
- [ ] Security policy defined (SECURITY.md)
- [ ] First-time contributor guidance

### Legal
- [ ] License chosen and applied
- [ ] Copyright notices in files
- [ ] Third-party licenses reviewed
- [ ] Trademark considerations addressed
- [ ] Export compliance checked (if applicable)

### Marketing/Launch
- [ ] Landing page or website (optional)
- [ ] Logo and branding
- [ ] Screenshots and demos
- [ ] Video walkthrough (optional)
- [ ] Blog post announcing launch
- [ ] Social media posts prepared
- [ ] Submit to:
  - [ ] Hacker News
  - [ ] Reddit (r/golang, r/claude, r/programming)
  - [ ] Product Hunt (optional)
  - [ ] Dev.to
  - [ ] Lobste.rs
- [ ] Package manager submissions:
  - [ ] Homebrew formula
  - [ ] Go install (works automatically)
  - [ ] apt/yum packages (future)

### Support Infrastructure
- [ ] Issue triage process
- [ ] PR review process
- [ ] Release process automated
- [ ] Version numbering strategy
- [ ] Deprecation policy
- [ ] Backward compatibility guarantees
- [ ] Monitoring/telemetry (opt-in)
- [ ] Error reporting (opt-in)

### Launch Day Checklist
- [ ] All tests passing
- [ ] Documentation proofread
- [ ] Version tagged (v1.0.0)
- [ ] GitHub Release created
- [ ] Binaries uploaded
- [ ] Homebrew formula submitted
- [ ] Announcement blog post published
- [ ] Social media posts sent
- [ ] Monitor GitHub notifications
- [ ] Respond to initial feedback
- [ ] Celebrate! 🎉

---

## Conclusion

This architecture provides a comprehensive roadmap for taking CAMI from a powerful internal tool to a successful open source project. The strategy balances:

- **Quality:** Production-ready testing before launch
- **Flexibility:** Multiple workflows for different user types
- **Clarity:** Clean separation between tool and content
- **Community:** Easy contributions and ecosystem growth
- **Pragmatism:** Phased rollout, backward compatibility
- **Vision:** Path to registry and marketplace

The "agent repo as working directory" philosophy is preserved for developers while enabling simpler workflows for consumers. The multi-source architecture enables official, community, and private agents to coexist harmoniously.

**Next Steps:**
1. Review and refine this architecture with the team
2. Begin Phase 0 (testing) immediately
3. Set up project tracking (GitHub Projects)
4. Assign owners to each phase
5. Start building toward v1.0.0 launch

**Success Metrics:**
- 1000+ GitHub stars in first month
- 50+ community contributions in first quarter
- 100+ projects using CAMI agents
- Active community in Discussions/Discord
- Featured in Go weekly newsletters
- Homebrew installs >500/month

The architecture is ambitious but achievable. CAMI has the potential to become the standard tool for Claude Code agent management. Let's make it happen.
