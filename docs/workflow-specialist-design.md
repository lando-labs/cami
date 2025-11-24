# Workflow Specialist Design - Flexibility Spectrum

**Date**: 2025-01-23
**Status**: Active Design Decision
**Version**: v0.4.0+

## Problem Statement

Workflow Specialist agents (Task Automators) need to execute user-defined processes, but workflows vary widely in their rigidity and automation level. Some users want fully automated scripts, others want checklists with guidelines. We need a flexible system that supports different workflow types.

## Design Decision

**Support a spectrum of workflow flexibility (levels 2-4), letting agent-architect decide the appropriate level based on the workflow provided by the user.**

## The Workflow Flexibility Spectrum

### Level 1: Fully Hardcoded ❌ (Not Recommended)
```yaml
workflow:
  - step: "Run tests"
    command: "npm test"
  - step: "Build image"
    command: "docker build -t myapp:latest ."
```

**Why not recommended**: Too brittle, no reusability across projects

---

### Level 2: Parameterized ✅ (Recommended for automation)
```yaml
parameters:
  namespace:
    type: string
    required: true
    description: "Kubernetes namespace to check"
  pod_pattern:
    type: string
    required: false
    default: "*"
    description: "Pod name pattern to filter"

workflow:
  - step: "List pods"
    command: "kubectl get pods -n {{namespace}} | grep {{pod_pattern}}"
  - step: "Check logs"
    command: "kubectl logs {{pod_name}} -n {{namespace}}"
```

**Use cases**:
- Kubernetes pod health checks
- Database backup workflows
- Deployment automation
- CI/CD pipeline steps

**Input gathering**: CAMI/Claude gathers parameter values from user before invoking agent

---

### Level 3: Semi-Structured ✅ (Recommended for guidelines)
```yaml
workflow:
  - step: "Verify tests pass"
    action: "Run the project's test suite and ensure all tests pass"
    success_criteria: "All tests green, no failures"

  - step: "Build artifact"
    action: "Build production artifact using project's build system"
    success_criteria: "Build completes without errors, artifact created"
```

**Use cases**:
- Pre-deployment checklists
- Code review workflows
- Security audit processes
- Onboarding procedures

**Execution**: Agent interprets actions, uses available tools, reports success/failure

---

### Level 4: Fully Flexible ✅ (Recommended for high-level checklists)
```markdown
## Release Workflow

1. ✓ All tests passing
2. ✓ Code reviewed and approved
3. ✓ Changelog updated
4. ✓ Version bumped
5. ✓ Deploy to staging
6. ✓ Smoke tests pass
7. ✓ Deploy to production
```

**Use cases**:
- Release checklists
- Incident response procedures
- Project setup workflows
- Documentation processes

**Execution**: Agent works through checklist, uses judgment on implementation

---

## Workflow Gathering Process

When CAMI/Claude gathers workflow details for a Workflow Specialist agent, ask:

### 1. "What's the workflow?"
Offer three input methods:
- **Describe it**: Walk through steps conversationally
- **Provide a file**: Share markdown checklist, shell script, YAML workflow
- **Point to docs**: Link or file path to procedure documentation

### 2. "Does this workflow need inputs each time?"
Examples:
- Namespace, environment, pod name (Kubernetes workflows)
- Database name, backup location (backup workflows)
- Version number, release notes (deployment workflows)
- None (fully self-contained workflows)

If yes, define parameters:
```yaml
parameters:
  param_name:
    type: string|number|boolean
    required: true|false
    default: "value"  # if not required
    description: "What this parameter is for"
```

### 3. "Should steps be automated commands or guidelines?"
- **Automated commands** → Level 2 (Parameterized)
  - User provides shell script or specific commands
  - Agent executes exact commands with parameter substitution

- **Guidelines** → Level 3 (Semi-Structured)
  - User provides action descriptions
  - Agent uses available tools to accomplish each step

- **Checklist** → Level 4 (Fully Flexible)
  - User provides high-level steps
  - Agent exercises judgment on implementation

