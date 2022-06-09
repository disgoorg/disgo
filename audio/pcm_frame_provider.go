package audio

import (
	"bytes"
	"encoding/binary"
	"io"
)

// PCMFrameProvider is an interface for providing PCM frames.
type PCMFrameProvider interface {
	// ProvidePCMFrame is called to get a PCM frame.
	ProvidePCMFrame() []int16

	// Close is called when the provider is no longer needed. It should close any open resources.
	Close()
}

// NewPCMStreamProvider creates a new PCMFrameProvider which reads PCM frames from the given io.Reader.
func NewPCMStreamProvider(r io.Reader) PCMFrameProvider {
	return &pcmStreamProvider{
		r: r,
	}
}

type pcmStreamProvider struct {
	r           io.Reader
	bytePCMBuff [1920 * 4]byte
	pcmBuff     [960 * 2]int16
}

func (p *pcmStreamProvider) ProvidePCMFrame() []int16 {
	_, err := p.r.Read(p.bytePCMBuff[:])
	if err != nil {
		if err != io.EOF {
			panic("ProvidePCMFrame: " + err.Error())
		}
		return nil
	}

	r := bytes.NewReader(p.bytePCMBuff[:])
	if err = binary.Read(r, binary.LittleEndian, p.pcmBuff[:]); err != nil {
		if err != io.EOF {
			panic("ProvidePCMFrame: " + err.Error())
		}
		return nil
	}
	return p.pcmBuff[:]
}

func (*pcmStreamProvider) Close() {}
