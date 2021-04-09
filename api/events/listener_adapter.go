package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// ListenerAdapter lets you override the handles for receiving events
type ListenerAdapter struct {
	OnGenericApplicationCommandEvent func(event *GenericApplicationCommandEvent)
	OnApplicationCommandCreate       func(event *ApplicationCommandCreateEvent)
	OnApplicationCommandUpdate       func(event *ApplicationCommandUpdateEvent)
	OnApplicationCommandDelete       func(event *ApplicationCommandDeleteEvent)

	OnGenericCategoryEvent func(event *GenericCategoryEvent)
	OnCategoryCreate       func(event *CategoryCreateEvent)
	OnCategoryUpdate       func(event *CategoryUpdateEvent)
	OnCategoryDelete       func(event *CategoryDeleteEvent)

	OnGenericChannelEvent func(event *GenericChannelEvent)

	OnGenericDMChannelEvent func(event *GenericDMChannelEvent)
	OnDMChannelCreate       func(event *DMChannelCreateEvent)
	OnDMChannelUpdate       func(event *DMChannelUpdateEvent)
	OnDMChannelDelete       func(event *DMChannelDeleteEvent)

	OnGenericDMMessageReactionEventEvent func(event *GenericDMMessageReactionEvent)
	OnDMMessageReactionAddEvent          func(event *DMMessageReactionAddEvent)
	OnDMMessageReactionRemove            func(event *DMMessageReactionRemoveEvent)
	OnDMMessageReactionRemoveEmote       func(event *DMMessageReactionRemoveEmoteEvent)
	OnDMMessageReactionRemoveAll         func(event *DMMessageReactionRemoveAllEvent)

	OnGenericEmoteEvent func(event *GenericEmoteEvent)
	OnEmoteCreate       func(event *EmoteCreateEvent)
	OnEmoteUpdate       func(event *EmoteUpdateEvent)
	OnEmoteDelete       func(event *EmoteDeleteEvent)

	OnException func(event *ExceptionEvent)

	OnGenericGatewayStatusEvent func(event *GenericGatewayStatusEvent)
	OnConnected                 func(event *ConnectedEvent)
	OnReconnected               func(event *ReconnectedEvent)
	OnResumed                   func(event *ResumedEvent)
	OnDisconnected              func(event *DisconnectedEvent)
	OnShutdown                  func(event *ShutdownEvent)

	OnGenericEvent func(event api.Event)

	// Guild Events
	OnGenericGuildEvent func(event *GenericGuildEvent)
	OnGuildUpdate       func(event *GuildUpdateEvent)
	OnGuildAvailable    func(event *GuildAvailableEvent)
	OnGuildUnavailable  func(event *GuildUnavailableEvent)
	OnGuildJoin         func(event *GuildJoinEvent)
	OnGuildLeave        func(event *GuildLeaveEvent)
	OnGuildReady        func(event *GuildReadyEvent)
	OnGuildBan          func(event *GuildBanEvent)
	OnGuildUnban        func(event *GuildUnbanEvent)

	OnGenericGuildInviteEvent func(event *GenericGuildInviteEvent)
	OnGuildInviteCreate       func(event *GuildInviteCreateEvent)
	OnGuildInviteDelete       func(event *GuildInviteDeleteEvent)

	// Member Events
	OnGenericGuildMemberEvent func(event *GenericGuildMemberEvent)
	OnGuildMemberJoin         func(event *GuildMemberJoinEvent)
	OnGuildMemberUpdate       func(event *GuildMemberUpdateEvent)
	OnGuildMemberLeave        func(event *GuildMemberLeaveEvent)

	OnGenericGuildMessageEvent func(event *GenericGuildMessageEvent)
	OnGuildMessageReceived     func(event *GuildMessageReceivedEvent)
	OnGuildMessageUpdate       func(event *GuildMessageUpdateEvent)
	OnGuildMessageDelete       func(event *GuildMessageDeleteEvent)

	OnGenericGuildMessageReactionEvent func(event *GenericGuildMessageReactionEvent)
	OnGuildMessageReactionAdd          func(event *GuildMessageReactionAddEvent)
	OnGuildMessageReactionRemove       func(event *GuildMessageReactionRemoveEvent)
	OnGuildMessageReactionRemoveEmote  func(event *GuildMessageReactionRemoveEmoteEvent)
	OnGuildMessageReactionRemoveAll    func(event *GuildMessageReactionRemoveAllEvent)

	OnGenericGuildVoiceEvent func(event *GenericGuildVoiceEvent)
	OnGuildVoiceUpdate       func(event *GuildVoiceUpdateEvent)
	OnGuildVoiceJoin         func(event *GuildVoiceJoinEvent)
	OnGuildVoiceLeave        func(event *GuildVoiceLeaveEvent)

	OnHeartbeat func(event *HeartbeatEvent)

	OnHttpRequest func(event *HttpRequestEvent)

	// Interaction Events
	OnGenericInteractionEvent func(event *GenericInteractionEvent)
	OnSlashCommand            func(event *SlashCommandEvent)

	OnGenericMessageEvent func(event *GenericMessageEvent)
	OnMessageDelete       func(event *MessageDeleteEvent)
	OnMessageReceived     func(event *MessageReceivedEvent)
	OnMessageUpdate       func(event *MessageUpdateEvent)

	OnGenericReactionEvent       func(event *GenericReactionEvents)
	OnMessageReactionAdd         func(event *MessageReactionAddEvent)
	OnMessageReactionRemove      func(event *MessageReactionRemoveEvent)
	OnMessageReactionRemoveEmote func(event *MessageReactionRemoveEmoteEvent)
	OnMessageReactionRemoveAll   func(event *MessageReactionRemoveAllEvent)

	OnRawGateway func(event *RawGatewayEvent)

	OnReadyEvent func(event *ReadyEvent)

	// Guild Role Events
	OnGenericRoleEvent func(event *GenericRoleEvent)
	OnRoleCreate       func(event *RoleCreateEvent)
	OnRoleUpdate       func(event *RoleUpdateEvent)
	OnRoleDelete       func(event *RoleDeleteEvent)

	OnSelfUpdate func(event *SelfUpdateEvent)

	OnGenericStoreChannelEvent func(event *StoreChannelCreateEvent)
	OnStoreChannelCreate       func(event *StoreChannelCreateEvent)
	OnStoreChannelUpdate       func(event *StoreChannelUpdateEvent)
	OnStoreChannelDelete       func(event *StoreChannelDeleteEvent)

	OnGenericTextChannelEvent func(event *GenericTextChannelEvent)
	OnTextChannelCreate       func(event *TextChannelCreateEvent)
	OnTextChannelUpdate       func(event *TextChannelUpdateEvent)
	OnTextChannelDelete       func(event *TextChannelDeleteEvent)

	OnGenericUserActivityEvent func(event *GenericUserActivityEvent)
	OnUserActivityStart        func(event *UserActivityStartEvent)
	OnUserActivityUpdate       func(event *UserActivityUpdateEvent)
	OnUserActivityEnd          func(event *UserActivityEndEvent)

	OnGenericUserEvent func(event *GenericUserEvent)
	OnUserUpdate       func(event *UserUpdateEvent)
	OnUserTyping       func(event *UserTypingEvent)
	OnGuildUserTyping  func(event *GuildUserTypingEvent)
	OnDMUserTyping     func(event *DMUserTypingEvent)

	OnGenericVoiceChannelEvent func(event *GenericVoiceChannelEvent)
	OnVoiceChannelCreate       func(event *VoiceChannelCreateEvent)
	OnVoiceChannelUpdate       func(event *VoiceChannelUpdateEvent)
	OnVoiceChannelDelete       func(event *VoiceChannelDeleteEvent)
}

