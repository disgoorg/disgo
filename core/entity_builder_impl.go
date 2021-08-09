package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewEntityBuilder(disgo Disgo) EntityBuilder {
	return &EntityBuilderImpl{disgo: disgo}
}

// EntityBuilderImpl is used for creating structs used by Disgo
type EntityBuilderImpl struct {
	disgo Disgo
}

// Disgo returns the discord.Disgo client
func (b *EntityBuilderImpl) Disgo() Disgo {
	return b.disgo
}

// CreateInteraction creates an Interaction from the discord.UnmarshalInteraction response
func (b EntityBuilderImpl) CreateInteraction(unmarshalInteraction discord.UnmarshalInteraction, c chan discord.InteractionResponse, updateCache CacheStrategy) *Interaction {
	coreInteraction := &Interaction{
		UnmarshalInteraction: unmarshalInteraction,
		Disgo:                b.disgo,
		ResponseChannel:      c,
		Replied:              false,
		Data: &InteractionData{
			UnmarshalInteractionData: unmarshalInteraction.Data,
		},
	}

	if coreInteraction.Member != nil {
		coreInteraction.Member.Member.User = unmarshalInteraction.User // fuck u discord why not give the user here
		coreInteraction.Member = b.CreateMember(*unmarshalInteraction.GuildID, *unmarshalInteraction.Member, updateCache)
		coreInteraction.User = coreInteraction.Member.User
	} else {
		coreInteraction.User = b.CreateUser(unmarshalInteraction.User, updateCache)
	}

	return coreInteraction
}

// CreateCommandInteraction creates a CommandInteraction from the discord.UnmarshalInteraction response
func (b *EntityBuilderImpl) CreateCommandInteraction(interaction *Interaction, updateCache CacheStrategy) *CommandInteraction {
	commandInteraction := &CommandInteraction{
		Interaction: interaction,
		Data: &CommandInteractionData{
			InteractionData: interaction.Data,
		},
	}

	var subCommandName *string
	var subCommandGroupName *string

	unmarshalOptions := commandInteraction.Data.InteractionData.UnmarshalInteractionData.Options
	if len(unmarshalOptions) > 0 {
		unmarshalOption := unmarshalOptions[0]
		if unmarshalOption.Type == discord.CommandOptionTypeSubCommandGroup {
			subCommandGroupName = &unmarshalOption.Name
			unmarshalOptions = unmarshalOption.Options
			unmarshalOption = unmarshalOption.Options[0]
		}
		if unmarshalOption.Type == discord.CommandOptionTypeSubCommand {
			subCommandName = &unmarshalOption.Name
			unmarshalOptions = unmarshalOption.Options
		}
	}

	options := make([]CommandOption, len(unmarshalOptions))
	for i, optionData := range options {
		options[i] = CommandOption{
			Resolved: commandInteraction.Data.Resolved,
			Name:     optionData.Name,
			Type:     optionData.Type,
			Value:    optionData.Value,
		}
	}

	commandInteraction.Data.Options = options
	commandInteraction.Data.SubCommandName = subCommandName
	commandInteraction.Data.SubCommandGroupName = subCommandGroupName

	resolved := &Resolved{
		Resolved: interaction.Data.Resolved,
		Users:    map[discord.Snowflake]*User{},
		Members:  map[discord.Snowflake]*Member{},
		Roles:    map[discord.Snowflake]*Role{},
		Channels: map[discord.Snowflake]Channel{},
	}
	for id, user := range interaction.Data.Resolved.Users {
		resolved.Users[id] = b.CreateUser(user, updateCache)
	}

	for id, member := range interaction.Data.Resolved.Members {
		// discord omits the user field Oof
		member.User = interaction.Data.Resolved.Users[id]
		resolved.Members[id] = b.CreateMember(member.GuildID, member, updateCache)
	}

	for id, role := range interaction.Data.Resolved.Roles {
		resolved.Roles[id] = b.CreateRole(role.GuildID, role, updateCache)
	}

	for id, channel := range interaction.Data.Resolved.Channels {
		resolved.Channels[id] = b.CreateChannel(channel, updateCache)
	}

	commandInteraction.Data.Resolved = resolved

	return commandInteraction
}

