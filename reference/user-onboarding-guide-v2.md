# CAMI User Onboarding Guide v2

**Version:** 2.0 (embraces agent-architect power)
**Date:** 2025-11-10

---

## Philosophy

> **Agent-architect is not an advanced feature. It's the core value proposition.**
>
> CAMI empowers you to create sophisticated agents on demand using agent-architect's research capabilities. Templates are weak. Research-driven agent creation is powerful.

---

## First-Time User Experience

### Recommended: Run CAMI from Source Directory

**Why?**
- You get agent-architect immediately (in `vc-agents/meta/`)
- You can create/edit agents in `vc-agents/`
- Your local edits work immediately (local source priority 200)
- Standard Git workflow for contributing back

### Installation

```bash
# Clone CAMI
git clone git@github.com:lando-labs/cami.git
cd cami

# Build
go build -o cami ./cmd/cami

# Initialize
./cami init
```

### Interactive Onboarding

```
Welcome to CAMI! 🚀

I see you're running from the CAMI source directory.
This is the best setup for creating and managing agents!

✓ Found 29 agents in vc-agents/
✓ Found agent-architect in vc-agents/meta/

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CAMI empowers you to CREATE AGENTS ON DEMAND.

The secret? Agent-architect researches domains, analyzes
Claude Code capabilities, and generates sophisticated agents
tailored to your needs.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

? What do you want to do?
  > Create agents on demand with agent-architect (recommended)
    Use existing agents only
    I'll explore on my own

✓ Great choice! Let's set up agent-architect.

? Where should we deploy agent-architect?
  > Current directory (work here, create agents here)
    Dedicated workspace (~/agent-workspace)
    Specific project

✓ Deploying agent-architect to /Users/lando/lando-labs/cami

? Configure vc-agents/ as local source (priority 200)?
  This lets you edit agents and use them immediately.
  > Yes (recommended for development)
    No

✓ Added local source: vc-agents (priority 200)

? Add official lando-labs/lando-agents as backup source?
  This provides baseline agents and updates.
  > Yes (recommended)
    No

✓ Added remote source: lando-labs/lando-agents (priority 100)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Setup complete! 🎉

You're ready to create agents.

Quick Start:
  1. Run: claude
  2. Ask: @agent-architect create a Terraform specialist agent
  3. Agent-architect will research and create a sophisticated agent
  4. Test it: ./cami deploy -a terraform -l ~/test-project
  5. Refine as needed
  6. Contribute: git add, git commit, git push (standard Git)

Try it now:
  claude

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Primary Workflow: Agent Creation

### Step 1: Ask Agent-Architect

```bash
cd ~/lando-labs/cami
claude
```

**In Claude Code:**
```
> @agent-architect I need an agent specialized in Terraform infrastructure as code

[Agent-architect researches:]
- Terraform best practices
- IaC patterns
- Claude Code tool capabilities
- Related MCPs (if available)
- Integration patterns

[Agent-architect creates:]
vc-agents/infrastructure/terraform-specialist.md
- Sophisticated system prompt
- Terraform-specific guidance
- IaC workflow patterns
- MCP awareness
- Best practices
```

### Step 2: Test the Agent

```bash
./cami deploy -a terraform-specialist -l ~/test-terraform-project

cd ~/test-terraform-project
claude
> @terraform-specialist help me design a multi-region AWS setup
```

### Step 3: Refine

**If it needs improvement:**
```bash
cd ~/lando-labs/cami
claude
> @agent-architect the terraform agent needs more focus on state management
```

**Or edit directly:**
```bash
vim vc-agents/infrastructure/terraform-specialist.md
# Make changes

./cami deploy -a terraform-specialist -l ~/test-terraform-project
# Test again
```

### Step 4: Contribute (Optional)

**Standard Git workflow:**
```bash
git status
# On branch main
# Untracked files:
#   vc-agents/infrastructure/terraform-specialist.md

git add vc-agents/infrastructure/terraform-specialist.md
git commit -m "Add Terraform specialist agent

- IaC best practices
- Multi-region support
- State management focus
- Workspace patterns"

