package cache

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

func NewChannelCache(flags Flags, neededFlags Flags, policy Policy[discord.Channel]) ChannelCache {
	return &ChannelCacheImpl{
		Cache:         NewCache[discord.Channel](flags, neededFlags, policy),
		guildChannels: map[snowflake.Snowflake]snowflake.Snowflake{},
		dmChannels:    map[snowflake.Snowflake]snowflake.Snowflake{},
	}
}

type ChannelCache interface {
	Cache[discord.Channel]

	GuildChannels(guildID snowflake.Snowflake) []discord.GuildChannel
	GuildThreadsInChannel(channelID snowflake.Snowflake) []discord.GuildThread

	GetGuildChannel(channelID snowflake.Snowflake) (discord.GuildChannel, bool)
	GetMessageChannel(channelID snowflake.Snowflake) (discord.MessageChannel, bool)
	GetGuildMessageChannel(channelID snowflake.Snowflake) (discord.GuildMessageChannel, bool)
	GetGuildThread(channelID snowflake.Snowflake) (discord.GuildThread, bool)
	GetGuildAudioChannel(channelID snowflake.Snowflake) (discord.GuildAudioChannel, bool)

	GetGuildTextChannel(channelID snowflake.Snowflake) (discord.GuildTextChannel, bool)
	GetDMChannel(channelID snowflake.Snowflake) (discord.DMChannel, bool)
	GetGuildVoiceChannel(channelID snowflake.Snowflake) (discord.GuildVoiceChannel, bool)
	GetGuildCategoryChannel(channelID snowflake.Snowflake) (discord.GuildCategoryChannel, bool)
	GetGuildNewsChannel(channelID snowflake.Snowflake) (discord.GuildNewsChannel, bool)
	GetGuildNewsThread(channelID snowflake.Snowflake) (discord.GuildThread, bool)
	GetGuildPublicThread(channelID snowflake.Snowflake) (discord.GuildThread, bool)
	GetGuildPrivateThread(channelID snowflake.Snowflake) (discord.GuildThread, bool)
	GetGuildStageVoiceChannel(channelID snowflake.Snowflake) (discord.GuildStageVoiceChannel, bool)
}

type ChannelCacheImpl struct {
	Cache[discord.Channel]
	guildChannels map[snowflake.Snowflake]snowflake.Snowflake
	dmChannels    map[snowflake.Snowflake]snowflake.Snowflake
}

func (c *ChannelCacheImpl) GuildChannels(guildID snowflake.Snowflake) []discord.GuildChannel {
	channels := c.FindAll(func(channel discord.Channel) bool {
		if ch, ok := channel.(discord.GuildChannel); ok {
			return ch.GuildID() == guildID
		}
		return false
	})
	guildChannels := make([]discord.GuildChannel, len(channels))
	for i, channel := range channels {
		guildChannels[i] = channel.(discord.GuildChannel)
	}
	return guildChannels
}

func (c *ChannelCacheImpl) GuildThreadsInChannel(channelID snowflake.Snowflake) []discord.GuildThread {
	channels := c.FindAll(func(channel discord.Channel) bool {
		if thread, ok := channel.(discord.GuildThread); ok {
			return *thread.ParentID() == channelID
		}
		return false
	})
	threads := make([]discord.GuildThread, len(channels))
	for i, channel := range channels {
		threads[i] = channel.(discord.GuildThread)
	}
	return threads
}

func (c *ChannelCacheImpl) GetGuildChannel(channelID snowflake.Snowflake) (discord.GuildChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetMessageChannel(channelID snowflake.Snowflake) (discord.MessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.MessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildMessageChannel(channelID snowflake.Snowflake) (discord.GuildMessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildMessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildThread(channelID snowflake.Snowflake) (discord.GuildThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildThread); ok {
			return cCh, true
		}
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildAudioChannel(channelID snowflake.Snowflake) (discord.GuildAudioChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildTextChannel(channelID snowflake.Snowflake) (discord.GuildTextChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildTextChannel{}, false
}

func (c *ChannelCacheImpl) GetDMChannel(channelID snowflake.Snowflake) (discord.DMChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.DMChannel); ok {
			return cCh, true
		}
	}
	return discord.DMChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildVoiceChannel(channelID snowflake.Snowflake) (discord.GuildVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildVoiceChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildCategoryChannel(channelID snowflake.Snowflake) (discord.GuildCategoryChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildCategoryChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildNewsChannel(channelID snowflake.Snowflake) (discord.GuildNewsChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildNewsChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildNewsThread(channelID snowflake.Snowflake) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildNewsThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildPublicThread(channelID snowflake.Snowflake) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPublicThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildPrivateThread(channelID snowflake.Snowflake) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPrivateThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildStageVoiceChannel(channelID snowflake.Snowflake) (discord.GuildStageVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildStageVoiceChannel{}, false
}
