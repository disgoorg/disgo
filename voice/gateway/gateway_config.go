package gateway

import (
	"time"

	"github.com/disgoorg/log"
	"github.com/gorilla/websocket"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		Dialer:            websocket.DefaultDialer,
		AutoReconnect:     true,
		ReconnectTimeout:  5 * time.Second,
		MaxReconnectTries: 10,
	}
}

type Config struct {
	Logger            log.Logger
	Dialer            *websocket.Dialer
	AutoReconnect     bool
	ReconnectTimeout  time.Duration
	MaxReconnectTries int
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

func WithDialer(dialer *websocket.Dialer) ConfigOpt {
	return func(config *Config) {
		config.Dialer = dialer
	}
}

func WithAutoReconnect(autoReconnect bool) ConfigOpt {
	return func(config *Config) {
		config.AutoReconnect = autoReconnect
	}
}

func WithReconnectTimeout(reconnectTimeout time.Duration) ConfigOpt {
	return func(config *Config) {
		config.ReconnectTimeout = reconnectTimeout
	}
}

func WithMaxReconnectTries(maxReconnectTries int) ConfigOpt {
	return func(config *Config) {
		config.MaxReconnectTries = maxReconnectTries
	}
}
