package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/snekROmonoro/disgo"
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/discord"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
	"github.com/snekROmonoro/disgo/sharding"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.Any("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithShardManagerConfigOpts(
			sharding.WithShardIDs(0, 1),
			sharding.WithShardCount(2),
			sharding.WithAutoScaling(true),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages),
				gateway.WithCompress(true),
			),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
			OnGuildReady: func(event *events.GuildReady) {
				slog.Info("guild %s ready", event.GuildID)
			},
			OnGuildsReady: func(event *events.GuildsReady) {
				slog.Info("guilds on shard %d ready", event.ShardID)
			},
		}),
	)
	if err != nil {
		slog.Error("error while building disgo", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenShardManager(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
		return
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}
	_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(event.Message.Content).Build())
}
