package rest

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var _ error = (*Error)(nil)

var (
	ErrTooManyRequests = ErrStatusCode(http.StatusTooManyRequests)
	ErrBadGateway      = ErrStatusCode(http.StatusBadGateway)
	ErrBadRequest      = ErrStatusCode(http.StatusBadRequest)
	ErrUnauthorized    = ErrStatusCode(http.StatusUnauthorized)
)

func ErrStatusCode(statusCode int) error {
	return NewErrorStatus(statusCode, errors.New(""))
}

// Error holds the http.Response & an error related to a REST request
type Error struct {
	StatusCode int
	Err        error
}

// NewError returns a new Error with the given http.Response & error
//goland:noinspection GoUnusedExportedFunction
func NewError(rs *http.Response, err error) error {
	return &Error{
		StatusCode: rs.StatusCode,
		Err:        err,
	}
}

func NewErrorStatus(statusCode int, err error) error {
	return &Error{
		StatusCode: statusCode,
		Err:        err,
	}
}

func (e Error) Is(target error) bool {
	err, ok := target.(*Error)
	if !ok {
		return false
	}
	return err.StatusCode == e.StatusCode
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d: err %v", e.StatusCode, e.Err)
}

func (e Error) String() string {
	return e.Error()
}
