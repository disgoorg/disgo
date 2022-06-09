package main

import (
	"context"
	"encoding/binary"
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

func play(client bot.Client, reader io.ReadCloser) {
	connection, err := client.ConnectChannel(context.Background(), 817327181659111454, 982083072067530762, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	println("starting playback")

	if err = connection.Speaking(voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	writeOpus(connection.UDP(), reader)
}

func writeOpus(w io.Writer, reader io.ReadCloser) {
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenbuf [4]byte
	for range ticker.C {
		_, err := io.ReadFull(reader, lenbuf[:])
		if err != nil {
			if err == io.EOF {
				reader.Close()
				return
			}
			return
		}

		// Read the integer
		framelen := int64(binary.LittleEndian.Uint32(lenbuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, reader, framelen)
		if err != nil && err != io.EOF {
			reader.Close()
			return
		}
	}
}
