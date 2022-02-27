package sharding

import (
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	GatewayCreateFunc: func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
		return gateway.New(token, url, eventHandlerFunc, config)
	},
	GatewayConfig: &gateway.DefaultConfig,
	RateLimiter:   srate.NewLimiter(&srate.DefaultConfig),
}

type Config struct {
	Logger            log.Logger
	AutoScaling       bool
	CustomShards      bool
	ShardIDs          map[int]struct{}
	ShardCount        int
	GatewayCreateFunc func(token string, url string, shardID int, shardCount int, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway
	GatewayConfig     *gateway.Config
	RateLimiter       srate.Limiter
	RateLimiterConfig *srate.Config
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
func WithAutoScaling(autoScaling bool) ConfigOpt {
	return func(config *Config) {
		config.AutoScaling = autoScaling
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShards(shardIDs ...int) ConfigOpt {
	return func(config *Config) {
		config.CustomShards = true
		if config.ShardIDs == nil {
			config.ShardIDs = make(map[int]struct{}, len(shardIDs))
		}
		for _, shardID := range shardIDs {
			config.ShardIDs[shardID] = struct{}{}
		}
		if config.ShardCount == 0 {
			config.ShardCount = len(config.ShardIDs)
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
func WithRateLimiter(rateLimiter srate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfig(rateConfig srate.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfigOpt(opts ...srate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &srate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}
