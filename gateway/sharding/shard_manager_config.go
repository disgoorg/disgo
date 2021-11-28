package sharding

import (
	"github.com/DisgoOrg/disgo/gateway"
	srate2 "github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	CustomShards: false,
	GatewayCreateFunc: func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
		return gateway.New(token, url, shardID, shardCount, eventHandlerFunc, config)
	},
	GatewayConfig: &gateway.DefaultConfig,
	RateLimiter:   srate2.NewLimiter(&srate2.DefaultConfig),
}

type Config struct {
	Logger            log.Logger
	CustomShards      bool
	Shards            *IntSet
	ShardCount        int
	GatewayCreateFunc func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway
	GatewayConfig     *gateway.Config
	RateLimiter       srate2.Limiter
	RateLimiterConfig *srate2.Config
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
		config.CustomShards = true
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
		config.CustomShards = true
		config.ShardCount = shardCount
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayCreateFunc(gatewayCreateFunc func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfig(gatewayConfig gateway.Config) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfig = &gatewayConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.GatewayConfig == nil {
			config.GatewayConfig = &gateway.DefaultConfig
		}
		config.GatewayConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiter(rateLimiter srate2.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfig(rateConfig srate2.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfigOpt(opts ...srate2.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &srate2.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}
