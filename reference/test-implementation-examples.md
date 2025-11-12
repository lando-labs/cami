# CAMI Test Implementation Examples

**Purpose:** Practical examples to kickstart Phase 0 testing
**Target:** Week 1-2 critical path tests
**Related:** [open-source-strategy.md#1-testing-architecture](./open-source-strategy.md#1-testing-architecture)

## Setup

### 1. Install testify

```bash
go get github.com/stretchr/testify
```

Update `go.mod`:
```go
require (
    github.com/stretchr/testify v1.9.0
)
```

### 2. Create test data structure

```bash
mkdir -p internal/agent/testdata
mkdir -p internal/deploy/testdata
mkdir -p internal/docs/testdata
```

### 3. Set up CI/CD

Create `.github/workflows/test.yml`:
```yaml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt
```

## Example 1: internal/agent Package Tests

### Test Data Files

**internal/agent/testdata/valid-agent.md**
```markdown
---
name: test-agent
version: "1.0.0"
description: A test agent for unit tests
---

This is the agent content.

## Role

Test agent role definition.

## Examples

Example usage here.
```

**internal/agent/testdata/invalid-frontmatter.md**
```markdown
---
name: test-agent
version: "1.0.0"
description: A test agent
invalid_yaml: [unclosed bracket
---

Content here.
```

**internal/agent/testdata/missing-version.md**
```markdown
---
name: test-agent
description: Missing version field
---

Content here.
```

**internal/agent/testdata/no-frontmatter.md**
```markdown
# Just a markdown file

No YAML frontmatter at all.
```

**internal/agent/testdata/empty-content.md**
```markdown
---
name: test-agent
version: "1.0.0"
description: Valid frontmatter but no content
---
```

### Unit Tests

**internal/agent/agent_test.go**
```go
package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lando/cami/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAgent_ValidFrontmatter(t *testing.T) {
	// Arrange
	agentPath := filepath.Join("testdata", "valid-agent.md")

	// Act
	ag, err := agent.LoadAgent(agentPath)

	// Assert
	require.NoError(t, err, "LoadAgent should not return an error for valid agent")
	assert.Equal(t, "test-agent", ag.Name, "Name should match frontmatter")
	assert.Equal(t, "1.0.0", ag.Version, "Version should match frontmatter")
	assert.Equal(t, "A test agent for unit tests", ag.Description, "Description should match")
	assert.Contains(t, ag.Content, "This is the agent content", "Content should be present")
	assert.Contains(t, ag.Content, "## Role", "Content should include role section")
}

func TestLoadAgent_InvalidFrontmatter(t *testing.T) {
	agentPath := filepath.Join("testdata", "invalid-frontmatter.md")
	_, err := agent.LoadAgent(agentPath)
	assert.Error(t, err, "LoadAgent should return error for invalid YAML")
	assert.Contains(t, err.Error(), "failed to parse frontmatter", "Error should mention frontmatter")
}

func TestLoadAgent_MissingVersion(t *testing.T) {
	agentPath := filepath.Join("testdata", "missing-version.md")
	ag, err := agent.LoadAgent(agentPath)

	// Current implementation: no validation, so this passes
	// Future: might want to validate required fields
	require.NoError(t, err)
	assert.Empty(t, ag.Version, "Version should be empty if not provided")
}

func TestLoadAgent_NoFrontmatter(t *testing.T) {
	agentPath := filepath.Join("testdata", "no-frontmatter.md")
	_, err := agent.LoadAgent(agentPath)
	assert.Error(t, err, "LoadAgent should fail without frontmatter")
	assert.Contains(t, err.Error(), "missing frontmatter", "Error should mention missing frontmatter")
}

func TestLoadAgent_EmptyContent(t *testing.T) {
	agentPath := filepath.Join("testdata", "empty-content.md")
	ag, err := agent.LoadAgent(agentPath)
	require.NoError(t, err, "Valid frontmatter with empty content should succeed")
	assert.Empty(t, ag.Content, "Content should be empty")
}

func TestLoadAgent_FileNotFound(t *testing.T) {
	agentPath := filepath.Join("testdata", "nonexistent.md")
	_, err := agent.LoadAgent(agentPath)
	assert.Error(t, err, "LoadAgent should fail for nonexistent file")
	assert.Contains(t, err.Error(), "failed to open file")
}

func TestAgent_FullContent(t *testing.T) {
	// Arrange
	ag := &agent.Agent{
		Name:        "test-agent",
		Version:     "1.0.0",
		Description: "Test description",
		Content:     "Agent content here",
	}

	// Act
	fullContent := ag.FullContent()

	// Assert
	assert.Contains(t, fullContent, "---", "Should include frontmatter delimiters")
	assert.Contains(t, fullContent, "name: test-agent", "Should include name")
	assert.Contains(t, fullContent, "version: 1.0.0", "Should include version")
	assert.Contains(t, fullContent, "description: Test description", "Should include description")
	assert.Contains(t, fullContent, "Agent content here", "Should include content")
}

func TestAgent_FileName(t *testing.T) {
	ag := &agent.Agent{
		FilePath: "/path/to/agents/architect.md",
	}
	assert.Equal(t, "architect.md", ag.FileName(), "FileName should return basename")
}

func TestLoadAgents_ValidDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create category directories
	coreDir := filepath.Join(tmpDir, "core")
	specializedDir := filepath.Join(tmpDir, "specialized")
	require.NoError(t, os.MkdirAll(coreDir, 0755))
	require.NoError(t, os.MkdirAll(specializedDir, 0755))

	// Create test agents
	agent1 := `---
name: architect
version: "1.0.0"
description: Architecture specialist
---
Content for architect`

	agent2 := `---
name: backend
version: "1.1.0"
description: Backend specialist
---
Content for backend`

	require.NoError(t, os.WriteFile(filepath.Join(coreDir, "architect.md"), []byte(agent1), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(specializedDir, "backend.md"), []byte(agent2), 0644))

	// Act
	agents, err := agent.LoadAgents(tmpDir)

	// Assert
	require.NoError(t, err)
	assert.Len(t, agents, 2, "Should load 2 agents")

	// Find architect agent
	var architectAgent *agent.Agent
	for _, ag := range agents {
		if ag.Name == "architect" {
			architectAgent = ag
			break
		}
	}
	require.NotNil(t, architectAgent, "Should find architect agent")
	assert.Equal(t, "core", architectAgent.Category, "Category should be extracted from directory")
}

func TestLoadAgents_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	agents, err := agent.LoadAgents(tmpDir)
	require.NoError(t, err, "Should not error on empty directory")
	assert.Empty(t, agents, "Should return empty slice")
}

func TestLoadAgents_NonexistentDirectory(t *testing.T) {
	_, err := agent.LoadAgents("/nonexistent/path")
	assert.Error(t, err, "Should error for nonexistent directory")
}
```

**Run tests:**
```bash
cd internal/agent
go test -v
go test -cover
```

## Example 2: internal/deploy Package Tests

### Test Data

**internal/deploy/testdata/setup.go** (helper for tests)
```go
package testdata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lando/cami/internal/agent"
	"github.com/stretchr/testify/require"
)

// CreateTestAgent creates a test agent for deployment tests
func CreateTestAgent(t *testing.T, name, version, description string) *agent.Agent {
	t.Helper()

	return &agent.Agent{
		Name:        name,
		Version:     version,
		Description: description,
		FilePath:    filepath.Join("/tmp", name+".md"),
		Content:     "Test agent content for " + name,
	}
}

// CreateProjectWithAgents sets up a test project with deployed agents
func CreateProjectWithAgents(t *testing.T, agents ...*agent.Agent) string {
	t.Helper()

	projectDir := t.TempDir()
	agentsDir := filepath.Join(projectDir, ".claude", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	for _, ag := range agents {
		content := ag.FullContent()
		path := filepath.Join(agentsDir, ag.Name+".md")
		require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	}

	return projectDir
}
```

### Unit Tests

**internal/deploy/deploy_test.go**
```go
package deploy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/deploy"
	"github.com/lando/cami/internal/deploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployAgent_Success(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	ag := testdata.CreateTestAgent(t, "architect", "1.0.0", "Test agent")

	// Act
	result, err := deploy.DeployAgent(ag, tmpDir, false)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "Deployed successfully", result.Message)
	assert.False(t, result.Conflict)

	// Verify file exists
	expectedPath := filepath.Join(tmpDir, ".claude", "agents", "architect.md")
	assert.FileExists(t, expectedPath)

	// Verify content
	content, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: architect")
	assert.Contains(t, string(content), "version: 1.0.0")
}

func TestDeployAgent_Conflict_NoOverwrite(t *testing.T) {
	// Arrange
	ag := testdata.CreateTestAgent(t, "backend", "1.1.0", "Backend agent")
	projectDir := testdata.CreateProjectWithAgents(t, ag) // Pre-deploy agent

	// Try to deploy same agent without overwrite
	result, err := deploy.DeployAgent(ag, projectDir, false)

	// Assert
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.True(t, result.Conflict)
	assert.Equal(t, "File already exists", result.Message)
}

func TestDeployAgent_Conflict_WithOverwrite(t *testing.T) {
	// Arrange
	oldAgent := testdata.CreateTestAgent(t, "frontend", "1.0.0", "Old version")
	projectDir := testdata.CreateProjectWithAgents(t, oldAgent)

	newAgent := testdata.CreateTestAgent(t, "frontend", "1.1.0", "New version")

	// Act - deploy with overwrite
	result, err := deploy.DeployAgent(newAgent, projectDir, true)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.False(t, result.Conflict)

	// Verify updated content
	agentPath := filepath.Join(projectDir, ".claude", "agents", "frontend.md")
	content, err := os.ReadFile(agentPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "version: 1.1.0", "Should have new version")
	assert.Contains(t, string(content), "New version", "Should have new description")
}

func TestDeployAgents_Multiple(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	agents := []*agent.Agent{
		testdata.CreateTestAgent(t, "architect", "1.0.0", "Architect"),
		testdata.CreateTestAgent(t, "backend", "1.1.0", "Backend"),
		testdata.CreateTestAgent(t, "frontend", "1.2.0", "Frontend"),
	}

	// Act
	results, err := deploy.DeployAgents(agents, tmpDir, false)

	// Assert
	require.NoError(t, err)
	assert.Len(t, results, 3)

	for _, result := range results {
		assert.True(t, result.Success, "All deployments should succeed")
	}

	// Verify all files exist
	agentsDir := filepath.Join(tmpDir, ".claude", "agents")
	assert.FileExists(t, filepath.Join(agentsDir, "architect.md"))
	assert.FileExists(t, filepath.Join(agentsDir, "backend.md"))
	assert.FileExists(t, filepath.Join(agentsDir, "frontend.md"))
}

func TestValidateTargetPath_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	err := deploy.ValidateTargetPath(tmpDir)
	assert.NoError(t, err)
}

func TestValidateTargetPath_NotExists(t *testing.T) {
	err := deploy.ValidateTargetPath("/nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestValidateTargetPath_NotDirectory(t *testing.T) {
	// Create a file, not directory
	tmpFile := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))

	err := deploy.ValidateTargetPath(tmpFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}

func TestCheckConflicts(t *testing.T) {
	// Arrange
	existingAgent := testdata.CreateTestAgent(t, "architect", "1.0.0", "Existing")
	projectDir := testdata.CreateProjectWithAgents(t, existingAgent)

	agentsToCheck := []*agent.Agent{
		testdata.CreateTestAgent(t, "architect", "1.1.0", "Will conflict"),
		testdata.CreateTestAgent(t, "backend", "1.0.0", "No conflict"),
	}

	// Act
	conflicts := deploy.CheckConflicts(agentsToCheck, projectDir)

	// Assert
	assert.True(t, conflicts["architect"], "Should detect conflict for architect")
	assert.False(t, conflicts["backend"], "Should not detect conflict for backend")
}
```

**Run tests:**
```bash
cd internal/deploy
go test -v -cover
```

## Example 3: Integration Test

**internal/deploy/integration_test.go**
```go
//go:build integration

package deploy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/deploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullDeploymentWorkflow(t *testing.T) {
	// This test uses real agent files from vc-agents/
	// Skip if vc-agents doesn't exist
	vcAgentsDir := "../../vc-agents"
	if _, err := os.Stat(vcAgentsDir); os.IsNotExist(err) {
		t.Skip("Skipping: vc-agents directory not found")
	}

	// Load real agents
	agents, err := agent.LoadAgents(vcAgentsDir)
	require.NoError(t, err)
	require.NotEmpty(t, agents, "Should load agents from vc-agents")

	// Find architect agent
	var architectAgent *agent.Agent
	for _, ag := range agents {
		if ag.Name == "architect" {
			architectAgent = ag
			break
		}
	}
	require.NotNil(t, architectAgent, "Should find architect agent")

	// Create temp project directory
	projectDir := t.TempDir()

	// Deploy agent
	result, err := deploy.DeployAgent(architectAgent, projectDir, false)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// Verify deployment
	deployedPath := filepath.Join(projectDir, ".claude", "agents", "architect.md")
	assert.FileExists(t, deployedPath)

	// Load deployed agent and verify
	deployedAgent, err := agent.LoadAgent(deployedPath)
	require.NoError(t, err)
	assert.Equal(t, architectAgent.Name, deployedAgent.Name)
	assert.Equal(t, architectAgent.Version, deployedAgent.Version)
	assert.Equal(t, architectAgent.Description, deployedAgent.Description)

	// Test conflict detection
	result2, err := deploy.DeployAgent(architectAgent, projectDir, false)
	require.NoError(t, err)
	assert.False(t, result2.Success)
	assert.True(t, result2.Conflict)

	// Test overwrite
	result3, err := deploy.DeployAgent(architectAgent, projectDir, true)
	require.NoError(t, err)
	assert.True(t, result3.Success)
}
```

**Run integration tests:**
```bash
go test -v -tags=integration ./internal/deploy/
```

## Example 4: E2E CLI Test

**test/e2e/cli_test.go**
```go
//go:build e2e

package e2e_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI_List_JSON(t *testing.T) {
	// Build CLI first
	buildCmd := exec.Command("go", "build", "-o", "cami-test", "../../cmd/cami/main.go")
	require.NoError(t, buildCmd.Run(), "Failed to build CLI")
	defer os.Remove("cami-test")

	// Run list command
	cmd := exec.Command("./cami-test", "list", "--output", "json")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI command failed: %s", string(output))

	// Parse JSON
	var result struct {
		Count  int `json:"count"`
		Agents []struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Category    string `json:"category"`
		} `json:"agents"`
	}

	err = json.Unmarshal(output, &result)
	require.NoError(t, err, "Failed to parse JSON output")

	// Verify
	assert.Greater(t, result.Count, 0, "Should list agents")
	assert.Len(t, result.Agents, result.Count, "Count should match array length")
}

func TestCLI_Deploy_Full_Workflow(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "cami-test", "../../cmd/cami/main.go")
	require.NoError(t, buildCmd.Run())
	defer os.Remove("cami-test")

	// Create temp project
	projectDir := t.TempDir()

	// Deploy architect agent
	cmd := exec.Command("./cami-test", "deploy",
		"--agents", "architect",
		"--location", projectDir)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Deploy failed: %s", string(output))

	// Verify file exists
	agentPath := filepath.Join(projectDir, ".claude", "agents", "architect.md")
	assert.FileExists(t, agentPath, "Agent should be deployed")

	// Try to deploy again (should conflict)
	cmd2 := exec.Command("./cami-test", "deploy",
		"--agents", "architect",
		"--location", projectDir)

	output2, err2 := cmd2.CombinedOutput()
	// This should fail due to conflict
	assert.Error(t, err2, "Should fail on conflict")
	assert.Contains(t, string(output2), "conflict", "Output should mention conflict")

	// Deploy with overwrite (should succeed)
	cmd3 := exec.Command("./cami-test", "deploy",
		"--agents", "architect",
		"--location", projectDir,
		"--overwrite")

	output3, err3 := cmd3.CombinedOutput()
	require.NoError(t, err3, "Deploy with overwrite failed: %s", string(output3))
}

func TestCLI_Scan(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "cami-test", "../../cmd/cami/main.go")
	require.NoError(t, buildCmd.Run())
	defer os.Remove("cami-test")

	// Create project with deployed agent
	projectDir := t.TempDir()
	cmd := exec.Command("./cami-test", "deploy",
		"--agents", "architect",
		"--location", projectDir)
	require.NoError(t, cmd.Run())

	// Scan project
	scanCmd := exec.Command("./cami-test", "scan",
		"--location", projectDir,
		"--output", "json")

	output, err := scanCmd.CombinedOutput()
	require.NoError(t, err, "Scan failed: %s", string(output))

	// Parse output
	var result struct {
		Statuses []struct {
			Name             string `json:"name"`
			DeployedVersion  string `json:"deployed_version"`
			AvailableVersion string `json:"available_version"`
			Status           string `json:"status"`
		} `json:"statuses"`
	}

	err = json.Unmarshal(output, &result)
	require.NoError(t, err)

	// Find architect in results
	var architectStatus *struct {
		Name             string `json:"name"`
		DeployedVersion  string `json:"deployed_version"`
		AvailableVersion string `json:"available_version"`
		Status           string `json:"status"`
	}

	for _, status := range result.Statuses {
		if status.Name == "architect" {
			architectStatus = &status
			break
		}
	}

	require.NotNil(t, architectStatus, "Should find architect in scan results")
	assert.NotEmpty(t, architectStatus.DeployedVersion, "Should have deployed version")
	assert.Equal(t, "up-to-date", architectStatus.Status, "Should be up-to-date")
}
```

**Run E2E tests:**
```bash
go test -v -tags=e2e ./test/e2e/
```

## Makefile for Testing

**Makefile** (root of project)
```makefile
.PHONY: test test-unit test-integration test-e2e test-all test-cover

# Run unit tests only (default)
test:
	go test -v ./...

# Run unit tests with coverage
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	@go tool cover -func=coverage.out | grep total

# Run integration tests
test-integration:
	go test -v -tags=integration ./...

# Run E2E tests
test-e2e:
	go build -o cami cmd/cami/main.go
	go test -v -tags=e2e ./test/e2e/
	rm -f cami

# Run all tests
test-all:
	@echo "Running unit tests..."
	@make test
	@echo "\nRunning integration tests..."
	@make test-integration
	@echo "\nRunning E2E tests..."
	@make test-e2e

# Clean test artifacts
clean-test:
	rm -f coverage.out coverage.html
	rm -f cami cami-test
```

## Usage

```bash
# Day-to-day development
make test

# Check coverage
make test-cover

# Before committing
make test-all

# CI/CD
make test-cover
```

## Coverage Target Verification

**scripts/check-coverage.sh**
```bash
#!/bin/bash
set -e

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Extract overall coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "Current coverage: ${COVERAGE}%"

# Check if meets target (80%)
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "❌ Coverage below 80% target"
    exit 1
else
    echo "✅ Coverage meets 80% target"
fi
```

Make executable:
```bash
chmod +x scripts/check-coverage.sh
```

Add to CI:
```yaml
- name: Check coverage target
  run: ./scripts/check-coverage.sh
```

## Next Steps

1. **Create test data** (Week 1 Day 1)
   - Set up `testdata/` directories
   - Create sample agent files
   - Create test fixtures

2. **Write agent tests** (Week 1 Day 2-3)
   - Copy examples above
   - Run tests: `go test -v ./internal/agent/`
   - Fix any failures

3. **Write deploy tests** (Week 1 Day 4-5)
   - Copy examples above
   - Run tests: `go test -v ./internal/deploy/`
   - Achieve 90% coverage

4. **Set up CI/CD** (Week 2 Day 1)
   - Copy GitHub Actions workflow
   - Push and verify tests run
   - Set up Codecov

5. **Continue with other packages** (Week 2-3)
   - Follow same pattern for `internal/docs`, `internal/config`, etc.

---

**These examples provide a solid foundation for Phase 0 testing. Copy, adapt, and iterate!**
