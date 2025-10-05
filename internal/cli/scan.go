package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lando/cami/internal/docs"
	"github.com/spf13/cobra"
)

// ScanOutput represents the JSON output for scan command
type ScanOutput struct {
	Location string      `json:"location"`
	Count    int         `json:"count"`
	Agents   []AgentInfo `json:"agents"`
}

// NewScanCommand creates the scan subcommand
func NewScanCommand() *cobra.Command {
	var (
		location     string
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan deployed agents at a location",
		Long:  `Scan a project location and list all deployed agents in the .claude/agents directory.`,
		Example: `  cami scan --location ~/projects/my-app
  cami scan -l ~/projects/my-app --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScan(location, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&location, "location", "l", "", "Target project path (required)")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "Output format: text or json")

	cmd.MarkFlagRequired("location")

	return cmd
}

func runScan(location, outputFormat string) error {
	// Validate output format
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'text' or 'json')", outputFormat)
	}

	// Scan deployed agents
	agents, err := docs.ScanDeployedAgentsInfo(location)
	if err != nil {
		return fmt.Errorf("failed to scan agents: %w", err)
	}

	if len(agents) == 0 {
		fmt.Printf("No agents found at %s\n", location)
		return nil
	}

	// Prepare output
	if outputFormat == "json" {
		output := ScanOutput{
			Location: location,
			Count:    len(agents),
			Agents:   make([]AgentInfo, len(agents)),
		}

		for i, ag := range agents {
			output.Agents[i] = AgentInfo{
				Name:        ag.Name,
				Version:     ag.Version,
				Description: ag.Description,
			}
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output
		fmt.Printf("Deployed Agents at %s (%d):\n\n", location, len(agents))
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
