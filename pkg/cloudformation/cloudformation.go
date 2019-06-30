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

package cloudformation

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/iancoleman/strcase"

	cloudformation "awsoperator.io/pkg/apis/cloudformation/v1alpha1"
	metav1alpha1 "awsoperator.io/pkg/apis/meta/v1alpha1"
	"awsoperator.io/pkg/controllerutils"
	"awsoperator.io/pkg/event"
	awsclientset "awsoperator.io/pkg/generated/clientset/versioned"
	cloudformationInformers "awsoperator.io/pkg/generated/informers/externalversions/cloudformation/v1alpha1"
	cloudformationListers "awsoperator.io/pkg/generated/listers/cloudformation/v1alpha1"
	"awsoperator.io/pkg/queue"
	"awsoperator.io/pkg/token"

	"github.com/aws/aws-sdk-go/aws"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

const (
	controllerName = "cloudformation-controller"

	stackDeletionFinalizerName = "deletion.stack.cloudformation.awsoperator.io"

	// SuccessSynced is used as part of the Event 'reason' when a Stack is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Stack fails
	// to sync due to a Stack of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Stack already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Stack"
	// MessageResourceSynced is the message used for an Event fired when a Stack
	// is synced successfully
	MessageResourceSynced = "Stack synced successfully"
)

type Controller struct {
	notificationARN string

	kubeclientset kubernetes.Interface
	awsclientset  awsclientset.Interface

	clients map[string]cloudformationiface.CloudFormationAPI

	syncHandler  func(string) error
	enqueueStack func(*cloudformation.Stack)

	stackLister cloudformationListers.StackLister
	stackSynced cache.InformerSynced
	stackIndex  cache.Indexer

	clientToken token.Token

	queue    workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

// New returns the ApiGateway controller
func New(
	kubeclientset kubernetes.Interface,
	awsclientset awsclientset.Interface,
	stackInformer cloudformationInformers.StackInformer,
	queueInformer queue.InformerFactory,
	clients map[string]cloudformationiface.CloudFormationAPI) (*Controller, error) {
	klog.V(4).Info("Creating Cloudformation event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})

	ctrl := &Controller{
		notificationARN: queueInformer.GetTopicARN(),
		kubeclientset:   kubeclientset,
		awsclientset:    awsclientset,
		clients:         clients,
		stackLister:     stackInformer.Lister(),
		stackSynced:     stackInformer.Informer().HasSynced,
		stackIndex:      stackInformer.Informer().GetIndexer(),
		queue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cloudformation"),
		recorder:        recorder,
	}

	klog.Info("Setting up Cloudformation event handlers")

	stackInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctrl.addStack,
		UpdateFunc: ctrl.updateStack,
		DeleteFunc: ctrl.deleteStack,
	})

	queueInformer.AddMessageHander(queue.ControllerHandlerFuncs{
		MessageFunc: ctrl.onMessage,
	})

	ctrl.syncHandler = ctrl.syncStack
	ctrl.enqueueStack = ctrl.enqueue
	ctrl.clientToken = token.New()

	err := stackInformer.Informer().GetIndexer().AddIndexers(cache.Indexers{
		"stackID": indexStacksByID,
	})
	if err != nil {
		return nil, fmt.Errorf("Error to add stackID index: %s", err.Error())
	}

	return ctrl, nil
}

// Run will kick off a workqueue and start processing items
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Info("Starting CloudFormation controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.stackSynced); !ok {
		return fmt.Errorf("Failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")
	return nil
}

func indexStacksByID(obj interface{}) ([]string, error) {
	stack, ok := obj.(*cloudformation.Stack)
	if !ok {
		return []string{}, nil
	}
	if len(stack.Status.StatusMeta.StackID) == 0 {
		return []string{}, nil
	}
	return []string{stack.Status.StatusMeta.StackID}, nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.queue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.queue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.queue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.queue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// onMessage will parse the message and enqueue the stack
func (c *Controller) onMessage(obj interface{}) {
	event := obj.(*event.Event)
	var stack *cloudformation.Stack
	var ok bool

	objects, err := c.stackIndex.ByIndex("stackID", event.StackId)
	if err != nil {
		klog.Errorf("could not find stack but stack id '%s' with error '%s'", event.StackId, err.Error())
		return
	}

	if len(objects) > 0 {
		for _, obj := range objects {
			stack, ok = obj.(*cloudformation.Stack)
			if ok {
				klog.V(4).Infof("updating Cloudformation Stack status %s", stack.Name)
				c.enqueueStack(stack)
				return
			}
			klog.V(4).Infof("couldn't convert object into stack '%+v'", stack)
			return
		}
	}
	return
}

// addStack will enqueue the creation of an Cloudformation Stack
func (c *Controller) addStack(obj interface{}) {
	stack := obj.(*cloudformation.Stack)
	klog.V(4).Infof("adding Cloudformation Stack %s", stack.Name)
	c.enqueueStack(stack)
}

// updateStack will enqueue the updating of an Cloudformation Stack
func (c *Controller) updateStack(old, cur interface{}) {
	oldS := old.(*cloudformation.Stack)
	curS := cur.(*cloudformation.Stack)
	klog.V(4).Infof("updating Cloudformation Stack %s", oldS.Name)
	c.enqueueStack(curS)
}

// deleteStack will enqueue the deletion of an Cloudformation Stack
func (c *Controller) deleteStack(obj interface{}) {
	s, ok := obj.(*cloudformation.Stack)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		s, ok = tombstone.Obj.(*cloudformation.Stack)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an ApiGateway Stack %#v", obj))
			return
		}
	}
	klog.V(4).Infof("deleting Cloudformation Stack %s", s.Name)

	c.enqueueStack(s)
}

