# CAMI Workflow Enhancement Plan

**Goal**: Transform CAMI into a guru-guided experience with robust workflows, agent classification, and intelligent onboarding.

**Philosophy**: Provide *guidance* not *prescription*. Meet users where they are technically, keep them engaged with vision over details.

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
                  cfg.AgentSources[0].Git != nil &&
                  !cfg.AgentSources[0].Git.Enabled &&
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

**User Experience**:
```
User: "Help me get started with CAMI"
Claude: *uses onboard tool*
"I see this is a fresh installation. Let me guide you..."
[Presents 3 paths, asks which interests them]
User: "I already have some agents in my projects"
Claude: "Great! Let's scan your existing projects and bring them into CAMI's tracking system..."
```

---

### 1.2 Track External Agents (Import Workflow)

**Problem**: Agents deployed outside CAMI aren't tracked in manifests.

**Solution**: New `import_agents` tool that scans, analyzes, and manifests external agents.

**Implementation**:

```go
// New MCP tool: import_agents
type ImportAgentsArgs struct {
    ProjectPath string `json:"project_path"`  // Path to scan
    DryRun      bool   `json:"dry_run"`       // Preview only
}

type ImportedAgent struct {
    Name         string
    Version      string
    FilePath     string
    Origin       string  // "external" vs "cami"
    SourceMatch  string  // If matches a known source
    ContentHash  string
    MetadataHash string
}

// Workflow:
// 1. Scan .claude/agents/ directory
// 2. Parse each agent's frontmatter
// 3. Try to match to known sources (by name + version)
// 4. Create manifest entries with origin="external"
// 5. Update both local and central manifests
```

**New Manifest Field**:
```go
type DeployedAgent struct {
    Name           string    `yaml:"name"`
    Version        string    `yaml:"version"`
    Source         string    `yaml:"source"`
    SourcePath     string    `yaml:"source_path"`
    Priority       int       `yaml:"priority"`
    DeployedAt     time.Time `yaml:"deployed_at"`
    ContentHash    string    `yaml:"content_hash"`
    MetadataHash   string    `yaml:"metadata_hash"`
    CustomOverride bool      `yaml:"custom_override"`
    NeedsUpgrade   bool      `yaml:"needs_upgrade,omitempty"`
    Origin         string    `yaml:"origin"`  // NEW: "cami", "external", "manual"
}
```

**Tool Description**:
```go
Description: "Import and track agents that were deployed outside of CAMI. " +
    "WHEN TO USE: User has existing .claude/agents/ from manual deployment or other tools. " +
    "WORKFLOW: " +
    "1) Scan the specified project's .claude/agents/ directory " +
    "2) Parse agent frontmatter (name, version, description) " +
    "3) Try to match agents to known sources (by name + version hash) " +
    "4) Show preview of what will be imported " +
    "5) Ask user to confirm import " +
    "6) Create manifest entries with origin='external' " +
    "7) Update local and central manifests " +
    "BENEFITS: Brings existing agents into CAMI's tracking system for updates, version management, etc."
```

**User Experience**:
```
User: "I have agents in ~/projects/my-app from before CAMI"
Claude: "Let me scan that project and import those agents into CAMI's tracking system..."
*uses import_agents with dry_run=true*
"I found 5 agents:
 - frontend (v1.0.0) - Matches known source 'team-agents'
 - backend (v1.1.0) - Matches known source 'team-agents'
 - custom-agent (no version) - No source match (external)
 - old-agent (v0.5.0) - Outdated (v1.2.0 available)
 - my-tool (v1.0.0) - No source match (external)

Should I import these into CAMI's tracking system?"
User: "Yes"
Claude: *uses import_agents with dry_run=false*
"✓ Imported 5 agents
 ✓ Created manifest
 ✓ 2 updates available (old-agent, custom-agent)
 ✓ 3 already up to date"
```

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

