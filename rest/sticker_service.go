package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Service        = (*stickerServiceImpl)(nil)
	_ StickerService = (*stickerServiceImpl)(nil)
)

func NewStickerService(restClient Client) StickerService {
	return &stickerServiceImpl{restClient: restClient}
}

type StickerService interface {
	GetNitroStickerPacks(opts ...RequestOpt) ([]discord.StickerPack, error)
	GetSticker(stickerID snowflake.Snowflake, opts ...RequestOpt) (*discord.Sticker, error)
	GetStickers(guildID snowflake.Snowflake, opts ...RequestOpt) ([]discord.Sticker, error)
	CreateSticker(guildID snowflake.Snowflake, createSticker discord.StickerCreate, opts ...RequestOpt) (*discord.Sticker, error)
	UpdateSticker(guildID snowflake.Snowflake, stickerID snowflake.Snowflake, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (*discord.Sticker, error)
	DeleteSticker(guildID snowflake.Snowflake, stickerID snowflake.Snowflake, opts ...RequestOpt) error
}

type stickerServiceImpl struct {
	restClient Client
}

func (s *stickerServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *stickerServiceImpl) GetNitroStickerPacks(opts ...RequestOpt) (stickerPacks []discord.StickerPack, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetNitroStickerPacks.Compile(nil)
	if err != nil {
		return
	}
	var stickerPacksRs discord.StickerPacks
	err = s.restClient.Do(compiledRoute, nil, &stickerPacksRs, opts...)
	if err == nil {
		stickerPacks = stickerPacksRs.StickerPacks
	}
	return
}

func (s *stickerServiceImpl) GetSticker(stickerID snowflake.Snowflake, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetSticker.Compile(nil, stickerID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) GetStickers(guildID snowflake.Snowflake, opts ...RequestOpt) (stickers []discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildStickers.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &stickers, opts...)
	return
}

func (s *stickerServiceImpl) CreateSticker(guildID snowflake.Snowflake, createSticker discord.StickerCreate, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildSticker.Compile(nil, guildID)
	if err != nil {
		return
	}
	body, err := createSticker.ToBody()
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, body, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) UpdateSticker(guildID snowflake.Snowflake, stickerID snowflake.Snowflake, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, stickerUpdate, &sticker, opts...)
	return
}

func (s *stickerServiceImpl) DeleteSticker(guildID snowflake.Snowflake, stickerID snowflake.Snowflake, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return err
	}
	return s.restClient.Do(compiledRoute, nil, nil, opts...)
}
