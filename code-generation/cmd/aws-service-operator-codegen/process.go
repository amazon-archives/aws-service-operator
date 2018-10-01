package main

import (
	"github.com/awslabs/aws-service-operator/code-generation/pkg/codegen"
	"github.com/spf13/cobra"
)

var modelPath, rootPath string

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the operator code based on the models files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		generation := codegen.New(modelPath, rootPath)
		generation.Run()
		return
	},
}

func init() {
	processCmd.Flags().StringVar(&modelPath, "model-path", "models/", "Model path used for regenerating the codebase")
	processCmd.Flags().StringVar(&rootPath, "root-path", "./", "Root path used for regenerating the codebase")
	rootCmd.AddCommand(processCmd)
}