// CreateComponentInteraction creates a ComponentInteraction from the discord.UnmarshalInteraction response
func (b *EntityBuilderImpl) CreateComponentInteraction(interaction *Interaction, updateCache CacheStrategy) *ComponentInteraction {
	return &ComponentInteraction{
		Interaction: interaction,
		Message:     b.CreateMessage(interaction.Message, updateCache),
		Data: &ComponentInteractionData{
			InteractionData: interaction.Data,
		},
	}
}

// CreateButtonInteraction creates a ButtonInteraction from the discord.UnmarshalInteraction response
func (b *EntityBuilderImpl) CreateButtonInteraction(componentInteraction *ComponentInteraction) *ButtonInteraction {
	return &ButtonInteraction{
		ComponentInteraction: componentInteraction,
		Data: &ButtonInteractionData{
			ComponentInteractionData: componentInteraction.Data,
		},
	}
}

// CreateSelectMenuInteraction creates a SelectMenuInteraction from the discord.UnmarshalInteraction response
func (b *EntityBuilderImpl) CreateSelectMenuInteraction(componentInteraction *ComponentInteraction) *SelectMenuInteraction {
	return &SelectMenuInteraction{
		ComponentInteraction: componentInteraction,
		Data: &SelectMenuInteractionData{
			ComponentInteractionData: componentInteraction.Data,
		},
	}
}

// CreateUser returns a new User entity
func (b *EntityBuilderImpl) CreateUser(user discord.User, updateCache CacheStrategy) *User {
	coreUser := &User{
		User:  user,
		Disgo: b.disgo,
	}
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().UserCache().Cache(coreUser)
	}
	return coreUser
}

// CreateSelfUser returns a new SelfUser entity
func (b *EntityBuilderImpl) CreateSelfUser(selfUser discord.SelfUser, updateCache CacheStrategy) *SelfUser {
	coreSelfUser := &SelfUser{
		SelfUser: selfUser,
		Disgo:    b.disgo,
		User:     b.CreateUser(selfUser.User, updateCache),
	}
	b.Disgo().SetSelfUser(coreSelfUser)
	return coreSelfUser
}

// CreateMessage returns a new discord.Message entity
func (b *EntityBuilderImpl) CreateMessage(message discord.Message, updateCache CacheStrategy) *Message {
	coreMsg := &Message{
		Message: message,
		Disgo:   b.disgo,
	}

	if message.Member != nil {
		coreMsg.Member.Member.User = message.Author // set the underlying discord data which is stored in cache
		coreMsg.Member = b.CreateMember(*message.GuildID, *message.Member, updateCache)
		coreMsg.Author = coreMsg.Member.User
	} else {
		coreMsg.Author = b.CreateUser(message.Author, updateCache)
	}

	if len(message.Components) > 0 {
		coreMsg.Components = b.CreateComponents(message.Components, updateCache)
	}

	// TODO: should we cache mentioned users, members, etc?
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().MessageCache().Cache(coreMsg)
	}
	return coreMsg
}

// CreateComponents returns a new slice of Component entities
func (b *EntityBuilderImpl) CreateComponents(unmarshalComponents []discord.UnmarshalComponent, updateCache CacheStrategy) []Component {
	components := make([]Component, len(unmarshalComponents))
	for i, component := range unmarshalComponents {
		switch component.Type {
		case discord.ComponentTypeActionRow:
			actionRow := ActionRow{
				UnmarshalComponent: component,
			}
			if len(component.Components) > 0 {
				actionRow.Components = b.CreateComponents(component.Components, updateCache)
			}
			components[i] = actionRow

		case discord.ComponentTypeButton:
			components[i] = Button{
				UnmarshalComponent: component,
			}

		case discord.ComponentTypeSelectMenu:
			components[i] = SelectMenu{
				UnmarshalComponent: component,
			}
		}
	}
	return components
}

