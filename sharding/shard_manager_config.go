package sharding

import (
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding/srate"
	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		GatewayCreateFunc: gateway.New,
	}
}

type Config struct {
	Logger                log.Logger
	ShardIDs              map[int]struct{}
	ShardCount            int
	AutoScaling           bool
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
	if c.RateLimiter == nil {
		c.RateLimiter = srate.NewLimiter(c.RateLimiterConfigOpts...)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithShards(shards ...int) ConfigOpt {
	return func(config *Config) {
		if config.ShardIDs == nil {
			config.ShardIDs = map[int]struct{}{}
		}
		for _, shardID := range shards {
			config.ShardIDs[shardID] = struct{}{}
		}
	}
}

func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

func WithAutoScaling(autoScaling bool) ConfigOpt {
	return func(config *Config) {
		config.AutoScaling = autoScaling
	}
}

func WithGatewayCreateFunc(gatewayCreateFunc gateway.CreateFunc) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

func WithRateLimiter(rateLimiter srate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfigOpt(opts ...srate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}
