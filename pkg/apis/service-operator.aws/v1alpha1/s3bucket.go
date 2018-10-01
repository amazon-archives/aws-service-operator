package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// S3Bucket defines the base resource
type S3Bucket struct {
	metav1.TypeMeta     `json:",inline"`
	metav1.ObjectMeta   `json:"metadata"`
	Spec                S3BucketSpec                `json:"spec"`
	Status              S3BucketStatus              `json:"status"`
	Output              S3BucketOutput              `json:"output"`
	AdditionalResources S3BucketAdditionalResources `json:"additionalResources"`
}

// S3BucketLogging defines the Logging resource for S3Bucket
type S3BucketLogging struct {
	Enabled bool   `json:"enabled"`
	Prefix  string `json:"prefix"`
}

// S3BucketSpec defines the Spec resource for S3Bucket
type S3BucketSpec struct {
	CloudFormationTemplateName      string          `json:"cloudFormationTemplateName"`
	CloudFormationTemplateNamespace string          `json:"cloudFormationTemplateNamespace"`
	RollbackCount                   int             `json:"rollbackCount"`
	Versioning                      bool            `json:"versioning"`
	AccessControl                   string          `json:"accessControl"`
	Logging                         S3BucketLogging `json:"logging"`
	Website                         S3BucketWebsite `json:"website"`
}

// S3BucketWebsite defines the Website resource for S3Bucket
type S3BucketWebsite struct {
	Enabled   bool   `json:"enabled"`
	IndexPage string `json:"indexPage"`
	ErrorPage string `json:"errorPage"`
}

// S3BucketOutput defines the output resource for S3Bucket
type S3BucketOutput struct {
	BucketName string `json:"bucketName"`
	BucketARN  string `json:"bucketARN"`
	WebsiteURL string `json:"websiteURL"`
}

// S3BucketStatus holds the status of the Cloudformation template
type S3BucketStatus struct {
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
	StackID              string `json:"stackID"`
}

// S3BucketAdditionalResources holds the additional resources
type S3BucketAdditionalResources struct {
	Services   []string `json:"services"`
	ConfigMaps []string `json:"configMaps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// S3BucketList defines the list attribute for the S3Bucket type
type S3BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []S3Bucket `json:"items"`
}

func init() {
	localSchemeBuilder.Register(addS3BucketTypes)
}

func addS3BucketTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&S3Bucket{},
		&S3BucketList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
