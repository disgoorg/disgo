package route

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
)

// NewRoute generates a new discord path struct
//goland:noinspection GoUnusedExportedFunction
func NewRoute(path string, queryParams ...string) *Route {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &Route{
		basePath:      API,
		path:          path,
		queryParams:   params,
		urlParamCount: countURLParams(path),
	}
}

// NewCustomRoute generates a new custom path struct
//goland:noinspection GoUnusedExportedFunction
func NewCustomRoute(basePath string, path string, queryParams ...string) *Route {
	route := NewRoute(path, queryParams...)
	route.basePath = basePath
	return route
}

// Route is a basic struct containing Method and URL
type Route struct {
	basePath      string
	path          string
	queryParams   map[string]struct{}
	urlParamCount int
}

// Compile returns a CompiledRoute
func (r *Route) Compile(queryValues QueryValues, params ...interface{}) (*CompiledRoute, error) {
	if len(params) != r.urlParamCount {
		return nil, discord.ErrInvalidArgCount(r.urlParamCount, len(params))
	}
	path := r.path
	for _, param := range params {
		start := strings.Index(path, "{")
		end := strings.Index(path, "}")
		path = path[:start] + fmt.Sprint(param) + path[end+1:]
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

	return &CompiledRoute{
		route:       r,
		path:        path,
		queryParams: queryParamsStr,
	}, nil
}

// Path returns the request path used by the path
func (r *Route) Path() string {
	return r.path
}

// CompiledRoute is Route compiled with all URL args
type CompiledRoute struct {
	route       *Route
	path        string
	queryParams string
}

func (r *CompiledRoute) URL() string {
	u := r.route.basePath + r.path
	if r.queryParams != "" {
		u += "?" + r.queryParams
	}
	return u
}
