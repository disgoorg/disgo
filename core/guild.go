package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Guild struct {
	discord.Guild
	Bot *Bot
}

// Update updates the current Guild
func (g *Guild) Update(updateGuild discord.GuildUpdate, opts ...rest.RequestOpt) (*Guild, error) {
	guild, err := g.Bot.RestServices.GuildService().UpdateGuild(g.ID, updateGuild, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateGuild(*guild, CacheStrategyNoWs), nil
}

// Delete deletes the current Guild
func (g *Guild) Delete(opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteGuild(g.ID, opts...)
}

// PublicRole returns the @everyone Role
func (g *Guild) PublicRole() *Role {
	return g.Bot.Caches.RoleCache().Get(g.ID, g.ID)
}

// CreateRole allows you to create a new Role
func (g *Guild) CreateRole(roleCreate discord.RoleCreate, opts ...rest.RequestOpt) (*Role, error) {
	role, err := g.Bot.RestServices.GuildService().CreateRole(g.ID, roleCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// UpdateRole allows you to update a Role
func (g *Guild) UpdateRole(roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, error) {
	role, err := g.Bot.RestServices.GuildService().UpdateRole(g.ID, roleID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// DeleteRole allows you to delete a Role
func (g *Guild) DeleteRole(roleID discord.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteRole(g.ID, roleID, opts...)
}

// Roles return all Role(s) in this Guild
func (g *Guild) Roles() []*Role {
	return g.Bot.Caches.RoleCache().GuildAll(g.ID)
}

// RoleCache return all Role(s) in this Guild
func (g *Guild) RoleCache() map[discord.Snowflake]*Role {
	return g.Bot.Caches.RoleCache().GuildCache(g.ID)
}

// CreateEmoji allows you to create a new Emoji
func (g *Guild) CreateEmoji(emojiCreate discord.EmojiCreate, opts ...rest.RequestOpt) (*Emoji, error) {
	emoji, err := g.Bot.RestServices.EmojiService().CreateEmoji(g.ID, emojiCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateEmoji(g.ID, *emoji, CacheStrategyNoWs), nil
}

// UpdateEmoji allows you to update an Emoji
func (g *Guild) UpdateEmoji(emojiID discord.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...rest.RequestOpt) (*Emoji, error) {
	emoji, err := g.Bot.RestServices.EmojiService().UpdateEmoji(g.ID, emojiID, emojiUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateEmoji(g.ID, *emoji, CacheStrategyNoWs), nil
}

// DeleteEmoji allows you to delete an Emoji
func (g *Guild) DeleteEmoji(emojiID discord.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.EmojiService().DeleteEmoji(g.ID, emojiID, opts...)
}

// CreateSticker allows you to create a new Sticker
func (g *Guild) CreateSticker(stickerCreate discord.StickerCreate, opts ...rest.RequestOpt) (*Sticker, error) {
	sticker, err := g.Bot.RestServices.StickerService().CreateSticker(g.ID, stickerCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

// UpdateSticker allows you to update a Sticker
func (g *Guild) UpdateSticker(stickerID discord.Snowflake, stickerUpdate discord.StickerUpdate, opts ...rest.RequestOpt) (*Sticker, error) {
	sticker, err := g.Bot.RestServices.StickerService().UpdateSticker(g.ID, stickerID, stickerUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

// DeleteSticker allows you to delete a Sticker
func (g *Guild) DeleteSticker(stickerID discord.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.StickerService().DeleteSticker(g.ID, stickerID, opts...)
}

// SelfMember returns the Member for the current logged-in User for this Guild
func (g *Guild) SelfMember() *Member {
	return g.Bot.Caches.MemberCache().Get(g.ID, g.Bot.ClientID)
}

// Leave leaves the Guild
func (g *Guild) Leave(opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.UserService().LeaveGuild(g.ID, opts...)
}

// Disconnect sends an GatewayCommand to disconnect from this Guild
func (g *Guild) Disconnect() error {
	return g.Bot.AudioController.Disconnect(g.ID)
}

func (g *Guild) RequestMembers(userIDs ...discord.Snowflake) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembers(g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQuery(query string, limit int) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQuery(g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilter(memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilter(g.ID, memberFilterFunc)
}

func (g *Guild) RequestMembersCtx(ctx context.Context, userIDs ...discord.Snowflake) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersCtx(ctx, g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQueryCtx(ctx context.Context, query string, limit int) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQueryCtx(ctx, g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilterCtx(ctx context.Context, memberFilterFunc func(member *Member) bool) ([]*Member, error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilterCtx(ctx, g.ID, memberFilterFunc)
}

func (g *Guild) RequestMembersChan(userIDs []discord.Snowflake) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersChan(g.ID, userIDs...)
}

func (g *Guild) RequestMembersWithQueryChan(query string, limit int) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithQueryChan(g.ID, query, limit)
}

func (g *Guild) RequestMembersWithFilterChan(memberFilterFunc func(member *Member) bool) (<-chan *Member, func(), error) {
	return g.Bot.MemberChunkingManager.RequestMembersWithFilterChan(g.ID, memberFilterFunc)
}

// GetMember returns the specific Member for this Guild
func (g *Guild) GetMember(userID discord.Snowflake) *Member {
	return g.Bot.Caches.MemberCache().Get(g.ID, userID)
}

// AddMember adds a member to the Guild with the oauth2 access token
func (g *Guild) AddMember(userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...rest.RequestOpt) (*Member, error) {
	member, err := g.Bot.RestServices.GuildService().AddMember(g.ID, userID, memberAdd, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// UpdateMember updates an existing member of the Guild
func (g *Guild) UpdateMember(userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, error) {
	member, err := g.Bot.RestServices.GuildService().UpdateMember(g.ID, userID, memberUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// KickMember kicks an existing member from the Guild
func (g *Guild) KickMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().RemoveMember(g.ID, userID, opts...)
}

// BanMember bans a Member from the Guild
func (g *Guild) BanMember(userID discord.Snowflake, deleteMessageDays int, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().AddBan(g.ID, userID, deleteMessageDays, opts...)
}

// UnbanMember unbans a Member from the Guild
func (g *Guild) UnbanMember(userID discord.Snowflake, opts ...rest.RequestOpt) error {
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
func (g *Guild) GetBan(userID discord.Snowflake, opts ...rest.RequestOpt) (*Ban, error) {
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
func (g *Guild) GetIntegrations(opts ...rest.RequestOpt) ([]*Integration, error) {
	guildIntegrations, err := g.Bot.RestServices.GuildService().GetIntegrations(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	integrations := make([]*Integration, len(guildIntegrations))
	for i, guildIntegration := range guildIntegrations {
		integrations[i] = g.Bot.EntityBuilder.CreateIntegration(g.ID, guildIntegration, CacheStrategyNoWs)
	}
	return integrations, nil
}

// DeleteIntegration deletes a specific Integration from the Guild. Requires PermissionManageServer
func (g *Guild) DeleteIntegration(integrationID discord.Snowflake, opts ...rest.RequestOpt) error {
	return g.Bot.RestServices.GuildService().DeleteIntegration(g.ID, integrationID, opts...)
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (g *Guild) GetGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (ApplicationCommand, error) {
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
func (g *Guild) UpdateGuildCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (ApplicationCommand, error) {
	command, err := g.Bot.RestServices.ApplicationService().UpdateGuildCommand(g.Bot.ApplicationID, g.ID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Bot.EntityBuilder.CreateApplicationCommand(command), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) DeleteGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) error {
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
func (g *Guild) GetGuildCommandPermissions(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
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
func (g *Guild) SetGuildCommandPermissions(commandID discord.Snowflake, permissions []discord.ApplicationCommandPermission, opts ...rest.RequestOpt) (*ApplicationCommandPermissions, error) {
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
