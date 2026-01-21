package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/finder"
)

// getBasePath returns the base path for module discovery based on cfg.Root
func getBasePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	if cfg.Root == "" {
		return wd, nil
	}

	if filepath.IsAbs(cfg.Root) {
		return cfg.Root, nil
	}

	return filepath.Join(wd, cfg.Root), nil
}

// getModuleType determines the module type based on its path
func getModuleType(path string) string {
	switch {
	case strings.Contains(path, "/"+DirComponents+"/") || strings.Contains(path, "\\"+DirComponents+"\\"):
		return TypeComponent
	case strings.Contains(path, "/"+DirBases+"/") || strings.Contains(path, "\\"+DirBases+"\\"):
		return TypeBase
	case strings.Contains(path, "/"+DirProjects+"/") || strings.Contains(path, "\\"+DirProjects+"\\"):
		return TypeProject
	default:
		return ""
	}
}

// resolveTargetPath resolves the target path based on args and flags
func resolveTargetPath(args []string) (string, error) {
	// Check if both module name and --path are specified
	if len(args) > 0 && pathFlag != "" {
		return "", fmt.Errorf("--path is mutually exclusive with module name argument")
	}

	// Check if neither module name nor --path is specified
	if len(args) == 0 && pathFlag == "" {
		return "", fmt.Errorf("must specify either a module name or --path")
	}

	// If explicit path is provided, use it directly
	if pathFlag != "" {
		return resolveExplicitPath(pathFlag)
	}

	// Use the module name from args
	moduleName := args[0]

	// Search for the module in all directories
	return findModuleInAllDirs(moduleName)
}

// resolveExplicitPath resolves an explicit path (can be relative or absolute)
func resolveExplicitPath(path string) (string, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	}

	return absPath, nil
}

// findModuleInAllDirs searches for a module across all three directories (components, bases, projects)
func findModuleInAllDirs(moduleName string) (string, error) {
	basePath, err := getBasePath()
	if err != nil {
		return "", err
	}

	var allMatches []string

	for _, moduleDir := range ModuleDirs {
		searchPath := filepath.Join(basePath, moduleDir)

		// Skip if directory doesn't exist
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		// Find the module
		matches, err := finder.FindModule(searchPath, moduleName)
		if err != nil {
			return "", fmt.Errorf("failed to search for module in %s: %w", moduleDir, err)
		}

		allMatches = append(allMatches, matches...)
	}

	if len(allMatches) == 0 {
		return "", fmt.Errorf("module '%s' not found in components, bases, or projects", moduleName)
	}

	if len(allMatches) > 1 {
		// Name clash detected across multiple directories
		fmt.Fprintf(os.Stderr, "Error: multiple modules named '%s' found - name clash detected:\n", moduleName)
		for i, match := range allMatches {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, match)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Please use --path to specify the exact path")
		return "", fmt.Errorf("name clash detected")
	}

	return allMatches[0], nil
}

// resolveTargetWithExample resolves the target path, optionally switching to an example directory
func resolveTargetWithExample(args []string, exampleName string) (string, error) {
	modulePath, err := resolveTargetPath(args)
	if err != nil {
		return "", err
	}

	// If no example specified, return the module path
	if exampleName == "" {
		return modulePath, nil
	}

	// Resolve the example path
	examplePath := filepath.Join(modulePath, "examples", exampleName)

	// Check if the example directory exists
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		return "", fmt.Errorf("example '%s' not found in %s/examples", exampleName, modulePath)
	}

	// Check if it contains any .tf file (valid terraform module)
	tfFiles, err := filepath.Glob(filepath.Join(examplePath, "*.tf"))
	if err != nil || len(tfFiles) == 0 {
		return "", fmt.Errorf("example '%s' is not a valid terraform module (no .tf files found)", exampleName)
	}

	return examplePath, nil
}

// readModuleVersion reads the module_version from .spacelift/config.yml
func readModuleVersion(modulePath string) string {
	spaceliftConfig := filepath.Join(modulePath, ".spacelift", "config.yml")
	data, err := os.ReadFile(spaceliftConfig)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module_version:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
