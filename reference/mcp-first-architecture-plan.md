# CAMI: MCP-First Architecture Plan

**Date:** 2025-11-13
**Status:** Proposal
**Goal:** Rethink CAMI with MCP as the primary interface, CLI as secondary convenience

---

## Current State Analysis

### Architecture Overview

```
cami/
├── cmd/
│   ├── cami/          # CLI binary (~400 LOC)
│   └── cami-mcp/      # MCP server (~860 LOC)
├── internal/
│   ├── agent/         # Core: Load agents from sources
│   ├── config/        # Core: Manage ~/.cami/config.yaml
│   ├── deploy/        # Core: Deploy agents to projects
│   ├── docs/          # Core: Update CLAUDE.md
│   ├── discovery/     # Core: Scan projects for agents
│   ├── cli/           # CLI commands (~2760 LOC total)
│   └── tui/           # TUI for discovery view
└── vc-agents/         # Local workspace (per-project)
```

### Current Workflow Assumptions

**CLI-First Design:**
- User runs `cami init` in project root
- Creates `vc-agents/` in current directory
- Config at `~/.cami/config.yaml` references local paths
- Each project can have its own `vc-agents/`

**MCP Design:**
- MCP server uses same internal packages
- Expects `vc-agents/` in working directory
- 12 tools wrap CLI functionality

### Code Distribution

| Component | Lines of Code | Purpose |
|-----------|---------------|---------|
| CLI Commands | ~2760 | Terminal interface |
| MCP Server | ~860 | Claude Code interface |
| Internal Packages | ~1500 | Shared business logic |
| **Total** | **~5120** | |

**Key Insight:** MCP is 30% the size of CLI but handles 90% of usage.

---

## Philosophy Shift: MCP-First

### Core Principles

1. **MCP is the Primary Interface**
   - Most users interact through Claude Code
   - MCP tools are the "main" API
   - CLI wraps MCP functionality (not vice versa)

2. **Global-First, Local Optional**
   - Default: `~/.cami/sources/` (global agent storage)
   - Optional: Project-local `vc-agents/` (for special cases)

3. **Claude-Native Workflows**
   - Onboarding happens through Claude
   - Discovery through Claude
   - Deployment through Claude

4. **CLI as Convenience**
   - Quick checks without opening Claude
   - Scripting and automation
   - Power user workflows

---

## Proposed Architecture

### Primary: MCP Server

```
~/.cami/
├── config.yaml          # Global configuration
├── sources/             # Global agent sources
│   ├── lando-agents/   # Cloned from GitHub
│   ├── company-agents/ # Cloned from company repo
│   └── my-agents/      # Local custom agents
└── cami                 # Single binary (MCP + CLI)
```

**Configuration:**
```yaml
version: "1"
agent_sources:
  - name: "lando-agents"
    type: "local"
    path: "~/.cami/sources/lando-agents"
    priority: 100
    git:
      enabled: true
      remote: "git@github.com:lando-labs/lando-agents.git"

  - name: "my-agents"
    type: "local"
    path: "~/.cami/sources/my-agents"
    priority: 200
    git:
      enabled: false

deploy_locations:
  - name: "my-project"
    path: "/Users/lando/projects/my-project"
  - name: "client-project"
    path: "/Users/lando/clients/acme-app"
```

### Secondary: CLI Helper

Same binary, different mode:

```bash
# MCP mode (default)
$ cami-mcp
# Runs as MCP server on stdio

# CLI mode (when invoked as 'cami')
$ cami list
$ cami deploy frontend backend ~/projects/my-app
$ cami scan ~/projects/my-app
$ cami source add git@github.com:company/agents.git
```

### Installation

**Step 1: Install Binary**
```bash
# Homebrew (future)
brew install lando/tap/cami

# Or manual
curl -L https://github.com/lando/cami/releases/latest/download/cami-macos \
  -o ~/.cami/cami
chmod +x ~/.cami/cami
```

**Step 2: Configure MCP**
```json
// Claude Code MCP settings
{
  "mcpServers": {
    "cami": {
      "command": "~/.cami/cami",
      "args": ["--mcp"]
    }
  }
}
```

**Step 3: Optional CLI**
```bash
# Add to PATH for CLI convenience
ln -s ~/.cami/cami /usr/local/bin/cami
```

---

## Migration Strategy

### Phase 1: Dual Mode Binary (Week 1)

**Goal:** Single binary with MCP + CLI modes

**Tasks:**
1. Refactor cmd/cami/main.go to detect mode
   - Check if stdio is a pipe → MCP mode
   - Check args → CLI mode
2. Move MCP server code into main binary
3. Update MCP server to use global `~/.cami/sources/`
4. Keep backward compatibility with `vc-agents/`

