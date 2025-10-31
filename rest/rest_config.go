package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func defaultConfig() config {
	return config{
		DefaultAllowedMentions: discord.AllowedMentions{
			Parse:       []discord.AllowedMentionType{discord.AllowedMentionTypeUsers, discord.AllowedMentionTypeRoles, discord.AllowedMentionTypeEveryone},
			Roles:       []snowflake.ID{},
			Users:       []snowflake.ID{},
			RepliedUser: true,
		},
	}
}

type config struct {
	DefaultAllowedMentions discord.AllowedMentions
}

// ConfigOpt can be used to supply optional parameters to New
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithDefaultAllowedMentions(mentions discord.AllowedMentions) ConfigOpt {
	return func(config *config) {
		config.DefaultAllowedMentions = mentions
	}
}