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
	SelfUser() (discord.OAuth2User, error)
	SetSelfUser(selfUser discord.OAuth2User) error
}

func NewSelfUserCache() SelfUserCache {
	return &selfUserCacheImpl{}
}

type selfUserCacheImpl struct {
	selfUserMu sync.Mutex
	selfUser   *discord.OAuth2User
}

func (c *selfUserCacheImpl) SelfUser() (discord.OAuth2User, error) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	if c.selfUser == nil {
		return discord.OAuth2User{}, ErrNotFound
	}
	return *c.selfUser, nil
}

func (c *selfUserCacheImpl) SetSelfUser(user discord.OAuth2User) error {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	c.selfUser = &user
	return nil
}

type GuildCache interface {
	GuildCache() Cache[discord.Guild]

	IsGuildUnready(guildID snowflake.ID, opts ...AccessOpt) bool
	SetGuildUnready(guildID snowflake.ID, unready bool, opts ...AccessOpt)
	UnreadyGuildIDs(opts ...AccessOpt) []snowflake.ID

	IsGuildUnavailable(guildID snowflake.ID, opts ...AccessOpt) bool
	SetGuildUnavailable(guildID snowflake.ID, unavailable bool, opts ...AccessOpt)
	UnavailableGuildIDs(opts ...AccessOpt) []snowflake.ID

	Guild(guildID snowflake.ID, opts ...AccessOpt) (discord.Guild, error)
	Guilds(opts ...AccessOpt) (iter.Seq[discord.Guild], error)
	GuildsLen(opts ...AccessOpt) (int, error)
	AddGuild(guild discord.Guild, opts ...AccessOpt) error
	RemoveGuild(guildID snowflake.ID, opts ...AccessOpt) (discord.Guild, error)
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

func (c *guildCacheImpl) IsGuildUnready(guildID snowflake.ID, opts ...AccessOpt) bool {
	return c.unreadyGuilds.Has(guildID, opts...)
}

func (c *guildCacheImpl) SetGuildUnready(guildID snowflake.ID, unready bool, opts ...AccessOpt) {
	if c.unreadyGuilds.Has(guildID) && !unready {
		c.unreadyGuilds.Remove(guildID, opts...)
	} else if !c.unreadyGuilds.Has(guildID) && unready {
		c.unreadyGuilds.Add(guildID, opts...)
	}
}

func (c *guildCacheImpl) UnreadyGuildIDs(opts ...AccessOpt) []snowflake.ID {
	var guilds []snowflake.ID
	for guildID := range c.unreadyGuilds.All(opts...) {
		guilds = append(guilds, guildID)
	}
	return guilds
}

func (c *guildCacheImpl) IsGuildUnavailable(guildID snowflake.ID, opts ...AccessOpt) bool {
	return c.unavailableGuilds.Has(guildID, opts...)
}

func (c *guildCacheImpl) SetGuildUnavailable(guildID snowflake.ID, unavailable bool, opts ...AccessOpt) {
	if c.unavailableGuilds.Has(guildID) && !unavailable {
		c.unavailableGuilds.Remove(guildID, opts...)
	} else if !c.unavailableGuilds.Has(guildID) && unavailable {
		c.unavailableGuilds.Add(guildID, opts...)
	}
}

func (c *guildCacheImpl) UnavailableGuildIDs(opts ...AccessOpt) []snowflake.ID {
	var guilds []snowflake.ID
	for guildId := range c.unavailableGuilds.All(opts...) {
		guilds = append(guilds, guildId)
	}
	return guilds
}

func (c *guildCacheImpl) Guild(guildID snowflake.ID, opts ...AccessOpt) (discord.Guild, error) {
	return c.cache.Get(guildID, opts...)
}

func (c *guildCacheImpl) Guilds(opts ...AccessOpt) (iter.Seq[discord.Guild], error) {
	return c.cache.All(opts...)
}

func (c *guildCacheImpl) GuildsLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *guildCacheImpl) AddGuild(guild discord.Guild, opts ...AccessOpt) error {
	return c.cache.Put(guild.ID, guild, opts...)
}

func (c *guildCacheImpl) RemoveGuild(guildID snowflake.ID, opts ...AccessOpt) (discord.Guild, error) {
	return c.cache.Remove(guildID, opts...)
}

