package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
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
	GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Emoji, error)
	GetEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) (*discord.Emoji, error)
	CreateEmoji(guildID discord.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (*discord.Emoji, error)
	UpdateEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (*discord.Emoji, error)
	DeleteEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) error
}

type emojiServiceImpl struct {
	restClient Client
}

func (s *emojiServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *emojiServiceImpl) GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *emojiServiceImpl) GetEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) CreateEmoji(guildID discord.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateEmoji.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiCreate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) UpdateEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (emoji *discord.Emoji, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, emojiUpdate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) DeleteEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
