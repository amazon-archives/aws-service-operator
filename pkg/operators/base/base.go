package base

import (
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/operators/cloudformationtemplate"
	"github.com/awslabs/aws-service-operator/pkg/operators/dynamodb"
	"github.com/awslabs/aws-service-operator/pkg/operators/ecrrepository"
	"github.com/awslabs/aws-service-operator/pkg/operators/s3bucket"
	"github.com/awslabs/aws-service-operator/pkg/operators/snssubscription"
	"github.com/awslabs/aws-service-operator/pkg/operators/snstopic"
	"github.com/awslabs/aws-service-operator/pkg/operators/sqsqueue"
)

type base struct {
	config *config.Config
}

func New(
	config *config.Config,
) *base {
	return &base{
		config: config,
	}
}

func (b *base) Watch(namespace string, stopCh chan struct{}) (err error) {
	if b.config.Resources["cloudformationtemplate"] {
		cloudformationtemplateoperator := cloudformationtemplate.NewOperator(b.config)
		err = cloudformationtemplateoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["dynamodb"] {
		dynamodboperator := dynamodb.NewOperator(b.config)
		err = dynamodboperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["ecrrepository"] {
		ecrrepositoryoperator := ecrrepository.NewOperator(b.config)
		err = ecrrepositoryoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["s3bucket"] {
		s3bucketoperator := s3bucket.NewOperator(b.config)
		err = s3bucketoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["snssubscription"] {
		snssubscriptionoperator := snssubscription.NewOperator(b.config)
		err = snssubscriptionoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["snstopic"] {
		snstopicoperator := snstopic.NewOperator(b.config)
		err = snstopicoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}
	if b.config.Resources["sqsqueue"] {
		sqsqueueoperator := sqsqueue.NewOperator(b.config)
		err = sqsqueueoperator.StartWatch(namespace, stopCh)
		if err != nil {
			return err
		}
	}

	return nil
}
