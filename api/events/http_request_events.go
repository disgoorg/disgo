package events

import (
	"net/http"

	"github.com/DisgoOrg/disgo/api/endpoints"
)

type HttpRequestEvent struct {
	Request  endpoints.Request
	Response endpoints.Response
}

func (e HttpRequestEvent) RateLimited() bool {
	return e.Response.StatusCode == http.StatusTooManyRequests
}
