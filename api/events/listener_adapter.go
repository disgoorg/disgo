package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// ListenerAdapter lets you override the handles for receiving events
type ListenerAdapter struct {
	OnGenericEvent func(*api.Event)

	// Guild Events
	OnGenericGuildEvent func(event *GenericGuildEvent)
	OnGuildJoin         func(event *GuildJoinEvent)
	OnGuildUpdate       func(event *GuildUpdateEvent)
	OnGuildLeave        func(event *GuildLeaveEvent)
	OnGuildAvailable    func(event *GuildAvailableEvent)
	OnGuildUnavailable  func(event *GuildUnavailableEvent)

	// Guild Role Events
	OnGenericRole func(event *GenericRoleEvent)
	OnRoleCreate  func(event *RoleCreateEvent)
	OnRoleUpdate  func(event *RoleUpdateEvent)
	OnRoleDelete  func(event *RoleDeleteEvent)

	// Message Events
	OnMessageReceived      func(event *MessageReceivedEvent)
	OnGuildMessageReceived func(event *GuildMessageReceivedEvent)

	// Interaction Events
	OnGenericInteraction func(event *GenericInteractionEvent)
	OnSlashCommand       func(event *SlashCommandEvent)
}

// OnEvent is getting called everytime we receive an event
func (l ListenerAdapter) OnEvent(event interface{}) {
	if e, ok := event.(api.Event); ok {
		if l.OnGenericEvent != nil {
			l.OnGenericEvent(&e)
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
		if l.OnGenericRole != nil {
			l.OnGenericRole(&e)
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
		if l.OnGenericInteraction != nil {
			l.OnGenericInteraction(&e)
		}
	case SlashCommandEvent:
		if l.OnSlashCommand != nil {
			l.OnSlashCommand(&e)
		}
	default:
		//log.Errorf("unexpected event received: %#e", event)
	}
}
