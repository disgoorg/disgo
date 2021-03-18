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
		unavailableGuilds: map[api.Snowflake]*api.UnavailableGuild{},
		members:           map[api.Snowflake]map[api.Snowflake]*api.Member{},
		roles:             map[api.Snowflake]map[api.Snowflake]*api.Role{},
		dmChannels:        map[api.Snowflake]*api.DMChannel{},
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
	unavailableGuilds map[api.Snowflake]*api.UnavailableGuild
	members           map[api.Snowflake]map[api.Snowflake]*api.Member
	roles             map[api.Snowflake]map[api.Snowflake]*api.Role
	dmChannels        map[api.Snowflake]*api.DMChannel
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
		if user.Tag() == tag {
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
func (c CacheImpl) FindUser(check func(u *api.User) bool) *api.User {
	for _, user := range c.users {
		if check(user) {
			return user
		}
	}
	return nil
}
func (c CacheImpl) FindUsers(check func(u *api.User) bool) []*api.User {
	users := make([]*api.User, 1)
	for _, user := range c.users {
		if check(user) {
			users = append(users, user)
		}
	}
	return users
}

// guild cache
func (c CacheImpl) Guild(guildID api.Snowflake) *api.Guild {
	return c.guilds[guildID]
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
		// update old guild_events
		return
	}
	// guild_events was not yet cached so cache it directly
	c.guilds[guild.ID] = guild
}
func (c CacheImpl) UncacheGuild(guildID api.Snowflake) {
	delete(c.guilds, guildID)
}
func (c CacheImpl) FindGuild(check func(g *api.Guild) bool) *api.Guild {
	for _, guild := range c.guilds {
		if check(guild) {
			return guild
		}
	}
	return nil
}
func (c CacheImpl) FindGuilds(check func(g *api.Guild) bool) []*api.Guild {
	guilds := make([]*api.Guild, 1)
	for _, guild := range c.guilds {
		if check(guild) {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}

// unavailable guild cache
func (c CacheImpl) UnavailableGuild(guildID api.Snowflake) *api.UnavailableGuild {
	return c.unavailableGuilds[guildID]
}
func (c CacheImpl) CacheUnavailableGuild(unavailableGuild *api.UnavailableGuild) {
	*c.unavailableGuilds[unavailableGuild.ID] = *unavailableGuild
}
func (c CacheImpl) UncacheUnavailableGuild(guildID api.Snowflake) {
	delete(c.unavailableGuilds, guildID)
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
			// update old guild_events
			return
		}
		guildMembers[member.ID] = member
	}
}
func (c CacheImpl) UncacheMember(guildID api.Snowflake, userID api.Snowflake) {
	delete(c.members[guildID], userID)
}
func (c CacheImpl) FindMember(guildID api.Snowflake, check func(u *api.Member) bool) *api.Member {
	for _, member := range c.Members(guildID) {
		if check(member) {
			return member
		}
	}
	return nil
}
func (c CacheImpl) FindMembers(guildID api.Snowflake, check func(u *api.Member) bool) []*api.Member {
	members := make([]*api.Member, 1)
	for _, member := range c.Members(guildID) {
		if check(member) {
			members = append(members, member)
		}
	}
	return members
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
func (c CacheImpl) FindRole(guildID api.Snowflake, check func(u *api.Role) bool) *api.Role {
	for _, role := range c.Roles(guildID) {
		if check(role) {
			return role
		}
	}
	return nil
}
func (c CacheImpl) FindRoles(guildID api.Snowflake, check func(u *api.Role) bool) []*api.Role {
	roles := make([]*api.Role, 1)
	for _, role := range c.Roles(guildID) {
		if check(role) {
			roles = append(roles, role)
		}
	}
	return roles
}

// other channel cache
func (c CacheImpl) Channel(channelID api.Snowflake) *api.Channel {
	dmChannel := c.DMChannel(channelID)
	if dmChannel != nil {
		return &dmChannel.Channel
	}
	category := c.Category(channelID)
	if category != nil {
		return &category.Channel
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return &textChannel.MessageChannel.Channel
	}
	voiceChannel := c.VoiceChannel(channelID)
	if voiceChannel != nil {
		return &voiceChannel.Channel
	}
	storeChannel := c.StoreChannel(channelID)
	if storeChannel != nil {
		return &storeChannel.Channel
	}
	return nil
}
func (c CacheImpl) MessageChannel(channelID api.Snowflake) *api.MessageChannel {
	dmChannel := c.DMChannel(channelID)
	if dmChannel != nil {
		return &dmChannel.MessageChannel
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return &textChannel.MessageChannel
	}
	return nil
}
func (c CacheImpl) GuildChannel(channelID api.Snowflake) *api.GuildChannel {
	category := c.Category(channelID)
	if category != nil {
		return &category.GuildChannel
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return &textChannel.GuildChannel
	}
	voiceChannel := c.VoiceChannel(channelID)
	if voiceChannel != nil {
		return &voiceChannel.GuildChannel
	}
	storeChannel := c.StoreChannel(channelID)
	if storeChannel != nil {
		return &storeChannel.GuildChannel
	}
	return nil
}

// dm channel cache
func (c CacheImpl) DMChannel(dmChannelID api.Snowflake) *api.DMChannel {
	return c.dmChannels[dmChannelID]
}
func (c CacheImpl) DMChannels() []*api.DMChannel {
	channels := make([]*api.DMChannel, len(c.dmChannels))
	i := 0
	for _, channel := range c.dmChannels {
		channels[i] = channel
		i++
	}
	return channels
}
func (c CacheImpl) DMChannelCache() map[api.Snowflake]*api.DMChannel {
	return c.dmChannels
}
func (c CacheImpl) CacheDMChannel(dmChannel *api.DMChannel) {
	if oldChannel, ok := c.dmChannels[dmChannel.ID]; ok {
		*oldChannel = *dmChannel
		return
	}
	c.dmChannels[dmChannel.ID] = dmChannel

}
func (c CacheImpl) UncacheDMChannel(channelID api.Snowflake) {
	delete(c.dmChannels, channelID)
}
func (c CacheImpl) FindDMChannel(check func(u *api.DMChannel) bool) *api.DMChannel {
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			return dmChannel
		}
	}
	return nil
}
func (c CacheImpl) FindDMChannels(check func(u *api.DMChannel) bool) []*api.DMChannel {
	dmChannels := make([]*api.DMChannel, 1)
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			dmChannels = append(dmChannels, dmChannel)
		}
	}
	return dmChannels
}

