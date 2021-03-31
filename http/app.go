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

func (s *Server) makeAppHandler() http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(s.Logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getAppsHandler := kithttp.NewServer(
		makeGetAppsEndpoint(s.AppService),
		decodeGetAppsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/apps/", getAppsHandler).Methods("GET")

	return r
}

//////////////
// Get apps //
//////////////

type getAppsRequest struct{}

func decodeGetAppsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getAppsRequest{}, nil
}

func makeGetAppsEndpoint(s flow.AppService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		apps, total, err := s.GetApps(ctx, flow.GetAppsRequest{})
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"data":        apps,
			"total_items": total,
		}, nil
	}
}
