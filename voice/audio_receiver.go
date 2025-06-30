package voice

import (
	"context"
	"errors"
	"log/slog"
	"net"

	"github.com/disgoorg/snowflake/v2"
)

type (
	// AudioReceiverCreateFunc is used to create a new AudioReceiver reading audio from the given Conn.
	AudioReceiverCreateFunc func(logger *slog.Logger, receiver OpusFrameReceiver, connection Conn) AudioReceiver

	// UserFilterFunc is used as a filter for which users to receive audio from.
	UserFilterFunc func(userID snowflake.ID) bool

	// AudioReceiver is used to receive audio from a voice connection and pass it to an OpusFrameReceiver.
	AudioReceiver interface {
		// Open starts receiving audio from the voice connection.
		Open()

		// CleanupUser cleans up any audio resources for the given user.
		CleanupUser(userID snowflake.ID)

		// Close stops receiving audio from the voice connection.
		Close()
	}

	// OpusFrameReceiver is an interface used to receive opus frames from an AudioReceiver.
	OpusFrameReceiver interface {
		// ReceiveOpusFrame receives an opus frame.
		ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error

		// CleanupUser cleans up any audio resources for the given user.
		CleanupUser(userID snowflake.ID)

		// Close stops receiving audio from the voice connection.
		Close()
	}
)

// NewAudioReceiver creates a new AudioReceiver reading audio to the given OpusFrameReceiver from the given Conn.
func NewAudioReceiver(logger *slog.Logger, opusReceiver OpusFrameReceiver, conn Conn) AudioReceiver {
	return &defaultAudioReceiver{
		logger:       logger,
		opusReceiver: opusReceiver,
		conn:         conn,
	}
}

type defaultAudioReceiver struct {
	logger       *slog.Logger
	cancelFunc   context.CancelFunc
	opusReceiver OpusFrameReceiver
	conn         Conn
}

func (s *defaultAudioReceiver) Open() {
	go s.open()
}

func (s *defaultAudioReceiver) open() {
	defer s.logger.Debug("closing audio receiver")
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			s.receive()
		}
	}
}

func (s *defaultAudioReceiver) CleanupUser(userID snowflake.ID) {
	s.opusReceiver.CleanupUser(userID)
}

func (s *defaultAudioReceiver) receive() {
	packet, err := s.conn.UDP().ReadPacket()
	if errors.Is(err, net.ErrClosed) {
		s.Close()
		return
	}
	if err != nil {
		s.logger.Error("error while reading packet", slog.Any("err", err))
		return
	}
	if s.opusReceiver != nil {
		if err = s.opusReceiver.ReceiveOpusFrame(s.conn.UserIDBySSRC(packet.SSRC), packet); err != nil {
			s.logger.Error("error while receiving opus frame", slog.Any("err", err))
		}
	}

}

func (s *defaultAudioReceiver) Close() {
	s.cancelFunc()
	s.opusReceiver.Close()
}
