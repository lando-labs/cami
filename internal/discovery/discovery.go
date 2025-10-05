package discovery

import (
	"fmt"
	"os"
	"path/filepath"
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
