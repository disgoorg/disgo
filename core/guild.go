package core

import (
	"strings"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Guild interface {
	Entity

}

type guildImpl struct {
	discord.Guild
	disgo Disgo
}

// Update updates the current Guild
func (g *guildImpl) Update(updateGuild discord.GuildUpdate, opts ...rest.RequestOpt) (Guild, rest.Error) {
	guild, err := g.disgo.RestServices().GuildService().UpdateGuild(g.ID, updateGuild, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

// Delete deletes the current Guild
func (g *guildImpl) Delete(opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().DeleteGuild(g.ID, opts...)
}

// PublicRole returns the @everyone Role
func (g *guildImpl) PublicRole() *Role {
	return g.disgo.Caches().RoleCache().Get(g.ID)
}

// CreateRole allows you to create a new Role
func (g *guildImpl) CreateRole(roleCreate discord.RoleCreate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := g.disgo.RestServices().GuildService().CreateRole(g.ID, roleCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// UpdateRole allows you to update a Role
func (g *guildImpl) UpdateRole(roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := g.disgo.RestServices().GuildService().UpdateRole(g.ID, roleID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// DeleteRole allows you to delete a Role
func (g *guildImpl) DeleteRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().DeleteRole(g.ID, roleID, opts...)
}

// Roles return all Role(s) in this Guild
func (g *guildImpl) Roles() []*Role {
	return g.disgo.Caches().RoleCache().All(g.ID)
}

// SelfMember returns the Member for the current logged-in User for this Guild
func (g *guildImpl) SelfMember() *Member {
	return g.disgo.Caches().MemberCache().Get(g.ID, g.disgo.ClientID())
}

// Leave leaves the Guild
func (g *guildImpl) Leave(opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().UserService().LeaveGuild(g.ID, opts...)
}

// Disconnect sends an GatewayCommand to disconnect from this Guild
func (g *guildImpl) Disconnect() error {
	return g.disgo.AudioController().Disconnect(g.ID)
}

// GetMember returns the specific Member for this Guild
func (g *guildImpl) GetMember(userID discord.Snowflake) *Member {
	return g.disgo.Caches().MemberCache().Get(g.ID, userID)
}

// AddMember adds a member to the Guild with the oauth2 access token
func (g *guildImpl) AddMember(userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := g.disgo.RestServices().GuildService().AddMember(g.ID, userID, memberAdd, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// UpdateMember updates an existing member of the Guild
func (g *guildImpl) UpdateMember(userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := g.disgo.RestServices().GuildService().UpdateMember(g.ID, userID, memberUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// KickMember kicks an existing member from the Guild
func (g *guildImpl) KickMember(userID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().RemoveMember(g.ID, userID, opts...)
}

// BanMember bans a Member from the Guild
func (g *guildImpl) BanMember(userID discord.Snowflake, deleteMessageDays int, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().AddBan(g.ID, userID, deleteMessageDays, opts...)
}

// UnbanMember unbans a Member from the Guild
func (g *guildImpl) UnbanMember(userID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().DeleteBan(g.ID, userID, opts...)
}

// GetBans fetches all bans for this Guild
func (g *guildImpl) GetBans(opts ...rest.RequestOpt) ([]*Ban, rest.Error) {
	guildBans, err := g.disgo.RestServices().GuildService().GetBans(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	bans := make([]*Ban, len(guildBans))
	for i, guildBan := range guildBans {
		bans[i] = g.disgo.EntityBuilder().CreateBan(g.ID, guildBan, CacheStrategyNoWs)
	}
	return bans, nil
}

// GetBan fetches a ban for a User for this Guild
func (g *guildImpl) GetBan(userID discord.Snowflake, opts ...rest.RequestOpt) (*Ban, rest.Error) {
	ban, err := g.disgo.RestServices().GuildService().GetBan(g.ID, userID, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateBan(g.ID, *ban, CacheStrategyNoWs), nil
}

// IconURL returns the Icon of a Guild
func (g *guildImpl) IconURL(size int) *string {
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
func (g *guildImpl) GetAuditLogs(filterOptions AuditLogFilterOptions, opts ...rest.RequestOpt) (*AuditLog, rest.Error) {
	auditLog, err := g.disgo.RestServices().AuditLogService().GetAuditLog(g.ID, filterOptions.UserID, filterOptions.ActionType, filterOptions.Before, filterOptions.Limit, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateAuditLog(g.ID, *auditLog, filterOptions, CacheStrategyNoWs), nil
}

// GetIntegrations gets all Integration(s) from the Guild. Requires PermissionManageServer
func (g *guildImpl) GetIntegrations(opts ...rest.RequestOpt) ([]*Integration, rest.Error) {
	guildIntegrations, err := g.disgo.RestServices().GuildService().GetIntegrations(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	integrations := make([]*Integration, len(guildIntegrations))
	for i, guildIntegration := range guildIntegrations {
		integrations[i] = g.disgo.EntityBuilder().CreateIntegration(g.ID, guildIntegration, CacheStrategyNoWs)
	}
	return integrations, nil
}

// DeleteIntegration deletes a specific Integration from the Guild. Requires PermissionManageServer
func (g *guildImpl) DeleteIntegration(integrationID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().GuildService().DeleteIntegration(g.ID, integrationID, opts...)
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (g *guildImpl) GetGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.disgo.RestServices().ApplicationService().GetGuildCommand(g.disgo.ApplicationID(), g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (g *guildImpl) GetGuildCommands(opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.disgo.RestServices().ApplicationService().GetGuildCommands(g.disgo.ApplicationID(), g.ID, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.disgo.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// CreateGuildCommand creates a new Guild discord.ApplicationCommand
func (g *guildImpl) CreateGuildCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.disgo.RestServices().ApplicationService().CreateGuildCommand(g.disgo.ApplicationID(), g.ID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// UpdateGuildCommand edits a specific Guild discord.ApplicationCommand
func (g *guildImpl) UpdateGuildCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.disgo.RestServices().ApplicationService().UpdateGuildCommand(g.disgo.ApplicationID(), g.ID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (g *guildImpl) DeleteGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.disgo.RestServices().ApplicationService().DeleteGuildCommand(g.disgo.ApplicationID(), g.ID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (g *guildImpl) SetGuildCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.disgo.RestServices().ApplicationService().SetGuildCommands(g.disgo.ApplicationID(), g.ID, commandCreates, opts...)
	if err != nil {
		return nil, err
	}
	commands := make([]*ApplicationCommand, len(cmds))
	for i, command := range cmds {
		commands[i] = g.disgo.EntityBuilder().CreateCommand(command, CacheStrategyNoWs)
	}
	return commands, nil
}

// GetGuildCommandsPermissions returns the api.GuildCommandPermissions for a all discord.ApplicationCommand(s) in an api.Guild
func (g *guildImpl) GetGuildCommandsPermissions(opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.disgo.RestServices().ApplicationService().GetGuildCommandsPermissions(g.disgo.ApplicationID(), g.ID, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.disgo.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// GetGuildCommandPermissions returns the api.GuildCommandPermissions for a specific discord.ApplicationCommand in an api.Guild
func (g *guildImpl) GetGuildCommandPermissions(commandID discord.Snowflake, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	permissions, err := g.disgo.RestServices().ApplicationService().GetGuildCommandPermissions(g.disgo.ApplicationID(), g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// SetGuildCommandsPermissions sets the discord.GuildCommandPermissions for all discord.ApplicationCommand(s)
func (g *guildImpl) SetGuildCommandsPermissions(commandPermissions []discord.GuildCommandPermissionsSet, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.disgo.RestServices().ApplicationService().SetGuildCommandsPermissions(g.disgo.ApplicationID(), g.ID, commandPermissions, opts...)
	if err != nil {
		return nil, err
	}
	permissions := make([]*GuildCommandPermissions, len(perms))
	for i, permission := range perms {
		permissions[i] = g.disgo.EntityBuilder().CreateCommandPermissions(permission, CacheStrategyNoWs)
	}
	return permissions, nil
}

// SetGuildCommandPermissions sets the api.GuildCommandPermissions for a specific discord.ApplicationCommand
func (g *guildImpl) SetGuildCommandPermissions(commandID discord.Snowflake, permissions []discord.CommandPermission, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	perms, err := g.disgo.RestServices().ApplicationService().SetGuildCommandPermissions(g.disgo.ApplicationID(), g.ID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateCommandPermissions(*perms, CacheStrategyNoWs), nil
}

// GetTemplates gets a specific GuildTemplate
func (g *guildImpl) GetTemplates(opts ...rest.RequestOpt) ([]*GuildTemplate, rest.Error) {
	guildTemplates, err := g.disgo.RestServices().GuildTemplateService().GetGuildTemplates(g.ID, opts...)
	if err != nil {
		return nil, err
	}
	templates := make([]*GuildTemplate, len(guildTemplates))
	for i, guildTemplate := range guildTemplates {
		templates[i] = g.disgo.EntityBuilder().CreateGuildTemplate(guildTemplate, CacheStrategyNoWs)
	}
	return templates, nil
}

// CreateTemplate creates a new GuildTemplate
func (g *guildImpl) CreateTemplate(guildTemplateCreate discord.GuildTemplateCreate, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.disgo.RestServices().GuildTemplateService().CreateGuildTemplate(g.ID, guildTemplateCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// SyncTemplate syncs the current Guild status to an existing GuildTemplate
func (g *guildImpl) SyncTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.disgo.RestServices().GuildTemplateService().SyncGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// UpdateTemplate updates a specific GuildTemplate
func (g *guildImpl) UpdateTemplate(templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.disgo.RestServices().GuildTemplateService().UpdateGuildTemplate(g.ID, templateCode, guildTemplateUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// DeleteTemplate deletes a specific GuildTemplate
func (g *guildImpl) DeleteTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.disgo.RestServices().GuildTemplateService().DeleteGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}
