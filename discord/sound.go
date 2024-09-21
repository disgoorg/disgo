package discord

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/disgoorg/json"
)

type SoundType string

const (
	SoundTypeMP3     SoundType = "audio/mpeg"
	SoundTypeOGG     SoundType = "audio/ogg"
	SoundTypeWAV     SoundType = "audio/wav"
	SoundTypeUnknown           = SoundTypeMP3
)

func (t SoundType) MIME() string {
	return string(t)
}

func (t SoundType) Header() string {
	return "data:" + string(t) + ";base64"
}

var _ json.Marshaler = (*Sound)(nil)
var _ fmt.Stringer = (*Sound)(nil)

func NewSound(soundType SoundType, reader io.Reader) (*Sound, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewSoundRaw(soundType, data), nil
}

func NewSoundRaw(soundType SoundType, src []byte) *Sound {
	data := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(data, src)
	return &Sound{Type: soundType, Data: data}
}

type Sound struct {
	Type SoundType
	Data []byte
}

func (s Sound) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s Sound) String() string {
	if len(s.Data) == 0 {
		return ""
	}
	return s.Type.Header() + "," + string(s.Data)
}
