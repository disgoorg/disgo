package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ Stickers = (*stickerImpl)(nil)

func NewStickers(client Client) Stickers {
	return &stickerImpl{client: client}
}

type Stickers interface {
	GetNitroStickerPacks(opts ...RequestOpt) ([]discord.StickerPack, error)
	GetSticker(stickerID snowflake.ID, opts ...RequestOpt) (*discord.Sticker, error)
	GetStickers(guildID snowflake.ID, opts ...RequestOpt) ([]discord.Sticker, error)
	CreateSticker(guildID snowflake.ID, createSticker discord.StickerCreate, opts ...RequestOpt) (*discord.Sticker, error)
	UpdateSticker(guildID snowflake.ID, stickerID snowflake.ID, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (*discord.Sticker, error)
	DeleteSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) error
}

type stickerImpl struct {
	client Client
}

func (s *stickerImpl) GetNitroStickerPacks(opts ...RequestOpt) (stickerPacks []discord.StickerPack, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetNitroStickerPacks.Compile(nil)
	if err != nil {
		return
	}
	var stickerPacksRs discord.StickerPacks
	err = s.client.Do(compiledRoute, nil, &stickerPacksRs, opts...)
	if err == nil {
		stickerPacks = stickerPacksRs.StickerPacks
	}
	return
}

func (s *stickerImpl) GetSticker(stickerID snowflake.ID, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetSticker.Compile(nil, stickerID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &sticker, opts...)
	return
}

func (s *stickerImpl) GetStickers(guildID snowflake.ID, opts ...RequestOpt) (stickers []discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGuildStickers.Compile(nil, guildID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &stickers, opts...)
	return
}

func (s *stickerImpl) CreateSticker(guildID snowflake.ID, createSticker discord.StickerCreate, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.CreateGuildSticker.Compile(nil, guildID)
	if err != nil {
		return
	}
	body, err := createSticker.ToBody()
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, body, &sticker, opts...)
	return
}

func (s *stickerImpl) UpdateSticker(guildID snowflake.ID, stickerID snowflake.ID, stickerUpdate discord.StickerUpdate, opts ...RequestOpt) (sticker *discord.Sticker, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.UpdateGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, stickerUpdate, &sticker, opts...)
	return
}

func (s *stickerImpl) DeleteSticker(guildID snowflake.ID, stickerID snowflake.ID, opts ...RequestOpt) error {
	compiledRoute, err := route.DeleteGuildSticker.Compile(nil, guildID, stickerID)
	if err != nil {
		return err
	}
	return s.client.Do(compiledRoute, nil, nil, opts...)
}
