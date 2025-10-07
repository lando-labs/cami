# CAMI UX Recommendations
## User Experience Design for Agent Discovery, Deployment Strategies & Claude.md Orchestration

---

## Executive Summary

This document provides comprehensive UX recommendations for CAMI's next evolution: agent discovery, deployment strategies, and configuration file orchestration. The recommendations prioritize **cognitive simplicity**, **keyboard efficiency**, and **progressive disclosure** - revealing complexity only when users need it.

**Core UX Principles Applied:**
1. **Recognition over Recall** - Show what's deployed, don't make users remember
2. **Visibility of System State** - Always show current deployment status
3. **User Control & Freedom** - Easy to preview, deploy, and rollback
4. **Consistency** - Maintain vim-style navigation patterns
5. **Error Prevention** - Show conflicts before they happen

---

## 1. User Workflows & Interaction Patterns

### 1.1 Identified User Personas

**Persona 1: The Project Bootstrapper**
- Goal: Quickly set up agents for a new project
- Pain Point: Doesn't know which agents are needed
- Need: Pre-configured deployment strategies

**Persona 2: The Multi-Project Maintainer**
- Goal: Keep agents synchronized across multiple projects
- Pain Point: Losing track of what's deployed where
- Need: Discovery view showing deployment status

**Persona 3: The Power User / Customizer**
- Goal: Fine-tune agent deployments per project type
- Pain Point: Strategies are too rigid
- Need: Customizable strategies and manual overrides

### 1.2 Recommended Information Architecture

**Expand from 4 views to 6 views:**

```
Current Views:
1. Agent Selection       (ViewAgentSelection)
2. Location Management   (ViewLocationManagement)
3. Deployment           (ViewDeployment)
4. Results              (ViewResults)

NEW Views:
5. Discovery Dashboard   (ViewDiscovery)         - NEW PRIMARY VIEW
6. Strategy Selection    (ViewStrategySelection)  - NEW
```

**Navigation Flow:**

```
┌─────────────────────────────────────────────────────────────┐
│                    Discovery Dashboard                       │
│  (Default/Home View - Shows deployment status at a glance)  │
│                                                              │
│  Press 'd' → Deploy Mode                                     │
│  Press 'l' → Locations                                       │
│  Press 'a' → Advanced/Manual Agent Selection                 │
│  Press 's' → Strategy Management                             │
└─────────────────────────────────────────────────────────────┘
                              │
                  ┌───────────┴────────────┐
                  │                        │
                  ▼                        ▼
        ┌─────────────────┐      ┌──────────────────┐
        │ Strategy Select │      │ Manual Selection │
        │   (Quick Path)  │      │  (Power Users)   │
        └─────────────────┘      └──────────────────┘
                  │                        │
                  └───────────┬────────────┘
                              ▼
                    ┌──────────────────┐
                    │   Choose Location│
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   Deploy/Results │
                    └──────────────────┘
```

---

## 2. Agent Discovery UX

### 2.1 Discovery Dashboard Design (New Primary View)

**Mental Model:** "Show me what's where" - A bird's-eye view of all deployments

**Layout:**
```
╔══════════════════════════════════════════════════════════════╗
║ CAMI - Deployment Dashboard                                  ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Locations (3)                                                ║
║                                                              ║
║ > Web App Project (/Users/me/projects/webapp)               ║
║   ├─ architect v1.0.0    ✓                                  ║
║   ├─ frontend  v2.1.0    ✓                                  ║
║   ├─ backend   v2.1.0    ✓                                  ║
║   ├─ qa        v1.5.0    ✓                                  ║
║   └─ ux        v1.0.0    ⚠ v1.2.0 available                 ║
║                                                              ║
║   Mobile App (/Users/me/projects/mobile)                    ║
║   ├─ architect v1.0.0    ✓                                  ║
║   ├─ designer  v1.1.0    ✓                                  ║
║   └─ frontend  v2.0.0    ⚠ v2.1.0 available                 ║
║                                                              ║
║   CLI Tool (/Users/me/projects/cli-tool)                    ║
║   └─ No agents deployed                                      ║
║                                                              ║
║ Available Updates: 2 agents can be updated                   ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ d: deploy  │  u: update all  │  s: strategies  │  a: agents ║
║ l: locations  │  ?: help  │  q: quit                        ║
╚══════════════════════════════════════════════════════════════╝
```