type ChannelCache interface {
	ChannelCache() Cache[discord.GuildChannel]

	Channel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildChannel, error)
	Channels(opts ...AccessOpt) (iter.Seq[discord.GuildChannel], error)
	ChannelsForGuild(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.GuildChannel], error)
	ChannelsLen(opts ...AccessOpt) (int, error)
	AddChannel(channel discord.GuildChannel, opts ...AccessOpt) error
	RemoveChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildChannel, error)
	RemoveChannelsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *channelCacheImpl) Channel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildChannel, error) {
	return c.cache.Get(channelID, opts...)
}

func (c *channelCacheImpl) Channels(opts ...AccessOpt) (iter.Seq[discord.GuildChannel], error) {
	return c.cache.All(opts...)
}

func (c *channelCacheImpl) ChannelsForGuild(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.GuildChannel], error) {
	seq, err := c.Channels(opts...)
	if err != nil {
		return nil, err
	}
	return func(yield func(discord.GuildChannel) bool) {
		for channel := range seq {
			if channel.GuildID() == guildID {
				if !yield(channel) {
					return
				}
			}
		}
	}, nil
}

func (c *channelCacheImpl) ChannelsLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *channelCacheImpl) AddChannel(channel discord.GuildChannel, opts ...AccessOpt) error {
	return c.cache.Put(channel.ID(), channel, opts...)
}

func (c *channelCacheImpl) RemoveChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildChannel, error) {
	return c.cache.Remove(channelID, opts...)
}

func (c *channelCacheImpl) RemoveChannelsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.RemoveIf(func(channel discord.GuildChannel) bool {
		return channel.GuildID() == guildID
	}, opts...)
}

type StageInstanceCache interface {
	StageInstanceCache() GroupedCache[discord.StageInstance]

	StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...AccessOpt) (discord.StageInstance, error)
	StageInstances(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.StageInstance], error)
	StageInstancesAllLen(opts ...AccessOpt) (int, error)
	StageInstancesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddStageInstance(stageInstance discord.StageInstance, opts ...AccessOpt) error
	RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...AccessOpt) (discord.StageInstance, error)
	RemoveStageInstancesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *stageInstanceCacheImpl) StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...AccessOpt) (discord.StageInstance, error) {
	return c.cache.Get(guildID, stageInstanceID, opts...)
}

func (c *stageInstanceCacheImpl) StageInstances(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.StageInstance], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *stageInstanceCacheImpl) StageInstancesAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *stageInstanceCacheImpl) StageInstancesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *stageInstanceCacheImpl) AddStageInstance(stageInstance discord.StageInstance, opts ...AccessOpt) error {
	return c.cache.Put(stageInstance.GuildID, stageInstance.ID, stageInstance, opts...)
}

func (c *stageInstanceCacheImpl) RemoveStageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID, opts ...AccessOpt) (discord.StageInstance, error) {
	return c.cache.Remove(guildID, stageInstanceID, opts...)
}

func (c *stageInstanceCacheImpl) RemoveStageInstancesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type GuildScheduledEventCache interface {
	GuildScheduledEventCache() GroupedCache[discord.GuildScheduledEvent]

	GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...AccessOpt) (discord.GuildScheduledEvent, error)
	GuildScheduledEvents(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.GuildScheduledEvent], error)
	GuildScheduledEventsAllLen(opts ...AccessOpt) (int, error)
	GuildScheduledEventsLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, opts ...AccessOpt) error
	RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...AccessOpt) (discord.GuildScheduledEvent, error)
	RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *guildScheduledEventCacheImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...AccessOpt) (discord.GuildScheduledEvent, error) {
	return c.cache.Get(guildID, guildScheduledEventID, opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEvents(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.GuildScheduledEvent], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *guildScheduledEventCacheImpl) GuildScheduledEventsLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *guildScheduledEventCacheImpl) AddGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, opts ...AccessOpt) error {
	return c.cache.Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent, opts...)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID, opts ...AccessOpt) (discord.GuildScheduledEvent, error) {
	return c.cache.Remove(guildID, guildScheduledEventID, opts...)
}

