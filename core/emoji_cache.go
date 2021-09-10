package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	EmojiFindFunc func(emoji *Emoji) bool

	EmojiCache interface {
		Get(guildID discord.Snowflake, emojiID discord.Snowflake) *Emoji
		GetCopy(guildID discord.Snowflake, emojiID discord.Snowflake) *Emoji
		Set(emoji *Emoji) *Emoji
		Remove(guildID discord.Snowflake, emojiID discord.Snowflake)

		Cache() map[discord.Snowflake]map[discord.Snowflake]*Emoji
		All() map[discord.Snowflake][]*Emoji

		GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Emoji
		GuildAll(guildID discord.Snowflake) []*Emoji

		FindFirst(emojiFindFunc EmojiFindFunc) *Emoji
		FindAll(emojiFindFunc EmojiFindFunc) []*Emoji
	}

	emojiCacheImpl struct {
		cacheFlags CacheFlags
		emojis     map[discord.Snowflake]map[discord.Snowflake]*Emoji
	}
)

func NewEmojiCache(cacheFlags CacheFlags) EmojiCache {
	return &emojiCacheImpl{
		cacheFlags: cacheFlags,
		emojis:     map[discord.Snowflake]map[discord.Snowflake]*Emoji{},
	}
}

func (c *emojiCacheImpl) Get(guildID discord.Snowflake, emojiID discord.Snowflake) *Emoji {
	if _, ok := c.emojis[guildID]; !ok {
		return nil
	}
	return c.emojis[guildID][emojiID]
}

func (c *emojiCacheImpl) GetCopy(guildID discord.Snowflake, emojiID discord.Snowflake) *Emoji {
	if emoji := c.Get(guildID, emojiID); emoji != nil {
		em := *emoji
		return &em
	}
	return nil
}

func (c *emojiCacheImpl) Set(emoji *Emoji) *Emoji {
	if c.cacheFlags.Missing(CacheFlagEmojis) {
		return emoji
	}
	if _, ok := c.emojis[emoji.GuildID]; !ok {
		c.emojis[emoji.GuildID] = map[discord.Snowflake]*Emoji{}
	}
	rol, ok := c.emojis[emoji.GuildID][emoji.ID]
	if ok {
		*rol = *emoji
		return rol
	}
	c.emojis[emoji.GuildID][emoji.ID] = emoji

	return emoji
}

func (c *emojiCacheImpl) Remove(guildID discord.Snowflake, emojiID discord.Snowflake) {
	if _, ok := c.emojis[guildID]; !ok {
		return
	}
	delete(c.emojis[guildID], emojiID)
}

func (c *emojiCacheImpl) Cache() map[discord.Snowflake]map[discord.Snowflake]*Emoji {
	return c.emojis
}

func (c *emojiCacheImpl) All() map[discord.Snowflake][]*Emoji {
	emojis := make(map[discord.Snowflake][]*Emoji, len(c.emojis))
	for guildID, guildEmojis := range c.emojis {
		emojis[guildID] = make([]*Emoji, len(guildEmojis))
		i := 0
		for _, guildEmoji := range guildEmojis {
			emojis[guildID] = append(emojis[guildID], guildEmoji)
		}
		i++
	}
	return emojis
}

func (c *emojiCacheImpl) GuildCache(guildID discord.Snowflake) map[discord.Snowflake]*Emoji {
	if _, ok := c.emojis[guildID]; !ok {
		return nil
	}
	return c.emojis[guildID]
}

func (c *emojiCacheImpl) GuildAll(guildID discord.Snowflake) []*Emoji {
	if _, ok := c.emojis[guildID]; !ok {
		return nil
	}
	emojis := make([]*Emoji, len(c.emojis[guildID]))
	i := 0
	for _, emoji := range c.emojis[guildID] {
		emojis = append(emojis, emoji)
		i++
	}
	return emojis
}

func (c *emojiCacheImpl) FindFirst(emojiFindFunc EmojiFindFunc) *Emoji {
	for _, guildEmojis := range c.emojis {
		for _, emoji := range guildEmojis {
			if emojiFindFunc(emoji) {
				return emoji
			}
		}
	}
	return nil
}

func (c *emojiCacheImpl) FindAll(emojiFindFunc EmojiFindFunc) []*Emoji {
	var emojis []*Emoji
	for _, guildEmojis := range c.emojis {
		for _, emoji := range guildEmojis {
			if emojiFindFunc(emoji) {
				emojis = append(emojis, emoji)
			}
		}
	}
	return emojis
}
