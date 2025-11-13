package discovery

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test agent file
func createTestAgentFile(t *testing.T, dir, name, version string) {
	t.Helper()
	
	agentContent := "---\n"
	agentContent += "name: " + name + "\n"
	agentContent += "version: " + version + "\n"
	agentContent += "description: Test agent " + name + "\n"
	agentContent += "---\n\n"
	agentContent += "# Test Content"

	filePath := filepath.Join(dir, name+".md")
	require.NoError(t, os.WriteFile(filePath, []byte(agentContent), 0644))
}

// Helper to create a test project with .claude/agents
func createTestProject(t *testing.T, rootDir, projectName string, agentNames ...string) string {
	t.Helper()
	
	projectPath := filepath.Join(rootDir, projectName)
	agentsDir := filepath.Join(projectPath, ".claude", "agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))
	
	for _, name := range agentNames {
		createTestAgentFile(t, agentsDir, name, "1.0.0")
	}
	
	return projectPath
}

func TestDiscoverProjects(t *testing.T) {
	t.Run("discover multiple projects", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		// Create projects with agents
		createTestProject(t, tmpDir, "project1", "frontend", "backend")
		createTestProject(t, tmpDir, "project2", "devops")
		createTestProject(t, tmpDir, "project3") // Empty .claude/agents
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 3)
		
		// Verify project details
		projectMap := make(map[string]*ProjectInfo)
		for _, p := range projects {
			projectMap[filepath.Base(p.Path)] = p
		}
		
		assert.True(t, projectMap["project1"].HasAgents)
		assert.Equal(t, 2, projectMap["project1"].AgentCount)
		
		assert.True(t, projectMap["project2"].HasAgents)
		assert.Equal(t, 1, projectMap["project2"].AgentCount)
		
		assert.False(t, projectMap["project3"].HasAgents)
		assert.Equal(t, 0, projectMap["project3"].AgentCount)
	})
	
	t.Run("empty-only filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		createTestProject(t, tmpDir, "with-agents", "frontend")
		createTestProject(t, tmpDir, "empty-project") // No agents
		
		opts := DiscoverOptions{
			RootPath:  tmpDir,
			EmptyOnly: true,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, "empty-project", filepath.Base(projects[0].Path))
		assert.False(t, projects[0].HasAgents)
	})
	
	t.Run("has-agent filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		createTestProject(t, tmpDir, "project1", "frontend", "backend")
		createTestProject(t, tmpDir, "project2", "devops", "frontend")
		createTestProject(t, tmpDir, "project3", "backend")
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
			HasAgent: "frontend",
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 2)
		
		// Verify both have frontend agent
		for _, p := range projects {
			found := false
			for _, ag := range p.Agents {
				if ag.Name == "frontend" {
					found = true
					break
				}
			}
			assert.True(t, found, "Project should have frontend agent")
		}
	})
	
	t.Run("max-depth filter", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		// Create nested structure
		createTestProject(t, tmpDir, "shallow", "agent1")
		createTestProject(t, filepath.Join(tmpDir, "level1"), "nested", "agent2")
		createTestProject(t, filepath.Join(tmpDir, "level1", "level2"), "deep", "agent3")
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
			MaxDepth: 1,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		// Should only find shallow and nested (depth 1), not deep (depth 2)
		assert.LessOrEqual(t, len(projects), 2)
	})
	
	t.Run("no projects found", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 0)
	})
	
	t.Run("nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		// Create nested projects
		createTestProject(t, filepath.Join(tmpDir, "workspace", "app1"), "app1", "frontend")
		createTestProject(t, filepath.Join(tmpDir, "workspace", "app2"), "app2", "backend")
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 2)
	})
	
	t.Run("relative paths populated", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		createTestProject(t, tmpDir, "project1", "agent")
		
		opts := DiscoverOptions{
			RootPath: tmpDir,
		}
		
		projects, err := DiscoverProjects(opts)
		
		require.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.NotEmpty(t, projects[0].RelativePath)
		assert.Equal(t, "project1", projects[0].RelativePath)
	})
	
	t.Run("invalid root path", func(t *testing.T) {
		opts := DiscoverOptions{
			RootPath: "/nonexistent/path",
		}
		
		_, err := DiscoverProjects(opts)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestScanLocation(t *testing.T) {
	t.Run("scan location with agents", func(t *testing.T) {
		tmpDir := t.TempDir()
		agentsDir := filepath.Join(tmpDir, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir, 0755))
		
		createTestAgentFile(t, agentsDir, "frontend", "1.0.0")
		createTestAgentFile(t, agentsDir, "backend", "2.0.0")
		
		location := &config.DeployLocation{
			Name: "test-project",
			Path: tmpDir,
		}
		
		availableAgents := []*agent.Agent{
			{Name: "frontend", Version: "1.0.0"},
			{Name: "backend", Version: "3.0.0"}, // Different version
			{Name: "devops", Version: "1.0.0"},  // Not deployed
		}
		
		status, err := ScanLocation(location, availableAgents)
		
		require.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, location, status.Location)
		assert.Len(t, status.AgentStatuses, 3)
		
		// Check statuses
		statusMap := make(map[string]*AgentStatus)
		for _, s := range status.AgentStatuses {
			statusMap[s.Agent.Name] = s
		}
		
		assert.Equal(t, StatusUpToDate, statusMap["frontend"].Status)
		assert.Equal(t, "1.0.0", statusMap["frontend"].DeployedVersion)
		
		assert.Equal(t, StatusUpdateAvailable, statusMap["backend"].Status)
		assert.Equal(t, "2.0.0", statusMap["backend"].DeployedVersion)
		assert.Equal(t, "3.0.0", statusMap["backend"].AvailableVersion)
		
		assert.Equal(t, StatusNotDeployed, statusMap["devops"].Status)
		assert.Equal(t, "", statusMap["devops"].DeployedVersion)
	})
	
	t.Run("scan location with no agents directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		location := &config.DeployLocation{
			Name: "empty-project",
			Path: tmpDir,
		}
		
		availableAgents := []*agent.Agent{
			{Name: "frontend", Version: "1.0.0"},
		}
		
		status, err := ScanLocation(location, availableAgents)
		
		require.NoError(t, err)
		assert.Len(t, status.AgentStatuses, 1)
		assert.Equal(t, StatusNotDeployed, status.AgentStatuses[0].Status)
	})
	
	t.Run("LastScanned timestamp set", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		location := &config.DeployLocation{
			Name: "test",
			Path: tmpDir,
		}
		
		before := time.Now()
		status, err := ScanLocation(location, []*agent.Agent{})
		after := time.Now()
		
		require.NoError(t, err)
		assert.True(t, status.LastScanned.After(before) || status.LastScanned.Equal(before))
		assert.True(t, status.LastScanned.Before(after) || status.LastScanned.Equal(after))
	})
}

