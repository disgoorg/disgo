package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	StickerFindFunc func(sticker *Sticker) bool

	StickerCache interface {
		Get(guildID discord.Snowflake, stickerID discord.Snowflake) *Sticker
		GetCopy(guildID discord.Snowflake, stickerID discord.Snowflake) *Sticker
		Set(sticker *Sticker) *Sticker
		Remove(guildID discord.Snowflake, stickerID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Sticker
		All() map[discord.Snowflake][]*Sticker

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Sticker
		GuildAll(guildID discord.Snowflake) []*Sticker

		FindFirst(stickerFindFunc StickerFindFunc) *Sticker
		FindAll(stickerFindFunc StickerFindFunc) []*Sticker
	}

	stickerCacheImpl struct {
		cacheFlags CacheFlags
		stickers   map[discord.Snowflake]map[discord.Snowflake]*Sticker
	}
)

func NewStickerCache(cacheFlags CacheFlags) StickerCache {
	return &stickerCacheImpl{
		cacheFlags: cacheFlags,
		stickers:   map[discord.Snowflake]map[discord.Snowflake]*Sticker{},
	}
}

func (c *stickerCacheImpl) Get(guildID discord.Snowflake, stickerID discord.Snowflake) *Sticker {
	if _, ok := c.stickers[guildID]; !ok {
		return nil
	}
	return c.stickers[guildID][stickerID]
}

func (c *stickerCacheImpl) GetCopy(guildID discord.Snowflake, stickerID discord.Snowflake) *Sticker {
	if sticker := c.Get(guildID, stickerID); sticker != nil {
		st := *sticker
		return &st
	}
	return nil
}

func (c *stickerCacheImpl) Set(sticker *Sticker) *Sticker {
	if sticker.GuildID == nil {
		return sticker
	}
	if c.cacheFlags.Missing(CacheFlagStickers) {
		return sticker
	}
	if _, ok := c.stickers[*sticker.GuildID]; !ok {
		c.stickers[*sticker.GuildID] = map[discord.Snowflake]*Sticker{}
	}
	st, ok := c.stickers[*sticker.GuildID][sticker.ID]
	if ok {
		*st = *sticker
		return st
	}
	c.stickers[*sticker.GuildID][sticker.ID] = sticker

	return sticker
}

func (c *stickerCacheImpl) Remove(guildID discord.Snowflake, stickerID discord.Snowflake) {
	if _, ok := c.stickers[guildID]; !ok {
		return
	}
	delete(c.stickers[guildID], stickerID)
}

func (c *stickerCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Sticker {
	return c.stickers
}

func (c *stickerCacheImpl) All() map[discord.Snowflake][]*Sticker {
	stickers := make(map[discord.Snowflake][]*Sticker, len(c.stickers))
	for guildID, guildStickers := range c.stickers {
		stickers[guildID] = make([]*Sticker, len(guildStickers))
		i := 0
		for _, guildSticker := range guildStickers {
			stickers[guildID] = append(stickers[guildID], guildSticker)
		}
		i++
	}
	return stickers
}

func (c *stickerCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Sticker {
	if _, ok := c.stickers[guildID]; !ok {
		return nil
	}
	return c.stickers[guildID]
}

func (c *stickerCacheImpl) GuildAll(guildID discord.Snowflake) []*Sticker {
	if _, ok := c.stickers[guildID]; !ok {
		return nil
	}
	stickers := make([]*Sticker, len(c.stickers[guildID]))
	i := 0
	for _, sticker := range c.stickers[guildID] {
		stickers = append(stickers, sticker)
		i++
	}
	return stickers
}

func (c *stickerCacheImpl) FindFirst(stickerFindFunc StickerFindFunc) *Sticker {
	for _, guildStickers := range c.stickers {
		for _, sticker := range guildStickers {
			if stickerFindFunc(sticker) {
				return sticker
			}
		}
	}
	return nil
}

func (c *stickerCacheImpl) FindAll(stickerFindFunc StickerFindFunc) []*Sticker {
	var stickers []*Sticker
	for _, guildStickers := range c.stickers {
		for _, sticker := range guildStickers {
			if stickerFindFunc(sticker) {
				stickers = append(stickers, sticker)
			}
		}
	}
	return stickers
}
