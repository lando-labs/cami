# CAMI: Clean MCP-First Architecture

**Date:** 2025-11-13
**Status:** Approved for Implementation
**Assumption:** No backward compatibility needed - clean slate

---

## Core Philosophy

**CAMI is an MCP server that manages Claude Code agents globally.**

- Primary interface: **MCP tools** (via Claude Code)
- Secondary interface: **CLI** (for scripting/quick checks)
- Agent storage: **Global only** (`~/.cami/sources/`)
- No per-project `vc-agents/` directories

---

## Clean Architecture

```
~/.cami/
├── config.yaml              # Single source of truth
├── sources/                 # Global agent sources
│   ├── lando-agents/       # git clone from lando-labs/lando-agents
│   ├── my-agents/          # Local custom agents
│   └── company-agents/     # git clone from company repo (optional)
└── cami                     # Single binary (MCP + CLI modes)
```

**That's it. No vc-agents anywhere.**

### Single Binary Design

**One binary, two modes:**

```go
// Entry point
func main() {
    if len(os.Args) == 1 {
        // No args - show helpful message
        showUsage()
        return
    }

    if os.Args[1] == "--mcp" {
        // MCP server mode
        runMCPServer()
    } else {
        // CLI mode
        runCLI()
    }
}
```

**Why single binary?**
- ✅ **Simpler installation**: One file vs two
- ✅ **Better security**: Version consistency, single checksum
- ✅ **User-friendly**: Less confusion, easier updates
- ✅ **Industry standard**: Docker, Git, Caddy all use this pattern

### Configuration

```yaml
# ~/.cami/config.yaml
version: "1"

agent_sources:
  - name: "lando-agents"
    path: "~/.cami/sources/lando-agents"
    priority: 100
    git:
      enabled: true
      remote: "git@github.com:lando-labs/lando-agents.git"

  - name: "my-agents"
    path: "~/.cami/sources/my-agents"
    priority: 200

deploy_locations:
  - name: "cami"
    path: "/Users/lando/lando-labs/cami"
  - name: "lando-agents"
    path: "/Users/lando/lando-labs/lando-agents"
```

---

## User Workflows

### First Time Setup (via Claude)

```
# 1. Install binary
$ curl -L https://github.com/lando-labs/cami/releases/latest/download/cami-macos -o ~/.cami/cami
$ chmod +x ~/.cami/cami

# 2. Configure MCP in Claude Code settings
{
  "mcpServers": {
    "cami": {
      "command": "/Users/lando/.cami/cami",
      "args": ["--mcp"]
    }
  }
}

# 3. Optional: Add to PATH for CLI
$ ln -s ~/.cami/cami /usr/local/bin/cami

# Then in Claude Code:
User: "Help me set up CAMI"

Claude: *uses mcp__cami__onboard*
→ No config found, guides setup

User: "Add the official agent library"

Claude: *uses mcp__cami__add_source*
URL: git@github.com:lando-labs/lando-agents.git
→ Clones to ~/.cami/sources/lando-agents
→ 29 agents available

User: "Deploy frontend and backend here"

Claude: *uses mcp__cami__deploy_agents*
→ Deploys to /Users/lando/lando-labs/cami/.claude/agents/
```

### Ongoing Usage (via Claude)

```
User: "What agents do I have?"
Claude: *uses mcp__cami__list_agents*
→ Shows 29 from lando-agents

User: "Update my agent sources"
Claude: *uses mcp__cami__update_source*
→ git pull in ~/.cami/sources/lando-agents

User: "Deploy qa agent to my other project"
Claude: *uses mcp__cami__deploy_agents*
→ Deploys to tracked location
```

### Power User (via CLI)

```bash
# Quick list
$ cami list | grep backend

# Scripted deployment
$ cami deploy frontend backend ~/projects/my-app

# Direct source management
$ cd ~/.cami/sources/my-agents
$ vim new-agent.md
$ cami deploy new-agent ~/projects/current
```

---

## Migration from Current State

### Your Current Setup

```
/Users/lando/lando-labs/cami/
├── vc-agents/
│   └── lando-agents/     # 29 agents from extracted repo
└── .claude/agents/       # 14 deployed agents
```

### Migration Steps

**1. Move lando-agents to global location**
```bash
mkdir -p ~/.cami/sources
mv /Users/lando/lando-labs/cami/vc-agents/lando-agents ~/.cami/sources/
```

