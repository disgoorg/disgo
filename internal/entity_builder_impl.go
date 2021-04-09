package internal

import (
	"github.com/DisgoOrg/disgo/api"
)

func newEntityBuilderImpl(disgo api.Disgo) api.EntityBuilder {
	return &EntityBuilderImpl{disgo: disgo}
}

type EntityBuilderImpl struct {
	disgo api.Disgo
}

func (b EntityBuilderImpl) Disgo() api.Disgo {
	return b.disgo
}

func (b EntityBuilderImpl) CreateTextChannel(channel *api.Channel) *api.TextChannel {
	channel.Disgo = b.Disgo()
	textChannel := &api.TextChannel{
		MessageChannel: api.MessageChannel{
			Channel: *channel,
		},
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	b.Disgo().Cache().CacheTextChannel(textChannel)
	return textChannel
}

func (b EntityBuilderImpl) CreateVoiceChannel(channel *api.Channel) *api.VoiceChannel {
	channel.Disgo = b.Disgo()
	voiceChannel := &api.VoiceChannel{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	b.Disgo().Cache().CacheVoiceChannel(voiceChannel)
	return voiceChannel
}

func (b EntityBuilderImpl) CreateStoreChannel(channel *api.Channel) *api.StoreChannel {
	channel.Disgo = b.Disgo()
	storeChannel := &api.StoreChannel{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	b.Disgo().Cache().CacheStoreChannel(storeChannel)
	return storeChannel
}

func (b EntityBuilderImpl) CreateCategory(channel *api.Channel) *api.Category {
	channel.Disgo = b.Disgo()
	category := &api.Category{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	b.Disgo().Cache().CacheCategory(category)
	return category
}

func (b EntityBuilderImpl) CreateDMChannel(channel *api.Channel) *api.DMChannel {
	channel.Disgo = b.Disgo()
	dmChannel := &api.DMChannel{
		MessageChannel: api.MessageChannel{
			Channel: *channel,
		},
	}
	b.Disgo().Cache().CacheDMChannel(dmChannel)
	return dmChannel
}
