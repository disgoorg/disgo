package core

import "github.com/DisgoOrg/snowflake"

type (
	GuildScheduledEventFindFunc func(guildScheduledEvent *GuildScheduledEvent) bool

	GuildScheduledEventCache interface {
		Get(guildScheduledEventID snowflake.Snowflake) *GuildScheduledEvent
		GetCopy(guildScheduledEventID snowflake.Snowflake) *GuildScheduledEvent
		Set(guildScheduledEvent *GuildScheduledEvent) *GuildScheduledEvent
		Remove(guildScheduledEventID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]*GuildScheduledEvent
		All() []*GuildScheduledEvent

		FindFirst(guildScheduledEventFindFunc GuildScheduledEventFindFunc) *GuildScheduledEvent
		FindAll(guildScheduledEventFindFunc GuildScheduledEventFindFunc) []*GuildScheduledEvent
	}

	guildScheduledEventCacheImpl struct {
		cacheFlags           CacheFlags
		guildScheduledEvents map[snowflake.Snowflake]*GuildScheduledEvent
	}
)

func NewGuildScheduledEventCache(cacheFlags CacheFlags) GuildScheduledEventCache {
	return &guildScheduledEventCacheImpl{
		cacheFlags:           cacheFlags,
		guildScheduledEvents: map[snowflake.Snowflake]*GuildScheduledEvent{},
	}
}

func (c *guildScheduledEventCacheImpl) Get(guildScheduledEventID snowflake.Snowflake) *GuildScheduledEvent {
	return c.guildScheduledEvents[guildScheduledEventID]
}

func (c *guildScheduledEventCacheImpl) GetCopy(guildScheduledEventID snowflake.Snowflake) *GuildScheduledEvent {
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
	gse, ok := c.guildScheduledEvents[guildScheduledEvent.ID]
	if ok {
		*gse = *guildScheduledEvent
		return gse
	}
	c.guildScheduledEvents[guildScheduledEvent.ID] = guildScheduledEvent
	return guildScheduledEvent
}

func (c *guildScheduledEventCacheImpl) Remove(id snowflake.Snowflake) {
	delete(c.guildScheduledEvents, id)
}

func (c *guildScheduledEventCacheImpl) Cache() map[snowflake.Snowflake]*GuildScheduledEvent {
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
	for _, gse := range c.guildScheduledEvents {
		if guildScheduledEventFindFunc(gse) {
			return gse
		}
	}
	return nil
}

func (c *guildScheduledEventCacheImpl) FindAll(guildScheduledEventFindFunc GuildScheduledEventFindFunc) []*GuildScheduledEvent {
	var guildScheduledEvents []*GuildScheduledEvent
	for _, gse := range c.guildScheduledEvents {
		if guildScheduledEventFindFunc(gse) {
			guildScheduledEvents = append(guildScheduledEvents, gse)
		}
	}
	return guildScheduledEvents
}
