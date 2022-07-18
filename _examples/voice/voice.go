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
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Info("starting up")

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
	)
	if err != nil {
		log.Fatal("error creating client: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error connecting to gateway: ", err)
	}

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func play(client bot.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := client.OpenVoice(ctx, guildID, channelID, false, false)
	if err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err = conn.WaitUntilConnected(ctx); err != nil {
		panic("error waiting for voice connection: " + err.Error())
	}

	println("starting playback")

	if err = conn.Speaking(voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	writeOpus(conn.UDPConn())
}

func writeOpus(w io.Writer) {
	file, err := os.Open("nico.dca")
	if err != nil {
		panic("error opening file: " + err.Error())
	}
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenBuf [4]byte
	for range ticker.C {
		_, err = io.ReadFull(file, lenBuf[:])
		if err != nil {
			if err == io.EOF {
				_ = file.Close()
				return
			}
			panic("error reading file: " + err.Error())
			return
		}

		// Read the integer
		frameLen := int64(binary.LittleEndian.Uint32(lenBuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, file, frameLen)
		if err != nil && err != io.EOF {
			_ = file.Close()
			return
		}
	}
}
