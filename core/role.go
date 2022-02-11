package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Role struct {
	discord.Role
	Bot *Bot
}

func (r *Role) IconURL(size int) *string {
	if r.Icon == nil {
		return nil
	}
	compiledRoute, _ := route.RoleIcon.Compile(nil, route.PNG, size, r.ID, *r.Icon)
	url := compiledRoute.URL()
	return &url
}

// Guild returns the Guild this Role belongs to.
// This will only check cached guilds!
func (r *Role) Guild() *Guild {
	return r.Bot.Caches.Guilds().Get(r.GuildID)
}

//ComparePositionTo compares the position of this Role to another Role.
func (r *Role) ComparePositionTo(role *Role) int {
	if role == nil {
		panic("Role argument cannot be nil")
	}
	if r.Position == role.Position {
		return int(r.ID.Time().Sub(role.ID.Time()).Nanoseconds())
	}
	return r.Position - role.Position
}

// Update updates this Role with the properties provided in discord.RoleUpdate
func (r *Role) Update(roleUpdate discord.RoleUpdate, opts ...rest.RequestOpt) (*Role, error) {
	role, err := r.Bot.RestServices.GuildService().UpdateRole(r.GuildID, r.ID, roleUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return r.Bot.EntityBuilder.CreateRole(r.GuildID, *role, CacheStrategyNoWs), nil
}

// SetPosition sets the position of this Role
func (r *Role) SetPosition(rolePositionUpdate discord.RolePositionUpdate, opts ...rest.RequestOpt) ([]*Role, error) {
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

// Delete deletes this Role
func (r *Role) Delete(opts ...rest.RequestOpt) error {
	return r.Bot.RestServices.GuildService().DeleteRole(r.GuildID, r.ID, opts...)
}
