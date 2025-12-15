package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SkillSource represents a source with its priority (mirrors config.SkillSource)
type SkillSource struct {
	Name     string
	Path     string
	Priority int
}

// LoadSkillsFromSources loads skills from multiple sources with priority-based deduplication
// Lower priority numbers override higher priority numbers when skill names conflict
func LoadSkillsFromSources(sources []SkillSource) ([]*Skill, error) {
	// Map to track highest priority (lowest number) skill for each name
	skillMap := make(map[string]*Skill)
	priorityMap := make(map[string]int)

	// Load skills from all sources
	for _, source := range sources {
		skills, err := LoadSkillsFromPath(source.Path)
		if err != nil {
			// Log error but continue with other sources
			fmt.Fprintf(os.Stderr, "Warning: failed to load skills from %s: %v\n", source.Path, err)
			continue
		}

		// Process each skill
		for _, skill := range skills {
			skill.SourceName = source.Name

			existingPriority, exists := priorityMap[skill.Name]

			// Add or replace skill based on priority (lower number = higher priority)
			if !exists || source.Priority < existingPriority {
				skillMap[skill.Name] = skill
				priorityMap[skill.Name] = source.Priority
			}
		}
	}

	// Convert map to slice
	var result []*Skill
	for _, skill := range skillMap {
		result = append(result, skill)
	}

	return result, nil
}

// LoadSkillsetsFromSources loads skillsets from multiple sources
func LoadSkillsetsFromSources(sources []SkillSource) ([]*Skillset, error) {
	// Map to track highest priority skillset for each name
	skillsetMap := make(map[string]*Skillset)
	priorityMap := make(map[string]int)

	for _, source := range sources {
		skillsets, err := LoadSkillsetsFromPath(source.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load skillsets from %s: %v\n", source.Path, err)
			continue
		}

		for _, skillset := range skillsets {
			skillset.SourceName = source.Name
			skillset.SourcePath = source.Path

			// Also set source name on all skills within
			for _, skill := range skillset.Skills {
				skill.SourceName = source.Name
			}

			existingPriority, exists := priorityMap[skillset.Name]

			if !exists || source.Priority < existingPriority {
				skillsetMap[skillset.Name] = skillset
				priorityMap[skillset.Name] = source.Priority
			}
		}
	}

	var result []*Skillset
	for _, skillset := range skillsetMap {
		result = append(result, skillset)
	}

	return result, nil
}

// LoadSkillsFromPath loads all skills from a source directory
// Supports both:
// - sources/guild/skillsets/<skill-name>/SKILL.md (inside guild)
// - sources/skill-source/<skill-name>/SKILL.md (standalone skill source)
func LoadSkillsFromPath(sourcePath string) ([]*Skill, error) {
	var skills []*Skill

	// Load .camiignore patterns if they exist
	ignorePatterns, err := loadCamiIgnore(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load .camiignore for skills: %v\n", err)
		ignorePatterns = []string{}
	}

	// Check for skillsets/ directory (guild-integrated)
	skillsetsDir := filepath.Join(sourcePath, "skillsets")
	if info, err := os.Stat(skillsetsDir); err == nil && info.IsDir() {
		dirSkills, err := loadSkillsFromDir(skillsetsDir, ignorePatterns)
		if err == nil {
			skills = append(skills, dirSkills...)
		}
	}

	// Also check for top-level skill directories (standalone source)
	// Each directory with a SKILL.md is a skill
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return skills, nil // Return what we have
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "skillsets" {
			continue
		}

		// Skip ignored directories
		if shouldIgnore(entry.Name(), ignorePatterns) {
			continue
		}

		skillDir := filepath.Join(sourcePath, entry.Name())
		skillPath := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillPath); err == nil {
			skill, err := LoadSkill(skillPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load skill from %s: %v\n", skillPath, err)
				continue
			}
			skills = append(skills, skill)
		}
	}

	return skills, nil
}

