# CAMI Architecture

Technical documentation for CAMI's architecture and design decisions.

## Overview

CAMI is a Go-based terminal application using the Bubbletea TUI framework for managing and deploying Claude Code agents across projects.

## Architecture Principles

1. **Clean Separation**: Business logic separated from UI concerns
2. **Terminal-First**: Built for keyboard-driven workflows
3. **Minimal Dependencies**: Only essential libraries
4. **Graceful Degradation**: Clear error messages and fallbacks

## Project Structure

```
cami/
├── cmd/cami/          # Application entry point
│   └── main.go        # CLI initialization, flags, and TUI startup
├── internal/          # Internal packages (not importable)
│   ├── agent/         # Agent parsing and management
│   ├── config/        # Configuration persistence
│   ├── deploy/        # Deployment engine
│   └── tui/           # Bubbletea UI components
└── vc-agents/         # Version-controlled agent definitions
```

## Core Components

### 1. Agent Parser (`internal/agent`)

**Responsibility**: Parse and load agent definitions from markdown files with YAML frontmatter.

**Key Types**:
- `Agent`: Represents a parsed agent with metadata and content
- `Metadata`: YAML frontmatter structure

**Key Functions**:
- `LoadAgents(dir)`: Loads all agents from a directory
- `LoadAgent(path)`: Parses a single agent file
- `FullContent()`: Reconstructs agent with frontmatter

**Design Decisions**:
- YAML frontmatter for structured metadata
- Markdown for agent content (human-readable)
- Tolerant parsing (logs warnings, continues on errors)
- File-based versioning (each agent has a version field)

### 2. Configuration (`internal/config`)

**Responsibility**: Manage deployment locations and persist configuration.

**Key Types**:
- `Config`: Application configuration
- `DeployLocation`: Deployment target with name and path

**Key Functions**:
- `Load()`: Read config from disk
- `Save()`: Write config to disk
- `AddDeployLocation()`: Add new location with validation
- `RemoveDeployLocation()`: Remove location by index

**Design Decisions**:
- JSON format for simplicity and human-readability
- Stored in home directory (`~/.cami.json`)
- Validation on add (path must exist, no duplicates)

### 3. Deployment Engine (`internal/deploy`)

**Responsibility**: Deploy agents to target projects.

**Key Types**:
- `Result`: Deployment result with success/conflict/error status

**Key Functions**:
- `DeployAgent()`: Deploy single agent
- `DeployAgents()`: Deploy multiple agents
- `ValidateTargetPath()`: Ensure target is valid
- `CheckConflicts()`: Detect existing files

**Design Decisions**:
- Creates `.claude/agents/` directory structure
- Conflict detection before deployment
- Optional overwrite mode
- Individual results for batch operations
- Full content preservation (frontmatter + content)

### 4. TUI (`internal/tui`)

**Responsibility**: Terminal user interface using Bubbletea.

**Key Components**:
- `Model`: Application state
- `ViewState`: Current view enum (agent selection, locations, deployment, results)
- `keyMap`: Keyboard bindings

**Views**:
1. **Agent Selection**: Browse and select agents to deploy
2. **Location Management**: Add/remove deployment locations
3. **Deployment**: Choose target location
4. **Results**: Show deployment results

**Design Decisions**:
- State machine architecture (ViewState enum)
- Keyboard-first navigation (vim-style hjkl + arrows)
- Lipgloss for styling (colors, formatting)
- No mouse dependency (but could be added)
- Single file responsibility (model.go = logic, views.go = rendering)

## Data Flow

### Deployment Flow

```
User selects agents → User selects location → Deploy → Show results
       ↓                      ↓                  ↓           ↓
   Agent list          DeployLocation      DeployEngine   Results view
   (TUI state)         (from config)       (creates files) (TUI state)
```

### Configuration Flow

```
User adds location → Validate path → Save to disk → Update TUI
        ↓                  ↓              ↓             ↓
    Input fields       os.Stat()      JSON write    Refresh list
    (TUI state)        (validation)   (persist)     (TUI state)
```

## Key Design Patterns

### 1. Bubbletea Pattern (ELM Architecture)

```go
type Model struct { /* state */ }

func (m Model) Init() tea.Cmd { /* initialization */ }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* state updates */ }
func (m Model) View() string { /* rendering */ }
```

