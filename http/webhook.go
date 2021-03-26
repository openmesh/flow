package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/openmesh/flow"
	"net/http"
)

func makeWebhookHandlers(evb flow.EventBus, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	ingestWebhookHandler := kithttp.NewServer(
		makeIngestWebhookEndpoint(evb),
		decodeIngestWebhookRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/webhooks/{topic}", ingestWebhookHandler).Methods("POST")

	return r
}

func makeIngestWebhookEndpoint(evb flow.EventBus) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(flow.Event)
		err := evb.Publish(req.Topic, req.Payload)
		if err != nil {
			return map[string]string{"status": "failed"}, err
		}
		return map[string]string{"status": "success"}, nil
	}
}

func decodeIngestWebhookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	topic, ok := vars["topic"]
	if !ok {
		return nil, flow.Errorf(flow.EINVALID, "bad route")
	}

	var payload interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return flow.Event{
		Payload: payload,
		Topic:   topic,
	}, nil
}
