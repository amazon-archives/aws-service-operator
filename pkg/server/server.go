package server

import (
	"fmt"
	"net/http"

	awsscheme "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/scheme"
	"github.com/awslabs/aws-service-operator/pkg/config"
	opBase "github.com/awslabs/aws-service-operator/pkg/operators/base"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const controllerName = "aws-service-operator"

// New creates a new server from a config
func New(config *config.Config) *Server {
	return &Server{
		Config: config,
	}
}

func (c *Server) exposeMetrics(errChan chan error) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		errChan <- fmt.Errorf("unable to expose metrics: %v", err)
	}
}

func (c *Server) watchOperatorResources(errChan chan error, stopChan chan struct{}) {

	logger := c.Config.Logger

	logger.Info("getting kubernetes context")
	awsscheme.AddToScheme(scheme.Scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logger.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: c.Config.KubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})
	c.Config.Recorder = recorder

	// start watching the aws operator resources
	logger.WithFields(logrus.Fields{"resources": c.Config.Resources}).Info("Watching")
	operators := opBase.New(c.Config) // TODO: remove context and Clientset
	err := operators.Watch(corev1.NamespaceAll, stopChan)
	if err != nil {
		errChan <- fmt.Errorf("unable to watch resources: %v", err)
	}
}

// Run starts the server to listen to Kubernetes
func (c *Server) Run(stopChan chan struct{}) {
	config := c.Config
	logger := config.Logger
	errChan := make(chan error, 1)

	logger.Info("starting metrics server")
	go c.exposeMetrics(errChan)

	logger.Info("starting resource watcher")
	go c.watchOperatorResources(errChan, stopChan)

	err := <-errChan
	if err != nil {
		logger.Fatal(err)
	}

}
