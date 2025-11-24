# CAMI Workflow Enhancement Plan V2

**Goal**: Transform CAMI into a guru-guided experience with robust workflows, agent classification, and intelligent onboarding.

**Philosophy**: Provide *guidance* not *prescription*. Meet users where they are technically, keep them engaged with vision over details.

**CRITICAL UPDATE**: This plan incorporates the three-class agent system with phase-weighted methodology and auto-classification.

---

## Agent Classification System (Core Foundation)

### Three Agent Classes

#### 1. Workflow Specialists
**Purpose**: Execute specific, user-defined workflows reliably
**Cognitive Model**: Single-purpose checklist executors
**Phase Weights**: Research (10-15%) → Execute (70-80%) → Validate (10-15%)

**Key Characteristics**:
- User provides the workflow/checklist during creation
- Agent validates workflow makes sense
- Agent executes that specific workflow reliably
- Keeps agents focused and file sizes manageable

**Examples**:
- `k8s-pod-checker`: Checks pod status with defined health check steps
- `jira-issue-updater`: Updates JIRA issues following specific workflow
- `deployment-to-staging`: Executes deployment checklist
- `component-builder`: Builds specific components with defined steps
- `api-health-checker`: Runs defined API health check sequence

**NOT Examples**:
- "Build any button" (too broad)
- "Handle all deployments" (too general)
- "Fix bugs" (no specific workflow)

#### 2. Technology Implementers
**Purpose**: Build complete capabilities in specific domains
**Cognitive Model**: Domain specialists who implement features
**Phase Weights**: Research (25-30%) → Execute (50-60%) → Validate (15-20%)

**Key Characteristics**:
- Domain experts (frontend, backend, database, etc.)
- Build complete features/screens/APIs
- Coordinate multiple concepts
- Make implementation decisions

**Examples**:
- `frontend`: Build complete UI screens and components
- `backend`: Build full APIs with multiple endpoints
- `database`: Design schemas, migrations, queries
- `auth-system`: Complete authentication implementation
- `payment-integration`: Full payment processing

#### 3. Strategic Planners
**Purpose**: Architect systems, research technologies, optimize at scale
**Cognitive Model**: System architects and researchers
**Phase Weights**: Research (40-50%) → Execute (25-35%) → Validate (20-25%)

**Key Characteristics**:
- High-level planning and architecture
- Technology evaluation and research
- Cross-system integration
- Strategic decision-making

**Examples**:
- `architect`: System design, API contracts, integration patterns
- `researcher`: Deep research, technology evaluation, feasibility studies
- `security`: Cross-system security analysis and recommendations
- `performance`: System-wide performance optimization strategies
- `devops`: Infrastructure planning, deployment strategies

---

## Phase-Weighted Methodology

All agents have all three phases, but weighted differently:

### Research/Analyze Phase
- Understanding the problem
- Gathering context
- Reading relevant code/docs
- Identifying patterns

### Build/Execute Phase
- Core implementation work
- Writing code
- Running commands
- Making changes

### Validate/Follow-up Phase
- Testing implementation
- Verifying success criteria
- Cleanup and documentation
- Error handling

**Key Insight**: Phases are modular but integrate agent traits. User can request focus shift:
- "Research this first, then build" → Temporarily boost research weight
- "Just build it quickly" → Boost execution weight
- "Make sure it's thoroughly tested" → Boost validation weight

---

## Phase 1: Foundation Fixes (Immediate)

### 1.1 Fix Onboarding Detection

**Problem**: Onboarding reports "everything configured" on fresh install because default `my-agents` source exists but is empty.

**Solution**: Detect "fresh install" vs "actually configured"

**Implementation**:
```go
// In onboard tool
type OnboardingState struct {
    ConfigExists     bool
    IsFreshInstall   bool  // NEW
    SourceCount      int
    TotalAgents      int
    DeployedAgents   int
    LocationCount    int
    RecommendedNext  string
}

// Detection logic:
isFreshInstall := configExists &&
                  len(cfg.AgentSources) == 1 &&
                  cfg.AgentSources[0].Name == "my-agents" &&
                  (cfg.AgentSources[0].Git == nil || !cfg.AgentSources[0].Git.Enabled) &&
                  totalAgents == 0 &&
                  len(cfg.Locations) == 0

if isFreshInstall {
    responseText = "# Welcome to CAMI! 🌊\n\n"
    responseText += "I see this is a fresh installation. Let me guide you through getting started.\n\n"
    responseText += "## Three Paths Forward:\n\n"
    responseText += "**Path 1: Add Agent Library** (Recommended for new users)\n"
    responseText += "  → Get pre-built professional agents\n"
    responseText += "  → I can help you find and add agent sources\n\n"
    responseText += "**Path 2: Create Custom Agents**\n"
    responseText += "  → Work with agent-architect to build specialized agents\n"
    responseText += "  → Best if you have specific needs\n\n"
    responseText += "**Path 3: Import Existing Agents**\n"
    responseText += "  → Already have agents deployed elsewhere?\n"
    responseText += "  → I can scan and track them\n\n"
    responseText += "Which path interests you?\n"

    state.RecommendedNext = "Choose your onboarding path"
}
```

