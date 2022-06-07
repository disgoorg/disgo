package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/log"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Info("starting up")

	file, _ := os.Open("nico.dca")
	defer file.Close()
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuildVoiceStates)),
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
	connection, err := client.ConnectChannel(context.Background(), 817327181659111454, 982083072067530762, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	println("starting playback")

	//connection.SetSendHandler(newReaderSendHandler(reader))

	if err = connection.Speaking(voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	writeOpus(connection.UDP(), reader)

	/*echo := &echoHandler{}
	connection.SetSendHandler(echo)
	connection.SetReceiveHandler(echo)
	*/

	//newEcho2(connection)
}
