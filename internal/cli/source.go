package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/spf13/cobra"
)

// NewSourceCommand creates the source management command
func NewSourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage agent sources",
		Long: `Manage agent sources (remote repositories or local directories).

Agent sources are where CAMI discovers agents. You can have multiple sources
with different priorities. Lower priority numbers override higher ones
when agent names conflict (1 = highest priority, 100 = lowest).`,
	}

	cmd.AddCommand(NewSourceAddCommand())
	cmd.AddCommand(NewSourceListCommand())
	cmd.AddCommand(NewSourceUpdateCommand())
	cmd.AddCommand(NewSourceStatusCommand())
	cmd.AddCommand(NewSourceRemoveCommand())

	return cmd
}

// NewSourceAddCommand creates the source add command
func NewSourceAddCommand() *cobra.Command {
	var priority int
	var name string

	cmd := &cobra.Command{
		Use:   "add <git-url>",
		Short: "Add a new agent source",
		Long: `Add a new agent source by cloning a Git repository.

The repository will be cloned to vc-agents/<name>/ and added to your configuration.

Examples:
  cami source add git@github.com:company/agents.git
  cami source add git@github.com:yourorg/team-agents.git --name official --priority 10
  cami source add git@github.com:mycompany/custom-agents.git --priority 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			// Derive name from URL if not provided
			if name == "" {
				name = deriveNameFromURL(url)
			}

			// Default priority for remote sources (50 = middle priority)
			if priority == 0 {
				priority = 50
			}

			return SourceAddCommand(url, name, priority)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name for the source (derived from URL if not specified)")
	cmd.Flags().IntVarP(&priority, "priority", "p", 0, "Priority (lower = higher precedence, default: 50)")

	return cmd
}

// NewSourceListCommand creates the source list command
func NewSourceListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all agent sources",
		Long:  `List all configured agent sources with their priorities and agent counts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return SourceListCommand()
		},
	}

	return cmd
}

// NewSourceUpdateCommand creates the source update command
func NewSourceUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Update agent sources",
		Long: `Update (git pull) agent sources.

If no name is specified, updates all sources with git remotes.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceName := ""
			if len(args) > 0 {
				sourceName = args[0]
			}
			return SourceUpdateCommand(sourceName)
		},
	}

	return cmd
}

// NewSourceStatusCommand creates the source status command
func NewSourceStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show git status of agent sources",
		Long:  `Show git status for all agent sources with git remotes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return SourceStatusCommand()
		},
	}

	return cmd
}

