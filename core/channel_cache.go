package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type (
	ChannelFindFunc func(channel *Channel) bool

	ChannelCache interface {
		Get(channelID discord.Snowflake) *Channel
		GetCopy(channelID discord.Snowflake) *Channel
		Set(channel *Channel) *Channel
		Remove(channelID discord.Snowflake)

		Cache() map[discord.Snowflake]*Channel
		All() []*Channel

		FindFirst(channelFindFunc ChannelFindFunc) *Channel
		FindAll(channelFindFunc ChannelFindFunc) []*Channel
	}

	channelCacheImpl struct {
		channels   map[discord.Snowflake]*Channel
		cacheFlags CacheFlags
	}
)

func NewChannelCache(cacheFlags CacheFlags) ChannelCache {
	return &channelCacheImpl{channels: map[discord.Snowflake]*Channel{}, cacheFlags: cacheFlags}
}

func (c *channelCacheImpl) Get(channelID discord.Snowflake) *Channel {
	return c.channels[channelID]
}

func (c *channelCacheImpl) GetCopy(channelID discord.Snowflake) *Channel {
	return &*c.Get(channelID)
}

func (c *channelCacheImpl) Set(channel *Channel) *Channel {
	// TODO: check cache flags for specific channel
	ch, ok := c.channels[channel.ID]
	if ok {
		*ch = *channel
		return ch
	}
	c.channels[channel.ID] = channel
	return channel
}

func (c *channelCacheImpl) Remove(id discord.Snowflake) {
	delete(c.channels, id)
}

func (c *channelCacheImpl) Cache() map[discord.Snowflake]*Channel {
	return c.channels
}

func (c *channelCacheImpl) All() []*Channel {
	channels := make([]*Channel, len(c.channels))
	i := 0
	for _, channel := range c.channels {
		channels[i] = channel
		i++
	}
	return channels
}

func (c *channelCacheImpl) FindFirst(channelFindFunc ChannelFindFunc) *Channel {
	for _, channel := range c.channels {
		if channelFindFunc(channel) {
			return channel
		}
	}
	return nil
}

func (c *channelCacheImpl) FindAll(channelFindFunc ChannelFindFunc) []*Channel {
	var channels []*Channel
	for _, channel := range c.channels {
		if channelFindFunc(channel) {
			channels = append(channels, channel)
		}
	}
	return channels
}
