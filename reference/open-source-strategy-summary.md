# CAMI Open Source Strategy - Executive Summary

**Version:** 2.0
**Date:** 2025-11-10
**Full Document:** [open-source-strategy.md](./open-source-strategy.md)

## TL;DR

Comprehensive 17-week plan to take CAMI from internal tool to successful open source project. Key decisions:

- **Philosophy:** Agent-architect is the core value prop, not an advanced feature
- **Contribution Model:** One-way pull (automated) + Git push (standard workflows)
- **Testing:** 80% coverage target, Go testing + testify, 6 weeks to production-ready
- **Agents:** Separate `lando-labs/lando-agents` repo, independent versioning
- **Sources:** Multi-source architecture (GitHub, Git, local), priority-based conflicts
- **Workflows:** Support 4 user types (library dev, consumer, team, polyglot)
- **CLI:** Major UX overhaul with interactive modes, better errors, shell completion
- **Launch:** v1.0.0 in Week 17, registry and marketplace in v2.0.0

## Core Philosophy

> **CAMI facilitates consumption, not contribution.**

**What CAMI does:**
- ✅ **Pull:** Effortless discovery, deployment, updates (automated)
- ✅ **Guidance:** Documentation, status, diff (helpful)

**What CAMI doesn't do:**
- ❌ **Push:** Use standard Git (`git push`, `gh pr create`)
- ❌ **Auth:** Use Git credentials
- ❌ **Conflicts:** Use Git merge/rebase

**Why?**
- Avoids reinventing Git/GitHub
- Users learn transferable skills
- Fewer edge cases to maintain
- Works with any Git platform

## Key Architectural Decisions

### 1. Repository Strategy: SEPARATE REPOS

**Decision:** Split into `lando-labs/cami` (tool) + `lando-labs/lando-agents` (content)

**Why:**
- Clear separation of concerns
- Independent versioning (tool v1.0.0, agents v2.3.0)
- Easier contributions (add agents without touching Go code)
- Different audiences (tool users vs agent creators)

**Structure:**
```
lando-labs/cami           → CLI/TUI/MCP tool (Go code)
lando-labs/lando-agents    → Official agent library (markdown files)
```

### 2. Remote Agent Sources: MULTI-SOURCE PRIORITY

**Decision:** Support multiple agent sources with priority-based conflict resolution

**Config format:**
```yaml
agent_sources:
  - name: local-dev          # Highest priority (experiments)
    type: local
    path: ~/my-agents
    priority: 200

  - name: acme-internal      # Team agents
    type: git
    url: git@github.com:acme-corp/agents.git
    priority: 150

  - name: official           # Lando Labs (baseline)
    type: github
    repo: lando-labs/lando-agents
    priority: 100

  - name: community          # Lowest priority
    type: github
    repo: awesome-claude/agents
    priority: 50
```

**Why:**
- Enables official + community + private agents
- Priority resolves name conflicts predictably
- Local dev overrides everything (good for testing)
- Easy to enable/disable sources

### 3. Testing Strategy: PRACTICAL 80%

**Decision:** Target 80% overall coverage with pragmatic focus on critical paths

**Coverage by package:**
- `internal/agent`: 95% (critical - data integrity)
- `internal/deploy`: 90% (critical - file operations)
- `internal/docs`: 90% (critical - CLAUDE.md manipulation)
- `internal/config`: 85% (high - configuration)
- `internal/discovery`: 85% (high - version logic)
- `internal/cli`: 70% (medium - glue code)
- `internal/tui`: 50% (low - UI, hard to test)
- `cmd/*`: 40% (low - integration tests cover)

**Tools:**
- Go standard `testing` package
- `testify/assert` for readable assertions
- `testify/mock` for mocking
- GitHub Actions CI/CD
- Codecov for coverage reporting

### 4. CLI UX: INTERACTIVE + MODERN

**Decision:** Major UX overhaul with interactive fallbacks and better feedback

**New commands:**
- `cami init` - Interactive setup wizard
- `cami search <query>` - Search agents
- `cami show <agent>` - Agent details
- `cami sources` - Source management
- `cami update` - Update agents (replaces update-docs)

