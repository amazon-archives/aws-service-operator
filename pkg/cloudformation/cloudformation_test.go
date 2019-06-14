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
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	cloudformation "awsoperator.io/pkg/apis/cloudformation/v1alpha1"
	meta "awsoperator.io/pkg/apis/meta/v1alpha1"
	"awsoperator.io/pkg/controllerutils"
	"awsoperator.io/pkg/generated/clientset/versioned/fake"
	informers "awsoperator.io/pkg/generated/informers/externalversions"
	"awsoperator.io/pkg/queue"

	"awsoperator.io/pkg/testutils"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
	stackID            = "arn:aws:cloudformation:us-west-2:XXXXXXXXXXXX:stack/dyanmodb-stack/ca1ee2f3-ea89-4d14-9302-8d99b07cb3b1"
)

type fixture struct {
	t *testing.T

	kubeclient *k8sfake.Clientset
	client     *fake.Clientset

	stackLister []*cloudformation.Stack

	actions     []core.Action
	kubeactions []core.Action

	objects     []runtime.Object
	kubeobjects []runtime.Object
}

type mockToken struct {
	staticToken string
}

func (m mockToken) Generate() string {
	return m.staticToken
}

type mockCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m *mockCloudFormationClient) CreateStack(input *cfn.CreateStackInput) (*cfn.CreateStackOutput, error) {
	output := &cfn.CreateStackOutput{}
	output.SetStackId(stackID)
	return output, nil
}

func (m *mockCloudFormationClient) UpdateStack(input *cfn.UpdateStackInput) (*cfn.UpdateStackOutput, error) {
	output := &cfn.UpdateStackOutput{}
	output.SetStackId(stackID)
	return output, nil
}

func (m *mockCloudFormationClient) DescribeStacks(*cfn.DescribeStacksInput) (*cfn.DescribeStacksOutput, error) {
	describeStackOutput := &cfn.DescribeStacksOutput{}
	stack := &cfn.Stack{}
	output := &cfn.Output{}
	output.SetOutputKey("Name")
	output.SetOutputValue("test")
	stack.SetOutputs([]*cfn.Output{output})
	stack.SetStackStatus("UPDATE_IN_PROGRESS")
	stack.SetStackStatusReason("User initiated")
	stack.SetStackId(stackID)
	describeStackOutput.SetStacks([]*cfn.Stack{stack})
	return describeStackOutput, nil
}

type queueInformer struct {
	queue.InformerFactory
}

func mockQueueInformer() queue.InformerFactory {
	return &queueInformer{}
}

func (m queueInformer) AddMessageHander(queue.ControllerHandler) {
	return
}

func (m queueInformer) Start(stopCh <-chan struct{}) {
	return
}

func (m queueInformer) GetTopicARN() string {
	return "topic-arn"
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}
	return f
}

func newStack(name, cwName, cwNamespace, templateBody string) *cloudformation.Stack {
	stack := &cloudformation.Stack{
		TypeMeta: metav1.TypeMeta{APIVersion: cloudformation.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: cloudformation.StackSpec{
			CloudFormationMeta: meta.CloudFormationMeta{
				Region: "us-west-2",
			},
			TemplateBody: templateBody,
			Parameters: map[string]string{
				"HashKeyElementName": "test",
				"HashKeyElementType": "S",
				"ReadCapacityUnits":  "5",
				"WriteCapacityUnits": "5",
			},
		},
	}
	controllerutils.AddFinalizer(stack, stackDeletionFinalizerName)
	return stack
}

func (f *fixture) newController() (*Controller, informers.SharedInformerFactory) {
	f.client = fake.NewSimpleClientset(f.objects...)
	f.kubeclient = k8sfake.NewSimpleClientset(f.kubeobjects...)

	mockCloudformation := &mockCloudFormationClient{}

	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc())

	c, _ := New(f.kubeclient, f.client, i.Cloudformation().V1alpha1().Stacks(), mockQueueInformer(), map[string]cloudformationiface.CloudFormationAPI{"us-west-2": mockCloudformation})
	c.stackSynced = alwaysReady
	c.clientToken = mockToken{staticToken: "8f8be20f-9a4d-46a1-8936-b6fadf2a7734"}
	c.recorder = &record.FakeRecorder{}

	for _, a := range f.stackLister {
		i.Cloudformation().V1alpha1().Stacks().Informer().GetIndexer().Add(a)
	}

	return c, i
}

func (f *fixture) run(stackName string) {
	f.runController(stackName, true, false)
}

