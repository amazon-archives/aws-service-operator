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

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	goversion "github.com/christopherhein/go-version"
	"k8s.io/component-base/logs"

	rCmd "awsoperator.io/cmd/awsoperator/root"
	sCmd "awsoperator.io/cmd/awsoperator/server"
	vCmd "awsoperator.io/cmd/awsoperator/version"
)

const (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// cfgFile store the path to the config file
	cfgFile = ".awsoperator.yaml"
)

// main will initialized the config and Execute the cobra commands
func main() {
	rand.Seed(time.Now().UnixNano())

	gversion := goversion.New(version, commit, date)

	rootCmd := rCmd.New()

	serverCmd := sCmd.New(gversion)
	versionCmd := vCmd.New(gversion)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine.Parse([]string{})
}
