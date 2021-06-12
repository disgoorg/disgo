package internal

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/DisgoOrg/disgo/api"
)

func newCacheImpl(disgo api.Disgo, memberCachePolicy api.MemberCachePolicy, threadMemberCachePolicy api.ThreadMemberCachePolicy, messageCachePolicy api.MessageCachePolicy, cacheFlags api.CacheFlags) api.Cache {
	cache := &CacheImpl{
		disgo:                   disgo,
		quit:                    make(chan interface{}),
		memberCachePolicy:       memberCachePolicy,
		threadMemberCachePolicy: threadMemberCachePolicy,
		messageCachePolicy:      messageCachePolicy,
		cacheFlags:              cacheFlags,
		globalCommands:          map[api.Snowflake]*api.Command{},
		guildCommands:           map[api.Snowflake]map[api.Snowflake]*api.Command{},
		users:                   map[api.Snowflake]*api.User{},
		guilds:                  map[api.Snowflake]*api.Guild{},
		members:                 map[api.Snowflake]map[api.Snowflake]*api.Member{},
		threadMembers:           map[api.Snowflake]map[api.Snowflake]map[api.Snowflake]*api.ThreadMember{},
		voiceStates:             map[api.Snowflake]map[api.Snowflake]*api.VoiceState{},
		roles:                   map[api.Snowflake]map[api.Snowflake]*api.Role{},
		dmChannels:              map[api.Snowflake]api.DMChannel{},
		categories:              map[api.Snowflake]map[api.Snowflake]api.Category{},
		textChannels:            map[api.Snowflake]map[api.Snowflake]api.TextChannel{},
		voiceChannels:           map[api.Snowflake]map[api.Snowflake]api.VoiceChannel{},
		storeChannels:           map[api.Snowflake]map[api.Snowflake]api.StoreChannel{},
		emojis:                  map[api.Snowflake]map[api.Snowflake]*api.Emoji{},
	}
	go cache.startCleanup(10 * time.Second)
	return cache
}

// CacheImpl is used for api.Disgo's api.Cache
type CacheImpl struct {
	disgo                   api.Disgo
	quit                    chan interface{}
	memberCachePolicy       api.MemberCachePolicy
	threadMemberCachePolicy api.ThreadMemberCachePolicy
	messageCachePolicy      api.MessageCachePolicy
	cacheFlags              api.CacheFlags
	globalCommands          map[api.Snowflake]*api.Command
	guildCommands           map[api.Snowflake]map[api.Snowflake]*api.Command
	users                   map[api.Snowflake]*api.User
	guilds                  map[api.Snowflake]*api.Guild
	messages                map[api.Snowflake]map[api.Snowflake]*api.Message
	members                 map[api.Snowflake]map[api.Snowflake]*api.Member
	threadMembers           map[api.Snowflake]map[api.Snowflake]map[api.Snowflake]*api.ThreadMember
	voiceStates             map[api.Snowflake]map[api.Snowflake]*api.VoiceState
	roles                   map[api.Snowflake]map[api.Snowflake]*api.Role
	dmChannels              map[api.Snowflake]api.DMChannel
	categories              map[api.Snowflake]map[api.Snowflake]api.Category
	textChannels            map[api.Snowflake]map[api.Snowflake]api.TextChannel
	threads                 map[api.Snowflake]map[api.Snowflake]api.Thread
	voiceChannels           map[api.Snowflake]map[api.Snowflake]api.VoiceChannel
	storeChannels           map[api.Snowflake]map[api.Snowflake]api.StoreChannel
	emojis                  map[api.Snowflake]map[api.Snowflake]*api.Emoji
}

// Disgo returns the current api.Disgo instance
func (c *CacheImpl) Disgo() api.Disgo {
	return c.disgo
}

// Close cleans up the cache and it's internal tasks
func (c *CacheImpl) Close() {
	c.Disgo().Logger().Info("closing cache goroutines...")
	c.quit <- true
	c.Disgo().Logger().Info("closed cache goroutines")
}

func (c CacheImpl) startCleanup(cleanupInterval time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			c.Disgo().Logger().Panicf("recovered cache cleanup goroutine error: %s", r)
			debug.PrintStack()
			c.startCleanup(cleanupInterval)
			return
		}
		c.Disgo().Logger().Info("shut down cache cleanup goroutine")
	}()

	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.DoCleanup()
		case <-c.quit:
			ticker.Stop()
			return
		}
	}
}

// DoCleanup removes items from the cache that no longer meet their policy
func (c *CacheImpl) DoCleanup() {
	for channelID, channelMessages := range c.messages {
		for _, message := range channelMessages {
			if !c.messageCachePolicy(message) {
				c.UncacheMessage(channelID, message.ID)
			}
		}
	}
}

// CacheFlags returns the current active api.CacheFlags
func (c CacheImpl) CacheFlags() api.CacheFlags {
	return c.cacheFlags
}

// Command returns an api.Command from cache
func (c *CacheImpl) Command(commandID api.Snowflake) *api.Command {
	if command, ok := c.globalCommands[commandID]; ok {
		return command
	}
	for _, guildCommands := range c.guildCommands {
		if command, ok := guildCommands[commandID]; ok {
			return command
		}
	}
	return nil
}

// GuildCommandCache returns the cache of commands in a Guild
func (c *CacheImpl) GuildCommandCache(guildID api.Snowflake) map[api.Snowflake]*api.Command {
	return c.guildCommands[guildID]
}

// AllGuildCommandCache returns the cache of all Guild Command(s)
func (c *CacheImpl) AllGuildCommandCache() map[api.Snowflake]map[api.Snowflake]*api.Command {
	return c.guildCommands
}

// GlobalCommandCache returns the cache of global Command(s)
func (c *CacheImpl) GlobalCommandCache() map[api.Snowflake]*api.Command {
	return c.globalCommands
}

