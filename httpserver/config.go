package httpserver

import (
	"log/slog"
	"net/http"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:     slog.Default(),
		HTTPServer: &http.Server{},
		ServeMux:   http.NewServeMux(),
		URL:        "/interactions/callback",
		Address:    ":80",
	}
}

// Config lets you configure your Server instance.
type Config struct {
	Logger     *slog.Logger
	HTTPServer *http.Server
	ServeMux   *http.ServeMux
	URL        string
	Address    string
	CertFile   string
	KeyFile    string
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the Logger of the Config.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithHTTPServer sets the http.Server of the Config.
func WithHTTPServer(httpServer *http.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

// WithServeMux sets the http.ServeMux of the Config.
func WithServeMux(serveMux *http.ServeMux) ConfigOpt {
	return func(config *Config) {
		config.ServeMux = serveMux
	}
}

// WithURL sets the URL of the Config.
func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

// WithAddress sets the Address of the Config.
func WithAddress(address string) ConfigOpt {
	return func(config *Config) {
		config.Address = address
	}
}

// WithTLS sets the CertFile of the Config.
func WithTLS(certFile string, keyFile string) ConfigOpt {
	return func(config *Config) {
		config.CertFile = certFile
		config.KeyFile = keyFile
	}
}
