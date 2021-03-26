package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/openmesh/flow"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Addr   string
	Domain string

	Logger log.Logger

	EventBus        flow.EventBus
	WebhookService  flow.WebhookService
	WorkflowService flow.WorkflowService
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
		server: &http.Server{},
	}

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overridding methods.
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	return s
}

// UseTLS returns true if the cert & key file are specified.
func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Healthy"))
}

func (s *Server) Open() (err error) {
	// Assign all the base handlers
	s.configureHandlers()
	s.router.HandleFunc("/health", healthCheck)

	// Open a listener on our bind address.
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	err = s.server.Serve(s.ln)
	if err != nil {
		return err
	}

	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// RegisterRoute allows additional routes to be registered to the router. This allows instrumenting middleware to be
// implemented without the Server knowing about the implementation.
func (s *Server) RegisterRoute(path string, handler http.Handler) {
	s.router.Handle(path, handler)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Override method for forms passing "_method" value.
	if r.Method == http.MethodPost {
		switch v := r.PostFormValue("_method"); v {
		case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
			r.Method = v
		}
	}

	// Allow CORS
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// Delegate remaining HTTP handling to the gorilla router.
	handlers.CORS(
		allowedOrigins,
		allowedHeaders,
		allowedMethods,
		handlers.AllowCredentials(),
	)(s.router).ServeHTTP(w, r)
}

func (s *Server) configureHandlers() {
	s.router.Handle("/v1/workflows/", makeWorkflowHandler(s.WorkflowService, s.Logger))
	s.router.Handle("/v1/webhooks/", makeWebhookHandlers(s.EventBus, s.Logger))
}