// CacheGlobalCommand adds a global command to the cache
func (c *CacheImpl) CacheGlobalCommand(command *api.Command) *api.Command {
	if c.CacheFlags().Missing(api.CacheFlagCommands) {
		return command
	}
	if _, ok := c.globalCommands[command.ID]; ok {
		*c.globalCommands[command.ID] = *command
		return c.globalCommands[command.ID]
	}
	c.globalCommands[command.ID] = command
	return command
}

// CacheGuildCommand adds a Guild Command to the cache
func (c *CacheImpl) CacheGuildCommand(command *api.Command) *api.Command {
	if c.CacheFlags().Missing(api.CacheFlagCommands) {
		return command
	}
	if guildCommands, ok := c.guildCommands[command.ID]; ok {
		if _, ok = guildCommands[command.ID]; ok {
			*guildCommands[command.ID] = *command
			return guildCommands[command.ID]
		}
		guildCommands[command.ID] = command
	}
	return command
}

// UncacheCommand removes a global Command from the cache
func (c *CacheImpl) UncacheCommand(commandID api.Snowflake) {
	if _, ok := c.globalCommands[commandID]; ok {
		delete(c.globalCommands, commandID)
		return
	}
	for _, guildCommands := range c.guildCommands {
		if _, ok := guildCommands[commandID]; ok {
			delete(c.guildCommands, commandID)
			return
		}
	}
}

// User allows you to get a user from the cache by ID
func (c *CacheImpl) User(id api.Snowflake) *api.User {
	return c.users[id]
}

// UserByTag allows you to get a user from the cache by their Tag
func (c *CacheImpl) UserByTag(tag string) *api.User {
	for _, user := range c.users {
		if user.Tag() == tag {
			return user
		}
	}
	return nil
}

// UsersByName allows you to get users from the cache by username
func (c *CacheImpl) UsersByName(name string, ignoreCase bool) []*api.User {
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

// Users returns all users from the cache as a slice
func (c *CacheImpl) Users() []*api.User {
	users := make([]*api.User, len(c.users))
	i := 0
	for _, user := range c.users {
		users[i] = user
		i++
	}
	return users
}

// UserCache returns all users from the cache a map
func (c *CacheImpl) UserCache() map[api.Snowflake]*api.User {
	return c.users
}

// CacheUser adds a user to the cache
func (c *CacheImpl) CacheUser(user *api.User) *api.User {
	// TODO: only cache user if we have a mutal guild & always cache self user
	if globalUser, ok := c.users[user.ID]; ok {
		*globalUser = *user
		return globalUser
	}
	c.users[user.ID] = user
	return user
}

// UncacheUser removes a user from the cache
func (c *CacheImpl) UncacheUser(id api.Snowflake) {
	delete(c.users, id)
}

// FindUser finds a user from the cache with a custom method
func (c *CacheImpl) FindUser(check func(u *api.User) bool) *api.User {
	for _, user := range c.users {
		if check(user) {
			return user
		}
	}
	return nil
}

// FindUsers finds several users from the cache with a custom method
func (c *CacheImpl) FindUsers(check func(u *api.User) bool) []*api.User {
	users := make([]*api.User, 1)
	for _, user := range c.users {
		if check(user) {
			users = append(users, user)
		}
	}
	return users
}

// Guild allows you to get a guild from the cache by ID
func (c *CacheImpl) Guild(guildID api.Snowflake) *api.Guild {
	return c.guilds[guildID]
}

// GuildsByName allows you to get guilds from the cache by name
func (c *CacheImpl) GuildsByName(name string, ignoreCase bool) []*api.Guild {
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

// Guilds returns the guild cache as a slice
func (c *CacheImpl) Guilds() []*api.Guild {
	guilds := make([]*api.Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}

// GuildCache returns the guild cache as a map
func (c *CacheImpl) GuildCache() map[api.Snowflake]*api.Guild {
	return c.guilds
}

// CacheGuild adds a guild to the cache
func (c *CacheImpl) CacheGuild(guild *api.Guild) *api.Guild {
	if _, ok := c.guilds[guild.ID]; ok {
		// update old guild
		*c.guilds[guild.ID] = *guild
		return c.guilds[guild.ID]
	}
	// guild was not yet cached so cache it directly
	c.guilds[guild.ID] = guild
	c.guildCommands[guild.ID] = map[api.Snowflake]*api.Command{}
	c.members[guild.ID] = map[api.Snowflake]*api.Member{}
	c.voiceStates[guild.ID] = map[api.Snowflake]*api.VoiceState{}
	c.roles[guild.ID] = map[api.Snowflake]*api.Role{}
	c.categories[guild.ID] = map[api.Snowflake]api.Category{}
	c.textChannels[guild.ID] = map[api.Snowflake]api.TextChannel{}
	c.voiceChannels[guild.ID] = map[api.Snowflake]api.VoiceChannel{}
	c.storeChannels[guild.ID] = map[api.Snowflake]api.StoreChannel{}
	return guild
}

//UncacheGuild removes a guild and all of it's children from the cache
func (c *CacheImpl) UncacheGuild(guildID api.Snowflake) {
	delete(c.guilds, guildID)
	delete(c.guildCommands, guildID)
	delete(c.members, guildID)
	delete(c.voiceStates, guildID)
	delete(c.roles, guildID)
	delete(c.categories, guildID)
	delete(c.textChannels, guildID)
	delete(c.voiceChannels, guildID)
	delete(c.storeChannels, guildID)
}

//FindGuild finds a guild by a custom method
func (c *CacheImpl) FindGuild(check func(g *api.Guild) bool) *api.Guild {
	for _, guild := range c.guilds {
		if check(guild) {
			return guild
		}
	}
	return nil
}

//FindGuilds finds multiple guilds with a custom method
func (c *CacheImpl) FindGuilds(check func(g *api.Guild) bool) []*api.Guild {
	guilds := make([]*api.Guild, 1)
	for _, guild := range c.guilds {
		if check(guild) {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}

// Message returns a message from cache
func (c *CacheImpl) Message(channelID api.Snowflake, messageID api.Snowflake) *api.Message {
	if channelMessages, ok := c.messages[channelID]; ok {
		return channelMessages[messageID]
	}
	return nil
}

// Messages returns the messages of a channel from cache
func (c *CacheImpl) Messages(channelID api.Snowflake) []*api.Message {
	if channelMessages, ok := c.messages[channelID]; ok {
		messages := make([]*api.Message, len(channelMessages))
		i := 0
		for _, message := range channelMessages {
			messages[i] = message
			i++
		}
		return messages
	}
	return nil
}

// MessageCache returns the entire cache of Message(s) for a Channel
func (c *CacheImpl) MessageCache(channelID api.Snowflake) map[api.Snowflake]*api.Message {
	return c.messages[channelID]
}

// AllMessageCache returns the entire cache of messages
func (c *CacheImpl) AllMessageCache() map[api.Snowflake]map[api.Snowflake]*api.Message {
	return c.messages
}

//CacheMessage adds a message to the cache
func (c *CacheImpl) CacheMessage(message *api.Message) *api.Message {
	// only cache message if we want to
	if !c.messageCachePolicy(message) {
		return message
	}
	if channelMessages, ok := c.messages[message.ChannelID]; ok {
		if channelMessage, ok := channelMessages[message.ID]; ok {
			*channelMessage = *message
			return channelMessage
		}
		channelMessages[message.ID] = message
	}
	return message
}

// UncacheMessage removes a Message from the cache if the cache policy allows it
func (c *CacheImpl) UncacheMessage(channelID api.Snowflake, messageID api.Snowflake) {
	if channelMessages, ok := c.messages[channelID]; ok {
		if message, ok := channelMessages[messageID]; ok {
			// check if we really want to uncache that message according to our policy
			if !c.messageCachePolicy(message) {
				delete(channelMessages, messageID)
			}
		}
	}
}

// Member returns a member from cache by guild ID and user ID
func (c *CacheImpl) Member(guildID api.Snowflake, userID api.Snowflake) *api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		return guildMembers[userID]
	}
	return nil
}

// MemberByTag returns a member from cache by guild ID and user tag
func (c *CacheImpl) MemberByTag(guildID api.Snowflake, tag string) *api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		for _, member := range guildMembers {
			if member.User.Tag() == tag {
				return member
			}
		}
	}
	return nil
}

// MembersByName returns members from cache by guild ID and username
func (c *CacheImpl) MembersByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		members := make([]*api.Member, 1)
		for _, member := range guildMembers {
			if ignoreCase && strings.ToLower(member.User.Username) == name || !ignoreCase && member.User.Username == name {
				members = append(members, member)
			}
		}
		return members
	}
	return nil
}

