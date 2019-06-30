/*
Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
not use this file except in compliance with the License. A copy of the
License is located at

     http://aws.amazon.com/apache2.0/

or in the "license" file accompanying this file. This file is distributed
on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing
permissions and limitations under the License.
*/

package root

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// verbose logging for the operator
	verbose int
)

// New returns a new root command to add subcomands to
func New() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "awsoperator",
		Short: "AWS Service Operator is an operator which manages the lifecycle of AWS resources using Kubernetes.",
		Long: `AWS Service Operator (awsoperator) allows you to use Custom Resource
Definitions (CRDs) to declare the resources your applications consume.
Giving you a single place to model your applications and their dependencies.

Explore the CLI by using:

	$ awsoperator --help

You can read more at https://awsoperator.io`,
		Version: "",
		Example: "awsoperator [command] [subcommand] [flags]",
	}

	rootCmd.PersistentFlags().IntVarP(&verbose, "verbose", "v", 3, "Log level for the CLI")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("logLevel"))

	return rootCmd
}
