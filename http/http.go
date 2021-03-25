package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/openmesh/flow"
)

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

// encodeError prints & optionally logs an error message.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// Extract error code & message.
	code, message := flow.ErrorCode(err), flow.ErrorMessage(err)

	// Track metrics by code.
	// errorCount.WithLabelValues(code).Inc()

	// Log & report internal errors.
	//if code == template.EINTERNAL {
	//	template.ReportError(r.Context(), err, r)
	//	LogError(r, err)
	//}

	// Print user message to response.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	_ = json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	flow.ECONFLICT:       http.StatusConflict,
	flow.EINVALID:        http.StatusBadRequest,
	flow.ENOTFOUND:       http.StatusNotFound,
	flow.ENOTIMPLEMENTED: http.StatusNotImplemented,
	flow.EUNAUTHORIZED:   http.StatusUnauthorized,
	flow.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for an error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
