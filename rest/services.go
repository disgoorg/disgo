package rest

import (
	"context"

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

type ApplicationService interface {
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
	GetAuditLog(guildID discord.Snowflake, userID discord.Snowflake, actionType discord.AuditLogEvent, before discord.Snowflake, limit int) (*discord.AuditLog, Error)
}

type GatewayService interface {
	GetGateway() (*discord.Gateway, Error)
	GetGatewayBot() (*discord.GatewayBot, Error)
}

type ChannelsService interface {
	GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, Error)
	UpdateChannel(ctx context.Context, channelID discord.Snowflake, channelUpdate discord.ChannelUpdate) (*discord.Channel, Error)
	DeleteChannel(ctx context.Context, channelID discord.Snowflake) Error

	GetWebhooks(ctx context.Context, channelID discord.Snowflake) ([]discord.Webhook, Error)
	CreateWebhook(ctx context.Context, channelID discord.Snowflake, update discord.WebhookCreate) (*discord.Webhook, Error)

	UpdatePermissionOverride(ctx context.Context, channelID discord.Snowflake, overwriteID discord.Snowflake, permissionOverwrite discord.PermissionOverwriteUpdate) Error
	DeletePermissionOverride(ctx context.Context, channelID discord.Snowflake, overwriteID discord.Snowflake) Error

	SendTyping(ctx context.Context, channelID discord.Snowflake) Error

	GetMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)
	CreateMessage(ctx context.Context, channelID discord.Snowflake, message discord.MessageCreate) (*discord.Message, Error)
	UpdateMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) Error
	BulkDeleteMessages(ctx context.Context, channelID discord.Snowflake, messageIDs ...discord.Snowflake) Error
	CrosspostMessage(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake) (*discord.Message, Error)

	AddReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveOwnReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string) Error
	RemoveUserReaction(ctx context.Context, channelID discord.Snowflake, messageID discord.Snowflake, emoji string, userID discord.Snowflake) Error
}

