package audio

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/audio/opus"
)

// PCMFrameProvider is an interface for providing PCM frames.
type PCMFrameProvider interface {
	// ProvidePCMFrame is called to get a PCM frame.
	ProvidePCMFrame() ([]int16, error)

	// Close is called when the provider is no longer needed. It should close any open resources.
	Close()
}

// NewPCMStreamProvider creates a new PCMFrameProvider which reads PCM frames from the given io.Reader.
func NewPCMStreamProvider(r io.Reader) PCMFrameProvider {
	return NewCustomPCMStreamProvider(r, 48000, 2)
}

func NewCustomPCMStreamProvider(r io.Reader, rate int, channels int) PCMFrameProvider {
	return &pcmStreamProvider{
		r:           r,
		bytePCMBuff: make([]byte, opus.GetOutputBuffSize(rate, channels)*2),
		pcmBuff:     make([]int16, opus.GetOutputBuffSize(rate, channels)),
	}
}

type pcmStreamProvider struct {
	r           io.Reader
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *pcmStreamProvider) ProvidePCMFrame() ([]int16, error) {
	_, err := p.r.Read(p.bytePCMBuff)
	if err != nil {
		return nil, err
	}

	if err = binary.Read(bytes.NewReader(p.bytePCMBuff), binary.LittleEndian, p.pcmBuff); err != nil {
		return nil, err
	}
	return p.pcmBuff, nil
}

func (*pcmStreamProvider) Close() {}
