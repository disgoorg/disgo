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
