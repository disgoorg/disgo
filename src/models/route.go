package models

import (
	"fmt"
)

type Route struct {
	Method string
	Url    string
}

func NewRoute(method string, url string) Route {
	return Route{
		Method: method,
		Url:    url,
	}
}

func (r Route) compile(args ...string) string {
	return fmt.Sprintf(r.Url, args)
}
