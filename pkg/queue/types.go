package queue

import (
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
)

// Queue wraps the config object for updating
type Queue struct {
	config       config.Config
	queueURL     string
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
