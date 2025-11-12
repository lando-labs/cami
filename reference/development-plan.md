# CAMI Development Plan (Simplified Architecture)

**Version:** 1.0
**Date:** 2025-11-10
**Target:** v1.0.0 Open Source Release

---

## Overview

Development plan based on simplified architecture:
- Minimal CLI (init, help, version)
- MCP-first design (Claude guides everything)
- Workspace model (vc-agents/ with multiple sources)
- Git basics only (clone, pull, status)
- Agent-architect as core value prop

---

## Pre-Development: Agent Extraction

### Current State
```
cami/
├── vc-agents/
│   ├── core/
│   │   ├── architect.md
│   │   ├── backend.md
│   │   ├── frontend.md
│   │   ├── mobile-native.md
│   │   └── qa.md
│   ├── specialized/ (11 agents)
│   ├── infrastructure/ (5 agents)
│   ├── integration/ (2 agents)
│   ├── design/ (4 agents)
│   └── meta/
│       └── agent-architect.md    ← Keep in CAMI repo
```

### Target State
```
lando-labs/cami/
├── .claude/
│   └── agents/
│       └── agent-architect.md    ← Moved here, shipped with CAMI
├── vc-agents/
│   ├── .gitignore               ← **/  (ignore all subdirs)
│   ├── README.md                ← Explains workspace
│   └── my-agents/               ← Default workspace (empty initially)
│       └── .gitkeep

lando-labs/lando-agents/         ← NEW REPO
├── core/ (5 agents)
├── specialized/ (11 agents)
├── infrastructure/ (5 agents)
├── integration/ (2 agents)
└── design/ (4 agents)
```

### Extraction Steps

**Phase -1: Extract Agents (Pre-Development)**

**Timeline:** 1 day

**Tasks:**
1. ✅ Create `lando-labs/lando-agents` repo on GitHub
2. ✅ Move all agents EXCEPT agent-architect
   ```bash
   # In new lando-agents repo
   mkdir -p core specialized infrastructure integration design

   # Copy from cami/vc-agents (preserving structure)
   cp -r ../cami/vc-agents/core/* core/
   cp -r ../cami/vc-agents/specialized/* specialized/
   # ... etc (exclude meta/)

   git add .
   git commit -m "Initial agent library (29 agents)"
   git push
   ```

3. ✅ In CAMI repo, move agent-architect
   ```bash
   cd cami
   git mv vc-agents/meta/agent-architect.md .claude/agents/
   git commit -m "Move agent-architect to .claude/agents (ships with CAMI)"
   ```

4. ✅ Clean up vc-agents in CAMI repo
   ```bash
   # Remove all agents
   git rm -r vc-agents/core
   git rm -r vc-agents/specialized
   git rm -r vc-agents/infrastructure
   git rm -r vc-agents/integration
   git rm -r vc-agents/design
   git rm -r vc-agents/meta

   # Create clean workspace structure
   mkdir -p vc-agents/my-agents
   touch vc-agents/my-agents/.gitkeep

   # Add .gitignore
   echo "**/" > vc-agents/.gitignore
   echo "!.gitignore" >> vc-agents/.gitignore
   echo "!README.md" >> vc-agents/.gitignore

   git add vc-agents/
   git commit -m "Reset vc-agents to clean workspace structure"
   ```

5. ✅ Create `vc-agents/README.md`
   ```markdown
   # Agent Workspace

   This directory houses multiple agent sources.

   ## Structure

   Each subdirectory can be:
   - Local workspace (e.g., `my-agents/`)
   - Cloned repository (e.g., `lando-agents/`, `company-agents/`)

   ## Priority

   When multiple sources have the same agent, highest priority wins:
   - 200: Local experiments
   - 150: Team/company sources
   - 100: Community sources

   ## Adding Sources

   Use CAMI MCP tools via Claude Code:
   ```
   @claude add this agent source: git@github.com:company/agents.git
   ```

   Or manually:
   ```bash
   cd vc-agents
   git clone git@github.com:company/agents.git company
   ```

   Then update `~/.cami/config.yaml`.
   ```

6. ✅ Update CAMI CLAUDE.md
   ```markdown
   ## Agent-Architect

   CAMI ships with agent-architect in `.claude/agents/agent-architect.md`.
   This meta-agent creates sophisticated agents via research.

   The agent library has been extracted to `lando-labs/lando-agents`
   and is no longer bundled with CAMI. Users can add it as a source
   if desired.
   ```

**Deliverables:**
- ✅ `lando-labs/lando-agents` repo (29 agents)
- ✅ CAMI ships only agent-architect
- ✅ Clean vc-agents workspace structure

