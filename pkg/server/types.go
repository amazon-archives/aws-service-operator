package server

import (
	"github.com/awslabs/aws-service-operator/pkg/config"
)

// Server defines the bas construct for the operator
type Server struct {
	Config *config.Config
}
