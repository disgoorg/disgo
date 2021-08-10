package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

// Services is a manager for all of disgo's HTTP requests
type Services interface {
	Close()
	Logger() log.Logger
	HTTPClient() HTTPClient
	ApplicationService() ApplicationService
	AuditLogService() AuditLogService
	GatewayService() GatewayService
	GuildService() GuildService
	ChannelsService() ChannelsService
	InteractionService() InteractionService
	InviteService() InviteService
	GuildTemplateService() GuildTemplateService
	UserService() UserService
	VoiceService() VoiceService
	WebhookService() WebhookService
	StageInstanceService() StageInstanceService
}

type Service interface {
	Logger() log.Logger
	HTTPClient() HTTPClient
	Services() Services
}

type ApplicationService interface {
	Service
	GetGlobalCommands(applicationID discord.Snowflake) ([]discord.ApplicationCommand, Error)
	GetGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake) (*discord.ApplicationCommand, Error)
	CreateGlobalCommand(applicationID discord.Snowflake, command discord.ApplicationCommandCreate) (*discord.ApplicationCommand, Error)
	SetGlobalCommands(applicationID discord.Snowflake, commands ...discord.ApplicationCommandCreate) ([]discord.ApplicationCommand, Error)
	UpdateGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate) (*discord.ApplicationCommand, Error)
	DeleteGlobalCommand(applicationID discord.Snowflake, commandID discord.Snowflake) Error

	GetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake) ([]discord.ApplicationCommand, Error)
	GetGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (*discord.ApplicationCommand, Error)
	CreateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, command discord.ApplicationCommandCreate) (*discord.ApplicationCommand, Error)
	SetGuildCommands(applicationID discord.Snowflake, guildID discord.Snowflake, commands ...discord.ApplicationCommandCreate) ([]discord.ApplicationCommand, Error)
	UpdateGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, command discord.ApplicationCommandUpdate) (*discord.ApplicationCommand, Error)
	DeleteGuildCommand(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) Error

	GetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake) ([]discord.GuildCommandPermissions, Error)
	GetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake) (*discord.GuildCommandPermissions, Error)
	SetGuildCommandsPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandPermissions ...discord.GuildCommandPermissionsSet) ([]discord.GuildCommandPermissions, Error)
	SetGuildCommandPermissions(applicationID discord.Snowflake, guildID discord.Snowflake, commandID discord.Snowflake, commandPermissions ...discord.CommandPermission) (*discord.GuildCommandPermissions, Error)
}

type AuditLogService interface {
	Service
	GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int) (*discord.AuditLog, Error)
}

type GatewayService interface {
	Service
	GetGateway() (*discord.Gateway, Error)
	GetGatewayBot() (*discord.GatewayBot, Error)
}

type ChannelsService interface {
	Service
	GetChannel()
	UpdateChannel()
	DeleteChannel()

	GetWebhooks()
	CreateWebhook()

	GetPermissionOverrides()
	GetPermissionOverride()
	CreatePermissionOverride()
	UpdatePermissionOverride()
	DeletePermissionOverride()

	SendTyping()

	GetMessage(channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)
	CreateMessage(channelID discord.Snowflake, message discord.MessageCreate) (*discord.Message, Error)
	UpdateMessage(channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(channelID discord.Snowflake, messageID discord.Snowflake) Error
	BulkDeleteMessages(channelID discord.Snowflake, messageIDs ...discord.Snowflake) Error
	CrosspostMessage(channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)

	AddReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveOwnReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveUserReaction(channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake) Error
}

type GuildService interface {
	Service
	GetGuild(guildID discord.Snowflake, withCounts bool) (*discord.Guild, Error)
	GetGuildPreview(guildID discord.Snowflake) (*discord.GuildPreview, Error)
	CreateGuild(guildCreate discord.GuildCreate) (*discord.Guild, Error)
	UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate) (*discord.Guild, Error)
	DeleteGuild(guildID discord.Snowflake) Error

	GetRoles(guildID discord.Snowflake) ([]discord.Role, Error)
	CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate) (*discord.Role, Error)
	UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate) (*discord.Role, Error)
	UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates ...discord.RolePositionUpdate) ([]discord.Role, Error)
	DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake) Error

	GetMember(guildID discord.Snowflake, userID discord.Snowflake) (*discord.Member, Error)
	GetMembers(guildID discord.Snowflake) ([]discord.Member, Error)
	SearchMembers(guildID discord.Snowflake, query string, limit int) ([]discord.Member, Error)
	AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd) (*discord.Member, Error)
	RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, reason string) Error
	UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate) (*discord.Member, Error)
	MoveMember(guildID discord.Snowflake, userID discord.Snowflake, channelID *discord.Snowflake) (*discord.Member, Error)
	AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) Error
	RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) Error

	UpdateSelfNick(guildID discord.Snowflake, nick string) (*string, Error)

	GetBans(guildID discord.Snowflake) ([]discord.Ban, Error)
	GetBan(guildID discord.Snowflake, userID discord.Snowflake) (*discord.Ban, Error)
	AddBan(guildID discord.Snowflake, userID discord.Snowflake, reason string, deleteMessageDays int) Error
	DeleteBan(guildID discord.Snowflake, userID discord.Snowflake) Error

	GetIntegrations(guildID discord.Snowflake) ([]discord.Integration, Error)
	DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake) Error

	GetEmojis() ([]discord.Emoji, Error)

	GetGuildTemplate(templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, createGuildTemplate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, updateGuildTemplate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}

type InteractionService interface {
	Service
	CreateInteractionResponse(interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) Error
	UpdateInteractionResponse(applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteInteractionResponse(applicationID discord.Snowflake, interactionToken string) Error

	CreateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate) (*discord.Message, Error)
	UpdateFollowupMessage(applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteFollowupMessage(applicationID discord.Snowflake, interactionToken string, followupMessageID discord.Snowflake) Error
}

type InviteService interface {
	Service
	GetInvite(code string) (*discord.Invite, Error)
	CreateInvite(channelID discord.Snowflake, inviteCreate discord.CreateChannelInvite)
	DeleteInvite(code string) Error
	GetGuildInvite(guildID discord.Snowflake) Error
	GetChannelInvites(channelID discord.Snowflake) ([]discord.Invite, Error)
}

type GuildTemplateService interface {
	Service
	GetGuildTemplate(templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}

type UserService interface {
	Service
	GetUser(userID discord.Snowflake) (*discord.User, Error)
	GetSelfUser() (*discord.SelfUser, Error)
	UpdateSelfUser(updateSelfUser discord.UpdateSelfUser) (*discord.SelfUser, Error)
	GetGuilds(before int, after int, limit int) ([]discord.PartialGuild, Error)
	LeaveGuild(guildID discord.Snowflake) Error
	GetDMChannels() ([]discord.Channel, Error)
	CreateDMChannel(userID discord.Snowflake) (*discord.Channel, Error)
}

type VoiceService interface {
	Service
	GetVoiceRegions() []discord.VoiceRegion
}

type StageInstanceService interface {
	Service
	GetStageInstance(stageInstanceID discord.Snowflake) (*discord.StageInstance, Error)
	CreateStageInstance(stageInstanceCreate discord.StageInstanceCreate) (*discord.StageInstance, Error)
	UpdateStageInstance(stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (*discord.StageInstance, Error)
	DeleteStageInstance(stageInstanceID discord.Snowflake) Error
}
