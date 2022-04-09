package route

import (
	"strings"
)

// MajorParameters is a list of url parameters which decide in which bucket a route belongs (https://discord.com/developers/docs/topics/rate-limits#rate-limits)
const MajorParameters = "guild.id:channel.id:webhook.id:interaction.token"

func countURLParams(url string) int {
	return strings.Count(url, "{")
}

// Method is an HTTP request Method
type Method string

// HTTP Methods used by Discord
const (
	DELETE Method = "DELETE"
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
)

// String returns the string representation of the Method
func (m Method) String() string {
	return string(m)
}

// QueryValues is used to supply query param value pairs to Route.Compile
type QueryValues map[string]any

// ImageFormat is the type of image on Discord's CDN (https://discord.com/developers/docs/reference#image-formatting-image-formats)
type ImageFormat string

// The available ImageFormat(s)
const (
	PNG    ImageFormat = "png"
	JPEG   ImageFormat = "jpg"
	WebP   ImageFormat = "webp"
	GIF    ImageFormat = "gif"
	Lottie ImageFormat = "json"
	BLANK  ImageFormat = ""
)

func (f ImageFormat) String() string {
	return string(f)
}

func (f ImageFormat) CanBeAnimated() bool {
	switch f {
	case WebP, GIF:
		return true
	default:
		return false
	}
}
