# CAMI User Onboarding Guide

**Version:** 1.0 (for v1.0.0 open source release)
**Audience:** New CAMI users (first-time setup)

---

## New User Journey: From Zero to Deployed Agents in 5 Minutes

### The Bootstrap Problem

**Chicken-and-egg situation:**
- New users want to use `agent-architect` to create custom agents
- But `agent-architect` itself is an agent in `lando-labs/lando-agents`
- How do they get it initially?

**Solution:** CAMI provides **multiple paths** based on user type:

---

## Path 1: Agent Consumer (Most Users) - "Just Use Agents"

**Goal:** Deploy existing agents to projects, no customization needed

### Step 1: Install CAMI

```bash
# Install via Go
go install github.com/lando-labs/cami/cmd/cami@latest

# Or via Homebrew (future)
brew install lando-labs/tap/cami
```

### Step 2: Initialize CAMI

```bash
cami init
```

**Interactive wizard:**
```
Welcome to CAMI! 🚀

Let's set up your environment.

? What's your primary use case?
  > Use existing agents in my projects (recommended)
    Develop custom agents
    Use company-internal agents
    Mix of everything (advanced)

✓ You selected: Use existing agents

? Add official Lando Labs agent library?
  > Yes (recommended)
    No, I'll add sources manually

✓ Adding source: lando-labs/lando-agents

? Do you have a project to deploy agents to?
  > Yes, let me select one
    No, I'll add projects later

? Select project directory:
  > /Users/dev/my-awesome-app

✓ Added location: my-awesome-app

Setup complete! 🎉

Try these commands:
  cami list              # See all available agents
  cami search frontend   # Find specific agents
  cami deploy            # Deploy agents (interactive)
```

### Step 3: Browse Available Agents

```bash
# TUI (interactive)
cami

# Or CLI
cami list
```

**Output:**
```
Available Agents (29 total)

── core ──
  architect       v1.1.0  System architecture and design
  backend         v1.1.0  Backend APIs and databases
  frontend        v1.1.0  React/Next.js UI development
  mobile-native   v1.1.0  iOS/Android native development
  qa              v1.1.0  Testing and quality assurance

── specialized ──
  ai-ml-specialist      v1.1.0  AI/ML integration
  design-system-spec... v1.0.0  Component libraries
  react-specialist      v1.0.0  Advanced React patterns
  ...

── meta ──
  agent-architect  v1.1.0  Create and optimize agents
```

### Step 4: Deploy Agents to Project

```bash
# Interactive (recommended for first time)
cami deploy

# Or with flags
cami deploy -a architect,frontend,backend -l my-awesome-app
```

**Interactive flow:**
```
? Select agents to deploy (space to select, enter to confirm):
  [x] architect
  [x] frontend
  [ ] backend
  [ ] qa

? Select deployment location:
  > my-awesome-app (/Users/dev/my-awesome-app)
    Add new location...

✓ Deployed 2 agents to my-awesome-app
  → architect v1.1.0
  → frontend v1.1.0

✓ Updated CLAUDE.md with agent documentation
```

### Step 5: Use Agents in Claude Code

```bash
cd /Users/dev/my-awesome-app
claude
```

**In Claude Code:**
```
> @architect help me design a new API

[Agent loads from .claude/agents/architect.md]
[Architect agent provides guidance...]
```

**That's it! No need to touch agent-architect yet.**

---

## Path 2: Agent Developer - "I Want to Create Agents"

**Goal:** Create custom agents, potentially contribute to lando-labs/lando-agents

### Step 1-3: Same as Path 1 (Install, Init, Browse)

### Step 4: Access Meta-Agent

**Option A: Deploy agent-architect to a project**
```bash
cami deploy -a agent-architect -l ~/my-project

cd ~/my-project
claude
> @agent-architect create a new data-science agent for me
```

**Option B: Create agent workspace**
```bash
# Create dedicated workspace for agent development
mkdir ~/my-agents
cd ~/my-agents

# Deploy agent-architect here
cami deploy -a agent-architect -l ~/my-agents

# Now use it to create agents
claude
> @agent-architect help me create a blockchain-specialist agent
```

**Option C: Use CAMI template (no agent-architect needed)**
```bash
# CAMI provides a starter template
cami create custom-agent --template basic

# Opens editor with:
# ~/my-agents/custom-agent.md with full frontmatter template
```

