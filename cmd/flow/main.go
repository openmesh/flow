package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/openmesh/flow/eventbus"
	"github.com/openmesh/flow/pg"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

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

	// Parse command line flags & load configuration.
	if err := m.ParseFlags(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	// Configuration path and parsed config data.
	Config     Config
	ConfigPath string

	// Postgres database used by the pg service implementations.
	DB *pg.DB

	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	HTTPServer *http.Server
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,

		DB:         pg.NewDB(""),
		HTTPServer: http.NewServer(),
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	// Close server if it has a value.
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	// Close DB connection if it has a value.
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ParseFlags parses the command line arguments & loads the config.
//
// This exists separately from the Run() function so that we can skip it
// during end-to-end tests. Those tests will configure manually and call Run().
func (m *Main) ParseFlags(ctx context.Context, args []string) error {
	// Our flag set is very simple. It only includes a config path.
	fs := flag.NewFlagSet("flow", flag.ContinueOnError)
	fs.StringVar(&m.ConfigPath, "config", DefaultConfigPath, "config path")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// The expand() function is here to automatically expand "~" to the user's
	// home directory. This is a common task as configuration files are typing
	// under the home directory during local development.
	configPath, err := expand(m.ConfigPath)
	if err != nil {
		return err
	}

	// Read our TOML formatted configuration file.
	config, err := ReadConfigFile(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", m.ConfigPath)
	} else if err != nil {
		return err
	}
	m.Config = config

	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Run(ctx context.Context) (err error) {
	m.DB.DSN = m.Config.DB.DSN

	if err := m.DB.Connect(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	m.HTTPServer.Logger = createLogger()
	// requestCount, errorCount, requestDuration := setupMetrics()

	// Initialize services.
	eventBus := eventbus.New()
	workflowService := pg.NewWorkflowService(m.DB)
	authService := pg.NewAuthService(m.DB)

	// Attach underlying service to the HTTP server.
	m.HTTPServer.EventBus = eventBus
	m.HTTPServer.WorkflowService = workflowService
	m.HTTPServer.AuthService = authService

	m.HTTPServer.RegisterRoute("/metrics", promhttp.Handler())

	// Copy configuration settings to the HTTP server.
	m.HTTPServer.Addr = m.Config.HTTP.Addr
	m.HTTPServer.Domain = m.Config.HTTP.Domain
	m.HTTPServer.HashKey = m.Config.HTTP.HashKey
	m.HTTPServer.BlockKey = m.Config.HTTP.BlockKey

	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

const (
	// DefaultConfigPath is the default path to the application configuration.
	DefaultConfigPath = "config.toml"

	// DefaultDSN is the default datasource name.
	DefaultDSN = "user=postgres password=postgres dbname=okount port=5432 sslmode=false host=localhost"
)

// Config represents the CLI configuration file.
type Config struct {
	DB struct {
		DSN string `toml:"dsn"`
	} `toml:"db"`

	HTTP struct {
		Addr     string `toml:"addr"`
		Domain   string `toml:"domain"`
		HashKey  string `toml:"hash-key"`
		BlockKey string `toml:"block-key"`
	} `toml:"http"`
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	var config Config
	config.DB.DSN = DefaultDSN
	return config
}

// ReadConfigFile unmarshalls config from config file
func ReadConfigFile(filename string) (Config, error) {
	config := DefaultConfig()
	if buf, err := ioutil.ReadFile(filename); err != nil {
		return config, err
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}
	return config, nil
}

// expand returns path using tilde expansion. This means that a file path that
// begins with the "~" will be expanded to prefix the user's home directory.
func expand(path string) (string, error) {
	// Ignore if path has no leading tilde.
	if path != "~" && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	// Fetch the current user to determine the home path.
	u, err := user.Current()
	if err != nil {
		return path, err
	} else if u.HomeDir == "" {
		return path, fmt.Errorf("home directory unset")
	}

	if path == "~" {
		return u.HomeDir, nil
	}
	return filepath.Join(u.HomeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
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
