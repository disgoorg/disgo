package endpoints

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
	return APIRoute{
		Route: Route{
			baseRoute:  "",
			route:      url,
			paramCount: countParams(url),
		},
		method: method,
	}
}
