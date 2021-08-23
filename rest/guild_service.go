package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewGuildService(client Client) GuildService {
	return nil
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

	GetEmojis(opts ...rest.RequestOpt) ([]discord.Emoji, Error)

	GetGuildTemplate(templateCode string) (*discord.GuildTemplate, Error)
	GetGuildTemplates(guildID discord.Snowflake) ([]discord.GuildTemplate, Error)
	CreateGuildTemplate(guildID discord.Snowflake, createGuildTemplate discord.GuildTemplateCreate) (*discord.GuildTemplate, Error)
	CreateGuildFromTemplate(templateCode string, createGuildFromTemplate discord.GuildFromTemplateCreate) (*discord.Guild, Error)
	SyncGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
	UpdateGuildTemplate(guildID discord.Snowflake, templateCode string, updateGuildTemplate discord.GuildTemplateUpdate) (*discord.GuildTemplate, Error)
	DeleteGuildTemplate(guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, Error)
}