func (c *guildScheduledEventCacheImpl) RemoveGuildScheduledEventsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type GuildSoundboardSoundCache interface {
	GuildSoundboardSoundCache() GroupedCache[discord.SoundboardSound]
	GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...AccessOpt) (discord.SoundboardSound, error)
	GuildSoundboardSounds(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.SoundboardSound], error)
	GuildSoundboardSoundsAllLen(opts ...AccessOpt) (int, error)
	GuildSoundboardSoundsLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddGuildSoundboardSound(sound discord.SoundboardSound, opts ...AccessOpt) error
	RemoveGuildSoundboardSound(guildID snowflake.ID, sound snowflake.ID, opts ...AccessOpt) (discord.SoundboardSound, error)
	RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...AccessOpt) (discord.SoundboardSound, error) {
	return c.cache.Get(guildID, soundID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSounds(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.SoundboardSound], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *guildSoundboardSoundCacheImpl) GuildSoundboardSoundsLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) AddGuildSoundboardSound(sound discord.SoundboardSound, opts ...AccessOpt) error {
	return c.cache.Put(*sound.GuildID, sound.SoundID, sound, opts...)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSound(guildID snowflake.ID, soundID snowflake.ID, opts ...AccessOpt) (discord.SoundboardSound, error) {
	return c.cache.Remove(guildID, soundID, opts...)
}

func (c *guildSoundboardSoundCacheImpl) RemoveGuildSoundboardSoundsByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type RoleCache interface {
	RoleCache() GroupedCache[discord.Role]

	Role(guildID snowflake.ID, roleID snowflake.ID, opts ...AccessOpt) (discord.Role, error)
	Roles(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Role], error)
	RolesAllLen(opts ...AccessOpt) (int, error)
	RolesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddRole(role discord.Role, opts ...AccessOpt) error
	RemoveRole(guildID snowflake.ID, roleID snowflake.ID, opts ...AccessOpt) (discord.Role, error)
	RemoveRolesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *roleCacheImpl) Role(guildID snowflake.ID, roleID snowflake.ID, opts ...AccessOpt) (discord.Role, error) {
	return c.cache.Get(guildID, roleID, opts...)
}

func (c *roleCacheImpl) Roles(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Role], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *roleCacheImpl) RolesAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *roleCacheImpl) RolesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *roleCacheImpl) AddRole(role discord.Role, opts ...AccessOpt) error {
	return c.cache.Put(role.GuildID, role.ID, role, opts...)
}

func (c *roleCacheImpl) RemoveRole(guildID snowflake.ID, roleID snowflake.ID, opts ...AccessOpt) (discord.Role, error) {
	return c.cache.Remove(guildID, roleID, opts...)
}

func (c *roleCacheImpl) RemoveRolesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type MemberCache interface {
	MemberCache() GroupedCache[discord.Member]

	Member(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Member, error)
	Members(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Member], error)
	MembersAllLen(opts ...AccessOpt) (int, error)
	MembersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddMember(member discord.Member, opts ...AccessOpt) error
	RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Member, error)
	RemoveMembersByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *memberCacheImpl) Member(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Member, error) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *memberCacheImpl) Members(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Member], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *memberCacheImpl) MembersAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *memberCacheImpl) MembersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *memberCacheImpl) AddMember(member discord.Member, opts ...AccessOpt) error {
	return c.cache.Put(member.GuildID, member.User.ID, member, opts...)
}

func (c *memberCacheImpl) RemoveMember(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Member, error) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *memberCacheImpl) RemoveMembersByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type ThreadMemberCache interface {
	ThreadMemberCache() GroupedCache[discord.ThreadMember]

	ThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.ThreadMember, error)
	ThreadMembers(threadID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.ThreadMember], error)
	ThreadMembersAllLen(opts ...AccessOpt) (int, error)
	ThreadMembersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddThreadMember(threadMember discord.ThreadMember, opts ...AccessOpt) error
	RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.ThreadMember, error)
	RemoveThreadMembersByThreadID(threadID snowflake.ID, opts ...AccessOpt) error
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

func (c *threadMemberCacheImpl) ThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.ThreadMember, error) {
	return c.cache.Get(threadID, userID, opts...)
}

