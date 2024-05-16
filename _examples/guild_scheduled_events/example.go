package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	slog.Info("starting example...")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildScheduledEvents|gateway.IntentGuilds|gateway.IntentGuildMessages),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagsAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterNone),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnGuildScheduledEventCreate: func(event *events.GuildScheduledEventCreate) {
				slog.Info("OnGuildScheduledEventCreate")
			},
			OnGuildScheduledEventUpdate: func(event *events.GuildScheduledEventUpdate) {
				slog.Info("OnGuildScheduledEventUpdate")
			},
			OnGuildScheduledEventDelete: func(event *events.GuildScheduledEventDelete) {
				slog.Info("OnGuildScheduledEventDelete")
			},
			OnGuildScheduledEventUserAdd: func(event *events.GuildScheduledEventUserAdd) {
				slog.Info("OnGuildScheduledEventUserAdd")
			},
			OnGuildScheduledEventUserRemove: func(event *events.GuildScheduledEventUserRemove) {
				slog.Info("OnGuildScheduledEventUserRemove")
			},
			OnMessageCreate: func(event *events.MessageCreate) {
				slog.Info("OnMessageCreate")
				if event.Message.Content != "test" {
					return
				}
				gse, _ := event.Client().Rest().CreateGuildScheduledEvent(*event.GuildID, discord.GuildScheduledEventCreate{
					ChannelID:          885677988916641802,
					Name:               "test",
					PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
					ScheduledStartTime: time.Now().Add(time.Hour),
					Description:        "test",
					EntityType:         discord.ScheduledEventEntityTypeVoice,
				})

				status := discord.ScheduledEventStatusActive
				gse, _ = event.Client().Rest().UpdateGuildScheduledEvent(gse.GuildID, gse.ID, discord.GuildScheduledEventUpdate{
					Status: &status,
				})

				time.Sleep(time.Second * 10)

				status = discord.ScheduledEventStatusCompleted
				gse, _ = event.Client().Rest().UpdateGuildScheduledEvent(gse.GuildID, gse.ID, discord.GuildScheduledEventUpdate{
					Status: &status,
				})
			},
		}),
	)
	if err != nil {
		slog.Error("error while building bot instance", slog.Any("err", err))
		return
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to discord", slog.Any("err", err))
	}

	slog.Info("Example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
