package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	MemberFindFunc func(member *Member) bool

	MemberCache interface {
		Get(guildID discord.Snowflake, userID discord.Snowflake) *Member
		GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *Member
		Set(member *Member) *Member
		Remove(guildID discord.Snowflake, userID discord.Snowflake)
		RemoveAll(guildID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Member
		All() map[discord.Snowflake][]*Member

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Member
		GuildAll(guildID discord.Snowflake) []*Member

		FindFirst(memberFindFunc MemberFindFunc) *Member
		FindAll(memberFindFunc MemberFindFunc) []*Member
	}

	memberCacheImpl struct {
		memberCachePolicy MemberCachePolicy
		members           map[discord.Snowflake]map[discord.Snowflake]*Member
	}
)

func NewMemberCache(memberCachePolicy MemberCachePolicy) MemberCache {
	if memberCachePolicy == nil {
		memberCachePolicy = MemberCachePolicyDefault
	}
	return &memberCacheImpl{
		memberCachePolicy: memberCachePolicy,
		members:           map[discord.Snowflake]map[discord.Snowflake]*Member{},
	}
}

func (c *memberCacheImpl) Get(guildID discord.Snowflake, userID discord.Snowflake) *Member {
	if _, ok := c.members[guildID]; !ok {
		return nil
	}
	return c.members[guildID][userID]
}

func (c *memberCacheImpl) GetCopy(guildID discord.Snowflake, userID discord.Snowflake) *Member {
	if member := c.Get(guildID, userID); member != nil {
		m := *member
		return &m
	}
	return nil
}

func (c *memberCacheImpl) Set(member *Member) *Member {
	// always cache self members
	if !c.memberCachePolicy(member) && member.User.ID != member.Bot.ClientID {
		return member
	}
	if _, ok := c.members[member.GuildID]; !ok {
		c.members[member.GuildID] = map[discord.Snowflake]*Member{}
	}
	rol, ok := c.members[member.GuildID][member.User.ID]
	if ok {
		*rol = *member
		return rol
	}
	c.members[member.GuildID][member.User.ID] = member

	return member
}

func (c *memberCacheImpl) Remove(guildID discord.Snowflake, userID discord.Snowflake) {
	if _, ok := c.members[guildID]; !ok {
		return
	}
	delete(c.members[guildID], userID)
}

func (c *memberCacheImpl) RemoveAll(guildID discord.Snowflake) {
	delete(c.members, guildID)
}

func (c *memberCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Member {
	return c.members
}

func (c *memberCacheImpl) All() map[discord.Snowflake][]*Member {
	members := make(map[discord.Snowflake][]*Member, len(c.members))
	for guildID, guildMembers := range c.members {
		members[guildID] = make([]*Member, len(guildMembers))
		i := 0
		for _, member := range guildMembers {
			members[guildID] = append(members[guildID], member)
		}
		i++
	}
	return members
}

func (c *memberCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Member {
	if _, ok := c.members[guildID]; !ok {
		return nil
	}
	return c.members[guildID]
}

func (c *memberCacheImpl) GuildAll(guildID discord.Snowflake) []*Member {
	if _, ok := c.members[guildID]; !ok {
		return nil
	}
	members := make([]*Member, len(c.members[guildID]))
	i := 0
	for _, member := range c.members[guildID] {
		members = append(members, member)
		i++
	}
	return members
}

func (c *memberCacheImpl) FindFirst(memberFindFunc MemberFindFunc) *Member {
	for _, guildMembers := range c.members {
		for _, member := range guildMembers {
			if memberFindFunc(member) {
				return member
			}
		}
	}
	return nil
}

func (c *memberCacheImpl) FindAll(memberFindFunc MemberFindFunc) []*Member {
	var members []*Member
	for _, guildMembers := range c.members {
		for _, member := range guildMembers {
			if memberFindFunc(member) {
				members = append(members, member)
			}
		}
	}
	return members
}
