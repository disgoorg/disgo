package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

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
	Bot *Bot
}

func (e *Emoji) URL(size int) string {
	fileExtension := route.PNG
	if e.Animated {
		fileExtension = route.PNG
	}
	compiledRoute, _ := route.CustomEmoji.Compile(nil, fileExtension, size, e.ID)
	return compiledRoute.URL()
}

// Guild returns the Guild of the Emoji from the Caches
func (e *Emoji) Guild() *Guild {
	return e.Bot.Caches.GuildCache().Get(e.GuildID)
}

func (e *Emoji) Update(emojiUpdate discord.EmojiUpdate, opts ...rest.RequestOpt) (*Emoji, error) {
	emoji, err := e.Bot.RestServices.EmojiService().UpdateEmoji(e.GuildID, e.ID, emojiUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return e.Bot.EntityBuilder.CreateEmoji(e.GuildID, *emoji, CacheStrategyNoWs), nil
}

func (e *Emoji) Delete(opts ...rest.RequestOpt) error {
	return e.Bot.RestServices.EmojiService().DeleteEmoji(e.GuildID, e.ID, opts...)
}

// Reaction returns the identifier used for adding and removing reactions for messages in discord
func (e *Emoji) Reaction() string {
	return ":" + e.Name + ":" + e.ID.String()
}
