package main

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/audio"
	"github.com/disgoorg/disgo/audio/opus"
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

	connection.Speaking(voice.SpeakingFlagMicrophone)

	connection.UDP().Write(voice.SilenceAudioFrames)

	buff := &bytes.Buffer{}

	encoder, err := opus.NewEncoder(24000, 2, opus.ApplicationAudio)
	if err != nil {
		panic("NewPCMOpusProvider: " + err.Error())
	}
	if err = encoder.Ctl(opus.SetBitrate(64000)); err != nil {
		panic("SetBitrate: " + err.Error())
	}

	connection.SetOpusFrameProvider(audio.NewPCMOpusProvider(encoder, audio.NewPCMStreamProvider(buff)))
	connection.SetOpusFraneReceiver(
		audio.NewPCMOpusReceiver(
			nil,
			audio.NewPCMCombinerReceiver(
				audio.NewSampleRateCombinedReceiver(
					nil,
					48000,
					24000,
					audio.NewPCMCombinedStreamReceiver(buff),
				),
			),
			nil,
		),
	)

	println("voice: ready")
}