// LoadSkillsetsFromPath loads all skillsets from a source directory
func LoadSkillsetsFromPath(sourcePath string) ([]*Skillset, error) {
	var skillsets []*Skillset

	ignorePatterns, _ := loadCamiIgnore(sourcePath)

	// Check for skillsets/ directory (guild-integrated)
	skillsetsDir := filepath.Join(sourcePath, "skillsets")
	if info, err := os.Stat(skillsetsDir); err == nil && info.IsDir() {
		dirSkillsets, err := loadSkillsetsFromDir(skillsetsDir, ignorePatterns)
		if err == nil {
			skillsets = append(skillsets, dirSkillsets...)
		}
	}

	// Also check for top-level skill directories (standalone source)
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return skillsets, nil
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "skillsets" {
			continue
		}

		if shouldIgnore(entry.Name(), ignorePatterns) {
			continue
		}

		skillDir := filepath.Join(sourcePath, entry.Name())
		skillPath := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillPath); err == nil {
			skillset, err := LoadSkillset(skillDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load skillset from %s: %v\n", skillDir, err)
				continue
			}
			skillsets = append(skillsets, skillset)
		}
	}

	return skillsets, nil
}

// loadSkillsFromDir loads skills from a specific directory
func loadSkillsFromDir(dir string, ignorePatterns []string) ([]*Skill, error) {
	var skills []*Skill

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if shouldIgnore(entry.Name(), ignorePatterns) {
			continue
		}

		skillDir := filepath.Join(dir, entry.Name())
		skillPath := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillPath); err == nil {
			skill, err := LoadSkill(skillPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load skill %s: %v\n", entry.Name(), err)
				continue
			}
			skills = append(skills, skill)
		}
	}

	return skills, nil
}

// loadSkillsetsFromDir loads skillsets from a specific directory
func loadSkillsetsFromDir(dir string, ignorePatterns []string) ([]*Skillset, error) {
	var skillsets []*Skillset

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if shouldIgnore(entry.Name(), ignorePatterns) {
			continue
		}

		skillsetDir := filepath.Join(dir, entry.Name())
		skillPath := filepath.Join(skillsetDir, "SKILL.md")

		if _, err := os.Stat(skillPath); err == nil {
			skillset, err := LoadSkillset(skillsetDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load skillset %s: %v\n", entry.Name(), err)
				continue
			}
			skillsets = append(skillsets, skillset)
		}
	}

	return skillsets, nil
}

// loadCamiIgnore reads and parses a .camiignore file
func loadCamiIgnore(dir string) ([]string, error) {
	ignorePath := filepath.Join(dir, ".camiignore")

	file, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open .camiignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .camiignore: %w", err)
	}

	return patterns, nil
}

// shouldIgnore checks if a path matches any ignore patterns
func shouldIgnore(relPath string, patterns []string) bool {
	fileName := filepath.Base(relPath)

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, fileName)
		if err == nil && matched {
			return true
		}

		matched, err = filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}

		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(relPath, dirPattern+"/") || relPath == dirPattern {
				return true
			}
		}
	}

	return false
}

// FindSkillByName finds a skill by name from a list of skills
func FindSkillByName(skills []*Skill, name string) *Skill {
	for _, skill := range skills {
		if skill.Name == name {
			return skill
		}
	}
	return nil
}

// FindSkillsetByName finds a skillset by name from a list of skillsets
func FindSkillsetByName(skillsets []*Skillset, name string) *Skillset {
	for _, skillset := range skillsets {
		if skillset.Name == name {
			return skillset
		}
	}
	return nil
}

// FilterSkillsByTags filters skills by tags (any match)
func FilterSkillsByTags(skills []*Skill, tags []string) []*Skill {
	if len(tags) == 0 {
		return skills
	}

	var filtered []*Skill
	for _, skill := range skills {
		for _, skillTag := range skill.Tags {
			for _, filterTag := range tags {
				if strings.EqualFold(skillTag, filterTag) {
					filtered = append(filtered, skill)
					goto next
				}
			}
		}
	next:
	}

	return filtered
}

