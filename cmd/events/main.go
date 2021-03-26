package main

import (
	"context"
	"fmt"
	"github.com/openmesh/flow/eventbus"
	"os"
	"os/signal"

	"github.com/openmesh/flow/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()
	m.HTTPServer.Addr = ":8080"

	// Execute program.
	if err := m.Run(ctx); err != nil {
		_ = m.Close()
		_, _ = fmt.Fprintln(os.Stderr, err)
		// wtf.ReportError(ctx, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program.
type Main struct {
	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	HTTPServer *http.Server
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		HTTPServer: http.NewServer(),
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Run(ctx context.Context) (err error) {
	m.HTTPServer.Logger = createLogger()
	// requestCount, errorCount, requestDuration := setupMetrics()

	// Initialize services.

	// todoService := logmw.NewTodoLoggingMiddleware(m.HTTPServer.Logger)(inmem.NewService())
	// todoService = instrmw.NewTodoInstrumentingMiddleware(requestCount, errorCount, requestDuration)(todoService)

	// Attach underlying service to the HTTP server.
	m.HTTPServer.EventBus = eventbus.New()
	// m.HTTPServer.WorkflowService = pg.NewWorkflowService

	m.HTTPServer.RegisterRoute("/metrics", promhttp.Handler())

	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

func createLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}

func setupMetrics() (metrics.Counter, metrics.Counter, metrics.Histogram) {
	fieldKeys := []string{"method", "error"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "todo",
		Subsystem: "todo_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	errorCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "todo",
		Subsystem: "todo_service",
		Name:      "error_count",
		Help:      "Number of errors that have occurred.",
	}, fieldKeys)

	requestDuration := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "todo",
		Subsystem: "todo_service",
		Name:      "request_duration_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	return requestCount, errorCount, requestDuration
}
