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
	Disgo() Disgo

	UserAgent() string
	Request(route endpoints.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) error

	SendMessage(channelID Snowflake, message MessageCreate) (*Message, error)
	OpenDMChannel(userID Snowflake) (*DMChannel, error)

	UpdateSelfNick(guildID Snowflake, nick *string) (*string, error)

	GetUser(userID Snowflake) (*User, error)
	GetMember(guildID Snowflake, userID Snowflake) (*Member, error)
	GetMembers(guildID Snowflake) ([]*Member, error)
	AddMember(guildID Snowflake, userID Snowflake, addGuildMemberData AddGuildMemberData) (*Member, error)
	KickMember(guildID Snowflake, userID Snowflake, reason *string) error
	UpdateMember(guildID Snowflake, userID Snowflake, updateGuildMemberData UpdateGuildMemberData) (*Member, error)
	MoveMember(guildID Snowflake, userID Snowflake, channelID *Snowflake) (*Member, error)
	AddMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error
	RemoveMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error

	GetRoles(guildID Snowflake) ([]*Role, error)
	CreateRole(guildID Snowflake, role UpdateRole) (*Role, error)
	UpdateRole(guildID Snowflake, roleID Snowflake, role UpdateRole) (*Role, error)
	UpdateRolePositions(guildID Snowflake, roleUpdates ...UpdateRolePosition) ([]*Role, error)
	DeleteRole(guildID Snowflake, roleID Snowflake) error

	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) error

	GetGlobalCommands(applicationID Snowflake) ([]*SlashCommand, error)
	CreateGlobalCommand(applicationID Snowflake, command SlashCommand) (*SlashCommand, error)
	SetGlobalCommands(applicationID Snowflake, commands ...SlashCommand) ([]*SlashCommand, error)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*SlashCommand, error)
	EditGlobalCommand(applicationID Snowflake, commandID Snowflake, command SlashCommand) (*SlashCommand, error)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*SlashCommand, error)
	CreateGuildGuildCommand(applicationID Snowflake, guildID Snowflake, command SlashCommand) (*SlashCommand, error)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...SlashCommand) ([]*SlashCommand, error)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*SlashCommand, error)
	EditGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command SlashCommand) (*SlashCommand, error)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

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
