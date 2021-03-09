package disgo

import (
	"errors"

	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/endpoints"
	"github.com/DiscoOrg/disgo/models"
)

var (
	BadGatewayError = errors.New("bad gateway could not reach discord")
	UnauthorizedError = errors.New("not authorized for this endpoint")
	RatelimitedError = errors.New("too many requests")
)

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	Disgo() Disgo
	Close()
	UserAgent() string
	Request(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) error
	RequestAsync(route endpoints.Route, rqBody interface{}, v interface{}, args ...interface{}) *promise.Promise
	GetUserById(models.Snowflake) *promise.Promise
	GetMemberById(models.Snowflake, models.Snowflake) *promise.Promise
	SendMessage(models.Snowflake, models.Message) *promise.Promise
	OpenDMChannel(models.Snowflake) *promise.Promise
	AddReaction(models.Snowflake, models.Snowflake, string) *promise.Promise
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}
