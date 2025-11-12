<!--
AI-Generated Documentation
Created by: agent-architect
Date: 2025-11-03
Purpose: Architecture and implementation guide for making Claude Code agents MCP-aware
-->

# MCP-Awareness Architecture for Claude Code Agents

## Executive Summary

This document defines a comprehensive approach for instructing Claude Code agents to intelligently discover, prioritize, and use Model Context Protocol (MCP) servers based on their domain expertise. The architecture enables agents to leverage domain-specific MCPs without hardcoding server names, creating a flexible and scalable multi-agent + multi-MCP ecosystem.

## 1. Core Architecture

### 1.1 MCP Type Taxonomy

MCPs are categorized by their primary function, enabling agents to discover and prioritize relevant tools:

#### **Design & UI Category**
- **Purpose**: Design systems, component libraries, UI frameworks, accessibility tools
- **Examples**:
  - `lando-labs-design-system` - Component library and design tokens
  - `figma-mcp` - Figma design file access
  - `storybook-mcp` - Component documentation
  - `axe-core-mcp` - Accessibility validation
- **Tool Patterns**: `design_system__*`, `component__*`, `token__*`, `a11y__*`
- **Resource Patterns**: `design://`, `component://`, `token://`

#### **Development Category**
- **Purpose**: File operations, code analysis, browser automation, testing
- **Examples**:
  - `filesystem` - File read/write operations
  - `playwright` - Browser automation and testing
  - `git-mcp` - Git operations
  - `code-analysis-mcp` - Static analysis tools
- **Tool Patterns**: `mcp__filesystem__*`, `mcp__playwright__*`, `mcp__git__*`
- **Resource Patterns**: `file://`, `git://`, `browser://`

#### **Infrastructure Category**
- **Purpose**: Cloud providers, Kubernetes, databases, monitoring, deployment
- **Examples**:
  - `firebase` - Firebase backend services
  - `kubernetes-mcp` - K8s cluster management
  - `aws-mcp` - AWS service integration
  - `terraform-mcp` - Infrastructure as code
  - `datadog-mcp` - Monitoring and observability
- **Tool Patterns**: `mcp__firebase__*`, `mcp__k8s__*`, `mcp__aws__*`
- **Resource Patterns**: `firebase://`, `k8s://`, `aws://`

#### **Knowledge Category**
- **Purpose**: Documentation, library references, API specifications
- **Examples**:
  - `context7` - Library documentation provider
  - `api-docs-mcp` - API specification access
  - `wiki-mcp` - Company wiki/knowledge base
  - `confluence-mcp` - Confluence documentation
- **Tool Patterns**: `mcp__context7__*`, `mcp__docs__*`
- **Resource Patterns**: `docs://`, `api://`, `wiki://`

#### **Data & Analytics Category**
- **Purpose**: Databases, data warehouses, analytics, ETL
- **Examples**:
  - `postgres-mcp` - PostgreSQL database access
  - `mongodb-mcp` - MongoDB operations
  - `bigquery-mcp` - BigQuery analytics
  - `elasticsearch-mcp` - Search and analytics
- **Tool Patterns**: `mcp__postgres__*`, `mcp__mongodb__*`, `mcp__bigquery__*`
- **Resource Patterns**: `db://`, `postgres://`, `mongodb://`

#### **Testing & QA Category**
- **Purpose**: Test frameworks, coverage tools, performance testing
- **Examples**:
  - `jest-mcp` - Jest test runner integration
  - `cypress-mcp` - E2E testing framework
  - `lighthouse-mcp` - Performance auditing
  - `coverage-mcp` - Code coverage analysis
- **Tool Patterns**: `mcp__jest__*`, `mcp__cypress__*`, `mcp__lighthouse__*`
- **Resource Patterns**: `test://`, `coverage://`, `perf://`

#### **Communication & Workflow Category**
- **Purpose**: Slack, email, calendar, project management
- **Examples**:
  - `slack-mcp` - Slack messaging
  - `linear-mcp` - Project management
  - `calendar-mcp` - Calendar integration
  - `email-mcp` - Email operations
- **Tool Patterns**: `mcp__slack__*`, `mcp__linear__*`, `mcp__calendar__*`
- **Resource Patterns**: `slack://`, `linear://`, `calendar://`

