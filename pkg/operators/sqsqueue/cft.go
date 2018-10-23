// >>>>>>> DO NOT EDIT THIS FILE <<<<<<<<<<
// This file is autogenerated via `aws-operator-codegen process`
// If you'd like the change anything about this file make edits to the .templ
// file in the pkg/codegen/assets directory.

package sqsqueue

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awsV1alpha1 "github.com/awslabs/aws-service-operator/pkg/apis/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/helpers"
)

// New generates a new object
func New(config *config.Config, sqsqueue *awsV1alpha1.SQSQueue, topicARN string) *Cloudformation {
	return &Cloudformation{
		SQSQueue: sqsqueue,
		config:   config,
		topicARN: topicARN,
	}
}

// Cloudformation defines the sqsqueue cfts
type Cloudformation struct {
	config   *config.Config
	SQSQueue *awsV1alpha1.SQSQueue
	topicARN string
}

// StackName returns the name of the stack based on the aws-operator-config
func (s *Cloudformation) StackName() string {
	return helpers.StackName(s.config.ClusterName, "sqsqueue", s.SQSQueue.Name, s.SQSQueue.Namespace)
}

// GetOutputs return the stack outputs from the DescribeStacks call
func (s *Cloudformation) GetOutputs() (map[string]string, error) {
	outputs := map[string]string{}
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DescribeStacksInput{
		StackName: aws.String(s.StackName()),
	}

	output, err := svc.DescribeStacks(&stackInputs)
	if err != nil {
		return nil, err
	}
	// Not sure if this is even possible
	if len(output.Stacks) != 1 {
		return nil, errors.New("no stacks returned with that stack name")
	}

	for _, out := range output.Stacks[0].Outputs {
		outputs[*out.OutputKey] = *out.OutputValue
	}

	return outputs, err
}

