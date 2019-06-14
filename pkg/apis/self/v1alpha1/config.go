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

package v1alpha1

import (
	"fmt"
	"strings"
	"time"

	goversion "github.com/christopherhein/go-version"
)

// New will return a Config object initialized with params from the flags or
// config file
func New() Config {
	return Config{
		ClusterName: "aws-service-operator",
		Kubernetes:  ConfigKubernetes{},
		AWS: ConfigAWS{
			SupportedRegions: []string{USWest2Region},
		},
		Server: ConfigServer{
			Address: "0.0.0.0:10355",
			Log: ConfigLog{
				Level: 4,
			},
			Metrics: ConfigMetrics{
				Enable: true,
			},
		},
		Resources: AllResources(),
		Version: &goversion.Info{
			Version: "test",
			Commit:  "test",
			Date:    time.Now().String(),
		},
	}
}

// Validate will check params on the config
func (s Config) Validate() error {
	// Check Resources Length
	if len(s.Resources) == 0 {
		return fmt.Errorf(".resources cannot be empty, please provide a whitelist or '*'")
	}

	// Check for supported regions
	if len(s.AWS.SupportedRegions) == 0 {
		return fmt.Errorf(".aws.supportedRegions cannot be empty, please supply at least 1 region")
	}

	return nil
}

// ResourcesContain will check for specific attributes in the resources key
func (s Config) ResourcesContain(name string) error {
	resourcemapper := map[string]bool{}
	for _, v := range s.Resources {
		if v == "*" {
			return nil
		}

		namer := strings.Split(v, ".")
		if len(namer) == 2 {
			resourcemapper[namer[0]] = true
		}
		resourcemapper[v] = true
	}

	namer := strings.Split(name, ".")
	if len(namer) == 2 {
		if _, ok := resourcemapper[namer[0]]; ok {
			return nil
		}
	}

	if _, ok := resourcemapper[name]; !ok {
		return fmt.Errorf("resource '%s' not in config %+v", name, resourcemapper)
	}

	return nil
}

func (s Config) Complete() Config {
	cfg := s.DeepCopy()
	if err := s.ResourcesContain("*"); err == nil {
		cfg.Resources = AllResources()
	}
	return *cfg
}