#### **AI & ML Category**
- **Purpose**: Model serving, training, MLOps, vector databases
- **Examples**:
  - `openai-mcp` - OpenAI API integration
  - `huggingface-mcp` - HuggingFace models
  - `pinecone-mcp` - Vector database
  - `mlflow-mcp` - ML experiment tracking
- **Tool Patterns**: `mcp__openai__*`, `mcp__hf__*`, `mcp__pinecone__*`
- **Resource Patterns**: `model://`, `vector://`, `ml://`

#### **Utilities Category**
- **Purpose**: Time, calculations, formatting, validation
- **Examples**:
  - `time` - Time zone conversions
  - `sequential-thinking` - Complex reasoning
  - `calculator-mcp` - Mathematical operations
  - `validator-mcp` - Data validation
- **Tool Patterns**: `mcp__time__*`, `mcp__calc__*`, `mcp__validate__*`
- **Resource Patterns**: N/A (typically tool-focused)

### 1.2 Agent-MCP Affinity Matrix

This matrix defines which agent categories should prioritize which MCP categories:

| Agent Category | Primary MCPs | Secondary MCPs | Tertiary MCPs |
|---|---|---|---|
| **Design Agents** (designer, ux, accessibility-expert) | Design & UI | Knowledge, Testing & QA | Development |
| **Frontend Agents** (frontend, react-specialist, design-system-specialist) | Design & UI, Development | Testing & QA, Knowledge | Infrastructure |
| **Backend Agents** (backend, api-integrator, data-engineer) | Data & Analytics, Development | Infrastructure, Knowledge | Testing & QA |
| **Infrastructure Agents** (devops, deploy, gcp-firebase) | Infrastructure, Development | Testing & QA, Communication | Data & Analytics |
| **Testing Agents** (qa, performance-optimizer) | Testing & QA, Development | Infrastructure, Design & UI | Knowledge |
| **Research Agents** (research-synthesizer, use-case-specialist) | Knowledge, Utilities | Communication, Development | All others |
| **AI/ML Agents** (ai-ml-specialist) | AI & ML, Data & Analytics | Infrastructure, Development | Knowledge |
| **Mobile Agents** (mobile-native) | Development, Design & UI | Testing & QA, Infrastructure | Knowledge |
| **Security Agents** (security-specialist) | Infrastructure, Development | Testing & QA, Data & Analytics | Communication |

### 1.3 MCP Discovery Pattern

Agents should discover and prioritize MCPs using this four-step approach:

#### Step 1: List Available MCPs
```
Use ListMcpResourcesTool to discover all configured MCPs
Parse tool names to identify MCP servers (pattern: mcp__[server]__[tool])
Build internal registry of available MCPs by category
```

#### Step 2: Categorize MCPs
```
For each discovered MCP:
  - Extract server name from tool pattern
  - Match against known category patterns
  - Tag with category (Design & UI, Infrastructure, etc.)
  - Note tool capabilities and resource types
```

#### Step 3: Prioritize by Affinity
```
Based on agent's role (from Agent-MCP Affinity Matrix):
  - Assign priority scores: Primary (3), Secondary (2), Tertiary (1)
  - Sort MCPs by: Affinity Score > Tool Count > Alphabetical
  - Create prioritized tool list for the current task
```

#### Step 4: Select and Use
```
For the current task:
  1. Identify task requirements (e.g., "read design tokens")
  2. Match requirements to MCP categories
  3. Select highest-priority MCP with matching capability
  4. Invoke MCP tool with appropriate parameters
  5. Fallback to lower-priority MCP if first fails
  6. Use Claude Code built-in tools if no MCP matches
```

## 2. Agent System Prompt Template

### 2.1 MCP Awareness Section

Add this section to agent system prompts after "Core Philosophy" and before "Three-Phase Methodology":

