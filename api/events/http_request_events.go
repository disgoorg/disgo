package events

import (
	"net/http"
)

type HttpRequestEvent struct {
	GenericEvent
	Request  *http.Request
	Response *http.Response
}

func (e HttpRequestEvent) RateLimited() bool {
	return e.Response.StatusCode == http.StatusTooManyRequests
}
