package cache

import (
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type SelfUserCache interface {
	SelfUser() (discord.OAuth2User, bool)
	SetSelfUser(selfUser discord.OAuth2User)
}

func NewSelfUserCache() SelfUserCache {
	return &selfUserCacheImpl{}
}

type selfUserCacheImpl struct {
	selfUserMu sync.Mutex
	selfUser   *discord.OAuth2User
}

func (c *selfUserCacheImpl) SelfUser() (discord.OAuth2User, bool) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	if c.selfUser == nil {
		return discord.OAuth2User{}, false
	}
	return *c.selfUser, true
}

func (c *selfUserCacheImpl) SetSelfUser(user discord.OAuth2User) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	c.selfUser = &user
}

type GuildCache interface {
	GuildCache() Cache[discord.Guild]

	IsGuildUnready(guildID snowflake.ID, opts ...RequestOpt) bool
	SetGuildUnready(guildID snowflake.ID, unready bool, opts ...RequestOpt)
	UnreadyGuildIDs(opts ...RequestOpt) []snowflake.ID

	IsGuildUnavailable(guildID snowflake.ID, opts ...RequestOpt) bool
	SetGuildUnavailable(guildID snowflake.ID, unavailable bool, opts ...RequestOpt)
	UnavailableGuildIDs(opts ...RequestOpt) []snowflake.ID

	Guild(guildID snowflake.ID, opts ...RequestOpt) (discord.Guild, bool)
	Guilds(opts ...RequestOpt) iter.Seq[discord.Guild]
	GuildsLen(opts ...RequestOpt) int
	AddGuild(guild discord.Guild, opts ...RequestOpt)
	RemoveGuild(guildID snowflake.ID, opts ...RequestOpt) (discord.Guild, bool)
}

func NewGuildCache(cache Cache[discord.Guild], unreadyGuilds Set[snowflake.ID], unavailableGuilds Set[snowflake.ID]) GuildCache {
	return &guildCacheImpl{
		cache:             cache,
		unreadyGuilds:     unreadyGuilds,
		unavailableGuilds: unavailableGuilds,
	}
}

type guildCacheImpl struct {
	cache             Cache[discord.Guild]
	unreadyGuilds     Set[snowflake.ID]
	unavailableGuilds Set[snowflake.ID]
}

func (c *guildCacheImpl) GuildCache() Cache[discord.Guild] {
	return c.cache
}

func (c *guildCacheImpl) IsGuildUnready(guildID snowflake.ID, opts ...RequestOpt) bool {
	return c.unreadyGuilds.Has(guildID, opts...)
}

func (c *guildCacheImpl) SetGuildUnready(guildID snowflake.ID, unready bool, opts ...RequestOpt) {
	if c.unreadyGuilds.Has(guildID) && !unready {
		c.unreadyGuilds.Remove(guildID, opts...)
	} else if !c.unreadyGuilds.Has(guildID) && unready {
		c.unreadyGuilds.Add(guildID, opts...)
	}
}

func (c *guildCacheImpl) UnreadyGuildIDs(opts ...RequestOpt) []snowflake.ID {
	var guilds []snowflake.ID
	for guildID := range c.unreadyGuilds.All(opts...) {
		guilds = append(guilds, guildID)
	}
	return guilds
}

func (c *guildCacheImpl) IsGuildUnavailable(guildID snowflake.ID, opts ...RequestOpt) bool {
	return c.unavailableGuilds.Has(guildID, opts...)
}

func (c *guildCacheImpl) SetGuildUnavailable(guildID snowflake.ID, unavailable bool, opts ...RequestOpt) {
	if c.unavailableGuilds.Has(guildID) && !unavailable {
		c.unavailableGuilds.Remove(guildID, opts...)
	} else if !c.unavailableGuilds.Has(guildID) && unavailable {
		c.unavailableGuilds.Add(guildID, opts...)
	}
}

func (c *guildCacheImpl) UnavailableGuildIDs(opts ...RequestOpt) []snowflake.ID {
	var guilds []snowflake.ID
	for guildId := range c.unavailableGuilds.All(opts...) {
		guilds = append(guilds, guildId)
	}
	return guilds
}

