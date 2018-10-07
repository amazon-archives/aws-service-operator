package base

import (
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/operators/cloudformationtemplate"
	"github.com/awslabs/aws-service-operator/pkg/operators/dynamodb"
	"github.com/awslabs/aws-service-operator/pkg/operators/ecrrepository"
	"github.com/awslabs/aws-service-operator/pkg/operators/s3bucket"
	"github.com/awslabs/aws-service-operator/pkg/operators/snssubscription"
	"github.com/awslabs/aws-service-operator/pkg/operators/snstopic"
	"github.com/awslabs/aws-service-operator/pkg/operators/sqsqueue"
	opkit "github.com/christopherhein/operator-kit"
)

type base struct {
	config       *config.Config
	context      *opkit.Context
	awsClientset awsclient.ServiceoperatorV1alpha1Interface
}

func New(
	config *config.Config,
	context *opkit.Context,
	awsClientset awsclient.ServiceoperatorV1alpha1Interface,
) *base {
	return &base{
		config:       config,
		context:      context,
		awsClientset: awsClientset,
	}
}

func (b *base) Watch(namespace string, stopCh chan struct{}) (err error) {
	cloudformationtemplateoperator := cloudformationtemplate.NewOperator(b.config, b.context, b.awsClientset)
	err = cloudformationtemplateoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	dynamodboperator := dynamodb.NewOperator(b.config, b.context, b.awsClientset)
	err = dynamodboperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	ecrrepositoryoperator := ecrrepository.NewOperator(b.config, b.context, b.awsClientset)
	err = ecrrepositoryoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	s3bucketoperator := s3bucket.NewOperator(b.config, b.context, b.awsClientset)
	err = s3bucketoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	snssubscriptionoperator := snssubscription.NewOperator(b.config, b.context, b.awsClientset)
	err = snssubscriptionoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	snstopicoperator := snstopic.NewOperator(b.config, b.context, b.awsClientset)
	err = snstopicoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}
	sqsqueueoperator := sqsqueue.NewOperator(b.config, b.context, b.awsClientset)
	err = sqsqueueoperator.StartWatch(namespace, stopCh)
	if err != nil {
		return err
	}

	return nil
}
