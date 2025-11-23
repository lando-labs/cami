package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/backup"
	"github.com/lando/cami/internal/cli"
	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/deploy"
	"github.com/lando/cami/internal/docs"
	"github.com/lando/cami/internal/manifest"
	"github.com/lando/cami/internal/normalize"
	"github.com/lando/cami/internal/tui"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.3.0"

// MCP server constants
const (
	serverName    = "cami"
	serverVersion = "0.3.0"
)

// ========== CLI FUNCTIONS ==========

func printHelp() {
	fmt.Printf("CAMI v%s - Claude Agent Management Interface\n\n", version)
	fmt.Println("A tool for managing and deploying Claude Code agents.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  cami --mcp               Start MCP server (for Claude Code integration)")
	fmt.Println("  cami list                List available agents")
	fmt.Println("  cami deploy              Deploy agents to a project")
	fmt.Println("  cami scan                Scan deployed agents at a location")
	fmt.Println("  cami update-docs         Update CLAUDE.md with agent info")
	fmt.Println("  cami source              Manage agent sources")
	fmt.Println("  cami locations           Manage deployment locations")
	fmt.Println("  cami -v, --version       Show version information")
	fmt.Println("  cami -h, --help          Show this help message")
	fmt.Println()
	fmt.Println("MCP SERVER:")
	fmt.Println("  The --mcp flag starts CAMI as an MCP server for use with Claude Code.")
	fmt.Println("  This is the primary interface - use natural language with Claude!")
	fmt.Println()
	fmt.Println("  Configure in Claude Code (.mcp.json):")
	fmt.Println("    {")
	fmt.Println("      \"mcpServers\": {")
	fmt.Println("        \"cami\": {")
	fmt.Println("          \"command\": \"cami\",")
	fmt.Println("          \"args\": [\"--mcp\"]")
	fmt.Println("        }")
	fmt.Println("      }")
	fmt.Println("    }")
	fmt.Println()
	fmt.Println("CLI COMMANDS:")
	fmt.Println("  Run 'cami <command> --help' for more information on a specific command.")
	fmt.Println()
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Config file: ~/cami/config.yaml")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/lando-labs/cami")
}

// loadAllAgents loads agents from all configured sources
func loadAllAgents() ([]*agent.Agent, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w - run 'cami source add <git-url>' to add agent sources", err)
	}

	if len(cfg.AgentSources) == 0 {
		return nil, fmt.Errorf("no agent sources configured - run 'cami source add <git-url>' to add agent sources")
	}

	// Convert config sources to agent sources
	agentSources := make([]agent.AgentSource, len(cfg.AgentSources))
	for i, src := range cfg.AgentSources {
		agentSources[i] = agent.AgentSource{
			Path:     src.Path,
			Priority: src.Priority,
		}
	}

	return agent.LoadAgentsFromSources(agentSources)
}

// updateDeploymentManifests updates both project and central manifests after deployment
func updateDeploymentManifests(projectPath string, agents []*agent.Agent, results []*deploy.Result) error {
	// Load config to get source information
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create a map of file paths to source info
	sourceMap := make(map[string]*config.AgentSource)
	for _, src := range cfg.AgentSources {
		sourceMap[src.Path] = &src
	}

	// Only track successfully deployed agents
	var deployedAgents []manifest.DeployedAgent
	now := time.Now()

	for _, result := range results {
		if !result.Success {
			continue
		}

		// Calculate content and metadata hashes
		deployedPath := filepath.Join(projectPath, ".claude", "agents", result.Agent.Name+".md")
		contentHash, err := manifest.CalculateContentHash(deployedPath)
		if err != nil {
			log.Printf("Warning: failed to calculate content hash for %s: %v", result.Agent.Name, err)
			contentHash = ""
		}

		metadataHash, err := manifest.CalculateMetadataHash(deployedPath)
		if err != nil {
			log.Printf("Warning: failed to calculate metadata hash for %s: %v", result.Agent.Name, err)
			metadataHash = ""
		}

		// Find source information from agent's file path
		sourceName := ""
		sourcePath := result.Agent.FilePath
		priority := 999 // Default priority if not found

		for srcPath, src := range sourceMap {
			if strings.HasPrefix(result.Agent.FilePath, srcPath) {
				sourceName = src.Name
				priority = src.Priority
				break
			}
		}

		deployedAgents = append(deployedAgents, manifest.DeployedAgent{
			Name:         result.Agent.Name,
			Version:      result.Agent.Version,
			Source:       sourceName,
			SourcePath:   sourcePath,
			Priority:     priority,
			DeployedAt:   now,
			ContentHash:  contentHash,
			MetadataHash: metadataHash,
		})
	}

	// Create/update project manifest
	projectManifest := &manifest.ProjectManifest{
		Version:      "1",
		State:        manifest.StateCAMINative,
		NormalizedAt: now,
		Agents:       deployedAgents,
	}

	if err := manifest.WriteProjectManifest(projectPath, projectManifest); err != nil {
		return fmt.Errorf("failed to write project manifest: %w", err)
	}

	// Update central manifest
	centralManifest, err := manifest.ReadCentralManifest()
	if err != nil {
		return fmt.Errorf("failed to read central manifest: %w", err)
	}

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	centralManifest.Deployments[absPath] = manifest.ProjectDeployment{
		State:        manifest.StateCAMINative,
		NormalizedAt: now,
		LastScanned:  now,
		Agents:       deployedAgents,
	}

	if err := manifest.WriteCentralManifest(centralManifest); err != nil {
		return fmt.Errorf("failed to write central manifest: %w", err)
	}

	return nil
}