// Members returns the member cache of a guild by snowflake
func (c *CacheImpl) Members(guildID api.Snowflake) []*api.Member {
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

// AllMembers returns the entire cache of members
func (c *CacheImpl) AllMembers() []*api.Member {
	members := make([]*api.Member, len(c.guilds))
	for _, guildMembers := range c.members {
		for _, member := range guildMembers {
			members = append(members, member)
		}
	}
	return members
}

//MemberCache returns the cache of a guild as a map
func (c *CacheImpl) MemberCache(guildID api.Snowflake) map[api.Snowflake]*api.Member {
	return c.members[guildID]
}

// AllMemberCache returns the entire cache as a map of maps
func (c *CacheImpl) AllMemberCache() map[api.Snowflake]map[api.Snowflake]*api.Member {
	return c.members
}

// CacheMember adds a member to the cache
func (c *CacheImpl) CacheMember(member *api.Member) *api.Member {
	// only cache member if we want to & always cache self member!
	if !c.memberCachePolicy(member) && member.User.ID != member.Disgo.ApplicationID() {
		return member
	}
	if guildMembers, ok := c.members[member.GuildID]; ok {
		if guildMember, ok := guildMembers[member.User.ID]; ok {
			*guildMember = *member
			return guildMember
		}
		guildMembers[member.User.ID] = member
	}
	return member
}

// UncacheMember removes a guild member from the cache
func (c *CacheImpl) UncacheMember(guildID api.Snowflake, userID api.Snowflake) {
	// TODO: add UncacheUser call!
	if guildMembers, ok := c.members[guildID]; ok {
		if member, ok := guildMembers[userID]; ok {
			// check if we really want to uncache that member according to our policy
			if !c.memberCachePolicy(member) {
				delete(guildMembers, userID)
			}
		}
	}
}

// FindMember allows you to find a member in a guild by custom method
func (c *CacheImpl) FindMember(guildID api.Snowflake, check func(u *api.Member) bool) *api.Member {
	for _, member := range c.Members(guildID) {
		if check(member) {
			return member
		}
	}
	return nil
}

// FindMembers allows you to find api.Member(s) in a guild by custom method
func (c *CacheImpl) FindMembers(guildID api.Snowflake, check func(u *api.Member) bool) []*api.Member {
	members := make([]*api.Member, 1)
	for _, member := range c.Members(guildID) {
		if check(member) {
			members = append(members, member)
		}
	}
	return members
}

// ThreadMember returns a api.ThreadMember from api.Cache by guild ID and thread ID
func (c *CacheImpl) ThreadMember(guildID api.Snowflake, threadID api.Snowflake, userID api.Snowflake) *api.ThreadMember {
	if guildThreadMembers, ok := c.threadMembers[guildID]; ok {
		if threadMembers, ok := guildThreadMembers[threadID]; ok {
			return threadMembers[userID]
		}
	}
	return nil
}

func (c *CacheImpl) ThreadMembers(guildID api.Snowflake, threadID api.Snowflake) []*api.ThreadMember {
	if guildThreadMembers, ok := c.threadMembers[guildID]; ok {
		if threadMembers, ok := guildThreadMembers[threadID]; ok {
			members := make([]*api.ThreadMember, len(threadMembers))
			i := 0
			for _, member := range threadMembers {
				members[i] = member
				i++
			}
			return members
		}
	}
	return nil
}

func (c *CacheImpl) ThreadMemberCache(guildID api.Snowflake) map[api.Snowflake]map[api.Snowflake]*api.ThreadMember {
	if guildThreadMembers, ok := c.threadMembers[guildID]; ok {
		return guildThreadMembers
	}
	return nil
}

func (c *CacheImpl) AllThreadMemberCache() map[api.Snowflake]map[api.Snowflake]map[api.Snowflake]*api.ThreadMember {
	return c.threadMembers
}

func (c *CacheImpl) CacheThreadMember(member *api.ThreadMember) *api.ThreadMember {
	// only cache member if we want to & always cache self ThreadMember!
	if !c.threadMemberCachePolicy(member) && member.UserID != member.Disgo.ApplicationID() {
		return member
	}
	if guildThreadMembers, ok := c.threadMembers[member.GuildID]; ok {
		if threadMembers, ok := guildThreadMembers[member.ThreadID]; ok {
			if threadMember, ok := threadMembers[member.UserID]; ok {
				*threadMember = *member
				return threadMember
			}
			threadMembers[member.UserID] = member
		}
	}
	return member
}
func (c *CacheImpl) UncacheThreadMember(guildID api.Snowflake, threadID api.Snowflake, userID api.Snowflake) {
	// TODO: add UncacheUser call?
	if guildThreadMembers, ok := c.threadMembers[guildID]; ok {
		if threadMembers, ok := guildThreadMembers[threadID]; ok {
			if threadMember, ok := threadMembers[userID]; ok {
				// check if we really want to uncache that member according to our policy
				if !c.threadMemberCachePolicy(threadMember) {
					delete(threadMembers, userID)
				}
			}
		}
	}
}

// UncacheThreadMembers removes all ThreadsMembers for a Guild from the Cache
func (c *CacheImpl) UncacheThreadMembers(guildID api.Snowflake) {
	delete(c.threadMembers, guildID)
}

// VoiceState returns a Member's api.VoiceState for a api.Guild
func (c *CacheImpl) VoiceState(guildID api.Snowflake, userID api.Snowflake) *api.VoiceState {
	if voiceStates, ok := c.voiceStates[guildID]; ok {
		return voiceStates[userID]
	}
	return nil
}

// VoiceStates returns the member cache of a guild by snowflake
func (c *CacheImpl) VoiceStates(guildID api.Snowflake) []*api.VoiceState {
	if guildVoiceStates, ok := c.voiceStates[guildID]; ok {
		voiceStates := make([]*api.VoiceState, len(guildVoiceStates))
		i := 0
		for _, voiceState := range guildVoiceStates {
			voiceStates[i] = voiceState
			i++
		}
		return voiceStates
	}
	return nil
}

// VoiceStateCache returns the api.VoiceState api.Cache of a api.Guild as a map
func (c *CacheImpl) VoiceStateCache(guildID api.Snowflake) map[api.Snowflake]*api.VoiceState {
	return c.voiceStates[guildID]
}

// CacheVoiceState adds a api.VoiceState from the api.Cache
func (c *CacheImpl) CacheVoiceState(voiceState *api.VoiceState) *api.VoiceState {
	// only cache voice states for ourself or member is cached & cache flag activated
	if c.cacheFlags.Missing(api.CacheFlagVoiceState) && voiceState.UserID != c.disgo.ApplicationID() {
		return voiceState
	}
	if guildVoiceStates, ok := c.voiceStates[voiceState.GuildID]; ok {
		if guildVoiceState, ok := guildVoiceStates[voiceState.UserID]; ok {
			*guildVoiceState = *voiceState
			return guildVoiceState
		}
		guildVoiceStates[voiceState.UserID] = voiceState
	}
	return voiceState
}

// UncacheVoiceState removes a api.VoiceState from the api.Cache
func (c *CacheImpl) UncacheVoiceState(guildID api.Snowflake, userID api.Snowflake) {
	delete(c.voiceStates[guildID], userID)
}

// Role returns a role from cache by guild ID and role ID
func (c *CacheImpl) Role(roleID api.Snowflake) *api.Role {
	for _, guildRoles := range c.roles {
		if role, ok := guildRoles[roleID]; ok {
			return role
		}
	}
	return nil
}

// RolesByName returns roles from cache by guild ID and name
func (c *CacheImpl) RolesByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.Role {
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

// Roles returns the role cache of a guild
func (c *CacheImpl) Roles(guildID api.Snowflake) []*api.Role {
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

// AllRoles returns the entire role cache
func (c *CacheImpl) AllRoles() []*api.Role {
	roles := make([]*api.Role, len(c.guilds))
	for _, guildRoles := range c.roles {
		for _, role := range guildRoles {
			roles = append(roles, role)
		}
	}
	return roles
}

// RoleCache returns the role cache of a guild by ID
func (c *CacheImpl) RoleCache(guildID api.Snowflake) map[api.Snowflake]*api.Role {
	return c.roles[guildID]
}

// AllRoleCache returns the entire role cache
func (c *CacheImpl) AllRoleCache() map[api.Snowflake]map[api.Snowflake]*api.Role {
	return c.roles
}

// CacheRole adds a role to the cache
func (c *CacheImpl) CacheRole(role *api.Role) *api.Role {
	if c.cacheFlags.Missing(api.CacheFlagRoles) {
		return role
	}
	if guildRoles, ok := c.roles[role.GuildID]; ok {
		if guildRole, ok := guildRoles[role.ID]; ok {
			*guildRole = *role
			return guildRole
		}
		guildRoles[role.ID] = role
	}
	return role
}

// UncacheRole removes a role from cache
func (c *CacheImpl) UncacheRole(guildID api.Snowflake, roleID api.Snowflake) {
	delete(c.roles[guildID], roleID)
}

// FindRole allows you to find a role in a guild by custom method
func (c *CacheImpl) FindRole(guildID api.Snowflake, check func(u *api.Role) bool) *api.Role {
	for _, role := range c.Roles(guildID) {
		if check(role) {
			return role
		}
	}
	return nil
}

// FindRoles allows you to find roles in a guild by custom method
func (c *CacheImpl) FindRoles(guildID api.Snowflake, check func(u *api.Role) bool) []*api.Role {
	roles := make([]*api.Role, 1)
	for _, role := range c.Roles(guildID) {
		if check(role) {
			roles = append(roles, role)
		}
	}
	return roles
}

// Channel returns a channel from any channel cache by ID
func (c *CacheImpl) Channel(channelID api.Snowflake) api.Channel {
	dmChannel := c.DMChannel(channelID)
	if dmChannel != nil {
		return dmChannel
	}
	category := c.Category(channelID)
	if category != nil {
		return category
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return textChannel
	}
	voiceChannel := c.VoiceChannel(channelID)
	if voiceChannel != nil {
		return voiceChannel
	}
	storeChannel := c.StoreChannel(channelID)
	if storeChannel != nil {
		return storeChannel
	}
	thread := c.Thread(channelID)
	if thread != nil {
		return thread
	}
	return nil
}

// MessageChannel returns a channel from dm or text channel cache by ID
func (c *CacheImpl) MessageChannel(channelID api.Snowflake) api.MessageChannel {
	dmChannel := c.DMChannel(channelID)
	if dmChannel != nil {
		return dmChannel
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return textChannel
	}
	thread := c.Thread(channelID)
	if thread != nil {
		return thread
	}
	return nil
}

// GuildChannel returns a channel from a guild by ID
func (c *CacheImpl) GuildChannel(channelID api.Snowflake) api.GuildChannel {
	category := c.Category(channelID)
	if category != nil {
		return category
	}
	textChannel := c.TextChannel(channelID)
	if textChannel != nil {
		return textChannel
	}
	voiceChannel := c.VoiceChannel(channelID)
	if voiceChannel != nil {
		return voiceChannel
	}
	storeChannel := c.StoreChannel(channelID)
	if storeChannel != nil {
		return storeChannel
	}
	thread := c.Thread(channelID)
	if thread != nil {
		return thread
	}
	return nil
}

// DMChannel returns a DM channel by ID
func (c *CacheImpl) DMChannel(dmChannelID api.Snowflake) api.DMChannel {
	return c.dmChannels[dmChannelID]
}

// DMChannels return all DM channels as a slice
func (c *CacheImpl) DMChannels() []api.DMChannel {
	channels := make([]api.DMChannel, len(c.dmChannels))
	i := 0
	for _, channel := range c.dmChannels {
		channels[i] = channel
		i++
	}
	return channels
}

// DMChannelCache returns the DM channels as a map
func (c *CacheImpl) DMChannelCache() map[api.Snowflake]api.DMChannel {
	return c.dmChannels
}

// CacheDMChannel adds a DM channel to the cache
func (c *CacheImpl) CacheDMChannel(dmChannel api.DMChannel) api.DMChannel {
	if c.cacheFlags.Missing(api.CacheFlagDMChannels) {
		return dmChannel
	}
	if oldChannel, ok := c.dmChannels[dmChannel.ID()]; ok {
		*oldChannel.(*api.ChannelImpl) = *dmChannel.(*api.ChannelImpl)
		return oldChannel
	}
	c.dmChannels[dmChannel.ID()] = dmChannel
	return dmChannel
}

// UncacheDMChannel removes a DM channel from cache
func (c *CacheImpl) UncacheDMChannel(dmChannelID api.Snowflake) {
	// TODO: check this
	// should be okay to just uncache all messages if the channel gets uncached as that should mean it got deleted
	if _, ok := c.messages[dmChannelID]; ok {
		delete(c.messages, dmChannelID)
	}
	delete(c.dmChannels, dmChannelID)
}

// FindDMChannel finds a DM channel in cache with a custom method
func (c *CacheImpl) FindDMChannel(check func(u api.DMChannel) bool) api.DMChannel {
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			return dmChannel
		}
	}
	return nil
}

// FindDMChannels finds DM Channels in cache with a custom method
func (c *CacheImpl) FindDMChannels(check func(u api.DMChannel) bool) []api.DMChannel {
	dmChannels := make([]api.DMChannel, 1)
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			dmChannels = append(dmChannels, dmChannel)
		}
	}
	return dmChannels
}

