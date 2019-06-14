/*
Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
not use this file except in compliance with the License. A copy of the
License is located at

     http://aws.amazon.com/apache2.0/

or in the "license" file accompanying this file. This file is distributed
on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing
permissions and limitations under the License.
*/

package controllermanager

import (
	"context"
	"net/http"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	self "awsoperator.io/pkg/apis/self/v1alpha1"
	ctrl "awsoperator.io/pkg/cloudformation"
	clientset "awsoperator.io/pkg/generated/clientset/versioned"
	"awsoperator.io/pkg/generated/controllers"
	informers "awsoperator.io/pkg/generated/informers/externalversions"
	"awsoperator.io/pkg/queue"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// Interface will wrap all controller running
type Interface interface {
	Run(int, <-chan struct{}) error
}

// ControllerManager implements all ctrl actions
type ControllerManager struct {
	config             self.Config
	mux                *http.Server
	sessions           map[string]*session.Session
	clients            map[string]cloudformationiface.CloudFormationAPI
	sqsclient          sqsiface.SQSAPI
	snsclient          snsiface.SNSAPI
	kubeclientset      kubernetes.Interface
	awsclientset       clientset.Interface
	awsinformerfactory informers.SharedInformerFactory
	queueFactory       queue.InformerFactory
}

// New will return a new instance of a Controller Manager
func New(c self.Config, httpServer *http.Server) *ControllerManager {
	// Needs to support in cluster methods still
	cfg, err := clientcmd.BuildConfigFromFlags(c.Kubernetes.URL, c.Kubernetes.Kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	awsClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building aws clientset: %s", err.Error())
	}

	awsInformerFactory := informers.NewSharedInformerFactory(awsClient, time.Second*3600)

	// Allows you to connect to any session from the SupportedRegions
	sessions := map[string]*session.Session{}
	clients := map[string]cloudformationiface.CloudFormationAPI{}
	for _, region := range c.AWS.SupportedRegions {
		sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			klog.Fatalf("Error building aws session: %s", err.Error())
		}
		sessions[region] = sess
		clients[region] = cloudformation.New(sess)
	}

	var sqsclient *sqs.SQS
	if c.AWS.Queue.Region != "" {
		sqsclient = sqs.New(sessions[c.AWS.Queue.Region])
	}
	var snsclient *sns.SNS
	if c.AWS.Queue.Region != "" {
		snsclient = sns.New(sessions[c.AWS.Queue.Region])
	}

	queueFactory, err := queue.New(sqsclient, snsclient, c.ClusterName)
	if err != nil {
		klog.Fatalf("Error creating queue: %s", err.Error())
	}

	return &ControllerManager{
		config:             c,
		sessions:           sessions,
		clients:            clients,
		sqsclient:          sqsclient,
		snsclient:          snsclient,
		mux:                httpServer,
		kubeclientset:      kubeClient,
		awsclientset:       awsClient,
		awsinformerfactory: awsInformerFactory,
		queueFactory:       queueFactory,
	}
}

// Run will start listing across the controllers depending on which ones have
// been configured in your config
func (cm *ControllerManager) Run(stopCh <-chan struct{}) error {
	klog.Info("Starting AWS Operator ControllerManager")

	klog.Infof("Version: %+v", cm.config.Version.Version)

	klog.Info("Starting Controllers")

	cfnctrl, err := ctrl.New(cm.kubeclientset,
		cm.awsclientset,
		cm.awsinformerfactory.Cloudformation().V1alpha1().Stacks(),
		cm.queueFactory,
		cm.clients)

	if err != nil {
		klog.Fatalf("Error starting controller: %s", err.Error())
		return err
	}

	cm.queueFactory.Start(stopCh)

	go func() {
		if err := cfnctrl.Run(2, stopCh); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}()

	ctrls := []controllers.Controller{}

	for _, resource := range cm.config.Resources {
		splitResource := strings.Split(resource, ".")
		controller := splitResource[0]

		if err := cm.config.ResourcesContain(controller); err == nil {
			ctrl, err := controllers.Get(controller,
				cm.kubeclientset,
				cm.awsclientset,
				cm.awsinformerfactory)
			if err != nil {
				klog.Fatalf("Error finding controller: %s", err.Error())
			}
			ctrls = append(ctrls, ctrl)

		}
	}

	cm.awsinformerfactory.Start(stopCh)

	for _, ctrl := range ctrls {
		go func() {
			if err := ctrl.Run(2, stopCh); err != nil {
				klog.Fatalf("Error running controller: %s", err.Error())
			}
		}()
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	// Shutdown metrics server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := cm.mux.Shutdown(ctx); err != nil {
		klog.Fatalf("Failed to shutdown error '%s", err.Error())
	}

	cancel()

	return nil
}