```markdown
## MCP Awareness: Specialized Tool Priority

You have access to Model Context Protocol (MCP) servers that extend your capabilities. As a [AGENT_DOMAIN] specialist, you should prioritize MCPs in this order:

### Primary MCPs (Use First)
**[PRIMARY_CATEGORY_1]**: [Description of why this category is relevant]
- **When to use**: [Specific scenarios for this agent]
- **Tool patterns**: `[pattern1]__*`, `[pattern2]__*`
- **Examples**: [Specific MCP examples if known]

**[PRIMARY_CATEGORY_2]**: [Description]
- **When to use**: [Scenarios]
- **Tool patterns**: `[patterns]`
- **Examples**: [MCPs]

### Secondary MCPs (Use When Relevant)
**[SECONDARY_CATEGORY]**: [Description]
- **When to use**: [Scenarios]
- **Tool patterns**: `[patterns]`

### MCP Discovery Protocol

At the start of each task:
1. Use `ListMcpResourcesTool` to discover available MCPs
2. Identify MCPs matching your primary categories
3. Prefer domain-specific MCP tools over generic alternatives
4. Fall back to Claude Code built-in tools if no suitable MCP exists

### MCP Usage Guidelines

**DO**:
- Check for specialized MCPs before using generic tools
- Use design system MCPs for component/token work (if design agent)
- Use infrastructure MCPs for cloud/deployment work (if infra agent)
- Use knowledge MCPs for documentation research (all agents)
- Invoke MCP tools with full, descriptive parameters

**DON'T**:
- Hardcode specific MCP server names in your logic
- Assume MCPs are available - always discover first
- Use MCPs outside your domain without justification
- Skip MCP discovery for complex or specialized tasks
```

### 2.2 Example: Frontend Agent with MCP Awareness

```markdown
## MCP Awareness: Specialized Tool Priority

You have access to Model Context Protocol (MCP) servers that extend your capabilities. As a Frontend specialist, you should prioritize MCPs in this order:

### Primary MCPs (Use First)
**Design & UI**: Access design systems, component libraries, design tokens
- **When to use**: Building components, implementing designs, accessing design tokens
- **Tool patterns**: `design_system__*`, `component__*`, `token__*`, `mcp__figma__*`
- **Examples**: `lando-labs-design-system`, `storybook-mcp`, `figma-mcp`

**Development**: File operations, code analysis, browser automation
- **When to use**: Reading/writing code, analyzing project structure, testing in browser
- **Tool patterns**: `mcp__filesystem__*`, `mcp__playwright__*`, `mcp__git__*`
- **Examples**: `filesystem`, `playwright`, `git-mcp`

### Secondary MCPs (Use When Relevant)
**Testing & QA**: Test runners, accessibility validation, performance auditing
- **When to use**: Verifying component quality, accessibility compliance, performance
- **Tool patterns**: `mcp__jest__*`, `mcp__axe__*`, `mcp__lighthouse__*`

**Knowledge**: Documentation for libraries, frameworks, APIs
- **When to use**: Learning library APIs, checking framework best practices
- **Tool patterns**: `mcp__context7__*`, `mcp__docs__*`

### MCP Discovery Protocol

At the start of each task:
1. Use `ListMcpResourcesTool` to discover available MCPs
2. Check for design system MCPs if building UI components
3. Use filesystem MCPs for reading project structure
4. Fall back to Claude Code built-in Read/Write tools if needed

### MCP Usage Guidelines

**DO**:
- Check for design system MCP before implementing components
- Use design token MCPs to ensure consistency
- Use knowledge MCPs (context7) to fetch React/framework docs
- Use filesystem MCP for multi-file reads when available

**DON'T**:
- Hardcode `lando-labs-design-system` - discover generically
- Skip MCP discovery when building design-system components
- Use backend/infrastructure MCPs without clear justification
```

### 2.3 Example: DevOps Agent with MCP Awareness

