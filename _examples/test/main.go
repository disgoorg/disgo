package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(logger))
	slog.Info("starting example...")
	slog.Info("disgo version", slog.Any("version", disgo.Version))

	r := handler.New()
	r.SlashCommand("/test", handleTestCommand)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithPresenceOpts(gateway.WithListeningActivity("your bullshit", gateway.WithActivityState("lol")), gateway.WithOnlineStatus(discord.OnlineStatusDND)),
		),
		bot.WithEventListeners(r),
	)
	if err != nil {
		slog.Error("error while building bot instance", slog.Any("err", err))
		return
	}

	if err = handler.SyncCommands(client, commands, []snowflake.ID{guildID}); err != nil {
		slog.Error("error while syncing commands", slog.Any("err", err))
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to discord", slog.Any("err", err))
	}

	defer client.Close(context.TODO())

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
