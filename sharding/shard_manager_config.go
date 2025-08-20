package sharding

import (
	"log/slog"

	"github.com/disgoorg/disgo/gateway"
)

func defaultConfig() config {
	return config{
		Logger:            slog.Default(),
		GatewayCreateFunc: gateway.New,
		ShardSplitCount:   DefaultShardSplitCount,
	}
}

// ShardState is used to tell a [gateway.Gateway] managed by the [ShardManager] which session & sequence it should use when starting the shard.
// This is useful for resuming shards when using the [ShardManager].
type ShardState struct {
	// SessionID is the session ID of the shard. This is used to resume the shard.
	SessionID string
	// Sequence is the sequence number of the shard. This is used to resume the shard.
	Sequence int
	// ResumeURL is the resume url to use for the shard. This is used to resume the shard.
	ResumeURL string
}

type config struct {
	// Logger is the logger of the ShardManager. Defaults to log.Default()
	Logger *slog.Logger
	// ShardIDs is a map of shardIDs the ShardManager should manage. Leave this nil to manage all shards.
	ShardIDs map[int]ShardState
	// ShardCount is the total shard count of the ShardManager. Leave this at 0 to let Discord calculate the shard count for you.
	ShardCount int
	// ShardSplitCount is the count a shard should be split into if it is too large. This is only used if AutoScaling is enabled.
	ShardSplitCount int
	// AutoScaling will automatically re-shard shards if they are too large. This is disabled by default.
	AutoScaling bool
	// GatewayCreateFunc is the function which is used by the ShardManager to create a new gateway.Gateway. Defaults to gateway.New.
	GatewayCreateFunc gateway.CreateFunc
	// GatewayConfigOpts are the ConfigOpt(s) which are applied to the gateway.Gateway.
	GatewayConfigOpts []gateway.ConfigOpt
	// IdentifyRateLimiter is the RateLimiter which is used by the ShardManager. Defaults to NewRateLimiter()
	IdentifyRateLimiter gateway.IdentifyRateLimiter
	// IdentifyRateLimiterConfigOpts are the gateway.IdentifyRateLimiterConfigOpt(s) which are applied to the gateway.IdentifyRateLimiter.
	IdentifyRateLimiterConfigOpts []gateway.IdentifyRateLimiterConfigOpt
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Server.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "sharding"))
	if c.IdentifyRateLimiter == nil {
		c.IdentifyRateLimiter = gateway.NewIdentifyRateLimiter(c.IdentifyRateLimiterConfigOpts...)
	}
}

// WithDefault returns a ConfigOpt that sets the default values for the ShardManager.
func WithDefault() ConfigOpt {
	return func(config *config) {}
}

// WithLogger sets the logger of the ShardManager.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithShardIDs sets the shardIDs the ShardManager should manage.
func WithShardIDs(shardIDs ...int) ConfigOpt {
	return func(config *config) {
		config.ShardIDs = map[int]ShardState{}
		for _, shardID := range shardIDs {
			config.ShardIDs[shardID] = ShardState{}
		}
	}
}

// WithShardIDsWithStates sets the shardIDs and their [ShardState] the ShardManager should manage.
func WithShardIDsWithStates(shards map[int]ShardState) ConfigOpt {
	return func(config *config) {
		config.ShardIDs = shards
	}
}

// WithShardCount sets the shard count of the ShardManager.
func WithShardCount(shardCount int) ConfigOpt {
	return func(config *config) {
		config.ShardCount = shardCount
	}
}

// WithShardSplitCount sets the count a shard should be split into if it is too large.
// This is only used if AutoScaling is enabled.
func WithShardSplitCount(shardSplitCount int) ConfigOpt {
	return func(config *config) {
		config.ShardSplitCount = shardSplitCount
	}
}

// WithAutoScaling sets whether the ShardManager should automatically re-shard shards if they are too large. This is disabled by default.
func WithAutoScaling(autoScaling bool) ConfigOpt {
	return func(config *config) {
		config.AutoScaling = autoScaling
	}
}

// WithGatewayCreateFunc sets the function which is used by the ShardManager to create a new gateway.Gateway.
func WithGatewayCreateFunc(gatewayCreateFunc gateway.CreateFunc) ConfigOpt {
	return func(config *config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

// WithGatewayConfigOpts lets you configure the gateway.Gateway created by the ShardManager.
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithIdentifyRateLimiter lets you inject your own RateLimiter into the ShardManager.
func WithIdentifyRateLimiter(rateLimiter gateway.IdentifyRateLimiter) ConfigOpt {
	return func(config *config) {
		config.IdentifyRateLimiter = rateLimiter
	}
}

// WithIdentifyRateLimiterConfigOpt lets you configure the default gateway.IdentifyRateLimiter used by the ShardManager.
func WithIdentifyRateLimiterConfigOpt(opts ...gateway.IdentifyRateLimiterConfigOpt) ConfigOpt {
	return func(config *config) {
		config.IdentifyRateLimiterConfigOpts = append(config.IdentifyRateLimiterConfigOpts, opts...)
	}
}

// WithDefaultIdentifyRateLimiterConfigOpt lets you configure the default gateway.IdentifyRateLimiter used by the ShardManager and prepend the options to the existing ones.
func WithDefaultIdentifyRateLimiterConfigOpt(opts ...gateway.IdentifyRateLimiterConfigOpt) ConfigOpt {
	return func(config *config) {
		config.IdentifyRateLimiterConfigOpts = append(opts, config.IdentifyRateLimiterConfigOpts...)
	}
}
