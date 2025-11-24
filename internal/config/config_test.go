package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()

	require.NoError(t, err)
	assert.Contains(t, path, "cami-workspace")
	assert.Contains(t, path, "config.yaml")
}

func TestLoadAndSave(t *testing.T) {
	// Use a temporary directory for config
	tmpDir := t.TempDir()

	// Override config path for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	t.Run("save and load config", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			AgentSources: []AgentSource{
				{
					Name:     "test-source",
					Type:     "local",
					Path:     "/test/path",
					Priority: 100,
					Git: &GitConfig{
						Enabled: true,
						Remote:  "git@github.com:test/repo.git",
					},
				},
			},
			Locations: []DeployLocation{
				{
					Name: "test-project",
					Path: "/test/project",
				},
			},
		}

		// Create config directory
		require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "cami-workspace"), 0755))

		// Save
		err := cfg.Save()
		require.NoError(t, err)

		// Load
		loaded, err := Load()
		require.NoError(t, err)

		// Verify
		assert.Equal(t, "1", loaded.Version)
		assert.Len(t, loaded.AgentSources, 1)
		assert.Equal(t, "test-source", loaded.AgentSources[0].Name)
		assert.Equal(t, "local", loaded.AgentSources[0].Type)
		assert.Equal(t, "/test/path", loaded.AgentSources[0].Path)
		assert.Equal(t, 100, loaded.AgentSources[0].Priority)
		assert.NotNil(t, loaded.AgentSources[0].Git)
		assert.True(t, loaded.AgentSources[0].Git.Enabled)
		assert.Equal(t, "git@github.com:test/repo.git", loaded.AgentSources[0].Git.Remote)

		assert.Len(t, loaded.Locations, 1)
		assert.Equal(t, "test-project", loaded.Locations[0].Name)
		assert.Equal(t, "/test/project", loaded.Locations[0].Path)
	})

	t.Run("load non-existent config creates default", func(t *testing.T) {
		// Remove config file
		os.Remove(filepath.Join(tmpDir, "cami-workspace", "config.yaml"))

		cfg, err := Load()

		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "1", cfg.Version)
		assert.Len(t, cfg.AgentSources, 0)
		assert.Len(t, cfg.Locations, 0)
	})
}

func TestAddAgentSource(t *testing.T) {
	t.Run("add new source", func(t *testing.T) {
		cfg := &Config{
			Version:      "1",
			AgentSources: []AgentSource{},
		}

		source := AgentSource{
			Name:     "new-source",
			Type:     "local",
			Path:     "/new/path",
			Priority: 150,
		}

		err := cfg.AddAgentSource(source)

		require.NoError(t, err)
		assert.Len(t, cfg.AgentSources, 1)
		assert.Equal(t, "new-source", cfg.AgentSources[0].Name)
	})

	t.Run("error on duplicate name", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			AgentSources: []AgentSource{
				{Name: "existing", Path: "/path", Priority: 100},
			},
		}

		duplicate := AgentSource{
			Name:     "existing",
			Path:     "/other/path",
			Priority: 200,
		}

		err := cfg.AddAgentSource(duplicate)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		assert.Len(t, cfg.AgentSources, 1) // Should not be added
	})
}

func TestRemoveAgentSource(t *testing.T) {
	t.Run("remove existing source", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			AgentSources: []AgentSource{
				{Name: "source1", Path: "/path1", Priority: 100},
				{Name: "source2", Path: "/path2", Priority: 150},
				{Name: "source3", Path: "/path3", Priority: 200},
			},
		}

		err := cfg.RemoveAgentSource("source2")

		require.NoError(t, err)
		assert.Len(t, cfg.AgentSources, 2)
		assert.Equal(t, "source1", cfg.AgentSources[0].Name)
		assert.Equal(t, "source3", cfg.AgentSources[1].Name)
	})

	t.Run("error when source not found", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			AgentSources: []AgentSource{
				{Name: "source1", Path: "/path1", Priority: 100},
			},
		}

		err := cfg.RemoveAgentSource("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.Len(t, cfg.AgentSources, 1) // Should not be removed
	})
}

