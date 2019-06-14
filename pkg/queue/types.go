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

package queue

// ControllerHandler implements the OnMessage handlers which allow the CFN controller
// to push messages to the workqueue
type ControllerHandler interface {
	OnMessage(interface{})
}

// ControllerHandlerFuncs implements the stuck which supports controller handler
type ControllerHandlerFuncs struct {
	MessageFunc func(interface{})
}

// OnMessage is called anytime there is a new message on the queue
func (c ControllerHandlerFuncs) OnMessage(obj interface{}) {
	if c.MessageFunc != nil {
		c.MessageFunc(obj)
	}
}

// Policy wraps the JSON policy
type Policy struct {
	Version   string      `json:"Version"`
	ID        string      `json:"Id"`
	Statement []Statement `json:"Statement"`
}

// Statement defines the QueuePolicy Statement
type Statement struct {
	Sid       string    `json:"Sid"`
	Effect    string    `json:"Effect"`
	Principal string    `json:"Principal"`
	Action    []string  `json:"Action"`
	Resource  string    `json:"Resource"`
	Condition Condition `json:"Condition"`
}

// Condition defines the Condition for Statments
type Condition struct {
	ArnEquals ArnEquals `json:"ArnEquals"`
}

// ArnEquals is a mapping for the SourceArn
type ArnEquals struct {
	AwsSourceArn string `json:"aws:SourceArn"`
}
