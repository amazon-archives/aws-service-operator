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

package controllerutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
)

// ContainsFinalizer will check if the finalizer exists
func ContainsFinalizer(obj interface{}, finalizer string) (bool, error) {
	metaobj, err := meta.Accessor(obj)
	if err != nil {
		return false, fmt.Errorf("object has no meta: %v", err)
	}
	for _, f := range metaobj.GetFinalizers() {
		if f == finalizer {
			return true, nil
		}
	}
	return false, nil
}

// AddFinalizer will add the finalizer from the list
func AddFinalizer(obj interface{}, finalizer string) error {
	metaobj, err := meta.Accessor(obj)
	if err != nil {
		return fmt.Errorf("object has no meta: %v", err)
	}
	metaobj.SetFinalizers(append(metaobj.GetFinalizers(), finalizer))
	return nil
}

// RemoveFinalizer will remove the finalizer from the list
func RemoveFinalizer(obj interface{}, finalizer string) error {
	metaobj, err := meta.Accessor(obj)
	if err != nil {
		return fmt.Errorf("object has no meta: %v", err)
	}

	var finalizers []string
	for _, f := range metaobj.GetFinalizers() {
		if f == finalizer {
			continue
		}
		finalizers = append(finalizers, f)
	}
	metaobj.SetFinalizers(finalizers)
	return nil
}
