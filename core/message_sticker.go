package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type MessageSticker struct {
	discord.MessageSticker
	Bot *Bot
}

func (s *MessageSticker) URL(size int) string {
	fileExtension := route.PNG
	if s.FormatType == discord.StickerFormatTypeLottie {
		fileExtension = route.Lottie
	}
	compiledRoute, _ := route.CustomSticker.Compile(nil, fileExtension, size, s.ID)
	return compiledRoute.URL()
}

func (s *MessageSticker) GetSticker() (*Sticker, rest.Error) {
	coreSticker := s.Bot.Caches.StickerCache().FindFirst(func(sticker *Sticker) bool { return sticker.ID == s.ID })
	if coreSticker != nil {
		return coreSticker, nil
	}

	sticker, err := s.Bot.RestServices.StickerService().GetSticker(s.ID)
	if err != nil {
		return nil, err
	}
	return s.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}