func (c *guildCacheImpl) Guild(guildID snowflake.ID, opts ...RequestOpt) (discord.Guild, bool) {
	return c.cache.Get(guildID, opts...)
}

func (c *guildCacheImpl) Guilds(opts ...RequestOpt) iter.Seq[discord.Guild] {
	return c.cache.All(opts...)
}

func (c *guildCacheImpl) GuildsLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *guildCacheImpl) AddGuild(guild discord.Guild, opts ...RequestOpt) {
	c.cache.Put(guild.ID, guild, opts...)
}

func (c *guildCacheImpl) RemoveGuild(guildID snowflake.ID, opts ...RequestOpt) (discord.Guild, bool) {
	return c.cache.Remove(guildID, opts...)
}

type ChannelCache interface {
	ChannelCache() Cache[discord.GuildChannel]

	Channel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildChannel, bool)
	Channels(opts ...RequestOpt) iter.Seq[discord.GuildChannel]
	ChannelsForGuild(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.GuildChannel]
	ChannelsLen(opts ...RequestOpt) int
	AddChannel(channel discord.GuildChannel, opts ...RequestOpt)
	RemoveChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildChannel, bool)
	RemoveChannelsByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewChannelCache(cache Cache[discord.GuildChannel]) ChannelCache {
	return &channelCacheImpl{
		cache: cache,
	}
}

type channelCacheImpl struct {
	cache Cache[discord.GuildChannel]
}

func (c *channelCacheImpl) ChannelCache() Cache[discord.GuildChannel] {
	return c.cache
}

func (c *channelCacheImpl) Channel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildChannel, bool) {
	return c.cache.Get(channelID, opts...)
}

func (c *channelCacheImpl) Channels(opts ...RequestOpt) iter.Seq[discord.GuildChannel] {
	return c.cache.All(opts...)
}

func (c *channelCacheImpl) ChannelsForGuild(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.GuildChannel] {
	return func(yield func(discord.GuildChannel) bool) {
		for channel := range c.Channels(opts...) {
			if channel.GuildID() == guildID {
				if !yield(channel) {
					return
				}
			}
		}
	}
}

func (c *channelCacheImpl) ChannelsLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *channelCacheImpl) AddChannel(channel discord.GuildChannel, opts ...RequestOpt) {
	c.cache.Put(channel.ID(), channel, opts...)
}

func (c *channelCacheImpl) RemoveChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildChannel, bool) {
	return c.cache.Remove(channelID, opts...)
}

func (c *channelCacheImpl) RemoveChannelsByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.RemoveIf(func(channel discord.GuildChannel) bool {
		return channel.GuildID() == guildID
	}, opts...)
}

type StageInstanceCache interface {
	StageInstanceCache() GroupedCache[discord.StageInstance]

	StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...RequestOpt) (discord.StageInstance, bool)
	StageInstances(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.StageInstance]
	StageInstancesAllLen(opts ...RequestOpt) int
	StageInstancesLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddStageInstance(stageInstance discord.StageInstance, opts ...RequestOpt)
	RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...RequestOpt) (discord.StageInstance, bool)
	RemoveStageInstancesByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewStageInstanceCache(cache GroupedCache[discord.StageInstance]) StageInstanceCache {
	return &stageInstanceCacheImpl{
		cache: cache,
	}
}

type stageInstanceCacheImpl struct {
	cache GroupedCache[discord.StageInstance]
}

func (c *stageInstanceCacheImpl) StageInstanceCache() GroupedCache[discord.StageInstance] {
	return c.cache
}

func (c *stageInstanceCacheImpl) StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...RequestOpt) (discord.StageInstance, bool) {
	return c.cache.Get(guildID, stageInstanceID, opts...)
}

func (c *stageInstanceCacheImpl) StageInstances(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.StageInstance] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *stageInstanceCacheImpl) StageInstancesAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *stageInstanceCacheImpl) StageInstancesLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *stageInstanceCacheImpl) AddStageInstance(stageInstance discord.StageInstance, opts ...RequestOpt) {
	c.cache.Put(stageInstance.GuildID, stageInstance.ID, stageInstance, opts...)
}

func (c *stageInstanceCacheImpl) RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...RequestOpt) (discord.StageInstance, bool) {
	return c.cache.Remove(guildID, stageInstanceID, opts...)
}

