package endpoints

import (
	"fmt"
)

const (
	APIVersion = "8"
	Base       = "https://discord.com/"
	CDN        = "https://cdn.discordapp.com/"
	API        = Base + "api/v" + APIVersion + "/"
)

type Route struct {
	Method Method
	Url    string
}

func NewRoute(method Method, url string) Route {
	return Route{
		Method: method,
		Url:    url,
	}
}

func (r Route) Compile(args ...interface{}) string {
	if len(args) == 0 {
		return API + r.Url
	}
	return API + fmt.Sprintf(r.Url, args)
}
