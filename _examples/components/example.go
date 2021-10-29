package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/core"

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

	disgo, err := bot.New(token,
		bot.WithGatewayOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuilds, discord.GatewayIntentGuildMessages, discord.GatewayIntentDirectMessages)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreateEvent) {
				if event.Message.Author.IsBot || event.Message.Author.IsSystem {
					return
				}
				if event.Message.Content == "test" {
					_, _ = event.Message.Reply(core.NewMessageCreateBuilder().
						SetActionRows(discord.NewActionRow(discord.NewDangerButton("danger", "danger"))).
						Build(),
					)
				}
			},
			OnButtonClick: func(event *events.ButtonClickEvent) {
				if event.CustomID == "danger" {
					_ = event.Create(core.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Ey that was danger").Build())
				}
			},
		}),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer disgo.Close()

	if err = disgo.ConnectGateway(); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
