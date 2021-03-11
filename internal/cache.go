package internal

import (
	"strings"

	"github.com/DiscoOrg/disgo/api"
)

func newCacheImpl(memberCachePolicy api.MemberCachePolicy) api.Cache {
	return &CacheImpl{
		memberCachePolicy: memberCachePolicy,
		users:             map[api.Snowflake]*api.User{},
		guilds:            map[api.Snowflake]*api.Guild{},
		members:           map[api.Snowflake]map[api.Snowflake]*api.Member{},
		roles:             map[api.Snowflake]map[api.Snowflake]*api.Role{},
		dmChannels:        map[api.Snowflake]map[api.Snowflake]*api.DMChannel{},
		categories:        map[api.Snowflake]map[api.Snowflake]*api.CategoryChannel{},
		textChannels:      map[api.Snowflake]map[api.Snowflake]*api.TextChannel{},
		voiceChannels:     map[api.Snowflake]map[api.Snowflake]*api.VoiceChannel{},
		storeChannels:     map[api.Snowflake]map[api.Snowflake]*api.StoreChannel{},
	}
}

type CacheImpl struct {
	memberCachePolicy api.MemberCachePolicy
	users             map[api.Snowflake]*api.User
	guilds            map[api.Snowflake]*api.Guild
	members           map[api.Snowflake]map[api.Snowflake]*api.Member
	roles             map[api.Snowflake]map[api.Snowflake]*api.Role
	dmChannels        map[api.Snowflake]map[api.Snowflake]*api.DMChannel
	categories        map[api.Snowflake]map[api.Snowflake]*api.CategoryChannel
	textChannels      map[api.Snowflake]map[api.Snowflake]*api.TextChannel
	voiceChannels     map[api.Snowflake]map[api.Snowflake]*api.VoiceChannel
	storeChannels     map[api.Snowflake]map[api.Snowflake]*api.StoreChannel
}

// user cache
func (c CacheImpl) User(id api.Snowflake) *api.User {
	return c.users[id]
}
func (c CacheImpl) UserByTag(tag string) *api.User {
	for _, user := range c.users {
		if user.Username+"#"+user.Discriminator == tag {
			return user
		}
	}
	return nil
}
func (c CacheImpl) UsersByName(name string, ignoreCase bool) []*api.User {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	users := make([]*api.User, 1)
	for _, user := range c.users {
		if ignoreCase && strings.ToLower(user.Username) == name || !ignoreCase && user.Username == name {
			users = append(users, user)
		}
	}
	return users
}
func (c CacheImpl) Users() []*api.User {
	users := make([]*api.User, len(c.users))
	i := 0
	for _, user := range c.users {
		users[i] = user
		i++
	}
	return users
}
func (c CacheImpl) UserCache() map[api.Snowflake]*api.User {
	return c.users
}
func (c CacheImpl) CacheUser(user *api.User) {
	if _, ok := c.guilds[user.ID]; ok {
		// update old user
		return
	}
	c.users[user.ID] = user
}
func (c CacheImpl) UncacheUser(id api.Snowflake) {
	delete(c.users, id)
}


// guild cache
func (c CacheImpl) Guild(id api.Snowflake) *api.Guild {
	return c.guilds[id]
}
func (c CacheImpl) GuildsByName(name string, ignoreCase bool) []*api.Guild {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	guilds := make([]*api.Guild, 1)
	for _, guild := range c.guilds {
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}
func (c CacheImpl) Guilds() []*api.Guild {
	guilds := make([]*api.Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}
func (c CacheImpl) GuildCache() map[api.Snowflake]*api.Guild {
	return c.guilds
}
func (c CacheImpl) CacheGuild(guild *api.Guild) {
	if _, ok := c.guilds[guild.ID]; ok {
		// update old guild
		return
	}
	// guild was not yet cached so cache it directly
	c.guilds[guild.ID] = guild
}
func (c CacheImpl) UncacheGuild(id api.Snowflake) {
	delete(c.guilds, id)
}


// member cache
func (c CacheImpl) Member(guildID api.Snowflake, userID api.Snowflake) *api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		return guildMembers[userID]
	}
	return nil
}
func (c CacheImpl) MemberByTag(guildID api.Snowflake, tag string) *api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		for _, member := range guildMembers {
			if member.Username+"#"+member.Discriminator == tag {
				return member
			}
		}
	}
	return nil
}
func (c CacheImpl) MembersByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		members := make([]*api.Member, 1)
		for _, member := range guildMembers {
			if ignoreCase && strings.ToLower(member.Username) == name || !ignoreCase && member.Username == name {
				members = append(members, member)
			}
		}
		return members
	}
	return nil
}
func (c CacheImpl) Members(guildID api.Snowflake) []*api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		members := make([]*api.Member, len(guildMembers))
		i := 0
		for _, member := range guildMembers {
			members[i] = member
			i++
		}
		return members
	}
	return nil
}
func (c CacheImpl) AllMembers() []*api.Member {
	members := make([]*api.Member, len(c.guilds))
	for _, guildMembers := range c.members {
		for _, member := range guildMembers {
			members = append(members, member)
		}
	}
	return members
}
func (c CacheImpl) MemberCache(guildID api.Snowflake) map[api.Snowflake]*api.Member {
	return c.members[guildID]
}
func (c CacheImpl) AllMemberCache() map[api.Snowflake]map[api.Snowflake]*api.Member {
	return c.members
}
func (c CacheImpl) CacheMember(member *api.Member) {
	if guildMembers, ok := c.members[member.GuildID]; ok {
		if _, ok := guildMembers[member.ID]; ok {
			// update old guild
			return
		}
		guildMembers[member.ID] = member
	}
}
func (c CacheImpl) UncacheMember(guildID api.Snowflake, userID api.Snowflake) {
	delete(c.members[guildID], userID)
}


// role cache
func (c CacheImpl) Role(guildID api.Snowflake, userID api.Snowflake) *api.Role {
	if guildRoles, ok := c.roles[guildID]; ok {
		return guildRoles[userID]
	}
	return nil
}
func (c CacheImpl) RolesByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.Role {
	if guildRoles, ok := c.roles[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		roles := make([]*api.Role, 1)
		for _, role := range guildRoles {
			if ignoreCase && strings.ToLower(role.Name) == name || !ignoreCase && role.Name == name {
				roles = append(roles, role)
			}
		}
		return roles
	}
	return nil
}
func (c CacheImpl) Roles(guildID api.Snowflake) []*api.Role {
	if guildRoles, ok := c.roles[guildID]; ok {
		roles := make([]*api.Role, len(guildRoles))
		i := 0
		for _, role := range guildRoles {
			roles[i] = role
			i++
		}
		return roles
	}
	return nil
}
func (c CacheImpl) AllRoles() []*api.Role {
	roles := make([]*api.Role, len(c.guilds))
	for _, guildRoles := range c.roles {
		for _, role := range guildRoles {
			roles = append(roles, role)
		}
	}
	return roles
}
func (c CacheImpl) RoleCache(guildID api.Snowflake) map[api.Snowflake]*api.Role {
	return c.roles[guildID]
}
func (c CacheImpl) AllRoleCache() map[api.Snowflake]map[api.Snowflake]*api.Role {
	return c.roles
}
func (c CacheImpl) CacheRole(role *api.Role) {
	if guildRoles, ok := c.roles[role.GuildID]; ok {
		if _, ok := guildRoles[role.ID]; ok {
			// update old role
			return
		}
		guildRoles[role.ID] = role
	}
}
func (c CacheImpl) UncacheRole(guildID api.Snowflake, userID api.Snowflake) {
	delete(c.roles[guildID], userID)
}
