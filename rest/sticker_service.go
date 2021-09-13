package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ StickerService = (*stickerServiceImpl)(nil)

func NewStickerService(restClient Client) StickerService {
	return &stickerServiceImpl{restClient: restClient}
}

type StickerService interface {
	GetNitroStickerPacks(opts ...RequestOpt) ([]discord.StickerPack, Error)
	GetSticker(stickerID discord.Snowflake, opts ...RequestOpt) (*discord.Sticker, Error)
	GetGuildStickers(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Sticker, Error)
	CreateGuildSticker(guildID discord.Snowflake, createSticker discord.StickerCreate, opts ...RequestOpt) (*discord.Sticker, Error)
	UpdateGuildSticker(guildID discord.Snowflake, stickerID discord.Snowflake, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (*discord.Sticker, Error)
	DeleteGuildSticker(guildID discord.Snowflake, stickerID discord.Snowflake, opts ...RequestOpt) Error
}

type stickerServiceImpl struct {
	restClient Client
}

func (s *stickerServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *stickerServiceImpl) GetNitroStickerPacks(opts ...RequestOpt) (stickerPacks []discord.StickerPack, rErr Error) {
	compiledRoute, err := route.GetNitroStickerPacks.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	var stickerPacksRs discord.StickerPacks
	rErr = s.restClient.Do(compiledRoute, nil, &stickerPacksRs, opts...)
	if rErr == nil {
		stickerPacks = stickerPacksRs.StickerPacks
	}
	return
}

func (s *stickerServiceImpl) GetSticker(stickerID discord.Snowflake, opts ...RequestOpt) (sticker *discord.Sticker, rErr Error) {
	compiledRoute, err := route.GetSticker.Compile(nil, stickerID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) GetGuildStickers(guildID discord.Snowflake, opts ...RequestOpt) (stickers []discord.Sticker, rErr Error) {
	compiledRoute, err := route.GetGuildStickers.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &stickers, opts...)
	return
}

func (s *stickerServiceImpl) CreateGuildSticker(guildID discord.Snowflake, createSticker discord.StickerCreate, opts ...RequestOpt) (sticker *discord.Sticker, rErr Error) {
	compiledRoute, err := route.CreateGuildSticker.Compile(nil, guildID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	body, err := createSticker.ToBody()
	if err != nil {
		rErr = NewError(nil, err)

		return
	}
	rErr = s.restClient.Do(compiledRoute, body, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) UpdateGuildSticker(guildID discord.Snowflake, stickerID discord.Snowflake, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (sticker *discord.Sticker, rErr Error) {
	compiledRoute, err := route.UpdateGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, stickerUpdate, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) DeleteGuildSticker(guildID discord.Snowflake, stickerID discord.Snowflake, opts ...RequestOpt) Error {
	compiledRoute, err := route.DeleteGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return NewError(nil, err)
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