**Config Addition**:
```go
type Config struct {
    Version            string           `yaml:"version"`
    AgentSources       []AgentSource    `yaml:"agent_sources"`
    Locations          []DeployLocation `yaml:"deploy_locations"`
    DefaultProjectsDir string           `yaml:"default_projects_dir,omitempty"` // NEW
}

// Default detection logic
func GetDefaultProjectsDir() (string, error) {
    cfg, err := Load()
    if err == nil && cfg.DefaultProjectsDir != "" {
        return cfg.DefaultProjectsDir, nil
    }

    // Fall back to ~/projects if not configured
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(homeDir, "projects"), nil
}
```

**Installer Enhancement**:
```bash
# In install.sh, after workspace setup:
echo ""
echo "Where do you want CAMI to create new projects by default?"
echo "Examples: ~/projects, ~/dev, ~/workspace, ~/code"
read -p "Default projects directory (press Enter for ~/projects): " PROJECTS_DIR

if [ -z "$PROJECTS_DIR" ]; then
    PROJECTS_DIR="$HOME/projects"
fi

# Expand ~ if present
PROJECTS_DIR="${PROJECTS_DIR/#\~/$HOME}"

# Add to config.yaml
sed -i '' "/^deploy_locations:/i\\
default_projects_dir: $PROJECTS_DIR
" "$INSTALL_DIR/config.yaml"

# Create directory if it doesn't exist
mkdir -p "$PROJECTS_DIR"

echo ""
print_success "Default projects directory: $PROJECTS_DIR"
```

**Onboarding Enhancement**:
```go
// In onboard tool, add workspace clarity
responseText += "## Your CAMI Setup\n\n"
responseText += "**CAMI Workspace:** " + configDir + "\n"
responseText += "  → Where agent sources and config live\n"
responseText += "  → Manage agents globally\n\n"

if cfg.DefaultProjectsDir != "" {
    responseText += "**Development Workspace:** " + cfg.DefaultProjectsDir + "\n"
    responseText += "  → Where your projects live\n"
    responseText += "  → Where agents get deployed\n\n"
}

responseText += "**Tracked Locations:** " + fmt.Sprintf("%d", len(cfg.Locations)) + "\n"
responseText += "  → Projects CAMI knows about\n\n"
```

---

## Phase 2: Agent Size Classification (sm/md/lg)

### 2.1 Agent Size Taxonomy

**Concept**: Classify agents by scope and complexity, not just by domain.

```yaml
# Agent frontmatter enhancement
---
name: button-specialist
version: 1.0.0
description: Creates and styles button components with accessibility
size: sm              # NEW
complexity: task      # NEW: task | builder | meta
scope: component      # NEW: component | feature | system
focus_areas:
  - React components
  - Tailwind styling
  - Accessibility (a11y)
---

# Examples by size:

## Small (sm) - Task-Efficient Agents
- Scope: Single component, single feature, specific task
- Complexity: Low cognitive load
- Examples:
  - button-specialist: Create button components
  - form-validator: Validate form inputs
  - error-handler: Handle and log errors
  - api-endpoint: Create single API endpoint
  - test-writer: Write unit tests for functions

## Medium (md) - Builder Agents
- Scope: Full feature, complete screen, entire API surface
- Complexity: Medium cognitive load, coordinates multiple concepts
- Examples:
  - frontend: Build complete UI screens
  - backend: Build full API with multiple endpoints
  - database: Design schema + migrations + queries
  - auth-system: Complete authentication flow
  - payment-integration: Full payment processing

## Large (lg) - Meta Agents
- Scope: System architecture, cross-cutting concerns, research
- Complexity: High cognitive load, connects many concepts
- Examples:
  - architect: System design, API contracts, integration
  - researcher: Deep research, technology evaluation
  - security: Cross-system security analysis
  - performance: System-wide performance optimization
  - devops: Infrastructure, deployment, monitoring
```

### 2.2 Size-Aware Workflows

**Agent Creation Workflow** (uses size to guide conversation):

