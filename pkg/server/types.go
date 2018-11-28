package server

import (
	"net/http"

	"github.com/awslabs/aws-service-operator/pkg/config"
)

// Server defines the bas construct for the operator
type Server struct {
	http.ServeMux
	Config config.Config
}
