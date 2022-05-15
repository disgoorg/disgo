package main

import (
	"context"
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
	"github.com/disgoorg/log"
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
			gateway.WithGatewayIntents(discord.GatewayIntentGuildScheduledEvents|discord.GatewayIntentGuilds|discord.GatewayIntentGuildMessages),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagsAll),
			cache.WithMemberCachePolicy(cache.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterNone),
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

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	log.Info("Example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
