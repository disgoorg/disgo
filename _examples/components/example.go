package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", info.Version)

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				if event.Message.Author.BotUser || event.Message.Author.System {
					return
				}
				if event.Message.Content == "test" {
					_, _ = event.Message.Reply(discord.NewMessageCreateBuilder().
						AddActionRow(discord.NewDangerButton("danger", "danger")).
						Build(),
					)
				}
			},
			OnComponentInteraction: func(event *events.ComponentInteractionEvent) {
				if event.ButtonInteractionData().CustomID == "danger" {
					_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Ey that was danger").Build())
				}
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer disgo.Close(context.TODO())

	if err = disgo.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
