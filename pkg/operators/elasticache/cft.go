// >>>>>>> DO NOT EDIT THIS FILE <<<<<<<<<<
// This file is autogenerated via `aws-operator-codegen process`
// If you'd like the change anything about this file make edits to the .templ
// file in the pkg/codegen/assets directory.

package elasticache

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awsV1alpha1 "github.com/awslabs/aws-service-operator/pkg/apis/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/helpers"
)

// New generates a new object
func New(config config.Config, elasticache *awsV1alpha1.ElastiCache, topicARN string) *Cloudformation {
	return &Cloudformation{
		ElastiCache: elasticache,
		config:      config,
		topicARN:    topicARN,
	}
}

// Cloudformation defines the elasticache cfts
type Cloudformation struct {
	config      config.Config
	ElastiCache *awsV1alpha1.ElastiCache
	topicARN    string
}

// StackName returns the name of the stack based on the aws-operator-config
func (s *Cloudformation) StackName() string {
	return helpers.StackName(s.config.ClusterName, "elasticache", s.ElastiCache.Name, s.ElastiCache.Namespace)
}

// GetOutputs return the stack outputs from the DescribeStacks call
func (s *Cloudformation) GetOutputs() (map[string]string, error) {
	outputs := map[string]string{}
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DescribeStacksInput{
		StackName: aws.String(s.StackName()),
	}

	output, err := svc.DescribeStacks(&stackInputs)
	if err != nil {
		return nil, err
	}
	// Not sure if this is even possible
	if len(output.Stacks) != 1 {
		return nil, errors.New("no stacks returned with that stack name")
	}

	for _, out := range output.Stacks[0].Outputs {
		outputs[*out.OutputKey] = *out.OutputValue
	}

	return outputs, err
}

