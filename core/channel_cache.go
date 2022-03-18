package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type ChannelCache struct {
	Cache[discord.Channel]
}

func (c *ChannelCache) GetGuildChannel(channelID snowflake.Snowflake) (discord.GuildChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetMessageChannel(channelID snowflake.Snowflake) (discord.MessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.MessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetBaseGuildMessageChannel(channelID snowflake.Snowflake) (discord.BaseGuildMessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.BaseGuildMessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetGuildMessageChannel(channelID snowflake.Snowflake) (discord.GuildMessageChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildMessageChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetGuildThread(channelID snowflake.Snowflake) (discord.GuildThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildThread); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetGuildAudioChannel(channelID snowflake.Snowflake) (discord.GuildAudioChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildAudioChannel); ok {
			return cCh, true
		}
	}
	return nil, false
}

func (c *ChannelCache) GetGuildTextChannel(channelID snowflake.Snowflake) (discord.GuildTextChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildTextChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildTextChannel{}, false
}

func (c *ChannelCache) GetDMChannel(channelID snowflake.Snowflake) (discord.DMChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.DMChannel); ok {
			return cCh, true
		}
	}
	return discord.DMChannel{}, false
}

func (c *ChannelCache) GetGuildVoiceChannel(channelID snowflake.Snowflake) (discord.GuildVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildVoiceChannel{}, false
}

func (c *ChannelCache) GetGuildCategoryChannel(channelID snowflake.Snowflake) (discord.GuildCategoryChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildCategoryChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildCategoryChannel{}, false
}

func (c *ChannelCache) GetGuildNewsChannel(channelID snowflake.Snowflake) (discord.GuildNewsChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildNewsChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsChannel{}, false
}

func (c *ChannelCache) GetGuildNewsThread(channelID snowflake.Snowflake) (discord.GuildNewsThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildNewsThread); ok {
			return cCh, true
		}
	}
	return discord.GuildNewsThread{}, false
}

func (c *ChannelCache) GetGuildPublicThread(channelID snowflake.Snowflake) (discord.GuildPublicThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildPublicThread); ok {
			return cCh, true
		}
	}
	return discord.GuildPublicThread{}, false
}

func (c *ChannelCache) GetGuildPrivateThread(channelID snowflake.Snowflake) (discord.GuildPrivateThread, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildPrivateThread); ok {
			return cCh, true
		}
	}
	return discord.GuildPrivateThread{}, false
}

func (c *ChannelCache) GetGuildStageVoiceChannel(channelID snowflake.Snowflake) (discord.GuildStageVoiceChannel, bool) {
	if ch, ok := c.Get(channelID); ok {
		if cCh, ok := ch.(discord.GuildStageVoiceChannel); ok {
			return cCh, true
		}
	}
	return discord.GuildStageVoiceChannel{}, false
}

func (c *ChannelCache) GuildChannels(guildID snowflake.Snowflake) []discord.GuildChannel {
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
