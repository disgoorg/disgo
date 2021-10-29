package events

import (
	"reflect"

	"github.com/DisgoOrg/disgo/core"
)

// ListenerAdapter lets you override the handles for receiving events
type ListenerAdapter struct {
	// Other events
	OnHeartbeat   func(event *HeartbeatEvent)
	OnHTTPRequest func(event *HTTPRequestEvent)
	OnRaw         func(event *RawEvent)

	// Channel Events
	OnGuildChannelCreate func(event *GuildChannelCreateEvent)
	OnGuildChannelUpdate func(event *GuildChannelUpdateEvent)
	OnGuildChannelDelete func(event *GuildChannelDeleteEvent)

	// Channel Events
	OnDMChannelCreate func(event *DMChannelCreateEvent)
	OnDMChannelUpdate func(event *DMChannelUpdateEvent)
	OnDMChannelDelete func(event *DMChannelDeleteEvent)

	// Channel Message Events
	OnDMMessageCreate func(event *DMMessageCreateEvent)
	OnDMMessageUpdate func(event *DMMessageUpdateEvent)
	OnDMMessageDelete func(event *DMMessageDeleteEvent)

	// Channel Reaction Events
	OnDMMessageReactionAdd         func(event *DMMessageReactionAddEvent)
	OnDMMessageReactionRemove      func(event *DMMessageReactionRemoveEvent)
	OnDMMessageReactionRemoveEmoji func(event *DMMessageReactionRemoveEmojiEvent)
	OnDMMessageReactionRemoveAll   func(event *DMMessageReactionRemoveAllEvent)

	// Emoji Events
	OnEmojiCreate func(event *EmojiCreateEvent)
	OnEmojiUpdate func(event *EmojiUpdateEvent)
	OnEmojiDelete func(event *EmojiDeleteEvent)

	// Sticker Events
	OnStickerCreate func(event *StickerCreateEvent)
	OnStickerUpdate func(event *StickerUpdateEvent)
	OnStickerDelete func(event *StickerDeleteEvent)

	// gateway status Events
	OnReady          func(event *ReadyEvent)
	OnResumed        func(event *ResumedEvent)
	OnInvalidSession func(event *InvalidSessionEvent)
	OnDisconnected   func(event *DisconnectedEvent)

	// Guild Events
	OnGuildJoin        func(event *GuildJoinEvent)
	OnGuildUpdate      func(event *GuildUpdateEvent)
	OnGuildLeave       func(event *GuildLeaveEvent)
	OnGuildAvailable   func(event *GuildAvailableEvent)
	OnGuildUnavailable func(event *GuildUnavailableEvent)
	OnGuildReady       func(event *GuildReadyEvent)
	OnGuildsReady      func(event *GuildsReadyEvent)
	OnGuildBan         func(event *GuildBanEvent)
	OnGuildUnban       func(event *GuildUnbanEvent)

	// Guild Invite Events
	OnGuildInviteCreate func(event *GuildInviteCreateEvent)
	OnGuildInviteDelete func(event *GuildInviteDeleteEvent)

	// Guild Member Events
	OnGuildMemberJoin   func(event *GuildMemberJoinEvent)
	OnGuildMemberUpdate func(event *GuildMemberUpdateEvent)
	OnGuildMemberLeave  func(event *GuildMemberLeaveEvent)

	// Guild Message Events
	OnGuildMessageCreate func(event *GuildMessageCreateEvent)
	OnGuildMessageUpdate func(event *GuildMessageUpdateEvent)
	OnGuildMessageDelete func(event *GuildMessageDeleteEvent)

	// Guild Message Reaction Events
	OnGuildMessageReactionAdd         func(event *GuildMessageReactionAddEvent)
	OnGuildMessageReactionRemove      func(event *GuildMessageReactionRemoveEvent)
	OnGuildMessageReactionRemoveEmoji func(event *GuildMessageReactionRemoveEmojiEvent)
	OnGuildMessageReactionRemoveAll   func(event *GuildMessageReactionRemoveAllEvent)

	// Guild Voice Events
	OnVoiceServerUpdate     func(event *VoiceServerUpdateEvent)
	OnGuildVoiceStateUpdate func(event *GuildVoiceStateUpdateEvent)
	OnGuildVoiceJoin        func(event *GuildVoiceJoinEvent)
	OnGuildVoiceMove        func(event *GuildVoiceMoveEvent)
	OnGuildVoiceLeave       func(event *GuildVoiceLeaveEvent)

	// Guild StageInstance Events
	OnStageInstanceCreate func(event *StageInstanceCreateEvent)
	OnStageInstanceUpdate func(event *StageInstanceUpdateEvent)
	OnStageInstanceDelete func(event *StageInstanceDeleteEvent)

	// Guild Role Events
	OnRoleCreate func(event *RoleCreateEvent)
	OnRoleUpdate func(event *RoleUpdateEvent)
	OnRoleDelete func(event *RoleDeleteEvent)

	// Interaction Events
	OnInteractionCreate                   func(event *InteractionCreateEvent)
	OnApplicationCommandInteractionCreate func(event *ApplicationCommandInteractionCreateEvent)
	OnSlashCommand                        func(event *SlashCommandEvent)
	OnUserCommand                         func(event *UserCommandEvent)
	OnMessageCommand                      func(event *MessageCommandEvent)
	OnComponentInteractionCreate          func(event *ComponentInteractionCreateEvent)
	OnButtonClick                         func(event *ButtonClickEvent)
	OnSelectMenuSubmit                    func(event *SelectMenuSubmitEvent)
	OnAutocomplete                        func(event *AutocompleteEvent)

	// Message Events
	OnMessageCreate func(event *MessageCreateEvent)
	OnMessageUpdate func(event *MessageUpdateEvent)
	OnMessageDelete func(event *MessageDeleteEvent)

	// Message Reaction Events
	OnMessageReactionAdd         func(event *MessageReactionAddEvent)
	OnMessageReactionRemove      func(event *MessageReactionRemoveEvent)
	OnMessageReactionRemoveEmoji func(event *MessageReactionRemoveEmojiEvent)
	OnMessageReactionRemoveAll   func(event *MessageReactionRemoveAllEvent)

	// Self Events
	OnSelfUpdate func(event *SelfUpdateEvent)

	// User Events
	OnUserUpdate      func(event *UserUpdateEvent)
	OnUserTyping      func(event *UserTypingEvent)
	OnGuildUserTyping func(event *GuildMemberTypingEvent)
	OnDMUserTyping    func(event *DMChannelUserTypingEvent)

	// User Activity Events
	OnUserActivityStart  func(event *UserActivityStartEvent)
	OnUserActivityUpdate func(event *UserActivityUpdateEvent)
	OnUserActivityStop   func(event *UserActivityStopEvent)

	OnUserStatusUpdate       func(event *UserStatusUpdateEvent)
	OnUserClientStatusUpdate func(event *UserClientStatusUpdateEvent)

	OnIntegrationCreate       func(event *IntegrationCreateEvent)
	OnIntegrationUpdate       func(event *IntegrationUpdateEvent)
	OnIntegrationDelete       func(event *IntegrationDeleteEvent)
	OnGuildIntegrationsUpdate func(event *GuildIntegrationsUpdateEvent)

	OnGuildWebhooksUpdate func(event *WebhooksUpdateEvent)
}

