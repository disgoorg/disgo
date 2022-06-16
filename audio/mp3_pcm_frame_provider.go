package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/disgoorg/disgo/audio/mp3"
)

// NewMP3PCMFrameProvider returns a PCMFrameProvider that reads mp3 and converts it into pcm frames.
// Write the Mp3 data to the returned Mp3Writer.
func NewMP3PCMFrameProvider(decoder *mp3.Decoder) (PCMFrameProvider, Mp3Writer, error) {
	if decoder == nil {
		var err error
		decoder, err = mp3.CreateDecoder()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create mp3 decoder: %w", err)
		}

		if err = decoder.Param(mp3.ForceRate, 48000, 48000); err != nil {
			return nil, nil, fmt.Errorf("failed to set param: %w", err)
		}
	}

	if err := decoder.OpenFeed(); err != nil {
		return nil, nil, fmt.Errorf("failed to open feed for mp3 decoder: %w", err)
	}

	writeFunc := func(p []byte) (int, error) {
		return decoder.Write(p)
	}

	return &mp3PCMFrameProvider{
		decoder: decoder,
	}, writeFunc, nil
}

type mp3PCMFrameProvider struct {
	decoder     *mp3.Decoder
	bytePCMBuff [1920 * 2]byte
	pcmBuff     [960 * 2]int16
}

func (p *mp3PCMFrameProvider) ProvidePCMFrame() ([]int16, error) {
	_, err := p.decoder.Read(p.bytePCMBuff[:])
	if err != nil {
		return nil, err
	}

	if err = binary.Read(bytes.NewReader(p.bytePCMBuff[:]), binary.LittleEndian, p.pcmBuff[:]); err != nil {
		return nil, err
	}
	return p.pcmBuff[:], nil
}

func (p *mp3PCMFrameProvider) Close() {
	_ = p.decoder.Close()
}

type Mp3Writer func(p []byte) (int, error)

func (w Mp3Writer) Write(p []byte) (int, error) {
	return w(p)
}