// text channel cache
func (c CacheImpl) TextChannel(textChannelID api.Snowflake) *api.TextChannel {
	for _, guild := range c.textChannels {
		if channel, ok := guild[textChannelID]; ok {
			return channel
		}
	}
	return nil
}
func (c CacheImpl) TextChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		textChannels := make([]*api.TextChannel, 1)
		for _, channel := range guildTextChannels {
			if ignoreCase && strings.ToLower(*channel.MessageChannel.Name) == name || !ignoreCase && *channel.MessageChannel.Name == name {
				textChannels = append(textChannels, channel)
			}
		}
		return textChannels
	}
	return nil
}
func (c CacheImpl) TextChannels(guildID api.Snowflake) []*api.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		textChannels := make([]*api.TextChannel, len(guildTextChannels))
		i := 0
		for _, textChannel := range guildTextChannels {
			textChannels[i] = textChannel
			i++
		}
		return textChannels
	}
	return nil
}
func (c CacheImpl) AllTextChannels() []*api.TextChannel {
	textChannels := make([]*api.TextChannel, len(c.textChannels))
	for _, guildTextChannels := range c.textChannels {
		for _, textChannel := range guildTextChannels {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}
func (c CacheImpl) TextChannelCache(guildID api.Snowflake) map[api.Snowflake]*api.TextChannel {
	return c.textChannels[guildID]
}
func (c CacheImpl) AllTextChannelCache() map[api.Snowflake]map[api.Snowflake]*api.TextChannel {
	return c.textChannels
}
func (c CacheImpl) CacheTextChannel(textChannel *api.TextChannel) {
	if guildTextChannels, ok := c.textChannels[textChannel.GuildID]; ok {
		if guildTextChannel, ok := guildTextChannels[textChannel.MessageChannel.ID]; ok {
			*guildTextChannel = *textChannel
			return
		}
		guildTextChannels[textChannel.MessageChannel.ID] = textChannel
	}
}
func (c CacheImpl) UncacheTextChannel(guildID api.Snowflake, textChannelID api.Snowflake) {
	delete(c.textChannels[guildID], textChannelID)
}
func (c CacheImpl) FindTextChannel(guildID api.Snowflake, check func(u *api.TextChannel) bool) *api.TextChannel {
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			return textChannel
		}
	}
	return nil
}
func (c CacheImpl) FindTextChannels(guildID api.Snowflake, check func(u *api.TextChannel) bool) []*api.TextChannel {
	textChannels := make([]*api.TextChannel, 1)
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}