```markdown
## MCP Awareness: Specialized Tool Priority

You have access to Model Context Protocol (MCP) servers that extend your capabilities. As a DevOps specialist, you should prioritize MCPs in this order:

### Primary MCPs (Use First)
**Infrastructure**: Cloud providers, Kubernetes, databases, monitoring
- **When to use**: Deploying apps, managing infrastructure, monitoring systems
- **Tool patterns**: `mcp__firebase__*`, `mcp__k8s__*`, `mcp__aws__*`, `mcp__datadog__*`
- **Examples**: `firebase`, `kubernetes-mcp`, `aws-mcp`, `terraform-mcp`

**Development**: File operations, git, code analysis
- **When to use**: Reading configs, modifying YAML/HCL, checking git history
- **Tool patterns**: `mcp__filesystem__*`, `mcp__git__*`
- **Examples**: `filesystem`, `git-mcp`

### Secondary MCPs (Use When Relevant)
**Testing & QA**: Performance testing, smoke tests, integration tests
- **When to use**: Validating deployments, running health checks
- **Tool patterns**: `mcp__lighthouse__*`, `mcp__cypress__*`

**Communication**: Slack, Linear, PagerDuty for deployment notifications
- **When to use**: Notifying team of deployments, creating incidents
- **Tool patterns**: `mcp__slack__*`, `mcp__pagerduty__*`

### MCP Discovery Protocol

At the start of each task:
1. Use `ListMcpResourcesTool` to discover available MCPs
2. Check for cloud provider MCPs (Firebase, AWS, GCP, Azure)
3. Check for infrastructure MCPs (Kubernetes, Terraform)
4. Fall back to Bash for manual operations if no MCP exists

### MCP Usage Guidelines

**DO**:
- Use Firebase MCP for Firebase-specific operations
- Use Kubernetes MCP for cluster management
- Use monitoring MCPs for observability setup
- Use communication MCPs for deployment notifications

**DON'T**:
- Hardcode specific cloud provider names
- Skip MCP discovery for cloud operations
- Use design/frontend MCPs without justification
```

## 3. Concrete Use Case Examples

### Use Case 1: Design System Specialist Building a Button Component

**Scenario**: User asks design-system-specialist to create a Button component.

**MCP Discovery Flow**:
1. Agent invokes `ListMcpResourcesTool`
2. Discovers `lando-labs-design-system` MCP with tools:
   - `design_system__get_tokens`
   - `design_system__get_component`
   - `design_system__list_components`
3. Identifies as "Design & UI" category (Primary for this agent)
4. Invokes `design_system__get_tokens` to fetch color/spacing tokens
5. Uses tokens to build Button component with consistent styling
6. Falls back to `Read` tool for examining existing components

**Without MCP Awareness**: Agent would use generic `Read` tool, might miss design tokens, risk inconsistency.

### Use Case 2: Backend Agent Setting Up Firebase

**Scenario**: User asks backend agent to initialize Firestore database.

**MCP Discovery Flow**:
1. Agent invokes `ListMcpResourcesTool`
2. Discovers `firebase` MCP with resources:
   - `firebase://guides/init/firestore`
   - `firebase://guides/init/backend`
3. Identifies as "Infrastructure" category (Secondary for backend agent)
4. Reads `firebase://guides/init/firestore` resource for setup steps
5. Invokes `mcp__firebase__firebase_init` with Firestore feature
6. Uses Firebase MCP for login, project creation

**Without MCP Awareness**: Agent would use Bash commands, might miss Firebase CLI shortcuts, error-prone manual setup.

### Use Case 3: Frontend Agent with Multiple Design System MCPs

**Scenario**: Project has two design system MCPs:
- `company-design-system` (company-wide)
- `project-ui-library` (project-specific)

**MCP Discovery Flow**:
1. Agent discovers both MCPs (both "Design & UI" category)
2. Prioritizes `project-ui-library` (more specific > more generic)
3. Falls back to `company-design-system` if component not found
4. Documents which MCP was used for future consistency

**Pattern Recognition**:
- Project-specific > Company-wide > Public library
- Local > Remote
- Versioned > Unversioned

### Use Case 4: DevOps Agent Deploying to Multiple Clouds

**Scenario**: Multi-cloud deployment to Firebase and AWS.

**MCP Discovery Flow**:
1. Agent discovers `firebase` and `aws-mcp` MCPs
2. Uses `firebase` MCP for Firebase Hosting deployment
3. Uses `aws-mcp` for AWS Lambda deployment
4. Uses `datadog-mcp` for unified monitoring setup
5. Documents multi-cloud architecture in reference/

**Cross-MCP Coordination**:
- Different MCPs for different concerns
- Agent orchestrates multiple MCPs in single workflow
- Maintains awareness of which MCP does what

### Use Case 5: QA Agent with Testing Framework MCPs

**Scenario**: User asks QA agent to run accessibility and performance tests.

