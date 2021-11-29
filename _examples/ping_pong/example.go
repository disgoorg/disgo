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
	"github.com/DisgoOrg/log"
)

func main() {
	disgo, err := bot.New(os.Getenv("token"),
		bot.WithGatewayOpts(
			gateway.WithGatewayIntents(
				discord.GatewayIntentGuilds,
				discord.GatewayIntentGuildMessages,
				discord.GatewayIntentDirectMessages,
			),
		),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsDefault)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: onMessageCreate,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	defer disgo.Close(context.TODO())

	if err = disgo.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onMessageCreate(event *events.MessageCreateEvent) {
	var message string
	if event.Message.Content == "ping" {
		message = "pong"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}
	if message != "" {
		_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}
