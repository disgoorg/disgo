package cache

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func NewChannelCache(flags Flags, neededFlags Flags, policy Policy[discord.Channel]) ChannelCache {
	return &ChannelCacheImpl{
		Cache: NewCache[discord.Channel](flags, neededFlags, policy),
	}
}

type ChannelCache interface {
	Cache[discord.Channel]

	GuildChannels(guildID snowflake.ID) []discord.GuildChannel
	GuildThreadsInChannel(channelID snowflake.ID) []discord.GuildThread

	GetGuildChannel(channelID snowflake.ID) (discord.GuildChannel, bool)
	GetMessageChannel(channelID snowflake.ID) (discord.MessageChannel, bool)
	GetGuildMessageChannel(channelID snowflake.ID) (discord.GuildMessageChannel, bool)
	GetGuildThread(channelID snowflake.ID) (discord.GuildThread, bool)
	GetGuildAudioChannel(channelID snowflake.ID) (discord.GuildAudioChannel, bool)

	GetGuildTextChannel(channelID snowflake.ID) (discord.GuildTextChannel, bool)
	GetDMChannel(channelID snowflake.ID) (discord.DMChannel, bool)
	GetGuildVoiceChannel(channelID snowflake.ID) (discord.GuildVoiceChannel, bool)
	GetGuildCategoryChannel(channelID snowflake.ID) (discord.GuildCategoryChannel, bool)
	GetGuildNewsChannel(channelID snowflake.ID) (discord.GuildNewsChannel, bool)
	GetGuildNewsThread(channelID snowflake.ID) (discord.GuildThread, bool)
	GetGuildPublicThread(channelID snowflake.ID) (discord.GuildThread, bool)
	GetGuildPrivateThread(channelID snowflake.ID) (discord.GuildThread, bool)
	GetGuildStageVoiceChannel(channelID snowflake.ID) (discord.GuildStageVoiceChannel, bool)
}

type ChannelCacheImpl struct {
	Cache[discord.Channel]
}

func (c *ChannelCacheImpl) GuildChannels(guildID snowflake.ID) []discord.GuildChannel {
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

func (c *ChannelCacheImpl) GuildThreadsInChannel(channelID snowflake.ID) []discord.GuildThread {
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

func (c *ChannelCacheImpl) GetGuildChannel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetMessageChannel(channelID snowflake.ID) (discord.MessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.MessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildMessageChannel(channelID snowflake.ID) (discord.GuildMessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if chM, ok := ch.(discord.GuildMessageChannel); ok {
			return chM, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildThread); ok {
			return cCh, true
		}
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildAudioChannel(channelID snowflake.ID) (discord.GuildAudioChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCacheImpl) GetGuildTextChannel(channelID snowflake.ID) (discord.GuildTextChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildTextChannel{}, false
}

func (c *ChannelCacheImpl) GetDMChannel(channelID snowflake.ID) (discord.DMChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.DMChannel); ok {
			return cCh, true
		}
	}
	return discord.DMChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildVoiceChannel(channelID snowflake.ID) (discord.GuildVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildVoiceChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildCategoryChannel(channelID snowflake.ID) (discord.GuildCategoryChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildCategoryChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildNewsChannel(channelID snowflake.ID) (discord.GuildNewsChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildNewsChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsChannel{}, false
}

func (c *ChannelCacheImpl) GetGuildNewsThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildNewsThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildPublicThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPublicThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildPrivateThread(channelID snowflake.ID) (discord.GuildThread, bool) {
	if ch, ok := c.GetGuildThread(channelID); ok && ch.Type() == discord.ChannelTypeGuildPrivateThread {
		return ch, true
	}
	return discord.GuildThread{}, false
}

func (c *ChannelCacheImpl) GetGuildStageVoiceChannel(channelID snowflake.ID) (discord.GuildStageVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildStageVoiceChannel{}, false
}
