package cli

import (
	"fmt"

	"github.com/lando/cami/internal/docs"
	"github.com/spf13/cobra"
)

// NewUpdateDocsCommand creates the update-docs subcommand
func NewUpdateDocsCommand() *cobra.Command {
	var (
		location    string
		sectionName string
		dryRun      bool
	)

	cmd := &cobra.Command{
		Use:   "update-docs",
		Short: "Update CLAUDE.md with deployed agent information",
		Long: `Update the CLAUDE.md file in a project with information about deployed agents.
This command scans the .claude/agents directory and generates a documentation section
listing all deployed agents with their metadata.`,
		Example: `  cami update-docs --location ~/projects/my-app
  cami update-docs -l ~/projects/my-app --section "Available Agents"
  cami update-docs -l ~/projects/my-app --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateDocs(location, sectionName, dryRun)
		},
	}

	cmd.Flags().StringVarP(&location, "location", "l", "", "Target project path (required)")
	cmd.Flags().StringVarP(&sectionName, "section", "s", "Deployed Agents", "Section name in CLAUDE.md")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without writing")

	cmd.MarkFlagRequired("location")

	return cmd
}

func runUpdateDocs(location, sectionName string, dryRun bool) error {
	content, err := docs.UpdateCLAUDEmd(location, sectionName, dryRun)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println("Dry run - CLAUDE.md would be updated with the following content:")
		fmt.Println()
		fmt.Println(content)
		return nil
	}

	fmt.Printf("✓ Successfully updated CLAUDE.md at %s\n", location)
	return nil
}
