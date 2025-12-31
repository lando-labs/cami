package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skillset represents a collection of related skills
type Skillset struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags,omitempty"` // For discovery and tech stack matching
	Skills      []*Skill `yaml:"-"`              // Loaded skills within this skillset
	FilePath    string   `yaml:"-"`              // Path to SKILLSET.yaml or skillset directory
	SourceName  string   `yaml:"-"`              // Which source it came from
	SourcePath  string   `yaml:"-"`              // Full path to source
}

// SkillsetMetadata contains YAML frontmatter for SKILLSET.yaml (optional file)
type SkillsetMetadata struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags,omitempty"`
}

// LoadSkillset loads a skillset from a directory
// It looks for either SKILLSET.yaml for explicit metadata or derives from the skill
func LoadSkillset(dirPath string) (*Skillset, error) {
	// Check if this is a valid skillset directory (must have SKILL.md)
	skillPath := filepath.Join(dirPath, "SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no SKILL.md found in %s", dirPath)
	}

	skillset := &Skillset{
		FilePath: dirPath,
	}

	// Try to load optional SKILLSET.yaml for metadata
	skillsetYamlPath := filepath.Join(dirPath, "SKILLSET.yaml")
	if data, err := os.ReadFile(skillsetYamlPath); err == nil {
		var metadata SkillsetMetadata
		if err := yaml.Unmarshal(data, &metadata); err == nil {
			skillset.Name = metadata.Name
			skillset.Version = metadata.Version
			skillset.Description = metadata.Description
			skillset.Tags = metadata.Tags
		}
	}

	// Load the primary skill
	skill, err := LoadSkill(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load SKILL.md: %w", err)
	}

	// If no SKILLSET.yaml, derive metadata from the skill
	if skillset.Name == "" {
		skillset.Name = skill.Name
	}
	if skillset.Version == "" {
		skillset.Version = skill.Version
	}
	if skillset.Description == "" {
		skillset.Description = skill.Description
	}
	if len(skillset.Tags) == 0 {
		skillset.Tags = skill.Tags
	}

	skillset.Skills = []*Skill{skill}

	return skillset, nil
}

// HasExplicitMetadata returns true if the skillset has a SKILLSET.yaml file
func (s *Skillset) HasExplicitMetadata() bool {
	if s.FilePath == "" {
		return false
	}
	skillsetYamlPath := filepath.Join(s.FilePath, "SKILLSET.yaml")
	_, err := os.Stat(skillsetYamlPath)
	return err == nil
}

// PrimarySkill returns the main skill in this skillset
func (s *Skillset) PrimarySkill() *Skill {
	if len(s.Skills) > 0 {
		return s.Skills[0]
	}
	return nil
}

// SkillCount returns the number of skills in this skillset
func (s *Skillset) SkillCount() int {
	return len(s.Skills)
}

// SupportFileCount returns the total number of support files across all skills
func (s *Skillset) SupportFileCount() int {
	count := 0
	for _, skill := range s.Skills {
		count += len(skill.SupportFiles)
	}
	return count
}

// DirectoryName returns the name of the skillset directory
func (s *Skillset) DirectoryName() string {
	return filepath.Base(s.FilePath)
}

// AllFiles returns all files in the skillset (SKILL.md + support files)
func (s *Skillset) AllFiles() []string {
	var files []string
	for _, skill := range s.Skills {
		files = append(files, skill.FilePath)
		files = append(files, skill.SupportFiles...)
	}

	// Also include SKILLSET.yaml if it exists
	skillsetYamlPath := filepath.Join(s.FilePath, "SKILLSET.yaml")
	if _, err := os.Stat(skillsetYamlPath); err == nil {
		files = append(files, skillsetYamlPath)
	}

	return files
}

// MatchesTechStack checks if this skillset matches any of the given technologies
// Used for STRATEGIES.yaml integration
func (s *Skillset) MatchesTechStack(technologies []string) bool {
	if len(technologies) == 0 {
		return false
	}

	// Normalize skill name and description for matching
	normalizedName := strings.ToLower(s.Name)
	normalizedTags := make([]string, len(s.Tags))
	for i, tag := range s.Tags {
		normalizedTags[i] = strings.ToLower(tag)
	}

	for _, tech := range technologies {
		normalizedTech := normalizeTechName(tech)
		if normalizedTech == "" {
			continue
		}

		// Check name - must be a word match, not just substring
		if containsWord(normalizedName, normalizedTech) {
			return true
		}

		// Check tags - exact match or contains as word
		for _, tag := range normalizedTags {
			if tag == normalizedTech || containsWord(tag, normalizedTech) {
				return true
			}
		}
	}

	return false
}

// containsWord checks if s contains word as a whole word (not just substring)
func containsWord(s, word string) bool {
	// Split on common separators
	separators := []string{"-", "_", " ", "."}
	words := []string{s}

	for _, sep := range separators {
		var newWords []string
		for _, w := range words {
			newWords = append(newWords, strings.Split(w, sep)...)
		}
		words = newWords
	}

	for _, w := range words {
		if w == word {
			return true
		}
	}

	return false
}

// normalizeTechName normalizes a technology name for matching
// "React 19+" -> "react", "Tailwind CSS 4+" -> "tailwind"
func normalizeTechName(tech string) string {
	tech = strings.ToLower(tech)

	// Remove version numbers and common suffixes
	parts := strings.Fields(tech) // Split on whitespace
	if len(parts) == 0 {
		return ""
	}

	// Take first word (the main tech name)
	result := parts[0]

	// Remove common suffixes that aren't meaningful for matching
	result = strings.TrimSuffix(result, ".js")

	return strings.TrimSpace(result)
}
