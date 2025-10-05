package deploy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lando/cami/internal/agent"
)

// Result represents the result of a deployment operation
type Result struct {
	Agent    *agent.Agent
	Success  bool
	Message  string
	Conflict bool
}

// DeployAgent deploys a single agent to a target location
func DeployAgent(ag *agent.Agent, targetPath string, overwrite bool) (*Result, error) {
	// Ensure .claude/agents directory exists
	agentsDir := filepath.Join(targetPath, ".claude", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return &Result{
			Agent:   ag,
			Success: false,
			Message: fmt.Sprintf("Failed to create agents directory: %v", err),
		}, nil
	}

	// Determine target file path
	targetFile := filepath.Join(agentsDir, ag.FileName())

	// Check for conflicts
	if _, err := os.Stat(targetFile); err == nil && !overwrite {
		return &Result{
			Agent:    ag,
			Success:  false,
			Conflict: true,
			Message:  "File already exists",
		}, nil
	}

	// Write agent file
	content := ag.FullContent()
	if err := os.WriteFile(targetFile, []byte(content), 0644); err != nil {
		return &Result{
			Agent:   ag,
			Success: false,
			Message: fmt.Sprintf("Failed to write file: %v", err),
		}, nil
	}

	return &Result{
		Agent:   ag,
		Success: true,
		Message: "Deployed successfully",
	}, nil
}

// DeployAgents deploys multiple agents to a target location
func DeployAgents(agents []*agent.Agent, targetPath string, overwrite bool) ([]*Result, error) {
	var results []*Result

	for _, ag := range agents {
		result, err := DeployAgent(ag, targetPath, overwrite)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// ValidateTargetPath ensures the target path is valid for deployment
func ValidateTargetPath(path string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist")
		}
		return fmt.Errorf("cannot access path: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	return nil
}

// CheckConflicts checks for existing agent files that would conflict
func CheckConflicts(agents []*agent.Agent, targetPath string) map[string]bool {
	conflicts := make(map[string]bool)
	agentsDir := filepath.Join(targetPath, ".claude", "agents")

	for _, ag := range agents {
		targetFile := filepath.Join(agentsDir, ag.FileName())
		if _, err := os.Stat(targetFile); err == nil {
			conflicts[ag.Name] = true
		}
	}

	return conflicts
}