func (c *stageInstanceCacheImpl) RemoveStageInstancesByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type GuildScheduledEventCache interface {
	GuildScheduledEventCache() GroupedCache[discord.GuildScheduledEvent]

	GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) (discord.GuildScheduledEvent, bool)
	GuildScheduledEvents(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.GuildScheduledEvent]
	GuildScheduledEventsAllLen(opts ...RequestOpt) int
	GuildScheduledEventsLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, opts ...RequestOpt)
	RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) (discord.GuildScheduledEvent, bool)
	RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewGuildScheduledEventCache(cache GroupedCache[discord.GuildScheduledEvent]) GuildScheduledEventCache {
	return &guildScheduledEventCacheImpl{
		cache: cache,
	}
}

type guildScheduledEventCacheImpl struct {
	cache GroupedCache[discord.GuildScheduledEvent]
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventCache() GroupedCache[discord.GuildScheduledEvent] {
	return c.cache
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) (discord.GuildScheduledEvent, bool) {
	return c.cache.Get(guildID, guildScheduledEventID, opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvents(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.GuildScheduledEvent] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *guildScheduledEventCacheImpl) AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, opts ...RequestOpt) {
	c.cache.Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent, opts...)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...RequestOpt) (discord.GuildScheduledEvent, bool) {
	return c.cache.Remove(guildID, guildScheduledEventID, opts...)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type GuildSoundboardSoundCache interface {
	GuildSoundboardSoundCache() GroupedCache[discord.SoundboardSound]
	GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) (discord.SoundboardSound, bool)
	GuildSoundboardSounds(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.SoundboardSound]
	GuildSoundboardSoundsAllLen(opts ...RequestOpt) int
	GuildSoundboardSoundsLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddGuildSoundboardSound(sound discord.SoundboardSound, opts ...RequestOpt)
	RemoveGuildSoundboardSound(guildID snowflake.ID, sound snowflake.ID, opts ...RequestOpt) (discord.SoundboardSound, bool)
	RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewGuildSoundboardSoundCache(cache GroupedCache[discord.SoundboardSound]) GuildSoundboardSoundCache {
	return &guildSoundboardSoundCacheImpl{
		cache: cache,
	}
}

type guildSoundboardSoundCacheImpl struct {
	cache GroupedCache[discord.SoundboardSound]
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundCache() GroupedCache[discord.SoundboardSound] {
	return c.cache
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) (discord.SoundboardSound, bool) {
	return c.cache.Get(guildID, soundID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSounds(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.SoundboardSound] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) AddGuildSoundboardSound(sound discord.SoundboardSound, opts ...RequestOpt) {
	c.cache.Put(*sound.GuildID, sound.SoundID, sound, opts...)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...RequestOpt) (discord.SoundboardSound, bool) {
	return c.cache.Remove(guildID, soundID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type RoleCache interface {
	RoleCache() GroupedCache[discord.Role]

	Role(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (discord.Role, bool)
	Roles(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Role]
	RolesAllLen(opts ...RequestOpt) int
	RolesLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddRole(role discord.Role, opts ...RequestOpt)
	RemoveRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (discord.Role, bool)
	RemoveRolesByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewRoleCache(cache GroupedCache[discord.Role]) RoleCache {
	return &roleCacheImpl{
		cache: cache,
	}
}

type roleCacheImpl struct {
	cache GroupedCache[discord.Role]
}

func (c *roleCacheImpl) RoleCache() GroupedCache[discord.Role] {
	return c.cache
}

func (c *roleCacheImpl) Role(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (discord.Role, bool) {
	return c.cache.Get(guildID, roleID, opts...)
}

func (c *roleCacheImpl) Roles(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Role] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *roleCacheImpl) RolesAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *roleCacheImpl) RolesLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *roleCacheImpl) AddRole(role discord.Role, opts ...RequestOpt) {
	c.cache.Put(role.GuildID, role.ID, role, opts...)
}

func (c *roleCacheImpl) RemoveRole(guildID snowflake.ID, roleID snowflake.ID, opts ...RequestOpt) (discord.Role, bool) {
	return c.cache.Remove(guildID, roleID, opts...)
}

func (c *roleCacheImpl) RemoveRolesByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type MemberCache interface {
	MemberCache() GroupedCache[discord.Member]

	Member(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Member, bool)
	Members(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Member]
	MembersAllLen(opts ...RequestOpt) int
	MembersLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddMember(member discord.Member, opts ...RequestOpt)
	RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Member, bool)
	RemoveMembersByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewMemberCache(cache GroupedCache[discord.Member]) MemberCache {
	return &memberCacheImpl{
		cache: cache,
	}
}

type memberCacheImpl struct {
	cache GroupedCache[discord.Member]
}

func (c *memberCacheImpl) MemberCache() GroupedCache[discord.Member] {
	return c.cache
}

func (c *memberCacheImpl) Member(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Member, bool) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *memberCacheImpl) Members(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Member] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *memberCacheImpl) MembersAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *memberCacheImpl) MembersLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *memberCacheImpl) AddMember(member discord.Member, opts ...RequestOpt) {
	c.cache.Put(member.GuildID, member.User.ID, member, opts...)
}

func (c *memberCacheImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Member, bool) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *memberCacheImpl) RemoveMembersByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type ThreadMemberCache interface {
	ThreadMemberCache() GroupedCache[discord.ThreadMember]

	ThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.ThreadMember, bool)
	ThreadMembers(threadID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.ThreadMember]
	ThreadMembersAllLen(opts ...RequestOpt) int
	ThreadMembersLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddThreadMember(threadMember discord.ThreadMember, opts ...RequestOpt)
	RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.ThreadMember, bool)
	RemoveThreadMembersByThreadID(threadID snowflake.ID, opts ...RequestOpt)
}

