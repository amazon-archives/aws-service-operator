/*
Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
not use this file except in compliance with the License. A copy of the
License is located at

     http://aws.amazon.com/apache2.0/

or in the "license" file accompanying this file. This file is distributed
on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing
permissions and limitations under the License.
*/

package server

import (
	"fmt"
	"net/http"
	"time"

	self "awsoperator.io/pkg/apis/self/v1alpha1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler wraps handler functions
type Handler struct {
	http.ServeMux
}

// New returns an http.Server for exposing endpoints
func New(config self.Config) *http.Server {
	return &http.Server{
		Handler:      newHandler(config),
		Addr:         config.Server.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func newHandler(config self.Config) *Handler {
	h := &Handler{}
	h.HandleFunc("/healthz", healthzFunc)
	if config.Server.Metrics.Enable {
		h.Handle(config.Server.Metrics.Endpoint, promhttp.Handler())
	}
	return h
}

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
