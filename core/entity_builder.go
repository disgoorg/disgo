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

	CreateGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, updateCache CacheStrategy) *GuildScheduledEvent
	CreateGuildScheduledEventUser(guildID discord.Snowflake, guildScheduledEventUser discord.GuildScheduledEventUser, updateCache CacheStrategy) *GuildScheduledEventUser

	CreateRole(guildID discord.Snowflake, role discord.Role, updateCache CacheStrategy) *Role
	CreateMember(guildID discord.Snowflake, member discord.Member, updateCache CacheStrategy) *Member
	CreateBan(guildID discord.Snowflake, ban discord.Ban, updateCache CacheStrategy) *Ban
	CreateVoiceState(voiceState discord.VoiceState, updateCache CacheStrategy) *VoiceState

	CreateApplicationCommand(applicationCommand discord.ApplicationCommand) ApplicationCommand
	CreateApplicationCommandPermissions(guildCommandPermissions discord.ApplicationCommandPermissions) *ApplicationCommandPermissions

	CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog
	CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) Integration

	CreateChannel(channel discord.Channel, updateCache CacheStrategy) Channel
	CreateThreadMember(threadMember discord.ThreadMember, updateCache CacheStrategy) *ThreadMember

	CreateInvite(invite discord.Invite, updateCache CacheStrategy) *Invite

	CreateEmoji(guildID discord.Snowflake, emoji discord.Emoji, updateCache CacheStrategy) *Emoji
	CreateStickerPack(stickerPack discord.StickerPack, updateCache CacheStrategy) *StickerPack
	CreateSticker(sticker discord.Sticker, updateCache CacheStrategy) *Sticker
	CreateMessageSticker(sticker discord.MessageSticker) *MessageSticker

	CreateWebhook(webhook discord.Webhook, updateCache CacheStrategy) Webhook
}

// entityBuilderImpl is used for creating structs used by Disgo
type entityBuilderImpl struct {
	bot *Bot
}

// Bot returns the discord.Bot client
func (b *entityBuilderImpl) Bot() *Bot {
	return b.bot
}

func (b *entityBuilderImpl) baseInteraction(baseInteraction discord.BaseInteraction, c chan<- discord.InteractionResponse, updateCache CacheStrategy) *BaseInteraction {
	member, user := b.parseMemberOrUser(baseInteraction.GuildID, baseInteraction.Member, baseInteraction.User, updateCache)
	return &BaseInteraction{
		ID:              baseInteraction.ID,
		ApplicationID:   baseInteraction.ApplicationID,
		Token:           baseInteraction.Token,
		Version:         baseInteraction.Version,
		GuildID:         baseInteraction.GuildID,
		ChannelID:       baseInteraction.ChannelID,
		Locale:          baseInteraction.Locale,
		GuildLocale:     baseInteraction.GuildLocale,
		Member:          member,
		User:            user,
		ResponseChannel: c,
		Acknowledged:    false,
		Bot:             b.bot,
	}
}

