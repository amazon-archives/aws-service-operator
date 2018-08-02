package server

import (
	"github.com/christopherhein/aws-operator/pkg/config"
)

// Server defines the bas construct for the operator
type Server struct {
	Config *config.Config
}
