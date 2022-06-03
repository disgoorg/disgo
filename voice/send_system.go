package voice

import (
	"time"

	"github.com/disgoorg/log"
)

var SilenceFrames = []byte{0xF8, 0xFF, 0xFE}

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

	silentFrames      int
	sentSpeakingStop  bool
	sentSpeakingStart bool
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
	if len(opus) == 0 {
		s.logger.Debug("opus data is empty")
		if s.silentFrames > 0 {
			println("sending silence")
			if _, err = s.connection.UDPConn().Write(SilenceFrames); err != nil {
				s.logger.Errorf("failed to send silence frames: %s", err)
			}
			s.silentFrames--
		} else if !s.sentSpeakingStop {
			println("sending speaking stop")
			if err = s.connection.Speaking(0); err != nil {
				s.logger.Errorf("failed to send speaking stop: %s", err)
			}
			s.sentSpeakingStop = true
			s.sentSpeakingStart = false
		}
		return
	}

	if !s.sentSpeakingStart {
		println("sending speaking start")
		if err = s.connection.Speaking(SpeakingFlagMicrophone | SpeakingFlagPriority); err != nil {
			s.logger.Errorf("failed to send speaking start: %s", err)
		}
		s.sentSpeakingStart = true
		s.sentSpeakingStop = false
		s.silentFrames = 5
	}

	println("sending opus")
	if _, err = s.connection.UDPConn().Write(opus); err != nil {
		s.logger.Errorf("failed to send opus data: %s", err)
	}
}

func (s *defaultSendSystem) Stop() {
	s.ticker.Stop()
}