// TextChannel returns a text channel from cache by ID
func (c *CacheImpl) TextChannel(textChannelID api.Snowflake) api.TextChannel {
	for _, guild := range c.textChannels {
		if channel, ok := guild[textChannelID]; ok {
			return channel
		}
	}
	return nil
}

// TextChannelsByName returns text channel from cache by guild ID and name
func (c *CacheImpl) TextChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []api.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		textChannels := make([]api.TextChannel, 1)
		for _, channel := range guildTextChannels {
			if ignoreCase && strings.ToLower(channel.Name()) == name || !ignoreCase && channel.Name() == name {
				textChannels = append(textChannels, channel)
			}
		}
		return textChannels
	}
	return nil
}

// TextChannels returns the text channel cache from a guild
func (c *CacheImpl) TextChannels(guildID api.Snowflake) []api.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		textChannels := make([]api.TextChannel, len(guildTextChannels))
		i := 0
		for _, textChannel := range guildTextChannels {
			textChannels[i] = textChannel
			i++
		}
		return textChannels
	}
	return nil
}

// AllTextChannels returns the text channel cache as a slice
func (c *CacheImpl) AllTextChannels() []api.TextChannel {
	textChannels := make([]api.TextChannel, len(c.textChannels))
	for _, guildTextChannels := range c.textChannels {
		for _, textChannel := range guildTextChannels {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}

// TextChannelCache returns the channel cache as a map
func (c *CacheImpl) TextChannelCache(guildID api.Snowflake) map[api.Snowflake]api.TextChannel {
	return c.textChannels[guildID]
}

// AllTextChannelCache returns the text channel cache as a map of maps
func (c *CacheImpl) AllTextChannelCache() map[api.Snowflake]map[api.Snowflake]api.TextChannel {
	return c.textChannels
}

// CacheTextChannel adds a channel to the cache
func (c *CacheImpl) CacheTextChannel(textChannel api.TextChannel) api.TextChannel {
	if c.cacheFlags.Missing(api.CacheFlagTextChannels) {
		return textChannel
	}
	if guildTextChannels, ok := c.textChannels[textChannel.GuildID()]; ok {
		if guildTextChannel, ok := guildTextChannels[textChannel.ID()]; ok {
			*guildTextChannel.(*api.ChannelImpl) = *textChannel.(*api.ChannelImpl)
			return guildTextChannel
		}
		guildTextChannels[textChannel.ID()] = textChannel
	}
	return textChannel
}

// UncacheTextChannel removes a text channel from the cache
func (c *CacheImpl) UncacheTextChannel(guildID api.Snowflake, textChannelID api.Snowflake) {
	// TODO: check this
	// should be okay to just uncache all messages if the channel gets uncached as that should mean it got deleted
	if _, ok := c.messages[textChannelID]; ok {
		delete(c.messages, textChannelID)
	}
	delete(c.textChannels[guildID], textChannelID)
}

// FindTextChannel finds a text channel in a guild by custom method
func (c *CacheImpl) FindTextChannel(guildID api.Snowflake, check func(u api.TextChannel) bool) api.TextChannel {
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			return textChannel
		}
	}
	return nil
}

