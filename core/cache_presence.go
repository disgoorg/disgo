package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	PresenceFindFunc func(presence *Presence) bool

	PresenceCache interface {
		Get(guildID discord.Snowflake, userID discord.Snowflake) *Presence
		GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *Presence
		Set(presence *Presence) *Presence
		Remove(guildID discord.Snowflake, userID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Presence
		All() map[discord.Snowflake][]*Presence

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Presence
		GuildAll(guildID discord.Snowflake) []*Presence

		FindFirst(presenceFindFunc PresenceFindFunc) *Presence
		FindAll(presenceFindFunc PresenceFindFunc) []*Presence
	}

	presenceCacheImpl struct {
		cacheFlags CacheFlags
		presences  map[discord.Snowflake]map[discord.Snowflake]*Presence
	}
)

func NewPresenceCache(cacheFlags CacheFlags) PresenceCache {
	return &presenceCacheImpl{
		cacheFlags: cacheFlags,
		presences:  map[discord.Snowflake]map[discord.Snowflake]*Presence{},
	}
}

func (c *presenceCacheImpl) Get(guildID discord.Snowflake, userID discord.Snowflake) *Presence {
	if _, ok := c.presences[guildID]; !ok {
		return nil
	}
	return c.presences[guildID][userID]
}

func (c *presenceCacheImpl) GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *Presence {
	if presence := c.Get(guildID, userID); presence != nil {
		m := *presence
		return &m
	}
	return nil
}

func (c *presenceCacheImpl) Set(presence *Presence) *Presence {
	if c.cacheFlags.Missing(CacheFlagPresences) {
		return presence
	}
	if _, ok := c.presences[presence.GuildID]; !ok {
		c.presences[presence.GuildID] = map[discord.Snowflake]*Presence{}
	}
	rol, ok := c.presences[presence.GuildID][presence.PresenceUser.ID]
	if ok {
		*rol = *presence
		return rol
	}
	c.presences[presence.GuildID][presence.PresenceUser.ID] = presence

	return presence
}

func (c *presenceCacheImpl) Remove(guildID discord.Snowflake, userID discord.Snowflake) {
	if _, ok := c.presences[guildID]; !ok {
		return
	}
	delete(c.presences[guildID], userID)
}

func (c *presenceCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Presence {
	return c.presences
}

func (c *presenceCacheImpl) All() map[discord.Snowflake][]*Presence {
	presences := make(map[discord.Snowflake][]*Presence, len(c.presences))
	for guildID, guildPresences := range c.presences {
		presences[guildID] = make([]*Presence, len(guildPresences))
		i := 0
		for _, presence := range guildPresences {
			presences[guildID] = append(presences[guildID], presence)
		}
		i++
	}
	return presences
}

func (c *presenceCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Presence {
	if _, ok := c.presences[guildID]; !ok {
		return nil
	}
	return c.presences[guildID]
}

func (c *presenceCacheImpl) GuildAll(guildID discord.Snowflake) []*Presence {
	if _, ok := c.presences[guildID]; !ok {
		return nil
	}
	presences := make([]*Presence, len(c.presences[guildID]))
	i := 0
	for _, presence := range c.presences[guildID] {
		presences = append(presences, presence)
		i++
	}
	return presences
}

func (c *presenceCacheImpl) FindFirst(presenceFindFunc PresenceFindFunc) *Presence {
	for _, guildPresences := range c.presences {
		for _, presence := range guildPresences {
			if presenceFindFunc(presence) {
				return presence
			}
		}
	}
	return nil
}

func (c *presenceCacheImpl) FindAll(presenceFindFunc PresenceFindFunc) []*Presence {
	var presences []*Presence
	for _, guildPresences := range c.presences {
		for _, presence := range guildPresences {
			if presenceFindFunc(presence) {
				presences = append(presences, presence)
			}
		}
	}
	return presences
}
