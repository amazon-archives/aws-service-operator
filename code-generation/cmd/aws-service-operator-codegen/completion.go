package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates the bash completion scripts",
	Long: `To load completion run
 . <(aws-service-operator-codegen completion)`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
