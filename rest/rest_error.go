package rest

import (
	"net/http"

	"github.com/DisgoOrg/disgo/discord"
)

var _ error = (*Error)(nil)

// Error holds the http.Response & an error related to a REST request
type Error struct {
	discord.APIError
	Request  *http.Request
	RqBody   []byte
	Response *http.Response
	RsBody   []byte
	Err      error
}

// NewErrorErr returns a new Error with the given http.Request, http.Response & error
//goland:noinspection GoUnusedExportedFunction
func NewErrorErr(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte, err error) error {
	return &Error{
		Request:  rq,
		RqBody:   rqBody,
		Response: rs,
		RsBody:   rsBody,
		Err:      err,
	}
}

func NewErrorAPIErr(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte, apiError discord.APIError) error {
	return &Error{
		APIError: apiError,
		Request:  rq,
		RqBody:   rqBody,
		Response: rs,
		RsBody:   rsBody,
	}
}

func NewError(rq *http.Request, rqBody []byte, rs *http.Response, rsBody []byte) error {
	return &Error{
		Request:  rq,
		RqBody:   rqBody,
		Response: rs,
		RsBody:   rsBody,
	}
}

func (e Error) Is(target error) bool {
	err, ok := target.(*Error)
	if !ok {
		return false
	}
	return err.Response != nil && e.Response != nil && err.Response.StatusCode == e.Response.StatusCode
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Response != nil {
		return e.Response.Status
	}
	return "unknown error"
}

func (e Error) String() string {
	return e.Error()
}