---

## Phase 0: Testing & Open Source Prep

**Timeline:** 2-3 weeks

**Goal:** Get CAMI to 80% test coverage and prepare for open source

### Week 1: Critical Path Tests

**Packages (target 90-95% coverage):**
- `internal/agent/` - Agent loading, parsing
- `internal/deploy/` - File operations, deployment
- `internal/docs/` - CLAUDE.md manipulation

**Tasks:**
1. Set up testing infrastructure
   ```bash
   go get github.com/stretchr/testify
   ```

2. Write unit tests for `internal/agent/`
   - `LoadAgents()` with nested folders
   - Agent frontmatter parsing
   - Category extraction
   - Priority handling

3. Write unit tests for `internal/deploy/`
   - Agent deployment to projects
   - File copying
   - Overwrite detection
   - Error handling

4. Write unit tests for `internal/docs/`
   - CLAUDE.md reading/writing
   - Agent section management
   - Markdown generation

5. Set up CI/CD (GitHub Actions)
   ```yaml
   name: Test
   on: [push, pull_request]
   jobs:
     test:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v3
         - uses: actions/setup-go@v4
           with:
             go-version: '1.21'
         - run: go test ./... -cover
   ```

**Deliverables:**
- 90%+ coverage on critical packages
- CI/CD running on every push

### Week 2: Configuration & Discovery Tests

**Packages (target 85% coverage):**
- `internal/config/` - Configuration management
- `internal/discovery/` - Agent scanning

**Tasks:**
1. Test `internal/config/`
   - Config loading/saving
   - Source management
   - Priority handling
   - Migration from old format

2. Test `internal/discovery/`
   - Agent scanning in workspace
   - Multi-source discovery
   - Priority-based deduplication
   - Git status detection

3. Integration tests
   - End-to-end workflows
   - Multi-source scenarios
   - Deployment workflows

**Deliverables:**
- 85%+ coverage on config/discovery
- Integration test suite

### Week 3: Open Source Prep

**Tasks:**

1. **License**
   - Add MIT or Apache 2.0 license
   - License headers in all Go files

2. **Documentation**
   - README.md (getting started)
   - CONTRIBUTING.md
   - CODE_OF_CONDUCT.md
   - SECURITY.md

3. **Repository cleanup**
   - Remove sensitive data
   - Clean up .gitignore
   - Verify no secrets in history

4. **GitHub setup**
   - Issue templates
   - PR template
   - GitHub Actions badges
   - Topics/tags

5. **Release preparation**
   - Semantic versioning
   - CHANGELOG.md
   - Release notes template

**Deliverables:**
- 80%+ overall test coverage
- Complete documentation
- Open source ready repository

---

## Phase 1: Workspace Architecture (Minimal CLI)

**Timeline:** 1 week

**Goal:** Implement minimal CLI with workspace initialization

### Tasks

1. **Update `internal/config/` for workspace model**
   ```go
   type Config struct {
       Version      string         `yaml:"version"`
       AgentSources []AgentSource  `yaml:"agent_sources"`
       Locations    []Location     `yaml:"deploy_locations"`
   }

   type AgentSource struct {
       Name     string  `yaml:"name"`
       Type     string  `yaml:"type"`      // "local"
       Path     string  `yaml:"path"`
       Priority int     `yaml:"priority"`
       Git      *GitConfig `yaml:"git,omitempty"`
   }

   type GitConfig struct {
       Enabled bool   `yaml:"enabled"`
       Remote  string `yaml:"remote,omitempty"`
   }
   ```

2. **Implement `cami init` (minimal)**
   ```go
   // cmd/cami/init.go
   func InitCommand() error {
       // Ask single question: storage location
       fmt.Println("Welcome to CAMI! 🚀")
       fmt.Println("\nWhere should we store your agents?")
       fmt.Println("  1. ./vc-agents (default)")
       fmt.Println("  2. Custom path")

       choice := promptChoice(1, 2)

       var vcAgentsPath string
       if choice == 1 {
           vcAgentsPath = "./vc-agents"
       } else {
           vcAgentsPath = promptString("Path: ")
       }

       // Create workspace structure
       myAgentsPath := filepath.Join(vcAgentsPath, "my-agents")
       os.MkdirAll(myAgentsPath, 0755)

       // Create config
       config := Config{
           Version: "1",
           AgentSources: []AgentSource{
               {
                   Name:     "my-agents",
                   Type:     "local",
                   Path:     myAgentsPath,
                   Priority: 200,
                   Git:      &GitConfig{Enabled: false},
               },
           },
       }

       SaveConfig(config)

       fmt.Println("\n✓ Created", myAgentsPath)
       fmt.Println("✓ Configuration saved")
       fmt.Println("\nNext: Start Claude and ask for onboarding help")

       return nil
   }
   ```

