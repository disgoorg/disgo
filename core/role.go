package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Role struct {
	discord.Role
	Bot *Bot
}

func (r *Role) String() string {
	return "<@&" + r.ID.String() + ">"
}

// Guild returns the Guild of this role from the Caches
func (r *Role) Guild() *Guild {
	return r.Bot.Caches.GuildCache().Get(r.GuildID)
}

// Update updates the Role with specific values
func (r *Role) Update(roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := r.Bot.RestServices.GuildService().UpdateRole(r.GuildID, r.ID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return r.Bot.EntityBuilder.CreateRole(r.GuildID, *role, CacheStrategyNoWs), nil
}

// SetPosition sets the position of the Role
func (r *Role) SetPosition(rolePositionUpdate discord.RolePositionUpdate, opts ...rest.RequestOpt) ([]*Role, rest.Error) {
	roles, err := r.Bot.RestServices.GuildService().UpdateRolePositions(r.GuildID, []discord.RolePositionUpdate{rolePositionUpdate}, opts...)
	if err != nil {
		return nil, err
	}
	coreRoles := make([]*Role, len(roles))
	for i, role := range roles {
		coreRoles[i] = r.Bot.EntityBuilder.CreateRole(r.GuildID, role, CacheStrategyNoWs)
	}
	return coreRoles, nil
}

// Delete deletes the Role
func (r *Role) Delete(opts ...rest.RequestOpt) rest.Error {
	return r.Bot.RestServices.GuildService().DeleteRole(r.GuildID, r.ID, opts...)
}
