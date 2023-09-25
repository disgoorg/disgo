package main

import (
	"context"
	"encoding/binary"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")
)

func main() {
	slog.Info("starting up")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	s := make(chan os.Signal, 1)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client(), s)
		}),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		client.Close(ctx)
	}()

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error connecting to gateway", slog.Any("error", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func play(client bot.Client, closeChan chan os.Signal) {
	conn := client.VoiceManager().CreateConn(guildID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := conn.Open(ctx, channelID, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer closeCancel()
		conn.Close(closeCtx)
	}()

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	writeOpus(conn.UDP())
	closeChan <- syscall.SIGTERM
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
