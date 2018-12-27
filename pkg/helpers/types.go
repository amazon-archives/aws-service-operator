package helpers

import "github.com/awslabs/aws-service-operator/pkg/config"

// Data wrapps the object that is needed for the services
type Data struct {
	Helpers Helpers
	Obj     interface{}
	Config  config.Config
}
