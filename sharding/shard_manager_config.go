package sharding

import (
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	Logger   :         log.Default(),
	Shards            : []int{0},
	ShardCount: 1,
	GatewayCreateFunc: func(token string, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway {
		return gateway.New(token, eventHandlerFunc, config)
	},
	GatewayConfig: &gateway.DefaultConfig,
}

type Config struct {
	Logger            log.Logger
	Shards            []int
	ShardCount        int
	GatewayCreateFunc func(token string, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway
	GatewayConfig *gateway.Config
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
		config.Shards = shards
	}
}

func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

func WithGatewayCreateFunc(gatewayCreateFunc func(token string, eventHandlerFunc gateway.EventHandlerFunc, config *gateway.Config) gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}
