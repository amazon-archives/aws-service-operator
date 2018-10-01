package queue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/helpers"
	opkit "github.com/christopherhein/operator-kit"
	"os"
	"strings"
)

// HandlerFunc allows you to define a custom function for when a message is stored
type HandlerFunc func(config *config.Config, msg *MessageBody) error

// HandleMessage will stub the handler for processing messages
func (f HandlerFunc) HandleMessage(config *config.Config, msg *MessageBody) error {
	return f(config, msg)
}

// Handler allows a custom function to be passed
type Handler interface {
	HandleMessage(config *config.Config, msg *MessageBody) error
}

// ParseMessage will take the message attribute and make it readable
func (m *MessageBody) ParseMessage() error {
	m.Updatable = false
	resp := make(map[string]string)
	items := strings.Split(m.Message, "\n")
	for _, item := range items {
		x := strings.Split(item, "=")
		key := x[0]
		if key != "" {
			s := x[1]
			s = s[1 : len(s)-1]
			resp[key] = s
		}
	}
	m.ParsedMessage = resp

	var resourceProperties ResourceProperties
	if resp["ResourceProperties"] != "null" {
		err := json.Unmarshal([]byte(resp["ResourceProperties"]), &resourceProperties)
		if err != nil {
			return err
		}
		m.ResourceProperties = resourceProperties
		for _, tag := range resourceProperties.Tags {
			switch tag.Key {
			case "Namespace":
				m.Namespace = tag.Value
			case "ResourceName":
				m.ResourceName = tag.Value
			}
		}
		if m.Namespace != "" && m.ResourceName != "" {
			m.Updatable = true
		}
	}
	return nil
}

// IsComplete returns a simple status instead of the raw CFT resp
func (m *MessageBody) IsComplete() bool {
	return helpers.IsStackComplete(m.ParsedMessage["ResourceStatus"], true)
}

// New will initialize the Queue object for watching
func New(config *config.Config, context *opkit.Context, awsclientset awsclient.ServiceoperatorV1alpha1Interface, timeout int) *Queue {
	return &Queue{
		config:       config,
		context:      context,
		awsclientset: awsclientset,
		timeout:      int64(timeout),
	}
}

// Register will create all the affilate resources
func (q *Queue) Register(name string, obj interface{}) (topicARN string, queueURL string, queueAttrs map[string]*string, subARN string) {
	config := q.config
	logger := config.Logger
	keyname := keyName(config.ClusterName, name)
	svc := sns.New(config.AWSSession)
	topicInputs := sns.CreateTopicInput{Name: aws.String(keyname)}
	output, err := svc.CreateTopic(&topicInputs)
	if err != nil {
		logger.Errorf("Error creating SNS Topic with error '%s'", err.Error())
	}
	topicARN = *output.TopicArn
	logger.Infof("Created sns topic '%s'", topicARN)

	sqsSvc := sqs.New(config.AWSSession)
	queueInputs := sqs.CreateQueueInput{
		QueueName: aws.String(keyname),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("10"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	}
	sqsOutput, err := sqsSvc.CreateQueue(&queueInputs)
	if err != nil {
		logger.Errorf("Error creating SQS Queue with error '%s'", err.Error())
	}
	queueURL = *sqsOutput.QueueUrl
	logger.Infof("Created sqs queue '%s'", queueURL)
	q.queueURL = queueURL

	// Get Queue Information
	// This will later happen with the CRD cft subscription
	queueQueryInputs := sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []*string{
			aws.String("All"),
		},
	}
	sqsQueueOutput, err := sqsSvc.GetQueueAttributes(&queueQueryInputs)
	if err != nil {
		logger.Errorf("Error getting SQS Information with error '%s'", err.Error())
	}
	queueAttrs = sqsQueueOutput.Attributes

	policy := newPolicy(*queueAttrs["QueueArn"], keyname, topicARN)
	policyb, err := json.Marshal(policy)
	if err != nil {
		logger.WithError(err).Error("error encoding policy to json")
	}
	addPermInput := &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		Attributes: map[string]*string{
			"Policy": aws.String(string(policyb)),
		}}
	_, err = sqsSvc.SetQueueAttributes(addPermInput)
	if err != nil {
		logger.WithError(err).Error("error setting queue policy")
	}

	// Move to somewhere else.
	subInput := sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Endpoint: aws.String(*queueAttrs["QueueArn"]),
		Protocol: aws.String("sqs"),
	}
	subOutput, err := svc.Subscribe(&subInput)
	if err != nil {
		logger.Errorf("Error creating SNS -> SQS Subscription with error '%s'", err.Error())
	}
	subARN = *subOutput.SubscriptionArn
	logger.Infof("Created sns -> sqs subscription '%s'", subARN)

	return
}

// StartWatch will start a go routine to watch for new events to store in etcd
func (q *Queue) StartWatch(h Handler, stopCh <-chan struct{}) {
	config := q.config
	logger := config.Logger
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
	if err != nil {
		logger.WithError(err).Infof("error creating AWS session")
		os.Exit(1)
	}

	svc := sqs.New(sess)

	process(q, svc, h, stopCh)
}

func keyName(clusterName string, name string) string {
	return clusterName + "-" + name
}

func newPolicy(queueARN string, keyname string, topicARN string) Policy {
	return Policy{
		Version: "2012-10-17",
		ID:      queueARN + "/" + keyname + "-policy",
		Statement: []Statement{
			Statement{
				Sid:       keyname,
				Effect:    "Allow",
				Principal: "*",
				Action:    []string{"SQS:ReceiveMessage", "SQS:SendMessage"},
				Resource:  queueARN,
				Condition: Condition{
					ArnEquals: ArnEquals{
						AwsSourceArn: topicARN,
					},
				},
			},
		},
	}
}

func process(q *Queue, svc *sqs.SQS, h Handler, stopCh <-chan struct{}) error {
	config := q.config
	logger := config.Logger
	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(q.queueURL),
			AttributeNames: aws.StringSlice([]string{
				"SentTimestamp",
			}),
			MaxNumberOfMessages: aws.Int64(1),
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
			WaitTimeSeconds: aws.Int64(q.timeout),
		})
		if err != nil {
			logger.WithError(err).Error("unable to receive messages from queue")
			os.Exit(1)
		}

		for _, message := range result.Messages {
			mb := &MessageBody{}
			err := json.Unmarshal([]byte(*message.Body), mb)
			if err != nil {
				logger.WithError(err).Error("error parsing message body")
			}
			mb.ParseMessage()
			logger.Debugf("%+v", mb)

			err = h.HandleMessage(config, mb)
			if err != nil {
				logger.WithError(err).Error("error processing message")
			}
			logger.Infof("stackID %v updated status to %v", mb.ParsedMessage["StackId"], mb.ParsedMessage["ResourceStatus"])
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(q.queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				logger.WithError(err).Error("delete error")
				os.Exit(1)
			}
		}

		select {
		case <-stopCh:
			os.Exit(1)
		default:
		}
	}
}
