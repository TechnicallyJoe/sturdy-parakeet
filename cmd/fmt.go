package cmd

import (
	"github.com/spf13/cobra"
)

// fmtCmd represents the fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt [module-name]",
	Short: "Run terraform/tofu fmt on a component, base, or project",
	Args:  cobra.MaximumNArgs(1),
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

		return runner.RunFmt(targetPath, argsFlag...)
	},
}

func init() {
	fmtCmd.Flags().BoolVarP(&initFlag, "init", "i", false, "Run init before the command")
	rootCmd.AddCommand(fmtCmd)
}
