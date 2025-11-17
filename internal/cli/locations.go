package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/lando/cami/internal/config"
	"github.com/spf13/cobra"
)

// LocationsOutput represents the JSON output format for locations command
type LocationsOutput struct {
	Locations []config.DeployLocation `json:"locations"`
	Count     int                     `json:"count"`
}

// NewLocationsCommand creates the locations command with subcommands
func NewLocationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locations",
		Short: "List all configured deployment locations",
		Long: `List all configured deployment locations.
Locations are stored in ~/cami-workspace/config.yaml and can be used for deployment targets.`,
		Example: `  cami locations
  cami locations --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, _ := cmd.Flags().GetString("output")
			return runListLocations(outputFormat)
		},
	}

	cmd.Flags().String("output", "text", "Output format: text or json")

	return cmd
}

// NewLocationCommand creates the location command with add/remove subcommands
func NewLocationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "Manage deployment locations",
		Long: `Manage deployment locations for agent deployment.
Use subcommands to add or remove deployment locations.`,
	}

	cmd.AddCommand(NewLocationAddCommand())
	cmd.AddCommand(NewLocationRemoveCommand())

	return cmd
}

// NewLocationAddCommand creates the location add subcommand
func NewLocationAddCommand() *cobra.Command {
	var (
		name string
		path string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new deployment location",
		Long: `Add a new deployment location to the configuration.
The location name must be unique and the path must exist.`,
		Example: `  cami location add --name my-project --path /path/to/project
  cami location add -n my-project -p ~/projects/my-app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddLocation(name, path)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Unique name for the location (required)")
	cmd.Flags().StringVarP(&path, "path", "p", "", "Absolute path to the project directory (required)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("path")

	return cmd
}

// NewLocationRemoveCommand creates the location remove subcommand
func NewLocationRemoveCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a deployment location",
		Long: `Remove a deployment location from the configuration.
The location is identified by its name.`,
		Example: `  cami location remove --name my-project
  cami location remove -n my-project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemoveLocation(name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the location to remove (required)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runListLocations(outputFormat string) error {
	// Validate output format
	if outputFormat != "text" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'text' or 'json')", outputFormat)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Prepare output
	output := LocationsOutput{
		Locations: cfg.Locations,
		Count:     len(cfg.Locations),
	}

	if outputFormat == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
		return nil
	}

	// Text output
	if len(cfg.Locations) == 0 {
		fmt.Println("No deployment locations configured.")
		fmt.Println("\nTo add a location:")
		fmt.Println("  cami location add --name <name> --path <path>")
		return nil
	}

	fmt.Printf("Configured Deployment Locations (%d):\n\n", len(cfg.Locations))

	// Use tabwriter for aligned columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tPATH")
	fmt.Fprintln(w, "----\t----")

	for _, loc := range cfg.Locations {
		fmt.Fprintf(w, "%s\t%s\n", loc.Name, loc.Path)
	}

	w.Flush()

	return nil
}

func runAddLocation(name, path string) error {
	// Validate name is not empty (should be caught by required flag, but defensive)
	if name == "" {
		return fmt.Errorf("location name cannot be empty")
	}

	// Validate and normalize path
	if path == "" {
		return fmt.Errorf("location path cannot be empty")
	}

	// Expand home directory if present
	if path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Validate path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", absPath)
		}
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Add location (this validates uniqueness)
	if err := cfg.AddDeployLocation(name, absPath); err != nil {
		return err
	}

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Successfully added location '%s' -> %s\n", name, absPath)

	return nil
}

func runRemoveLocation(name string) error {
	// Validate name is not empty (should be caught by required flag, but defensive)
	if name == "" {
		return fmt.Errorf("location name cannot be empty")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if any locations exist
	if len(cfg.Locations) == 0 {
		return fmt.Errorf("no deployment locations configured")
	}

	// Remove location by name
	if err := cfg.RemoveDeployLocationByName(name); err != nil {
		return err
	}

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Successfully removed location '%s'\n", name)

	return nil
}
