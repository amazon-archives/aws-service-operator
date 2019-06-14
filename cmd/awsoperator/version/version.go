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

package version

import (
	"fmt"

	goversion "github.com/christopherhein/go-version"
	"github.com/spf13/cobra"
)

// New returns a version subcommand
func New(info *goversion.Info) *cobra.Command {
	var shortened bool
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			var response string

			if shortened {
				response = info.ToShortened()
			} else {
				response = info.ToJSON()
			}
			fmt.Printf("%+v", response)
			return
		},
	}
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Use shortened output for version data.")

	return versionCmd
}
