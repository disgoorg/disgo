package events

import (
	"net/http"
)

type HttpRequestEvent struct {
	GenericEvent
	Request  *http.Request
	Response *http.Response
}
