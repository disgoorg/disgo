package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/disgoorg/disgo/audio/mp3"
	"github.com/disgoorg/disgo/audio/opus"
)

// NewMP3PCMFrameProvider returns a PCMFrameProvider that reads mp3 and converts it into pcm frames.
// Write the Mp3 data to the returned mp3Writer.
func NewMP3PCMFrameProvider(decoder *mp3.Decoder) (PCMFrameProvider, io.Writer, error) {
	return NewCustomMP3PCMFrameProvider(decoder, 48000, 2)
}

func NewCustomMP3PCMFrameProvider(decoder *mp3.Decoder, rate int, channels int) (PCMFrameProvider, io.Writer, error) {
	if decoder == nil {
		var err error
		decoder, err = mp3.CreateDecoder()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create mp3 decoder: %w", err)
		}
	}

	if err := decoder.Param(mp3.ForceRate, rate, float64(rate)); err != nil {
		return nil, nil, fmt.Errorf("failed to set param: %w", err)
	}

	if err := decoder.OpenFeed(); err != nil {
		return nil, nil, fmt.Errorf("failed to open feed for mp3 decoder: %w", err)
	}

	writeFunc := mp3Writer(func(p []byte) (int, error) {
		return decoder.Write(p)
	})

	return &mp3PCMFrameProvider{
		decoder:     decoder,
		bytePCMBuff: make([]byte, opus.GetOutputBuffSize(rate, channels)*2),
		pcmBuff:     make([]int16, opus.GetOutputBuffSize(rate, channels)),
	}, writeFunc, nil
}

type mp3PCMFrameProvider struct {
	decoder     *mp3.Decoder
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *mp3PCMFrameProvider) ProvidePCMFrame() ([]int16, error) {
	_, err := p.decoder.Read(p.bytePCMBuff)
	if err != nil {
		return nil, err
	}

	if err = binary.Read(bytes.NewReader(p.bytePCMBuff), binary.LittleEndian, p.pcmBuff); err != nil {
		return nil, err
	}
	return p.pcmBuff, nil
}

func (p *mp3PCMFrameProvider) Close() {
	_ = p.decoder.Close()
}

type mp3Writer func(p []byte) (int, error)

func (w mp3Writer) Write(p []byte) (int, error) {
	return w(p)
}