func NewThreadMemberCache(cache GroupedCache[discord.ThreadMember]) ThreadMemberCache {
	return &threadMemberCacheImpl{
		cache: cache,
	}
}

type threadMemberCacheImpl struct {
	cache GroupedCache[discord.ThreadMember]
}

func (c *threadMemberCacheImpl) ThreadMemberCache() GroupedCache[discord.ThreadMember] {
	return c.cache
}

func (c *threadMemberCacheImpl) ThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.ThreadMember, bool) {
	return c.cache.Get(threadID, userID, opts...)
}

func (c *threadMemberCacheImpl) ThreadMembers(threadID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.ThreadMember] {
	return c.cache.GroupAll(threadID, opts...)
}

func (c *threadMemberCacheImpl) ThreadMembersAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *threadMemberCacheImpl) ThreadMembersLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *threadMemberCacheImpl) AddThreadMember(threadMember discord.ThreadMember, opts ...RequestOpt) {
	c.cache.Put(threadMember.ThreadID, threadMember.UserID, threadMember, opts...)
}

func (c *threadMemberCacheImpl) RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.ThreadMember, bool) {
	return c.cache.Remove(threadID, userID, opts...)
}

func (c *threadMemberCacheImpl) RemoveThreadMembersByThreadID(threadID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(threadID, opts...)
}

type PresenceCache interface {
	PresenceCache() GroupedCache[discord.Presence]

	Presence(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Presence, bool)
	Presences(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Presence]
	PresencesAllLen(opts ...RequestOpt) int
	PresencesLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddPresence(presence discord.Presence, opts ...RequestOpt)
	RemovePresence(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Presence, bool)
	RemovePresencesByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewPresenceCache(cache GroupedCache[discord.Presence]) PresenceCache {
	return &presenceCacheImpl{
		cache: cache,
	}
}

type presenceCacheImpl struct {
	cache GroupedCache[discord.Presence]
}

func (c *presenceCacheImpl) PresenceCache() GroupedCache[discord.Presence] {
	return c.cache
}

func (c *presenceCacheImpl) Presence(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Presence, bool) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *presenceCacheImpl) Presences(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Presence] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *presenceCacheImpl) PresencesAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *presenceCacheImpl) PresencesLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *presenceCacheImpl) AddPresence(presence discord.Presence, opts ...RequestOpt) {
	c.cache.Put(presence.GuildID, presence.PresenceUser.ID, presence, opts...)
}