```go
// create_agent tool enhancement
Description: "Create a new agent with agent-architect. " +
    "WORKFLOW: " +
    "1) Determine agent size/complexity first " +
    "2) Ask size-appropriate questions " +
    "3) Invoke agent-architect with size context " +

    "SIZE DETECTION: " +
    "Ask: 'What will this agent do?' " +
    "  - Single task/component → sm " +
    "  - Full feature/screen/API → md " +
    "  - Architecture/research/cross-cutting → lg " +

    "SIZE-SPECIFIC QUESTIONS: " +
    "Small (sm): " +
    "  - What specific task? (be concrete) " +
    "  - What inputs/outputs? " +
    "  - Any edge cases? " +
    "  → Focus on details, completion " +

    "Medium (md): " +
    "  - What feature/capability? (higher level) " +
    "  - What are the main components? " +
    "  - How does it integrate? " +
    "  → Focus on architecture, integration " +

    "Large (lg): " +
    "  - What system concerns? (highest level) " +
    "  - What decisions need making? " +
    "  - What tradeoffs exist? " +
    "  → Focus on strategy, research, planning " +

    "IMPORTANT: Match question depth to agent size. " +
    "Don't ask sm-agent questions for lg-agent creation."
```

**User Experience Examples**:

```
# Creating Small Agent
User: "I need an agent to create button components"
Claude: "That sounds like a focused, task-efficient agent. Let me ask a few specific questions:
  - What component library? (React, Vue, etc.)
  - What styling approach? (Tailwind, CSS-in-JS, etc.)
  - Any specific accessibility requirements?"
User: [answers]
Claude: *creates sm agent with agent-architect*

# Creating Medium Agent
User: "I need an agent to build authentication"
Claude: "Authentication is a full feature with multiple parts. Let me understand the scope:
  - What auth method? (JWT, sessions, OAuth)
  - What features? (login, signup, password reset, 2FA)
  - Where does it integrate? (frontend, backend, or both)"
User: [answers]
Claude: *creates md agent with agent-architect*

# Creating Large Agent
User: "I need an agent to help with system architecture"
Claude: "An architecture agent will help make high-level decisions. Let me understand your needs:
  - What decisions do you need help with? (monolith vs microservices, data flow, etc.)
  - What are your key constraints? (team size, timeline, scale, budget)
  - What tradeoffs matter most? (speed, cost, complexity)"
User: [answers]
Claude: *creates lg agent with agent-architect*
```

### 2.3 Size in STRATEGIES.yaml

