# Agent Creation Interview System

**Version:** 1.0
**Date:** 2025-11-10

---

## Problem

**Current (weak):**
```
User: "Create a Terraform agent"
Claude: [Creates generic Terraform agent]
```

**Better:**
```
User: "Create a Terraform agent"
Claude: [Asks discovery questions to understand philosophy, constraints, workflow]
```

---

## Philosophy

> **Agent-architect's strength is research + philosophy, not template filling.**

Good agents have:
1. **Clear philosophy** (how should this domain be approached?)
2. **Specific constraints** (what NOT to do?)
3. **Workflow patterns** (how does work actually get done?)
4. **Context awareness** (what tools, MCPs, environment?)
5. **Opinionated guidance** (best practices, not just capability lists)

---

## CAMI's Role

CAMI provides structure for the interview process via:

### 1. MCP Tool: `agent_create` (New)

**Purpose:** Initiate guided agent creation

**Parameters:**
- `initial_request` (string): User's initial description
- `interactive` (bool): Enable interview mode (default: true)

**What it does:**
1. Validates initial request
2. Returns interview questions based on domain
3. Tracks interview state
4. Guides agent-architect through discovery
5. Provides final context for agent creation

### 2. Agent-Architect Enhancement

**Add to `.claude/agents/agent-architect.md`:**

```markdown
## Agent Creation Interview Process

When creating an agent, DO NOT immediately create based on a simple request.
Instead, conduct a discovery interview to understand:

1. **Philosophy & Approach**
2. **Constraints & Boundaries**
3. **Workflow Patterns**
4. **Context & Environment**
5. **User's Experience Level**

### Interview Framework

Use the CAMI `agent_create` MCP tool to guide the interview.

#### Step 1: Initial Understanding

User provides initial request (e.g., "Create a Terraform agent")

Ask:
- "What will you primarily use this agent for?"
- "What's your philosophy on [domain]?"
- "What should this agent NEVER do?"

#### Step 2: Domain Research

Based on answers, research:
- Best practices in the domain
- Common pitfalls
- Tool ecosystem
- Related MCPs (if applicable)

#### Step 3: Philosophy Refinement

Present research findings and ask:
- "Which approach resonates with you?"
- "What's your team's stance on [key decision]?"
- "Any specific patterns or anti-patterns?"

#### Step 4: Workflow Mapping

Understand how work gets done:
- "Walk me through a typical workflow"
- "What tools do you use?"
- "What's your approval/review process?"

#### Step 5: Context Gathering

Environmental factors:
- "What MCPs are available in your setup?"
- "What's your experience level with [domain]?"
- "Are you working solo or with a team?"

#### Step 6: Create Agent

Now you have:
- ✅ Clear philosophy
- ✅ Specific constraints
- ✅ Workflow patterns
- ✅ Context awareness
- ✅ User's needs

Create a **sophisticated, opinionated agent** based on this understanding.

### Example: Terraform Agent Interview

**Bad (no interview):**
```
User: "Create a Terraform agent"
[Creates generic agent with Terraform commands]
```

**Good (with interview):**
```
User: "Create a Terraform agent"

You: "I'll create a Terraform specialist agent. Let me understand your approach first.

What will you primarily use this agent for?
  - Infrastructure provisioning
  - State management
  - Module development
  - Multi-environment deployments
  - All of the above

User: "Multi-environment deployments"

You: "Got it. What's your philosophy on environment separation?"

User: "We use workspaces and keep state in S3, separate accounts for prod"

You: "Perfect. What should this agent NEVER do?"

User: "Never directly modify production without approval"

You: "Understood. Walk me through your typical deployment workflow."

User: "We use Atlantis for PR-based workflows. Plan on PR, apply on merge."

You: "Excellent context. What MCPs do you have available?"

[Checks available MCPs]

You: "I see you have the GitHub MCP. Should the agent integrate with PR workflows?"

User: "Yes, that would be helpful"

You: "Great! Based on our conversation, I'll create a Terraform agent with:

**Philosophy:**
- Multi-environment safety first
- Workspace-based separation
- PR-driven changes (Atlantis integration)
- Production protection

**Key Features:**
- Workspace awareness
- State management guidance
- GitHub PR integration
- Approval gates for production
- S3 backend patterns

**Constraints:**
- Never direct production changes
- Always plan before apply
- Require reviews for prod

Sound good?"

User: "Perfect"

