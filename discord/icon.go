package discord

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/disgoorg/json/v2"
)

// IconType represents the mimetype of icons discord supports.
// https://docs.discord.com/developers/reference#image-formatting
type IconType string

// IconTypes is a list of supported icon types by discord.
var IconTypes = []IconType{
	IconTypeJPEG,
	IconTypePNG,
	IconTypeWEBP,
	IconTypeAVIF,
	IconTypeGIF,
}

const (
	IconTypeJPEG    IconType = "image/jpeg"
	IconTypePNG     IconType = "image/png"
	IconTypeWEBP    IconType = "image/webp"
	IconTypeAVIF    IconType = "image/avif"
	IconTypeGIF     IconType = "image/gif"
	IconTypeUnknown          = IconTypeJPEG
)

func (t IconType) MIME() string {
	return string(t)
}

func (t IconType) Header() string {
	return "data:" + string(t) + ";base64"
}

// NewIcon reads the image data from the provided reader and creates a new Icon with the specified mime type.
// The image data will be base64 encoded.
func NewIcon(iconType IconType, r io.Reader) (*Icon, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewIconRaw(iconType, data), nil
}

// NewIconRaw creates a new Icon with the specified mime type and raw image data.
// The image data will be base64 encoded.
func NewIconRaw(iconType IconType, src []byte) *Icon {
	data := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(data, src)
	return &Icon{
		Type: iconType,
		Data: data,
	}
}

// ParseIcon reads the image data from the provided reader, detects the mime type, and creates a new Icon.
// The image data will be base64 encoded.
// http.DetectContentType does not seem to support IconTypeAVIF.
func ParseIcon(r io.Reader) (*Icon, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseIconRaw(data)
}

// ParseIconRaw detects the mime type of the provided image data and creates a new Icon.
// The image data will be base64 encoded.
// http.DetectContentType does not seem to support IconTypeAVIF.
func ParseIconRaw(src []byte) (*Icon, error) {
	mime := IconType(http.DetectContentType(src))
	if slices.Index(IconTypes, mime) == -1 {
		return nil, fmt.Errorf("unsupported icon type: %s", mime)
	}
	return NewIconRaw(mime, src), nil
}

var (
	_ json.Marshaler = (*Icon)(nil)
	_ fmt.Stringer   = (*Icon)(nil)
)

// Icon represents a base64 encoded image with its mimetype, used for icons in discord.
// Use NewIcon, NewIconRaw, ParseIcon, or ParseIconRaw to create a new Icon.
type Icon struct {
	// Type is the mimetype of the image.
	Type IconType
	// Data is the base64 encoded image data.
	Data []byte
}

func (i Icon) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// String formats the icon into the correct format for discord.
// https://docs.discord.com/developers/reference#image-data
func (i Icon) String() string {
	if len(i.Data) == 0 {
		return ""
	}
	return i.Type.Header() + "," + string(i.Data)
}