3. **Update `internal/agent/` for multi-source**
   ```go
   func LoadAgents(config *Config) ([]*Agent, error) {
       var allAgents []*Agent

       for _, source := range config.AgentSources {
           agents, err := loadAgentsFromSource(source)
           if err != nil {
               return nil, err
           }
           allAgents = append(allAgents, agents...)
       }

       // Deduplicate by priority
       return deduplicateByPriority(allAgents), nil
   }

   func deduplicateByPriority(agents []*Agent) []*Agent {
       // Sort by priority (descending)
       sort.Slice(agents, func(i, j int) bool {
           return agents[i].Priority > agents[j].Priority
       })

       // Keep highest priority for each name
       seen := make(map[string]bool)
       result := []*Agent{}

       for _, agent := range agents {
           if !seen[agent.Name] {
               result = append(result, agent)
               seen[agent.Name] = true
           }
       }

       return result
   }
   ```

4. **Tests**
   - Test init command
   - Test multi-source loading
   - Test priority deduplication

**Deliverables:**
- ✅ `cami init` working (minimal)
- ✅ Config format updated
- ✅ Multi-source agent loading
- ✅ Priority-based deduplication

---

## Phase 2: MCP Onboarding & Source Management

**Timeline:** 2 weeks

**Goal:** Add MCP tools for onboarding, source management, and agent creation interviews

### Week 1: Core MCP Tools

**Tasks:**

1. **`onboard` tool**
   ```go
   // cmd/cami-mcp/onboard.go
   func HandleOnboard(params map[string]interface{}) (string, error) {
       config := LoadConfig()

       state := OnboardingState{
           InitComplete: config != nil,
           HasAgentArch: checkAgentArchitect(),
           SourceCount:  len(config.AgentSources),
       }

       // Return guided text based on state
       if !state.InitComplete {
           return "Please run `cami init` first...", nil
       }

       if state.SourceCount == 1 {
           // Only my-agents, explain workflow
           return generateOnboardingText(state), nil
       }

       // Has sources, provide quick reference
       return generateQuickReference(state), nil
   }
   ```

2. **`source_add` tool**
   ```go
   func HandleSourceAdd(params map[string]interface{}) (string, error) {
       url := params["url"].(string)
       name := params["name"].(string)  // optional
       priority := params["priority"].(int)  // optional, default 150

       if name == "" {
           name = deriveNameFromURL(url)
       }

       // Clone to vc-agents/
       vcAgentsDir := getVCAgentsDir()
       targetPath := filepath.Join(vcAgentsDir, name)

       cmd := exec.Command("git", "clone", url, targetPath)
       if err := cmd.Run(); err != nil {
           return "", err
       }

       // Update config
       config := LoadConfig()
       config.AgentSources = append(config.AgentSources, AgentSource{
           Name:     name,
           Type:     "local",
           Path:     targetPath,
           Priority: priority,
           Git: &GitConfig{
               Enabled: true,
               Remote:  url,
           },
       })
       SaveConfig(config)

       // Count agents
       agents, _ := loadAgentsFromPath(targetPath)

       return fmt.Sprintf("✓ Cloned %s to vc-agents/%s\n✓ Added source with priority %d\n✓ Found %d agents",
           name, name, priority, len(agents)), nil
   }
   ```

3. **`source_list` tool**
   ```go
   func HandleSourceList(params map[string]interface{}) (string, error) {
       config := LoadConfig()

       var result strings.Builder
       result.WriteString("Agent Sources:\n\n")

       for _, source := range config.AgentSources {
           agents, _ := loadAgentsFromPath(source.Path)

           result.WriteString(fmt.Sprintf("  %s (priority %d)\n", source.Name, source.Priority))
           result.WriteString(fmt.Sprintf("    Path: %s\n", source.Path))
           result.WriteString(fmt.Sprintf("    Agents: %d\n", len(agents)))

           if source.Git != nil && source.Git.Enabled {
               result.WriteString(fmt.Sprintf("    Git: %s\n", source.Git.Remote))
           }

           result.WriteString("\n")
       }

       return result.String(), nil
   }
   ```

