package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/deploy"
	"github.com/spf13/cobra"
)

// DeployOutput represents the JSON output format for deploy command
type DeployOutput struct {
	Success   bool         `json:"success"`
	Deployed  []string     `json:"deployed"`
	Failed    []string     `json:"failed"`
	Conflicts []string     `json:"conflicts"`
	Results   []ResultItem `json:"results"`
}

// ResultItem represents a single deployment result
type ResultItem struct {
	Agent   string `json:"agent"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewDeployCommand creates the deploy subcommand
func NewDeployCommand(vcAgentsDir string) *cobra.Command {
	var (
		agentNames   string
		location     string
		overwrite    bool
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy agents to a target project",
		Long: `Deploy one or more agents to a target project location.
Agents are deployed to the .claude/agents directory in the target location.`,
		Example: `  cami deploy --agents frontend,backend --location ~/projects/my-app
  cami deploy -a frontend,backend -l ~/projects/my-app --overwrite
  cami deploy -a frontend,backend -l ~/projects/my-app --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(vcAgentsDir, agentNames, location, overwrite, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&agentNames, "agents", "a", "", "Comma-separated list of agent names (required)")
	cmd.Flags().StringVarP(&location, "location", "l", "", "Target project path (required)")
	cmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "Output format: text or json")

	cmd.MarkFlagRequired("agents")
	cmd.MarkFlagRequired("location")

	return cmd
}

func runDeploy(vcAgentsDir, agentNames, location string, overwrite bool, outputFormat string) error {
	// Validate output format
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'text' or 'json')", outputFormat)
	}

	// Validate target path
	if err := deploy.ValidateTargetPath(location); err != nil {
		return fmt.Errorf("invalid location: %w", err)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Load all available agents from all sources
	var allAgents []*agent.Agent
	if len(cfg.AgentSources) == 0 {
		// Fall back to legacy behavior
		allAgents, err = agent.LoadAgents(vcAgentsDir)
		if err != nil {
			return fmt.Errorf("failed to load agents: %w", err)
		}
	} else {
		// Convert config sources to agent sources
		var sources []agent.AgentSource
		for _, src := range cfg.AgentSources {
			sources = append(sources, agent.AgentSource{
				Path:     src.Path,
				Priority: src.Priority,
			})
		}

		// Load from all sources with priority
		allAgents, err = agent.LoadAgentsFromSources(sources)
		if err != nil {
			return fmt.Errorf("failed to load agents: %w", err)
		}
	}

	// Parse requested agent names
	requestedNames := strings.Split(agentNames, ",")
	for i := range requestedNames {
		requestedNames[i] = strings.TrimSpace(requestedNames[i])
	}

	// Build map of available agents
	agentMap := make(map[string]*agent.Agent)
	for _, ag := range allAgents {
		agentMap[ag.Name] = ag
	}

	// Select agents to deploy
	var agentsToDeploy []*agent.Agent
	var notFound []string

	for _, name := range requestedNames {
		if ag, ok := agentMap[name]; ok {
			agentsToDeploy = append(agentsToDeploy, ag)
		} else {
			notFound = append(notFound, name)
		}
	}

	// Report agents not found
	if len(notFound) > 0 {
		return fmt.Errorf("agents not found: %s", strings.Join(notFound, ", "))
	}

	// Deploy agents
	results, err := deploy.DeployAgents(agentsToDeploy, location, overwrite)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	// Process results
	output := DeployOutput{
		Success:   true,
		Deployed:  []string{},
		Failed:    []string{},
		Conflicts: []string{},
		Results:   []ResultItem{},
	}

	for _, result := range results {
		item := ResultItem{
			Agent:   result.Agent.Name,
			Message: result.Message,
		}

		if result.Success {
			item.Status = "success"
			output.Deployed = append(output.Deployed, result.Agent.Name)
		} else if result.Conflict {
			item.Status = "conflict"
			output.Conflicts = append(output.Conflicts, result.Agent.Name)
			output.Success = false
		} else {
			item.Status = "failed"
			output.Failed = append(output.Failed, result.Agent.Name)
			output.Success = false
		}

		output.Results = append(output.Results, item)
	}

	// Output results
	if outputFormat == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output
		fmt.Printf("Deployment Results:\n\n")
		for _, item := range output.Results {
			statusIcon := "✓"
			if item.Status == "failed" {
				statusIcon = "✗"
			} else if item.Status == "conflict" {
				statusIcon = "⚠"
			}
			fmt.Printf("  %s %s: %s\n", statusIcon, item.Agent, item.Message)
		}

		fmt.Printf("\nSummary:\n")
		fmt.Printf("  Deployed: %d\n", len(output.Deployed))
		if len(output.Conflicts) > 0 {
			fmt.Printf("  Conflicts: %d\n", len(output.Conflicts))
		}
		if len(output.Failed) > 0 {
			fmt.Printf("  Failed: %d\n", len(output.Failed))
		}
	}

	// Return non-zero exit code if deployment was not successful
	if !output.Success {
		os.Exit(1)
	}

	return nil
}