func (c *threadMemberCacheImpl) ThreadMembers(threadID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.ThreadMember], error) {
	return c.cache.GroupAll(threadID, opts...)
}

func (c *threadMemberCacheImpl) ThreadMembersAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *threadMemberCacheImpl) ThreadMembersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *threadMemberCacheImpl) AddThreadMember(threadMember discord.ThreadMember, opts ...AccessOpt) error {
	return c.cache.Put(threadMember.ThreadID, threadMember.UserID, threadMember, opts...)
}

func (c *threadMemberCacheImpl) RemoveThreadMember(threadID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.ThreadMember, error) {
	return c.cache.Remove(threadID, userID, opts...)
}

func (c *threadMemberCacheImpl) RemoveThreadMembersByThreadID(threadID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(threadID, opts...)
}

type PresenceCache interface {
	PresenceCache() GroupedCache[discord.Presence]

	Presence(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Presence, error)
	Presences(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Presence], error)
	PresencesAllLen(opts ...AccessOpt) (int, error)
	PresencesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddPresence(presence discord.Presence, opts ...AccessOpt) error
	RemovePresence(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Presence, error)
	RemovePresencesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *presenceCacheImpl) Presence(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Presence, error) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *presenceCacheImpl) Presences(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Presence], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *presenceCacheImpl) PresencesAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *presenceCacheImpl) PresencesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *presenceCacheImpl) AddPresence(presence discord.Presence, opts ...AccessOpt) error {
	return c.cache.Put(presence.GuildID, presence.PresenceUser.ID, presence, opts...)
}

func (c *presenceCacheImpl) RemovePresence(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.Presence, error) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *presenceCacheImpl) RemovePresencesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type VoiceStateCache interface {
	VoiceStateCache() GroupedCache[discord.VoiceState]

	VoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.VoiceState, error)
	VoiceStates(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.VoiceState], error)
	VoiceStatesAllLen(opts ...AccessOpt) (int, error)
	VoiceStatesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddVoiceState(voiceState discord.VoiceState, opts ...AccessOpt) error
	RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.VoiceState, error)
	RemoveVoiceStatesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *voiceStateCacheImpl) VoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.VoiceState, error) {
	return c.cache.Get(guildID, userID, opts...)
}

func (c *voiceStateCacheImpl) VoiceStates(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.VoiceState], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *voiceStateCacheImpl) VoiceStatesAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *voiceStateCacheImpl) VoiceStatesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *voiceStateCacheImpl) AddVoiceState(voiceState discord.VoiceState, opts ...AccessOpt) error {
	return c.cache.Put(voiceState.GuildID, voiceState.UserID, voiceState, opts...)
}

func (c *voiceStateCacheImpl) RemoveVoiceState(guildID snowflake.ID, userID snowflake.ID, opts ...AccessOpt) (discord.VoiceState, error) {
	return c.cache.Remove(guildID, userID, opts...)
}

func (c *voiceStateCacheImpl) RemoveVoiceStatesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type MessageCache interface {
	MessageCache() GroupedCache[discord.Message]

	Message(channelID snowflake.ID, messageID snowflake.ID, opts ...AccessOpt) (discord.Message, error)
	Messages(channelID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Message], error)
	MessagesAllLen(opts ...AccessOpt) (int, error)
	MessagesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddMessage(message discord.Message, opts ...AccessOpt) error
	RemoveMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...AccessOpt) (discord.Message, error)
	RemoveMessagesByChannelID(channelID snowflake.ID, opts ...AccessOpt) error
	RemoveMessagesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *messageCacheImpl) Message(channelID snowflake.ID, messageID snowflake.ID, opts ...AccessOpt) (discord.Message, error) {
	return c.cache.Get(channelID, messageID, opts...)
}

func (c *messageCacheImpl) Messages(channelID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Message], error) {
	return c.cache.GroupAll(channelID, opts...)
}

func (c *messageCacheImpl) MessagesAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *messageCacheImpl) MessagesLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *messageCacheImpl) AddMessage(message discord.Message, opts ...AccessOpt) error {
	return c.cache.Put(message.ChannelID, message.ID, message, opts...)
}

