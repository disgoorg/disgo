package core

import (
	"context"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Guild struct {
	discord.Guild
	Disgo Disgo
}

// Update updates the current Guild
func (g *Guild) Update(ctx context.Context, updateGuild discord.GuildUpdate) (*Guild, rest.Error) {
	guild, err := g.Disgo.RestServices().GuildService().UpdateGuild(ctx, g.ID, updateGuild)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

// Delete deletes the current Guild
func (g *Guild) Delete(ctx context.Context, ) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteGuild(ctx, g.ID)
}

// PublicRole returns the @everyone Role
func (g *Guild) PublicRole() *Role {
	return g.Disgo.Cache().RoleCache().Get(g.ID)
}

// CreateRole allows you to create a new Role
func (g *Guild) CreateRole(ctx context.Context, roleCreate discord.RoleCreate) (*Role, rest.Error) {
	role, err := g.Disgo.RestServices().GuildService().CreateRole(ctx, g.ID, roleCreate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// UpdateRole allows you to update a Role
func (g *Guild) UpdateRole(ctx context.Context, roleID discord.Snowflake, roleUpdate discord.RoleUpdate) (*Role, rest.Error) {
	role, err := g.Disgo.RestServices().GuildService().UpdateRole(ctx, g.ID, roleID, roleUpdate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// DeleteRole allows you to delete a Role
func (g *Guild) DeleteRole(ctx context.Context, roleID discord.Snowflake) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteRole(ctx, g.ID, roleID)
}

// Roles return all Role(s) in this Guild
func (g *Guild) Roles() []*Role {
	return g.Disgo.Cache().RoleCache().All(g.ID)
}

// SelfMember returns the Member for the current logged-in User for this Guild
func (g *Guild) SelfMember() *Member {
	return g.Disgo.Cache().MemberCache().Get(g.ID, g.Disgo.ClientID())
}

// Leave leaves the Guild
func (g *Guild) Leave(ctx context.Context, ) rest.Error {
	return g.Disgo.RestServices().UserService().LeaveGuild(ctx, g.ID)
}

// Disconnect sends an GatewayCommand to disconnect from this Guild
func (g *Guild) Disconnect() error {
	return g.Disgo.AudioController().Disconnect(g.ID)
}

// GetMember returns the specific Member for this Guild
func (g *Guild) GetMember(userID discord.Snowflake) *Member {
	return g.Disgo.Cache().MemberCache().Get(g.ID, userID)
}

// AddMember adds a member to the Guild with the oauth2 access token
func (g *Guild) AddMember(ctx context.Context, userID discord.Snowflake, memberAdd discord.MemberAdd) (*Member, rest.Error) {
	member, err := g.Disgo.RestServices().GuildService().AddMember(ctx, g.ID, userID, memberAdd)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// UpdateMember updates an existing member of the Guild
func (g *Guild) UpdateMember(ctx context.Context, userID discord.Snowflake, memberUpdate discord.MemberUpdate) (*Member, rest.Error) {
	member, err := g.Disgo.RestServices().GuildService().UpdateMember(ctx, g.ID, userID, memberUpdate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// KickMember kicks an existing member from the Guild
func (g *Guild) KickMember(ctx context.Context, userID discord.Snowflake, reason string) rest.Error {
	return g.Disgo.RestServices().GuildService().RemoveMember(ctx, g.ID, userID, reason)
}

// BanMember bans a Member from the Guild
func (g *Guild) BanMember(ctx context.Context, userID discord.Snowflake, reason string, deleteMessageDays int) rest.Error {
	return g.Disgo.RestServices().GuildService().AddBan(ctx, g.ID, userID, reason, deleteMessageDays)
}

// UnbanMember unbans a Member from the Guild
func (g *Guild) UnbanMember(ctx context.Context, userID discord.Snowflake) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteBan(ctx, g.ID, userID)
}

// GetBans fetches all bans for this Guild
func (g *Guild) GetBans(ctx context.Context) ([]*Ban, rest.Error) {
	guildBans, err := g.Disgo.RestServices().GuildService().GetBans(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	bans := make([]*Ban, len(guildBans))
	for i, guildBan := range guildBans {
		bans[i] = g.Disgo.EntityBuilder().CreateBan(g.ID, guildBan, CacheStrategyNoWs)
	}
	return bans, nil
}

// GetBan fetches a ban for a User for this Guild
func (g *Guild) GetBan(ctx context.Context, userID discord.Snowflake) (*Ban, rest.Error) {
	ban, err := g.Disgo.RestServices().GuildService().GetBan(ctx, g.ID, userID)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateBan(g.ID, *ban, CacheStrategyNoWs), nil
}

// IconURL returns the Icon of a Guild
func (g *Guild) IconURL(size int) *string {
	if g.Icon == nil {
		return nil
	}
	animated := strings.HasPrefix(*g.Icon, "a_")
	format := route.PNG
	if animated {
		format = route.GIF
	}
	compiledRoute, err := route.GuildIcon.Compile(nil, format, size, g.ID.String(), *g.Icon)
	if err != nil {
		return nil
	}
	u := compiledRoute.URL()
	return &u
}

// GetAuditLogs gets AuditLog(s) for this Guild
func (g *Guild) GetAuditLogs(ctx context.Context, filterOptions AuditLogFilterOptions) (*AuditLog, rest.Error) {
	auditLog, err := g.Disgo.RestServices().AuditLogService().GetAuditLog(ctx, g.ID, filterOptions.UserID, filterOptions.ActionType, filterOptions.Before, filterOptions.Limit)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateAuditLog(g.ID, *auditLog, filterOptions, CacheStrategyNoWs), nil
}

// GetIntegrations gets all Integration(s) from the Guild. Requires PermissionManageServer
func (g *Guild) GetIntegrations(ctx context.Context, ) ([]*Integration, rest.Error) {
	guildIntegrations, err := g.Disgo.RestServices().GuildService().GetIntegrations(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	integrations := make([]*Integration, len(guildIntegrations))
	for i, guildIntegration := range guildIntegrations {
		integrations[i] = g.Disgo.EntityBuilder().CreateIntegration(g.ID, guildIntegration, CacheStrategyNoWs)
	}
	return integrations, nil
}

// DeleteIntegration deletes a specific Integration from the Guild. Requires PermissionManageServer
func (g *Guild) DeleteIntegration(ctx context.Context, integrationID discord.Snowflake) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteIntegration(ctx, g.ID, integrationID)
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (g *Guild) GetGuildCommand(ctx context.Context, commandID discord.Snowflake) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().GetGuildCommand(ctx, g.Disgo.ApplicationID(), g.ID, commandID)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (g *Guild) GetGuildCommands(ctx context.Context, ) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.Disgo.RestServices().ApplicationService().GetGuildCommands(ctx, g.Disgo.ApplicationID(), g.ID)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.Disgo.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) CreateGuildCommand(ctx context.Context, commandCreate discord.ApplicationCommandCreate) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().CreateGuildCommand(ctx, g.Disgo.ApplicationID(), g.ID, commandCreate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// UpdateGuildCommand edits a specific Guild discord.ApplicationCommand
func (g *Guild) UpdateGuildCommand(ctx context.Context, commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().UpdateGuildCommand(ctx, g.Disgo.ApplicationID(), g.ID, commandID, commandUpdate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) DeleteGuildCommand(ctx context.Context, commandID discord.Snowflake) rest.Error {
	return g.Disgo.RestServices().ApplicationService().DeleteGuildCommand(ctx, g.Disgo.ApplicationID(), g.ID, commandID)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommands(ctx context.Context, commandCreates ...discord.ApplicationCommandCreate) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.Disgo.RestServices().ApplicationService().SetGuildCommands(ctx, g.Disgo.ApplicationID(), g.ID, commandCreates...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.Disgo.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the api.GuildCommandPermissions for a all discord.ApplicationCommand(s) in an api.Guild
func (g *Guild) GetGuildCommandsPermissions(ctx context.Context, ) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().GetGuildCommandsPermissions(ctx, g.Disgo.ApplicationID(), g.ID)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.Disgo.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the api.GuildCommandPermissions for a specific discord.ApplicationCommand in an api.Guild
func (g *Guild) GetGuildCommandPermissions(ctx context.Context, commandID discord.Snowflake) (*GuildCommandPermissions, rest.Error) {
	permissions, err := g.Disgo.RestServices().ApplicationService().GetGuildCommandPermissions(ctx, g.Disgo.ApplicationID(), g.ID, commandID)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// SetGuildCommandsPermissions sets the discord.GuildCommandPermissions for all discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommandsPermissions(ctx context.Context, commandPermissions ...discord.GuildCommandPermissionsSet) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().SetGuildCommandsPermissions(ctx, g.Disgo.ApplicationID(), g.ID, commandPermissions...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.Disgo.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific discord.ApplicationCommand
func (g *Guild) SetGuildCommandPermissions(ctx context.Context, commandID discord.Snowflake, permissions ...discord.CommandPermission) (*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().SetGuildCommandPermissions(ctx, g.Disgo.ApplicationID(), g.ID, commandID, permissions...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommandPermissions(*perms, CacheStrategyNoWs), nil
}

// GetTemplates gets a specific GuildTemplate
func (g *Guild) GetTemplates(ctx context.Context, ) ([]*GuildTemplate, rest.Error) {
	guildTemplates, err := g.Disgo.RestServices().GuildTemplateService().GetGuildTemplates(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	templates := make([]*GuildTemplate, len(guildTemplates))
	for i, guildTemplate := range guildTemplates {
		templates[i] = g.Disgo.EntityBuilder().CreateGuildTemplate(guildTemplate, CacheStrategyNoWs)
	}
	return templates, nil
}

// CreateTemplate creates a new GuildTemplate
func (g *Guild) CreateTemplate(ctx context.Context, guildTemplateCreate discord.GuildTemplateCreate) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().CreateGuildTemplate(ctx, g.ID, guildTemplateCreate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// SyncTemplate syncs the current Guild status to an existing GuildTemplate
func (g *Guild) SyncTemplate(ctx context.Context, templateCode string) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().SyncGuildTemplate(ctx, g.ID, templateCode)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// UpdateTemplate updates a specific GuildTemplate
func (g *Guild) UpdateTemplate(ctx context.Context, templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().UpdateGuildTemplate(ctx, g.ID, templateCode, guildTemplateUpdate)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// DeleteTemplate deletes a specific GuildTemplate
func (g *Guild) DeleteTemplate(ctx context.Context, templateCode string) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().DeleteGuildTemplate(ctx, g.ID, templateCode)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}
