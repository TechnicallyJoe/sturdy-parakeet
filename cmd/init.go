package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [module-name]",
	Short: "Run terraform/tofu init on a component, base, or project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath, err := resolveTargetPath(args)
		if err != nil {
			return err
		}

		return runner.RunInit(targetPath, argsFlag...)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
