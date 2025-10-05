package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

// NewRootCommand creates the root command with all subcommands
func NewRootCommand(vcAgentsDir string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cami",
		Short: "Claude Agent Management Interface",
		Long: `CAMI - Claude Agent Management Interface

A terminal application for managing and deploying Claude Code agents.
When run without arguments, launches the interactive TUI.`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			// This will be handled by main.go to launch TUI
			// If we get here, it means no subcommand was specified
			cmd.Help()
		},
	}

	// Add subcommands
	rootCmd.AddCommand(NewDeployCommand(vcAgentsDir))
	rootCmd.AddCommand(NewUpdateDocsCommand())
	rootCmd.AddCommand(NewListCommand(vcAgentsDir))
	rootCmd.AddCommand(NewScanCommand())
	rootCmd.AddCommand(NewLocationsCommand())
	rootCmd.AddCommand(NewLocationCommand())

	// Set custom version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("CAMI v%s\nClaude Agent Management Interface\n", version))

	return rootCmd
}
