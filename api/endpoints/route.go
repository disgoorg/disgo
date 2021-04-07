package endpoints

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Route the base struct for routes used in disgo
type Route struct {
	baseRoute  string
	route      string
	paramCount int
}

// Compile builds a full request URL based on provided arguments
func (r Route) Compile(args ...interface{}) CompiledRoute {
	if len(args) != r.paramCount {
		log.Errorf("invalid amount of arguments received. expected: %d, received: %d", r.paramCount, len(args))
	}
	route := r.route
	if len(args) > 0 {
		for _, arg := range args {
			start := strings.Index(route, "{")
			end := strings.Index(route, "}")
			var value string
			if t, ok := arg.(Token); ok {
				value = string(t)
			} else {
				value = fmt.Sprint(arg)
			}
			route = route[:start] + value + route[end+1:]
		}
	}

	return CompiledRoute{route: r.baseRoute + route}
}

func countParams(url string) int {
	paramCount := strings.Count(url, "{")
	if paramCount != strings.Count(url, "}") {
		log.Errorf("invalid format for route provided: %s", url)
	}
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
