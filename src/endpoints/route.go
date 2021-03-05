package endpoints

import (
	"fmt"
	"strings"
)

// Discord Route Constants
const (
	APIVersion = "8"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com/"
	API        = Base + "api/v" + APIVersion + "/"
	WS         = "wss://gateway.discord.gg/"
)

// Discord Route Methods
var (
	CDNGuildIcon = func(guildID string, hash string, size int) string {
		animated := strings.HasPrefix(hash, "a_")
		format := "png"
		if animated {
			format = "gif"
		}
		return fmt.Sprintf(CDN + "icons/%s/%s.%s?size=%d", guildID, hash, format, size)
	}
)

// Route is a basic struct containing Method and URL
type Route struct {
	Method Method
	URL    string
}

// NewRoute generates a new Route struct
func NewRoute(method Method, url string) Route {
	return Route{
		Method: method,
		URL:    url,
	}
}

// Compile builds a full request URL based on arguments
func (r Route) Compile(args ...interface{}) string {
	if len(args) == 0 {
		return API + r.URL
	}
	return API + fmt.Sprintf(r.URL, args)
}
