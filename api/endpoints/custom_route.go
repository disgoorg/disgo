package endpoints

// CustomRoute is APIRoute but custom for you
type CustomRoute struct {
	*APIRoute
}

// Compile returns a CompiledAPIRoute
func (r CustomRoute) Compile(queryParams map[string]string, args ...interface{}) (*CompiledAPIRoute, error) {
	compiledRoute, err := r.Route.Compile(queryParams, args...)
	if err != nil {
		return nil, err
	}
	return &CompiledAPIRoute{
		CompiledRoute: compiledRoute,
		method:        r.method,
	}, nil
}

// NewCustomRoute generates a new custom route struct
func NewCustomRoute(method Method, url string, queryParams ...string) *APIRoute {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &APIRoute{
		Route: &Route{
			baseRoute:   "",
			route:       url,
			queryParams: params,
			paramCount:  countParams(url),
		},
		method: method,
	}
}