// store channel cache
func (c CacheImpl) StoreChannel(storeChannelID api.Snowflake) *api.StoreChannel {
	for _, guild := range c.storeChannels {
		if channel, ok := guild[storeChannelID]; ok {
			return channel
		}
	}
	return nil
}
func (c CacheImpl) StoreChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		storeChannels := make([]*api.StoreChannel, 1)
		for _, channel := range guildStoreChannels {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				storeChannels = append(storeChannels, channel)
			}
		}
		return storeChannels
	}
	return nil
}
func (c CacheImpl) StoreChannels(guildID api.Snowflake) []*api.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		storeChannels := make([]*api.StoreChannel, len(guildStoreChannels))
		i := 0
		for _, storeChannel := range guildStoreChannels {
			storeChannels[i] = storeChannel
			i++
		}
		return storeChannels
	}
	return nil
}
func (c CacheImpl) AllStoreChannel() []*api.StoreChannel {
	storeChannels := make([]*api.StoreChannel, len(c.storeChannels))
	for _, guildStoreChannels := range c.storeChannels {
		for _, storeChannel := range guildStoreChannels {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}
func (c CacheImpl) StoreChannelCache(guildID api.Snowflake) map[api.Snowflake]*api.StoreChannel {
	return c.storeChannels[guildID]
}
func (c CacheImpl) AllStoreChannelCache() map[api.Snowflake]map[api.Snowflake]*api.StoreChannel {
	return c.storeChannels
}
func (c CacheImpl) CacheStoreChannel(storeChannel *api.StoreChannel) {
	if guildStoreChannels, ok := c.storeChannels[storeChannel.GuildID]; ok {
		if guildStoreChannel, ok := guildStoreChannels[storeChannel.ID]; ok {
			*guildStoreChannel = *storeChannel
			return
		}
		guildStoreChannels[storeChannel.ID] = storeChannel
	}
}
func (c CacheImpl) UncacheStoreChannel(guildID api.Snowflake, storeChannelID api.Snowflake) {
	delete(c.storeChannels[guildID], storeChannelID)
}
func (c CacheImpl) FindStoreChannel(guildID api.Snowflake, check func(u *api.StoreChannel) bool) *api.StoreChannel {
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			return storeChannel
		}
	}
	return nil
}
func (c CacheImpl) FindStoreChannels(guildID api.Snowflake, check func(u *api.StoreChannel) bool) []*api.StoreChannel {
	storeChannels := make([]*api.StoreChannel, 1)
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}

