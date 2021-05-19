package endpoints

// APIRoute is a basic struct containing Method and URL
type APIRoute struct {
	*Route
	method Method
}

// Compile returns a CompiledAPIRoute
func (r *APIRoute) Compile(queryParams map[string]interface{}, args ...interface{}) (*CompiledAPIRoute, error) {
	compiledRoute, err := r.Route.Compile(queryParams, args...)
	if err != nil {
		return nil, err
	}
	return &CompiledAPIRoute{
		CompiledRoute: compiledRoute,
		method:        r.method,
	}, nil
}

// Method returns the request method used by the route
func (r *APIRoute) Method() Method {
	return r.method
}

// NewAPIRoute generates a new discord api route struct
func NewAPIRoute(method Method, url string, queryParams ...string) *APIRoute {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &APIRoute{
		Route: &Route{
			baseRoute:   API,
			route:       url,
			queryParams: params,
			paramCount:  countParams(url),
		},
		method: method,
	}
}

// CompiledAPIRoute is APIRoute compiled with all URL args
type CompiledAPIRoute struct {
	*CompiledRoute
	method Method
}

// Method returns the request method used by the route
func (r *CompiledAPIRoute) Method() Method {
	return r.method
}
