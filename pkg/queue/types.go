package queue

import (
	awsclient "github.com/christopherhein/aws-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/christopherhein/aws-operator/pkg/config"
	opkit "github.com/christopherhein/operator-kit"
)

// Queue wraps the config object for updating
type Queue struct {
	config       *config.Config
	queueURL     string
	context      *opkit.Context
	awsclientset awsclient.ServiceoperatorV1alpha1Interface
	timeout      int64
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

// MessageBody will parse the message from the Body of SQS
type MessageBody struct {
	Type               string `json:"Type"`
	TopicArn           string `json:"TopicArn"`
	Message            string `json:"Message"`
	ParsedMessage      map[string]string
	Namespace          string
	ResourceName       string
	ResourceProperties ResourceProperties
	Updatable          bool
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
