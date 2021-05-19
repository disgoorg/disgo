package endpoints

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Route the base struct for routes used in disgo
type Route struct {
	baseRoute   string
	route       string
	queryParams map[string]struct{}
	paramCount  int
}

// Compile builds a full request URL based on provided arguments
func (r *Route) Compile(queryParams map[string]interface{}, args ...interface{}) (*CompiledRoute, error) {
	if len(args) != r.paramCount {
		return nil, errors.New("invalid amount of arguments received. expected: " + strconv.Itoa(len(args)) + ", received: " + strconv.Itoa(r.paramCount))
	}
	route := r.route
	if len(args) > 0 {
		for _, arg := range args {
			start := strings.Index(route, "{")
			end := strings.Index(route, "}")
			route = route[:start] + fmt.Sprint(arg) + route[end+1:]
		}
	}

	compiledRoute := r.baseRoute + route
	if queryParams != nil {
		query := url.Values{}
		for param, value := range queryParams {
			if _, ok := r.queryParams[param]; !ok {
				return nil, errors.New("unexpected query param '" + param + "' received")
			}
			query.Add(param, fmt.Sprint(value))
		}
		if len(query) > 0 {
			compiledRoute += "?" + query.Encode()
		}
	}

	return &CompiledRoute{route: compiledRoute}, nil
}

// NewRoute generates a Route when given a URL
func NewRoute(url string, queryParams ...string) *Route {
	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}
	return &Route{
		baseRoute:   "",
		route:       url,
		queryParams: params,
		paramCount:  countParams(url),
	}
}

func countParams(url string) int {
	paramCount := strings.Count(url, "{")
	return paramCount
}

// CompiledRoute is Route compiled with all URL args
type CompiledRoute struct {
	route string
}

// Route returns the full request url
func (r *CompiledRoute) Route() string {
	return r.route
}
