package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	awsclient "github.com/christopherhein/aws-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	opkit "github.com/christopherhein/operator-kit"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
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
}

// LoggingConfig defines the attributes for the logger
type LoggingConfig struct {
	File              string
	Level             string
	DisableTimestamps bool
	FullTimestamps    bool
}
