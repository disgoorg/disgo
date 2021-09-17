package sharding

import (
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

type Config struct {
	Logger            log.Logger
	Shards            []int
	ShardCount        int
	GatewayCreateFunc func() gateway.Gateway
	EventHandlerFunc  gateway.EventHandlerFunc
}
