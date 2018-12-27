package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElastiCache defines the base resource
type ElastiCache struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                ElastiCacheSpec                `json:"spec"`
	Status              ElastiCacheStatus              `json:"status"`
	Output              ElastiCacheOutput              `json:"output"`
	AdditionalResources ElastiCacheAdditionalResources `json:"additionalResources"`
}

// ElastiCacheSpec defines the Spec resource for ElastiCache
type ElastiCacheSpec struct {
	CloudFormationTemplateName      string `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int    `json:"rollbackCount"`
	AutoMinorVersionUpgrade         bool   `json:"AutoMinorVersionUpgrade"`
	AZMode                          string `json:"AZMode"`
	CacheNodeType                   string `json:"CacheNodeType"`
	CacheParameterGroupName         string `json:"CacheParameterGroupName"`
	CacheSubnetGroupName            string `json:"CacheSubnetGroupName"`
	Engine                          string `json:"Engine"`
	EngineVersion                   string `json:"EngineVersion"`
	NotificationTopicArn            string `json:"NotificationTopicArn"`
	NumCacheNodes                   int    `json:"NumCacheNodes"`
	Port                            int    `json:"Port"`
	PreferredMaintenanceWindow      string `json:"PreferredMaintenanceWindow"`
	PreferredAvailabilityZone       string `json:"PreferredAvailabilityZone"`
	PreferredAvailabilityZones      string `json:"PreferredAvailabilityZones"`
	SnapshotWindow                  string `json:"SnapshotWindow"`
	VpcSecurityGroupIds             string `json:"VpcSecurityGroupIds"`
}

// ElastiCacheOutput defines the output resource for ElastiCache
type ElastiCacheOutput struct {
	RedisEndpointAddress         string `json:"RedisEndpointAddress"`
	RedisEndpointPort            string `json:"RedisEndpointPort"`
	ConfigurationEndpointAddress string `json:"ConfigurationEndpoint"`
	ConfigurationEndpointPort    string `json:"ConfigurationEndpointPort"`
}

// ElastiCacheStatus holds the status of the Cloudformation template
type ElastiCacheStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// ElastiCacheAdditionalResources holds the additional resources
type ElastiCacheAdditionalResources struct {
	Services []string `json:"services"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElastiCacheList defines the list attribute for the ElastiCache type
type ElastiCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ElastiCache `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addElastiCacheTypes)
}

func addElastiCacheTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ElastiCache{},
		&ElastiCacheList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
