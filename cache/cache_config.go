package cache

import (
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

	GuildCache       Cache[discord.Guild]
	GuildCachePolicy Policy[discord.Guild]

	ChannelCache       Cache[discord.GuildChannel]
	ChannelCachePolicy Policy[discord.GuildChannel]

	StageInstanceCache       GroupedCache[discord.StageInstance]
	StageInstanceCachePolicy Policy[discord.StageInstance]

	GuildScheduledEventCache       GroupedCache[discord.GuildScheduledEvent]
	GuildScheduledEventCachePolicy Policy[discord.GuildScheduledEvent]

	RoleCache       GroupedCache[discord.Role]
	RoleCachePolicy Policy[discord.Role]

	MemberCache       GroupedCache[discord.Member]
	MemberCachePolicy Policy[discord.Member]

	ThreadMemberCache       GroupedCache[discord.ThreadMember]
	ThreadMemberCachePolicy Policy[discord.ThreadMember]

	PresenceCache       GroupedCache[discord.Presence]
	PresenceCachePolicy Policy[discord.Presence]

	VoiceStateCache       GroupedCache[discord.VoiceState]
	VoiceStateCachePolicy Policy[discord.VoiceState]

	MessageCache       GroupedCache[discord.Message]
	MessageCachePolicy Policy[discord.Message]

	EmojiCache       GroupedCache[discord.Emoji]
	EmojiCachePolicy Policy[discord.Emoji]

	StickerCache       GroupedCache[discord.Sticker]
	StickerCachePolicy Policy[discord.Sticker]
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Caches.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.GuildCache == nil {
		c.GuildCache = NewCache[discord.Guild](c.CacheFlags, FlagGuilds, c.GuildCachePolicy)
	}
	if c.ChannelCache == nil {
		c.ChannelCache = NewCache[discord.GuildChannel](c.CacheFlags, FlagChannels, c.ChannelCachePolicy)
	}
	if c.StageInstanceCache == nil {
		c.StageInstanceCache = NewGroupedCache[discord.StageInstance](c.CacheFlags, FlagStageInstances, c.StageInstanceCachePolicy)
	}
	if c.GuildScheduledEventCache == nil {
		c.GuildScheduledEventCache = NewGroupedCache[discord.GuildScheduledEvent](c.CacheFlags, FlagGuildScheduledEvents, c.GuildScheduledEventCachePolicy)
	}
	if c.RoleCache == nil {
		c.RoleCache = NewGroupedCache[discord.Role](c.CacheFlags, FlagRoles, c.RoleCachePolicy)
	}
	if c.MemberCache == nil {
		c.MemberCache = NewGroupedCache[discord.Member](c.CacheFlags, FlagMembers, c.MemberCachePolicy)
	}
	if c.ThreadMemberCache == nil {
		c.ThreadMemberCache = NewGroupedCache[discord.ThreadMember](c.CacheFlags, FlagThreadMembers, c.ThreadMemberCachePolicy)
	}
	if c.PresenceCache == nil {
		c.PresenceCache = NewGroupedCache[discord.Presence](c.CacheFlags, FlagPresences, c.PresenceCachePolicy)
	}
	if c.VoiceStateCache == nil {
		c.VoiceStateCache = NewGroupedCache[discord.VoiceState](c.CacheFlags, FlagVoiceStates, c.VoiceStateCachePolicy)
	}
	if c.MessageCache == nil {
		c.MessageCache = NewGroupedCache[discord.Message](c.CacheFlags, FlagMessages, c.MessageCachePolicy)
	}
	if c.EmojiCache == nil {
		c.EmojiCache = NewGroupedCache[discord.Emoji](c.CacheFlags, FlagEmojis, c.EmojiCachePolicy)
	}
	if c.StickerCache == nil {
		c.StickerCache = NewGroupedCache[discord.Sticker](c.CacheFlags, FlagStickers, c.StickerCachePolicy)
	}
}

// WithCacheFlags sets the Flags of the Config.
func WithCacheFlags(flags ...Flags) ConfigOpt {
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

// WithChannelCachePolicy sets the Policy[discord.Channel] of the Config.
func WithChannelCachePolicy(policy Policy[discord.GuildChannel]) ConfigOpt {
	return func(config *Config) {
		config.ChannelCachePolicy = policy
	}
}

// WithStageInstanceCachePolicy sets the Policy[discord.Guild] of the Config.
func WithStageInstanceCachePolicy(policy Policy[discord.StageInstance]) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCachePolicy = policy
	}
}

// WithGuildScheduledEventCachePolicy sets the Policy[discord.GuildScheduledEvent] of the Config.
func WithGuildScheduledEventCachePolicy(policy Policy[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

// WithRoleCachePolicy sets the Policy[discord.Role] of the Config.
func WithRoleCachePolicy(policy Policy[discord.Role]) ConfigOpt {
	return func(config *Config) {
		config.RoleCachePolicy = policy
	}
}

// WithMemberCachePolicy sets the Policy[discord.Member] of the Config.
func WithMemberCachePolicy(policy Policy[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCachePolicy = policy
	}
}

// WithThreadMemberCachePolicy sets the Policy[discord.ThreadMember] of the Config.
func WithThreadMemberCachePolicy(policy Policy[discord.ThreadMember]) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCachePolicy = policy
	}
}

// WithPresenceCachePolicy sets the Policy[discord.Presence] of the Config.
func WithPresenceCachePolicy(policy Policy[discord.Presence]) ConfigOpt {
	return func(config *Config) {
		config.PresenceCachePolicy = policy
	}
}

// WithVoiceStateCachePolicy sets the Policy[discord.VoiceState] of the Config.
func WithVoiceStateCachePolicy(policy Policy[discord.VoiceState]) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCachePolicy = policy
	}
}

// WithMessageCachePolicy sets the Policy[discord.Message] of the Config.
func WithMessageCachePolicy(policy Policy[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCachePolicy = policy
	}
}

// WithEmojiCachePolicy sets the Policy[discord.Emoji] of the Config.
func WithEmojiCachePolicy(policy Policy[discord.Emoji]) ConfigOpt {
	return func(config *Config) {
		config.EmojiCachePolicy = policy
	}
}

// WithStickerCachePolicy sets the Policy[discord.Sticker] of the Config.
func WithStickerCachePolicy(policy Policy[discord.Sticker]) ConfigOpt {
	return func(config *Config) {
		config.StickerCachePolicy = policy
	}
}
