package core

type Caches interface {
	Config() CacheConfig

	UserCache() UserCache
	RoleCache() RoleCache
	MemberCache() MemberCache
	VoiceStateCache() VoiceStateCache
	MessageCache() MessageCache
	EmojiCache() EmojiCache
	GuildCache() GuildCache
	ChannelCache() ChannelCache
	StageInstanceCache() StageInstanceCache
}

func NewCaches(config CacheConfig) Caches {
	return &cachesImpl{
		config: config,

		userCache:          NewUserCache(config.CacheFlags),
		roleCache:          NewRoleCache(config.CacheFlags),
		memberCache:        NewMemberCache(config.MemberCachePolicy),
		voiceStateCache:    NewVoiceStateCache(config.CacheFlags),
		messageCache:       NewMessageCache(config.MessageCachePolicy),
		emojiCache:         NewEmojiCache(config.CacheFlags),
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
	voiceStateCache    VoiceStateCache
	messageCache       MessageCache
	emojiCache         EmojiCache
	guildCache         GuildCache
	channelCache       ChannelCache
	stageInstanceCache StageInstanceCache
}

func (c *cachesImpl) Config() CacheConfig {
	return c.config
}

func (c *cachesImpl) UserCache() UserCache {
	return c.userCache
}
func (c *cachesImpl) RoleCache() RoleCache {
	return c.roleCache
}
func (c *cachesImpl) MemberCache() MemberCache {
	return c.memberCache
}
func (c *cachesImpl) VoiceStateCache() VoiceStateCache {
	return c.voiceStateCache
}
func (c *cachesImpl) MessageCache() MessageCache {
	return c.messageCache
}
func (c *cachesImpl) EmojiCache() EmojiCache {
	return c.emojiCache
}
func (c *cachesImpl) GuildCache() GuildCache {
	return c.guildCache
}
func (c *cachesImpl) ChannelCache() ChannelCache {
	return c.channelCache
}
func (c *cachesImpl) StageInstanceCache() StageInstanceCache {
	return c.stageInstanceCache
}
