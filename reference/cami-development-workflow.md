<!--
AI-Generated Documentation
Created by: Claude Code
Date: 2025-10-05
Purpose: Document CAMI project development workflow and version control practices
-->

# CAMI Development Workflow

This document outlines the specific development workflow for the CAMI project itself, which differs from how CAMI is used in other projects.

## Directory Structure

```
cami/
├── vc-agents/               # SOURCE OF TRUTH - Version controlled agents
│   ├── accessibility-expert.md
│   ├── agent-architect.md
│   ├── ...
│   └── ux.md
│
├── .claude/agents/          # TESTING ONLY - Deployed agents for CAMI development
│   └── (deployed from vc-agents/)
│
├── internal/                # CAMI application code
└── cmd/cami/                # CAMI CLI/TUI entry points
```

## The Golden Rule

**NEVER modify files in `.claude/agents/` directly in the CAMI project.**

### Why This Matters

1. **vc-agents/ is the source of truth** - This is where all agent definitions are version controlled
2. **.claude/agents/ is ephemeral** - These are test deployments that can be deleted and recreated
3. **Version control requires consistency** - All changes must go through vc-agents/ first

## Development Workflow

### Making Changes to Agents

```bash
# 1. ALWAYS modify vc-agents/ (source)
vim vc-agents/backend.md

# 2. Increment version in frontmatter
version: "1.1.0" → "1.2.0"

# 3. Test by redeploying to .claude/agents/
./cami deploy --agents backend --location . --overwrite

# 4. Verify the changes work
# (Use the deployed agent in Claude Code)

# 5. Commit to git
git add vc-agents/backend.md
git commit -m "Update backend agent: add new methodology"
```

### Creating New Agents

```bash
# 1. Create in vc-agents/ (source)
vim vc-agents/new-specialist.md

# 2. Start with version "1.0.0"

# 3. Deploy to test
./cami deploy --agents new-specialist --location . --overwrite

# 4. Test and iterate (always edit vc-agents/, redeploy to test)

# 5. Commit when ready
git add vc-agents/new-specialist.md
git commit -m "Add new-specialist agent v1.0.0"
```

### Resetting Test Environment

```bash
# Clean deployed agents
rm -rf .claude/agents

# Redeploy from source
./cami deploy --agents <agent-list> --location . --overwrite
```

## Version Management

### When to Increment Versions

- **v1.0.0 → v1.0.1** (patch): Bug fixes, typos, minor clarifications
- **v1.0.0 → v1.1.0** (minor): New capabilities, additional documentation, methodology improvements
- **v1.0.0 → v2.0.0** (major): Breaking changes, complete rewrites, fundamental approach changes

### Current Version Standard

All agents are currently at **v1.1.0** (as of 2025-10-05) after adding:
- Documentation Strategy section
- AI-Generated Documentation Marking convention
- reference/ folder pattern

## How Other Projects Use CAMI

In contrast to the CAMI project itself, when CAMI is used to deploy agents to other projects:

```bash
# In ANY other project (not CAMI)
cd ~/my-app

# Deploy agents - this MODIFIES .claude/agents/ (correct!)
cami deploy --agents backend,frontend --location .

# Update docs
cami update-docs --location .
```

**In other projects**:
- ✅ `.claude/agents/` IS the source of truth
- ✅ Modifying deployed agents is expected
- ✅ CAMI manages the files

**In the CAMI project**:
- ✅ `vc-agents/` IS the source of truth
- ❌ NEVER modify `.claude/agents/` directly
- ✅ CAMI deploys from `vc-agents/` for testing

## Common Mistakes to Avoid

### ❌ Wrong: Editing Deployed Agents in CAMI

```bash
# In CAMI project - DON'T DO THIS
vim .claude/agents/backend.md  # ❌ Changes will be lost!
```

### ✅ Correct: Edit Source, Then Deploy

```bash
# In CAMI project - DO THIS
vim vc-agents/backend.md       # ✅ Edit source
./cami deploy --agents backend --location . --overwrite  # ✅ Deploy to test
```

### ❌ Wrong: Forgetting to Increment Version

```bash
# Making changes without version bump
vim vc-agents/backend.md
# (make changes but forget to update version field)
git commit  # ❌ Version doesn't reflect changes
```

### ✅ Correct: Always Increment Version

```bash
vim vc-agents/backend.md
# Change: version: "1.1.0" → "1.2.0"
git commit -m "Update backend v1.2.0: add new feature"  # ✅
```

## CI/CD Considerations (Future)

When CAMI gets automated testing:

```yaml
# Example: .github/workflows/test-agents.yml
on: [push]
jobs:
  test:
    steps:
      - name: Validate agent YAML
        run: ./cami list --output json  # Should load all 23 agents

      - name: Deploy test suite
        run: ./cami deploy --agents qa --location ./test-project

      - name: Version check
        run: |
          # Ensure all agents have versions
          # Ensure no duplicate names
```

## Summary

| Aspect | CAMI Project | Other Projects |
|--------|--------------|----------------|
| Source of truth | `vc-agents/` | `.claude/agents/` |
| Edit agents | Only in `vc-agents/` | Directly in `.claude/agents/` OK |
| Deploy command | For testing only | Primary workflow |
| Version control | Git tracks `vc-agents/` | Git tracks `.claude/agents/` |
| CAMI's role | Source code + tests | Agent deployment manager |

## Questions?

If you're unsure whether to edit `vc-agents/` or `.claude/agents/`, ask:

**"Am I working on the CAMI project itself?"**
- Yes → Edit `vc-agents/`, deploy to `.claude/agents/` for testing
- No → Deploy with CAMI, edit `.claude/agents/` as needed

---

Last updated: 2025-10-05 (v1.0)