func (f *fixture) runController(stackName string, startInformers bool, expectError bool) {
	c, i := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		i.Start(stopCh)
	}

	err := c.syncHandler(stackName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing stack: %v", err)
	} else if expectError && err == nil {
		f.t.Error("expected error syncing stack, got nil")
	}

	actions := filterInformerActions(f.client.Actions())
	for i, action := range actions {
		if len(f.actions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(actions)-len(f.actions), actions[i:])
			break
		}

		expectedAction := f.actions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}

	k8sActions := filterInformerActions(f.kubeclient.Actions())
	for i, action := range k8sActions {
		if len(f.kubeactions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(k8sActions)-len(f.kubeactions), k8sActions[i:])
			break
		}

		expectedAction := f.kubeactions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.kubeactions) > len(k8sActions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.kubeactions)-len(k8sActions), f.kubeactions[len(k8sActions):])
	}
}

func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.CreateAction:
		e, _ := expected.(core.CreateAction)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expObject, object))
		}
	case core.UpdateAction:
		e, _ := expected.(core.UpdateAction)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expObject, object))
		}
	case core.PatchAction:
		e, _ := expected.(core.PatchAction)
		expPatch := e.GetPatch()
		patch := a.GetPatch()

		if !reflect.DeepEqual(expPatch, patch) {
			t.Errorf("Action %s %s has wrong patch\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expPatch, patch))
		}
	}
}

func (f *fixture) expectUpdateStackStatusAction(stack *cloudformation.Stack) {
	action := core.NewUpdateAction(schema.GroupVersionResource{Resource: "stacks"}, stack.Namespace, stack)
	f.actions = append(f.actions, action)
}

func filterInformerActions(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			(action.Matches("list", "stacks") ||
				action.Matches("watch", "stacks")) {
			continue
		}
		ret = append(ret, action)
	}

	return ret
}

func getKey(stack *cloudformation.Stack, t *testing.T) string {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(stack)
	if err != nil {
		t.Errorf("Unexpected error getting key for foo %v: %v", stack.Name, err)
		return ""
	}
	return key
}

func TestCreateStackWithoutClientRequestToken(t *testing.T) {
	f := newFixture(t)
	templateBody, _ := testutils.Asset("dynamodb.yaml")

	stack := newStack("test", "test-cw-name", "test-cw-namespace", string(templateBody))

	f.stackLister = append(f.stackLister, stack)
	f.objects = append(f.objects, stack)

	expStack := stack.DeepCopy()
	expStack.Spec.ClientRequestToken = "8f8be20f-9a4d-46a1-8936-b6fadf2a7734"

	f.expectUpdateStackStatusAction(expStack)
	f.run(getKey(stack, t))
}

func TestCreateStackWithClientRequestToken(t *testing.T) {
	f := newFixture(t)
	templateBody, _ := testutils.Asset("dynamodb.yaml")

	stack := newStack("test", "test-cw-name", "test-cw-namespace", string(templateBody))
	stack.Spec.ClientRequestToken = "8f8be20f-9a4d-46a1-8936-b6fadf2a7734"

	f.stackLister = append(f.stackLister, stack)
	f.objects = append(f.objects, stack)

	expStack := stack.DeepCopy()
	status := cloudformation.StackStatus{
		StatusMeta: meta.StatusMeta{
			Status:  meta.CreateInProgressStatus,
			StackID: stackID,
		},
	}
	expStack.Status = status
	f.expectUpdateStackStatusAction(expStack)
	f.run(getKey(stack, t))
}

func TestUpdateExistingStack(t *testing.T) {
	f := newFixture(t)
	templateBody, _ := testutils.Asset("dynamodb.yaml")

	stack := newStack("test", "test-cw-name", "test-cw-namespace", string(templateBody))
	stack.Spec.ClientRequestToken = "8f8be20f-9a4d-46a1-8936-b6fadf2a7734"
	status := cloudformation.StackStatus{
		StatusMeta: meta.StatusMeta{
			Status:  meta.CreateCompleteStatus,
			StackID: stackID,
		},
	}
	stack.Status = status

	f.stackLister = append(f.stackLister, stack)
	f.objects = append(f.objects, stack)

	expStack := stack.DeepCopy()
	expStack.Status.Status = meta.UpdateInProgressStatus

	f.expectUpdateStackStatusAction(expStack)

	expStack2 := expStack.DeepCopy()
	expStack2.Status.Message = aws.String("User initiated")

	expStack2.Outputs = map[string]string{"Name": "test"}
	f.expectUpdateStackStatusAction(expStack2)

	f.run(getKey(stack, t))
}

func TestCreateStackInputs(t *testing.T) {
	templateBody, _ := testutils.Asset("dynamodb.yaml")
	stack := newStack("test", "test-cw-name", "test-cw-namespace", string(templateBody))
	input := &cfn.CreateStackInput{}
	createStackInputs(stack, "", input)

	if *input.StackName != "test-default" {
		t.Errorf("StackName '%s' not equal to '%s'", *input.StackName, "test-default")
	}

	if len(input.Capabilities) != 0 {
		t.Errorf("Capabililities length '%d' not equal to '0", len(input.Capabilities))
	}

	if len(input.Parameters) != 4 {
		t.Errorf("Parameters length '%d' not equal to '4", len(input.Parameters))
	}
}