git push origin main

# Or create PR if working on fork
gh pr create --title "Add Terraform specialist agent" \
  --body "Sophisticated agent for Terraform IaC workflows..."
```

---

## Secondary Workflow: Using Existing Agents

### For Users Who Just Want to Deploy

```bash
# Install CAMI globally (not from source)
go install github.com/lando-labs/cami/cmd/cami@latest

# Initialize
cami init
# Select: "Use existing agents only"

# Browse
cami list

# Deploy
cami deploy -a architect,frontend,backend -l ~/my-project

# Update periodically
cami update
```

**This workflow is simpler but less powerful.** No agent creation, just consumption.

---

## Advanced Workflow: Local + Remote Sources

### Scenario: Use Official + Create Custom

```yaml
# ~/.cami/config.yaml
agent_sources:
  - name: my-custom
    type: local
    path: ~/my-agents
    priority: 200        # Your custom agents win

  - name: official
    type: github
    repo: lando-labs/lando-agents
    priority: 100        # Official as baseline
```

**Workflow:**
```bash
# Use agent-architect from official source
cami deploy -a agent-architect -l ~/my-agents

cd ~/my-agents
claude
> @agent-architect create a blockchain specialist

# Agent created in ~/my-agents/blockchain-specialist.md
# CAMI sees it immediately (priority 200)

./cami deploy -a blockchain-specialist -l ~/crypto-project
```

---

## Understanding Priorities

### How Conflict Resolution Works

If the same agent exists in multiple sources, **highest priority wins**.

**Example:**
```yaml
agent_sources:
  - name: my-experiments
    path: ~/experiments
    priority: 200        # LOCAL: Your edits

  - name: team
    url: git@github.com:acme/agents.git
    priority: 150        # TEAM: Company agents

  - name: official
    repo: lando-labs/lando-agents
    priority: 100        # OFFICIAL: Baseline
```

**If `frontend.md` exists in all three:**
- CAMI uses `~/experiments/frontend.md` (priority 200)
- Remove from experiments → falls back to team's version
- Remove from team → falls back to official

**This enables:**
- ✅ Local experimentation without breaking anything
- ✅ Team customization without forking official
- ✅ Fallback to official when custom not available
- ✅ Predictable, transparent behavior

---

## Common Questions

### "Do I need to know Git?"

**For consuming agents:** No
- `cami update` pulls agents automatically
- No Git knowledge needed

**For creating/contributing agents:** Yes (basic)
- `git add`, `git commit`, `git push`
- Standard Git workflow
- LLMs (like Claude) can guide you

### "Should I run CAMI from source or install globally?"

**From source (recommended for developers):**
- ✅ Access to agent-architect immediately
- ✅ Edit agents in vc-agents/
- ✅ Local source (priority 200) for instant testing
- ✅ Contribute back with standard Git

**Global install (simpler for consumers):**
- ✅ Easier updates (`go install ...@latest`)
- ✅ Works from any directory
- ❌ No local agent development
- ❌ Harder to contribute back

### "How do I contribute my custom agent?"

**If working in lando-labs/cami (or fork):**
```bash
git add vc-agents/category/your-agent.md
git commit -m "Add your-agent"
git push origin main
# Or: gh pr create
```

**If working in separate directory:**
```bash
# Copy agent to lando-agents repo
cp ~/my-agents/cool-agent.md ~/lando-labs/lando-agents/vc-agents/specialized/

cd ~/lando-labs/lando-agents
git add vc-agents/specialized/cool-agent.md
git commit -m "Add cool-agent"
git push
```

**CAMI doesn't manage this - use standard Git.**

### "Can agent-architect improve existing agents?"

**Absolutely!**
```bash
claude
> @agent-architect review and improve the frontend agent

[Agent-architect analyzes:]
- Current frontend.md
- Latest React/Next.js patterns
- MCP integrations
- Gaps or outdated guidance

[Agent-architect suggests improvements]

> @agent-architect apply those improvements to frontend.md

[Updates vc-agents/core/frontend.md]

