package events

import (
	"net/http"

	"github.com/DisgoOrg/disgo/api/endpoints"
)

type GenericHttpEvent struct {
	Request endpoints.Request
	Response endpoints.Response
}

func (e GenericHttpEvent) RateLimited() bool {
	return e.Response.StatusCode == http.StatusTooManyRequests
}