// FindTextChannels finds text channels in a guild by custom method
func (c *CacheImpl) FindTextChannels(guildID api.Snowflake, check func(u api.TextChannel) bool) []api.TextChannel {
	textChannels := make([]api.TextChannel, 1)
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}

// Thread returns a thread from cache by ID
func (c *CacheImpl) Thread(threadID api.Snowflake) api.Thread {
	for _, guild := range c.threads {
		if channel, ok := guild[threadID]; ok {
			return channel
		}
	}
	return nil
}

// ThreadsByName returns text channel from cache by guild ID and name
func (c *CacheImpl) ThreadsByName(guildID api.Snowflake, name string, ignoreCase bool) []api.Thread {
	if guildThreads, ok := c.threads[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		threads := make([]api.Thread, 1)
		for _, channel := range guildThreads {
			if ignoreCase && strings.ToLower(channel.Name()) == name || !ignoreCase && channel.Name() == name {
				threads = append(threads, channel)
			}
		}
		return threads
	}
	return nil
}

// Threads returns the text channel cache from a guild
func (c *CacheImpl) Threads(guildID api.Snowflake) []api.Thread {
	if guildThreads, ok := c.threads[guildID]; ok {
		threads := make([]api.Thread, len(guildThreads))
		i := 0
		for _, thread := range guildThreads {
			threads[i] = thread
			i++
		}
		return threads
	}
	return nil
}