// voice channel cache
func (c CacheImpl) VoiceChannel(voiceChannelID api.Snowflake) *api.VoiceChannel {
	for _, guild := range c.voiceChannels {
		if channel, ok := guild[voiceChannelID]; ok {
			return channel
		}
	}
	return nil
}
func (c CacheImpl) VoiceChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		voiceChannels := make([]*api.VoiceChannel, 1)
		for _, channel := range guildVoiceChannels {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				voiceChannels = append(voiceChannels, channel)
			}
		}
		return voiceChannels
	}
	return nil
}
func (c CacheImpl) VoiceChannels(guildID api.Snowflake) []*api.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		voiceChannels := make([]*api.VoiceChannel, len(guildVoiceChannels))
		i := 0
		for _, voiceChannel := range guildVoiceChannels {
			voiceChannels[i] = voiceChannel
			i++
		}
		return voiceChannels
	}
	return nil
}
func (c CacheImpl) AllVoiceChannel() []*api.VoiceChannel {
	voiceChannels := make([]*api.VoiceChannel, len(c.voiceChannels))
	for _, guildVoiceChannels := range c.voiceChannels {
		for _, voiceChannel := range guildVoiceChannels {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}
func (c CacheImpl) VoiceChannelCache(guildID api.Snowflake) map[api.Snowflake]*api.VoiceChannel {
	return c.voiceChannels[guildID]
}
func (c CacheImpl) AllVoiceChannelCache() map[api.Snowflake]map[api.Snowflake]*api.VoiceChannel {
	return c.voiceChannels
}
func (c CacheImpl) CacheVoiceChannel(voiceChannel *api.VoiceChannel) {
	if guildVoiceChannels, ok := c.voiceChannels[voiceChannel.GuildID]; ok {
		if guildVoiceChannel, ok := guildVoiceChannels[voiceChannel.ID]; ok {
			*guildVoiceChannel = *voiceChannel
			return
		}
		guildVoiceChannels[voiceChannel.ID] = voiceChannel
	}
}
func (c CacheImpl) UncacheVoiceChannel(guildID api.Snowflake, voiceChannelID api.Snowflake) {
	delete(c.voiceChannels[guildID], voiceChannelID)
}
func (c CacheImpl) FindVoiceChannel(guildID api.Snowflake, check func(u *api.VoiceChannel) bool) *api.VoiceChannel {
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			return voiceChannel
		}
	}
	return nil
}
func (c CacheImpl) FindVoiceChannels(guildID api.Snowflake, check func(u *api.VoiceChannel) bool) []*api.VoiceChannel {
	voiceChannels := make([]*api.VoiceChannel, 1)
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}

// category channel cache
func (c CacheImpl) Category(categoryID api.Snowflake) *api.CategoryChannel {
	for _, guild := range c.categories {
		if channel, ok := guild[categoryID]; ok {
			return channel
		}
	}
	return nil
}
func (c CacheImpl) CategoriesByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.CategoryChannel {
	if guildCategories, ok := c.categories[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		categories := make([]*api.CategoryChannel, 1)
		for _, channel := range guildCategories {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				categories = append(categories, channel)
			}
		}
		return categories
	}
	return nil
}
func (c CacheImpl) Categories(guildID api.Snowflake) []*api.CategoryChannel {
	if guildCategories, ok := c.categories[guildID]; ok {
		categories := make([]*api.CategoryChannel, len(guildCategories))
		i := 0
		for _, category := range guildCategories {
			categories[i] = category
			i++
		}
		return categories
	}
	return nil
}
func (c CacheImpl) AllCategories() []*api.CategoryChannel {
	categories := make([]*api.CategoryChannel, len(c.categories))
	for _, guildCategories := range c.categories {
		for _, category := range guildCategories {
			categories = append(categories, category)
		}
	}
	return categories
}
func (c CacheImpl) CategoryCache(guildID api.Snowflake) map[api.Snowflake]*api.CategoryChannel {
	return c.categories[guildID]
}
func (c CacheImpl) AllCategoryCache() map[api.Snowflake]map[api.Snowflake]*api.CategoryChannel {
	return c.categories
}
func (c CacheImpl) CacheCategory(category *api.CategoryChannel) {
	if guildCategories, ok := c.categories[category.GuildID]; ok {
		if guildCategory, ok := guildCategories[category.ID]; ok {
			*guildCategory = *category
			return
		}
		guildCategories[category.ID] = category
	}
}
func (c CacheImpl) UncacheCategory(guildID api.Snowflake, categoryID api.Snowflake) {
	delete(c.categories[guildID], categoryID)
}
func (c CacheImpl) FindCategory(guildID api.Snowflake, check func(u *api.CategoryChannel) bool) *api.CategoryChannel {
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			return category
		}
	}
	return nil
}
func (c CacheImpl) FindCategories(guildID api.Snowflake, check func(u *api.CategoryChannel) bool) []*api.CategoryChannel {
	categories := make([]*api.CategoryChannel, 1)
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			categories = append(categories, category)
		}
	}
	return categories
}
