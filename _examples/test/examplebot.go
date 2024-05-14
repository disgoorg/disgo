package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	//go:embed gopher.png
	gopher []byte
)

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.Any("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsNonPrivileged, gateway.IntentMessageContent),
			gateway.WithPresenceOpts(gateway.WithListeningActivity("your bullshit", gateway.WithActivityState("lol")), gateway.WithOnlineStatus(discord.OnlineStatusDND)),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagsAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterNone),
		bot.WithEventListeners(listener),
	)
	if err != nil {
		slog.Error("error while building bot instance", slog.Any("err", err))
		return
	}

	registerCommands(client)

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to discord", slog.Any("err", err))
	}

	defer client.Close(context.TODO())

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
