package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Sticker struct {
	discord.Sticker
	Bot  Bot
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

// Guild returns the Guild this Sticker was created for or nil if this isn't a Guild-specific Sticker.
// This will only check cached guilds!
func (s *Sticker) Guild() (guild Guild, ok bool) {
	if s.Type != discord.StickerTypeGuild {
		return
	}
	return s.Bot.Caches().Guilds().Get(*s.GuildID)
}

// Update updates this Sticker with the properties provided in discord.StickerUpdate
func (s *Sticker) Update(stickerUpdate discord.StickerUpdate, opts ...rest.RequestOpt) (*Sticker, error) {
	if s.Type != discord.StickerTypeGuild {
		return nil, discord.ErrStickerTypeGuild
	}

	sticker, err := s.Bot.RestServices().StickerService().UpdateSticker(*s.GuildID, s.ID, stickerUpdate, opts...)
	if err != nil {
		return nil, err
	}
	return s.Bot.EntityBuilder().CreateSticker(*sticker, CacheStrategyNoWs), nil
}

// Delete deletes this Sticker
func (s *Sticker) Delete(opts ...rest.RequestOpt) error {
	if s.Type != discord.StickerTypeGuild {
		return discord.ErrStickerTypeGuild
	}
	return s.Bot.RestServices().StickerService().DeleteSticker(*s.GuildID, s.ID, opts...)
}
