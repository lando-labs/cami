package cli

import (
	"github.com/spf13/cobra"
)

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize CAMI configuration",
		Long: `Initialize CAMI configuration with workspace setup.

This command sets up:
  - Agent sources directory (sources/)
  - Default local source (my-agents/)
  - Configuration file ($CAMI_DIR/config.yaml, defaults to ~/cami-workspace/)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return InitCommand()
		},
	}

	return cmd
}