// ThreadCache returns the channel cache as a map
func (c *CacheImpl) ThreadCache(guildID api.Snowflake) map[api.Snowflake]api.Thread {
	return c.threads[guildID]
}

// AllThreadCache returns the thread cache as a map
func (c *CacheImpl) AllThreadCache() map[api.Snowflake]map[api.Snowflake]api.Thread {
	return c.threads
}

// CacheThread adds a channel to the cache
func (c *CacheImpl) CacheThread(thread api.Thread) api.Thread {
	if c.cacheFlags.Missing(api.CacheFlagThreads) {
		return thread
	}
	if guildThreads, ok := c.threads[thread.GuildID()]; ok {
		if guildThread, ok := guildThreads[thread.ID()]; ok {
			*guildThread.(*api.ChannelImpl) = *thread.(*api.ChannelImpl)
			return guildThread
		}
		guildThreads[thread.ID()] = thread
	}
	return thread
}

// UncacheThread removes a text channel from the cache
func (c *CacheImpl) UncacheThread(guildID api.Snowflake, threadID api.Snowflake) {
	// TODO: check this
	// should be okay to just uncache all messages if the channel gets uncached as that should mean it got deleted
	if _, ok := c.messages[threadID]; ok {
		delete(c.messages, threadID)
	}
	delete(c.threads[guildID], threadID)
}

// UncacheThreads removes all Threads for a Guild from the Cache
func (c *CacheImpl) UncacheThreads(guildID api.Snowflake) {
	delete(c.threads, guildID)
}

// FindThread finds a text channel in a guild by custom method
func (c *CacheImpl) FindThread(guildID api.Snowflake, check func(u api.Thread) bool) api.Thread {
	for _, thread := range c.ThreadCache(guildID) {
		if check(thread) {
			return thread
		}
	}
	return nil
}

// FindThreads finds text channels in a guild by custom method
func (c *CacheImpl) FindThreads(guildID api.Snowflake, check func(u api.Thread) bool) []api.Thread {
	threads := make([]api.Thread, 1)
	for _, thread := range c.ThreadCache(guildID) {
		if check(thread) {
			threads = append(threads, thread)
		}
	}
	return threads
}

//StoreChannel returns a store channel from cache by ID
func (c *CacheImpl) StoreChannel(storeChannelID api.Snowflake) api.StoreChannel {
	for _, guild := range c.storeChannels {
		if channel, ok := guild[storeChannelID]; ok {
			return channel
		}
	}
	return nil
}

//StoreChannelsByName returns store channels from cache by name
func (c *CacheImpl) StoreChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []api.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		storeChannels := make([]api.StoreChannel, 1)
		for _, channel := range guildStoreChannels {
			if ignoreCase && strings.ToLower(channel.Name()) == name || !ignoreCase && channel.Name() == name {
				storeChannels = append(storeChannels, channel)
			}
		}
		return storeChannels
	}
	return nil
}

