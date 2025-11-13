<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
## Deployed Agents

The following Claude Code agents are available in this project:

### accessibility-expert (v1.1.0)
Use this agent when ensuring WCAG compliance, implementing accessible UI, testing with screen readers, or creating inclusive designs. Invoke for accessibility audits, ARIA implementation, keyboard navigation, screen reader testing, or WCAG 2.1/2.2 compliance validation.

### agent-architect (v1.1.0)
Use this agent PROACTIVELY when you need to create, refine, or optimize Claude Code agent configurations. This includes designing new agents from scratch, improving existing agent system prompts, establishing agent interaction patterns, defining agent responsibilities and boundaries, or architecting multi-agent systems with clear separation of concerns.

### ai-ml-specialist (v1.1.0)
Use this agent when integrating AI/ML capabilities, deploying models, building ML pipelines, or implementing intelligent features. Invoke for LLM integration, model serving, ML pipeline development, prompt engineering, or AI-powered feature implementation.

### api-integrator (v1.1.0)
Use this agent when integrating third-party APIs, designing RESTful/GraphQL APIs, implementing webhooks, or building API gateways. Invoke for API client development, webhook handling, API design, authentication flows, rate limiting, or API documentation.

### architect (v1.1.0)
Use this agent when you need to analyze system requirements, plan technical architecture, or guide the evolution of existing systems. Invoke for architecture decisions, system design, technology stack selection, scalability planning, or refactoring guidance.

### backend (v1.1.0)
Use this agent when building APIs, databases, server-side logic, or optimizing backend performance. Invoke for REST/GraphQL API development, database schema design, authentication/authorization, data processing, caching strategies, or performance tuning.

### deploy (v1.1.0)
Use this agent when configuring deployment infrastructure, building CI/CD pipelines, containerizing applications, or monitoring production systems. Invoke for Docker/Kubernetes setup, deployment automation, infrastructure as code, monitoring configuration, or production troubleshooting.

### designer (v1.1.0)
Use this agent when evaluating visual design, crafting design systems, ensuring aesthetic quality, or creating timeless visual experiences. Invoke for color palette design, typography systems, visual hierarchy, design system creation, brand consistency, or aesthetic refinement.

### devops (v1.1.0)
Use this agent when building CI/CD pipelines, implementing infrastructure as code, setting up monitoring/observability, or automating deployment workflows. Invoke for GitHub Actions, Terraform, container orchestration, log aggregation, or GitOps implementations.

### frontend (v1.1.0)
Use this agent when building or maintaining user interface components, setting up styling systems, managing React applications, or ensuring frontend consistency. Invoke for component creation, state management, CSS/styling, accessibility implementation, or frontend performance optimization.

### mcp-specialist (v1.1.0)
Use this agent when building Model Context Protocol (MCP) servers, creating tools/resources/prompts, or integrating MCP capabilities into applications. Invoke for MCP server development, tool design, resource providers, protocol implementation, TypeScript/Node.js MCP projects, or extending Claude Desktop with custom context.

### qa (v1.1.0)
Use this agent when writing tests, analyzing test coverage, creating testing documentation, or maintaining testing standards. Invoke for unit tests, integration tests, E2E tests, test coverage analysis, test strategy planning, or quality assurance automation.

### terminal-specialist (v1.1.0)
Use this agent when you need to build, analyze, or optimize terminal applications and CLI tools. This includes designing interactive terminal UIs, implementing terminal-based workflows, working with ANSI escape codes, handling terminal compatibility, or optimizing terminal rendering performance.

### ux (v1.1.0)
Use this agent when designing user experiences, analyzing user workflows, creating interaction patterns, or validating usability. Invoke for user journey mapping, information architecture, interaction design, usability analysis, or UX research and validation.

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->

## Using CAMI (Claude Agent Management Interface)

CAMI is the tool that manages the agents listed above. As Claude Code, you can use CAMI's MCP server to help users manage agents across their projects.

### Available MCP Tools

You have access to these CAMI MCP tools:

1. **mcp__cami__list_agents** - List all available agents from configured sources
   - Use when: User asks "what agents are available?" or wants to discover agents
   - Returns: Agent names, versions, descriptions, and categories

2. **mcp__cami__deploy_agents** - Deploy agents to a project's `.claude/agents/` directory
   - Use when: User wants to add agents to a project
   - Parameters: `agent_names` (array), `target_path` (absolute path), `overwrite` (optional)
   - Handles: Conflict detection, directory creation

3. **mcp__cami__scan_deployed_agents** - Scan a project to see what agents are deployed
   - Use when: User asks "what agents are installed?" or wants to audit agents
   - Parameters: `target_path` (absolute path)
   - Returns: Agent status (up-to-date, update-available, not-deployed)

