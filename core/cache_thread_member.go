package core

import "github.com/DisgoOrg/snowflake"

type (
	ThreadMemberFindFunc func(threadMember *ThreadMember) bool

	ThreadMemberCache interface {
		Get(threadID snowflake.Snowflake, userID snowflake.Snowflake) *ThreadMember
		GetCopy(threadID snowflake.Snowflake, userID snowflake.Snowflake) *ThreadMember
		Set(threadMember *ThreadMember) *ThreadMember
		Remove(threadID snowflake.Snowflake, userID snowflake.Snowflake)
		RemoveAll(threadID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*ThreadMember
		All() map[snowflake.Snowflake][]*ThreadMember

		ThreadCache(threadID snowflake.Snowflake) map[snowflake.Snowflake]*ThreadMember
		ThreadAll(threadID snowflake.Snowflake) []*ThreadMember

		FindFirst(threadMemberFindFunc ThreadMemberFindFunc) *ThreadMember
		FindAll(threadMemberFindFunc ThreadMemberFindFunc) []*ThreadMember
	}

	threadMemberCacheImpl struct {
		cacheFlags    CacheFlags
		threadMembers map[snowflake.Snowflake]map[snowflake.Snowflake]*ThreadMember
	}
)

func NewThreadMemberCache(cacheFlags CacheFlags) ThreadMemberCache {
	return &threadMemberCacheImpl{
		cacheFlags:    cacheFlags,
		threadMembers: map[snowflake.Snowflake]map[snowflake.Snowflake]*ThreadMember{},
	}
}

func (c *threadMemberCacheImpl) Get(threadID snowflake.Snowflake, userID snowflake.Snowflake) *ThreadMember {
	if _, ok := c.threadMembers[threadID]; !ok {
		return nil
	}
	return c.threadMembers[threadID][userID]
}

func (c *threadMemberCacheImpl) GetCopy(threadID snowflake.Snowflake, userID snowflake.Snowflake) *ThreadMember {
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
	if _, ok := c.threadMembers[threadMember.ThreadID]; !ok {
		c.threadMembers[threadMember.ThreadID] = map[snowflake.Snowflake]*ThreadMember{}
	}
	rol, ok := c.threadMembers[threadMember.ThreadID][threadMember.UserID]
	if ok {
		*rol = *threadMember
		return rol
	}
	c.threadMembers[threadMember.ThreadID][threadMember.UserID] = threadMember

	return threadMember
}

func (c *threadMemberCacheImpl) Remove(threadID snowflake.Snowflake, userID snowflake.Snowflake) {
	if _, ok := c.threadMembers[threadID]; !ok {
		return
	}
	delete(c.threadMembers[threadID], userID)
}

func (c *threadMemberCacheImpl) RemoveAll(threadID snowflake.Snowflake) {
	delete(c.threadMembers, threadID)
}

func (c *threadMemberCacheImpl) Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*ThreadMember {
	return c.threadMembers
}

func (c *threadMemberCacheImpl) All() map[snowflake.Snowflake][]*ThreadMember {
	threadMembers := make(map[snowflake.Snowflake][]*ThreadMember, len(c.threadMembers))
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

func (c *threadMemberCacheImpl) ThreadCache(threadID snowflake.Snowflake) map[snowflake.Snowflake]*ThreadMember {
	if _, ok := c.threadMembers[threadID]; !ok {
		return nil
	}
	return c.threadMembers[threadID]
}

func (c *threadMemberCacheImpl) ThreadAll(threadID snowflake.Snowflake) []*ThreadMember {
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
