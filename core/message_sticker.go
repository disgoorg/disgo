package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type MessageSticker struct {
	discord.MessageSticker
	Bot *Bot
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
