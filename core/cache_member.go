package core

import "github.com/DisgoOrg/snowflake"

type (
	MemberFindFunc func(member *Member) bool

	MemberCache interface {
		Get(guildID snowflake.Snowflake, userID snowflake.Snowflake) *Member
		GetCopy(guildID snowflake.Snowflake, userID snowflake.Snowflake) *Member
		Set(member *Member) *Member
		Remove(guildID snowflake.Snowflake, userID snowflake.Snowflake)
		RemoveAll(guildID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Member
		All() map[snowflake.Snowflake][]*Member

		GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*Member
		GuildAll(guildID snowflake.Snowflake) []*Member

		FindFirst(memberFindFunc MemberFindFunc) *Member
		FindAll(memberFindFunc MemberFindFunc) []*Member
	}

	memberCacheImpl struct {
		memberCachePolicy MemberCachePolicy
		members           map[snowflake.Snowflake]map[snowflake.Snowflake]*Member
	}
)

func NewMemberCache(memberCachePolicy MemberCachePolicy) MemberCache {
	if memberCachePolicy == nil {
		memberCachePolicy = MemberCachePolicyDefault
	}
	return &memberCacheImpl{
		memberCachePolicy: memberCachePolicy,
		members:           map[snowflake.Snowflake]map[snowflake.Snowflake]*Member{},
	}
}

func (c *memberCacheImpl) Get(guildID snowflake.Snowflake, userID snowflake.Snowflake) *Member {
	if _, ok := c.members[guildID]; !ok {
		return nil
	}
	return c.members[guildID][userID]
}

func (c *memberCacheImpl) GetCopy(guildID snowflake.Snowflake, userID snowflake.Snowflake) *Member {
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
		c.members[member.GuildID] = map[snowflake.Snowflake]*Member{}
	}
	rol, ok := c.members[member.GuildID][member.User.ID]
	if ok {
		*rol = *member
		return rol
	}
	c.members[member.GuildID][member.User.ID] = member

	return member
}

func (c *memberCacheImpl) Remove(guildID snowflake.Snowflake, userID snowflake.Snowflake) {
	if _, ok := c.members[guildID]; !ok {
		return
	}
	delete(c.members[guildID], userID)
}

func (c *memberCacheImpl) RemoveAll(guildID snowflake.Snowflake) {
	delete(c.members, guildID)
}

func (c *memberCacheImpl) Cache() map[snowflake.Snowflake]map[snowflake.Snowflake]*Member {
	return c.members
}

func (c *memberCacheImpl) All() map[snowflake.Snowflake][]*Member {
	members := make(map[snowflake.Snowflake][]*Member, len(c.members))
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

func (c *memberCacheImpl) GuildCache(guildID snowflake.Snowflake) map[snowflake.Snowflake]*Member {
	if _, ok := c.members[guildID]; !ok {
		return nil
	}
	return c.members[guildID]
}

func (c *memberCacheImpl) GuildAll(guildID snowflake.Snowflake) []*Member {
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