**Improvements:**
- Short flags: `-a` (agents), `-l` (location), `-o` (output)
- Interactive mode: Omit required flags → CLI prompts
- Colors: ✓ Green, ✗ Red, ⚠ Yellow, ℹ Blue
- Progress: Spinners for long operations
- Completion: Bash, Zsh, Fish shell completion
- Better errors: Clear problem + actionable solutions

### 5. Workflows: SUPPORT ALL TYPES

**Decision:** Enable 4 distinct workflows with clear documentation

**Workflow A: Agent Library Developer**
- Works inside `lando-agents` repo
- CAMI installed globally
- Local source points to agent files
- **Use case:** Contributing official agents

**Workflow B: Agent Consumer** (DEFAULT for open source)
- CAMI installed: `go install github.com/lando-labs/cami/cmd/cami@latest`
- Agents pulled from GitHub: `lando-labs/lando-agents`
- Deploy to projects
- **Use case:** Using CAMI in projects

**Workflow C: Team Library**
- Private Git repo: `acme-corp/claude-agents`
- Company-specific agents
- Priority over official when conflicts
- **Use case:** Company internal agents

**Workflow D: Polyglot**
- Mix of local + team + official + community
- Priority-based resolution
- Explicit source selection when needed
- **Use case:** Power users, experimentation

### 6. Agent-Architect: CORE VALUE PROPOSITION

**Decision:** Open source agent-architect and make it the **recommended onboarding path**

**Philosophy:**
- ❌ Templates are weak
- ✅ Research-driven agent creation is powerful
- ✅ Agent-architect should be the **first thing** new users try

**Onboarding Flow:**
```bash
cami init
> Create agents on demand with agent-architect (recommended) ← Default!
> Use existing agents only
```

**Why:**
- Showcases CAMI's real power (research → sophisticated agents)
- Users create exactly what they need
- No settling for "close enough" from templates
- Empowers community to expand ecosystem

**Implementation:**
- Place in `meta/` category for organization
- Deploy during `cami init` if user selects agent creation
- Recommend running CAMI from source directory
- Standard Git workflow for contributing created agents

## Contribution Workflow

### The One-Way Street with a Roadmap

**CAMI's role:**
1. **Pull (automated):** `cami update` pulls latest agents from sources
2. **Status (helpful):** `cami status` shows Git state if available
3. **Diff (helpful):** `cami diff agent-name` compares local vs official
4. **Guidance (docs):** Point to Git/GitHub workflow

**User's role:**
1. **Create/edit:** Use agent-architect or edit directly
2. **Test:** Deploy with CAMI to project
3. **Contribute:** Standard Git workflow
   ```bash
   git add vc-agents/category/agent.md
   git commit -m "Add agent"
   git push origin main
   gh pr create --title "..." --body "..."
   ```

**Benefits:**
- ✅ Clean separation (CAMI = consume, Git = contribute)
- ✅ Standard tools (git, gh CLI)
- ✅ No edge cases (Git handles complexity)
- ✅ LLM-friendly (Claude already knows Git)
- ✅ Platform agnostic (GitHub, GitLab, Bitbucket)

**See:** [contribution-philosophy.md](./contribution-philosophy.md)

## Implementation Timeline

### Phase 0: Testing (Weeks 1-6)
**Goal:** Production-ready test coverage

- Week 1-2: Critical path tests (agent, deploy)
- Week 3: High priority (docs, config)
- Week 4: Medium priority (discovery, CLI)
- Week 5: Polish (TUI, documentation)
- Week 6: QA (cross-platform, performance)

**Deliverable:** 80%+ coverage, CI/CD working

### Phase 1: Repository Split (Weeks 7-8)
**Goal:** Separate tool from agents

- Create `lando-labs/lando-agents` repo
- Migrate `vc-agents/` content
- Agent validation CI/CD
- Update tool to use remote agents
- Migration tooling

**Deliverable:** Two repos, migration documented

### Phase 2: Multi-Source (Weeks 9-11)
**Goal:** Remote agent sources

- `internal/sources` package
- GitHub and local source types
- Cache management
- Priority-based conflicts
- Config migration (JSON → YAML)

**Deliverable:** v0.3.0-beta

