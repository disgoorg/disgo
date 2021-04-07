package api

import (
	"errors"

	"github.com/DisgoOrg/disgo/api/endpoints"
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
	EditMessage(channelID Snowflake, messageID Snowflake, message MessageUpdate) (*Message, error)
	DeleteMessage(channelID Snowflake, messageID Snowflake) error
	BulkDeleteMessages(channelID Snowflake, messageIDs ...Snowflake) error
	CrosspostMessage(channelID Snowflake, messageID Snowflake) (*Message, error)

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

	GetGlobalCommands(applicationID Snowflake) ([]*Command, error)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*Command, error)
	CreateGlobalCommand(applicationID Snowflake, command Command) (*Command, error)
	SetGlobalCommands(applicationID Snowflake, commands ...Command) ([]*Command, error)
	EditGlobalCommand(applicationID Snowflake, commandID Snowflake, command UpdateCommand) (*Command, error)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*Command, error)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*Command, error)
	CreateGuildCommand(applicationID Snowflake, guildID Snowflake, command Command) (*Command, error)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...Command) ([]*Command, error)
	EditGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command UpdateCommand) (*Command, error)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

	GetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake) ([]*GuildCommandPermissions, error)
	GetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake, commandPermissions ...SetGuildCommandPermissions) ([]*GuildCommandPermissions, error)
	SetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake, commandPermissions SetGuildCommandPermissions) (*GuildCommandPermissions, error)

	SendInteractionResponse(interactionID Snowflake, interactionToken endpoints.Token, interactionResponse InteractionResponse) error
	EditInteractionResponse(applicationID Snowflake, interactionToken endpoints.Token, followupMessage FollowupMessage) (*Message, error)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken endpoints.Token) error

	SendFollowupMessage(applicationID Snowflake, interactionToken endpoints.Token, followupMessage FollowupMessage) (*Message, error)
	EditFollowupMessage(applicationID Snowflake, interactionToken endpoints.Token, messageID Snowflake, followupMessage FollowupMessage) (*Message, error)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken endpoints.Token, followupMessageID Snowflake) error
}

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}
