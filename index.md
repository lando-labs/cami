---
layout: default
title: CAMI - Claude Agent Management Interface
---

# Your AI Agent Guild Headquarters

**CAMI enables Claude Code to manage specialized AI agents across all your projects.** Single binary, clean workspace, natural language interface.

```bash
# Talk to Claude like a person
"Add the frontend agent to this project"
"What agents are available?"
"Update my agents"

# CAMI handles the rest
```

[Download v0.3.0](#installation){: .btn .btn-primary .btn-lg} [View on GitHub](https://github.com/lando-labs/cami){: .btn .btn-outline-secondary .btn-lg}

---

## The Problem: AI Agents Are Scattered

You're building with Claude Code. You've created custom agents for frontend, backend, DevOps, testing. Your team has their own. Open source projects offer specialized agents.

**Where do you keep them all? How do you discover new ones? How do you share with your team?**

Every project becomes a dumping ground for copy-pasted agent files. Version control is manual. Discovery is word-of-mouth. Collaboration is email attachments.

## The Solution: Centralized Agent Management

CAMI gives you a single source of truth for all your AI agents:

```
~/cami/sources/
├── official-agents/     # Public agent libraries (Git)
├── team-agents/         # Your company's agents (Git)
└── my-agents/           # Personal overrides (local)
```

Deploy agents to any project with a conversation:

```
You: "Add the frontend and testing agents to this project"

Claude: ✓ Deployed frontend (v1.1.0)
        ✓ Deployed testing (v1.2.0)
        ✓ Updated CLAUDE.md with agent documentation
```

CAMI runs as an MCP server, giving Claude Code 13 specialized tools for agent management. No manual file copying, no version mismatches, no scattered configuration.

---

## Key Features

### Natural Language Interface

Talk to CAMI through Claude Code. No commands to memorize, no config files to edit manually.

```
"Help me get started with CAMI"
"Add the official agent library"
"What agents do I have deployed?"
"Update agents in this project"
```

### Priority-Based Agent Sources

Override default agents with your custom versions:

```yaml
agent_sources:
  - name: my-agents
    priority: 10         # Highest priority (your custom versions)

  - name: team-agents
    priority: 50         # Medium priority (company standard)

  - name: official-agents
    priority: 100        # Lowest priority (public library)
```

When the same agent exists in multiple sources, lower priority numbers win. Customize without forking.

### Version Tracking & Updates

CAMI tracks agent versions and notifies you of updates:

```
You: "What agents do I have?"

Claude: Found 3 deployed agents:
        • frontend (v1.0.0) → Update available to v1.1.0
        • backend (v1.1.0) → Up to date
        • custom-agent (v1.0.0) → Not found in sources
```

### Git-Trackable Workspace

Your CAMI setup is version-controllable. Share agent sources and configuration with your team:

```bash
cd ~/cami
git init
git add .
git commit -m "Team CAMI setup with company agents"
git remote add origin git@github.com:yourteam/cami-setup.git
git push
```

New team members clone your setup, run the installer, and instantly have access to all team agents.

### Smart Documentation

CAMI automatically maintains your project's `CLAUDE.md` file with deployed agent information:

```markdown
## Deployed Agents

### frontend (v1.1.0)
Use this agent when building user interfaces, React components,
styling, and frontend architecture...

### backend (v1.1.0)
Use this agent for API development, database design, server
architecture, and backend services...
```

No manual documentation maintenance. Always in sync with deployed agents.

### Cross-Platform Single Binary

One binary for macOS, Linux, and Windows. No dependencies, no runtime, no virtual environments.

```bash
# Download, extract, run installer
$ ./install.sh

# CAMI is ready
$ cami --version
CAMI v0.3.0
```

---

## How It Works

### 1. Install CAMI

Download the release for your platform and run the installer:

```bash
# Extract and install
$ cd cami-0.3.0-darwin-arm64
$ ./install.sh

# Creates ~/cami/ workspace
# Installs binary to /usr/local/bin/cami
```

### 2. Add Agent Sources

Add agent libraries from Git repositories:

```bash
$ cd ~/cami
$ claude

You: "Add the official Lando agent library"

Claude: ✓ Cloned lando-agents to ~/cami/sources/lando-agents
        ✓ Found 29 agents
```

### 3. Deploy Agents to Projects

Navigate to any project and deploy agents:

```bash
$ cd ~/projects/my-app
$ claude

You: "Deploy frontend, backend, and qa agents"

Claude: ✓ Deployed frontend (v1.1.0)
        ✓ Deployed backend (v1.1.0)
        ✓ Deployed qa (v1.1.0)
        ✓ Updated CLAUDE.md
```

### 4. Work With Specialized Agents

Agents appear in `.claude/agents/` and are automatically available in Claude Code:

```bash
~/projects/my-app/.claude/agents/
├── frontend.md
├── backend.md
└── qa.md
```

Use agent commands or natural language to invoke specialists.

---

## Installation

### Direct Download (Recommended)

Download pre-built binaries from [GitHub Releases](https://github.com/lando-labs/cami/releases):

**macOS**
- [cami-0.3.0-darwin-amd64.tar.gz](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-darwin-amd64.tar.gz) (Intel)
- [cami-0.3.0-darwin-arm64.tar.gz](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-darwin-arm64.tar.gz) (Apple Silicon)

**Linux**
- [cami-0.3.0-linux-amd64.tar.gz](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-linux-amd64.tar.gz) (x86_64)
- [cami-0.3.0-linux-arm64.tar.gz](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-linux-arm64.tar.gz) (ARM64)

**Windows**
- [cami-0.3.0-windows-amd64.zip](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-windows-amd64.zip) (x86_64)
- [cami-0.3.0-windows-arm64.zip](https://github.com/lando-labs/cami/releases/download/v0.3.0/cami-0.3.0-windows-arm64.zip) (ARM64)

**Extract and Install**

```bash
# macOS/Linux
tar -xzf cami-0.3.0-darwin-arm64.tar.gz
cd cami-0.3.0-darwin-arm64
./install.sh

# Windows (PowerShell)
Expand-Archive cami-0.3.0-windows-amd64.zip
cd cami-0.3.0-windows-amd64
.\install.sh
```

The installer creates:
- `~/cami/` - Your CAMI workspace
- `/usr/local/bin/cami` - Binary on your PATH

### Homebrew (Coming Soon)

```bash
brew install lando-labs/tap/cami
```

### Build From Source

Requires Go 1.21+:

```bash
git clone https://github.com/lando-labs/cami.git
cd cami
make install
```

---

## Getting Started

### First-Time Setup

Open your CAMI workspace and ask Claude for help:

```bash
cd ~/cami
claude
```

```
You: "Help me get started with CAMI"

Claude: I see CAMI isn't configured yet. Let me help you set it up.
        I'll add the official Lando agent library with 29 professional agents.

        ✓ Cloned lando-agents to ~/cami/sources/lando-agents
        ✓ Found 29 agents

        Which agents would you like to add to your projects?
```

### Discover Available Agents

```
You: "What agents are available?"

Claude: You have 29 agents available from lando-agents:

        Development:
        • frontend (v1.1.0) - UI and React development
        • backend (v1.1.0) - API and server development
        • fullstack (v1.1.0) - End-to-end application development

        Quality & Testing:
        • qa (v1.1.0) - Testing and quality assurance
        • code-reviewer (v1.0.0) - Code review and best practices

        DevOps & Infrastructure:
        • devops (v1.1.0) - CI/CD and deployment
        • database (v1.0.0) - Database design and optimization

        [... 22 more agents]
```

### Deploy to a Project

Navigate to your project and deploy agents:

```bash
cd ~/projects/my-web-app
claude
```

```
You: "Add frontend, backend, and qa agents to this project"

Claude: ✓ Created .claude/agents/ directory
        ✓ Deployed frontend (v1.1.0)
        ✓ Deployed backend (v1.1.0)
        ✓ Deployed qa (v1.1.0)
        ✓ Updated CLAUDE.md with agent documentation

        You now have 3 specialized agents available in this project!
```

### Update Agents

Keep your agents current:

```
You: "Update my agents"

Claude: ✓ Updated lando-agents (3 new commits)

        Checking deployed agents...
        • frontend (v1.1.0) → Update available to v1.2.0

        Would you like to update the frontend agent?
```

### Create Custom Agents

Build your own agents and manage them with CAMI:

```
You: "Help me create a custom database agent for PostgreSQL"

Claude: [Uses agent-architect to design the agent]

        I've created a PostgreSQL specialist agent. Let me save it to your
        custom agents directory.

        ✓ Created ~/cami/sources/my-agents/postgres-dba.md
        ✓ Agent discovered in my-agents source

        Would you like to deploy it to this project?
```

---

## CLI for Power Users

CAMI's CLI provides scripting and automation capabilities:

```bash
# Agent management
cami list                           # List available agents
cami deploy frontend backend ~/app  # Deploy agents to project
cami scan ~/app                     # Check deployed agents
cami update-docs ~/app              # Update CLAUDE.md

# Source management
cami source list                    # List agent sources
cami source add <git-url>           # Add new source
cami source update                  # Pull latest from all sources
cami source status                  # Check git status

# Location tracking
cami locations add my-app ~/app     # Track project location
cami locations list                 # List tracked projects
```

Use in CI/CD, deployment scripts, or automation workflows.

---

## Why Developers Love CAMI

### "Single Source of Truth"

> "Before CAMI, I had 5 copies of my frontend agent across different projects. Now I have one source and deploy it everywhere. Updates propagate in seconds."
>
> — **Sarah M., Full-Stack Developer**

### "Team Collaboration That Actually Works"

> "Our team shares a Git repo with our CAMI setup. New developers clone it, run the installer, and have access to all 15 of our company agents. No onboarding docs, no manual setup."
>
> — **Michael T., Engineering Manager**

### "Priority System Is Brilliant"

> "I can use the official agents as defaults but override with my customized versions for specific projects. No forking, no merge conflicts, just priority numbers."
>
> — **Alex K., DevOps Engineer**

### "Natural Language Is The Interface"

> "I don't remember CLI commands. I just ask Claude 'add the testing agent' and it happens. That's how all tools should work."
>
> — **Jordan P., Frontend Developer**

### "Version Tracking Saves Time"

> "CAMI tells me when agent updates are available and exactly which projects need updating. No more wondering if I'm using the latest version."
>
> — **Casey L., Tech Lead**

---

## Architecture

### Dual-Mode Binary

```bash
# MCP Server Mode (primary)
$ cami --mcp
# Runs as MCP server on stdio for Claude Code integration

# CLI Mode (secondary)
$ cami list
$ cami deploy frontend ~/my-app
# Direct commands for scripting and automation
```

### Workspace Structure

```
~/cami/                          # Your CAMI workspace
├── CLAUDE.md                    # CAMI documentation
├── README.md                    # Quick start guide
├── .mcp.json                    # MCP server configuration
├── .gitignore                   # Git ignore rules
├── config.yaml                  # CAMI configuration
├── .claude/
│   └── agents/                  # CAMI's own agents
├── sources/                     # Agent sources
│   ├── official-agents/        # Cloned from Git
│   ├── team-agents/            # Cloned from Git
│   └── my-agents/              # Your custom agents

/usr/local/bin/cami             # Binary on PATH
```

### Configuration Format

`~/cami/config.yaml`:

```yaml
version: "1"
agent_sources:
  - name: team-agents
    type: local
    path: ~/cami/sources/team-agents
    priority: 50
    git:
      enabled: true
      remote: git@github.com:yourorg/team-agents.git

  - name: my-agents
    type: local
    path: ~/cami/sources/my-agents
    priority: 10
    git:
      enabled: false

deploy_locations:
  - name: my-project
    path: ~/projects/my-project
```

### Agent File Format

Agents are markdown files with YAML frontmatter:

```markdown
---
name: frontend
version: "1.1.0"
description: Use this agent when building user interfaces...
---

# Frontend Development Specialist

You are a specialized frontend development expert focusing on React,
modern CSS, accessibility, and performance optimization...
```

---

## MCP Tools Reference

CAMI provides 13 MCP tools that enable Claude Code to manage agents through natural language:

**Project Creation**
- `create_project` - Create new project with agents and documentation

**Agent Management**
- `list_agents` - List all available agents from sources
- `deploy_agents` - Deploy agents to `.claude/agents/`
- `scan_deployed_agents` - Check deployed agents and versions
- `update_claude_md` - Update CLAUDE.md with agent info

**Source Management**
- `list_sources` - List configured agent sources
- `add_source` - Add source by cloning Git repository
- `update_source` - Pull latest from Git sources
- `source_status` - Check Git status of sources

**Location Management**
- `add_location` - Register project directory
- `list_locations` - List tracked projects
- `remove_location` - Unregister project

**Onboarding**
- `onboard` - Get personalized setup guidance

Claude Code uses these tools automatically when you have natural conversations about agent management.

---

## Use Cases

### Open Source Project Maintainers

Provide specialized agents for contributors:

```bash
# In your project's .github/ or docs/
git clone https://github.com/yourproject/agents.git

# Contributors add your agent source
cami source add https://github.com/yourproject/agents.git

# Deploy project-specific agents
cami deploy yourproject-contributing yourproject-testing ~/project
```

Contributors get context-aware AI assistance specific to your project's standards.

### Development Teams

Share company-standard agents across the team:

```bash
# Create team agent repository
~/company-agents/
├── .camiignore
├── frontend.md
├── backend.md
├── devops.md
└── security.md

# Team members add the source
cami source add git@github.com:company/agents.git

# Auto-deploy to all projects
cami deploy frontend backend devops ~/projects/*
```

Ensure consistent AI assistance aligned with company practices.

### Solo Developers

Organize personal agents across projects:

```bash
# Custom agents for different stacks
~/cami/sources/my-agents/
├── nextjs-specialist.md
├── python-django.md
├── rust-backend.md
└── ios-swiftui.md

# Deploy relevant specialists per project
cami deploy nextjs-specialist ~/web-app
cami deploy python-django ~/api-server
cami deploy rust-backend ~/microservices
```

Context-switch faster with project-appropriate AI specialists.

### CI/CD Automation

Automate agent deployment in build pipelines:

```bash
# .github/workflows/deploy.yml
- name: Deploy CAMI agents
  run: |
    cami source add git@github.com:company/agents.git
    cami deploy frontend backend qa $GITHUB_WORKSPACE
    cami update-docs $GITHUB_WORKSPACE
```

Ensure every deployment has up-to-date agents and documentation.

---

## Roadmap

### v0.4.0 (Planned)

- **Agent Classification System** - Categorize agents by domain, complexity, and use case
- **Remote Agent Sources** - HTTP endpoints and direct Git URL support
- **Enhanced Agent Discovery** - Search, filter, and recommendation engine
- **Team Collaboration Features** - Shared agent libraries and team workflows

### Future Considerations

- **Agent Analytics** - Track which agents are most used
- **Agent Templates** - Scaffolding for creating new agents
- **Web UI** - Browser-based agent management dashboard
- **Agent Marketplace** - Public registry of community agents

---

## Contributing

CAMI is open source and welcomes contributions!

**Development Setup:**

```bash
git clone https://github.com/lando-labs/cami.git
cd cami
claude  # Opens in Claude Code with dev mode enabled
```

**Build & Test:**

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linters
make install    # Install locally for testing
```

**Project Structure:**

```
cami/
├── cmd/cami/          # Binary entry point
├── internal/
│   ├── agent/         # Agent loading and parsing
│   ├── config/        # Configuration management
│   ├── deploy/        # Agent deployment
│   ├── mcp/           # MCP server implementation
│   └── cli/           # CLI commands
└── install/
    ├── templates/     # User workspace templates
    └── install.sh     # Installation script
```

See [CONTRIBUTING.md](https://github.com/lando-labs/cami/blob/main/CONTRIBUTING.md) for detailed guidelines.

---

## Documentation

- **[GitHub Repository](https://github.com/lando-labs/cami)** - Source code and issues
- **[README](https://github.com/lando-labs/cami/blob/main/README.md)** - Quick start guide
- **[CLAUDE.md](https://github.com/lando-labs/cami/blob/main/CLAUDE.md)** - Complete MCP tool documentation
- **[Releases](https://github.com/lando-labs/cami/releases)** - Download binaries

---

## Get Started Now

Stop copy-pasting agent files. Start managing them like infrastructure.

[Download CAMI v0.3.0](#installation){: .btn .btn-primary .btn-lg} [View on GitHub](https://github.com/lando-labs/cami){: .btn .btn-outline-secondary .btn-lg}

**Quick Start:**

```bash
# Download and extract for your platform
tar -xzf cami-0.3.0-darwin-arm64.tar.gz

# Run installer
cd cami-0.3.0-darwin-arm64
./install.sh

# Open CAMI workspace
cd ~/cami
claude

# Ask Claude to help you get started
"Help me get started with CAMI"
```

Within minutes, you'll have a centralized agent management system with access to professional agent libraries.

---

## License

MIT License - See [LICENSE](https://github.com/lando-labs/cami/blob/main/LICENSE) for details.

---

## Support

- **GitHub Issues** - [Bug reports and feature requests](https://github.com/lando-labs/cami/issues)
- **Discussions** - [Questions and community support](https://github.com/lando-labs/cami/discussions)
- **Twitter** - [@lando_labs](https://twitter.com/lando_labs) for updates

Built with ❤️ by [Lando Labs](https://github.com/lando-labs)