// OnEvent is getting called everytime we receive an event
func (l ListenerAdapter) OnEvent(event core.Event) {
	switch e := event.(type) {
	case *HeartbeatEvent:
		if listener := l.OnHeartbeat; listener != nil {
			listener(e)
		}
	case *HTTPRequestEvent:
		if listener := l.OnHTTPRequest; listener != nil {
			listener(e)
		}
	case *RawEvent:
		if listener := l.OnRaw; listener != nil {
			listener(e)
		}

	// GetGuildChannel Events
	case *GuildChannelCreateEvent:
		if listener := l.OnGuildChannelCreate; listener != nil {
			listener(e)
		}
	case *GuildChannelUpdateEvent:
		if listener := l.OnGuildChannelUpdate; listener != nil {
			listener(e)
		}
	case *GuildChannelDeleteEvent:
		if listener := l.OnGuildChannelDelete; listener != nil {
			listener(e)
		}

	// DMChannel Events
	case *DMChannelCreateEvent:
		if listener := l.OnDMChannelCreate; listener != nil {
			listener(e)
		}
	case *DMChannelUpdateEvent:
		if listener := l.OnDMChannelUpdate; listener != nil {
			listener(e)
		}
	case *DMChannelDeleteEvent:
		if listener := l.OnDMChannelDelete; listener != nil {
			listener(e)
		}

	// DMChannel Message Events
	case *DMMessageCreateEvent:
		if listener := l.OnDMMessageCreate; listener != nil {
			listener(e)
		}
	case *DMMessageUpdateEvent:
		if listener := l.OnDMMessageUpdate; listener != nil {
			listener(e)
		}
	case *DMMessageDeleteEvent:
		if listener := l.OnDMMessageDelete; listener != nil {
			listener(e)
		}

	// DMChannel Events// Category Events
	case *DMMessageReactionAddEvent:
		if listener := l.OnDMMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemoveEvent:
		if listener := l.OnDMMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemoveEmojiEvent:
		if listener := l.OnDMMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *DMMessageReactionRemoveAllEvent:
		if listener := l.OnDMMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Emoji Events
	case *EmojiCreateEvent:
		if listener := l.OnEmojiCreate; listener != nil {
			listener(e)
		}
	case *EmojiUpdateEvent:
		if listener := l.OnEmojiUpdate; listener != nil {
			listener(e)
		}
	case *EmojiDeleteEvent:
		if listener := l.OnEmojiDelete; listener != nil {
			listener(e)
		}

	// Sticker Events
	case *StickerCreateEvent:
		if listener := l.OnStickerCreate; listener != nil {
			listener(e)
		}
	case *StickerUpdateEvent:
		if listener := l.OnStickerUpdate; listener != nil {
			listener(e)
		}
	case *StickerDeleteEvent:
		if listener := l.OnStickerDelete; listener != nil {
			listener(e)
		}

	// gateway Status Events
	case *ReadyEvent:
		if listener := l.OnReady; listener != nil {
			listener(e)
		}
	case *ResumedEvent:
		if listener := l.OnResumed; listener != nil {
			listener(e)
		}
	case *InvalidSessionEvent:
		if listener := l.OnInvalidSession; listener != nil {
			listener(e)
		}
	case *DisconnectedEvent:
		if listener := l.OnDisconnected; listener != nil {
			listener(e)
		}

	// Guild Events
	case *GuildJoinEvent:
		if listener := l.OnGuildJoin; listener != nil {
			listener(e)
		}
	case *GuildUpdateEvent:
		if listener := l.OnGuildUpdate; listener != nil {
			listener(e)
		}
	case *GuildLeaveEvent:
		if listener := l.OnGuildLeave; listener != nil {
			listener(e)
		}
	case *GuildAvailableEvent:
		if listener := l.OnGuildAvailable; listener != nil {
			listener(e)
		}
	case *GuildUnavailableEvent:
		if listener := l.OnGuildUnavailable; listener != nil {
			listener(e)
		}
	case *GuildReadyEvent:
		if listener := l.OnGuildReady; listener != nil {
			listener(e)
		}
	case *GuildsReadyEvent:
		if listener := l.OnGuildsReady; listener != nil {
			listener(e)
		}
	case *GuildBanEvent:
		if listener := l.OnGuildBan; listener != nil {
			listener(e)
		}
	case *GuildUnbanEvent:
		if listener := l.OnGuildUnban; listener != nil {
			listener(e)
		}

	// Guild Invite Events
	case *GuildInviteCreateEvent:
		if listener := l.OnGuildInviteCreate; listener != nil {
			listener(e)
		}
	case *GuildInviteDeleteEvent:
		if listener := l.OnGuildInviteDelete; listener != nil {
			listener(e)
		}

	// Member Events
	case *GuildMemberJoinEvent:
		if listener := l.OnGuildMemberJoin; listener != nil {
			listener(e)
		}
	case *GuildMemberUpdateEvent:
		if listener := l.OnGuildMemberUpdate; listener != nil {
			listener(e)
		}
	case *GuildMemberLeaveEvent:
		if listener := l.OnGuildMemberLeave; listener != nil {
			listener(e)
		}

	// Guild Message Events
	case *GuildMessageCreateEvent:
		if listener := l.OnGuildMessageCreate; listener != nil {
			listener(e)
		}
	case *GuildMessageUpdateEvent:
		if listener := l.OnGuildMessageUpdate; listener != nil {
			listener(e)
		}
	case *GuildMessageDeleteEvent:
		if listener := l.OnGuildMessageDelete; listener != nil {
			listener(e)
		}

	// Guild Message Reaction Events
	case *GuildMessageReactionAddEvent:
		if listener := l.OnGuildMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemoveEvent:
		if listener := l.OnGuildMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemoveEmojiEvent:
		if listener := l.OnGuildMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *GuildMessageReactionRemoveAllEvent:
		if listener := l.OnGuildMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Guild Voice Events
	case *VoiceServerUpdateEvent:
		if listener := l.OnVoiceServerUpdate; listener != nil {
			listener(e)
		}
	case *GuildVoiceStateUpdateEvent:
		if listener := l.OnGuildVoiceStateUpdate; listener != nil {
			listener(e)
		}
	case *GuildVoiceJoinEvent:
		if listener := l.OnGuildVoiceJoin; listener != nil {
			listener(e)
		}
	case *GuildVoiceMoveEvent:
		if listener := l.OnGuildVoiceMove; listener != nil {
			listener(e)
		}
	case *GuildVoiceLeaveEvent:
		if listener := l.OnGuildVoiceLeave; listener != nil {
			listener(e)
		}

	// Guild StageInstance Events
	case *StageInstanceCreateEvent:
		if listener := l.OnStageInstanceCreate; listener != nil {
			listener(e)
		}
	case *StageInstanceUpdateEvent:
		if listener := l.OnStageInstanceUpdate; listener != nil {
			listener(e)
		}
	case *StageInstanceDeleteEvent:
		if listener := l.OnStageInstanceDelete; listener != nil {
			listener(e)
		}

	// Guild Role Events
	case *RoleCreateEvent:
		if listener := l.OnRoleCreate; listener != nil {
			listener(e)
		}
	case *RoleUpdateEvent:
		if listener := l.OnRoleUpdate; listener != nil {
			listener(e)
		}
	case *RoleDeleteEvent:
		if listener := l.OnRoleDelete; listener != nil {
			listener(e)
		}

	// Interaction Events
	case *InteractionCreateEvent:
		if listener := l.OnInteractionCreate; listener != nil {
			listener(e)
		}
	case *ApplicationCommandInteractionCreateEvent:
		if listener := l.OnApplicationCommandInteractionCreate; listener != nil {
			listener(e)
		}
	case *SlashCommandEvent:
		if listener := l.OnSlashCommand; listener != nil {
			listener(e)
		}
	case *UserCommandEvent:
		if listener := l.OnUserCommand; listener != nil {
			listener(e)
		}
	case *MessageCommandEvent:
		if listener := l.OnMessageCommand; listener != nil {
			listener(e)
		}
	case *ComponentInteractionCreateEvent:
		if listener := l.OnComponentInteractionCreate; listener != nil {
			listener(e)
		}
	case *ButtonClickEvent:
		if listener := l.OnButtonClick; listener != nil {
			listener(e)
		}
	case *SelectMenuSubmitEvent:
		if listener := l.OnSelectMenuSubmit; listener != nil {
			listener(e)
		}
	case *AutocompleteEvent:
		if listener := l.OnAutocomplete; listener != nil {
			listener(e)
		}

	// Message Events
	case *MessageCreateEvent:
		if listener := l.OnMessageCreate; listener != nil {
			listener(e)
		}
	case *MessageUpdateEvent:
		if listener := l.OnMessageUpdate; listener != nil {
			listener(e)
		}
	case *MessageDeleteEvent:
		if listener := l.OnMessageDelete; listener != nil {
			listener(e)
		}

	// Message Reaction Events
	case *MessageReactionAddEvent:
		if listener := l.OnMessageReactionAdd; listener != nil {
			listener(e)
		}
	case *MessageReactionRemoveEvent:
		if listener := l.OnMessageReactionRemove; listener != nil {
			listener(e)
		}
	case *MessageReactionRemoveEmojiEvent:
		if listener := l.OnMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case *MessageReactionRemoveAllEvent:
		if listener := l.OnMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Self Events
	case *SelfUpdateEvent:
		if listener := l.OnSelfUpdate; listener != nil {
			listener(e)
		}

	// User Events
	case *UserUpdateEvent:
		if listener := l.OnUserUpdate; listener != nil {
			listener(e)
		}
	case *UserTypingEvent:
		if listener := l.OnUserTyping; listener != nil {
			listener(e)
		}
	case *GuildMemberTypingEvent:
		if listener := l.OnGuildUserTyping; listener != nil {
			listener(e)
		}
	case *DMChannelUserTypingEvent:
		if listener := l.OnDMUserTyping; listener != nil {
			listener(e)
		}

	// User Activity Events
	case *UserActivityStartEvent:
		if listener := l.OnUserActivityStart; listener != nil {
			listener(e)
		}
	case *UserActivityUpdateEvent:
		if listener := l.OnUserActivityUpdate; listener != nil {
			listener(e)
		}
	case *UserActivityStopEvent:
		if listener := l.OnUserActivityStop; listener != nil {
			listener(e)
		}

	// User Status Events
	case *UserStatusUpdateEvent:
		if listener := l.OnUserStatusUpdate; listener != nil {
			listener(e)
		}
	case *UserClientStatusUpdateEvent:
		if listener := l.OnUserClientStatusUpdate; listener != nil {
			listener(e)
		}

	// Integration Events
	case *IntegrationCreateEvent:
		if listener := l.OnIntegrationCreate; listener != nil {
			listener(e)
		}
	case *IntegrationUpdateEvent:
		if listener := l.OnIntegrationUpdate; listener != nil {
			listener(e)
		}
	case *IntegrationDeleteEvent:
		if listener := l.OnIntegrationDelete; listener != nil {
			listener(e)
		}
	case *GuildIntegrationsUpdateEvent:
		if listener := l.OnGuildIntegrationsUpdate; listener != nil {
			listener(e)
		}

	case *WebhooksUpdateEvent:
		if listener := l.OnGuildWebhooksUpdate; listener != nil {
			listener(e)
		}

	default:
		if e, ok := e.(core.Event); ok {
			var name string
			if t := reflect.TypeOf(e); t.Kind() == reflect.Ptr {
				name = "*" + t.Elem().Name()
			} else {
				name = t.Name()
			}
			e.Bot().Logger.Errorf("unexpected event received: \"%s\", event: \"%#e\"", name, event)
		}
	}
}
