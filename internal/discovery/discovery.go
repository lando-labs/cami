package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
)

// DeploymentStatus represents the status of an agent at a location
type DeploymentStatus string

const (
	StatusUpToDate      DeploymentStatus = "up-to-date"
	StatusUpdateAvailable DeploymentStatus = "update-available"
	StatusNotDeployed   DeploymentStatus = "not-deployed"
	StatusUnknown       DeploymentStatus = "unknown"
)

// AgentStatus represents an agent's deployment status at a location
type AgentStatus struct {
	Agent          *agent.Agent
	DeployedVersion string
	AvailableVersion string
	Status         DeploymentStatus
	Location       *config.DeployLocation
}

// LocationStatus represents all agent statuses for a single location
type LocationStatus struct {
	Location     *config.DeployLocation
	AgentStatuses []*AgentStatus
	LastScanned  time.Time
}

// DiscoveryResult contains all discovered agent statuses across locations
type DiscoveryResult struct {
	LocationStatuses []*LocationStatus
	AvailableAgents  []*agent.Agent
}

// ScanLocation scans a deployment location for deployed agents
func ScanLocation(location *config.DeployLocation, availableAgents []*agent.Agent) (*LocationStatus, error) {
	agentsDir := filepath.Join(location.Path, ".claude", "agents")

	// Check if .claude/agents directory exists
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		// Directory doesn't exist - no agents deployed
		statuses := make([]*AgentStatus, 0, len(availableAgents))
		for _, ag := range availableAgents {
			statuses = append(statuses, &AgentStatus{
				Agent:            ag,
				DeployedVersion:  "",
				AvailableVersion: ag.Version,
				Status:           StatusNotDeployed,
				Location:         location,
			})
		}

		return &LocationStatus{
			Location:      location,
			AgentStatuses: statuses,
			LastScanned:   time.Now(),
		}, nil
	}

	// Read deployed agent files
	deployedAgents := make(map[string]*agent.Agent)
	files, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read agents directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		agentPath := filepath.Join(agentsDir, file.Name())
		deployedAgent, err := agent.LoadAgent(agentPath)
		if err != nil {
			// Skip agents we can't parse
			continue
		}

		deployedAgents[deployedAgent.Name] = deployedAgent
	}

	// Compare with available agents
	statuses := make([]*AgentStatus, 0, len(availableAgents))
	for _, availableAgent := range availableAgents {
		status := &AgentStatus{
			Agent:            availableAgent,
			AvailableVersion: availableAgent.Version,
			Location:         location,
		}

		if deployedAgent, exists := deployedAgents[availableAgent.Name]; exists {
			status.DeployedVersion = deployedAgent.Version

			if deployedAgent.Version == availableAgent.Version {
				status.Status = StatusUpToDate
			} else {
				status.Status = StatusUpdateAvailable
			}
		} else {
			status.DeployedVersion = ""
			status.Status = StatusNotDeployed
		}

		statuses = append(statuses, status)
	}

	return &LocationStatus{
		Location:      location,
		AgentStatuses: statuses,
		LastScanned:   time.Now(),
	}, nil
}

// ScanAllLocations scans all configured locations
func ScanAllLocations(locations []config.DeployLocation, availableAgents []*agent.Agent) (*DiscoveryResult, error) {
	locationStatuses := make([]*LocationStatus, 0, len(locations))

	for i := range locations {
		status, err := ScanLocation(&locations[i], availableAgents)
		if err != nil {
			// Skip locations that can't be scanned but continue with others
			continue
		}
		locationStatuses = append(locationStatuses, status)
	}

	return &DiscoveryResult{
		LocationStatuses: locationStatuses,
		AvailableAgents:  availableAgents,
	}, nil
}

// GetStatusSymbol returns the symbol to display for a status
func GetStatusSymbol(status DeploymentStatus) string {
	switch status {
	case StatusUpToDate:
		return "✓"
	case StatusUpdateAvailable:
		return "⚠"
	case StatusNotDeployed:
		return "○"
	default:
		return "?"
	}
}