**IMPROVEMENT** (from agent-architect review):
Add `install_timestamp` and `setup_complete` flag to config for more robust detection:
```go
type Config struct {
    Version            string           `yaml:"version"`
    InstallTimestamp   time.Time        `yaml:"install_timestamp,omitempty"`  // NEW
    SetupComplete      bool             `yaml:"setup_complete,omitempty"`     // NEW
    AgentSources       []AgentSource    `yaml:"agent_sources"`
    Locations          []DeployLocation `yaml:"deploy_locations"`
    DefaultProjectsDir string           `yaml:"default_projects_dir,omitempty"`
}
```

---

### 1.2 Track External Agents (Import Workflow)

**Problem**: Agents deployed outside CAMI aren't tracked in manifests.

**Solution**: New `import_agents` tool that scans, analyzes, and manifests external agents.

**Implementation** (COMPLETED in Phase 1):
- Tool exists and works
- Uses composite matching (name + version + content hash) per agent-architect recommendation
- Origin tracking implemented

---

### 1.3 Configurable Development Workspace

**Problem**: Need default location for project creation, but it's confusing terminology.

**Solution**: Add `DefaultProjectsDir` to config, clarify terminology everywhere.

**Terminology Clarification**:
```
1. CAMI Source Repo (github.com/lando-labs/cami)
   - Clone this only to build/contribute
   - Users generally don't interact with this directly

2. CAMI Workspace (~/cami-workspace/)
   - Created by installer
   - Contains config.yaml, agent sources, CAMI's own agents
   - User interacts here to manage agents globally
   - Configurable via CAMI_DIR env var (already supported!)

3. Development Workspace (NEW: configurable)
   - Where user's projects live
   - Where agents get deployed
   - Default: ~/projects/ or ~/dev/ or ~/workspace/ (user choice)
   - Config field: default_projects_dir
```

**Implementation** (COMPLETED in Phase 1)

---

## Phase 2: Agent Classification Implementation

### 2.1 Agent Frontmatter Schema

**Minimal frontmatter approach** (per agent-architect recommendation):

```yaml
---
name: k8s-pod-checker
version: 1.0.0
description: Checks Kubernetes pod status following defined health check workflow
class: workflow-specialist          # NEW: workflow-specialist | technology-implementer | strategic-planner
specialty: kubernetes-operations    # NEW: Domain/specialty (free-form string)
---
```

**Why minimal?**
- Class determines phase weights automatically
- Specialty provides context without rigid taxonomy
- Avoids frontmatter bloat
- Easy to understand and maintain

### 2.2 Agent Schema Parser Update

```go
// internal/agent/agent.go
type Agent struct {
    Name        string
    Version     string
    Description string
    Model       string
    Color       string
    Class       string   // NEW: workflow-specialist, technology-implementer, strategic-planner
    Specialty   string   // NEW: Domain/specialty (e.g., "kubernetes-operations", "react-development")
    Content     string
    FilePath    string
}

// Phase weight lookup
var PhaseWeights = map[string]struct {
    Research  int
    Execute   int
    Validate  int
}{
    "workflow-specialist": {
        Research:  15,
        Execute:   70,
        Validate:  15,
    },
    "technology-implementer": {
        Research:  30,
        Execute:   55,
        Validate:  15,
    },
    "strategic-planner": {
        Research:  45,
        Execute:   30,
        Validate:  25,
    },
}
```

### 2.3 Agent-Architect Enhancement

**Auto-Classification** (not menu-driven):