func runTUI() error {
	// Load agents from all configured sources
	agents, err := loadAllAgents()
	if err != nil {
		return fmt.Errorf("error loading agents: %v", err)
	}

	if len(agents) == 0 {
		return fmt.Errorf("no agents found - check your configured agent sources")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	// Create and run TUI
	model := tui.NewModel(agents, cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %v", err)
	}

	return nil
}

func runCLI() {
	// Check for version/help flags first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("CAMI v%s\n", version)
			fmt.Println("Claude Agent Management Interface")
			os.Exit(0)
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		}
	}

	// Create root command (legacy CLI - kept for backward compatibility)
	// TODO: Refactor CLI commands to use config-based loading
	cfg, err := config.Load()
	var vcAgentsDir string
	if err == nil && len(cfg.AgentSources) > 0 {
		vcAgentsDir = cfg.AgentSources[0].Path
	}
	rootCmd := cli.NewRootCommand(vcAgentsDir)

	// Execute CLI command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ========== MCP SERVER FUNCTIONS ==========

// MCP type definitions

type DeployAgentsArgs struct {
	AgentNames []string `json:"agent_names" jsonschema_description:"Array of agent names to deploy (e.g. ['architect', 'backend'])"`
	TargetPath string   `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
	Overwrite  bool     `json:"overwrite,omitempty" jsonschema_description:"Whether to overwrite existing agent files (default: false)"`
}

type UpdateClaudeMdArgs struct {
	TargetPath string `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
}

type ScanDeployedAgentsArgs struct {
	TargetPath string `json:"target_path" jsonschema_description:"Absolute path to target project directory"`
}

