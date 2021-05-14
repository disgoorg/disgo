package internal

import (
	"github.com/DisgoOrg/disgo/api"
)

func newEntityBuilderImpl(disgo api.Disgo) api.EntityBuilder {
	return &EntityBuilderImpl{disgo: disgo}
}

// EntityBuilderImpl is used for creating structs used by Disgo
type EntityBuilderImpl struct {
	disgo api.Disgo
}

// Disgo returns the api.Disgo client
func (b EntityBuilderImpl) Disgo() api.Disgo {
	return b.disgo
}

// CreateInteraction returns a new api.Interaction entity
func (b EntityBuilderImpl) CreateInteraction(interaction *api.Interaction, updateCache api.CacheStrategy) *api.Interaction {
	if interaction.Member != nil {
		interaction.Member = b.CreateMember(*interaction.GuildID, interaction.Member, api.CacheStrategyYes)
	}
	if interaction.User != nil {
		interaction.User = b.CreateUser(interaction.User, updateCache)
	}

	if interaction.Data != nil && interaction.Data.Resolved != nil {
		resolved := interaction.Data.Resolved
		if resolved.Users != nil {
			for _, user := range resolved.Users {
				user = b.CreateUser(user, updateCache)
			}
		}
		if resolved.Members != nil {
			for id, member := range resolved.Members {
				member.User = resolved.Users[id]
				member = b.CreateMember(*interaction.GuildID, member, updateCache)
			}
		}
		if resolved.Roles != nil {
			for _, role := range resolved.Roles {
				role = b.CreateRole(*interaction.GuildID, role, updateCache)
			}
		}
		// TODO how do we cache partial channels?
		/*if resolved.Channels != nil {
			for _, channel := range resolved.Channels {
				channel.Disgo = disgo
				disgo.Cache().CacheChannel(channel)
			}
		}*/
	}
	return interaction
}

// CreateGlobalCommand returns a new api.Command entity
func (b EntityBuilderImpl) CreateGlobalCommand(command *api.Command, updateCache api.CacheStrategy) *api.Command {
	command.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGlobalCommand(command)
	}
	return command
}

// CreateUser returns a new api.User entity
func (b EntityBuilderImpl) CreateUser(user *api.User, updateCache api.CacheStrategy) *api.User {
	user.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheUser(user)
	}
	return user
}

// CreateMessage returns a new api.Message entity
func (b EntityBuilderImpl) CreateMessage(message *api.Message, updateCache api.CacheStrategy) *api.Message {
	message.Disgo = b.Disgo()
	if message.Member != nil {
		message.Member = b.CreateMember(*message.GuildID, message.Member, updateCache)
	}
	if message.Author != nil {
		message.Author = b.CreateUser(message.Author, updateCache)
	}
	// TODO: should we cache mentioned users, members, etc?
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheMessage(message)
	}
	return message
}

// CreateGuild returns a new api.Guild entity
func (b EntityBuilderImpl) CreateGuild(guild *api.Guild, updateCache api.CacheStrategy) *api.Guild {
	guild.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGuild(guild)
	}
	return guild
}

// CreateMember returns a new api.Member entity
func (b EntityBuilderImpl) CreateMember(guildID api.Snowflake, member *api.Member, updateCache api.CacheStrategy) *api.Member {
	member.Disgo = b.Disgo()
	member.GuildID = guildID
	member.User = b.CreateUser(member.User, updateCache)
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheMember(member)
	}
	return member
}

// CreateThreadMember returns a new api.ThreadMember entity
func (b EntityBuilderImpl) CreateThreadMember(guildID api.Snowflake, member *api.ThreadMember, updateCache api.CacheStrategy) *api.ThreadMember {
	member.Disgo = b.Disgo()
	member.GuildID = guildID

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheThreadMember(member)
	}
	return member
}

// CreateVoiceState returns a new api.VoiceState entity
func (b EntityBuilderImpl) CreateVoiceState(guildID api.Snowflake, voiceState *api.VoiceState, updateCache api.CacheStrategy) *api.VoiceState {
	voiceState.Disgo = b.Disgo()
	voiceState.GuildID = guildID
	b.Disgo().Logger().Infof("voiceState: %+v", voiceState)

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheVoiceState(voiceState)
	}
	return voiceState
}

// CreateGuildCommand returns a new api.Command entity
func (b EntityBuilderImpl) CreateGuildCommand(guildID api.Snowflake, command *api.Command, updateCache api.CacheStrategy) *api.Command {
	command.Disgo = b.Disgo()
	command.GuildID = &guildID
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGuildCommand(command)
	}
	return command
}

// CreateGuildCommandPermissions returns a new api.GuildCommandPermissions entity
func (b EntityBuilderImpl) CreateGuildCommandPermissions(guildCommandPermissions *api.GuildCommandPermissions, updateCache api.CacheStrategy) *api.GuildCommandPermissions {
	guildCommandPermissions.Disgo = b.Disgo()
	if updateCache(b.Disgo()) && b.Disgo().Cache().CacheFlags().Has(api.CacheFlagCommandPermissions) {
		if cmd := b.Disgo().Cache().Command(guildCommandPermissions.ID); cmd != nil {
			cmd.GuildPermissions[guildCommandPermissions.GuildID] = guildCommandPermissions
		}
	}
	return guildCommandPermissions
}

// CreateRole returns a new api.Role entity
func (b EntityBuilderImpl) CreateRole(guildID api.Snowflake, role *api.Role, updateCache api.CacheStrategy) *api.Role {
	role.Disgo = b.Disgo()
	role.GuildID = guildID
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheRole(role)
	}
	return role
}

// CreateTextChannel returns a new api.TextChannel entity
func (b EntityBuilderImpl) CreateTextChannel(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.TextChannel {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheTextChannel(channel)
	}
	return channel
}

// CreateThread returns a new api.Thread entity
func (b EntityBuilderImpl) CreateThread(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.Thread {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheThread(channel)
	}
	return channel
}

// CreateVoiceChannel returns a new api.VoiceChannel entity
func (b EntityBuilderImpl) CreateVoiceChannel(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.VoiceChannel {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheVoiceChannel(channel)
	}
	return channel
}

// CreateStoreChannel returns a new api.StoreChannel entity
func (b EntityBuilderImpl) CreateStoreChannel(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.StoreChannel {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheStoreChannel(channel)
	}
	return channel
}

// CreateCategory returns a new api.Category entity
func (b EntityBuilderImpl) CreateCategory(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.Category {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheCategory(channel)
	}
	return channel
}

// CreateDMChannel returns a new api.DMChannel entity
func (b EntityBuilderImpl) CreateDMChannel(channel *api.ChannelImpl, updateCache api.CacheStrategy) api.DMChannel {
	channel.Disgo_ = b.Disgo()

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheDMChannel(channel)
	}
	return channel
}

// CreateEmote returns a new api.Emote entity
func (b EntityBuilderImpl) CreateEmote(guildID api.Snowflake, emote *api.Emote, updateCache api.CacheStrategy) *api.Emote {
	emote.Disgo = b.Disgo()
	emote.GuildID = guildID
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheEmote(emote)
	}
	return emote
}
