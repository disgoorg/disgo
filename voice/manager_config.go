package voice

import "github.com/disgoorg/log"

// DefaultManagerConfig returns the default ManagerConfig with sensible defaults.
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Logger:         log.Default(),
		ConnCreateFunc: NewConn,
	}
}

// ManagerConfig is a function that configures a Manager.
type ManagerConfig struct {
	Logger log.Logger

	ConnCreateFunc ConnCreateFunc
	ConnOpts       []ConnConfigOpt
}

// ManagerConfigOpt is used to functionally configure a ManagerConfig.
type ManagerConfigOpt func(ManagerConfig *ManagerConfig)

// Apply applies the given ManagerConfigOpts to the ManagerConfig.
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

// WithConnCreateFunc sets the ConnCreateFunc for the Manager
func WithConnCreateFunc(connectionCreateFunc ConnCreateFunc) ManagerConfigOpt {
	return func(ManagerConfig *ManagerConfig) {
		ManagerConfig.ConnCreateFunc = connectionCreateFunc
	}
}

// WithConnConfigOpts sets the ConnConfigOpt(s) for the Manager
func WithConnConfigOpts(opts ...ConnConfigOpt) ManagerConfigOpt {
	return func(ManagerConfig *ManagerConfig) {
		ManagerConfig.ConnOpts = append(ManagerConfig.ConnOpts, opts...)
	}
}
