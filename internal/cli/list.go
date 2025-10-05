package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lando/cami/internal/agent"
	"github.com/spf13/cobra"
)

// AgentInfo represents agent information for JSON output
type AgentInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
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
		Long:  `List all available agents from the vc-agents directory.`,
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

	// Load agents
	agents, err := agent.LoadAgents(vcAgentsDir)
	if err != nil {
		return fmt.Errorf("failed to load agents: %w", err)
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
				FilePath:    ag.FilePath,
			}
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output
		fmt.Printf("Available Agents (%d):\n\n", len(agents))
		for _, ag := range agents {
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

	return nil
}
