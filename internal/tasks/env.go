package tasks

import (
	"os"
	"path/filepath"
)

// Environment variable names for built-in variables
const (
	EnvGitRoot    = "MOTF_GIT_ROOT"
	EnvModulePath = "MOTF_MODULE_PATH"
	EnvModuleName = "MOTF_MODULE_NAME"
	EnvConfigPath = "MOTF_CONFIG_PATH"
	EnvBinary     = "MOTF_BINARY"
)

// EnvBuilder constructs environment variables for task execution.
// It starts with the current process environment and adds MOTF_* built-in variables.
type EnvBuilder struct {
	vars map[string]string
}

// NewEnvBuilder creates a new EnvBuilder with empty built-in variables.
func NewEnvBuilder() *EnvBuilder {
	return &EnvBuilder{
		vars: make(map[string]string),
	}
}

// WithGitRoot sets the MOTF_GIT_ROOT variable.
func (b *EnvBuilder) WithGitRoot(gitRoot string) *EnvBuilder {
	b.vars[EnvGitRoot] = gitRoot
	return b
}

// WithModulePath sets the MOTF_MODULE_PATH variable.
func (b *EnvBuilder) WithModulePath(modulePath string) *EnvBuilder {
	b.vars[EnvModulePath] = modulePath
	return b
}

// WithModuleName sets the MOTF_MODULE_NAME variable.
// If empty, it will be derived from the module path if available.
func (b *EnvBuilder) WithModuleName(moduleName string) *EnvBuilder {
	b.vars[EnvModuleName] = moduleName
	return b
}

// WithConfigPath sets the MOTF_CONFIG_PATH variable.
func (b *EnvBuilder) WithConfigPath(configPath string) *EnvBuilder {
	b.vars[EnvConfigPath] = configPath
	return b
}

// WithBinary sets the MOTF_BINARY variable.
func (b *EnvBuilder) WithBinary(binary string) *EnvBuilder {
	b.vars[EnvBinary] = binary
	return b
}

// Build returns the complete environment for task execution.
// It includes the current process environment plus all MOTF_* built-in variables.
func (b *EnvBuilder) Build() []string {
	// Start with current environment
	env := os.Environ()

	// Add built-in variables
	for key, value := range b.vars {
		env = append(env, key+"="+value)
	}

	return env
}

// ModuleNameFromPath extracts the module name from an absolute module path.
// It returns the last component of the path (e.g., "/path/to/storage-account" -> "storage-account").
func ModuleNameFromPath(modulePath string) string {
	if modulePath == "" {
		return ""
	}
	return filepath.Base(modulePath)
}
