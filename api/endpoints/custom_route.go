package endpoints

import "strings"

// CustomRoute is APIRoute but custom for you
type CustomRoute struct {
	APIRoute
}

// Compile returns a CompiledAPIRoute
func (r CustomRoute) Compile(args ...interface{}) CompiledAPIRoute {
	return CompiledAPIRoute{
		CompiledRoute: r.Route.Compile(args...),
		method:        r.method,
	}
}

// NewCustomRoute generates a new custom route struct
func NewCustomRoute(method Method, url string) APIRoute {
	urls := strings.SplitN(url, "/", 2)
	return APIRoute{
		Route: Route{
			baseRoute:  urls[0],
			route:      urls[1],
			paramCount: countParams(url),
		},
		method: method,
	}
}
