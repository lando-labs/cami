package normalize

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/backup"
	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/manifest"
)

// SourceIssue represents a problem with an agent in a source
type SourceIssue struct {
	AgentFile string
	Problems  []string // "missing version", "no description", etc.
}

// SourceAnalysis represents the compliance state of a source
type SourceAnalysis struct {
	SourceName        string
	Path              string
	IsCompliant       bool
	AgentCount        int
	Issues            []SourceIssue
	MissingCAMIIgnore bool
}

// SourceNormalizationOptions specifies what to fix
type SourceNormalizationOptions struct {
	AddVersions      bool // Add v1.0.0 to agents missing versions
	AddDescriptions  bool // Generate descriptions (use agent-architect)
	AddCategories    bool // Auto-categorize agents
	CreateCAMIIgnore bool // Add .camiignore template
}

// SourceNormalizationResult represents the outcome
type SourceNormalizationResult struct {
	Success       bool
	Changes       []string
	AgentsUpdated int
	BackupPath    string
}

// AgentAnalysis represents a single agent in a project
type AgentAnalysis struct {
	Name          string
	HasVersion    bool
	Version       string
	MatchesSource string // Source name if matches, empty if no match
	IsTracked     bool   // In manifest?
	NeedsUpgrade  bool
	FilePath      string
	ContentHash   string
	MetadataHash  string
}

// ProjectRecommendations suggests normalization actions
type ProjectRecommendations struct {
	MinimalRequired     bool // Create manifests
	StandardRecommended bool // Link sources, add versions
	FullOptional        bool // Rewrite with agent-architect
}

// ProjectAnalysis represents the state of a project for normalization
type ProjectAnalysis struct {
	Path            string
	State           manifest.ProjectState
	HasAgentsDir    bool
	HasManifest     bool
	AgentCount      int
	Agents          []AgentAnalysis
	Recommendations ProjectRecommendations
}

// ProjectNormalizationLevel specifies depth of normalization
type ProjectNormalizationLevel string

const (
	LevelMinimal  ProjectNormalizationLevel = "minimal"  // Just manifests
	LevelStandard ProjectNormalizationLevel = "standard" // Manifests + source links
	LevelFull     ProjectNormalizationLevel = "full"     // Complete rewrite
)

// ProjectNormalizationOptions specifies what to do
type ProjectNormalizationOptions struct {
	Level           ProjectNormalizationLevel
	UpgradeAgents   []string          // Agents to add versions to
	CopyToSource    map[string]string // agent name -> source name
	SkipAgents      []string          // Leave these alone
	CustomOverrides []string          // Mark as intentionally customized
}

// ProjectNormalizationResult represents the outcome
type ProjectNormalizationResult struct {
	Success      bool
	StateBefore  manifest.ProjectState
	StateAfter   manifest.ProjectState
	Changes      []string
	BackupPath   string
	UndoAvailable bool
}

// AnalyzeSource analyzes a source for CAMI compliance
func AnalyzeSource(sourceName string, sourcePath string) (*SourceAnalysis, error) {
	analysis := &SourceAnalysis{
		SourceName: sourceName,
		Path:       sourcePath,
	}

	// Check if path exists
	if _, err := os.Stat(sourcePath); err != nil {
		return nil, fmt.Errorf("source path does not exist: %w", err)
	}

	// Load all agents from source
	agents, err := agent.LoadAgentsFromPath(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load agents from source: %w", err)
	}

	analysis.AgentCount = len(agents)

	// Check for .camiignore
	camiIgnorePath := filepath.Join(sourcePath, ".camiignore")
	if _, err := os.Stat(camiIgnorePath); os.IsNotExist(err) {
		analysis.MissingCAMIIgnore = true
	}

	// Analyze each agent for issues
	for _, ag := range agents {
		var problems []string

		// Check for missing version
		if ag.Version == "" {
			problems = append(problems, "missing version")
		}

		// Check for missing description
		if ag.Description == "" {
			problems = append(problems, "missing description")
		}

		// Check for missing name
		if ag.Name == "" {
			problems = append(problems, "missing name")
		}

		// If there are problems, add to issues
		if len(problems) > 0 {
			analysis.Issues = append(analysis.Issues, SourceIssue{
				AgentFile: ag.FileName(),
				Problems:  problems,
			})
		}
	}

	// Source is compliant if no issues and has .camiignore
	analysis.IsCompliant = len(analysis.Issues) == 0 && !analysis.MissingCAMIIgnore

	return analysis, nil
}