**MCP Discovery Flow**:
1. Agent discovers `axe-core-mcp`, `lighthouse-mcp`, `jest-mcp`
2. Categorizes as "Testing & QA" (Primary for QA agent)
3. Uses `axe-core-mcp` for accessibility audit
4. Uses `lighthouse-mcp` for performance audit
5. Uses `jest-mcp` for unit test execution
6. Aggregates results into comprehensive report

**MCP Composition**: Agent uses multiple MCPs for comprehensive testing strategy.

### Use Case 6: Research Agent Using Knowledge MCPs

**Scenario**: User asks research-synthesizer to learn about React Server Components.

**MCP Discovery Flow**:
1. Agent discovers `context7` MCP (Knowledge category, Primary for research)
2. Invokes `mcp__context7__resolve-library-id` with "React"
3. Gets library ID `/facebook/react/v19.0.0`
4. Invokes `mcp__context7__get-library-docs` with topic "server components"
5. Synthesizes documentation into research report

**Without MCP Awareness**: Agent would use WebSearch, might get outdated/incorrect information.

### Use Case 7: AI/ML Specialist Using Vector Database

**Scenario**: User asks ai-ml-specialist to implement semantic search.

**MCP Discovery Flow**:
1. Agent discovers `pinecone-mcp` (AI & ML category, Primary)
2. Uses `mcp__pinecone__create_index` to set up vector index
3. Uses `mcp__pinecone__upsert_vectors` to add embeddings
4. Uses `mcp__pinecone__query` for similarity search
5. Integrates with `openai-mcp` for generating embeddings

**Multi-MCP AI Workflow**: Combines embedding generation + vector storage MCPs.

### Use Case 8: Security Specialist Auditing Infrastructure

**Scenario**: User asks security-specialist to audit Kubernetes security.

**MCP Discovery Flow**:
1. Agent discovers `kubernetes-mcp`, `trivy-mcp` (Infrastructure + Testing)
2. Uses `mcp__k8s__get_security_policies` to review policies
3. Uses `mcp__trivy__scan_cluster` for vulnerability scan
4. Uses `mcp__aws__audit_iam` to check cloud permissions
5. Generates comprehensive security report

**Cross-Domain MCP Usage**: Security agent uses Infrastructure, Testing, and Development MCPs.

### Use Case 9: Mobile Agent with Cross-Platform Testing

**Scenario**: User asks mobile-native agent to test on multiple devices.

**MCP Discovery Flow**:
1. Agent discovers `browserstack-mcp`, `firebase-test-lab-mcp`
2. Uses `mcp__browserstack__test_ios` for iOS devices
3. Uses `mcp__firebase_test_lab__test_android` for Android devices
4. Aggregates test results across platforms

**Platform-Specific MCPs**: Agent knows which MCP for which platform.

### Use Case 10: No Suitable MCP Available

**Scenario**: User asks frontend agent to minify CSS, but no CSS minification MCP exists.

**MCP Discovery Flow**:
1. Agent invokes `ListMcpResourcesTool`
2. Searches for CSS-related tools in Development category
3. Finds no suitable MCP
4. Falls back to Claude Code built-in Bash tool
5. Uses `npx postcss` or similar CLI tool
6. Documents manual approach

**Graceful Degradation**: Agent uses built-in tools when no MCP available.

## 4. Implementation Recommendations

### 4.1 Rollout Strategy

#### Phase 1: Core Agents (Week 1-2)
- **Agents**: frontend, backend, devops, designer, qa
- **Action**: Add MCP Awareness section to system prompts
- **Testing**: Verify MCP discovery works in sample tasks
- **Validation**: Agents prefer MCPs over built-in tools when available

#### Phase 2: Specialized Agents (Week 3-4)
- **Agents**: react-specialist, design-system-specialist, gcp-firebase, api-integrator
- **Action**: Add domain-specific MCP priorities
- **Testing**: Test with actual project MCPs
- **Validation**: Agents use domain MCPs correctly

#### Phase 3: Research & AI Agents (Week 5-6)
- **Agents**: research-synthesizer, ai-ml-specialist, use-case-specialist
- **Action**: Emphasize Knowledge and AI/ML MCP categories
- **Testing**: Verify context7, sequential-thinking usage
- **Validation**: Research agents leverage knowledge MCPs