### 2. Result Pattern

Instead of error-only returns, deployment returns structured results:

```go
type Result struct {
    Agent    *agent.Agent
    Success  bool
    Message  string
    Conflict bool
}
```

This allows:
- Batch operations with partial failures
- Detailed status reporting
- Differentiation between errors and conflicts

### 3. View State Machine

```go
type ViewState int

const (
    ViewAgentSelection
    ViewLocationManagement
    ViewDeployment
    ViewResults
)
```

Clean separation of UI states with dedicated update/render functions.

## Technology Choices

### Why Go?
- Fast compilation and execution
- Single binary deployment (no runtime)
- Excellent terminal library ecosystem
- Strong standard library for file operations

### Why Bubbletea?
- Clean, functional TUI framework
- ELM architecture (predictable state management)
- Active maintenance and ecosystem
- Built-in terminal abstractions

### Why Lipgloss?
- Declarative styling
- Terminal-safe color handling
- Consistent with Bubbletea

### Why YAML frontmatter?
- Human-readable metadata
- Standard in documentation tools
- Easy to parse (gopkg.in/yaml.v3)
- Compatible with markdown

## Testing Strategy

Current test files (can be run manually):
- `test_agents.go`: Verify agent parsing
- `test_deploy.go`: Verify deployment engine
- `test_config.go`: Verify configuration management

Note: TUI testing requires a TTY and is best done interactively.

Future: Add unit tests for each package.

## Extension Points

### Adding New Agent Metadata

1. Update `Metadata` struct in `internal/agent/agent.go`
2. Update YAML frontmatter in agent files
3. Update display in `internal/tui/views.go`

### Adding New Views

1. Add to `ViewState` enum in `internal/tui/model.go`
2. Create `update{ViewName}()` function
3. Create `view{ViewName}()` function
4. Add state transitions in existing update functions

### Adding New Deployment Targets

Currently deploys to `.claude/agents/`. To support other targets:

1. Add target type to `DeployLocation` in `internal/config/config.go`
2. Update deployment logic in `internal/deploy/deploy.go`
3. Add UI selection in `internal/tui/views.go`

## Performance Considerations

- **Agent Loading**: O(n) on startup, cached in memory
- **File Operations**: Synchronous (acceptable for typical use)
- **Rendering**: Bubbletea handles efficient terminal updates
- **Memory**: All agents loaded in memory (fine for expected scale <100 agents)

## Security Considerations

- **Path Validation**: Validates all paths before operations
- **No Eval**: No dynamic code execution
- **File Permissions**: Creates files with 0644, dirs with 0755
- **No Network**: Pure local file operations

## Future Enhancements

### Milestone 2 (Planned)
- Claude.md section updates
- Agent orchestration guidance
- Deployment history
- Rollback capability

### Potential Features
- Agent dependency resolution
- Template-based agent generation
- Remote agent repositories
- Multi-project batch deployments
- Agent update notifications

## Contributing

When adding features:

1. **Maintain separation**: Keep business logic out of TUI
2. **Follow patterns**: Use existing Result/ViewState patterns
3. **Test manually**: Run test files for affected components
4. **Update docs**: Keep this file and README in sync
5. **Keyboard-first**: All features must be keyboard accessible

## Debugging

Run with verbose output:
```bash
# Build with debug info
go build -gcflags="-N -l" -o cami cmd/cami/main.go

# Add debug prints (use stderr, not stdout)
fmt.Fprintf(os.Stderr, "Debug: %+v\n", value)
```

For TUI debugging, redirect to a file:
```bash
./cami 2>debug.log
```

## Build & Release

```bash
# Development build
make build

# Production build with version
go build -ldflags="-s -w -X main.version=1.0.0" -o cami cmd/cami/main.go

# Cross-compile for other platforms
GOOS=linux GOARCH=amd64 go build -o cami-linux cmd/cami/main.go
GOOS=darwin GOARCH=arm64 go build -o cami-darwin-arm64 cmd/cami/main.go
```

## Dependencies

Core:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/bubbles` - TUI components (key bindings)
- `gopkg.in/yaml.v3` - YAML parsing

All dependencies are vendored in go.mod.

---

Last updated: 2025-10-04 (v0.1.0)
