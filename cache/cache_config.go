package cache

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		GuildCachePolicy:               PolicyAll[discord.Guild],
		ChannelCachePolicy:             PolicyAll[discord.GuildChannel],
		StageInstanceCachePolicy:       PolicyAll[discord.StageInstance],
		GuildScheduledEventCachePolicy: PolicyAll[discord.GuildScheduledEvent],
		RoleCachePolicy:                PolicyAll[discord.Role],
		MemberCachePolicy:              PolicyAll[discord.Member],
		ThreadMemberCachePolicy:        PolicyAll[discord.ThreadMember],
		PresenceCachePolicy:            PolicyAll[discord.Presence],
		VoiceStateCachePolicy:          PolicyAll[discord.VoiceState],
		MessageCachePolicy:             PolicyAll[discord.Message],
		EmojiCachePolicy:               PolicyAll[discord.Emoji],
		StickerCachePolicy:             PolicyAll[discord.Sticker],
	}
}

// Config lets you configure your Caches instance.
type Config struct {
	CacheFlags Flags

	SelfUserCache SelfUserCache

	GuildCache       GuildCache
	GuildCachePolicy Policy[discord.Guild]

	ChannelCache       ChannelCache
	ChannelCachePolicy Policy[discord.GuildChannel]

	StageInstanceCache       StageInstanceCache
	StageInstanceCachePolicy Policy[discord.StageInstance]

	GuildScheduledEventCache       GuildScheduledEventCache
	GuildScheduledEventCachePolicy Policy[discord.GuildScheduledEvent]

	RoleCache       RoleCache
	RoleCachePolicy Policy[discord.Role]

	MemberCache       MemberCache
	MemberCachePolicy Policy[discord.Member]

	ThreadMemberCache       ThreadMemberCache
	ThreadMemberCachePolicy Policy[discord.ThreadMember]

	PresenceCache       PresenceCache
	PresenceCachePolicy Policy[discord.Presence]

	VoiceStateCache       VoiceStateCache
	VoiceStateCachePolicy Policy[discord.VoiceState]

	MessageCache       MessageCache
	MessageCachePolicy Policy[discord.Message]

	EmojiCache       EmojiCache
	EmojiCachePolicy Policy[discord.Emoji]

	StickerCache       StickerCache
	StickerCachePolicy Policy[discord.Sticker]
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Caches.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.SelfUserCache == nil {
		c.SelfUserCache = NewSelfUserCache()
	}
	if c.GuildCache == nil {
		c.GuildCache = NewGuildCache(NewCache[discord.Guild](c.CacheFlags, FlagGuilds, c.GuildCachePolicy), NewSet[snowflake.ID](), NewSet[snowflake.ID]())
	}
	if c.ChannelCache == nil {
		c.ChannelCache = NewChannelCache(NewCache[discord.GuildChannel](c.CacheFlags, FlagChannels, c.ChannelCachePolicy))
	}
	if c.StageInstanceCache == nil {
		c.StageInstanceCache = NewStageInstanceCache(NewGroupedCache[discord.StageInstance](c.CacheFlags, FlagStageInstances, c.StageInstanceCachePolicy))
	}
	if c.GuildScheduledEventCache == nil {
		c.GuildScheduledEventCache = NewGuildScheduledEventCache(NewGroupedCache[discord.GuildScheduledEvent](c.CacheFlags, FlagGuildScheduledEvents, c.GuildScheduledEventCachePolicy))
	}
	if c.RoleCache == nil {
		c.RoleCache = NewRoleCache(NewGroupedCache[discord.Role](c.CacheFlags, FlagRoles, c.RoleCachePolicy))
	}
	if c.MemberCache == nil {
		c.MemberCache = NewMemberCache(NewGroupedCache[discord.Member](c.CacheFlags, FlagMembers, c.MemberCachePolicy))
	}
	if c.ThreadMemberCache == nil {
		c.ThreadMemberCache = NewThreadMemberCache(NewGroupedCache[discord.ThreadMember](c.CacheFlags, FlagThreadMembers, c.ThreadMemberCachePolicy))
	}
	if c.PresenceCache == nil {
		c.PresenceCache = NewPresenceCache(NewGroupedCache[discord.Presence](c.CacheFlags, FlagPresences, c.PresenceCachePolicy))
	}
	if c.VoiceStateCache == nil {
		c.VoiceStateCache = NewVoiceStateCache(NewGroupedCache[discord.VoiceState](c.CacheFlags, FlagVoiceStates, c.VoiceStateCachePolicy))
	}
	if c.MessageCache == nil {
		c.MessageCache = NewMessageCache(NewGroupedCache[discord.Message](c.CacheFlags, FlagMessages, c.MessageCachePolicy))
	}
	if c.EmojiCache == nil {
		c.EmojiCache = NewEmojiCache(NewGroupedCache[discord.Emoji](c.CacheFlags, FlagEmojis, c.EmojiCachePolicy))
	}
	if c.StickerCache == nil {
		c.StickerCache = NewStickerCache(NewGroupedCache[discord.Sticker](c.CacheFlags, FlagStickers, c.StickerCachePolicy))
	}
}

