package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var token = os.Getenv("token")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Infof("bot version: %s", disgo.Version)

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
						println("MessageCreateEvent")
					}
				}
			},
			OnThreadCreate: func(event *events.ThreadCreate) {
				println("ThreadCreateEvent")
			},
			OnThreadUpdate: func(event *events.ThreadUpdate) {
				println("ThreadUpdateEvent")
			},
			OnThreadDelete: func(event *events.ThreadDelete) {
				println("ThreadDeleteEvent")
			},
			OnThreadHide: func(event *events.ThreadHide) {
				println("ThreadHideEvent")
			},
			OnThreadShow: func(event *events.ThreadShow) {
				println("ThreadShowEvent")
			},
			OnThreadMemberAdd: func(event *events.ThreadMemberAdd) {
				println("ThreadMemberAddEvent")
			},
			OnThreadMemberUpdate: func(event *events.ThreadMemberUpdate) {
				println("ThreadMemberUpdateEvent")
			},
			OnThreadMemberRemove: func(event *events.ThreadMemberRemove) {
				println("ThreadMemberRemoveEvent")
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
