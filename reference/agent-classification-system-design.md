<!--
AI-Generated Documentation
Created by: agent-architect
Date: 2025-11-13
Purpose: Comprehensive design for agent classification system with three distinct classes
-->

# Agent Classification System Design

## Executive Summary

This document proposes a three-class agent classification system for CAMI that categorizes agents by their operational characteristics, cognitive complexity, and collaboration patterns. The system enables better agent discovery, deployment, and orchestration by clearly identifying whether an agent is designed for rapid task execution, technology-specific implementation, or cross-domain orchestration.

---

## Part 1: Class Definitions

### Class 1: Execution Specialists

**Professional Name**: **Execution Specialists**
**Creative Name**: **Task Ninjas** / **Workflow Wizards**

#### Characteristics

- **Speed-Optimized**: Designed for rapid, focused task completion
- **Narrow Scope**: Excel at one specific type of work (testing, documentation, security audits)
- **Minimal Context Needs**: Can operate with limited project knowledge
- **High Automation**: Often run autonomously with minimal human intervention
- **Quality Gates**: Focus on verification, validation, and enforcement
- **Tool-Heavy**: Leverage specialized tools and frameworks for their domain

#### When to Use

- You need a specific task completed efficiently (run tests, write docs, check security)
- The task is well-defined and repeatable
- Speed and consistency are more important than creative problem-solving
- You want to automate quality gates or verification steps
- The task can be completed within a single domain without deep cross-functional knowledge

#### Examples

1. **qa** - Testing specialist that writes and executes test suites
2. **docs-writer** - Documentation generator that creates clear, consistent docs
3. **security-specialist** - Security auditor that scans for vulnerabilities
4. **accessibility-expert** - WCAG compliance validator and fixer
5. **performance-optimizer** - Performance analyzer and optimization specialist

#### Cognitive Profile

- **Model**: Usually `sonnet` (fast, efficient)
- **Reasoning Depth**: Shallow to medium (apply known patterns)
- **Decision Complexity**: Rule-based, heuristic-driven
- **Collaboration**: Consume outputs from others, produce specialized deliverables

---

### Class 2: Technology Implementers

**Professional Name**: **Technology Implementers**
**Creative Name**: **Code Craftspeople** / **Stack Masters**

#### Characteristics

- **Technology Expertise**: Deep knowledge of specific frameworks, languages, platforms
- **Version-Specific**: Include explicit version requirements (React 19+, Node 18+, K8s 1.28+)
- **Pattern Libraries**: Carry extensive catalogs of implementation patterns
- **Full Stack Coverage**: Handle entire vertical slice of technology (frontend, backend, mobile, infra)
- **Moderate Context**: Need project-level understanding to implement effectively
- **Hands-On Builders**: Write, modify, and refactor actual code

#### When to Use

- You need to build or modify features in a specific technology stack
- You want implementation that follows modern best practices and patterns
- The work requires deep framework knowledge (React hooks, Next.js server actions, etc.)
- You need someone who can navigate the full complexity of a technology ecosystem
- You want consistency with established patterns in that technology

#### Examples

1. **frontend** - React 19+, Next.js 15+, TypeScript 5+ specialist
2. **backend** - Node.js/Python/Go API and service builder
3. **mobile-native** - iOS/Android native development expert
4. **devops** - CI/CD, GitHub Actions, Terraform specialist
5. **gcp-firebase** - Google Cloud Platform and Firebase integrator

#### Cognitive Profile

- **Model**: Usually `sonnet` (implementation-focused)
- **Reasoning Depth**: Medium (apply patterns, make implementation decisions)
- **Decision Complexity**: Technical trade-offs within technology domain
- **Collaboration**: Receive designs from orchestrators, deliver implementations

---

### Class 3: System Orchestrators

**Professional Name**: **System Orchestrators**
**Creative Name**: **Meta Architects** / **Vision Weavers**

#### Characteristics