**Key Information Displayed:**
1. **Agent name + version** (e.g., "frontend v2.1.0")
2. **Status indicator:**
   - `✓` Deployed and up-to-date
   - `⚠` Update available (show new version)
   - `✗` Deployed but source deleted
   - `○` Not deployed
3. **Quick stats:** Updates available, total locations
4. **Last sync date** (optional, shown on detail view)

**Interaction Patterns:**
- **j/k or ↓/↑**: Navigate between locations
- **h/l or ←/→**: Expand/collapse location details
- **Enter on location**: Show detailed agent list for that location
- **Enter on agent**: Show agent diff (deployed vs. available)
- **u**: Update all outdated agents
- **d**: Quick deploy to selected location
- **s**: Open strategy selection

### 2.2 Detailed Agent View (Drill-down)

When user presses Enter on a location, expand to show detailed agent info:

```
╔══════════════════════════════════════════════════════════════╗
║ Web App Project - Deployed Agents (5)                       ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ > frontend v2.1.0                                           ║
║   Status:      Up-to-date                                   ║
║   Deployed:    2025-10-01 14:32                             ║
║   File:        .claude/agents/frontend.md                   ║
║   Description: React/TypeScript specialist...               ║
║                                                              ║
║   [Press 'v' to view content, 'r' to remove, 'u' to update] ║
║                                                              ║
║   ux v1.0.0                                     ⚠ Update    ║
║   Status:      v1.2.0 available                             ║
║   Deployed:    2025-09-15 10:21                             ║
║   Changes:     +45 lines, -12 lines                         ║
║                                                              ║
║   [Press 'u' to update, 'd' to view diff]                   ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ esc: back  │  u: update selected  │  r: remove              ║
╚══════════════════════════════════════════════════════════════╝
```

### 2.3 Agent Diff View

When comparing deployed vs. available versions:

```
╔══════════════════════════════════════════════════════════════╗
║ ux Agent - Version Comparison                                ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Deployed: v1.0.0 (2025-09-15)    │    Available: v1.2.0     ║
║                                                              ║
║ Changes:                                                     ║
║   + Added accessibility validation section                   ║
║   + New interaction patterns for mobile                      ║
║   ~ Updated cognitive walkthrough steps                      ║
║   - Removed deprecated heuristics                            ║
║                                                              ║
║ File size: 15.2KB → 16.8KB (+1.6KB)                         ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ u: update now  │  i: ignore this version  │  esc: back      ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 3. Deployment Strategies UX

### 3.1 Strategy Mental Model

**Concept:** "Agent Stacks" - Pre-configured bundles of agents for common project types

**Strategy Definition:**
```yaml
name: "Web Application Stack"
description: "Full-stack web development with React/Vue + Node/Python backend"
agents:
  - architect   # Always include for architectural decisions
  - frontend    # React, Vue, Angular specialists
  - backend     # Node, Python, Go specialists
  - ux          # User experience design
  - qa          # Testing and quality assurance
  - deploy      # CI/CD and deployment
optional:
  - designer    # If design system needed
