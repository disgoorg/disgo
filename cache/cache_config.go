package cache

import (
	"github.com/disgoorg/disgo/discord"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		CacheFlags:                     FlagsDefault,
		GuildCachePolicy:               PolicyDefault[discord.Guild],
		ChannelCachePolicy:             PolicyDefault[discord.Channel],
		StageInstanceCachePolicy:       PolicyDefault[discord.StageInstance],
		GuildScheduledEventCachePolicy: PolicyDefault[discord.GuildScheduledEvent],
		RoleCachePolicy:                PolicyDefault[discord.Role],
		MemberCachePolicy:              PolicyDefault[discord.Member],
		ThreadMemberCachePolicy:        PolicyDefault[discord.ThreadMember],
		PresenceCachePolicy:            PolicyDefault[discord.Presence],
		VoiceStateCachePolicy:          PolicyDefault[discord.VoiceState],
		MessageCachePolicy:             PolicyDefault[discord.Message],
		EmojiCachePolicy:               PolicyDefault[discord.Emoji],
		StickerCachePolicy:             PolicyDefault[discord.Sticker],
	}
}

// Config lets you configure your Caches instance.
type Config struct {
	CacheFlags Flags

	GuildCache       GuildCache
	GuildCachePolicy Policy[discord.Guild]

	ChannelCache       ChannelCache
	ChannelCachePolicy Policy[discord.Channel]

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
}

// WithCacheFlags sets the Flags of the Config.
func WithCacheFlags(flags ...Flags) ConfigOpt {
	return func(config *Config) {
		config.CacheFlags = config.CacheFlags.Add(flags...)
	}
}

func WithGuildCache(cache GuildCache) ConfigOpt {
	return func(config *Config) {
		config.GuildCache = cache
	}
}

// WithGuildCachePolicy sets the Policy[discord.Guild] of the Config.
func WithGuildCachePolicy(policy Policy[discord.Guild]) ConfigOpt {
	return func(config *Config) {
		config.GuildCachePolicy = policy
	}
}

func WithChannelCache(cache ChannelCache) ConfigOpt {
	return func(config *Config) {
		config.ChannelCache = cache
	}
}

// WithChannelCachePolicy sets the Policy[discord.Channel] of the Config.
func WithChannelCachePolicy(policy Policy[discord.Channel]) ConfigOpt {
	return func(config *Config) {
		config.ChannelCachePolicy = policy
	}
}

func WithStageInstanceCache(cache GroupedCache[discord.StageInstance]) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCache = cache
	}
}

// WithStageInstanceCachePolicy sets the Policy[discord.Guild] of the Config.
func WithStageInstanceCachePolicy(policy Policy[discord.StageInstance]) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCachePolicy = policy
	}
}

func WithGuildScheduledEventCache(cache GroupedCache[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCache = cache
	}
}

// WithGuildScheduledEventCachePolicy sets the Policy[discord.GuildScheduledEvent] of the Config.
func WithGuildScheduledEventCachePolicy(policy Policy[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

func WithRoleCache(cache GroupedCache[discord.Role]) ConfigOpt {
	return func(config *Config) {
		config.RoleCache = cache
	}
}

// WithRoleCachePolicy sets the Policy[discord.Role] of the Config.
func WithRoleCachePolicy(policy Policy[discord.Role]) ConfigOpt {
	return func(config *Config) {
		config.RoleCachePolicy = policy
	}
}

func WithMemberCache(cache GroupedCache[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCache = cache
	}
}

// WithMemberCachePolicy sets the Policy[discord.Member] of the Config.
func WithMemberCachePolicy(policy Policy[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCachePolicy = policy
	}
}

func WithThreadMemberCache(cache GroupedCache[discord.ThreadMember]) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCache = cache
	}
}

// WithThreadMemberCachePolicy sets the Policy[discord.ThreadMember] of the Config.
func WithThreadMemberCachePolicy(policy Policy[discord.ThreadMember]) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCachePolicy = policy
	}
}

func WithPresenceCache(cache GroupedCache[discord.Presence]) ConfigOpt {
	return func(config *Config) {
		config.PresenceCache = cache
	}
}

// WithPresenceCachePolicy sets the Policy[discord.Presence] of the Config.
func WithPresenceCachePolicy(policy Policy[discord.Presence]) ConfigOpt {
	return func(config *Config) {
		config.PresenceCachePolicy = policy
	}
}

func WithVoiceStateCache(cache GroupedCache[discord.VoiceState]) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCache = cache
	}
}

// WithVoiceStateCachePolicy sets the Policy[discord.VoiceState] of the Config.
func WithVoiceStateCachePolicy(policy Policy[discord.VoiceState]) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCachePolicy = policy
	}
}

func WithMessageCache(cache GroupedCache[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCache = cache
	}
}

// WithMessageCachePolicy sets the Policy[discord.Message] of the Config.
func WithMessageCachePolicy(policy Policy[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCachePolicy = policy
	}
}

func WithEmojiCache(cache GroupedCache[discord.Emoji]) ConfigOpt {
	return func(config *Config) {
		config.EmojiCache = cache
	}
}

// WithEmojiCachePolicy sets the Policy[discord.Emoji] of the Config.
func WithEmojiCachePolicy(policy Policy[discord.Emoji]) ConfigOpt {
	return func(config *Config) {
		config.EmojiCachePolicy = policy
	}
}

func WithStickerCache(cache GroupedCache[discord.Sticker]) ConfigOpt {
	return func(config *Config) {
		config.StickerCache = cache
	}
}

// WithStickerCachePolicy sets the Policy[discord.Sticker] of the Config.
func WithStickerCachePolicy(policy Policy[discord.Sticker]) ConfigOpt {
	return func(config *Config) {
		config.StickerCachePolicy = policy
	}
}
