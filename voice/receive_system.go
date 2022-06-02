package voice

import (
	"context"

	"github.com/disgoorg/log"
)

func NewReceiveSystem(receiveHandler ReceiveHandler, connection Connection) ReceiveSystem {
	return &defaultReceiveSystem{
		receiveHandler: receiveHandler,
		connection:     connection,
		opusBuffer:     make([]byte, 1400),
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
	connection     Connection

	opusBuffer []byte
}

func (s *defaultReceiveSystem) Start() {
	defer s.logger.Debugf("closing receive system")
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	defer s.cancelFunc()
	go func() {
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
	i, err := s.connection.UDPConn().Read(s.opusBuffer)
	if err != nil {
		return
	}

	if s.receiveHandler == nil || !s.receiveHandler.CanReceive() {
		return
	}

	s.receiveHandler.HandleOpus(s.opusBuffer[:i])
}

func (s *defaultReceiveSystem) Stop() {
	s.cancelFunc()
}
