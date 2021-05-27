package internal

import (
	"encoding/json"

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
func (b *EntityBuilderImpl) Disgo() api.Disgo {
	return b.disgo
}

func (b EntityBuilderImpl) createInteraction(fullInteraction *api.FullInteraction, updateCache api.CacheStrategy) *api.Interaction {
	interaction := &api.Interaction{
		Disgo:     b.disgo,
		ID:        fullInteraction.ID,
		Type:      fullInteraction.Type,
		GuildID:   fullInteraction.GuildID,
		ChannelID: fullInteraction.ChannelID,
		Token:     fullInteraction.Token,
		Version:   fullInteraction.Version,
	}

	if fullInteraction.Member != nil {
		interaction.Member = b.CreateMember(*fullInteraction.GuildID, fullInteraction.Member, api.CacheStrategyYes)
	}
	if fullInteraction.User != nil {
		interaction.User = b.CreateUser(fullInteraction.User, updateCache)
	}
	return interaction
}

// CreateButtonInteraction creates a api.ButtonInteraction from the full interaction response
func (b *EntityBuilderImpl) CreateButtonInteraction(fullInteraction *api.FullInteraction, updateCache api.CacheStrategy) *api.ButtonInteraction {
	var data *api.ButtonInteractionData
	_ = json.Unmarshal(fullInteraction.Data, &data)

	return &api.ButtonInteraction{
		Interaction: b.createInteraction(fullInteraction, updateCache),
		Message:     b.CreateMessage(fullInteraction.FullMessage, updateCache),
		Data:        data,
	}
}

// CreateCommandInteraction creates a api.CommandInteraction from the full interaction response
func (b *EntityBuilderImpl) CreateCommandInteraction(fullInteraction *api.FullInteraction, updateCache api.CacheStrategy) *api.CommandInteraction {
	var data *api.CommandInteractionData
	_ = json.Unmarshal(fullInteraction.Data, &data)

	if data.Resolved != nil {
		resolved := data.Resolved
		if resolved.Users != nil {
			for _, user := range resolved.Users {
				user = b.CreateUser(user, updateCache)
			}
		}
		if resolved.Members != nil {
			for id, member := range resolved.Members {
				member.User = resolved.Users[id]
				member = b.CreateMember(member.GuildID, member, updateCache)
			}
		}
		if resolved.Roles != nil {
			for _, role := range resolved.Roles {
				role = b.CreateRole(role.GuildID, role, updateCache)
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

	return &api.CommandInteraction{
		Interaction: b.createInteraction(fullInteraction, updateCache),
		Data:        data,
	}
}

// CreateGlobalCommand returns a new api.Command entity
func (b *EntityBuilderImpl) CreateGlobalCommand(command *api.Command, updateCache api.CacheStrategy) *api.Command {
	command.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGlobalCommand(command)
	}
	return command
}

// CreateUser returns a new api.User entity
func (b *EntityBuilderImpl) CreateUser(user *api.User, updateCache api.CacheStrategy) *api.User {
	user.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheUser(user)
	}
	return user
}

func (b *EntityBuilderImpl) createComponent(unmarshalComponent *api.UnmarshalComponent, updateCache api.CacheStrategy) api.Component {
	switch unmarshalComponent.ComponentType {
	case api.ComponentTypeActionRow:
		components := make([]api.Component, len(unmarshalComponent.Components))
		for i, unmarshalC := range unmarshalComponent.Components {
			components[i] = b.createComponent(unmarshalC, updateCache)
		}
		return &api.ActionRow{
			ComponentImpl: api.ComponentImpl{
				ComponentType: api.ComponentTypeActionRow,
			},
			Components: components,
		}

	case api.ComponentTypeButton:
		return &api.Button{
			ComponentImpl: api.ComponentImpl{
				ComponentType: api.ComponentTypeButton,
			},
			Style:    unmarshalComponent.Style,
			Label:    unmarshalComponent.Label,
			Emote:    b.CreateEmote("", unmarshalComponent.Emote, updateCache),
			CustomID: unmarshalComponent.CustomID,
			URL:      unmarshalComponent.URL,
			Disabled: unmarshalComponent.Disabled,
		}

	default:
		b.Disgo().Logger().Errorf("unexpected component type %d received", unmarshalComponent.ComponentType)
		return nil
	}
}

// CreateMessage returns a new api.Message entity
func (b *EntityBuilderImpl) CreateMessage(fullMessage *api.FullMessage, updateCache api.CacheStrategy) *api.Message {
	message := fullMessage.Message
	message.Disgo = b.Disgo()

	if message.Member != nil {
		message.Member.User = message.Author
		message.Member = b.CreateMember(*message.GuildID, message.Member, updateCache)
	}
	if message.Author != nil {
		message.Author = b.CreateUser(message.Author, updateCache)
	}

	if fullMessage.UnmarshalComponents != nil {
		for _, component := range fullMessage.UnmarshalComponents {
			message.Components = append(message.Components, b.createComponent(component, updateCache))
		}
	}

	// TODO: should we cache mentioned users, members, etc?
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheMessage(message)
	}
	return message
}

// CreateGuild returns a new api.Guild entity
func (b *EntityBuilderImpl) CreateGuild(guild *api.Guild, updateCache api.CacheStrategy) *api.Guild {
	guild.Disgo = b.Disgo()
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGuild(guild)
	}
	return guild
}

// CreateMember returns a new api.Member entity
func (b *EntityBuilderImpl) CreateMember(guildID api.Snowflake, member *api.Member, updateCache api.CacheStrategy) *api.Member {
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
func (b *EntityBuilderImpl) CreateVoiceState(guildID api.Snowflake, voiceState *api.VoiceState, updateCache api.CacheStrategy) *api.VoiceState {
	voiceState.Disgo = b.Disgo()
	voiceState.GuildID = guildID
	b.Disgo().Logger().Infof("voiceState: %+v", voiceState)

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheVoiceState(voiceState)
	}
	return voiceState
}

// CreateGuildCommand returns a new api.Command entity
func (b *EntityBuilderImpl) CreateGuildCommand(guildID api.Snowflake, command *api.Command, updateCache api.CacheStrategy) *api.Command {
	command.Disgo = b.Disgo()
	command.GuildID = &guildID
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheGuildCommand(command)
	}
	return command
}

// CreateGuildCommandPermissions returns a new api.GuildCommandPermissions entity
func (b *EntityBuilderImpl) CreateGuildCommandPermissions(guildCommandPermissions *api.GuildCommandPermissions, updateCache api.CacheStrategy) *api.GuildCommandPermissions {
	guildCommandPermissions.Disgo = b.Disgo()
	if updateCache(b.Disgo()) && b.Disgo().Cache().CacheFlags().Has(api.CacheFlagCommandPermissions) {
		if cmd := b.Disgo().Cache().Command(guildCommandPermissions.ID); cmd != nil {
			cmd.GuildPermissions[guildCommandPermissions.GuildID] = guildCommandPermissions
		}
	}
	return guildCommandPermissions
}

// CreateRole returns a new api.Role entity
func (b *EntityBuilderImpl) CreateRole(guildID api.Snowflake, role *api.Role, updateCache api.CacheStrategy) *api.Role {
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
func (b *EntityBuilderImpl) CreateEmote(guildID api.Snowflake, emote *api.Emote, updateCache api.CacheStrategy) *api.Emote {
	emote.Disgo = b.Disgo()
	emote.GuildID = guildID
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().CacheEmote(emote)
	}
	return emote
}
