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
	"encoding/json"
	"fmt"

	"k8s.io/klog"

	"awsoperator.io/pkg/event"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// InformerFactory implements all the handler and starting for watching sqs
type InformerFactory interface {
	AddMessageHander(ControllerHandler)
	Start(<-chan struct{})
	GetTopicARN() string
}

// Informer implements the controller
type Informer struct {
	sqsclient sqsiface.SQSAPI
	snsclient snsiface.SNSAPI

	queueName string
	queueURL  string
	queueARN  string

	topicARN string
	subARN   string

	handlers ControllerHandler
}

// New will return a queue controller and set up any resource if they do not exist
func New(
	sqsclient sqsiface.SQSAPI,
	snsclient snsiface.SNSAPI,
	name string) (*Informer, error) {

	ctrl := &Informer{
		sqsclient: sqsclient,
		snsclient: snsclient,
		queueName: name,
	}

	queueURL, queueARN, err := ctrl.createSQSQueue(name)
	if err != nil {
		return ctrl, err
	}
	ctrl.queueURL = queueURL
	ctrl.queueARN = queueARN

	topicARN, subARN, err := ctrl.createSNSTopic(name)
	if err != nil {
		return ctrl, err
	}
	ctrl.topicARN = topicARN
	ctrl.subARN = subARN

	err = ctrl.setPolicy()
	if err != nil {
		return ctrl, err
	}

	return ctrl, nil
}

// GetTopicARN will return your topic arn
func (c *Informer) GetTopicARN() string {
	return c.topicARN
}

// AddMessageHander will allow you add the handler functions to the controller
// object for operating on.
func (c *Informer) AddMessageHander(handlers ControllerHandler) {
	c.handlers = handlers
}

// Start will load up messages and call the message handler function
func (c *Informer) Start(stopCh <-chan struct{}) {
	go func() {
		for {
			select {
			case <-stopCh:
				klog.Infof("Shutting down SQS listener")
				return
			default:
			}
			input := &sqs.ReceiveMessageInput{}
			input.SetQueueUrl(c.queueURL)
			input.SetAttributeNames(aws.StringSlice([]string{"SentTimestamp"}))
			input.SetMaxNumberOfMessages(1)
			input.SetMessageAttributeNames(aws.StringSlice([]string{"All"}))
			input.SetWaitTimeSeconds(10)

			output, err := c.sqsclient.ReceiveMessage(input)
			if err != nil {
				klog.Errorf("Error pulling messages off the queue '%s'", err.Error())
				return
			}

			for _, message := range output.Messages {
				evtMessage := &event.Message{}
				err := json.Unmarshal([]byte(*message.Body), evtMessage)
				if err != nil {
					klog.Errorf("Error unmarshalling the message body '%s'", err.Error())
					break
				}

				evt := &event.Event{}
				err = event.Unmarshal(evtMessage.Message, evt)
				if err != nil {
					klog.Errorf("Error unmarshalling the message '%s'", err.Error())
					break
				}

				c.handlers.OnMessage(evt)

				deleteInput := &sqs.DeleteMessageInput{}
				deleteInput.SetQueueUrl(c.queueURL)
				deleteInput.SetReceiptHandle(*message.ReceiptHandle)
				_, err = c.sqsclient.DeleteMessage(deleteInput)
				if err != nil {
					klog.Errorf("Error deleting the message '%s'", err.Error())
					break
				}
			}
		}
	}()
}

func (c *Informer) getQueueURL(name string) (string, error) {
	input := &sqs.GetQueueUrlInput{}
	input.SetQueueName(name)
	output, err := c.sqsclient.GetQueueUrl(input)
	if err != nil {
		return "", err
	}
	return *output.QueueUrl, nil
}

func (c *Informer) createQueue(name string) (string, error) {
	input := &sqs.CreateQueueInput{}
	input.SetQueueName(name)
	output, err := c.sqsclient.CreateQueue(input)
	if err != nil {
		return "", err
	}
	return *output.QueueUrl, nil
}

func (c *Informer) createSQSQueue(name string) (string, string, error) {
	if name == "" {
		return "", "", fmt.Errorf("SQS Queue name can't be blank")
	}

	var queueURL string
	var err error
	queueURL, err = c.getQueueURL(name)
	if err != nil {
		klog.Infof("Error getting Queue URL '%s'", err.Error())
		queueURL, err = c.createQueue(name)
		if err != nil {
			return "", "", err
		}
	}

	queueQueryInputs := &sqs.GetQueueAttributesInput{}
	queueQueryInputs.SetQueueUrl(queueURL)
	queueQueryInputs.SetAttributeNames([]*string{aws.String("All")})

	sqsQueueOutput, err := c.sqsclient.GetQueueAttributes(queueQueryInputs)
	if err != nil {
		return "", "", err
	}
	queueARN := *sqsQueueOutput.Attributes["QueueArn"]

	return string(queueURL), string(queueARN), nil
}

func (c *Informer) createSNSTopic(name string) (string, string, error) {
	topicInputs := &sns.CreateTopicInput{}
	topicInputs.SetName(name)
	output, err := c.snsclient.CreateTopic(topicInputs)
	if err != nil {
		return "", "", err
	}
	topicARN := *output.TopicArn

	subInput := &sns.SubscribeInput{}
	subInput.SetTopicArn(topicARN)
	subInput.SetEndpoint(c.queueARN)
	subInput.SetProtocol("sqs")

	subOutput, err := c.snsclient.Subscribe(subInput)
	if err != nil {
		return "", "", err
	}
	subARN := *subOutput.SubscriptionArn

	return topicARN, subARN, nil
}

func (c *Informer) setPolicy() error {
	policy := newPolicy(c.queueARN, c.queueName, []string{c.topicARN})
	policyb, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	input := &sqs.SetQueueAttributesInput{}
	input.SetQueueUrl(c.queueURL)
	input.SetAttributes(map[string]*string{"Policy": aws.String(string(policyb))})
	_, err = c.sqsclient.SetQueueAttributes(input)
	if err != nil {
		return err
	}
	return nil
}

func newPolicy(queueARN string, name string, topicARNs []string) Policy {
	statements := []Statement{}
	for _, topicARN := range topicARNs {
		statements = append(statements, Statement{
			Sid:       topicARN + " statment",
			Effect:    "Allow",
			Principal: "*",
			Action:    []string{"SQS:ReceiveMessage", "SQS:SendMessage"},
			Resource:  queueARN,
			Condition: Condition{
				ArnEquals: ArnEquals{
					AwsSourceArn: topicARN,
				},
			},
		})
	}

	return Policy{
		Version:   "2012-10-17",
		ID:        queueARN + "/" + name + "-policy",
		Statement: statements,
	}
}
