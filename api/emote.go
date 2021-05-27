package api

// NewEmote creates a new custom Emote with the given parameters
func NewEmote(name string, emoteID Snowflake) *Emote {
	return &Emote{Name: name, ID: emoteID, Animated: false}
}

// NewAnimatedEmote creates a new animated custom Emote with the given parameters
func NewAnimatedEmote(name string, emoteID Snowflake) *Emote {
	return &Emote{Name: name, ID: emoteID, Animated: true}
}

// NewEmoji creates a new emoji with the given unicode
func NewEmoji(name string) *Emote {
	return &Emote{Name: name}
}

// Emote allows you to interact with emojis & emotes
type Emote struct {
	Disgo    Disgo
	GuildID  Snowflake `json:"guild_id,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       Snowflake `json:"id,omitempty"`
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
