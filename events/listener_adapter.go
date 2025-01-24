package events

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
)

var _ bot.EventListener = (*ListenerAdapter)(nil)

// ListenerAdapter lets you override the handles for receiving events
type ListenerAdapter struct {
	// raw event
	OnRaw func(event *Raw)

	// heartbeat ack event
	OnHeartbeatAck func(event *HeartbeatAck)

	// GuildApplicationCommandPermissionsUpdate
	OnGuildApplicationCommandPermissionsUpdate func(event *GuildApplicationCommandPermissionsUpdate)

	// AutoModeration
	OnAutoModerationRuleCreate      func(event *AutoModerationRuleCreate)
	OnAutoModerationRuleUpdate      func(event *AutoModerationRuleUpdate)
	OnAutoModerationRuleDelete      func(event *AutoModerationRuleDelete)
	OnAutoModerationActionExecution func(event *AutoModerationActionExecution)

	// Thread Events
	OnThreadCreate func(event *ThreadCreate)
	OnThreadUpdate func(event *ThreadUpdate)
	OnThreadDelete func(event *ThreadDelete)
	OnThreadShow   func(event *ThreadShow)
	OnThreadHide   func(event *ThreadHide)

	// ThreadMember Events
	OnThreadMemberAdd    func(event *ThreadMemberAdd)
	OnThreadMemberUpdate func(event *ThreadMemberUpdate)
	OnThreadMemberRemove func(event *ThreadMemberRemove)

	// Guild Channel Events
	OnGuildChannelCreate     func(event *GuildChannelCreate)
	OnGuildChannelUpdate     func(event *GuildChannelUpdate)
	OnGuildChannelDelete     func(event *GuildChannelDelete)
	OnGuildChannelPinsUpdate func(event *GuildChannelPinsUpdate)

	// DM Channel Events
	OnDMChannelPinsUpdate func(event *DMChannelPinsUpdate)

	// Channel Message Events
	OnDMMessageCreate func(event *DMMessageCreate)
	OnDMMessageUpdate func(event *DMMessageUpdate)
	OnDMMessageDelete func(event *DMMessageDelete)

	// Channel Reaction Events
	OnDMMessageReactionAdd         func(event *DMMessageReactionAdd)
	OnDMMessageReactionRemove      func(event *DMMessageReactionRemove)
	OnDMMessageReactionRemoveEmoji func(event *DMMessageReactionRemoveEmoji)
	OnDMMessageReactionRemoveAll   func(event *DMMessageReactionRemoveAll)

	// Emoji Events
	OnEmojisUpdate func(event *EmojisUpdate)
	OnEmojiCreate  func(event *EmojiCreate)
	OnEmojiUpdate  func(event *EmojiUpdate)
	OnEmojiDelete  func(event *EmojiDelete)

	// Entitlement Events
	OnEntitlementCreate func(event *EntitlementCreate)
	OnEntitlementUpdate func(event *EntitlementUpdate)
	OnEntitlementDelete func(event *EntitlementDelete)

	// Subscription Events
	OnSubscriptionCreate func(event *SubscriptionCreate)
	OnSubscriptionUpdate func(event *SubscriptionUpdate)
	OnSubscriptionDelete func(event *SubscriptionDelete)

	// Sticker Events
	OnStickersUpdate func(event *StickersUpdate)
	OnStickerCreate  func(event *StickerCreate)
	OnStickerUpdate  func(event *StickerUpdate)
	OnStickerDelete  func(event *StickerDelete)

	// gateway status Events
	OnReady   func(event *Ready)
	OnResumed func(event *Resumed)

	// Guild Events
	OnGuildJoin                func(event *GuildJoin)
	OnGuildUpdate              func(event *GuildUpdate)
	OnGuildLeave               func(event *GuildLeave)
	OnGuildAvailable           func(event *GuildAvailable)
	OnGuildUnavailable         func(event *GuildUnavailable)
	OnGuildReady               func(event *GuildReady)
	OnGuildsReady              func(event *GuildsReady)
	OnGuildBan                 func(event *GuildBan)
	OnGuildUnban               func(event *GuildUnban)
	OnGuildAuditLogEntryCreate func(event *GuildAuditLogEntryCreate)

	// Guild Invite Events
	OnGuildInviteCreate func(event *InviteCreate)
	OnGuildInviteDelete func(event *InviteDelete)

	// Guild Member Events
	OnGuildMemberJoin   func(event *GuildMemberJoin)
	OnGuildMemberUpdate func(event *GuildMemberUpdate)
	OnGuildMemberLeave  func(event *GuildMemberLeave)

	// Guild Message Events
	OnGuildMessageCreate func(event *GuildMessageCreate)
	OnGuildMessageUpdate func(event *GuildMessageUpdate)
	OnGuildMessageDelete func(event *GuildMessageDelete)

	// Guild Message Reaction Events
	OnGuildMessageReactionAdd         func(event *GuildMessageReactionAdd)
	OnGuildMessageReactionRemove      func(event *GuildMessageReactionRemove)
	OnGuildMessageReactionRemoveEmoji func(event *GuildMessageReactionRemoveEmoji)
	OnGuildMessageReactionRemoveAll   func(event *GuildMessageReactionRemoveAll)

	// Guild Soundboard Events
	OnGuildSoundboardSoundCreate  func(event *GuildSoundboardSoundCreate)
	OnGuildSoundboardSoundUpdate  func(event *GuildSoundboardSoundUpdate)
	OnGuildSoundboardSoundDelete  func(event *GuildSoundboardSoundDelete)
	OnGuildSoundboardSoundsUpdate func(event *GuildSoundboardSoundsUpdate)
	OnSoundboardSounds            func(event *SoundboardSounds)

	// Guild Voice Events
	OnVoiceServerUpdate           func(event *VoiceServerUpdate)
	OnGuildVoiceChannelEffectSend func(event *GuildVoiceChannelEffectSend)
	OnGuildVoiceStateUpdate       func(event *GuildVoiceStateUpdate)
	OnGuildVoiceJoin              func(event *GuildVoiceJoin)
	OnGuildVoiceMove              func(event *GuildVoiceMove)
	OnGuildVoiceLeave             func(event *GuildVoiceLeave)

	// Guild StageInstance Events
	OnStageInstanceCreate func(event *StageInstanceCreate)
	OnStageInstanceUpdate func(event *StageInstanceUpdate)
	OnStageInstanceDelete func(event *StageInstanceDelete)

	// Guild Role Events
	OnRoleCreate func(event *RoleCreate)
	OnRoleUpdate func(event *RoleUpdate)
	OnRoleDelete func(event *RoleDelete)

	// Guild Scheduled Events
	OnGuildScheduledEventCreate     func(event *GuildScheduledEventCreate)
	OnGuildScheduledEventUpdate     func(event *GuildScheduledEventUpdate)
	OnGuildScheduledEventDelete     func(event *GuildScheduledEventDelete)
	OnGuildScheduledEventUserAdd    func(event *GuildScheduledEventUserAdd)
	OnGuildScheduledEventUserRemove func(event *GuildScheduledEventUserRemove)

	// Interaction Events
	OnInteraction                   func(event *InteractionCreate)
	OnApplicationCommandInteraction func(event *ApplicationCommandInteractionCreate)
	OnComponentInteraction          func(event *ComponentInteractionCreate)
	OnAutocompleteInteraction       func(event *AutocompleteInteractionCreate)
	OnModalSubmit                   func(event *ModalSubmitInteractionCreate)

	// Message Events
	OnMessageCreate func(event *MessageCreate)
	OnMessageUpdate func(event *MessageUpdate)
	OnMessageDelete func(event *MessageDelete)

	// Message Poll Events
	OnMessagePollVoteAdd    func(event *MessagePollVoteAdd)
	OnMessagePollVoteRemove func(event *MessagePollVoteRemove)

	// DM Message Poll Events
	OnDMMessagePollVoteAdd    func(event *DMMessagePollVoteAdd)
	OnDMMessagePollVoteRemove func(event *DMMessagePollVoteRemove)

	// Guild Message Poll Events
	OnGuildMessagePollVoteAdd    func(event *GuildMessagePollVoteAdd)
	OnGuildMessagePollVoteRemove func(event *GuildMessagePollVoteRemove)

	// Message Reaction Events
	OnMessageReactionAdd         func(event *MessageReactionAdd)
	OnMessageReactionRemove      func(event *MessageReactionRemove)
	OnMessageReactionRemoveEmoji func(event *MessageReactionRemoveEmoji)
	OnMessageReactionRemoveAll   func(event *MessageReactionRemoveAll)

	// Self Events
	OnSelfUpdate func(event *SelfUpdate)

	// User Events
	OnUserUpdate             func(event *UserUpdate)
	OnUserTypingStart        func(event *UserTypingStart)
	OnGuildMemberTypingStart func(event *GuildMemberTypingStart)
	OnDMUserTypingStart      func(event *DMUserTypingStart)

	OnPresenceUpdate func(event *PresenceUpdate)

	// User Activity Events
	OnUserActivityStart  func(event *UserActivityStart)
	OnUserActivityUpdate func(event *UserActivityUpdate)
	OnUserActivityStop   func(event *UserActivityStop)

	OnUserStatusUpdate       func(event *UserStatusUpdate)
	OnUserClientStatusUpdate func(event *UserClientStatusUpdate)

	OnIntegrationCreate       func(event *IntegrationCreate)
	OnIntegrationUpdate       func(event *IntegrationUpdate)
	OnIntegrationDelete       func(event *IntegrationDelete)
	OnGuildIntegrationsUpdate func(event *GuildIntegrationsUpdate)

	OnGuildWebhooksUpdate func(event *WebhooksUpdate)
}

