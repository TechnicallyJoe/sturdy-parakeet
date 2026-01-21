package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/finder"
	"github.com/spf13/cobra"
)

// listJsonFlag controls JSON output for list command
var listJsonFlag bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all modules (components, bases, and projects)",
	Long: `List all modules found in components, bases, and projects directories.

Use the --search/-s flag to filter modules using wildcards.
Use the --json flag to output in JSON format for scripting.

Examples:
  tfpl list                    # List all modules
  tfpl list -s storage         # List modules containing "storage"
  tfpl list -s *account*       # List modules with "account" anywhere in the name
  tfpl list --json             # Output as JSON`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVarP(&searchFlag, "search", "s", "", "Filter modules using wildcards (e.g., *storage*)")
	listCmd.Flags().BoolVar(&listJsonFlag, "json", false, "Output in JSON format")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	basePath, err := getBasePath()
	if err != nil {
		return err
	}

	modules, err := collectModules(basePath, searchFlag)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		if listJsonFlag {
			fmt.Println("[]")
			return nil
		}
		if searchFlag != "" {
			fmt.Printf("No modules found matching '%s'\n", searchFlag)
		} else {
			fmt.Println("No modules found")
		}
		return nil
	}

	sortModules(modules)

	if listJsonFlag {
		return printModulesJSON(modules)
	}

	printModules(modules)
	return nil
}

// collectModules discovers all modules across components, bases, and projects directories
func collectModules(basePath, searchFilter string) ([]ModuleInfo, error) {
	var allModules []ModuleInfo

	for _, moduleDir := range ModuleDirs {
		searchPath := filepath.Join(basePath, moduleDir)

		// Skip if directory doesn't exist
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		// List all modules in this directory
		modules, err := finder.ListAllModules(searchPath)
		if err != nil {
			return nil, fmt.Errorf("failed to list modules in %s: %w", moduleDir, err)
		}

		// Process each module
		for name, path := range modules {
			// Apply search filter if specified
			if searchFilter != "" && !finder.MatchesWildcard(name, searchFilter) {
				continue
			}

			// Make path relative to basePath
			relativePath, err := filepath.Rel(basePath, path)
			if err != nil {
				relativePath = path // Fallback to full path if relative fails
			}

			allModules = append(allModules, ModuleInfo{
				Name:    name,
				Type:    getModuleType(path),
				Path:    relativePath,
				Version: readModuleVersion(path),
			})
		}
	}

	return allModules, nil
}

// sortModules sorts modules by type (component, base, project) then alphabetically by name
func sortModules(modules []ModuleInfo) {
	sort.Slice(modules, func(i, j int) bool {
		// First compare by type order
		orderI := ModuleTypeOrder[modules[i].Type]
		orderJ := ModuleTypeOrder[modules[j].Type]
		if orderI != orderJ {
			return orderI < orderJ
		}
		// Then compare by name
		return modules[i].Name < modules[j].Name
	})
}

// printModules outputs the list of modules to stdout
func printModules(modules []ModuleInfo) {
	fmt.Println("Found modules:")

	for _, mod := range modules {
		versionStr := ""
		if mod.Version != "" {
			versionStr = fmt.Sprintf(" (v%s)", mod.Version)
		}
		fmt.Printf("  %-20s [%-9s]  %s%s\n", mod.Name, mod.Type, mod.Path, versionStr)
	}
}

// printModulesJSON outputs the list of modules in JSON format
func printModulesJSON(modules []ModuleInfo) error {
	output, err := json.MarshalIndent(modules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
