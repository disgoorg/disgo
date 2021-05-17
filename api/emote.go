package api

// An Emote allows you to interact with custom emojis in discord.
type Emote struct {
	Disgo    Disgo
	ID       Snowflake `json:"id"`
	GuildID  Snowflake `json:"guild_id"`
	Name     string    `json:"name"`
	Animated bool      `json:"animated,omitempty"`
}

// Guild returns the Guild of the Emote from the Cache
func (e Emote) Guild() *Guild {
	return e.Disgo.Cache().Guild(e.GuildID)
}

// Mention returns the string used to send the emoji
func (e Emote) Mention() string {
	start := "<:"
	if e.Animated {
		start = "<a:"
	}
	return start + e.Name + ":" + e.ID.String() + ">"
}

// String formats the Emote as string
func (e Emote) String() string {
	return e.Mention()
}

// Reaction returns the identifier used for adding and removing reactions for messages in discord
func (e Emote) Reaction() string {
	return ":" + e.Name + ":" + e.ID.String()
}