# Test
./cami deploy -a frontend -l ~/test-project

# Commit if good
git add vc-agents/core/frontend.md
git commit -m "Improve frontend agent with React 19 patterns"
```

### "What if I mess up an official agent?"

**No worries!** Your local source (priority 200) overrides, but doesn't break official.

```bash
# You broke frontend.md locally
vim vc-agents/core/frontend.md
# Made bad edits

# Option 1: Fix it
vim vc-agents/core/frontend.md
# Correct the edits

# Option 2: Revert to official
rm vc-agents/core/frontend.md
# CAMI falls back to official (priority 100)

# Option 3: Git revert
git checkout vc-agents/core/frontend.md
# Restore from Git history
```

---

## Decision Tree

### "Which workflow am I?"

```
Do you want to create custom agents?
├─ Yes → Run CAMI from source directory
│        Deploy agent-architect
│        Create agents on demand
│        Contribute via Git
│
└─ No → Install CAMI globally
         Use existing agents only
         Simple consumption model
```

### "Which source priority should I use?"

```
What type of source?
├─ My local experiments → 200 (highest, overrides everything)
├─ Team/company agents → 150 (overrides official)
├─ Official lando-labs → 100 (baseline)
└─ Community/untrusted → 50 (lowest, fallback only)
```

---

## Example Sessions

### Session 1: First-Time Creator

```bash
# Clone and build CAMI
git clone git@github.com:lando-labs/cami.git
cd cami
go build -o cami ./cmd/cami

# Initialize
./cami init
# Select: "Create agents on demand"
# Deploy agent-architect to current directory
# Configure vc-agents as local source

# Create first agent
claude
> @agent-architect I need a Kubernetes specialist agent

# Test it
./cami deploy -a kubernetes-specialist -l ~/k8s-project
cd ~/k8s-project
claude
> @kubernetes-specialist help me design a production cluster

# Refine
cd ~/lando-labs/cami
claude
> @agent-architect add Helm chart patterns to kubernetes-specialist

# Contribute
git add vc-agents/infrastructure/kubernetes-specialist.md
git commit -m "Add Kubernetes specialist agent"
git push origin main
```

### Session 2: Existing Agent Consumer

```bash
# Install globally
go install github.com/lando-labs/cami/cmd/cami@latest

# Initialize
cami init
# Select: "Use existing agents only"

# Browse
cami list
# Shows all official agents

# Deploy to project
cami deploy -a architect,backend,frontend -l ~/my-app

# Use in project
cd ~/my-app
claude
> @architect help me design the API layer

# Update agents periodically
cami update
# Pulls latest from lando-labs/lando-agents
```

### Session 3: Team Custom Agents

```bash
# Set up team source
cami sources add team git git@github.com:acme/agents.git --priority 150
cami sources add official github lando-labs/lando-agents --priority 100

# Deploy agent-architect
cami deploy -a agent-architect -l ~/team-workspace

# Create team-specific agent
cd ~/team-workspace
claude
> @agent-architect create an agent for our internal API gateway
  Context: We use Kong with custom plugins...

# Agent created: acme-api-gateway-specialist.md

# Add to team repo
cp .claude/agents/acme-api-gateway-specialist.md ~/acme-agents/
cd ~/acme-agents
git add acme-api-gateway-specialist.md
git commit -m "Add internal API gateway specialist"
git push

# Team members get it
cami update
# Pulls from acme/agents
```

---

## Summary

**CAMI's Power:**
- ✅ Agent-architect creates sophisticated agents via research
- ✅ Local source (priority 200) enables rapid iteration
- ✅ Multiple sources enable team customization
- ✅ Standard Git workflow for contribution

**Not Like Other Tools:**
- ❌ No weak templates
- ❌ No complex Git wrappers
- ❌ No vendor lock-in

**Core Workflow:**
```
Ask agent-architect → Create agent → Test → Refine → Contribute (Git)
```

**Get Started:**
```bash
git clone git@github.com:lando-labs/cami.git
cd cami
./cami init
claude
> @agent-architect let's create my first agent
```