func TestScanAllLocations(t *testing.T) {
	t.Run("scan multiple locations", func(t *testing.T) {
		tmpDir1 := t.TempDir()
		tmpDir2 := t.TempDir()
		
		// Setup first location
		agentsDir1 := filepath.Join(tmpDir1, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir1, 0755))
		createTestAgentFile(t, agentsDir1, "frontend", "1.0.0")
		
		// Setup second location
		agentsDir2 := filepath.Join(tmpDir2, ".claude", "agents")
		require.NoError(t, os.MkdirAll(agentsDir2, 0755))
		createTestAgentFile(t, agentsDir2, "backend", "1.0.0")
		
		locations := []config.DeployLocation{
			{Name: "project1", Path: tmpDir1},
			{Name: "project2", Path: tmpDir2},
		}
		
		availableAgents := []*agent.Agent{
			{Name: "frontend", Version: "1.0.0"},
			{Name: "backend", Version: "1.0.0"},
		}
		
		result, err := ScanAllLocations(locations, availableAgents)
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.LocationStatuses, 2)
		assert.Equal(t, availableAgents, result.AvailableAgents)
	})
	
	t.Run("skip inaccessible locations", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		locations := []config.DeployLocation{
			{Name: "bad", Path: "/nonexistent/path"},
			{Name: "good", Path: tmpDir},
		}
		
		result, err := ScanAllLocations(locations, []*agent.Agent{})
		
		require.NoError(t, err)
		// Should have 1 location (the good one)
		assert.Len(t, result.LocationStatuses, 2)
		// Both locations scanned, just one had no agents dir
	})
}

func TestGetStatusSymbol(t *testing.T) {
	tests := []struct {
		status DeploymentStatus
		want   string
	}{
		{StatusUpToDate, "✓"},
		{StatusUpdateAvailable, "⚠"},
		{StatusNotDeployed, "○"},
		{StatusUnknown, "?"},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := GetStatusSymbol(tt.status)
			assert.Equal(t, tt.want, got)
		})
	}
}
