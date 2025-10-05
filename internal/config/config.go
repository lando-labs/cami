package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	DeployLocations []DeployLocation `json:"deploy_locations"`
}

// DeployLocation represents a deployment target
type DeployLocation struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const configFileName = ".cami.json"

// GetConfigPath returns the path to the config file in the user's home directory
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
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
			DeployLocations: []DeployLocation{},
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddDeployLocation adds a new deployment location
func (c *Config) AddDeployLocation(name, path string) error {
	// Check if location already exists
	for _, loc := range c.DeployLocations {
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

	c.DeployLocations = append(c.DeployLocations, DeployLocation{
		Name: name,
		Path: path,
	})

	return nil
}

// RemoveDeployLocation removes a deployment location by index
func (c *Config) RemoveDeployLocation(index int) error {
	if index < 0 || index >= len(c.DeployLocations) {
		return fmt.Errorf("invalid index: %d", index)
	}

	c.DeployLocations = append(c.DeployLocations[:index], c.DeployLocations[index+1:]...)
	return nil
}

// RemoveDeployLocationByName removes a deployment location by name
func (c *Config) RemoveDeployLocationByName(name string) error {
	for i, loc := range c.DeployLocations {
		if loc.Name == name {
			c.DeployLocations = append(c.DeployLocations[:i], c.DeployLocations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("location with name %q not found", name)
}
