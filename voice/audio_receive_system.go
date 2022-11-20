package voice

import (
	"context"
	"net"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type (
	AudioReceiveSystemCreateFunc func(logger log.Logger, receiver OpusFrameReceiver, connection Conn) AudioReceiveSystem
	UserFilterFunc               func(userID snowflake.ID) bool

	AudioReceiveSystem interface {
		Open()
		CleanupUser(userID snowflake.ID)
		Close()
	}

	OpusFrameReceiver interface {
		ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error
		CleanupUser(userID snowflake.ID)
		Close()
	}
)

func NewAudioReceiveSystem(logger log.Logger, opusFrameReceiver OpusFrameReceiver, connection Conn) AudioReceiveSystem {
	return &defaultAudioReceiveSystem{
		logger:            logger,
		opusFrameReceiver: opusFrameReceiver,
		connection:        connection,
	}
}

type defaultAudioReceiveSystem struct {
	cancelFunc context.CancelFunc

	logger            log.Logger
	opusFrameReceiver OpusFrameReceiver
	connection        Conn
}

func (s *defaultAudioReceiveSystem) Open() {
	defer s.logger.Debugf("closing receive system")
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

func (s *defaultAudioReceiveSystem) CleanupUser(userID snowflake.ID) {
	s.opusFrameReceiver.CleanupUser(userID)
}

func (s *defaultAudioReceiveSystem) receive() {
	packet, err := s.connection.UDPConn().ReadPacket()
	if err == net.ErrClosed {
		s.Close()
		return
	}
	if err != nil {
		s.logger.Errorf("error while reading packet: %s", err)
		return
	}
	if s.opusFrameReceiver != nil {
		if err = s.opusFrameReceiver.ReceiveOpusFrame(s.connection.UserIDBySSRC(packet.SSRC), packet); err != nil {
			s.logger.Errorf("error while receiving opus frame: %s", err)
		}
	}

}

func (s *defaultAudioReceiveSystem) Close() {
	s.cancelFunc()
	s.opusFrameReceiver.Close()
}
