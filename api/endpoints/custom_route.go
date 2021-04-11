package endpoints

// CustomRoute is APIRoute but custom for you
type CustomRoute struct {
	APIRoute
}

// Compile returns a CompiledAPIRoute
func (r CustomRoute) Compile(args ...interface{}) (*CompiledAPIRoute, error) {
	compiledRoute, err := r.Route.Compile(args...)
	if err != nil {
		return nil, err
	}
	return &CompiledAPIRoute{
		CompiledRoute: compiledRoute,
		method:        r.method,
	}, nil
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