// CreateGuildTemplate returns a new discord.GuildTemplate entity
func (b *EntityBuilderImpl) CreateGuildTemplate(guildTemplate discord.GuildTemplate, updateCache CacheStrategy) *GuildTemplate {
	coreTemplate := &GuildTemplate{
		GuildTemplate: guildTemplate,
	}

	if coreTemplate.Creator != nil {
		coreTemplate.Creator = b.CreateUser(guildTemplate.Creator, updateCache)
	}
	return coreTemplate
}

// CreateGuild returns a new discord.Guild entity
func (b *EntityBuilderImpl) CreateGuild(guild discord.Guild, updateCache CacheStrategy) *Guild {
	coreGuild := &Guild{
		Guild: guild,
		Disgo: b.disgo,
	}
	for _, channel := range guild.Channels {
		channel.GuildID = &guild.ID
		b.CreateChannel(channel, updateCache)
	}

	for _, role := range guild.Roles {
		role.GuildID = guild.ID
		b.CreateRole(guild.ID, role, updateCache)
	}

	for _, member := range guild.Members {
		b.CreateMember(guild.ID, member, updateCache)
	}

	for _, voiceState := range guild.VoiceStates {
		b.CreateVoiceState(guild.ID, voiceState, updateCache)
	}

	for _, emote := range guild.Emojis {
		b.CreateEmoji(guild.ID, emote, updateCache)
	}

	// TODO: presence
	/*for i := range fullGuild.Presences {
		presence := fullGuild.Presences[i]
		presence.Disgo = disgo
		b.Disgo().Cache().CachePresence(presence)
	}*/

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().GuildCache().Cache(coreGuild)
	}
	return coreGuild
}

// CreateMember returns a new discord.Member entity
func (b *EntityBuilderImpl) CreateMember(guildID discord.Snowflake, member discord.Member, updateCache CacheStrategy) *Member {
	coreMember := &Member{
		Member: member,
		Disgo:  b.disgo,
	}

	coreMember.GuildID = guildID
	coreMember.User = b.CreateUser(member.User, updateCache)
	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().MemberCache().Cache(coreMember)
	}
	return coreMember
}

func (b *EntityBuilderImpl) CreateBan(guildID discord.Snowflake, ban discord.Ban, updateCache CacheStrategy) *Ban {
	return &Ban{
		Ban:     ban,
		Disgo:   b.disgo,
		User:    b.CreateUser(ban.User, updateCache),
		GuildID: guildID,
	}
}

// CreateVoiceState returns a new discord.VoiceState entity
func (b *EntityBuilderImpl) CreateVoiceState(guildID discord.Snowflake, voiceState discord.VoiceState, updateCache CacheStrategy) *VoiceState {
	coreState := &VoiceState{
		VoiceState: voiceState,
		Disgo:      b.disgo,
	}
	if voiceState.Member != nil {
		coreState.Member = b.CreateMember(guildID, *voiceState.Member, updateCache)
	}

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().VoiceStateCache().Cache(coreState)
	}
	return coreState
}

// CreateCommand returns a new discord.Command entity
func (b *EntityBuilderImpl) CreateCommand(command discord.Command, updateCache CacheStrategy) *Command {
	coreCommand := &Command{
		Command: command,
		Disgo:   b.disgo,
	}
	if updateCache(b.Disgo()) {
		if command.GuildID == nil {
			return b.Disgo().Cache().GlobalCommandCache().Cache(coreCommand)
		} else {
			return b.Disgo().Cache().GuildCommandCache().Cache(coreCommand)
		}
	}
	return coreCommand
}

// CreateCommandPermissions returns a new discord.GuildCommandPermissions entity
func (b *EntityBuilderImpl) CreateCommandPermissions(guildCommandPermissions discord.GuildCommandPermissions, updateCache CacheStrategy) *GuildCommandPermissions {
	coreGuildCommandPermissions := &GuildCommandPermissions{
		GuildCommandPermissions: guildCommandPermissions,
		Disgo:                   b.disgo,
	}

	if updateCache(b.Disgo()) && b.Disgo().Cache().CacheFlags().Has(CacheFlagCommandPermissions) {
		if cmd := b.Disgo().Cache().GuildCommandCache().Get(guildCommandPermissions.ID); cmd != nil {
			cmd.GuildPermissions[guildCommandPermissions.GuildID] = coreGuildCommandPermissions
		}
	}
	return coreGuildCommandPermissions
}

