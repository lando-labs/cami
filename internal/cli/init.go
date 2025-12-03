package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lando/cami/internal/config"
)

// InitCommand initializes CAMI configuration
func InitCommand() error {
	fmt.Println("Welcome to CAMI! 🚀")
	fmt.Println()
	fmt.Println("CAMI manages Claude Code agents in a workspace directory.")
	fmt.Println()

	// Check if config already exists
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠ Configuration already exists at: %s\n", configPath)
		fmt.Print("Overwrite? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Initialization cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Ask where to store agents
	fmt.Println("? Where should we store your agents?")
	fmt.Println("  1. ./sources (default, in this directory)")
	fmt.Println("  2. Specify custom path")
	fmt.Print("\nChoice (1): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var sourcesPath string
	if choice == "" || choice == "1" {
		// Default: ./sources
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		sourcesPath = filepath.Join(cwd, "sources")
	} else if choice == "2" {
		fmt.Print("Path: ")
		path, _ := reader.ReadString('\n')
		sourcesPath = strings.TrimSpace(path)

		// Expand ~ to home directory
		if strings.HasPrefix(sourcesPath, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			sourcesPath = filepath.Join(home, sourcesPath[2:])
		}
	} else {
		return fmt.Errorf("invalid choice: %s", choice)
	}

	// Create sources structure
	sourcesPath, err = filepath.Abs(sourcesPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	myAgentsPath := filepath.Join(sourcesPath, "my-agents")
	if err := os.MkdirAll(myAgentsPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create .gitkeep
	gitkeepPath := filepath.Join(myAgentsPath, ".gitkeep")
	if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create .gitkeep: %w", err)
	}

	// Create config with my-agents source
	cfg = &config.Config{
		Version: "1",
		AgentSources: []config.AgentSource{
			{
				Name:     "my-agents",
				Type:     "local",
				Path:     myAgentsPath,
				Priority: 200,
				Git: &config.GitConfig{
					Enabled: false,
				},
			},
		},
		Locations: []config.DeployLocation{},
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath, _ = config.GetConfigPath()

	fmt.Println()
	fmt.Printf("✓ Created %s\n", myAgentsPath)
	fmt.Printf("✓ Configuration saved to %s\n", configPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Ensure MCP is configured in Claude Code settings")
	fmt.Println("  2. Start Claude in this directory: claude")
	fmt.Println("  3. Ask: @claude help me onboard with CAMI")
	fmt.Println()
	fmt.Println("For more info: cat README.md")

	return nil
}
