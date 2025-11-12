# CAMI Contribution Philosophy

**Version:** 1.0
**Date:** 2025-11-10

---

## Core Principle

> **CAMI facilitates consumption, not contribution.**
>
> CAMI makes it **dead simple** to discover, deploy, and use agents.
> Contributing agents back uses **standard Git workflows** that developers already know.

---

## What CAMI Does (Pull Model)

### ✅ Discovery
```bash
cami list                    # Browse all agents
cami search terraform        # Find specific agents
cami show architect          # View agent details
```

### ✅ Consumption
```bash
cami update                  # Pull latest agents from sources
cami deploy -a terraform -l ~/project  # Deploy to project
```

### ✅ Organization
```bash
cami sources list            # Show configured sources
cami sources add ...         # Add new source
cami status                  # Check source states
```

### ✅ Guidance
```bash
cami help contribute         # Show contribution guide
cami diff architect          # Compare local vs official
```

---

## What CAMI Doesn't Do (Push Model)

### ❌ Git Operations
```bash
# Use Git directly:
git add vc-agents/new-agent.md
git commit -m "Add new agent"
git push origin main
```

### ❌ Pull Requests
```bash
# Use GitHub CLI:
gh pr create --title "Add Terraform agent" --body "..."
```

### ❌ Authentication
```bash
# Use Git credentials:
git config --global user.name "..."
gh auth login
```

### ❌ Conflict Resolution
```bash
# Use Git:
git merge
git rebase
```

---

## Why This Separation?

### Avoids
- ❌ Building half-assed Git wrapper
- ❌ Million edge cases (auth, conflicts, merge strategies)
- ❌ Reinventing GitHub CLI
- ❌ Complex Git state management
- ❌ Platform lock-in (GitHub, GitLab, Bitbucket)