func (c *presenceCacheImpl) RemovePresence(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.Presence, bool) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *presenceCacheImpl) RemovePresencesByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type VoiceStateCache interface {
	VoiceStateCache() GroupedCache[discord.VoiceState]

	VoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.VoiceState, bool)
	VoiceStates(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.VoiceState]
	VoiceStatesAllLen(opts ...RequestOpt) int
	VoiceStatesLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddVoiceState(voiceState discord.VoiceState, opts ...RequestOpt)
	RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.VoiceState, bool)
	RemoveVoiceStatesByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewVoiceStateCache(cache GroupedCache[discord.VoiceState]) VoiceStateCache {
	return &voiceStateCacheImpl{
		cache: cache,
	}
}

type voiceStateCacheImpl struct {
	cache GroupedCache[discord.VoiceState]
}

func (c *voiceStateCacheImpl) VoiceStateCache() GroupedCache[discord.VoiceState] {
	return c.cache
}

func (c *voiceStateCacheImpl) VoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.VoiceState, bool) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *voiceStateCacheImpl) VoiceStates(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.VoiceState] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *voiceStateCacheImpl) VoiceStatesAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *voiceStateCacheImpl) VoiceStatesLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *voiceStateCacheImpl) AddVoiceState(voiceState discord.VoiceState, opts ...RequestOpt) {
	c.cache.Put(voiceState.GuildID, voiceState.UserID, voiceState, opts...)
}

func (c *voiceStateCacheImpl) RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...RequestOpt) (discord.VoiceState, bool) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *voiceStateCacheImpl) RemoveVoiceStatesByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type MessageCache interface {
	MessageCache() GroupedCache[discord.Message]

	Message(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (discord.Message, bool)
	Messages(channelID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Message]
	MessagesAllLen(opts ...RequestOpt) int
	MessagesLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddMessage(message discord.Message, opts ...RequestOpt)
	RemoveMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (discord.Message, bool)
	RemoveMessagesByChannelID(channelID snowflake.ID, opts ...RequestOpt)
	RemoveMessagesByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewMessageCache(cache GroupedCache[discord.Message]) MessageCache {
	return &messageCacheImpl{
		cache: cache,
	}
}

type messageCacheImpl struct {
	cache GroupedCache[discord.Message]
}

func (c *messageCacheImpl) MessageCache() GroupedCache[discord.Message] {
	return c.cache
}

func (c *messageCacheImpl) Message(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (discord.Message, bool) {
	return c.cache.Get(channelID, messageID, opts...)
}

func (c *messageCacheImpl) Messages(channelID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Message] {
	return c.cache.GroupAll(channelID, opts...)
}

func (c *messageCacheImpl) MessagesAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *messageCacheImpl) MessagesLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *messageCacheImpl) AddMessage(message discord.Message, opts ...RequestOpt) {
	c.cache.Put(message.ChannelID, message.ID, message, opts...)
}

func (c *messageCacheImpl) RemoveMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...RequestOpt) (discord.Message, bool) {
	return c.cache.Remove(channelID, messageID, opts...)
}

func (c *messageCacheImpl) RemoveMessagesByChannelID(channelID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(channelID, opts...)
}

func (c *messageCacheImpl) RemoveMessagesByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.RemoveIf(func(_ snowflake.ID, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == guildID
	}, opts...)
}

type EmojiCache interface {
	EmojiCache() GroupedCache[discord.Emoji]

	Emoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (discord.Emoji, bool)
	Emojis(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Emoji]
	EmojisAllLen(opts ...RequestOpt) int
	EmojisLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddEmoji(emoji discord.Emoji, opts ...RequestOpt)
	RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (discord.Emoji, bool)
	RemoveEmojisByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewEmojiCache(cache GroupedCache[discord.Emoji]) EmojiCache {
	return &emojiCacheImpl{
		cache: cache,
	}
}

type emojiCacheImpl struct {
	cache GroupedCache[discord.Emoji]
}

func (c *emojiCacheImpl) EmojiCache() GroupedCache[discord.Emoji] {
	return c.cache
}

func (c *emojiCacheImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (discord.Emoji, bool) {
	return c.cache.Get(guildID, emojiID, opts...)
}

func (c *emojiCacheImpl) Emojis(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Emoji] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *emojiCacheImpl) EmojisAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *emojiCacheImpl) EmojisLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *emojiCacheImpl) AddEmoji(emoji discord.Emoji, opts ...RequestOpt) {
	c.cache.Put(emoji.GuildID, emoji.ID, emoji, opts...)
}

