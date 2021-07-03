package events

import (
	"reflect"

	"github.com/DisgoOrg/disgo/api"
)

// ListenerAdapter lets you override the handles for receiving events
type ListenerAdapter struct {
	// Other events
	OnGenericEvent func(event GenericEvent)
	OnHeartbeat    func(event HeartbeatEvent)
	OnHTTPRequest  func(event HTTPRequestEvent)
	OnRawGateway   func(event RawGatewayEvent)
	OnReadyEvent   func(event ReadyEvent)

	// api.Command Events
	OnGenericCommandEvent func(event GenericCommandEvent)
	OnCommandCreate       func(event CommandCreateEvent)
	OnCommandUpdate       func(event CommandUpdateEvent)
	OnCommandDelete       func(event CommandDeleteEvent)

	// api.Channel Events
	OnGenericChannelEvent func(event GenericChannelEvent)

	// api.GuildChannel Events
	OnGenericGuildChannelEvent func(event GenericGuildChannelEvent)
	OnGuildChannelCreate       func(event GuildChannelCreateEvent)
	OnGuildChannelUpdate       func(event GuildChannelUpdateEvent)
	OnGuildChannelDelete       func(event GuildChannelDeleteEvent)

	// api.Category Events
	OnGenericCategoryEvent func(event GenericCategoryEvent)
	OnCategoryCreate       func(event CategoryCreateEvent)
	OnCategoryUpdate       func(event CategoryUpdateEvent)
	OnCategoryDelete       func(event CategoryDeleteEvent)

	// api.DMChannel Events
	OnGenericDMChannelEvent func(event GenericDMChannelEvent)
	OnDMChannelCreate       func(event DMChannelCreateEvent)
	OnDMChannelUpdate       func(event DMChannelUpdateEvent)
	OnDMChannelDelete       func(event DMChannelDeleteEvent)

	// api.DMChannel Reaction Events
	OnGenericDMMessageReactionEventEvent func(event GenericDMMessageReactionEvent)
	OnDMMessageReactionAdd               func(event DMMessageReactionAddEvent)
	OnDMMessageReactionRemove            func(event DMMessageReactionRemoveEvent)
	OnDMMessageReactionRemoveEmoji       func(event DMMessageReactionRemoveEmojiEvent)
	OnDMMessageReactionRemoveAll         func(event DMMessageReactionRemoveAllEvent)

	// api.StoreChannel Events
	OnGenericStoreChannelEvent func(event GenericStoreChannelEvent)
	OnStoreChannelCreate       func(event StoreChannelCreateEvent)
	OnStoreChannelUpdate       func(event StoreChannelUpdateEvent)
	OnStoreChannelDelete       func(event StoreChannelDeleteEvent)

	// api.TextChannel Events
	OnGenericTextChannelEvent func(event GenericTextChannelEvent)
	OnTextChannelCreate       func(event TextChannelCreateEvent)
	OnTextChannelUpdate       func(event TextChannelUpdateEvent)
	OnTextChannelDelete       func(event TextChannelDeleteEvent)

	// api.VoiceChannel Events
	OnGenericVoiceChannelEvent func(event GenericVoiceChannelEvent)
	OnVoiceChannelCreate       func(event VoiceChannelCreateEvent)
	OnVoiceChannelUpdate       func(event VoiceChannelUpdateEvent)
	OnVoiceChannelDelete       func(event VoiceChannelDeleteEvent)

	// api.Emoji Events
	OnGenericEmoteEvent func(event GenericEmoteEvent)
	OnEmoteCreate       func(event EmoteCreateEvent)
	OnEmoteUpdate       func(event EmoteUpdateEvent)
	OnEmoteDelete       func(event EmoteDeleteEvent)

	// api.GatewayStatus Events
	OnGenericGatewayStatusEvent func(event GenericGatewayStatusEvent)
	OnConnected                 func(event ConnectedEvent)
	OnReconnected               func(event ReconnectedEvent)
	OnResumed                   func(event ResumedEvent)
	OnDisconnected              func(event DisconnectedEvent)

	// api.Guild Events
	OnGenericGuildEvent func(event GenericGuildEvent)
	OnGuildJoin         func(event GuildJoinEvent)
	OnGuildUpdate       func(event GuildUpdateEvent)
	OnGuildLeave        func(event GuildLeaveEvent)
	OnGuildAvailable    func(event GuildAvailableEvent)
	OnGuildUnavailable  func(event GuildUnavailableEvent)
	OnGuildReady        func(event GuildReadyEvent)
	OnGuildBan          func(event GuildBanEvent)
	OnGuildUnban        func(event GuildUnbanEvent)

	// api.Guild api.Invite Events
	OnGenericGuildInviteEvent func(event GenericGuildInviteEvent)
	OnGuildInviteCreate       func(event GuildInviteCreateEvent)
	OnGuildInviteDelete       func(event GuildInviteDeleteEvent)

	// api.Guild api.Member Events
	OnGenericGuildMemberEvent func(event GenericGuildMemberEvent)
	OnGuildMemberJoin         func(event GuildMemberJoinEvent)
	OnGuildMemberUpdate       func(event GuildMemberUpdateEvent)
	OnGuildMemberLeave        func(event GuildMemberLeaveEvent)

	// api.Guild api.Message Events
	OnGenericGuildMessageEvent func(event GenericGuildMessageEvent)
	OnGuildMessageCreate       func(event GuildMessageCreateEvent)
	OnGuildMessageUpdate       func(event GuildMessageUpdateEvent)
	OnGuildMessageDelete       func(event GuildMessageDeleteEvent)

	// api.Guild api.Message Reaction Events
	OnGenericGuildMessageReactionEvent func(event GenericGuildMessageReactionEvent)
	OnGuildMessageReactionAdd          func(event GuildMessageReactionAddEvent)
	OnGuildMessageReactionRemove       func(event GuildMessageReactionRemoveEvent)
	OnGuildMessageReactionRemoveEmoji  func(event GuildMessageReactionRemoveEmojiEvent)
	OnGuildMessageReactionRemoveAll    func(event GuildMessageReactionRemoveAllEvent)

	// api.Guild Voice Events
	OnGenericGuildVoiceEvent func(event GenericGuildVoiceEvent)
	OnGuildVoiceUpdate       func(event GuildVoiceUpdateEvent)
	OnGuildVoiceJoin         func(event GuildVoiceJoinEvent)
	OnGuildVoiceLeave        func(event GuildVoiceLeaveEvent)

	// api.Guild api.Role Events
	OnGenericRoleEvent func(event GenericRoleEvent)
	OnRoleCreate       func(event RoleCreateEvent)
	OnRoleUpdate       func(event RoleUpdateEvent)
	OnRoleDelete       func(event RoleDeleteEvent)

	// api.Interaction Events
	OnGenericInteractionEvent func(event GenericInteractionEvent)
	OnCommand                 func(event CommandEvent)
	OnGenericComponentEvent   func(event GenericComponentEvent)
	OnButtonClick             func(event ButtonClickEvent)
	OnSelectMenuSubmit        func(event SelectMenuSubmitEvent)

	// api.Message Events
	OnGenericMessageEvent func(event GenericMessageEvent)
	OnMessageCreate       func(event MessageCreateEvent)
	OnMessageUpdate       func(event MessageUpdateEvent)
	OnMessageDelete       func(event MessageDeleteEvent)

	// api.Message Reaction Events
	OnGenericReactionEvent       func(event GenericReactionEvents)
	OnMessageReactionAdd         func(event MessageReactionAddEvent)
	OnMessageReactionRemove      func(event MessageReactionRemoveEvent)
	OnMessageReactionRemoveEmoji func(event MessageReactionRemoveEmojiEvent)
	OnMessageReactionRemoveAll   func(event MessageReactionRemoveAllEvent)

	// Self Events
	OnSelfUpdate func(event SelfUpdateEvent)

	// api.User Events
	OnGenericUserEvent func(event GenericUserEvent)
	OnUserUpdate       func(event UserUpdateEvent)
	OnUserTyping       func(event UserTypingEvent)
	OnGuildUserTyping  func(event GuildMemberTypingEvent)
	OnDMUserTyping     func(event DMUserTypingEvent)

	// api.User api.Activity Events
	OnGenericUserActivityEvent func(event GenericUserActivityEvent)
	OnUserActivityStart        func(event UserActivityStartEvent)
	OnUserActivityUpdate       func(event UserActivityUpdateEvent)
	OnUserActivityEnd          func(event UserActivityEndEvent)
}

