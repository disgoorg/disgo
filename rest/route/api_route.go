package route

import (
	"fmt"
	"net/url"
	"strings"
)

// NewAPIRoute generates a new discord api path struct
func NewAPIRoute(method Method, path string, queryParams ...string) *APIRoute {
	return newAPIRoute(method, path, queryParams, true)
}

func NewAPIRouteNoAuth(method Method, path string, queryParams ...string) *APIRoute {
	return newAPIRoute(method, path, queryParams, false)
}

func newAPIRoute(method Method, path string, queryParams []string, needsAuth bool) *APIRoute {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &APIRoute{
		basePath:      API,
		path:          path,
		queryParams:   params,
		urlParamCount: countURLParams(path),
		method:        method,
		needsAuth:     needsAuth,
	}
}

// NewCustomAPIRoute generates a new custom path struct
//goland:noinspection GoUnusedExportedFunction
func NewCustomAPIRoute(method Method, basePath string, path string, queryParams ...string) *APIRoute {
	route := NewAPIRoute(method, path, queryParams...)
	route.basePath = basePath
	return route
}

// APIRoute is a basic struct containing Method and URL
type APIRoute struct {
	basePath      string
	path          string
	queryParams   map[string]struct{}
	urlParamCount int
	method        Method
	needsAuth     bool
}

// Compile returns a CompiledAPIRoute
func (r *APIRoute) Compile(queryValues QueryValues, params ...interface{}) (*CompiledAPIRoute, error) {
	if len(params) != r.urlParamCount {
		return nil, ErrInvalidArgCount(r.urlParamCount, len(params))
	}
	path := r.path
	var majorParams []string
	for _, param := range params {
		start := strings.Index(path, "{")
		end := strings.Index(path, "}")
		paramName := path[start+1 : end]
		paramValue := fmt.Sprint(param)
		if strings.Contains(MajorParameters, paramName) {
			majorParams = append(majorParams, paramName+"="+paramValue)
		}
		path = path[:start] + paramValue + path[end+1:]
	}

	queryParamsStr := ""
	if queryValues != nil {
		query := url.Values{}
		for param, value := range queryValues {
			if _, ok := r.queryParams[param]; !ok {
				return nil, ErrUnexpectedQueryParam(param)
			}
			query.Add(param, fmt.Sprint(value))
		}
		if len(query) > 0 {
			queryParamsStr = query.Encode()
		}
	}

	return &CompiledAPIRoute{
		APIRoute:    r,
		path:        path,
		queryParams: queryParamsStr,
		majorParams: strings.Join(majorParams, ":"),
	}, nil
}

// Method returns the request method used by the path
func (r *APIRoute) Method() Method {
	return r.method
}

// Path returns the request path used by the path
func (r *APIRoute) Path() string {
	return r.path
}

// NeedsAuth returns whether the route requires authentication
func (r *APIRoute) NeedsAuth() bool {
	return r.needsAuth
}

// CompiledAPIRoute is APIRoute compiled with all URL args
type CompiledAPIRoute struct {
	APIRoute    *APIRoute
	path        string
	queryParams string
	majorParams string
}

// MajorParams returns the major parameter from the request
func (r *CompiledAPIRoute) MajorParams() string {
	return r.majorParams
}

// URL returns the full URL for the request
func (r *CompiledAPIRoute) URL() string {
	u := r.APIRoute.basePath + r.path
	if r.queryParams != "" {
		u += "?" + r.queryParams
	}
	return u
}
