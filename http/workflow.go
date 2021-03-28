package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/openmesh/flow"
	"net/http"
	"strconv"
)

func (s *Server) makeWorkflowHandler() http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(s.Logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createWorkflowHandler := kithttp.NewServer(
		makeCreateWorkflowEndpoint(s.WorkflowService),
		decodeCreateWorkflowRequest,
		encodeResponse,
		opts...,
	)

	updateWorkflowHandler := kithttp.NewServer(
		makeUpdateWorkflowEndpoint(s.WorkflowService),
		decodeUpdateWorkflowRequest,
		encodeResponse,
		opts...,
	)

	deleteWorkflowHandler := kithttp.NewServer(
		makeDeleteWorkflowEndpoint(s.WorkflowService),
		decodeDeleteWorkflowRequest,
		encodeEmptyResponse,
		opts...,
	)

	getWorkflowByIDHandler := kithttp.NewServer(
		makeGetWorkflowByIDEndpoint(s.WorkflowService),
		decodeGetWorkflowByIDRequest,
		encodeResponse,
		opts...,
	)

	getWorkflowsHandler := kithttp.NewServer(
		makeGetWorkflowsEndpoint(s.WorkflowService),
		decodeGetWorkflowsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/workflows/", s.authenticate(createWorkflowHandler)).Methods("POST")
	r.Handle("/v1/workflows/{id}", s.authenticate(updateWorkflowHandler)).Methods("PUT")
	r.Handle("/v1/workflows/{id}", s.authenticate(deleteWorkflowHandler)).Methods("DELETE")
	r.Handle("/v1/workflows/{id}", s.authenticate(getWorkflowByIDHandler)).Methods("GET")
	r.Handle("/v1/workflows/", s.authenticate(getWorkflowsHandler)).Methods("GET")

	return r
}

/////////////////////
// Create workflow //
/////////////////////

type createWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// makeCreateWorkflowEndpoint returns an endpoint that calls CreateWorkflow on a flow.WorkflowService.
func makeCreateWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createWorkflowRequest)
		workflow := flow.Workflow{
			Name:        req.Name,
			Description: req.Description,
		}
		err := s.CreateWorkflow(ctx, &workflow)
		return workflow, err
	}
}

// decodeCreateWorkflowRequest takes a http.Request and converts it into a flow.CreateWorkflowRequest. It returns an
// error if the JSON body cannot be encoded.
func decodeCreateWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createWorkflowRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, flow.Errorf(flow.EINVALID, "Failed to encode JSON body.")
	}

	return req, nil
}

/////////////////////
// Update workflow //
/////////////////////

type updateWorkflowRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
}

// makeUpdateWorkflowEndpoint returns an endpoint that calls UpdateWorkflow on a flow.WorkflowService.
func makeUpdateWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateWorkflowRequest)
		upd := flow.WorkflowUpdate{
			Name:        req.Name,
			Description: req.Description,
		}
		return s.UpdateWorkflow(ctx, req.ID, upd)
	}
}

// decodeUpdateWorkflowRequest takes a http.Request and converts it into a flow.UpdateWorkflowRequest. It returns an
// error if the JSON body cannot be encoded or the ID cannot be parsed.
func decodeUpdateWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateWorkflowRequest
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

type deleteWorkflowRequest struct {
	ID uuid.UUID
}

// makeDeleteWorkflowEndpoint returns an endpoint that calls DeleteWorkflow on a flow.WorkflowService.
func makeDeleteWorkflowEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteWorkflowRequest)
		return s.DeleteWorkflow(ctx, req.ID), nil
	}
}

// decodeDeleteWorkflowRequest takes a http.Request and converts it into a flow.DeleteWorkflowRequest. It returns an
// error if the ID cannot be parsed.
func decodeDeleteWorkflowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req deleteWorkflowRequest
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

type getWorkflowByIDRequest struct {
	ID uuid.UUID
}

// makeGetWorkflowByIDEndpoint returns an endpoint that calls GetWorkflowByID on a flow.WorkflowService.
func makeGetWorkflowByIDEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getWorkflowByIDRequest)
		return s.GetWorkflowByID(ctx, req.ID)
	}
}

// decodeGetWorkflowByIDRequest takes a http.Request and converts it into a flow.GetWorkflowByIDRequest. It returns an
// error if the ID cannot be parsed.
func decodeGetWorkflowByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getWorkflowByIDRequest
	var err error

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

///////////////////
// Get workflows //
///////////////////

type getWorkflowsRequest struct {
	ID          *uuid.UUID `json:"id"`
	Page        int        `json:"page"`
	Limit       int        `json:"limit"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
}

type getWorkflowsResponse struct {
	Data       []*flow.Workflow `json:"data"`
	TotalItems int              `json:"total_items"`
}

// makeGetWorkflowsEndpoint returns an endpoint that calls GetWorkflows on a flow.WorkflowService.
func makeGetWorkflowsEndpoint(s flow.WorkflowService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getWorkflowsRequest)
		filter := flow.WorkflowFilter{
			ID:          req.ID,
			Page:        req.Page,
			Limit:       req.Limit,
			Name:        req.Name,
			Description: req.Description,
		}
		workflows, total, err := s.GetWorkflows(ctx, filter)
		if err != nil {
			return nil, err
		}

		return getWorkflowsResponse{
			Data:       workflows,
			TotalItems: total,
		}, nil
	}
}

func decodeGetWorkflowsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getWorkflowsRequest{}

	vars := mux.Vars(r)
	if val, ok := vars["page"]; ok {
		if parsed, err := strconv.Atoi(val); err != nil {
			return nil, flow.Errorf(flow.EINVALID, "Invalid value provided for parameter 'page'.")
		} else {
			req.Page = parsed
		}
	}
	if val, ok := vars["limit"]; ok {
		if parsed, err := strconv.Atoi(val); err != nil {
			return nil, flow.Errorf(flow.EINVALID, "Invalid value provided for parameter 'limit'.")
		} else {
			req.Page = parsed
		}
	}
	if val, ok := vars["id"]; ok {
		if id, err := uuid.Parse(val); err != nil {
			return nil, flow.Errorf(flow.EINVALID, "Invalid value provided for parameter 'id'.")
		} else {
			req.ID = &id
		}
	}
	if val, ok := vars["name"]; ok {
		req.Name = &val
	}
	if val, ok := vars["description"]; ok {
		req.Description = &val
	}

	return req, nil
}
