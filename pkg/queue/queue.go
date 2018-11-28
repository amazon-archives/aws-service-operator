package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/queuemanager"
)

// New will initialize the Queue object for watching
func New(config config.Config, awsclientset awsclient.ServiceoperatorV1alpha1Interface, timeout int) *Queue {
	return &Queue{
		config:       config,
		awsclientset: awsclientset,
		timeout:      int64(timeout),
	}
}

// RegisterQueue wkll create the Queue so it is accessible to SNS.
func RegisterQueue(awsSession *session.Session, clusterName, name string) (queueURL, queueARN string, err error) {
	sqsSvc := sqs.New(awsSession)
	keyname := keyName(clusterName, name)
	queueInputs := sqs.CreateQueueInput{
		QueueName: aws.String(keyname),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("10"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	}
	sqsOutput, err := sqsSvc.CreateQueue(&queueInputs)
	if err != nil {
		return "", "", err
	}
	queueURL = *sqsOutput.QueueUrl

	queueQueryInputs := sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []*string{
			aws.String("All"),
		},
	}
	sqsQueueOutput, err := sqsSvc.GetQueueAttributes(&queueQueryInputs)
	if err != nil {
		return "", "", err
	}
	queueARN = *sqsQueueOutput.Attributes["QueueArn"]
	return queueURL, queueARN, nil
}

// Subscribe will listen to the global queue and distribute messages
func Subscribe(config config.Config, manager *queuemanager.QueueManager, ctx context.Context) {
	logger := config.Logger
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
	if err != nil {
		logger.WithError(err).Error("error creating AWS session")
		os.Exit(1)
	}

	svc := sqs.New(sess)

	process(config, svc, manager, ctx)
}

// Register will create all the affilate resources
func (q *Queue) Register(name string) (topicARN string, subARN string) {
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

	// Move to somewhere else.
	subInput := sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Endpoint: aws.String(config.QueueARN),
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

func SetQueuePolicy(config config.Config, manager *queuemanager.QueueManager) error {
	sqsSvc := sqs.New(config.AWSSession)
	topicARNs := manager.Keys()

	policy := newPolicy(config.QueueARN, config.ClusterName, topicARNs)
	policyb, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	addPermInput := &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(config.QueueURL),
		Attributes: map[string]*string{
			"Policy": aws.String(string(policyb)),
		}}
	_, err = sqsSvc.SetQueueAttributes(addPermInput)
	if err != nil {
		return err
	}
	return nil
}

func keyName(clusterName string, name string) string {
	return clusterName + "-" + name
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

func process(config config.Config, svc *sqs.SQS, manager *queuemanager.QueueManager, ctx context.Context) error {
	logger := config.Logger
	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(config.QueueURL),
			AttributeNames: aws.StringSlice([]string{
				"SentTimestamp",
			}),
			MaxNumberOfMessages: aws.Int64(1),
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
			WaitTimeSeconds: aws.Int64(10),
		})
		if err != nil {
			logger.WithError(err).Error("unable to receive messages from queue")
			os.Exit(1)
		}

		for _, message := range result.Messages {
			mb := &queuemanager.MessageBody{}
			err := json.Unmarshal([]byte(*message.Body), mb)
			if err != nil {
				logger.WithError(err).Error("error parsing message body")
			}
			mb.ParseMessage()
			logger.Debugf("%+v", mb)

			h, ok := manager.Get(mb.TopicARN)
			if !ok {
				logger.Errorf("error missing handler function for topic %s", mb.TopicARN)
			}

			err = h.HandleMessage(config, mb)
			if err != nil {
				logger.WithError(err).Error("error processing message")
			}
			logger.Debugf("stackID %v updated status to %v", mb.ParsedMessage["StackId"], mb.ParsedMessage["ResourceStatus"])
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(config.QueueURL),
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				logger.WithError(err).Error("delete error")
				os.Exit(1)
			}
		}
	}
}
