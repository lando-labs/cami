package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper to create a temporary agent file
func createTestAgent(t *testing.T, dir, name, version, description, content string) string {
	t.Helper()
	
	agentContent := "---\n"
	agentContent += "name: " + name + "\n"
	agentContent += "version: " + version + "\n"
	agentContent += "description: " + description + "\n"
	agentContent += "---\n\n"
	agentContent += content

	filePath := filepath.Join(dir, name+".md")
	err := os.WriteFile(filePath, []byte(agentContent), 0644)
	require.NoError(t, err)
	
	return filePath
}

func TestLoadAgent(t *testing.T) {
	t.Run("valid agent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := createTestAgent(t, tmpDir, "test-agent", "1.0.0", "Test agent", "# Test Content")

		agent, err := LoadAgent(filePath)
		
		require.NoError(t, err)
		assert.Equal(t, "test-agent", agent.Name)
		assert.Equal(t, "1.0.0", agent.Version)
		assert.Equal(t, "Test agent", agent.Description)
		assert.Contains(t, agent.Content, "# Test Content")
		assert.Equal(t, filePath, agent.FilePath)
	})

	t.Run("missing frontmatter delimiter", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "bad.md")
		
		err := os.WriteFile(filePath, []byte("no frontmatter here"), 0644)
		require.NoError(t, err)

		_, err = LoadAgent(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing frontmatter delimiter")
	})

	t.Run("empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.md")
		
		err := os.WriteFile(filePath, []byte(""), 0644)
		require.NoError(t, err)

		_, err = LoadAgent(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty file")
	})

	t.Run("invalid yaml frontmatter", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "bad-yaml.md")
		
		content := "---\ninvalid: yaml: content:\n---\n"
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)

		_, err = LoadAgent(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse frontmatter")
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := LoadAgent("/nonexistent/file.md")
		assert.Error(t, err)
	})
}

func TestLoadAgents(t *testing.T) {
	t.Run("load multiple agents", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		createTestAgent(t, tmpDir, "agent1", "1.0.0", "First agent", "Content 1")
		createTestAgent(t, tmpDir, "agent2", "2.0.0", "Second agent", "Content 2")
		createTestAgent(t, tmpDir, "agent3", "3.0.0", "Third agent", "Content 3")

		agents, err := LoadAgents(tmpDir)
		
		require.NoError(t, err)
		assert.Len(t, agents, 3)
		
		// Verify all agents loaded
		names := make(map[string]bool)
		for _, agent := range agents {
			names[agent.Name] = true
		}
		assert.True(t, names["agent1"])
		assert.True(t, names["agent2"])
		assert.True(t, names["agent3"])
	})

	t.Run("load agents from nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		coreDir := filepath.Join(tmpDir, "core")
		specializedDir := filepath.Join(tmpDir, "specialized")
		require.NoError(t, os.MkdirAll(coreDir, 0755))
		require.NoError(t, os.MkdirAll(specializedDir, 0755))
		
		createTestAgent(t, coreDir, "core-agent", "1.0.0", "Core", "Core content")
		createTestAgent(t, specializedDir, "specialized-agent", "1.0.0", "Specialized", "Specialized content")

		agents, err := LoadAgents(tmpDir)
		
		require.NoError(t, err)
		assert.Len(t, agents, 2)
		
		// Verify categories are extracted
		for _, agent := range agents {
			if agent.Name == "core-agent" {
				assert.Equal(t, "core", agent.Category)
			} else if agent.Name == "specialized-agent" {
				assert.Equal(t, "specialized", agent.Category)
			}
		}
	})

	t.Run("skip non-markdown files", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		createTestAgent(t, tmpDir, "valid", "1.0.0", "Valid", "Content")
		
		// Create non-markdown files
		os.WriteFile(filepath.Join(tmpDir, "README.txt"), []byte("readme"), 0644)
		os.WriteFile(filepath.Join(tmpDir, "data.json"), []byte("{}"), 0644)

		agents, err := LoadAgents(tmpDir)
		
		require.NoError(t, err)
		assert.Len(t, agents, 1)
		assert.Equal(t, "valid", agents[0].Name)
	})

	t.Run("empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		agents, err := LoadAgents(tmpDir)
		
		require.NoError(t, err)
		assert.Len(t, agents, 0)
	})

	t.Run("directory does not exist", func(t *testing.T) {
		_, err := LoadAgents("/nonexistent/directory")
		assert.Error(t, err)
	})
}

func TestLoadAgentsFromSources(t *testing.T) {
	t.Run("load from multiple sources", func(t *testing.T) {
		source1 := t.TempDir()
		source2 := t.TempDir()
		
		createTestAgent(t, source1, "agent1", "1.0.0", "From source 1", "Content 1")
		createTestAgent(t, source2, "agent2", "2.0.0", "From source 2", "Content 2")

		sources := []AgentSource{
			{Path: source1, Priority: 100},
			{Path: source2, Priority: 150},
		}

		agents, err := LoadAgentsFromSources(sources)
		
		require.NoError(t, err)
		assert.Len(t, agents, 2)
	})

	t.Run("priority deduplication", func(t *testing.T) {
		source1 := t.TempDir()
		source2 := t.TempDir()
		
		// Same agent name, different versions and priorities
		createTestAgent(t, source1, "frontend", "1.0.0", "Low priority", "Old version")
		createTestAgent(t, source2, "frontend", "2.0.0", "High priority", "New version")

		sources := []AgentSource{
			{Path: source1, Priority: 100},
			{Path: source2, Priority: 200}, // Higher priority
		}

		agents, err := LoadAgentsFromSources(sources)
		
		require.NoError(t, err)
		assert.Len(t, agents, 1)
		
		// Should get the higher priority version
		assert.Equal(t, "frontend", agents[0].Name)
		assert.Equal(t, "2.0.0", agents[0].Version)
		assert.Equal(t, "High priority", agents[0].Description)
	})

	t.Run("handle source loading errors gracefully", func(t *testing.T) {
		validSource := t.TempDir()
		createTestAgent(t, validSource, "valid", "1.0.0", "Valid agent", "Content")

		sources := []AgentSource{
			{Path: "/nonexistent/path", Priority: 100},
			{Path: validSource, Priority: 200},
		}

		agents, err := LoadAgentsFromSources(sources)
		
		// Should still succeed with valid source
		require.NoError(t, err)
		assert.Len(t, agents, 1)
		assert.Equal(t, "valid", agents[0].Name)
	})

	t.Run("empty sources list", func(t *testing.T) {
		sources := []AgentSource{}

		agents, err := LoadAgentsFromSources(sources)
		
		require.NoError(t, err)
		assert.Len(t, agents, 0)
	})
}

func TestAgentFullContent(t *testing.T) {
	agent := &Agent{
		Name:        "test",
		Version:     "1.0.0",
		Description: "Test agent",
		Content:     "# Agent Content\n\nThis is the content.",
	}

	fullContent := agent.FullContent()

	assert.Contains(t, fullContent, "---")
	assert.Contains(t, fullContent, "name: test")
	assert.Contains(t, fullContent, "version: 1.0.0")
	assert.Contains(t, fullContent, "description: Test agent")
	assert.Contains(t, fullContent, "# Agent Content")
}

func TestAgentFileName(t *testing.T) {
	agent := &Agent{
		FilePath: "/path/to/agents/my-agent.md",
	}

	fileName := agent.FileName()
	assert.Equal(t, "my-agent.md", fileName)
}
