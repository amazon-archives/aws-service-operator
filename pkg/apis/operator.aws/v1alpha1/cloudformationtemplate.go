package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFormationTemplate defines the base resource
type CloudFormationTemplate struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Data                CloudFormationTemplateData                `json:"data"`
	Status              CloudFormationTemplateStatus              `json:"status"`
	Output              CloudFormationTemplateOutput              `json:"output"`
	AdditionalResources CloudFormationTemplateAdditionalResources `json:"additionalResources"`
}

// CloudFormationTemplateData defines the Data resource for CloudFormationTemplate
type CloudFormationTemplateData struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
	Key                             string `json:"key"`
	Template                        string `json:"template"`
}

// CloudFormationTemplateOutput defines the output resource for CloudFormationTemplate
type CloudFormationTemplateOutput struct {
	URL string `json:"url"`
}

// CloudFormationTemplateStatus holds the status of the Cloudformation template
type CloudFormationTemplateStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
}

// CloudFormationTemplateAdditionalResources holds the additional resources
type CloudFormationTemplateAdditionalResources struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudFormationTemplateList defines the list attribute for the CloudFormationTemplate type
type CloudFormationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CloudFormationTemplate `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addCloudFormationTemplateTypes)
}

func addCloudFormationTemplateTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CloudFormationTemplate{},
		&CloudFormationTemplateList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
