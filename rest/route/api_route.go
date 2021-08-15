package route

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
)

const MajorParameters = "guild.id:channel.id:webhook.id:interaction.token"

// NewAPIRoute generates a new discord api route struct
func NewAPIRoute(method Method, url string, queryParams ...string) *APIRoute {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &APIRoute{
		baseRoute:     baseRoute,
		route:         url,
		queryParams:   params,
		urlParamCount: countURLParams(url),
		method:        method,
	}
}

// APIRoute is a basic struct containing Method and URL
type APIRoute struct {
	baseRoute     string
	route         string
	queryParams   map[string]struct{}
	urlParamCount int
	method        Method
}

// Compile returns a CompiledAPIRoute
func (r *APIRoute) Compile(queryValues QueryValues, params ...interface{}) (*CompiledAPIRoute, error) {
	if len(params) != r.urlParamCount {
		return nil, discord.ErrInvalidArgCount(len(params), r.urlParamCount)
	}
	route := r.route
	var major []string
	for _, param := range params {
		start := strings.Index(route, "{")
		end := strings.Index(route, "}")
		paramName := route[start+1 : end-1]
		paramValue := fmt.Sprint(param)
		if strings.Contains(MajorParameters, paramName) {
			major = append(major, paramName+"="+paramValue)
		}
		route = route[:start] + paramValue + route[end+1:]
	}


	compiledRoute := r.baseRoute + route
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
		route:       route,
		queryParams: queryParamsStr,
	}, nil
}

// Method returns the request method used by the route
func (r *APIRoute) Method() Method {
	return r.method
}

// CompiledAPIRoute is APIRoute compiled with all URL args
type CompiledAPIRoute struct {
	*APIRoute
	route       string
	queryParams string
	major       string
}

// Method returns the request method used by the route
func (r *CompiledAPIRoute) Method() Method {
	return r.method
}
