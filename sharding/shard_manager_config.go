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

type config struct {
	// Logger is the logger of the ShardManager. Defaults to log.Default()
	Logger *slog.Logger
	// ShardIDs is a map of shardIDs the ShardManager should manage. Leave this nil to manage all shards.
	ShardIDs map[int]struct{}
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
	// RateLimiter is the RateLimiter which is used by the ShardManager. Defaults to NewRateLimiter()
	RateLimiter RateLimiter
	// RateLimiterConfigOpts are the RateLimiterConfigOpt(s) which are applied to the RateLimiter.
	RateLimiterConfigOpts []RateLimiterConfigOpt
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Server.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "sharding"))
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateLimiterConfigOpts...)
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
		if config.ShardIDs == nil {
			config.ShardIDs = map[int]struct{}{}
		}
		for _, shardID := range shardIDs {
			config.ShardIDs[shardID] = struct{}{}
		}
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

// WithRateLimiter lets you inject your own RateLimiter into the ShardManager.
func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfigOpt lets you configure the default RateLimiter used by the ShardManager.
func WithRateLimiterConfigOpt(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}

// WithDefaultRateLimiterConfigOpt lets you configure the default RateLimiter used by the ShardManager and prepend the options to the existing ones.
func WithDefaultRateLimiterConfigOpt(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *config) {
		config.RateLimiterConfigOpts = append(opts, config.RateLimiterConfigOpts...)
	}
}
