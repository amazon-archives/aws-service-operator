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
	"testing"
)

func TestAllControllersAndResources(t *testing.T) {
	ctrls, resources := AllControllersAndResources()
	if len(ctrls) < 1 {
		t.Errorf("expected '%d' to be greater than 1", len(ctrls))
	}

	if len(resources) < 1 {
		t.Errorf("expected '%d' to be greater than 1", len(resources))
	}
}

func TestAllControllers(t *testing.T) {
	ctrls := AllControllers()
	if len(ctrls) < 1 {
		t.Errorf("expected '%d' to be greater than 1", len(ctrls))
	}
}

func TestAllResources(t *testing.T) {
	resources := AllResources()
	if len(resources) < 1 {
		t.Errorf("expected '%d' to be greater than 1", len(resources))
	}
}
