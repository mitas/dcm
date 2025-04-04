package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/mitas/dcm/internal/model"
)

// ManagedConfig represents the configuration for managed projects
type ManagedConfig struct {
	Projects []model.ManagedProject `yaml:"projects"`
}

// GetDefaultConfigPath returns the default config file path
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configDir := filepath.Join(homeDir, ".config", "dcm")
	return filepath.Join(configDir, "config.yaml")
}

// LoadManagedConfig loads the managed projects configuration from a file
func LoadManagedConfig(configPath string) (*ManagedConfig, error) {
	if configPath == "" {
		configPath = GetDefaultConfigPath()
	}

	// If config file doesn't exist, create default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found at %s, creating default config\n", configPath)
		return createDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config ManagedConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// SaveManagedConfig saves the managed projects configuration to a file
func SaveManagedConfig(config *ManagedConfig, configPath string) error {
	if configPath == "" {
		configPath = GetDefaultConfigPath()
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default config file
func createDefaultConfig(configPath string) (*ManagedConfig, error) {
	config := &ManagedConfig{
		Projects: []model.ManagedProject{},
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}

	// Save empty config
	if err := SaveManagedConfig(config, configPath); err != nil {
		return nil, err
	}

	return config, nil
}
