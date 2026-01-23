package spacelift

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Directory and file constants for Spacelift configuration
const (
	DirSpacelift = ".spacelift"
	FileConfig   = "config.yml"
)

// config represents the structure of .spacelift/config.yml
type config struct {
	ModuleVersion string `yaml:"module_version"`
}

// ReadModuleVersion reads the module_version from .spacelift/config.yml
// Returns empty string if the file doesn't exist or can't be parsed.
func ReadModuleVersion(modulePath string) string {
	configPath := filepath.Join(modulePath, DirSpacelift, FileConfig)
	data, err := os.ReadFile(configPath) //nolint:gosec // configPath is constructed from known constants
	if err != nil {
		return ""
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return ""
	}

	return cfg.ModuleVersion
}
