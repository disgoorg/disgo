package events

import (
	"net/http"
)

// HTTPRequest indicates a new http.Request was made and can be used to collect data of StatusCodes
type HTTPRequest struct {
	*GenericEvent
	Request  *http.Request
	Response *http.Response
}