```markdown
# Agent-Architect Enhancement

When creating agents, automatically classify based on user's request:

## Classification Signals

**Workflow Specialist**:
- User mentions "checklist", "workflow", "steps", "process"
- Single-purpose verbs: "check", "update", "deploy", "validate"
- Specific tools: "k8s", "jira", "github action"
- Examples: "check pod status", "update JIRA ticket", "deploy to staging"

**Technology Implementer**:
- User mentions "build", "create", "implement"
- Domain nouns: "frontend", "backend", "API", "database", "authentication"
- Feature scope: "complete screen", "full API", "entire flow"
- Examples: "build authentication", "create product API", "implement checkout"

**Strategic Planner**:
- User mentions "architect", "design", "research", "evaluate", "optimize"
- System-level scope: "architecture", "integration", "strategy", "performance"
- Decision-making: "should we", "what's the best way", "how do we scale"
- Examples: "design the system", "evaluate frameworks", "plan deployment strategy"

## Agent Creation Workflow

1. **Listen to user's request**
2. **Auto-classify** based on signals above
3. **Gather class-appropriate details**:

### For Workflow Specialists:
- "What specific workflow should this agent execute?"
- "What are the steps in order?"
- "What commands or tools will it use?"
- "What indicates success vs failure?"
- Offer three input methods:
  a) Describe interactively
  b) Provide file (markdown, shell script, YAML)
  c) Point to existing documentation

### For Technology Implementers:
- "What technology/framework are you using?"
- "What are the main features/capabilities needed?"
- "How does this integrate with the rest of the system?"
- "Any specific patterns or conventions to follow?"

### For Strategic Planners:
- "What decisions need to be made?"
- "What are the key constraints?" (timeline, budget, scale, team)
- "What tradeoffs matter most?" (speed vs quality, cost vs features)
- "What's the current state vs desired state?"

4. **Generate agent with appropriate phase weights**
5. **Include workflow in Execute phase** (for Workflow Specialists)

## Workflow Import/Parsing

Support three input formats for Workflow Specialists:

### Interactive Description
User describes workflow, agent-architect structures it:
```
User: "It should check if the pod is running, then check resource usage, then verify endpoints are responding"
Agent: *structures into numbered steps with validation*
```

### File Upload
Parse common formats:
- **Markdown**: Extract ordered/unordered lists, code blocks
- **Shell scripts**: Parse commands, comments become descriptions
- **YAML**: Support structured workflow definitions
- **Docs**: Extract procedures from documentation

### Point to Documentation
User provides URL or file path, agent-architect:
- Reads the documentation
- Extracts workflow/procedure
- Structures into agent format
- Asks user to confirm interpretation

## Agent File Structure

All agents follow this structure regardless of class:

```markdown
---
name: agent-name
version: 1.0.0
description: What this agent does
class: workflow-specialist
specialty: kubernetes-operations
model: sonnet  # or haiku for simple workflows
---

# Philosophy

[Why this agent exists, what problem it solves]

# Research Phase (15% for Workflow Specialists)

When invoked:
1. Read the current context
2. Understand the specific request
3. [Class-specific research steps]

# Execute Phase (70% for Workflow Specialists)

## Workflow: [Workflow Name]

**Purpose**: [What this workflow accomplishes]

**Prerequisites**:
- [What must be true before starting]

**Steps**:
1. [Step 1 with specific command/action]
   - Success criteria: [How to know it worked]
   - On failure: [What to do]

2. [Step 2 with specific command/action]
   - Success criteria: [How to know it worked]
   - On failure: [What to do]

3. [Continue...]

**Completion Criteria**:
- [How to know the entire workflow succeeded]

# Validate Phase (15% for Workflow Specialists)

After execution:
1. Verify all success criteria met
2. Report results to user
3. [Class-specific validation steps]

# When to Use This Agent

Use this agent when:
- [Specific trigger condition 1]
- [Specific trigger condition 2]

# When NOT to Use This Agent

Do not use this agent when:
- [Out of scope condition 1]
- [Out of scope condition 2]
```

## Validation

Before generating agent, validate:
- Workflow steps are clear and actionable
- Success/failure criteria are measurable
- Commands/tools are available in expected environment
- Workflow scope matches single-purpose principle

Ask user: "I've structured your workflow into X steps. Does this match your intent?"
```

---

## Phase 3: Guru-Guided Workflows

### 3.1 Philosophy: Guidance Not Prescription

**Current Problem**: Workflows are too prescriptive ("Do exactly these 7 steps").

**Solution**: Conversational, adaptive workflows that meet users where they are.

**Key Principles**:

1. **Start Broad, Narrow Contextually**
   - Don't assume user knows tech
   - Offer tech choices only when relevant
   - Provide escape hatches ("not sure" option)