// WithCaches sets the Flags of the Config.
func WithCaches(flags ...Flags) ConfigOpt {
	return func(config *Config) {
		config.CacheFlags = config.CacheFlags.Add(flags...)
	}
}

// WithGuildCachePolicy sets the Policy[discord.Guild] of the Config.
func WithGuildCachePolicy(policy Policy[discord.Guild]) ConfigOpt {
	return func(config *Config) {
		config.GuildCachePolicy = policy
	}
}

// WithGuildCache sets the GuildCache of the Config.
func WithGuildCache(guildCache GuildCache) ConfigOpt {
	return func(config *Config) {
		config.GuildCache = guildCache
	}
}

// WithChannelCachePolicy sets the Policy[discord.Channel] of the Config.
func WithChannelCachePolicy(policy Policy[discord.GuildChannel]) ConfigOpt {
	return func(config *Config) {
		config.ChannelCachePolicy = policy
	}
}

// WithChannelCache sets the ChannelCache of the Config.
func WithChannelCache(channelCache ChannelCache) ConfigOpt {
	return func(config *Config) {
		config.ChannelCache = channelCache
	}
}

// WithStageInstanceCachePolicy sets the Policy[discord.Guild] of the Config.
func WithStageInstanceCachePolicy(policy Policy[discord.StageInstance]) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCachePolicy = policy
	}
}

// WithStageInstanceCache sets the StageInstanceCache of the Config.
func WithStageInstanceCache(stageInstanceCache StageInstanceCache) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCache = stageInstanceCache
	}
}

// WithGuildScheduledEventCachePolicy sets the Policy[discord.GuildScheduledEvent] of the Config.
func WithGuildScheduledEventCachePolicy(policy Policy[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

// WithGuildScheduledEventCache sets the GuildScheduledEventCache of the Config.
func WithGuildScheduledEventCache(guildScheduledEventCache GuildScheduledEventCache) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCache = guildScheduledEventCache
	}
}

// WithRoleCachePolicy sets the Policy[discord.Role] of the Config.
func WithRoleCachePolicy(policy Policy[discord.Role]) ConfigOpt {
	return func(config *Config) {
		config.RoleCachePolicy = policy
	}
}

// WithRoleCache sets the RoleCache of the Config.
func WithRoleCache(roleCache RoleCache) ConfigOpt {
	return func(config *Config) {
		config.RoleCache = roleCache
	}
}

// WithMemberCachePolicy sets the Policy[discord.Member] of the Config.
func WithMemberCachePolicy(policy Policy[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCachePolicy = policy
	}
}

// WithMemberCache sets the MemberCache of the Config.
func WithMemberCache(memberCache MemberCache) ConfigOpt {
	return func(config *Config) {
		config.MemberCache = memberCache
	}
}

// WithThreadMemberCachePolicy sets the Policy[discord.ThreadMember] of the Config.
func WithThreadMemberCachePolicy(policy Policy[discord.ThreadMember]) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCachePolicy = policy
	}
}

// WithThreadMemberCache sets the ThreadMemberCache of the Config.
func WithThreadMemberCache(threadMemberCache ThreadMemberCache) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCache = threadMemberCache
	}
}

// WithPresenceCachePolicy sets the Policy[discord.Presence] of the Config.
func WithPresenceCachePolicy(policy Policy[discord.Presence]) ConfigOpt {
	return func(config *Config) {
		config.PresenceCachePolicy = policy
	}
}

// WithPresenceCache sets the PresenceCache of the Config.
func WithPresenceCache(presenceCache PresenceCache) ConfigOpt {
	return func(config *Config) {
		config.PresenceCache = presenceCache
	}
}

// WithVoiceStateCachePolicy sets the Policy[discord.VoiceState] of the Config.
func WithVoiceStateCachePolicy(policy Policy[discord.VoiceState]) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCachePolicy = policy
	}
}

// WithVoiceStateCache sets the VoiceStateCache of the Config.
func WithVoiceStateCache(voiceStateCache VoiceStateCache) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCache = voiceStateCache
	}
}

// WithMessageCachePolicy sets the Policy[discord.Message] of the Config.
func WithMessageCachePolicy(policy Policy[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCachePolicy = policy
	}
}

// WithMessageCache sets the MessageCache of the Config.
func WithMessageCache(messageCache MessageCache) ConfigOpt {
	return func(config *Config) {
		config.MessageCache = messageCache
	}
}

// WithEmojiCachePolicy sets the Policy[discord.Emoji] of the Config.
func WithEmojiCachePolicy(policy Policy[discord.Emoji]) ConfigOpt {
	return func(config *Config) {
		config.EmojiCachePolicy = policy
	}
}

// WithEmojiCache sets the EmojiCache of the Config.
func WithEmojiCache(emojiCache EmojiCache) ConfigOpt {
	return func(config *Config) {
		config.EmojiCache = emojiCache
	}
}

// WithStickerCachePolicy sets the Policy[discord.Sticker] of the Config.
func WithStickerCachePolicy(policy Policy[discord.Sticker]) ConfigOpt {
	return func(config *Config) {
		config.StickerCachePolicy = policy
	}
}

// WithStickerCache sets the StickerCache of the Config.
func WithStickerCache(stickerCache StickerCache) ConfigOpt {
	return func(config *Config) {
		config.StickerCache = stickerCache
	}
}
