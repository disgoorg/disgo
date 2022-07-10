package rpc

import (
	"github.com/disgoorg/log"
)

// DefaultConfig is the configuration which is used by default
func DefaultConfig() *Config {
	return &Config{
		Logger:          log.Default(),
		TransportCreate: NewIPCTransport,
	}
}

type Config struct {
	Logger          log.Logger
	Transport       Transport
	TransportCreate TransportCreate
	Origin          string
}

// ConfigOpt can be used to supply optional parameters to NewIPCClient or NewWSClient
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger applies a custom logger to the rpc Client
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithTransport applies your own Transport to the rpc Client
func WithTransport(transport Transport) ConfigOpt {
	return func(config *Config) {
		config.Transport = transport
	}
}

// WithTransportCreate applies a custom logger to the rpc Client
func WithTransportCreate(transportCreate TransportCreate) ConfigOpt {
	return func(config *Config) {
		config.TransportCreate = transportCreate
	}
}

// WithIPCTransport applies the ipc Transport to the rpc Client
func WithIPCTransport() ConfigOpt {
	return func(config *Config) {
		config.TransportCreate = NewIPCTransport
	}
}

// WithWSTransport applies the ws Transport to the rpc Client
func WithWSTransport(origin string) ConfigOpt {
	return func(config *Config) {
		config.TransportCreate = NewWSTransport
		config.Origin = origin
	}
}

// WithOrigin sets the origin for the ws Transport
func WithOrigin(origin string) ConfigOpt {
	return func(config *Config) {
		config.Origin = origin
	}
}
