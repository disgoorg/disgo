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

	GetUser(userID Snowflake) (*User, error)
	GetSelfUser() (*User, error)
	UpdateSelfUser() (*User, error)
	GetGuilds() ([]*Guild, error)
	LeaveGuild(guildID Snowflake) error
	GetDMChannels() ([]DMChannel, error)
	CreateDMChannel(userID Snowflake) (DMChannel, error)

	GetGuild(guildID Snowflake) (*Guild, error)
	CreateGuild(guildCreate GuildCreate) (*Guild, error)
	UpdateGuild(guildUpdate GuildUpdate) (*Guild, error)
	DeleteGuild(guildID Snowflake) error
	GetGuildVanityURL(guildID Snowflake) (*string, error)

	CreateGuildChannel(guildID Snowflake, channelCreate ChannelCreate) (GuildChannel, error)
	GetGuildChannels() ([]GuildChannel, error)
	UpdateGuildChannelPositions() error

	GetBans(guildID Snowflake) ([]*Ban, error)
	GetBan(guildID Snowflake, userID Snowflake) (*Ban, error)
	CreateBan(guildID Snowflake, userID Snowflake, delDays int, reason string) error
	DeleteBan(guildID Snowflake, userID Snowflake) error

	GetMember(guildID Snowflake, userID Snowflake) (*Member, error)
	GetMembers(guildID Snowflake) ([]*Member, error)
	AddMember(guildID Snowflake, userID Snowflake, addGuildMemberData AddGuildMember) (*Member, error)
	KickMember(guildID Snowflake, userID Snowflake, reason *string) error
	UpdateMember(guildID Snowflake, userID Snowflake, updateGuildMemberData UpdateGuildMember) (*Member, error)
	MoveMember(guildID Snowflake, userID Snowflake, channelID *Snowflake) (*Member, error)
	AddMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error
	RemoveMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error

	UpdateSelfNick(guildID Snowflake, nick *string) (*string, error)

	GetPruneMembersCount(guildID Snowflake, days int, includeRoles []Snowflake) (*int, error)
	PruneMembers(guildID Snowflake, days int, computePruneCount bool, includeRoles []Snowflake, reason string) (*int, error)

	GetGuildWebhooks(guildID Snowflake)

	GetAuditLogs(guildID Snowflake)

	GetGuildVoiceRegions(guildID Snowflake) ([]*VoiceRegion, error)

	GetGuildIntegrations(guildID Snowflake) ([]*Integration, error)
	CreateGuildIntegration(guildID Snowflake) (*Integration, error)
	UpdateGuildIntegration(guildID Snowflake) (*Integration, error)
	DeleteGuildIntegration(guildID Snowflake) error
	SyncIntegration(guildID Snowflake)

	GetRoles(guildID Snowflake) ([]*Role, error)
	CreateRole(guildID Snowflake, roleCreate RoleCreate) (*Role, error)
	UpdateRole(guildID Snowflake, roleID Snowflake, roleUpdate RoleUpdate) (*Role, error)
	UpdateRolePositions(guildID Snowflake, rolePositionUpdates ...RolePositionUpdate) ([]*Role, error)
	DeleteRole(guildID Snowflake, roleID Snowflake) error

	CreateMessage(channelID Snowflake, messageCreate MessageCreate) (*Message, error)
	UpdateMessage(channelID Snowflake, messageID Snowflake, messageUpdate MessageUpdate) (*Message, error)
	DeleteMessage(channelID Snowflake, messageID Snowflake) error
	BulkDeleteMessages(channelID Snowflake, messageIDs ...Snowflake) error
	CrosspostMessage(channelID Snowflake, messageID Snowflake) (*Message, error)

	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) error

	GetGlobalCommands(applicationID Snowflake) ([]*Command, error)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*Command, error)
	CreateGlobalCommand(applicationID Snowflake, commandCreate CommandCreate) (*Command, error)
	SetGlobalCommands(applicationID Snowflake, commandCreates ...CommandCreate) ([]*Command, error)
	UpdateGlobalCommand(applicationID Snowflake, commandID Snowflake, commandUpdate CommandUpdate) (*Command, error)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*Command, error)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*Command, error)
	CreateGuildCommand(applicationID Snowflake, guildID Snowflake, command CommandCreate) (*Command, error)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...CommandCreate) ([]*Command, error)
	UpdateGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command CommandUpdate) (*Command, error)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

	GetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake) ([]*GuildCommandPermissions, error)
	GetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake, commandPermissions ...SetGuildCommandPermissions) ([]*GuildCommandPermissions, error)
	SetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake, commandPermissions SetGuildCommandPermissions) (*GuildCommandPermissions, error)

	CreateInteractionResponse(interactionID Snowflake, interactionToken string, interactionResponse InteractionResponse) error
	UpdateInteractionResponse(applicationID Snowflake, interactionToken string, messageUpdate MessageUpdate) (*Message, error)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken string) error

	CreateFollowupMessage(applicationID Snowflake, interactionToken string, messageCreate MessageCreate) (*Message, error)
	UpdateFollowupMessage(applicationID Snowflake, interactionToken string, messageID Snowflake, messageUpdate MessageUpdate) (*Message, error)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken string, followupMessageID Snowflake) error
}
