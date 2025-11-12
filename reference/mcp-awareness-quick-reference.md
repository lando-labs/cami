<!--
AI-Generated Documentation
Created by: agent-architect
Date: 2025-11-03
Purpose: Quick reference guide for MCP-aware agent design
-->

# MCP Awareness Quick Reference

## Agent Prompt Template: MCP Section

Add this section after "Core Philosophy" and before "Three-Phase Methodology":

```markdown
## MCP Awareness: Specialized Tool Priority

You have access to Model Context Protocol (MCP) servers that extend your capabilities. As a [AGENT_DOMAIN] specialist, you should prioritize MCPs in this order:

### Primary MCPs (Use First)
**[CATEGORY_NAME]**: [Description]
- **When to use**: [Specific scenarios]
- **Tool patterns**: `[pattern]__*`
- **Resource patterns**: `[scheme]://`
- **Examples**: [MCP examples]

### Secondary MCPs (Use When Relevant)
**[CATEGORY_NAME]**: [Description]
- **When to use**: [Scenarios]
- **Tool patterns**: `[pattern]__*`

### MCP Discovery Protocol

At the start of each task:
1. Invoke `ListMcpResourcesTool` to discover available MCPs
2. Identify MCPs matching your primary categories
3. Prefer domain-specific MCP tools over generic alternatives
4. Fall back to Claude Code built-in tools if no suitable MCP exists

### MCP Usage Guidelines

**DO**:
- ✓ Discover MCPs at task start
- ✓ Use domain-specific MCPs first
- ✓ Document which MCPs were used
- ✓ Fall back gracefully to built-in tools

**DON'T**:
- ✗ Hardcode MCP server names
- ✗ Assume MCPs are available
- ✗ Skip MCP discovery for specialized tasks
```

## MCP Categories

| Category | Purpose | Tool Patterns | Use Cases |
|---|---|---|---|
| **Design & UI** | Design systems, components, tokens | `design_system__*`, `component__*`, `token__*` | UI development, design system work |
| **Development** | Files, git, code analysis, browser | `mcp__filesystem__*`, `mcp__git__*`, `mcp__playwright__*` | Code operations, testing |
| **Infrastructure** | Cloud, K8s, databases, monitoring | `mcp__firebase__*`, `mcp__k8s__*`, `mcp__aws__*` | Deployment, infrastructure |
| **Knowledge** | Docs, API specs, wikis | `mcp__context7__*`, `mcp__docs__*` | Research, learning |
| **Data & Analytics** | Databases, warehouses, ETL | `mcp__postgres__*`, `mcp__bigquery__*` | Data operations |
| **Testing & QA** | Test frameworks, coverage, perf | `mcp__jest__*`, `mcp__axe__*`, `mcp__lighthouse__*` | Testing, validation |
| **Communication** | Slack, email, calendar | `mcp__slack__*`, `mcp__calendar__*` | Notifications, scheduling |
| **AI & ML** | Models, vectors, MLOps | `mcp__openai__*`, `mcp__pinecone__*` | ML workflows |
| **Utilities** | Time, calculations, validation | `mcp__time__*`, `mcp__calc__*` | General utilities |

## Agent-MCP Affinity

| Agent Type | Primary MCPs | Secondary MCPs |
|---|---|---|
| **Design Agents** | Design & UI | Knowledge, Testing & QA |
| **Frontend Agents** | Design & UI, Development | Testing & QA, Knowledge |
| **Backend Agents** | Data & Analytics, Development | Infrastructure, Knowledge |
| **Infrastructure Agents** | Infrastructure, Development | Testing & QA, Communication |
| **Testing Agents** | Testing & QA, Development | Infrastructure, Design & UI |
| **Research Agents** | Knowledge, Utilities | Communication, Development |
| **AI/ML Agents** | AI & ML, Data & Analytics | Infrastructure, Development |
| **Mobile Agents** | Development, Design & UI | Testing & QA, Infrastructure |
| **Security Agents** | Infrastructure, Development | Testing & QA, Data & Analytics |

## 4-Step MCP Discovery Pattern

### Step 1: Discover
```
Invoke ListMcpResourcesTool
Parse tool names (pattern: mcp__[server]__[tool])
Build MCP registry
```

### Step 2: Categorize
```
Match MCPs to categories (Design & UI, Infrastructure, etc.)
Tag with capabilities
Note tool patterns
```

### Step 3: Prioritize
```
Rank by affinity (Primary=3, Secondary=2, Tertiary=1)
Sort by: Affinity > Tool Count > Alphabetical
```

### Step 4: Select & Use
```
Match task to MCP category
Choose highest-priority MCP
Invoke tool
Fall back if needed
```

## Integration Checklist

When adding MCP awareness to an agent:

- [ ] Add "MCP Awareness" section after "Core Philosophy"
- [ ] Define 2-3 Primary MCP categories for this agent
- [ ] Define 1-2 Secondary MCP categories
- [ ] Include MCP Discovery Protocol
- [ ] Include MCP Usage Guidelines (DO/DON'T)
- [ ] Update Phase 1 to include MCP discovery
- [ ] Mark MCP-enhanced steps with "**MCP-Enhanced**:" prefix
- [ ] Add MCP tools to "Tools" sections
- [ ] Update Quality Standards to mention MCP usage
- [ ] Add MCP checks to Self-Verification Checklist
- [ ] Include concrete MCP usage examples (optional but helpful)

## Common Patterns

### Pattern 1: Design System Access
```
Agent: frontend, designer, design-system-specialist
MCP Category: Design & UI
Typical Flow:
  1. Discover design system MCPs
  2. Fetch design tokens: design_system__get_tokens()
  3. Get components: design_system__get_component()
  4. Build UI with consistent tokens
