package voice

import (
	"time"

	"github.com/disgoorg/log"
)

func NewSendSystem(sendHandler SendHandler, connection Connection, tickInterval time.Duration) SendSystem {
	return &defaultSendSystem{
		tickInterval: tickInterval,
		sendHandler:  sendHandler,
		connection:   connection,
	}
}

type SendSystem interface {
	Start()
	Stop()
}

type defaultSendSystem struct {
	logger       log.Logger
	tickInterval time.Duration
	ticker       *time.Ticker
	sendHandler  SendHandler
	connection   Connection

	silentFrames int
}

func (s *defaultSendSystem) Start() {
	defer s.logger.Debug("closing send system")
	s.ticker = time.NewTicker(s.tickInterval)
	defer s.ticker.Stop()
	go func() {
		for range s.ticker.C {
			s.send()
		}
	}()
}

func (s *defaultSendSystem) send() {
	if !s.sendHandler.CanProvide() && s.silentFrames > 5 {

		return
	}
	opus, err := s.sendHandler.ProvideOpus(s.tickInterval)
	if err != nil {
		s.logger.Errorf("failed to provide opus data: %s", err)
		return
	}
	if _, err = s.connection.UDPConn().Write(s.tickInterval).Write(opus); err != nil {
		s.logger.Errorf("failed to send opus data: %s", err)
	}
}

func (s *defaultSendSystem) Stop() {
	s.ticker.Stop()
}
