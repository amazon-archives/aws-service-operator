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

package token

import (
  uuid "github.com/satori/go.uuid"
)

type Token interface {
  // Generate creates a random string
  Generate() string
}

type token struct {}

// New initialized a new token
func New() token {
  return token{}
}

// Generate will generate a UUID V4 string
func (t token) Generate() string {
  return uuid.NewV4().String()
}