//StoreChannels returns store channels from cache by guild ID
func (c *CacheImpl) StoreChannels(guildID api.Snowflake) []api.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		storeChannels := make([]api.StoreChannel, len(guildStoreChannels))
		i := 0
		for _, storeChannel := range guildStoreChannels {
			storeChannels[i] = storeChannel
			i++
		}
		return storeChannels
	}
	return nil
}

// AllStoreChannels returns all store channels from cache as a map
func (c *CacheImpl) AllStoreChannels() []api.StoreChannel {
	storeChannels := make([]api.StoreChannel, len(c.storeChannels))
	for _, guildStoreChannels := range c.storeChannels {
		for _, storeChannel := range guildStoreChannels {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}

//StoreChannelCache returns the store channels of a guild by ID
func (c *CacheImpl) StoreChannelCache(guildID api.Snowflake) map[api.Snowflake]api.StoreChannel {
	return c.storeChannels[guildID]
}

//AllStoreChannelCache returns all store channels from cache as a map of maps
func (c *CacheImpl) AllStoreChannelCache() map[api.Snowflake]map[api.Snowflake]api.StoreChannel {
	return c.storeChannels
}

// CacheStoreChannel adds a store channel to the cache
func (c *CacheImpl) CacheStoreChannel(storeChannel api.StoreChannel) api.StoreChannel {
	if c.cacheFlags.Missing(api.CacheFlagStoreChannels) {
		return storeChannel
	}
	if guildStoreChannels, ok := c.storeChannels[storeChannel.GuildID()]; ok {
		if guildStoreChannel, ok := guildStoreChannels[storeChannel.ID()]; ok {
			*guildStoreChannel.(*api.ChannelImpl) = *storeChannel.(*api.ChannelImpl)
			return guildStoreChannel
		}
		guildStoreChannels[storeChannel.ID()] = storeChannel
	}
	return storeChannel
}

// UncacheStoreChannel removes a store channel from cache
func (c *CacheImpl) UncacheStoreChannel(guildID api.Snowflake, storeChannelID api.Snowflake) {
	delete(c.storeChannels[guildID], storeChannelID)
}

// FindStoreChannel returns a store channel from cache by custom method
func (c *CacheImpl) FindStoreChannel(guildID api.Snowflake, check func(u api.StoreChannel) bool) api.StoreChannel {
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			return storeChannel
		}
	}
	return nil
}

// FindStoreChannels returns store channels from cache by custom method
func (c *CacheImpl) FindStoreChannels(guildID api.Snowflake, check func(u api.StoreChannel) bool) []api.StoreChannel {
	storeChannels := make([]api.StoreChannel, 1)
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}

// VoiceChannel returns a voice channel from cache by ID
func (c *CacheImpl) VoiceChannel(voiceChannelID api.Snowflake) api.VoiceChannel {
	for _, guild := range c.voiceChannels {
		if channel, ok := guild[voiceChannelID]; ok {
			return channel
		}
	}
	return nil
}

// VoiceChannelsByName returns voice channels from cache by name
func (c *CacheImpl) VoiceChannelsByName(guildID api.Snowflake, name string, ignoreCase bool) []api.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		voiceChannels := make([]api.VoiceChannel, 1)
		for _, channel := range guildVoiceChannels {
			if ignoreCase && strings.ToLower(channel.Name()) == name || !ignoreCase && channel.Name() == name {
				voiceChannels = append(voiceChannels, channel)
			}
		}
		return voiceChannels
	}
	return nil
}

// VoiceChannels returns voice channels from a guild's cache
func (c *CacheImpl) VoiceChannels(guildID api.Snowflake) []api.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		voiceChannels := make([]api.VoiceChannel, len(guildVoiceChannels))
		i := 0
		for _, voiceChannel := range guildVoiceChannels {
			voiceChannels[i] = voiceChannel
			i++
		}
		return voiceChannels
	}
	return nil
}

// AllVoiceChannels returns all of the voice channels from cache as a slice
func (c *CacheImpl) AllVoiceChannels() []api.VoiceChannel {
	voiceChannels := make([]api.VoiceChannel, len(c.voiceChannels))
	for _, guildVoiceChannels := range c.voiceChannels {
		for _, voiceChannel := range guildVoiceChannels {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}

// VoiceChannelCache returns all of the voice channels from cache as a map
func (c *CacheImpl) VoiceChannelCache(guildID api.Snowflake) map[api.Snowflake]api.VoiceChannel {
	return c.voiceChannels[guildID]
}

// AllVoiceChannelCache returns all of the voice channels from cache as a map of maps
func (c *CacheImpl) AllVoiceChannelCache() map[api.Snowflake]map[api.Snowflake]api.VoiceChannel {
	return c.voiceChannels
}

// CacheVoiceChannel adds a voice channel to cache
func (c *CacheImpl) CacheVoiceChannel(voiceChannel api.VoiceChannel) api.VoiceChannel {
	if c.cacheFlags.Missing(api.CacheFlagVoiceChannels) {
		return voiceChannel
	}
	if guildVoiceChannels, ok := c.voiceChannels[voiceChannel.GuildID()]; ok {
		if guildVoiceChannel, ok := guildVoiceChannels[voiceChannel.ID()]; ok {
			*guildVoiceChannel.(*api.ChannelImpl) = *voiceChannel.(*api.ChannelImpl)
			return guildVoiceChannel
		}
		guildVoiceChannels[voiceChannel.ID()] = voiceChannel
	}
	return voiceChannel
}

// UncacheVoiceChannel removes a voice channel from cache
func (c *CacheImpl) UncacheVoiceChannel(guildID api.Snowflake, voiceChannelID api.Snowflake) {
	delete(c.voiceChannels[guildID], voiceChannelID)
}

// FindVoiceChannel returns a voice channel from cache by custom method
func (c *CacheImpl) FindVoiceChannel(guildID api.Snowflake, check func(u api.VoiceChannel) bool) api.VoiceChannel {
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			return voiceChannel
		}
	}
	return nil
}

