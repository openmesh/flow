package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/openmesh/flow"
	"net/http"
)

func makeAuthHandler(s flow.AuthService, sc *securecookie.SecureCookie, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	signUpHandler := kithttp.NewServer(
		makeSignUpEndpoint(s),
		decodeSignUpRequest,
		makeEncodeSignUpResponse(sc),
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/auth/signup/", signUpHandler).Methods("POST")

	return r
}

/////////////
// Sign up //
/////////////

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func makeSignUpEndpoint(s flow.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signUpRequest)
		user, err := s.SignUp(ctx, req.Email, req.Name, req.Password)
		return user, err
	}
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req signUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, flow.Errorf(flow.EINVALID, "Failed to encode JSON body.")
	}

	return req, nil
}

func makeEncodeSignUpResponse(sc *securecookie.SecureCookie) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(errorer); ok && e.error() != nil {
			// Not a Go kit transport error, but a business-logic error.
			// Provide those as HTTP errors.
			encodeError(ctx, e.error(), w)
			return nil
		}

		// Cast response into *flow.User
		user := response.(*flow.User)

		// Assign cookie values
		cookieVal := Session{
			UserID:      user.ID,
			RedirectURL: "",
			State:       "",
		}
		// Encode cookie values
		encodedValues, err := sc.Encode(SessionCookieName, cookieVal)
		if err != nil {
			val := err.Error()
			fmt.Println(val)
			encodeError(ctx, flow.Errorf(flow.EINTERNAL, "failed to encode cookie: %w", err), w)
			return nil
		}

		// Create and set session cookie
		cookie := &http.Cookie{
			Name:  SessionCookieName,
			Value: encodedValues,
			Path:  "/",
			// TODO make this secure only
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		// Encode user into response body
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json.NewEncoder(w).Encode(response)
	}
}
