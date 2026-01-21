package cmd

import (
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [module-name]",
	Short: "Run tests on a component, base, or project",
	Long: `Run tests on a component, base, or project using the configured test engine.

The test engine (e.g., terratest, terraform, tofu) is configured in .tfpl.yml under the 'test' section.
By default, terratest is used, which runs 'go test ./...' in the module directory.

Examples:
  tfpl test storage-account                    # Run tests on storage-account module
  tfpl test storage-account -a -v              # Run tests with verbose output
  tfpl test storage-account -a -timeout=30m    # Run tests with custom timeout`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath, err := resolveTargetPath(args)
		if err != nil {
			return err
		}

		return runner.RunTest(targetPath, argsFlag...)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