func (c *messageCacheImpl) RemoveMessage(channelID snowflake.ID, messageID snowflake.ID, opts ...AccessOpt) (discord.Message, error) {
	return c.cache.Remove(channelID, messageID, opts...)
}

func (c *messageCacheImpl) RemoveMessagesByChannelID(channelID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(channelID, opts...)
}

func (c *messageCacheImpl) RemoveMessagesByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.RemoveIf(func(_ snowflake.ID, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == guildID
	}, opts...)
}

type EmojiCache interface {
	EmojiCache() GroupedCache[discord.Emoji]

	Emoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...AccessOpt) (discord.Emoji, error)
	Emojis(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Emoji], error)
	EmojisAllLen(opts ...AccessOpt) (int, error)
	EmojisLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddEmoji(emoji discord.Emoji, opts ...AccessOpt) error
	RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...AccessOpt) (discord.Emoji, error)
	RemoveEmojisByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *emojiCacheImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...AccessOpt) (discord.Emoji, error) {
	return c.cache.Get(guildID, emojiID, opts...)
}

func (c *emojiCacheImpl) Emojis(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Emoji], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *emojiCacheImpl) EmojisAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *emojiCacheImpl) EmojisLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *emojiCacheImpl) AddEmoji(emoji discord.Emoji, opts ...AccessOpt) error {
	return c.cache.Put(emoji.GuildID, emoji.ID, emoji, opts...)
}

func (c *emojiCacheImpl) RemoveEmoji(guildID snowflake.ID, emojiID snowflake.ID, opts ...AccessOpt) (discord.Emoji, error) {
	return c.cache.Remove(guildID, emojiID, opts...)
}

func (c *emojiCacheImpl) RemoveEmojisByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
}

type StickerCache interface {
	StickerCache() GroupedCache[discord.Sticker]

	Sticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...AccessOpt) (discord.Sticker, error)
	Stickers(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Sticker], error)
	StickersAllLen(opts ...AccessOpt) (int, error)
	StickersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error)
	AddSticker(sticker discord.Sticker, opts ...AccessOpt) error
	RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...AccessOpt) (discord.Sticker, error)
	RemoveStickersByGuildID(guildID snowflake.ID, opts ...AccessOpt) error
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

func (c *stickerCacheImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...AccessOpt) (discord.Sticker, error) {
	return c.cache.Get(guildID, stickerID, opts...)
}

func (c *stickerCacheImpl) Stickers(guildID snowflake.ID, opts ...AccessOpt) (iter.Seq[discord.Sticker], error) {
	return c.cache.GroupAll(guildID, opts...)
}

func (c *stickerCacheImpl) StickersAllLen(opts ...AccessOpt) (int, error) {
	return c.cache.Len(opts...)
}

func (c *stickerCacheImpl) StickersLen(guildID snowflake.ID, opts ...AccessOpt) (int, error) {
	return c.cache.GroupLen(guildID, opts...)
}

func (c *stickerCacheImpl) AddSticker(sticker discord.Sticker, opts ...AccessOpt) error {
	if sticker.GuildID == nil {
		return nil
	}
	return c.cache.Put(*sticker.GuildID, sticker.ID, sticker, opts...)
}

func (c *stickerCacheImpl) RemoveSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...AccessOpt) (discord.Sticker, error) {
	return c.cache.Remove(guildID, stickerID, opts...)
}

