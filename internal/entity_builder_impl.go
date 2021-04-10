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

func (b EntityBuilderImpl) CreateGlobalCommand(command *api.Command, updateCache bool) *api.Command {
	command.Disgo = b.Disgo()
	if updateCache {
		return b.Disgo().Cache().CacheGlobalCommand(command)
	}
	return command
}

func (b EntityBuilderImpl) CreateUser(user *api.User, updateCache bool) *api.User {
	user.Disgo = b.Disgo()
	if updateCache {
		return b.Disgo().Cache().CacheUser(user)
	}
	return user
}

func (b EntityBuilderImpl) CreateMessage(message *api.Message, updateCache bool) *api.Message {
	message.Disgo = b.Disgo()
	if message.Member != nil {
		message.Member = b.CreateMember(*message.GuildID, message.Member, true)
	}
	if message.Author != nil {
		
	}
	if updateCache {
		return b.Disgo().Cache().CacheMessage(message)
	}
	return message
}

func (b EntityBuilderImpl) CreateGuild(guild *api.Guild, updateCache bool) *api.Guild {
	guild.Disgo = b.Disgo()
	if updateCache {
		return b.Disgo().Cache().CacheGuild(guild)
	}
	return guild
}

func (b EntityBuilderImpl) CreateMember(guildID api.Snowflake, member *api.Member, updateCache bool) *api.Member {
	member.Disgo = b.Disgo()
	member.GuildID = guildID
	member.User = b.CreateUser(member.User, updateCache)
	if updateCache {
		return b.Disgo().Cache().CacheMember(member)
	}
	return member
}

func (b EntityBuilderImpl) CreateVoiceState(voiceState *api.VoiceState, updateCache bool) *api.VoiceState {
	voiceState.Disgo = b.Disgo()
	if updateCache {
		return b.Disgo().Cache().CacheVoiceState(voiceState)
	}
	return voiceState
}

func (b EntityBuilderImpl) CreateGuildCommand(guildID api.Snowflake, command *api.Command, updateCache bool) *api.Command {
	command.Disgo = b.Disgo()
	command.GuildID = &guildID
	if updateCache {
		return b.Disgo().Cache().CacheGuildCommand(command)
	}
	return command
}

func (b EntityBuilderImpl) CreateRole(guildID api.Snowflake, role *api.Role, updateCache bool) *api.Role {
	role.Disgo = b.Disgo()
	role.GuildID = guildID
	if updateCache {
		return b.Disgo().Cache().CacheRole(role)
	}
	return role
}

func (b EntityBuilderImpl) CreateTextChannel(channel *api.Channel, updateCache bool) *api.TextChannel {
	channel.Disgo = b.Disgo()
	textChannel := &api.TextChannel{
		MessageChannel: api.MessageChannel{
			Channel: *channel,
		},
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	if updateCache {
		return b.Disgo().Cache().CacheTextChannel(textChannel)
	}
	return textChannel
}

func (b EntityBuilderImpl) CreateVoiceChannel(channel *api.Channel, updateCache bool) *api.VoiceChannel {
	channel.Disgo = b.Disgo()
	voiceChannel := &api.VoiceChannel{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	if updateCache {
		return b.Disgo().Cache().CacheVoiceChannel(voiceChannel)
	}
	return voiceChannel
}

func (b EntityBuilderImpl) CreateStoreChannel(channel *api.Channel, updateCache bool) *api.StoreChannel {
	channel.Disgo = b.Disgo()
	storeChannel := &api.StoreChannel{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	if updateCache {
		return b.Disgo().Cache().CacheStoreChannel(storeChannel)
	}
	return storeChannel
}

func (b EntityBuilderImpl) CreateCategory(channel *api.Channel, updateCache bool) *api.Category {
	channel.Disgo = b.Disgo()
	category := &api.Category{
		GuildChannel: api.GuildChannel{
			Channel: *channel,
		},
	}
	if updateCache {
		return b.Disgo().Cache().CacheCategory(category)
	}
	return category
}

func (b EntityBuilderImpl) CreateDMChannel(channel *api.Channel, updateCache bool) *api.DMChannel {
	channel.Disgo = b.Disgo()
	dmChannel := &api.DMChannel{
		MessageChannel: api.MessageChannel{
			Channel: *channel,
		},
	}
	if updateCache {
		return b.Disgo().Cache().CacheDMChannel(dmChannel)
	}
	return dmChannel
}

func (b EntityBuilderImpl) CreateEmote(guildId api.Snowflake, emote *api.Emote, updateCache bool) *api.Emote {
	emote.Disgo = b.Disgo()
	emote.GuildID = guildId
	if updateCache {
		return b.Disgo().Cache().CacheEmote(emote)
	}
	return emote
}
