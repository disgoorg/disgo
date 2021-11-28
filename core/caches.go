package core

type Caches interface {
	Config() CacheConfig

	Users() UserCache
	Roles() RoleCache
	Members() MemberCache
	ThreadMembers() ThreadMemberCache
	Presences() PresenceCache
	VoiceStates() VoiceStateCache
	Messages() MessageCache
	Emojis() EmojiCache
	Stickers() StickerCache
	Guilds() GuildCache
	Channels() ChannelCache
	StageInstances() StageInstanceCache
}

func NewCaches(config CacheConfig) Caches {
	return &cachesImpl{
		config: config,

		userCache:          NewUserCache(config.CacheFlags),
		roleCache:          NewRoleCache(config.CacheFlags),
		memberCache:        NewMemberCache(config.MemberCachePolicy),
		threadMemberCache:  NewThreadMemberCache(config.CacheFlags),
		presenceCache:      NewPresenceCache(config.CacheFlags),
		voiceStateCache:    NewVoiceStateCache(config.CacheFlags),
		messageCache:       NewMessageCache(config.MessageCachePolicy),
		emojiCache:         NewEmojiCache(config.CacheFlags),
		stickerCache:       NewStickerCache(config.CacheFlags),
		guildCache:         NewGuildCache(config.CacheFlags),
		channelCache:       NewChannelCache(config.CacheFlags),
		stageInstanceCache: NewStageInstanceCache(config.CacheFlags),
	}
}

type cachesImpl struct {
	config CacheConfig

	userCache          UserCache
	roleCache          RoleCache
	memberCache        MemberCache
	threadMemberCache  ThreadMemberCache
	presenceCache      PresenceCache
	voiceStateCache    VoiceStateCache
	messageCache       MessageCache
	emojiCache         EmojiCache
	stickerCache       StickerCache
	guildCache         GuildCache
	channelCache       ChannelCache
	stageInstanceCache StageInstanceCache
}

func (c *cachesImpl) Config() CacheConfig {
	return c.config
}

func (c *cachesImpl) Users() UserCache {
	return c.userCache
}

func (c *cachesImpl) Roles() RoleCache {
	return c.roleCache
}

func (c *cachesImpl) Members() MemberCache {
	return c.memberCache
}

func (c *cachesImpl) ThreadMembers() ThreadMemberCache {
	return c.threadMemberCache
}

func (c *cachesImpl) Presences() PresenceCache {
	return c.presenceCache
}

func (c *cachesImpl) VoiceStates() VoiceStateCache {
	return c.voiceStateCache
}

func (c *cachesImpl) Messages() MessageCache {
	return c.messageCache
}

func (c *cachesImpl) Emojis() EmojiCache {
	return c.emojiCache
}

func (c *cachesImpl) Stickers() StickerCache {
	return c.stickerCache
}

func (c *cachesImpl) Guilds() GuildCache {
	return c.guildCache
}

func (c *cachesImpl) Channels() ChannelCache {
	return c.channelCache
}

func (c *cachesImpl) StageInstances() StageInstanceCache {
	return c.stageInstanceCache
}
