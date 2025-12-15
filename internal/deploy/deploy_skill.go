package deploy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lando/cami/internal/manifest"
	"github.com/lando/cami/internal/skill"
)

// SkillDeployResult represents the result of deploying a single skill
type SkillDeployResult struct {
	SkillName     string
	SkillsetName  string
	Success       bool
	Message       string
	FilesDeployed int
	LinkedAgent   string
}

// DeploySkillsOptions contains options for skill deployment
type DeploySkillsOptions struct {
	SkillNames    []string // Specific skills to deploy
	SkillsetName  string   // Deploy all skills from a skillset (alternative)
	TargetPath    string   // Project directory
	Overwrite     bool     // Overwrite existing skills
	LinkToAgent   string   // Optional: link skills to an agent
	AvailableSkills []*skill.Skill // Pre-loaded skills to deploy from
}

// DeploySkills deploys skills to a project's .claude/skills/ directory
func DeploySkills(opts DeploySkillsOptions) ([]SkillDeployResult, error) {
	if len(opts.SkillNames) == 0 && opts.SkillsetName == "" {
		return nil, fmt.Errorf("either skill_names or skillset_name must be provided")
	}

	if opts.TargetPath == "" {
		return nil, fmt.Errorf("target_path is required")
	}

	// Ensure target path exists
	if _, err := os.Stat(opts.TargetPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("target path does not exist: %s", opts.TargetPath)
	}

	// Create .claude/skills directory if it doesn't exist
	skillsDir := filepath.Join(opts.TargetPath, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create skills directory: %w", err)
	}

	var results []SkillDeployResult

	for _, skillName := range opts.SkillNames {
		result := deploySkill(skillName, skillsDir, opts)
		results = append(results, result)
	}

	return results, nil
}

// deploySkill deploys a single skill
func deploySkill(skillName string, skillsDir string, opts DeploySkillsOptions) SkillDeployResult {
	result := SkillDeployResult{
		SkillName: skillName,
	}

	// Find the skill in available skills
	var skillToDeploy *skill.Skill
	for _, s := range opts.AvailableSkills {
		if s.Name == skillName {
			skillToDeploy = s
			break
		}
	}

	if skillToDeploy == nil {
		result.Success = false
		result.Message = fmt.Sprintf("skill '%s' not found in available skills", skillName)
		return result
	}

	result.SkillsetName = skillToDeploy.SkillsetName()

	// Target directory for this skill
	targetDir := filepath.Join(skillsDir, result.SkillsetName)

	// Check if skill already exists
	if _, err := os.Stat(targetDir); err == nil {
		if !opts.Overwrite {
			result.Success = false
			result.Message = fmt.Sprintf("skill '%s' already exists at %s (use overwrite to replace)", skillName, targetDir)
			return result
		}
		// Remove existing skill directory
		if err := os.RemoveAll(targetDir); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("failed to remove existing skill: %v", err)
			return result
		}
	}

	// Create skill directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("failed to create skill directory: %v", err)
		return result
	}

	// Copy SKILL.md
	if err := copyFile(skillToDeploy.FilePath, filepath.Join(targetDir, "SKILL.md")); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("failed to copy SKILL.md: %v", err)
		return result
	}
	result.FilesDeployed++

	// Copy support files
	for _, supportFile := range skillToDeploy.SupportFiles {
		// Get relative path from skillset directory
		relPath, err := filepath.Rel(skillToDeploy.SkillsetPath, supportFile)
		if err != nil {
			continue
		}

		targetFile := filepath.Join(targetDir, relPath)

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
			continue
		}

		if err := copyFile(supportFile, targetFile); err != nil {
			continue
		}
		result.FilesDeployed++
	}

	result.Success = true
	result.Message = fmt.Sprintf("deployed skill '%s' (%d files)", skillName, result.FilesDeployed)
	result.LinkedAgent = opts.LinkToAgent

	return result
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// CreateDeployedSkill creates a DeployedSkill manifest entry from a deployment result
func CreateDeployedSkill(result SkillDeployResult, skill *skill.Skill, priority int) manifest.DeployedSkill {
	contentHash := ""
	if skill != nil {
		contentHash, _ = manifest.CalculateContentHash(skill.FilePath)
	}

	deployed := manifest.DeployedSkill{
		Name:         result.SkillName,
		SkillsetName: result.SkillsetName,
		DeployedAt:   time.Now(),
		FileCount:    result.FilesDeployed,
		ContentHash:  contentHash,
		Origin:       "cami",
	}

	if skill != nil {
		deployed.Version = skill.Version
		deployed.Source = skill.SourceName
		deployed.SourcePath = skill.SkillsetPath
		deployed.Priority = priority
	}

	if result.LinkedAgent != "" {
		deployed.LinkedAgents = []string{result.LinkedAgent}
	}

	return deployed
}
