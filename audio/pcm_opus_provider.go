package audio

import (
	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/voice"
)

func NewPCMOpusProvider(encoder *opus.Encoder, pcmProvider PCMFrameProvider) voice.OpusFrameProvider {
	if encoder == nil {
		encoder, _ = opus.NewEncoder(48000, 2, opus.ApplicationAudio)
	}
	return &pcmOpusProvider{
		encoder:     encoder,
		pcmProvider: pcmProvider,
		opusBuff:    make([]byte, 4000),
	}
}

type pcmOpusProvider struct {
	encoder     *opus.Encoder
	pcmProvider PCMFrameProvider
	opusBuff    []byte
}

func (p *pcmOpusProvider) ProvideOpusFrame() []byte {
	pcm := p.pcmProvider.ProvidePCMFrame()
	if len(pcm) == 0 {
		return nil
	}

	n, err := p.encoder.Encode(pcm, p.opusBuff)
	if err != nil {
		return nil
	}
	return p.opusBuff[:n]
}

func (p *pcmOpusProvider) Close() {
	p.encoder.Destroy()
	p.pcmProvider.Close()
}
