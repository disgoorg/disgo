package endpoints

import (
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
func (r Route) Compile(args ...string) string {
	if len(args) != r.paramCount {
		log.Errorf("invalid amount of arguments received. expected: %d, received: %d", r.paramCount, len(args))
	}
	route := r.route
	if len(args) > 0 {
		for _, arg := range args {
			start := strings.Index(route, "{")
			end := strings.Index(route, "}")
			route = route[:start] + arg + route[end+1:]
		}
	}

	return r.baseRoute + route
}

func countParams(url string) int {
	paramCount := strings.Count(url, "{")
	if paramCount != strings.Count(url, "}") {
		log.Errorf("invalid format for route provided: %s", url)
	}
	return paramCount
}
