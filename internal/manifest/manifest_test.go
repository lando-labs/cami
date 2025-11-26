package manifest

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectManifestReadWrite(t *testing.T) {
	t.Run("write and read project manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test manifest
		manifest := &ProjectManifest{
			Version:      "2",
			State:        StateCAMINative,
			NormalizedAt: time.Now().Truncate(time.Second),
			Agents: []DeployedAgent{
				{
					Name:           "frontend",
					Version:        "1.1.0",
					Source:         "team-agents",
					SourcePath:     "/home/user/cami/sources/team-agents/frontend.md",
					Priority:       50,
					DeployedAt:     time.Now().Truncate(time.Second),
					ContentHash:    "sha256:abc123",
					MetadataHash:   "sha256:def456",
					CustomOverride: false,
					NeedsUpgrade:   false,
				},
			},
		}

		// Write manifest
		err := WriteProjectManifest(tmpDir, manifest)
		require.NoError(t, err)

		// Verify .claude directory was created
		claudeDir := filepath.Join(tmpDir, ".claude")
		assert.DirExists(t, claudeDir)

		// Verify manifest file exists
		manifestPath := filepath.Join(tmpDir, ProjectManifestFilename)
		assert.FileExists(t, manifestPath)

		// Read manifest back
		loaded, err := ReadProjectManifest(tmpDir)
		require.NoError(t, err)

		// Verify contents
		assert.Equal(t, manifest.Version, loaded.Version)
		assert.Equal(t, manifest.State, loaded.State)
		assert.True(t, manifest.NormalizedAt.Equal(loaded.NormalizedAt))
		assert.Len(t, loaded.Agents, 1)
		assert.Equal(t, "frontend", loaded.Agents[0].Name)
		assert.Equal(t, "1.1.0", loaded.Agents[0].Version)
		assert.Equal(t, "team-agents", loaded.Agents[0].Source)
		assert.Equal(t, 50, loaded.Agents[0].Priority)
		assert.False(t, loaded.Agents[0].CustomOverride)
	})

	t.Run("read non-existent manifest returns error", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := ReadProjectManifest(tmpDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("write creates .claude directory if missing", func(t *testing.T) {
		tmpDir := t.TempDir()

		manifest := &ProjectManifest{
			Version:      "2",
			State:        StateCAMINative,
			NormalizedAt: time.Now(),
			Agents:       []DeployedAgent{},
		}

		err := WriteProjectManifest(tmpDir, manifest)
		require.NoError(t, err)

		claudeDir := filepath.Join(tmpDir, ".claude")
		assert.DirExists(t, claudeDir)
	})
}

func TestCentralManifestReadWrite(t *testing.T) {
	t.Run("write and read central manifest", func(t *testing.T) {
		// Override home directory for testing
		tmpDir := t.TempDir()
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		// Create test manifest
		manifest := &CentralManifest{
			Version:               "2",
			LastUpdated:           time.Now().Truncate(time.Second),
			ManifestFormatVersion: ManifestFormatVersion,
			Deployments: map[string]ProjectDeployment{
				"/home/user/projects/my-app": {
					State:        StateCAMINative,
					NormalizedAt: time.Now().Truncate(time.Second),
					LastScanned:  time.Now().Truncate(time.Second),
					Agents: []DeployedAgent{
						{
							Name:         "frontend",
							Version:      "1.1.0",
							Source:       "team-agents",
							SourcePath:   "/home/user/cami/sources/team-agents/frontend.md",
							Priority:     50,
							DeployedAt:   time.Now().Truncate(time.Second),
							ContentHash:  "sha256:abc123",
							MetadataHash: "sha256:def456",
						},
					},
				},
			},
		}

		// Write manifest
		err := WriteCentralManifest(manifest)
		require.NoError(t, err)

		// Verify cami-workspace directory was created
		camiDir := filepath.Join(tmpDir, "cami-workspace")
		assert.DirExists(t, camiDir)

		// Verify manifest file exists
		manifestPath := filepath.Join(camiDir, CentralManifestFilename)
		assert.FileExists(t, manifestPath)

		// Read manifest back
		loaded, err := ReadCentralManifest()
		require.NoError(t, err)

		// Verify contents
		assert.Equal(t, "2", loaded.Version)
		assert.Equal(t, ManifestFormatVersion, loaded.ManifestFormatVersion)
		assert.Len(t, loaded.Deployments, 1)

		deployment, exists := loaded.Deployments["/home/user/projects/my-app"]
		assert.True(t, exists)
		assert.Equal(t, StateCAMINative, deployment.State)
		assert.Len(t, deployment.Agents, 1)
		assert.Equal(t, "frontend", deployment.Agents[0].Name)
	})

	t.Run("read non-existent central manifest returns empty", func(t *testing.T) {
		// Override home directory for testing
		tmpDir := t.TempDir()
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		manifest, err := ReadCentralManifest()

		require.NoError(t, err)
		assert.NotNil(t, manifest)
		assert.Equal(t, "2", manifest.Version)
		assert.Equal(t, ManifestFormatVersion, manifest.ManifestFormatVersion)
		assert.Empty(t, manifest.Deployments)
	})

	t.Run("write updates LastUpdated timestamp", func(t *testing.T) {
		// Override home directory for testing
		tmpDir := t.TempDir()
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)
		os.Setenv("HOME", tmpDir)

		oldTime := time.Now().Add(-1 * time.Hour).Truncate(time.Second)
		manifest := &CentralManifest{
			Version:               "2",
			LastUpdated:           oldTime,
			ManifestFormatVersion: ManifestFormatVersion,
			Deployments:           make(map[string]ProjectDeployment),
		}

		err := WriteCentralManifest(manifest)
		require.NoError(t, err)

		// Read back and verify timestamp was updated
		loaded, err := ReadCentralManifest()
		require.NoError(t, err)
		assert.True(t, loaded.LastUpdated.After(oldTime))
	})
}

