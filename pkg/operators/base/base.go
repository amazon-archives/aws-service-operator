package base

import (
	"context"
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

func (b *base) Watch(ctx context.Context, namespace string) {
	if b.config.Resources["cloudformationtemplate"] {
		cloudformationtemplateoperator := cloudformationtemplate.NewOperator(b.config)
		go cloudformationtemplateoperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["dynamodb"] {
		dynamodboperator := dynamodb.NewOperator(b.config)
		go dynamodboperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["ecrrepository"] {
		ecrrepositoryoperator := ecrrepository.NewOperator(b.config)
		go ecrrepositoryoperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["s3bucket"] {
		s3bucketoperator := s3bucket.NewOperator(b.config)
		go s3bucketoperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["snssubscription"] {
		snssubscriptionoperator := snssubscription.NewOperator(b.config)
		go snssubscriptionoperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["snstopic"] {
		snstopicoperator := snstopic.NewOperator(b.config)
		go snstopicoperator.StartWatch(ctx, namespace)
	}
	if b.config.Resources["sqsqueue"] {
		sqsqueueoperator := sqsqueue.NewOperator(b.config)
		go sqsqueueoperator.StartWatch(ctx, namespace)
	}

	<-ctx.Done()
}
