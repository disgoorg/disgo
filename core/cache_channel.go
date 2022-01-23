package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type (
	ChannelFindFunc func(channel Channel) bool

	ChannelCache interface {
		Get(channelID snowflake.Snowflake) Channel
		GetCopy(channelID snowflake.Snowflake) Channel
		Set(channel Channel) Channel
		Remove(channelID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]Channel
		All() []Channel

		FindFirst(channelFindFunc ChannelFindFunc) Channel
		FindAll(channelFindFunc ChannelFindFunc) []Channel

		ForAll(channelFunc func(channel Channel))
	}

	channelCacheImpl struct {
		channels   map[snowflake.Snowflake]Channel
		cacheFlags CacheFlags
	}
)

func NewChannelCache(cacheFlags CacheFlags) ChannelCache {
	return &channelCacheImpl{channels: map[snowflake.Snowflake]Channel{}, cacheFlags: cacheFlags}
}

func (c *channelCacheImpl) Get(channelID snowflake.Snowflake) Channel {
	return c.channels[channelID]
}

func (c *channelCacheImpl) GetCopy(channelID snowflake.Snowflake) Channel {
	if channel := c.Get(channelID); channel != nil {
		ch := &channel
		return *ch
	}
	return nil
}

func (c *channelCacheImpl) Set(channel Channel) Channel {
	if c.cacheFlags.Missing(getCacheFlagForChannelType(channel.Type())) {
		return channel
	}
	ch, ok := c.channels[channel.ID()]
	if ok {
		ch = ch.set(channel)
		return ch
	}
	c.channels[channel.ID()] = channel
	return channel
}

func (c *channelCacheImpl) Remove(id snowflake.Snowflake) {
	delete(c.channels, id)
}

func (c *channelCacheImpl) Cache() map[snowflake.Snowflake]Channel {
	return c.channels
}

func (c *channelCacheImpl) All() []Channel {
	channels := make([]Channel, len(c.channels))
	i := 0
	for _, channel := range c.channels {
		channels[i] = channel
		i++
	}
	return channels
}

func (c *channelCacheImpl) FindFirst(channelFindFunc ChannelFindFunc) Channel {
	for _, channel := range c.channels {
		if channelFindFunc(channel) {
			return channel
		}
	}
	return nil
}

func (c *channelCacheImpl) FindAll(channelFindFunc ChannelFindFunc) []Channel {
	var channels []Channel
	for _, channel := range c.channels {
		if channelFindFunc(channel) {
			channels = append(channels, channel)
		}
	}
	return channels
}

func (c *channelCacheImpl) ForAll(channelFunc func(channel Channel)) {
	for _, channel := range c.channels {
		channelFunc(channel)
	}
}

func getCacheFlagForChannelType(channelType discord.ChannelType) CacheFlags {
	switch channelType {
	case discord.ChannelTypeGuildText:
		return CacheFlagGuildTextChannels
	case discord.ChannelTypeDM:
		return CacheFlagDMChannels
	case discord.ChannelTypeGuildVoice:
		return CacheFlagGuildVoiceChannels
	case discord.ChannelTypeGroupDM:
		return CacheFlagGroupDMChannels
	case discord.ChannelTypeGuildCategory:
		return CacheFlagGuildCategories
	case discord.ChannelTypeGuildNews:
		return CacheFlagGuildNewsChannels
	case discord.ChannelTypeGuildStore:
		return CacheFlagGuildStoreChannels
	case discord.ChannelTypeGuildNewsThread:
		return CacheFlagGuildNewsThreads
	case discord.ChannelTypeGuildPublicThread:
		return CacheFlagGuildPublicThreads
	case discord.ChannelTypeGuildPrivateThread:
		return CacheFlagGuildPrivateThreads
	case discord.ChannelTypeGuildStageVoice:
		return CacheFlagGuildStageVoiceChannels
	default:
		return CacheFlagsNone
	}
}
