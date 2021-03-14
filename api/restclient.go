package api

import (
	"errors"

	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

// Errors when connecting to discord
var (
	ErrBadGateway                 = errors.New("bad gateway could not reach discord")
	ErrUnauthorized               = errors.New("not authorized for this endpoint")
	ErrBadRequest                 = errors.New("bad request please check your request")
	ErrRatelimited                = errors.New("too many requests")
	ErrTooMuchApplicationCommands = errors.New("you can provide a max of 100 application commands")
)

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	Close()
	UserAgent() string
	Request(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) error
	RequestAsync(route endpoints.APIRoute, rqBody interface{}, v interface{}, args ...string) *promise.Promise
	GetUser(Snowflake) *promise.Promise
	GetMember(Snowflake, Snowflake) *promise.Promise
	SendMessage(Snowflake, Message) *promise.Promise
	OpenDMChannel(Snowflake) *promise.Promise
	AddReaction(Snowflake, Snowflake, string) *promise.Promise
	RemoveOwnReaction(Snowflake, Snowflake, string) *promise.Promise
	RemoveUserReaction(Snowflake, Snowflake, string, Snowflake) *promise.Promise

	GetGlobalApplicationCommands(Snowflake) *promise.Promise
	CreateGlobalApplicationGlobalCommand(Snowflake, ApplicationCommand) *promise.Promise
	SetGlobalApplicationCommands(Snowflake, ...ApplicationCommand) *promise.Promise
	GetGlobalApplicationCommand(Snowflake, Snowflake) *promise.Promise
	EditGlobalApplicationCommand(Snowflake, Snowflake, ApplicationCommand) *promise.Promise
	DeleteGlobalApplicationCommand(Snowflake, Snowflake) *promise.Promise

	GetGuildApplicationCommands(Snowflake, Snowflake) *promise.Promise
	CreateGuildApplicationGuildCommand(Snowflake, Snowflake, ApplicationCommand) *promise.Promise
	SetGuildApplicationCommands(Snowflake, Snowflake, ...ApplicationCommand) *promise.Promise
	GetGuildApplicationCommand(Snowflake, Snowflake, Snowflake) *promise.Promise
	EditGuildApplicationCommand(Snowflake, Snowflake, Snowflake, ApplicationCommand) *promise.Promise
	DeleteGuildApplicationCommand(Snowflake, Snowflake, Snowflake) *promise.Promise

	SendInteractionResponse(Snowflake, string, InteractionResponse) *promise.Promise
	EditInteractionResponse(Snowflake, string, InteractionResponse) *promise.Promise
	DeleteInteractionResponse(Snowflake, string) *promise.Promise
	SendFollowupMessage(Snowflake, string, FollowupMessage) *promise.Promise
	EditFollowupMessage(Snowflake, string, Snowflake, InteractionResponse) *promise.Promise
	DeleteFollowupMessage(Snowflake, string, Snowflake) *promise.Promise
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}
