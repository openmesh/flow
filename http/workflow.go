package http

import (
	"context"
	"encoding/json"
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

	updateWorkflowHandler := kithttp.NewServer(
		makeUpdateWorkflowEndpoint(s),
		decodeUpdateWorkflowRequest,
		encodeResponse,
		opts...,
	)

	deleteWorkflowHandler := kithttp.NewServer(
		makeDeleteWorkflowEndpoint(s),
		decodeDeleteWorkflowRequest,
		encodeEmptyResponse,
		opts...,
	)

	getWorkflowByIDHandler := kithttp.NewServer(
		makeGetWorkflowByIDEndpoint(s),
		decodeGetWorkflowByIDRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/workflows/", createWorkflowHandler).Methods("POST")
	r.Handle("/v1/workflows/{id}", updateWorkflowHandler).Methods("PUT")
	r.Handle("/v1/workflows/{id}", deleteWorkflowHandler).Methods("DELETE")
	r.Handle("/v1/workflows/{id}", getWorkflowByIDHandler).Methods("GET")

	return r
}

/////////////////////
// Create workflow //
/////////////////////

// makeCreateWorkflowEndpoint returns an endpoint that calls CreateWorkflow on a flow.WorkflowService.
func makeCreateWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.CreateWorkflowRequest)
		return s.CreateWorkflow(ctx, req)
	}
}

// decodeCreateWorkflowRequest takes a http.Request and converts it into a flow.CreateWorkflowRequest. It returns an
// error if the JSON body cannot be encoded.
func decodeCreateWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req flow.CreateWorkflowRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, flow.Errorf(flow.EINVALID, "Failed to encode JSON body.")
	}

	return req, nil
}

/////////////////////
// Update workflow //
/////////////////////

// makeUpdateWorkflowEndpoint returns an endpoint that calls UpdateWorkflow on a flow.WorkflowService.
func makeUpdateWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.UpdateWorkflowRequest)
		return s.UpdateWorkflow(ctx, req)
	}
}

// decodeUpdateWorkflowRequest takes a http.Request and converts it into a flow.UpdateWorkflowRequest. It returns an
// error if the JSON body cannot be encoded or the ID cannot be parsed.
func decodeUpdateWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req flow.UpdateWorkflowRequest
	var err error

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, flow.Errorf(flow.EINVALID, "Failed to encode JSON body.")
	}

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

/////////////////////
// Delete workflow //
/////////////////////

// makeDeleteWorkflowEndpoint returns an endpoint that calls DeleteWorkflow on a flow.WorkflowService.
func makeDeleteWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.DeleteWorkflowRequest)
		return s.DeleteWorkflow(ctx, req), nil
	}
}

// decodeDeleteWorkflowRequest takes a http.Request and converts it into a flow.DeleteWorkflowRequest. It returns an
// error if the ID cannot be parsed.
func decodeDeleteWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req flow.DeleteWorkflowRequest
	var err error

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

////////////////////////
// Get workflow by ID //
////////////////////////

// makeGetWorkflowByIDEndpoint returns an endpoint that calls GetWorkflowByID on a flow.WorkflowService.
func makeGetWorkflowByIDEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.GetWorkflowByIDRequest)
		return s.GetWorkflowByID(ctx, req)
	}
}

// decodeGetWorkflowByIDRequest takes a http.Request and converts it into a flow.GetWorkflowByIDRequest. It returns an
// error if the ID cannot be parsed.
func decodeGetWorkflowByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req flow.GetWorkflowByIDRequest
	var err error

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}
