package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/log"

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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")

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
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUpdate: func(event *events.GuildScheduledEventUpdate) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventDelete: func(event *events.GuildScheduledEventDelete) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUserAdd: func(event *events.GuildScheduledEventUserAdd) {
				log.Infof("%T\n", event)
			},
			OnGuildScheduledEventUserRemove: func(event *events.GuildScheduledEventUserRemove) {
				log.Infof("%T\n", event)
			},
			OnMessageCreate: func(event *events.MessageCreate) {
				log.Infof("%T\n", event)
				if event.Message.Content != "test" {
					return
				}
				gse, _ := event.Client().Rest().CreateGuildScheduledEvent(*event.GuildID, discord.GuildScheduledEventCreate{
					ChannelID:          885677988916641802,
					Name:               "test",
					PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
					ScheduledStartTime: time.Now().Add(time.Hour),
					Description:        "",
					EntityType:         discord.ScheduledEventEntityTypeVoice,
				})

				status := discord.ScheduledEventStatusActive
				gse, _ = event.Client().Rest().UpdateGuildScheduledEvent(gse.GuildID, gse.ID, discord.GuildScheduledEventUpdate{
					Status: &status,
				})
				//_ = gse.AudioChannel().Connect()

				time.Sleep(time.Second * 10)

				status = discord.ScheduledEventStatusCompleted
				gse, _ = event.Client().Rest().UpdateGuildScheduledEvent(gse.GuildID, gse.ID, discord.GuildScheduledEventUpdate{
					Status: &status,
				})
				//_ = gse.Guilds().Disconnect()
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	log.Info("Example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
