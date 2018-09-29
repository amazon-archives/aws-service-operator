package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "aws-operator-codegen",
		Short: "Processes AWS Operator Model files and outputs codegened operators",
		Long:  `TODO WRITE THIS`,
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