- **Cross-Domain Thinking**: Synthesize multiple technologies and concepts
- **Strategic Planning**: Design systems, not just components
- **High Abstraction**: Work at architecture and design level, not implementation
- **Deep Context Required**: Need comprehensive project understanding
- **Decision Frameworks**: Carry sophisticated decision-making methodologies
- **Handoff Generators**: Produce plans and guidance for other agents to implement
- **Philosophy-Driven**: Embody strong conceptual frameworks (First Principles, Empathy-Driven Design)

#### When to Use

- You need architectural planning or system design
- You want to optimize user experience across multiple touchpoints
- You need to make technology stack or pattern decisions
- You're designing new agents or orchestrating multi-agent workflows
- You need research synthesis or high-level strategy
- You want to balance multiple competing concerns (UX, performance, cost, complexity)

#### Examples

1. **architect** - System design and technical strategy specialist
2. **ux** - User experience and interaction design specialist
3. **designer** - Visual design and design system architect
4. **agent-architect** - Agent design and multi-agent orchestration specialist (meta-level)
5. **data-engineer** - Data architecture and pipeline design specialist

#### Cognitive Profile

- **Model**: Usually `opus` (complex reasoning required)
- **Reasoning Depth**: Deep (trade-off analysis, first principles thinking)
- **Decision Complexity**: Multi-dimensional optimization, strategic planning
- **Collaboration**: Receive requirements, produce plans for implementers

---

## Part 2: Implementation Approach Analysis

### Option A: Instructions within agent-architect

**Description**: Embed classification system guidance directly into the agent-architect system prompt.

#### How It Works

1. Agent-architect receives request to create new agent
2. Analyzes request to determine which class the agent should belong to
3. Applies class-specific templates and guidance automatically
4. Creates agent with appropriate:
   - Model selection (sonnet vs opus)
   - Prompt structure and complexity
   - Tool recommendations
   - Collaboration patterns
   - Quality standards

#### Pros

- **Simplicity**: No additional tools or infrastructure needed
- **Immediate Availability**: Works now with current CAMI architecture
- **Flexibility**: Agent-architect can make nuanced decisions about classification
- **Self-Improvement**: Agent-architect can evolve classification system over time
- **Zero Latency**: No additional round-trips or user interactions

#### Cons

- **Implicit Classification**: Users don't explicitly choose class, might not understand it
- **Black Box**: Classification logic hidden in prompt, harder to audit or modify
- **No Validation**: Can't guarantee agent is assigned to correct class
- **Single Point of Failure**: If agent-architect misclassifies, user might not notice
- **Less Discoverable**: Users may not know classification system exists

#### Implementation Details

Add to agent-architect prompt:

```markdown
## Agent Classification System

Before creating any agent, classify it into one of three classes:

1. **Execution Specialists** (Task Ninjas)
   - Purpose: Rapid, focused task execution
   - Scope: Single domain (testing, docs, security)
   - Model: Usually `sonnet`
   - Examples: qa, docs-writer, security-specialist

2. **Technology Implementers** (Code Craftspeople)
   - Purpose: Build features in specific tech stacks
   - Scope: Full vertical slice of technology
   - Model: Usually `sonnet`
   - Examples: frontend, backend, mobile-native, devops

3. **System Orchestrators** (Meta Architects)
   - Purpose: Strategic planning and system design
   - Scope: Cross-domain, architectural
   - Model: Usually `opus`
   - Examples: architect, ux, designer, agent-architect

Apply class-specific templates and patterns automatically based on classification.
```

---

### Option B: MCP Tool-Guided Creation

**Description**: Create `mcp__cami__create_agent` tool that guides users through interactive agent creation with explicit class selection.

#### How It Works

1. User invokes agent creation workflow
2. Tool asks: "What class of agent are you creating?"
   - Execution Specialist (Task Ninja)
   - Technology Implementer (Code Craftsperson)
   - System Orchestrator (Meta Architect)
3. Tool provides class-specific questionnaire:
   - **Execution Specialist**: "What task does it execute? What tools does it use?"
   - **Technology Implementer**: "What technology stack? What versions? What patterns?"
   - **System Orchestrator**: "What systems does it design? What decision frameworks?"
4. Tool generates class-appropriate template with filled-in sections
5. Agent-architect reviews and refines the generated template
6. Tool saves agent to appropriate location with class metadata