// NormalizeSource fixes source agents to meet CAMI standards
func NormalizeSource(sourceName string, sourcePath string, options SourceNormalizationOptions) (*SourceNormalizationResult, error) {
	result := &SourceNormalizationResult{}

	// Create backup first
	backupPath, err := backup.CreateBackup(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}
	result.BackupPath = backupPath

	// Load all agents
	agents, err := agent.LoadAgentsFromPath(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load agents: %w", err)
	}

	// Process each agent
	for _, ag := range agents {
		updated := false

		// Add version if missing and requested
		if options.AddVersions && ag.Version == "" {
			ag.Version = "1.0.0"
			updated = true
			result.Changes = append(result.Changes, fmt.Sprintf("Added version 1.0.0 to %s", ag.FileName()))
		}

		// Add description placeholder if missing and requested
		if options.AddDescriptions && ag.Description == "" {
			ag.Description = fmt.Sprintf("Description for %s agent", ag.Name)
			updated = true
			result.Changes = append(result.Changes, fmt.Sprintf("Added description placeholder to %s", ag.FileName()))
		}

		// Write updated agent if changes were made
		if updated {
			if err := writeAgent(ag); err != nil {
				return nil, fmt.Errorf("failed to write agent %s: %w", ag.FileName(), err)
			}
			result.AgentsUpdated++
		}
	}

	// Create .camiignore if requested and missing
	if options.CreateCAMIIgnore {
		camiIgnorePath := filepath.Join(sourcePath, ".camiignore")
		if _, err := os.Stat(camiIgnorePath); os.IsNotExist(err) {
			if err := createCAMIIgnoreTemplate(camiIgnorePath); err != nil {
				return nil, fmt.Errorf("failed to create .camiignore: %w", err)
			}
			result.Changes = append(result.Changes, "Created .camiignore file")
		}
	}

	result.Success = true
	return result, nil
}

// AnalyzeProject analyzes a project's state for normalization
func AnalyzeProject(projectPath string, availableSources []config.AgentSource) (*ProjectAnalysis, error) {
	analysis := &ProjectAnalysis{
		Path: projectPath,
	}

	// Check if path exists
	if _, err := os.Stat(projectPath); err != nil {
		return nil, fmt.Errorf("project path does not exist: %w", err)
	}

	// Check for .claude/agents/ directory
	agentsDir := filepath.Join(projectPath, ".claude", "agents")
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		analysis.State = manifest.StateNonCAMI
		analysis.HasAgentsDir = false
		return analysis, nil
	}
	analysis.HasAgentsDir = true

	// Check for manifest
	manifestPath := filepath.Join(projectPath, manifest.ProjectManifestFilename)
	if _, err := os.Stat(manifestPath); err == nil {
		analysis.HasManifest = true

		// Read manifest to determine state
		projectManifest, err := manifest.ReadProjectManifest(projectPath)
		if err == nil {
			analysis.State = projectManifest.State
		}
	} else {
		// No manifest - determine if legacy or aware
		analysis.HasManifest = false
		// For now, assume cami-aware (has agents but not tracked)
		// TODO: Add logic to detect old CAMI format (legacy)
		analysis.State = manifest.StateCAMIAware
	}

	// Load all agents from project
	deployedAgents, err := agent.LoadAgentsFromPath(agentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load deployed agents: %w", err)
	}
	analysis.AgentCount = len(deployedAgents)

	// Build map of available agents from sources
	sourceAgentsMap := make(map[string]*agent.Agent)
	for _, source := range availableSources {
		sourceAgents, err := agent.LoadAgentsFromPath(source.Path)
		if err != nil {
			// Log but continue
			continue
		}
		for _, ag := range sourceAgents {
			// Use first occurrence (priority already handled by caller)
			if _, exists := sourceAgentsMap[ag.Name]; !exists {
				sourceAgentsMap[ag.Name] = ag
			}
		}
	}

	// Analyze each deployed agent
	for _, deployedAgent := range deployedAgents {
		agentAnalysis := AgentAnalysis{
			Name:       deployedAgent.Name,
			HasVersion: deployedAgent.Version != "",
			Version:    deployedAgent.Version,
			FilePath:   deployedAgent.FilePath,
		}

		// Calculate hashes
		if contentHash, err := manifest.CalculateContentHash(deployedAgent.FilePath); err == nil {
			agentAnalysis.ContentHash = contentHash
		}
		if metadataHash, err := manifest.CalculateMetadataHash(deployedAgent.FilePath); err == nil {
			agentAnalysis.MetadataHash = metadataHash
		}

		// Check if agent matches a source
		if sourceAgent, exists := sourceAgentsMap[deployedAgent.Name]; exists {
			// Found matching source
			agentAnalysis.MatchesSource = "available-sources" // Could be more specific

			// Check if version matches
			if deployedAgent.Version != sourceAgent.Version {
				agentAnalysis.NeedsUpgrade = true
			}
		}

		// Check if tracked in manifest
		if analysis.HasManifest {
			// TODO: Check if agent is in manifest
			// For now, assume tracked if manifest exists
			agentAnalysis.IsTracked = true
		}

		analysis.Agents = append(analysis.Agents, agentAnalysis)
	}

	// Generate recommendations
	if !analysis.HasManifest {
		analysis.Recommendations.MinimalRequired = true
	}

	hasUnversioned := false
	hasUnmatched := false
	for _, ag := range analysis.Agents {
		if !ag.HasVersion {
			hasUnversioned = true
		}
		if ag.MatchesSource == "" {
			hasUnmatched = true
		}
	}

	if hasUnversioned || hasUnmatched {
		analysis.Recommendations.StandardRecommended = true
	}

	return analysis, nil
}