4. **mcp__cami__update_claude_md** - Update a project's CLAUDE.md with agent documentation
   - Use when: After deploying agents, to keep documentation in sync
   - Parameters: `target_path` (absolute path)
   - Auto-generates the "Deployed Agents" section

5. **mcp__cami__add_location** - Register a project directory for agent deployment
   - Use when: User wants to track a project in CAMI
   - Parameters: `name` (friendly name), `path` (absolute path)

6. **mcp__cami__list_locations** - List all registered project locations
   - Use when: User asks "what projects am I tracking?"
   - Returns: Location names and paths

7. **mcp__cami__remove_location** - Unregister a project directory
   - Use when: User wants to stop tracking a project
   - Parameters: `name` (location name)

8. **mcp__cami__list_sources** - List all configured agent sources
   - Use when: User asks "where are my agents from?" or wants to see sources
   - Returns: Source names, paths, priorities, agent counts, git info

9. **mcp__cami__add_source** - Add a new agent source by cloning a Git repository
   - Use when: User wants to add official agents (lando-agents) or company sources
   - Parameters: `url` (git URL), `name` (optional), `priority` (optional, default 100)
   - Clones to vc-agents/<name>/ and updates config

10. **mcp__cami__update_source** - Update agent sources with git pull
   - Use when: User wants to get latest agents from sources
   - Parameters: `name` (optional, updates all if not specified)
   - Updates sources that have git remotes

11. **mcp__cami__source_status** - Show git status of agent sources
   - Use when: User wants to check if sources have uncommitted changes
   - Shows which sources are clean or have local modifications

### Common Workflows

**When user wants to add the official Lando agent library:**
```
1. Use mcp__cami__add_source with URL: git@github.com:lando-labs/lando-agents.git
2. Use mcp__cami__list_agents to show newly available agents
3. Ask user which agents they want to deploy
4. Use mcp__cami__deploy_agents to add them to current project
5. Use mcp__cami__update_claude_md to document the deployment
```

**When user wants to add agents to their current project:**
```
1. Use mcp__cami__list_agents to show available agents
2. Ask user which agents they want
3. Use mcp__cami__deploy_agents with current directory
4. Use mcp__cami__update_claude_md to document the deployment
```

**When user asks "what agents do I have?":**
```
1. Use mcp__cami__scan_deployed_agents with current directory
2. Show which agents are deployed and their status
3. Suggest updates if any agents are outdated
```

**When user wants to update agents from sources:**
```
1. Use mcp__cami__update_source (updates all sources)
2. Use mcp__cami__scan_deployed_agents to see which deployed agents have updates
3. Ask if user wants to redeploy updated agents
4. Use mcp__cami__deploy_agents with overwrite=true if needed
```

**When working with agent-architect to create new agents:**
```
1. Agent-architect creates the agent file
2. Save it to vc-agents/my-agents/ or appropriate source directory
3. Use mcp__cami__deploy_agents to deploy it to projects
4. Use mcp__cami__update_claude_md to document it
```

### CLI Commands (for reference)

Users can also run CAMI commands directly:

- `./cami list` - List available agents
- `./cami deploy <agents...> <path>` - Deploy agents
- `./cami scan <path>` - Scan deployed agents
- `./cami update-docs <path>` - Update CLAUDE.md
- `./cami locations list` - List tracked locations
- `./cami source list` - List agent sources
- `./cami source add <git-url>` - Add a new agent source

### Multi-Source Architecture

CAMI supports multiple agent sources with priority-based deduplication:

- **vc-agents/my-agents/** (priority 200) - User's local agents
- **vc-agents/lando-agents/** (priority 100) - Official Lando agents
- **vc-agents/company-agents/** (priority 150) - Company/team agents

When the same agent exists in multiple sources, the highest priority wins.

## Architecture Documentation

### Open Source Strategy
Comprehensive planning document covering CAMI's path to open source release. See [reference/open-source-strategy.md](reference/open-source-strategy.md) for:
- Testing architecture and coverage targets
- Remote agent source design (GitHub, Git, local)
- Multi-workflow support (developers, consumers, teams)
- Repository strategy (tool vs agents separation)
- CLI UX improvements and command structure
- Development workflow documentation
- Phase-by-phase implementation roadmap
- Open source readiness checklist

### Quick Links
- [Testing Strategy](reference/open-source-strategy.md#1-testing-architecture)
- [Remote Sources](reference/open-source-strategy.md#2-remote-agent-sources-design)
- [Workflows](reference/open-source-strategy.md#3-multi-workflow-support)
- [Repository Strategy](reference/open-source-strategy.md#4-repository-strategy)
- [CLI UX](reference/open-source-strategy.md#5-cli-ux-improvements)
- [Roadmap](reference/open-source-strategy.md#implementation-roadmap)