func TestGetAgentSource(t *testing.T) {
	cfg := &Config{
		Version: "1",
		AgentSources: []AgentSource{
			{Name: "source1", Path: "/path1", Priority: 100},
			{Name: "source2", Path: "/path2", Priority: 150},
		},
	}

	t.Run("get existing source", func(t *testing.T) {
		source, err := cfg.GetAgentSource("source2")

		require.NoError(t, err)
		assert.Equal(t, "source2", source.Name)
		assert.Equal(t, "/path2", source.Path)
		assert.Equal(t, 150, source.Priority)
	})

	t.Run("error when not found", func(t *testing.T) {
		_, err := cfg.GetAgentSource("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestAddDeployLocation(t *testing.T) {
	t.Run("add new location", func(t *testing.T) {
		tmpDir := t.TempDir() // Create real directory

		cfg := &Config{
			Version:   "1",
			Locations: []DeployLocation{},
		}

		err := cfg.AddDeployLocation("project1", tmpDir)

		require.NoError(t, err)
		assert.Len(t, cfg.Locations, 1)
		assert.Equal(t, "project1", cfg.Locations[0].Name)
		assert.Equal(t, tmpDir, cfg.Locations[0].Path)
	})

	t.Run("error on duplicate name", func(t *testing.T) {
		tmpDir := t.TempDir()

		cfg := &Config{
			Version: "1",
			Locations: []DeployLocation{
				{Name: "existing", Path: tmpDir},
			},
		}

		err := cfg.AddDeployLocation("existing", tmpDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		assert.Len(t, cfg.Locations, 1)
	})

	t.Run("error when path does not exist", func(t *testing.T) {
		cfg := &Config{
			Version:   "1",
			Locations: []DeployLocation{},
		}

		err := cfg.AddDeployLocation("bad", "/nonexistent/path")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestRemoveDeployLocationByName(t *testing.T) {
	t.Run("remove existing location", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			Locations: []DeployLocation{
				{Name: "loc1", Path: "/path1"},
				{Name: "loc2", Path: "/path2"},
				{Name: "loc3", Path: "/path3"},
			},
		}

		err := cfg.RemoveDeployLocationByName("loc2")

		require.NoError(t, err)
		assert.Len(t, cfg.Locations, 2)
		assert.Equal(t, "loc1", cfg.Locations[0].Name)
		assert.Equal(t, "loc3", cfg.Locations[1].Name)
	})

	t.Run("error when not found", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			Locations: []DeployLocation{
				{Name: "loc1", Path: "/path1"},
			},
		}

		err := cfg.RemoveDeployLocationByName("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.Len(t, cfg.Locations, 1)
	})
}

func TestRemoveDeployLocation(t *testing.T) {
	t.Run("remove by index", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			Locations: []DeployLocation{
				{Name: "loc1", Path: "/path1"},
				{Name: "loc2", Path: "/path2"},
			},
		}

		err := cfg.RemoveDeployLocation(0)

		require.NoError(t, err)
		assert.Len(t, cfg.Locations, 1)
		assert.Equal(t, "loc2", cfg.Locations[0].Name)
	})

	t.Run("error on invalid index", func(t *testing.T) {
		cfg := &Config{
			Version: "1",
			Locations: []DeployLocation{
				{Name: "loc1", Path: "/path1"},
			},
		}

		err := cfg.RemoveDeployLocation(5)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid index")
	})
}

func TestIsFreshInstall(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "fresh install - no sources",
			config: &Config{
				Version:            "1",
				InstallTimestamp:   time.Now(),
				SetupComplete:      false,
				AgentSources:       []AgentSource{},
				Locations:          []DeployLocation{},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: true,
		},
		{
			name: "fresh install - only my-agents with no git",
			config: &Config{
				Version:          "1",
				InstallTimestamp: time.Now(),
				SetupComplete:    false,
				AgentSources: []AgentSource{
					{
						Name:     "my-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/my-agents",
						Priority: 10,
						Git: &GitConfig{
							Enabled: false,
						},
					},
				},
				Locations:          []DeployLocation{},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: true,
		},
		{
			name: "not fresh - setup complete",
			config: &Config{
				Version:          "1",
				InstallTimestamp: time.Now(),
				SetupComplete:    true,
				AgentSources: []AgentSource{
					{
						Name:     "my-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/my-agents",
						Priority: 10,
					},
				},
				Locations:          []DeployLocation{},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: false,
		},
		{
			name: "not fresh - has real source",
			config: &Config{
				Version:          "1",
				InstallTimestamp: time.Now(),
				SetupComplete:    false,
				AgentSources: []AgentSource{
					{
						Name:     "my-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/my-agents",
						Priority: 10,
					},
					{
						Name:     "lando-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/lando-agents",
						Priority: 100,
						Git: &GitConfig{
							Enabled: true,
							Remote:  "git@github.com:lando-labs/lando-agents.git",
						},
					},
				},
				Locations:          []DeployLocation{},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: false,
		},
		{
			name: "not fresh - my-agents has git",
			config: &Config{
				Version:          "1",
				InstallTimestamp: time.Now(),
				SetupComplete:    false,
				AgentSources: []AgentSource{
					{
						Name:     "my-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/my-agents",
						Priority: 10,
						Git: &GitConfig{
							Enabled: true,
							Remote:  "git@github.com:company/agents.git",
						},
					},
				},
				Locations:          []DeployLocation{},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: false,
		},
		{
			name: "not fresh - has locations",
			config: &Config{
				Version:          "1",
				InstallTimestamp: time.Now(),
				SetupComplete:    false,
				AgentSources: []AgentSource{
					{
						Name:     "my-agents",
						Type:     "local",
						Path:     "/home/user/cami-workspace/sources/my-agents",
						Priority: 10,
					},
				},
				Locations: []DeployLocation{
					{
						Name: "my-project",
						Path: "/home/user/projects/my-project",
					},
				},
				DefaultProjectsDir: "/home/user/projects",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsFreshInstall()
			if result != tt.expected {
				t.Errorf("IsFreshInstall() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
