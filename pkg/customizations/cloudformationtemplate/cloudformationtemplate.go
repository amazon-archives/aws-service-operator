package cloudformationtemplate

import (
	"bytes"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	awsV1alpha1 "github.com/awslabs/aws-service-operator/pkg/apis/service-operator.aws/v1alpha1"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OnAdd will be fired when you add a new CFT
func OnAdd(config config.Config, cft *awsV1alpha1.CloudFormationTemplate) {
	logger := config.Logger

	updateOutput(config, cft, "UPLOAD_IN_PROGRESS", "")
	err := addFileToS3(config, config.Bucket, cft.Data.Key, cft.Data.Template)
	if err != nil {
		logger.WithError(err).Error("error uploading cloudformation")
		updateOutput(config, cft, "UPLOAD_FAILED", err.Error())
	}
	logger.Infof("added cloudformationtemplate '%s'", cft.Name)
	updateOutput(config, cft, "UPLOAD_COMPLETE", "")
}

// OnUpdate will be fired when you update a CFT
func OnUpdate(config config.Config, oldcft *awsV1alpha1.CloudFormationTemplate, newcft *awsV1alpha1.CloudFormationTemplate) {
	if !reflect.DeepEqual(oldcft.Data, newcft.Data) {
		logger := config.Logger

		updateOutput(config, newcft, "UPLOAD_IN_PROGRESS", "")
		err := addFileToS3(config, config.Bucket, oldcft.Data.Key, newcft.Data.Template)
		if err != nil {
			logger.WithError(err).Error("error uploading cloudformation")
			updateOutput(config, newcft, "UPLOAD_FAILED", err.Error())
		}
		logger.Infof("updated cloudformationtemplate '%s'", oldcft.Name)
		updateOutput(config, newcft, "UPDATE_COMPLETE", "")
	}
}

// OnDelete will be fired when you delete a CFT
func OnDelete(config config.Config, cft *awsV1alpha1.CloudFormationTemplate) {
	logger := config.Logger
	sess := config.AWSSession
	_, err := s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(cft.Data.Key),
	})
	if err != nil {
		logger.WithError(err).Error("error deleting cloudformation")
	}
	logger.Infof("deleted cloudformationtemplate '%s'", cft.Name)
}

func addFileToS3(config config.Config, bucket string, filename string, template string) error {
	buffer := []byte(template)
	sess := config.AWSSession
	svc := s3.New(sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(buffer),
	})
	return err
}

func updateOutput(config config.Config, cft *awsV1alpha1.CloudFormationTemplate, status string, reason string) error {
	logger := config.Logger
	clientSet, _ := awsclient.NewForConfig(config.RESTConfig)
	resource, err := clientSet.CloudFormationTemplates(cft.Namespace).Get(cft.Name, metav1.GetOptions{})
	if err != nil {
		logger.WithError(err).Error("error getting cloudformation template")
		return err
	}

	resourceCopy := resource.DeepCopy()
	resourceCopy.Status.ResourceStatus = status
	resourceCopy.Status.ResourceStatusReason = reason
	resourceCopy.Output.URL = "https://s3." + config.Region + ".amazonaws.com/" + config.Bucket + "/" + cft.Data.Key

	_, err = clientSet.CloudFormationTemplates(cft.Namespace).Update(resourceCopy)
	if err != nil {
		logger.WithError(err).Error("error updating resource")
		return err
	}
	return nil
}