2. **Vision Before Details**
   - Focus on what problem they're solving
   - Talk about outcomes, not implementation
   - Technical details emerge from conversation

3. **Tree Decisions Not Linear Steps**
   - Branch based on user responses
   - Adapt depth based on expertise signals
   - Allow backtracking and refinement

4. **Respect STRATEGIES.yaml**
   - Reference tech stack when needed
   - Don't re-ask what's already configured
   - Use as tiebreaker, not requirement
   - **EXPLICITLY tell user when referencing it** so they can update and restart

### 3.2 Enhanced create_project Workflow

```go
Description: "Create a new project with proper CAMI setup. " +
    "Use this when user wants to start a new project. " +

    "## Philosophy: Guru-Guided Creation\n" +
    "Be a wise guide, not a prescriptive checklist. Meet the user where they are. " +
    "Focus on vision and outcomes, let tech details emerge naturally.\n\n" +

    "## Conversational Workflow (Adaptive)\n\n" +

    "### Opening: Understand Intent\n" +
    "Ask: 'What are you building?' (open-ended)\n" +
    "Listen for:\n" +
    "  - Problem being solved\n" +
    "  - Who it's for\n" +
    "  - Key outcomes desired\n" +
    "Don't jump to tech! Stay high-level.\n\n" +

    "### Branching: Assess Technical Context\n" +
    "If user mentions specific tech → Note it, move forward\n" +
    "If user is vague → Offer gentle guidance:\n" +
    "  'I see you have [tech from STRATEGIES.yaml] configured. Is this for that stack?'\n" +
    "  **IMPORTANT**: Explicitly mention STRATEGIES.yaml so they know they can update it.\n" +
    "  'Not sure yet? That's fine - we can figure it out as we go.'\n\n" +

    "### Agent Recommendation: Capability-Based with Class Awareness\n" +
    "Based on what they're building, suggest agent classes:\n\n" +

    "**Workflow Specialists** (when they need specific processes):\n" +
    "  - 'Need to deploy regularly?' → deployment-workflow agent\n" +
    "  - 'Checking service health?' → health-check-workflow agent\n" +
    "  - 'Building specific components?' → component-builder agent\n\n" +

    "**Technology Implementers** (when they need full features):\n" +
    "  - 'Building user interfaces?' → frontend agent\n" +
    "  - 'Need data storage?' → database agent\n" +
    "  - 'Processing payments?' → payment-integration agent\n\n" +

    "**Strategic Planners** (when they need guidance):\n" +
    "  - 'Not sure about architecture?' → architect agent (offer proactively)\n" +
    "  - 'Need to evaluate technologies?' → researcher agent\n" +
    "  - 'Optimizing performance?' → performance agent\n\n" +

    "Use mcp__cami__list_agents to find matches.\n" +
    "If no match → Offer agent creation with agent-architect\n\n" +

    "### Scope Awareness: Match Agent Classes to Project Complexity\n" +
    "Simple app → Workflow Specialists + Technology Implementers\n" +
    "Complex system → Add Strategic Planner for integration/architecture\n" +
    "'Not sure about architecture?' → Suggest architect agent proactively\n\n" +

    "### Confirmation: Vision-Focused Summary\n" +
    "Show:\n" +
    "  - Project vision (what problem it solves)\n" +
    "  - Suggested agent guild (by capability and class)\n" +
    "  - Where it will be created\n" +
    "Ask: 'Does this feel right for what you're building?'\n" +
    "Allow refinement before executing.\n\n" +

    "### Execution: Create & Guide Next Steps\n" +
    "1. Create project directory\n" +
    "2. Deploy agreed agents\n" +
    "3. Write vision doc (based on conversation, not template)\n" +
    "4. Create manifest\n" +
    "5. Guide user: 'Your [project name] is ready. Want to start with [first capability]?'\n\n" +

    "## Adaptive Patterns\n\n" +

    "**Expert User** (mentions specific tech, knows stack):\n" +
    "  → Keep it brief, confirm tech, deploy agents, done\n" +
    "  → Don't over-explain\n\n" +

    "**Intermediate User** (knows some tech, unsure on others):\n" +
    "  → Reference STRATEGIES.yaml for tiebreakers (mention it by name)\n" +
    "  → Offer options: 'I see you have Next.js configured - use that?'\n" +
    "  → Explain tradeoffs when asked\n\n" +

    "**Beginner User** (vague tech, focused on outcomes):\n" +
    "  → Stay vision-focused: 'What should it do?' not 'What tech?'\n" +
    "  → Use STRATEGIES.yaml tech as defaults silently\n" +
    "  → Deploy agents that handle tech for them\n" +
    "  → Guide: 'Don't worry about [tech details], your agents will handle that'\n\n" +

    "## Tech Stack Handling\n\n" +

    "**If STRATEGIES.yaml exists**:\n" +
    "  - Use as defaults\n" +
    "  - Explicitly mention when referencing: 'Based on your STRATEGIES.yaml...'\n" +
    "  - This lets users know they can update it and restart\n" +
    "  - Don't re-ask about configured tech\n" +
    "  - Only surface tech questions for NEW decisions\n\n" +

    "**If no STRATEGIES.yaml**:\n" +
    "  - Still start with vision\n" +
    "  - Ask tech questions only when needed for agent selection\n" +
    "  - Use sensible defaults (React, Node.js, PostgreSQL)\n" +
    "  - Focus on capabilities: 'You'll need frontend, backend, database capabilities'\n\n" +

    "## Key Insight\n" +
    "The tech stack is a MEANS to an END. Focus on the end (what they're building), " +
    "let the means (tech) emerge from STRATEGIES.yaml or defaults. " +
    "Only ask tech questions when there's genuine ambiguity."
```

