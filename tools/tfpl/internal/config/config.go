package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the .tfpl.yml configuration
type Config struct {
	Root   string `yaml:"root"`
	Binary string `yaml:"binary"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Root:   "",
		Binary: "terraform",
	}
}

// Load attempts to find and load .tfpl.yml starting from startDir and walking up
func Load(startDir string) (Config, error) {
	cfg := DefaultConfig()

	// Walk up the directory tree to find .tfpl.yml
	dir := startDir
	for {
		configPath := filepath.Join(dir, ".tfpl.yml")
		if _, err := os.Stat(configPath); err == nil {
			// Found config file, load it
			data, err := os.ReadFile(configPath)
			if err != nil {
				return cfg, fmt.Errorf("failed to read config file: %w", err)
			}

			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return cfg, fmt.Errorf("failed to parse config file: %w", err)
			}

			// Validate binary
			if cfg.Binary != "terraform" && cfg.Binary != "tofu" {
				return cfg, fmt.Errorf("invalid binary '%s': must be 'terraform' or 'tofu'", cfg.Binary)
			}

			return cfg, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, no config file found
			break
		}
		dir = parent
	}

	// No config file found, return defaults
	return cfg, nil
}
