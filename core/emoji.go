package core

import "github.com/DisgoOrg/disgo/discord"

// NewEmote creates a new custom Emoji with the given parameters
//goland:noinspection GoUnusedExportedFunction
func NewEmote(name string, emoteID discord.Snowflake, animated bool) *Emoji {
	return &Emoji{
		Emoji: discord.Emoji{
			Name:     name,
			ID:       emoteID,
			Animated: animated,
		},
	}
}

// NewEmoji creates a new emoji with the given unicode
//goland:noinspection GoUnusedExportedFunction
func NewEmoji(name string) *Emoji {
	return &Emoji{
		Emoji: discord.Emoji{
			Name: name,
		},
	}
}

type Emoji struct {
	discord.Emoji
	Disgo   Disgo
	GuildID discord.Snowflake
}

// Guild returns the Guild of the Emoji from the Cache
func (e *Emoji) Guild() *Guild {
	return e.Disgo.Cache().GuildCache().Get(e.GuildID)
}

// Mention returns the string used to send the Emoji
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
