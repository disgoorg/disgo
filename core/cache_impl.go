package core

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

func NewCache(disgo Disgo, config CacheConfig) Cache {
	cache := &CacheImpl{
		disgo:  disgo,
		quit:   make(chan struct{}),
		config: config,
	}
	go cache.startCleanup(10 * time.Second)
	return cache
}

// CacheImpl is used for api.Disgo's api.Cache
type CacheImpl struct {
	disgo Disgo
	quit  chan struct{}

	config CacheConfig

	globalCommandCache GlobalCommandCache
	guildCommandCache  GuildCommandCache
	userCache          UserCache
	roleCache          RoleCache
	memberCache        MemberCache
	voiceStateCache    VoiceStateCache
	messageCache       MessageCache
	emojiCache         EmojiCache
	guildCache         GuildCache
	channelCache       ChannelCache
	textChannelCache   TextChannelCache
	voiceChannelCache  VoiceChannelCache
	dmChannelCache     DMChannelCache
	categoryCache      CategoryCache
	storeChannelCache  StoreChannelCache
	stageChannelCache  StageChannelCache
}

// Disgo returns the current api.Disgo instance
func (c *CacheImpl) Disgo() Disgo {
	return c.disgo
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

func (c *CacheImpl) DoCleanup() {
	c.globalCommandCache.DoCleanup()
	c.guildCommandCache.DoCleanup()
	c.userCache.DoCleanup()
	c.roleCache.DoCleanup()
	c.memberCache.DoCleanup()
	c.voiceStateCache.DoCleanup()
	c.messageCache.DoCleanup()
	c.emojiCache.DoCleanup()
	c.guildCache.DoCleanup()
	c.channelCache.DoCleanup()
	c.textChannelCache.DoCleanup()
	c.voiceChannelCache.DoCleanup()
	c.dmChannelCache.DoCleanup()
	c.categoryCache.DoCleanup()
	c.storeChannelCache.DoCleanup()
	c.stageChannelCache.DoCleanup()
}

// Close cleans up the cache and it's internal tasks
func (c *CacheImpl) Close() {
	c.Disgo().Logger().Info("closing cache goroutines...")
	c.quit <- struct{}{}
	c.Disgo().Logger().Info("closed cache goroutines")
}

// Config returns the current active api.CacheFlags
func (c CacheImpl) Config() CacheConfig {
	return c.config
}

var _ BaseCache = (*BaseCacheImpl)(nil)

type BaseCacheImpl struct {
	disgo Disgo
}

func (c *BaseCacheImpl) DoCleanup() {
	panic("implement me")
}

func (c *BaseCacheImpl) Disgo() Disgo {
	return c.disgo
}

var _ MessageCache = (*MessageCacheImpl)(nil)

type MessageCacheImpl struct {
	BaseCacheImpl
	messages    map[discord.Snowflake]map[discord.Snowflake]*Message
	cachePolicy MessageCachePolicy
}

func (c *MessageCacheImpl) DoCleanup() {
	for channelID, channelMessages := range c.messages {
		for _, message := range channelMessages {
			if !c.cachePolicy(message) {
				c.Uncache(channelID, message.ID)
			}
		}
	}
}

func (c *MessageCacheImpl) Get(channelID discord.Snowflake, messageID discord.Snowflake) *Message {
	return c.messages[channelID][messageID]
}

func (c *MessageCacheImpl) GetAll(channelID discord.Snowflake) []*Message {
	if _, ok := c.messages[channelID]; !ok {
		return nil
	}

	messages := make([]*Message, len(c.messages[channelID]))
	i := 0
	for _, message := range c.messages[channelID] {
		messages[i] = message
		i++
	}
	return messages
}

func (c *MessageCacheImpl) GetCache(channelID discord.Snowflake) map[discord.Snowflake]*Message {
	return c.messages[channelID]
}

func (c *MessageCacheImpl) GetAllCache() map[discord.Snowflake]map[discord.Snowflake]*Message {
	return c.messages
}

func (c *MessageCacheImpl) Cache(message *Message) *Message {
	if _, ok := c.messages[message.ChannelID]; !ok {
		c.messages[message.ChannelID] = map[discord.Snowflake]*Message{}
	}

	cache := c.messages[message.ChannelID]
	if _, ok := cache[message.ID]; ok {
		*cache[message.ID] = *message
	} else {
		cache[message.ID] = message
	}
	return cache[message.ID]
}

func (c *MessageCacheImpl) Uncache(channelID discord.Snowflake, messageID discord.Snowflake) {
	if _, ok := c.messages[channelID]; !ok {
		return
	}

	delete(c.messages[channelID], messageID)
}

var _ GuildCache = (*GuildCacheImpl)(nil)

type GuildCacheImpl struct {
	BaseCacheImpl
	disgo      Disgo
	cacheFlags CacheFlags
	guilds     map[discord.Snowflake]*Guild
}

func (c *GuildCacheImpl) Disgo() Disgo {
	return c.disgo
}

func (c *GuildCacheImpl) Guild(guildID discord.Snowflake) *Guild {
	return c.guilds[guildID]
}

func (c *GuildCacheImpl) GuildsByName(name string, ignoreCase bool) []*Guild {
	if ignoreCase {
		name = strings.ToLower(name)
	}

	var guilds []*Guild
	for _, guild := range c.guilds {
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}

func (c *GuildCacheImpl) Guilds() []*Guild {
	out := make([]*Guild, len(c.guilds))
	i := 0

	for _, guild := range c.guilds {
		out[i] = guild
		i++
	}

	return out
}

func (c *GuildCacheImpl) GuildCache() map[discord.Snowflake]*Guild {
	return c.guilds
}

func (c *GuildCacheImpl) CacheGuild(guild *Guild) *Guild {
	if c.cacheFlags.Missing(CacheFlagGuilds) {
		return guild
	}

	if _, ok := c.guilds[guild.ID]; ok {
		*c.guilds[guild.ID] = *guild
		return c.guilds[guild.ID]
	} else {
		c.guilds[guild.ID] = guild
		return guild
	}
}
func (c *GuildCacheImpl) UncacheGuild(guildID discord.Snowflake) {
	delete(c.guilds, guildID)
}

/*
globalCommands     map[discord.Snowflake]*entities.Command
guildCommands      map[discord.Snowflake]map[discord.Snowflake]*entities.Command
users              map[discord.Snowflake]*entities.User
guilds             map[discord.Snowflake]*entities.Guild
members            map[discord.Snowflake]map[discord.Snowflake]*entities.Member
voiceStates        map[discord.Snowflake]map[discord.Snowflake]*entities.VoiceState
roles              map[discord.Snowflake]map[discord.Snowflake]*entities.Role
dmChannels         map[discord.Snowflake]*entities.DMChannel
categories         map[discord.Snowflake]map[discord.Snowflake]*entities.Category
textChannels       map[discord.Snowflake]map[discord.Snowflake]*entities.TextChannel
voiceChannels      map[discord.Snowflake]map[discord.Snowflake]*entities.VoiceChannel
storeChannels      map[discord.Snowflake]map[discord.Snowflake]*entities.StoreChannel
emotes             map[discord.Snowflake]map[discord.Snowflake]*entities.Emoji
*/

/*

// Command returns an discord.Command from cache
func (c *CacheImpl) Command(commandID discord.Snowflake) *entities.Command {
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
func (c *CacheImpl) GuildCommandCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.Command {
	return c.guildCommands[guildID]
}

// AllGuildCommandCache returns the cache of all Guild Command(s)
func (c *CacheImpl) AllGuildCommandCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Command {
	return c.guildCommands
}

// GlobalCommandCache returns the cache of global Command(s)
func (c *CacheImpl) GlobalCommandCache() map[discord.Snowflake]*entities.Command {
	return c.globalCommands
}

// CacheGlobalCommand adds a global command to the cache
func (c *CacheImpl) CacheGlobalCommand(command *entities.Command) *entities.Command {
	if c.CacheFlags().Missing(CacheFlagCommands) {
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
func (c *CacheImpl) CacheGuildCommand(command *entities.Command) *entities.Command {
	if c.CacheFlags().Missing(CacheFlagCommands) {
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

// Uncache removes a global Command from the cache
func (c *CacheImpl) Uncache(commandID discord.Snowflake) {
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
func (c *CacheImpl) User(id discord.Snowflake) *entities.User {
	return c.users[id]
}

// UserByTag allows you to get a user from the cache by their Tag
func (c *CacheImpl) UserByTag(tag string) *entities.User {
	for _, user := range c.users {
		if user.Tag() == tag {
			return user
		}
	}
	return nil
}

// UsersByName allows you to get users from the cache by username
func (c *CacheImpl) UsersByName(name string, ignoreCase bool) []*entities.User {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	users := make([]*entities.User, 1)
	for _, user := range c.users {
		if ignoreCase && strings.ToLower(user.Username) == name || !ignoreCase && user.Username == name {
			users = append(users, user)
		}
	}
	return users
}

// Users returns all users from the cache as a slice
func (c *CacheImpl) Users() []*entities.User {
	users := make([]*entities.User, len(c.users))
	i := 0
	for _, user := range c.users {
		users[i] = user
		i++
	}
	return users
}

// UserCache returns all users from the cache a map
func (c *CacheImpl) UserCache() map[discord.Snowflake]*entities.User {
	return c.users
}

// Cache adds a user to the cache
func (c *CacheImpl) Cache(user *entities.User) *entities.User {
	// TODO: only cache user if we have a mutal guild & always cache self user
	if globalUser, ok := c.users[user.ID]; ok {
		*globalUser = *user
		return globalUser
	}
	c.users[user.ID] = user
	return user
}

// UncacheUser removes a user from the cache
func (c *CacheImpl) UncacheUser(id discord.Snowflake) {
	delete(c.users, id)
}

// FindUser finds a user from the cache with a custom method
func (c *CacheImpl) FindUser(check func(u *entities.User) bool) *entities.User {
	for _, user := range c.users {
		if check(user) {
			return user
		}
	}
	return nil
}

// FindUsers finds several users from the cache with a custom method
func (c *CacheImpl) FindUsers(check func(u *entities.User) bool) []*entities.User {
	users := make([]*entities.User, 1)
	for _, user := range c.users {
		if check(user) {
			users = append(users, user)
		}
	}
	return users
}

// Guild allows you to get a guild from the cache by ID
func (c *CacheImpl) Guild(guildID discord.Snowflake) *entities.Guild {
	return c.guilds[guildID]
}

// GuildsByName allows you to get guilds from the cache by name
func (c *CacheImpl) GuildsByName(name string, ignoreCase bool) []*entities.Guild {
	if ignoreCase {
		name = strings.ToLower(name)
	}
	guilds := make([]*entities.Guild, 1)
	for _, guild := range c.guilds {
		if ignoreCase && strings.ToLower(guild.Name) == name || !ignoreCase && guild.Name == name {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}

// Guilds returns the guild cache as a slice
func (c *CacheImpl) Guilds() []*entities.Guild {
	guilds := make([]*entities.Guild, len(c.guilds))
	i := 0
	for _, guild := range c.guilds {
		guilds[i] = guild
		i++
	}
	return guilds
}

// GuildCache returns the guild cache as a map
func (c *CacheImpl) GuildCache() map[discord.Snowflake]*entities.Guild {
	return c.guilds
}

// CacheGuild adds a guild to the cache
func (c *CacheImpl) CacheGuild(guild *entities.Guild) *entities.Guild {
	if _, ok := c.guilds[guild.ID]; ok {
		// update old guild
		*c.guilds[guild.ID] = *guild
		return c.guilds[guild.ID]
	}
	// guild was not yet cached so cache it directly
	c.guilds[guild.ID] = guild
	c.guildCommands[guild.ID] = map[discord.Snowflake]*entities.Command{}
	c.members[guild.ID] = map[discord.Snowflake]*entities.Member{}
	c.voiceStates[guild.ID] = map[discord.Snowflake]*entities.VoiceState{}
	c.roles[guild.ID] = map[discord.Snowflake]*entities.Role{}
	c.categories[guild.ID] = map[discord.Snowflake]*entities.Category{}
	c.textChannels[guild.ID] = map[discord.Snowflake]*entities.TextChannel{}
	c.voiceChannels[guild.ID] = map[discord.Snowflake]*entities.VoiceChannel{}
	c.storeChannels[guild.ID] = map[discord.Snowflake]*entities.StoreChannel{}
	return guild
}

//UncacheGuild removes a guild and all of it's children from the cache
func (c *CacheImpl) UncacheGuild(guildID discord.Snowflake) {
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
func (c *CacheImpl) FindGuild(check func(g *entities.Guild) bool) *entities.Guild {
	for _, guild := range c.guilds {
		if check(guild) {
			return guild
		}
	}
	return nil
}

//FindGuilds finds multiple guilds with a custom method
func (c *CacheImpl) FindGuilds(check func(g *entities.Guild) bool) []*entities.Guild {
	guilds := make([]*entities.Guild, 1)
	for _, guild := range c.guilds {
		if check(guild) {
			guilds = append(guilds, guild)
		}
	}
	return guilds
}

// Message returns a message from cache
func (c *CacheImpl) Message(channelID discord.Snowflake, messageID discord.Snowflake) *entities.Message {
	if channelMessages, ok := c.messages[channelID]; ok {
		return channelMessages[messageID]
	}
	return nil
}

// Messages returns the messages of a channel from cache
func (c *CacheImpl) Messages(channelID discord.Snowflake) []*entities.Message {
	if channelMessages, ok := c.messages[channelID]; ok {
		messages := make([]*entities.Message, len(channelMessages))
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
func (c *CacheImpl) MessageCache(channelID discord.Snowflake) map[discord.Snowflake]*entities.Message {
	return c.messages[channelID]
}

// AllMessageCache returns the entire cache of messages
func (c *CacheImpl) AllMessageCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Message {
	return c.messages
}

//CacheMessage adds a message to the cache
func (c *CacheImpl) CacheMessage(message *entities.Message) *entities.Message {
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
func (c *CacheImpl) UncacheMessage(channelID discord.Snowflake, messageID discord.Snowflake) {
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
func (c *CacheImpl) Member(guildID discord.Snowflake, userID discord.Snowflake) *entities.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		return guildMembers[userID]
	}
	return nil
}

// MemberByTag returns a member from cache by guild ID and user tag
func (c *CacheImpl) MemberByTag(guildID discord.Snowflake, tag string) *entities.Member {
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
func (c *CacheImpl) MembersByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		members := make([]*entities.Member, 1)
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
func (c *CacheImpl) Members(guildID discord.Snowflake) []*entities.Member {
	if guildMembers, ok := c.members[guildID]; ok {
		members := make([]*entities.Member, len(guildMembers))
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
func (c *CacheImpl) AllMembers() []*entities.Member {
	members := make([]*entities.Member, len(c.guilds))
	for _, guildMembers := range c.members {
		for _, member := range guildMembers {
			members = append(members, member)
		}
	}
	return members
}

//MemberCache returns the cache of a guild as a map
func (c *CacheImpl) MemberCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.Member {
	return c.members[guildID]
}

// AllMemberCache returns the entire cache as a map of maps
func (c *CacheImpl) AllMemberCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Member {
	return c.members
}

// Cache adds a member to the cache
func (c *CacheImpl) Cache(member *entities.Member) *entities.Member {
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

// Uncache removes a guild member from the cache
func (c *CacheImpl) Uncache(guildID discord.Snowflake, userID discord.Snowflake) {
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
func (c *CacheImpl) FindMember(guildID discord.Snowflake, check func(u *entities.Member) bool) *entities.Member {
	for _, member := range c.Members(guildID) {
		if check(member) {
			return member
		}
	}
	return nil
}

// FindMembers allows you to find api.Member(s) in a guild by custom method
func (c *CacheImpl) FindMembers(guildID discord.Snowflake, check func(u *entities.Member) bool) []*entities.Member {
	members := make([]*entities.Member, 1)
	for _, member := range c.Members(guildID) {
		if check(member) {
			members = append(members, member)
		}
	}
	return members
}

// VoiceState returns a Member's discord.VoiceState for an discord.Guild
func (c *CacheImpl) VoiceState(guildID discord.Snowflake, userID discord.Snowflake) *entities.VoiceState {
	if voiceStates, ok := c.voiceStates[guildID]; ok {
		return voiceStates[userID]
	}
	return nil
}

// VoiceStates returns the member cache of a guild by snowflake
func (c *CacheImpl) VoiceStates(guildID discord.Snowflake) []*entities.VoiceState {
	if guildVoiceStates, ok := c.voiceStates[guildID]; ok {
		voiceStates := make([]*entities.VoiceState, len(guildVoiceStates))
		i := 0
		for _, voiceState := range guildVoiceStates {
			voiceStates[i] = voiceState
			i++
		}
		return voiceStates
	}
	return nil
}

// VoiceStateCache returns the discord.VoiceState api.Cache of an discord.Guild as a map
func (c *CacheImpl) VoiceStateCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.VoiceState {
	return c.voiceStates[guildID]
}

// Cache adds an discord.VoiceState from the api.Cache
func (c *CacheImpl) Cache(voiceState *entities.VoiceState) *entities.VoiceState {
	// only cache voice states for ourself or member is cached & cache flag activated
	if c.cacheFlags.Missing(CacheFlagVoiceState) && voiceState.UserID != c.disgo.ApplicationID() {
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

// Uncache removes an discord.VoiceState from the api.Cache
func (c *CacheImpl) Uncache(guildID discord.Snowflake, userID discord.Snowflake) {
	delete(c.voiceStates[guildID], userID)
}

// Role returns a role from cache by guild ID and role ID
func (c *CacheImpl) Role(roleID discord.Snowflake) *entities.Role {
	for _, guildRoles := range c.roles {
		if role, ok := guildRoles[roleID]; ok {
			return role
		}
	}
	return nil
}

// RolesByName returns roles from cache by guild ID and name
func (c *CacheImpl) RolesByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.Role {
	if guildRoles, ok := c.roles[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		roles := make([]*entities.Role, 1)
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
func (c *CacheImpl) Roles(guildID discord.Snowflake) []*entities.Role {
	if guildRoles, ok := c.roles[guildID]; ok {
		roles := make([]*entities.Role, len(guildRoles))
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
func (c *CacheImpl) AllRoles() []*entities.Role {
	roles := make([]*entities.Role, len(c.guilds))
	for _, guildRoles := range c.roles {
		for _, role := range guildRoles {
			roles = append(roles, role)
		}
	}
	return roles
}

// RoleCache returns the role cache of a guild by ID
func (c *CacheImpl) RoleCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.Role {
	return c.roles[guildID]
}

// AllRoleCache returns the entire role cache
func (c *CacheImpl) AllRoleCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Role {
	return c.roles
}

// Cache adds a role to the cache
func (c *CacheImpl) Cache(role *entities.Role) *entities.Role {
	if c.cacheFlags.Missing(CacheFlagRoles) {
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

// Uncache removes a role from cache
func (c *CacheImpl) Uncache(guildID discord.Snowflake, roleID discord.Snowflake) {
	delete(c.roles[guildID], roleID)
}

// FindRole allows you to find a role in a guild by custom method
func (c *CacheImpl) FindRole(guildID discord.Snowflake, check func(u *entities.Role) bool) *entities.Role {
	for _, role := range c.Roles(guildID) {
		if check(role) {
			return role
		}
	}
	return nil
}

// FindRoles allows you to find roles in a guild by custom method
func (c *CacheImpl) FindRoles(guildID discord.Snowflake, check func(u *entities.Role) bool) []*entities.Role {
	roles := make([]*entities.Role, 1)
	for _, role := range c.Roles(guildID) {
		if check(role) {
			roles = append(roles, role)
		}
	}
	return roles
}

// Channel returns a channel from any channel cache by ID
func (c *CacheImpl) Channel(channelID discord.Snowflake) *entities.Channel {
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

// MessageChannel returns a channel from dm or text channel cache by ID
func (c *CacheImpl) MessageChannel(channelID discord.Snowflake) *entities.MessageChannel {
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

// GuildChannel returns a channel from a guild by ID
func (c *CacheImpl) GuildChannel(channelID discord.Snowflake) *entities.GuildChannel {
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

// DMChannel returns a DM channel by ID
func (c *CacheImpl) DMChannel(dmChannelID discord.Snowflake) *entities.DMChannel {
	return c.dmChannels[dmChannelID]
}

// DMChannels return all DM channels as a slice
func (c *CacheImpl) DMChannels() []*entities.DMChannel {
	channels := make([]*entities.DMChannel, len(c.dmChannels))
	i := 0
	for _, channel := range c.dmChannels {
		channels[i] = channel
		i++
	}
	return channels
}

// DMChannelCache returns the DM channels as a map
func (c *CacheImpl) DMChannelCache() map[discord.Snowflake]*entities.DMChannel {
	return c.dmChannels
}

// CacheDMChannel adds a DM channel to the cache
func (c *CacheImpl) CacheDMChannel(dmChannel *entities.DMChannel) *entities.DMChannel {
	if c.cacheFlags.Missing(CacheFlagDMChannels) {
		return dmChannel
	}
	if oldChannel, ok := c.dmChannels[dmChannel.ID]; ok {
		*oldChannel = *dmChannel
		return oldChannel
	}
	c.dmChannels[dmChannel.ID] = dmChannel
	return dmChannel
}

// UncacheDMChannel removes a DM channel from cache
func (c *CacheImpl) UncacheDMChannel(dmChannelID discord.Snowflake) {
	// TODO: check this
	// should be okay to just uncache all messages if the channel gets uncached as that should mean it got deleted
	if _, ok := c.messages[dmChannelID]; ok {
		delete(c.messages, dmChannelID)
	}
	delete(c.dmChannels, dmChannelID)
}

// FindDMChannel finds a DM channel in cache with a custom method
func (c *CacheImpl) FindDMChannel(check func(u *entities.DMChannel) bool) *entities.DMChannel {
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			return dmChannel
		}
	}
	return nil
}

// FindDMChannels finds DM Channels in cache with a custom method
func (c *CacheImpl) FindDMChannels(check func(u *entities.DMChannel) bool) []*entities.DMChannel {
	dmChannels := make([]*entities.DMChannel, 1)
	for _, dmChannel := range c.dmChannels {
		if check(dmChannel) {
			dmChannels = append(dmChannels, dmChannel)
		}
	}
	return dmChannels
}

// TextChannel returns a text channel from cache by ID
func (c *CacheImpl) TextChannel(textChannelID discord.Snowflake) *entities.TextChannel {
	for _, guild := range c.textChannels {
		if channel, ok := guild[textChannelID]; ok {
			return channel
		}
	}
	return nil
}

// TextChannelsByName returns text channel from cache by guild ID and name
func (c *CacheImpl) TextChannelsByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		textChannels := make([]*entities.TextChannel, 1)
		for _, channel := range guildTextChannels {
			if ignoreCase && strings.ToLower(*channel.MessageChannel.Name) == name || !ignoreCase && *channel.MessageChannel.Name == name {
				textChannels = append(textChannels, channel)
			}
		}
		return textChannels
	}
	return nil
}

// TextChannels returAllCategories() []Categoryns the text channel cache from a guild
func (c *CacheImpl) TextChannels(guildID discord.Snowflake) []*entities.TextChannel {
	if guildTextChannels, ok := c.textChannels[guildID]; ok {
		textChannels := make([]*entities.TextChannel, len(guildTextChannels))
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
func (c *CacheImpl) AllTextChannels() []*entities.TextChannel {
	textChannels := make([]*entities.TextChannel, len(c.textChannels))
	for _, guildTextChannels := range c.textChannels {
		for _, textChannel := range guildTextChannels {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}

// TextChannelCache returns the channel cache as a map
func (c *CacheImpl) TextChannelCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.TextChannel {
	return c.textChannels[guildID]
}

// AllTextChannelCache returns the text channel cache as a map of maps
func (c *CacheImpl) AllTextChannelCache() map[discord.Snowflake]map[discord.Snowflake]*entities.TextChannel {
	return c.textChannels
}

// CacheTextChannel adds a channel to the cache
func (c *CacheImpl) CacheTextChannel(textChannel *entities.TextChannel) *entities.TextChannel {
	if c.cacheFlags.Missing(CacheFlagTextChannels) {
		return textChannel
	}
	if guildTextChannels, ok := c.textChannels[*textChannel.GuildChannel.GuildID]; ok {
		if guildTextChannel, ok := guildTextChannels[textChannel.MessageChannel.ID]; ok {
			*guildTextChannel = *textChannel
			return guildTextChannel
		}
		guildTextChannels[textChannel.MessageChannel.ID] = textChannel
	}
	return textChannel
}

// UncacheTextChannel removes a text channel from the cache
func (c *CacheImpl) UncacheTextChannel(guildID discord.Snowflake, textChannelID discord.Snowflake) {
	// TODO: check this
	// should be okay to just uncache all messages if the channel gets uncached as that should mean it got deleted
	if _, ok := c.messages[textChannelID]; ok {
		delete(c.messages, textChannelID)
	}
	delete(c.textChannels[guildID], textChannelID)
}

// FindTextChannel finds a text channel in a guild by custom method
func (c *CacheImpl) FindTextChannel(guildID discord.Snowflake, check func(u *entities.TextChannel) bool) *entities.TextChannel {
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			return textChannel
		}
	}
	return nil
}

// FindTextChannels finds text channels in a guild by custom method
func (c *CacheImpl) FindTextChannels(guildID discord.Snowflake, check func(u *entities.TextChannel) bool) []*entities.TextChannel {
	textChannels := make([]*entities.TextChannel, 1)
	for _, textChannel := range c.TextChannelCache(guildID) {
		if check(textChannel) {
			textChannels = append(textChannels, textChannel)
		}
	}
	return textChannels
}

//StoreChannel returns a store channel from cache by ID
func (c *CacheImpl) StoreChannel(storeChannelID discord.Snowflake) *entities.StoreChannel {
	for _, guild := range c.storeChannels {
		if channel, ok := guild[storeChannelID]; ok {
			return channel
		}
	}
	return nil
}

//StoreChannelsByName returns store channels from cache by name
func (c *CacheImpl) StoreChannelsByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		storeChannels := make([]*entities.StoreChannel, 1)
		for _, channel := range guildStoreChannels {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				storeChannels = append(storeChannels, channel)
			}
		}
		return storeChannels
	}
	return nil
}

//StoreChannels returns store channels from cache by guild ID
func (c *CacheImpl) StoreChannels(guildID discord.Snowflake) []*entities.StoreChannel {
	if guildStoreChannels, ok := c.storeChannels[guildID]; ok {
		storeChannels := make([]*entities.StoreChannel, len(guildStoreChannels))
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
func (c *CacheImpl) AllStoreChannels() []*entities.StoreChannel {
	storeChannels := make([]*entities.StoreChannel, len(c.storeChannels))
	for _, guildStoreChannels := range c.storeChannels {
		for _, storeChannel := range guildStoreChannels {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}

//StoreChannelCache returns the store channels of a guild by ID
func (c *CacheImpl) StoreChannelCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.StoreChannel {
	return c.storeChannels[guildID]
}

//AllStoreChannelCache returns all store channels from cache as a map of maps
func (c *CacheImpl) AllStoreChannelCache() map[discord.Snowflake]map[discord.Snowflake]*entities.StoreChannel {
	return c.storeChannels
}

// CacheStoreChannel adds a store channel to the cache
func (c *CacheImpl) CacheStoreChannel(storeChannel *entities.StoreChannel) *entities.StoreChannel {
	if c.cacheFlags.Missing(CacheFlagStoreChannels) {
		return storeChannel
	}
	if guildStoreChannels, ok := c.storeChannels[*storeChannel.GuildID]; ok {
		if guildStoreChannel, ok := guildStoreChannels[storeChannel.ID]; ok {
			*guildStoreChannel = *storeChannel
			return guildStoreChannel
		}
		guildStoreChannels[storeChannel.ID] = storeChannel
	}
	return storeChannel
}

// UncacheStoreChannel removes a store channel from cache
func (c *CacheImpl) UncacheStoreChannel(guildID discord.Snowflake, storeChannelID discord.Snowflake) {
	delete(c.storeChannels[guildID], storeChannelID)
}

// FindStoreChannel returns a store channel from cache by custom method
func (c *CacheImpl) FindStoreChannel(guildID discord.Snowflake, check func(u *entities.StoreChannel) bool) *entities.StoreChannel {
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			return storeChannel
		}
	}
	return nil
}

// FindStoreChannels returns store channels from cache by custom method
func (c *CacheImpl) FindStoreChannels(guildID discord.Snowflake, check func(u *entities.StoreChannel) bool) []*entities.StoreChannel {
	storeChannels := make([]*entities.StoreChannel, 1)
	for _, storeChannel := range c.StoreChannelCache(guildID) {
		if check(storeChannel) {
			storeChannels = append(storeChannels, storeChannel)
		}
	}
	return storeChannels
}

// VoiceChannel returns a voice channel from cache by ID
func (c *CacheImpl) VoiceChannel(voiceChannelID discord.Snowflake) *entities.VoiceChannel {
	for _, guild := range c.voiceChannels {
		if channel, ok := guild[voiceChannelID]; ok {
			return channel
		}
	}
	return nil
}

// VoiceChannelsByName returns voice channels from cache by name
func (c *CacheImpl) VoiceChannelsByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		voiceChannels := make([]*entities.VoiceChannel, 1)
		for _, channel := range guildVoiceChannels {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				voiceChannels = append(voiceChannels, channel)
			}
		}
		return voiceChannels
	}
	return nil
}

// VoiceChannels returns voice channels from a guild's cache
func (c *CacheImpl) VoiceChannels(guildID discord.Snowflake) []*entities.VoiceChannel {
	if guildVoiceChannels, ok := c.voiceChannels[guildID]; ok {
		voiceChannels := make([]*entities.VoiceChannel, len(guildVoiceChannels))
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
func (c *CacheImpl) AllVoiceChannels() []*entities.VoiceChannel {
	voiceChannels := make([]*entities.VoiceChannel, len(c.voiceChannels))
	for _, guildVoiceChannels := range c.voiceChannels {
		for _, voiceChannel := range guildVoiceChannels {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}

// VoiceChannelCache returns all of the voice channels from cache as a map
func (c *CacheImpl) VoiceChannelCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.VoiceChannel {
	return c.voiceChannels[guildID]
}

// AllVoiceChannelCache returns all of the voice channels from cache as a map of maps
func (c *CacheImpl) AllVoiceChannelCache() map[discord.Snowflake]map[discord.Snowflake]*entities.VoiceChannel {
	return c.voiceChannels
}

// CacheVoiceChannel adds a voice channel to cache
func (c *CacheImpl) CacheVoiceChannel(voiceChannel *entities.VoiceChannel) *entities.VoiceChannel {
	if c.cacheFlags.Missing(CacheFlagVoiceChannels) {
		return voiceChannel
	}
	if guildVoiceChannels, ok := c.voiceChannels[*voiceChannel.GuildID]; ok {
		if guildVoiceChannel, ok := guildVoiceChannels[voiceChannel.ID]; ok {
			*guildVoiceChannel = *voiceChannel
			return guildVoiceChannel
		}
		guildVoiceChannels[voiceChannel.ID] = voiceChannel
	}
	return voiceChannel
}

// UncacheVoiceChannel removes a voice channel from cache
func (c *CacheImpl) UncacheVoiceChannel(guildID discord.Snowflake, voiceChannelID discord.Snowflake) {
	delete(c.voiceChannels[guildID], voiceChannelID)
}

// FindVoiceChannel returns a voice channel from cache by custom method
func (c *CacheImpl) FindVoiceChannel(guildID discord.Snowflake, check func(u *entities.VoiceChannel) bool) *entities.VoiceChannel {
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			return voiceChannel
		}
	}
	return nil
}

// FindVoiceChannels returns voice channels from cache by custom method
func (c *CacheImpl) FindVoiceChannels(guildID discord.Snowflake, check func(u *entities.VoiceChannel) bool) []*entities.VoiceChannel {
	voiceChannels := make([]*entities.VoiceChannel, 1)
	for _, voiceChannel := range c.VoiceChannelCache(guildID) {
		if check(voiceChannel) {
			voiceChannels = append(voiceChannels, voiceChannel)
		}
	}
	return voiceChannels
}

// Category returns a category from cache by ID
func (c *CacheImpl) Category(categoryID discord.Snowflake) *entities.Category {
	for _, guildCategories := range c.categories {
		if channel, ok := guildCategories[categoryID]; ok {
			return channel
		}
	}
	return nil
}

// CategoriesByName returns categories from cache by name
func (c *CacheImpl) CategoriesByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.Category {
	if guildCategories, ok := c.categories[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		categories := make([]*entities.Category, 1)
		for _, channel := range guildCategories {
			if ignoreCase && strings.ToLower(*channel.Name) == name || !ignoreCase && *channel.Name == name {
				categories = append(categories, channel)
			}
		}
		return categories
	}
	return nil
}

// Categories returns the categories of a guild by ID
func (c *CacheImpl) Categories(guildID discord.Snowflake) []*entities.Category {
	if guildCategories, ok := c.categories[guildID]; ok {
		categories := make([]*entities.Category, len(guildCategories))
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
func (c *CacheImpl) AllCategories() []*entities.Category {
	categories := make([]*entities.Category, len(c.categories))
	for _, guildCategories := range c.categories {
		for _, category := range guildCategories {
			categories = append(categories, category)
		}
	}
	return categories
}

// CategoryCache returns all categories from a guild's cache as a map
func (c *CacheImpl) CategoryCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.Category {
	return c.categories[guildID]
}

// AllCategoryCache returns all categories as a map of maps
func (c *CacheImpl) AllCategoryCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Category {
	return c.categories
}

// CacheCategory adds a category to the cache
func (c *CacheImpl) CacheCategory(category *entities.Category) *entities.Category {
	if c.cacheFlags.Missing(CacheFlagCategories) {
		return category
	}
	if guildCategories, ok := c.categories[*category.GuildID]; ok {
		if guildCategory, ok := guildCategories[category.ID]; ok {
			*guildCategory = *category
			return guildCategory
		}
		guildCategories[category.ID] = category
	}
	return category
}

// UncacheCategory removes a category from cache
func (c *CacheImpl) UncacheCategory(guildID discord.Snowflake, categoryID discord.Snowflake) {
	delete(c.categories[guildID], categoryID)
}

// FindCategory finds a category in a guild by custom method
func (c *CacheImpl) FindCategory(guildID discord.Snowflake, check func(u *entities.Category) bool) *entities.Category {
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			return category
		}
	}
	return nil
}

// FindCategories finds categories in a guild by custom method
func (c *CacheImpl) FindCategories(guildID discord.Snowflake, check func(u *entities.Category) bool) []*entities.Category {
	categories := make([]*entities.Category, 1)
	for _, category := range c.CategoryCache(guildID) {
		if check(category) {
			categories = append(categories, category)
		}
	}
	return categories
}

// Emote returns a specific emote from the cache
func (c *CacheImpl) Emote(emoteID discord.Snowflake) *entities.Emoji {
	for _, guildEmotes := range c.emotes {
		if emote, ok := guildEmotes[emoteID]; ok {
			return emote
		}
	}
	return nil
}

// EmotesByName returns all emotes for a guild by name
func (c *CacheImpl) EmotesByName(guildID discord.Snowflake, name string, ignoreCase bool) []*entities.Emoji {
	if guildEmotes, ok := c.emotes[guildID]; ok {
		if ignoreCase {
			name = strings.ToLower(name)
		}
		emotes := make([]*entities.Emoji, 1)
		for _, emote := range guildEmotes {
			if ignoreCase && strings.ToLower(emote.Name) == name || !ignoreCase && emote.Name == name {
				emotes = append(emotes, emote)
			}
		}
		return emotes
	}
	return nil
}

// Emotes returns all cached emotes for a guild
func (c *CacheImpl) Emotes(guildID discord.Snowflake) []*entities.Emoji {
	if guildEmotes, ok := c.emotes[guildID]; ok {
		emotes := make([]*entities.Emoji, len(guildEmotes))
		i := 0
		for _, emote := range guildEmotes {
			emotes[i] = emote
			i++
		}
		return emotes
	}
	return nil
}

// EmojiCache returns the emote cache for a specific guild
func (c *CacheImpl) EmoteCache(guildID discord.Snowflake) map[discord.Snowflake]*entities.Emoji {
	return c.emotes[guildID]
}

// AllEmoteCache returns the full emote cache
func (c *CacheImpl) AllEmoteCache() map[discord.Snowflake]map[discord.Snowflake]*entities.Emoji {
	return c.emotes
}

// CacheEmote adds an Emote to the api.Cache if emoji caches are used
func (c *CacheImpl) CacheEmote(emote *entities.Emoji) *entities.Emoji {
	if c.cacheFlags.Missing(CacheFlagEmotes) {
		return emote
	}
	if guildEmotes, ok := c.emotes[emote.GuildID]; ok {
		if guildEmote, ok := guildEmotes[emote.ID]; ok {
			*guildEmote = *emote
			return guildEmote
		}
		guildEmotes[emote.ID] = emote
	}
	return emote
}

// UncacheEmote removes an Emote from api.Cache
func (c *CacheImpl) UncacheEmote(guildID discord.Snowflake, emoteID discord.Snowflake) {
	delete(c.emotes[guildID], emoteID)
}
*/
