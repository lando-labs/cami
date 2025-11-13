package agent

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Agent represents a Claude agent with metadata and content
type Agent struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Category    string `yaml:"-"` // Folder name (e.g., "core", "specialized")
	FilePath    string `yaml:"-"`
	Content     string `yaml:"-"`
}

// Metadata contains the YAML frontmatter data
type Metadata struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// LoadAgentsFromPath reads all agents from a directory (supports nested folders)
func LoadAgentsFromPath(dir string) ([]*Agent, error) {
	return LoadAgents(dir)
}

// AgentSource represents a source with its priority
type AgentSource struct {
	Path     string
	Priority int
}

// LoadAgentsFromSources loads agents from multiple sources with priority-based deduplication
// Higher priority sources override lower priority sources when agent names conflict
func LoadAgentsFromSources(sources []AgentSource) ([]*Agent, error) {
	// Map to track highest priority agent for each name
	agentMap := make(map[string]*Agent)
	priorityMap := make(map[string]int)

	// Load agents from all sources
	for _, source := range sources {
		agents, err := LoadAgentsFromPath(source.Path)
		if err != nil {
			// Log error but continue with other sources
			fmt.Fprintf(os.Stderr, "Warning: failed to load agents from %s: %v\n", source.Path, err)
			continue
		}

		// Process each agent
		for _, agent := range agents {
			existingPriority, exists := priorityMap[agent.Name]

			// Add or replace agent based on priority
			if !exists || source.Priority > existingPriority {
				agentMap[agent.Name] = agent
				priorityMap[agent.Name] = source.Priority
			}
		}
	}

	// Convert map to slice
	var result []*Agent
	for _, agent := range agentMap {
		result = append(result, agent)
	}

	return result, nil
}

// LoadAgents reads all agents from the vc-agents directory (supports nested folders)
func LoadAgents(vcAgentsDir string) ([]*Agent, error) {
	var agents []*Agent

	// Walk the directory tree to support categorized agents in subdirectories
	err := filepath.Walk(vcAgentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a .md file
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Load the agent
		agent, err := LoadAgent(path)
		if err != nil {
			// Log error but continue loading other agents
			fmt.Fprintf(os.Stderr, "Warning: failed to load agent %s: %v\n", info.Name(), err)
			return nil
		}

		// Extract category from the folder structure
		// If agent is in vcAgentsDir/category/agent.md, category is extracted
		// If agent is in vcAgentsDir/agent.md, category is empty (uncategorized)
		relPath, err := filepath.Rel(vcAgentsDir, filepath.Dir(path))
		if err != nil {
			return err
		}

		if relPath != "." {
			// Extract the first directory level as the category
			parts := strings.Split(relPath, string(filepath.Separator))
			agent.Category = parts[0]
		}

		agents = append(agents, agent)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk vc-agents directory: %w", err)
	}

	return agents, nil
}

// LoadAgent parses a single agent file
func LoadAgent(filePath string) (*Agent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

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

	// Read the rest of the file (agent content)
	var contentLines []string
	for scanner.Scan() {
		contentLines = append(contentLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	content := strings.Join(contentLines, "\n")

	return &Agent{
		Name:        metadata.Name,
		Version:     metadata.Version,
		Description: metadata.Description,
		FilePath:    filePath,
		Content:     content,
	}, nil
}

// FullContent returns the complete agent file content including frontmatter
func (a *Agent) FullContent() string {
	frontmatter := fmt.Sprintf(`---
name: %s
version: %s
description: %s
---
`, a.Name, a.Version, a.Description)

	return frontmatter + a.Content
}

// FileName returns just the filename without path
func (a *Agent) FileName() string {
	return filepath.Base(a.FilePath)
}
