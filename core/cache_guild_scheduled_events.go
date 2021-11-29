package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	GuildScheduledEventFindFunc func(guildScheduledEvent *GuildScheduledEvent) bool

	GuildScheduledEventCache interface {
		Get(guildScheduledEventID discord.Snowflake) *GuildScheduledEvent
		GetCopy(guildScheduledEventID discord.Snowflake) *GuildScheduledEvent
		Set(guildScheduledEvent *GuildScheduledEvent) *GuildScheduledEvent
		Remove(guildScheduledEventID discord.Snowflake)

		Cache() map[discord.Snowflake]*GuildScheduledEvent
		All() []*GuildScheduledEvent

		FindFirst(guildScheduledEventFindFunc GuildScheduledEventFindFunc) *GuildScheduledEvent
		FindAll(guildScheduledEventFindFunc GuildScheduledEventFindFunc) []*GuildScheduledEvent
	}

	guildScheduledEventCacheImpl struct {
		cacheFlags           CacheFlags
		guildScheduledEvents map[discord.Snowflake]*GuildScheduledEvent
	}
)

func NewGuildScheduledEventCache(cacheFlags CacheFlags) GuildScheduledEventCache {
	return &guildScheduledEventCacheImpl{
		cacheFlags:           cacheFlags,
		guildScheduledEvents: map[discord.Snowflake]*GuildScheduledEvent{},
	}
}

func (c *guildScheduledEventCacheImpl) Get(guildScheduledEventID discord.Snowflake) *GuildScheduledEvent {
	return c.guildScheduledEvents[guildScheduledEventID]
}

func (c *guildScheduledEventCacheImpl) GetCopy(guildScheduledEventID discord.Snowflake) *GuildScheduledEvent {
	if guildScheduledEvent := c.Get(guildScheduledEventID); guildScheduledEvent != nil {
		st := *guildScheduledEvent
		return &st
	}
	return nil
}

func (c *guildScheduledEventCacheImpl) Set(guildScheduledEvent *GuildScheduledEvent) *GuildScheduledEvent {
	if c.cacheFlags.Missing(CacheFlagGuildScheduledEvents) {
		return guildScheduledEvent
	}
	stI, ok := c.guildScheduledEvents[guildScheduledEvent.ID]
	if ok {
		*stI = *guildScheduledEvent
		return stI
	}
	c.guildScheduledEvents[guildScheduledEvent.ID] = guildScheduledEvent
	return guildScheduledEvent
}

func (c *guildScheduledEventCacheImpl) Remove(id discord.Snowflake) {
	delete(c.guildScheduledEvents, id)
}

func (c *guildScheduledEventCacheImpl) Cache() map[discord.Snowflake]*GuildScheduledEvent {
	return c.guildScheduledEvents
}

func (c *guildScheduledEventCacheImpl) All() []*GuildScheduledEvent {
	guildScheduledEvents := make([]*GuildScheduledEvent, len(c.guildScheduledEvents))
	i := 0
	for _, guildScheduledEvent := range c.guildScheduledEvents {
		guildScheduledEvents[i] = guildScheduledEvent
		i++
	}
	return guildScheduledEvents
}

func (c *guildScheduledEventCacheImpl) FindFirst(guildScheduledEventFindFunc GuildScheduledEventFindFunc) *GuildScheduledEvent {
	for _, stI := range c.guildScheduledEvents {
		if guildScheduledEventFindFunc(stI) {
			return stI
		}
	}
	return nil
}

func (c *guildScheduledEventCacheImpl) FindAll(guildScheduledEventFindFunc GuildScheduledEventFindFunc) []*GuildScheduledEvent {
	var guildScheduledEvents []*GuildScheduledEvent
	for _, stI := range c.guildScheduledEvents {
		if guildScheduledEventFindFunc(stI) {
			guildScheduledEvents = append(guildScheduledEvents, stI)
		}
	}
	return guildScheduledEvents
}
