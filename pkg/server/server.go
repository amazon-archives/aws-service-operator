package server

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"

	awsscheme "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/scheme"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	opBase "github.com/awslabs/aws-service-operator/pkg/operators/base"
	opkit "github.com/christopherhein/operator-kit"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	context, restConfig, kubeclientset, awsClientset, err := createContext(config.Kubeconfig)
	if err != nil {
		logger.Fatalf("failed to create context. %+v\n", err)
	}
	config.Context = context
	config.AWSClientset = awsClientset
	config.RESTConfig = restConfig

	// TODO: Figure out how to get the node tag so I can store the
	// `KubernetesCluster` attribute so that all resources can be cleaned up
	// using that tag. Also we can create an inventory of all aws resources that
	// are modifiable at any time and deletable using the kubectl cli

	region := ""
	ec2Session, err := session.NewSession()
	metadata := ec2metadata.New(ec2Session)
	if config.Region == "" {
		region, _ = metadata.Region()
	} else {
		region = config.Region
	}

	logger.Infof("region: %+v", region)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		logger.Infof("error creating AWS session '%s'\n", err)
	}
	config.AWSSession = sess

	awsscheme.AddToScheme(scheme.Scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logger.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})
	config.Recorder = recorder

	// start watching the aws operator resources
	logger.Info("Watching the resources")
	operators := opBase.New(config, context, awsClientset)
	err = operators.Watch(corev1.NamespaceAll, stopChan)
	if err != nil {
	  logger.Infof("error watching operators '%s'\n", err)
	}
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func createContext(kubeconfig string) (*opkit.Context, *rest.Config, kubernetes.Interface, awsclient.ServiceoperatorV1alpha1Interface, error) {
	config, err := getConfig(kubeconfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
	}

	awsclientset, err := awsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create object store clientset. %+v", err)
	}

	context := &opkit.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		Interval:              500 * time.Millisecond,
		Timeout:               60 * time.Second,
	}
	return context, config, clientset, awsclientset, nil
}