---

## Phase 4: Implementation Checklist

### 4.1 Code Changes

**Config Enhancements** (`internal/config/config.go`):
- [x] Add `DefaultProjectsDir` field (DONE)
- [ ] Add `InstallTimestamp` field (NEW from review)
- [ ] Add `SetupComplete` field (NEW from review)
- [ ] Add config migration for existing users

**Manifest Enhancements** (`internal/manifest/manifest.go`):
- [x] Add `Origin` field to `DeployedAgent` (DONE)
- [ ] Consider manifest format version to v3 if needed
- [ ] Add migration for v2 → v3 manifests if schema changes

**Agent Schema** (`internal/agent/agent.go`):
- [ ] Add `Class` field (workflow-specialist, technology-implementer, strategic-planner)
- [ ] Add `Specialty` field (free-form string)
- [ ] Update parser to handle new fields
- [ ] Add phase weight lookup function

**Enhanced Onboarding** (`cmd/cami/main.go`):
- [x] Add `IsFreshInstall` detection (DONE)
- [ ] Add install_timestamp checking (NEW)
- [ ] Add path-based guidance
- [ ] Add workspace location clarity
- [ ] Add existing agent scanning

**Enhanced create_project** (`cmd/cami/main.go`):
- [ ] Rewrite description with guru-guided philosophy
- [ ] Add adaptive workflow patterns
- [ ] Add class-awareness guidance
- [ ] Add STRATEGIES.yaml transparency (mention by name)
- [ ] Add tech stack handling rules

**Agent-Architect Updates** (`.claude/agents/agent-architect.md`):
- [ ] Add auto-classification framework
- [ ] Add class-based question sets
- [ ] Add workflow gathering (3 input methods)
- [ ] Add workflow validation
- [ ] Add phase-weighted agent template
- [ ] Update examples with all three classes

**Installer Enhancement** (`install/install.sh`):
- [x] Add default projects directory prompt (DONE)
- [ ] Add install_timestamp to config (NEW)
- [ ] Set setup_complete=false initially (NEW)

### 4.2 Documentation Updates

**CLAUDE.md** (CAMI development docs):
- [ ] Update with agent classification system
- [ ] Add class descriptions and examples
- [ ] Update onboarding workflow documentation
- [ ] Update create_project examples
- [ ] Add STRATEGIES.yaml transparency note

**User CLAUDE.md** (`install/templates/CLAUDE.md`):
- [ ] Add "Three Agent Classes" explanation
- [ ] Add "Three Locations" explanation
- [ ] Update onboarding paths section
- [ ] Add workflow specialist examples
- [ ] Add import workflow example

**README.md**:
- [ ] Add agent classification section
- [ ] Add phase-weighted methodology explanation
- [ ] Update feature list
- [ ] Add workflow specialist examples

**STRATEGIES.yaml Template** (`install/templates/sources/my-agents/STRATEGIES.yaml`):
- [ ] Add agent_class_philosophy section
- [ ] Add examples of each class
- [ ] Add guidance on when to use each class
- [ ] Explain phase weights

### 4.3 Testing Plan

**Unit Tests**:
- [x] Test fresh install detection (DONE)
- [x] Test import_agents scanning (DONE)
- [x] Test origin tracking in manifests (DONE)
- [ ] Test agent class parsing
- [ ] Test phase weight lookup
- [ ] Test install_timestamp logic