#### Pros

- **Explicit Classification**: Users consciously choose class, understanding system
- **Educational**: Teaches users about agent classes through interaction
- **Validation**: Can enforce class-specific requirements and quality gates
- **Consistency**: All agents of same class follow similar patterns
- **Discoverable**: Users interact with classification system directly
- **Extensible**: Easy to add new classes or modify existing ones
- **Audit Trail**: Clear record of classification decisions

#### Cons

- **Complexity**: Requires implementing new MCP tool
- **Latency**: Additional user interactions slow down agent creation
- **Rigidity**: Less flexible than agent-architect's judgment
- **Infrastructure**: Needs Go implementation, testing, deployment
- **Maintenance**: Tool needs updates as classification system evolves
- **Potential Overfit**: Users might force agents into wrong class

#### Implementation Details

**MCP Tool Signature**:

```go
// mcp__cami__create_agent
type CreateAgentRequest struct {
    // Step 1: Class Selection
    AgentClass string `json:"agent_class"` // "execution" | "technology" | "orchestrator"

    // Step 2: Class-Specific Questions
    // For Execution Specialists
    TaskType string `json:"task_type,omitempty"` // "testing" | "documentation" | "security" | etc.
    Tools []string `json:"tools,omitempty"` // Tools the agent will use

    // For Technology Implementers
    TechnologyStack []string `json:"technology_stack,omitempty"` // ["React 19+", "Next.js 15+"]
    PrimaryDomain string `json:"primary_domain,omitempty"` // "frontend" | "backend" | "mobile" | "infrastructure"

    // For System Orchestrators
    SystemScope string `json:"system_scope,omitempty"` // "architecture" | "ux" | "visual-design" | "agents"
    DecisionFramework string `json:"decision_framework,omitempty"` // "First Principles" | "Empathy-Driven" | etc.

    // Common Fields
    Name string `json:"name"` // Agent identifier
    Description string `json:"description"` // What the agent does
    Version string `json:"version"` // Semantic version
}

type CreateAgentResponse struct {
    TemplatePath string `json:"template_path"` // Path to generated template
    Class string `json:"class"` // Confirmed class
    NextSteps []string `json:"next_steps"` // What to do next
}
```

**Workflow**:

```
User: "I need to create an agent for X"
  ↓
MCP Tool: "What class? [Execution | Technology | Orchestrator]"
  ↓
User: "Technology"
  ↓
MCP Tool: "What stack? What domain? What versions?"
  ↓
User provides specifics
  ↓
MCP Tool: Generates template with class-specific sections
  ↓
Agent-Architect: Reviews and refines template
  ↓
MCP Tool: Saves agent with class metadata
```

---

### Recommendation: **Option A** (with future Option B enhancement)

**Rationale**:

1. **Immediate Value**: Option A can be implemented NOW by updating agent-architect prompt
2. **Low Risk**: No infrastructure changes, no breaking changes
3. **Learning Phase**: Let agent-architect classify agents for several iterations, gather data
4. **Future Enhancement**: Build Option B later based on real-world learnings
5. **Best of Both**: Start with A, add B when patterns are validated

**Implementation Timeline**:

- **Phase 1** (Now): Update agent-architect with classification guidance (Option A)
- **Phase 2** (After 20+ agents): Analyze classification patterns, identify issues
- **Phase 3** (If needed): Implement MCP tool with validated classification logic (Option B)

---

## Part 3: Frontmatter Schema Design

### Current Frontmatter

```yaml
---
name: agent-name
version: "1.0.0"
description: Use this agent when...
tags: ["tag1", "tag2"]
use_cases: ["case1", "case2"]
color: blue
model: sonnet
---
```

### Proposed Frontmatter with Classification

