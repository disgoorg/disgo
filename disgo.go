package disgo

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/handlers"
	"github.com/DisgoOrg/disgo/httpserver"
)

// New creates a new core.Client instance with the provided bot token & ConfigOpt(s)
//goland:noinspection GoUnusedExportedFunction
func New(token string, opts ...bot.ConfigOpt) (bot.Client, error) {
	config := &bot.Config{}
	config.Apply(opts)

	if config.EventManagerConfig == nil {
		config.EventManagerConfig = &bot.DefaultEventManagerConfig
	}

	if config.EventManagerConfig.GatewayHandlers == nil {
		config.EventManagerConfig.GatewayHandlers = handlers.GetGatewayHandlers()
	}
	if config.EventManagerConfig.HTTPServerHandler == nil {
		config.EventManagerConfig.HTTPServerHandler = handlers.GetHTTPServerHandler()
	}

	return bot.BuildClient(token, *config,
		func(client bot.Client) gateway.EventHandlerFunc {
			return handlers.DefaultGatewayEventHandler(client)
		},
		func(client bot.Client) httpserver.EventHandlerFunc {
			return handlers.DefaultHTTPServerEventHandler(client)
		},
	)
}
