package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/disgo/sharding/rate"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", info.Version)

	disgo, err := core.NewBot(token,
		core.WithShardManagerConfigOpts(
			sharding.WithShards(0, 1, 2),
			sharding.WithShardCount(3),
			sharding.WithGatewayConfigOpts(
				gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages),
				gateway.WithCompress(true),
			),
			sharding.WithRateLimiterConfigOpt(
				rate.WithMaxConcurrency(2),
			),
		),
		core.WithCacheConfig(core.CacheConfig{CacheFlags: core.CacheFlagsDefault}),
		core.WithEventListeners(&core.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgo: %s", err)
	}

	defer disgo.Close()

	if errs := disgo.Connect(); errs != nil {
		log.Fatal("error while connecting to gateway: ", errs)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *core.MessageCreateEvent) {
	if event.Message.Author.IsBot {
		return
	}
	_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
