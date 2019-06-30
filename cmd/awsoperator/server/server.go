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

package server

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	goversion "github.com/christopherhein/go-version"

	self "awsoperator.io/pkg/apis/self/v1alpha1"
	controllerManager "awsoperator.io/pkg/controller-manager"
	"awsoperator.io/pkg/server"

	"awsoperator.io/pkg/signals"
	"github.com/spf13/cobra"

	"k8s.io/klog"
)

var (
	cfgFile, kubeconfig string
	info                *goversion.Info
	serverCmd           = &cobra.Command{
		Use:   "server",
		Short: "AWS Service Operator server runs the operator.",
		Long: `AWS Service Operator server (awsoperator) runs the informer loops
that watch for CRDs and create the CloudFormation stacks.

Explore the CLI by using:

	$ awsoperator server --help

You can read more at https://awsoperator.io`,
		Run: func(_ *cobra.Command, _ []string) {
			cfg, err := getConfig()
			if err != nil {
				klog.Fatalf("%s\n", err.Error())
				os.Exit(1)
			}

			// Validate the config
			err = cfg.Validate()
			if err != nil {

				klog.Fatalf("%s\n", err.Error())
				os.Exit(1)
			}

			cfg.Complete()

			stopCh := signals.SetupSignalHandler()

			httpServer := server.New(cfg)
			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					klog.Fatalf("%s\n", err.Error())
					os.Exit(1)
				}
			}()

			// Run controller manager
			cm := controllerManager.New(cfg, httpServer)
			if err := cm.Run(stopCh); err != nil {
				klog.Fatalf("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
)

// New returns a new server command
func New(ginfo *goversion.Info) *cobra.Command {
	info = ginfo
	return serverCmd
}

func init() {
	serverCmd.Flags().StringVarP(&cfgFile, "config", "f", "", "Config file (default is $HOME/.awsoperator.yaml)")
	serverCmd.MarkFlagRequired("config")
	serverCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path for local development")
}

func getConfig() (config self.Config, err error) {
	configFile, err := os.Open(cfgFile)
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	config = self.New()
	if kubeconfig != "" {
		config.Kubernetes.Kubeconfig = kubeconfig
	}
	config.Version = info
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