**Deliverables:**
- Single `cami` binary
- Works as MCP server (default)
- Works as CLI (with args)

### Phase 2: Global Sources Default (Week 2)

**Goal:** Change defaults to global storage

**Tasks:**
1. Update `cami init` to create `~/.cami/sources/my-agents/`
2. Update config defaults to use `~/.cami/sources/`
3. Add migration helper for existing users
4. Update all documentation

**Deliverables:**
- Global sources by default
- Migration guide for existing users
- Updated README and docs

### Phase 3: Documentation Overhaul (Week 1)

**Goal:** Document MCP-first philosophy

**Tasks:**
1. Update README with MCP-first approach
2. Create "Installation" guide (MCP focus)
3. Create "CLI Reference" (secondary docs)
4. Update CLAUDE.md with new mental model

**Deliverables:**
- MCP-first README
- Clear installation instructions
- Updated user guides

### Phase 4: Optional CLI Streamlining (Week 1)

**Goal:** Simplify CLI to essential commands only

**Tasks:**
1. Audit CLI commands for necessity
2. Remove or consolidate rarely-used commands
3. Focus CLI on quick checks and scripting
4. Update help text to suggest MCP for complex tasks

**Deliverables:**
- Leaner CLI surface
- Clear MCP vs CLI guidance

---

## Key Changes

### What Changes

| Aspect | Current | MCP-First |
|--------|---------|-----------|
| **Default Storage** | `./vc-agents/` per project | `~/.cami/sources/` global |
| **Primary Interface** | CLI commands | MCP tools via Claude |
| **Binary Count** | 2 binaries (cami, cami-mcp) | 1 binary (dual mode) |
| **Mental Model** | "CLI with MCP support" | "MCP with CLI helper" |
| **Onboarding** | "Run cami init" | "Ask Claude to help" |
| **Documentation** | CLI-focused | MCP-focused |

### What Stays the Same

- All internal packages (`agent`, `config`, `deploy`, `docs`, `discovery`)
- Configuration format (backward compatible)
- MCP tool functionality (12 tools)
- Core workflows (deploy, scan, update)
- Test suite
- CI/CD

### Backward Compatibility

**Support for Old Workflows:**
```bash
# Still works: project-local vc-agents
$ cd my-project
$ cami init  # Creates ./vc-agents/ if user prefers

# Config can have both global and local sources
agent_sources:
  - name: "global-lando"
    path: "~/.cami/sources/lando-agents"
    priority: 100
  - name: "project-local"
    path: "./vc-agents/my-agents"
    priority: 200
```

**Migration Path:**
```bash
# Helper command to migrate existing projects
$ cami migrate-to-global

# Shows what would change
# Offers to move vc-agents to ~/.cami/sources/
# Updates config
```

---

## Benefits of MCP-First

### For Users

1. **Simpler Mental Model**
   - "CAMI is an MCP server that helps Claude manage agents"
   - No confusion about where to run commands

2. **Better Discovery**
   - Claude guides you through setup
   - Personalized recommendations via `onboard` tool
   - No need to remember commands

3. **Global Sources**
   - Agents available everywhere
   - No per-project setup overhead
   - Easier to maintain single source of truth

4. **Natural Workflows**
   - "Claude, add the frontend agent to this project"
   - "Claude, what agents do I have?"
   - "Claude, update my agent sources"

### For Developers

1. **Less Code to Maintain**
   - One primary interface (MCP)
   - CLI is thin wrapper
   - No duplicate logic

2. **Clearer Architecture**
   - MCP server is the core
   - Internal packages serve MCP
   - CLI is auxiliary

3. **Better Testing**
   - Test MCP tools primarily
   - CLI tests are minimal

4. **Easier Evolution**
   - New features go in MCP first
   - CLI follows if needed

---

## User Journey Comparison

### Current (CLI-First)

```
1. User discovers CAMI
2. Installs binary via GitHub releases
3. Reads README → "Run cami init"
4. Runs `cami init` in project
5. Creates ./vc-agents/my-agents/
6. Runs `cami source add ...` to get agents
7. Runs `cami list` to see agents
8. Runs `cami deploy ...` to deploy
9. Maybe discovers MCP later
```

### Proposed (MCP-First)

```
1. User discovers CAMI
2. Installs binary via Homebrew/releases
3. Configures MCP in Claude Code settings
4. Opens Claude Code
5. Claude: "I see you have CAMI! Let me help you set it up."
6. Claude uses `onboard` tool → personalized guidance
7. User: "Add the official agents"
8. Claude uses `add_source` → clones lando-agents
9. User: "Deploy frontend and backend"
10. Claude uses `deploy_agents` → agents deployed
11. User: "Update my CLAUDE.md"
12. Claude uses `update_claude_md` → documentation updated
```

