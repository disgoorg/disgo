package route

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
)

// NewAPIRoute generates a new discord api path struct
//goland:noinspection GoUnusedExportedFunction
func NewAPIRoute(method Method, path string, queryParams ...string) *APIRoute {
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
}

// Compile returns a CompiledAPIRoute
func (r *APIRoute) Compile(queryValues QueryValues, params ...interface{}) (*CompiledAPIRoute, error) {
	if len(params) != r.urlParamCount {
		return nil, discord.ErrInvalidArgCount(len(params), r.urlParamCount)
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
				return nil, discord.ErrUnexpectedQueryParam(param)
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

func (r *CompiledAPIRoute) URL() string {
	u := r.APIRoute.basePath + r.path
	if r.queryParams != "" {
		u += "?" + r.queryParams
	}
	return u
}
