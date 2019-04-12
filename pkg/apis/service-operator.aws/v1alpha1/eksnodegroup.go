package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroup defines the base resource
type EKSNodeGroup struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                EKSNodeGroupSpec                `json:"spec"`
	Status              EKSNodeGroupStatus              `json:"status"`
	Output              EKSNodeGroupOutput              `json:"output"`
	AdditionalResources EKSNodeGroupAdditionalResources `json:"additionalResources"`
}

// EKSNodeGroupAutoScalingGroup defines the AutoScalingGroup resource for EKSNodeGroup
type EKSNodeGroupAutoScalingGroup struct {
	NodeAutoScalingGroupMinSize         int `json:"minSize"`
	NodeAutoScalingGroupMaxSize         int `json:"maxSize"`
	NodeAutoScalingGroupDesiredCapacity int `json:"desiredCapacity"`
}

// EKSNodeGroupNetworking defines the Networking resource for EKSNodeGroup
type EKSNodeGroupNetworking struct {
	ClusterControlPlaneSecurityGroup string `json:"controlPlaneSecurityGroup"`
	VpcId                            string `json:"vpcID"`
	Subnets                          string `json:"subnets"`
}

// EKSNodeGroupNode defines the Node resource for EKSNodeGroup
type EKSNodeGroupNode struct {
	NodeImageId        string `json:"imageID"`
	KeyName            string `json:"keyName"`
	NodeInstanceType   string `json:"instanceType"`
	NodeVolumeSize     int    `json:"volumeSize"`
	BootstrapArguments string `json:"bootstrapArguments"`
}

// EKSNodeGroupSpec defines the Spec resource for EKSNodeGroup
type EKSNodeGroupSpec struct {
	CloudFormationTemplateName      string                       `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string                       `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int                          `json:"rollbackCount"`
	Node                            EKSNodeGroupNode             `json:"node"`
	AutoScalingGroup                EKSNodeGroupAutoScalingGroup `json:"autoScalingGroup"`
	Networking                      EKSNodeGroupNetworking       `json:"networking"`
}

// EKSNodeGroupOutput defines the output resource for EKSNodeGroup
type EKSNodeGroupOutput struct {
	NodeInstanceRole  string `json:"nodeInstanceRole"`
	NodeSecurityGroup string `json:"nodeSecurityGroup"`
}

// EKSNodeGroupStatus holds the status of the Cloudformation template
type EKSNodeGroupStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// EKSNodeGroupAdditionalResources holds the additional resources
type EKSNodeGroupAdditionalResources struct {
	ConfigMaps []string `json:"configMaps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EKSNodeGroupList defines the list attribute for the EKSNodeGroup type
type EKSNodeGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EKSNodeGroup `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addEKSNodeGroupTypes)
}

func addEKSNodeGroupTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&EKSNodeGroup{},
		&EKSNodeGroupList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