// CreateRole returns a new discord.Role entity
func (b *EntityBuilderImpl) CreateRole(guildID discord.Snowflake, role discord.Role, updateCache CacheStrategy) *Role {
	coreRole := &Role{
		Role: role,
	}

	coreRole.GuildID = guildID

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().RoleCache().Cache(coreRole)
	}
	return coreRole
}

// CreateAuditLog returns a new discord.AuditLog entity
func (b *EntityBuilderImpl) CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog {
	coreAuditLog := &AuditLog{
		AuditLog:      auditLog,
		Disgo:         b.disgo,
		GuildID:       guildID,
		FilterOptions: filterOptions,
	}
	for _, user := range auditLog.Users {
		coreAuditLog.Users[user.ID] = b.CreateUser(user, updateCache)
	}
	for _, integration := range auditLog.Integrations {
		coreAuditLog.Integrations[integration.ID] = b.CreateIntegration(guildID, integration, updateCache)
	}
	for _, webhook := range auditLog.Webhooks {
		coreAuditLog.Webhooks[webhook.ID] = b.CreateWebhook(guildID, webhook)
	}
	return coreAuditLog
}

// CreateIntegration returns a new discord.Integration entity
func (b *EntityBuilderImpl) CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) *Integration {
	coreIntegration := &Integration{
		Integration: integration,
		Disgo:       b.disgo,
		GuildID:     guildID,
	}

	coreIntegration.User = b.CreateUser(*integration.User, updateCache)

	if integration.Application != nil {
		coreIntegration.Application = &IntegrationApplication{
			IntegrationApplication: *integration.Application,
			Bot:                    b.CreateUser(integration.Application.Bot, updateCache),
		}
	}
	return coreIntegration
}

// CreateWebhook returns a new Webhook entity
func (b *EntityBuilderImpl) CreateWebhook(guildID discord.Snowflake, webhook discord.Webhook) *Webhook {
	coreWebhook := &Webhook{
		Webhook: webhook,
		Disgo:   b.disgo,
		GuildID: guildID,
	}
	return coreWebhook
}

// CreateChannel returns a new Channel entity
func (b *EntityBuilderImpl) CreateChannel(discordChannel discord.Channel, updateCache CacheStrategy) Channel {
	channel := &channelImpl{
		Channel: discordChannel,
		disgo:   b.disgo,
	}
	if updateCache(b.Disgo()) {
		switch channel.Type() {
		case discord.ChannelTypeText:
			return b.Disgo().Cache().TextChannelCache().Cache(channel)
		case discord.ChannelTypeCategory:
			return b.Disgo().Cache().CategoryCache().Cache(channel)
		case discord.ChannelTypeDM:
			return b.Disgo().Cache().DMChannelCache().Cache(channel)
		case discord.ChannelTypeNews:
			return b.Disgo().Cache().NewsChannelCache().Cache(channel)
		case discord.ChannelTypeStore:
			return b.Disgo().Cache().StoreChannelCache().Cache(channel)
		case discord.ChannelTypeStage:
			return b.Disgo().Cache().StageChannelCache().Cache(channel)
		case discord.ChannelTypeVoice:
			return b.Disgo().Cache().VoiceChannelCache().Cache(channel)
		}
	}
	return channel
}

func (b *EntityBuilderImpl) CreateStageInstance(stageInstance discord.StageInstance, updateCache CacheStrategy) *StageInstance {
	// TODO
	panic("implement me")
}

// CreateEmoji returns a new discord.Emoji entity
func (b *EntityBuilderImpl) CreateEmoji(guildID discord.Snowflake, emoji discord.Emoji, updateCache CacheStrategy) *Emoji {
	coreEmoji := &Emoji{
		Emoji: emoji,
		Disgo: b.disgo,
	}
	if emoji.ID == "" { // return if emoji is no custom emote
		return coreEmoji
	}

	coreEmoji.GuildID = guildID

	if updateCache(b.Disgo()) {
		return b.Disgo().Cache().EmojiCache().Cache(coreEmoji)
	}
	return coreEmoji
}
