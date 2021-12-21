package httpserver

import (
	"net/http"

	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	URL:  "/interactions/callback",
	Port: ":80",
}

type Config struct {
	Logger           log.Logger
	HTTPServer       *http.Server
	ServeMux         *http.ServeMux
	EventHandlerFunc EventHandlerFunc
	URL              string
	Port             string
	PublicKey        string
	CertFile         string
	KeyFile          string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServer(httpServer *http.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithServeMux(serveMux *http.ServeMux) ConfigOpt {
	return func(config *Config) {
		config.ServeMux = serveMux
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithPort(port string) ConfigOpt {
	return func(config *Config) {
		config.Port = port
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithPublicKey(publicKey string) ConfigOpt {
	return func(config *Config) {
		config.PublicKey = publicKey
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithTLS(certFile string, keyFile string) ConfigOpt {
	return func(config *Config) {
		config.CertFile = certFile
		config.KeyFile = keyFile
	}
}
