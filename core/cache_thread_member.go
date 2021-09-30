package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	ThreadMemberFindFunc func(threadMember *ThreadMember) bool

	ThreadMemberCache interface {
		Get(threadID discord.Snowflake, userID discord.Snowflake) *ThreadMember
		GetCopy(threadID discord.Snowflake, userID discord.Snowflake) *ThreadMember
		Set(threadMember *ThreadMember) *ThreadMember
		Remove(threadID discord.Snowflake, userID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*ThreadMember
		All() map[discord.Snowflake][]*ThreadMember

		ThreadCache(threadID discord.Snowflake) map[discord.Snowflake]*ThreadMember
		ThreadAll(threadID discord.Snowflake) []*ThreadMember

		FindFirst(threadMemberFindFunc ThreadMemberFindFunc) *ThreadMember
		FindAll(threadMemberFindFunc ThreadMemberFindFunc) []*ThreadMember
	}

	threadMemberCacheImpl struct {
		cacheFlags    CacheFlags
		threadMembers map[discord.Snowflake]map[discord.Snowflake]*ThreadMember
	}
)

func NewThreadMemberCache(cacheFlags CacheFlags) ThreadMemberCache {
	return &threadMemberCacheImpl{
		cacheFlags:    cacheFlags,
		threadMembers: map[discord.Snowflake]map[discord.Snowflake]*ThreadMember{},
	}
}

func (c *threadMemberCacheImpl) Get(threadID discord.Snowflake, userID discord.Snowflake) *ThreadMember {
	if _, ok := c.threadMembers[threadID]; !ok {
		return nil
	}
	return c.threadMembers[threadID][userID]
}

func (c *threadMemberCacheImpl) GetCopy(threadID discord.Snowflake, userID discord.Snowflake) *ThreadMember {
	if threadMember := c.Get(threadID, userID); threadMember != nil {
		m := *threadMember
		return &m
	}
	return nil
}

func (c *threadMemberCacheImpl) Set(threadMember *ThreadMember) *ThreadMember {
	// always cache self threadMembers
	if c.cacheFlags.Missing(CacheFlagRoles) && threadMember.UserID != threadMember.Bot.ClientID {
		return threadMember
	}
	if _, ok := c.threadMembers[threadMember.ID]; !ok {
		c.threadMembers[threadMember.ID] = map[discord.Snowflake]*ThreadMember{}
	}
	rol, ok := c.threadMembers[threadMember.ID][threadMember.UserID]
	if ok {
		*rol = *threadMember
		return rol
	}
	c.threadMembers[threadMember.ID][threadMember.UserID] = threadMember

	return threadMember
}

func (c *threadMemberCacheImpl) Remove(threadID discord.Snowflake, userID discord.Snowflake) {
	if _, ok := c.threadMembers[threadID]; !ok {
		return
	}
	delete(c.threadMembers[threadID], userID)
}

func (c *threadMemberCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*ThreadMember {
	return c.threadMembers
}

func (c *threadMemberCacheImpl) All() map[discord.Snowflake][]*ThreadMember {
	threadMembers := make(map[discord.Snowflake][]*ThreadMember, len(c.threadMembers))
	for threadID, guildThreadMembers := range c.threadMembers {
		threadMembers[threadID] = make([]*ThreadMember, len(guildThreadMembers))
		i := 0
		for _, threadMember := range guildThreadMembers {
			threadMembers[threadID] = append(threadMembers[threadID], threadMember)
		}
		i++
	}
	return threadMembers
}

func (c *threadMemberCacheImpl) ThreadCache(threadID discord.Snowflake) map[discord.Snowflake]*ThreadMember {
	if _, ok := c.threadMembers[threadID]; !ok {
		return nil
	}
	return c.threadMembers[threadID]
}

func (c *threadMemberCacheImpl) ThreadAll(threadID discord.Snowflake) []*ThreadMember {
	if _, ok := c.threadMembers[threadID]; !ok {
		return nil
	}
	threadMembers := make([]*ThreadMember, len(c.threadMembers[threadID]))
	i := 0
	for _, threadMember := range c.threadMembers[threadID] {
		threadMembers = append(threadMembers, threadMember)
		i++
	}
	return threadMembers
}

func (c *threadMemberCacheImpl) FindFirst(threadMemberFindFunc ThreadMemberFindFunc) *ThreadMember {
	for _, guildThreadMembers := range c.threadMembers {
		for _, threadMember := range guildThreadMembers {
			if threadMemberFindFunc(threadMember) {
				return threadMember
			}
		}
	}
	return nil
}

func (c *threadMemberCacheImpl) FindAll(threadMemberFindFunc ThreadMemberFindFunc) []*ThreadMember {
	var threadMembers []*ThreadMember
	for _, guildThreadMembers := range c.threadMembers {
		for _, threadMember := range guildThreadMembers {
			if threadMemberFindFunc(threadMember) {
				threadMembers = append(threadMembers, threadMember)
			}
		}
	}
	return threadMembers
}
