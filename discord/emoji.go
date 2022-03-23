package discord

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/rest/route"
)

var _ Mentionable = (*Emoji)(nil)

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	ID            snowflake.Snowflake   `json:"id,omitempty"`
	Name          string                `json:"name,omitempty"` // may be empty for deleted emojis
	Roles         []snowflake.Snowflake `json:"roles,omitempty"`
	Creator       *User                 `json:"creator,omitempty"`
	RequireColons bool                  `json:"require_colons,omitempty"`
	Managed       bool                  `json:"managed,omitempty"`
	Animated      bool                  `json:"animated,omitempty"`
	Available     bool                  `json:"available,omitempty"`

	GuildID snowflake.Snowflake `json:"guild_id,omitempty"`
}

// Mention returns the string used to send the Emoji
func (e Emoji) Mention() string {
	if e.Animated {
		return AnimatedEmojiMention(e.ID, e.Name)
	}
	return EmojiMention(e.ID, e.Name)
}

// String formats the Emoji as string
func (e Emoji) String() string {
	return e.Mention()
}

func (e Emoji) URL(opts ...CDNOpt) string {
	if avatar := formatAssetURL(route.CustomEmoji, opts, e.ID); avatar != nil {
		return *avatar
	}
	return ""
}

type EmojiCreate struct {
	Name  string                `json:"name"`
	Image Icon                  `json:"image"`
	Roles []snowflake.Snowflake `json:"roles,omitempty"`
}

type EmojiUpdate struct {
	Name  string                `json:"name,omitempty"`
	Roles []snowflake.Snowflake `json:"roles,omitempty"`
}

type ReactionEmoji struct {
	ID       snowflake.Snowflake `json:"id,omitempty"`
	Name     string              `json:"name,omitempty"`
	Animated bool                `json:"animated"`
}
