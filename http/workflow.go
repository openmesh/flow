package http

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openmesh/flow"
)

func makeWorkflowHandler(s flow.WorkflowService, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createWorkflowHandler := kithttp.NewServer(
		makeCreateWorkflowEndpoint(s),
		decodeCreateWorkflowRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/workflows", createWorkflowHandler).Methods("POST")

	return r
}

func makeCreateWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.CreateWorkflowRequest)
		return s.CreateWorkflow(ctx, req)
	}
}

func decodeCreateWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return flow.CreateWorkflowRequest{}, nil
}