// NormalizeProject creates manifests and links agents to sources
func NormalizeProject(projectPath string, options ProjectNormalizationOptions, availableSources []config.AgentSource) (*ProjectNormalizationResult, error) {
	result := &ProjectNormalizationResult{}

	// Analyze current state
	analysis, err := AnalyzeProject(projectPath, availableSources)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze project: %w", err)
	}
	result.StateBefore = analysis.State

	// Create backup
	backupPath, err := backup.CreateBackup(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}
	result.BackupPath = backupPath
	result.UndoAvailable = true

	// Build source map for lookup
	sourceMap := make(map[string]config.AgentSource)
	for _, source := range availableSources {
		sourceMap[source.Name] = source
	}

	// Process based on level
	switch options.Level {
	case LevelMinimal:
		// Just create manifests
		if err := createMinimalManifests(projectPath, analysis); err != nil {
			return nil, fmt.Errorf("failed to create manifests: %w", err)
		}
		result.Changes = append(result.Changes, "Created project manifest")
		result.StateAfter = manifest.StateCAMINative

	case LevelStandard:
		// Create manifests with source links
		if err := createStandardManifests(projectPath, analysis, availableSources); err != nil {
			return nil, fmt.Errorf("failed to create manifests: %w", err)
		}
		result.Changes = append(result.Changes, "Created project manifest with source links")
		result.StateAfter = manifest.StateCAMINative

	case LevelFull:
		// Full normalization (not implemented yet)
		return nil, fmt.Errorf("full normalization not yet implemented")
	}

	result.Success = true
	return result, nil
}

// writeAgent writes an agent back to its file
func writeAgent(ag *agent.Agent) error {
	content := ag.FullContent()
	return os.WriteFile(ag.FilePath, []byte(content), 0644)
}

// createCAMIIgnoreTemplate creates a .camiignore file with default patterns
func createCAMIIgnoreTemplate(path string) error {
	template := `# CAMI Ignore File
# Patterns to exclude from agent loading

# Common patterns
*.draft.md
*.tmp.md
.DS_Store
README.md
CHANGELOG.md

# Directories
examples/
templates/
`
	return os.WriteFile(path, []byte(template), 0644)
}

