package deploy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lando/cami/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test agent
func createTestAgent(name, version string) *agent.Agent {
	return &agent.Agent{
		Name:        name,
		Version:     version,
		Description: "Test agent for " + name,
		FilePath:    "/fake/path/" + name + ".md",
		Content:     "# Test Content\n\nThis is test content.",
	}
}

func TestValidateTargetPath(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := ValidateTargetPath(tmpDir)
		assert.NoError(t, err)
	})

	t.Run("path does not exist", func(t *testing.T) {
		err := ValidateTargetPath("/nonexistent/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("path is not a directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "file.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("test"), 0644))

		err := ValidateTargetPath(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})
}

func TestDeployAgent(t *testing.T) {
	t.Run("successful deployment", func(t *testing.T) {
		tmpDir := t.TempDir()
		ag := createTestAgent("test-agent", "1.0.0")

		result, err := DeployAgent(ag, tmpDir, false)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.False(t, result.Conflict)
		assert.Equal(t, "Deployed successfully", result.Message)
		assert.Equal(t, ag, result.Agent)

		// Verify file was created
		agentPath := filepath.Join(tmpDir, ".claude", "agents", "test-agent.md")
		assert.FileExists(t, agentPath)

		// Verify content
		content, err := os.ReadFile(agentPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "name: test-agent")
		assert.Contains(t, string(content), "version: 1.0.0")
		assert.Contains(t, string(content), "# Test Content")
	})

	t.Run("creates .claude/agents directory if missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		ag := createTestAgent("new-agent", "1.0.0")

		// Ensure .claude/agents doesn't exist yet
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		_, err := os.Stat(agentsDir)
		assert.True(t, os.IsNotExist(err))

		result, err := DeployAgent(ag, tmpDir, false)

		require.NoError(t, err)
		assert.True(t, result.Success)

		// Verify directory was created
		assert.DirExists(t, agentsDir)
	})

	t.Run("conflict when file exists and overwrite is false", func(t *testing.T) {
		tmpDir := t.TempDir()
		ag := createTestAgent("existing-agent", "1.0.0")

		// Deploy once
		result1, err := DeployAgent(ag, tmpDir, false)
		require.NoError(t, err)
		require.True(t, result1.Success)

		// Try to deploy again without overwrite
		result2, err := DeployAgent(ag, tmpDir, false)

		require.NoError(t, err)
		assert.False(t, result2.Success)
		assert.True(t, result2.Conflict)
		assert.Equal(t, "File already exists", result2.Message)
	})

	t.Run("overwrite existing file when overwrite is true", func(t *testing.T) {
		tmpDir := t.TempDir()
		ag1 := createTestAgent("overwrite-test", "1.0.0")
		ag2 := createTestAgent("overwrite-test", "2.0.0")

		// Deploy first version
		result1, err := DeployAgent(ag1, tmpDir, false)
		require.NoError(t, err)
		require.True(t, result1.Success)

		// Deploy second version with overwrite
		result2, err := DeployAgent(ag2, tmpDir, true)

		require.NoError(t, err)
		assert.True(t, result2.Success)
		assert.False(t, result2.Conflict)

		// Verify new version is deployed
		agentPath := filepath.Join(tmpDir, ".claude", "agents", "overwrite-test.md")
		content, err := os.ReadFile(agentPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "version: 2.0.0")
	})
}

