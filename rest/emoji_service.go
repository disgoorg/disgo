package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ EmojiService = (*emojiServiceImpl)(nil)

func NewEmojiService(restClient Client) EmojiService {
	return &emojiServiceImpl{restClient: restClient}
}

type EmojiService interface {
	Service
	GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Emoji, Error)
	GetEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) (*discord.Emoji, Error)
	CreateEmoji(guildID discord.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (*discord.Emoji, Error)
	UpdateEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (*discord.Emoji, Error)
	DeleteEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) Error
}

type emojiServiceImpl struct {
	restClient Client
}

func (s *emojiServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *emojiServiceImpl) GetEmojis(guildID discord.Snowflake, opts ...RequestOpt) (emojis []discord.Emoji, rErr Error) {
	compiledRoute, err := route.GetEmojis.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &emojis, opts...)
	return
}

func (s *emojiServiceImpl) GetEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) (emoji *discord.Emoji, rErr Error) {
	compiledRoute, err := route.GetEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) CreateEmoji(guildID discord.Snowflake, emojiCreate discord.EmojiCreate, opts ...RequestOpt) (emoji *discord.Emoji, rErr Error) {
	compiledRoute, err := route.CreateEmoji.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, emojiCreate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) UpdateEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, emojiUpdate discord.EmojiUpdate, opts ...RequestOpt) (emoji *discord.Emoji, rErr Error) {
	compiledRoute, err := route.UpdateEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, emojiUpdate, &emoji, opts...)
	return
}

func (s *emojiServiceImpl) DeleteEmoji(guildID discord.Snowflake, emojiID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteEmoji.Compile(nil, guildID, emojiID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
