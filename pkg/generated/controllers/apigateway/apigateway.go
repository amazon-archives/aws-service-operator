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

package apigateway

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	metav1alpha1 "awsoperator.io/pkg/apis/meta/v1alpha1"
	cfnencoder "awsoperator.io/pkg/encoding/cloudformation"

	apigateway "awsoperator.io/pkg/apis/apigateway/v1alpha1"
	awsclientset "awsoperator.io/pkg/generated/clientset/versioned"
	apigatewayInformers "awsoperator.io/pkg/generated/informers/externalversions/apigateway/v1alpha1"
	apigatewayListers "awsoperator.io/pkg/generated/listers/apigateway/v1alpha1"

	cloudformation "awsoperator.io/pkg/apis/cloudformation/v1alpha1"
	cloudformationInformers "awsoperator.io/pkg/generated/informers/externalversions/cloudformation/v1alpha1"
	cloudformationListers "awsoperator.io/pkg/generated/listers/cloudformation/v1alpha1"
)

const (
	controllerName = "apigateway-controller"

	accountDeletionFinalizerName = "deletion.account.apigateway.awsoperator.io"

	// SuccessSynced is used as part of the Event 'reason' when a Stack is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Stack fails
	// to sync due to a Stack of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Stack already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Stack"
	// MessageResourceSynced is the message used for an Event fired when a resource
	// is synced successfully
	MessageResourceSynced = "Resource synced successfully"
)

// Interface exposes the methods to run the controller
type Interface interface {
	Run(int, <-chan struct{})
}

// Controller will return the controller for operating resources
type Controller struct {
	kubeclientset kubernetes.Interface
	awsclientset  awsclientset.Interface

	stackLister cloudformationListers.StackLister
	stackSynced cache.InformerSynced
	stackIndex  cache.Indexer

	syncAccountHandler    func(string) error
	enqueueAccountHandler func(*apigateway.Account)
	accountLister         apigatewayListers.AccountLister
	accountSynced         cache.InformerSynced
	accountIndex          cache.Indexer
	accountQueue          workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

// New will generate a new controller
func New(
	kubeclientset kubernetes.Interface,
	awsclientset awsclientset.Interface,
	stackInformer cloudformationInformers.StackInformer,
	accountInformer apigatewayInformers.AccountInformer) (*Controller, error) {
	klog.V(4).Info("Creating APIGateway event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})

	ctrl := &Controller{
		kubeclientset: kubeclientset,
		awsclientset:  awsclientset,

		stackLister: stackInformer.Lister(),
		stackSynced: stackInformer.Informer().HasSynced,
		stackIndex:  stackInformer.Informer().GetIndexer(),

		accountLister: accountInformer.Lister(),
		accountSynced: accountInformer.Informer().HasSynced,

		accountQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "account"),

		recorder: recorder,
	}

	klog.Info("Setting up Account event handlers")

	accountInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctrl.addAccount,
		UpdateFunc: ctrl.updateAccount,
		DeleteFunc: ctrl.deleteAccount,
	})

	ctrl.syncAccountHandler = ctrl.syncAccount
	ctrl.enqueueAccountHandler = ctrl.enqueueAccount

	err := accountInformer.Informer().GetIndexer().AddIndexers(cache.Indexers{
		"stackID": indexAccountsByID,
	})
	if err != nil {
		return nil, fmt.Errorf("Error to add stackID index: %s", err.Error())
	}

	return ctrl, nil
}

// Run will kick off a workqueue and start processing items
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.accountQueue.ShutDown()

	klog.Info("Starting Apigateway controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.stackSynced, c.accountSynced); !ok {
		return fmt.Errorf("Failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runAccountWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")
	return nil
}

func indexAccountsByID(obj interface{}) ([]string, error) {
	resource, ok := obj.(*apigateway.Account)
	if !ok {
		return []string{}, nil
	}
	if len(resource.Status.StatusMeta.StackID) == 0 {
		return []string{}, nil
	}
	return []string{resource.Status.StatusMeta.StackID}, nil
}

func (c *Controller) runAccountWorker() {
	for c.processNextAccountWorkItem() {
	}
}

func (c *Controller) processNextAccountWorkItem() bool {
	obj, shutdown := c.accountQueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.accountQueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.accountQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncAccountHandler(key); err != nil {
			c.accountQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.accountQueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) addAccount(obj interface{}) {
	resource := obj.(*apigateway.Account)
	klog.V(4).Infof("adding Apigateway Account %s", resource.Name)
	c.enqueueAccountHandler(resource)
}

func (c *Controller) updateAccount(old, cur interface{}) {
	oldR := old.(*apigateway.Account)
	curR := cur.(*apigateway.Account)
	klog.V(4).Infof("updating Apigateway Account %s", oldR.Name)
	c.enqueueAccountHandler(curR)
}

func (c *Controller) deleteAccount(obj interface{}) {
	r, ok := obj.(*apigateway.Account)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		r, ok = tombstone.Obj.(*apigateway.Account)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an ApiGateway Account %#v", obj))
			return
		}
	}
	klog.V(4).Infof("deleting Apigateway Account %s", r.Name)

	c.enqueueAccountHandler(r)
}

func (c *Controller) enqueueAccount(resource *apigateway.Account) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(resource)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", resource, err))
		return
	}

	c.accountQueue.Add(key)
}

func (c *Controller) syncAccount(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	resource, err := c.accountLister.Accounts(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("account '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	var stack *cloudformation.Stack

	stackName := resource.Spec.CloudFormationMeta.StackName
	if stackName == "" {
		stack, err = c.awsclientset.CloudformationV1alpha1().Stacks(namespace).Create(newAccountStack(resource))
		if err != nil {
			return err
		}

		stackName = stack.Name
	}

	return nil
}

func newAccountStack(resource *apigateway.Account) *cloudformation.Stack {
	params := map[string]string{}
	cfnencoder.MarshalTypes(params, resource.Spec, "Parameter")

	return &cloudformation.Stack{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{"apigateway", "account", resource.Namespace, resource.Name}, "-"),
			Namespace: resource.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(resource, apigateway.SchemeGroupVersion.WithKind("Account")),
			},
		},
		Spec: cloudformation.StackSpec{
			CloudFormationMeta: metav1alpha1.CloudFormationMeta{
				Region:                resource.Spec.Region,
				NotificationARNs:      resource.Spec.NotificationARNs,
				OnFailure:             resource.Spec.OnFailure,
				Tags:                  resource.Spec.Tags,
				TerminationProtection: resource.Spec.TerminationProtection,
			},
			Parameters:   params,
			TemplateBody: resource.GetTemplate(),
		},
	}
}
