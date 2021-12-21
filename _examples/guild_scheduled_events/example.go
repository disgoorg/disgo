package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentGuildScheduledEvents|discord.GatewayIntentGuilds|discord.GatewayIntentGuildMessages),
		),
		bot.WithCacheOpts(
			core.WithCacheFlags(core.CacheFlagsAll),
			core.WithMemberCachePolicy(core.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(core.MemberChunkingFilterNone),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnGuildScheduledEventCreate: func(event *events.GuildScheduledEventCreateEvent) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUpdate: func(event *events.GuildScheduledEventUpdateEvent) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventDelete: func(event *events.GuildScheduledEventDeleteEvent) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUserAdd: func(event *events.GuildScheduledEventUserAddEvent) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUserRemove: func(event *events.GuildScheduledEventUserRemoveEvent) {
				log.Infof("%T\n", event)
			},
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				log.Infof("%T\n", event)
				if event.Message.Content != "test" {
					return
				}
				gse, _ := event.Guild().CreateGuildScheduledEvent(discord.GuildScheduledEventCreate{
					ChannelID:    "885677988916641802",
					Name:         "test",
					PrivacyLevel: discord.ScheduledEventPrivacyLevelGuildOnly,
					ScheduledStartTime: discord.Time{
						Time: time.Now().Add(time.Hour),
					},
					Description: "",
					EntityType:  discord.ScheduledEventEntityTypeVoice,
				})

				status := discord.ScheduledEventStatusActive
				gse, _ = gse.Update(discord.GuildScheduledEventUpdate{
					Status: &status,
				})
				//_ = gse.AudioChannel().Connect()

				time.Sleep(time.Second * 10)

				status = discord.ScheduledEventStatusCompleted
				gse, _ = gse.Update(discord.GuildScheduledEventUpdate{
					Status: &status,
				})
				//_ = gse.Guild().Disconnect()
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	if err = disgo.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	log.Info("Example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
