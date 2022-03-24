package cache

import (
	"github.com/disgoorg/disgo/discord"
)

func DefaultConfig() *Config {
	return &Config{
		CacheFlags:         FlagsDefault,
		MemberCachePolicy:  MemberCachePolicyDefault,
		MessageCachePolicy: MessageCachePolicyDefault,
	}
}

type Config struct {
	CacheFlags         Flags
	MemberCachePolicy  Policy[discord.Member]
	MessageCachePolicy Policy[discord.Message]
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithCacheFlags(cacheFlags ...Flags) ConfigOpt {
	return func(config *Config) {
		var flags Flags
		for _, flag := range cacheFlags {
			flags = flags.Add(flag)
		}
		config.CacheFlags = flags
	}
}

func WithMemberCachePolicy(memberCachePolicy Policy[discord.Member]) ConfigOpt {
	return func(config *Config) {
		config.MemberCachePolicy = memberCachePolicy
	}
}

func WithMessageCachePolicy(messageCachePolicy Policy[discord.Message]) ConfigOpt {
	return func(config *Config) {
		config.MessageCachePolicy = messageCachePolicy
	}
}