// createMinimalManifests creates basic project and central manifests
func createMinimalManifests(projectPath string, analysis *ProjectAnalysis) error {
	// Create project manifest
	projectManifest := &manifest.ProjectManifest{
		Version:      "2",
		State:        manifest.StateCAMINative,
		NormalizedAt: time.Now(),
		Agents:       []manifest.DeployedAgent{},
	}

	// Add agents from analysis
	for _, ag := range analysis.Agents {
		deployedAgent := manifest.DeployedAgent{
			Name:           ag.Name,
			Version:        ag.Version,
			Source:         "unknown", // Minimal doesn't link sources
			SourcePath:     "",
			Priority:       0,
			DeployedAt:     time.Now(),
			ContentHash:    ag.ContentHash,
			MetadataHash:   ag.MetadataHash,
			CustomOverride: false,
			NeedsUpgrade:   false,
		}
		projectManifest.Agents = append(projectManifest.Agents, deployedAgent)
	}

	// Write project manifest
	if err := manifest.WriteProjectManifest(projectPath, projectManifest); err != nil {
		return fmt.Errorf("failed to write project manifest: %w", err)
	}

	// Update central manifest
	if err := updateCentralManifest(projectPath, projectManifest); err != nil {
		return fmt.Errorf("failed to update central manifest: %w", err)
	}

	return nil
}

// createStandardManifests creates manifests with source links
func createStandardManifests(projectPath string, analysis *ProjectAnalysis, availableSources []config.AgentSource) error {
	// Build map of source agents
	sourceAgentsMap := make(map[string]*agent.Agent)
	sourcePriorityMap := make(map[string]int)
	sourceNameMap := make(map[string]string)

	for _, source := range availableSources {
		sourceAgents, err := agent.LoadAgentsFromPath(source.Path)
		if err != nil {
			continue
		}
		for _, ag := range sourceAgents {
			if _, exists := sourceAgentsMap[ag.Name]; !exists {
				sourceAgentsMap[ag.Name] = ag
				sourcePriorityMap[ag.Name] = source.Priority
				sourceNameMap[ag.Name] = source.Name
			}
		}
	}

	// Create project manifest
	projectManifest := &manifest.ProjectManifest{
		Version:      "2",
		State:        manifest.StateCAMINative,
		NormalizedAt: time.Now(),
		Agents:       []manifest.DeployedAgent{},
	}

	// Add agents with source links
	for _, ag := range analysis.Agents {
		deployedAgent := manifest.DeployedAgent{
			Name:           ag.Name,
			Version:        ag.Version,
			DeployedAt:     time.Now(),
			ContentHash:    ag.ContentHash,
			MetadataHash:   ag.MetadataHash,
			CustomOverride: false,
		}

		// Try to find source
		if sourceAgent, exists := sourceAgentsMap[ag.Name]; exists {
			deployedAgent.Source = sourceNameMap[ag.Name]
			deployedAgent.SourcePath = sourceAgent.FilePath
			deployedAgent.Priority = sourcePriorityMap[ag.Name]

			// Check if needs upgrade
			if ag.Version != sourceAgent.Version {
				deployedAgent.NeedsUpgrade = true
			}
		} else {
			// No source found - mark as unknown
			deployedAgent.Source = "unknown"
			deployedAgent.SourcePath = ""
			deployedAgent.Priority = 999
		}

		projectManifest.Agents = append(projectManifest.Agents, deployedAgent)
	}

	// Write project manifest
	if err := manifest.WriteProjectManifest(projectPath, projectManifest); err != nil {
		return fmt.Errorf("failed to write project manifest: %w", err)
	}

	// Update central manifest
	if err := updateCentralManifest(projectPath, projectManifest); err != nil {
		return fmt.Errorf("failed to update central manifest: %w", err)
	}

	return nil
}

// updateCentralManifest updates the central manifest with project info
func updateCentralManifest(projectPath string, projectManifest *manifest.ProjectManifest) error {
	// Read central manifest
	centralManifest, err := manifest.ReadCentralManifest()
	if err != nil {
		return err
	}

	// Get absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Update or add deployment
	centralManifest.Deployments[absPath] = manifest.ProjectDeployment{
		State:        projectManifest.State,
		NormalizedAt: projectManifest.NormalizedAt,
		LastScanned:  time.Now(),
		Agents:       projectManifest.Agents,
	}

	// Write central manifest
	return manifest.WriteCentralManifest(centralManifest)
}
