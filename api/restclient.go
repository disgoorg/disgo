package api

import (
	"errors"

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
	Request(route endpoints.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) error
	GetUser(userID Snowflake) (*User, error)
	GetMember(guildID Snowflake, userID Snowflake) (*Member, error)
	SendMessage(channelID Snowflake, message Message) (*Message, error)
	OpenDMChannel(userID Snowflake) (*DMChannel, error)
	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) error

	GetGlobalApplicationCommands(applicationID Snowflake) ([]*ApplicationCommand, error)
	CreateGlobalApplicationGlobalCommand(applicationID Snowflake, command ApplicationCommand) (*ApplicationCommand, error)
	SetGlobalApplicationCommands(applicationID Snowflake, commands ...ApplicationCommand) ([]*ApplicationCommand, error)
	GetGlobalApplicationCommand(applicationID Snowflake, commandID Snowflake) (*ApplicationCommand, error)
	EditGlobalApplicationCommand(applicationID Snowflake, commandID Snowflake, command ApplicationCommand) (*ApplicationCommand, error)
	DeleteGlobalApplicationCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildApplicationCommands(applicationID Snowflake, guildID Snowflake) ([]*ApplicationCommand, error)
	CreateGuildApplicationGuildCommand(applicationID Snowflake, guildID Snowflake, command ApplicationCommand) (*ApplicationCommand, error)
	SetGuildApplicationCommands(applicationID Snowflake, guildID Snowflake, commands ...ApplicationCommand) ([]*ApplicationCommand, error)
	GetGuildApplicationCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*ApplicationCommand, error)
	EditGuildApplicationCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command ApplicationCommand) (*ApplicationCommand, error)
	DeleteGuildApplicationCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

	SendInteractionResponse(interactionID Snowflake, interactionToken string, interactionResponse InteractionResponse) error
	EditInteractionResponse(applicationID Snowflake, interactionToken string, interactionResponse InteractionResponse) (*Message, error)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken string) error

	SendFollowupMessage(applicationID Snowflake, interactionToken string, followupMessage FollowupMessage) (*Message, error)
	EditFollowupMessage(applicationID Snowflake, interactionToken string, messageID Snowflake, followupMessage FollowupMessage) (*Message, error)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken string, followupMessageID Snowflake) error
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}