// enqueue will add an stack record to the queue
func (c *Controller) enqueue(stack *cloudformation.Stack) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(stack)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", stack, err))
		return
	}

	c.queue.Add(key)
}

func (c *Controller) syncStack(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	stack, err := c.stackLister.Stacks(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("cloudformation stack '%s' in work queue no longer exists", key))
			return nil
		}
		klog.Error("error getting stack")
		return err
	}

	if !stack.DeletionTimestamp.IsZero() {
		return c.deleteCFNStack(stack)
	}

	if ok, _ := controllerutils.ContainsFinalizer(stack, stackDeletionFinalizerName); !ok {
		err := c.addCFNFinalizer(stack)
		if err != nil {
			return err
		}
		utilruntime.HandleError(fmt.Errorf("cloudformation finalizer required to process '%s'", stack.Name))
		return nil
	}

	if stack.Spec.ClientRequestToken == "" {
		err := c.generateClientRequestToken(stack)
		if err != nil {
			return err
		}
		utilruntime.HandleError(fmt.Errorf("cloudformation client request token required to process '%s'", stack.Name))
		return nil
	}

	if stack.Status.StackID == "" {
		err := c.createCFNStack(stack)
		if err != nil {
			return err
		}
	}

	if stack.Status.Status == metav1alpha1.CreateCompleteStatus || stack.Status.Status == metav1alpha1.UpdateCompleteStatus {
		err := c.updateCFNStack(stack)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("stack errored updating '%s' reason '%s", stack.Name, err.Error()))
			return nil
		}
	}

	if stack.Status.StackID != "" {
		err := c.describeCFNStackStatus(stack)
		if err != nil {
			return err
		}
	}

	objRef := &corev1.ObjectReference{
		Kind:            stack.Kind,
		Namespace:       stack.Namespace,
		Name:            stack.Name,
		UID:             stack.UID,
		APIVersion:      stack.APIVersion,
		ResourceVersion: stack.ResourceVersion,
	}

	c.recorder.Event(objRef, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// generateClientRequestToken will generate the client request token
func (c *Controller) generateClientRequestToken(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()
	stackCopy.Spec.ClientRequestToken = c.clientToken.Generate()

	_, err := c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).Update(stackCopy)
	return err
}

// addCFNFinalizer will add the finalizer string for the controller
func (c *Controller) addCFNFinalizer(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()
	if err := controllerutils.AddFinalizer(stackCopy, stackDeletionFinalizerName); err != nil {
		return err
	}

	_, err := c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).Update(stackCopy)
	return err
}

// createCFNStack will create the cloudformation stack
func (c *Controller) createCFNStack(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()
	region := stackCopy.Spec.Region

	input := &cfn.CreateStackInput{}
	createStackInputs(stackCopy, c.notificationARN, input)

	output, err := c.clients[region].CreateStack(input)
	if err != nil {
		return err
	}

	status := cloudformation.StackStatus{
		StatusMeta: metav1alpha1.StatusMeta{
			Status:  metav1alpha1.CreateInProgressStatus,
			StackID: string(*output.StackId),
		},
	}
	stackCopy.Status = status

	_, err = c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).UpdateStatus(stackCopy)
	return err
}

// updateCFNStack will update the cloudformation stack
func (c *Controller) updateCFNStack(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()

	input := &cfn.UpdateStackInput{}

	updateStackInputs(stackCopy, c.notificationARN, input)
	region := stackCopy.Spec.Region

	_, err := c.clients[region].UpdateStack(input)
	if err != nil {
		return err
	}

	stackCopy.Status.Status = metav1alpha1.UpdateInProgressStatus

	_, err = c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).UpdateStatus(stackCopy)
	return err
}

func (c *Controller) deleteCFNStack(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()
	region := stackCopy.Spec.Region

	input := &cfn.DeleteStackInput{}
	input.SetStackName(namer(stackCopy.ObjectMeta.Name, stackCopy.ObjectMeta.Namespace))

	_, err := c.clients[region].DeleteStack(input)
	if err != nil {
		return err
	}

	deleteWaiter := &cfn.DescribeStacksInput{}
	deleteWaiter.SetStackName(stackCopy.Status.StackID)

	err = c.clients[region].WaitUntilStackDeleteComplete(deleteWaiter)
	if err != nil {
		return err
	}

	if err := controllerutils.RemoveFinalizer(stackCopy, stackDeletionFinalizerName); err != nil {
		return err
	}

	_, err = c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).Update(stackCopy)
	return err
}

