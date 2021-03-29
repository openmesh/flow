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
	"time"
)

func (s *Server) makeNodeHandler() http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(s.Logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createNodeHandler := kithttp.NewServer(
		makeCreateNodeEndpoint(s.NodeService),
		decodeCreateNodeRequest,
		encodeResponse,
		opts...,
	)

	updateNodeHandler := kithttp.NewServer(
		makeUpdateNodeEndpoint(s.NodeService),
		decodeUpdateNodeRequest,
		encodeResponse,
		opts...,
	)

	deleteNodeHandler := kithttp.NewServer(
		makeDeleteNodeEndpoint(s.NodeService),
		decodeDeleteNodeRequest,
		encodeResponse,
		opts...,
	)

	getNodeByIDHandler := kithttp.NewServer(
		makeGetNodeByIDEndpoint(s.NodeService),
		decodeDeleteNodeRequest,
		encodeResponse,
		opts...,
	)

	getNodesHandler := kithttp.NewServer(
		makeGetNodesEndpoint(s.NodeService),
		decodeGetNodesRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/nodes/", s.authenticate(createNodeHandler)).Methods("POST")
	r.Handle("/v1/nodes/{id}", s.authenticate(updateNodeHandler)).Methods("PUT")
	r.Handle("/v1/nodes/{id}", s.authenticate(deleteNodeHandler)).Methods("DELETE")
	r.Handle("/v1/nodes/{id}", s.authenticate(getNodeByIDHandler)).Methods("GET")
	r.Handle("/v1/nodes/", s.authenticate(getNodesHandler)).Methods("GET")

	return r
}

/////////////////
// Create node //
/////////////////

type createNodeRequest struct {
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	WorkflowID  uuid.UUID    `json:"workflow_id" db:"workflow_id"`
	Integration string       `json:"integration" db:"integration"`
	Action      string       `json:"action" db:"action"`
	ParentIDs   []*uuid.UUID `json:"parent_ids" db:"-"`
	ChildrenIDs []*uuid.UUID `json:"children_ids" db:"-"`
}

type createParamRequest struct {
	Key   string         `json:"key" db:"key"`
	Value string         `json:"value" db:"value"`
	Type  flow.ParamType `json:"type" db:"type"`
}

func makeCreateNodeEndpoint(s flow.NodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createNodeRequest)
		node := flow.Node{
			Integration: req.Integration,
			Action:      req.Action,
			ParentIDs:   req.ParentIDs,
			ChildrenIDs: req.ChildrenIDs,
		}

		//for _, param := range req.Params {
		//	node.Params = append(node.Params, &flow.Param{
		//		Key:   param.Key,
		//		Value: param.Value,
		//		Type:  param.Type,
		//	})
		//}

		err := s.CreateNode(ctx, &node)
		return node, err
	}
}

func decodeCreateNodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createNodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, flow.Errorf(flow.EINVALID, "Failed to encode JSON body.")
	}

	return req, nil
}

/////////////////////
// Update workflow //
/////////////////////

type updateNodeRequest struct {
	ID          uuid.UUID    `json:"id"`
	WorkflowID  uuid.UUID    `json:"workflow_id"`
	Integration string       `json:"integration"`
	Action      string       `json:"action"`
	ParentIDs   []*uuid.UUID `json:"parent_ids" db:"-"`
	ChildrenIDs []*uuid.UUID `json:"children_ids" db:"-"`
}

// makeUpdateNodeEndpoint returns an endpoint that calls UpdateNode on a flow.NodeService.
func makeUpdateNodeEndpoint(s flow.NodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateNodeRequest)
		upd := flow.NodeUpdate{
			Integration: req.Integration,
			Action:      req.Action,
			ParentIDs:   req.ParentIDs,
			ChildrenIDs: req.ChildrenIDs,
		}
		return s.UpdateNode(ctx, req.ID, upd)
	}
}

// decodeUpdateNodeRequest takes a http.Request and converts it into a flow.UpdateNodeRequest. It
// returns an error if the JSON body cannot be encoded or the ID cannot be parsed.
func decodeUpdateNodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateNodeRequest
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

/////////////////
// Delete node //
/////////////////

type deleteNodeRequest struct {
	ID uuid.UUID
}

// makeDeleteNodeEndpoint returns an endpoint that calls DeleteNode on a flow.NodeService.
func makeDeleteNodeEndpoint(s flow.NodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteNodeRequest)
		return s.DeleteNode(ctx, req.ID), nil
	}
}

// decodeDeleteNodeRequest takes a http.Request and converts it into a flow.DeleteNodeRequest. It
// returns an error if the ID cannot be parsed.
func decodeDeleteNodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req deleteNodeRequest
	var err error

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

////////////////////
// Get node by ID //
////////////////////

type getNodeByIDRequest struct {
	ID uuid.UUID
}

// makeGetNodeByIDEndpoint returns an endpoint that calls GetNodeByID on a flow.NodeService.
func makeGetNodeByIDEndpoint(s flow.NodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getNodeByIDRequest)
		return s.GetNodeByID(ctx, req.ID)
	}
}

// decodeGetNodeByIDRequest takes a http.Request and converts it into a flow.GetNodeByIDRequest. It
// returns an error if the ID cannot be parsed.
func decodeGetNodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req getNodeByIDRequest
	var err error

	req.ID, err = uuidFromVar(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

///////////////
// Get nodes //
///////////////

type getNodesRequest struct {
	WorkflowID *uuid.UUID `json:"id"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}

type getNodesResponse struct {
	Data       []*flow.Node `json:"data"`
	TotalItems int          `json:"total_items"`
}

// makeGetNodesEndpoint returns an endpoint that calls GetNodes on a flow.NodeService.
func makeGetNodesEndpoint(s flow.NodeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getNodesRequest)
		filter := flow.NodeFilter{
			Page:       req.Page,
			Limit:      req.Limit,
			WorkflowID: req.WorkflowID,
		}
		nodes, total, err := s.GetNodes(ctx, filter)
		if err != nil {
			return nil, err
		}

		return getNodesResponse{
			Data:       nodes,
			TotalItems: total,
		}, nil
	}
}

func decodeGetNodesRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
	if val, ok := vars["workflow_id"]; ok {
		if id, err := uuid.Parse(val); err != nil {
			return nil, flow.Errorf(flow.EINVALID, "Invalid value provided for parameter 'workflow_id'.")
		} else {
			req.ID = &id
		}
	}

	return req, nil
}