[Creates sophisticated, contextual agent]
```

### Interview Templates by Domain

CAMI can provide domain-specific interview templates:

#### Infrastructure/DevOps (Terraform, K8s, etc.)
1. Environment strategy? (workspaces, accounts, clusters)
2. State management philosophy?
3. Approval/review process?
4. Disaster recovery stance?
5. Multi-region considerations?

#### Backend Development
1. API design philosophy? (REST, GraphQL, gRPC)
2. Database approach? (ORM, query builder, raw SQL)
3. Error handling strategy?
4. Testing philosophy?
5. Authentication/authorization patterns?

#### Frontend Development
1. State management preference?
2. Styling approach? (CSS-in-JS, modules, Tailwind)
3. Component philosophy? (Atomic, feature-based)
4. Testing strategy?
5. Accessibility priority?

#### Data/Analytics
1. Data quality philosophy?
2. Pipeline architecture? (batch, streaming, hybrid)
3. Schema evolution strategy?
4. Testing approach for data?
5. Monitoring/observability stance?

#### AI/ML
1. Model development philosophy?
2. Training vs inference priorities?
3. Experimentation tracking?
4. Production deployment strategy?
5. Monitoring/drift detection?

[More templates...]
```

---

## Implementation Plan

### Phase 2A: Add to MCP Tools (New Priority)

**Timeline:** 1 week (within Phase 2)

**Tasks:**

1. **Create `agent_create` MCP tool**

```go
// cmd/cami-mcp/agent_create.go

type AgentCreationRequest struct {
    InitialRequest string            `json:"initial_request"`
    Domain         string            `json:"domain"`          // derived or specified
    Interactive    bool              `json:"interactive"`     // default: true
    Responses      map[string]string `json:"responses"`       // interview answers
}

type InterviewQuestion struct {
    Question string   `json:"question"`
    Options  []string `json:"options"`  // optional multiple choice
    Required bool     `json:"required"`
}

type InterviewTemplate struct {
    Domain    string               `json:"domain"`
    Questions []InterviewQuestion  `json:"questions"`
}

func HandleAgentCreate(params map[string]interface{}) (string, error) {
    request := parseRequest(params)

    if request.Interactive {
        // Detect domain from initial request
        domain := detectDomain(request.InitialRequest)

        // Load interview template for domain
        template := loadInterviewTemplate(domain)

        // Return interview questions
        return formatInterviewQuestions(template), nil
    }

    // Non-interactive fallback (use initial request only)
    return "Creating agent based on request...", nil
}

func detectDomain(request string) string {
    // Simple keyword detection (can be improved)
    keywords := map[string][]string{
        "infrastructure": {"terraform", "kubernetes", "k8s", "docker", "ansible"},
        "backend":        {"api", "server", "backend", "database", "graphql"},
        "frontend":       {"react", "vue", "angular", "component", "ui"},
        "data":           {"etl", "pipeline", "analytics", "warehouse", "spark"},
        "ai-ml":          {"model", "training", "ml", "ai", "tensorflow"},
        "devops":         {"ci", "cd", "pipeline", "deploy", "jenkins"},
    }

    lower := strings.ToLower(request)
    for domain, keys := range keywords {
        for _, key := range keys {
            if strings.Contains(lower, key) {
                return domain
            }
        }
    }

    return "general"
}

func loadInterviewTemplate(domain string) InterviewTemplate {
    // Load from embedded templates or config
    templates := map[string]InterviewTemplate{
        "infrastructure": {
            Domain: "infrastructure",
            Questions: []InterviewQuestion{
                {
                    Question: "What will you primarily use this agent for?",
                    Options: []string{
                        "Multi-environment deployments",
                        "Infrastructure provisioning",
                        "State management",
                        "Module development",
                        "All of the above",
                    },
                    Required: true,
                },
                {
                    Question: "What's your philosophy on environment separation?",
                    Required: true,
                },
                {
                    Question: "What should this agent NEVER do?",
                    Required: true,
                },
                {
                    Question: "Walk me through your typical workflow",
                    Required: true,
                },
                {
                    Question: "Are you working solo or with a team?",
                    Options: []string{"Solo", "Small team (2-5)", "Large team (6+)"},
                    Required: false,
                },
            },
        },
        "backend": {
            Domain: "backend",
            Questions: []InterviewQuestion{
                {
                    Question: "What's your API design philosophy?",
                    Options: []string{"REST", "GraphQL", "gRPC", "Mix"},
                    Required: true,
                },
                {
                    Question: "What's your database approach?",
                    Options: []string{"ORM", "Query builder", "Raw SQL", "Mix"},
                    Required: true,
                },
                {
                    Question: "What's your error handling strategy?",
                    Required: true,
                },
                {
                    Question: "What's your testing philosophy?",
                    Required: true,
                },
            },
        },
        // More templates...
    }

    if template, exists := templates[domain]; exists {
        return template
    }

    return templates["general"]
}
```

