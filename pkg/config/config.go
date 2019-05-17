package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	goVersion "github.com/christopherhein/go-version"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

const controllerName = "aws-service-operator"

// Config defines the configuration for the operator
type Config struct {
	Version          *goVersion.Info
	Region           string
	Kubeconfig       string
	MasterURL        string
	QueueURL         string
	QueueARN         string
	AWSSession       *session.Session
	AWSClientset     awsclient.ServiceoperatorV1alpha1Interface
	KubeClientset    kubernetes.Interface
	RESTConfig       *rest.Config
	LoggingConfig    LoggingConfig
	Logger           *logrus.Entry
	Resources        map[string]bool
	ClusterName      string
	Bucket           string
	AccountID        string
	DefaultNamespace string
	K8sNamespace     string
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

// CreateContext will create all the contexts for the informers
func CreateContext(masterURL, kubeconfig string) (awsclient.ServiceoperatorV1alpha1Interface, kubernetes.Interface, *rest.Config, error) {
	config, err := getKubeconfig(masterURL, kubeconfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	awsclientset, err := awsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create object store clientset. %+v", err)
	}

	return awsclientset, clientset, config, nil
}

func CreateRecorder(logger *logrus.Entry, kubeclientset kubernetes.Interface) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logger.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})
}
