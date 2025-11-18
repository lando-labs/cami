package normalize

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeSource(t *testing.T) {
	t.Run("compliant source with no issues", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create valid agent files
		createTestAgent(t, tmpDir, "agent1.md", "agent1", "1.0.0", "Description for agent1")
		createTestAgent(t, tmpDir, "agent2.md", "agent2", "1.1.0", "Description for agent2")

		// Create .camiignore
		camiIgnorePath := filepath.Join(tmpDir, ".camiignore")
		require.NoError(t, os.WriteFile(camiIgnorePath, []byte("*.draft.md\n"), 0644))

		analysis, err := AnalyzeSource("test-source", tmpDir)

		require.NoError(t, err)
		assert.Equal(t, "test-source", analysis.SourceName)
		assert.Equal(t, tmpDir, analysis.Path)
		assert.True(t, analysis.IsCompliant)
		assert.Equal(t, 2, analysis.AgentCount)
		assert.Empty(t, analysis.Issues)
		assert.False(t, analysis.MissingCAMIIgnore)
	})

	t.Run("source with missing versions", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agent without version
		createTestAgent(t, tmpDir, "agent1.md", "agent1", "", "Description")

		analysis, err := AnalyzeSource("test-source", tmpDir)

		require.NoError(t, err)
		assert.False(t, analysis.IsCompliant)
		assert.Len(t, analysis.Issues, 1)
		assert.Equal(t, "agent1.md", analysis.Issues[0].AgentFile)
		assert.Contains(t, analysis.Issues[0].Problems, "missing version")
	})

	t.Run("source with missing descriptions", func(t *testing.T) {
		tmpDir := t.TempDir()

		createTestAgent(t, tmpDir, "agent1.md", "agent1", "1.0.0", "")

		analysis, err := AnalyzeSource("test-source", tmpDir)

		require.NoError(t, err)
		assert.False(t, analysis.IsCompliant)
		assert.Len(t, analysis.Issues, 1)
		assert.Contains(t, analysis.Issues[0].Problems, "missing description")
	})

	t.Run("source missing .camiignore", func(t *testing.T) {
		tmpDir := t.TempDir()

		createTestAgent(t, tmpDir, "agent1.md", "agent1", "1.0.0", "Description")

		analysis, err := AnalyzeSource("test-source", tmpDir)

		require.NoError(t, err)
		assert.False(t, analysis.IsCompliant)
		assert.True(t, analysis.MissingCAMIIgnore)
	})

	t.Run("error on non-existent path", func(t *testing.T) {
		_, err := AnalyzeSource("test-source", "/nonexistent/path")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestNormalizeSource(t *testing.T) {
	t.Run("add versions to agents", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agents without versions
		createTestAgent(t, tmpDir, "agent1.md", "agent1", "", "Description 1")
		createTestAgent(t, tmpDir, "agent2.md", "agent2", "", "Description 2")

		options := SourceNormalizationOptions{
			AddVersions: true,
		}

		result, err := NormalizeSource("test-source", tmpDir, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, 2, result.AgentsUpdated)
		assert.NotEmpty(t, result.BackupPath)
		assert.Len(t, result.Changes, 2)

		// Verify agents were updated
		content, err := os.ReadFile(filepath.Join(tmpDir, "agent1.md"))
		require.NoError(t, err)
		assert.Contains(t, string(content), "version: 1.0.0")
	})

	t.Run("add description placeholders", func(t *testing.T) {
		tmpDir := t.TempDir()

		createTestAgent(t, tmpDir, "agent1.md", "agent1", "1.0.0", "")

		options := SourceNormalizationOptions{
			AddDescriptions: true,
		}

		result, err := NormalizeSource("test-source", tmpDir, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, 1, result.AgentsUpdated)

		// Verify description was added
		content, err := os.ReadFile(filepath.Join(tmpDir, "agent1.md"))
		require.NoError(t, err)
		assert.Contains(t, string(content), "description: Description for agent1 agent")
	})

	t.Run("create .camiignore file", func(t *testing.T) {
		tmpDir := t.TempDir()

		createTestAgent(t, tmpDir, "agent1.md", "agent1", "1.0.0", "Description")

		options := SourceNormalizationOptions{
			CreateCAMIIgnore: true,
		}

		result, err := NormalizeSource("test-source", tmpDir, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Contains(t, result.Changes, "Created .camiignore file")

		// Verify .camiignore exists
		camiIgnorePath := filepath.Join(tmpDir, ".camiignore")
		assert.FileExists(t, camiIgnorePath)

		content, err := os.ReadFile(camiIgnorePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "*.draft.md")
	})

	t.Run("backup is created", func(t *testing.T) {
		tmpDir := t.TempDir()

		createTestAgent(t, tmpDir, "agent1.md", "agent1", "", "Description")

		options := SourceNormalizationOptions{
			AddVersions: true,
		}

		result, err := NormalizeSource("test-source", tmpDir, options)

		require.NoError(t, err)
		assert.NotEmpty(t, result.BackupPath)
		assert.DirExists(t, result.BackupPath)

		// Verify backup contains original files
		backupAgent := filepath.Join(result.BackupPath, "agent1.md")
		assert.FileExists(t, backupAgent)
	})
}

func TestAnalyzeProject(t *testing.T) {
	t.Run("non-cami project without agents directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		analysis, err := AnalyzeProject(tmpDir, []config.AgentSource{})

		require.NoError(t, err)
		assert.Equal(t, manifest.StateNonCAMI, analysis.State)
		assert.False(t, analysis.HasAgentsDir)
		assert.False(t, analysis.HasManifest)
		assert.Equal(t, 0, analysis.AgentCount)
	})

	t.Run("cami-aware project with agents but no manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agents directory
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))

		// Add agents
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")
		createTestAgent(t, agentsDir, "agent2.md", "agent2", "1.1.0", "Description")

		analysis, err := AnalyzeProject(tmpDir, []config.AgentSource{})

		require.NoError(t, err)
		assert.Equal(t, manifest.StateCAMIAware, analysis.State)
		assert.True(t, analysis.HasAgentsDir)
		assert.False(t, analysis.HasManifest)
		assert.Equal(t, 2, analysis.AgentCount)
		assert.Len(t, analysis.Agents, 2)
		assert.True(t, analysis.Recommendations.MinimalRequired)
	})

	t.Run("cami-native project with manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agents directory
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		// Create manifest
		projectManifest := &manifest.ProjectManifest{
			Version:      "2",
			State:        manifest.StateCAMINative,
			NormalizedAt: mustParseTime("2025-01-01T00:00:00Z"),
			Agents:       []manifest.DeployedAgent{},
		}
		require.NoError(t, manifest.WriteProjectManifest(tmpDir, projectManifest))

		analysis, err := AnalyzeProject(tmpDir, []config.AgentSource{})

		require.NoError(t, err)
		assert.Equal(t, manifest.StateCAMINative, analysis.State)
		assert.True(t, analysis.HasManifest)
		assert.Equal(t, 1, analysis.AgentCount)
	})

	t.Run("detect agents matching sources", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := t.TempDir()

		// Create source agent
		createTestAgent(t, sourceDir, "agent1.md", "agent1", "2.0.0", "Description")

		// Create deployed agent with older version
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		sources := []config.AgentSource{
			{
				Name:     "test-source",
				Path:     sourceDir,
				Priority: 50,
			},
		}

		analysis, err := AnalyzeProject(tmpDir, sources)

		require.NoError(t, err)
		assert.Len(t, analysis.Agents, 1)
		assert.Equal(t, "available-sources", analysis.Agents[0].MatchesSource)
		assert.True(t, analysis.Agents[0].NeedsUpgrade)
	})

	t.Run("error on non-existent path", func(t *testing.T) {
		_, err := AnalyzeProject("/nonexistent/path", []config.AgentSource{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestNormalizeProject(t *testing.T) {
	// Override HOME for central manifest
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpHome)

	t.Run("minimal normalization creates manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agents directory
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		options := ProjectNormalizationOptions{
			Level: LevelMinimal,
		}

		result, err := NormalizeProject(tmpDir, options, []config.AgentSource{})

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, manifest.StateCAMIAware, result.StateBefore)
		assert.Equal(t, manifest.StateCAMINative, result.StateAfter)
		assert.NotEmpty(t, result.BackupPath)
		assert.True(t, result.UndoAvailable)

		// Verify manifest was created
		manifestPath := filepath.Join(tmpDir, manifest.ProjectManifestFilename)
		assert.FileExists(t, manifestPath)

		// Read and verify manifest
		projectManifest, err := manifest.ReadProjectManifest(tmpDir)
		require.NoError(t, err)
		assert.Equal(t, manifest.StateCAMINative, projectManifest.State)
		assert.Len(t, projectManifest.Agents, 1)
		assert.Equal(t, "agent1", projectManifest.Agents[0].Name)
	})

	t.Run("standard normalization links sources", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := t.TempDir()

		// Create source agent
		createTestAgent(t, sourceDir, "agent1.md", "agent1", "1.0.0", "Description")

		// Create deployed agent
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		sources := []config.AgentSource{
			{
				Name:     "test-source",
				Path:     sourceDir,
				Priority: 50,
			},
		}

		options := ProjectNormalizationOptions{
			Level: LevelStandard,
		}

		result, err := NormalizeProject(tmpDir, options, sources)

		require.NoError(t, err)
		assert.True(t, result.Success)

		// Verify manifest has source links
		projectManifest, err := manifest.ReadProjectManifest(tmpDir)
		require.NoError(t, err)
		assert.Len(t, projectManifest.Agents, 1)
		assert.Equal(t, "test-source", projectManifest.Agents[0].Source)
		assert.NotEmpty(t, projectManifest.Agents[0].SourcePath)
		assert.Equal(t, 50, projectManifest.Agents[0].Priority)
	})

	t.Run("central manifest is updated", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create agents directory
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		options := ProjectNormalizationOptions{
			Level: LevelMinimal,
		}

		_, err := NormalizeProject(tmpDir, options, []config.AgentSource{})
		require.NoError(t, err)

		// Read central manifest
		centralManifest, err := manifest.ReadCentralManifest()
		require.NoError(t, err)

		// Verify project is in central manifest
		absPath, _ := filepath.Abs(tmpDir)
		deployment, exists := centralManifest.Deployments[absPath]
		assert.True(t, exists)
		assert.Equal(t, manifest.StateCAMINative, deployment.State)
		assert.Len(t, deployment.Agents, 1)
	})

	t.Run("backup is created before normalization", func(t *testing.T) {
		tmpDir := t.TempDir()

		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		options := ProjectNormalizationOptions{
			Level: LevelMinimal,
		}

		result, err := NormalizeProject(tmpDir, options, []config.AgentSource{})

		require.NoError(t, err)
		assert.NotEmpty(t, result.BackupPath)
		assert.DirExists(t, result.BackupPath)
	})

	t.Run("full normalization returns not implemented error", func(t *testing.T) {
		tmpDir := t.TempDir()

		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		createTestAgent(t, agentsDir, "agent1.md", "agent1", "1.0.0", "Description")

		options := ProjectNormalizationOptions{
			Level: LevelFull,
		}

		_, err := NormalizeProject(tmpDir, options, []config.AgentSource{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")
	})
}

// Helper functions

func createTestAgent(t *testing.T, dir, filename, name, version, description string) {
	t.Helper()

	content := "---\n"
	if name != "" {
		content += "name: " + name + "\n"
	}
	if version != "" {
		content += "version: " + version + "\n"
	}
	if description != "" {
		content += "description: " + description + "\n"
	}
	content += "---\n\n# Agent Content\n\nThis is the agent content.\n"

	filePath := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
