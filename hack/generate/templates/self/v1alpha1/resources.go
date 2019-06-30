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

import "strings"

// AvailableControllers lists all the resources this controller manager handles
var AvailableControllers = map[string][]string{
	<%= available_controllers %>
}

// AllResources returns all possible resources
func AllResources() []string {
	_, resources := AllControllersAndResources()
	return resources
}

// AllControllers returns all possible controllers
func AllControllers() []string {
	ctrls, _ := AllControllersAndResources()
	return ctrls
}

// AllControllersAndResources returns controllers and resources as slices
func AllControllersAndResources() ([]string, []string) {
	controllersResp := []string{}
	resourcesResp := []string{}

	for key, resources := range AvailableControllers {
		controllersResp = append(controllersResp, key)
		for _, resource := range resources {
			resourcesResp = append(resourcesResp, strings.Join([]string{key, resource}, "."))
		}
	}
	return controllersResp, resourcesResp
}