**2. Create initial config**
```bash
mkdir -p ~/.cami
cat > ~/.cami/config.yaml <<EOF
version: "1"
agent_sources:
  - name: "lando-agents"
    path: "~/.cami/sources/lando-agents"
    priority: 100
    git:
      enabled: true
      remote: "git@github.com:lando-labs/lando-agents.git"
  - name: "my-agents"
    path: "~/.cami/sources/my-agents"
    priority: 200

deploy_locations:
  - name: "cami"
    path: "/Users/lando/lando-labs/cami"
EOF

mkdir -p ~/.cami/sources/my-agents
```

**3. Update MCP server to use global sources**
- Change `getVCAgentsDir()` to always return `~/.cami/sources`
- Update agent loading to use config sources
- Remove all `vc-agents/` references

**4. Clean up old structure**
```bash
# Remove old vc-agents directory
rm -rf /Users/lando/lando-labs/cami/vc-agents

# Update .gitignore to remove vc-agents references
```

**Done! All agents now global.**

---

## Code Changes Required

### 1. Merge Binaries into Single Entry Point

**Create new unified entry point:**
- Merge `cmd/cami-mcp/main.go` into `cmd/cami/main.go`
- Add mode detection logic (--mcp flag)
- Keep all internal packages unchanged

**Files to remove:**
- `cmd/cami-mcp/` directory entirely

### 2. Remove vc-agents Concept

**Files to modify:**
- `cmd/cami/main.go` - Remove `getVCAgentsDir()`, use config sources, add MCP mode
- `internal/cli/init.go` - Create `~/.cami/sources/my-agents/` instead
- `.gitignore` - Remove vc-agents patterns
- `README.md` - Update all references

**Files to remove:**
- `vc-agents/README.md` - No longer needed
- `vc-agents/.gitignore` - No longer needed

### 3. Update Agent Loading

**Current:**
```go
func getVCAgentsDir() (string, error) {
    // Check working directory for vc-agents/
    // Check executable path for vc-agents/
    // ...
}

vcAgentsDir, _ := getVCAgentsDir()
agents, _ := agent.LoadAgents(vcAgentsDir)
```

**New:**
```go
func getAgentSources() ([]config.AgentSource, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }
    return cfg.AgentSources, nil
}

sources, _ := getAgentSources()
agents, _ := agent.LoadAgentsFromSources(sources)
```

### 4. Update Init Command

**Current:**
```go
// Creates ./vc-agents/my-agents in current directory
```

**New:**
```go
// Creates ~/.cami/sources/my-agents globally
// Creates ~/.cami/config.yaml with default config
```

### 5. Unified Main Entry Point

**New structure:**
```go
// cmd/cami/main.go
package main

func main() {
    // Show helpful message if no args
    if len(os.Args) == 1 {
        fmt.Println("CAMI - Claude Agent Management Interface")
        fmt.Println("\nUsage:")
        fmt.Println("  cami [command]     Run CLI command (try: cami --help)")
        fmt.Println("  cami --mcp         Start MCP server")
        fmt.Println("\nInstallation:")
        fmt.Println("  MCP: Configure with --mcp flag in Claude Code settings")
        fmt.Println("  CLI: Add to PATH for command-line usage")
        os.Exit(0)
    }

    // MCP server mode
    if os.Args[1] == "--mcp" {
        runMCPServer()
        return
    }

    // CLI mode (default for all other args)
    runCLI()
}

func runMCPServer() {
    // All logic from cmd/cami-mcp/main.go
    // Use config-based loading
    cfg, err := config.Load()
    if err != nil {
        log.Printf("No config found - user needs to run cami init or use onboard")
    }

    // In each tool handler
    allAgents, _ := agent.LoadAgentsFromSources(cfg.AgentSources)
    // ...
}

func runCLI() {
    // Existing CLI logic from current cmd/cami/main.go
    // Use config-based loading
    cfg, _ := config.Load()
    allAgents, _ := agent.LoadAgentsFromSources(cfg.AgentSources)
    // ...
}
```

---

## Implementation Plan

### Phase 1: Merge Binaries & Config-Based Loading (Day 1, ~2 hours)

**Goal:** Single binary with both modes, using config sources

**Tasks:**
1. Merge `cmd/cami-mcp/main.go` into `cmd/cami/main.go`
   - Add mode detection (`--mcp` flag vs CLI args)
   - Create `runMCPServer()` and `runCLI()` functions
   - Add helpful message for no-args case
2. Update both modes to use config-based loading
   - Change from `getVCAgentsDir()` to `config.Load()`
   - Use `agent.LoadAgentsFromSources(cfg.AgentSources)`
3. Update build to produce single `cami` binary
4. Test both modes work with current setup

**Result:** Single binary, both modes functional, uses config

### Phase 2: Remove vc-agents & Clean Up (Day 1, ~30 minutes)

**Goal:** Eliminate all vc-agents references