// OnEvent is getting called everytime we receive an event
func (l *ListenerAdapter) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *Raw:
		if listener := l.OnRaw; listener != nil {
			listener(e)
		}

	case *HeartbeatAck:
		if listener := l.OnHeartbeatAck; listener != nil {
			listener(e)
		}

	case *GuildApplicationCommandPermissionsUpdate:
		if listener := l.OnGuildApplicationCommandPermissionsUpdate; listener != nil {
			listener(e)
		}

	// Automoderation Events
	case *AutoModerationRuleCreate:
		if listener := l.OnAutoModerationRuleCreate; listener != nil {
			listener(e)
		}
	case *AutoModerationRuleUpdate:
		if listener := l.OnAutoModerationRuleUpdate; listener != nil {
			listener(e)
		}
	case *AutoModerationRuleDelete:
		if listener := l.OnAutoModerationRuleDelete; listener != nil {
			listener(e)
		}
	case *AutoModerationActionExecution:
		if listener := l.OnAutoModerationActionExecution; listener != nil {
			listener(e)
		}

	// Thread Events
	case *ThreadCreate:
		if listener := l.OnThreadCreate; listener != nil {
			listener(e)
		}
	case *ThreadUpdate:
		if listener := l.OnThreadUpdate; listener != nil {
			listener(e)
		}
	case *ThreadDelete:
		if listener := l.OnThreadDelete; listener != nil {
			listener(e)
		}
	case *ThreadShow:
		if listener := l.OnThreadShow; listener != nil {
			listener(e)
		}
	case *ThreadHide:
		if listener := l.OnThreadHide; listener != nil {
			listener(e)
		}

	// ThreadMember Events
	case *ThreadMemberAdd:
		if listener := l.OnThreadMemberAdd; listener != nil {
			listener(e)
		}
	case *ThreadMemberUpdate:
		if listener := l.OnThreadMemberUpdate; listener != nil {
			listener(e)
		}
	case *ThreadMemberRemove:
		if listener := l.OnThreadMemberRemove; listener != nil {
			listener(e)
		}

	// GuildChannel Events
	case *GuildChannelCreate:
		if listener := l.OnGuildChannelCreate; listener != nil {
			listener(e)
		}
	case *GuildChannelUpdate:
		if listener := l.OnGuildChannelUpdate; listener != nil {
			listener(e)
		}
	case *GuildChannelDelete:
		if listener := l.OnGuildChannelDelete; listener != nil {
			listener(e)
		}
	case *GuildChannelPinsUpdate:
		if listener := l.OnGuildChannelPinsUpdate; listener != nil {
			listener(e)
		}

	// DMChannel Events
	case *DMChannelPinsUpdate:
		if listener := l.OnDMChannelPinsUpdate; listener != nil {
			listener(e)
		}

	// DMChannel Message Events
	case *DMMessageCreate:
		if listener := l.OnDMMessageCreate; listener != nil {
			listener(e)
		}
	case *DMMessageUpdate:
		if listener := l.OnDMMessageUpdate; listener != nil {
			listener(e)
		}
	case *DMMessageDelete:
		if listener := l.OnDMMessageDelete; listener != nil {
			listener(e)
		}

	// DMChannel Events// Category Events
	case *DMMessageReactionAdd:
		if listener := l.OnDMMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemove:
		if listener := l.OnDMMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemoveEmoji:
		if listener := l.OnDMMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemoveAll:
		if listener := l.OnDMMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Emoji Events
	case *EmojisUpdate:
		if listener := l.OnEmojisUpdate; listener != nil {
			listener(e)
		}
	case *EmojiCreate:
		if listener := l.OnEmojiCreate; listener != nil {
			listener(e)
		}
	case *EmojiUpdate:
		if listener := l.OnEmojiUpdate; listener != nil {
			listener(e)
		}
	case *EmojiDelete:
		if listener := l.OnEmojiDelete; listener != nil {
			listener(e)
		}

	// Entitlement Events
	case *EntitlementCreate:
		if listener := l.OnEntitlementCreate; listener != nil {
			listener(e)
		}
	case *EntitlementUpdate:
		if listener := l.OnEntitlementUpdate; listener != nil {
			listener(e)
		}
	case *EntitlementDelete:
		if listener := l.OnEntitlementDelete; listener != nil {
			listener(e)
		}

	// Subscription Events
	case *SubscriptionCreate:
		if listener := l.OnSubscriptionCreate; listener != nil {
			listener(e)
		}
	case *SubscriptionUpdate:
		if listener := l.OnSubscriptionUpdate; listener != nil {
			listener(e)
		}
	case *SubscriptionDelete:
		if listener := l.OnSubscriptionDelete; listener != nil {
			listener(e)
		}

	// Sticker Events
	case *StickersUpdate:
		if listener := l.OnStickersUpdate; listener != nil {
			listener(e)
		}
	case *StickerCreate:
		if listener := l.OnStickerCreate; listener != nil {
			listener(e)
		}
	case *StickerUpdate:
		if listener := l.OnStickerUpdate; listener != nil {
			listener(e)
		}
	case *StickerDelete:
		if listener := l.OnStickerDelete; listener != nil {
			listener(e)
		}

	// gateway Status Events
	case *Ready:
		if listener := l.OnReady; listener != nil {
			listener(e)
		}
	case *Resumed:
		if listener := l.OnResumed; listener != nil {
			listener(e)
		}

	// Guild Events
	case *GuildJoin:
		if listener := l.OnGuildJoin; listener != nil {
			listener(e)
		}
	case *GuildUpdate:
		if listener := l.OnGuildUpdate; listener != nil {
			listener(e)
		}
	case *GuildLeave:
		if listener := l.OnGuildLeave; listener != nil {
			listener(e)
		}
	case *GuildAvailable:
		if listener := l.OnGuildAvailable; listener != nil {
			listener(e)
		}
	case *GuildUnavailable:
		if listener := l.OnGuildUnavailable; listener != nil {
			listener(e)
		}
	case *GuildReady:
		if listener := l.OnGuildReady; listener != nil {
			listener(e)
		}
	case *GuildsReady:
		if listener := l.OnGuildsReady; listener != nil {
			listener(e)
		}
	case *GuildBan:
		if listener := l.OnGuildBan; listener != nil {
			listener(e)
		}
	case *GuildUnban:
		if listener := l.OnGuildUnban; listener != nil {
			listener(e)
		}
	case *GuildAuditLogEntryCreate:
		if listener := l.OnGuildAuditLogEntryCreate; listener != nil {
			listener(e)
		}

	// Guild Invite Events
	case *InviteCreate:
		if listener := l.OnGuildInviteCreate; listener != nil {
			listener(e)
		}
	case *InviteDelete:
		if listener := l.OnGuildInviteDelete; listener != nil {
			listener(e)
		}

	// Member Events
	case *GuildMemberJoin:
		if listener := l.OnGuildMemberJoin; listener != nil {
			listener(e)
		}
	case *GuildMemberUpdate:
		if listener := l.OnGuildMemberUpdate; listener != nil {
			listener(e)
		}
	case *GuildMemberLeave:
		if listener := l.OnGuildMemberLeave; listener != nil {
			listener(e)
		}

	// Guild Message Events
	case *GuildMessageCreate:
		if listener := l.OnGuildMessageCreate; listener != nil {
			listener(e)
		}
	case *GuildMessageUpdate:
		if listener := l.OnGuildMessageUpdate; listener != nil {
			listener(e)
		}
	case *GuildMessageDelete:
		if listener := l.OnGuildMessageDelete; listener != nil {
			listener(e)
		}

	// Guild Message Reaction Events
	case *GuildMessageReactionAdd:
		if listener := l.OnGuildMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemove:
		if listener := l.OnGuildMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemoveEmoji:
		if listener := l.OnGuildMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemoveAll:
		if listener := l.OnGuildMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Guild Soundboard Sound Events
	case *GuildSoundboardSoundCreate:
		if listener := l.OnGuildSoundboardSoundCreate; listener != nil {
			listener(e)
		}
	case *GuildSoundboardSoundUpdate:
		if listener := l.OnGuildSoundboardSoundUpdate; listener != nil {
			listener(e)
		}
	case *GuildSoundboardSoundDelete:
		if listener := l.OnGuildSoundboardSoundDelete; listener != nil {
			listener(e)
		}
	case *GuildSoundboardSoundsUpdate:
		if listener := l.OnGuildSoundboardSoundsUpdate; listener != nil {
			listener(e)
		}
	case *SoundboardSounds:
		if listener := l.OnSoundboardSounds; listener != nil {
			listener(e)
		}

	// Guild Voice Events
	case *VoiceServerUpdate:
		if listener := l.OnVoiceServerUpdate; listener != nil {
			listener(e)
		}
	case *GuildVoiceChannelEffectSend:
		if listener := l.OnGuildVoiceChannelEffectSend; listener != nil {
			listener(e)
		}
	case *GuildVoiceStateUpdate:
		if listener := l.OnGuildVoiceStateUpdate; listener != nil {
			listener(e)
		}
	case *GuildVoiceJoin:
		if listener := l.OnGuildVoiceJoin; listener != nil {
			listener(e)
		}
	case *GuildVoiceMove:
		if listener := l.OnGuildVoiceMove; listener != nil {
			listener(e)
		}
	case *GuildVoiceLeave:
		if listener := l.OnGuildVoiceLeave; listener != nil {
			listener(e)
		}

	// Guild StageInstance Events
	case *StageInstanceCreate:
		if listener := l.OnStageInstanceCreate; listener != nil {
			listener(e)
		}
	case *StageInstanceUpdate:
		if listener := l.OnStageInstanceUpdate; listener != nil {
			listener(e)
		}
	case *StageInstanceDelete:
		if listener := l.OnStageInstanceDelete; listener != nil {
			listener(e)
		}

	// Guild Role Events
	case *RoleCreate:
		if listener := l.OnRoleCreate; listener != nil {
			listener(e)
		}
	case *RoleUpdate:
		if listener := l.OnRoleUpdate; listener != nil {
			listener(e)
		}
	case *RoleDelete:
		if listener := l.OnRoleDelete; listener != nil {
			listener(e)
		}

	// Guild ScheduledEvents
	case *GuildScheduledEventCreate:
		if listener := l.OnGuildScheduledEventCreate; listener != nil {
			listener(e)
		}
	case *GuildScheduledEventUpdate:
		if listener := l.OnGuildScheduledEventUpdate; listener != nil {
			listener(e)
		}
	case *GuildScheduledEventDelete:
		if listener := l.OnGuildScheduledEventDelete; listener != nil {
			listener(e)
		}
	case *GuildScheduledEventUserAdd:
		if listener := l.OnGuildScheduledEventUserAdd; listener != nil {
			listener(e)
		}
	case *GuildScheduledEventUserRemove:
		if listener := l.OnGuildScheduledEventUserRemove; listener != nil {
			listener(e)
		}

	// Interaction Events
	case *InteractionCreate:
		if listener := l.OnInteraction; listener != nil {
			listener(e)
		}
	case *ApplicationCommandInteractionCreate:
		if listener := l.OnApplicationCommandInteraction; listener != nil {
			listener(e)
		}
	case *ComponentInteractionCreate:
		if listener := l.OnComponentInteraction; listener != nil {
			listener(e)
		}
	case *AutocompleteInteractionCreate:
		if listener := l.OnAutocompleteInteraction; listener != nil {
			listener(e)
		}
	case *ModalSubmitInteractionCreate:
		if listener := l.OnModalSubmit; listener != nil {
			listener(e)
		}

	// Message Events
	case *MessageCreate:
		if listener := l.OnMessageCreate; listener != nil {
			listener(e)
		}
	case *MessageUpdate:
		if listener := l.OnMessageUpdate; listener != nil {
			listener(e)
		}
	case *MessageDelete:
		if listener := l.OnMessageDelete; listener != nil {
			listener(e)
		}

	// Message Poll Events
	case *MessagePollVoteAdd:
		if listener := l.OnMessagePollVoteAdd; listener != nil {
			listener(e)
		}
	case *MessagePollVoteRemove:
		if listener := l.OnMessagePollVoteRemove; listener != nil {
			listener(e)
		}
	case *DMMessagePollVoteAdd:
		if listener := l.OnDMMessagePollVoteAdd; listener != nil {
			listener(e)
		}
	case *DMMessagePollVoteRemove:
		if listener := l.OnDMMessagePollVoteRemove; listener != nil {
			listener(e)
		}
	case *GuildMessagePollVoteAdd:
		if listener := l.OnGuildMessagePollVoteAdd; listener != nil {
			listener(e)
		}
	case *GuildMessagePollVoteRemove:
		if listener := l.OnGuildMessagePollVoteRemove; listener != nil {
			listener(e)
		}

	// Message Reaction Events
	case *MessageReactionAdd:
		if listener := l.OnMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *MessageReactionRemove:
		if listener := l.OnMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *MessageReactionRemoveEmoji:
		if listener := l.OnMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *MessageReactionRemoveAll:
		if listener := l.OnMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Self Events
	case *SelfUpdate:
		if listener := l.OnSelfUpdate; listener != nil {
			listener(e)
		}

	// User Events
	case *UserUpdate:
		if listener := l.OnUserUpdate; listener != nil {
			listener(e)
		}
	case *UserTypingStart:
		if listener := l.OnUserTypingStart; listener != nil {
			listener(e)
		}
	case *GuildMemberTypingStart:
		if listener := l.OnGuildMemberTypingStart; listener != nil {
			listener(e)
		}
	case *DMUserTypingStart:
		if listener := l.OnDMUserTypingStart; listener != nil {
			listener(e)
		}

	case *PresenceUpdate:
		if listener := l.OnPresenceUpdate; listener != nil {
			listener(e)
		}

	// User Activity Events
	case *UserActivityStart:
		if listener := l.OnUserActivityStart; listener != nil {
			listener(e)
		}
	case *UserActivityUpdate:
		if listener := l.OnUserActivityUpdate; listener != nil {
			listener(e)
		}
	case *UserActivityStop:
		if listener := l.OnUserActivityStop; listener != nil {
			listener(e)
		}

	// User Status Events
	case *UserStatusUpdate:
		if listener := l.OnUserStatusUpdate; listener != nil {
			listener(e)
		}
	case *UserClientStatusUpdate:
		if listener := l.OnUserClientStatusUpdate; listener != nil {
			listener(e)
		}

	// Integration Events
	case *IntegrationCreate:
		if listener := l.OnIntegrationCreate; listener != nil {
			listener(e)
		}
	case *IntegrationUpdate:
		if listener := l.OnIntegrationUpdate; listener != nil {
			listener(e)
		}
	case *IntegrationDelete:
		if listener := l.OnIntegrationDelete; listener != nil {
			listener(e)
		}
	case *GuildIntegrationsUpdate:
		if listener := l.OnGuildIntegrationsUpdate; listener != nil {
			listener(e)
		}

	case *WebhooksUpdate:
		if listener := l.OnGuildWebhooksUpdate; listener != nil {
			listener(e)
		}

	default:
		e.Client().Logger().Error("unexpected event received", slog.String("type", fmt.Sprintf("%T", event)), slog.String("data", fmt.Sprintf("%+v", event)))
	}
}