func (c *Controller) describeCFNStackStatus(stack *cloudformation.Stack) error {
	stackCopy := stack.DeepCopy()
	region := stackCopy.Spec.Region

	input := &cfn.DescribeStacksInput{}
	input.SetStackName(stackCopy.Status.StackID)

	outputs, err := c.clients[region].DescribeStacks(input)
	if err != nil {
		return err
	}

	if len(outputs.Stacks) == 0 {
		return fmt.Errorf("could not find stack with name '%s'", stackCopy.Name)
	}

	outputsMap := map[string]string{}
	for _, output := range outputs.Stacks[0].Outputs {
		outputsMap[string(*output.OutputKey)] = string(*output.OutputValue)
	}
	stackCopy.Outputs = outputsMap

	_, err = c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).Update(stackCopy)
	if err != nil {
		return err
	}

	statusCopy := stackCopy.DeepCopy()

	status := cloudformation.StackStatus{
		StatusMeta: metav1alpha1.StatusMeta{
			// Figure out proper way to get the stack status and use my enum
			Status:  strcase.ToCamel(strings.ToLower(string(*outputs.Stacks[0].StackStatus))),
			Message: outputs.Stacks[0].StackStatusReason,
			StackID: string(*outputs.Stacks[0].StackId),
		},
	}
	statusCopy.Status = status
	_, err = c.awsclientset.CloudformationV1alpha1().Stacks(stackCopy.Namespace).UpdateStatus(statusCopy)
	return err
}

func createStackInputs(stack *cloudformation.Stack, notificationARN string, input *cfn.CreateStackInput) {
	input.SetCapabilities(stack.Spec.Capabilities)

	parameters := []*cfn.Parameter{}
	for key, value := range stack.Spec.Parameters {
		param := &cfn.Parameter{}
		param.SetParameterKey(key)
		param.SetParameterValue(value)
		parameters = append(parameters, param)
	}
	input.SetParameters(parameters)

	input.SetClientRequestToken(stack.Spec.ClientRequestToken)
	input.SetEnableTerminationProtection(stack.Spec.TerminationProtection)

	notificationARNs := []*string{}
	for _, notificationARN := range stack.Spec.NotificationARNs {
		notificationARNs = append(notificationARNs, notificationARN)
	}
	if notificationARN != "" {
		notificationARNs = append(notificationARNs, aws.String(notificationARN))
	}
	input.SetNotificationARNs(notificationARNs)

	onFailure := stack.Spec.OnFailure
	if onFailure == "" {
		onFailure = "DELETE"
	}
	input.SetOnFailure(onFailure)

	// TODO: Comeback and make this more thought out with cluster name and a hash
	input.SetStackName(namer(stack.ObjectMeta.Name, stack.ObjectMeta.Namespace))
	input.SetTemplateBody(stack.Spec.TemplateBody)

	tags := []*cfn.Tag{}
	for key, value := range stack.Spec.Tags {
		tag := &cfn.Tag{}
		tag.SetKey(key)
		tag.SetValue(*value)
		tags = append(tags, tag)
	}
	input.SetTags(tags)
}

func updateStackInputs(stack *cloudformation.Stack, notificationARN string, input *cfn.UpdateStackInput) {
	input.SetCapabilities(stack.Spec.Capabilities)

	parameters := []*cfn.Parameter{}
	for key, value := range stack.Spec.Parameters {
		param := &cfn.Parameter{}
		param.SetParameterKey(key)
		param.SetParameterValue(value)
		parameters = append(parameters, param)
	}
	input.SetParameters(parameters)

	// ClientRequestToken I expected to abe a stack specific token which would
	// help stop other stacks from using the same name but it doesn't seem to
	// work like that
	// input.SetClientRequestToken(stack.Spec.ClientRequestToken)

	notificationARNs := []*string{}
	for _, notificationARN := range stack.Spec.NotificationARNs {
		notificationARNs = append(notificationARNs, notificationARN)
	}
	if notificationARN != "" {
		notificationARNs = append(notificationARNs, aws.String(notificationARN))
	}
	input.SetNotificationARNs(notificationARNs)

	onFailure := stack.Spec.CloudFormationMeta.OnFailure
	if onFailure == "" {
		onFailure = "DELETE"
	}

	// TODO: Comeback and make this more thought out with cluster name and a hash
	input.SetStackName(namer(stack.ObjectMeta.Name, stack.ObjectMeta.Namespace))
	input.SetTemplateBody(stack.Spec.TemplateBody)

	tags := []*cfn.Tag{}
	for key, value := range stack.Spec.Tags {
		tag := &cfn.Tag{}
		tag.SetKey(key)
		tag.SetValue(*value)
		tags = append(tags, tag)
	}
	input.SetTags(tags)
}

func namer(name, namespace string) string {
	namerArr := []string{name, namespace}
	return strings.Join(namerArr, "-")
}
