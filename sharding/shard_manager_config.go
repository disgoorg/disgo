package sharding

import (
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/sharding/rate"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	CustomShards: false,
	GatewayCreateFunc: func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
		return gateway.New(token, url, shardID, shardCount, eventHandlerFunc, config)
	},
	GatewayConfig: &gateway.DefaultConfig,
	RateLimiter:   rate.NewLimiter(&rate.DefaultConfig),
}

type Config struct {
	Logger            log.Logger
	CustomShards      bool
	Shards            *IntSet
	ShardCount        int
	GatewayCreateFunc func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway
	GatewayConfig     *gateway.Config
	RateLimiter       rate.Limiter
	RateLimiterConfig *rate.Config
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

func WithShards(shards ...int) ConfigOpt {
	return func(config *Config) {
		config.CustomShards = true
		if config.Shards == nil {
			config.Shards = NewIntSet(shards...)
		}
		for _, shardID := range shards {
			config.Shards.Add(shardID)
		}
	}
}

func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.CustomShards = true
		config.ShardCount = shardCount
	}
}

func WithGatewayCreateFunc(gatewayCreateFunc func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

func WithGatewayConfig(gatewayConfig gateway.Config) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfig = &gatewayConfig
	}
}

func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Apply(opts)
	}
}

func WithRateLimiter(rateLimiter rate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfig(rateConfig rate.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateConfig
	}
}

func WithRateLimiterConfigOpt(opts ...rate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &rate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}