2. **Update agent-architect.md**

Add interview framework to system prompt (see above).

3. **Create interview templates**

Embed templates for common domains (infrastructure, backend, frontend, data, ai-ml, devops).

4. **Test interview flow**

```
User: @agent-architect create a Terraform agent

Claude (via agent_create MCP):
I'll create a Terraform specialist agent. Let me understand your approach first.

[Loads infrastructure interview template]

What will you primarily use this agent for?
  1. Multi-environment deployments
  2. Infrastructure provisioning
  3. State management
  4. Module development
  5. All of the above

User: 1

Claude:
What's your philosophy on environment separation?

User: We use workspaces with S3 backend, separate AWS accounts for prod

Claude:
What should this agent NEVER do?

User: Never modify production without explicit approval

Claude:
Walk me through your typical workflow

User: PR-based with Atlantis. Plan on PR creation, apply on merge to main

Claude:
[Researches Terraform + Atlantis + multi-env patterns]

Based on our conversation, I'll create a Terraform agent optimized for:
- Multi-environment safety (workspace-aware)
- PR-driven workflows (Atlantis integration)
- Production protection (approval gates)
- S3 state management

[Creates sophisticated agent with this context]
```

---

## Interview Question Database

### Core Questions (All Domains)

1. **Purpose:** "What will you primarily use this agent for?"
2. **Philosophy:** "What's your philosophy on [domain key concept]?"
3. **Constraints:** "What should this agent NEVER do?"
4. **Workflow:** "Walk me through your typical workflow"
5. **Team:** "Are you working solo or with a team?"
6. **Experience:** "What's your experience level with [domain]?"

### Domain-Specific Questions

**Infrastructure (Terraform, K8s, Ansible):**
- Environment strategy?
- State/config management approach?
- Approval process for changes?
- Disaster recovery philosophy?
- Multi-region/zone considerations?

**Backend (APIs, Databases):**
- API design philosophy? (REST, GraphQL, gRPC)
- Database approach? (ORM, query builder, raw SQL)
- Error handling strategy?
- Testing philosophy?
- Auth/authz patterns?
- Caching strategy?

**Frontend (React, Vue, Angular):**
- State management preference?
- Styling approach?
- Component architecture?
- Testing strategy?
- Accessibility priority?
- Performance optimization stance?

**Data (ETL, Analytics):**
- Data quality philosophy?
- Pipeline architecture? (batch, streaming, hybrid)
- Schema evolution strategy?
- Testing approach for data?
- Monitoring/observability stance?

**AI/ML:**
- Model development philosophy?
- Training vs inference priorities?
- Experimentation tracking approach?
- Production deployment strategy?
- Monitoring/drift detection?

**DevOps (CI/CD):**
- Deployment frequency goal?
- Rollback strategy?
- Testing gates?
- Environment promotion strategy?
- Secrets management approach?

**Security:**
- Threat modeling approach?
- Defense-in-depth philosophy?
- Incident response priorities?
- Compliance requirements?
- Risk tolerance?

---

## User Experience Flow

### 1. Simple Request (Enhanced)

**Before (weak):**
```
User: @agent-architect create a Terraform agent
Claude: [Creates agent immediately]
```

**After (guided):**
```
User: @agent-architect create a Terraform agent

Claude: I'll create a Terraform specialist agent. Let me ask a few questions to
make it perfect for your workflow.

[5-7 targeted questions]

Claude: Great! Based on our conversation:
- Multi-env with workspaces
- PR-driven via Atlantis
- Production protection
- S3 state management

Creating your agent now...

[Creates sophisticated, contextualized agent]

✓ Created: vc-agents/my-agents/infrastructure/terraform-specialist.md
```

### 2. Detailed Request (Streamlined)

**If user provides context upfront:**
```
User: @agent-architect create a Terraform agent for multi-environment deployments
using Atlantis PR workflows, workspaces, and S3 backend. Never touch production
without approval.

Claude: Perfect! I have everything I need. Let me verify:
- Multi-env: workspaces + S3
- Atlantis PR workflow
- Production protection

Anything else to add?

User: Nope

Claude: Creating your agent...
[Skips interview, uses provided context]
```

