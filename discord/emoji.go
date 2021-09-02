package discord

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	ID            Snowflake   `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"` // may be empty for deleted emojis
	Roles         []Snowflake `json:"roles,omitempty"`
	Creator       *User       `json:"creator,omitempty"`
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     bool        `json:"available,omitempty"`

	GuildID Snowflake `json:"guild_id,omitempty"`
}

type EmojiCreate struct {
	Name  string      `json:"name"`
	Image *Icon       `json:"image"`
	Roles []Snowflake `json:"roles,omitempty"`
}

type EmojiUpdate struct {
	Name  string      `json:"name,omitempty"`
	Roles []Snowflake `json:"roles,omitempty"`
}