// CreateStack will create the stack with the supplied params
func (s *Cloudformation) CreateStack() (output *cloudformation.CreateStackOutput, err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	cftemplate := helpers.GetCloudFormationTemplate(s.config, "elasticache", s.ElastiCache.Spec.CloudFormationTemplateName, s.ElastiCache.Spec.CloudFormationTemplateNamespace)

	stackInputs := cloudformation.CreateStackInput{
		StackName:   aws.String(s.StackName()),
		TemplateURL: aws.String(cftemplate),
		NotificationARNs: []*string{
			aws.String(s.topicARN),
		},
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
	}

	resourceName := helpers.CreateParam("ResourceName", s.ElastiCache.Name)
	resourceVersion := helpers.CreateParam("ResourceVersion", s.ElastiCache.ResourceVersion)
	namespace := helpers.CreateParam("Namespace", s.ElastiCache.Namespace)
	clusterName := helpers.CreateParam("ClusterName", s.config.ClusterName)
	elastiCacheClusterName := helpers.CreateParam("ClusterName", helpers.Stringify(s.ElastiCache.Name))
	autoMinorVersionUpgradeTemp := "{{.Obj.Spec.AutoMinorVersionUpgrade}}"
	autoMinorVersionUpgradeValue, err := helpers.Templatize(autoMinorVersionUpgradeTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	autoMinorVersionUpgrade := helpers.CreateParam("AutoMinorVersionUpgrade", helpers.Stringify(autoMinorVersionUpgradeValue))
	azModeTemp := "{{.Obj.Spec.AZMode}}"
	azModeValue, err := helpers.Templatize(azModeTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	azMode := helpers.CreateParam("AZMode", helpers.Stringify(azModeValue))
	cacheNodeTypeTemp := "{{.Obj.Spec.CacheNodeType}}"
	cacheNodeTypeValue, err := helpers.Templatize(cacheNodeTypeTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheNodeType := helpers.CreateParam("CacheNodeType", helpers.Stringify(cacheNodeTypeValue))
	cacheParameterGroupNameTemp := "{{.Obj.Spec.CacheParameterGroupName}}"
	cacheParameterGroupNameValue, err := helpers.Templatize(cacheParameterGroupNameTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheParameterGroupName := helpers.CreateParam("CacheParameterGroupName", helpers.Stringify(cacheParameterGroupNameValue))
	cacheSubnetGroupNameTemp := "{{.Obj.Spec.CacheSubnetGroupName}}"
	cacheSubnetGroupNameValue, err := helpers.Templatize(cacheSubnetGroupNameTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheSubnetGroupName := helpers.CreateParam("CacheSubnetGroupName", helpers.Stringify(cacheSubnetGroupNameValue))
	engineTemp := "{{.Obj.Spec.Engine}}"
	engineValue, err := helpers.Templatize(engineTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	engine := helpers.CreateParam("Engine", helpers.Stringify(engineValue))
	engineVersionTemp := "{{.Obj.Spec.EngineVersion}}"
	engineVersionValue, err := helpers.Templatize(engineVersionTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	engineVersion := helpers.CreateParam("EngineVersion", helpers.Stringify(engineVersionValue))
	notificationTopicArnTemp := "{{.Obj.Spec.NotificationTopicArn}}"
	notificationTopicArnValue, err := helpers.Templatize(notificationTopicArnTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	notificationTopicArn := helpers.CreateParam("NotificationTopicArn", helpers.Stringify(notificationTopicArnValue))
	numCacheNodesTemp := "{{.Obj.Spec.NumCacheNodes}}"
	numCacheNodesValue, err := helpers.Templatize(numCacheNodesTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	numCacheNodes := helpers.CreateParam("NumCacheNodes", helpers.Stringify(numCacheNodesValue))
	portTemp := "{{.Obj.Spec.Port}}"
	portValue, err := helpers.Templatize(portTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	port := helpers.CreateParam("Port", helpers.Stringify(portValue))
	preferredMaintenanceWindowTemp := "{{.Obj.Spec.PreferredMaintenanceWindow}}"
	preferredMaintenanceWindowValue, err := helpers.Templatize(preferredMaintenanceWindowTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredMaintenanceWindow := helpers.CreateParam("PreferredMaintenanceWindow", helpers.Stringify(preferredMaintenanceWindowValue))
	preferredAvailabilityZoneTemp := "{{.Obj.Spec.PreferredAvailabilityZone}}"
	preferredAvailabilityZoneValue, err := helpers.Templatize(preferredAvailabilityZoneTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredAvailabilityZone := helpers.CreateParam("PreferredAvailabilityZone", helpers.Stringify(preferredAvailabilityZoneValue))
	preferredAvailabilityZonesTemp := "{{.Obj.Spec.PreferredAvailabilityZones}}"
	preferredAvailabilityZonesValue, err := helpers.Templatize(preferredAvailabilityZonesTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredAvailabilityZones := helpers.CreateParam("PreferredAvailabilityZones", helpers.Stringify(preferredAvailabilityZonesValue))
	snapshotWindowTemp := "{{.Obj.Spec.SnapshotWindow}}"
	snapshotWindowValue, err := helpers.Templatize(snapshotWindowTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	snapshotWindow := helpers.CreateParam("SnapshotWindow", helpers.Stringify(snapshotWindowValue))
	vpcSecurityGroupIdsTemp := "{{.Obj.Spec.VpcSecurityGroupIds}}"
	vpcSecurityGroupIdsValue, err := helpers.Templatize(vpcSecurityGroupIdsTemp, helpers.Data{Obj: s.ElastiCache, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	vpcSecurityGroupIds := helpers.CreateParam("VpcSecurityGroupIds", helpers.Stringify(vpcSecurityGroupIdsValue))

	parameters := []*cloudformation.Parameter{}
	parameters = append(parameters, resourceName)
	parameters = append(parameters, resourceVersion)
	parameters = append(parameters, namespace)
	parameters = append(parameters, clusterName)
	parameters = append(parameters, elastiCacheClusterName)
	parameters = append(parameters, autoMinorVersionUpgrade)
	parameters = append(parameters, azMode)
	parameters = append(parameters, cacheNodeType)
	parameters = append(parameters, cacheParameterGroupName)
	parameters = append(parameters, cacheSubnetGroupName)
	parameters = append(parameters, engine)
	parameters = append(parameters, engineVersion)
	parameters = append(parameters, notificationTopicArn)
	parameters = append(parameters, numCacheNodes)
	parameters = append(parameters, port)
	parameters = append(parameters, preferredMaintenanceWindow)
	parameters = append(parameters, preferredAvailabilityZone)
	parameters = append(parameters, preferredAvailabilityZones)
	parameters = append(parameters, snapshotWindow)
	parameters = append(parameters, vpcSecurityGroupIds)

	stackInputs.SetParameters(parameters)

	resourceNameTag := helpers.CreateTag("ResourceName", s.ElastiCache.Name)
	resourceVersionTag := helpers.CreateTag("ResourceVersion", s.ElastiCache.ResourceVersion)
	namespaceTag := helpers.CreateTag("Namespace", s.ElastiCache.Namespace)
	clusterNameTag := helpers.CreateTag("ClusterName", s.config.ClusterName)

	tags := []*cloudformation.Tag{}
	tags = append(tags, resourceNameTag)
	tags = append(tags, resourceVersionTag)
	tags = append(tags, namespaceTag)
	tags = append(tags, clusterNameTag)

	stackInputs.SetTags(tags)

	output, err = svc.CreateStack(&stackInputs)
	return
}

// UpdateStack will update the existing stack
func (s *Cloudformation) UpdateStack(updated *awsV1alpha1.ElastiCache) (output *cloudformation.UpdateStackOutput, err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	cftemplate := helpers.GetCloudFormationTemplate(s.config, "elasticache", updated.Spec.CloudFormationTemplateName, updated.Spec.CloudFormationTemplateNamespace)

	stackInputs := cloudformation.UpdateStackInput{
		StackName:   aws.String(s.StackName()),
		TemplateURL: aws.String(cftemplate),
		NotificationARNs: []*string{
			aws.String(s.topicARN),
		},
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
	}

	resourceName := helpers.CreateParam("ResourceName", s.ElastiCache.Name)
	resourceVersion := helpers.CreateParam("ResourceVersion", s.ElastiCache.ResourceVersion)
	namespace := helpers.CreateParam("Namespace", s.ElastiCache.Namespace)
	clusterName := helpers.CreateParam("ClusterName", s.config.ClusterName)
	elastiCacheClusterName := helpers.CreateParam("ClusterName", helpers.Stringify(s.ElastiCache.Name))
	autoMinorVersionUpgradeTemp := "{{.Obj.Spec.AutoMinorVersionUpgrade}}"
	autoMinorVersionUpgradeValue, err := helpers.Templatize(autoMinorVersionUpgradeTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	autoMinorVersionUpgrade := helpers.CreateParam("AutoMinorVersionUpgrade", helpers.Stringify(autoMinorVersionUpgradeValue))
	azModeTemp := "{{.Obj.Spec.AZMode}}"
	azModeValue, err := helpers.Templatize(azModeTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	azMode := helpers.CreateParam("AZMode", helpers.Stringify(azModeValue))
	cacheNodeTypeTemp := "{{.Obj.Spec.CacheNodeType}}"
	cacheNodeTypeValue, err := helpers.Templatize(cacheNodeTypeTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheNodeType := helpers.CreateParam("CacheNodeType", helpers.Stringify(cacheNodeTypeValue))
	cacheParameterGroupNameTemp := "{{.Obj.Spec.CacheParameterGroupName}}"
	cacheParameterGroupNameValue, err := helpers.Templatize(cacheParameterGroupNameTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheParameterGroupName := helpers.CreateParam("CacheParameterGroupName", helpers.Stringify(cacheParameterGroupNameValue))
	cacheSubnetGroupNameTemp := "{{.Obj.Spec.CacheSubnetGroupName}}"
	cacheSubnetGroupNameValue, err := helpers.Templatize(cacheSubnetGroupNameTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	cacheSubnetGroupName := helpers.CreateParam("CacheSubnetGroupName", helpers.Stringify(cacheSubnetGroupNameValue))
	engineTemp := "{{.Obj.Spec.Engine}}"
	engineValue, err := helpers.Templatize(engineTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	engine := helpers.CreateParam("Engine", helpers.Stringify(engineValue))
	engineVersionTemp := "{{.Obj.Spec.EngineVersion}}"
	engineVersionValue, err := helpers.Templatize(engineVersionTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	engineVersion := helpers.CreateParam("EngineVersion", helpers.Stringify(engineVersionValue))
	notificationTopicArnTemp := "{{.Obj.Spec.NotificationTopicArn}}"
	notificationTopicArnValue, err := helpers.Templatize(notificationTopicArnTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	notificationTopicArn := helpers.CreateParam("NotificationTopicArn", helpers.Stringify(notificationTopicArnValue))
	numCacheNodesTemp := "{{.Obj.Spec.NumCacheNodes}}"
	numCacheNodesValue, err := helpers.Templatize(numCacheNodesTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	numCacheNodes := helpers.CreateParam("NumCacheNodes", helpers.Stringify(numCacheNodesValue))
	portTemp := "{{.Obj.Spec.Port}}"
	portValue, err := helpers.Templatize(portTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	port := helpers.CreateParam("Port", helpers.Stringify(portValue))
	preferredMaintenanceWindowTemp := "{{.Obj.Spec.PreferredMaintenanceWindow}}"
	preferredMaintenanceWindowValue, err := helpers.Templatize(preferredMaintenanceWindowTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredMaintenanceWindow := helpers.CreateParam("PreferredMaintenanceWindow", helpers.Stringify(preferredMaintenanceWindowValue))
	preferredAvailabilityZoneTemp := "{{.Obj.Spec.PreferredAvailabilityZone}}"
	preferredAvailabilityZoneValue, err := helpers.Templatize(preferredAvailabilityZoneTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredAvailabilityZone := helpers.CreateParam("PreferredAvailabilityZone", helpers.Stringify(preferredAvailabilityZoneValue))
	preferredAvailabilityZonesTemp := "{{.Obj.Spec.PreferredAvailabilityZones}}"
	preferredAvailabilityZonesValue, err := helpers.Templatize(preferredAvailabilityZonesTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	preferredAvailabilityZones := helpers.CreateParam("PreferredAvailabilityZones", helpers.Stringify(preferredAvailabilityZonesValue))
	snapshotWindowTemp := "{{.Obj.Spec.SnapshotWindow}}"
	snapshotWindowValue, err := helpers.Templatize(snapshotWindowTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	snapshotWindow := helpers.CreateParam("SnapshotWindow", helpers.Stringify(snapshotWindowValue))
	vpcSecurityGroupIdsTemp := "{{.Obj.Spec.VpcSecurityGroupIds}}"
	vpcSecurityGroupIdsValue, err := helpers.Templatize(vpcSecurityGroupIdsTemp, helpers.Data{Obj: updated, Config: s.config, Helpers: helpers.New()})
	if err != nil {
		return output, err
	}
	vpcSecurityGroupIds := helpers.CreateParam("VpcSecurityGroupIds", helpers.Stringify(vpcSecurityGroupIdsValue))

	parameters := []*cloudformation.Parameter{}
	parameters = append(parameters, resourceName)
	parameters = append(parameters, resourceVersion)
	parameters = append(parameters, namespace)
	parameters = append(parameters, clusterName)
	parameters = append(parameters, elastiCacheClusterName)
	parameters = append(parameters, autoMinorVersionUpgrade)
	parameters = append(parameters, azMode)
	parameters = append(parameters, cacheNodeType)
	parameters = append(parameters, cacheParameterGroupName)
	parameters = append(parameters, cacheSubnetGroupName)
	parameters = append(parameters, engine)
	parameters = append(parameters, engineVersion)
	parameters = append(parameters, notificationTopicArn)
	parameters = append(parameters, numCacheNodes)
	parameters = append(parameters, port)
	parameters = append(parameters, preferredMaintenanceWindow)
	parameters = append(parameters, preferredAvailabilityZone)
	parameters = append(parameters, preferredAvailabilityZones)
	parameters = append(parameters, snapshotWindow)
	parameters = append(parameters, vpcSecurityGroupIds)

	stackInputs.SetParameters(parameters)

	resourceNameTag := helpers.CreateTag("ResourceName", s.ElastiCache.Name)
	resourceVersionTag := helpers.CreateTag("ResourceVersion", s.ElastiCache.ResourceVersion)
	namespaceTag := helpers.CreateTag("Namespace", s.ElastiCache.Namespace)
	clusterNameTag := helpers.CreateTag("ClusterName", s.config.ClusterName)

	tags := []*cloudformation.Tag{}
	tags = append(tags, resourceNameTag)
	tags = append(tags, resourceVersionTag)
	tags = append(tags, namespaceTag)
	tags = append(tags, clusterNameTag)

	stackInputs.SetTags(tags)

	output, err = svc.UpdateStack(&stackInputs)
	return
}

// DeleteStack will delete the stack
func (s *Cloudformation) DeleteStack() (err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DeleteStackInput{}
	stackInputs.SetStackName(s.StackName())

	_, err = svc.DeleteStack(&stackInputs)
	return
}

// WaitUntilStackDeleted will delete the stack
func (s *Cloudformation) WaitUntilStackDeleted() (err error) {
	sess := s.config.AWSSession
	svc := cloudformation.New(sess)

	stackInputs := cloudformation.DescribeStacksInput{
		StackName: aws.String(s.StackName()),
	}

	err = svc.WaitUntilStackDeleteComplete(&stackInputs)
	return
}
