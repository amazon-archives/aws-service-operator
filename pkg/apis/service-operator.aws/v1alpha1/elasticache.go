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
	AutoMinorVersionUpgrade         bool   `json:"autoMinorVersionUpgrade"`
	AZMode                          string `json:"azMode"`
	CacheNodeType                   string `json:"cacheNodeType"`
	CacheParameterGroupName         string `json:"cacheParameterGroupName"`
	CacheSubnetGroupName            string `json:"cacheSubnetGroupName"`
	Engine                          string `json:"engine"`
	EngineVersion                   string `json:"engineVersion"`
	NotificationTopicArn            string `json:"notificationTopicArn"`
	NumCacheNodes                   int    `json:"numCacheNodes"`
	Port                            int    `json:"port"`
	PreferredMaintenanceWindow      string `json:"preferredMaintenanceWindow"`
	PreferredAvailabilityZone       string `json:"preferredAvailabilityZone"`
	PreferredAvailabilityZones      string `json:"preferredAvailabilityZones"`
	SnapshotWindow                  string `json:"snapshotWindow"`
	VpcSecurityGroupIds             string `json:"vpcSecurityGroupIds"`
}

// ElastiCacheOutput defines the output resource for ElastiCache
type ElastiCacheOutput struct {
	RedisEndpointAddress         string `json:"redisEndpointAddress"`
	RedisEndpointPort            string `json:"redisEndpointPort"`
	ConfigurationEndpointAddress string `json:"configurationEndpoint"`
	ConfigurationEndpointPort    string `json:"configurationEndpointPort"`
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
