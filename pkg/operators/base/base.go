package base

import (
	"context"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/operators/cloudformationtemplate"
	"github.com/awslabs/aws-service-operator/pkg/operators/dynamodb"
	"github.com/awslabs/aws-service-operator/pkg/operators/ecrrepository"
	"github.com/awslabs/aws-service-operator/pkg/operators/elasticache"
	"github.com/awslabs/aws-service-operator/pkg/operators/s3bucket"
	"github.com/awslabs/aws-service-operator/pkg/operators/snssubscription"
	"github.com/awslabs/aws-service-operator/pkg/operators/snstopic"
	"github.com/awslabs/aws-service-operator/pkg/operators/sqsqueue"
	"github.com/awslabs/aws-service-operator/pkg/queuemanager"
)

type base struct {
	config                 config.Config
	queueManager           *queuemanager.QueueManager
	cloudformationtemplate *cloudformationtemplate.Operator
	dynamodb               *dynamodb.Operator
	ecrrepository          *ecrrepository.Operator
	elasticache            *elasticache.Operator
	s3bucket               *s3bucket.Operator
	snssubscription        *snssubscription.Operator
	snstopic               *snstopic.Operator
	sqsqueue               *sqsqueue.Operator
}

func New(
	config config.Config,
	queueManager *queuemanager.QueueManager,
) *base {
	return &base{
		config:                 config,
		queueManager:           queueManager,
		cloudformationtemplate: cloudformationtemplate.NewOperator(config, queueManager),
		dynamodb:               dynamodb.NewOperator(config, queueManager),
		ecrrepository:          ecrrepository.NewOperator(config, queueManager),
		elasticache:            elasticache.NewOperator(config, queueManager),
		s3bucket:               s3bucket.NewOperator(config, queueManager),
		snssubscription:        snssubscription.NewOperator(config, queueManager),
		snstopic:               snstopic.NewOperator(config, queueManager),
		sqsqueue:               sqsqueue.NewOperator(config, queueManager),
	}
}

func (b *base) Watch(ctx context.Context, namespace string) {
	if b.config.Resources["cloudformationtemplate"] {
		go b.cloudformationtemplate.StartWatch(ctx, namespace)
	}
	if b.config.Resources["dynamodb"] {
		go b.dynamodb.StartWatch(ctx, namespace)
	}
	if b.config.Resources["ecrrepository"] {
		go b.ecrrepository.StartWatch(ctx, namespace)
	}
	if b.config.Resources["elasticache"] {
		go b.elasticache.StartWatch(ctx, namespace)
	}
	if b.config.Resources["s3bucket"] {
		go b.s3bucket.StartWatch(ctx, namespace)
	}
	if b.config.Resources["snssubscription"] {
		go b.snssubscription.StartWatch(ctx, namespace)
	}
	if b.config.Resources["snstopic"] {
		go b.snstopic.StartWatch(ctx, namespace)
	}
	if b.config.Resources["sqsqueue"] {
		go b.sqsqueue.StartWatch(ctx, namespace)
	}

	<-ctx.Done()
}