func (c *emojiCacheImpl) RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...RequestOpt) (discord.Emoji, bool) {
	return c.cache.Remove(guildID, emojiID, opts...)
}

func (c *emojiCacheImpl) RemoveEmojisByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

type StickerCache interface {
	StickerCache() GroupedCache[discord.Sticker]

	Sticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) (discord.Sticker, bool)
	Stickers(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Sticker]
	StickersAllLen(opts ...RequestOpt) int
	StickersLen(guildID snowflake.ID, opts ...RequestOpt) int
	AddSticker(sticker discord.Sticker, opts ...RequestOpt)
	RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) (discord.Sticker, bool)
	RemoveStickersByGuildID(guildID snowflake.ID, opts ...RequestOpt)
}

func NewStickerCache(cache GroupedCache[discord.Sticker]) StickerCache {
	return &stickerCacheImpl{
		cache: cache,
	}
}

type stickerCacheImpl struct {
	cache GroupedCache[discord.Sticker]
}

func (c *stickerCacheImpl) StickerCache() GroupedCache[discord.Sticker] {
	return c.cache
}

func (c *stickerCacheImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) (discord.Sticker, bool) {
	return c.cache.Get(guildID, stickerID, opts...)
}

func (c *stickerCacheImpl) Stickers(guildID snowflake.ID, opts ...RequestOpt) iter.Seq[discord.Sticker] {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *stickerCacheImpl) StickersAllLen(opts ...RequestOpt) int {
	return c.cache.Len(opts...)
}

func (c *stickerCacheImpl) StickersLen(guildID snowflake.ID, opts ...RequestOpt) int {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *stickerCacheImpl) AddSticker(sticker discord.Sticker, opts ...RequestOpt) {
	if sticker.GuildID == nil {
		return
	}
	c.cache.Put(*sticker.GuildID, sticker.ID, sticker, opts...)
}

func (c *stickerCacheImpl) RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) (discord.Sticker, bool) {
	return c.cache.Remove(guildID, stickerID, opts...)
}

func (c *stickerCacheImpl) RemoveStickersByGuildID(guildID snowflake.ID, opts ...RequestOpt) {
	c.cache.GroupRemove(guildID, opts...)
}

// Caches combines all different entity caches into one with some utility methods.
type Caches interface {
	SelfUserCache
	GuildCache
	ChannelCache
	StageInstanceCache
	GuildScheduledEventCache
	GuildSoundboardSoundCache
	RoleCache
	MemberCache
	ThreadMemberCache
	PresenceCache
	VoiceStateCache
	MessageCache
	EmojiCache
	StickerCache

	// CacheFlags returns the current configured FLags of the caches.
	CacheFlags() Flags

	// MemberPermissions returns the calculated permissions of the given member.
	// This requires the FlagRoles to be set.
	MemberPermissions(member discord.Member, opts ...RequestOpt) discord.Permissions

	// MemberPermissionsInChannel returns the calculated permissions of the given member in the given channel.
	// This requires the FlagRoles and FlagChannels to be set.
	MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member, opts ...RequestOpt) discord.Permissions

	// MemberRoles returns all roles of the given member.
	// This requires the FlagRoles to be set.
	MemberRoles(member discord.Member, opts ...RequestOpt) []discord.Role

	// AudioChannelMembers returns all members which are in the given audio channel.
	// This requires the FlagVoiceStates to be set.
	AudioChannelMembers(channel discord.GuildAudioChannel, opts ...RequestOpt) []discord.Member

	// SelfMember returns the current bot member from the given guildID.
	// This is only available after we received the gateway.EventTypeGuildCreate event for the given guildID.
	SelfMember(guildID snowflake.ID, opts ...RequestOpt) (discord.Member, bool)

	// GuildThreadsInChannel returns all discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildThreadsInChannel(channelID snowflake.ID, opts ...RequestOpt) []discord.GuildThread

	// GuildMessageChannel returns a discord.GuildMessageChannel from the ChannelCache and a bool indicating if it exists.
	GuildMessageChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildMessageChannel, bool)

	// GuildThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool)

	// GuildAudioChannel returns a discord.GetGuildAudioChannel from the ChannelCache and a bool indicating if it exists.
	GuildAudioChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildAudioChannel, bool)

	// GuildTextChannel returns a discord.GuildTextChannel from the ChannelCache and a bool indicating if it exists.
	GuildTextChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildTextChannel, bool)

	// GuildVoiceChannel returns a discord.GuildVoiceChannel from the ChannelCache and a bool indicating if it exists.
	GuildVoiceChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildVoiceChannel, bool)

	// GuildCategoryChannel returns a discord.GuildCategoryChannel from the ChannelCache and a bool indicating if it exists.
	GuildCategoryChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildCategoryChannel, bool)

	// GuildNewsChannel returns a discord.GuildNewsChannel from the ChannelCache and a bool indicating if it exists.
	GuildNewsChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildNewsChannel, bool)

	// GuildNewsThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildNewsThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool)

	// GuildPublicThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildPublicThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool)

	// GuildPrivateThread returns a discord.GuildThread from the ChannelCache and a bool indicating if it exists.
	GuildPrivateThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool)

	// GuildStageVoiceChannel returns a discord.GuildStageVoiceChannel from the ChannelCache and a bool indicating if it exists.
	GuildStageVoiceChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildStageVoiceChannel, bool)

	// GuildForumChannel returns a discord.GuildForumChannel from the ChannelCache and a bool indicating if it exists.
	GuildForumChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildForumChannel, bool)

	// GuildMediaChannel returns a discord.GuildMediaChannel from the ChannelCache and a bool indicating if it exists.
	GuildMediaChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildMediaChannel, bool)
}