### 3. Skip Interview (Advanced Users)

**If user wants to skip:**
```
User: @agent-architect create a Terraform agent (skip interview)

Claude: Creating Terraform agent with default patterns...

[Creates good-but-generic agent]

Note: For a more tailored agent, try without "(skip interview)"
```

---

## Configuration

**`~/.cami/config.yaml`:**

```yaml
agent_creation:
  interactive_default: true        # Enable interview by default
  interview_timeout: 300          # 5 minutes max for interview
  skip_interview_keywords:        # Keywords to skip interview
    - "skip interview"
    - "no questions"
    - "quick agent"
```

---

## Benefits

### For Users

1. ✅ **Better agents** - Contextual, not generic
2. ✅ **Learn domain** - Interview teaches best practices
3. ✅ **Clear constraints** - Agent knows what NOT to do
4. ✅ **Workflow integration** - Matches actual work patterns

### For Agent Quality

1. ✅ **Philosophy-driven** - Clear stance on key decisions
2. ✅ **Opinionated** - Guides toward best practices
3. ✅ **Context-aware** - Knows user's environment
4. ✅ **Constraint-explicit** - Clear boundaries

### For Ecosystem

1. ✅ **Higher quality agents** - Better contributions
2. ✅ **Consistent methodology** - Interview framework reusable
3. ✅ **Knowledge capture** - Interview responses inform future agents

---

## Implementation Priority

**Phase 2A (within Phase 2): Agent Creation Interview**

**Week 1:**
- Day 1-2: `agent_create` MCP tool
- Day 3-4: Interview templates (5 domains)
- Day 5: Update agent-architect.md

**Week 2:**
- Day 1-2: Integration testing
- Day 3-4: Polish UX
- Day 5: Documentation

**Deliverables:**
- ✅ `agent_create` MCP tool
- ✅ 5+ domain interview templates
- ✅ Agent-architect enhanced
- ✅ Tested interview flows

---

## Success Metrics

**Qualitative:**
- Agents created with interview > agents without
- User feedback: "Agent understands my workflow"
- Contributions have clear philosophy sections

**Quantitative:**
- 80%+ of agent creations use interview
- Average 5-7 questions per interview
- Interview completion rate >70%
- Agent quality ratings higher with interview

---

## Example Interview Templates

### Template: Infrastructure Agent

```yaml
domain: infrastructure
questions:
  - question: "What will you primarily use this agent for?"
    type: multiple_choice
    options:
      - "Multi-environment deployments"
      - "Infrastructure provisioning"
      - "State management"
      - "Module development"
      - "All of the above"
    required: true

  - question: "What's your philosophy on environment separation?"
    type: open
    prompt: "e.g., workspaces, separate accounts, mono-repo, etc."
    required: true

  - question: "What should this agent NEVER do?"
    type: open
    prompt: "e.g., modify production directly, delete resources, etc."
    required: true

  - question: "Walk me through your typical deployment workflow"
    type: open
    prompt: "e.g., PR-based, manual approval, automated on merge, etc."
    required: true

  - question: "Are you working solo or with a team?"
    type: multiple_choice
    options:
      - "Solo"
      - "Small team (2-5)"
      - "Large team (6+)"
    required: false
```

### Template: Backend API Agent

```yaml
domain: backend
questions:
  - question: "What's your API design philosophy?"
    type: multiple_choice
    options:
      - "REST (resource-oriented)"
      - "GraphQL (query-based)"
      - "gRPC (performance-critical)"
      - "Mix of approaches"
    required: true

  - question: "What's your database approach?"
    type: multiple_choice
    options:
      - "ORM (e.g., Prisma, TypeORM)"
      - "Query builder (e.g., Knex)"
      - "Raw SQL"
      - "Mix depending on use case"
    required: true

  - question: "What's your error handling strategy?"
    type: open
    prompt: "e.g., exceptions, result types, error codes, etc."
    required: true

  - question: "What's your testing philosophy?"
    type: open
    prompt: "e.g., unit + integration, E2E focus, TDD, etc."
    required: true

  - question: "Any specific frameworks or constraints?"
    type: open
    prompt: "e.g., Express, Fastify, company patterns, etc."
    required: false
```

---

## Summary

**The interview system ensures agent-architect leverages its research capabilities to create sophisticated, philosophy-driven agents instead of generic templates.**

**Key principle:**
> Don't create agents. **Craft** agents through discovery.
