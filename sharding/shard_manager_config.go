package sharding

import (
	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding/srate"
)

//goland:noinspection GoUnusedGlobalVariable
func DefaultConfig() *Config {
	return &Config{
		CustomShards:      false,
		GatewayCreateFunc: gateway.New,
	}
}

type Config struct {
	Logger                log.Logger
	CustomShards          bool
	Shards                *IntSet
	ShardCount            int
	GatewayCreateFunc     gateway.CreateFunc
	GatewayConfigOpts     []gateway.ConfigOpt
	RateLimiter           srate.Limiter
	RateLimiterConfigOpts []srate.ConfigOpt
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
func WithShards(shards ...int) ConfigOpt {
	return func(config *Config) {
		if config.Shards == nil {
			config.Shards = NewIntSet(shards...)
		}
		for _, shardID := range shards {
			config.Shards.Add(shardID)
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayCreateFunc(gatewayCreateFunc gateway.CreateFunc) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiter(rateLimiter srate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfigOpt(opts ...srate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}
