package audio

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/audio/mp3"
)

func NewMP3PCMFrameProvider(decoder *mp3.Decoder) (PCMFrameProvider, Mp3Writer) {
	if decoder == nil {
		var err error
		decoder, err = mp3.CreateDecoder()
		if err != nil {
			panic("NewMP3PCMFrameProvider: " + err.Error())
		}
	}

	if err := decoder.OpenFeed(); err != nil {
		panic("NewMP3PCMFrameProvider: " + err.Error())
	}

	writeFunc := func(p []byte) (int, error) {
		return decoder.Write(p)
	}

	return &mp3PCMFrameProvider{
		decoder: decoder,
	}, writeFunc
}

type mp3PCMFrameProvider struct {
	decoder     *mp3.Decoder
	bytePCMBuff [1920 * 2]byte
	pcmBuff     [960 * 2]int16
}

func (p *mp3PCMFrameProvider) ProvidePCMFrame() []int16 {
	_, err := p.decoder.Read(p.bytePCMBuff[:])
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

func (p *mp3PCMFrameProvider) Close() {
	_ = p.decoder.Close()
}

type Mp3Writer func(p []byte) (int, error)

func (w Mp3Writer) Write(p []byte) (int, error) {
	return w(p)
}