```yaml
# STRATEGIES.yaml enhancement
agent_sizing_philosophy: |
  When creating agents for this guild, consider scope and user expertise:

  **Small Agents (sm)**
  - When user is focused on specific details
  - When task is well-defined and narrow
  - When quick iteration is needed
  - Examples: Component creators, validators, formatters

  **Medium Agents (md)**
  - When user needs a full feature built
  - When multiple concepts must integrate
  - When scope is a complete capability
  - Examples: Feature builders, API creators, screen designers

  **Large Agents (lg)**
  - When user needs strategic guidance
  - When decisions affect multiple systems
  - When research/planning is required
  - Examples: Architects, researchers, optimizers

# Agent roster by size
agent_roster:
  small:
    - button-specialist
    - form-validator
    - error-handler
  medium:
    - frontend
    - backend
    - database
  large:
    - architect
    - researcher

# When to use which size
sizing_guidance: |
  Ask yourself: "How much context does this agent need to hold?"

  - Can explain in 1 sentence → sm
  - Needs a paragraph → md
  - Needs a whole discussion → lg
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

4. **Respect Strategies.yaml**
   - Reference tech stack when needed
   - Don't re-ask what's already configured
   - Use as tiebreaker, not requirement

### 3.2 Enhanced create_project Workflow

**Before (Prescriptive)**:
```
1) Use AskUserQuestion to gather: project name, description, tech stack, key features
2) Use mcp__cami__list_agents to see available agents
3) Recommend agents based on requirements and get user confirmation
4) If agents don't exist, use Task tool to invoke agent-architect (parallel if multiple)
5) Write a focused vision_doc (200-300 words, vision NOT implementation)
6) Invoke this tool with name, description, agent_names, and vision_doc
7) Confirm success and guide user to next steps
```

**After (Guided)**:
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
    "  'I notice you have [tech from STRATEGIES.yaml] configured. Is this for that stack?'\n" +
    "  'Not sure yet? That's fine - we can figure it out as we go.'\n\n" +

    "### Agent Recommendation: Capability-Based\n" +
    "Based on what they're building (not just tech), suggest agent types:\n" +
    "  - 'Building user interfaces?' → frontend-type agent\n" +
    "  - 'Need data storage?' → database-type agent\n" +
    "  - 'Processing payments?' → integration-type agent\n" +
    "Use mcp__cami__list_agents to find matches.\n" +
    "If no match → Offer agent creation with agent-architect\n\n" +

    "### Size Awareness: Match Agent to Scope\n" +
    "Simple app → Suggest sm+md agents (task + feature builders)\n" +
    "Complex system → Add lg agent (architect for integration)\n" +
    "'Not sure about architecture?' → Suggest architect agent proactively\n\n" +

    "### Confirmation: Vision-Focused Summary\n" +
    "Show:\n" +
    "  - Project vision (what problem it solves)\n" +
    "  - Suggested agent guild (by capability, not tech)\n" +
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
    "  → Reference STRATEGIES.yaml for tiebreakers\n" +
    "  → Offer options: 'I see you have Next.js configured - use that?'\n" +
    "  → Explain tradeoffs when asked\n\n" +

    "**Beginner User** (vague tech, focused on outcomes):\n" +
    "  → Stay vision-focused: 'What should it do?' not 'What tech?'\n" +
    "  → Use STRATEGIES.yaml tech as defaults silently\n" +
    "  → Deploy agents that handle tech for them\n" +
    "  → Guide: 'Don't worry about [tech details], your agents will handle that'\n\n" +

    "## Examples\n\n" +

    "**Expert User**:\n" +
    "User: 'Create a Next.js app with Drizzle ORM and PostgreSQL'\n" +
    "Claude: 'Got it. I'll set up your Next.js project with frontend, backend, and database agents. Deploy to ~/projects/[name]?'\n" +
    "User: 'Yes'\n" +
    "Claude: *deploys and confirms*\n\n" +

    "**Beginner User**:\n" +
    "User: 'I want to build a todo app'\n" +
    "Claude: 'A todo app - nice! Let me understand:\n" +
    "  - Is this for learning or for real use?\n" +
    "  - Just you or multiple users?\n" +
    "  - Need it on mobile too, or just web?'\n" +
    "User: 'Just web, for learning'\n" +
    "Claude: 'Perfect for learning! I'll set you up with:\n" +
    "  - frontend agent (handles the interface)\n" +
    "  - backend agent (handles data)\n" +
    "  - qa agent (helps test as you build)\n" +
    "They'll guide you through building it step by step. Sound good?'\n" +
    "User: 'Yes!'\n" +
    "Claude: *creates project with vision: \"Learning project: Todo app for personal use\"*\n\n" +

    "**Intermediate User**:\n" +
    "User: 'I need an e-commerce site'\n" +
    "Claude: 'E-commerce - there's a lot involved! Let me break it down:\n" +
    "  - Product catalog\n" +
    "  - Shopping cart\n" +
    "  - Payment processing\n" +
    "  - User accounts\n" +
    "Do you want all of these, or starting with a subset?'\n" +
    "User: 'Start with catalog and cart, add payments later'\n" +
    "Claude: 'Smart approach! I'll set you up with:\n" +
    "  - frontend agent (product pages, cart UI)\n" +
    "  - backend agent (product API, cart logic)\n" +
    "  - database agent (product catalog schema)\n" +
    "  - architect agent (helps plan payment integration when ready)\n" +
    "I see you have Next.js + PostgreSQL configured - use those?'\n" +
    "User: 'Yes'\n" +
    "Claude: *creates project*\n\n" +

    "## Tech Stack Handling\n\n" +

    "**If STRATEGIES.yaml exists**:\n" +
    "  - Use as defaults\n" +
    "  - Don't re-ask about configured tech\n" +
    "  - Only surface tech questions for NEW decisions\n" +
    "  - Example: 'Need payment processing - should we use Stripe (recommended) or other?'\n\n" +

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

### 3.3 Enhanced Onboarding Workflow

```go
Description: "Get personalized onboarding guidance for CAMI. " +
    "WHEN TO USE: User is new to CAMI, asks 'what should I do next?', seems lost, or just installed. " +

    "## Workflow: Discovery & Guidance\n\n" +

    "### Phase 1: Detect Context\n" +
    "1. Check if fresh install (default my-agents, no real sources, no locations)\n" +
    "2. Check if has existing agents (scan projects for .claude/agents/)\n" +
    "3. Check if has configured sources\n" +
    "4. Determine user's starting point\n\n" +

    "### Phase 2: Present Contextual Paths\n\n" +

    "**Fresh Install**:\n" +
    "Present 3 paths:\n" +
    "  1. Add Agent Library (get pre-built agents)\n" +
    "  2. Create Custom Agents (build from scratch)\n" +
    "  3. Import Existing Agents (already have .claude/agents/ elsewhere)\n" +
    "Ask: 'Which path interests you?'\n\n" +

    "**Has Existing Agents**:\n" +
    "  → Proactively offer import: 'I see you have agents in [locations]. Want to bring them into CAMI?'\n\n" +

    "**Has Sources But Nothing Deployed**:\n" +
    "  → Guide to first deployment: 'You have X agents available. Want to start a project?'\n\n" +

    "**Fully Configured**:\n" +
    "  → Offer next-level features: 'Looks good! Want to explore agent creation, updates, or normalization?'\n\n" +

    "### Phase 3: Execute Chosen Path\n" +
    "Based on user choice:\n" +
    "  - Path 1 → Use add_source tool, suggest sources\n" +
    "  - Path 2 → Start agent creation with agent-architect\n" +
    "  - Path 3 → Use import_agents tool\n" +
    "  - Next level → Explain available workflows\n\n" +

    "## Workspace Clarity\n" +
    "Always explain the three locations clearly:\n" +
    "  1. CAMI Workspace (~/cami-workspace/) - Where agent sources live\n" +
    "  2. Development Workspace ([default_projects_dir]) - Where projects get created\n" +
    "  3. Tracked Locations - Specific projects CAMI knows about\n\n" +

    "Make this clear in response:\n" +
    "'Your Setup:\n" +
    " • CAMI Workspace: ~/cami-workspace/ (agent management HQ)\n" +
    " • Projects Directory: ~/projects/ (where new projects go)\n" +
    " • Tracked Projects: 3 (my-app, client-site, cami)'\n\n" +

    "## Tone: Welcoming Guide\n" +
    "- Use 'Let me guide you...' not 'You must do...'\n" +
    "- Ask questions, don't prescribe\n" +
    "- Offer options, not orders\n" +
    "- Celebrate progress: '✓ Great! You're set up'\n" +
    "- Encourage next steps: 'Ready to try creating a project?'"
