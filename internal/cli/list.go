package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/spf13/cobra"
)

// AgentInfo represents agent information for JSON output
type AgentInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FilePath    string `json:"file_path,omitempty"`
}

// ListOutput represents the JSON output for list command
type ListOutput struct {
	Count  int         `json:"count"`
	Agents []AgentInfo `json:"agents"`
}

// NewListCommand creates the list subcommand
func NewListCommand(vcAgentsDir string) *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available agents",
		Long:  `List all available agents from configured sources.`,
		Example: `  cami list
  cami list --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(vcAgentsDir, outputFormat)
		},
	}

	cmd.Flags().StringVar(&outputFormat, "output", "text", "Output format: text or json")

	return cmd
}

func runList(vcAgentsDir, outputFormat string) error {
	// Validate output format
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'text' or 'json')", outputFormat)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If no sources configured, fall back to legacy behavior
	var agents []*agent.Agent
	if len(cfg.AgentSources) == 0 {
		agents, err = agent.LoadAgents(vcAgentsDir)
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
		agents, err = agent.LoadAgentsFromSources(sources)
		if err != nil {
			return fmt.Errorf("failed to load agents: %w", err)
		}
	}

	if len(agents) == 0 {
		fmt.Println("No agents found")
		return nil
	}

	// Prepare output
	if outputFormat == "json" {
		output := ListOutput{
			Count:  len(agents),
			Agents: make([]AgentInfo, len(agents)),
		}

		for i, ag := range agents {
			output.Agents[i] = AgentInfo{
				Name:        ag.Name,
				Version:     ag.Version,
				Description: ag.Description,
				Category:    ag.Category,
				FilePath:    ag.FilePath,
			}
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output - group by category
		fmt.Printf("Available Agents (%d):\n\n", len(agents))

		// Group agents by category
		categoryMap := make(map[string][]*agent.Agent)
		for _, ag := range agents {
			category := ag.Category
			if category == "" {
				category = "uncategorized"
			}
			categoryMap[category] = append(categoryMap[category], ag)
		}

		// Display in category order
		categoryOrder := []string{"core", "specialized", "infrastructure", "integration", "design", "meta", "uncategorized"}

		for _, category := range categoryOrder {
			categoryAgents, exists := categoryMap[category]
			if !exists || len(categoryAgents) == 0 {
				continue
			}

			// Capitalize category name for display
			displayCategory := category
			if category != "uncategorized" {
				displayCategory = string(category[0]-32) + category[1:]
			} else {
				displayCategory = "Uncategorized"
			}

			fmt.Printf("## %s (%d agents)\n\n", displayCategory, len(categoryAgents))

			for _, ag := range categoryAgents {
				fmt.Printf("  %s", ag.Name)
				if ag.Version != "" {
					fmt.Printf(" (v%s)", ag.Version)
				}
				fmt.Println()
				if ag.Description != "" {
					fmt.Printf("    %s\n", ag.Description)
				}
				fmt.Println()
			}
		}
	}

	return nil
}
