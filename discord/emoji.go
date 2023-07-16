package discord

import (
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var _ Mentionable = (*Emoji)(nil)

// Emoji allows you to interact with emojis & emotes
type Emoji struct {
	ID            snowflake.ID   `json:"id,omitempty"`
	GuildID       snowflake.ID   `json:"guild_id,omitempty"` // not present in the API but we need it
	Name          string         `json:"name,omitempty"`     // may be empty for deleted emojis
	Roles         []snowflake.ID `json:"roles,omitempty"`
	Creator       *User          `json:"creator,omitempty"`
	RequireColons bool           `json:"require_colons,omitempty"`
	Managed       bool           `json:"managed,omitempty"`
	Animated      bool           `json:"animated,omitempty"`
	Available     bool           `json:"available,omitempty"`
}

// Reaction returns a string used for manipulating with reactions. May be empty if the Name is empty
func (e Emoji) Reaction() string {
	if e.Name == "" {
		return ""
	}
	return reaction(e.Name, e.ID)
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
	return formatAssetURL(CustomEmoji, opts, e.ID)
}

func (e Emoji) CreatedAt() time.Time {
	if e.ID == 0 {
		return time.Time{}
	}
	return e.ID.Time()
}

type EmojiCreate struct {
	Name  string         `json:"name"`
	Image Icon           `json:"image"`
	Roles []snowflake.ID `json:"roles,omitempty"`
}

type EmojiUpdate struct {
	Name  *string         `json:"name,omitempty"`
	Roles *[]snowflake.ID `json:"roles,omitempty"`
}

type PartialEmoji struct {
	ID       *snowflake.ID `json:"id"`
	Name     *string       `json:"name"`
	Animated bool          `json:"animated"`
}

// Reaction returns a string used for manipulating with reactions. May be empty if the Name is nil
func (e PartialEmoji) Reaction() string {
	if e.Name == nil {
		return ""
	}
	if e.ID == nil {
		return *e.Name
	}
	return reaction(*e.Name, *e.ID)
}

func reaction(name string, id snowflake.ID) string {
	return fmt.Sprintf("%s:%s", name, id)
}