```

### Pattern 2: Cloud Deployment
```
Agent: devops, deploy, gcp-firebase
MCP Category: Infrastructure
Typical Flow:
  1. Discover cloud provider MCPs (firebase, aws, k8s)
  2. Initialize service: mcp__firebase__firebase_init()
  3. Deploy: mcp__firebase__firebase_deploy()
  4. Monitor with datadog-mcp or similar
```

### Pattern 3: Documentation Research
```
Agent: research-synthesizer, all agents
MCP Category: Knowledge
Typical Flow:
  1. Discover knowledge MCPs (context7)
  2. Resolve library: mcp__context7__resolve-library-id()
  3. Fetch docs: mcp__context7__get-library-docs()
  4. Synthesize findings
```

### Pattern 4: Accessibility Validation
```
Agent: frontend, qa, accessibility-expert
MCP Category: Testing & QA
Typical Flow:
  1. Discover accessibility MCPs (axe-core, lighthouse)
  2. Scan component: mcp__axe__scan_component()
  3. Fix violations
  4. Re-validate
  5. Document compliance
```

### Pattern 5: Graceful Fallback
```
Any Agent:
Typical Flow:
  1. Discover MCPs
  2. No suitable MCP found
  3. Fall back to Claude Code built-in tools (Read, Write, Bash)
  4. Document manual approach
  5. Suggest MCP that would help (optional)
```

## DO's and DON'Ts

### DO
✓ Discover MCPs at start of task with `ListMcpResourcesTool`
✓ Prioritize domain-specific MCPs (design agents → design MCPs)
✓ Use MCP tools for specialized capabilities
✓ Fall back gracefully to built-in tools
✓ Document which MCPs were used
✓ Check for multiple MCPs in same category
✓ Use knowledge MCPs (context7) for documentation
✓ Validate with testing MCPs when available

### DON'T
✗ Hardcode MCP server names (e.g., "lando-labs-design-system")
✗ Assume MCPs are always available
✗ Skip MCP discovery for specialized tasks
✗ Use MCPs outside your domain without justification
✗ Invoke MCPs for tasks better suited to built-in tools
✗ Fail task if MCP unavailable - fall back gracefully
✗ Use infrastructure MCPs if you're a design agent (without reason)
✗ Ignore MCP capabilities when they match your task

## Example: Frontend Agent MCP Usage

```markdown
Task: Build a Button component

