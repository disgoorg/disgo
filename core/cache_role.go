package core

import "github.com/DisgoOrg/snowflake"

type (
	RoleFindFunc func(role *Role) bool

	RoleCache interface {
		Get(guildID snowflake.Snowflake, roleID snowflake.Snowflake) *Role
		GetCopy(guildID snowflake.Snowflake, roleID snowflake.Snowflake) *Role
		Set(role *Role) *Role
		Remove(guildID snowflake.Snowflake, roleID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Role
		All() map[snowflake.Snowflake][]*Role

		GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*Role
		GuildAll(guildID snowflake.Snowflake) []*Role

		FindFirst(roleFindFunc RoleFindFunc) *Role
		FindAll(roleFindFunc RoleFindFunc) []*Role
	}

	roleCacheImpl struct {
		cacheFlags CacheFlags
		roles      map[snowflake.Snowflake]map[snowflake.Snowflake]*Role
	}
)

func NewRoleCache(cacheFlags CacheFlags) RoleCache {
	return &roleCacheImpl{
		cacheFlags: cacheFlags,
		roles:      map[snowflake.Snowflake]map[snowflake.Snowflake]*Role{},
	}
}

func (c *roleCacheImpl) Get(guildID snowflake.Snowflake, roleID snowflake.Snowflake) *Role {
	if _, ok := c.roles[guildID]; !ok {
		return nil
	}
	return c.roles[guildID][roleID]
}

func (c *roleCacheImpl) GetCopy(guildID snowflake.Snowflake, roleID snowflake.Snowflake) *Role {
	if role := c.Get(guildID, roleID); role != nil {
		ro := *role
		return &ro
	}
	return nil
}

func (c *roleCacheImpl) Set(role *Role) *Role {
	if c.cacheFlags.Missing(CacheFlagRoles) {
		return role
	}
	if _, ok := c.roles[role.GuildID]; !ok {
		c.roles[role.GuildID] = map[snowflake.Snowflake]*Role{}
	}
	rol, ok := c.roles[role.GuildID][role.ID]
	if ok {
		*rol = *role
		return rol
	}
	c.roles[role.GuildID][role.ID] = role

	return role
}

func (c *roleCacheImpl) Remove(guildID snowflake.Snowflake, roleID snowflake.Snowflake) {
	if _, ok := c.roles[guildID]; !ok {
		return
	}
	delete(c.roles[guildID], roleID)
}

func (c *roleCacheImpl) Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Role {
	return c.roles
}

func (c *roleCacheImpl) All() map[snowflake.Snowflake][]*Role {
	roles := make(map[snowflake.Snowflake][]*Role, len(c.roles))
	for guildID, guildRoles := range c.roles {
		roles[guildID] = make([]*Role, len(guildRoles))
		i := 0
		for _, guildRole := range guildRoles {
			roles[guildID] = append(roles[guildID], guildRole)
		}
		i++
	}
	return roles
}

func (c *roleCacheImpl) GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*Role {
	if _, ok := c.roles[guildID]; !ok {
		return nil
	}
	return c.roles[guildID]
}

func (c *roleCacheImpl) GuildAll(guildID snowflake.Snowflake) []*Role {
	if _, ok := c.roles[guildID]; !ok {
		return nil
	}
	roles := make([]*Role, len(c.roles[guildID]))
	i := 0
	for _, role := range c.roles[guildID] {
		roles = append(roles, role)
		i++
	}
	return roles
}

func (c *roleCacheImpl) FindFirst(roleFindFunc RoleFindFunc) *Role {
	for _, guildRoles := range c.roles {
		for _, role := range guildRoles {
			if roleFindFunc(role) {
				return role
			}
		}
	}
	return nil
}

func (c *roleCacheImpl) FindAll(roleFindFunc RoleFindFunc) []*Role {
	var roles []*Role
	for _, guildRoles := range c.roles {
		for _, role := range guildRoles {
			if roleFindFunc(role) {
				roles = append(roles, role)
			}
		}
	}
	return roles
}