```

---

## Phase 4: Implementation Checklist

### 4.1 Code Changes

**Config Enhancements** (`internal/config/config.go`):
- [ ] Add `DefaultProjectsDir` field
- [ ] Add `GetDefaultProjectsDir()` function
- [ ] Add config migration for existing users

**Manifest Enhancements** (`internal/manifest/manifest.go`):
- [ ] Add `Origin` field to `DeployedAgent`
- [ ] Update manifest format version to 3
- [ ] Add migration for v2 → v3 manifests

**New Tool: import_agents** (`cmd/cami/main.go`):
- [ ] Implement scan logic
- [ ] Implement source matching
- [ ] Implement manifest creation with origin tracking
- [ ] Add dry-run preview mode

**Enhanced Onboarding** (`cmd/cami/main.go`):
- [ ] Add `IsFreshInstall` detection
- [ ] Add path-based guidance
- [ ] Add workspace location clarity
- [ ] Add existing agent scanning

**Enhanced create_project** (`cmd/cami/main.go`):
- [ ] Rewrite description with guru-guided philosophy
- [ ] Add adaptive workflow patterns
- [ ] Add size-awareness guidance
- [ ] Add tech stack handling rules

**Agent Frontmatter** (`internal/agent/agent.go`):
- [ ] Add `Size` field (sm/md/lg)
- [ ] Add `Complexity` field (task/builder/meta)
- [ ] Add `Scope` field (component/feature/system)
- [ ] Update parser to handle new fields

**Installer Enhancement** (`install/install.sh`):
- [ ] Add default projects directory prompt
- [ ] Add to config.yaml during install
- [ ] Create directory if doesn't exist
- [ ] Add to success message

### 4.2 Documentation Updates

**CLAUDE.md**:
- [ ] Update onboarding workflow documentation
- [ ] Add import_agents tool documentation
- [ ] Add size classification explanation
- [ ] Update create_project examples
- [ ] Add workspace location clarity

**User CLAUDE.md** (`install/templates/CLAUDE.md`):
- [ ] Add "Three Locations" explanation
- [ ] Update onboarding paths section
- [ ] Add agent size explanation
- [ ] Add import workflow example

**README.md**:
- [ ] Add workspace configuration section
- [ ] Add agent size classification section
- [ ] Update feature list
- [ ] Add import workflow to examples

**STRATEGIES.yaml Template** (`install/templates/sources/my-agents/STRATEGIES.yaml`):
- [ ] Add agent_sizing_philosophy section
- [ ] Add agent_roster by size
- [ ] Add sizing_guidance

### 4.3 Testing Plan

**Unit Tests**:
- [ ] Test fresh install detection
- [ ] Test import_agents scanning
- [ ] Test source matching logic
- [ ] Test origin tracking in manifests
- [ ] Test default projects dir resolution

**Integration Tests**:
- [ ] Test full onboarding flow
- [ ] Test import workflow
- [ ] Test create_project with size awareness
- [ ] Test manifest migration v2 → v3

**User Testing**:
- [ ] Fresh install onboarding
- [ ] Import existing agents workflow
- [ ] Create project with beginner persona
- [ ] Create project with expert persona
- [ ] Multi-size agent roster creation

---

## Phase 5: Rollout Strategy

### 5.1 Version Targeting

**v0.3.1** (Quick Fixes):
- Fix onboarding fresh install detection
- Add default_projects_dir to config
- Update installer to prompt for projects dir
- Documentation clarity on three locations

**v0.4.0** (Import & Size Classification):
- Add import_agents tool
- Add Origin field to manifests
- Add Size/Complexity/Scope fields to agents
- Update agent-architect to use size classification
- Enhanced create_project workflow

**v0.5.0** (Guru-Guided Workflows):
- Fully adaptive create_project workflow
- Enhanced onboarding with path selection
- Size-aware agent creation workflows
- Comprehensive documentation updates

### 5.2 Migration Path

**For Existing Users**:
1. Auto-detect if `default_projects_dir` missing → prompt during next use
2. Auto-migrate manifests from v2 → v3 (add origin="cami" to all)
3. Suggest running import_agents on tracked locations
4. Update CLAUDE.md with new workflows

**Communication**:
- Release notes explain new workflows
- Show before/after examples
- Highlight improved onboarding
- Explain import feature for existing agents

---

## Success Metrics

### Quantitative
- Onboarding completion rate (did user add sources and deploy agents?)
- Import usage (how many external agents imported?)
- Project creation time (faster with guided workflow?)
- Agent roster diversity (using sm/md/lg effectively?)

### Qualitative
- User feedback on onboarding experience
- Reported confusion about workspaces (should decrease)
- User understanding of agent sizes
- Satisfaction with guru-guided approach

---

## Open Questions

1. **Size Classification**
   - Should agents self-report size or should it be auto-detected?
   - Can we validate size against actual agent content?
   - Should size affect tool selection in agents?

2. **Import Workflow**
   - Should import be automatic on first onboarding?
   - How to handle naming conflicts during import?
   - Should import update agents to match sources?

3. **Default Projects Directory**
   - Should we auto-detect common locations (~/projects, ~/dev, ~/src)?
   - Allow multiple default locations for different types?
   - Integrate with IDEs' workspace directories?

4. **Workflow Philosophy**
   - How prescriptive should size-based questioning be?
   - Balance between guidance and prescription?
   - When to override user choice for their benefit?

---

## Next Steps

1. **Review & Refine** - Get feedback on this plan
2. **Prioritize** - Which phase/feature first?
3. **Prototype** - Build onboarding fix first (v0.3.1)
4. **Test** - Early user testing on guided workflows
5. **Iterate** - Refine based on feedback
6. **Ship** - Roll out in phases

