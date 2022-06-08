package opus

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/opus/opus"
	"github.com/disgoorg/disgo/voice"
)

var _ voice.AudioSendHandler = (*PCMSendHandler)(nil)

func NewPCMSendHandler(encoder *opus.Encoder, pcmProvider PCMProvider) *PCMSendHandler {
	return &PCMSendHandler{
		encoder:     encoder,
		pcmProvider: pcmProvider,
		opusBuff:    make([]byte, 4000),
	}
}

type PCMSendHandler struct {
	encoder     *opus.Encoder
	pcmProvider PCMProvider
	opusBuff    []byte
}

func (h *PCMSendHandler) ProvideOpus() []byte {
	pcm := h.pcmProvider.ProvidePCM()
	if len(pcm) == 0 {
		return nil
	}

	n, err := h.encoder.Encode(pcm, h.opusBuff)
	if err != nil {
		return nil
	}
	return h.opusBuff[:n]
}

var _ PCMProvider = (*PCMReaderProvider)(nil)

type PCMProvider interface {
	ProvidePCM() []int16
}

func NewPCMReaderProvider(r io.Reader) *PCMReaderProvider {
	return &PCMReaderProvider{
		r:           r,
		bytePCMBuff: make([]byte, 960*4),
		pcmBuff:     make([]int16, 960),
	}
}

type PCMReaderProvider struct {
	r           io.Reader
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *PCMReaderProvider) ProvidePCM() []int16 {
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