Step 1: Discover MCPs
  → ListMcpResourcesTool
  → Found: design-system-mcp, filesystem, playwright, context7

Step 2: Identify relevant MCPs
  → design-system-mcp (Primary: Design & UI)
  → filesystem (Primary: Development)
  → context7 (Secondary: Knowledge)

Step 3: Fetch design tokens
  → design-system-mcp has: design_system__get_tokens
  → Invoke: design_system__get_tokens(category: "colors,spacing")
  → Result: { primary: "#3b82f6", spacing: { md: "16px" } }

Step 4: Check existing components
  → design-system-mcp has: design_system__get_component
  → Invoke: design_system__get_component(name: "Button")
  → Result: existing Button implementation or "not found"

Step 5: Build component
  → Use tokens from design-system-mcp
  → Match patterns from existing components
  → Ensure consistency

Step 6: Validate
  → Use playwright MCP for browser testing (if needed)
  → Document: "Built with design-system-mcp for consistency"
```

## Rollout Phases

### Phase 1: Core Agents (Week 1-2)
- frontend, backend, devops, designer, qa
- Add MCP awareness sections
- Test with available MCPs

### Phase 2: Specialized Agents (Week 3-4)
- react-specialist, design-system-specialist, gcp-firebase, api-integrator
- Add domain-specific MCP priorities
- Test cross-MCP workflows

### Phase 3: Research & AI (Week 5-6)
- research-synthesizer, ai-ml-specialist, use-case-specialist
- Emphasize Knowledge and AI/ML MCPs
- Test documentation workflows

### Phase 4: Infrastructure & Security (Week 7-8)
- security-specialist, performance-optimizer, deploy
- Add Infrastructure MCP priorities
- Test cloud provider integrations

## Files to Create/Update

### New Files
- `reference/mcp-awareness-architecture.md` - Full architecture
- `reference/mcp-awareness-example-frontend.md` - Concrete example
- `reference/mcp-awareness-quick-reference.md` - This file
- `reference/mcp-configuration.md` - Project MCP config

### Updated Files
- `.claude/agents/frontend.md` - Add MCP awareness
- `.claude/agents/backend.md` - Add MCP awareness
- `.claude/agents/devops.md` - Add MCP awareness
- `.claude/agents/designer.md` - Add MCP awareness
- `.claude/agents/qa.md` - Add MCP awareness
- (And all other agents in phases 2-4)

### CLAUDE.md Updates
Add section:
```markdown
## Model Context Protocol (MCP) Servers

Agents automatically discover and use MCP servers for specialized capabilities:

- **Design Agents**: Design system MCPs, component libraries
- **Infrastructure Agents**: Cloud provider MCPs (Firebase, AWS, K8s)
- **All Agents**: Knowledge MCPs (context7) for documentation

See `reference/mcp-configuration.md` for configured MCPs.
See `reference/mcp-awareness-architecture.md` for full details.
```

## Testing Checklist

- [ ] Agent discovers MCPs with `ListMcpResourcesTool`
- [ ] Agent prioritizes domain-relevant MCPs
- [ ] Agent uses MCP tools correctly
- [ ] Agent falls back gracefully when MCPs unavailable
- [ ] Agent documents which MCPs were used
- [ ] Multiple MCPs in same category handled correctly
- [ ] Cross-category fallback works (primary → secondary → built-in)
- [ ] Agent doesn't hardcode MCP names
- [ ] Agent handles MCP errors gracefully

---

**For full details**: See `reference/mcp-awareness-architecture.md`
**For example**: See `reference/mcp-awareness-example-frontend.md`
**For project MCPs**: See `reference/mcp-configuration.md` (create as needed)
