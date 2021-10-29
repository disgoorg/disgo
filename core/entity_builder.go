package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// CacheStrategy is used to determine whether something should be cached when making an api request. When using the
// gateway, you'll receive the event shortly afterwards if you have the correct GatewayIntents.
type CacheStrategy func(bot *Bot) bool

// Default caches strategy choices
var (
	CacheStrategyYes  CacheStrategy = func(bot *Bot) bool { return true }
	CacheStrategyNo   CacheStrategy = func(bot *Bot) bool { return true }
	CacheStrategyNoWs CacheStrategy = func(bot *Bot) bool { return bot.HasGateway() }
)

var _ EntityBuilder = (*entityBuilderImpl)(nil)

func NewEntityBuilder(bot *Bot) EntityBuilder {
	return &entityBuilderImpl{bot: bot}
}

// EntityBuilder is used to create structs for disgo's caches
type EntityBuilder interface {
	Bot() *Bot

	CreateInteraction(interaction discord.Interaction, responseChannel chan<- discord.InteractionResponse, updateCache CacheStrategy) Interaction

	CreateUser(user discord.User, updateCache CacheStrategy) *User
	CreateSelfUser(selfUser discord.OAuth2User, updateCache CacheStrategy) *SelfUser
	CreatePresence(presence discord.Presence, updateCache CacheStrategy) *Presence

	CreateMessage(message discord.Message, updateCache CacheStrategy) *Message

	CreateGuild(guild discord.Guild, updateCache CacheStrategy) *Guild
	CreateGuildTemplate(guildTemplate discord.GuildTemplate, updateCache CacheStrategy) *GuildTemplate
	CreateStageInstance(stageInstance discord.StageInstance, updateCache CacheStrategy) *StageInstance

	CreateRole(guildID discord.Snowflake, role discord.Role, updateCache CacheStrategy) *Role
	CreateMember(guildID discord.Snowflake, member discord.Member, updateCache CacheStrategy) *Member
	CreateBan(guildID discord.Snowflake, ban discord.Ban, updateCache CacheStrategy) *Ban
	CreateVoiceState(voiceState discord.VoiceState, updateCache CacheStrategy) *VoiceState

	CreateApplicationCommand(applicationCommand discord.ApplicationCommand) ApplicationCommand
	CreateApplicationCommandPermissions(guildCommandPermissions discord.ApplicationCommandPermissions) *ApplicationCommandPermissions

	CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog
	CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) *Integration

	CreateChannel(channel discord.Channel, updateCache CacheStrategy) *Channel

	CreateInvite(invite discord.Invite, updateCache CacheStrategy) *Invite

	CreateEmoji(guildID discord.Snowflake, emoji discord.Emoji, updateCache CacheStrategy) *Emoji
	CreateStickerPack(stickerPack discord.StickerPack, updateCache CacheStrategy) *StickerPack
	CreateSticker(sticker discord.Sticker, updateCache CacheStrategy) *Sticker
	CreateMessageSticker(sticker discord.MessageSticker) *MessageSticker

	CreateWebhook(webhook discord.Webhook) *Webhook
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
func (b *entityBuilderImpl) CreateInteraction(interaction discord.Interaction, c chan<- discord.InteractionResponse, updateCache CacheStrategy) Interaction {
	interactionFields := &InteractionFields{
		Bot:             b.Bot(),
		User:            nil,
		Member:          nil,
		ResponseChannel: c,
		Acknowledged:    false,
	}

	switch i := interaction.(type) {
	case discord.AutocompleteInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		autocompleteInteraction := &AutocompleteInteraction{
			InteractionFields: interactionFields,
			CommandID:         i.Data.CommandID,
			CommandName:       i.Data.CommandName,
		}

		unmarshalOptions := i.Data.Options
		if len(unmarshalOptions) > 0 {
			unmarshalOption := unmarshalOptions[0]
			if option, ok := unmarshalOption.(discord.AutocompleteOptionSubCommandGroup); ok {
				autocompleteInteraction.SubCommandGroupName = &option.Name
				unmarshalOptions = make([]discord.AutocompleteOption, len(option.Options))
				for i := range option.Options {
					unmarshalOptions[i] = option.Options[i]
				}
				unmarshalOption = option.Options[0]
			}
			if option, ok := unmarshalOption.(discord.AutocompleteOptionSubCommand); ok {
				autocompleteInteraction.SubCommandName = &option.Name
				unmarshalOptions = option.Options
			}
		}

		autocompleteInteraction.Options = make(map[string]discord.AutocompleteOption, len(unmarshalOptions))
		for _, option := range unmarshalOptions {
			var name string
			switch o := option.(type) {
			case discord.AutocompleteOptionString:
				name = o.Name

			case discord.AutocompleteOptionInt:
				name = o.Name

			case discord.AutocompleteOptionBool:
				name = o.Name

			case discord.AutocompleteOptionUser:
				name = o.Name

			case discord.AutocompleteOptionChannel:
				name = o.Name

			case discord.AutocompleteOptionRole:
				name = o.Name

			case discord.AutocompleteOptionMentionable:
				name = o.Name

			case discord.AutocompleteOptionFloat:
				name = o.Name

			default:
				b.Bot().Logger.Errorf("unknown application command autocomplete option with type %d received", option.Type)
				continue
			}
			autocompleteInteraction.Options[name] = option
		}

		return autocompleteInteraction

	case discord.SlashCommandInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		slashCommandInteraction := &SlashCommandInteraction{
			InteractionFields: interactionFields,
			CommandID:         i.Data.CommandID,
			CommandName:       i.Data.CommandName,
		}

		resolved := &SlashCommandResolved{
			Users:    map[discord.Snowflake]*User{},
			Members:  map[discord.Snowflake]*Member{},
			Roles:    map[discord.Snowflake]*Role{},
			Channels: map[discord.Snowflake]*Channel{},
		}
		slashCommandInteraction.Resolved = resolved
		for id, user := range i.Data.Resolved.Users {
			resolved.Users[id] = b.CreateUser(user, updateCache)
		}

		for id, member := range i.Data.Resolved.Members {
			// discord omits the user field Oof
			member.User = i.Data.Resolved.Users[id]
			resolved.Members[id] = b.CreateMember(*i.GuildID, member, updateCache)
		}

		for id, role := range i.Data.Resolved.Roles {
			resolved.Roles[id] = b.CreateRole(*i.GuildID, role, updateCache)
		}

		for id, channel := range i.Data.Resolved.Channels {
			resolved.Channels[id] = b.CreateChannel(channel, updateCache)
		}

		unmarshalOptions := i.Data.Options
		if len(unmarshalOptions) > 0 {
			unmarshalOption := unmarshalOptions[0]
			if option, ok := unmarshalOption.(discord.SlashCommandOptionSubCommandGroup); ok {
				slashCommandInteraction.SubCommandGroupName = &option.Name
				unmarshalOptions = make([]discord.SlashCommandOption, len(option.Options))
				for i := range option.Options {
					unmarshalOptions[i] = option.Options[i]
				}
				unmarshalOption = option.Options[0]
			}
			if option, ok := unmarshalOption.(discord.SlashCommandOptionSubCommand); ok {
				slashCommandInteraction.SubCommandName = &option.Name
				unmarshalOptions = option.Options
			}
		}

		slashCommandInteraction.Options = make(map[string]SlashCommandOption, len(unmarshalOptions))
		for _, option := range unmarshalOptions {
			switch o := option.(type) {
			case discord.SlashCommandOptionString:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionString{
					SlashCommandOptionString: o,
					Resolved:                 resolved,
				}

			case discord.SlashCommandOptionInt:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionInt{
					SlashCommandOptionInt: o,
				}

			case discord.SlashCommandOptionBool:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionBool{
					SlashCommandOptionBool: o,
				}

			case discord.SlashCommandOptionUser:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionUser{
					SlashCommandOptionUser: o,
					Resolved:               resolved,
				}

			case discord.SlashCommandOptionChannel:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionChannel{
					SlashCommandOptionChannel: o,
					Resolved:                  resolved,
				}

			case discord.SlashCommandOptionRole:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionRole{
					SlashCommandOptionRole: o,
					Resolved:               resolved,
				}

			case discord.SlashCommandOptionMentionable:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionMentionable{
					SlashCommandOptionMentionable: o,
					Resolved:                      resolved,
				}

			case discord.SlashCommandOptionFloat:
				slashCommandInteraction.Options[o.Name] = SlashCommandOptionFloat{
					SlashCommandOptionFloat: o,
				}

			default:
				b.Bot().Logger.Errorf("unknown application command autocomplete option with type %d received", option.Type)
				continue
			}
		}

		return slashCommandInteraction

	case discord.UserCommandInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		userCommandInteraction := &UserCommandInteraction{
			InteractionFields: interactionFields,
			CommandID:         i.Data.CommandID,
			CommandName:       i.Data.CommandName,
			Resolved: &UserCommandResolved{
				Users:   map[discord.Snowflake]*User{},
				Members: map[discord.Snowflake]*Member{},
			},
			TargetID: i.Data.TargetID,
		}

		for id, user := range i.Data.Resolved.Users {
			userCommandInteraction.Resolved.Users[id] = b.CreateUser(user, updateCache)
		}

		for id, member := range i.Data.Resolved.Members {
			// discord omits the user field Oof
			member.User = i.Data.Resolved.Users[id]
			userCommandInteraction.Resolved.Members[id] = b.CreateMember(*i.GuildID, member, updateCache)
		}

		return userCommandInteraction

	case discord.MessageCommandInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		messageCommandInteraction := &MessageCommandInteraction{
			InteractionFields: interactionFields,
			CommandID:         i.Data.CommandID,
			CommandName:       i.Data.CommandName,
			Resolved: &MessageCommandResolved{
				Messages: map[discord.Snowflake]*Message{},
			},
			TargetID: i.Data.TargetID,
		}

		for id, message := range i.Data.Resolved.Messages {
			messageCommandInteraction.Resolved.Messages[id] = b.CreateMessage(message, updateCache)
		}

		return messageCommandInteraction

	case discord.ButtonInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		message := b.CreateMessage(i.Message, updateCache)

		return &ButtonInteraction{
			InteractionFields: interactionFields,
			Message:           message,
			CustomID:          i.Data.CustomID,
		}

	case discord.SelectMenuInteraction:
		interactionFields.InteractionFields = i.InteractionFields
		interactionFields.Member, interactionFields.User = b.parseMemberOrUser(i.GuildID, i.Member, i.User, updateCache)

		message := b.CreateMessage(i.Message, updateCache)

		return &SelectMenuInteraction{
			InteractionFields: interactionFields,
			Message:           message,
			CustomID:          i.Data.CustomID,
			Values:            i.Data.Values,
		}

	default:
		b.Bot().Logger.Error("unknown interaction type %d received", interaction.InteractionType())
		return nil
	}
}