// CreateInteraction creates an Interaction from the discord.Interaction response
func (b *entityBuilderImpl) CreateInteraction(interaction discord.Interaction, c chan<- discord.InteractionResponse, updateCache CacheStrategy) Interaction {
	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		var interactionData ApplicationCommandInteractionData
		switch d := i.Data.(type) {
		case discord.SlashCommandInteractionData:
			data := &SlashCommandInteractionData{
				SlashCommandInteractionData: d,
				Resolved: &SlashCommandResolved{
					Users:    map[discord.Snowflake]*User{},
					Members:  map[discord.Snowflake]*Member{},
					Roles:    map[discord.Snowflake]*Role{},
					Channels: map[discord.Snowflake]Channel{},
				},
			}
			for id, u := range d.Resolved.Users {
				data.Resolved.Users[id] = b.CreateUser(u, updateCache)
			}

			for id, m := range d.Resolved.Members {
				// discord omits the user field Oof
				m.User = d.Resolved.Users[id]
				data.Resolved.Members[id] = b.CreateMember(*i.GuildID, m, updateCache)
			}

			for id, r := range d.Resolved.Roles {
				data.Resolved.Roles[id] = b.CreateRole(*i.GuildID, r, updateCache)
			}

			for id, c := range d.Resolved.Channels {
				data.Resolved.Channels[id] = b.CreateChannel(c, updateCache)
			}

			unmarshalOptions := d.Options
			if len(unmarshalOptions) > 0 {
				unmarshalOption := unmarshalOptions[0]
				if option, ok := unmarshalOption.(discord.SlashCommandOptionSubCommandGroup); ok {
					data.SubCommandGroupName = &option.OptionName
					unmarshalOptions = make([]discord.SlashCommandOption, len(option.Options))
					for ii := range option.Options {
						unmarshalOptions[ii] = option.Options[ii]
					}
					unmarshalOption = option.Options[0]
				}
				if option, ok := unmarshalOption.(discord.SlashCommandOptionSubCommand); ok {
					data.SubCommandName = &option.OptionName
					unmarshalOptions = option.Options
				}
			}

			data.Options = make(map[string]SlashCommandOption, len(unmarshalOptions))
			for _, option := range unmarshalOptions {
				var slashCommandOption SlashCommandOption
				switch o := option.(type) {
				case discord.SlashCommandOptionString:
					slashCommandOption = SlashCommandOptionString{
						SlashCommandOptionString: o,
						Resolved:                 data.Resolved,
					}

				case discord.SlashCommandOptionInt:
					slashCommandOption = SlashCommandOptionInt{
						SlashCommandOptionInt: o,
					}

				case discord.SlashCommandOptionBool:
					slashCommandOption = SlashCommandOptionBool{
						SlashCommandOptionBool: o,
					}

				case discord.SlashCommandOptionUser:
					slashCommandOption = SlashCommandOptionUser{
						SlashCommandOptionUser: o,
						Resolved:               data.Resolved,
					}

				case discord.SlashCommandOptionChannel:
					slashCommandOption = SlashCommandOptionChannel{
						SlashCommandOptionChannel: o,
						Resolved:                  data.Resolved,
					}

				case discord.SlashCommandOptionRole:
					slashCommandOption = SlashCommandOptionRole{
						SlashCommandOptionRole: o,
						Resolved:               data.Resolved,
					}

				case discord.SlashCommandOptionMentionable:
					slashCommandOption = SlashCommandOptionMentionable{
						SlashCommandOptionMentionable: o,
						Resolved:                      data.Resolved,
					}

				case discord.SlashCommandOptionFloat:
					slashCommandOption = SlashCommandOptionFloat{
						SlashCommandOptionFloat: o,
					}

				default:
					b.Bot().Logger.Errorf("unknown slash command option with type %d received", option.Type())
					continue
				}
				data.Options[option.Name()] = slashCommandOption
			}
			interactionData = data

		case discord.UserCommandInteractionData:
			data := &UserCommandInteractionData{
				UserCommandInteractionData: d,
				Resolved: &UserCommandResolved{
					Users:   map[discord.Snowflake]*User{},
					Members: map[discord.Snowflake]*Member{},
				},
			}
			for id, u := range d.Resolved.Users {
				data.Resolved.Users[id] = b.CreateUser(u, updateCache)
			}

			for id, m := range d.Resolved.Members {
				// discord omits the user field Oof
				m.User = d.Resolved.Users[id]
				data.Resolved.Members[id] = b.CreateMember(*i.GuildID, m, updateCache)
			}
			interactionData = data

		case discord.MessageCommandInteractionData:
			data := &MessageCommandInteractionData{
				MessageCommandInteractionData: d,
				Resolved: &MessageCommandResolved{
					Messages: map[discord.Snowflake]*Message{},
				},
			}
			for id, message := range d.Resolved.Messages {
				data.Resolved.Messages[id] = b.CreateMessage(message, updateCache)
			}
			interactionData = data
		}
		return &ApplicationCommandInteraction{
			ReplyInteraction: &ReplyInteraction{BaseInteraction: b.baseInteraction(i.BaseInteraction, c, updateCache)},
			Data:             interactionData,
		}

	case discord.ComponentInteraction:
		componentInteraction := &ComponentInteraction{
			ReplyInteraction: &ReplyInteraction{BaseInteraction: b.baseInteraction(i.BaseInteraction, c, updateCache)},
			Message:          b.CreateMessage(i.Message, updateCache),
		}
		switch d := i.Data.(type) {
		case discord.ButtonInteractionData:
			componentInteraction.Data = &ButtonInteractionData{
				ButtonInteractionData: d,
				interaction:           componentInteraction,
			}

		case discord.SelectMenuInteractionData:
			componentInteraction.Data = &SelectMenuInteractionData{
				SelectMenuInteractionData: d,
				interaction:               componentInteraction,
			}
		}
		return componentInteraction

	case discord.AutocompleteInteraction:
		autocompleteInteraction := &AutocompleteInteraction{
			BaseInteraction: b.baseInteraction(i.BaseInteraction, c, updateCache),
			Data: AutocompleteInteractionData{
				AutocompleteInteractionData: i.Data,
			},
		}

		unmarshalOptions := i.Data.Options
		if len(unmarshalOptions) > 0 {
			unmarshalOption := unmarshalOptions[0]
			if option, ok := unmarshalOption.(discord.AutocompleteOptionSubCommandGroup); ok {
				autocompleteInteraction.Data.SubCommandGroupName = &option.GroupName
				unmarshalOptions = make([]discord.AutocompleteOption, len(option.Options))
				for i := range option.Options {
					unmarshalOptions[i] = option.Options[i]
				}
				unmarshalOption = option.Options[0]
			}
			if option, ok := unmarshalOption.(discord.AutocompleteOptionSubCommand); ok {
				autocompleteInteraction.Data.SubCommandName = &option.CommandName
				unmarshalOptions = option.Options
			}
		}

		autocompleteInteraction.Data.Options = make(map[string]discord.AutocompleteOption, len(unmarshalOptions))
		for _, option := range unmarshalOptions {
			autocompleteInteraction.Data.Options[option.Name()] = option
		}

		return autocompleteInteraction

	default:
		b.Bot().Logger.Error("unknown interaction type %d received", interaction.Type())
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
		return b.Bot().Caches.Users().Set(coreUser)
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
		return b.Bot().Caches.Presences().Set(corePresence)
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
		return b.Bot().Caches.Messages().Set(coreMsg)
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
		return b.Bot().Caches.Guilds().Set(coreGuild)
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
		return b.Bot().Caches.Members().Set(coreMember)
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
		return b.Bot().Caches.VoiceStates().Set(coreState)
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
		b.Bot().Logger.Errorf("unknown application command type %d received", applicationCommand.Type())
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
		return b.Bot().Caches.Roles().Set(coreRole)
	}
	return coreRole
}

