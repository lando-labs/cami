# Contributing to CAMI

Thank you for your interest in contributing to CAMI! This document provides guidelines and information for contributors.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/lando-labs/cami/issues)
2. If not, create a new issue using the bug report template
3. Include:
   - CAMI version (`cami --version`)
   - Operating system and version
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs or error messages

### Suggesting Features

1. Check existing [Issues](https://github.com/lando-labs/cami/issues) and [Discussions](https://github.com/lando-labs/cami/discussions)
2. Open a new issue with the feature request template
3. Describe the use case and why it would benefit users

### Pull Requests

1. Fork the repository
2. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes following our coding standards
4. Write or update tests as needed
5. Run tests and linting:
   ```bash
   make test
   make lint
   ```
6. Commit with clear messages (see commit conventions below)
7. Push and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.21+
- Make
- Git

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/cami.git
cd cami

# Build
make build

# Run tests
make test

# Run linter
make lint

# Local development (uses go run)
go run ./cmd/cami --version
```

### Project Structure

```
cami/
├── cmd/cami/          # Main entry point
├── internal/          # Internal packages
│   ├── agent/         # Agent loading and parsing
│   ├── config/        # Configuration management
│   ├── deploy/        # Agent deployment
│   ├── mcp/           # MCP server implementation
│   └── ...
├── install/           # Installation scripts and templates
└── .claude/agents/    # Bundled agents
```

## Coding Standards

### Go Code

- Follow standard Go conventions and idioms
- Run `gofmt -s` before committing
- Keep functions focused and reasonably sized
- Add comments for exported functions
- Handle errors explicitly

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting (no code change)
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:
```
feat(mcp): add new deploy_agents tool
fix(config): handle missing config file gracefully
docs: update installation instructions
```

### Agent Contributions

If contributing new agents to the bundled set:

1. Follow the agent frontmatter format:
   ```yaml
   ---
   name: agent-name
   version: "1.0.0"
   description: Clear description of when to use this agent
   class: workflow-specialist|technology-implementer|strategic-planner
   specialty: domain-specialty
   model: haiku|sonnet|opus
   ---
   ```

2. Use the three-phase methodology appropriate for the agent class
3. Include clear boundaries and limitations
4. Add self-verification checklist

## Testing

- Write unit tests for new functionality
- Ensure existing tests pass
- Test on multiple platforms if possible (macOS, Linux)

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/agent/...

# Run with coverage
go test -cover ./...
```

## Questions?

- Open a [Discussion](https://github.com/lando-labs/cami/discussions) for questions
- Join our community channels (coming soon)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
