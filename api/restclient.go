package api

import (
	"errors"

	"github.com/DisgoOrg/restclient"
)

// ErrMaxCommands returned if a Guild reached max of 100 Command(s)
var ErrMaxCommands = errors.New("you can provide a max of 100 application commands")

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

	GetUser(userID Snowflake) (*User, restclient.RestError)
	GetSelfUser() (*SelfUser, restclient.RestError)
	UpdateSelfUser(updateSelfUser UpdateSelfUser) (*SelfUser, restclient.RestError)
	GetGuilds(before int, after int, limit int) ([]*PartialGuild, restclient.RestError)
	LeaveGuild(guildID Snowflake) restclient.RestError
	GetDMChannels() ([]*DMChannel, restclient.RestError)
	CreateDMChannel(userID Snowflake) (*DMChannel, restclient.RestError)

	GetMessage(channelID Snowflake, messageID Snowflake) (*Message, restclient.RestError)
	CreateMessage(channelID Snowflake, message MessageCreate) (*Message, restclient.RestError)
	UpdateMessage(channelID Snowflake, messageID Snowflake, messageUpdate MessageUpdate) (*Message, restclient.RestError)
	DeleteMessage(channelID Snowflake, messageID Snowflake) restclient.RestError
	BulkDeleteMessages(channelID Snowflake, messageIDs ...Snowflake) restclient.RestError
	CrosspostMessage(channelID Snowflake, messageID Snowflake) (*Message, restclient.RestError)

	GetGuild(guildID Snowflake, withCounts bool) (*Guild, restclient.RestError)
	GetGuildPreview(guildID Snowflake) (*GuildPreview, restclient.RestError)
	CreateGuild(createGuild CreateGuild) (*Guild, restclient.RestError)
	UpdateGuild(guildID Snowflake, updateGuild UpdateGuild) (*Guild, restclient.RestError)
	DeleteGuild(guildID Snowflake) restclient.RestError

	GetMember(guildID Snowflake, userID Snowflake) (*Member, restclient.RestError)
	GetMembers(guildID Snowflake) ([]*Member, restclient.RestError)
	SearchMembers(guildID Snowflake, query string, limit int) ([]*Member, restclient.RestError)
	AddMember(guildID Snowflake, userID Snowflake, addMember AddMember) (*Member, restclient.RestError)
	RemoveMember(guildID Snowflake, userID Snowflake, reason string) restclient.RestError
	UpdateMember(guildID Snowflake, userID Snowflake, updateMember UpdateMember) (*Member, restclient.RestError)
	UpdateSelfNick(guildID Snowflake, nick string) (*string, restclient.RestError)
	MoveMember(guildID Snowflake, userID Snowflake, channelID *Snowflake) (*Member, restclient.RestError)
	AddMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) restclient.RestError
	RemoveMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) restclient.RestError

	GetRoles(guildID Snowflake) ([]*Role, restclient.RestError)
	CreateRole(guildID Snowflake, createRole CreateRole) (*Role, restclient.RestError)
	UpdateRole(guildID Snowflake, roleID Snowflake, updateRole UpdateRole) (*Role, restclient.RestError)
	UpdateRolePositions(guildID Snowflake, roleUpdates ...UpdateRolePosition) ([]*Role, restclient.RestError)
	DeleteRole(guildID Snowflake, roleID Snowflake) restclient.RestError

	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) restclient.RestError
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) restclient.RestError
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) restclient.RestError

	GetGlobalCommands(applicationID Snowflake) ([]*Command, restclient.RestError)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*Command, restclient.RestError)
	CreateGlobalCommand(applicationID Snowflake, command CommandCreate) (*Command, restclient.RestError)
	SetGlobalCommands(applicationID Snowflake, commands ...CommandCreate) ([]*Command, restclient.RestError)
	UpdateGlobalCommand(applicationID Snowflake, commandID Snowflake, command CommandUpdate) (*Command, restclient.RestError)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) restclient.RestError

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*Command, restclient.RestError)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*Command, restclient.RestError)
	CreateGuildCommand(applicationID Snowflake, guildID Snowflake, command CommandCreate) (*Command, restclient.RestError)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...CommandCreate) ([]*Command, restclient.RestError)
	UpdateGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command CommandUpdate) (*Command, restclient.RestError)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) restclient.RestError

	GetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake) ([]*GuildCommandPermissions, restclient.RestError)
	GetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, restclient.RestError)
	SetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake, commandPermissions ...SetGuildCommandPermissions) ([]*GuildCommandPermissions, restclient.RestError)
	SetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake, commandPermissions SetGuildCommandPermissions) (*GuildCommandPermissions, restclient.RestError)

	SendInteractionResponse(interactionID Snowflake, interactionToken string, interactionResponse InteractionResponse) restclient.RestError
	UpdateInteractionResponse(applicationID Snowflake, interactionToken string, messageUpdate MessageUpdate) (*Message, restclient.RestError)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken string) restclient.RestError

	SendFollowupMessage(applicationID Snowflake, interactionToken string, messageCreate MessageCreate) (*Message, restclient.RestError)
	UpdateFollowupMessage(applicationID Snowflake, interactionToken string, messageID Snowflake, messageUpdate MessageUpdate) (*Message, restclient.RestError)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken string, followupMessageID Snowflake) restclient.RestError

	GetGuildTemplate(templateCode string) (*GuildTemplate, restclient.RestError)
	GetGuildTemplates(guildID Snowflake) ([]*GuildTemplate, restclient.RestError)
	CreateGuildTemplate(guildID Snowflake, createGuildTemplate CreateGuildTemplate) (*GuildTemplate, restclient.RestError)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate CreateGuildFromTemplate) (*Guild, restclient.RestError)
	SyncGuildTemplate(guildID Snowflake, templateCode string) (*GuildTemplate, restclient.RestError)
	UpdateGuildTemplate(guildID Snowflake, templateCode string, updateGuildTemplate UpdateGuildTemplate) (*GuildTemplate, restclient.RestError)
	DeleteGuildTemplate(guildID Snowflake, templateCode string) (*GuildTemplate, restclient.RestError)
}
