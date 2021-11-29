package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token   = os.Getenv("token")
	//guildID = discord.Snowflake(os.Getenv("guild_id"))
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Infof("bot version: %s", info.Version)

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(discord.GatewayIntentsAll),
		),
		bot.WithCacheOpts(
			core.WithCacheFlags(core.CacheFlagsAll),
			core.WithMemberCachePolicy(core.MemberCachePolicyAll),
		),
		bot.WithMemberChunkingFilter(core.MemberChunkingFilterAll),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				if _, ok := event.Channel().(core.GuildThread); ok {
					println("MessageCreateEvent")
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

	defer disgo.Close(context.TODO())

	if err = disgo.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
