package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages)),
		bot.WithEventListenerFunc(func(event *events.MessageCreate) {
			if event.Message.Author.Bot || event.Message.Author.System {
				return
			}
			if event.Message.Content == "test" {
				_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
					AddActionRow(discord.NewDangerButton("danger", "danger")).
					SetMessageReferenceByID(event.Message.ID).
					Build(),
				)
			}
		}),
		bot.WithEventListenerFunc(func(event *events.ComponentInteractionCreate) {
			if event.ButtonInteractionData().CustomID() == "danger" {
				_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetEphemeral(true).SetContent("Ey that was danger").Build())
			}
		}),
	)
	if err != nil {
		slog.Error("error while building bot", slog.Any("err", err))
		return
	}
	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
		return
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
