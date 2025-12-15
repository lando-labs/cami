package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skill represents a Claude Code skill with metadata and content
type Skill struct {
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	Description  string   `yaml:"description"`
	Tags         []string `yaml:"tags,omitempty"`          // For discovery/matching
	AllowedTools []string `yaml:"allowed-tools,omitempty"` // Optional tool restrictions
	FilePath     string   `yaml:"-"`                       // Path to SKILL.md
	SkillsetPath string   `yaml:"-"`                       // Parent skillset directory
	SourceName   string   `yaml:"-"`                       // Which source it came from
	Content      string   `yaml:"-"`                       // SKILL.md content (after frontmatter)
	SupportFiles []string `yaml:"-"`                       // List of supporting .md files in skillset
}

// Metadata contains the YAML frontmatter data for skills
type Metadata struct {
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	Description  string   `yaml:"description"`
	Tags         []string `yaml:"tags,omitempty"`
	AllowedTools []string `yaml:"allowed-tools,omitempty"`
}

// LoadSkill parses a single SKILL.md file
func LoadSkill(filePath string) (*Skill, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	// Read first line, should be "---"
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}

	// Read frontmatter
	var frontmatterLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		frontmatterLines = append(frontmatterLines, line)
	}

	// Parse YAML frontmatter
	var metadata Metadata
	frontmatterYAML := strings.Join(frontmatterLines, "\n")
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Read the rest of the file (skill content)
	var contentLines []string
	for scanner.Scan() {
		contentLines = append(contentLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	content := strings.Join(contentLines, "\n")

	// Get skillset directory path
	skillsetPath := filepath.Dir(filePath)

	// Find support files in the same directory
	supportFiles, err := findSupportFiles(skillsetPath, filePath)
	if err != nil {
		// Log warning but continue
		fmt.Fprintf(os.Stderr, "Warning: failed to find support files for %s: %v\n", filePath, err)
		supportFiles = []string{}
	}

	return &Skill{
		Name:         metadata.Name,
		Version:      metadata.Version,
		Description:  metadata.Description,
		Tags:         metadata.Tags,
		AllowedTools: metadata.AllowedTools,
		FilePath:     filePath,
		SkillsetPath: skillsetPath,
		Content:      content,
		SupportFiles: supportFiles,
	}, nil
}

// findSupportFiles finds all .md files in the skillset directory except SKILL.md
func findSupportFiles(skillsetPath, skillMdPath string) ([]string, error) {
	var supportFiles []string

	err := filepath.Walk(skillsetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip SKILL.md itself
		if path == skillMdPath {
			return nil
		}

		// Skip SKILLSET.md
		if strings.ToUpper(info.Name()) == "SKILLSET.MD" {
			return nil
		}

		supportFiles = append(supportFiles, path)
		return nil
	})

	return supportFiles, err
}

// FullContent returns the complete skill file content including frontmatter
func (s *Skill) FullContent() string {
	frontmatter := fmt.Sprintf("---\nname: %s\nversion: %s\ndescription: %s\n",
		s.Name, s.Version, s.Description)

	// Add optional fields if present
	if len(s.Tags) > 0 {
		frontmatter += fmt.Sprintf("tags: [%s]\n", strings.Join(s.Tags, ", "))
	}
	if len(s.AllowedTools) > 0 {
		frontmatter += fmt.Sprintf("allowed-tools: [%s]\n", strings.Join(s.AllowedTools, ", "))
	}

	frontmatter += "---\n"

	return frontmatter + s.Content
}

// FileName returns just the filename without path (usually "SKILL.md")
func (s *Skill) FileName() string {
	return filepath.Base(s.FilePath)
}

// SkillsetName returns the name of the skillset (directory name)
func (s *Skill) SkillsetName() string {
	return filepath.Base(s.SkillsetPath)
}
