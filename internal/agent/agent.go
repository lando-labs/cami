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
	FilePath    string `yaml:"-"`
	Content     string `yaml:"-"`
}

// Metadata contains the YAML frontmatter data
type Metadata struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// LoadAgents reads all agents from the vc-agents directory
func LoadAgents(vcAgentsDir string) ([]*Agent, error) {
	files, err := os.ReadDir(vcAgentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read vc-agents directory: %w", err)
	}

	var agents []*Agent
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		agentPath := filepath.Join(vcAgentsDir, file.Name())
		agent, err := LoadAgent(agentPath)
		if err != nil {
			// Log error but continue loading other agents
			fmt.Fprintf(os.Stderr, "Warning: failed to load agent %s: %v\n", file.Name(), err)
			continue
		}
		agents = append(agents, agent)
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
