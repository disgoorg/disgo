package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

var _ EntityBuilder = (*entityBuilderImpl)(nil)

func NewEntityBuilder(bot *Bot) EntityBuilder {
	return &entityBuilderImpl{bot: bot}
}

// entityBuilderImpl is used for creating structs used by Disgo
type entityBuilderImpl struct {
	bot *Bot
}

// Bot returns the discord.Bot client
func (b *entityBuilderImpl) Bot() *Bot {
	return b.bot
}

// CreateInteraction creates an Interaction from the discord.Interaction response
func (b *entityBuilderImpl) CreateInteraction(interaction discord.Interaction, c chan<- discord.InteractionResponse, updateCache CacheStrategy) *Interaction {
	coreInteraction := &Interaction{
		Interaction:     interaction,
		Bot:             b.Bot(),
		ResponseChannel: c,
		Responded:       false,
	}

	if interaction.Member != nil {
		interaction.Member.User = interaction.User // fuck u discord why not give the user here
		coreInteraction.Member = b.CreateMember(*interaction.GuildID, *interaction.Member, updateCache)
		coreInteraction.User = coreInteraction.Member.User
	} else {
		coreInteraction.User = b.CreateUser(interaction.User, updateCache)
	}

	return coreInteraction
}

func (b *entityBuilderImpl) CreateApplicationCommandInteraction(interaction *Interaction, updateCache CacheStrategy) *ApplicationCommandInteraction {
	commandInteraction := &ApplicationCommandInteraction{
		Interaction: interaction,
		ApplicationCommandInteractionData: ApplicationCommandInteractionData{
			CommandName: interaction.Data.Name,
		},
	}

	resolved := &Resolved{
		Users:    map[discord.Snowflake]*User{},
		Members:  map[discord.Snowflake]*Member{},
		Roles:    map[discord.Snowflake]*Role{},
		Channels: map[discord.Snowflake]*Channel{},
		Messages: map[discord.Snowflake]*Message{},
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

	for id, message := range interaction.Data.Resolved.Messages {
		resolved.Messages[id] = b.CreateMessage(message, updateCache)
	}

	commandInteraction.ApplicationCommandInteractionData.Resolved = resolved

	return commandInteraction
}

func (b *entityBuilderImpl) CreateSlashCommandInteraction(applicationInteraction *ApplicationCommandInteraction) *SlashCommandInteraction {
	slashCommandInteraction := &SlashCommandInteraction{
		ApplicationCommandInteraction: applicationInteraction,
	}

	unmarshalOptions := slashCommandInteraction.Data.Options
	if len(unmarshalOptions) > 0 {
		unmarshalOption := unmarshalOptions[0]
		if unmarshalOption.Type == discord.CommandOptionTypeSubCommandGroup {
			slashCommandInteraction.SubCommandGroupName = &unmarshalOption.Name
			unmarshalOptions = unmarshalOption.Options
			unmarshalOption = unmarshalOption.Options[0]
		}
		if unmarshalOption.Type == discord.CommandOptionTypeSubCommand {
			slashCommandInteraction.SubCommandName = &unmarshalOption.Name
			unmarshalOptions = unmarshalOption.Options
		}
	}

	optionMap := make(map[string]SlashCommandOption, len(unmarshalOptions))
	for _, option := range unmarshalOptions {
		optionMap[option.Name] = SlashCommandOption{
			Resolved: slashCommandInteraction.Resolved,
			Name:     option.Name,
			Type:     option.Type,
			Value:    option.Value,
		}
	}
	slashCommandInteraction.Options = optionMap

	return slashCommandInteraction
}

func (b *entityBuilderImpl) CreateContextCommandInteraction(applicationInteraction *ApplicationCommandInteraction) *ContextCommandInteraction {
	return &ContextCommandInteraction{
		ApplicationCommandInteraction: applicationInteraction,
		ContextCommandInteractionData: ContextCommandInteractionData{
			TargetID: applicationInteraction.Data.TargetID,
		},
	}
}

func (b *entityBuilderImpl) CreateUserCommandInteraction(contextCommandInteraction *ContextCommandInteraction) *UserCommandInteraction {
	return &UserCommandInteraction{
		ContextCommandInteraction: contextCommandInteraction,
	}
}

func (b *entityBuilderImpl) CreateMessageCommandInteraction(contextCommandInteraction *ContextCommandInteraction) *MessageCommandInteraction {
	return &MessageCommandInteraction{
		ContextCommandInteraction: contextCommandInteraction,
	}
}

// CreateComponentInteraction creates a ComponentInteraction from the discord.Interaction response
func (b *entityBuilderImpl) CreateComponentInteraction(interaction *Interaction, updateCache CacheStrategy) *ComponentInteraction {
	return &ComponentInteraction{
		Interaction: interaction,
		Message:     b.CreateMessage(interaction.Message, updateCache),
	}
}

// CreateButtonInteraction creates a ButtonInteraction from the discord.Interaction response
func (b *entityBuilderImpl) CreateButtonInteraction(componentInteraction *ComponentInteraction) *ButtonInteraction {
	return &ButtonInteraction{
		ComponentInteraction: componentInteraction,
	}
}

// CreateSelectMenuInteraction creates a SelectMenuInteraction from the discord.Interaction response
func (b *entityBuilderImpl) CreateSelectMenuInteraction(componentInteraction *ComponentInteraction) *SelectMenuInteraction {
	return &SelectMenuInteraction{
		ComponentInteraction: componentInteraction,
		SelectMenuInteractionData: SelectMenuInteractionData{
			Values: componentInteraction.Data.Values,
		},
	}
}

// CreateUser returns a new User entity
func (b *entityBuilderImpl) CreateUser(user discord.User, updateCache CacheStrategy) *User {
	coreUser := &User{
		User: user,
		Bot:  b.Bot(),
	}
	if updateCache(b.Bot()) {
		return b.Bot().Caches.UserCache().Set(coreUser)
	}
	return coreUser
}

// CreateSelfUser returns a new SelfUser entity
func (b *entityBuilderImpl) CreateSelfUser(selfUser discord.OAuth2User, updateCache CacheStrategy) *SelfUser {
	coreSelfUser := &SelfUser{
		OAuth2User: selfUser,
		Bot:        b.Bot(),
		User:       b.CreateUser(selfUser.User, updateCache),
	}
	b.Bot().SelfUser = coreSelfUser
	return coreSelfUser
}

// CreateMessage returns a new discord.Message entity
func (b *entityBuilderImpl) CreateMessage(message discord.Message, updateCache CacheStrategy) *Message {
	coreMsg := &Message{
		Message: message,
		Bot:     b.Bot(),
	}

	if message.Member != nil {
		message.Member.User = message.Author // set the underlying discord data which is stored in caches
		coreMsg.Member = b.CreateMember(*message.GuildID, *message.Member, updateCache)
		coreMsg.Author = coreMsg.Member.User
	} else {
		coreMsg.Author = b.CreateUser(message.Author, updateCache)
	}

	if len(message.Components) > 0 {
		coreMsg.Components = b.CreateComponents(message.Components, updateCache)
	}

	if len(message.Stickers) > 0 {
		coreMsg.Stickers = make([]*MessageSticker, len(message.Stickers))
	}

	for i, sticker := range message.Stickers {
		coreMsg.Stickers[i] = b.CreateMessageSticker(sticker)
	}

	// TODO: should we caches mentioned users, members, etc?
	if updateCache(b.Bot()) {
		return b.Bot().Caches.MessageCache().Set(coreMsg)
	}
	return coreMsg
}

// CreateComponents returns a new slice of Component entities
func (b *entityBuilderImpl) CreateComponents(unmarshalComponents []discord.Component, updateCache CacheStrategy) []Component {
	components := make([]Component, len(unmarshalComponents))
	for i, component := range unmarshalComponents {
		switch component.Type {
		case discord.ComponentTypeActionRow:
			actionRow := ActionRow{
				Component: component,
			}
			if len(component.Components) > 0 {
				actionRow.Components = b.CreateComponents(component.Components, updateCache)
			}
			components[i] = actionRow

		case discord.ComponentTypeButton:
			components[i] = Button{
				Component: component,
			}

		case discord.ComponentTypeSelectMenu:
			components[i] = SelectMenu{
				Component: component,
			}
		}
	}
	return components
}

// CreateGuildTemplate returns a new discord.GuildTemplate entity
func (b *entityBuilderImpl) CreateGuildTemplate(guildTemplate discord.GuildTemplate, updateCache CacheStrategy) *GuildTemplate {
	coreTemplate := &GuildTemplate{
		GuildTemplate: guildTemplate,
	}

	if coreTemplate.Creator != nil {
		coreTemplate.Creator = b.CreateUser(guildTemplate.Creator, updateCache)
	}
	return coreTemplate
}

// CreateGuild returns a new discord.Guild entity
func (b *entityBuilderImpl) CreateGuild(guild discord.Guild, updateCache CacheStrategy) *Guild {
	coreGuild := &Guild{
		Guild: guild,
		Bot:   b.Bot(),
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.GuildCache().Set(coreGuild)
	}
	return coreGuild
}

// CreateMember returns a new discord.Member entity
func (b *entityBuilderImpl) CreateMember(guildID discord.Snowflake, member discord.Member, updateCache CacheStrategy) *Member {
	coreMember := &Member{
		Member: member,
		Bot:    b.Bot(),
	}

	coreMember.GuildID = guildID
	coreMember.User = b.CreateUser(member.User, updateCache)
	if updateCache(b.Bot()) {
		return b.Bot().Caches.MemberCache().Set(coreMember)
	}
	return coreMember
}

func (b *entityBuilderImpl) CreateBan(guildID discord.Snowflake, ban discord.Ban, updateCache CacheStrategy) *Ban {
	return &Ban{
		Ban:     ban,
		Bot:     b.Bot(),
		User:    b.CreateUser(ban.User, updateCache),
		GuildID: guildID,
	}
}

// CreateVoiceState returns a new discord.VoiceState entity
func (b *entityBuilderImpl) CreateVoiceState(voiceState discord.VoiceState, updateCache CacheStrategy) *VoiceState {
	coreState := &VoiceState{
		VoiceState: voiceState,
		Bot:        b.Bot(),
	}
	if voiceState.Member != nil {
		coreState.Member = b.CreateMember(voiceState.GuildID, *voiceState.Member, updateCache)
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.VoiceStateCache().Set(coreState)
	}
	return coreState
}

// CreateApplicationCommand returns a new discord.ApplicationCommand entity
func (b *entityBuilderImpl) CreateApplicationCommand(command discord.ApplicationCommand) *ApplicationCommand {
	return &ApplicationCommand{
		ApplicationCommand: command,
		Bot:                b.Bot(),
	}
}

// CreateApplicationCommandPermissions returns a new discord.ApplicationCommandPermissions entity
func (b *entityBuilderImpl) CreateApplicationCommandPermissions(guildCommandPermissions discord.ApplicationCommandPermissions) *ApplicationCommandPermissions {
	coreGuildCommandPermissions := &ApplicationCommandPermissions{
		ApplicationCommandPermissions: guildCommandPermissions,
		Bot:                           b.Bot(),
	}

	return coreGuildCommandPermissions
}

// CreateRole returns a new discord.Role entity
func (b *entityBuilderImpl) CreateRole(guildID discord.Snowflake, role discord.Role, updateCache CacheStrategy) *Role {
	coreRole := &Role{
		Role: role,
	}

	coreRole.GuildID = guildID

	if updateCache(b.Bot()) {
		return b.Bot().Caches.RoleCache().Set(coreRole)
	}
	return coreRole
}

// CreateAuditLog returns a new discord.AuditLog entity
func (b *entityBuilderImpl) CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog {
	coreAuditLog := &AuditLog{
		AuditLog:      auditLog,
		Bot:           b.Bot(),
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
		coreAuditLog.Webhooks[webhook.ID] = b.CreateWebhook(webhook)
	}
	return coreAuditLog
}

// CreateIntegration returns a new discord.Integration entity
func (b *entityBuilderImpl) CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) *Integration {
	coreIntegration := &Integration{
		Integration: integration,
		Bot:         b.Bot(),
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
func (b *entityBuilderImpl) CreateWebhook(webhook discord.Webhook) *Webhook {
	coreWebhook := &Webhook{
		Webhook: webhook,
		Bot:     b.Bot(),
	}
	return coreWebhook
}

// CreateChannel returns a new Channel entity
func (b *entityBuilderImpl) CreateChannel(channel discord.Channel, updateCache CacheStrategy) *Channel {
	coreChannel := &Channel{
		Channel: channel,
		Bot:     b.Bot(),
	}
	if channel.Type == discord.ChannelTypeVoice || channel.Type == discord.ChannelTypeStage {
		coreChannel.ConnectedMemberIDs = map[discord.Snowflake]struct{}{}
	}
	if updateCache(b.Bot()) {
		return b.Bot().Caches.ChannelCache().Set(coreChannel)
	}
	return coreChannel
}

func (b *entityBuilderImpl) CreateStageInstance(stageInstance discord.StageInstance, updateCache CacheStrategy) *StageInstance {

	coreStageInstance := &StageInstance{StageInstance: stageInstance, Bot: b.Bot()}

	if channel := b.Bot().Caches.ChannelCache().Get(stageInstance.ChannelID); channel != nil {
		channel.StageInstanceID = &stageInstance.ID
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.StageInstanceCache().Set(coreStageInstance)
	}
	return coreStageInstance
}

func (b *entityBuilderImpl) CreateInvite(invite discord.Invite, updateCache CacheStrategy) *Invite {
	coreInvite := &Invite{
		Invite: invite,
		Bot:    b.Bot(),
	}

	if invite.Inviter != nil {
		coreInvite.Inviter = b.CreateUser(*invite.Inviter, updateCache)
	}

	if invite.TargetUser != nil {
		coreInvite.TargetUser = b.CreateUser(*invite.TargetUser, updateCache)
	}

	return coreInvite
}

// CreateEmoji returns a new discord.Emoji entity
func (b *entityBuilderImpl) CreateEmoji(guildID discord.Snowflake, emoji discord.Emoji, updateCache CacheStrategy) *Emoji {
	coreEmoji := &Emoji{
		Emoji: emoji,
		Bot:   b.Bot(),
	}
	if emoji.ID == "" { // return if emoji is no custom emote
		return coreEmoji
	}

	coreEmoji.GuildID = guildID

	if updateCache(b.Bot()) {
		return b.Bot().Caches.EmojiCache().Set(coreEmoji)
	}
	return coreEmoji
}

func (b *entityBuilderImpl) CreateStickerPack(stickerPack discord.StickerPack, updateCache CacheStrategy) *StickerPack {
	coreStickerPack := &StickerPack{
		StickerPack: stickerPack,
		Bot:         b.Bot(),
		Stickers:    make([]*Sticker, len(stickerPack.Stickers)),
	}
	for i, sticker := range stickerPack.Stickers {
		coreStickerPack.Stickers[i] = b.CreateSticker(sticker, updateCache)
	}

	return coreStickerPack
}

// CreateSticker returns a new discord.Sticker entity
func (b *entityBuilderImpl) CreateSticker(sticker discord.Sticker, updateCache CacheStrategy) *Sticker {
	coreSticker := &Sticker{
		Sticker: sticker,
		Bot:     b.Bot(),
	}

	if sticker.User != nil {
		coreSticker.User = b.Bot().EntityBuilder.CreateUser(*sticker.User, CacheStrategyNo)
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.StickerCache().Set(coreSticker)
	}
	return coreSticker
}

// CreateMessageSticker returns a new discord.Sticker entity
func (b *entityBuilderImpl) CreateMessageSticker(messageSticker discord.MessageSticker) *MessageSticker {
	return &MessageSticker{
		MessageSticker: messageSticker,
		Bot:            b.Bot(),
	}
}
