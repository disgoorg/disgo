package sharding

import (
	"log/slog"

	"github.com/disgoorg/disgo/gateway"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:            slog.Default(),
		GatewayCreateFunc: gateway.New,
		ShardSplitCount:   ShardSplitCount,
	}
}

// Config lets you configure your ShardManager instance.
type Config struct {
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

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateLimiterConfigOpts...)
	}
}

// WithLogger sets the logger of the ShardManager.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithShardIDs sets the shardIDs the ShardManager should manage.
func WithShardIDs(shardIDs ...int) ConfigOpt {
	return func(config *Config) {
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
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

// WithShardSplitCount sets the count a shard should be split into if it is too large.
// This is only used if AutoScaling is enabled.
func WithShardSplitCount(shardSplitCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardSplitCount = shardSplitCount
	}
}

// WithAutoScaling sets whether the ShardManager should automatically re-shard shards if they are too large. This is disabled by default.
func WithAutoScaling(autoScaling bool) ConfigOpt {
	return func(config *Config) {
		config.AutoScaling = autoScaling
	}
}

// WithGatewayCreateFunc sets the function which is used by the ShardManager to create a new gateway.Gateway.
func WithGatewayCreateFunc(gatewayCreateFunc gateway.CreateFunc) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

// WithGatewayConfigOpts lets you configure the gateway.Gateway created by the ShardManager.
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithRateLimiter lets you inject your own RateLimiter into the ShardManager.
func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfigOpt lets you configure the default RateLimiter used by the ShardManager.
func WithRateLimiterConfigOpt(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}