// ProjectInfo represents a discovered Claude Code project
type ProjectInfo struct {
	Path         string         `json:"path"`
	HasClaude    bool           `json:"has_claude"`
	HasAgents    bool           `json:"has_agents"`
	AgentCount   int            `json:"agent_count"`
	Agents       []*agent.Agent `json:"agents,omitempty"`
	RelativePath string         `json:"relative_path,omitempty"`
}

// DiscoverOptions configures the discovery process
type DiscoverOptions struct {
	RootPath  string
	EmptyOnly bool   // Only return projects with .claude/ but no agents
	HasAgent  string // Only return projects that have a specific agent
	MaxDepth  int    // Maximum directory depth (0 = unlimited)
}

// DiscoverProjects finds all Claude Code projects in a directory tree
func DiscoverProjects(opts DiscoverOptions) ([]*ProjectInfo, error) {
	var projects []*ProjectInfo

	// Normalize root path
	rootPath, err := filepath.Abs(opts.RootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Verify root path exists
	if _, err := os.Stat(rootPath); err != nil {
		return nil, fmt.Errorf("path does not exist: %s", rootPath)
	}

	// Walk the directory tree
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories we can't access
			return nil
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this is a .claude directory
		if info.Name() == ".claude" {
			projectPath := filepath.Dir(path)

			// Check depth limit
			if opts.MaxDepth > 0 {
				relPath, err := filepath.Rel(rootPath, projectPath)
				if err == nil {
					depth := len(strings.Split(relPath, string(filepath.Separator)))
					if depth > opts.MaxDepth {
						return filepath.SkipDir
					}
				}
			}

			// Scan this project
			project, err := scanProject(projectPath, rootPath)
			if err != nil {
				// Log error but continue
				fmt.Fprintf(os.Stderr, "Warning: failed to scan %s: %v\n", projectPath, err)
				return filepath.SkipDir
			}

			// Apply filters
			if opts.EmptyOnly && project.HasAgents {
				return filepath.SkipDir
			}

			if opts.HasAgent != "" {
				found := false
				for _, ag := range project.Agents {
					if ag.Name == opts.HasAgent {
						found = true
						break
					}
				}
				if !found {
					return filepath.SkipDir
				}
			}

			projects = append(projects, project)

			// Skip descending into this .claude directory
			return filepath.SkipDir
		}

		// Skip common directories that won't have projects
		skipDirs := []string{"node_modules", ".git", "vendor", "dist", "build", ".next", ".venv", "__pycache__", "target", ".cache"}
		for _, skip := range skipDirs {
			if info.Name() == skip {
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

	return projects, nil
}

// scanProject scans a single project directory
func scanProject(projectPath, rootPath string) (*ProjectInfo, error) {
	info := &ProjectInfo{
		Path:      projectPath,
		HasClaude: false,
		HasAgents: false,
	}

	// Calculate relative path for display
	relPath, err := filepath.Rel(rootPath, projectPath)
	if err == nil && relPath != "." {
		info.RelativePath = relPath
	}

	// Check if .claude directory exists
	claudeDir := filepath.Join(projectPath, ".claude")
	if stat, err := os.Stat(claudeDir); err == nil && stat.IsDir() {
		info.HasClaude = true
	} else {
		return info, nil
	}

	// Check if .claude/agents directory exists
	agentsDir := filepath.Join(claudeDir, "agents")
	if stat, err := os.Stat(agentsDir); err != nil || !stat.IsDir() {
		return info, nil
	}

	// Load agents
	agents, err := agent.LoadAgentsFromPath(agentsDir)
	if err != nil {
		return info, fmt.Errorf("failed to load agents: %w", err)
	}

	if len(agents) > 0 {
		info.HasAgents = true
		info.AgentCount = len(agents)
		info.Agents = agents
	}

	return info, nil
}
