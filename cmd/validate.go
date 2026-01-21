package cmd

import (
	"github.com/spf13/cobra"
)

// valCmd represents the validate command
var valCmd = &cobra.Command{
	Use:     "val [module-name]",
	Aliases: []string{"validate"},
	Short:   "Run terraform/tofu validate on a component, base, or project",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath, err := resolveTargetPath(args)
		if err != nil {
			return err
		}

		// Run init first if flag is set
		if initFlag {
			if err := runner.RunInit(targetPath); err != nil {
				return err
			}
		}

		return runner.RunValidate(targetPath, argsFlag...)
	},
}

func init() {
	valCmd.Flags().BoolVarP(&initFlag, "init", "i", false, "Run init before the command")
	rootCmd.AddCommand(valCmd)
}
