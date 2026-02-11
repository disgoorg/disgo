package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/godave"
	"github.com/disgoorg/snowflake/v2"
	"github.com/hajimehoshi/go-mp3"
	"github.com/kazzmir/opus-go/opus"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/voice"
)

const (
	SampleRate = 48000
	Channels   = 2
)

var (
	token     = os.Getenv("disgo_token")
	guildID   = snowflake.GetEnv("disgo_guild_id")
	channelID = snowflake.GetEnv("disgo_channel_id")

	// 48kHz stereo mp3 data
	//go:embed test.mp3
	testDataMP3 []byte
)

func main() {
	slog.Info("starting up")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
		bot.WithVoiceManagerConfigOpts(
			voice.WithDaveSessionCreateFunc(godave.NewNoopSession),
		),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
		return
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
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func play(client *bot.Client) {
	conn := client.VoiceManager.CreateConn(guildID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := conn.Open(ctx, channelID, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}

	provider, err := newOpusFrameProvider(bytes.NewReader(testDataMP3))
	if err != nil {
		panic("error creating opus frame provider: " + err.Error())
	}

	conn.SetOpusFrameProvider(provider)
}

func newOpusFrameProvider(r io.Reader) (voice.OpusFrameProvider, error) {
	rc := newReadCanceller(r)

	decoder, err := mp3.NewDecoder(rc)
	if err != nil {
		return nil, err
	}

	encoder, err := opus.NewEncoder(SampleRate, Channels, opus.ApplicationVoIP)
	if err != nil {
		return nil, err
	}

	return &mp3Reader{
		decoder: decoder,
		encoder: encoder,
		// 16bit pcm stereo = 4 bytes per sample
		buf:     make([]byte, 4096),
		opusBuf: make([]byte, 4000),
	}, nil
}

var _ voice.OpusFrameProvider = (*mp3Reader)(nil)

type mp3Reader struct {
	rc      io.ReadCloser
	decoder *mp3.Decoder
	encoder *opus.Encoder
	eof     bool
	buf     []byte
	opusBuf []byte
}

func (m *mp3Reader) ProvideOpusFrame() ([]byte, error) {
	if m.eof {
		return nil, io.EOF
	}

	neededBytes := voice.OpusFrameSize * Channels * 2

	pcmBytes := make([]byte, neededBytes)
	n, err := io.ReadFull(m.decoder, pcmBytes)
	if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
		if n == 0 {
			slog.Info("reached end of mp3 data")
			m.eof = true
			return nil, io.EOF
		}
		// zero-pad remainder
		for i := n; i < neededBytes; i++ {
			pcmBytes[i] = 0
		}
		m.eof = true
	} else if err != nil {
		return nil, err
	}

	pcm := bytesToInt16LE(pcmBytes)

	nn, err := m.encoder.Encode(pcm, voice.OpusFrameSize, m.opusBuf)
	if err != nil {
		return nil, err
	}

	return m.opusBuf[:nn], nil
}

func (m *mp3Reader) Close() {
	m.rc.Close()
	m.encoder.Close()
}

func bytesToInt16LE(b []byte) []int16 {
	ints := make([]int16, len(b)/2)
	for i := 0; i < len(ints); i++ {
		ints[i] = int16(b[i*2]) | int16(b[i*2+1])<<8
	}
	return ints
}

type readCanceller struct {
	r      io.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

func newReadCanceller(r io.Reader) *readCanceller {
	ctx, cancel := context.WithCancel(context.Background())
	return &readCanceller{
		r:      r,
		ctx:    ctx,
		cancel: cancel,
	}
}
func (rc *readCanceller) Read(p []byte) (n int, err error) {
	select {
	case <-rc.ctx.Done():
		return 0, io.EOF
	default:
		return rc.r.Read(p)
	}
}

func (rc *readCanceller) Close() {
	rc.cancel()
}
