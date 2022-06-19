package voice

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/disgoorg/log"
)

var (
	SilenceAudioFrames       = []byte{0xF8, 0xFF, 0xFE}
	OpusFrameSize      int64 = 20
)

func NewAudioSendSystem(logger log.Logger, opusProvider OpusFrameProvider, connection Connection) AudioSendSystem {
	return &defaultAudioSendSystem{
		logger:       logger,
		opusProvider: opusProvider,
		connection:   connection,
		silentFrames: 5,
	}
}

type AudioSendSystem interface {
	Open()
	Close()
}

type defaultAudioSendSystem struct {
	logger       log.Logger
	cancelFunc   context.CancelFunc
	opusProvider OpusFrameProvider
	connection   Connection

	silentFrames      int
	sentSpeakingStop  bool
	sentSpeakingStart bool
}

func (s *defaultAudioSendSystem) Open() {
	defer s.logger.Debug("closing send system")
	lastFrameSent := time.Now().UnixMilli()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		default:
			s.send()
			sleepTime := time.Duration(OpusFrameSize - (time.Now().UnixMilli() - lastFrameSent))
			if sleepTime > 0 {
				time.Sleep(sleepTime * time.Millisecond)
			}
			if time.Now().UnixMilli() < lastFrameSent+OpusFrameSize*3 {
				lastFrameSent += OpusFrameSize
			} else {
				lastFrameSent = time.Now().UnixMilli()
			}
		}
	}
}

func (s *defaultAudioSendSystem) send() {
	if s.opusProvider == nil {
		return
	}
	opus, err := s.opusProvider.ProvideOpusFrame()
	if err == net.ErrClosed {
		s.Close()
		return
	}
	if err != nil && err != io.EOF {
		s.logger.Errorf("error while reading opus frame: %s", err)
		return
	}
	if len(opus) == 0 {
		if s.silentFrames > 0 {
			if _, err = s.connection.UDP().Write(SilenceAudioFrames); err != nil {
				s.logger.Errorf("failed to send silence frames: %s", err)
			}
			s.silentFrames--
		} else if !s.sentSpeakingStop {
			if err = s.connection.Speaking(0); err != nil {
				s.logger.Errorf("failed to send speaking stop: %s", err)
			}
			s.sentSpeakingStop = true
			s.sentSpeakingStart = false
		}
		return
	}

	if !s.sentSpeakingStart {
		if err = s.connection.Speaking(SpeakingFlagMicrophone); err != nil {
			s.logger.Errorf("failed to send speaking start: %s", err)
		}
		s.sentSpeakingStart = true
		s.sentSpeakingStop = false
		s.silentFrames = 5
	}

	if _, err = s.connection.UDP().Write(opus); err != nil {
		s.logger.Errorf("failed to send audio data: %s", err)
	}
}

func (s *defaultAudioSendSystem) Close() {
	s.cancelFunc()
}

type OpusFrameProvider interface {
	ProvideOpusFrame() ([]byte, error)
	Close()
}

var _ OpusFrameProvider = (*OpusStreamProvider)(nil)

func NewOpusStreamProvider(r io.Reader) *OpusStreamProvider {
	return &OpusStreamProvider{
		r:    r,
		buff: make([]byte, 4000),
	}
}

type OpusStreamProvider struct {
	r       io.Reader
	lenBuff [4]byte
	buff    []byte
}

func (h *OpusStreamProvider) ProvideOpusFrame() ([]byte, error) {
	_, err := h.r.Read(h.lenBuff[:])
	if err != nil {
		return nil, err
	}

	buff := make([]byte, int64(binary.LittleEndian.Uint32(h.lenBuff[:])))
	var n int
	n, err = h.r.Read(buff)
	if err != nil {
		return nil, err
	}
	return buff[:n], nil
}

func (*OpusStreamProvider) Close() {}