type DeployResult struct {
	AgentName string `json:"agent_name"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Conflict  bool   `json:"conflict,omitempty"`
}

type DeployAgentsResponse struct {
	Results []DeployResult `json:"results"`
}

type AgentInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FileName    string `json:"file_name"`
}

type ListAgentsResponse struct {
	Agents []AgentInfo `json:"agents"`
}

type AgentStatusInfo struct {
	Name             string `json:"name"`
	DeployedVersion  string `json:"deployed_version"`
	AvailableVersion string `json:"available_version"`
	Status           string `json:"status"`
}

type ScanDeployedAgentsResponse struct {
	Statuses []AgentStatusInfo `json:"statuses"`
}

type AddLocationArgs struct {
	Name string `json:"name" jsonschema_description:"Friendly name for the location (e.g. 'my-project')"`
	Path string `json:"path" jsonschema_description:"Absolute path to project directory"`
}

type RemoveLocationArgs struct {
	Name string `json:"name" jsonschema_description:"Name of location to remove"`
}

type LocationInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ListLocationsResponse struct {
	Locations []LocationInfo `json:"locations"`
}

type AddSourceArgs struct {
	URL      string `json:"url" jsonschema_description:"Git URL to clone (e.g. 'git@github.com:yourorg/your-agents.git')"`
	Name     string `json:"name,omitempty" jsonschema_description:"Name for the source (derived from URL if not specified)"`
	Priority int    `json:"priority,omitempty" jsonschema_description:"Priority (lower = higher precedence, 1 = highest, default: 50)"`
}

type UpdateSourceArgs struct {
	Name string `json:"name,omitempty" jsonschema_description:"Name of source to update (updates all if not specified)"`
}

type SourceInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Priority     int    `json:"priority"`
	AgentCount   int    `json:"agent_count"`
	GitRemote    string `json:"git_remote,omitempty"`
	GitEnabled   bool   `json:"git_enabled"`
	IsCompliant  bool   `json:"is_compliant"`
	IssueCount   int    `json:"issue_count,omitempty"`
	IssuesSummary string `json:"issues_summary,omitempty"`
}

type ListSourcesResponse struct {
	Sources []SourceInfo `json:"sources"`
}

type CreateProjectArgs struct {
	Name        string   `json:"name" jsonschema_description:"Project name (kebab-case for directory)"`
	Path        string   `json:"path,omitempty" jsonschema_description:"Project directory path (defaults to ~/projects/{name})"`
	Description string   `json:"description" jsonschema_description:"High-level project description (2-3 paragraphs)"`
	AgentNames  []string `json:"agent_names" jsonschema_description:"List of agent names to deploy to the project"`
	VisionDoc   string   `json:"vision_doc,omitempty" jsonschema_description:"Focused CLAUDE.md content (vision, not implementation details)"`
}

type CreateProjectResponse struct {
	ProjectPath    string   `json:"project_path"`
	AgentsDeployed []string `json:"agents_deployed"`
	Success        bool     `json:"success"`
}

type OnboardingState struct {
	ConfigExists    bool   `json:"config_exists"`
	SourceCount     int    `json:"source_count"`
	LocationCount   int    `json:"location_count"`
	HasAgentArch    bool   `json:"has_agent_architect"`
	TotalAgents     int    `json:"total_agents"`
	DeployedAgents  int    `json:"deployed_agents"`
	RecommendedNext string `json:"recommended_next"`
}

type OnboardResponse struct {
	State OnboardingState `json:"state"`
}

func runMCPServer() {
	// Initialize logger to stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check if config exists
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: No config found - some tools may not work until configured")
	} else if len(cfg.AgentSources) == 0 {
		log.Printf("Warning: No agent sources configured - some tools may not work until configured")
	} else {
		log.Printf("Loaded config with %d agent source(s)", len(cfg.AgentSources))
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register all MCP tools
	registerMCPTools(server)

	// Start server with stdio transport
	log.Printf("Starting CAMI MCP server v%s", serverVersion)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerMCPTools(server *mcp.Server) {
	// MCP tool registration - uses config-based loading for all agent operations

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

		// Load all available agents from configured sources
		allAgents, err := loadAllAgents()
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

		// Update manifests to track deployment
		if err := updateDeploymentManifests(args.TargetPath, agentsToDeploy, results); err != nil {
			log.Printf("Warning: failed to update deployment manifests: %v", err)
			// Don't fail the deployment if manifest update fails
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
		// Load all available agents from configured sources
		agents, err := loadAllAgents()
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
		if err := deploy.ValidateTargetPath(args.TargetPath); err != nil {
			return nil, nil, fmt.Errorf("invalid target path: %w", err)
		}

		availableAgents, err := loadAllAgents()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load agents: %w", err)
		}

		agentsDir := filepath.Join(args.TargetPath, ".claude", "agents")
		var statusInfos []AgentStatusInfo
		responseText := fmt.Sprintf("Scanning %s\n\n", args.TargetPath)

		if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
			responseText += "No .claude/agents directory found.\n"
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
			}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
		}

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
			if status == "up-to-date" {
				statusSymbol = "✓"
			} else if status == "update-available" {
				statusSymbol = "⚠"
			}

			versionInfo := ""
			if deployedVersion != "" {
				versionInfo = fmt.Sprintf(" (deployed: v%s, available: v%s)", deployedVersion, availableAgent.Version)
			}

			responseText += fmt.Sprintf("%s %s: %s%s\n", statusSymbol, availableAgent.Name, status, versionInfo)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, &ScanDeployedAgentsResponse{Statuses: statusInfos}, nil
	})

	// Register add_location, list_locations, remove_location tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_location",
		Description: "Add a new deployment location to CAMI's configuration. Use this to register a project directory for agent deployment.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddLocationArgs) (*mcp.CallToolResult, any, error) {
		if args.Name == "" || args.Path == "" {
			return nil, nil, fmt.Errorf("both name and path are required")
		}
		if !filepath.IsAbs(args.Path) {
			return nil, nil, fmt.Errorf("path must be absolute: %s", args.Path)
		}
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

		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.AddDeployLocation(args.Name, args.Path); err != nil {
			return nil, nil, fmt.Errorf("failed to add location: %w", err)
		}

		if err := cfg.Save(); err != nil {
			return nil, nil, fmt.Errorf("failed to save config: %w", err)
		}

		responseText := fmt.Sprintf("Added location '%s' at %s\n\nTotal locations: %d", args.Name, args.Path, len(cfg.Locations))
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: responseText}}}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_locations",
		Description: "List all configured deployment locations in CAMI. Use this to see what project directories are registered for agent deployment.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		var locationInfos []LocationInfo
		responseText := fmt.Sprintf("Configured locations (%d total):\n\n", len(cfg.Locations))

		if len(cfg.Locations) == 0 {
			responseText += "No locations configured yet. Use add_location to register a project directory.\n"
		} else {
			for _, loc := range cfg.Locations {
				locationInfos = append(locationInfos, LocationInfo{Name: loc.Name, Path: loc.Path})
				responseText += fmt.Sprintf("• %s\n  %s\n\n", loc.Name, loc.Path)
			}
		}

		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: responseText}}}, &ListLocationsResponse{Locations: locationInfos}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_location",
		Description: "Remove a deployment location from CAMI's configuration. Use this to unregister a project directory.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args RemoveLocationArgs) (*mcp.CallToolResult, any, error) {
		if args.Name == "" {
			return nil, nil, fmt.Errorf("location name is required")
		}

		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.RemoveDeployLocationByName(args.Name); err != nil {
			return nil, nil, fmt.Errorf("failed to remove location: %w", err)
		}

		if err := cfg.Save(); err != nil {
			return nil, nil, fmt.Errorf("failed to save config: %w", err)
		}

		responseText := fmt.Sprintf("Removed location '%s'\n\nRemaining locations: %d", args.Name, len(cfg.Locations))
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: responseText}}}, nil, nil
	})

	// Register list_sources tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_sources",
		Description: "List all configured agent sources in CAMI. " +
			"Shows source names, paths, priorities, and agent counts. " +
			"Use this to see what agent sources are configured.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		var sourceInfos []SourceInfo
		responseText := fmt.Sprintf("Agent Sources (%d total):\n\n", len(cfg.AgentSources))

		if len(cfg.AgentSources) == 0 {
			responseText += "No agent sources configured yet.\n\n"
			responseText += "Add a source with: mcp__cami__add_source\n"
			responseText += "Or run: cami source add <git-url>\n"
		} else {
			for _, source := range cfg.AgentSources {
				agents, err := agent.LoadAgentsFromPath(source.Path)
				agentCount := 0
				if err == nil {
					agentCount = len(agents)
				}

				// Check compliance status
				analysis, err := normalize.AnalyzeSource(source.Name, source.Path)
				isCompliant := false
				issueCount := 0
				issuesSummary := ""

				if err == nil {
					isCompliant = analysis.IsCompliant
					issueCount = len(analysis.Issues)

					if !isCompliant {
						var issues []string
						if analysis.MissingCAMIIgnore {
							issues = append(issues, "no .camiignore")
						}
						if issueCount > 0 {
							issues = append(issues, fmt.Sprintf("%d agents with issues", issueCount))
						}
						issuesSummary = strings.Join(issues, ", ")
					}
				}

				sourceInfo := SourceInfo{
					Name:          source.Name,
					Path:          source.Path,
					Priority:      source.Priority,
					AgentCount:    agentCount,
					GitEnabled:    source.Git != nil && source.Git.Enabled,
					IsCompliant:   isCompliant,
					IssueCount:    issueCount,
					IssuesSummary: issuesSummary,
				}

				if source.Git != nil && source.Git.Enabled {
					sourceInfo.GitRemote = source.Git.Remote
				}

				sourceInfos = append(sourceInfos, sourceInfo)

				// Format response with compliance status
				complianceIcon := "✓"
				if !isCompliant {
					complianceIcon = "⚠️"
				}

				responseText += fmt.Sprintf("• %s %s (priority %d)\n", complianceIcon, source.Name, source.Priority)
				responseText += fmt.Sprintf("  Path: %s\n", source.Path)
				responseText += fmt.Sprintf("  Agents: %d\n", agentCount)

				if source.Git != nil && source.Git.Enabled {
					responseText += fmt.Sprintf("  Git: %s\n", source.Git.Remote)
				}

				// Show compliance status
				if isCompliant {
					responseText += "  Compliance: ✓ Compliant\n"
				} else {
					responseText += fmt.Sprintf("  Compliance: ⚠️ Issues (%s)\n", issuesSummary)
				}

				responseText += "\n"
			}

			// Add summary at the end
			compliantCount := 0
			for _, info := range sourceInfos {
				if info.IsCompliant {
					compliantCount++
				}
			}

			if compliantCount < len(sourceInfos) {
				responseText += fmt.Sprintf("**Note:** %d/%d sources have compliance issues.\n", len(sourceInfos)-compliantCount, len(sourceInfos))
				responseText += "Use `detect_source_state` and `normalize_source` to fix issues.\n"
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, &ListSourcesResponse{Sources: sourceInfos}, nil
	})

	// Register add_source tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "add_source",
		Description: "Add a new agent source by cloning a Git repository. " +
			"The repository will be cloned to your CAMI workspace sources/ directory and added to configuration. " +
			"Use this to add official agent libraries or team/company agent sources.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddSourceArgs) (*mcp.CallToolResult, any, error) {
		name := args.Name
		if name == "" {
			name = strings.TrimSuffix(args.URL, ".git")
			parts := strings.Split(name, "/")
			name = parts[len(parts)-1]
			if idx := strings.LastIndex(name, ":"); idx != -1 {
				name = name[idx+1:]
			}
		}

		priority := args.Priority
		if priority == 0 {
			priority = 50 // Default: middle priority
		}

		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		for _, s := range cfg.AgentSources {
			if s.Name == name {
				return nil, nil, fmt.Errorf("source with name %q already exists", name)
			}
		}

		// Determine target path using config (respects CAMI_DIR)
		configDir, err := config.GetConfigDir()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get config directory: %w", err)
		}
		sourcesDir := filepath.Join(configDir, "sources")
		if err := os.MkdirAll(sourcesDir, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create sources directory: %w", err)
		}
		targetPath := filepath.Join(sourcesDir, name)

		if _, err := os.Stat(targetPath); err == nil {
			return nil, nil, fmt.Errorf("directory already exists: %s", targetPath)
		}

		cmd := exec.Command("git", "clone", args.URL, targetPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, string(output))
		}

		agents, err := agent.LoadAgentsFromPath(targetPath)
		agentCount := 0
		if err == nil {
			agentCount = len(agents)
		}

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

		responseText := fmt.Sprintf("✓ Cloned %s to %s/sources/%s\n", name, configDir, name)
		responseText += fmt.Sprintf("✓ Added source with priority %d\n", priority)
		responseText += fmt.Sprintf("✓ Found %d agents\n\n", agentCount)

		// Auto-detect compliance
		analysis, err := normalize.AnalyzeSource(name, targetPath)
		if err == nil && !analysis.IsCompliant {
			responseText += "## Source Compliance Check\n\n"
			responseText += "⚠️ **This source has compliance issues:**\n\n"

			if analysis.MissingCAMIIgnore {
				responseText += "- Missing .camiignore file\n"
			}

			issueCount := len(analysis.Issues)
			if issueCount > 0 {
				responseText += fmt.Sprintf("- %d agents with issues:\n", issueCount)
				for _, issue := range analysis.Issues {
					responseText += fmt.Sprintf("  - %s: %s\n", issue.AgentFile, strings.Join(issue.Problems, ", "))
				}
			}

			responseText += "\n**Recommendation:** Use `normalize_source` to fix these issues automatically.\n"
			responseText += "This will add missing versions, descriptions, and create .camiignore.\n"
		} else if err == nil {
			responseText += "✓ **Source is fully compliant!**\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
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
			if args.Name != "" && source.Name != args.Name {
				continue
			}

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
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
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
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, nil, nil
	})

	// Register create_project tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "create_project",
		Description: "Create a new project with proper setup. " +
			"WHEN TO USE: User says 'I want to create/start a new project' or similar. " +
			"WORKFLOW (FOLLOW THIS ORDER): " +
			"1) Use AskUserQuestion to gather: project name, description, tech stack, key features " +
			"2) Use mcp__cami__list_agents to see available agents " +
			"3) Recommend agents based on requirements and get user confirmation " +
			"4) If agents don't exist, use Task tool to invoke agent-architect (parallel if multiple) " +
			"5) Write a focused vision_doc (200-300 words, vision NOT implementation) " +
			"6) Invoke this tool with name, description, agent_names, and vision_doc " +
			"7) Confirm success and guide user to next steps. " +
			"IMPORTANT: NEVER skip steps 1-3. Always gather requirements and confirm agents BEFORE using this tool.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args CreateProjectArgs) (*mcp.CallToolResult, any, error) {
		// Validate project name
		if args.Name == "" {
			return nil, nil, fmt.Errorf("project name is required")
		}

		// Determine project path
		projectPath := args.Path
		if projectPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			projectPath = filepath.Join(homeDir, "projects", args.Name)
		}

		// Expand ~ in path if present
		if strings.HasPrefix(projectPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			projectPath = filepath.Join(homeDir, projectPath[2:])
		}

		// Check if directory already exists
		if _, err := os.Stat(projectPath); err == nil {
			return nil, nil, fmt.Errorf("project directory already exists: %s", projectPath)
		}

		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create project directory: %w", err)
		}

		// Deploy agents
		agentsPath := filepath.Join(projectPath, ".claude", "agents")
		if err := os.MkdirAll(agentsPath, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create agents directory: %w", err)
		}

		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Convert config sources to agent sources for loading
		agentSources := make([]agent.AgentSource, len(cfg.AgentSources))
		for i, src := range cfg.AgentSources {
			agentSources[i] = agent.AgentSource{
				Path:     src.Path,
				Priority: src.Priority,
			}
		}

		allAgents, err := agent.LoadAgentsFromSources(agentSources)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load agents: %w", err)
		}

		// Find agents to deploy
		var agentsToDeploy []*agent.Agent
		var notFound []string
		for _, agentName := range args.AgentNames {
			var found *agent.Agent
			for _, a := range allAgents {
				if a.Name == agentName {
					found = a
					break
				}
			}

			if found == nil {
				notFound = append(notFound, agentName)
			} else {
				agentsToDeploy = append(agentsToDeploy, found)
			}
		}

		if len(notFound) > 0 {
			return nil, nil, fmt.Errorf("agents not found: %v", notFound)
		}

		// Deploy agents using the proper deployment function
		results, err := deploy.DeployAgents(agentsToDeploy, projectPath, false)
		if err != nil {
			return nil, nil, fmt.Errorf("deployment failed: %w", err)
		}

		// Update manifests to track deployment
		if err := updateDeploymentManifests(projectPath, agentsToDeploy, results); err != nil {
			log.Printf("Warning: failed to update deployment manifests: %v", err)
			// Don't fail the deployment if manifest update fails
		}

		// Collect successfully deployed agent names
		deployedAgents := []string{}
		for _, result := range results {
			if result.Success {
				deployedAgents = append(deployedAgents, result.Agent.Name)
			}
		}

		// Write CLAUDE.md if vision doc provided
		if args.VisionDoc != "" {
			claudeMdPath := filepath.Join(projectPath, "CLAUDE.md")
			if err := os.WriteFile(claudeMdPath, []byte(args.VisionDoc), 0644); err != nil {
				return nil, nil, fmt.Errorf("failed to write CLAUDE.md: %w", err)
			}
		}

		// Always update CLAUDE.md with deployed agents (creates file if it doesn't exist)
		if _, err := docs.UpdateCLAUDEmd(projectPath, "Deployed Agents", false); err != nil {
			// Don't fail the whole operation if this fails, just warn
			fmt.Fprintf(os.Stderr, "Warning: failed to update CLAUDE.md with agents: %v\n", err)
		}

		// Register project location
		cfg.Locations = append(cfg.Locations, config.DeployLocation{
			Name: args.Name,
			Path: projectPath,
		})

		if err := cfg.Save(); err != nil {
			// Don't fail the whole operation, just warn
			fmt.Fprintf(os.Stderr, "Warning: failed to save location to config: %v\n", err)
		}

		responseText := fmt.Sprintf("✅ Project Created Successfully!\n\n")
		responseText += fmt.Sprintf("**Project**: %s\n", args.Name)
		responseText += fmt.Sprintf("**Location**: %s\n", projectPath)
		responseText += fmt.Sprintf("**Agents Deployed**: %d\n", len(deployedAgents))
		responseText += "\n**Deployed Agents:**\n"
		for _, name := range deployedAgents {
			responseText += fmt.Sprintf("- %s\n", name)
		}
		responseText += "\n**Next Steps:**\n"
		responseText += fmt.Sprintf("1. Navigate to project: `cd %s`\n", projectPath)
		responseText += "2. Review CLAUDE.md for project vision\n"
		responseText += "3. Start building with your specialized agents!\n"

		return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
			}, &CreateProjectResponse{
				ProjectPath:    projectPath,
				AgentsDeployed: deployedAgents,
				Success:        true,
			}, nil
	})

	// Register onboard tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "onboard",
		Description: "Get personalized onboarding guidance for CAMI based on current setup state. " +
			"WHEN TO USE: User is new to CAMI, asks 'what should I do next?', or seems lost. " +
			"WORKFLOW: " +
			"1) Invoke this tool to analyze current setup state " +
			"2) Review the analysis and recommended next step " +
			"3) If no agent sources: Offer to add sources with mcp__cami__add_source " +
			"4) If sources exist but no deployed agents: Help user create/start a project " +
			"5) If project exists: Guide user to use deployed agents. " +
			"IMPORTANT: Use this proactively when user seems uncertain about next steps.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		cfg, err := config.Load()
		configExists := err == nil

		state := OnboardingState{ConfigExists: configExists}
		var responseText string

		if !configExists {
			responseText = "# Welcome to CAMI! 🚀\n\n"
			responseText += "CAMI is not yet configured. Let me help you get started!\n\n"
			responseText += "## First Step: Add Agent Sources\n\n"
			responseText += "To use CAMI, you need to add an agent source (a Git repository containing agent definitions).\n\n"
			responseText += "**I can help you add an agent source!** Just tell me:\n"
			responseText += "  \"Add agent source from <git-url>\"\n\n"
			responseText += "Or you can do it manually:\n"
			responseText += "- CLI: `cami source add <git-url>`\n"
			responseText += "- MCP: Use `mcp__cami__add_source` with your Git repository URL\n\n"
			responseText += "This will:\n"
			configDir, _ := config.GetConfigDir()
			responseText += fmt.Sprintf("1. Create `%s/` directory for global configuration\n", configDir)
			responseText += fmt.Sprintf("2. Clone the agent repository to `%s/sources/<name>/`\n", configDir)
			responseText += "3. Make agents available across all your projects\n\n"
			responseText += "After adding a source, you can deploy agents to any project with `mcp__cami__deploy_agents`!\n"

			state.RecommendedNext = "Add an agent source"

			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
			}, &OnboardResponse{State: state}, nil
		}

		// Config exists - gather state
		state.SourceCount = len(cfg.AgentSources)
		state.LocationCount = len(cfg.Locations)

		// Count total available agents
		// Convert config sources to agent sources
		agentSources := make([]agent.AgentSource, len(cfg.AgentSources))
		for i, src := range cfg.AgentSources {
			agentSources[i] = agent.AgentSource{
				Path:     src.Path,
				Priority: src.Priority,
			}
		}
		allAgents, err := agent.LoadAgentsFromSources(agentSources)
		if err == nil {
			state.TotalAgents = len(allAgents)
		}

		// Try to count deployed agents in current directory
		wd, _ := os.Getwd()
		deployedAgentsPath := filepath.Join(wd, ".claude", "agents")
		if files, err := os.ReadDir(deployedAgentsPath); err == nil {
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
					state.DeployedAgents++
				}
			}
		}

		// Generate personalized guidance
		responseText = "# CAMI Setup Status\n\n"

		responseText += "## Agent Sources\n"
		if state.SourceCount == 0 {
			responseText += "⚠️ **No agent sources configured**\n\n"
			responseText += "**Recommended:** Add an agent source with `mcp__cami__add_source`\n"
			responseText += "- Provide a Git URL to your agent repository\n\n"
			state.RecommendedNext = "Add agent sources"
		} else if state.TotalAgents == 0 {
			responseText += fmt.Sprintf("✓ %d source(s) configured (but no agents found)\n\n", state.SourceCount)
			state.RecommendedNext = "Add agent sources or create agents"
		} else {
			responseText += fmt.Sprintf("✓ %d source(s) configured\n", state.SourceCount)
			responseText += fmt.Sprintf("✓ %d agents available\n\n", state.TotalAgents)
		}

		if state.TotalAgents > 0 {
			responseText += "## Available Agents\n"
			responseText += fmt.Sprintf("You have access to %d agents.\n\n", state.TotalAgents)
			responseText += "- Use `mcp__cami__list_agents` to see all available agents\n"
			responseText += "- Use `mcp__cami__deploy_agents` to add agents to projects\n\n"
		}

		if state.DeployedAgents > 0 {
			responseText += "## Deployed Agents (Current Project)\n"
			responseText += fmt.Sprintf("✓ %d agents deployed\n\n", state.DeployedAgents)
		} else if state.TotalAgents > 0 {
			responseText += "## Deployed Agents (Current Project)\n"
			responseText += "⚠️ **No agents deployed in this project yet**\n\n"
			if state.RecommendedNext == "" {
				state.RecommendedNext = "Deploy agents to current project"
			}
		}

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
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, &OnboardResponse{State: state}, nil
	})

	// Register detect_source_state tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "detect_source_state",
		Description: "Analyze an agent source for CAMI compliance. " +
			"Checks for missing versions, descriptions, and .camiignore file. " +
			"Use after adding a new source or to audit existing sources.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct {
		SourceName string `json:"source_name"`
	}) (*mcp.CallToolResult, any, error) {
		// Load config to get source path
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Find source
		source, err := cfg.GetAgentSource(args.SourceName)
		if err != nil {
			return nil, nil, fmt.Errorf("source not found: %w", err)
		}

		// Analyze source
		analysis, err := normalize.AnalyzeSource(args.SourceName, source.Path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to analyze source: %w", err)
		}

		// Format response
		responseText := fmt.Sprintf("# Source Analysis: %s\n\n", args.SourceName)
		responseText += fmt.Sprintf("**Path:** %s\n", analysis.Path)
		responseText += fmt.Sprintf("**Agent Count:** %d\n", analysis.AgentCount)
		responseText += fmt.Sprintf("**Compliant:** %v\n\n", analysis.IsCompliant)

		if analysis.MissingCAMIIgnore {
			responseText += "⚠️ Missing .camiignore file\n\n"
		}

		if len(analysis.Issues) > 0 {
			responseText += fmt.Sprintf("## Issues Found (%d)\n\n", len(analysis.Issues))
			for _, issue := range analysis.Issues {
				responseText += fmt.Sprintf("**%s:**\n", issue.AgentFile)
				for _, problem := range issue.Problems {
					responseText += fmt.Sprintf("  - %s\n", problem)
				}
				responseText += "\n"
			}

			responseText += "## Recommended Action\n\n"
			responseText += "Use `normalize_source` to fix these issues automatically.\n"
		} else {
			responseText += "✓ **All agents are compliant!**\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, analysis, nil
	})

	// Register normalize_source tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "normalize_source",
		Description: "Fix source agents to meet CAMI standards. " +
			"Can add missing versions (v1.0.0), description placeholders, and create .camiignore. " +
			"Creates backup before making changes.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct {
		SourceName       string `json:"source_name"`
		AddVersions      bool   `json:"add_versions"`
		AddDescriptions  bool   `json:"add_descriptions"`
		CreateCAMIIgnore bool   `json:"create_camiignore"`
	}) (*mcp.CallToolResult, any, error) {
		// Load config to get source path
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Find source
		source, err := cfg.GetAgentSource(args.SourceName)
		if err != nil {
			return nil, nil, fmt.Errorf("source not found: %w", err)
		}

		// Normalize source
		options := normalize.SourceNormalizationOptions{
			AddVersions:      args.AddVersions,
			AddDescriptions:  args.AddDescriptions,
			CreateCAMIIgnore: args.CreateCAMIIgnore,
		}

		result, err := normalize.NormalizeSource(args.SourceName, source.Path, options)
		if err != nil {
			return nil, nil, fmt.Errorf("normalization failed: %w", err)
		}

		// Format response
		responseText := fmt.Sprintf("# Source Normalization: %s\n\n", args.SourceName)
		responseText += fmt.Sprintf("**Status:** %s\n", func() string {
			if result.Success {
				return "✓ Success"
			}
			return "✗ Failed"
		}())
		responseText += fmt.Sprintf("**Agents Updated:** %d\n", result.AgentsUpdated)
		responseText += fmt.Sprintf("**Backup Created:** %s\n\n", result.BackupPath)

		if len(result.Changes) > 0 {
			responseText += "## Changes Made\n\n"
			for _, change := range result.Changes {
				responseText += fmt.Sprintf("- %s\n", change)
			}
		} else {
			responseText += "No changes were needed.\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, result, nil
	})

	// Register detect_project_state tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "detect_project_state",
		Description: "Analyze a project's normalization state. " +
			"Detects project type (non-cami, cami-aware, cami-legacy, cami-native), " +
			"checks for manifests, and provides normalization recommendations.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct {
		ProjectPath string `json:"project_path"`
	}) (*mcp.CallToolResult, any, error) {
		// Load config to get available sources
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Analyze project
		analysis, err := normalize.AnalyzeProject(args.ProjectPath, cfg.AgentSources)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to analyze project: %w", err)
		}

		// Format response
		responseText := fmt.Sprintf("# Project Analysis\n\n")
		responseText += fmt.Sprintf("**Path:** %s\n", analysis.Path)
		responseText += fmt.Sprintf("**State:** %s\n", analysis.State)
		responseText += fmt.Sprintf("**Has Agents Directory:** %v\n", analysis.HasAgentsDir)
		responseText += fmt.Sprintf("**Has Manifest:** %v\n", analysis.HasManifest)
		responseText += fmt.Sprintf("**Agent Count:** %d\n\n", analysis.AgentCount)

		if len(analysis.Agents) > 0 {
			responseText += "## Deployed Agents\n\n"
			for _, ag := range analysis.Agents {
				responseText += fmt.Sprintf("**%s**", ag.Name)
				if ag.Version != "" {
					responseText += fmt.Sprintf(" (v%s)", ag.Version)
				} else {
					responseText += " (no version)"
				}
				if ag.MatchesSource != "" {
					responseText += fmt.Sprintf(" - matches %s", ag.MatchesSource)
					if ag.NeedsUpgrade {
						responseText += " (update available)"
					}
				} else {
					responseText += " - not in sources"
				}
				responseText += "\n"
			}
			responseText += "\n"
		}

		// Show recommendations
		responseText += "## Recommendations\n\n"
		if analysis.Recommendations.MinimalRequired {
			responseText += "✓ **Minimal normalization required:** Create manifests for tracking\n"
		}
		if analysis.Recommendations.StandardRecommended {
			responseText += "✓ **Standard normalization recommended:** Link agents to sources\n"
		}
		if analysis.Recommendations.FullOptional {
			responseText += "✓ **Full normalization optional:** Rewrite agents with agent-architect\n"
		}

		if !analysis.Recommendations.MinimalRequired && !analysis.Recommendations.StandardRecommended {
			responseText += "✓ **Project is fully normalized!**\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, analysis, nil
	})

	// Register normalize_project tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "normalize_project",
		Description: "Normalize a project by creating manifests and linking agents to sources. " +
			"Supports minimal (just manifests) and standard (manifests + source links) levels. " +
			"Creates backup before making changes.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct {
		ProjectPath string `json:"project_path"`
		Level       string `json:"level"` // "minimal" or "standard"
	}) (*mcp.CallToolResult, any, error) {
		// Load config to get available sources
		cfg, err := config.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load config: %w", err)
		}

		// Parse level
		var level normalize.ProjectNormalizationLevel
		switch args.Level {
		case "minimal":
			level = normalize.LevelMinimal
		case "standard":
			level = normalize.LevelStandard
		case "full":
			level = normalize.LevelFull
		default:
			return nil, nil, fmt.Errorf("invalid level: %s (must be 'minimal', 'standard', or 'full')", args.Level)
		}

		// Normalize project
		options := normalize.ProjectNormalizationOptions{
			Level: level,
		}

		result, err := normalize.NormalizeProject(args.ProjectPath, options, cfg.AgentSources)
		if err != nil {
			return nil, nil, fmt.Errorf("normalization failed: %w", err)
		}

		// Format response
		responseText := fmt.Sprintf("# Project Normalization\n\n")
		responseText += fmt.Sprintf("**Status:** %s\n", func() string {
			if result.Success {
				return "✓ Success"
			}
			return "✗ Failed"
		}())
		responseText += fmt.Sprintf("**State Before:** %s\n", result.StateBefore)
		responseText += fmt.Sprintf("**State After:** %s\n", result.StateAfter)
		responseText += fmt.Sprintf("**Backup Created:** %s\n\n", result.BackupPath)

		if len(result.Changes) > 0 {
			responseText += "## Changes Made\n\n"
			for _, change := range result.Changes {
				responseText += fmt.Sprintf("- %s\n", change)
			}
			responseText += "\n"
		}

		if result.UndoAvailable {
			responseText += "**Undo available:** Use backup.RestoreFromBackup to revert changes\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, result, nil
	})

	// Register cleanup_backups tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "cleanup_backups",
		Description: "Clean up old backup directories, keeping only the N most recent. " +
			"Use when backup count exceeds threshold (10+) or to free up disk space. " +
			"Default keeps 3 most recent backups.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct {
		TargetPath string `json:"target_path"`
		KeepRecent int    `json:"keep_recent"` // Default: 3
	}) (*mcp.CallToolResult, any, error) {
		// Default to keeping 3 recent backups
		keepRecent := args.KeepRecent
		if keepRecent <= 0 {
			keepRecent = backup.DefaultKeepRecent
		}

		// Analyze archive first
		analysis, err := backup.AnalyzeArchive(args.TargetPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to analyze archive: %w", err)
		}

		// Format response
		responseText := fmt.Sprintf("# Backup Cleanup\n\n")
		responseText += fmt.Sprintf("**Total Backups:** %d\n", analysis.TotalBackups)
		responseText += fmt.Sprintf("**Total Size:** %.2f MB\n\n", float64(analysis.TotalSizeBytes)/(1024*1024))

		if analysis.TotalBackups <= keepRecent {
			responseText += fmt.Sprintf("No cleanup needed - only %d backups exist (keeping %d)\n", analysis.TotalBackups, keepRecent)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
			}, analysis, nil
		}

		// Perform cleanup
		result, err := backup.CleanupBackups(args.TargetPath, backup.CleanupOptions{
			KeepRecent: keepRecent,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("cleanup failed: %w", err)
		}

		responseText += fmt.Sprintf("## Cleanup Results\n\n")
		responseText += fmt.Sprintf("**Removed:** %d backups\n", result.RemovedCount)
		responseText += fmt.Sprintf("**Freed:** %.2f MB\n", float64(result.FreedBytes)/(1024*1024))
		responseText += fmt.Sprintf("**Kept:** %d backups\n\n", len(result.KeptBackups))

		if len(result.KeptBackups) > 0 {
			responseText += "**Remaining backups:**\n"
			for _, backupPath := range result.KeptBackups {
				responseText += fmt.Sprintf("- %s\n", filepath.Base(backupPath))
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: responseText}},
		}, result, nil
	})
}

// ========== MAIN ENTRY POINT ==========

func main() {
	// No args - show helpful usage message
	if len(os.Args) == 1 {
		fmt.Println("CAMI - Claude Agent Management Interface")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  cami --mcp         Start MCP server for Claude Code (primary interface)")
		fmt.Println("  cami [command]     Run CLI command (try: cami --help)")
		fmt.Println()
		fmt.Println("Recommended workflow:")
		fmt.Println("  Open your CAMI workspace in Claude Code and interact naturally:")
		fmt.Println("    cd ~/cami")
		fmt.Println("    claude")
		fmt.Println()
		fmt.Println("CLI Examples:")
		fmt.Println("  cami list                   List available agents")
		fmt.Println("  cami source add <git-url>   Add agent source")
		fmt.Println("  cami deploy <agents> <path> Deploy agents")
		fmt.Println("  cami --help                 Show full help")
		os.Exit(0)
	}

	// Check for MCP server mode
	if os.Args[1] == "--mcp" {
		runMCPServer()
		return
	}

	// Otherwise, run CLI mode
	runCLI()
}
