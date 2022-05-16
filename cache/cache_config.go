package cache

import (
	"github.com/disgoorg/disgo/discord"
)

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

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithCacheFlags(flags ...Flags) ConfigOpt {
	return func(config *Config) {
		config.CacheFlags = config.CacheFlags.Add(flags...)
	}
}

func WithGuildCachePolicy(policy Policy[discord.Guild]) ConfigOpt {
	return func(config *Config) {
		config.GuildCachePolicy = policy
	}
}

func WithChannelCachePolicy(policy Policy[discord.Channel]) ConfigOpt {
	return func(config *Config) {
		config.ChannelCachePolicy = policy
	}
}

func WithStageInstanceCachePolicy(policy Policy[discord.StageInstance]) ConfigOpt {
	return func(config *Config) {
		config.StageInstanceCachePolicy = policy
	}
}

func WithGuildScheduledEventCachePolicy(policy Policy[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *Config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

func WithRoleCachePolicy(policy Policy[discord.Role]) ConfigOpt {
	return func(config *Config) {
		config.RoleCachePolicy = policy
	}
}

func WithMemberCachePolicy(policy Policy[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCachePolicy = policy
	}
}

func WithThreadMemberCachePolicy(policy Policy[discord.ThreadMember]) ConfigOpt {
	return func(config *Config) {
		config.ThreadMemberCachePolicy = policy
	}
}

func WithPresenceCachePolicy(policy Policy[discord.Presence]) ConfigOpt {
	return func(config *Config) {
		config.PresenceCachePolicy = policy
	}
}

func WithVoiceStateCachePolicy(policy Policy[discord.VoiceState]) ConfigOpt {
	return func(config *Config) {
		config.VoiceStateCachePolicy = policy
	}
}

func WithMessageCachePolicy(policy Policy[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCachePolicy = policy
	}
}

func WithEmojiCachePolicy(policy Policy[discord.Emoji]) ConfigOpt {
	return func(config *Config) {
		config.EmojiCachePolicy = policy
	}
}

func WithStickerCachePolicy(policy Policy[discord.Sticker]) ConfigOpt {
	return func(config *Config) {
		config.StickerCachePolicy = policy
	}
}
