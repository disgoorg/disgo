package httpserver

import (
	"net/http"

	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		URL:        "/interactions/callback",
		Address:    ":80",
		HTTPServer: &http.Server{},
		ServeMux:   http.NewServeMux(),
	}
}

type Config struct {
	Logger     log.Logger
	HTTPServer *http.Server
	ServeMux   *http.ServeMux
	URL        string
	Address    string
	PublicKey  string
	CertFile   string
	KeyFile    string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithHTTPServer(httpServer *http.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

func WithServeMux(serveMux *http.ServeMux) ConfigOpt {
	return func(config *Config) {
		config.ServeMux = serveMux
	}
}

func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

func WithAddress(address string) ConfigOpt {
	return func(config *Config) {
		config.Address = address
	}
}

func WithPublicKey(publicKey string) ConfigOpt {
	return func(config *Config) {
		config.PublicKey = publicKey
	}
}

func WithTLS(certFile string, keyFile string) ConfigOpt {
	return func(config *Config) {
		config.CertFile = certFile
		config.KeyFile = keyFile
	}
}
