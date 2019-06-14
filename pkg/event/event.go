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

package event

import (
	"reflect"
	"strings"
)

const (
	tagName = "event"
)

// Message wraps the cloudformation SQS events
type Message struct {
	Type               string `json:"Type"`
	TopicARN           string `json:"TopicArn"`
	Message            string `json:"Message"`
	ResourceProperties ResourceProperties
}

// ResourceProperties will wrap the ResourceProperties object
type ResourceProperties struct {
	Tags []Tag `json:"Tags"`
}

// Tag represents a Tag
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Event allows the event logic
// StackId='arn:aws:cloudformation:us-west-2:915347744415:stack/test-ror/c769c5e0-8c24-11e9-8e2e-02cb67b6aa16'\n
// Timestamp='2019-06-11T08:44:19.385Z'\n
// EventId='WebServerSecurityGroup-CREATE_COMPLETE-2019-06-11T08:44:19.385Z'\n
// LogicalResourceId='WebServerSecurityGroup'\n
// Namespace='915347744415'\n
// PhysicalResourceId='sg-0f3131359f75d68b6'\n
// ResourceProperties='{\"GroupDescription\":\"Enable HTTP access locked down to the load balancer + SSH access\",\"VpcId\":\"vpc-244f4e43\",\"SecurityGroupIngress\":[{\"FromPort\":\"80\",\"ToPort\":\"80\",\"IpProtocol\":\"tcp\",\"SourceSecurityGroupId\":\"sg-b80c64c3\"},{\"CidrIp\":\"0.0.0.0/0\",\"FromPort\":\"22\",\"ToPort\":\"22\",\"IpProtocol\":\"tcp\"}]}'\n
// ResourceStatus='CREATE_COMPLETE'\n
// ResourceStatusReason=''\n
// ResourceType='AWS::EC2::SecurityGroup'\n
// StackName='test-ror'\n
// ClientRequestToken='Console-CreateStack-1b1012af-5ffe-7195-4dc8-9e5709e230b7'\n
type Event struct {
	StackId              string `json:"StackId"`
	Timestamp            string `json:"Timestamp"`
	EventId              string `json:"EventId"`
	LogicalResourceId    string `json:"LogicalResourceId"`
	Namespace            string `json:"Namespace"`
	PhysicalResourceId   string `json:"PhysicalResourceId"`
	ResourceProperties   string `json:"ResourceProperties"`
	ResourceStatus       string `json:"ResourceStatus"`
	ResourceStatusReason string `json:"ResourceStatusReason"`
	ResourceType         string `json:"ResourceType"`
	StackName            string `json:"StackName"`
	ClientRequestToken   string `jwon:"ClientRequestToken"`
}

// Unmarshal will parse the message body
func Unmarshal(message string, obj interface{}) error {
	valueOf := reflect.ValueOf(obj)
	items := strings.Split(message, "\n")
	for _, item := range items {
		x := strings.Split(item, "=")
		key := x[0]
		if key != "" {
			s := x[1]
			s = s[1 : len(s)-1]
			field := valueOf.Elem().FieldByName(key)

			if field.CanSet() {
				field.SetString(s)
			}
		}
	}

	return nil
}
