package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/cache"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("token")
	//guildID = snowflake.Snowflake(os.Getenv("guild_id"))
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Infof("bot version: %s", info.Version)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentsAll),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagsAll),
			cache.WithMemberCachePolicy(cache.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterAll),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				if channel, ok := event.Channel(); ok {
					if _, ok = channel.(discord.GuildThread); ok {
						println("MessageCreateEvent")
					}
				}
			},
			OnThreadCreate: func(event *events.ThreadCreateEvent) {
				println("ThreadCreateEvent")
			},
			OnThreadUpdate: func(event *events.ThreadUpdateEvent) {
				println("ThreadUpdateEvent")
			},
			OnThreadDelete: func(event *events.ThreadDeleteEvent) {
				println("ThreadDeleteEvent")
			},
			OnThreadHide: func(event *events.ThreadHideEvent) {
				println("ThreadHideEvent")
			},
			OnThreadShow: func(event *events.ThreadShowEvent) {
				println("ThreadShowEvent")
			},
			OnThreadMemberAdd: func(event *events.ThreadMemberAddEvent) {
				println("ThreadMemberAddEvent")
			},
			OnThreadMemberUpdate: func(event *events.ThreadMemberUpdateEvent) {
				println("ThreadMemberUpdateEvent")
			},
			OnThreadMemberRemove: func(event *events.ThreadMemberRemoveEvent) {
				println("ThreadMemberRemoveEvent")
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