## Agent-Architect's Role

Agent-architect detects the appropriate flexibility level based on input:

**Signals for Level 2 (Parameterized)**:
- User provides shell script with variables
- Specific commands mentioned (kubectl, docker, npm, etc.)
- Workflow file is .sh, .bash, or has command blocks

**Signals for Level 3 (Semi-Structured)**:
- User provides action descriptions without specific commands
- Mentions "verify", "ensure", "check" without exact commands
- Workflow needs project-specific adaptation

**Signals for Level 4 (Fully Flexible)**:
- User provides markdown checklist
- High-level steps without implementation details
- Focus on outcomes rather than methods

## Implementation in Agent Frontmatter

Workflow Specialists can include workflow definition in frontmatter or content:

### Example: Parameterized Workflow (Level 2)

```yaml
---
name: k8s-pod-checker
version: "1.0.0"
description: Check Kubernetes pod health following a diagnostic workflow
class: workflow-specialist
specialty: kubernetes-operations
parameters:
  namespace:
    type: string
    required: true
    description: "Kubernetes namespace to check"
  pod_pattern:
    type: string
    required: false
    default: "*"
    description: "Pod name pattern to filter (optional)"
---

# Kubernetes Pod Health Checker

## Workflow

1. **List Pod Status**
   ```bash
   kubectl get pods -n {{namespace}} -o wide | grep {{pod_pattern}}
   ```
   Success: Pods listed
   Failure: Connection issues or invalid namespace

2. **Check Pod Events**
   ```bash
   kubectl describe pod {{pod_name}} -n {{namespace}}
   ```
   Success: Events displayed
   Failure: Pod not found

[... rest of workflow ...]
```

### Example: Semi-Structured Workflow (Level 3)

```yaml
---
name: pre-deployment-checklist
version: "1.0.0"
description: Execute pre-deployment verification checklist
class: workflow-specialist
specialty: deployment-verification
---

# Pre-Deployment Checklist

## Workflow

1. **Verify tests pass**
   - Run the project's test suite
   - Success: All tests green, no failures
   - Failure: Any test failures or errors

2. **Build artifact**
   - Build production artifact using project's build system
   - Success: Build completes, artifact created
   - Failure: Build errors or warnings

[... rest of checklist ...]
```

## Key Principles

1. **Workflow Specialists execute, they don't decide** - Focus on execution fidelity, not strategic decisions

2. **Parameters are gathered upfront** - No interactive Q&A mid-workflow (breaks execution flow)

3. **Clear success/failure criteria** - Each step should have objective completion criteria

4. **Flexibility serves users** - Support different workflow types rather than forcing one model

5. **Agent-architect auto-detects** - Users don't need to understand levels, they just provide workflows

## Future Enhancements

### Workflow Import from Files
```
User: "Create a workflow agent from this script"
User: *provides deploy.sh*

CAMI: *reads script, identifies parameters, suggests workflow structure*
"I found these parameters: ENVIRONMENT, VERSION, REGION
 Should I create a parameterized workflow agent?"
```

### Workflow Validation
- Syntax validation for commands
- Parameter type checking
- Success criteria verification

### Workflow Execution Reporting
```
✓ Step 1: List pods (completed in 0.5s)
✓ Step 2: Check events (completed in 0.3s)
⚠ Step 3: Get logs (warning: no logs available)
✓ Step 4: Resource check (completed in 0.2s)

Summary: 4/4 steps completed (1 warning)
```

## Related Documents

- [WORKFLOW_ENHANCEMENT_PLAN_V2.md](/Users/lando/lando-labs/cami/WORKFLOW_ENHANCEMENT_PLAN_V2.md) - Original enhancement plan
- [Agent Classification System](/Users/lando/lando-labs/cami/CLAUDE.md#agent-classification-system) - Developer docs
- [User CLAUDE.md](/Users/lando/lando-labs/cami/install/templates/CLAUDE.md) - User-facing workflow guidance
