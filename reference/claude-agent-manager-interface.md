# Product Overview
Simple terminal app that can help with managing my Product Team Agents and deploying them as Claude Code Agents from a version controlled project to other projects. Readily available when a Claude Code session is initiated in that project. This will be a Claude Agent Mangement Interface or CAMI


## Core Features
- Agent Deployment to specific projects
- Deploy location configuration/selection
- Version controlled agents
- claude.md section updates

### Advanced features
- none for now

## Guidance

### Product
- versioned software, doesnt need a production branch, but does need tags
- Simple clean experience
- user triggers app
  - snazzy interface to manage agents and deployments
  - app copies source agents to `{givenProjectLocation}/.claude/agents`

### Technical
- Add any necessary fields for supporting the app like `version` to the YAML frontmatter of agents
- agents should be stored in `vc-agents` folder

### Roadmap 
- Milestone 1: 
  - Agent Deployment
  - Version control
- Mileston 2: 
  - Claude.md updates for guidance on agent orchestration