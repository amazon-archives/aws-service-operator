package server

import (
	// "fmt"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/ec2metadata"
	// "github.com/aws/aws-sdk-go/aws/session"

	awsscheme "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/scheme"
	//awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	opBase "github.com/awslabs/aws-service-operator/pkg/operators/base"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

const controllerName = "aws-service-operator"

// New creates a new server from a config
func New(config *config.Config) *Server {
	return &Server{
		Config: config,
	}
}

// Run starts the server to listen to Kubernetes
func (c *Server) Run(stopChan chan struct{}) {
	config := c.Config
	logger := config.Logger
	logger.Info("getting kubernetes context")

	awsscheme.AddToScheme(scheme.Scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logger.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: config.KubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})
	config.Recorder = recorder

	// start watching the aws operator resources
	logger.Info("Watching the resources")
	operators := opBase.New(config) // TODO: remove context and Clientset
	err := operators.Watch(corev1.NamespaceAll, stopChan)
	if err != nil {
		logger.Infof("error watching operators '%s'\n", err)
	}
}
