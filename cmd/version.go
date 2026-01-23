package cmd

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Version loading strategies for modules.
// Currently supports Spacelift registry versions.
// Future: Add support for other version sources (e.g., git tags, terraform registry).

// spaceliftConfig represents the structure of .spacelift/config.yml
type spaceliftConfig struct {
	ModuleVersion string `yaml:"module_version"`
}

// readModuleVersion reads the module version from available sources.
// Currently reads from .spacelift/config.yml if present.
func readModuleVersion(modulePath string) string {
	// Try Spacelift config first
	if version := readSpaceliftVersion(modulePath); version != "" {
		return version
	}

	// Future: Add other version sources here
	// e.g., git tags, terraform registry, etc.

	return ""
}

// readSpaceliftVersion reads the module_version from .spacelift/config.yml
func readSpaceliftVersion(modulePath string) string {
	configPath := filepath.Join(modulePath, DirSpacelift, FileSpaceliftConfig)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var cfg spaceliftConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}

	return cfg.ModuleVersion
}
