package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Version            string           `yaml:"version"`
	AgentSources       []AgentSource    `yaml:"agent_sources"`
	Locations          []DeployLocation `yaml:"deploy_locations"`
	DefaultProjectsDir string           `yaml:"default_projects_dir,omitempty"` // Where new projects are created by default
}

// AgentSource represents a source of agents
type AgentSource struct {
	Name     string     `yaml:"name"`
	Type     string     `yaml:"type"` // "local"
	Path     string     `yaml:"path"`
	Priority int        `yaml:"priority"`
	Git      *GitConfig `yaml:"git,omitempty"`
}

// GitConfig holds git-specific configuration
type GitConfig struct {
	Enabled bool   `yaml:"enabled"`
	Remote  string `yaml:"remote,omitempty"`
}

// DeployLocation represents a deployment target
type DeployLocation struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

const (
	configDirName  = "cami-workspace"
	configFileName = "config.yaml"
)

// GetConfigDir returns the path to the config directory
// Respects CAMI_DIR environment variable for custom workspace location
func GetConfigDir() (string, error) {
	// Check for custom CAMI_DIR environment variable first
	if camiDir := os.Getenv("CAMI_DIR"); camiDir != "" {
		return camiDir, nil
	}

	// Fall back to default ~/cami-workspace
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, configDirName), nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Version:      "1",
			AgentSources: []AgentSource{},
			Locations:    []DeployLocation{},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddAgentSource adds a new agent source
func (c *Config) AddAgentSource(source AgentSource) error {
	// Check if source with this name already exists
	for _, s := range c.AgentSources {
		if s.Name == source.Name {
			return fmt.Errorf("source with name %q already exists", source.Name)
		}
	}

	c.AgentSources = append(c.AgentSources, source)
	return nil
}

// RemoveAgentSource removes an agent source by name
func (c *Config) RemoveAgentSource(name string) error {
	for i, source := range c.AgentSources {
		if source.Name == name {
			c.AgentSources = append(c.AgentSources[:i], c.AgentSources[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("source with name %q not found", name)
}

// GetAgentSource retrieves an agent source by name
func (c *Config) GetAgentSource(name string) (*AgentSource, error) {
	for _, source := range c.AgentSources {
		if source.Name == name {
			return &source, nil
		}
	}
	return nil, fmt.Errorf("source with name %q not found", name)
}

// AddDeployLocation adds a new deployment location
func (c *Config) AddDeployLocation(name, path string) error {
	// Check if location already exists
	for _, loc := range c.Locations {
		if loc.Name == name {
			return fmt.Errorf("location with name %q already exists", name)
		}
		if loc.Path == path {
			return fmt.Errorf("location with path %q already exists", path)
		}
	}

	// Validate path exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	c.Locations = append(c.Locations, DeployLocation{
		Name: name,
		Path: path,
	})

	return nil
}

// RemoveDeployLocation removes a deployment location by index
func (c *Config) RemoveDeployLocation(index int) error {
	if index < 0 || index >= len(c.Locations) {
		return fmt.Errorf("invalid index: %d", index)
	}

	c.Locations = append(c.Locations[:index], c.Locations[index+1:]...)
	return nil
}

// RemoveDeployLocationByName removes a deployment location by name
func (c *Config) RemoveDeployLocationByName(name string) error {
	for i, loc := range c.Locations {
		if loc.Name == name {
			c.Locations = append(c.Locations[:i], c.Locations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("location with name %q not found", name)
}

// GetDefaultProjectsDir returns the default directory for creating new projects
// Falls back to ~/projects if not configured
func GetDefaultProjectsDir() (string, error) {
	cfg, err := Load()
	if err == nil && cfg.DefaultProjectsDir != "" {
		return cfg.DefaultProjectsDir, nil
	}

	// Fall back to ~/projects if not configured
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, "projects"), nil
}

