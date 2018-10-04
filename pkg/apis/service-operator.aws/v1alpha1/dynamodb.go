package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DynamoDB defines the base resource
type DynamoDB struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                DynamoDBSpec                `json:"spec"`
	Status              DynamoDBStatus              `json:"status"`
	Output              DynamoDBOutput              `json:"output"`
	AdditionalResources DynamoDBAdditionalResources `json:"additionalResources"`
}

// DynamoDBHashAttribute defines the HashAttribute resource for DynamoDB
type DynamoDBHashAttribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DynamoDBRangeAttribute defines the RangeAttribute resource for DynamoDB
type DynamoDBRangeAttribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DynamoDBSpec defines the Spec resource for DynamoDB
type DynamoDBSpec struct {
	CloudFormationTemplateName      string                 `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string                 `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int                    `json:"rollbackCount"`
	RangeAttribute                  DynamoDBRangeAttribute `json:"rangeAttribute"`
	ReadCapacityUnits               int                    `json:"readCapacityUnits"`
	WriteCapacityUnits              int                    `json:"writeCapacityUnits"`
	HashAttribute                   DynamoDBHashAttribute  `json:"hashAttribute"`
}

// DynamoDBOutput defines the output resource for DynamoDB
type DynamoDBOutput struct {
	TableName string `json:"tableName"`
	TableARN  string `json:"tableARN"`
}

// DynamoDBStatus holds the status of the Cloudformation template
type DynamoDBStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// DynamoDBAdditionalResources holds the additional resources
type DynamoDBAdditionalResources struct {
	ConfigMaps []string `json:"configMaps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DynamoDBList defines the list attribute for the DynamoDB type
type DynamoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DynamoDB `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addDynamoDBTypes)
}

func addDynamoDBTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&DynamoDB{},
		&DynamoDBList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
