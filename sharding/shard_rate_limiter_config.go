package sharding

import "github.com/DisgoOrg/log"

var DefaultRateLimitConfig = RateLimitConfig{
	Logger: log.Default(),
}

type RateLimitConfig struct {
	Logger         log.Logger
	MaxConcurrency int
}
