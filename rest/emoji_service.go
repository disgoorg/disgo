package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Service      = (*emojiServiceImpl)(nil)
	_ EmojiService = (*emojiServiceImpl)(nil)
)

func NewEmojiService(restClient Client) EmojiService {
	return &emojiServiceImpl{restClient: restClient}
}

type EmojiService interface {
	Service
	GetEmojis(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Emoji, error)
	GetEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) (*discord.Emoji, error)
	CreateEmoji(guildID snowflake.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (*discord.Emoji, error)
	UpdateEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (*discord.Emoji, error)
	DeleteEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) error
}

type emojiServiceImpl struct {
	restClient Client
}

func (s *emojiServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *emojiServiceImpl) GetEmojis(guildID snowflake.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *emojiServiceImpl) GetEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) CreateEmoji(guildID snowflake.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateEmoji.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiCreate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) UpdateEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiUpdate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) DeleteEmoji(guildID snowflake.Snowflake, emojiID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
