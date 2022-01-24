package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

type Guild struct {
	discord.Guild
	Bot *Bot
}

// Update updates this Guild
func (g *Guild) Update(updateGuild discord.GuildUpdate, opts ...rest.RequestOpt) (*Guild, error) {
	guild, err := g.Bot.RestServices.GuildService().UpdateGuild(g.ID, updateGuild, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuild(*guild, CacheStrategyNoWs), nil
}

// Delete deletes the this Guild
func (g *Guild) Delete(opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteGuild(g.ID, opts...)
}

// CreateChannel creates a new GuildChannel in this Guild
func (g *Guild) CreateChannel(guildChannelCreate discord.GuildChannelCreate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := g.Bot.RestServices.GuildService().CreateChannel(g.ID, guildChannelCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateChannel(channel, CacheStrategyNoWs).(GuildChannel), nil
}

// UpdateChannel updates a GuildChannel in this Guild
func (g *Guild) UpdateChannel(channelID snowflake.Snowflake, guildChannelUpdate discord.GuildChannelUpdate, opts ...rest.RequestOpt) (GuildChannel, error) {
	channel, err := g.Bot.RestServices.ChannelService().UpdateChannel(channelID, guildChannelUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateChannel(channel, CacheStrategyNoWs).(GuildChannel), nil
}

// DeleteChannel deletes a GuildChannel in this Guild
func (g *Guild) DeleteChannel(channelID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.ChannelService().DeleteChannel(channelID, opts...)
}

// PublicRole returns the @everyone Role
func (g *Guild) PublicRole() *Role {
	return g.Bot.Caches.Roles().Get(g.ID, g.ID)
}

// CreateRole creates a new Role with the properties provided in discord.RoleCreate
func (g *Guild) CreateRole(roleCreate discord.RoleCreate, opts ...rest.RequestOpt) (*Role, error) {
	role, err := g.Bot.RestServices.GuildService().CreateRole(g.ID, roleCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// UpdateRole updates a Role with the properties provided in discord.RoleUpdate
func (g *Guild) UpdateRole(roleID snowflake.Snowflake, roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, error) {
	role, err := g.Bot.RestServices.GuildService().UpdateRole(g.ID, roleID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// DeleteRole deletes a Role
func (g *Guild) DeleteRole(roleID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteRole(g.ID, roleID, opts...)
}

// Roles returns all Role(s) in this Guild
func (g *Guild) Roles() []*Role {
	return g.Bot.Caches.Roles().GuildAll(g.ID)
}

// RoleCache returns all Role(s) in this Guild
func (g *Guild) RoleCache() map[snowflake.Snowflake]*Role {
	return g.Bot.Caches.Roles().GuildCache(g.ID)
}

func (g *Guild) CreateGuildScheduledEvent(guildScheduledEventCreate discord.GuildScheduledEventCreate, opts ...rest.RequestOpt) (*GuildScheduledEvent, error) {
	guildScheduledEvent, err := g.Bot.RestServices.GuildScheduledEventService().CreateGuildScheduledEvent(g.ID, guildScheduledEventCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildScheduledEvent(*guildScheduledEvent, CacheStrategyNoWs), nil
}

func (g *Guild) UpdateGuildScheduledEvent(guildScheduledEventID snowflake.Snowflake, guildScheduledEventUpdate discord.GuildScheduledEventUpdate, opts ...rest.RequestOpt) (*GuildScheduledEvent, error) {
	guildScheduledEvent, err := g.Bot.RestServices.GuildScheduledEventService().UpdateGuildScheduledEvent(g.ID, guildScheduledEventID, guildScheduledEventUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildScheduledEvent(*guildScheduledEvent, CacheStrategyNoWs), nil
}

func (g *Guild) DeleteGuildScheduledEvent(guildScheduledEventID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildScheduledEventService().DeleteGuildScheduledEvent(g.ID, guildScheduledEventID, opts...)
}

func (g *Guild) GetGuildScheduledEventUsers(guildScheduledEventID snowflake.Snowflake, limit int, withMember bool, before snowflake.Snowflake, after snowflake.Snowflake, opts ...rest.RequestOpt) ([]*GuildScheduledEventUser, error) {
	users, err := g.Bot.RestServices.GuildScheduledEventService().GetGuildScheduledEventUsers(g.ID, guildScheduledEventID, limit, withMember, before, after, opts...)
	if err != nil {
		return nil, err
	}
	eventUsers := make([]*GuildScheduledEventUser, len(users))
	for i := range users {
		eventUsers[i] = g.Bot.EntityBuilder.CreateGuildScheduledEventUser(g.ID, users[i], CacheStrategyNoWs)
	}
	return eventUsers, nil
}

// CreateEmoji creates a new Emoji with the properties provided in discord.EmojiCreate
func (g *Guild) CreateEmoji(emojiCreate discord.EmojiCreate, opts ...rest.RequestOpt) (*Emoji, error) {
	emoji, err := g.Bot.RestServices.EmojiService().CreateEmoji(g.ID, emojiCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateEmoji(g.ID, *emoji, CacheStrategyNoWs), nil
}

// UpdateEmoji creates an Emoji with the properties provided in discord.EmojiUpdate
func (g *Guild) UpdateEmoji(emojiID snowflake.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...rest.RequestOpt) (*Emoji, error) {
	emoji, err := g.Bot.RestServices.EmojiService().UpdateEmoji(g.ID, emojiID, emojiUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateEmoji(g.ID, *emoji, CacheStrategyNoWs), nil
}

// DeleteEmoji deletes an Emoji
func (g *Guild) DeleteEmoji(emojiID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.EmojiService().DeleteEmoji(g.ID, emojiID, opts...)
}

// CreateSticker creates a new Sticker with the properties provided in discord.StickerCreate
func (g *Guild) CreateSticker(stickerCreate discord.StickerCreate, opts ...rest.RequestOpt) (*Sticker, error) {
	sticker, err := g.Bot.RestServices.StickerService().CreateSticker(g.ID, stickerCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

// UpdateSticker updates a Sticker with the properties provided in discord.StickerCreate
func (g *Guild) UpdateSticker(stickerID snowflake.Snowflake, stickerUpdate discord.StickerUpdate, opts ...rest.RequestOpt) (*Sticker, error) {
	sticker, err := g.Bot.RestServices.StickerService().UpdateSticker(g.ID, stickerID, stickerUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

// DeleteSticker deletes a Sticker
func (g *Guild) DeleteSticker(stickerID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.StickerService().DeleteSticker(g.ID, stickerID, opts...)
}

// SelfMember returns the Member for the current logged-in User for this Guild
func (g *Guild) SelfMember() *Member {
	return g.Bot.Caches.Members().Get(g.ID, g.Bot.ClientID)
}

// Leave leaves this Guild
func (g *Guild) Leave(opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.UserService().LeaveGuild(g.ID, opts...)
}

// Disconnect sends a GatewayCommand to disconnect from this Guild
func (g *Guild) Disconnect(ctx context.Context) error {
	return g.Bot.AudioController.Disconnect(ctx, g.ID)
}

func (g *Guild) RequestMembers(userIDs ...snowflake.Snowflake) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembers(g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQuery(query string, limit int) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQuery(g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilter(memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilter(g.ID, memberFilterFunc)
}

func (g *Guild) RequestMembersCtx(ctx context.Context, userIDs ...snowflake.Snowflake) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersCtx(ctx, g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQueryCtx(ctx context.Context, query string, limit int) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQueryCtx(ctx, g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilterCtx(ctx context.Context, memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilterCtx(ctx, g.ID, memberFilterFunc)
}

func (g *Guild) RequestMembersChan(userIDs []snowflake.Snowflake) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersChan(g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQueryChan(query string, limit int) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQueryChan(g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilterChan(memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilterChan(g.ID, memberFilterFunc)
}

// GetMember returns the specific Member for this Guild
func (g *Guild) GetMember(userID snowflake.Snowflake) *Member {
	return g.Bot.Caches.Members().Get(g.ID, userID)
}

// AddMember adds a member to the Guild with the oauth2 access token
func (g *Guild) AddMember(userID snowflake.Snowflake, memberAdd discord.MemberAdd, opts ...rest.RequestOpt) (*Member, error) {
	member, err := g.Bot.RestServices.GuildService().AddMember(g.ID, userID, memberAdd, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// UpdateMember updates an existing member of the Guild
func (g *Guild) UpdateMember(userID snowflake.Snowflake, memberUpdate discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, error) {
	member, err := g.Bot.RestServices.GuildService().UpdateMember(g.ID, userID, memberUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// KickMember kicks an existing member from the Guild
func (g *Guild) KickMember(userID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().RemoveMember(g.ID, userID, opts...)
}

// BanMember bans a Member from the Guild
func (g *Guild) BanMember(userID snowflake.Snowflake, deleteMessageDays int, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().AddBan(g.ID, userID, deleteMessageDays, opts...)
}

// UnbanMember unbans a Member from the Guild
func (g *Guild) UnbanMember(userID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteBan(g.ID, userID, opts...)
}

// GetBans fetches all bans for this Guild
func (g *Guild) GetBans(opts ...rest.RequestOpt) ([]*Ban, error) {
	guildBans, err := g.Bot.RestServices.GuildService().GetBans(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	bans := make([]*Ban, len(guildBans))
	for i, guildBan := range guildBans {
		bans[i] = g.Bot.EntityBuilder.CreateBan(g.ID, guildBan, CacheStrategyNoWs)
	}
	return bans, nil
}

// GetBan fetches a ban for a User for this Guild
func (g *Guild) GetBan(userID snowflake.Snowflake, opts ...rest.RequestOpt) (*Ban, error) {
	ban, err := g.Bot.RestServices.GuildService().GetBan(g.ID, userID, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateBan(g.ID, *ban, CacheStrategyNoWs), nil
}

// IconURL returns the Icon of a Guild
func (g *Guild) IconURL(size int) *string {
	return discord.FormatAssetURL(route.GuildIcon, g.ID, g.Icon, size)
}

// GetAuditLogs gets AuditLog(s) for this Guild
func (g *Guild) GetAuditLogs(filterOptions AuditLogFilterOptions, opts ...rest.RequestOpt) (*AuditLog, error) {
	auditLog, err := g.Bot.RestServices.AuditLogService().GetAuditLog(g.ID, filterOptions.UserID, filterOptions.ActionType, filterOptions.Before, filterOptions.Limit, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateAuditLog(g.ID, *auditLog, filterOptions, CacheStrategyNoWs), nil
}

// GetIntegrations gets all Integration(s) from the Guild. Requires PermissionManageServer
func (g *Guild) GetIntegrations(opts ...rest.RequestOpt) ([]Integration, error) {
	guildIntegrations, err := g.Bot.RestServices.GuildService().GetIntegrations(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	integrations := make([]Integration, len(guildIntegrations))
	for i := range guildIntegrations {
		integrations[i] = g.Bot.EntityBuilder.CreateIntegration(g.ID, guildIntegrations[i], CacheStrategyNoWs)
	}
	return integrations, nil
}

// DeleteIntegration deletes a specific Integration from the Guild. Requires PermissionManageServer
func (g *Guild) DeleteIntegration(integrationID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteIntegration(g.ID, integrationID, opts...)
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (g *Guild) GetGuildCommand(commandID snowflake.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := g.Bot.RestServices.ApplicationService().GetGuildCommand(g.Bot.ApplicationID, g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommand(command), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (g *Guild) GetGuildCommands(opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := g.Bot.RestServices.ApplicationService().GetGuildCommands(g.Bot.ApplicationID, g.ID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.Bot.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) CreateGuildCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := g.Bot.RestServices.ApplicationService().CreateGuildCommand(g.Bot.ApplicationID, g.ID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommand(command), nil
}

// UpdateGuildCommand edits a specific Guild discord.ApplicationCommand
func (g *Guild) UpdateGuildCommand(commandID snowflake.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := g.Bot.RestServices.ApplicationService().UpdateGuildCommand(g.Bot.ApplicationID, g.ID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommand(command), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) DeleteGuildCommand(commandID snowflake.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.ApplicationService().DeleteGuildCommand(g.Bot.ApplicationID, g.ID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]ApplicationCommand, error) {
	cmds, err := g.Bot.RestServices.ApplicationService().SetGuildCommands(g.Bot.ApplicationID, g.ID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.Bot.EntityBuilder.CreateApplicationCommand(command)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the core.ApplicationCommandPermissions for a all discord.ApplicationCommand(s) in an core.Guild
func (g *Guild) GetGuildCommandsPermissions(opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error) {
	perms, err := g.Bot.RestServices.ApplicationService().GetGuildCommandsPermissions(g.Bot.ApplicationID, g.ID, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.Bot.EntityBuilder.CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand in an core.Guild
func (g *Guild) GetGuildCommandPermissions(commandID snowflake.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	permissions, err := g.Bot.RestServices.ApplicationService().GetGuildCommandPermissions(g.Bot.ApplicationID, g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommandPermissions(*permissions), nil
}

// SetGuildCommandsPermissions sets the discord.ApplicationCommandPermissions for all discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommandsPermissions(commandPermissions []discord.ApplicationCommandPermissionsSet, opts ...rest.RequestOpt) ([]*ApplicationCommandPermissions, error) {
	perms, err := g.Bot.RestServices.ApplicationService().SetGuildCommandsPermissions(g.Bot.ApplicationID, g.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*ApplicationCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.Bot.EntityBuilder.CreateApplicationCommandPermissions(permission)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the core.ApplicationCommandPermissions for a specific discord.ApplicationCommand
func (g *Guild) SetGuildCommandPermissions(commandID snowflake.Snowflake, permissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
	perms, err := g.Bot.RestServices.ApplicationService().SetGuildCommandPermissions(g.Bot.ApplicationID, g.ID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommandPermissions(*perms), nil
}

// GetTemplates gets a specific GuildTemplate
func (g *Guild) GetTemplates(opts ...rest.RequestOpt) ([]*GuildTemplate, error) {
	guildTemplates, err := g.Bot.RestServices.GuildTemplateService().GetGuildTemplates(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	templates := make([]*GuildTemplate, len(guildTemplates))
	for i, guildTemplate := range guildTemplates {
		templates[i] = g.Bot.EntityBuilder.CreateGuildTemplate(guildTemplate, CacheStrategyNoWs)
	}
	return templates, nil
}

// CreateTemplate creates a new GuildTemplate
func (g *Guild) CreateTemplate(guildTemplateCreate discord.GuildTemplateCreate, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := g.Bot.RestServices.GuildTemplateService().CreateGuildTemplate(g.ID, guildTemplateCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// SyncTemplate syncs the current Guild status to an existing GuildTemplate
func (g *Guild) SyncTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := g.Bot.RestServices.GuildTemplateService().SyncGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// UpdateTemplate updates a specific GuildTemplate
func (g *Guild) UpdateTemplate(templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := g.Bot.RestServices.GuildTemplateService().UpdateGuildTemplate(g.ID, templateCode, guildTemplateUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// DeleteTemplate deletes a specific GuildTemplate
func (g *Guild) DeleteTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, error) {
	guildTemplate, err := g.Bot.RestServices.GuildTemplateService().DeleteGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}
