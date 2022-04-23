package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Guilds = (*guildImpl)(nil)

func NewGuilds(client Client) Guilds {
	return &guildImpl{client: client}
}

type Guilds interface {
	GetGuild(guildID snowflake.Snowflake, withCounts bool, opts ...RequestOpt) (*discord.RestGuild, error)
	GetGuildPreview(guildID snowflake.Snowflake, opts ...RequestOpt) (*discord.GuildPreview, error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.RestGuild, error)
	UpdateGuild(guildID snowflake.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.RestGuild, error)
	DeleteGuild(guildID snowflake.Snowflake, opts ...RequestOpt) error

	CreateGuildChannel(guildID snowflake.Snowflake, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (discord.GuildChannel, error)
	GetGuildChannels(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.GuildChannel, error)
	UpdateChannelPositions(guildID snowflake.Snowflake, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error

	GetRoles(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Role, error)
	GetRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Role, error)
	CreateRole(guildID snowflake.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRolePositions(guildID snowflake.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, error)
	DeleteRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error

	GetBans(guildID snowflake.Snowflake, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) ([]discord.Ban, error)
	GetBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (*discord.Ban, error)
	AddBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, deleteMessageDays int, opts ...RequestOpt) error
	DeleteBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error

	GetIntegrations(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Integration, error)
	DeleteIntegration(guildID snowflake.Snowflake, integrationID snowflake.Snowflake, opts ...RequestOpt) error

	GetAllWebhooks(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Webhook, error)

	GetAuditLog(guildID snowflake.Snowflake, userID snowflake.Snowflake, actionType discord.AuditLogEvent, before snowflake.Snowflake, limit int, opts ...RequestOpt) (*discord.AuditLog, error)
}

type guildImpl struct {
	client Client
}

func (s *guildImpl) GetGuild(guildID snowflake.Snowflake, withCounts bool, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	values := route.QueryValues{}
	if withCounts {
		values["withCounts"] = true
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuild.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &guild, opts...)
	return
}

func (s *guildImpl) GetGuildPreview(guildID snowflake.Snowflake, opts ...RequestOpt) (guildPreview *discord.GuildPreview, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildPreview.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &guildPreview, opts...)
	return
}

func (s *guildImpl) CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuild.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, guildCreate, &guild, opts...)
	return
}

func (s *guildImpl) UpdateGuild(guildID snowflake.Snowflake, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuild.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, guildUpdate, &guild, opts...)
	return
}

func (s *guildImpl) DeleteGuild(guildID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) CreateGuildChannel(guildID snowflake.Snowflake, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (guildChannel discord.GuildChannel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildChannel.Compile(nil, guildID)
	if err != nil {
		return
	}
	var ch discord.UnmarshalChannel
	err = s.client.Do(compiledRoute, guildChannelCreate, &ch, opts...)
	if err == nil {
		guildChannel = ch.Channel.(discord.GuildChannel)
	}
	return
}

func (s *guildImpl) GetGuildChannels(guildID snowflake.Snowflake, opts ...RequestOpt) (channels []discord.GuildChannel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildChannels.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &channels, opts...)
	return
}

func (s *guildImpl) UpdateChannelPositions(guildID snowflake.Snowflake, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateChannelPositions.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, guildChannelPositionUpdates, nil, opts...)
}

func (s *guildImpl) GetRoles(guildID snowflake.Snowflake, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRoles.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &roles, opts...)
	return
}

func (s *guildImpl) GetRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) (role []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &role, opts...)
	return
}

func (s *guildImpl) CreateRole(guildID snowflake.Snowflake, createRole discord.RoleCreate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateRole.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, createRole, &role, opts...)
	return
}

func (s *guildImpl) UpdateRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, roleUpdate, &role, opts...)
	return
}

func (s *guildImpl) UpdateRolePositions(guildID snowflake.Snowflake, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRolePositions.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, rolePositionUpdates, &roles, opts...)
	return
}

func (s *guildImpl) DeleteRole(guildID snowflake.Snowflake, roleID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteRole.Compile(nil, guildID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetBans(guildID snowflake.Snowflake, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) (bans []discord.Ban, err error) {
	values := route.QueryValues{}
	if before != "" {
		values["before"] = before
	}
	if after != "" {
		values["after"] = after
	}
	if limit != 0 {
		values["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBans.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &bans, opts...)
	return
}

func (s *guildImpl) GetBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) (ban *discord.Ban, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBan.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &ban, opts...)
	return
}

func (s *guildImpl) AddBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, deleteMessageDays int, opts ...RequestOpt) error {
	compiledRoute, err := route.AddBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, discord.AddBan{DeleteMessageDays: deleteMessageDays}, nil, opts...)
}

func (s *guildImpl) DeleteBan(guildID snowflake.Snowflake, userID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetIntegrations(guildID snowflake.Snowflake, opts ...RequestOpt) (integrations []discord.Integration, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetIntegrations.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &integrations, opts...)
	return
}

func (s *guildImpl) DeleteIntegration(guildID snowflake.Snowflake, integrationID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteIntegration.Compile(nil, guildID, integrationID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetAllWebhooks(guildID snowflake.Snowflake, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildWebhooks.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *guildImpl) GetEmojis(guildID snowflake.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *guildImpl) GetAuditLog(guildID snowflake.Snowflake, userID snowflake.Snowflake, actionType discord.AuditLogEvent, before snowflake.Snowflake, limit int, opts ...RequestOpt) (auditLog *discord.AuditLog, err error) {
	values := route.QueryValues{}
	if userID != "" {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != "" {
		values["before"] = guildID
	}
	if limit != 0 {
		values["limit"] = limit
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuditLogs.Compile(values, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &auditLog, opts...)
	return
}
