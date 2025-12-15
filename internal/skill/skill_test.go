package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSkill(t *testing.T) {
	// Create a temporary directory with a test skill
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skill directory: %v", err)
	}

	// Create SKILL.md
	skillContent := `---
name: test-skill
version: "1.0.0"
description: A test skill for unit testing
tags: [testing, example]
allowed-tools: [Read, Write]
---

# Test Skill

This is a test skill content.
`
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	// Create a support file
	supportContent := `# Components

Some component documentation.
`
	supportPath := filepath.Join(skillDir, "components.md")
	if err := os.WriteFile(supportPath, []byte(supportContent), 0644); err != nil {
		t.Fatalf("failed to write support file: %v", err)
	}

	// Load the skill
	skill, err := LoadSkill(skillPath)
	if err != nil {
		t.Fatalf("failed to load skill: %v", err)
	}

	// Verify metadata
	if skill.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got '%s'", skill.Name)
	}
	if skill.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", skill.Version)
	}
	if skill.Description != "A test skill for unit testing" {
		t.Errorf("unexpected description: %s", skill.Description)
	}
	if len(skill.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(skill.Tags))
	}
	if len(skill.AllowedTools) != 2 {
		t.Errorf("expected 2 allowed tools, got %d", len(skill.AllowedTools))
	}
	if len(skill.SupportFiles) != 1 {
		t.Errorf("expected 1 support file, got %d", len(skill.SupportFiles))
	}
}

func TestLoadSkillset(t *testing.T) {
	tmpDir := t.TempDir()
	skillsetDir := filepath.Join(tmpDir, "react-tailwind")
	if err := os.MkdirAll(skillsetDir, 0755); err != nil {
		t.Fatalf("failed to create skillset directory: %v", err)
	}

	// Create SKILL.md
	skillContent := `---
name: react-tailwind
version: "1.0.0"
description: React + Tailwind patterns
tags: [react, tailwind, frontend]
---

# React Tailwind Skill
`
	skillPath := filepath.Join(skillsetDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	// Load the skillset
	skillset, err := LoadSkillset(skillsetDir)
	if err != nil {
		t.Fatalf("failed to load skillset: %v", err)
	}

	if skillset.Name != "react-tailwind" {
		t.Errorf("expected name 'react-tailwind', got '%s'", skillset.Name)
	}
	if skillset.SkillCount() != 1 {
		t.Errorf("expected 1 skill, got %d", skillset.SkillCount())
	}
}

func TestLoadSkillsetWithMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	skillsetDir := filepath.Join(tmpDir, "react-tailwind")
	if err := os.MkdirAll(skillsetDir, 0755); err != nil {
		t.Fatalf("failed to create skillset directory: %v", err)
	}

	// Create SKILLSET.yaml
	skillsetYaml := `name: react-tailwind-skillset
version: "2.0.0"
description: A comprehensive React + Tailwind skillset
tags: [react, tailwind, styling]
`
	skillsetYamlPath := filepath.Join(skillsetDir, "SKILLSET.yaml")
	if err := os.WriteFile(skillsetYamlPath, []byte(skillsetYaml), 0644); err != nil {
		t.Fatalf("failed to write SKILLSET.yaml: %v", err)
	}

	// Create SKILL.md
	skillContent := `---
name: react-tailwind
version: "1.0.0"
description: React + Tailwind patterns
---

# React Tailwind Skill
`
	skillPath := filepath.Join(skillsetDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	// Load the skillset
	skillset, err := LoadSkillset(skillsetDir)
	if err != nil {
		t.Fatalf("failed to load skillset: %v", err)
	}

	// Should use SKILLSET.yaml metadata
	if skillset.Name != "react-tailwind-skillset" {
		t.Errorf("expected name from SKILLSET.yaml, got '%s'", skillset.Name)
	}
	if skillset.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got '%s'", skillset.Version)
	}
	if !skillset.HasExplicitMetadata() {
		t.Error("expected HasExplicitMetadata() to return true")
	}
}

func TestLoadSkillsFromPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a source with skillsets directory
	skillsetsDir := filepath.Join(tmpDir, "skillsets")
	if err := os.MkdirAll(skillsetsDir, 0755); err != nil {
		t.Fatalf("failed to create skillsets directory: %v", err)
	}

	// Create skill 1
	skill1Dir := filepath.Join(skillsetsDir, "skill1")
	if err := os.MkdirAll(skill1Dir, 0755); err != nil {
		t.Fatalf("failed to create skill1 directory: %v", err)
	}
	skill1Content := `---
name: skill1
version: "1.0.0"
description: Skill 1
---

# Skill 1
`
	if err := os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1Content), 0644); err != nil {
		t.Fatalf("failed to write skill1: %v", err)
	}

	// Create skill 2
	skill2Dir := filepath.Join(skillsetsDir, "skill2")
	if err := os.MkdirAll(skill2Dir, 0755); err != nil {
		t.Fatalf("failed to create skill2 directory: %v", err)
	}
	skill2Content := `---
name: skill2
version: "1.0.0"
description: Skill 2
tags: [backend]
---

# Skill 2
`
	if err := os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2Content), 0644); err != nil {
		t.Fatalf("failed to write skill2: %v", err)
	}

	// Load skills
	skills, err := LoadSkillsFromPath(tmpDir)
	if err != nil {
		t.Fatalf("failed to load skills: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

func TestLoadSkillsFromSources(t *testing.T) {
	// Create two sources with overlapping skills
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	// Source 1: high priority (low number)
	skill1Dir := filepath.Join(tmpDir1, "frontend")
	if err := os.MkdirAll(skill1Dir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	skill1Content := `---
name: frontend
version: "2.0.0"
description: Frontend from source 1 (higher priority)
---

# Frontend v2
`
	if err := os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1Content), 0644); err != nil {
		t.Fatalf("failed to write skill: %v", err)
	}

	// Source 2: low priority (high number)
	skill2Dir := filepath.Join(tmpDir2, "frontend")
	if err := os.MkdirAll(skill2Dir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	skill2Content := `---
name: frontend
version: "1.0.0"
description: Frontend from source 2 (lower priority)
---

# Frontend v1
`
	if err := os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2Content), 0644); err != nil {
		t.Fatalf("failed to write skill: %v", err)
	}

	// Load from both sources
	sources := []SkillSource{
		{Name: "source1", Path: tmpDir1, Priority: 10},
		{Name: "source2", Path: tmpDir2, Priority: 50},
	}

	skills, err := LoadSkillsFromSources(sources)
	if err != nil {
		t.Fatalf("failed to load skills: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill (deduplicated), got %d", len(skills))
	}

	// Should be version 2.0.0 from source1 (higher priority)
	frontend := FindSkillByName(skills, "frontend")
	if frontend == nil {
		t.Fatal("expected to find frontend skill")
	}
	if frontend.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0' from higher priority source, got '%s'", frontend.Version)
	}
}

func TestMatchesTechStack(t *testing.T) {
	skillset := &Skillset{
		Name:        "react-tailwind",
		Description: "React 19+ with Tailwind CSS implementation patterns",
		Tags:        []string{"react", "tailwind", "frontend"},
	}

	tests := []struct {
		name         string
		technologies []string
		expected     bool
	}{
		{
			name:         "matches react",
			technologies: []string{"React 19+"},
			expected:     true,
		},
		{
			name:         "matches tailwind",
			technologies: []string{"Tailwind CSS 4+"},
			expected:     true,
		},
		{
			name:         "matches multiple",
			technologies: []string{"React 19+", "Tailwind CSS 4+"},
			expected:     true,
		},
		{
			name:         "no match",
			technologies: []string{"Vue.js", "CSS Modules"},
			expected:     false,
		},
		{
			name:         "empty technologies",
			technologies: []string{},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skillset.MatchesTechStack(tt.technologies)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilterSkillsByTags(t *testing.T) {
	skills := []*Skill{
		{Name: "react-tailwind", Tags: []string{"react", "tailwind", "frontend"}},
		{Name: "react-mui", Tags: []string{"react", "mui", "frontend"}},
		{Name: "express-api", Tags: []string{"express", "backend", "api"}},
	}

	// Filter by frontend tag
	filtered := FilterSkillsByTags(skills, []string{"frontend"})
	if len(filtered) != 2 {
		t.Errorf("expected 2 frontend skills, got %d", len(filtered))
	}

	// Filter by react tag
	filtered = FilterSkillsByTags(skills, []string{"react"})
	if len(filtered) != 2 {
		t.Errorf("expected 2 react skills, got %d", len(filtered))
	}

	// Filter by backend tag
	filtered = FilterSkillsByTags(skills, []string{"backend"})
	if len(filtered) != 1 {
		t.Errorf("expected 1 backend skill, got %d", len(filtered))
	}

	// Filter by nonexistent tag
	filtered = FilterSkillsByTags(skills, []string{"nonexistent"})
	if len(filtered) != 0 {
		t.Errorf("expected 0 skills, got %d", len(filtered))
	}
}
