package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SQSQueue defines the base resource
type SQSQueue struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                SQSQueueSpec                `json:"spec"`
	Status              SQSQueueStatus              `json:"status"`
	Output              SQSQueueOutput              `json:"output"`
	AdditionalResources SQSQueueAdditionalResources `json:"additionalResources"`
}

// SQSQueueSpec defines the Spec resource for SQSQueue
type SQSQueueSpec struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
	ContentBasedDeduplication       bool   `json:"contentBasedDeduplication"`
	DelaySeconds                    int    `json:"delaySeconds"`
	MaximumMessageSize              int    `json:"maximumMessageSize"`
	MessageRetentionPeriod          int    `json:"messageRetentionPeriod"`
	ReceiveMessageWaitTimeSeconds   int    `json:"receiveMessageWaitTimeSeconds"`
	UsedeadletterQueue              bool   `json:"usedeadletterQueue"`
	VisibilityTimeout               int    `json:"visibilityTimeout"`
	FifoQueue                       bool   `json:"fifoQueue"`
}

// SQSQueueOutput defines the output resource for SQSQueue
type SQSQueueOutput struct {
	QueueURL            string `json:"queueURL"`
	QueueARN            string `json:"queueARN"`
	QueueName           string `json:"queueName"`
	DeadLetterQueueURL  string `json:"deadLetterQueueURL"`
	DeadLetterQueueARN  string `json:"deadLetterQueueARN"`
	DeadLetterQueueName string `json:"deadLetterQueueName"`
}

// SQSQueueStatus holds the status of the Cloudformation template
type SQSQueueStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// SQSQueueAdditionalResources holds the additional resources
type SQSQueueAdditionalResources struct {
	ConfigMaps []string `json:"configMaps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SQSQueueList defines the list attribute for the SQSQueue type
type SQSQueueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SQSQueue `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addSQSQueueTypes)
}

func addSQSQueueTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&SQSQueue{},
		&SQSQueueList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