func TestDeployAgents(t *testing.T) {
	t.Run("deploy multiple agents successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		agents := []*agent.Agent{
			createTestAgent("agent1", "1.0.0"),
			createTestAgent("agent2", "2.0.0"),
			createTestAgent("agent3", "3.0.0"),
		}

		results, err := DeployAgents(agents, tmpDir, false)

		require.NoError(t, err)
		assert.Len(t, results, 3)

		// All should succeed
		for i, result := range results {
			assert.True(t, result.Success, "Agent %d should succeed", i)
			assert.Equal(t, agents[i], result.Agent)
		}

		// Verify all files exist
		assert.FileExists(t, filepath.Join(tmpDir, ".claude", "agents", "agent1.md"))
		assert.FileExists(t, filepath.Join(tmpDir, ".claude", "agents", "agent2.md"))
		assert.FileExists(t, filepath.Join(tmpDir, ".claude", "agents", "agent3.md"))
	})

	t.Run("some agents succeed, some conflict", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Pre-deploy one agent
		existingAgent := createTestAgent("existing", "1.0.0")
		_, err := DeployAgent(existingAgent, tmpDir, false)
		require.NoError(t, err)

		// Try to deploy multiple including the existing one
		agents := []*agent.Agent{
			createTestAgent("new1", "1.0.0"),
			createTestAgent("existing", "2.0.0"), // This will conflict
			createTestAgent("new2", "1.0.0"),
		}

		results, err := DeployAgents(agents, tmpDir, false)

		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Check results
		assert.True(t, results[0].Success)
		assert.False(t, results[1].Success)
		assert.True(t, results[1].Conflict)
		assert.True(t, results[2].Success)
	})

	t.Run("deploy empty list", func(t *testing.T) {
		tmpDir := t.TempDir()

		results, err := DeployAgents([]*agent.Agent{}, tmpDir, false)

		require.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("deploy with overwrite", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Pre-deploy agents
		agent1 := createTestAgent("agent1", "1.0.0")
		agent2 := createTestAgent("agent2", "1.0.0")
		_, err := DeployAgents([]*agent.Agent{agent1, agent2}, tmpDir, false)
		require.NoError(t, err)

		// Deploy new versions with overwrite
		newAgents := []*agent.Agent{
			createTestAgent("agent1", "2.0.0"),
			createTestAgent("agent2", "2.0.0"),
		}

		results, err := DeployAgents(newAgents, tmpDir, true)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)

		// Verify versions were updated
		content1, _ := os.ReadFile(filepath.Join(tmpDir, ".claude", "agents", "agent1.md"))
		assert.Contains(t, string(content1), "version: 2.0.0")

		content2, _ := os.ReadFile(filepath.Join(tmpDir, ".claude", "agents", "agent2.md"))
		assert.Contains(t, string(content2), "version: 2.0.0")
	})
}

func TestCheckConflicts(t *testing.T) {
	t.Run("no conflicts", func(t *testing.T) {
		tmpDir := t.TempDir()
		agents := []*agent.Agent{
			createTestAgent("agent1", "1.0.0"),
			createTestAgent("agent2", "1.0.0"),
		}

		conflicts := CheckConflicts(agents, tmpDir)

		assert.Len(t, conflicts, 0)
	})

	t.Run("some conflicts", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Pre-deploy some agents
		existingAgents := []*agent.Agent{
			createTestAgent("existing1", "1.0.0"),
			createTestAgent("existing2", "1.0.0"),
		}
		_, err := DeployAgents(existingAgents, tmpDir, false)
		require.NoError(t, err)

		// Check conflicts with mixed list
		agents := []*agent.Agent{
			createTestAgent("existing1", "2.0.0"),
			createTestAgent("new", "1.0.0"),
			createTestAgent("existing2", "2.0.0"),
		}

		conflicts := CheckConflicts(agents, tmpDir)

		assert.Len(t, conflicts, 2)
		assert.True(t, conflicts["existing1"])
		assert.True(t, conflicts["existing2"])
		assert.False(t, conflicts["new"])
	})

	t.Run("all conflicts", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Pre-deploy agents
		agents := []*agent.Agent{
			createTestAgent("agent1", "1.0.0"),
			createTestAgent("agent2", "1.0.0"),
		}
		_, err := DeployAgents(agents, tmpDir, false)
		require.NoError(t, err)

		// Check same agents
		conflicts := CheckConflicts(agents, tmpDir)

		assert.Len(t, conflicts, 2)
		assert.True(t, conflicts["agent1"])
		assert.True(t, conflicts["agent2"])
	})

	t.Run("agents directory doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		agents := []*agent.Agent{
			createTestAgent("agent1", "1.0.0"),
		}

		conflicts := CheckConflicts(agents, tmpDir)

		assert.Len(t, conflicts, 0)
	})
}