```

### 3.2 Strategy Selection View

**Two Modes:**
1. **Quick Deploy** (Predefined Strategies)
2. **Custom Strategy** (Build Your Own)

```
╔══════════════════════════════════════════════════════════════╗
║ Deployment Strategies                                        ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Predefined Strategies:                                       ║
║                                                              ║
║ > ⚡ Web Application Stack                        (6 agents) ║
║   architect, frontend, backend, ux, qa, deploy              ║
║   Perfect for: React/Vue + Node/Python web apps              ║
║                                                              ║
║   📱 Mobile Application Stack                    (5 agents) ║
║   architect, frontend, designer, ux, qa                      ║
║   Perfect for: React Native, Flutter apps                    ║
║                                                              ║
║   🖥️  Terminal/CLI Application Stack            (4 agents) ║
║   architect, backend, qa, deploy                             ║
║   Perfect for: CLI tools, system utilities                   ║
║                                                              ║
║   🔧 API/Backend Service Stack                   (5 agents) ║
║   architect, backend, qa, deploy                             ║
║   Perfect for: REST/GraphQL APIs, microservices              ║
║                                                              ║
║   🎨 Design System Project                       (4 agents) ║
║   architect, designer, frontend, qa                          ║
║   Perfect for: Component libraries, design systems           ║
║                                                              ║
║   🏛️  Legacy Modernization Stack                (6 agents) ║
║   architect, backend, qa, deploy                             ║
║   Perfect for: Refactoring, tech debt reduction              ║
║                                                              ║
║   ➕ Custom Strategy...                                      ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: preview & deploy  │  e: edit  │  c: create custom    ║
║ esc: back  │  ?: help                                        ║
╚══════════════════════════════════════════════════════════════╝
```

### 3.3 Strategy Preview (Before Deploy)

**Always show what will be deployed:**

```
╔══════════════════════════════════════════════════════════════╗
║ Preview: Web Application Stack                              ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Will deploy 6 agents to: Web App Project                    ║
║                                                              ║
║ ✓ architect v1.0.0    (not currently deployed)              ║
║ ✓ frontend  v2.1.0    (will update from v2.0.0)            ║
║ ✓ backend   v2.1.0    (already up-to-date)                  ║
║ ✓ ux        v1.2.0    (will update from v1.0.0)            ║
║ ✓ qa        v1.5.0    (not currently deployed)              ║
║ ✓ deploy    v1.0.0    (not currently deployed)              ║
║                                                              ║
║ Summary:                                                     ║
║   2 updates, 3 new deployments, 1 unchanged                  ║
║                                                              ║
║ Claude.md:                                                   ║
║   ☐ Generate new Claude.md with agent instructions          ║
║   ☐ Update existing Claude.md (append agent section)        ║
║   ☑ Skip Claude.md orchestration                            ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: deploy  │  space: toggle agent  │  m: manage claude  ║
║ esc: cancel                                                  ║
╚══════════════════════════════════════════════════════════════╝
```

### 3.4 Custom Strategy Builder

For power users who want to create their own strategies:

```
╔══════════════════════════════════════════════════════════════╗
║ Create Custom Strategy                                       ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Strategy Name: My API Stack_                                 ║
║ Description:   Backend API with comprehensive testing        ║
║                                                              ║
║ Select Agents (7 available):                                 ║
║                                                              ║
║   [ ] architect    System architecture & design              ║
║   [✓] backend      Backend development specialist            ║
║   [ ] frontend     Frontend development specialist           ║
║   [ ] designer     Visual & design system expert             ║
║   [ ] ux           User experience designer                  ║
║   [✓] qa           Quality assurance & testing               ║
║   [✓] deploy       Deployment & CI/CD specialist             ║
║                                                              ║
║ Selected: 3 agents                                           ║
║                                                              ║
║ ☐ Save as reusable strategy                                 ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ space: toggle  │  enter: continue  │  s: save & continue    ║
║ esc: cancel                                                  ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 4. Claude.md Orchestration UX

### 4.1 Orchestration Mental Model

**Concept:** Claude.md is the "project manifest" that tells Claude Code how to work with the deployed agents.

**Three Orchestration Modes:**
1. **Generate New** - Create fresh Claude.md with agent instructions
2. **Update Existing** - Append/update agent section in existing file
3. **Skip** - Deploy agents only, no file changes

### 4.2 Claude.md Preview/Edit View

```
╔══════════════════════════════════════════════════════════════╗
║ Claude.md Orchestration                                      ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Target: /Users/me/projects/webapp/CLAUDE.md                 ║
║ Status: File exists (will update)                            ║
║                                                              ║
║ Deployment: Web Application Stack (6 agents)                 ║
║                                                              ║
║ Preview of changes:                                          ║
║ ┌────────────────────────────────────────────────────────┐  ║
║ │ # Available Agents                                     │  ║
║ │                                                        │  ║
║ │ This project has the following Claude Code agents:    │  ║
║ │                                                        │  ║
║ │ - **architect**: Use for system design and tech...    │  ║
║ │ - **frontend**: React/TypeScript specialist for...    │  ║
║ │ - **backend**: Node.js/Python backend development...  │  ║
║ │ - **ux**: User experience design and usability...     │  ║
║ │ - **qa**: Testing, quality assurance, and...          │  ║
║ │ - **deploy**: CI/CD, containerization, and...         │  ║
║ │                                                        │  ║
║ │ To use an agent, invoke with: @agent-name              │  ║
║ └────────────────────────────────────────────────────────┘  ║
║                                                              ║
║ Insert location:                                             ║
║   ( ) Top of file                                            ║
║   (•) After existing content                                 ║
║   ( ) Custom position...                                     ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: apply  │  e: edit content  │  s: skip  │  esc: back  ║
╚══════════════════════════════════════════════════════════════╝
```