4. **`source_update` tool**
   ```go
   func HandleSourceUpdate(params map[string]interface{}) (string, error) {
       config := LoadConfig()
       sourceName := params["source"].(string)  // optional

       var results []string

       for _, source := range config.AgentSources {
           if sourceName != "" && source.Name != sourceName {
               continue
           }

           if source.Git == nil || !source.Git.Enabled {
               results = append(results, fmt.Sprintf("⊗ %s: no git remote", source.Name))
               continue
           }

           // git pull
           cmd := exec.Command("git", "-C", source.Path, "pull")
           output, err := cmd.CombinedOutput()

           if err != nil {
               results = append(results, fmt.Sprintf("✗ %s: %s", source.Name, err))
           } else if strings.Contains(string(output), "Already up to date") {
               results = append(results, fmt.Sprintf("✓ %s: up to date", source.Name))
           } else {
               results = append(results, fmt.Sprintf("✓ %s: updated", source.Name))
           }
       }

       return strings.Join(results, "\n"), nil
   }
   ```

5. **`source_status` tool**
   ```go
   func HandleSourceStatus(params map[string]interface{}) (string, error) {
       config := LoadConfig()

       var result strings.Builder
       result.WriteString("Agent Source Status:\n\n")

       for _, source := range config.AgentSources {
           result.WriteString(fmt.Sprintf("  %s\n", source.Name))

           if source.Git == nil || !source.Git.Enabled {
               result.WriteString("    Git: not enabled\n\n")
               continue
           }

           // git status --porcelain
           cmd := exec.Command("git", "-C", source.Path, "status", "--porcelain")
           output, _ := cmd.Output()

           if len(output) == 0 {
               result.WriteString("    Git: ✓ clean\n")
           } else {
               lines := strings.Split(string(output), "\n")
               result.WriteString(fmt.Sprintf("    Git: ⚠ %d uncommitted changes\n", len(lines)-1))
               for _, line := range lines[:min(3, len(lines))] {
                   if line != "" {
                       result.WriteString(fmt.Sprintf("      %s\n", line))
                   }
               }
           }

           result.WriteString("\n")
       }

       return result.String(), nil
   }
   ```

**Deliverables:**
- ✅ 5 new MCP tools implemented
- ✅ Integration tests for MCP tools
- ✅ Error handling

### Week 2: Agent Creation Interview System

**Tasks:**

1. **Implement `agent_create` MCP tool**
   ```go
   // Detects domain from request
   // Loads interview template
   // Returns structured questions
   // Guides agent-architect through discovery
   ```

2. **Create interview templates**
   - Infrastructure (Terraform, K8s, Docker)
   - Backend (APIs, databases)
   - Frontend (React, Vue, components)
   - Data (ETL, analytics)
   - AI/ML (models, training)
   - DevOps (CI/CD, deployments)

3. **Update agent-architect.md**
   Add interview framework:
   ```markdown
   ## Agent Creation Interview Process

   When creating agents, conduct discovery interview:
   1. Purpose & use case
   2. Philosophy & approach
   3. Constraints & boundaries
   4. Workflow patterns
   5. Context & environment

   Use `agent_create` MCP tool to guide interview.

   [Detailed interview framework...]
   ```

4. **Workspace awareness**
   - Check available sources
   - Ask where to create (if multiple sources)
   - Default to my-agents
   - Place in appropriate category

5. **Test interview flows**
   - Test all domain templates
   - Verify agent quality with vs without interview
   - Test skip interview option

**Deliverables:**
- ✅ `agent_create` MCP tool
- ✅ 6+ interview templates
- ✅ Agent-architect interview-aware
- ✅ Agent-architect workspace-aware
- ✅ Tests for guided creation

**See:** [agent-creation-interview-system.md](./agent-creation-interview-system.md)

---

## Phase 3: Documentation & Polish

**Timeline:** 1 week

**Goal:** Complete documentation and user experience polish

### Tasks

1. **Update README.md**
   - Getting started (simplified)
   - MCP configuration
   - Quick reference
   - Link to docs

2. **Create CONTRIBUTING.md**
   ```markdown
   # Contributing to CAMI

   ## For Agent Contributors

   Agents live in `lando-labs/lando-agents`, not this repo.

   To contribute agents:
   1. Fork `lando-labs/lando-agents`
   2. Create agents using agent-architect
   3. Test with CAMI
   4. Submit PR

   ## For CAMI Tool Contributors

   [Standard Go contribution guide]
   ```

3. **Create user guides**
   - Onboarding guide (MCP-first)
   - Common workflows
   - Troubleshooting

4. **CLI polish**
   - Better error messages
   - Help text improvements
   - Version command

5. **MCP polish**
   - Consistent response formatting
   - Better error messages
   - Loading indicators (where applicable)

**Deliverables:**
- ✅ Complete documentation
- ✅ Polished user experience
- ✅ Ready for v1.0.0

