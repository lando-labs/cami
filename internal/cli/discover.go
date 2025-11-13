package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lando/cami/internal/discovery"
	"github.com/spf13/cobra"
)

// DiscoverOutput represents the JSON output for discover command
type DiscoverOutput struct {
	RootPath string                    `json:"root_path"`
	Count    int                       `json:"count"`
	Projects []*discovery.ProjectInfo  `json:"projects"`
}

// NewDiscoverCommand creates the discover subcommand
func NewDiscoverCommand() *cobra.Command {
	var (
		rootPath     string
		emptyOnly    bool
		hasAgent     string
		maxDepth     int
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover Claude Code projects in a directory tree",
		Long: `Discover all Claude Code projects in a directory tree.

This command recursively scans a directory to find all projects with .claude/ folders,
and reports on their agent deployment status.`,
		Example: `  # Find all Claude Code projects
  cami discover --path ~/Development

  # Find projects with .claude/ but no agents (onboarding opportunities)
  cami discover --path ~/Development --empty-only

  # Find projects with a specific agent deployed
  cami discover --path ~/Development --has-agent frontend

  # Limit search depth
  cami discover --path ~/Development --max-depth 3

  # JSON output for scripting
  cami discover --path ~/Development --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiscover(rootPath, emptyOnly, hasAgent, maxDepth, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&rootPath, "path", "p", ".", "Root directory to search")
	cmd.Flags().BoolVar(&emptyOnly, "empty-only", false, "Only show projects with .claude/ but no agents")
	cmd.Flags().StringVar(&hasAgent, "has-agent", "", "Only show projects with specific agent deployed")
	cmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Maximum directory depth (0 = unlimited)")
	cmd.Flags().StringVar(&outputFormat, "output", "text", "Output format: text or json")

	return cmd
}

func runDiscover(rootPath string, emptyOnly bool, hasAgent string, maxDepth int, outputFormat string) error {
	// Validate output format
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'text' or 'json')", outputFormat)
	}

	// Resolve root path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Run discovery
	opts := discovery.DiscoverOptions{
		RootPath:  absPath,
		EmptyOnly: emptyOnly,
		HasAgent:  hasAgent,
		MaxDepth:  maxDepth,
	}

	projects, err := discovery.DiscoverProjects(opts)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	// Handle no results
	if len(projects) == 0 {
		if emptyOnly {
			fmt.Println("No Claude Code projects found with .claude/ but no agents")
		} else if hasAgent != "" {
			fmt.Printf("No Claude Code projects found with agent: %s\n", hasAgent)
		} else {
			fmt.Println("No Claude Code projects found")
		}
		return nil
	}

	// Output results
	if outputFormat == "json" {
		output := DiscoverOutput{
			RootPath: absPath,
			Count:    len(projects),
			Projects: projects,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		// Text output
		fmt.Printf("Claude Code Projects in %s\n", absPath)
		fmt.Printf("Found %d project(s)\n\n", len(projects))

		for _, project := range projects {
			displayPath := project.Path
			if project.RelativePath != "" {
				displayPath = project.RelativePath
			}

			if project.HasAgents {
				fmt.Printf("✓ %s (%d agents)\n", displayPath, project.AgentCount)

				// List agents
				for _, agent := range project.Agents {
					fmt.Printf("    - %s", agent.Name)
					if agent.Version != "" {
						fmt.Printf(" (v%s)", agent.Version)
					}
					fmt.Println()
				}
			} else {
				fmt.Printf("○ %s (no agents - onboarding opportunity)\n", displayPath)
			}
			fmt.Println()
		}

		// Summary
		withAgents := 0
		withoutAgents := 0
		totalAgents := 0

		for _, p := range projects {
			if p.HasAgents {
				withAgents++
				totalAgents += p.AgentCount
			} else {
				withoutAgents++
			}
		}

		fmt.Println("Summary:")
		fmt.Printf("  Projects with agents: %d\n", withAgents)
		fmt.Printf("  Projects without agents: %d\n", withoutAgents)
		fmt.Printf("  Total agents deployed: %d\n", totalAgents)
	}

	return nil
}
