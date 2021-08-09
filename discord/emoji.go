package discord

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	GuildID  Snowflake  `json:"guild_id,omitempty"`
	Name     string     `json:"name,omitempty"`
	ID       Snowflake  `json:"id,omitempty"`
	Animated bool       `json:"animated,omitempty"`
}