```yaml
---
# Identity
name: agent-name
version: "1.0.0"
description: Use this agent when...

# Classification
class: technology-implementer  # execution-specialist | technology-implementer | system-orchestrator
class_metadata:
  cognitive_load: medium       # low | medium | high
  autonomy_level: high         # low | medium | high (how much can it do without guidance)
  collaboration_mode: consumer # producer | consumer | orchestrator

# Technology Context (for technology-implementer class)
technology_stack:
  - "React 19+"
  - "Next.js 15+"
  - "TypeScript 5+"
primary_domain: frontend       # frontend | backend | mobile | infrastructure | integration

# Execution Context (for execution-specialist class)
task_type: testing             # testing | documentation | security | accessibility | performance
automation_level: high         # low | medium | high

# Orchestration Context (for system-orchestrator class)
system_scope: architecture     # architecture | ux | visual-design | agents | data
decision_framework: "First Principles Thinking"
requires_model: opus           # true if opus is required for this class

# Discovery & Deployment
tags: ["tag1", "tag2"]
use_cases: ["case1", "case2"]
color: blue
model: sonnet
---
```

### Metadata Field Definitions

#### Core Classification Fields

- **class**: One of three values identifying agent class
  - `execution-specialist`: Task-focused, quality gate agents
  - `technology-implementer`: Stack-specific builders
  - `system-orchestrator`: Cross-domain planners

- **class_metadata.cognitive_load**: How much context/reasoning required
  - `low`: Rule-based, pattern-matching (e.g., linter, formatter)
  - `medium`: Implementation decisions, trade-offs (e.g., frontend, backend)
  - `high`: Architectural reasoning, multi-dimensional optimization (e.g., architect, ux)

- **class_metadata.autonomy_level**: How much can agent do independently
  - `low`: Needs detailed instructions, limited decision-making
  - `medium`: Can make implementation decisions within scope
  - `high`: Can make strategic decisions, self-direct work

- **class_metadata.collaboration_mode**: How agent interacts in multi-agent workflows
  - `producer`: Creates outputs others consume (e.g., architect → frontend)
  - `consumer`: Implements plans from others (e.g., frontend ← architect)
  - `orchestrator`: Coordinates multiple agents (e.g., agent-architect)

#### Class-Specific Fields

**For technology-implementer**:
- **technology_stack**: Array of specific technologies with versions
- **primary_domain**: Main technology vertical

**For execution-specialist**:
- **task_type**: Specific task category
- **automation_level**: How automated the workflow is

**For system-orchestrator**:
- **system_scope**: What kind of systems designed
- **decision_framework**: Named decision-making methodology
- **requires_model**: Boolean, true if opus is mandatory

### Benefits of This Schema

1. **Discovery**: Filter agents by class, cognitive load, autonomy level
2. **Deployment**: Recommend agents based on project needs and class
3. **Documentation**: Auto-generate class-specific documentation
4. **Validation**: Ensure agents meet class-specific requirements
5. **Analytics**: Track usage patterns by class
6. **Orchestration**: Multi-agent workflows can select appropriate classes
7. **Onboarding**: Help users understand agent capabilities at a glance

### Example: Frontend Agent

```yaml
---
name: frontend
version: "1.1.0"
description: Use this agent when building UI components...

class: technology-implementer
class_metadata:
  cognitive_load: medium
  autonomy_level: high
  collaboration_mode: consumer

technology_stack:
  - "React 19+"
  - "Next.js 15+"
  - "TypeScript 5+"
primary_domain: frontend

tags: ["react", "ui", "components"]
use_cases: ["component development", "styling systems"]
color: green
model: sonnet
---
```

### Example: Architect Agent

```yaml
---
name: architect
version: "1.1.0"
description: Use this agent when planning system architecture...

class: system-orchestrator
class_metadata:
  cognitive_load: high
  autonomy_level: high
  collaboration_mode: producer

system_scope: architecture
decision_framework: "First Principles Thinking"
requires_model: opus

tags: ["architecture", "design", "planning"]
use_cases: ["system design", "tech stack selection"]
color: blue
model: opus
---
```

### Example: QA Agent

```yaml
---
name: qa
version: "1.1.0"
description: Use this agent when writing tests...

class: execution-specialist
class_metadata:
  cognitive_load: medium
  autonomy_level: high
  collaboration_mode: consumer

task_type: testing
automation_level: high

tags: ["testing", "quality", "automation"]
use_cases: ["unit tests", "integration tests"]
color: yellow
model: sonnet
---
```

