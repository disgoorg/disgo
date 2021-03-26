package events

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

type ListenerAdapter struct {
	OnGenericEvent func(*api.GenericEvent)

	// Guild Events
	OnGenericGuildEvent func(*GenericGuildEvent)
	OnGuildJoin         func(*GuildJoinEvent)
	OnGuildUpdate       func(*GuildUpdateEvent)
	OnGuildLeave        func(*GuildLeaveEvent)
	OnGuildAvailable    func(*GuildAvailableEvent)
	OnGuildUnavailable  func(*GuildUnavailableEvent)

	// Guild Role Events
	OnGenericGuildRole func(*GenericGuildRoleEvent)
	OnGuildRoleCreate  func(*GuildRoleCreateEvent)
	OnGuildRoleUpdate  func(*GuildRoleUpdateEvent)
	OnGuildRoleDelete  func(*GuildRoleDeleteEvent)

	// Message Events
	OnMessageReceived      func(*MessageReceivedEvent)
	OnGuildMessageReceived func(*GuildMessageReceivedEvent)

	// Interaction Events
	OnGenericInteraction   func(*GenericInteractionEvent)
	OnSlashCommand         func(*SlashCommandEvent)
}

func (l ListenerAdapter) OnEvent(event interface{}) {
	if event, ok := event.(api.GenericEvent); ok {
		if l.OnGenericEvent != nil {
			l.OnGenericEvent(&event)
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
	case GenericGuildRoleEvent:
		if l.OnGenericGuildRole != nil {
			l.OnGenericGuildRole(&e)
		}
	case GuildRoleCreateEvent:
		if l.OnGuildRoleCreate != nil {
			l.OnGuildRoleCreate(&e)
		}
	case GuildRoleUpdateEvent:
		if l.OnGuildRoleUpdate != nil {
			l.OnGuildRoleUpdate(&e)
		}
	case GuildRoleDeleteEvent:
		if l.OnGuildRoleDelete != nil {
			l.OnGuildRoleDelete(&e)
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
		log.Errorf("unexpected event received: %#e", event)
	}
}
