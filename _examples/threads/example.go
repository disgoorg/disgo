package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var token = os.Getenv("token")

func main() {
	slog.Info("starting example...")
	slog.Info("bot version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsAll),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagsAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterAll),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreate) {
				if channel, ok := event.Channel(); ok {
					if _, ok = channel.(discord.GuildThread); ok {
						slog.Info("MessageCreateEvent")
					}
				}
			},
			OnThreadCreate: func(event *events.ThreadCreate) {
				slog.Info("ThreadCreateEvent")
			},
			OnThreadUpdate: func(event *events.ThreadUpdate) {
				slog.Info("ThreadUpdateEvent")
			},
			OnThreadDelete: func(event *events.ThreadDelete) {
				slog.Info("ThreadDeleteEvent")
			},
			OnThreadHide: func(event *events.ThreadHide) {
				slog.Info("ThreadHideEvent")
			},
			OnThreadShow: func(event *events.ThreadShow) {
				slog.Info("ThreadShowEvent")
			},
			OnThreadMemberAdd: func(event *events.ThreadMemberAdd) {
				slog.Info("ThreadMemberAddEvent")
			},
			OnThreadMemberUpdate: func(event *events.ThreadMemberUpdate) {
				slog.Info("ThreadMemberUpdateEvent")
			},
			OnThreadMemberRemove: func(event *events.ThreadMemberRemove) {
				slog.Info("ThreadMemberRemoveEvent")
			},
		}),
	)
	if err != nil {
		slog.Error("error while building bot instance", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to discord", slog.Any("err", err))
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
