package route

import "strings"

const MajorParameters = "guild.id:channel.id:webhook.id:interaction.token"

func countURLParams(url string) int {
	paramCount := strings.Count(url, "{")
	return paramCount
}

// Method is a HTTP request Method
type Method string

// HTTP Methods used by Discord
const (
	DELETE Method = "DELETE"
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
)

func (m Method) String() string {
	return string(m)
}

// QueryValues is used to supply query param value pairs to Route.Compile
type QueryValues map[string]interface{}

// FileExtension is the type of image on Discord's CDN
type FileExtension string

// The available FileExtension(s)
const (
	PNG    FileExtension = "png"
	JPEG   FileExtension = "jpg"
	WebP   FileExtension = "webp"
	GIF    FileExtension = "gif"
	Lottie FileExtension = "json"
	BLANK  FileExtension = ""
)

func (f FileExtension) String() string {
	return string(f)
}