func (c *stickerCacheImpl) RemoveStickersByGuildID(guildID snowflake.ID, opts ...AccessOpt) error {
	return c.cache.GroupRemove(guildID, opts...)
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
	MemberPermissions(member discord.Member, opts ...AccessOpt) (discord.Permissions, error)

	// MemberPermissionsInChannel returns the calculated permissions of the given member in the given channel.
	// This requires the FlagRoles and FlagChannels to be set.
	MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member, opts ...AccessOpt) (discord.Permissions, error)

	// MemberRoles returns all roles of the given member.
	// This requires the FlagRoles to be set.
	MemberRoles(member discord.Member, opts ...AccessOpt) ([]discord.Role, error)

	// AudioChannelMembers returns all members which are in the given audio channel.
	// This requires the FlagVoiceStates to be set.
	AudioChannelMembers(channel discord.GuildAudioChannel, opts ...AccessOpt) ([]discord.Member, error)

	// SelfMember returns the current bot member from the given guildID.
	// This is only available after we received the gateway.EventTypeGuildCreate event for the given guildID.
	SelfMember(guildID snowflake.ID, opts ...AccessOpt) (discord.Member, error)

	// GuildThreadsInChannel returns all discord.GuildThread from the ChannelCache.
	GuildThreadsInChannel(channelID snowflake.ID, opts ...AccessOpt) ([]discord.GuildThread, error)

	// GuildMessageChannel returns a discord.GuildMessageChannel from the ChannelCache.
	GuildMessageChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildMessageChannel, error)

	// GuildThread returns a discord.GuildThread from the ChannelCache.
	GuildThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error)

	// GuildAudioChannel returns a discord.GetGuildAudioChannel from the ChannelCache.
	GuildAudioChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildAudioChannel, error)

	// GuildTextChannel returns a discord.GuildTextChannel from the ChannelCache.
	GuildTextChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildTextChannel, error)

	// GuildVoiceChannel returns a discord.GuildVoiceChannel from the ChannelCache.
	GuildVoiceChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildVoiceChannel, error)

	// GuildCategoryChannel returns a discord.GuildCategoryChannel from the ChannelCache.
	GuildCategoryChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildCategoryChannel, error)

	// GuildNewsChannel returns a discord.GuildNewsChannel from the ChannelCache.
	GuildNewsChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildNewsChannel, error)

	// GuildNewsThread returns a discord.GuildThread from the ChannelCache.
	GuildNewsThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error)

	// GuildPublicThread returns a discord.GuildThread from the ChannelCache.
	GuildPublicThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error)

	// GuildPrivateThread returns a discord.GuildThread from the ChannelCache.
	GuildPrivateThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error)

	// GuildStageVoiceChannel returns a discord.GuildStageVoiceChannel from the ChannelCache.
	GuildStageVoiceChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildStageVoiceChannel, error)

	// GuildForumChannel returns a discord.GuildForumChannel from the ChannelCache.
	GuildForumChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildForumChannel, error)

	// GuildMediaChannel returns a discord.GuildMediaChannel from the ChannelCache.
	GuildMediaChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildMediaChannel, error)
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

func (c *cachesImpl) MemberPermissions(member discord.Member, opts ...AccessOpt) (discord.Permissions, error) {
	guild, err := c.Guild(member.GuildID, opts...)
	if err == nil && guild.OwnerID == member.User.ID {
		return discord.PermissionsAll, nil
	}

	var permissions discord.Permissions
	publicRole, err := c.Role(member.GuildID, member.GuildID, opts...)
	if err == nil {
		permissions = publicRole.Permissions
	}

	roles, err := c.MemberRoles(member, opts...)
	if err != nil {
		return permissions, err
	}
	for _, role := range roles {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll, nil
		}
	}
	if member.CommunicationDisabledUntil != nil && member.CommunicationDisabledUntil.After(time.Now()) {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions, nil
}

func (c *cachesImpl) MemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member, opts ...AccessOpt) (discord.Permissions, error) {
	permissions, err := c.MemberPermissions(member, opts...)
	if err != nil {
		return permissions, err
	}
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll, nil
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

	return permissions, nil
}

