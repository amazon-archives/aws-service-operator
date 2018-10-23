package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

// Config defines the configuration for the operator
type Config struct {
	Region           string
	Kubeconfig       string
	MasterURL        string
	QueueURL         string
	QueueARN         string
	AWSSession       *session.Session
	AWSClientset     awsclient.ServiceoperatorV1alpha1Interface
	KubeClientset    kubernetes.Interface
	RESTConfig       *rest.Config
	LoggingConfig    *LoggingConfig
	Logger           *logrus.Entry
	Resources        map[string]bool
	ClusterName      string
	Bucket           string
	AccountID        string
	DefaultNamespace string
	Recorder         record.EventRecorder
}

// LoggingConfig defines the attributes for the logger
type LoggingConfig struct {
	File              string
	Level             string
	DisableTimestamps bool
	FullTimestamps    bool
}

func getKubeconfig(masterURL, kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	return rest.InClusterConfig()
}

func (c *Config) CreateContext(masterURL, kubeconfig string) error {
	config, err := getKubeconfig(masterURL, kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to get k8s client. %+v", err)
	}

	awsclientset, err := awsclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create object store clientset. %+v", err)
	}

	c.AWSClientset = awsclientset
	c.KubeClientset = clientset
	c.RESTConfig = config

	return nil
}
