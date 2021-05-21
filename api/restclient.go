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

// ErrorResponse contains custom errors from discord
type ErrorResponse struct {
	Code    int
	Message string
}

// RestClient is a manager for all of disgo's HTTP requests
type RestClient interface {
	Close()
	Disgo() Disgo

	UserAgent() string
	Request(route *endpoints.CompiledAPIRoute, rqBody interface{}, rsBody interface{}) error

	GetGateway() (*GatewayRs, error)
	GetGatewayBot() (*GatewayBotRs, error)
	GetBotApplication() (*Application, error)
	GetVoiceRegions() ([]*VoiceRegion, error)

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
	AddMember(guildID Snowflake, userID Snowflake, addGuildMemberData *AddGuildMember) (*Member, error)
	UpdateMember(guildID Snowflake, userID Snowflake, updateGuildMemberData *UpdateGuildMember) (*Member, error)
	RemoveMember(guildID Snowflake, userID Snowflake, reason *string) error
	AddMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error
	RemoveMemberRole(guildID Snowflake, userID Snowflake, roleID Snowflake) error

	UpdateSelfNick(guildID Snowflake, nick *string) (*string, error)

	GetPruneMembersCount(guildID Snowflake, days int, includeRoles []Snowflake) (*int, error)
	PruneMembers(guildID Snowflake, days int, computePruneCount bool, includeRoles []Snowflake, reason string) (*int, error)

	GetGuildWebhooks(guildID Snowflake)

	GetAuditLogs(guildID Snowflake)

	GetGuildVoiceRegions(guildID Snowflake) ([]*VoiceRegion, error)

	GetGuildIntegrations(guildID Snowflake) ([]*Integration, error)
	CreateGuildIntegration(guildID Snowflake)
	UpdateGuildIntegration(guildID Snowflake)
	DeleteGuildIntegration(guildID Snowflake)
	SyncIntegration(guildID Snowflake)

	GetRoles(guildID Snowflake) ([]*Role, error)
	CreateRole(guildID Snowflake, role *RoleUpdate) (*Role, error)
	UpdateRole(guildID Snowflake, roleID Snowflake, role *RoleUpdate) (*Role, error)
	UpdateRolePositions(guildID Snowflake, roleUpdates ...*RolePositionUpdate) ([]*Role, error)
	DeleteRole(guildID Snowflake, roleID Snowflake) error

	CreateMessage(channelID Snowflake, message *MessageCreate) (*Message, error)
	UpdateMessage(channelID Snowflake, messageID Snowflake, message *MessageUpdate) (*Message, error)
	DeleteMessage(channelID Snowflake, messageID Snowflake) error
	BulkDeleteMessages(channelID Snowflake, messageIDs ...Snowflake) error
	CrosspostMessage(channelID Snowflake, messageID Snowflake) (*Message, error)

	AddReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveOwnReaction(channelID Snowflake, messageID Snowflake, emoji string) error
	RemoveUserReaction(channelID Snowflake, messageID Snowflake, emoji string, userID Snowflake) error

	GetGlobalCommands(applicationID Snowflake) ([]*Command, error)
	GetGlobalCommand(applicationID Snowflake, commandID Snowflake) (*Command, error)
	CreateGlobalCommand(applicationID Snowflake, command *CommandCreate) (*Command, error)
	SetGlobalCommands(applicationID Snowflake, commands ...*CommandCreate) ([]*Command, error)
	UpdateGlobalCommand(applicationID Snowflake, commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteGlobalCommand(applicationID Snowflake, commandID Snowflake) error

	GetGuildCommands(applicationID Snowflake, guildID Snowflake) ([]*Command, error)
	GetGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*Command, error)
	CreateGuildCommand(applicationID Snowflake, guildID Snowflake, command *CommandCreate) (*Command, error)
	SetGuildCommands(applicationID Snowflake, guildID Snowflake, commands ...*CommandCreate) ([]*Command, error)
	UpdateGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake, command *CommandUpdate) (*Command, error)
	DeleteGuildCommand(applicationID Snowflake, guildID Snowflake, commandID Snowflake) error

	GetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake) ([]*GuildCommandPermissions, error)
	GetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake) (*GuildCommandPermissions, error)
	SetGuildCommandsPermissions(applicationID Snowflake, guildID Snowflake, commandPermissions ...*SetGuildCommandPermissions) ([]*GuildCommandPermissions, error)
	SetGuildCommandPermissions(applicationID Snowflake, guildID Snowflake, commandID Snowflake, commandPermissions *SetGuildCommandPermissions) (*GuildCommandPermissions, error)

	CreateInteractionResponse(interactionID Snowflake, interactionToken string, interactionResponse *InteractionResponse) error
	UpdateInteractionResponse(applicationID Snowflake, interactionToken string, followupMessage *FollowupMessage) (*Message, error)
	DeleteInteractionResponse(applicationID Snowflake, interactionToken string) error

	CreateFollowupMessage(applicationID Snowflake, interactionToken string, followupMessage *FollowupMessage) (*Message, error)
	UpdateFollowupMessage(applicationID Snowflake, interactionToken string, messageID Snowflake, followupMessage *FollowupMessage) (*Message, error)
	DeleteFollowupMessage(applicationID Snowflake, interactionToken string, followupMessageID Snowflake) error
}
