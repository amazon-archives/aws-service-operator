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

import "testing"

func TestValidation(t *testing.T) {
	cfg := New()
	cfg.Resources = []string{}

	if err := cfg.Validate(); err == nil {
		t.Error("error expected when resources blank")
	}
	cfg.Resources = []string{"*"}
	cfg.AWS.SupportedRegions = []string{}

	if err := cfg.Validate(); err == nil {
		t.Error("error expected when supported regions blank")
	}
}

func TestResourcesContain(t *testing.T) {
	cfg := New()
	cfg.Resources = []string{}

	if err := cfg.ResourcesContain("cloudformation.stack"); err == nil {
		t.Error("error expect when resources blank")
	}

	cfg.Resources = []string{"*"}
	if err := cfg.ResourcesContain("cloudformation.stack"); err != nil {
		t.Errorf("unexpected error when resources features wildcard")
	}

	cfg.Resources = []string{"*"}
	if err := cfg.ResourcesContain("*"); err != nil {
		t.Errorf("unexpected error when resources features wildcard")
	}

	cfg.Resources = []string{"cloudformation"}
	if err := cfg.ResourcesContain("cloudformation.stack"); err != nil {
		t.Errorf("unexpected error when resources features grouping '%s'", err.Error())
	}
}

func TestComplete(t *testing.T) {
	cfg := New()
	cfg.Resources = []string{"*"}

	cfg = cfg.Complete()

	if len(cfg.Resources) < 3 {
		t.Errorf("expected more than one resource %d", len(cfg.Resources))
	}
}