#### Phase 4: Infrastructure & Security (Week 7-8)
- **Agents**: security-specialist, performance-optimizer, deploy
- **Action**: Add Infrastructure and Testing MCP priorities
- **Testing**: Test cloud provider MCPs, monitoring MCPs
- **Validation**: Infrastructure agents orchestrate multiple MCPs

### 4.2 Agent Prompt Template Integration

When creating/updating agents, use this checklist:

**Before "Three-Phase Methodology"**:
- [ ] Add "MCP Awareness: Specialized Tool Priority" section
- [ ] Define Primary MCPs (2-3 categories relevant to agent)
- [ ] Define Secondary MCPs (1-2 categories for cross-domain work)
- [ ] Include MCP Discovery Protocol (4-step process)
- [ ] Include MCP Usage Guidelines (DO/DON'T)

**In "Phase 1: Research/Analyze"**:
- [ ] Add "Discover available MCPs using ListMcpResourcesTool"
- [ ] Add "Identify domain-specific MCPs for this task"
- [ ] Reference MCP resources in tool list

**In "Phase 2: Build/Core Action"**:
- [ ] Add "Use [Primary MCP Category] tools when available"
- [ ] Add "Fall back to Claude Code built-in tools if needed"

**In "Phase 3: Verify/Follow-up"**:
- [ ] Add "Document which MCPs were used"
- [ ] Add "Note any MCP limitations encountered"

### 4.3 MCP Configuration Documentation

Create `reference/mcp-configuration.md` with:

```markdown
# MCP Configuration Guide

## Configured MCPs

List all MCPs configured in Claude Desktop or project:

### Design & UI
- **lando-labs-design-system** (v1.0.0)
  - Tools: design_system__get_tokens, design_system__get_component
  - Resources: design://tokens, design://components
  - Used by: designer, frontend, design-system-specialist

### Infrastructure
- **firebase** (v1.0.0)
  - Tools: firebase_init, firebase_deploy, firebase_get_environment
  - Resources: firebase://guides/*
  - Used by: backend, gcp-firebase, devops, deploy

### Knowledge
- **context7** (latest)
  - Tools: resolve-library-id, get-library-docs
  - Used by: All agents for library documentation

## Adding New MCPs

When adding a new MCP:
1. Install MCP in Claude Desktop config
2. Update this file with MCP details
3. Identify which agent categories should use it
4. Test with relevant agents
5. Update agent system prompts if needed (new categories)
```

### 4.4 User-Facing MCP Discovery

Add to CLAUDE.md:

```markdown
## Model Context Protocol (MCP) Servers

This project uses MCP servers to extend agent capabilities:

### Design & UI
- [MCP name] - Component library and design tokens

### Infrastructure
- [MCP name] - Firebase backend integration
- [MCP name] - Kubernetes deployment

### Knowledge
- context7 - Library documentation provider

### Configuring MCPs

MCPs are configured in Claude Desktop settings. See `reference/mcp-configuration.md` for details.

### Which Agent Uses Which MCP?

- **Design agents** (designer, ux, frontend): Design & UI MCPs
- **Backend agents** (backend, api-integrator): Infrastructure, Data MCPs
- **DevOps agents** (devops, deploy): Infrastructure, Testing MCPs
- **Research agents** (research-synthesizer): Knowledge MCPs
- **All agents**: Knowledge MCPs for library docs

Agents automatically discover and prioritize MCPs based on their domain.
```

### 4.5 Monitoring MCP Usage

Create analytics/logging for MCP adoption:

```markdown
## MCP Usage Metrics

Track in agent documentation:
- Which MCPs were discovered
- Which MCPs were actually used
- Which tasks benefited from MCPs
- Which agents used which MCPs most

Example log:
```
[frontend] Task: Build Button component
[frontend] Discovered MCPs: lando-labs-design-system, filesystem, context7
[frontend] Used: lando-labs-design-system (design_system__get_tokens)
[frontend] Fallback: Read (for existing components)
[frontend] Outcome: Success - component matches design system
```
```

## 5. Advanced Patterns

### 5.1 MCP Composition Workflows

Agents should compose multiple MCPs for complex workflows:

**Example: Full-Stack Feature Development**

```
User Request: "Build a user profile page with Firebase backend"

Architect Agent:
  - Uses sequential-thinking MCP for planning
  - Uses context7 MCP for Next.js docs
  - Designs architecture, delegates to specialists

Frontend Agent (delegated):
  - Discovers: design-system-mcp, react-mcp, filesystem
  - Uses design-system-mcp for UI components
  - Uses filesystem for reading project structure
  - Builds UI with design tokens

Backend Agent (delegated):
  - Discovers: firebase-mcp, postgres-mcp
  - Uses firebase-mcp for Firestore schema
  - Uses firebase-mcp for auth integration
  - Builds API endpoints

DevOps Agent (delegated):
  - Discovers: firebase-mcp, github-actions-mcp
  - Uses firebase-mcp for deployment
  - Uses github-actions-mcp for CI/CD
  - Configures hosting
```

### 5.2 MCP Preference Overrides

Allow users to override MCP preferences:

**In CLAUDE.md**:
```markdown
## MCP Preferences (User Overrides)

Override default MCP priorities for this project:

```yaml
mcp_preferences:
  design_agents:
    primary:
      - "project-ui-library"  # Use project-specific first
      - "company-design-system"  # Then company-wide

  infrastructure_agents:
    primary:
      - "aws-mcp"  # Prefer AWS over GCP for this project
      - "terraform-mcp"

  exclude:
    - "deprecated-mcp"  # Don't use this MCP
```
```

Agents should read this config during MCP discovery.

### 5.3 MCP Capability Negotiation

Agents should verify MCP capabilities before use:

```
1. Discover MCP with ListMcpResourcesTool
2. Check tool list for required capability (e.g., "design_system__get_tokens")
3. If capability missing, check alternative MCPs
4. If no MCP has capability, fall back to built-in tools
5. Document capability gaps for user awareness
```

### 5.4 MCP Versioning Strategy

Handle MCP version differences:

```markdown
## MCP Versioning Awareness

Agents should:
1. Check MCP version if exposed in tool metadata
2. Prefer stable versions over beta/experimental
3. Document version used in outputs
4. Warn if deprecated MCP features detected
5. Suggest upgrades if newer version available

Example:
```
[frontend] Using lando-labs-design-system v1.2.0
[frontend] Note: v2.0.0 available with new token structure
[frontend] Recommend upgrading for modern design tokens
```
```

### 5.5 Fallback Hierarchy

Define clear fallback paths:

```
1. Primary MCP in category (e.g., project-specific design system)
2. Secondary MCP in category (e.g., company design system)
3. Tertiary MCP in related category
4. Claude Code built-in tools (Read, Write, Bash, etc.)
5. Bash commands for CLI tools
6. Request user guidance if all fallbacks fail
```

## 6. Testing & Validation

### 6.1 Agent MCP Awareness Tests

Test each agent's MCP awareness:

**Test 1: MCP Discovery**
```
Given: Agent is invoked for domain-specific task
When: Agent begins task
Then: Agent calls ListMcpResourcesTool
And: Agent identifies MCPs in primary categories
```

**Test 2: MCP Prioritization**
```
Given: Multiple MCPs in same category available
When: Agent needs specific capability
Then: Agent uses most specific/relevant MCP first
And: Falls back to generic MCP if specific fails
```

**Test 3: Cross-Category Fallback**
```
Given: No primary category MCP available
When: Agent needs capability
Then: Agent checks secondary category MCPs
And: Uses tertiary or built-in tools if needed
```

**Test 4: Graceful Degradation**
```
Given: No suitable MCP available
When: Agent attempts MCP discovery
Then: Agent falls back to built-in tools without error
And: Documents manual approach used
```

### 6.2 Integration Test Scenarios

**Scenario 1: Design System Workflow**
```
Agents: designer, frontend, design-system-specialist
MCPs: lando-labs-design-system, storybook-mcp
Test: Build new component using design tokens
Validation: All agents use design-system MCP, tokens are consistent
```

**Scenario 2: Firebase Deployment**
```
Agents: backend, gcp-firebase, devops
MCPs: firebase, github-actions-mcp
Test: Deploy Firestore + Cloud Functions
Validation: Agents use firebase MCP, deployment succeeds
```

**Scenario 3: Multi-Cloud Architecture**
```
Agents: architect, backend, devops
MCPs: firebase, aws-mcp, kubernetes-mcp
Test: Design and deploy multi-cloud system
Validation: Agents use appropriate MCPs per cloud provider
```

## 7. Documentation Standards

### 7.1 Agent Documentation

Each agent should document:

```markdown
## MCP Usage

This agent prioritizes these MCP categories:
- **Primary**: [Categories]
- **Secondary**: [Categories]

### Typical MCP Workflows

**Task: [Common Task]**
- Discovers: [MCPs]
- Uses: [Specific tools]
- Falls back to: [Alternatives]

**Task: [Another Common Task]**
- [Similar documentation]
```

### 7.2 Project Documentation

In `reference/mcp-configuration.md`:

```markdown
# MCP Configuration

## Installed MCPs

### [MCP Name] (Category: [Category])
- **Version**: [Version]
- **Description**: [What it does]
- **Tools**: [List of tools]
- **Resources**: [List of resources]
- **Used by agents**: [Agent list]
- **Configuration**: [Any config needed]

## MCP Usage Patterns

Document how agents use MCPs together:

### Design System + Frontend Workflow
1. Designer creates design tokens with design-system MCP
2. Frontend agent reads tokens with same MCP
3. React-specialist builds components using tokens

### Backend + Infrastructure Workflow
1. Backend agent initializes Firebase with firebase MCP
2. DevOps agent deploys with firebase MCP + github-actions-mcp
3. Monitoring configured with datadog-mcp
```

## 8. Future Enhancements

### 8.1 MCP Discovery Service

Build dedicated MCP discovery service:

```typescript
interface MCPDiscoveryService {
  discoverAll(): MCP[];
  categorize(mcp: MCP): MCPCategory;
  prioritize(mcps: MCP[], agentRole: AgentRole): MCP[];
  selectBest(capability: string, mcps: MCP[]): MCP | null;
  recordUsage(agent: string, mcp: string, task: string): void;
}
```

### 8.2 MCP Recommendation Engine

Suggest MCPs to users based on project:

```
Analyze project structure:
  - If React project → Suggest react-mcp, storybook-mcp
  - If Firebase detected → Suggest firebase-mcp
  - If AWS detected → Suggest aws-mcp
  - If Kubernetes detected → Suggest kubernetes-mcp

Generate recommendation report:
  "Based on your project, consider installing these MCPs:
   - firebase-mcp (for backend agents)
   - lando-labs-design-system (for design agents)
   - github-actions-mcp (for devops agents)"
```

### 8.3 MCP Capability Registry

Maintain registry of MCP capabilities:

```yaml
mcp_registry:
  firebase:
    category: Infrastructure
    capabilities:
      - database_operations
      - authentication
      - hosting_deployment
      - cloud_functions
    tools:
      - firebase_init
      - firebase_deploy
    resources:
      - firebase://guides/*

  lando-labs-design-system:
    category: Design & UI
    capabilities:
      - design_tokens
      - component_library
      - theme_system
    tools:
      - design_system__get_tokens
      - design_system__get_component
    resources:
      - design://tokens
      - design://components
```

### 8.4 MCP Analytics Dashboard

Track MCP usage across agents:

```
Metrics to track:
- MCP discovery rate per agent
- MCP usage rate per agent
- Most-used MCPs overall
- MCP fallback frequency
- Task success rate with vs without MCPs
- Agent satisfaction with MCP capabilities

Dashboard views:
- MCP adoption by agent category
- MCP usage trends over time
- MCP capability gaps (requested but unavailable)
- Agent-MCP affinity heatmap
```

## 9. Conclusion

This MCP-awareness architecture enables Claude Code agents to:

1. **Discover MCPs dynamically** without hardcoded dependencies
2. **Prioritize domain-relevant MCPs** based on agent specialization
3. **Compose multiple MCPs** for complex workflows
4. **Degrade gracefully** when MCPs are unavailable
5. **Scale seamlessly** as new MCPs are added to the ecosystem

By embedding MCP awareness into agent system prompts and providing clear discovery patterns, we create a flexible, powerful multi-agent + multi-MCP system that amplifies both agent capabilities and user productivity.

---

**Next Steps**:
1. Review and approve architecture
2. Update agent-architect system prompt with MCP awareness guidelines
3. Roll out to Phase 1 agents (frontend, backend, devops, designer, qa)
4. Create reference/mcp-configuration.md for project
5. Monitor adoption and iterate based on feedback