func TestCalculateContentHash(t *testing.T) {
	t.Run("calculate hash of file content", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.md")

		content := `---
name: test-agent
version: 1.0.0
---

# Test Agent

This is a test agent.`

		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		hash, err := CalculateContentHash(testFile)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Contains(t, hash, "sha256:")
	})

	t.Run("normalized content produces same hash", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Content with different whitespace
		content1 := "line1\nline2\n\nline3"
		content2 := "line1\r\nline2\r\n\r\nline3\r\n"
		content3 := "line1  \nline2  \n\nline3  "

		file1 := filepath.Join(tmpDir, "file1.md")
		file2 := filepath.Join(tmpDir, "file2.md")
		file3 := filepath.Join(tmpDir, "file3.md")

		require.NoError(t, os.WriteFile(file1, []byte(content1), 0644))
		require.NoError(t, os.WriteFile(file2, []byte(content2), 0644))
		require.NoError(t, os.WriteFile(file3, []byte(content3), 0644))

		hash1, err := CalculateContentHash(file1)
		require.NoError(t, err)

		hash2, err := CalculateContentHash(file2)
		require.NoError(t, err)

		hash3, err := CalculateContentHash(file3)
		require.NoError(t, err)

		// All should produce same hash due to normalization
		assert.Equal(t, hash1, hash2)
		assert.Equal(t, hash1, hash3)
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		_, err := CalculateContentHash("/nonexistent/file.md")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
	})
}

func TestCalculateMetadataHash(t *testing.T) {
	t.Run("calculate hash of frontmatter only", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.md")

		content := `---
name: test-agent
version: 1.0.0
---

# Test Agent

This is a test agent.`

		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		hash, err := CalculateMetadataHash(testFile)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Contains(t, hash, "sha256:")
	})

	t.Run("same frontmatter produces same hash regardless of content", func(t *testing.T) {
		tmpDir := t.TempDir()

		frontmatter := `---
name: test-agent
version: 1.0.0
---`

		content1 := frontmatter + "\n\nContent A"
		content2 := frontmatter + "\n\nContent B (different)"

		file1 := filepath.Join(tmpDir, "file1.md")
		file2 := filepath.Join(tmpDir, "file2.md")

		require.NoError(t, os.WriteFile(file1, []byte(content1), 0644))
		require.NoError(t, os.WriteFile(file2, []byte(content2), 0644))

		hash1, err := CalculateMetadataHash(file1)
		require.NoError(t, err)

		hash2, err := CalculateMetadataHash(file2)
		require.NoError(t, err)

		// Metadata hashes should be equal (same frontmatter)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("error on file without frontmatter", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.md")

		content := "# No frontmatter here"
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		_, err = CalculateMetadataHash(testFile)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no frontmatter found")
	})
}

