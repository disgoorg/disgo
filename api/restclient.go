package api

import (
	"errors"

	"github.com/DisgoOrg/restclient"
)

// UserAgent is the global useragent disgo uses for all its requests
var UserAgent = "DiscordBot (" + Github + ", " + Version + ")"

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	restclient.RestClient
	Close()
	Disgo() Disgo

	SendMessage(channelID Snowflake, message *MessageCreate) (*Message, error)
	EditMessage(channelID Snowflake, messageID Snowflake, message *MessageUpdate) (*Message, error)
	DeleteMessage(channelID Snowflake, messageID Snowflake) error
	BulkDeleteMessages(channelID Snowflake, messageIDs ...Snowflake) error
	CrosspostMessage(channelID Snowflake, messageID Snowflake) (*Message, error)

	OpenDMChannel(userID Snowflake) (*DMChannel, error)

	UpdateSelfNick(guildID Snowflake, nick *string) (*string, error)

	GetUser(userID Snowflake) (*User, error)
	GetMember(guildID Snowflake, userID Snowflake) (*Member, error)
	GetMembers(guildID Snowflake) ([]*Member, error)
	AddMember(guildID Snowflake, userID Snowflake, addGuildMemberData *AddGuildMemberData) (*Member, error)
	KickMember(guildID Snowflake, userID Snowflake, reason *string) error
	UpdateMember(guildID Snowflake, userID Snowflake, updateGuildMemberData *UpdateGuildMemberData) (*Member, error)
	MoveMember(guildID Snowflake, userID Snowflake, channelID *Snowflake) (*Member, error)
	AddMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error
	RemoveMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error

	GetRoles(guildID Snowflake) ([]*Role, error)
	CreateRole(guildID Snowflake, role *UpdateRole) (*Role, error)
	UpdateRole(guildID Snowflake, roleID Snowflake, role *UpdateRole) (*Role, error)
	UpdateRolePositions(guildID Snowflake, roleUpdates ...*UpdateRolePosition) ([]*Role, error)
	DeleteRole(guildID Snowflake, roleID Snowflake) error

	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) error

	GetGlobalCommands(applicationID Snowflake) ([]*Command, error)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*Command, error)
	CreateGlobalCommand(applicationID Snowflake, command *CommandCreate) (*Command, error)
	SetGlobalCommands(applicationID Snowflake, commands ...*CommandCreate) ([]*Command, error)
	EditGlobalCommand(applicationID Snowflake, commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*Command, error)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*Command, error)
	CreateGuildCommand(applicationID Snowflake, guildID Snowflake, command *CommandCreate) (*Command, error)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...*CommandCreate) ([]*Command, error)
	EditGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

	GetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake) ([]*GuildCommandPermissions, error)
	GetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake, commandPermissions ...*SetGuildCommandPermissions) ([]*GuildCommandPermissions, error)
	SetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake, commandPermissions *SetGuildCommandPermissions) (*GuildCommandPermissions, error)

	SendInteractionResponse(interactionID Snowflake, interactionToken string, interactionResponse *InteractionResponse) error
	EditInteractionResponse(applicationID Snowflake, interactionToken string, followupMessage *FollowupMessage) (*Message, error)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken string) error

	SendFollowupMessage(applicationID Snowflake, interactionToken string, followupMessage *FollowupMessage) (*Message, error)
	EditFollowupMessage(applicationID Snowflake, interactionToken string, messageID Snowflake, followupMessage *FollowupMessage) (*Message, error)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken string, followupMessageID Snowflake) error
}
