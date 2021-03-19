package events

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

type ListenerAdapter struct {
	OnGenericEvent         func(*api.GenericEvent)
	OnGuildAvailable       func(*GuildAvailableEvent)
	OnGuildUnavailable     func(*GuildUnavailableEvent)
	OnGuildJoin            func(*GuildJoinEvent)
	OnGuildLeave           func(*GuildLeaveEvent)
	OnMessageReceived      func(*MessageReceivedEvent)
	OnGuildMessageReceived func(*GuildMessageReceivedEvent)
	OnSlashCommand         func(*SlashCommandEvent)
}

func (l ListenerAdapter) OnEvent(event interface{}) {
	switch v := event.(type) {
	case GuildAvailableEvent:
		if l.OnGuildAvailable != nil {
			l.OnGuildAvailable(&v)
		}
	case GuildUnavailableEvent:
		if l.OnGuildUnavailable != nil {
			l.OnGuildUnavailable(&v)
		}
	case GuildJoinEvent:
		if l.OnGuildJoin != nil {
			l.OnGuildJoin(&v)
		}
	case GuildLeaveEvent:
		if l.OnGuildLeave != nil {
			l.OnGuildLeave(&v)
		}
	case MessageReceivedEvent:
		if l.OnMessageReceived != nil {
			l.OnMessageReceived(&v)
		}
	case GuildMessageReceivedEvent:
		if l.OnGuildMessageReceived != nil {
			l.OnGuildMessageReceived(&v)
		}
	case SlashCommandEvent:
		if l.OnSlashCommand != nil {
			l.OnSlashCommand(&v)
		}
	default:
		log.Errorf("unexpected event received: %#v", event)
	}
	if event, ok := event.(api.GenericEvent); ok {
		if l.OnGenericEvent != nil {
			l.OnGenericEvent(&event)
		}
	}
}
