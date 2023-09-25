package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
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

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error connecting to voice gateway", slog.Any("err", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func play(client bot.Client) {
	conn := client.VoiceManager().CreateConn(guildID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := conn.Open(ctx, channelID, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	defer func() {
		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel2()
		conn.Close(ctx2)
	}()

	slog.Info("starting playback")

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}

	if _, err := conn.UDP().Write(voice.SilenceAudioFrame); err != nil {
		panic("error sending silence: " + err.Error())
	}
	for {
		packet, err := conn.UDP().ReadPacket()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("connection closed")
				return
			}
			slog.Info("error while reading from reader", slog.Any("err", err))
			continue
		}
		if _, err = conn.UDP().Write(packet.Opus); err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("connection closed")
				return
			}
			slog.Info("error while writing to UDPConn", slog.Any("err", err))
			continue
		}
	}
}
