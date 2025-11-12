---
description: Create a new CAMI agent using agent-architect
argument-hint: <domain/specialty>, [complexity], [requirements]
model: opus
color: cyan
---

# CAMI New Agent Creation

Create a new CAMI agent using the @agent-architect specialist.

## User Requirements
$ARGUMENTS

## Workflow

### Phase 1: Design Agent with @agent-architect

Invoke the @agent-architect agent with these instructions:

"Design a new CAMI agent with the following specifications:

**User Requirements**: $ARGUMENTS

**Technical Requirements**:
- Target Directory: /Users/lando/Development/cami/vc-agents
- Starting Version: 1.0.0
- Must follow CAMI agent architecture standards
- Must include complete YAML frontmatter (name, version, description, tags, use_cases, color, model)
- Must include 'Use this agent PROACTIVELY when...' in description
- Must follow three-phase specialist methodology
- Must select appropriate archetype (Technical Specialist, Quality Guardian, System Designer, or Integration Specialist)
- Include version numbers for technical specialists (e.g., React 19+, Node.js 18+)
- Set appropriate model (opus for complex reasoning, sonnet for most agents)

Please create a production-ready agent file and save it to the vc-agents directory."

Wait for @agent-architect to complete the agent design and save the file.

### Phase 2: Verification

After agent creation, verify:
1. Use `list_agents` MCP tool to confirm the new agent appears in the VC agents list
2. Read the agent file to show the user what was created

### Phase 3: Deployment Instructions

**DO NOT auto-deploy the agent.** Instead, provide the user with deployment instructions using the MCP tools:

```
✅ Agent Created Successfully!

📝 Agent Details:
  • Name: [agent-name]
  • Version: [version]
  • Description: [description]
  • Location: /Users/lando/Development/cami/vc-agents/[agent-name].md

🚀 To deploy this agent to a project, use the CAMI MCP tools:

1. Deploy the agent:
   deploy_agents(
     agent_names: ["[agent-name]"],
     target_path: "/path/to/your/project",
     overwrite: false
   )

2. Update CLAUDE.md:
   update_claude_md(
     target_path: "/path/to/your/project"
   )

The agent will be available via @[agent-name]
```

## Error Handling

If @agent-architect encounters issues:
1. Report the specific problem
2. Suggest corrective actions
3. Offer to retry with clarified requirements

## Quality Checklist

Before presenting results, verify:
- [ ] Agent file exists in vc-agents/
- [ ] Agent has valid YAML frontmatter with all required fields
- [ ] Agent includes PROACTIVELY language in description
- [ ] Agent follows appropriate archetype template
- [ ] Agent has three-phase methodology
- [ ] Version numbers included for technical specialists
- [ ] Appropriate model selected (opus/sonnet)

## Usage Examples

```
/cami-new-agent security auditing specialist, complex, focuses on OWASP and penetration testing

/cami-new-agent Python development expert, moderate complexity, includes FastAPI and async patterns

/cami-new-agent database migration specialist, handles PostgreSQL upgrades and data transformation
```

## Note
This command creates the agent but does NOT deploy it. Use the `deploy_agents` and `update_claude_md` MCP tools to deploy the agent to specific projects after reviewing the generated agent file.