func (b *entityBuilderImpl) parseMemberOrUser(guildID *discord.Snowflake, member *discord.Member, user *discord.User, updateCache CacheStrategy) (rMember *Member, rUser *User) {
	if member != nil {
		rMember = b.CreateMember(*guildID, *member, updateCache)
		rUser = rMember.User
	} else {
		rUser = b.CreateUser(*user, updateCache)
	}
	return
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

func (b *entityBuilderImpl) CreatePresence(presence discord.Presence, updateCache CacheStrategy) *Presence {
	corePresence := &Presence{
		Presence: presence,
		Bot:      b.Bot(),
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.PresenceCache().Set(corePresence)
	}
	return corePresence
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
func (b *entityBuilderImpl) CreateApplicationCommand(applicationCommand discord.ApplicationCommand) ApplicationCommand {
	switch c := applicationCommand.(type) {
	case discord.SlashCommand:
		return &SlashCommand{
			SlashCommand: c,
			Bot:          b.Bot(),
		}

	case discord.UserCommand:
		return &UserCommand{
			UserCommand: c,
			Bot:         b.Bot(),
		}

	case discord.MessageCommand:
		return &MessageCommand{
			MessageCommand: c,
			Bot:            b.Bot(),
		}
	default:
		b.Bot().Logger.Error("")
		return nil
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