**Natural. Guided. No commands to remember.**

---

## Risk Assessment

### Low Risk

- ✅ Internal packages don't change
- ✅ MCP tools already work well
- ✅ Configuration format stays same
- ✅ Can maintain backward compatibility

### Medium Risk

- ⚠️ User confusion during transition
- ⚠️ Documentation overhaul needed
- ⚠️ Migration path must be smooth

### Mitigation

1. **Clear Migration Guide**
   - Document old vs new approach
   - Provide migration helper tool
   - Show both patterns side-by-side

2. **Gradual Transition**
   - Phase 1: Make both work equally well
   - Phase 2: Document MCP as preferred
   - Phase 3: Soft-deprecate CLI-first workflow

3. **User Communication**
   - Blog post explaining philosophy
   - Update README with clear mental model
   - Highlight benefits of MCP-first

---

## Success Metrics

### User Adoption

- 80%+ of new users install via MCP
- 90%+ of interactions happen via MCP tools
- <5% support questions about "which mode to use"

### Code Quality

- Reduced CLI code by 50%
- 90%+ test coverage maintained
- Single binary < 10MB

### Documentation

- MCP-first README gets 90%+ positive feedback
- Installation guide has <10% bounce rate
- Support requests decrease 30%

---

## Timeline

| Week | Phase | Deliverables |
|------|-------|--------------|
| 1 | Dual Mode Binary | Single binary with MCP + CLI |
| 2 | Global Sources | Default to ~/.cami/sources/ |
| 3 | Documentation | MCP-first docs, migration guide |
| 4 | CLI Streamlining | Lean CLI, clear separation |

**Total:** 4 weeks to complete MCP-first transition

---

## Decision Points

### Should We Do This?

**Arguments For:**
- ✅ Aligns with actual usage (90% MCP)
- ✅ Simpler mental model for users
- ✅ Less code to maintain
- ✅ Better Claude Code integration
- ✅ Global sources more convenient

**Arguments Against:**
- ❌ Requires user education/migration
- ❌ Changes established workflow
- ❌ Documentation overhaul needed
- ❌ May confuse existing users

### Questions to Answer

1. **Do we want to maintain two equal interfaces?**
   - No → Go MCP-first
   - Yes → Keep status quo

2. **Is global storage better than per-project?**
   - For most users: Yes
   - For some workflows: Project-local needed
   - Solution: Support both, default to global

3. **Is now the right time?**
   - Pre-1.0: Yes, easier to change
   - Post-launch: Harder but still doable
   - Current: Pre-1.0, perfect timing

---

## Recommendation

**YES, transition to MCP-first architecture.**

**Reasoning:**
1. We're pre-1.0 (easier to change now)
2. Usage data shows MCP is primary interface
3. Simpler for users long-term
4. Less code to maintain
5. Better aligns with Claude Code ecosystem

**Approach:**
- Gradual, phased migration (4 weeks)
- Maintain backward compatibility
- Clear documentation
- Migration helper tools

**Next Steps:**
1. Get stakeholder buy-in on philosophy
2. Review and approve this plan
3. Begin Phase 1: Dual mode binary
4. Communicate changes to community

---

## Appendix: Examples

### Example 1: First-Time User

**MCP-First:**
```
User: "I want to start using CAMI"

Claude: *uses mcp__cami__onboard*
"Welcome! CAMI isn't configured yet. Let me help you set it up.
I'll add the official Lando agent library with 29 professional agents."

Claude: *uses mcp__cami__add_source*
"✓ Added lando-agents (29 agents available)"

User: "Deploy frontend and backend to this project"

Claude: *uses mcp__cami__deploy_agents*
"✓ Deployed frontend (v1.1.0)
 ✓ Deployed backend (v1.1.0)"

Claude: *uses mcp__cami__update_claude_md*
"✓ Updated CLAUDE.md with agent documentation"
```

**Natural, guided, no commands needed.**

### Example 2: Power User

**CLI Still Available:**
```bash
# Quick check
$ cami list | grep frontend
frontend (v1.1.0) - Use this agent when building...

# Scripted deployment
$ for project in ~/projects/*; do
    cami deploy frontend backend "$project"
  done

# Git workflow
$ cd ~/.cami/sources/my-agents
$ vim devops-specialist.md
$ git add . && git commit -m "Update devops agent"
$ cami deploy devops-specialist ~/active-project
```

**CLI remains useful for automation.**

---

## Conclusion

MCP-first is the right direction for CAMI. It aligns with actual usage, simplifies the user experience, and reduces maintenance burden. The transition is low-risk with proper planning, and we're at the perfect time (pre-1.0) to make this change.

**Recommended Action:** Approve plan and proceed with Phase 1.
