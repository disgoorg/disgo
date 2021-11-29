package events

import (
	"net/http"
)

// HTTPRequestEvent indicates a new http.Request was made and can be used to collect data of StatusCodes
type HTTPRequestEvent struct {
	*GenericEvent
	Request  *http.Request
	Response *http.Response
}
