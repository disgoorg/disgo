package audio

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/audio/mp3"
)

func NewMP3PCMFrameProvider(decoder *mp3.Decoder, mp3Provider Mp3FrameProvider) PCMFrameProvider {
	if decoder == nil {
		var err error
		decoder, err = mp3.CreateDecoder()
		if err != nil {
			panic("NewMP3PCMFrameProvider: " + err.Error())
		}
		if err = decoder.Param(mp3.ForceRate, 48000, 48000.0); err != nil {
			panic("NewMP3PCMFrameProvider: " + err.Error())
		}
		if err = decoder.Param(mp3.ParamFlagForceStereo, 2, 2.0); err != nil {
			panic("NewMP3PCMFrameProvider: " + err.Error())
		}
	}
	return &mp3PCMFrameProvider{
		decoder:     decoder,
		mp3Provider: mp3Provider,
	}
}

type mp3PCMFrameProvider struct {
	decoder     *mp3.Decoder
	mp3Provider Mp3FrameProvider
	bytePCMBuff [1920 * 4]byte
	pcmBuff     [960 * 2]int16
}

func (p *mp3PCMFrameProvider) ProvidePCMFrame() []int16 {
	frame := p.mp3Provider.ProvideMp3Frame()
	if len(frame) == 0 {
		return nil
	}

	_, err := p.decoder.Decode(frame, p.bytePCMBuff[:])
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
	p.mp3Provider.Close()
}

type Mp3FrameProvider interface {
	ProvideMp3Frame() []byte
	Close()
}

// NewMp3StreamProvider creates a new Mp3FrameProvider which reads Mp3 frames from the given io.Reader.
func NewMp3StreamProvider(r io.Reader) Mp3FrameProvider {
	return &mp3StreamProvider{
		r: r,
	}
}

type mp3StreamProvider struct {
	r io.Reader
}

func (p *mp3StreamProvider) ProvideMp3Frame() []byte {
	return nil
}

func (p *mp3StreamProvider) Close() {

}
