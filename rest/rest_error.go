package rest

import (
	"fmt"
	"net/http"

	"github.com/disgoorg/json"
)

// JSONErrorCode is the error code returned by the Discord API.
// See https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
type JSONErrorCode int

var _ error = (*Error)(nil)

// Error holds the http.Response & an error related to a REST request
type Error struct {
	Request  *http.Request  `json:"-"`
	RqBody   []byte         `json:"-"`
	Response *http.Response `json:"-"`
	RsBody   []byte         `json:"-"`

	Code    JSONErrorCode   `json:"code"`
	Errors  json.RawMessage `json:"errors"`
	Message string          `json:"message"`
}

// NewError returns a new Error with the given http.Request, http.Response
func NewError(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte) error {
	var err Error
	_ = json.Unmarshal(rsBody, &err)

	err.Request = rq
	err.RqBody = rqBody
	err.Response = rs
	err.RsBody = rsBody

	return err
}

// Is returns true if the error is a *Error with the same status code as the target error
func (e Error) Is(target error) bool {
	err, ok := target.(*Error)
	if !ok {
		return false
	}
	if e.Code != 0 && err.Code != 0 {
		return e.Code == err.Code
	}
	return err.Response != nil && e.Response != nil && err.Response.StatusCode == e.Response.StatusCode
}

// Error returns the error formatted as string
func (e Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("%d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("Status: %s, Body: %s", e.Response.Status, string(e.RsBody))
}

// Error returns the error formatted as string
func (e Error) String() string {
	return e.Error()
}