func TestNormalizeContent(t *testing.T) {
	t.Run("normalize line endings", func(t *testing.T) {
		input := "line1\r\nline2\rline3\n"
		expected := "line1\nline2\nline3"

		result := NormalizeContent([]byte(input))

		assert.Equal(t, expected, string(result))
	})

	t.Run("trim leading and trailing whitespace", func(t *testing.T) {
		input := "\n  \nline1\nline2\n  \n"
		expected := "line1\nline2"

		result := NormalizeContent([]byte(input))

		assert.Equal(t, expected, string(result))
	})

	t.Run("normalize multiple blank lines to single", func(t *testing.T) {
		input := "line1\n\n\n\nline2"
		expected := "line1\n\nline2"

		result := NormalizeContent([]byte(input))

		assert.Equal(t, expected, string(result))
	})

	t.Run("trim trailing spaces on lines", func(t *testing.T) {
		input := "line1  \nline2\t\nline3"
		expected := "line1\nline2\nline3"

		result := NormalizeContent([]byte(input))

		assert.Equal(t, expected, string(result))
	})

	t.Run("comprehensive normalization", func(t *testing.T) {
		input := "\r\n  line1  \r\n\r\n\r\nline2\t\r\n\r\nline3  \r\n\r\n"
		expected := "line1\n\nline2\n\nline3"

		result := NormalizeContent([]byte(input))

		assert.Equal(t, expected, string(result))
	})
}

func TestExtractFrontmatter(t *testing.T) {
	t.Run("extract valid frontmatter", func(t *testing.T) {
		content := []byte(`---
name: test-agent
version: 1.0.0
---

# Content here`)

		frontmatter, err := extractFrontmatter(content)

		require.NoError(t, err)
		assert.Contains(t, string(frontmatter), "name: test-agent")
		assert.Contains(t, string(frontmatter), "version: 1.0.0")
		assert.NotContains(t, string(frontmatter), "# Content here")
	})

	t.Run("error on missing opening delimiter", func(t *testing.T) {
		content := []byte(`name: test-agent
---

# Content`)

		_, err := extractFrontmatter(content)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no frontmatter found")
	})

	t.Run("error on missing closing delimiter", func(t *testing.T) {
		content := []byte(`---
name: test-agent
version: 1.0.0

# Content (no closing ---)`)

		_, err := extractFrontmatter(content)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not properly closed")
	})
}

func TestCalculateFileHash(t *testing.T) {
	t.Run("calculate raw file hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")

		content := "test content"
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		hash, err := CalculateFileHash(testFile)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Contains(t, hash, "sha256:")
	})

	t.Run("different content produces different hash", func(t *testing.T) {
		tmpDir := t.TempDir()

		file1 := filepath.Join(tmpDir, "file1.txt")
		file2 := filepath.Join(tmpDir, "file2.txt")

		require.NoError(t, os.WriteFile(file1, []byte("content A"), 0644))
		require.NoError(t, os.WriteFile(file2, []byte("content B"), 0644))

		hash1, err := CalculateFileHash(file1)
		require.NoError(t, err)

		hash2, err := CalculateFileHash(file2)
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		_, err := CalculateFileHash("/nonexistent/file.txt")

		assert.Error(t, err)
	})
}

func TestProjectState(t *testing.T) {
	t.Run("project states are distinct", func(t *testing.T) {
		states := []ProjectState{
			StateCAMINative,
			StateCAMILegacy,
			StateCAMIAware,
			StateNonCAMI,
		}

		// Ensure all states are unique
		seen := make(map[ProjectState]bool)
		for _, state := range states {
			assert.False(t, seen[state], "duplicate state: %s", state)
			seen[state] = true
		}

		assert.Len(t, seen, 4)
	})
}
