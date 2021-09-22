package rate

import (
	"github.com/DisgoOrg/log"
)

var DefaultConfig = Config{
	Logger:         log.Default(),
	MaxConcurrency: 1,
}

type Config struct {
	Logger         log.Logger
	MaxConcurrency int
}
