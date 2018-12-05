package queuemanager

import (
	"sync"

	"github.com/awslabs/aws-service-operator/pkg/config"
)

// HandlerFunc allows you to define a custom function for when a message is stored
type HandlerFunc func(config config.Config, msg *MessageBody) error

// Handler allows a custom function to be passed
type Handler interface {
	HandleMessage(config config.Config, msg *MessageBody) error
}

// Queue Manager allows you to register topics and a handler function
type QueueManager struct {
	lock     sync.RWMutex
	handlers map[string]Handler
}

// MessageBody will parse the message from the Body of SQS
type MessageBody struct {
	Type               string `json:"Type"`
	TopicARN           string `json:"TopicArn"`
	Message            string `json:"Message"`
	ParsedMessage      map[string]string
	Namespace          string
	ResourceName       string
	ResourceProperties ResourceProperties
	Updatable          bool
}

// ResourceProperties will wrap the ResourceProperties object
type ResourceProperties struct {
	Tags []Tag `json:"Tags"`
}

// Tag represents a Tag
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
