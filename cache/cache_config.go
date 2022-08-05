package cache

import (
	"github.com/disgoorg/disgo/discord"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		GuildCachePolicy:               PolicyAll[discord.Guild],
		ChannelCachePolicy:             PolicyAll[discord.Channel],
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

	GuildCachePolicy               Policy[discord.Guild]
	ChannelCachePolicy             Policy[discord.Channel]
	StageInstanceCachePolicy       Policy[discord.StageInstance]
	GuildScheduledEventCachePolicy Policy[discord.GuildScheduledEvent]
	RoleCachePolicy                Policy[discord.Role]
	MemberCachePolicy              Policy[discord.Member]
	ThreadMemberCachePolicy        Policy[discord.ThreadMember]
	PresenceCachePolicy            Policy[discord.Presence]
	VoiceStateCachePolicy          Policy[discord.VoiceState]
	MessageCachePolicy             Policy[discord.Message]
	EmojiCachePolicy               Policy[discord.Emoji]
	StickerCachePolicy             Policy[discord.Sticker]
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

// WithGuildCachePolicy sets the Policy[discord.Guild] of the Config.
func WithGuildCachePolicy(policy Policy[discord.Guild]) ConfigOpt {
	return func(config *Config) {
		config.GuildCachePolicy = policy
	}
}

// WithChannelCachePolicy sets the Policy[discord.Channel] of the Config.
func WithChannelCachePolicy(policy Policy[discord.Channel]) ConfigOpt {
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