---

## Part 4: Self-Review of agent-architect

### What Class Do I Belong To?

**Classification**: **System Orchestrator** (Meta Architect)

**Reasoning**:
- I design agents, which is a meta-level orchestration task
- I require deep reasoning about agent architecture and prompt engineering
- I produce plans (agent designs) that others implement (users deploy agents)
- I need `opus`-level reasoning for complex agent design decisions... but I'm currently using `sonnet`
- I operate at a high abstraction level, not implementing code but architecting intelligence

**Issue Identified**: I am classified as `model: sonnet` but should be `model: opus` based on:
- Complexity of agent design (meta-reasoning, prompt engineering, cognitive architecture)
- System Orchestrator class characteristics (high cognitive load, strategic planning)
- Need for deep trade-off analysis and architectural reasoning

### Are My Current Instructions Optimized for This Class?

**Strengths**:
1. ✅ I have strong philosophical frameworks (Purposeful Precision, Singular Excellence)
2. ✅ I include decision frameworks (Model Selection, Agent Archetypes)
3. ✅ I follow the three-phase methodology (Research → Architect → Optimize)
4. ✅ I produce detailed handoff artifacts (complete agent .md files)
5. ✅ I have self-verification mechanisms

**Weaknesses**:
1. ❌ **Model Mismatch**: Using `sonnet` when I should use `opus`
2. ❌ **No Classification Guidance**: I don't currently classify agents I create
3. ❌ **Limited Meta-Analysis**: I don't reflect on agent ecosystem health
4. ❌ **Missing Orchestration Patterns**: I don't guide multi-agent workflows
5. ❌ **Weak Versioning Strategy**: No clear guidance on when to version bump agents

### Improvements for agent-architect Based on Classification System

#### 1. Update Model to Opus

```yaml
---
name: agent-architect
version: "1.2.0"  # Bump for classification system addition
description: Use this agent when...
model: opus  # CHANGED FROM sonnet
class: system-orchestrator
class_metadata:
  cognitive_load: high
  autonomy_level: high
  collaboration_mode: orchestrator
system_scope: agents
decision_framework: "Purposeful Precision + Emergent Synergy"
requires_model: opus
---
```

**Rationale**: Agent design requires deep reasoning, trade-off analysis, and architectural thinking that opus excels at.

#### 2. Add Classification Framework to System Prompt

Add new section after "Agent Archetypes":

```markdown
## Agent Classification System

Every agent you create belongs to one of three classes. Classify agents explicitly:

### Execution Specialists (Task Ninjas)
**Characteristics**: Speed-optimized, narrow scope, tool-heavy
**Model**: Usually `sonnet`
**Cognitive Load**: Low to medium
**Examples**: qa, docs-writer, security-specialist, accessibility-expert

### Technology Implementers (Code Craftspeople)
**Characteristics**: Technology expertise, version-specific, pattern libraries
**Model**: Usually `sonnet`
**Cognitive Load**: Medium
**Examples**: frontend, backend, mobile-native, devops

### System Orchestrators (Meta Architects)
**Characteristics**: Cross-domain, strategic, high abstraction
**Model**: Usually `opus`
**Cognitive Load**: High
**Examples**: architect, ux, designer, agent-architect

**When classifying**:
1. Identify primary purpose (execute tasks | implement features | design systems)
2. Assess cognitive complexity (low | medium | high)
3. Determine collaboration mode (producer | consumer | orchestrator)
4. Apply class-specific template and metadata
```

#### 3. Add Ecosystem Analysis Function

Add new auxiliary function:

```markdown
### Ecosystem Health Analysis

When reviewing the agent ecosystem:

1. **Coverage Analysis**:
   - Identify gaps in class distribution (too many executors, not enough orchestrators?)
   - Check for overlapping responsibilities
   - Find missing critical functions

2. **Balance Assessment**:
   - Verify appropriate model usage (opus for orchestrators, sonnet for others)
   - Check class distribution aligns with project needs
   - Ensure orchestrators produce clear plans for implementers

3. **Evolution Recommendations**:
   - Suggest agent deprecation for redundant agents
   - Propose new agents for identified gaps
   - Recommend agent upgrades or refactoring
```