// CreateStack will create the stack with the supplied params
func (s *Cloudformation) CreateStack() (output *cloudformation.CreateStackOutput, err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	cftemplateURL := helpers.GetCloudFormationTemplate(s.config, "sqsqueue", s.SQSQueue.Spec.CloudFormationTemplateName, s.SQSQueue.Spec.CloudFormationTemplateNamespace)

	cftemplate, err := helpers.FetchAndProcessTemplate(cftemplateURL, s.SQSQueue)
	if err != nil {
		return output, err
	}

	stackInputs := cloudformation.CreateStackInput{
		StackName:    aws.String(s.StackName()),
		TemplateBody: cftemplate,
		NotificationARNs: []*string{
			aws.String(s.topicARN),
		},
	}

	resourceName := helpers.CreateParam("ResourceName", s.SQSQueue.Name)
	resourceVersion := helpers.CreateParam("ResourceVersion", s.SQSQueue.ResourceVersion)
	namespace := helpers.CreateParam("Namespace", s.SQSQueue.Namespace)
	clusterName := helpers.CreateParam("ClusterName", s.config.ClusterName)
	contentBasedDeduplicationTemp := "{{.Obj.Spec.ContentBasedDeduplication}}"
	contentBasedDeduplicationValue, err := helpers.Templatize(contentBasedDeduplicationTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	contentBasedDeduplication := helpers.CreateParam("ContentBasedDeduplication", helpers.Stringify(contentBasedDeduplicationValue))
	delaySecondsTemp := "{{.Obj.Spec.DelaySeconds}}"
	delaySecondsValue, err := helpers.Templatize(delaySecondsTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	delaySeconds := helpers.CreateParam("DelaySeconds", helpers.Stringify(delaySecondsValue))
	maximumMessageSizeTemp := "{{.Obj.Spec.MaximumMessageSize}}"
	maximumMessageSizeValue, err := helpers.Templatize(maximumMessageSizeTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	maximumMessageSize := helpers.CreateParam("MaximumMessageSize", helpers.Stringify(maximumMessageSizeValue))
	messageRetentionPeriodTemp := "{{.Obj.Spec.MessageRetentionPeriod}}"
	messageRetentionPeriodValue, err := helpers.Templatize(messageRetentionPeriodTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	messageRetentionPeriod := helpers.CreateParam("MessageRetentionPeriod", helpers.Stringify(messageRetentionPeriodValue))
	receiveMessageWaitTimeSecondsTemp := "{{.Obj.Spec.ReceiveMessageWaitTimeSeconds}}"
	receiveMessageWaitTimeSecondsValue, err := helpers.Templatize(receiveMessageWaitTimeSecondsTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	receiveMessageWaitTimeSeconds := helpers.CreateParam("ReceiveMessageWaitTimeSeconds", helpers.Stringify(receiveMessageWaitTimeSecondsValue))
	usedeadletterQueueTemp := "{{.Obj.Spec.UsedeadletterQueue}}"
	usedeadletterQueueValue, err := helpers.Templatize(usedeadletterQueueTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	usedeadletterQueue := helpers.CreateParam("UsedeadletterQueue", helpers.Stringify(usedeadletterQueueValue))
	visibilityTimeoutTemp := "{{.Obj.Spec.VisibilityTimeout}}"
	visibilityTimeoutValue, err := helpers.Templatize(visibilityTimeoutTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	visibilityTimeout := helpers.CreateParam("VisibilityTimeout", helpers.Stringify(visibilityTimeoutValue))
	fifoQueueTemp := "{{.Obj.Spec.FifoQueue}}"
	fifoQueueValue, err := helpers.Templatize(fifoQueueTemp, helpers.Data{Obj: s.SQSQueue, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	fifoQueue := helpers.CreateParam("FifoQueue", helpers.Stringify(fifoQueueValue))

	parameters := []*cloudformation.Parameter{}
	parameters = append(parameters, resourceName)
	parameters = append(parameters, resourceVersion)
	parameters = append(parameters, namespace)
	parameters = append(parameters, clusterName)
	parameters = append(parameters, contentBasedDeduplication)
	parameters = append(parameters, delaySeconds)
	parameters = append(parameters, maximumMessageSize)
	parameters = append(parameters, messageRetentionPeriod)
	parameters = append(parameters, receiveMessageWaitTimeSeconds)
	parameters = append(parameters, usedeadletterQueue)
	parameters = append(parameters, visibilityTimeout)
	parameters = append(parameters, fifoQueue)

	stackInputs.SetParameters(parameters)

	resourceNameTag := helpers.CreateTag("ResourceName", s.SQSQueue.Name)
	resourceVersionTag := helpers.CreateTag("ResourceVersion", s.SQSQueue.ResourceVersion)
	namespaceTag := helpers.CreateTag("Namespace", s.SQSQueue.Namespace)
	clusterNameTag := helpers.CreateTag("ClusterName", s.config.ClusterName)

	tags := []*cloudformation.Tag{}
	tags = append(tags, resourceNameTag)
	tags = append(tags, resourceVersionTag)
	tags = append(tags, namespaceTag)
	tags = append(tags, clusterNameTag)

	stackInputs.SetTags(tags)

	output, err = svc.CreateStack(&stackInputs)
	return
}

// UpdateStack will update the existing stack
func (s *Cloudformation) UpdateStack(updated *awsV1alpha1.SQSQueue) (output *cloudformation.UpdateStackOutput, err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	cftemplate := helpers.GetCloudFormationTemplate(s.config, "sqsqueue", updated.Spec.CloudFormationTemplateName, updated.Spec.CloudFormationTemplateNamespace)

	stackInputs := cloudformation.UpdateStackInput{
		StackName:   aws.String(s.StackName()),
		TemplateURL: aws.String(cftemplate),
		NotificationARNs: []*string{
			aws.String(s.topicARN),
		},
	}

	resourceName := helpers.CreateParam("ResourceName", s.SQSQueue.Name)
	resourceVersion := helpers.CreateParam("ResourceVersion", s.SQSQueue.ResourceVersion)
	namespace := helpers.CreateParam("Namespace", s.SQSQueue.Namespace)
	clusterName := helpers.CreateParam("ClusterName", s.config.ClusterName)
	contentBasedDeduplicationTemp := "{{.Obj.Spec.ContentBasedDeduplication}}"
	contentBasedDeduplicationValue, err := helpers.Templatize(contentBasedDeduplicationTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	contentBasedDeduplication := helpers.CreateParam("ContentBasedDeduplication", helpers.Stringify(contentBasedDeduplicationValue))
	delaySecondsTemp := "{{.Obj.Spec.DelaySeconds}}"
	delaySecondsValue, err := helpers.Templatize(delaySecondsTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	delaySeconds := helpers.CreateParam("DelaySeconds", helpers.Stringify(delaySecondsValue))
	maximumMessageSizeTemp := "{{.Obj.Spec.MaximumMessageSize}}"
	maximumMessageSizeValue, err := helpers.Templatize(maximumMessageSizeTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	maximumMessageSize := helpers.CreateParam("MaximumMessageSize", helpers.Stringify(maximumMessageSizeValue))
	messageRetentionPeriodTemp := "{{.Obj.Spec.MessageRetentionPeriod}}"
	messageRetentionPeriodValue, err := helpers.Templatize(messageRetentionPeriodTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	messageRetentionPeriod := helpers.CreateParam("MessageRetentionPeriod", helpers.Stringify(messageRetentionPeriodValue))
	receiveMessageWaitTimeSecondsTemp := "{{.Obj.Spec.ReceiveMessageWaitTimeSeconds}}"
	receiveMessageWaitTimeSecondsValue, err := helpers.Templatize(receiveMessageWaitTimeSecondsTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	receiveMessageWaitTimeSeconds := helpers.CreateParam("ReceiveMessageWaitTimeSeconds", helpers.Stringify(receiveMessageWaitTimeSecondsValue))
	usedeadletterQueueTemp := "{{.Obj.Spec.UsedeadletterQueue}}"
	usedeadletterQueueValue, err := helpers.Templatize(usedeadletterQueueTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	usedeadletterQueue := helpers.CreateParam("UsedeadletterQueue", helpers.Stringify(usedeadletterQueueValue))
	visibilityTimeoutTemp := "{{.Obj.Spec.VisibilityTimeout}}"
	visibilityTimeoutValue, err := helpers.Templatize(visibilityTimeoutTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	visibilityTimeout := helpers.CreateParam("VisibilityTimeout", helpers.Stringify(visibilityTimeoutValue))
	fifoQueueTemp := "{{.Obj.Spec.FifoQueue}}"
	fifoQueueValue, err := helpers.Templatize(fifoQueueTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	fifoQueue := helpers.CreateParam("FifoQueue", helpers.Stringify(fifoQueueValue))

	parameters := []*cloudformation.Parameter{}
	parameters = append(parameters, resourceName)
	parameters = append(parameters, resourceVersion)
	parameters = append(parameters, namespace)
	parameters = append(parameters, clusterName)
	parameters = append(parameters, contentBasedDeduplication)
	parameters = append(parameters, delaySeconds)
	parameters = append(parameters, maximumMessageSize)
	parameters = append(parameters, messageRetentionPeriod)
	parameters = append(parameters, receiveMessageWaitTimeSeconds)
	parameters = append(parameters, usedeadletterQueue)
	parameters = append(parameters, visibilityTimeout)
	parameters = append(parameters, fifoQueue)

	stackInputs.SetParameters(parameters)

	resourceNameTag := helpers.CreateTag("ResourceName", s.SQSQueue.Name)
	resourceVersionTag := helpers.CreateTag("ResourceVersion", s.SQSQueue.ResourceVersion)
	namespaceTag := helpers.CreateTag("Namespace", s.SQSQueue.Namespace)
	clusterNameTag := helpers.CreateTag("ClusterName", s.config.ClusterName)

	tags := []*cloudformation.Tag{}
	tags = append(tags, resourceNameTag)
	tags = append(tags, resourceVersionTag)
	tags = append(tags, namespaceTag)
	tags = append(tags, clusterNameTag)

	stackInputs.SetTags(tags)

	output, err = svc.UpdateStack(&stackInputs)
	return
}

// DeleteStack will delete the stack
func (s *Cloudformation) DeleteStack() (err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DeleteStackInput{}
	stackInputs.SetStackName(s.StackName())

	_, err = svc.DeleteStack(&stackInputs)
	return
}

// WaitUntilStackDeleted will delete the stack
func (s *Cloudformation) WaitUntilStackDeleted() (err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DescribeStacksInput{
		StackName: aws.String(s.StackName()),
	}

	err = svc.WaitUntilStackDeleteComplete(&stackInputs)
	return
}
