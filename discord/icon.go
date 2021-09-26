package discord

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/DisgoOrg/disgo/json"
)

type IconType string

const (
	IconTypeJPEG    IconType = "image/jpeg"
	IconTypePNG     IconType = "image/png"
	IconTypeWEBP    IconType = "image/webp"
	IconTypeGIF     IconType = "image/gif"
	IconTypeUnknown          = IconTypeJPEG
)

func (t IconType) GetMIME() string {
	return string(t)
}

func (t IconType) GetHeader() string {
	return "data:" + string(t) + ";base64"
}

var _ json.Marshaler = (*Icon)(nil)
var _ fmt.Stringer = (*Icon)(nil)

func NewIcon(iconType IconType, reader io.Reader) (Icon, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Icon{}, err
	}

	return NewIconRaw(iconType, data), nil
}

func NewIconRaw(iconType IconType, data []byte) Icon {
	return Icon{Type: iconType, Data: base64.StdEncoding.EncodeToString(data)}
}

type Icon struct {
	Type IconType
	Data string
}

func (i Icon) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i Icon) String() string {
	if i.Data == "" {
		return ""
	}
	return i.Type.GetHeader() + "," + i.Data
}