type GuildService interface {
	GetGuild(ctx context.Context, guildID discord.Snowflake, withCounts bool) (*discord.Guild, Error)
	GetGuildPreview(ctx context.Context, guildID discord.Snowflake) (*discord.GuildPreview, Error)
	CreateGuild(ctx context.Context, guildCreate discord.GuildCreate) (*discord.Guild, Error)
	UpdateGuild(ctx context.Context, guildID discord.Snowflake, guildUpdate discord.GuildUpdate) (*discord.Guild, Error)
	DeleteGuild(ctx context.Context, guildID discord.Snowflake) Error

	GetRoles(ctx context.Context, guildID discord.Snowflake) ([]discord.Role, Error)
	CreateRole(ctx context.Context, guildID discord.Snowflake, createRole discord.RoleCreate) (*discord.Role, Error)
	UpdateRole(ctx context.Context, guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate) (*discord.Role, Error)
	UpdateRolePositions(ctx context.Context, guildID discord.Snowflake, rolePositionUpdates ...discord.RolePositionUpdate) ([]discord.Role, Error)
	DeleteRole(ctx context.Context, guildID discord.Snowflake, roleID discord.Snowflake) Error

	GetMember(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake) (*discord.Member, Error)
	GetMembers(ctx context.Context, guildID discord.Snowflake) ([]discord.Member, Error)
	SearchMembers(ctx context.Context, guildID discord.Snowflake, query string, limit int) ([]discord.Member, Error)
	AddMember(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd) (*discord.Member, Error)
	RemoveMember(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, reason string) Error
	UpdateMember(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate) (*discord.Member, Error)
	MoveMember(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, channelID *discord.Snowflake) (*discord.Member, Error)
	AddMemberRole(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) Error
	RemoveMemberRole(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake) Error

	UpdateSelfNick(ctx context.Context, guildID discord.Snowflake, nick string) (*string, Error)

	GetBans(ctx context.Context, guildID discord.Snowflake) ([]discord.Ban, Error)
	GetBan(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake) (*discord.Ban, Error)
	AddBan(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake, reason string, deleteMessageDays int) Error
	DeleteBan(ctx context.Context, guildID discord.Snowflake, userID discord.Snowflake) Error

	GetIntegrations(ctx context.Context, guildID discord.Snowflake) ([]discord.Integration, Error)
	DeleteIntegration(ctx context.Context, guildID discord.Snowflake, integrationID discord.Snowflake) Error

	GetEmojis(ctx context.Context) ([]discord.Emoji, Error)

	GetGuildTemplate(ctx context.Context, templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(ctx context.Context, guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(ctx context.Context, guildID discord.Snowflake, createGuildTemplate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(ctx context.Context, templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string, updateGuildTemplate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}

type InteractionService interface {
	CreateInteractionResponse(ctx context.Context, interactionID discord.Snowflake, interactionToken string, interactionResponse discord.InteractionResponse) Error
	UpdateInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteInteractionResponse(ctx context.Context, applicationID discord.Snowflake, interactionToken string) Error

	CreateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageCreate discord.MessageCreate) (*discord.Message, Error)
	UpdateFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteFollowupMessage(ctx context.Context, applicationID discord.Snowflake, interactionToken string, followupMessageID discord.Snowflake) Error
}

type InviteService interface {
	GetInvite(ctx context.Context, code string) (*discord.Invite, Error)
	CreateInvite(ctx context.Context, channelID discord.Snowflake, inviteCreate discord.InviteCreate)
	DeleteInvite(ctx context.Context, code string) (*discord.Invite, Error)
	GetGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]discord.Invite, Error)
	GetChannelInvites(ctx context.Context, channelID discord.Snowflake) ([]discord.Invite, Error)
}

type GuildTemplateService interface {
	GetGuildTemplate(ctx context.Context, templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(ctx context.Context, guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(ctx context.Context, guildID discord.Snowflake, guildTemplateCreate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(ctx context.Context, templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}

type UserService interface {
	GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, Error)
	GetSelfUser(ctx context.Context) (*discord.SelfUser, Error)
	UpdateSelfUser(ctx context.Context, updateSelfUser discord.UpdateSelfUser) (*discord.SelfUser, Error)
	GetGuilds(ctx context.Context, before int, after int, limit int) ([]discord.PartialGuild, Error)
	LeaveGuild(ctx context.Context, guildID discord.Snowflake) Error
	GetDMChannels(ctx context.Context) ([]discord.Channel, Error)
	CreateDMChannel(ctx context.Context, userID discord.Snowflake) (*discord.Channel, Error)
}

type VoiceService interface {
	GetVoiceRegions(ctx context.Context) []discord.VoiceRegion
}

type StageInstanceService interface {
	GetStageInstance(ctx context.Context, stageInstanceID discord.Snowflake) (*discord.StageInstance, Error)
	CreateStageInstance(ctx context.Context, stageInstanceCreate discord.StageInstanceCreate) (*discord.StageInstance, Error)
	UpdateStageInstance(ctx context.Context, stageInstanceID discord.Snowflake, stageInstanceUpdate discord.StageInstanceUpdate) (*discord.StageInstance, Error)
	DeleteStageInstance(ctx context.Context, stageInstanceID discord.Snowflake) Error
}

type WebhookService interface {
	GetWebhook(ctx context.Context, webhookID discord.Snowflake) (*discord.Webhook, Error)
	UpdateWebhook(ctx context.Context, webhookID discord.Snowflake, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhook(ctx context.Context, webhookID discord.Snowflake) Error

	GetWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) (*discord.Webhook, Error)
	UpdateWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string, webhookUpdate discord.WebhookUpdate) (*discord.Webhook, Error)
	DeleteWebhookWithToken(ctx context.Context, webhookID discord.Snowflake, webhookToken string) Error

	CreateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageSlack(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	CreateMessageGitHub(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageCreate discord.MessageCreate, wait bool, threadID discord.Snowflake) (*discord.Message, Error)
	UpdateMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake, messageUpdate discord.MessageUpdate) (*discord.Message, Error)
	DeleteMessage(ctx context.Context, webhookID discord.Snowflake, webhookToken string, messageID discord.Snowflake) Error
}