// New returns a new default Caches instance with the given ConfigOpt(s) applied.
func New(opts ...ConfigOpt) Caches {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &cachesImpl{
		config:                    cfg,
		selfUserCache:             cfg.SelfUserCache,
		guildCache:                cfg.GuildCache,
		channelCache:              cfg.ChannelCache,
		stageInstanceCache:        cfg.StageInstanceCache,
		guildScheduledEventCache:  cfg.GuildScheduledEventCache,
		guildSoundboardSoundCache: cfg.GuildSoundboardSoundCache,
		roleCache:                 cfg.RoleCache,
		memberCache:               cfg.MemberCache,
		threadMemberCache:         cfg.ThreadMemberCache,
		presenceCache:             cfg.PresenceCache,
		voiceStateCache:           cfg.VoiceStateCache,
		messageCache:              cfg.MessageCache,
		emojiCache:                cfg.EmojiCache,
		stickerCache:              cfg.StickerCache,
	}
}

// these type aliases are needed to allow having the GuildCache, ChannelCache, etc. as methods on the cachesImpl struct
type (
	guildCache                = GuildCache
	channelCache              = ChannelCache
	stageInstanceCache        = StageInstanceCache
	guildScheduledEventCache  = GuildScheduledEventCache
	guildSoundboardSoundCache = GuildSoundboardSoundCache
	roleCache                 = RoleCache
	memberCache               = MemberCache
	threadMemberCache         = ThreadMemberCache
	presenceCache             = PresenceCache
	voiceStateCache           = VoiceStateCache
	messageCache              = MessageCache
	emojiCache                = EmojiCache
	stickerCache              = StickerCache
	selfUserCache             = SelfUserCache
)

type cachesImpl struct {
	config config

	guildCache
	channelCache
	stageInstanceCache
	guildScheduledEventCache
	guildSoundboardSoundCache
	roleCache
	memberCache
	threadMemberCache
	presenceCache
	voiceStateCache
	messageCache
	emojiCache
	stickerCache
	selfUserCache
}

func (c *cachesImpl) CacheFlags() Flags {
	return c.config.CacheFlags
}

func (c *cachesImpl) MemberPermissions(member discord.Member, opts ...RequestOpt) discord.Permissions {
	if guild, ok := c.Guild(member.GuildID, opts...); ok && guild.OwnerID == member.User.ID {
		return discord.PermissionsAll
	}

	var permissions discord.Permissions
	if publicRole, ok := c.Role(member.GuildID, member.GuildID, opts...); ok {
		permissions = publicRole.Permissions
	}

	for _, role := range c.MemberRoles(member, opts...) {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll
		}
	}
	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions
}

