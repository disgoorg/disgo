package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Role struct {
	discord.Role
	Disgo Disgo
}

// Mention parses the Role as a Mention
func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

// String parses the Role to a String representation
func (r *Role) String() string {
	return r.Mention()
}

// Guild returns the Guild of this role from the Cache
func (r *Role) Guild() *Guild {
	return r.Disgo.Cache().GuildCache().Get(r.GuildID)
}

// Update updates the Role with specific values
func (r *Role) Update(roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, rest.Error) {
	role, err := r.Disgo.RestServices().GuildService().UpdateRole(r.GuildID, r.ID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return r.Disgo.EntityBuilder().CreateRole(r.GuildID, *role, CacheStrategyNoWs), nil
}

// SetPosition sets the position of the Role
func (r *Role) SetPosition(rolePositionUpdate discord.RolePositionUpdate, opts ...rest.RequestOpt) ([]*Role, rest.Error) {
	roles, err := r.Disgo.RestServices().GuildService().UpdateRolePositions(r.GuildID, []discord.RolePositionUpdate{rolePositionUpdate}, opts...)
	if err != nil {
		return nil, err
	}
	coreRoles := make([]*Role, len(roles))
	for i, role := range roles {
		coreRoles[i] = r.Disgo.EntityBuilder().CreateRole(r.GuildID, role, CacheStrategyNoWs)
	}
	return coreRoles, nil
}

// Delete deletes the Role
func (r *Role) Delete(opts ...rest.RequestOpt) rest.Error {
	return r.Disgo.RestServices().GuildService().DeleteRole(r.GuildID, r.ID, opts...)
}
