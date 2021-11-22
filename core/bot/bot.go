package bot

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/collectors"
	"github.com/DisgoOrg/disgo/core/handlers"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/gateway/sharding"
	"github.com/DisgoOrg/disgo/gateway/sharding/srate"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/rrate"
	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"
)

func New(token string, opts ...ConfigOpt) (*core.Bot, error) {
	config := &Config{}
	config.Apply(opts)

	if config.EventManagerConfig.GatewayHandlers == nil {
		config.EventManagerConfig.GatewayHandlers = handlers.GetGatewayHandlers()
	}
	if config.EventManagerConfig.HTTPServerHandler == nil {
		config.EventManagerConfig.HTTPServerHandler = handlers.GetHTTPServerHandler()
	}

	return buildBot(token, *config)
}

func buildBot(token string, config Config) (*core.Bot, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := IDFromToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting application id from BotToken")
	}
	bot := &core.Bot{
		Token: token,
	}

	// TODO: figure out how we handle different application & client ids
	bot.ApplicationID = *id
	bot.ClientID = *id

	if config.Logger == nil {
		config.Logger = log.Default()
	}
	bot.Logger = config.Logger

	if config.RestClient == nil {
		if config.RestClientConfig == nil {
			config.RestClientConfig = &rest.DefaultConfig
		}
		if config.RestClientConfig.Logger == nil {
			config.RestClientConfig.Logger = config.Logger
		}
		if config.RestClientConfig.Headers == nil {
			config.RestClientConfig.Headers = http.Header{}
		}
		if _, ok := config.RestClientConfig.Headers["Authorization"]; !ok {
			config.RestClientConfig.Headers["Authorization"] = []string{discord.TokenTypeBot.Apply(token)}
		}

		if config.RestClientConfig.RateLimiterConfig == nil {
			config.RestClientConfig.RateLimiterConfig = &rrate.DefaultConfig
		}
		if config.RestClientConfig.RateLimiterConfig.Logger == nil {
			config.RestClientConfig.RateLimiterConfig.Logger = config.Logger
		}
		config.RestClient = rest.NewClient(config.RestClientConfig)
	}

	if config.RestServices == nil {
		config.RestServices = rest.NewServices(bot.Logger, config.RestClient)
	}
	bot.RestServices = config.RestServices

	if config.EventManager == nil {
		if config.EventManagerConfig == nil {
			config.EventManagerConfig = &core.DefaultEventManagerConfig
		}

		config.EventManager = core.NewEventManager(bot, config.EventManagerConfig)
	}
	bot.EventManager = config.EventManager

	if config.Collectors == nil {
		if config.CollectorsConfig == nil {
			config.CollectorsConfig = &collectors.DefaultConfig
		}
		config.Collectors = core.NewCollectors(bot, *config.CollectorsConfig)
	}
	bot.Collectors = config.Collectors

	if config.Gateway == nil && config.GatewayConfig != nil {
		var gatewayRs *discord.Gateway
		gatewayRs, err = bot.RestServices.GatewayService().GetGateway()
		if err != nil {
			return nil, err
		}
		config.Gateway = gateway.New(token, gatewayRs.URL, 0, 0, handlers.DefaultGatewayEventHandler(bot), config.GatewayConfig)
	}
	bot.Gateway = config.Gateway

	if config.ShardManager == nil && config.ShardManagerConfig != nil {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = bot.RestServices.GatewayService().GetGatewayBot()
		if err != nil {
			return nil, err
		}

		if config.ShardManagerConfig.RateLimiterConfig == nil {
			config.ShardManagerConfig.RateLimiterConfig = &srate.DefaultConfig
		}
		if config.ShardManagerConfig.RateLimiterConfig.Logger == nil {
			config.ShardManagerConfig.RateLimiterConfig.Logger = config.Logger
		}
		if config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency == 0 {
			config.ShardManagerConfig.RateLimiterConfig.MaxConcurrency = gatewayBotRs.SessionStartLimit.MaxConcurrency
		}

		// apply recommended shard count
		if !config.ShardManagerConfig.CustomShards {
			config.ShardManagerConfig.ShardCount = gatewayBotRs.Shards
			config.ShardManagerConfig.Shards = sharding.NewIntSet()
			for i := 0; i < gatewayBotRs.Shards; i++ {
				config.ShardManagerConfig.Shards.Add(i)
			}
		}
		if config.ShardManager == nil {
			config.ShardManager = sharding.New(token, gatewayBotRs.URL, handlers.DefaultGatewayEventHandler(bot), config.ShardManagerConfig)
		}
	}
	bot.ShardManager = config.ShardManager

	if config.HTTPServer == nil && config.HTTPServerConfig != nil {
		if config.HTTPServerConfig.Logger == nil {
			config.HTTPServerConfig.Logger = config.Logger
		}
		config.HTTPServer = httpserver.New(handlers.DefaultHTTPServerEventHandler(bot), config.HTTPServerConfig)
	}
	bot.HTTPServer = config.HTTPServer

	if config.AudioController == nil {
		config.AudioController = core.NewAudioController(bot)
	}
	bot.AudioController = config.AudioController

	if config.MemberChunkingManager == nil {
		if config.MemberChunkingFilter == nil {
			config.MemberChunkingFilter = &core.MemberChunkingFilterNone
		}
		config.MemberChunkingManager = core.NewMemberChunkingManager(bot, *config.MemberChunkingFilter)
	}
	bot.MemberChunkingManager = config.MemberChunkingManager

	if config.EntityBuilder == nil {
		config.EntityBuilder = core.NewEntityBuilder(bot)
	}
	bot.EntityBuilder = config.EntityBuilder

	if config.Caches == nil {
		if config.CacheConfig == nil {
			config.CacheConfig = &core.DefaultCacheConfig
		}
		config.Caches = core.NewCaches(*config.CacheConfig)
	}
	bot.Caches = config.Caches

	return bot, nil
}

// IDFromToken returns the applicationID from the BotToken
//goland:noinspection GoUnusedExportedFunction
func IDFromToken(token string) (*discord.Snowflake, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return nil, discord.ErrInvalidBotToken
	}
	byteID, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return nil, err
	}
	strID := discord.Snowflake(byteID)
	return &strID, nil
}
