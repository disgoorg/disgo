package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo"
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/events"

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

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				if event.Message.Author.BotUser || event.Message.Author.System {
					return
				}
				if event.Message.Content == "test" {
					_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
						AddActionRow(discord.NewDangerButton("danger", "danger")).
						SetMessageReferenceByID(event.Message.ID).
						Build(),
					)
				}
			},
			OnComponentInteraction: func(event *events.ComponentInteractionEvent) {
				if event.ButtonInteractionData().CustomID() == "danger" {
					_, _ = event.Client().Rest().ChannelService().CreateMessage(event.ChannelID(), discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Ey that was danger").Build())
				}
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