// FindVoiceChannels returns voice channels from cache by custom method
func (c *CacheImpl) FindVoiceChannels(guildID api.Snowflake, check func(u api.VoiceChannel) bool) []api.VoiceChannel {
	voiceChannels := make([]api.VoiceChannel, 1)
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}

// Category returns a category from cache by ID
func (c *CacheImpl) Category(categoryID api.Snowflake) api.Category {
	for _, guildCategories := range c.categories {
		if channel, ok := guildCategories[categoryID]; ok {
			return channel
		}
	}
	return nil
}

// CategoriesByName returns categories from cache by name
func (c *CacheImpl) CategoriesByName(guildID api.Snowflake, name string, ignoreCase bool) []api.Category {
	if guildCategories, ok := c.categories[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		categories := make([]api.Category, 1)
		for _, channel := range guildCategories {
			if ignoreCase && strings.ToLower(channel.Name()) == name || !ignoreCase && channel.Name() == name {
				categories = append(categories, channel)
			}
		}
		return categories
	}
	return nil
}

// Categories returns the categories of a guild by ID
func (c *CacheImpl) Categories(guildID api.Snowflake) []api.Category {
	if guildCategories, ok := c.categories[guildID]; ok {
		categories := make([]api.Category, len(guildCategories))
		i := 0
		for _, category := range guildCategories {
			categories[i] = category
			i++
		}
		return categories
	}
	return nil
}

// AllCategories returns all categories from cache as a slice
func (c *CacheImpl) AllCategories() []api.Category {
	categories := make([]api.Category, len(c.categories))
	for _, guildCategories := range c.categories {
		for _, category := range guildCategories {
			categories = append(categories, category)
		}
	}
	return categories
}

// CategoryCache returns all categories from a guild's cache as a map
func (c *CacheImpl) CategoryCache(guildID api.Snowflake) map[api.Snowflake]api.Category {
	return c.categories[guildID]
}

// AllCategoryCache returns all categories as a map of maps
func (c *CacheImpl) AllCategoryCache() map[api.Snowflake]map[api.Snowflake]api.Category {
	return c.categories
}

// CacheCategory adds a category to the cache
func (c *CacheImpl) CacheCategory(category api.Category) api.Category {
	if c.cacheFlags.Missing(api.CacheFlagCategories) {
		return category
	}
	if guildCategories, ok := c.categories[category.GuildID()]; ok {
		if guildCategory, ok := guildCategories[category.ID()]; ok {
			// todo this is ugly af topi
			*guildCategory.(*api.ChannelImpl) = *category.(*api.ChannelImpl)
			return guildCategory
		}
		guildCategories[category.ID()] = category
	}
	return category
}

// UncacheCategory removes a category from cache
func (c *CacheImpl) UncacheCategory(guildID api.Snowflake, categoryID api.Snowflake) {
	delete(c.categories[guildID], categoryID)
}

// FindCategory finds a category in a guild by custom method
func (c *CacheImpl) FindCategory(guildID api.Snowflake, check func(u api.Category) bool) api.Category {
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			return category
		}
	}
	return nil
}

// FindCategories finds categories in a guild by custom method
func (c *CacheImpl) FindCategories(guildID api.Snowflake, check func(u api.Category) bool) []api.Category {
	categories := make([]api.Category, 1)
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			categories = append(categories, category)
		}
	}
	return categories
}

// Emote returns a specific emoji from the cache
func (c *CacheImpl) Emote(emoteID api.Snowflake) *api.Emoji {
	for _, guildEmojis := range c.emojis {
		if emoji, ok := guildEmojis[emoteID]; ok {
			return emoji
		}
	}
	return nil
}

// EmotesByName returns all emojis for a guild by name
func (c *CacheImpl) EmotesByName(guildID api.Snowflake, name string, ignoreCase bool) []*api.Emoji {
	if guildEmojis, ok := c.emojis[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		emojis := make([]*api.Emoji, 1)
		for _, emoji := range guildEmojis {
			if ignoreCase && strings.ToLower(emoji.Name) == name || !ignoreCase && emoji.Name == name {
				emojis = append(emojis, emoji)
			}
		}
		return emojis
	}
	return nil
}

// Emotes returns all cached emojis for a guild
func (c *CacheImpl) Emotes(guildID api.Snowflake) []*api.Emoji {
	if guildEmojis, ok := c.emojis[guildID]; ok {
		emojis := make([]*api.Emoji, len(guildEmojis))
		i := 0
		for _, emoji := range guildEmojis {
			emojis[i] = emoji
			i++
		}
		return emojis
	}
	return nil
}

// EmoteCache returns the emoji cache for a specific guild
func (c *CacheImpl) EmoteCache(guildID api.Snowflake) map[api.Snowflake]*api.Emoji {
	return c.emojis[guildID]
}

// AllEmoteCache returns the full emoji cache
func (c *CacheImpl) AllEmoteCache() map[api.Snowflake]map[api.Snowflake]*api.Emoji {
	return c.emojis
}

// CacheEmote adds an Emote to the api.Cache if emoji caches are used
func (c *CacheImpl) CacheEmote(emoji *api.Emoji) *api.Emoji {
	if c.cacheFlags.Missing(api.CacheFlagEmotes) {
		return emoji
	}
	if guildEmojis, ok := c.emojis[emoji.GuildID]; ok {
		if guildEmoji, ok := guildEmojis[emoji.ID]; ok {
			*guildEmoji = *emoji
			return guildEmoji
		}
		guildEmojis[emoji.ID] = emoji
	}
	return emoji
}

// UncacheEmote removes an Emote from api.Cache
func (c *CacheImpl) UncacheEmote(guildID api.Snowflake, emoteID api.Snowflake) {
	delete(c.emojis[guildID], emoteID)
}