**Integration Tests**:
- [ ] Test full onboarding flow with install_timestamp
- [ ] Test agent-architect auto-classification
- [ ] Test workflow gathering (all 3 methods)
- [ ] Test create_project with class awareness
- [ ] Test STRATEGIES.yaml transparency

**User Testing**:
- [ ] Fresh install onboarding with new flags
- [ ] Create Workflow Specialist agent
- [ ] Create Technology Implementer agent
- [ ] Create Strategic Planner agent
- [ ] Test auto-classification accuracy

---

## Phase 5: Rollout Strategy

### 5.1 Version Targeting

**v0.3.2** (Config Improvements):
- Add install_timestamp and setup_complete to config
- Improve fresh install detection robustness
- Documentation updates for workspace clarity

**v0.4.0** (Agent Classification System):
- Add Class and Specialty fields to agent schema
- Update agent-architect with auto-classification
- Add phase weight system
- Workflow Specialist workflow gathering (3 methods)
- Enhanced create_project with class awareness

**v0.5.0** (Full Guru-Guided Experience):
- Fully adaptive create_project workflow
- Enhanced onboarding with class-aware suggestions
- STRATEGIES.yaml transparency (mention by name)
- Comprehensive documentation
- Example agents for all three classes

### 5.2 Migration Path

**For Existing Users**:
1. Auto-add install_timestamp (use file mtime of config.yaml)
2. Set setup_complete=true if sources exist
3. Existing agents work without class field (backward compatible)
4. Suggest running normalization to add class to existing agents
5. Update CLAUDE.md with new workflows

**Communication**:
- Release notes explain agent classification
- Show examples of each class
- Highlight workflow gathering feature
- Explain phase-weighted methodology

---

## Success Metrics

### Quantitative
- Onboarding completion rate (did user add sources and deploy agents?)
- Agent creation by class (are all three being used?)
- Workflow specialist adoption (how many workflow-driven agents created?)
- Project creation time (faster with guided workflow?)
- STRATEGIES.yaml update rate (are users updating it?)

### Qualitative
- User feedback on class system clarity
- Reported confusion about agent types (should decrease)
- User understanding of phase weights
- Satisfaction with guru-guided approach
- Feedback on workflow gathering methods

---

## Critical Clarifications (Based on User Feedback)

### What Workflow Specialists Are NOT:
❌ "Build any button" - Too broad, no specific workflow
❌ "Handle all deployments" - Too general, not checklist-driven
❌ "Fix bugs" - No defined process

### What Workflow Specialists ARE:
✅ "Check k8s pod status following these 5 steps"
✅ "Update JIRA ticket with this workflow"
✅ "Deploy to staging using this checklist"
✅ "Build button component following our design system steps"
✅ "Run API health checks in this sequence"

**Key Principle**: User provides the workflow, agent validates and executes it.

### Workflow Input Methods (All Three Supported):

1. **Interactive Description**
   - User describes workflow conversationally
   - Agent-architect structures it into steps
   - User confirms interpretation

2. **File Provision**
   - User uploads markdown checklist, shell script, or YAML
   - Agent-architect parses and structures
   - Supports: .md, .sh, .yaml, .txt

3. **Documentation Reference**
   - User points to existing docs (file or URL)
   - Agent-architect reads and extracts procedure
   - User confirms extracted workflow

---

## Open Questions

1. **Agent Classification**
   - Should we add validation that agent content matches declared class?
   - How to handle agents that span multiple classes?
   - Add linting tool to check class consistency?

2. **Workflow Gathering**
   - Which input method will users prefer?
   - Should we support workflow libraries (reusable workflows)?
   - How to version workflows separately from agents?

3. **Phase Weights**
   - Should phase weights be tunable by user?
   - How to communicate phase focus to users?
   - Should agents report which phase they're in during execution?

4. **STRATEGIES.yaml Transparency**
   - How often should we mention it (every reference or once per session)?
   - Should we offer to open/edit it directly?
   - Track when it was last updated?

---

## Next Steps

1. **Review V2 Plan** - Get agent-architect validation ✅ (Ready to do)
2. **Implement Config Improvements** (v0.3.2) - install_timestamp, setup_complete
3. **Update Agent-Architect** (v0.4.0) - Auto-classification, workflow gathering
4. **Test Classification** - Create sample agents in all three classes
5. **Roll Out Incrementally** - Phases as defined above