// OnEvent is getting called everytime we receive an event
func (l ListenerAdapter) OnEvent(event interface{}) {
	if e, ok := event.(api.Event); ok {
		if l.OnGenericEvent != nil {
			l.OnGenericEvent(e)
		}
	}
	switch e := event.(type) {
	// Guild Events
	case GenericGuildEvent:
		if l.OnGenericGuildEvent != nil {
			l.OnGenericGuildEvent(&e)
		}
	case GuildJoinEvent:
		if l.OnGuildJoin != nil {
			l.OnGuildJoin(&e)
		}
	case GuildUpdateEvent:
		if l.OnGuildUpdate != nil {
			l.OnGuildUpdate(&e)
		}
	case GuildLeaveEvent:
		if l.OnGuildLeave != nil {
			l.OnGuildLeave(&e)
		}
	case GuildAvailableEvent:
		if l.OnGuildAvailable != nil {
			l.OnGuildAvailable(&e)
		}
	case GuildUnavailableEvent:
		if l.OnGuildUnavailable != nil {
			l.OnGuildUnavailable(&e)
		}

	// Guild Role Events
	case GenericRoleEvent:
		if l.OnGenericRoleEvent != nil {
			l.OnGenericRoleEvent(&e)
		}
	case RoleCreateEvent:
		if l.OnRoleCreate != nil {
			l.OnRoleCreate(&e)
		}
	case RoleUpdateEvent:
		if l.OnRoleUpdate != nil {
			l.OnRoleUpdate(&e)
		}
	case RoleDeleteEvent:
		if l.OnRoleDelete != nil {
			l.OnRoleDelete(&e)
		}

	// Message Events
	case MessageReceivedEvent:
		if l.OnMessageReceived != nil {
			l.OnMessageReceived(&e)
		}
	case GuildMessageReceivedEvent:
		if l.OnGuildMessageReceived != nil {
			l.OnGuildMessageReceived(&e)
		}

	// Interaction Events
	case GenericInteractionEvent:
		if l.OnGenericInteractionEvent != nil {
			l.OnGenericInteractionEvent(&e)
		}
	case SlashCommandEvent:
		if l.OnSlashCommand != nil {
			l.OnSlashCommand(&e)
		}
	default:
		//log.Errorf("unexpected event received: %#e", event)
	}
}
