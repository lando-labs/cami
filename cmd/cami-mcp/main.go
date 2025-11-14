package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/deploy"
	"github.com/lando/cami/internal/docs"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "cami"
	serverVersion = "0.1.0"
)

// getVCAgentsDir determines the vc-agents directory location
func getVCAgentsDir() (string, error) {
	// First check environment variable
	if envDir := os.Getenv("CAMI_VC_AGENTS_DIR"); envDir != "" {
		if _, err := os.Stat(envDir); err == nil {
			return envDir, nil
		}
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %v", err)
	}

	// Try working directory first
	wd, _ := os.Getwd()
	vcAgentsPath := filepath.Join(wd, "vc-agents")
	if _, err := os.Stat(vcAgentsPath); err == nil {
		return vcAgentsPath, nil
	}

	// Try relative to executable
	execDir := filepath.Dir(execPath)
	vcAgentsPath = filepath.Join(execDir, "vc-agents")
	if _, err := os.Stat(vcAgentsPath); err == nil {
		return vcAgentsPath, nil
	}

	return "", fmt.Errorf("vc-agents directory not found")
}

// DeployAgentsArgs represents the arguments for deploy_agents tool
type DeployAgentsArgs struct {
	AgentNames []string `json:"agent_names" jsonschema_description:"Array of agent names to deploy (e.g. ['architect', 'backend'])"`
	TargetPath string   `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
	Overwrite  bool     `json:"overwrite,omitempty" jsonschema_description:"Whether to overwrite existing agent files (default: false)"`
}

// UpdateClaudeMdArgs represents the arguments for update_claude_md tool
type UpdateClaudeMdArgs struct {
	TargetPath string `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
}

// ScanDeployedAgentsArgs represents the arguments for scan_deployed_agents tool
type ScanDeployedAgentsArgs struct {
	TargetPath string `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
}

// DeployResult represents the result of a deployment operation
type DeployResult struct {
	AgentName string `json:"agent_name"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Conflict  bool   `json:"conflict,omitempty"`
}

// DeployAgentsResponse wraps the deployment results (MCP requires object, not array)
type DeployAgentsResponse struct {
	Results []DeployResult `json:"results"`
}

// AgentInfo represents basic agent information
type AgentInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Category    string `json:"category"` // Folder category (e.g., "core", "specialized")
	FileName    string `json:"file_name"`
}

// ListAgentsResponse wraps the agent list (MCP requires object, not array)
type ListAgentsResponse struct {
	Agents []AgentInfo `json:"agents"`
}

// AgentStatusInfo represents deployment status information
type AgentStatusInfo struct {
	Name             string `json:"name"`
	DeployedVersion  string `json:"deployed_version"`
	AvailableVersion string `json:"available_version"`
	Status           string `json:"status"`
}

// ScanDeployedAgentsResponse wraps the status info (MCP requires object, not array)
type ScanDeployedAgentsResponse struct {
	Statuses []AgentStatusInfo `json:"statuses"`
}

// AddLocationArgs represents the arguments for add_location tool
type AddLocationArgs struct {
	Name string `json:"name" jsonschema_description:"Friendly name for the location (e.g. 'my-project')"`
	Path string `json:"path" jsonschema_description:"Absolute path to project directory"`
}

// RemoveLocationArgs represents the arguments for remove_location tool
type RemoveLocationArgs struct {
	Name string `json:"name" jsonschema_description:"Name of location to remove"`
}

// LocationInfo represents a deployment location
type LocationInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ListLocationsResponse wraps the location list (MCP requires object, not array)
type ListLocationsResponse struct {
	Locations []LocationInfo `json:"locations"`
}

// AddSourceArgs represents the arguments for add_source tool
type AddSourceArgs struct {
	URL      string `json:"url" jsonschema_description:"Git URL to clone (e.g. 'git@github.com:lando/lando-agents.git')"`
	Name     string `json:"name,omitempty" jsonschema_description:"Name for the source (derived from URL if not specified)"`
	Priority int    `json:"priority,omitempty" jsonschema_description:"Priority (higher = higher precedence, default: 100)"`
}

// UpdateSourceArgs represents the arguments for update_source tool
type UpdateSourceArgs struct {
	Name string `json:"name,omitempty" jsonschema_description:"Name of source to update (updates all if not specified)"`
}

