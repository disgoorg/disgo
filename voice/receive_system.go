package voice

import (
	"context"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

type ReceiveHandler interface {
	HandleOpus(userID snowflake.ID, packet *Packet)
}

func NewReceiveSystem(receiveHandler ReceiveHandler, connection *Connection) ReceiveSystem {
	return &defaultReceiveSystem{
		logger:         log.Default(),
		receiveHandler: receiveHandler,
		connection:     connection,
	}
}

type ReceiveSystem interface {
	Start()
	Stop()
}

type defaultReceiveSystem struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	logger         log.Logger
	receiveHandler ReceiveHandler
	connection     *Connection
}

func (s *defaultReceiveSystem) Start() {
	go func() {
		defer s.logger.Debugf("closing receive system")
		s.ctx, s.cancelFunc = context.WithCancel(context.Background())
		defer s.cancelFunc()
	loop:
		for {
			select {
			case <-s.ctx.Done():
				break loop
			default:
				s.receive()
			}
		}
	}()
}

func (s *defaultReceiveSystem) receive() {
	packet, err := s.connection.UDP().ReadPacket()
	if err != nil {
		s.logger.Errorf("error while reading packet: %s", err)
		return
	}
	if s.receiveHandler != nil {
		s.receiveHandler.HandleOpus(s.connection.UserIDBySSRC(packet.SSRC), packet)
	}
}

func (s *defaultReceiveSystem) Stop() {
	s.cancelFunc()
}
