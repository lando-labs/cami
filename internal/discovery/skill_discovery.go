package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lando/cami/internal/skill"
)

// DeployedSkillStatus represents the status of a deployed skill
type DeployedSkillStatus struct {
	Name             string
	SkillsetName     string
	DeployedVersion  string
	AvailableVersion string
	Status           string // "up-to-date", "update-available", "not-in-sources"
	FilePath         string
	LinkedAgents     []string
}

// ScanDeployedSkills scans a project for deployed skills and compares with available versions
func ScanDeployedSkills(projectPath string, availableSkills []*skill.Skill) ([]DeployedSkillStatus, error) {
	skillsDir := filepath.Join(projectPath, ".claude", "skills")

	// Check if skills directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		return []DeployedSkillStatus{}, nil // No skills deployed
	}

	var results []DeployedSkillStatus

	// Scan skills directory
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillsetDir := filepath.Join(skillsDir, entry.Name())
		skillPath := filepath.Join(skillsetDir, "SKILL.md")

		// Check if SKILL.md exists
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		// Load the deployed skill
		deployedSkill, err := skill.LoadSkill(skillPath)
		if err != nil {
			// Log but continue
			continue
		}

		status := DeployedSkillStatus{
			Name:            deployedSkill.Name,
			SkillsetName:    entry.Name(),
			DeployedVersion: deployedSkill.Version,
			FilePath:        skillPath,
		}

		// Find corresponding skill in available sources
		var matchingSkill *skill.Skill
		for _, s := range availableSkills {
			if s.Name == deployedSkill.Name {
				matchingSkill = s
				break
			}
		}

		if matchingSkill == nil {
			status.Status = "not-in-sources"
			status.AvailableVersion = ""
		} else {
			status.AvailableVersion = matchingSkill.Version
			if deployedSkill.Version == matchingSkill.Version {
				status.Status = "up-to-date"
			} else {
				status.Status = "update-available"
			}
		}

		results = append(results, status)
	}

	return results, nil
}

// CountDeployedSkills returns the number of skills deployed to a project
func CountDeployedSkills(projectPath string) int {
	skillsDir := filepath.Join(projectPath, ".claude", "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			count++
		}
	}

	return count
}

// HasDeployedSkills checks if a project has any skills deployed
func HasDeployedSkills(projectPath string) bool {
	return CountDeployedSkills(projectPath) > 0
}

// GetSkillsDir returns the path to the skills directory
func GetSkillsDir(projectPath string) string {
	return filepath.Join(projectPath, ".claude", "skills")
}