func (c *cachesImpl) MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member, opts ...RequestOpt) discord.Permissions {
	permissions := c.MemberPermissions(member, opts...)
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll
	}

	var (
		allow discord.Permissions
		deny  discord.Permissions
	)

	if overwrite, ok := channel.PermissionOverwrites().Role(channel.GuildID()); ok {
		permissions |= overwrite.Allow
		permissions &= ^overwrite.Deny
	}

	for _, roleID := range member.RoleIDs {
		if roleID == channel.GuildID() {
			continue
		}

		if overwrite, ok := channel.PermissionOverwrites().Role(roleID); ok {
			allow |= overwrite.Allow
			deny |= overwrite.Deny
		}
	}

	if overwrite, ok := channel.PermissionOverwrites().Member(member.User.ID); ok {
		allow |= overwrite.Allow
		deny |= overwrite.Deny
	}

	permissions &= ^deny
	permissions |= allow

	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}

	return permissions
}

func (c *cachesImpl) MemberRoles(member discord.Member, opts ...RequestOpt) []discord.Role {
	var roles []discord.Role

	for role := range c.Roles(member.GuildID, opts...) {
		if slices.Contains(member.RoleIDs, role.ID) {
			roles = append(roles, role)
		}
	}
	return roles
}

func (c *cachesImpl) AudioChannelMembers(channel discord.GuildAudioChannel, opts ...RequestOpt) []discord.Member {
	var members []discord.Member
	for state := range c.VoiceStates(channel.GuildID(), opts...) {
		if member, ok := c.Member(channel.GuildID(), state.UserID, opts...); ok && state.ChannelID != nil && *state.ChannelID == channel.ID() {
			members = append(members, member)
		}
	}
	return members
}

func (c *cachesImpl) SelfMember(guildID snowflake.ID, opts ...RequestOpt) (discord.Member, bool) {
	selfUser, ok := c.SelfUser()
	if !ok {
		return discord.Member{}, false
	}
	return c.Member(guildID, selfUser.ID, opts...)
}

func (c *cachesImpl) GuildThreadsInChannel(channelID snowflake.ID, opts ...RequestOpt) []discord.GuildThread {
	var threads []discord.GuildThread
	for channel := range c.Channels(opts...) {
		if thread, ok := channel.(discord.GuildThread); ok && *thread.ParentID() == channelID {
			threads = append(threads, thread)
		}
	}
	return threads
}

func (c *cachesImpl) MessageChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.MessageChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.MessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildMessageChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildMessageChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if chM, ok := ch.(discord.GuildMessageChannel); ok {
			return chM, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildThread); ok {
			return cCh, true
		}
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildAudioChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildAudioChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *cachesImpl) GuildTextChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildTextChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildTextChannel{}, false
}

func (c *cachesImpl) GuildVoiceChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildVoiceChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildVoiceChannel{}, false
}

func (c *cachesImpl) GuildCategoryChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildCategoryChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildCategoryChannel{}, false
}

func (c *cachesImpl) GuildNewsChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildNewsChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildNewsChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsChannel{}, false
}

func (c *cachesImpl) GuildNewsThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID, opts...); ok && ch.Type() == discord.ChannelTypeGuildNewsThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildPublicThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID, opts...); ok && ch.Type() == discord.ChannelTypeGuildPublicThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildPrivateThread(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildThread, bool) {
	if ch, ok := c.GuildThread(channelID, opts...); ok && ch.Type() == discord.ChannelTypeGuildPrivateThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *cachesImpl) GuildStageVoiceChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildStageVoiceChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildStageVoiceChannel{}, false
}

func (c *cachesImpl) GuildForumChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildForumChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildForumChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildForumChannel{}, false
}

func (c *cachesImpl) GuildMediaChannel(channelID snowflake.ID, opts ...RequestOpt) (discord.GuildMediaChannel, bool) {
	if ch, ok := c.Channel(channelID, opts...); ok {
		if cCh, ok := ch.(discord.GuildMediaChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildMediaChannel{}, false
}