### Phase 3: CLI UX (Weeks 12-14)
**Goal:** Modern CLI experience

- New commands (init, search, show)
- Interactive modes
- Colors and progress
- Shell completion
- Better errors

**Deliverable:** v0.4.0

### Phase 4: Git Sources (Weeks 15-16)
**Goal:** Full Git support

- Generic Git source type
- Authentication (SSH, token)
- Private repos
- Team workflows

**Deliverable:** v0.5.0

### Phase 5: Launch (Week 17)
**Goal:** Public release

- Final QA and security audit
- License (MIT or Apache 2.0)
- Community guidelines
- Launch announcement
- Package managers (Homebrew)

**Deliverable:** v1.0.0 - Open Source Release! 🎉

### Phase 6: Post-Launch (Ongoing)
**Goal:** Community and features

- **Short-term:** Bug fixes, first contributions
- **Medium-term:** Deployment tracking, batch operations
- **Long-term:** Registry (`agents.cami.dev`), marketplace, web UI

## Critical Success Factors

### Before Launch
1. **80%+ test coverage** - No compromise
2. **Clean repository split** - Tool vs agents separate
3. **Multi-source working** - Official + community + private
4. **Documentation complete** - README, guides, API docs
5. **Migration path clear** - v0.2.0 → v1.0.0
6. **Community ready** - Guidelines, templates, issue labels

### Launch Metrics
- 1000+ GitHub stars in first month
- 50+ community contributions in first quarter
- 100+ projects using CAMI agents
- Active community (Discussions/Discord)
- Homebrew installs >500/month

### Long-term Success
- Agent marketplace with ratings
- Registry at `agents.cami.dev`
- CI/CD integrations
- Enterprise features (teams, analytics)
- Ecosystem of community agents

## Quick Decision Reference

| Question | Answer | Rationale |
|----------|--------|-----------|
| Monorepo or separate? | **Separate repos** | Clean separation, easier contributions |
| Open source agent-architect? | **Yes, as meta-agent** | Empowerment + education |
| Testing framework? | **Go testing + testify** | Standard, fast, easy onboarding |
| Coverage target? | **80% overall** | Pragmatic, focus on critical paths |
| Config format? | **YAML** | Better comments, human-readable |
| Default workflow? | **Agent Consumer (B)** | Simplest for end users |
| Agent versioning? | **Semantic (SemVer)** | Clear compatibility |
| CLI style? | **Interactive + flags** | Flexible for all users |
| Authentication? | **SSH keys + tokens** | Standard Git methods |
| Conflict resolution? | **Priority-based** | Predictable, configurable |

## Next Actions

1. **Review this strategy** with team
2. **Start Phase 0 testing** immediately
3. **Set up project tracking** (GitHub Projects)
4. **Assign phase owners** for accountability
5. **Create testing backlog** (Week 1-2 focus)

## Questions to Resolve

1. **License:** MIT or Apache 2.0?
2. **Repository name:** Keep `lando-labs` or move to `cami-dev`?
3. **Registry domain:** `agents.cami.dev` available?
4. **Launch date:** Flexible or hard deadline?
5. **Team size:** Who's assigned to testing phase?

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Testing takes longer | Start immediately, prioritize critical paths |
| Community low quality | Quality gates, maintainer review, templates |
| Support burden | Clear docs, community guidelines, triaging process |
| Breaking changes | Semantic versioning, deprecation policy |
| Private use only | Launch marketing, showcase examples |
| Competitor forks | Strong community, frequent releases, engagement |

## Resources

- **Full Strategy:** [open-source-strategy.md](./open-source-strategy.md) (24,000 words)
- **Current CAMI:** v0.2.0 (manual testing, local agents)
- **Target Launch:** v1.0.0 (Week 17, fully tested, multi-source)
- **GitHub:** lando-labs/cami + lando-labs/lando-agents (future)

---

**This is not just a plan - it's a roadmap to building a category-defining tool.**

The architecture preserves what works (your "agent repo as working directory" workflow) while enabling what's needed (community contributions, private agents, ecosystem growth). It's ambitious but achievable with clear phases and practical decisions.

**Let's build the future of Claude Code agent management.** 🚀
