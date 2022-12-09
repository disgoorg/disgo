package voice

import (
	"context"
	"net"

	"github.com/disgoorg/disgo/voice/voiceudp"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type (
	AudioReceiverCreateFunc func(logger log.Logger, receiver OpusFrameReceiver, connection Conn) AudioReceiver
	UserFilterFunc          func(userID snowflake.ID) bool

	AudioReceiver interface {
		Open()
		CleanupUser(userID snowflake.ID)
		Close()
	}

	OpusFrameReceiver interface {
		ReceiveOpusFrame(userID snowflake.ID, packet *voiceudp.Packet) error
		CleanupUser(userID snowflake.ID)
		Close()
	}
)

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
	packet, err := s.conn.Conn().ReadPacket()
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
