package audio

import (
	"fmt"
	"io"

	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/voice"
)

// NewPCMOpusProvider creates a new voice.OpusFrameProvider which gets PCM frames from the given PCMFrameProvider and encodes the PCM frames into Opus frames.
// You can pass your own *opus.Encoder or nil to use the default Opus encoder(48000hz sample rate, 2 channels, opus.ApplicationAudio & 64kbps bitrate).
func NewPCMOpusProvider(encoder *opus.Encoder, pcmProvider PCMFrameProvider) (voice.OpusFrameProvider, error) {
	if encoder == nil {
		var err error
		if encoder, err = opus.NewEncoder(48000, 2, opus.ApplicationAudio); err != nil {
			return nil, fmt.Errorf("failed to create opus encoder: %w", err)
		}
		if err = encoder.Ctl(opus.SetBitrate(64000)); err != nil {
			return nil, fmt.Errorf("failed to set opus bitrate: %w", err)
		}

	}
	return &pcmOpusProvider{
		encoder:     encoder,
		pcmProvider: pcmProvider,
		opusBuff:    make([]byte, 2048),
	}, nil
}

type pcmOpusProvider struct {
	encoder     *opus.Encoder
	pcmProvider PCMFrameProvider
	opusBuff    []byte
}

func (p *pcmOpusProvider) ProvideOpusFrame() ([]byte, error) {
	pcm, err := p.pcmProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}
	if len(pcm) == 0 {
		return nil, io.EOF
	}

	n, err := p.encoder.Encode(pcm, p.opusBuff)
	if err != nil {
		return nil, err
	}
	return p.opusBuff[:n], nil
}

func (p *pcmOpusProvider) Close() {
	p.encoder.Destroy()
	p.pcmProvider.Close()
}
