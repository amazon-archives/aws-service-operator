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

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type mockSQSClient struct {
	sqsiface.SQSAPI
}

func (s *mockSQSClient) GetQueueURL(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	output := &sqs.GetQueueUrlOutput{}
	output.SetQueueUrl("https://suburl/queue/")
	return output, nil
}

func TestRun(t *testing.T) {
	// q := New()
	// q.Run()
}
