package voice

import (
	"context"
	"net"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type (
	// AudioReceiverCreateFunc is used to create a new AudioReceiver reading audio from the given Conn.
	AudioReceiverCreateFunc func(logger log.Logger, receiver OpusFrameReceiver, connection Conn) AudioReceiver

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
func NewAudioReceiver(logger log.Logger, opusReceiver OpusFrameReceiver, conn Conn) AudioReceiver {
	return &defaultAudioReceiver{
		logger:       logger,
		opusReceiver: opusReceiver,
		conn:         conn,
	}
}

type defaultAudioReceiver struct {
	logger       log.Logger
	cancelFunc   context.CancelFunc
	opusReceiver OpusFrameReceiver
	conn         Conn
}

func (s *defaultAudioReceiver) Open() {
	go s.open()
}

func (s *defaultAudioReceiver) open() {
	defer s.logger.Debugf("closing audio receiver")
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
	if err == net.ErrClosed {
		s.Close()
		return
	}
	if err != nil {
		s.logger.Errorf("error while reading packet: %s", err)
		return
	}
	if s.opusReceiver != nil {
		if err = s.opusReceiver.ReceiveOpusFrame(s.conn.UserIDBySSRC(packet.SSRC), packet); err != nil {
			s.logger.Errorf("error while receiving opus frame: %s", err)
		}
	}

}

func (s *defaultAudioReceiver) Close() {
	s.cancelFunc()
	s.opusReceiver.Close()
}