### Step 5: Add Local Source

```bash
# Point CAMI to your custom agents
cami sources add my-agents local ~/my-agents --priority 200

# Now your agents are available
cami list
# Shows: official agents + your custom agents
```

### Step 6: Deploy Your Custom Agent

```bash
cami deploy -a custom-agent -l ~/some-project

# Your agent (priority 200) is used instead of official (priority 100)
```

---

## Path 3: Team Library - "We Need Company Agents"

**Goal:** Set up private agent repository for team

### Step 1-2: Install and Init CAMI

```bash
go install github.com/lando-labs/cami/cmd/cami@latest
cami init
# Select: "Use company-internal agents"
```

### Step 3: Create Company Agent Repo

```bash
# On GitHub/GitLab
# Create private repo: acme-corp/claude-agents

# Clone and structure it
git clone git@github.com:acme-corp/claude-agents.git
cd claude-agents

# Use official agents as templates
cami export architect > core/architect.md
cami export backend > core/backend.md

# Customize for your company
vim core/backend.md
# Add company-specific context, APIs, patterns

git add core/
git commit -m "Initial agent library"
git push
```

### Step 4: Team Setup

Share with team:
```bash
# Each team member runs:
cami sources add acme git git@github.com:acme-corp/claude-agents.git \
  --auth ssh-key --priority 150

# Also add official for baseline
cami sources add official github lando-labs/lando-agents --priority 100
```

**Team member config (`~/.cami/config.yaml`):**
```yaml
agent_sources:
  - name: acme
    type: git
    url: git@github.com:acme-corp/claude-agents.git
    auth: ssh-key
    priority: 150    # Company agents win

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100    # Official as fallback
```

### Step 5: Deploy Team Agents

```bash
# Team agents override official when names conflict
cami list
# Shows: acme agents (priority 150) + official agents (priority 100)

cami deploy -a backend -l ~/project
# Uses acme's backend (customized), not official
```

---

## Path 4: Contributor - "I Want to Improve Official Agents"

**Goal:** Contribute to lando-labs/lando-agents

### Step 1: Fork and Clone

```bash
# Fork lando-labs/lando-agents on GitHub

# Clone your fork
git clone git@github.com:YOUR_USERNAME/lando-agents.git
cd lando-agents
```

### Step 2: Point CAMI to Local Fork

```bash
cami sources add my-fork local ~/path/to/lando-agents --priority 200
```

### Step 3: Edit Agents

```bash
cd ~/path/to/lando-agents
vim core/frontend.md
# Improve the agent
```

### Step 4: Test Your Changes

```bash
# Deploy your edited version (priority 200)
cami deploy -a frontend -l ~/test-project

# Test in Claude Code
cd ~/test-project
claude
> @frontend help me build a component
```

### Step 5: Contribute Back

```bash
git add core/frontend.md
git commit -m "Improve frontend agent MCP awareness"
git push origin main

# Create PR to lando-labs/lando-agents
```

---

## Common Scenarios

### Scenario 1: "I Want to Create One Custom Agent"

**Fastest path:**
```bash
# Use CAMI template (no agent-architect needed)
cami create my-agent --template basic --category specialized

# Edit the generated file
vim ~/.cami/local-agents/specialized/my-agent.md

# Add local source
cami sources add local local ~/.cami/local-agents --priority 200

# Deploy it
cami deploy -a my-agent -l ~/project
```

### Scenario 2: "I Need agent-architect Frequently"

**Create a dedicated workspace:**
```bash
mkdir ~/agent-workspace
cami deploy -a agent-architect -l ~/agent-workspace

# Add to shell profile
echo 'alias agent-lab="cd ~/agent-workspace && claude"' >> ~/.zshrc

# Now from anywhere:
agent-lab
> @agent-architect create a new monitoring-specialist agent
```

### Scenario 3: "My Team Wants Official + Custom Agents"

**Multi-source setup:**
```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: team-custom
    type: git
    url: git@github.com:mycompany/custom-agents.git
    priority: 150

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100
```

**Result:**
- Team agents override official when names match
- Official agents fill in gaps
- Best of both worlds

### Scenario 4: "I'm Experimenting with Agent Ideas"

**Local override:**
```bash
mkdir ~/experiments
cami sources add experiments local ~/experiments --priority 200

# Create quick test agents
vim ~/experiments/test-agent.md

# Test immediately
cami deploy -a test-agent -l ~/sandbox
```