### Enables
- ✅ Clean, focused tool (one job, done well)
- ✅ Standard workflows (users know Git already)
- ✅ LLM-friendly (Claude knows Git/gh)
- ✅ Flexibility (works with any Git workflow)
- ✅ Maintainability (we don't manage Git)
- ✅ Transferable skills (learn Git, not CAMI-specific commands)

---

## User Workflows

### Workflow 1: Agent Consumer (Pull Only)

**Goal:** Use existing agents, never contribute

```bash
# Install CAMI
go install github.com/lando-labs/cami/cmd/cami@latest

# Initialize
cami init
# Selects: "Use existing agents"

# Browse and deploy
cami list
cami deploy -a architect,frontend -l ~/my-project

# Update periodically
cami update
```

**CAMI involvement:** 100% (discovery, deployment, updates)
**Git involvement:** 0% (never touches Git)

---

### Workflow 2: Agent Developer (Pull + Create + Push)

**Goal:** Create custom agents, potentially contribute back

```bash
# Work in CAMI/lando-agents repo
cd ~/lando-labs/lando-agents

# Use agent-architect to create
claude
> @agent-architect create a Terraform specialist agent

# Agent created in vc-agents/infrastructure/terraform.md

# Test it
./cami deploy -a terraform -l ~/test-project

# Standard Git workflow
git status
git add vc-agents/infrastructure/terraform.md
git commit -m "Add Terraform specialist agent"
git push origin feature/terraform-agent

# Create PR (optional)
gh pr create --title "Add Terraform specialist" \
  --body "Agent for managing Terraform IaC workflows"
```

**CAMI involvement:** Discovery, testing, deployment
**Git involvement:** All contribution steps (add, commit, push, PR)

---

### Workflow 3: Fork Contributor (External User)

**Goal:** Contribute to lando-labs/lando-agents from fork

```bash
# Fork on GitHub
gh repo fork lando-labs/lando-agents --clone

cd lando-agents

# Point CAMI to local fork
cami sources add my-fork local ~/lando-agents --priority 200

# Create agent with agent-architect
cami deploy -a agent-architect -l .
claude
> @agent-architect create a GraphQL specialist

# Standard fork workflow
git checkout -b add-graphql-specialist
git add vc-agents/specialized/graphql-specialist.md
git commit -m "Add GraphQL specialist agent"
git push origin add-graphql-specialist

# Create PR to upstream
gh pr create --repo lando-labs/lando-agents \
  --title "Add GraphQL specialist agent" \
  --body "Agent specialized in GraphQL API design and optimization"
```

**CAMI involvement:** Local source management, deployment
**Git/GitHub involvement:** All contribution mechanics

---

## Status Awareness (Optional Helper)

CAMI can **show** Git status if source is a Git repo, but doesn't **manage** it.

### `cami status`

**Implementation:**
```go
// internal/sources/git_status.go
func (s *LocalSource) GitStatus() (*GitInfo, error) {
    if !s.IsGitRepo() {
        return nil, nil  // Not a Git repo, skip
    }

    // Run: git status --porcelain
    // Parse output
    // Return summary
}
```

**Output:**
```bash
$ cami status

Agent Sources:
  ✓ official (lando-labs/lando-agents)
    Type: github
    Status: Up to date
    Last updated: 2 hours ago

  ⚠ local (~/lando-labs/lando-agents)
    Type: local
    Git status: 2 uncommitted changes
      M vc-agents/core/frontend.md
      ?? vc-agents/infrastructure/terraform.md

    Tip: Run 'git status' in ~/lando-labs/lando-agents for details

Deploy Locations:
  my-project (/Users/lando/my-project)
    6 agents deployed
    Last updated: 1 day ago
```

**What it does:**
- Shows Git status if available
- Provides helpful tips
- Doesn't interpret or manage

**What it doesn't do:**
- Stage files
- Commit changes
- Push to remote
- Resolve conflicts

---

## LLM-Friendly Contribution Guide

### For Claude Code (and other LLMs)

When helping users contribute agents, follow this pattern:

**✅ Do:**
```
User: "Can you create a Terraform agent and contribute it?"

Claude:
1. Let me create the agent using agent-architect
   [Creates agent in vc-agents/infrastructure/terraform.md]

2. I've created the agent. To contribute it to lando-labs/lando-agents:

   git add vc-agents/infrastructure/terraform.md
   git commit -m "Add Terraform specialist agent"
   git push origin main

   Then create a PR:
   gh pr create --title "Add Terraform specialist" \
     --body "Agent specialized in Terraform IaC workflows"
```

**❌ Don't:**
```
User: "Can you contribute this agent?"

Claude:
Let me use CAMI to push this...
  cami push terraform --message "..."  ← Doesn't exist!
```

### Include in lando-labs/lando-agents CLAUDE.md

```markdown
## Contributing Agents

### For LLMs Helping Users

When Claude Code (or other LLMs) help users contribute:

1. **Create/edit agent** in `vc-agents/category/name.md`
2. **Test locally** with CAMI deployment
3. **Guide user through Git workflow:**
   ```bash
   git status
   git add vc-agents/category/name.md
   git commit -m "Add/update agent description"
   git push origin branch-name
   ```
4. **Suggest PR creation:**
   ```bash
   gh pr create --title "..." --body "..."
   ```

**Important:**
- Use **standard Git commands** (not hypothetical CAMI commands)
- Use **`gh` CLI** for PRs (if available)
- Keep it **simple and standard**
- Don't invent CAMI features that don't exist

### For Human Contributors

1. Fork or clone `lando-labs/lando-agents`
2. Create/edit agent files in `vc-agents/`
3. Test with CAMI: `cami deploy -a your-agent -l ~/test-project`
4. Standard Git workflow: `git add`, `git commit`, `git push`
5. Create PR via GitHub UI or `gh pr create`

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.
```

---

## Implementation in CAMI

### Commands to Add

```bash
cami status                  # Show source states (includes Git status)
cami diff <agent>            # Compare local vs official agent
cami help contribute         # Display contribution guide
```

### Commands to NOT Add

```bash
cami push                    # ❌ Use git push
cami pr                      # ❌ Use gh pr create
cami commit                  # ❌ Use git commit
cami fork                    # ❌ Use gh repo fork
```

---

## Documentation Strategy

### CAMI Docs
- Focus on **consumption** (discovery, deployment, updates)
- Link to Git/GitHub docs for contribution
- Include "Contributing Agents" guide with Git workflow

### lando-agents Docs
- Standard `CONTRIBUTING.md` with Git workflow
- LLM-friendly guidance in `CLAUDE.md`
- Examples of good PRs

---

## Benefits of This Approach

### For Users
- **Learn transferable skills** (Git, not CAMI-specific)
- **Use familiar tools** (git, gh)
- **Flexible workflows** (works with any Git setup)
- **Clear separation** (CAMI = consume, Git = contribute)

### For CAMI Development
- **Smaller scope** (no Git management)
- **Fewer edge cases** (Git handles complexity)
- **Easier maintenance** (don't reinvent Git)
- **Better focus** (one job, done well)

### For Open Source Community
- **Standard patterns** (familiar to open source contributors)
- **Platform agnostic** (works with GitHub, GitLab, Bitbucket)
- **LLM-friendly** (Claude already knows Git)

---

## Summary

**CAMI's Role:**
- ✅ Make consumption **effortless** (discovery, deployment, updates)
- ✅ Provide **guidance** for contribution (docs, examples)
- ✅ Show **status** (what's changed locally)

**Not CAMI's Role:**
- ❌ Replace Git (users know it already)
- ❌ Manage PRs (gh CLI does this well)
- ❌ Handle authentication (Git credentials work)

**Philosophy:**
> **CAMI is a one-way street with a roadmap.**
>
> It pulls agents down effortlessly. It guides you back upstream with standard tools.

---

## Next Steps

1. ✅ Document this philosophy
2. ⬜ Update open source strategy docs
3. ⬜ Create `CONTRIBUTING.md` for lando-agents
4. ⬜ Add LLM guidance to lando-agents `CLAUDE.md`
5. ⬜ Implement `cami status` (shows Git status)
6. ⬜ Implement `cami diff` (compares versions)
7. ⬜ Add contribution guide to `cami help`
