package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ECRRepository defines the base resource
type ECRRepository struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                ECRRepositorySpec                `json:"spec"`
	Status              ECRRepositoryStatus              `json:"status"`
	Output              ECRRepositoryOutput              `json:"output"`
	AdditionalResources ECRRepositoryAdditionalResources `json:"additionalResources"`
}

// ECRRepositorySpec defines the Spec resource for ECRRepository
type ECRRepositorySpec struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
}

// ECRRepositoryOutput defines the output resource for ECRRepository
type ECRRepositoryOutput struct {
	RepositoryName string `json:"repositoryName"`
	RepositoryARN  string `json:"repositoryARN"`
	RepositoryURL  string `json:"repositoryURL"`
}

// ECRRepositoryStatus holds the status of the Cloudformation template
type ECRRepositoryStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// ECRRepositoryAdditionalResources holds the additional resources
type ECRRepositoryAdditionalResources struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ECRRepositoryList defines the list attribute for the ECRRepository type
type ECRRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ECRRepository `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addECRRepositoryTypes)
}

func addECRRepositoryTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ECRRepository{},
		&ECRRepositoryList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