---

## Decision Tree: Which Path Am I?

```
Are you creating custom agents?
├─ No → Path 1 (Agent Consumer)
│       Just use existing agents
│
└─ Yes → Do you need agent-architect's help?
    ├─ No → Use `cami create` templates
    │
    └─ Yes → Is this for a team?
        ├─ No → Path 2 (Agent Developer)
        │       Create workspace, deploy agent-architect
        │
        └─ Yes → Path 3 (Team Library)
                Set up private repo
```

---

## Meta-Agent Access Strategy

**Key insight:** You don't need agent-architect to GET STARTED.

**When you DO need agent-architect:**
1. Creating complex agents with specific patterns
2. Optimizing existing agents systematically
3. Learning best practices from meta-agent guidance

**How to access it:**
- **Option 1:** Deploy to a dedicated workspace (`~/agent-lab`)
- **Option 2:** Deploy to your main project when needed
- **Option 3:** Use CAMI templates for simple agents (no meta-agent needed)

**Not needed for:**
- Using existing agents (Path 1)
- Simple agent customization (templates work fine)
- First-time users (too advanced)

---

## Configuration Examples

### Minimal Config (Agent Consumer)

```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100

deploy_locations:
  - name: my-app
    path: /Users/dev/my-app
```

### Developer Config (Local Override)

```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: local-dev
    type: local
    path: /Users/dev/lando-agents
    priority: 200

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100
```

### Team Config (Company + Official)

```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: acme-internal
    type: git
    url: git@github.com:acme/agents.git
    auth: ssh-key
    priority: 150

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100

  - name: my-experiments
    type: local
    path: ~/my-agents
    priority: 200
```

---

## First-Run Experience

### v1.0.0 Default Behavior

When user runs `cami` for the first time:

```bash
cami
```

**Output:**
```
Welcome to CAMI! 🚀

No configuration found. Let's set you up.

? Run setup wizard?
  > Yes, guide me through setup
    No, I'll configure manually

✓ Running setup wizard...

? What's your primary use case?
  > Use existing agents (recommended)
    Develop custom agents
    Company-internal agents
    Advanced (multiple sources)

[... interactive setup ...]

Setup complete! 🎉

Configuration saved to: ~/.cami/config.yaml

Next steps:
  cami list     # Browse available agents
  cami deploy   # Deploy agents to projects
  cami --help   # See all commands
```

### Skip Wizard (Expert Mode)

```bash
CAMI_NO_WIZARD=1 cami
# Creates minimal config, jumps to TUI
```

---

## Summary: Solving the Meta-Agent Paradox

**The Problem:**
- agent-architect is itself an agent
- New users don't have it initially
- How do they bootstrap?

**The Solution:**
1. **Most users don't need it** - Use existing agents (Path 1)
2. **Simple customization** - Use `cami create` templates (no meta-agent)
3. **Complex customization** - Deploy agent-architect to workspace (Path 2)
4. **Advanced users** - Work inside lando-agents repo with local source

**Key Principle:**
> CAMI provides **progressive complexity**. Start simple, scale as needed.

**Onboarding Flow:**
```
New User
  → Install CAMI
  → Run `cami init` (wizard)
  → Deploy existing agents
  → Use in projects
  → (Optional) Create simple agents with templates
  → (Optional) Deploy agent-architect when needed
  → (Optional) Contribute back to lando-agents
```

---

## FAQ

### Q: Do I need agent-architect to use CAMI?
**A:** No! Most users just deploy existing agents. agent-architect is for creating/optimizing custom agents.

### Q: How do I get agent-architect if it's an agent?
**A:** It's in the official lando-labs/lando-agents repo. Deploy it like any other agent: `cami deploy -a agent-architect -l ~/workspace`

### Q: Can I create agents without agent-architect?
**A:** Yes! Use `cami create <name> --template basic` for simple agents. agent-architect helps with complex/optimized agents.

### Q: What if I work inside lando-agents repo?
**A:** Set local source with priority 200. Your edits override official agents immediately.

### Q: How do teams share custom agents?
**A:** Create private Git repo, add as source with priority 150. Team agents override official.

### Q: What happens when agent names conflict?
**A:** Higher priority wins. Local (200) > Team (150) > Official (100).

---

**This guide ensures every user type has a clear, logical onboarding path without hitting the meta-agent paradox.**