### 4.3 .clinerules Support (Alternative Format)

Some users may prefer .clinerules over CLAUDE.md:

```
╔══════════════════════════════════════════════════════════════╗
║ Configuration File Format                                    ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ Choose configuration format:                                 ║
║                                                              ║
║ > (•) CLAUDE.md     - Markdown format, human-readable       ║
║   ( ) .clinerules   - JSON/YAML format, machine-parseable   ║
║   ( ) Both          - Generate both formats                  ║
║   ( ) Skip          - Deploy agents only                     ║
║                                                              ║
║ ℹ️  CLAUDE.md is recommended for most projects               ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: confirm  │  esc: back                                 ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 5. Suggested New Agents

To enrich deployment strategies, consider these additional agents:

### 5.1 Foundational Agents

**1. data-engineer**
- Use Cases: Database design, data pipelines, ETL
- Strategies: Backend Service, Data Platform, Analytics
- Description: Database architecture, SQL optimization, data modeling

**2. security-specialist**
- Use Cases: Security audits, penetration testing, compliance
- Strategies: All production applications
- Description: Security best practices, OWASP, auth/authz

**3. devops**
- Use Cases: Infrastructure, Kubernetes, cloud platforms
- Strategies: Production deployment, microservices
- Description: IaC, container orchestration, monitoring

### 5.2 Specialized Agents

**4. mobile-native**
- Use Cases: iOS/Android native development
- Strategies: Mobile Application Stack
- Description: Swift, Kotlin, native platform APIs

**5. ai-ml-specialist**
- Use Cases: Machine learning, AI integration
- Strategies: AI/ML Projects, Data Science
- Description: Model training, ML pipelines, inference

**6. docs-writer**
- Use Cases: Documentation, technical writing
- Strategies: Open source projects, API documentation
- Description: API docs, user guides, README generation

**7. performance-optimizer**
- Use Cases: Performance tuning, profiling
- Strategies: Production optimization, legacy modernization
- Description: Benchmarking, memory optimization, caching

**8. accessibility-expert**
- Use Cases: WCAG compliance, inclusive design
- Strategies: Web apps, mobile apps, public-facing
- Description: Screen reader support, keyboard nav, ARIA

### 5.3 Domain-Specific Agents

**9. blockchain-specialist**
- Use Cases: Smart contracts, Web3, DeFi
- Strategies: Blockchain projects
- Description: Solidity, Web3.js, blockchain architecture

**10. game-dev**
- Use Cases: Game development, Unity, Unreal
- Strategies: Game projects
- Description: Game loops, physics, rendering

**11. embedded-systems**
- Use Cases: IoT, firmware, hardware integration
- Strategies: Embedded/IoT projects
- Description: C/C++, RTOS, hardware protocols

**12. api-integrator**
- Use Cases: Third-party API integration, webhooks
- Strategies: Integration-heavy projects
- Description: API design, webhook handling, rate limiting

---

## 6. Enhanced Deployment Strategies

With expanded agent library, here are comprehensive strategies:

### Strategy Matrix

| Strategy | Core Agents | Optional Agents | Use Cases |
|----------|-------------|-----------------|-----------|
| **Full-Stack Web** | architect, frontend, backend, ux, qa | designer, security-specialist, performance-optimizer | SaaS, web apps |
| **Mobile App** | architect, frontend, designer, ux, qa | mobile-native, accessibility-expert | iOS/Android apps |
| **Backend API** | architect, backend, qa, deploy | security-specialist, api-integrator, data-engineer | REST/GraphQL APIs |
| **CLI/Terminal** | architect, backend, qa | docs-writer | Command-line tools |
| **Design System** | architect, designer, frontend, qa | ux, accessibility-expert | Component libraries |
| **Legacy Modernization** | architect, backend, qa, deploy | security-specialist, performance-optimizer | Refactoring projects |
| **Microservices** | architect, backend, qa, devops | security-specialist, data-engineer | Distributed systems |
| **Data Platform** | architect, data-engineer, backend, devops | ai-ml-specialist | Data pipelines, ETL |
| **AI/ML Project** | architect, ai-ml-specialist, data-engineer | backend, devops | ML models, training |
| **Blockchain/Web3** | architect, blockchain-specialist, frontend, security-specialist | backend | DApps, smart contracts |
| **Game Development** | architect, game-dev, designer | performance-optimizer | Unity/Unreal games |
| **IoT/Embedded** | architect, embedded-systems, backend | devops, security-specialist | Firmware, IoT |
| **Open Source Library** | architect, backend, qa, docs-writer | - | npm/pip packages |
| **Enterprise SaaS** | architect, frontend, backend, ux, qa, security-specialist, deploy | performance-optimizer, accessibility-expert | B2B platforms |

---

## 7. User Journeys

### Journey 1: The Quick Bootstrapper

**Goal:** Set up agents for a new React web app in under 60 seconds

**Flow:**
1. Launch CAMI → Discovery Dashboard (shows no deployments)
2. Press `s` → Strategy Selection
3. Navigate to "Web Application Stack" → Press Enter
4. Preview shows 6 agents will deploy → Press Enter to confirm
5. Select "Web App Project" location → Press Enter
6. Claude.md orchestration → Select "Generate new" → Press Enter
7. Results view → 6 agents deployed, CLAUDE.md created
8. Total: 7 keystrokes, 30 seconds

**User Feedback:** "That was incredibly fast. I didn't have to think about which agents to use."

### Journey 2: The Maintainer Syncing Projects

**Goal:** Check which agents need updating across 5 projects

**Flow:**
1. Launch CAMI → Discovery Dashboard
2. See at a glance:
   - Project A: 2 agents need updates (⚠)
   - Project B: All up-to-date (✓)
   - Project C: 1 agent update available (⚠)
3. Press `u` → "Update all outdated agents across all locations?"
4. Confirm → All projects updated in one action
5. Dashboard refreshes → All locations show ✓
6. Total: 3 keystrokes, 15 seconds

**User Feedback:** "I love seeing everything at once. No more guessing what's deployed where."

### Journey 3: The Power User Creating Custom Strategy

**Goal:** Create a custom "API Testing Stack" strategy for microservices

**Flow:**
1. Launch CAMI → Press `s` → Strategy Selection
2. Navigate to "Custom Strategy..." → Press Enter
3. Strategy Builder opens:
   - Enter name: "API Testing Stack"
   - Enter description: "Comprehensive API testing suite"
4. Select agents:
   - Toggle architect ✓
   - Toggle backend ✓
   - Toggle qa ✓
   - Toggle security-specialist ✓
   - Toggle api-integrator ✓
5. Press `s` → "Save as reusable strategy"
6. Strategy saved, appears in strategy list
7. Deploy to 3 microservice projects using saved strategy
8. Total: ~2 minutes first time, then 30 seconds per deployment

**User Feedback:** "Perfect for my team's microservices. I can reuse this strategy across all our APIs."

---

## 8. Keyboard Navigation Enhancements

### 8.1 Global Shortcuts (Available Everywhere)

```
q           Quit application
?           Show context-sensitive help
esc         Back to previous view
ctrl+c      Emergency quit
```

### 8.2 Discovery Dashboard Shortcuts

```
j/k  ↓/↑    Navigate locations
h/l  ←/→    Collapse/expand location details
enter       View location details / agent details
d           Deploy to selected location
s           Open strategy selection
a           Advanced agent selection (current manual flow)
u           Update all outdated agents
l           Manage locations
r           Refresh deployment status
```

### 8.3 Strategy Selection Shortcuts

```
j/k  ↓/↑    Navigate strategies
enter       Preview & deploy selected strategy
e           Edit strategy (custom only)
c           Create custom strategy
d           Delete strategy (custom only)
esc         Back to dashboard
```

### 8.4 Agent Selection Shortcuts (Manual Mode)

```
j/k  ↓/↑    Navigate agents
space / x   Toggle agent selection
a           Select all
n           Select none
/           Filter/search agents (future)
enter       Continue to location selection
esc         Back to dashboard
```

---

## 9. Visual Design Enhancements

### 9.1 Status Indicators

```
✓  Deployed and up-to-date          (Green #42)
⚠  Update available                  (Yellow #214)
✗  Deployed but source missing       (Red #196)
○  Available but not deployed        (Gray #240)
⟳  Updating...                       (Cyan #51)
▶  Expanded location                 (Blue #63)
▸  Collapsed location                (Gray #240)
```

### 9.2 Color Scheme (Lipgloss Styles)

```go
// Status colors
upToDateStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
updateAvailStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Yellow
errorStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
neutralStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray

// UI elements
locationStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
agentNameStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
versionStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
strategyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
```

### 9.3 Progressive Disclosure Pattern

**Level 1 (Overview):**
```
Locations (3)
  Web App ✓  Mobile ⚠  CLI ○
```

**Level 2 (Location Details):**
```
▶ Web App Project
  ├─ architect v1.0.0 ✓
  ├─ frontend v2.1.0 ✓
  └─ backend v2.1.0 ✓
```

**Level 3 (Agent Details):**
```
▶ frontend v2.1.0
  Status: Up-to-date
  Deployed: 2025-10-01 14:32
  [Actions: view, remove, update]
```

---

## 10. Error Prevention & Recovery

### 10.1 Conflict Prevention

**Before deploying, show potential conflicts:**

```
╔══════════════════════════════════════════════════════════════╗
║ ⚠ Deployment Conflicts Detected                              ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ 3 agents already exist at this location:                     ║
║                                                              ║
║   frontend v2.0.0 → v2.1.0  (update available)              ║
║   backend  v2.1.0 → v2.1.0  (same version)                   ║
║   ux       v1.0.0 → v1.2.0  (update available)              ║
║                                                              ║
║ Choose action:                                               ║
║   (•) Update only if newer version                           ║
║   ( ) Overwrite all (force update)                           ║
║   ( ) Skip conflicting agents                                ║
║   ( ) Cancel deployment                                      ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: confirm  │  esc: cancel                               ║
╚══════════════════════════════════════════════════════════════╝
```

### 10.2 Validation Feedback

**Real-time validation during location entry:**

```
Add New Location:

Name: api-server_
Path: /invalid/path/does/not/exist
      ✗ Path does not exist

[enter: save  •  tab: next field  •  esc: cancel]
(Save disabled until path is valid)
```

### 10.3 Rollback Support

**Allow easy rollback of deployments:**

```
╔══════════════════════════════════════════════════════════════╗
║ Deployment History - Web App Project                        ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║ > 2025-10-04 15:30  Web Application Stack (6 agents)        ║
║   2025-10-01 14:32  Manual: frontend, backend (2 agents)    ║
║   2025-09-15 10:21  Mobile Stack (5 agents)                  ║
║                                                              ║
║ Select a deployment to rollback or view details.            ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║ enter: view details  │  r: rollback  │  esc: back           ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 11. Implementation Priorities

### Phase 1: Foundation (Week 1-2)
- [ ] Add Discovery Dashboard as new default view
- [ ] Implement agent discovery logic (scan .claude/agents/)
- [ ] Show deployed agents with version comparison
- [ ] Basic status indicators (✓ ⚠ ✗)

### Phase 2: Strategies (Week 3-4)
- [ ] Define 6-8 core deployment strategies (hardcoded initially)
- [ ] Implement Strategy Selection view
- [ ] Strategy preview before deployment
- [ ] Deploy strategy to location

### Phase 3: Orchestration (Week 5-6)
- [ ] Claude.md generation logic
- [ ] Update existing CLAUDE.md
- [ ] .clinerules support
- [ ] Preview orchestration changes

### Phase 4: Advanced Features (Week 7-8)
- [ ] Custom strategy builder
- [ ] Save/load custom strategies
- [ ] Update all outdated agents (bulk update)
- [ ] Agent diff view

### Phase 5: Polish (Week 9-10)
- [ ] Deployment history
- [ ] Rollback support
- [ ] Search/filter agents
- [ ] Performance optimization

---

## 12. Data Model Extensions

### 12.1 Deployed Agent Metadata

```go
// DeployedAgent represents an agent deployed to a location
type DeployedAgent struct {
    Name         string    `json:"name"`
    Version      string    `json:"version"`
    DeployedAt   time.Time `json:"deployed_at"`
    FilePath     string    `json:"file_path"`
    SourceHash   string    `json:"source_hash"`  // For change detection
}

// LocationStatus represents deployment status for a location
type LocationStatus struct {
    Location       DeployLocation    `json:"location"`
    DeployedAgents []DeployedAgent   `json:"deployed_agents"`
    LastScanned    time.Time         `json:"last_scanned"`
}
```

### 12.2 Strategy Definition

```go
// Strategy represents a deployment strategy
type Strategy struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    AgentIDs    []string `json:"agent_ids"`
    Optional    []string `json:"optional,omitempty"`
    Tags        []string `json:"tags,omitempty"`
    IsCustom    bool     `json:"is_custom"`
    Icon        string   `json:"icon,omitempty"`
}
```

### 12.3 Orchestration Config

```go
// OrchestrationConfig defines how Claude.md should be managed
type OrchestrationConfig struct {
    Mode           string `json:"mode"`           // "generate", "update", "skip"
    Format         string `json:"format"`         // "claude.md", ".clinerules", "both"
    InsertPosition string `json:"insert_position"` // "top", "bottom", "after-marker"
}
```

---

## 13. Accessibility Considerations

### 13.1 Screen Reader Support

- All status indicators have text alternatives
- Navigation states announced clearly
- Error messages are descriptive and actionable

### 13.2 Keyboard-Only Navigation

- No mouse required for any operation
- Tab/Shift+Tab for field navigation in forms
- Vim-style (hjkl) AND arrow keys supported
- Shortcuts are mnemonic (d=deploy, s=strategies, etc.)

### 13.3 Color Independence

- Status indicators use BOTH color AND symbols (✓ ⚠ ✗)
- Important information not conveyed by color alone
- High contrast text for terminal visibility

---

## 14. Future Enhancements (Beyond Current Scope)

### 14.1 Team Collaboration
- Share strategies across team (export/import)
- Remote strategy repository
- Versioned strategies

### 14.2 Agent Marketplace
- Discover community-built agents
- Agent ratings and reviews
- One-click install from registry

### 14.3 Intelligence Features
- Auto-detect project type, suggest strategies
- Agent usage analytics
- Conflict resolution suggestions

### 14.4 Integration
- Git integration (track agent changes in version control)
- CI/CD hooks (auto-deploy agents on pull)
- VS Code extension (GUI wrapper)

---

## 15. Success Metrics

**Measure UX success through:**

1. **Time to First Deployment**
   - Target: < 60 seconds for new users
   - Current: Unknown (measure baseline)

2. **Discovery Efficiency**
   - Users can answer "what's deployed where?" in < 10 seconds
   - Zero cognitive load to check deployment status

3. **Error Rate**
   - < 5% of deployments result in conflicts/errors
   - Conflict resolution success rate > 90%

4. **Strategy Adoption**
   - > 70% of deployments use strategies (vs. manual selection)
   - Custom strategy creation by power users (20-30%)

5. **User Satisfaction**
   - Keyboard navigation feels natural (qualitative)
   - Progressive disclosure reduces overwhelm (qualitative)

---

## Conclusion

These UX recommendations transform CAMI from a simple deployment tool into an intelligent agent orchestration platform. The key principles are:

1. **Visibility First** - Discovery Dashboard shows everything at a glance
2. **Speed & Efficiency** - Strategies enable one-command deployments
3. **Progressive Disclosure** - Simple by default, powerful when needed
4. **Error Prevention** - Conflicts and issues surfaced before they occur
5. **Keyboard Excellence** - Every action optimized for keyboard flow

The recommended information architecture maintains CAMI's terminal-first philosophy while adding the intelligence and automation users need to manage agents across multiple projects effortlessly.

**Next Steps:**
1. Review this UX design with stakeholders
2. Create low-fidelity prototypes of key views
3. Begin Phase 1 implementation (Discovery Dashboard)
4. Gather user feedback early and iterate

---

**Document prepared by:** UX Designer Agent
**Date:** 2025-10-04
**Version:** 1.0