// CreateAuditLog returns a new discord.AuditLog entity
func (b *entityBuilderImpl) CreateAuditLog(guildID discord.Snowflake, auditLog discord.AuditLog, filterOptions AuditLogFilterOptions, updateCache CacheStrategy) *AuditLog {
	coreAuditLog := &AuditLog{
		AuditLog:      auditLog,
		GuildID:       guildID,
		FilterOptions: filterOptions,
		Bot:           b.Bot(),
	}
	for _, guildScheduledEvent := range auditLog.GuildScheduledEvents {
		coreAuditLog.GuildScheduledEvents[guildScheduledEvent.ID] = b.CreateGuildScheduledEvent(guildScheduledEvent, updateCache)
	}
	for _, integration := range auditLog.Integrations {
		coreAuditLog.Integrations[integration.ID()] = b.CreateIntegration(guildID, integration, updateCache)
	}
	for _, thread := range auditLog.Threads {
		coreAuditLog.Threads[thread.ID()] = b.CreateChannel(thread, updateCache).(GuildThread)
	}
	for _, user := range auditLog.Users {
		coreAuditLog.Users[user.ID] = b.CreateUser(user, updateCache)
	}
	for _, webhook := range auditLog.Webhooks {
		coreAuditLog.Webhooks[webhook.ID()] = b.CreateWebhook(webhook, updateCache)
	}
	return coreAuditLog
}

// CreateIntegration returns a new discord.Integration entity
func (b *entityBuilderImpl) CreateIntegration(guildID discord.Snowflake, integration discord.Integration, updateCache CacheStrategy) Integration {
	var coreIntegration Integration

	switch i := integration.(type) {
	case discord.TwitchIntegration:
		coreIntegration = &TwitchIntegration{
			TwitchIntegration: i,
			Bot:               b.Bot(),
			GuildID:           guildID,
			User:              b.CreateUser(i.User, updateCache),
		}

	case discord.YouTubeIntegration:
		coreIntegration = &YouTubeIntegration{
			YouTubeIntegration: i,
			Bot:                b.Bot(),
			GuildID:            guildID,
			User:               b.CreateUser(i.User, updateCache),
		}

	case discord.BotIntegration:
		coreIntegration = &BotIntegration{
			BotIntegration: i,
			Bot:            b.Bot(),
			GuildID:        guildID,
			Application: &IntegrationApplication{
				IntegrationApplication: i.Application,
				Bot:                    b.CreateUser(i.Application.Bot, updateCache),
			},
		}

	default:
		b.Bot().Logger.Errorf("unknown integration type %d received", integration.Type())
		return nil
	}

	return coreIntegration
}

// CreateWebhook returns a new Webhook entity
func (b *entityBuilderImpl) CreateWebhook(webhook discord.Webhook, updateCache CacheStrategy) Webhook {
	var coreWebhook Webhook

	switch w := webhook.(type) {
	case discord.IncomingWebhook:
		coreWebhook = &IncomingWebhook{
			IncomingWebhook: w,
			Bot:             b.Bot(),
			User:            b.CreateUser(w.User, updateCache),
		}

	case discord.ChannelFollowerWebhook:
		coreWebhook = &ChannelFollowerWebhook{
			ChannelFollowerWebhook: w,
			Bot:                    b.Bot(),
			User:                   b.CreateUser(w.User, updateCache),
		}

	case discord.ApplicationWebhook:
		coreWebhook = &ApplicationWebhook{
			ApplicationWebhook: w,
			Bot:                b.Bot(),
		}

	default:
		b.Bot().Logger.Errorf("unknown webhook type %d received", webhook.Type())
		return nil
	}

	return coreWebhook
}

