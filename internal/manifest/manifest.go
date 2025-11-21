package manifest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ProjectState represents the normalization state of a project
type ProjectState string

const (
	StateCAMINative ProjectState = "cami-native" // Fully normalized
	StateCAMILegacy ProjectState = "cami-legacy" // Old CAMI format
	StateCAMIAware  ProjectState = "cami-aware"  // Has agents, not tracked
	StateNonCAMI    ProjectState = "non-cami"    // No agents directory
)

// DeployedAgent represents an agent in a manifest
type DeployedAgent struct {
	Name           string    `yaml:"name"`
	Version        string    `yaml:"version"`
	Source         string    `yaml:"source"`                   // Source name
	SourcePath     string    `yaml:"source_path"`              // Full path to source file
	Priority       int       `yaml:"priority"`
	DeployedAt     time.Time `yaml:"deployed_at"`
	ContentHash    string    `yaml:"content_hash"`             // SHA256 of normalized content
	MetadataHash   string    `yaml:"metadata_hash"`            // SHA256 of frontmatter only
	CustomOverride bool      `yaml:"custom_override"`          // Intentionally customized
	NeedsUpgrade   bool      `yaml:"needs_upgrade,omitempty"`  // Missing version, etc.
}

// ProjectManifest represents a project's deployment manifest (local)
type ProjectManifest struct {
	Version      string          `yaml:"version"`        // Schema version
	State        ProjectState    `yaml:"state"`
	NormalizedAt time.Time       `yaml:"normalized_at"`
	Agents       []DeployedAgent `yaml:"agents"`
}

// ProjectDeployment represents a project in the central manifest
type ProjectDeployment struct {
	State        ProjectState    `yaml:"state"`
	NormalizedAt time.Time       `yaml:"normalized_at"`
	LastScanned  time.Time       `yaml:"last_scanned"`
	Agents       []DeployedAgent `yaml:"agents"`
}

// CentralManifest represents the central deployments manifest
type CentralManifest struct {
	Version               string                       `yaml:"version"`
	LastUpdated           time.Time                    `yaml:"last_updated"`
	ManifestFormatVersion int                          `yaml:"manifest_format_version"`
	Deployments           map[string]ProjectDeployment `yaml:"deployments"` // Key: absolute project path
}

const (
	// ManifestFormatVersion is the current manifest format version
	ManifestFormatVersion = 2

	// ProjectManifestFilename is the local manifest filename
	ProjectManifestFilename = ".claude/cami-manifest.yaml"

	// CentralManifestFilename is the central manifest filename
	CentralManifestFilename = "deployments.yaml"
)

// ReadProjectManifest reads a project's local manifest
func ReadProjectManifest(projectPath string) (*ProjectManifest, error) {
	manifestPath := filepath.Join(projectPath, ProjectManifestFilename)

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project manifest not found at %s", manifestPath)
		}
		return nil, fmt.Errorf("failed to read project manifest: %w", err)
	}

	var manifest ProjectManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse project manifest: %w", err)
	}

	return &manifest, nil
}

// WriteProjectManifest writes a project's local manifest
func WriteProjectManifest(projectPath string, manifest *ProjectManifest) error {
	manifestPath := filepath.Join(projectPath, ProjectManifestFilename)

	// Ensure .claude directory exists
	claudeDir := filepath.Join(projectPath, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal project manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project manifest: %w", err)
	}

	return nil
}

// ReadCentralManifest reads the central deployments manifest
func ReadCentralManifest() (*CentralManifest, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	manifestPath := filepath.Join(homeDir, "cami-workspace", CentralManifestFilename)

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty manifest if it doesn't exist yet
			return &CentralManifest{
				Version:               "2",
				LastUpdated:           time.Now(),
				ManifestFormatVersion: ManifestFormatVersion,
				Deployments:           make(map[string]ProjectDeployment),
			}, nil
		}
		return nil, fmt.Errorf("failed to read central manifest: %w", err)
	}

	var manifest CentralManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse central manifest: %w", err)
	}

	return &manifest, nil
}

// WriteCentralManifest writes the central deployments manifest
func WriteCentralManifest(manifest *CentralManifest) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	camiDir := filepath.Join(homeDir, "cami-workspace")
	if err := os.MkdirAll(camiDir, 0755); err != nil {
		return fmt.Errorf("failed to create cami-workspace directory: %w", err)
	}

	manifestPath := filepath.Join(camiDir, CentralManifestFilename)

	// Update last updated timestamp
	manifest.LastUpdated = time.Now()

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal central manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write central manifest: %w", err)
	}

	return nil
}

// CalculateContentHash calculates SHA256 hash of normalized file content
func CalculateContentHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Normalize content (strip excess whitespace, normalize line endings)
	normalized := NormalizeContent(data)

	hash := sha256.Sum256(normalized)
	return fmt.Sprintf("sha256:%x", hash), nil
}

// CalculateMetadataHash calculates SHA256 hash of frontmatter only
func CalculateMetadataHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Extract frontmatter
	frontmatter, err := extractFrontmatter(data)
	if err != nil {
		return "", fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	// Normalize and hash
	normalized := NormalizeContent(frontmatter)
	hash := sha256.Sum256(normalized)
	return fmt.Sprintf("sha256:%x", hash), nil
}

// NormalizeContent strips whitespace and normalizes line endings for hashing
func NormalizeContent(content []byte) []byte {
	// Convert to string for easier manipulation
	text := string(content)

	// Normalize line endings to \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Normalize multiple blank lines to single blank line
	re := regexp.MustCompile(`\n\n+`)
	text = re.ReplaceAllString(text, "\n\n")

	// Trim trailing spaces on each line
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	text = strings.Join(lines, "\n")

	return []byte(text)
}

// extractFrontmatter extracts YAML frontmatter from markdown content
func extractFrontmatter(content []byte) ([]byte, error) {
	text := string(content)

	// Check for frontmatter delimiters
	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		return nil, fmt.Errorf("no frontmatter found")
	}

	// Find the closing delimiter
	// Skip the opening "---"
	afterOpening := text[4:]
	closingIndex := strings.Index(afterOpening, "\n---\n")
	if closingIndex == -1 {
		closingIndex = strings.Index(afterOpening, "\r\n---\r\n")
		if closingIndex == -1 {
			return nil, fmt.Errorf("frontmatter not properly closed")
		}
	}

	// Extract frontmatter (without delimiters)
	frontmatter := afterOpening[:closingIndex]
	return []byte(frontmatter), nil
}

// CalculateFileHash is a helper that calculates hash of file content as-is
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return fmt.Sprintf("sha256:%x", hash.Sum(nil)), nil
}
