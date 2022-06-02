package voice

import "github.com/disgoorg/log"

func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Logger:               log.Default(),
		ConnectionCreateFunc: NewConnection,
	}
}

type ManagerConfig struct {
	Logger log.Logger

	ConnectionCreateFunc ConnectionCreateFunc
	ConnectionOpts       []ConnectionConfigOpt
}

type ManagerConfigOpt func(ManagerConfig *ManagerConfig)

func (c *ManagerConfig) Apply(opts []ManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the logger for the webhook client
func WithLogger(logger log.Logger) ManagerConfigOpt {
	return func(ManagerConfig *ManagerConfig) {
		ManagerConfig.Logger = logger
	}
}
