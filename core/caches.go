package core

type Caches interface {
	Config() CacheConfig

	Users() Cache[User]
	Roles() GroupedCache[Role]
	Members() GroupedCache[Member]
	ThreadMembers() GroupedCache[ThreadMember]
	Presences() GroupedCache[Presence]
	VoiceStates() GroupedCache[VoiceState]
	Messages() GroupedCache[Message]
	Emojis() GroupedCache[Emoji]
	Stickers() GroupedCache[Sticker]
	Guilds() Cache[Guild]
	Channels() Cache[Channel]
	StageInstances() Cache[StageInstance]
	GuildScheduledEvents() Cache[GuildScheduledEvent]
}

func NewCaches(config CacheConfig) Caches {
	return &cachesImpl{
		config: config,

		userCache:                NewCache[User](config.CacheFlags, CacheFlagUsers),
		roleCache:                NewGroupedCache[Role](config.CacheFlags, CacheFlagRoles),
		memberCache:              NewGroupedCacheWithPolicy[Member](config.MemberCachePolicy),
		threadMemberCache:        NewGroupedCache[ThreadMember](config.CacheFlags, CacheFlagThreadMembers),
		presenceCache:            NewGroupedCache[Presence](config.CacheFlags, CacheFlagPresences),
		voiceStateCache:          NewGroupedCache[VoiceState](config.CacheFlags, CacheFlagVoiceStates),
		messageCache:             NewGroupedCacheWithPolicy[Message](config.MessageCachePolicy),
		emojiCache:               NewGroupedCache[Emoji](config.CacheFlags, CacheFlagEmojis),
		stickerCache:             NewGroupedCache[Sticker](config.CacheFlags, CacheFlagStickers),
		guildCache:               NewCache[Guild](config.CacheFlags, CacheFlagGuilds),
		channelCache:             NewCache[Channel](config.CacheFlags, CacheFlagsAllChannels),
		stageInstanceCache:       NewCache[StageInstance](config.CacheFlags, CacheFlagStageInstances),
		guildScheduledEventCache: NewCache[GuildScheduledEvent](config.CacheFlags, CacheFlagGuildScheduledEvents),
	}
}

type cachesImpl struct {
	config CacheConfig

	userCache                Cache[User]
	roleCache                GroupedCache[Role]
	memberCache              GroupedCache[Member]
	threadMemberCache        GroupedCache[ThreadMember]
	presenceCache            GroupedCache[Presence]
	voiceStateCache          GroupedCache[VoiceState]
	messageCache             GroupedCache[Message]
	emojiCache               GroupedCache[Emoji]
	stickerCache             GroupedCache[Sticker]
	guildCache               Cache[Guild]
	channelCache             Cache[Channel]
	stageInstanceCache       Cache[StageInstance]
	guildScheduledEventCache Cache[GuildScheduledEvent]
}

func (c *cachesImpl) Config() CacheConfig {
	return c.config
}

func (c *cachesImpl) Users() Cache[User] {
	return c.userCache
}

func (c *cachesImpl) Roles() GroupedCache[Role] {
	return c.roleCache
}

func (c *cachesImpl) Members() GroupedCache[Member] {
	return c.memberCache
}

func (c *cachesImpl) ThreadMembers() GroupedCache[ThreadMember] {
	return c.threadMemberCache
}

func (c *cachesImpl) Presences() GroupedCache[Presence] {
	return c.presenceCache
}

func (c *cachesImpl) VoiceStates() GroupedCache[VoiceState] {
	return c.voiceStateCache
}

func (c *cachesImpl) Messages() GroupedCache[Message] {
	return c.messageCache
}

func (c *cachesImpl) Emojis() GroupedCache[Emoji] {
	return c.emojiCache
}

func (c *cachesImpl) Stickers() GroupedCache[Sticker] {
	return c.stickerCache
}

func (c *cachesImpl) Guilds() Cache[Guild] {
	return c.guildCache
}

func (c *cachesImpl) Channels() Cache[Channel] {
	return c.channelCache
}

func (c *cachesImpl) StageInstances() Cache[StageInstance] {
	return c.stageInstanceCache
}

func (c *cachesImpl) GuildScheduledEvents() Cache[GuildScheduledEvent] {
	return c.guildScheduledEventCache
}
