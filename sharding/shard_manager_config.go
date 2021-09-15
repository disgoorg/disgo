package sharding

import "github.com/DisgoOrg/disgo/gateway"

type Config struct {
	Shards           []int
	ShardCount       int
	Gateway          func() gateway.Gateway
	EventHandlerFunc gateway.EventHandlerFunc
}
