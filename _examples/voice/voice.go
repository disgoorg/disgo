package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)

	logger := log.New(log.LstdFlags | log.Lshortfile)
	logger.SetLevel(log.LevelInfo)

	file, _ := os.Open("nico.dca")
	defer file.Close()
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuildVoiceStates), gateway.WithLogger(logger)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client(), file)
		}),
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

func play(client bot.Client, reader io.Reader) {
	connection, err := client.ConnectChannel(context.Background(), 817327181659111454, 982083072067530762)
	if err != nil {
		client.Logger().Error("error connecting to voice channel: ", err)
		return
	}

	time.Sleep(2 * time.Second)

	connection.SetSendHandler(newReaderSendHandler(reader))
	// echo := newEchoHandler()
	// connection.SetReceiveHandler(echo)
	// connection.SetSendHandler(echo)
}
