package endpoints

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Route the base struct for routes used in disgo
type Route struct {
	baseRoute  string
	route      string
	paramCount int
}

// Compile builds a full request URL based on provided arguments
func (r Route) Compile(args ...interface{}) (*CompiledRoute, error) {
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

	return &CompiledRoute{route: r.baseRoute + route}, nil
}

// NewRoute generates a Route when given a URL
func NewRoute(url string) Route {
	return Route{
		baseRoute:  "",
		route:      url,
		paramCount: countParams(url),
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
func (r CompiledRoute) Route() string {
	return r.route
}
