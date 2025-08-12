package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/disgoorg/json/v2"
)

// JSONErrorCode is the error code returned by the Discord API.
// See https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
type JSONErrorCode int

var _ error = (*Error)(nil)

// Error holds the *[http.Request], *[http.Response] & an error related to a REST request.
// It's always a pointer to *[Error] that is returned by the REST client.
type Error struct {
	Request  *http.Request  `json:"-"`
	RqBody   []byte         `json:"-"`
	Response *http.Response `json:"-"`
	RsBody   []byte         `json:"-"`

	Code    JSONErrorCode   `json:"code"`
	Errors  json.RawMessage `json:"errors"`
	Message string          `json:"message"`
}

// newError returns a new *Error with the given http.Request, http.Response
func newError(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte) *Error {
	err := &Error{
		Request:  rq,
		RqBody:   rqBody,
		Response: rs,
		RsBody:   rsBody,
	}
	_ = json.Unmarshal(rsBody, &err)

	return err
}

// Is returns true if the error is a *Error with the same status code as the target error
func (e *Error) Is(target error) bool {
	var err *Error
	if ok := errors.As(target, &err); !ok {
		return false
	}
	if e.Code != 0 && err.Code != 0 {
		return e.Code == err.Code
	}
	return err.Response != nil && e.Response != nil && err.Response.StatusCode == e.Response.StatusCode
}

// Error returns the error formatted as string
func (e *Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("%d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("Status: %s, Body: %s", e.Response.Status, string(e.RsBody))
}

// Error returns the error formatted as string
func (e *Error) String() string {
	return e.Error()
}
