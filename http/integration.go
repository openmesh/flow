package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openmesh/flow"
	"net/http"
)

func (s *Server) makeIntegrationHandler() http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(s.Logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getIntegrationsHandler := kithttp.NewServer(
		makeGetIntegrationsEndpoint(s.IntegrationService),
		decodeGetIntegrationsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/integrations", getIntegrationsHandler).Methods("GET")

	return r
}

//////////////////////
// Get integrations //
//////////////////////

type getIntegrationsRequest struct{}

func decodeGetIntegrationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getIntegrationsRequest{}, nil
}

func makeGetIntegrationsEndpoint(s flow.IntegrationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		apps, total, err := s.GetIntegrations(ctx, flow.GetIntegrationsRequest{})
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"data":        apps,
			"total_items": total,
		}, nil
	}
}
