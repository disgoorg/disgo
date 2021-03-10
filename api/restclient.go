package api

import (
	"errors"

	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

var (
	BadGatewayError = errors.New("bad gateway could not reach discord")
	UnauthorizedError = errors.New("not authorized for this endpoint")
	RatelimitedError = errors.New("too many requests")
)

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	Close()
	UserAgent() string
	Request(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) error
	RequestAsync(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) *promise.Promise
	GetUserById(Snowflake) *promise.Promise
	GetMemberById(Snowflake, Snowflake) *promise.Promise
	SendMessage(Snowflake, Message) *promise.Promise
	OpenDMChannel(Snowflake) *promise.Promise
	AddReaction(Snowflake, Snowflake, string) *promise.Promise
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}
