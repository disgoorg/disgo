package audio

import (
	"bytes"
	"encoding/binary"
	"io"
)

type PCMFrameProvider interface {
	ProvidePCMFrame() []int16
	Close()
}

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
