package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuildVoiceStates)),
		bot.WithEventListenerFunc(messageCreateListener),
	)
	if err != nil {
		log.Fatal("error creating client: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.ConnectGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func messageCreateListener(e *events.MessageCreate) {
	connection, err := e.Client().ConnectChannel(context.Background(), *e.Message.GuildID, e.Message.ChannelID)
	if err != nil {
		return

	}

	connection.SetSendHandler()
}