---

## Phase 4: Beta Testing & Refinement

**Timeline:** 1-2 weeks

**Goal:** Test with real users, gather feedback, fix issues

### Tasks

1. **Internal testing**
   - Test all workflows
   - Verify edge cases
   - Check error handling

2. **External beta**
   - Invite trusted users
   - Gather feedback
   - Document pain points

3. **Refinements**
   - Fix bugs
   - Improve UX based on feedback
   - Update docs

4. **Performance**
   - Profile hot paths
   - Optimize if needed
   - Test with large agent counts

**Deliverables:**
- ✅ Beta tested with 5-10 users
- ✅ Major issues resolved
- ✅ Ready for public release

---

## Phase 5: Launch v1.0.0

**Timeline:** 1 week

**Goal:** Public open source release

### Tasks

1. **Final QA**
   - Run full test suite
   - Manual testing of all flows
   - Security audit

2. **Release preparation**
   - Tag v1.0.0
   - Generate CHANGELOG
   - Create release notes

3. **Repository finalization**
   - Make repo public
   - Add GitHub topics
   - Configure settings

4. **Launch**
   - Publish release
   - Announce (Twitter, Reddit, etc.)
   - Monitor for issues

5. **Post-launch**
   - Respond to issues
   - Accept PRs
   - Community engagement

**Deliverables:**
- 🎉 CAMI v1.0.0 open source
- 🎉 lando-labs/lando-agents available
- 🎉 Community ready

---

## Timeline Summary

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| **Phase -1**: Agent Extraction | 1 day | lando-agents repo, clean CAMI |
| **Phase 0**: Testing & Prep | 2-3 weeks | 80% coverage, open source ready |
| **Phase 1**: Minimal CLI | 1 week | `cami init`, workspace model |
| **Phase 2**: MCP Tools | 1-2 weeks | 5 MCP tools, agent-architect integration |
| **Phase 3**: Documentation | 1 week | Complete docs, polished UX |
| **Phase 4**: Beta Testing | 1-2 weeks | User tested, refined |
| **Phase 5**: Launch | 1 week | v1.0.0 public release |
| **Total** | **7-10 weeks** | Open source CAMI + lando-agents |

---

## Success Criteria

### Before v1.0.0 Launch

**Must Have:**
- ✅ 80%+ test coverage
- ✅ All MCP tools working
- ✅ Agent-architect workspace-aware
- ✅ Complete documentation
- ✅ No critical bugs
- ✅ Tested by 5+ users
- ✅ Clean Git history (no secrets)
- ✅ License in place (MIT/Apache)

**Nice to Have:**
- ⭐ Homebrew formula
- ⭐ Pre-built binaries
- ⭐ Video walkthrough
- ⭐ Blog post

### Post-Launch Metrics

**First Month:**
- 100+ GitHub stars
- 10+ agent contributions to lando-agents
- 50+ projects using CAMI
- Active community (GitHub Discussions)

**First Quarter:**
- 500+ stars
- 50+ agents in lando-agents
- 200+ projects
- Community sources emerging

---

## Risk Management

### High Risk Items

1. **MCP complexity**
   - Risk: MCP tools harder to implement than expected
   - Mitigation: Start with simple tools, iterate
   - Fallback: Keep existing MCP tools, add new ones incrementally

2. **Testing timeline**
   - Risk: 2-3 weeks not enough for 80% coverage
   - Mitigation: Prioritize critical paths, accept 70% if needed
   - Fallback: Launch with 70%, hit 80% post-launch

3. **Beta feedback**
   - Risk: Major UX issues found late
   - Mitigation: Internal testing first, small beta group
   - Fallback: Extend beta period, delay launch

### Medium Risk Items

1. **Agent-architect changes**
   - Risk: Workspace awareness harder than expected
   - Mitigation: Start simple, iterate based on usage
   - Fallback: Manual workflow (user specifies path)

2. **Documentation**
   - Risk: Hard to explain new workflow
   - Mitigation: Video walkthrough, clear examples
   - Fallback: Better docs post-launch based on questions

---

## Open Questions

1. **License:** MIT or Apache 2.0?
2. **Homebrew:** Package immediately or wait for adoption?
3. **Pre-built binaries:** Provide or let users build?
4. **Video walkthrough:** Create before or after launch?
5. **Launch venue:** Just GitHub or also Reddit/HN?

---

## Next Steps

1. ✅ Review this plan
2. ⬜ Decide on open questions
3. ⬜ Execute Phase -1 (agent extraction)
4. ⬜ Create GitHub issues for all phases
5. ⬜ Start Phase 0 (testing)
