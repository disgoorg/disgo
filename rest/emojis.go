package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Emojis = (*emojiImpl)(nil)

func NewEmojis(restClient Client) Emojis {
	return &emojiImpl{restClient: restClient}
}

type Emojis interface {
	GetEmojis(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Emoji, error)
	GetEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) (*discord.Emoji, error)
	CreateEmoji(guildID snowflake.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (*discord.Emoji, error)
	UpdateEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (*discord.Emoji, error)
	DeleteEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) error
}

type emojiImpl struct {
	restClient Client
}

func (s *emojiImpl) GetEmojis(guildID snowflake.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *emojiImpl) GetEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emoji, opts...)
	return
}

func (s *emojiImpl) CreateEmoji(guildID snowflake.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateEmoji.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiCreate, &emoji, opts...)
	return
}

func (s *emojiImpl) UpdateEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiUpdate, &emoji, opts...)
	return
}

func (s *emojiImpl) DeleteEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
