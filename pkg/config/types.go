package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	opkit "github.com/christopherhein/operator-kit"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

// Config defines the configuration for the operator
type Config struct {
	Region        string
	Kubeconfig    string
	AWSSession    *session.Session
	AWSClientset  awsclient.ServiceoperatorV1alpha1Interface
	RESTConfig    *rest.Config
	Context       *opkit.Context
	LoggingConfig *LoggingConfig
	Logger        *logrus.Entry
	Resources     []string
	ClusterName   string
	Bucket        string
	AccountID     string
	Recorder	  record.EventRecorder
}

// LoggingConfig defines the attributes for the logger
type LoggingConfig struct {
	File              string
	Level             string
	DisableTimestamps bool
	FullTimestamps    bool
}
