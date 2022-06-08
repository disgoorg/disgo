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
		r:           r,
		bytePCMBuff: make([]byte, 960*4),
		pcmBuff:     make([]int16, 960),
	}
}

type pcmStreamProvider struct {
	r           io.Reader
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *pcmStreamProvider) ProvidePCMFrame() []int16 {
	n, err := p.r.Read(p.bytePCMBuff)
	if err != nil {
		return nil
	}

	r := bytes.NewReader(p.bytePCMBuff[:n])
	if err = binary.Read(r, binary.LittleEndian, p.pcmBuff); err != nil {
		return nil
	}

	return p.pcmBuff[:n/2]
}

func (*pcmStreamProvider) Close() {}