func (c *cachesImpl) MemberRoles(member discord.Member, opts ...AccessOpt) ([]discord.Role, error) {
	var roles []discord.Role

	rolesSeq, err := c.Roles(member.GuildID, opts...)
	if err != nil {
		return nil, err
	}
	for role := range rolesSeq {
		if slices.Contains(member.RoleIDs, role.ID) {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

func (c *cachesImpl) AudioChannelMembers(channel discord.GuildAudioChannel, opts ...AccessOpt) ([]discord.Member, error) {
	var members []discord.Member
	voiceStatesSeq, err := c.VoiceStates(channel.GuildID(), opts...)
	if err != nil {
		return nil, err
	}
	for state := range voiceStatesSeq {
		if state.ChannelID != nil && *state.ChannelID == channel.ID() {
			member, err := c.Member(channel.GuildID(), state.UserID, opts...)
			if err == nil {
				members = append(members, member)
			}
		}
	}
	return members, nil
}

func (c *cachesImpl) SelfMember(guildID snowflake.ID, opts ...AccessOpt) (discord.Member, error) {
	selfUser, err := c.SelfUser()
	if err != nil {
		return discord.Member{}, err
	}
	return c.Member(guildID, selfUser.ID, opts...)
}

func (c *cachesImpl) GuildThreadsInChannel(channelID snowflake.ID, opts ...AccessOpt) ([]discord.GuildThread, error) {
	var threads []discord.GuildThread
	channelsSeq, err := c.Channels(opts...)
	if err != nil {
		return nil, err
	}
	for channel := range channelsSeq {
		if thread, ok := channel.(discord.GuildThread); ok && *thread.ParentID() == channelID {
			threads = append(threads, thread)
		}
	}
	return threads, nil
}

func (c *cachesImpl) MessageChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.MessageChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return nil, err
	}
	if cCh, ok := ch.(discord.MessageChannel); ok {
		return cCh, nil
	}
	return nil, ErrNotFound
}

func (c *cachesImpl) GuildMessageChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildMessageChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return nil, err
	}
	if chM, ok := ch.(discord.GuildMessageChannel); ok {
		return chM, nil
	}
	return nil, ErrNotFound
}

func (c *cachesImpl) GuildThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildThread{}, err
	}
	if cCh, ok := ch.(discord.GuildThread); ok {
		return cCh, nil
	}
	return discord.GuildThread{}, ErrNotFound
}

func (c *cachesImpl) GuildAudioChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildAudioChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return nil, err
	}
	if cCh, ok := ch.(discord.GuildAudioChannel); ok {
		return cCh, nil
	}
	return nil, ErrNotFound
}

func (c *cachesImpl) GuildTextChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildTextChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildTextChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildTextChannel); ok {
		return cCh, nil
	}
	return discord.GuildTextChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildVoiceChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildVoiceChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildVoiceChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
		return cCh, nil
	}
	return discord.GuildVoiceChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildCategoryChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildCategoryChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildCategoryChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
		return cCh, nil
	}
	return discord.GuildCategoryChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildNewsChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildNewsChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildNewsChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildNewsChannel); ok {
		return cCh, nil
	}
	return discord.GuildNewsChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildNewsThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error) {
	ch, err := c.GuildThread(channelID, opts...)
	if err != nil {
		return discord.GuildThread{}, err
	}
	if ch.Type() == discord.ChannelTypeGuildNewsThread {
		return ch, nil
	}
	return discord.GuildThread{}, ErrNotFound
}

func (c *cachesImpl) GuildPublicThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error) {
	ch, err := c.GuildThread(channelID, opts...)
	if err != nil {
		return discord.GuildThread{}, err
	}
	if ch.Type() == discord.ChannelTypeGuildPublicThread {
		return ch, nil
	}
	return discord.GuildThread{}, ErrNotFound
}

func (c *cachesImpl) GuildPrivateThread(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildThread, error) {
	ch, err := c.GuildThread(channelID, opts...)
	if err != nil {
		return discord.GuildThread{}, err
	}
	if ch.Type() == discord.ChannelTypeGuildPrivateThread {
		return ch, nil
	}
	return discord.GuildThread{}, ErrNotFound
}

func (c *cachesImpl) GuildStageVoiceChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildStageVoiceChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildStageVoiceChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
		return cCh, nil
	}
	return discord.GuildStageVoiceChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildForumChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildForumChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildForumChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildForumChannel); ok {
		return cCh, nil
	}
	return discord.GuildForumChannel{}, ErrNotFound
}

func (c *cachesImpl) GuildMediaChannel(channelID snowflake.ID, opts ...AccessOpt) (discord.GuildMediaChannel, error) {
	ch, err := c.Channel(channelID, opts...)
	if err != nil {
		return discord.GuildMediaChannel{}, err
	}
	if cCh, ok := ch.(discord.GuildMediaChannel); ok {
		return cCh, nil
	}
	return discord.GuildMediaChannel{}, ErrNotFound
}