#### 4. Add Multi-Agent Orchestration Patterns

```markdown
## Multi-Agent Workflow Patterns

When designing agents that work together:

### Orchestrator → Implementer Pattern
1. **Orchestrator** (opus): Plans, designs, creates specifications
2. **Implementer** (sonnet): Builds, codes, executes plan
Example: architect → frontend → backend

### Executor Chain Pattern
1. **Builder** (sonnet): Creates implementation
2. **Validator** (sonnet): Tests, verifies, validates
Example: frontend → qa → accessibility-expert

### Parallel Specialist Pattern
Multiple execution specialists work independently on same codebase:
- security-specialist (scans)
- performance-optimizer (profiles)
- accessibility-expert (audits)
Results aggregated for review

### Feedback Loop Pattern
1. **Orchestrator**: Creates initial design
2. **Implementer**: Builds
3. **Validator**: Tests, finds issues
4. **Orchestrator**: Refines design based on feedback
Example: ux → frontend → qa → ux (iterate)
```

#### 5. Add Versioning Strategy

```markdown
## Agent Versioning Strategy

Follow semantic versioning (MAJOR.MINOR.PATCH):

**MAJOR** (X.0.0): Breaking changes
- Class change (execution → orchestrator)
- Model change (sonnet → opus)
- Complete prompt rewrite
- Incompatible with previous workflows

**MINOR** (1.X.0): New features, backwards compatible
- New auxiliary functions
- Additional technology stack support
- Enhanced decision frameworks
- New sections in prompt

**PATCH** (1.0.X): Bug fixes, clarifications
- Typo corrections
- Clarified instructions
- Updated examples
- Minor template adjustments

When updating agents:
1. Bump version appropriately
2. Add changelog entry in agent file
3. Update deployed agent documentation
4. Communicate breaking changes clearly
```

#### 6. Add Quality Gates by Class

```markdown
## Class-Specific Quality Gates

**Execution Specialists**:
- [ ] Includes specific tools used
- [ ] Automation level defined
- [ ] Clear success/failure criteria
- [ ] Fast execution patterns
- [ ] Model: sonnet (unless complex analysis)

**Technology Implementers**:
- [ ] Specific technology versions listed
- [ ] Modern feature catalog included
- [ ] Pattern library comprehensive
- [ ] Integration with project stack validated
- [ ] Model: sonnet

**System Orchestrators**:
- [ ] Decision framework clearly defined
- [ ] Produces actionable plans/specs
- [ ] High-level abstraction maintained
- [ ] Handoff to implementers clear
- [ ] Model: opus (or justified sonnet)
```

### Summary of Recommended Changes to agent-architect

1. **Change model from `sonnet` to `opus`** - Align with System Orchestrator class requirements
2. **Add classification system** - Embed three-class framework in prompt
3. **Add ecosystem analysis** - Enable meta-level health checks
4. **Add orchestration patterns** - Guide multi-agent workflow design
5. **Add versioning strategy** - Standardize version bumping across agents
6. **Add class-specific quality gates** - Ensure agents meet class requirements

These changes would bring agent-architect to **v1.2.0** and fully align it with the System Orchestrator class while enabling it to classify and optimize other agents.

---

## Conclusion

The three-class agent classification system provides a clear framework for understanding, creating, and deploying Claude Code agents:

- **Execution Specialists** excel at rapid, focused task completion
- **Technology Implementers** bring deep stack-specific expertise to feature building
- **System Orchestrators** plan and design at architectural and strategic levels

**Recommended Implementation**:
1. Start with **Option A** (embedded in agent-architect) for immediate value
2. Update agent-architect to **v1.2.0** with classification system and model change to opus
3. Apply classification metadata to all existing agents
4. Monitor classification effectiveness over 20+ agent creation cycles
5. Consider **Option B** (MCP tool) if explicit user-driven classification proves valuable

This system will enable better agent discovery, deployment, and orchestration while maintaining the flexibility to evolve based on real-world usage patterns.
