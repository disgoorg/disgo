package voice

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type ShouldReceiveUserFunc func(userID snowflake.ID) bool

type AudioReceiveSystem interface {
	Open()
	CleanupUser(userID snowflake.ID)
	Close()
}

func NewAudioReceiveSystem(opusFrameReceiver OpusFrameReceiver, connection Connection) AudioReceiveSystem {
	return &defaultAudioReceiveSystem{
		logger:            log.Default(),
		opusFrameReceiver: opusFrameReceiver,
		connection:        connection,
	}
}

type defaultAudioReceiveSystem struct {
	cancelFunc context.CancelFunc

	logger            log.Logger
	opusFrameReceiver OpusFrameReceiver
	connection        Connection
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
	packet, err := s.connection.UDP().ReadPacket()
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

type OpusFrameReceiver interface {
	ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error
	CleanupUser(userID snowflake.ID)
	Close()
}

func NewOpusStreamReceiver(w io.Writer, receiveUserFunc ShouldReceiveUserFunc) OpusFrameReceiver {
	return &opusStreamReceiver{
		w:               w,
		receiveUserFunc: receiveUserFunc,
	}
}

type opusStreamReceiver struct {
	w               io.Writer
	receiveUserFunc ShouldReceiveUserFunc
}

func (r *opusStreamReceiver) ReceiveOpusFrame(userID snowflake.ID, packet *Packet) error {
	if r.receiveUserFunc == nil || !r.receiveUserFunc(userID) {
		return nil
	}
	if err := binary.Write(r.w, binary.LittleEndian, uint32(len(packet.Opus))); err != nil {
		return err
	}
	if _, err := r.w.Write(packet.Opus); err != nil {
		return err
	}
	return nil
}

func (*opusStreamReceiver) CleanupUser(_ snowflake.ID) {}
func (*opusStreamReceiver) Close()                     {}