// OnEvent is getting called everytime we receive an event
func (l ListenerAdapter) OnEvent(event interface{}) {
	switch e := event.(type) {
	case GenericEvent:
		if listener := l.OnGenericEvent; listener != nil {
			listener(e)
		}
	case HeartbeatEvent:
		if listener := l.OnHeartbeat; listener != nil {
			listener(e)
		}
	case HTTPRequestEvent:
		if listener := l.OnHTTPRequest; listener != nil {
			listener(e)
		}
	case RawGatewayEvent:
		if listener := l.OnRawGateway; listener != nil {
			listener(e)
		}
	case ReadyEvent:
		if listener := l.OnReadyEvent; listener != nil {
			listener(e)
		}

	// api.Command Events
	case GenericCommandEvent:
		if listener := l.OnGenericCommandEvent; listener != nil {
			listener(e)
		}
	case CommandCreateEvent:
		if listener := l.OnCommandCreate; listener != nil {
			listener(e)
		}
	case CommandUpdateEvent:
		if listener := l.OnCommandUpdate; listener != nil {
			listener(e)
		}
	case CommandDeleteEvent:
		if listener := l.OnCommandDelete; listener != nil {
			listener(e)
		}

	// api.Channel Events
	case GenericChannelEvent:
		if listener := l.OnGenericChannelEvent; listener != nil {
			listener(e)
		}

	// api.GuildChannel Events
	case GenericGuildChannelEvent:
		if listener := l.OnGenericGuildChannelEvent; listener != nil {
			listener(e)
		}
	case GuildChannelCreateEvent:
		if listener := l.OnGuildChannelCreate; listener != nil {
			listener(e)
		}
	case GuildChannelUpdateEvent:
		if listener := l.OnGuildChannelUpdate; listener != nil {
			listener(e)
		}
	case GuildChannelDeleteEvent:
		if listener := l.OnGuildChannelDelete; listener != nil {
			listener(e)
		}

	// api.Category Events
	case GenericCategoryEvent:
		if listener := l.OnGenericCategoryEvent; listener != nil {
			listener(e)
		}
	case CategoryCreateEvent:
		if listener := l.OnCategoryCreate; listener != nil {
			listener(e)
		}
	case CategoryUpdateEvent:
		if listener := l.OnCategoryUpdate; listener != nil {
			listener(e)
		}
	case CategoryDeleteEvent:
		if listener := l.OnCategoryDelete; listener != nil {
			listener(e)
		}

	// api.DMChannel Events// api.Category Events
	case GenericDMChannelEvent:
		if listener := l.OnGenericDMChannelEvent; listener != nil {
			listener(e)
		}
	case DMChannelCreateEvent:
		if listener := l.OnDMChannelCreate; listener != nil {
			listener(e)
		}
	case DMChannelUpdateEvent:
		if listener := l.OnDMChannelUpdate; listener != nil {
			listener(e)
		}
	case DMChannelDeleteEvent:
		if listener := l.OnDMChannelDelete; listener != nil {
			listener(e)
		}

	// api.DMChannel Events// api.Category Events
	case GenericDMMessageReactionEvent:
		if listener := l.OnGenericDMMessageReactionEventEvent; listener != nil {
			listener(e)
		}
	case DMMessageReactionAddEvent:
		if listener := l.OnDMMessageReactionAdd; listener != nil {
			listener(e)
		}
	case DMMessageReactionRemoveEvent:
		if listener := l.OnDMMessageReactionRemove; listener != nil {
			listener(e)
		}
	case DMMessageReactionRemoveEmojiEvent:
		if listener := l.OnDMMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case DMMessageReactionRemoveAllEvent:
		if listener := l.OnDMMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// api.StoreChannel Events
	case GenericStoreChannelEvent:
		if listener := l.OnGenericStoreChannelEvent; listener != nil {
			listener(e)
		}
	case StoreChannelCreateEvent:
		if listener := l.OnStoreChannelCreate; listener != nil {
			listener(e)
		}
	case StoreChannelUpdateEvent:
		if listener := l.OnStoreChannelUpdate; listener != nil {
			listener(e)
		}
	case StoreChannelDeleteEvent:
		if listener := l.OnStoreChannelDelete; listener != nil {
			listener(e)
		}

	// api.TextChannel Events
	case GenericTextChannelEvent:
		if listener := l.OnGenericTextChannelEvent; listener != nil {
			listener(e)
		}
	case TextChannelCreateEvent:
		if listener := l.OnTextChannelCreate; listener != nil {
			listener(e)
		}
	case TextChannelUpdateEvent:
		if listener := l.OnTextChannelUpdate; listener != nil {
			listener(e)
		}
	case TextChannelDeleteEvent:
		if listener := l.OnTextChannelDelete; listener != nil {
			listener(e)
		}

	// api.VoiceChannel Events
	case GenericVoiceChannelEvent:
		if listener := l.OnGenericVoiceChannelEvent; listener != nil {
			listener(e)
		}
	case VoiceChannelCreateEvent:
		if listener := l.OnVoiceChannelCreate; listener != nil {
			listener(e)
		}
	case VoiceChannelUpdateEvent:
		if listener := l.OnVoiceChannelUpdate; listener != nil {
			listener(e)
		}
	case VoiceChannelDeleteEvent:
		if listener := l.OnVoiceChannelDelete; listener != nil {
			listener(e)
		}

	// api.emote Events
	case GenericEmoteEvent:
		if listener := l.OnGenericEmoteEvent; listener != nil {
			listener(e)
		}
	case EmoteCreateEvent:
		if listener := l.OnEmoteCreate; listener != nil {
			listener(e)
		}
	case EmoteUpdateEvent:
		if listener := l.OnEmoteUpdate; listener != nil {
			listener(e)
		}
	case EmoteDeleteEvent:
		if listener := l.OnEmoteDelete; listener != nil {
			listener(e)
		}

	// api.GatewayStatus Events
	case GenericGatewayStatusEvent:
		if listener := l.OnGenericGatewayStatusEvent; listener != nil {
			listener(e)
		}
	case ConnectedEvent:
		if listener := l.OnConnected; listener != nil {
			listener(e)
		}
	case ReconnectedEvent:
		if listener := l.OnReconnected; listener != nil {
			listener(e)
		}
	case ResumedEvent:
		if listener := l.OnResumed; listener != nil {
			listener(e)
		}
	case DisconnectedEvent:
		if listener := l.OnDisconnected; listener != nil {
			listener(e)
		}

	// api.Guild Events
	case GenericGuildEvent:
		if listener := l.OnGenericGuildEvent; listener != nil {
			listener(e)
		}
	case GuildJoinEvent:
		if listener := l.OnGuildJoin; listener != nil {
			listener(e)
		}
	case GuildUpdateEvent:
		if listener := l.OnGuildUpdate; listener != nil {
			listener(e)
		}
	case GuildLeaveEvent:
		if listener := l.OnGuildLeave; listener != nil {
			listener(e)
		}
	case GuildAvailableEvent:
		if listener := l.OnGuildAvailable; listener != nil {
			listener(e)
		}
	case GuildUnavailableEvent:
		if listener := l.OnGuildUnavailable; listener != nil {
			listener(e)
		}
	case GuildReadyEvent:
		if listener := l.OnGuildReady; listener != nil {
			listener(e)
		}
	case GuildBanEvent:
		if listener := l.OnGuildBan; listener != nil {
			listener(e)
		}
	case GuildUnbanEvent:
		if listener := l.OnGuildUnban; listener != nil {
			listener(e)
		}

	// api.Guild api.Invite Events
	case GenericGuildInviteEvent:
		if listener := l.OnGenericGuildInviteEvent; listener != nil {
			listener(e)
		}
	case GuildInviteCreateEvent:
		if listener := l.OnGuildInviteCreate; listener != nil {
			listener(e)
		}
	case GuildInviteDeleteEvent:
		if listener := l.OnGuildInviteDelete; listener != nil {
			listener(e)
		}

	// api.Member Events
	case GenericGuildMemberEvent:
		if listener := l.OnGenericGuildMemberEvent; listener != nil {
			listener(e)
		}
	case GuildMemberJoinEvent:
		if listener := l.OnGuildMemberJoin; listener != nil {
			listener(e)
		}
	case GuildMemberUpdateEvent:
		if listener := l.OnGuildMemberUpdate; listener != nil {
			listener(e)
		}
	case GuildMemberLeaveEvent:
		if listener := l.OnGuildMemberLeave; listener != nil {
			listener(e)
		}

	// api.Guild api.Message Events
	case GenericGuildMessageEvent:
		if listener := l.OnGenericGuildMessageEvent; listener != nil {
			listener(e)
		}
	case GuildMessageCreateEvent:
		if listener := l.OnGuildMessageCreate; listener != nil {
			listener(e)
		}
	case GuildMessageUpdateEvent:
		if listener := l.OnGuildMessageUpdate; listener != nil {
			listener(e)
		}
	case GuildMessageDeleteEvent:
		if listener := l.OnGuildMessageDelete; listener != nil {
			listener(e)
		}

	// api.Guild api.Message Reaction Events
	case GenericGuildMessageReactionEvent:
		if listener := l.OnGenericGuildMessageReactionEvent; listener != nil {
			listener(e)
		}
	case GuildMessageReactionAddEvent:
		if listener := l.OnGuildMessageReactionAdd; listener != nil {
			listener(e)
		}
	case GuildMessageReactionRemoveEvent:
		if listener := l.OnGuildMessageReactionRemove; listener != nil {
			listener(e)
		}
	case GuildMessageReactionRemoveEmojiEvent:
		if listener := l.OnGuildMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case GuildMessageReactionRemoveAllEvent:
		if listener := l.OnGuildMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// api.Guild Voice Events
	case GenericGuildVoiceEvent:
		if listener := l.OnGenericGuildVoiceEvent; listener != nil {
			listener(e)
		}
	case GuildVoiceUpdateEvent:
		if listener := l.OnGuildVoiceUpdate; listener != nil {
			listener(e)
		}
	case GuildVoiceJoinEvent:
		if listener := l.OnGuildVoiceJoin; listener != nil {
			listener(e)
		}
	case GuildVoiceLeaveEvent:
		if listener := l.OnGuildVoiceLeave; listener != nil {
			listener(e)
		}

	// api.Guild api.Role Events
	case GenericRoleEvent:
		if listener := l.OnGenericRoleEvent; listener != nil {
			listener(e)
		}
	case RoleCreateEvent:
		if listener := l.OnRoleCreate; listener != nil {
			listener(e)
		}
	case RoleUpdateEvent:
		if listener := l.OnRoleUpdate; listener != nil {
			listener(e)
		}
	case RoleDeleteEvent:
		if listener := l.OnRoleDelete; listener != nil {
			listener(e)
		}

	// Interaction Events
	case GenericInteractionEvent:
		if listener := l.OnGenericInteractionEvent; listener != nil {
			listener(e)
		}
	case CommandEvent:
		if listener := l.OnCommand; listener != nil {
			listener(e)
		}
	case GenericComponentEvent:
		if listener := l.OnGenericComponentEvent; listener != nil {
			listener(e)
		}
	case ButtonClickEvent:
		if listener := l.OnButtonClick; listener != nil {
			listener(e)
		}
	case SelectMenuSubmitEvent:
		if listener := l.OnSelectMenuSubmit; listener != nil {
			listener(e)
		}

	// api.Message Events
	case GenericMessageEvent:
		if listener := l.OnGenericMessageEvent; listener != nil {
			listener(e)
		}
	case MessageCreateEvent:
		if listener := l.OnMessageCreate; listener != nil {
			listener(e)
		}
	case MessageUpdateEvent:
		if listener := l.OnMessageUpdate; listener != nil {
			listener(e)
		}
	case MessageDeleteEvent:
		if listener := l.OnMessageDelete; listener != nil {
			listener(e)
		}

	// api.Message Reaction Events
	case GenericReactionEvents:
		if listener := l.OnGenericReactionEvent; listener != nil {
			listener(e)
		}
	case MessageReactionAddEvent:
		if listener := l.OnMessageReactionAdd; listener != nil {
			listener(e)
		}
	case MessageReactionRemoveEvent:
		if listener := l.OnMessageReactionRemove; listener != nil {
			listener(e)
		}
	case MessageReactionRemoveEmojiEvent:
		if listener := l.OnMessageReactionRemoveEmoji; listener != nil {
			listener(e)
		}
	case MessageReactionRemoveAllEvent:
		if listener := l.OnMessageReactionRemoveAll; listener != nil {
			listener(e)
		}

	// Self Events
	case SelfUpdateEvent:
		if listener := l.OnSelfUpdate; listener != nil {
			listener(e)
		}

	// api.User Events
	case GenericUserEvent:
		if listener := l.OnGenericUserEvent; listener != nil {
			listener(e)
		}
	case UserUpdateEvent:
		if listener := l.OnUserUpdate; listener != nil {
			listener(e)
		}
	case UserTypingEvent:
		if listener := l.OnUserTyping; listener != nil {
			listener(e)
		}
	case GuildMemberTypingEvent:
		if listener := l.OnGuildUserTyping; listener != nil {
			listener(e)
		}
	case DMUserTypingEvent:
		if listener := l.OnDMUserTyping; listener != nil {
			listener(e)
		}

	// api.User api.Activity Events
	case GenericUserActivityEvent:
		if listener := l.OnGenericUserActivityEvent; listener != nil {
			listener(e)
		}
	case UserActivityStartEvent:
		if listener := l.OnUserActivityStart; listener != nil {
			listener(e)
		}
	case UserActivityUpdateEvent:
		if listener := l.OnUserActivityUpdate; listener != nil {
			listener(e)
		}
	case UserActivityEndEvent:
		if listener := l.OnUserActivityEnd; listener != nil {
			listener(e)
		}

	default:
		if e, ok := e.(api.Event); ok {
			e.Disgo().Logger().Errorf("unexpected event received: \"%s\", event: \"%#e\"", reflect.TypeOf(event).Name(), event)
		}
	}
}