**Tasks:**
1. Delete `vc-agents/` directory (after backing up)
2. Delete `cmd/cami-mcp/` directory (merged into cmd/cami)
3. Remove `getVCAgentsDir()` function (if any remnants)
4. Update `cami init` to create `~/.cami/sources/my-agents/`
5. Remove vc-agents from `.gitignore`
6. Update build scripts/Makefile for single binary

**Result:** No more vc-agents, no more separate MCP binary

### Phase 3: Migration Helper (Day 1)

**Goal:** Script to help migrate (for future users)

**Tasks:**
1. Create `cami migrate` command
2. Detects old `vc-agents/` in projects
3. Offers to move to `~/.cami/sources/`
4. Updates config

**Result:** Easy migration path (even though you don't need it)

### Phase 4: Documentation (Day 1)

**Goal:** Update all docs for MCP-first, global-only approach

**Tasks:**
1. Update README - MCP-first installation
2. Update CLAUDE.md - global sources only
3. Update development-plan.md - reflect new architecture
4. Create MIGRATION.md - for future reference

**Result:** Documentation matches reality

---

## File Structure Changes

### Before
```
cami/
├── .gitignore              # Contains vc-agents/** patterns
├── vc-agents/
│   ├── .gitignore         # Ignores subdirs
│   ├── README.md          # Explains workspace
│   └── lando-agents/      # Cloned repo
├── cmd/
│   ├── cami/
│   └── cami-mcp/
└── internal/
    ├── cli/init.go        # Creates ./vc-agents
    └── ...
```

### After
```
cami/
├── .gitignore              # No vc-agents references
├── cmd/
│   └── cami/              # Single binary (MCP + CLI modes)
└── internal/
    ├── cli/init.go        # Creates ~/.cami/sources
    └── ...

# Separate (user's machine)
~/.cami/
├── config.yaml
├── cami                    # Installed binary
└── sources/
    ├── lando-agents/
    └── my-agents/
```

---

## Benefits of Clean Approach

### For You (Now)
1. **Simpler mental model** - One place for all agents, one binary
2. **Less clutter** - No vc-agents in every project
3. **Easier to manage** - Update sources once, available everywhere
4. **Git-friendly** - No need to gitignore per-project agent dirs
5. **Single installation** - Download one file, not two

### For Future Users
1. **Obvious setup** - "Install CAMI, add sources, deploy agents"
2. **No per-project init** - Just deploy agents where needed
3. **Centralized management** - One place to update agents
4. **Clear documentation** - No confusion about local vs global
5. **Transparent modes** - Don't need to know about MCP vs CLI

### For Codebase
1. **Less code** - No vc-agents detection logic, single entry point
2. **Clearer architecture** - Config is source of truth
3. **Easier testing** - Mock config, not filesystem
4. **Better separation** - Storage location is config concern
5. **Simpler builds** - One binary to compile and distribute

### Single Binary Benefits

**Security:**
- ✅ Version consistency (impossible to have mismatched MCP/CLI versions)
- ✅ Single checksum to verify
- ✅ Simpler code auditing
- ✅ No IPC vulnerabilities

**User Experience:**
- ✅ One file to install
- ✅ Easier updates (one file)
- ✅ Less confusion ("which binary do I need?")
- ✅ Transparent mode switching

**Distribution:**
- ✅ Smaller releases (one binary per platform)
- ✅ Simpler CI/CD
- ✅ Clearer trust model

---

## Testing Strategy

### Unit Tests
```go
// Before: Tests depend on vc-agents/ filesystem
func TestLoadAgents(t *testing.T) {
    vcDir := createVCAgentsDir(t)
    agents, _ := LoadAgents(vcDir)
    // ...
}

// After: Tests use config
func TestLoadAgentsFromSources(t *testing.T) {
    sources := []AgentSource{
        {Path: createTestAgentDir(t, "source1")},
        {Path: createTestAgentDir(t, "source2")},
    }
    agents, _ := LoadAgentsFromSources(sources)
    // ...
}
```

### Integration Tests
```bash
# Test global sources work
$ rm -rf ~/.cami
$ cami init
$ cami source add git@github.com:lando-labs/lando-agents.git
$ cami list | grep frontend
✓ frontend found

# Test deployment works
$ cami deploy frontend ~/test-project
$ ls ~/test-project/.claude/agents/
frontend.md
```

---

## Rollout Plan

### Step 1: Your Migration (15 minutes)

```bash
# 1. Backup current state
cp -r vc-agents vc-agents.backup

# 2. Create global structure
mkdir -p ~/.cami/sources
mv vc-agents/lando-agents ~/.cami/sources/
mkdir -p ~/.cami/sources/my-agents

# 3. Create config
cat > ~/.cami/config.yaml <<EOF
version: "1"
agent_sources:
  - name: "lando-agents"
    path: "~/.cami/sources/lando-agents"
    priority: 100
    git:
      enabled: true
      remote: "git@github.com:lando-labs/lando-agents.git"
  - name: "my-agents"
    path: "~/.cami/sources/my-agents"
    priority: 200
deploy_locations:
  - name: "cami"
    path: "/Users/lando/lando-labs/cami"
EOF
```

### Step 2: Code Updates (1-2 hours)

1. Merge cmd/cami-mcp into cmd/cami with mode detection
2. Update both modes to use config-based loading
3. Remove getVCAgentsDir() function
4. Update cami init to create ~/.cami/sources/my-agents
5. Remove vc-agents/ and cmd/cami-mcp/ directories
6. Update .gitignore
7. Update build to produce single binary
8. Run tests

### Step 3: Documentation (30 minutes)

1. Update README
2. Update CLAUDE.md
3. Update development-plan.md
4. Commit changes

### Step 4: Verification (10 minutes)

```bash
# Test single binary - no args
$ cami
# (Should show helpful usage message)

# Test MCP server mode
$ cami --mcp
# (Should start MCP server, load from ~/.cami/sources)

# Test CLI mode
$ cami list
# (Should show agents from global sources)

# Test deployment
$ cami deploy frontend .
# (Should work)

# Update MCP config to use --mcp flag
{
  "mcpServers": {
    "cami": {
      "command": "/Users/lando/.cami/cami",
      "args": ["--mcp"]
    }
  }
}
```

**Total time: ~2-3 hours for clean migration**

---

## Questions & Answers

### Q: What if I want project-specific agents?

**A:** Use priority in config:

```yaml
agent_sources:
  - name: "global-lando"
    path: "~/.cami/sources/lando-agents"
    priority: 100

  - name: "project-specific"
    path: "/path/to/project/custom-agents"
    priority: 200  # Overrides global
```

Still no vc-agents, just different source paths.

### Q: What about the lando-agents repo separation?

**A:** Perfect! It stays separate:

```
lando-labs/lando-agents/     # Upstream repo (29 agents)
  ↓ git clone
~/.cami/sources/lando-agents/ # Your local copy
```

### Q: How do I create new agents?

**A:** Via agent-architect (already deployed):

```
User: "Create a new devops agent"
Agent-architect: *creates agent*
Agent-architect: *saves to ~/.cami/sources/my-agents/devops.md*
User: "Deploy it to this project"
Claude: *uses deploy_agents*
```

### Q: What happens to existing deployed agents?

**A:** They stay! Deployed agents are in `.claude/agents/`:

```
project/.claude/agents/
├── frontend.md        # Stays
├── backend.md         # Stays
└── architect.md       # Stays
```

Only the **source** changes (vc-agents → ~/.cami/sources)

---

## Success Criteria

✅ **Single binary** (`cami`) with both MCP and CLI modes
✅ **Zero** vc-agents references in codebase
✅ **Zero** cmd/cami-mcp directory (merged into cmd/cami)
✅ All agents loaded from `~/.cami/sources/`
✅ All tests pass
✅ MCP server works with `--mcp` flag
✅ CLI works normally (without --mcp)
✅ Documentation reflects MCP-first, single-binary, global-only approach

---

## Recommendation

**Proceed with clean migration immediately.**

**Why:**
- You're the only user (no breaking changes for others)
- Simplifies architecture significantly
- 2-3 hours total work
- Better foundation for future

**Next Steps:**
1. Do your personal migration (15 min)
2. Update code (1-2 hours)
3. Update docs (30 min)
4. Commit and done

---

## Appendix: Before/After Comparison

### Before (Current)

```
User workflow:
1. cd my-project
2. cami init                    # Creates ./vc-agents
3. cami source add ...          # Clones to ./vc-agents/lando-agents
4. cami deploy frontend .

File structure:
project1/
├── vc-agents/
│   └── lando-agents/
└── .claude/agents/

project2/
├── vc-agents/
│   └── lando-agents/          # Duplicate!
└── .claude/agents/
```

### After (Clean)

```
User workflow:
1. cami init                    # Creates ~/.cami/sources (once)
2. cami source add ...          # Clones to ~/.cami/sources/lando-agents (once)
3. cami deploy frontend ~/project1
4. cami deploy frontend ~/project2

File structure:
~/.cami/sources/
└── lando-agents/               # Single copy

project1/.claude/agents/
project2/.claude/agents/
```

**Cleaner. Simpler. Better.**
