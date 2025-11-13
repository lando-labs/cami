# Contributing to CAMI

Thank you for your interest in contributing to CAMI (Claude Agent Management Interface)! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- A GitHub account

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cami.git
   cd cami
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/lando/cami.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the project:
   ```bash
   go build -o cami ./cmd/cami
   ```

6. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```

## Development Workflow

### Branching Strategy

- `main` - stable release branch
- `feature/*` - new features
- `fix/*` - bug fixes
- `docs/*` - documentation updates

### Making Changes

1. Create a new branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes in small, logical commits

3. Keep your branch up to date with upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

4. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

## Testing

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run tests with race detection:
```bash
go test -race ./...
```

### Test Coverage Requirements

- **Critical packages** (`internal/agent`, `internal/deploy`): 90-95% coverage
- **Core packages** (`internal/config`, `internal/discovery`): 85%+ coverage
- **Overall project**: 80%+ coverage

### Writing Tests

- Use table-driven tests for testing multiple cases
- Use `t.TempDir()` for test isolation
- Test both success and error paths
- Include edge cases and boundary conditions
- Use descriptive test names that explain what is being tested

Example:
```go
func TestMyFeature(t *testing.T) {
	t.Run("successful case", func(t *testing.T) {
		// Test implementation
	})

	t.Run("error case", func(t *testing.T) {
		// Test implementation
	})
}
```

## Code Style

### Go Standards

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Run `go vet` to catch common issues
- Use meaningful variable and function names
- Write clear comments for exported functions and types

### Running Linters

```bash
# Format code
gofmt -s -w .

# Run go vet
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks
- `perf`: Performance improvements

Examples:
```
feat(agent): add support for remote agent sources

fix(deploy): handle missing .claude/agents directory

docs: update installation instructions

test: add tests for discovery package
```

## Submitting Changes

### Pull Request Process

1. Ensure all tests pass and coverage requirements are met
2. Update documentation if needed
3. Add a clear description of your changes in the PR
4. Link any related issues
5. Request review from maintainers

### Pull Request Template

```markdown
## Description
Brief description of what this PR does

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] Coverage requirements met

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings introduced
```

## Reporting Bugs

### Before Reporting

- Check existing issues to avoid duplicates
- Verify the bug exists in the latest version
- Collect relevant information (OS, Go version, etc.)

### Bug Report Template

```markdown
**Describe the bug**
A clear and concise description of the bug

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen

**Environment**
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.24]
- CAMI version: [e.g., 0.1.0]

**Additional context**
Any other relevant information
```

## Requesting Features

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of the problem

**Describe the solution you'd like**
A clear description of what you want to happen

**Describe alternatives you've considered**
Alternative solutions or features you've considered

**Additional context**
Any other context or screenshots
```

## Development Resources

### Project Structure

```
cami/
├── cmd/cami/          # CLI entry point
├── internal/
│   ├── agent/        # Agent loading and management
│   ├── cli/          # CLI commands
│   ├── config/       # Configuration management
│   ├── deploy/       # Agent deployment
│   ├── discovery/    # Project discovery
│   ├── docs/         # Documentation generation
│   └── tui/          # Terminal UI
├── vc-agents/        # Version-controlled agent sources
└── .claude/          # Claude Code configuration
```

### Key Packages

- **agent**: Core agent data structures and loading
- **cli**: Command-line interface implementation (uses Cobra)
- **config**: YAML configuration management
- **deploy**: Agent deployment to projects
- **discovery**: Bulk project scanning and status tracking
- **tui**: Interactive terminal UI (uses Bubble Tea)

### External Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/stretchr/testify` - Testing utilities

## Questions?

If you have questions about contributing, feel free to:
- Open a discussion on GitHub
- Ask in an issue
- Reach out to maintainers

Thank you for contributing to CAMI!
