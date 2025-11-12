package docs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/lando/cami/internal/agent"
)

const (
	sectionMarkerStart = "<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->"
	sectionMarkerEnd   = "<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->"
	defaultSectionName = "Deployed Agents"
)

var (
	// Regex to match both old and new marker formats
	// Old: <!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
	// New: <!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: 2025-10-09T14:30:00-05:00 -->
	sectionMarkerPattern = regexp.MustCompile(`<!-- CAMI-MANAGED: DEPLOYED-AGENTS(?:\s*\|\s*Last Updated:\s*[^>]+)?\s*-->`)
)

// UpdateCLAUDEmd updates the CLAUDE.md file with deployed agent information
func UpdateCLAUDEmd(projectPath, sectionName string, dryRun bool) (string, error) {
	if sectionName == "" {
		sectionName = defaultSectionName
	}

	claudePath := filepath.Join(projectPath, "CLAUDE.md")
	agentsDir := filepath.Join(projectPath, ".claude", "agents")

	// Check if .claude/agents directory exists
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		return "", fmt.Errorf("no agents directory found at %s", agentsDir)
	}

	// Scan deployed agents
	deployedAgents, err := scanDeployedAgents(agentsDir)
	if err != nil {
		return "", fmt.Errorf("failed to scan agents: %w", err)
	}

	if len(deployedAgents) == 0 {
		return "", fmt.Errorf("no agents found in %s", agentsDir)
	}

	// Generate agent section content
	agentSection := generateAgentSection(sectionName, deployedAgents)

	// Read existing CLAUDE.md if it exists
	var existingContent string
	if _, err := os.Stat(claudePath); err == nil {
		data, err := os.ReadFile(claudePath)
		if err != nil {
			return "", fmt.Errorf("failed to read CLAUDE.md: %w", err)
		}
		existingContent = string(data)
	}

	// Merge content
	newContent := mergeContent(existingContent, agentSection)

	if dryRun {
		return newContent, nil
	}

	// Write updated CLAUDE.md
	if err := os.WriteFile(claudePath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}

	return newContent, nil
}

// scanDeployedAgents scans the .claude/agents directory and returns agent info
func scanDeployedAgents(agentsDir string) ([]*agent.Agent, error) {
	files, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, err
	}

	var agents []*agent.Agent
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		agentPath := filepath.Join(agentsDir, file.Name())
		ag, err := agent.LoadAgent(agentPath)
		if err != nil {
			// Skip invalid agent files
			continue
		}
		agents = append(agents, ag)
	}

	return agents, nil
}

// generateAgentSection creates the markdown section for deployed agents
func generateAgentSection(sectionName string, agents []*agent.Agent) string {
	var sb strings.Builder

	// Generate timestamp in RFC3339 format (ISO 8601)
	timestamp := time.Now().Format(time.RFC3339)

	// Write start marker with timestamp
	sb.WriteString(fmt.Sprintf("<!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: %s -->\n", timestamp))
	sb.WriteString(fmt.Sprintf("## %s\n\n", sectionName))
	sb.WriteString("The following Claude Code agents are available in this project:\n\n")

	for _, ag := range agents {
		sb.WriteString(fmt.Sprintf("### %s", ag.Name))
		if ag.Version != "" {
			sb.WriteString(fmt.Sprintf(" (v%s)", ag.Version))
		}
		sb.WriteString("\n")
		if ag.Description != "" {
			sb.WriteString(ag.Description)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString(sectionMarkerEnd)
	sb.WriteString("\n")

	return sb.String()
}

// mergeContent merges new agent section into existing CLAUDE.md content
func mergeContent(existing, newSection string) string {
	// If no existing content, return just the new section
	if existing == "" {
		return newSection
	}

	// Check if managed section already exists (handles both old and new marker formats)
	startMatch := sectionMarkerPattern.FindStringIndex(existing)
	endIdx := strings.Index(existing, sectionMarkerEnd)

	if startMatch == nil || endIdx == -1 {
		// No existing managed section, append to end
		if !strings.HasSuffix(existing, "\n") {
			existing += "\n"
		}
		existing += "\n"
		return existing + newSection
	}

	startIdx := startMatch[0]

	// Replace existing managed section
	// Find the end of the end marker line
	endOfMarker := endIdx + len(sectionMarkerEnd)

	// Find the next newline after the end marker
	nextNewline := strings.Index(existing[endOfMarker:], "\n")
	if nextNewline != -1 {
		endOfMarker += nextNewline + 1
	}

	// Construct new content
	var sb strings.Builder
	sb.WriteString(existing[:startIdx])
	sb.WriteString(newSection)
	if endOfMarker < len(existing) {
		sb.WriteString(existing[endOfMarker:])
	}

	return sb.String()
}

// ScanDeployedAgentsInfo returns information about deployed agents
func ScanDeployedAgentsInfo(projectPath string) ([]*agent.Agent, error) {
	agentsDir := filepath.Join(projectPath, ".claude", "agents")

	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no agents directory found at %s", agentsDir)
	}

	return scanDeployedAgents(agentsDir)
}

// ExtractExistingSection extracts the current managed section if it exists
func ExtractExistingSection(projectPath string) (string, error) {
	claudePath := filepath.Join(projectPath, "CLAUDE.md")

	data, err := os.ReadFile(claudePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	content := string(data)
	startMatch := sectionMarkerPattern.FindStringIndex(content)
	endIdx := strings.Index(content, sectionMarkerEnd)

	if startMatch == nil || endIdx == -1 {
		return "", nil
	}

	startIdx := startMatch[0]
	endOfMarker := endIdx + len(sectionMarkerEnd)
	scanner := bufio.NewScanner(strings.NewReader(content[endOfMarker:]))
	if scanner.Scan() {
		endOfMarker += len(scanner.Text()) + 1
	}

	return content[startIdx:endOfMarker], nil
}