// NewSourceRemoveCommand creates the source remove command
func NewSourceRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an agent source",
		Long: `Remove an agent source from configuration.

Note: This only removes the source from configuration, it does not delete
the directory. Use 'rm -rf vc-agents/<name>' to delete the files.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return SourceRemoveCommand(args[0])
		},
	}

	return cmd
}

// SourceAddCommand adds a new agent source
func SourceAddCommand(url, name string, priority int) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if source already exists
	for _, s := range cfg.AgentSources {
		if s.Name == name {
			return fmt.Errorf("source with name %q already exists", name)
		}
	}

	// Find vc-agents directory
	vcAgentsDir, err := findVCAgentsDir()
	if err != nil {
		return fmt.Errorf("failed to find vc-agents directory: %w", err)
	}

	targetPath := filepath.Join(vcAgentsDir, name)

	// Check if directory already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("directory already exists: %s", targetPath)
	}

	fmt.Printf("Cloning %s to vc-agents/%s...\n", url, name)

	// Clone repository
	cmd := exec.Command("git", "clone", url, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Count agents
	agents, err := agent.LoadAgentsFromPath(targetPath)
	if err != nil {
		fmt.Printf("Warning: failed to load agents: %v\n", err)
	}

	// Add to config
	source := config.AgentSource{
		Name:     name,
		Type:     "local",
		Path:     targetPath,
		Priority: priority,
		Git: &config.GitConfig{
			Enabled: true,
			Remote:  url,
		},
	}

	if err := cfg.AddAgentSource(source); err != nil {
		return fmt.Errorf("failed to add source: %w", err)
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\n✓ Cloned %s to vc-agents/%s\n", name, name)
	fmt.Printf("✓ Added source with priority %d\n", priority)
	if agents != nil {
		fmt.Printf("✓ Found %d agents\n", len(agents))
	}

	return nil
}

// SourceListCommand lists all agent sources
func SourceListCommand() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.AgentSources) == 0 {
		fmt.Println("No agent sources configured.")
		fmt.Println()
		fmt.Println("Add a source with: cami source add <git-url>")
		return nil
	}

	fmt.Println("Agent Sources:")
	fmt.Println()

	for _, source := range cfg.AgentSources {
		fmt.Printf("  %s (priority %d)\n", source.Name, source.Priority)
		fmt.Printf("    Path: %s\n", source.Path)

		// Count agents
		agents, err := agent.LoadAgentsFromPath(source.Path)
		if err == nil {
			fmt.Printf("    Agents: %d\n", len(agents))
		} else {
			fmt.Printf("    Agents: error loading (%v)\n", err)
		}

		if source.Git != nil && source.Git.Enabled {
			fmt.Printf("    Git: %s\n", source.Git.Remote)
		}

		fmt.Println()
	}

	return nil
}

// SourceUpdateCommand updates agent sources
func SourceUpdateCommand(sourceName string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var updated, skipped []string

	for _, source := range cfg.AgentSources {
		// Skip if specific source requested and this isn't it
		if sourceName != "" && source.Name != sourceName {
			continue
		}

		// Skip if no git remote
		if source.Git == nil || !source.Git.Enabled {
			skipped = append(skipped, source.Name)
			continue
		}

		fmt.Printf("Updating %s...\n", source.Name)

		cmd := exec.Command("git", "-C", source.Path, "pull")
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n", err)
			continue
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "Already up to date") {
			fmt.Printf("  ✓ Up to date\n")
		} else {
			fmt.Printf("  ✓ Updated\n")
		}

		updated = append(updated, source.Name)
	}

	fmt.Println()
	if len(updated) > 0 {
		fmt.Printf("Updated: %s\n", strings.Join(updated, ", "))
	}
	if len(skipped) > 0 {
		fmt.Printf("Skipped (no git remote): %s\n", strings.Join(skipped, ", "))
	}

	return nil
}

// SourceStatusCommand shows git status for sources
func SourceStatusCommand() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Agent Source Status:")
	fmt.Println()

	for _, source := range cfg.AgentSources {
		fmt.Printf("  %s\n", source.Name)

		if source.Git == nil || !source.Git.Enabled {
			fmt.Println("    Git: not enabled")
			fmt.Println()
			continue
		}

		// Get git status
		cmd := exec.Command("git", "-C", source.Path, "status", "--porcelain")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("    Git: error (%v)\n", err)
			fmt.Println()
			continue
		}

		if len(output) == 0 {
			fmt.Println("    Git: ✓ clean")
		} else {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			fmt.Printf("    Git: ⚠ %d uncommitted changes\n", len(lines))
			for i, line := range lines {
				if i >= 3 {
					fmt.Printf("      ... and %d more\n", len(lines)-3)
					break
				}
				fmt.Printf("      %s\n", line)
			}
		}

		fmt.Println()
	}

	return nil
}

// SourceRemoveCommand removes an agent source
func SourceRemoveCommand(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find source
	var found bool
	for _, s := range cfg.AgentSources {
		if s.Name == name {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("source %q not found", name)
	}

	// Remove from config
	if err := cfg.RemoveAgentSource(name); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Removed source %q from configuration\n", name)
	fmt.Println()
	fmt.Printf("Note: Directory vc-agents/%s still exists.\n", name)
	fmt.Printf("Remove it manually with: rm -rf vc-agents/%s\n", name)

	return nil
}

// Helper functions

func deriveNameFromURL(url string) string {
	// Remove .git suffix
	name := strings.TrimSuffix(url, ".git")

	// Get last part of URL
	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]

	// Remove any : prefix (for SSH URLs)
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		name = name[idx+1:]
	}

	return name
}

func findVCAgentsDir() (string, error) {
	// Try current directory first
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	vcAgentsPath := filepath.Join(cwd, "vc-agents")
	if _, err := os.Stat(vcAgentsPath); err == nil {
		return vcAgentsPath, nil
	}

	return "", fmt.Errorf("vc-agents directory not found (run from CAMI directory or run 'cami init' first)")
}
