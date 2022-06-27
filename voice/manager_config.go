package voice

import "github.com/disgoorg/log"

func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Logger:         log.Default(),
		ConnCreateFunc: NewConn,
	}
}

type ManagerConfig struct {
	Logger log.Logger

	ConnCreateFunc ConnCreateFunc
	ConnOpts       []ConnConfigOpt
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

func WithConnCreateFunc(connectionCreateFunc ConnCreateFunc) ManagerConfigOpt {
	return func(ManagerConfig *ManagerConfig) {
		ManagerConfig.ConnCreateFunc = connectionCreateFunc
	}
}

func WithConnConfigOpts(opts ...ConnConfigOpt) ManagerConfigOpt {
	return func(ManagerConfig *ManagerConfig) {
		ManagerConfig.ConnOpts = append(ManagerConfig.ConnOpts, opts...)
	}
}