// SourceInfo represents an agent source
type SourceInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Priority    int    `json:"priority"`
	AgentCount  int    `json:"agent_count"`
	GitRemote   string `json:"git_remote,omitempty"`
	GitEnabled  bool   `json:"git_enabled"`
}

// ListSourcesResponse wraps the source list (MCP requires object, not array)
type ListSourcesResponse struct {
	Sources []SourceInfo `json:"sources"`
}

// OnboardingState represents the current onboarding state
type OnboardingState struct {
	ConfigExists    bool   `json:"config_exists"`
	SourceCount     int    `json:"source_count"`
	LocationCount   int    `json:"location_count"`
	HasAgentArch    bool   `json:"has_agent_architect"`
	TotalAgents     int    `json:"total_agents"`
	DeployedAgents  int    `json:"deployed_agents"`
	RecommendedNext string `json:"recommended_next"`
}

// OnboardResponse wraps the onboarding state
type OnboardResponse struct {
	State OnboardingState `json:"state"`
}

func main() {
	// Initialize logger to stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get vc-agents directory
	vcAgentsDir, err := getVCAgentsDir()
	if err != nil {
		log.Fatalf("Failed to locate vc-agents directory: %v", err)
	}

	log.Printf("Using vc-agents directory: %s", vcAgentsDir)

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register deploy_agents tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "deploy_agents",
		Description: "Deploy selected agents to a target project's .claude/agents/ directory. " +
			"Use this when the user wants to add specific agents to a project. " +
			"Handles conflict detection and creates necessary directories.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args DeployAgentsArgs) (*mcp.CallToolResult, any, error) {
		// Validate target path
		if err := deploy.ValidateTargetPath(args.TargetPath); err != nil {
			return nil, nil, fmt.Errorf("invalid target path: %w", err)
		}

		// Load all available agents
		allAgents, err := agent.LoadAgents(vcAgentsDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load agents: %w", err)
		}

		// Filter requested agents
		agentMap := make(map[string]*agent.Agent)
		for _, ag := range allAgents {
			agentMap[ag.Name] = ag
		}

		var agentsToDeploy []*agent.Agent
		var notFound []string
		for _, name := range args.AgentNames {
			if ag, exists := agentMap[name]; exists {
				agentsToDeploy = append(agentsToDeploy, ag)
			} else {
				notFound = append(notFound, name)
			}
		}

		if len(notFound) > 0 {
			return nil, nil, fmt.Errorf("agents not found: %v", notFound)
		}

		// Deploy agents
		results, err := deploy.DeployAgents(agentsToDeploy, args.TargetPath, args.Overwrite)
		if err != nil {
			return nil, nil, fmt.Errorf("deployment failed: %w", err)
		}

		// Convert results to response format
		var deployResults []DeployResult
		for _, result := range results {
			deployResults = append(deployResults, DeployResult{
				AgentName: result.Agent.Name,
				Success:   result.Success,
				Message:   result.Message,
				Conflict:  result.Conflict,
			})
		}

		// Format response
		responseText := fmt.Sprintf("Deployed %d agents to %s\n\n", len(agentsToDeploy), args.TargetPath)
		for _, result := range deployResults {
			status := "✓"
			if !result.Success {
				status = "✗"
			}
			responseText += fmt.Sprintf("%s %s: %s\n", status, result.AgentName, result.Message)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &DeployAgentsResponse{Results: deployResults}, nil
	})

	// Register update_claude_md tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "update_claude_md",
		Description: "Update a project's CLAUDE.md file with documentation about deployed agents. " +
			"Adds or updates the 'Available Agents' section. " +
			"Use this after deploying agents to document them.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args UpdateClaudeMdArgs) (*mcp.CallToolResult, any, error) {
		// Validate target path
		if err := deploy.ValidateTargetPath(args.TargetPath); err != nil {
			return nil, nil, fmt.Errorf("invalid target path: %w", err)
		}

		// Update CLAUDE.md
		_, err := docs.UpdateCLAUDEmd(args.TargetPath, "Deployed Agents", false)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to update CLAUDE.md: %w", err)
		}

		// Get deployed agents for response
		deployedAgents, err := docs.ScanDeployedAgentsInfo(args.TargetPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan deployed agents: %w", err)
		}

		responseText := fmt.Sprintf("Updated CLAUDE.md at %s\n\n", args.TargetPath)
		responseText += fmt.Sprintf("Documented %d agents:\n", len(deployedAgents))
		for _, ag := range deployedAgents {
			responseText += fmt.Sprintf("  • %s (v%s)\n", ag.Name, ag.Version)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Register list_agents tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_agents",
		Description: "List all available agents from CAMI's version-controlled agent repository. " +
			"Returns agent names, versions, descriptions, and categories. " +
			"Use this to discover what agents are available for deployment.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		// Load all available agents
		agents, err := agent.LoadAgents(vcAgentsDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load agents: %w", err)
		}

		// Group agents by category
		categoryMap := make(map[string][]*agent.Agent)
		for _, ag := range agents {
			category := ag.Category
			if category == "" {
				category = "uncategorized"
			}
			categoryMap[category] = append(categoryMap[category], ag)
		}

		// Determine category display order
		categoryOrder := []string{"core", "specialized", "infrastructure", "integration", "design", "meta", "uncategorized"}

		// Convert to response format and build text output
		var agentInfos []AgentInfo
		responseText := fmt.Sprintf("Available agents (%d total):\n\n", len(agents))

		for _, category := range categoryOrder {
			categoryAgents, exists := categoryMap[category]
			if !exists || len(categoryAgents) == 0 {
				continue
			}

			// Capitalize category name for display
			displayCategory := category
			if category != "uncategorized" {
				displayCategory = strings.ToUpper(category[:1]) + category[1:]
			} else {
				displayCategory = "Uncategorized"
			}

			responseText += fmt.Sprintf("## %s (%d agents)\n\n", displayCategory, len(categoryAgents))

			for _, ag := range categoryAgents {
				agentInfos = append(agentInfos, AgentInfo{
					Name:        ag.Name,
					Version:     ag.Version,
					Description: ag.Description,
					Category:    ag.Category,
					FileName:    ag.FileName(),
				})
				responseText += fmt.Sprintf("• %s (v%s)\n  %s\n\n", ag.Name, ag.Version, ag.Description)
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &ListAgentsResponse{Agents: agentInfos}, nil
	})

	// Register scan_deployed_agents tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "scan_deployed_agents",
		Description: "Scan a project directory to find deployed agents and compare with available versions. " +
			"Returns agent status (current, outdated, unknown). " +
			"Use this to audit what agents are deployed.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ScanDeployedAgentsArgs) (*mcp.CallToolResult, any, error) {
		// Validate target path
		if err := deploy.ValidateTargetPath(args.TargetPath); err != nil {
			return nil, nil, fmt.Errorf("invalid target path: %w", err)
		}

		// Load available agents
		availableAgents, err := agent.LoadAgents(vcAgentsDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load agents: %w", err)
		}

		// Scan the agents directory directly
		agentsDir := filepath.Join(args.TargetPath, ".claude", "agents")

		var statusInfos []AgentStatusInfo
		responseText := fmt.Sprintf("Scanning %s\n\n", args.TargetPath)

		// Check if agents directory exists
		if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
			responseText += "No .claude/agents directory found.\n"
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: responseText},
				},
			}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
		}

		// Load deployed agents
		deployedAgents := make(map[string]*agent.Agent)
		files, err := os.ReadDir(agentsDir)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read agents directory: %w", err)
		}

		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
				continue
			}

			agentPath := filepath.Join(agentsDir, file.Name())
			deployedAgent, err := agent.LoadAgent(agentPath)
			if err != nil {
				continue
			}
			deployedAgents[deployedAgent.Name] = deployedAgent
		}

		responseText += fmt.Sprintf("Found %d deployed agents\n\n", len(deployedAgents))

		// Compare with available agents
		for _, availableAgent := range availableAgents {
			status := "not-deployed"
			deployedVersion := ""

			if deployedAgent, exists := deployedAgents[availableAgent.Name]; exists {
				deployedVersion = deployedAgent.Version
				if deployedAgent.Version == availableAgent.Version {
					status = "up-to-date"
				} else {
					status = "update-available"
				}
			}

			statusInfos = append(statusInfos, AgentStatusInfo{
				Name:             availableAgent.Name,
				DeployedVersion:  deployedVersion,
				AvailableVersion: availableAgent.Version,
				Status:           status,
			})

			statusSymbol := "○"
			switch status {
			case "up-to-date":
				statusSymbol = "✓"
			case "update-available":
				statusSymbol = "⚠"
			}

			versionInfo := ""
			if deployedVersion != "" {
				versionInfo = fmt.Sprintf(" (deployed: v%s, available: v%s)", deployedVersion, availableAgent.Version)
			}

			responseText += fmt.Sprintf("%s %s: %s%s\n", statusSymbol, availableAgent.Name, status, versionInfo)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
	})

	// Register add_location tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "add_location",
		Description: "Add a new deployment location to CAMI's configuration. " +
			"Use this to register a project directory for agent deployment.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddLocationArgs) (*mcp.CallToolResult, any, error) {
		// Validate inputs
		if args.Name == "" {
			return nil, nil, fmt.Errorf("location name is required")
		}
		if args.Path == "" {
			return nil, nil, fmt.Errorf("location path is required")
		}

		// Validate path is absolute
		if !filepath.IsAbs(args.Path) {
			return nil, nil, fmt.Errorf("path must be absolute: %s", args.Path)
		}

		// Validate path exists and is a directory
		info, err := os.Stat(args.Path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil, fmt.Errorf("path does not exist: %s", args.Path)
			}
			return nil, nil, fmt.Errorf("failed to validate path: %w", err)
		}
		if !info.IsDir() {
			return nil, nil, fmt.Errorf("path is not a directory: %s", args.Path)
		}

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Add location
		if err := cfg.AddDeployLocation(args.Name, args.Path); err != nil {
			return nil, nil, fmt.Errorf("failed to add location: %w", err)
		}

		// Save config
		if err := cfg.Save(); err != nil {
			return nil, nil, fmt.Errorf("failed to save config: %w", err)
		}

		responseText := fmt.Sprintf("Added location '%s' at %s\n", args.Name, args.Path)
		responseText += fmt.Sprintf("\nTotal locations: %d", len(cfg.Locations))

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Register list_locations tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_locations",
		Description: "List all configured deployment locations in CAMI. " +
			"Use this to see what project directories are registered for agent deployment.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Convert to response format
		var locationInfos []LocationInfo
		responseText := fmt.Sprintf("Configured locations (%d total):\n\n", len(cfg.Locations))

		if len(cfg.Locations) == 0 {
			responseText += "No locations configured yet. Use add_location to register a project directory.\n"
		} else {
			for _, loc := range cfg.Locations {
				locationInfos = append(locationInfos, LocationInfo{
					Name: loc.Name,
					Path: loc.Path,
				})
				responseText += fmt.Sprintf("• %s\n  %s\n\n", loc.Name, loc.Path)
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &ListLocationsResponse{Locations: locationInfos}, nil
	})

	// Register remove_location tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "remove_location",
		Description: "Remove a deployment location from CAMI's configuration. " +
			"Use this to unregister a project directory.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args RemoveLocationArgs) (*mcp.CallToolResult, any, error) {
		// Validate input
		if args.Name == "" {
			return nil, nil, fmt.Errorf("location name is required")
		}

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Remove location
		if err := cfg.RemoveDeployLocationByName(args.Name); err != nil {
			return nil, nil, fmt.Errorf("failed to remove location: %w", err)
		}

		// Save config
		if err := cfg.Save(); err != nil {
			return nil, nil, fmt.Errorf("failed to save config: %w", err)
		}

		responseText := fmt.Sprintf("Removed location '%s'\n", args.Name)
		responseText += fmt.Sprintf("\nRemaining locations: %d", len(cfg.Locations))

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Register onboard tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "onboard",
		Description: "Get personalized onboarding guidance for CAMI based on current setup state. " +
			"Analyzes configuration, available agents, and deployed agents to provide next steps. " +
			"Use this when user is new to CAMI or asks 'what should I do next?'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		// Check if config exists
		cfg, err := config.Load()
		configExists := err == nil

		state := OnboardingState{
			ConfigExists: configExists,
		}

		var responseText string

		if !configExists {
			// No config - needs to run cami init
			responseText = "# Welcome to CAMI! 🚀\n\n"
			responseText += "CAMI is not yet initialized. Let's get you started!\n\n"
			responseText += "## First Steps\n\n"
			responseText += "1. **Initialize CAMI**\n"
			responseText += "   Run: `./cami init`\n"
			responseText += "   This creates your agent workspace at `vc-agents/my-agents/`\n\n"
			responseText += "2. **Add the Official Agent Library** (optional but recommended)\n"
			responseText += "   After init, you can add 29 professional agents:\n"
			responseText += "   - Use `mcp__cami__add_source` with URL: `git@github.com:lando-labs/lando-agents.git`\n"
			responseText += "   - Or run: `./cami source add git@github.com:lando-labs/lando-agents.git`\n\n"
			responseText += "3. **Deploy Agents**\n"
			responseText += "   Once you have agents available, deploy them to projects\n\n"
			responseText += "Run `./cami init` to get started!\n"

			state.RecommendedNext = "Run ./cami init"

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: responseText},
				},
			}, &OnboardResponse{State: state}, nil
		}

		// Config exists - gather state
		state.SourceCount = len(cfg.AgentSources)
		state.LocationCount = len(cfg.Locations)

		// Check if agent-architect exists
		archPath := filepath.Join(vcAgentsDir, "..", ".claude", "agents", "agent-architect.md")
		if _, err := os.Stat(archPath); err == nil {
			state.HasAgentArch = true
		}

		// Count total available agents
		allAgents, err := agent.LoadAgents(vcAgentsDir)
		if err == nil {
			state.TotalAgents = len(allAgents)
		}

		// Try to count deployed agents in current directory (best effort)
		wd, _ := os.Getwd()
		deployedAgentsPath := filepath.Join(wd, ".claude", "agents")
		if files, err := os.ReadDir(deployedAgentsPath); err == nil {
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
					state.DeployedAgents++
				}
			}
		}

		// Generate personalized guidance based on state
		responseText = "# CAMI Setup Status\n\n"

		// Sources section
		responseText += "## Agent Sources\n"
		if state.SourceCount == 0 {
			responseText += "⚠️ **No agent sources configured**\n\n"
			responseText += "You won't have any agents available until you add a source.\n\n"
			responseText += "**Recommended:** Add the official Lando agent library (29 agents)\n"
			responseText += "- Use `mcp__cami__add_source` with URL: `git@github.com:lando-labs/lando-agents.git`\n"
			responseText += "- Or run: `./cami source add git@github.com:lando-labs/lando-agents.git`\n\n"
			state.RecommendedNext = "Add agent sources"
		} else if state.SourceCount == 1 && state.TotalAgents == 0 {
			responseText += fmt.Sprintf("✓ %d source configured (but no agents found)\n\n", state.SourceCount)
			responseText += "Your workspace is set up, but it's empty.\n\n"
			responseText += "**Options:**\n"
			responseText += "1. Add official agents: `mcp__cami__add_source` with `git@github.com:lando-labs/lando-agents.git`\n"
			responseText += "2. Create custom agents with agent-architect\n\n"
			state.RecommendedNext = "Add agent sources or create agents"
		} else {
			responseText += fmt.Sprintf("✓ %d source(s) configured\n", state.SourceCount)
			responseText += fmt.Sprintf("✓ %d agents available\n\n", state.TotalAgents)
		}

		// Agents section
		if state.TotalAgents > 0 {
			responseText += "## Available Agents\n"
			responseText += fmt.Sprintf("You have access to %d agents across your sources.\n\n", state.TotalAgents)
			responseText += "- Use `mcp__cami__list_agents` to see all available agents\n"
			responseText += "- Use `mcp__cami__deploy_agents` to add agents to projects\n\n"
		}

		// Deployed agents section
		if state.DeployedAgents > 0 {
			responseText += "## Deployed Agents (Current Project)\n"
			responseText += fmt.Sprintf("✓ %d agents deployed in this project\n\n", state.DeployedAgents)
			responseText += "- Use `mcp__cami__scan_deployed_agents` to check for updates\n"
			responseText += "- Use `mcp__cami__update_source` to get latest agent versions\n\n"
		} else if state.TotalAgents > 0 {
			responseText += "## Deployed Agents (Current Project)\n"
			responseText += "⚠️ **No agents deployed in this project yet**\n\n"
			responseText += "Deploy agents to get started:\n"
			responseText += "1. Use `mcp__cami__list_agents` to see available agents\n"
			responseText += "2. Use `mcp__cami__deploy_agents` to add them here\n"
			responseText += "3. Use `mcp__cami__update_claude_md` to document them\n\n"
			if state.RecommendedNext == "" {
				state.RecommendedNext = "Deploy agents to current project"
			}
		}

		// Locations section
		responseText += "## Tracked Locations\n"
		if state.LocationCount == 0 {
			responseText += "No locations tracked yet.\n\n"
			responseText += "Track projects to easily manage agents across multiple codebases:\n"
			responseText += "- Use `mcp__cami__add_location` to register project directories\n"
			responseText += "- Use `mcp__cami__list_locations` to see tracked projects\n\n"
		} else {
			responseText += fmt.Sprintf("✓ %d project(s) tracked\n\n", state.LocationCount)
			responseText += "- Use `mcp__cami__list_locations` to see tracked projects\n"
			responseText += "- Use `mcp__cami__scan_deployed_agents` to check agent status\n\n"
		}

		// Agent architect section
		if state.HasAgentArch {
			responseText += "## Agent Architect\n"
			responseText += "✓ Agent-architect is available\n\n"
			responseText += "Use agent-architect to create custom agents tailored to your needs.\n"
			responseText += "Created agents are saved to `vc-agents/my-agents/` and can be deployed anywhere.\n\n"
		}

		// Next steps
		responseText += "## Quick Commands\n\n"
		responseText += "**List agents:** `mcp__cami__list_agents`\n"
		responseText += "**Deploy agents:** `mcp__cami__deploy_agents`\n"
		responseText += "**Scan current project:** `mcp__cami__scan_deployed_agents`\n"
		responseText += "**Add agent source:** `mcp__cami__add_source`\n"
		responseText += "**Update sources:** `mcp__cami__update_source`\n\n"

		if state.RecommendedNext == "" {
			if state.TotalAgents > 0 && state.DeployedAgents == 0 {
				state.RecommendedNext = "Deploy agents to current project"
			} else if state.TotalAgents > 0 {
				state.RecommendedNext = "Explore and manage your agents"
			} else {
				state.RecommendedNext = "Add agent sources"
			}
		}

		responseText += fmt.Sprintf("**Recommended next step:** %s\n", state.RecommendedNext)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &OnboardResponse{State: state}, nil
	})

	// Register list_sources tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_sources",
		Description: "List all configured agent sources in CAMI. " +
			"Shows source names, paths, priorities, and agent counts. " +
			"Use this to see what agent sources are configured.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Convert to response format
		var sourceInfos []SourceInfo
		responseText := fmt.Sprintf("Agent Sources (%d total):\n\n", len(cfg.AgentSources))

		if len(cfg.AgentSources) == 0 {
			responseText += "No agent sources configured yet.\n\n"
			responseText += "Add a source with: mcp__cami__add_source\n"
			responseText += "Or run: cami source add <git-url>\n"
		} else {
			for _, source := range cfg.AgentSources {
				// Count agents
				agents, err := agent.LoadAgentsFromPath(source.Path)
				agentCount := 0
				if err == nil {
					agentCount = len(agents)
				}

				sourceInfo := SourceInfo{
					Name:       source.Name,
					Path:       source.Path,
					Priority:   source.Priority,
					AgentCount: agentCount,
					GitEnabled: source.Git != nil && source.Git.Enabled,
				}

				if source.Git != nil && source.Git.Enabled {
					sourceInfo.GitRemote = source.Git.Remote
				}

				sourceInfos = append(sourceInfos, sourceInfo)

				responseText += fmt.Sprintf("• %s (priority %d)\n", source.Name, source.Priority)
				responseText += fmt.Sprintf("  Path: %s\n", source.Path)
				responseText += fmt.Sprintf("  Agents: %d\n", agentCount)

				if source.Git != nil && source.Git.Enabled {
					responseText += fmt.Sprintf("  Git: %s\n", source.Git.Remote)
				}

				responseText += "\n"
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, &ListSourcesResponse{Sources: sourceInfos}, nil
	})

	// Register add_source tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "add_source",
		Description: "Add a new agent source by cloning a Git repository. " +
			"The repository will be cloned to vc-agents/<name>/ and added to configuration. " +
			"Use this to add official agent libraries or team/company agent sources.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddSourceArgs) (*mcp.CallToolResult, any, error) {
		// Derive name from URL if not provided
		name := args.Name
		if name == "" {
			// Remove .git suffix
			name = strings.TrimSuffix(args.URL, ".git")
			// Get last part of URL
			parts := strings.Split(name, "/")
			name = parts[len(parts)-1]
			// Remove any : prefix (for SSH URLs)
			if idx := strings.LastIndex(name, ":"); idx != -1 {
				name = name[idx+1:]
			}
		}

		// Default priority for remote sources
		priority := args.Priority
		if priority == 0 {
			priority = 100
		}

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Check if source already exists
		for _, s := range cfg.AgentSources {
			if s.Name == name {
				return nil, nil, fmt.Errorf("source with name %q already exists", name)
			}
		}

		// Find vc-agents directory
		targetPath := filepath.Join(vcAgentsDir, name)

		// Check if directory already exists
		if _, err := os.Stat(targetPath); err == nil {
			return nil, nil, fmt.Errorf("directory already exists: %s", targetPath)
		}

		// Clone repository
		cmd := exec.Command("git", "clone", args.URL, targetPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, string(output))
		}

		// Count agents
		agents, err := agent.LoadAgentsFromPath(targetPath)
		agentCount := 0
		if err == nil {
			agentCount = len(agents)
		}

		// Add to config
		source := config.AgentSource{
			Name:     name,
			Type:     "local",
			Path:     targetPath,
			Priority: priority,
			Git: &config.GitConfig{
				Enabled: true,
				Remote:  args.URL,
			},
		}

		if err := cfg.AddAgentSource(source); err != nil {
			return nil, nil, fmt.Errorf("failed to add source: %w", err)
		}

		if err := cfg.Save(); err != nil {
			return nil, nil, fmt.Errorf("failed to save config: %w", err)
		}

		responseText := fmt.Sprintf("✓ Cloned %s to vc-agents/%s\n", name, name)
		responseText += fmt.Sprintf("✓ Added source with priority %d\n", priority)
		responseText += fmt.Sprintf("✓ Found %d agents\n", agentCount)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Register update_source tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "update_source",
		Description: "Update (git pull) agent sources. " +
			"If no name is specified, updates all sources with git remotes. " +
			"Use this to get the latest agents from configured sources.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args UpdateSourceArgs) (*mcp.CallToolResult, any, error) {
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		var updated, skipped []string
		responseText := ""

		for _, source := range cfg.AgentSources {
			// Skip if specific source requested and this isn't it
			if args.Name != "" && source.Name != args.Name {
				continue
			}

			// Skip if no git remote
			if source.Git == nil || !source.Git.Enabled {
				skipped = append(skipped, source.Name)
				continue
			}

			responseText += fmt.Sprintf("Updating %s...\n", source.Name)

			cmd := exec.Command("git", "-C", source.Path, "pull")
			output, err := cmd.CombinedOutput()

			if err != nil {
				responseText += fmt.Sprintf("  ✗ Failed: %v\n", err)
				continue
			}

			outputStr := string(output)
			if strings.Contains(outputStr, "Already up to date") {
				responseText += "  ✓ Up to date\n"
			} else {
				responseText += "  ✓ Updated\n"
			}

			updated = append(updated, source.Name)
		}

		responseText += "\n"
		if len(updated) > 0 {
			responseText += fmt.Sprintf("Updated: %s\n", strings.Join(updated, ", "))
		}
		if len(skipped) > 0 {
			responseText += fmt.Sprintf("Skipped (no git remote): %s\n", strings.Join(skipped, ", "))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Register source_status tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "source_status",
		Description: "Show git status of agent sources. " +
			"Displays uncommitted changes in source repositories. " +
			"Use this to check if sources have local modifications.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		responseText := "Agent Source Status:\n\n"

		for _, source := range cfg.AgentSources {
			responseText += fmt.Sprintf("• %s\n", source.Name)

			if source.Git == nil || !source.Git.Enabled {
				responseText += "  Git: not enabled\n\n"
				continue
			}

			// Get git status
			cmd := exec.Command("git", "-C", source.Path, "status", "--porcelain")
			output, err := cmd.Output()
			if err != nil {
				responseText += fmt.Sprintf("  Git: error (%v)\n\n", err)
				continue
			}

			if len(output) == 0 {
				responseText += "  Git: ✓ clean\n"
			} else {
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				responseText += fmt.Sprintf("  Git: ⚠ %d uncommitted changes\n", len(lines))
				for i, line := range lines {
					if i >= 3 {
						responseText += fmt.Sprintf("    ... and %d more\n", len(lines)-3)
						break
					}
					responseText += fmt.Sprintf("    %s\n", line)
				}
			}

			responseText += "\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	})

	// Run server with stdio transport
	log.Printf("Starting CAMI MCP server v%s", serverVersion)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
