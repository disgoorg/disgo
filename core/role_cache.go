package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	RoleFindFunc func(role *Role) bool

	RoleCache interface {
		Get(guildID discord.Snowflake, roleID discord.Snowflake) *Role
		GetCopy(guildID discord.Snowflake, roleID discord.Snowflake) *Role
		Set(role *Role) *Role
		Remove(guildID discord.Snowflake, roleID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Role
		All() map[discord.Snowflake][]*Role

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Role
		GuildAll(guildID discord.Snowflake) []*Role

		FindFirst(roleFindFunc RoleFindFunc) *Role
		FindAll(roleFindFunc RoleFindFunc) []*Role
	}

	roleCacheImpl struct {
		cacheFlags CacheFlags
		roles      map[discord.Snowflake]map[discord.Snowflake]*Role
	}
)

func NewRoleCache(cacheFlags CacheFlags) RoleCache {
	return &roleCacheImpl{
		cacheFlags: cacheFlags,
		roles:      map[discord.Snowflake]map[discord.Snowflake]*Role{},
	}
}

func (c *roleCacheImpl) Get(guildID discord.Snowflake, roleID discord.Snowflake) *Role {
	if _, ok := c.roles[guildID]; !ok {
		return nil
	}
	return c.roles[guildID][roleID]
}

func (c *roleCacheImpl) GetCopy(guildID discord.Snowflake, roleID discord.Snowflake) *Role {
	return &*c.Get(guildID, roleID)
}

func (c *roleCacheImpl) Set(role *Role) *Role {
	if !c.cacheFlags.Missing(CacheFlagRoles) {
		return role
	}
	if _, ok := c.roles[role.GuildID]; !ok {
		c.roles[role.GuildID] = map[discord.Snowflake]*Role{}
	}
	rol, ok := c.roles[role.GuildID][role.ID]
	if ok {
		*rol = *role
		return rol
	}
	c.roles[role.GuildID][role.ID] = role

	return role
}

func (c *roleCacheImpl) Remove(guildID discord.Snowflake, roleID discord.Snowflake) {
	if _, ok := c.roles[guildID]; !ok {
		return
	}
	delete(c.roles[guildID], roleID)
}

func (c *roleCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Role {
	return c.roles
}

func (c *roleCacheImpl) All() map[discord.Snowflake][]*Role {
	roles := make(map[discord.Snowflake][]*Role, len(c.roles))
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

func (c *roleCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Role {
	if _, ok := c.roles[guildID]; !ok {
		return nil
	}
	return c.roles[guildID]
}

func (c *roleCacheImpl) GuildAll(guildID discord.Snowflake) []*Role {
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
