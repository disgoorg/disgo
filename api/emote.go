package api

// NewEmote creates a new custom Emoji with the given parameters
func NewEmote(name string, emoteID Snowflake) *Emoji {
	return &Emoji{Name: name, ID: emoteID, Animated: false}
}

// NewAnimatedEmote creates a new animated custom Emoji with the given parameters
func NewAnimatedEmote(name string, emoteID Snowflake) *Emoji {
	return &Emoji{Name: name, ID: emoteID, Animated: true}
}

// NewEmoji creates a new emoji with the given unicode
func NewEmoji(name string) *Emoji {
	return &Emoji{Name: name}
}

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	Disgo    Disgo
	GuildID  Snowflake `json:"guild_id,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       Snowflake `json:"id,omitempty"`
	Animated bool      `json:"animated,omitempty"`
}

// Guild returns the Guild of the Emoji from the Cache
func (e *Emoji) Guild() *Guild {
	return e.Disgo.Cache().Guild(e.GuildID)
}

// Mention returns the string used to send the emoji
func (e *Emoji) Mention() string {
	start := "<:"
	if e.Animated {
		start = "<a:"
	}
	return start + e.Name + ":" + e.ID.String() + ">"
}

// String formats the Emoji as string
func (e *Emoji) String() string {
	return e.Mention()
}

// Reaction returns the identifier used for adding and removing reactions for messages in discord
func (e *Emoji) Reaction() string {
	return ":" + e.Name + ":" + e.ID.String()
}
