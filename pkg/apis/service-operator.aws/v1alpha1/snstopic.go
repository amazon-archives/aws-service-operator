package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SNSTopic defines the base resource
type SNSTopic struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                SNSTopicSpec                `json:"spec"`
	Status              SNSTopicStatus              `json:"status"`
	Output              SNSTopicOutput              `json:"output"`
	AdditionalResources SNSTopicAdditionalResources `json:"additionalResources"`
}
type SNSTopicSpec struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
}

// SNSTopicOutput defines the output resource for SNSTopic
type SNSTopicOutput struct {
	TopicARN string `json:"topicARN"`
}

// SNSTopicStatus holds the status of the Cloudformation template
type SNSTopicStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// SNSTopicAdditionalResources holds the additional resources
type SNSTopicAdditionalResources struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SNSTopicList defines the list attribute for the SNSTopic type
type SNSTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SNSTopic `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addSNSTopicTypes)
}

func addSNSTopicTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&SNSTopic{},
		&SNSTopicList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
