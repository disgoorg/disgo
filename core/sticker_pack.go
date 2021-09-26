package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

type StickerPack struct {
	discord.StickerPack
	Bot      *Bot
	Stickers []*Sticker
}

func (p *StickerPack) CoverSticker() *Sticker {
	for _, sticker := range p.Stickers {
		if sticker.ID == p.CoverStickerID {
			return sticker
		}
	}
	return nil
}

func (p *StickerPack) BannerAsset(size int) string {
	compiledRoute, _ := route.StickerPackBanner.Compile(nil, route.PNG, size, p.BannerAssetID)
	return compiledRoute.URL()
}
