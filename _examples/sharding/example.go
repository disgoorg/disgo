package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"

	"github.com/DisgoOrg/disgo/events"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/sharding"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", info.Version)

	disgo, err := bot.New(token,
		bot.WithShardManagerOpts(
			sharding.WithShards(0, 1, 2),
			sharding.WithShardCount(3),
			sharding.WithGatewayConfigOpts(
				gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages),
				gateway.WithCompress(true),
			),
		),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsDefault)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
			OnGuildReady: func(event *events.GuildReadyEvent) {
				log.Infof("guild %s ready", event.GuildID)
			},
			OnGuildsReady: func(event *events.GuildsReadyEvent) {
				log.Infof("guilds on shard %d ready", event.ShardID)
			},
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgo: %s", err)
	}

	defer disgo.Close()

	if errs := disgo.ConnectShardManager(); errs != nil {
		log.Fatal("error while connecting to gateway: ", errs)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	if event.Message.Author.IsBot {
		return
	}
	_, _ = event.Message.Reply(core.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
