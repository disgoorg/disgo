package voice

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/disgoorg/log"
)

var SilenceFrames = []byte{0xF8, 0xFF, 0xFE}

type AudioSendHandler interface {
	ProvideOpus() []byte
}

func NewSendSystem(sendHandler AudioSendHandler, connection *Connection) SendSystem {
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
	closed      bool
	sendHandler AudioSendHandler
	connection  *Connection

	silentFrames      int
	sentSpeakingStop  bool
	sentSpeakingStart bool
}

func (s *defaultSendSystem) Start() {
	go func() {
		defer s.logger.Debug("closing send system")
		lastFrameSent := time.Now().UnixMilli()
		for !s.closed {
			s.send()
			sleepTime := time.Duration(20 - (time.Now().UnixMilli() - lastFrameSent))
			if sleepTime > 0 {
				time.Sleep(sleepTime * time.Millisecond)
			}
			if time.Now().UnixMilli() < lastFrameSent+60 {
				lastFrameSent += 20
			} else {
				lastFrameSent = time.Now().UnixMilli()
			}
		}
	}()
}

func (s *defaultSendSystem) send() {
	opus := s.sendHandler.ProvideOpus()
	if len(opus) == 0 {
		if s.silentFrames > 0 {
			if _, err := s.connection.UDP().Write(SilenceFrames); err != nil {
				s.logger.Errorf("failed to send silence frames: %s", err)
			}
			s.silentFrames--
		} else if !s.sentSpeakingStop {
			go func() {
				if err := s.connection.Speaking(0); err != nil {
					s.logger.Errorf("failed to send speaking stop: %s", err)
				}
			}()
			s.sentSpeakingStop = true
			s.sentSpeakingStart = false
		}
		return
	}

	if !s.sentSpeakingStart {
		go func() {
			if err := s.connection.Speaking(SpeakingFlagMicrophone); err != nil {
				s.logger.Errorf("failed to send speaking start: %s", err)
			}
		}()
		s.sentSpeakingStart = true
		s.sentSpeakingStop = false
		s.silentFrames = 5
	}

	if _, err := s.connection.UDP().Write(opus); err != nil {
		s.logger.Errorf("failed to send opus data: %s", err)
	}
}

func (s *defaultSendSystem) Stop() {
	s.closed = true
}

var _ AudioSendHandler = (*OpusSendHandler)(nil)

func NewOpusSendHandler(r io.Reader) *OpusSendHandler {
	return &OpusSendHandler{
		r:    r,
		buff: make([]byte, 4000),
	}
}

type OpusSendHandler struct {
	r       io.Reader
	lenBuff [4]byte
	buff    []byte
}

func (h *OpusSendHandler) ProvideOpus() []byte {
	_, err := h.r.Read(h.lenBuff[:])
	if err == io.EOF {
		return nil
	}

	buff := make([]byte, int64(binary.LittleEndian.Uint32(h.lenBuff[:])))
	var n int
	n, err = h.r.Read(buff)
	if err != nil {
		return nil
	}
	return buff[:n]
}
