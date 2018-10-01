package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SNSSubscription defines the base resource
type SNSSubscription struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                SNSSubscriptionSpec                `json:"spec"`
	Status              SNSSubscriptionStatus              `json:"status"`
	Output              SNSSubscriptionOutput              `json:"output"`
	AdditionalResources SNSSubscriptionAdditionalResources `json:"additionalResources"`
}

// SNSSubscriptionSpec defines the Spec resource for SNSSubscription
type SNSSubscriptionSpec struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
	TopicName                       string `json:"topicName"`
	Protocol                        string `json:"protocol"`
	Endpoint                        string `json:"endpoint"`
	QueueURL                        string `json:"queueURL"`
}

// SNSSubscriptionOutput defines the output resource for SNSSubscription
type SNSSubscriptionOutput struct {
	SubscriptionARN string `json:"subscriptionARN"`
}

// SNSSubscriptionStatus holds the status of the Cloudformation template
type SNSSubscriptionStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// SNSSubscriptionAdditionalResources holds the additional resources
type SNSSubscriptionAdditionalResources struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SNSSubscriptionList defines the list attribute for the SNSSubscription type
type SNSSubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SNSSubscription `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addSNSSubscriptionTypes)
}

func addSNSSubscriptionTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&SNSSubscription{},
		&SNSSubscriptionList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
