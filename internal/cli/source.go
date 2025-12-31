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
	cmd.AddCommand(NewSourceReconcileCommand())

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

The repository will be cloned to sources/<name>/ and added to your configuration.

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
the directory. Use 'rm -rf sources/<name>' to delete the files.`,
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

	// Find sources directory
	sourcesDir, err := findSourcesDir()
	if err != nil {
		return fmt.Errorf("failed to find sources directory: %w", err)
	}

	targetPath := filepath.Join(sourcesDir, name)

	// Check if directory already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("directory already exists: %s", targetPath)
	}

	fmt.Printf("Cloning %s to sources/%s...\n", url, name)

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

	fmt.Printf("\n✓ Cloned %s to sources/%s\n", name, name)
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
	fmt.Printf("Note: Directory sources/%s still exists.\n", name)
	fmt.Printf("Remove it manually with: rm -rf sources/%s\n", name)

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

func findSourcesDir() (string, error) {
	// Try current directory first
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	sourcesPath := filepath.Join(cwd, "sources")
	if _, err := os.Stat(sourcesPath); err == nil {
		return sourcesPath, nil
	}

	return "", fmt.Errorf("sources directory not found (run from CAMI workspace or run 'cami init' first)")
}

// NewSourceReconcileCommand creates the source reconcile command
func NewSourceReconcileCommand() *cobra.Command {
	var checkOnly bool
	var quiet bool
	var autoAdd bool

	cmd := &cobra.Command{
		Use:   "reconcile",
		Short: "Detect and fix untracked sources",
		Long: `Scan the sources directory and compare to config.yaml.

Detects sources that exist on disk but aren't tracked in configuration.
This can happen when sources are manually cloned or when config save fails.

Use --check-only to just report issues without prompting to fix.
Use --quiet to suppress output when no issues found (useful for hooks).
Use --auto-add to automatically add untracked sources without prompting.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return SourceReconcileCommand(checkOnly, quiet, autoAdd)
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for issues, don't offer to fix")
	cmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress output when no issues found")
	cmd.Flags().BoolVar(&autoAdd, "auto-add", false, "Automatically add untracked sources")

	return cmd
}

// ReconcileResult holds the result of a reconcile operation
type ReconcileResult struct {
	UntrackedSources []UntrackedSource
	OrphanedConfigs  []string // Sources in config but not on disk
	TotalOnDisk      int
	TotalInConfig    int
}

// UntrackedSource represents a source directory not in config
type UntrackedSource struct {
	Name       string
	Path       string
	AgentCount int
	HasGit     bool
	GitRemote  string
}

// SourceReconcileCommand reconciles sources directory with config
func SourceReconcileCommand(checkOnly, quiet, autoAdd bool) error {
	result, err := ReconcileSources()
	if err != nil {
		return err
	}

	hasIssues := len(result.UntrackedSources) > 0 || len(result.OrphanedConfigs) > 0

	// Quiet mode: exit silently if no issues
	if quiet && !hasIssues {
		return nil
	}

	// Report findings
	if !hasIssues {
		fmt.Println("✓ All sources are properly tracked")
		fmt.Printf("  %d sources on disk, %d in config\n", result.TotalOnDisk, result.TotalInConfig)
		return nil
	}

	// Report untracked sources
	if len(result.UntrackedSources) > 0 {
		fmt.Printf("⚠ Found %d untracked source(s):\n\n", len(result.UntrackedSources))
		for _, src := range result.UntrackedSources {
			fmt.Printf("  • %s\n", src.Name)
			fmt.Printf("    Path: %s\n", src.Path)
			fmt.Printf("    Agents: %d\n", src.AgentCount)
			if src.HasGit {
				fmt.Printf("    Git: %s\n", src.GitRemote)
			}
			fmt.Println()
		}
	}

	// Report orphaned configs
	if len(result.OrphanedConfigs) > 0 {
		fmt.Printf("⚠ Found %d orphaned config(s) (in config but not on disk):\n\n", len(result.OrphanedConfigs))
		for _, name := range result.OrphanedConfigs {
			fmt.Printf("  • %s\n", name)
		}
		fmt.Println()
	}

	// Check-only mode: just report
	if checkOnly {
		return nil
	}

	// Auto-add mode or prompt
	if len(result.UntrackedSources) > 0 {
		if autoAdd {
			return addUntrackedSources(result.UntrackedSources)
		}

		// Prompt user
		fmt.Print("Add untracked sources to config? (y/N): ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return addUntrackedSources(result.UntrackedSources)
		}
	}

	return nil
}

// ReconcileSources compares sources directory with config
func ReconcileSources() (*ReconcileResult, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	sourcesDir, err := findSourcesDir()
	if err != nil {
		return nil, err
	}

	result := &ReconcileResult{
		TotalInConfig: len(cfg.AgentSources),
	}

	// Build set of configured source paths
	configuredPaths := make(map[string]string) // path -> name
	configuredNames := make(map[string]bool)
	for _, src := range cfg.AgentSources {
		// Normalize path for comparison
		absPath, _ := filepath.Abs(src.Path)
		configuredPaths[absPath] = src.Name
		configuredNames[src.Name] = true
	}

	// Scan sources directory
	entries, err := os.ReadDir(sourcesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read sources directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		result.TotalOnDisk++

		dirPath := filepath.Join(sourcesDir, entry.Name())
		absPath, _ := filepath.Abs(dirPath)

		// Check if this directory is in config
		if _, exists := configuredPaths[absPath]; exists {
			continue
		}

		// Also check by name in case paths differ slightly
		if configuredNames[entry.Name()] {
			continue
		}

		// This is an untracked source
		untracked := UntrackedSource{
			Name: entry.Name(),
			Path: dirPath,
		}

		// Count agents
		agents, err := agent.LoadAgentsFromPath(dirPath)
		if err == nil {
			untracked.AgentCount = len(agents)
		}

		// Check for git
		gitDir := filepath.Join(dirPath, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			untracked.HasGit = true
			// Try to get remote URL
			cmd := exec.Command("git", "-C", dirPath, "remote", "get-url", "origin")
			if output, err := cmd.Output(); err == nil {
				untracked.GitRemote = strings.TrimSpace(string(output))
			}
		}

		result.UntrackedSources = append(result.UntrackedSources, untracked)
	}

	// Check for orphaned configs (in config but not on disk)
	for _, src := range cfg.AgentSources {
		if _, err := os.Stat(src.Path); os.IsNotExist(err) {
			result.OrphanedConfigs = append(result.OrphanedConfigs, src.Name)
		}
	}

	return result, nil
}

// addUntrackedSources adds untracked sources to config
func addUntrackedSources(sources []UntrackedSource) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	for _, src := range sources {
		// Determine priority (default 50 for discovered sources)
		priority := 50

		newSource := config.AgentSource{
			Name:     src.Name,
			Type:     "local",
			Path:     src.Path,
			Priority: priority,
		}

		if src.HasGit && src.GitRemote != "" {
			newSource.Git = &config.GitConfig{
				Enabled: true,
				Remote:  src.GitRemote,
			}
		} else {
			newSource.Git = &config.GitConfig{
				Enabled: false,
			}
		}

		if err := cfg.AddAgentSource(newSource); err != nil {
			fmt.Printf("  ✗ Failed to add %s: %v\n", src.Name, err)
			continue
		}

		fmt.Printf("  ✓ Added %s (priority %d, %d agents)\n", src.Name, priority, src.AgentCount)
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("\n✓ Config updated")
	return nil
}
