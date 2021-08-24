package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewGuildService(client Client) GuildService {
	return nil
}

type GuildService interface {
	Service
	GetGuild(guildID discord.Snowflake, withCounts bool, opts ...RequestOpt) (*discord.Guild, Error)
	GetGuildPreview(guildID discord.Snowflake, opts ...RequestOpt) (*discord.GuildPreview, Error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.Guild, Error)
	UpdateGuild(guildID discord.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.Guild, Error)
	DeleteGuild(guildID discord.Snowflake, opts ...RequestOpt) Error

	GetRoles(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Role, Error)
	CreateRole(guildID discord.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, Error)
	UpdateRole(guildID discord.Snowflake, roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, Error)
	UpdateRolePositions(guildID discord.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, Error)
	DeleteRole(guildID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error

	GetMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Member, Error)
	GetMembers(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Member, Error)
	SearchMembers(guildID discord.Snowflake, query string, limit int, opts ...RequestOpt) ([]discord.Member, Error)
	AddMember(guildID discord.Snowflake, userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...RequestOpt) (*discord.Member, Error)
	RemoveMember(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error
	UpdateMember(guildID discord.Snowflake, userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...RequestOpt) (*discord.Member, Error)
	MoveMember(guildID discord.Snowflake, userID discord.Snowflake, channelID *discord.Snowflake, opts ...RequestOpt) (*discord.Member, Error)
	AddMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error
	RemoveMemberRole(guildID discord.Snowflake, userID discord.Snowflake, roleID discord.Snowflake, opts ...RequestOpt) Error

	UpdateSelfNick(guildID discord.Snowflake, nick string, opts ...RequestOpt) (*string, Error)

	GetBans(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Ban, Error)
	GetBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) (*discord.Ban, Error)
	AddBan(guildID discord.Snowflake, userID discord.Snowflake, deleteMessageDays int, opts ...RequestOpt) Error
	DeleteBan(guildID discord.Snowflake, userID discord.Snowflake, opts ...RequestOpt) Error

	GetIntegrations(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Integration, Error)
	DeleteIntegration(guildID discord.Snowflake, integrationID discord.Snowflake, opts ...RequestOpt) Error

	GetEmojis(opts ...RequestOpt) ([]discord.Emoji, Error)

	GetGuildTemplate(templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, createGuildTemplate discord.GuildTemplateCreate, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate, opts ...RequestOpt) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, updateGuildTemplate discord.GuildTemplateUpdate, opts ...RequestOpt) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string, opts ...RequestOpt) (*discord.GuildTemplate, Error)
}