// CreateChannel returns a new Channel entity
func (b *entityBuilderImpl) CreateChannel(channel discord.Channel, updateCache CacheStrategy) Channel {
	var c Channel
	switch ch := channel.(type) {
	case discord.GuildTextChannel:
		c = &GuildTextChannel{
			GuildTextChannel: ch,
			Bot:              b.Bot(),
		}

	case discord.DMChannel:
		c = &DMChannel{
			DMChannel: ch,
			Bot:       b.Bot(),
		}

	case discord.GuildVoiceChannel:
		c = &GuildVoiceChannel{
			GuildVoiceChannel:  ch,
			Bot:                b.Bot(),
			ConnectedMemberIDs: map[discord.Snowflake]struct{}{},
		}

	case discord.GroupDMChannel:
		c = &GroupDMChannel{
			GroupDMChannel: ch,
			Bot:            b.Bot(),
		}

	case discord.GuildCategoryChannel:
		c = &GuildCategoryChannel{
			GuildCategoryChannel: ch,
			Bot:                  b.Bot(),
		}

	case discord.GuildNewsChannel:
		c = &GuildNewsChannel{
			GuildNewsChannel: ch,
			Bot:              b.Bot(),
		}

	case discord.GuildStoreChannel:
		c = &GuildStoreChannel{
			GuildStoreChannel: ch,
			Bot:               b.Bot(),
		}

	case discord.GuildNewsThread:
		c = &GuildNewsThread{
			GuildNewsThread: ch,
			Bot:             b.Bot(),
		}

	case discord.GuildPrivateThread:
		c = &GuildPrivateThread{
			GuildPrivateThread: ch,
			Bot:                b.Bot(),
		}

	case discord.GuildPublicThread:
		c = &GuildPublicThread{
			GuildPublicThread: ch,
			Bot:               b.Bot(),
		}

	case discord.GuildStageVoiceChannel:
		c = &GuildStageVoiceChannel{
			GuildStageVoiceChannel: ch,
			Bot:                    b.Bot(),
			StageInstanceID:        nil,
			ConnectedMemberIDs:     map[discord.Snowflake]struct{}{},
		}

	default:
		panic("unknown channel type")
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.Channels().Set(c)
	}
	return c
}

func (b *entityBuilderImpl) CreateThreadMember(threadMember discord.ThreadMember, updateCache CacheStrategy) *ThreadMember {
	coreThreadMember := &ThreadMember{
		ThreadMember: threadMember,
		Bot:          b.Bot(),
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.ThreadMembers().Set(coreThreadMember)
	}
	return coreThreadMember
}

func (b *entityBuilderImpl) CreateStageInstance(stageInstance discord.StageInstance, updateCache CacheStrategy) *StageInstance {
	coreStageInstance := &StageInstance{
		StageInstance: stageInstance,
		Bot:           b.Bot(),
	}

	if channel := b.Bot().Caches.Channels().Get(stageInstance.ChannelID); channel != nil {
		if ch, ok := channel.(*GuildStageVoiceChannel); ok {
			ch.StageInstanceID = &stageInstance.ID
		}
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.StageInstances().Set(coreStageInstance)
	}
	return coreStageInstance
}

func (b *entityBuilderImpl) CreateGuildScheduledEvent(guildScheduledEvent discord.GuildScheduledEvent, updateCache CacheStrategy) *GuildScheduledEvent {
	coreGuildScheduledEvent := &GuildScheduledEvent{
		GuildScheduledEvent: guildScheduledEvent,
		Creator:             b.CreateUser(guildScheduledEvent.Creator, updateCache),
		Bot:                 b.Bot(),
	}

	if updateCache(b.Bot()) {
		return b.Bot().Caches.GuildScheduledEvents().Set(coreGuildScheduledEvent)
	}
	return coreGuildScheduledEvent
}

func (b *entityBuilderImpl) CreateGuildScheduledEventUser(guildID discord.Snowflake, guildScheduledEventUser discord.GuildScheduledEventUser, updateCache CacheStrategy) *GuildScheduledEventUser {
	coreGuildScheduledEventUser := &GuildScheduledEventUser{
		GuildScheduledEventUser: guildScheduledEventUser, Bot: b.Bot(),
		User: b.CreateUser(guildScheduledEventUser.User, updateCache),
	}
	if guildScheduledEventUser.Member != nil {
		coreGuildScheduledEventUser.Member = b.CreateMember(guildID, *guildScheduledEventUser.Member, updateCache)
	}

	return coreGuildScheduledEventUser
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
		return b.Bot().Caches.Emojis().Set(coreEmoji)
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
		return b.Bot().Caches.Stickers().Set(coreSticker)
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
