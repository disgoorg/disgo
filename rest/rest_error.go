package rest

import "net/http"

var _ Error = (*errorImpl)(nil)
var _ error = (*errorImpl)(nil)

// Error holds the http.Response & an error related to a REST request
type Error interface {
	error
	Response() *http.Response
}

// NewError returns a new Error with the given http.Response & error
//goland:noinspection GoUnusedExportedFunction
func NewError(response *http.Response, err error) Error {
	return &errorImpl{
		err:      err,
		response: response,
	}
}

type errorImpl struct {
	err      error
	response *http.Response
}

// Error returns the specific error message
func (r *errorImpl) Error() string {
	return r.err.Error()
}

// Error returns the specific error message
func (r *errorImpl) String() string {
	return r.err.Error()
}

// Response returns the http.Response. May be nil depending on what broke during the request
func (r *errorImpl) Response() *http.Response {
	return r.response
}
