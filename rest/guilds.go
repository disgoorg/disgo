package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Guilds = (*guildImpl)(nil)

func NewGuilds(client Client) Guilds {
	return &guildImpl{client: client}
}

type Guilds interface {
	GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (*discord.RestGuild, error)
	GetGuildPreview(guildID snowflake.ID, opts ...RequestOpt) (*discord.GuildPreview, error)
	CreateGuild(guildCreate discord.GuildCreate, opts ...RequestOpt) (*discord.RestGuild, error)
	UpdateGuild(guildID snowflake.ID, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (*discord.RestGuild, error)
	DeleteGuild(guildID snowflake.ID, opts ...RequestOpt) error

	CreateGuildChannel(guildID snowflake.ID, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (discord.GuildChannel, error)
	GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) ([]discord.GuildChannel, error)
	UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error

	GetRoles(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Role, error)
	GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) ([]discord.Role, error)
	CreateRole(guildID snowflake.ID, createRole discord.RoleCreate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (*discord.Role, error)
	UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) ([]discord.Role, error)
	DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error

	GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]discord.Ban, error)
	GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (*discord.Ban, error)
	AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDays int, opts ...RequestOpt) error
	DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error

	GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Integration, error)
	DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error

	GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Webhook, error)

	GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, before snowflake.ID, limit int, opts ...RequestOpt) (*discord.AuditLog, error)
}

type guildImpl struct {
	client Client
}

func (s *guildImpl) GetGuild(guildID snowflake.ID, withCounts bool, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
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

func (s *guildImpl) GetGuildPreview(guildID snowflake.ID, opts ...RequestOpt) (guildPreview *discord.GuildPreview, err error) {
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

func (s *guildImpl) UpdateGuild(guildID snowflake.ID, guildUpdate discord.GuildUpdate, opts ...RequestOpt) (guild *discord.RestGuild, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuild.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, guildUpdate, &guild, opts...)
	return
}

func (s *guildImpl) DeleteGuild(guildID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuild.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) CreateGuildChannel(guildID snowflake.ID, guildChannelCreate discord.GuildChannelCreate, opts ...RequestOpt) (guildChannel discord.GuildChannel, err error) {
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

func (s *guildImpl) GetGuildChannels(guildID snowflake.ID, opts ...RequestOpt) (channels []discord.GuildChannel, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildChannels.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &channels, opts...)
	return
}

func (s *guildImpl) UpdateChannelPositions(guildID snowflake.ID, guildChannelPositionUpdates []discord.GuildChannelPositionUpdate, opts ...RequestOpt) error {
	compiledRoute, err := route.UpdateChannelPositions.Compile(nil, guildID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, guildChannelPositionUpdates, nil, opts...)
}

func (s *guildImpl) GetRoles(guildID snowflake.ID, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRoles.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &roles, opts...)
	return
}

func (s *guildImpl) GetRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (role []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &role, opts...)
	return
}

func (s *guildImpl) CreateRole(guildID snowflake.ID, createRole discord.RoleCreate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateRole.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, createRole, &role, opts...)
	return
}

func (s *guildImpl) UpdateRole(guildID snowflake.ID, roleID snowflake.ID, roleUpdate discord.RoleUpdate, opts ...RequestOpt) (role *discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRole.Compile(nil, guildID, roleID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, roleUpdate, &role, opts...)
	return
}

func (s *guildImpl) UpdateRolePositions(guildID snowflake.ID, rolePositionUpdates []discord.RolePositionUpdate, opts ...RequestOpt) (roles []discord.Role, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateRolePositions.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, rolePositionUpdates, &roles, opts...)
	return
}

func (s *guildImpl) DeleteRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteRole.Compile(nil, guildID, roleID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetBans(guildID snowflake.ID, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (bans []discord.Ban, err error) {
	values := route.QueryValues{}
	if before != 0 {
		values["before"] = before
	}
	if after != 0 {
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

func (s *guildImpl) GetBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (ban *discord.Ban, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBan.Compile(nil, guildID, userID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &ban, opts...)
	return
}

func (s *guildImpl) AddBan(guildID snowflake.ID, userID snowflake.ID, deleteMessageDays int, opts ...RequestOpt) error {
	compiledRoute, err := route.AddBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, discord.AddBan{DeleteMessageDays: deleteMessageDays}, nil, opts...)
}

func (s *guildImpl) DeleteBan(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteBan.Compile(nil, guildID, userID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetIntegrations(guildID snowflake.ID, opts ...RequestOpt) (integrations []discord.Integration, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetIntegrations.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &integrations, opts...)
	return
}

func (s *guildImpl) DeleteIntegration(guildID snowflake.ID, integrationID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteIntegration.Compile(nil, guildID, integrationID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}

func (s *guildImpl) GetAllWebhooks(guildID snowflake.ID, opts ...RequestOpt) (webhooks []discord.Webhook, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildWebhooks.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &webhooks, opts...)
	return
}

func (s *guildImpl) GetEmojis(guildID snowflake.ID, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *guildImpl) GetAuditLog(guildID snowflake.ID, userID snowflake.ID, actionType discord.AuditLogEvent, before snowflake.ID, limit int, opts ...RequestOpt) (auditLog *discord.AuditLog, err error) {
	values := route.QueryValues{}
	if userID != 0 {
		values["user_id"] = userID
	}
	if actionType != 0 {
		values["action_type"] = actionType
	}
	if before != 0 {
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
