package cmd

import (
	"fmt"
	"os"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/config"
	"github.com/TechnicallyJoe/sturdy-parakeet/internal/terraform"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var (
	cfg    *config.Config
	runner *terraform.Runner

	// Flags
	pathFlag   string
	initFlag   bool
	argsFlag   []string
	searchFlag string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     "tfpl",
	Short:   "Terraform Polylith CLI tool",
	Version: version,
	Long: `tfpl (Terraform Polylith) is a CLI tool for working with polylith-style Terraform repositories.

It supports running terraform/tofu commands on components, bases, and projects organized
in a polylith structure.`,
	Example: `  tfpl fmt storage-account         # Run fmt on storage-account (searches all types)
  tfpl val k8s-argocd              # Run validate on k8s-argocd
  tfpl val -i k8s-argocd           # Run init then validate on k8s-argocd
  tfpl init k8s-argocd             # Run init on k8s-argocd
  tfpl fmt --path iac/components/azurerm/storage-account  # Run fmt on explicit path
  tfpl init storage-account -a -upgrade -a -reconfigure  # Run init with extra args`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		cfg, err = config.Load(wd)
		if err != nil {
			return err
		}

		// Create terraform runner with config
		runner = terraform.NewRunner(cfg)

		return nil
	},
}

func init() {
	// Add persistent flags
	rootCmd.PersistentFlags().StringVar(&pathFlag, "path", "", "Explicit path (mutually exclusive with module name)")
	rootCmd.PersistentFlags().StringArrayVarP(&argsFlag, "args", "a", []string{}, "Extra arguments to pass to terraform/tofu (can be specified multiple times)")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
