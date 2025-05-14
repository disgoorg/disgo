package httpserver

import (
	"log/slog"
	"net/http"
)

func defaultConfig() config {
	return config{
		Logger:     slog.Default(),
		HTTPServer: &http.Server{},
		ServeMux:   http.NewServeMux(),
		URL:        "/interactions/callback",
		Address:    ":80",
		Verifier:   DefaultVerifier{},
	}
}

type config struct {
	Logger     *slog.Logger
	HTTPServer *http.Server
	ServeMux   *http.ServeMux
	URL        string
	Address    string
	CertFile   string
	KeyFile    string
	Verifier   Verifier
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Server.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "httpserver"))
}

// WithLogger sets the Logger of the config.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithHTTPServer sets the http.Server of the config.
func WithHTTPServer(httpServer *http.Server) ConfigOpt {
	return func(config *config) {
		config.HTTPServer = httpServer
	}
}

// WithServeMux sets the http.ServeMux of the config.
func WithServeMux(serveMux *http.ServeMux) ConfigOpt {
	return func(config *config) {
		config.ServeMux = serveMux
	}
}

// WithURL sets the URL of the config.
func WithURL(url string) ConfigOpt {
	return func(config *config) {
		config.URL = url
	}
}

// WithAddress sets the Address of the config.
func WithAddress(address string) ConfigOpt {
	return func(config *config) {
		config.Address = address
	}
}

// WithTLS sets the CertFile of the config.
func WithTLS(certFile string, keyFile string) ConfigOpt {
	return func(config *config) {
		config.CertFile = certFile
		config.KeyFile = keyFile
	}
}

// WithVerifier sets the Verifier of the config.
func WithVerifier(verifier Verifier) ConfigOpt {
	return func(config *config) {
		config.Verifier = verifier
	}
}
