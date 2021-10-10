package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Sticker struct {
	discord.Sticker
	Bot  *Bot
	User *User
}

func (s *Sticker) URL(size int) string {
	fileExtension := route.PNG
	if s.FormatType == discord.StickerFormatTypeLottie {
		fileExtension = route.Lottie
	}
	compiledRoute, _ := route.CustomSticker.Compile(nil, fileExtension, size, s.ID)
	return compiledRoute.URL()
}

func (s *Sticker) Guild() *Guild {
	if s.Type != discord.StickerTypeGuild {
		return nil
	}
	return s.Bot.Caches.GuildCache().Get(*s.GuildID)
}

func (s *Sticker) Update(stickerUpdate discord.StickerUpdate, opts ...rest.RequestOpt) (*Sticker, error) {
	if s.Type != discord.StickerTypeGuild {
		return nil, rest.NewError(nil, discord.ErrStickerTypeGuild)
	}

	sticker, err := s.Bot.RestServices.StickerService().UpdateSticker(*s.GuildID, s.ID, stickerUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return s.Bot.EntityBuilder.CreateSticker(*sticker, CacheStrategyNoWs), nil
}

func (s *Sticker) Delete(opts ...rest.RequestOpt) error {
	if s.Type != discord.StickerTypeGuild {
		return rest.NewError(nil, discord.ErrStickerTypeGuild)
	}
	return s.Bot.RestServices.StickerService().DeleteSticker(*s.GuildID, s.ID, opts...)
}
