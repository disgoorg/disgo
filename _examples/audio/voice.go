package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/audio"
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
	log.Info("starting up")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithGatewayIntents(discord.GatewayIntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
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

func play(client bot.Client) {
	connection, err := client.ConnectChannel(context.Background(), 817327181659111454, 982083072067530762, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	rs, err := http.Get("https://p.scdn.co/mp3-preview/ee121ca281c629bb4e99c33d877fe98fbb752289?cid=774b29d4f13844c495f206cafdad9c86")
	if err != nil {
		panic("error getting audio: " + err.Error())
	}

	provider, writeFunc := audio.NewMP3PCMFrameProvider(nil)

	go func() {
		defer rs.Body.Close()
		io.Copy(writeFunc, rs.Body)
	}()

	connection.SetOpusFrameProvider(audio.NewPCMOpusProvider(nil, provider))

	println("voice: ready")
}
