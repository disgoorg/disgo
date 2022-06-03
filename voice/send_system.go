package voice

import (
	"time"

	"github.com/disgoorg/log"
)

type SendHandler interface {
	ProvideOpus() ([]byte, error)
}

func NewSendSystem(sendHandler SendHandler, connection *Connection) SendSystem {
	return &defaultSendSystem{
		logger:       log.Default(),
		sendHandler:  sendHandler,
		connection:   connection,
		silentFrames: 5,
	}
}

type SendSystem interface {
	Start()
	Stop()
}

type defaultSendSystem struct {
	logger      log.Logger
	ticker      *time.Ticker
	sendHandler SendHandler
	connection  *Connection

	silentFrames int
	speaking     bool
}

func (s *defaultSendSystem) Start() {
	go func() {
		defer s.logger.Debug("closing send system")
		s.ticker = time.NewTicker(time.Millisecond * 20)
		defer s.ticker.Stop()
		for range s.ticker.C {
			s.send()
		}
	}()
}

func (s *defaultSendSystem) send() {
	opus, err := s.sendHandler.ProvideOpus()
	if err != nil {
		s.logger.Errorf("failed to provide opus data: %s", err)
		return
	}

	s.intercept(opus)
	if opus == nil {
		if s.silentFrames > 0 {
			println("sent silent frame")
			if _, err = s.connection.UDPConn().Write([]byte{0xF8, 0xFF, 0xFE}); err != nil {
				s.logger.Errorf("failed to send opus data: %s", err)
			}
		}

		return
	}

	println("sent opus frame")
	if _, err = s.connection.UDPConn().Write(opus); err != nil {
		s.logger.Errorf("failed to send opus data: %s", err)
	}
}

func (s *defaultSendSystem) intercept(opus []byte) {
	if opus != nil && !s.speaking {
		println("sent speaking start")

		s.silentFrames = 5
		s.speaking = true

		if err := s.connection.Speaking(SpeakingFlagMicrophone); err != nil {
			s.logger.Error("error sending speaking: ", err)
		}
	} else if opus == nil && s.speaking {
		s.silentFrames--
		if s.silentFrames == 0 {
			return
		}

		s.speaking = false
		if err := s.connection.Speaking(0); err != nil {
			s.logger.Error("error sending speaking: ", err)
		}
	}
}

func (s *defaultSendSystem) Stop() {
	s.ticker.Stop()
}
