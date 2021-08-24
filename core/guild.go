package core

import (
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
func (g *Guild) Update(updateGuild discord.GuildUpdate, opts ...rest.RequestOpt) (*Guild, rest.Error) {
	guild, err := g.Disgo.RestServices().GuildService().UpdateGuild(g.ID, updateGuild, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuild(*guild, CacheStrategyNoWs), nil
}

// Delete deletes the current Guild
func (g *Guild) Delete(opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteGuild(g.ID, opts...)
}

// PublicRole returns the @everyone Role
func (g *Guild) PublicRole() *Role {
	return g.Disgo.Cache().RoleCache().Get(g.ID)
}

// CreateRole allows you to create a new Role
func (g *Guild) CreateRole(roleCreate discord.RoleCreate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := g.Disgo.RestServices().GuildService().CreateRole(g.ID, roleCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// UpdateRole allows you to update a Role
func (g *Guild) UpdateRole(roleID discord.Snowflake, roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := g.Disgo.RestServices().GuildService().UpdateRole(g.ID, roleID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateRole(g.ID, *role, CacheStrategyNoWs), nil
}

// DeleteRole allows you to delete a Role
func (g *Guild) DeleteRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteRole(g.ID, roleID, opts...)
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
func (g *Guild) Leave(opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().UserService().LeaveGuild(g.ID, opts...)
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
func (g *Guild) AddMember(userID discord.Snowflake, memberAdd discord.MemberAdd, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := g.Disgo.RestServices().GuildService().AddMember(g.ID, userID, memberAdd, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// UpdateMember updates an existing member of the Guild
func (g *Guild) UpdateMember(userID discord.Snowflake, memberUpdate discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := g.Disgo.RestServices().GuildService().UpdateMember(g.ID, userID, memberUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateMember(g.ID, *member, CacheStrategyNoWs), nil
}

// KickMember kicks an existing member from the Guild
func (g *Guild) KickMember(userID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().RemoveMember(g.ID, userID, opts...)
}

// BanMember bans a Member from the Guild
func (g *Guild) BanMember(userID discord.Snowflake, deleteMessageDays int, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().AddBan(g.ID, userID, deleteMessageDays, opts...)
}

// UnbanMember unbans a Member from the Guild
func (g *Guild) UnbanMember(userID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteBan(g.ID, userID, opts...)
}

// GetBans fetches all bans for this Guild
func (g *Guild) GetBans(opts ...rest.RequestOpt) ([]*Ban, rest.Error) {
	guildBans, err := g.Disgo.RestServices().GuildService().GetBans(g.ID, opts...)
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
func (g *Guild) GetBan(userID discord.Snowflake, opts ...rest.RequestOpt) (*Ban, rest.Error) {
	ban, err := g.Disgo.RestServices().GuildService().GetBan(g.ID, userID, opts...)
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
func (g *Guild) GetAuditLogs(filterOptions AuditLogFilterOptions, opts ...rest.RequestOpt) (*AuditLog, rest.Error) {
	auditLog, err := g.Disgo.RestServices().AuditLogService().GetAuditLog(g.ID, filterOptions.UserID, filterOptions.ActionType, filterOptions.Before, filterOptions.Limit, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateAuditLog(g.ID, *auditLog, filterOptions, CacheStrategyNoWs), nil
}

// GetIntegrations gets all Integration(s) from the Guild. Requires PermissionManageServer
func (g *Guild) GetIntegrations(opts ...rest.RequestOpt) ([]*Integration, rest.Error) {
	guildIntegrations, err := g.Disgo.RestServices().GuildService().GetIntegrations(g.ID, opts...)
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
func (g *Guild) DeleteIntegration(integrationID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().GuildService().DeleteIntegration(g.ID, integrationID, opts...)
}

// GetGuildCommand fetches a specific Guild discord.ApplicationCommand
func (g *Guild) GetGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().GetGuildCommand(g.Disgo.ApplicationID(), g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// GetGuildCommands fetches all Guild discord.ApplicationCommand(s)
func (g *Guild) GetGuildCommands(opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.Disgo.RestServices().ApplicationService().GetGuildCommands(g.Disgo.ApplicationID(), g.ID, opts...)
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
func (g *Guild) CreateGuildCommand(commandCreate discord.ApplicationCommandCreate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().CreateGuildCommand(g.Disgo.ApplicationID(), g.ID, commandCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// UpdateGuildCommand edits a specific Guild discord.ApplicationCommand
func (g *Guild) UpdateGuildCommand(commandID discord.Snowflake, commandUpdate discord.ApplicationCommandUpdate, opts ...rest.RequestOpt) (*ApplicationCommand, rest.Error) {
	command, err := g.Disgo.RestServices().ApplicationService().UpdateGuildCommand(g.Disgo.ApplicationID(), g.ID, commandID, commandUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommand(*command, CacheStrategyNoWs), nil
}

// DeleteGuildCommand creates a new Guild discord.ApplicationCommand
func (g *Guild) DeleteGuildCommand(commandID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return g.Disgo.RestServices().ApplicationService().DeleteGuildCommand(g.Disgo.ApplicationID(), g.ID, commandID, opts...)
}

// SetGuildCommands overrides all Guild discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommands(commandCreates []discord.ApplicationCommandCreate, opts ...rest.RequestOpt) ([]*ApplicationCommand, rest.Error) {
	cmds, err := g.Disgo.RestServices().ApplicationService().SetGuildCommands(g.Disgo.ApplicationID(), g.ID, commandCreates, opts...)
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
func (g *Guild) GetGuildCommandsPermissions(opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().GetGuildCommandsPermissions(g.Disgo.ApplicationID(), g.ID, opts...)
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
func (g *Guild) GetGuildCommandPermissions(commandID discord.Snowflake, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	permissions, err := g.Disgo.RestServices().ApplicationService().GetGuildCommandPermissions(g.Disgo.ApplicationID(), g.ID, commandID, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommandPermissions(*permissions, CacheStrategyNoWs), nil
}

// SetGuildCommandsPermissions sets the discord.GuildCommandPermissions for all discord.ApplicationCommand(s)
func (g *Guild) SetGuildCommandsPermissions(commandPermissions []discord.GuildCommandPermissionsSet, opts ...rest.RequestOpt) ([]*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().SetGuildCommandsPermissions(g.Disgo.ApplicationID(), g.ID, commandPermissions, opts...)
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
func (g *Guild) SetGuildCommandPermissions(commandID discord.Snowflake, permissions []discord.CommandPermission, opts ...rest.RequestOpt) (*GuildCommandPermissions, rest.Error) {
	perms, err := g.Disgo.RestServices().ApplicationService().SetGuildCommandPermissions(g.Disgo.ApplicationID(), g.ID, commandID, permissions, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateCommandPermissions(*perms, CacheStrategyNoWs), nil
}

// GetTemplates gets a specific GuildTemplate
func (g *Guild) GetTemplates(opts ...rest.RequestOpt) ([]*GuildTemplate, rest.Error) {
	guildTemplates, err := g.Disgo.RestServices().GuildTemplateService().GetGuildTemplates(g.ID, opts...)
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
func (g *Guild) CreateTemplate(guildTemplateCreate discord.GuildTemplateCreate, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().CreateGuildTemplate(g.ID, guildTemplateCreate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// SyncTemplate syncs the current Guild status to an existing GuildTemplate
func (g *Guild) SyncTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().SyncGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// UpdateTemplate updates a specific GuildTemplate
func (g *Guild) UpdateTemplate(templateCode string, guildTemplateUpdate discord.GuildTemplateUpdate, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().UpdateGuildTemplate(g.ID, templateCode, guildTemplateUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}

// DeleteTemplate deletes a specific GuildTemplate
func (g *Guild) DeleteTemplate(templateCode string, opts ...rest.RequestOpt) (*GuildTemplate, rest.Error) {
	guildTemplate, err := g.Disgo.RestServices().GuildTemplateService().DeleteGuildTemplate(g.ID, templateCode, opts...)
	if err != nil {
		return nil, err
	}
	return g.Disgo.EntityBuilder().CreateGuildTemplate(*guildTemplate, CacheStrategyNoWs), nil
}
