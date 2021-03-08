package disgo

import (
	log "github.com/sirupsen/logrus"
)

type ListenerAdapter struct {}

func (l ListenerAdapter) OnGenericEvent(event GenericEvent) {
	println("OnGenericEvent")
}

func (l ListenerAdapter) OnGuildAvailable(GuildAvailableEvent) {}

func (l ListenerAdapter) OnGuildUnavailable(GuildUnavailableEvent) {}

func (l ListenerAdapter) OnGuildJoin(GuildJoinEvent) {}

func (l ListenerAdapter) OnGuildLeave(GuildLeaveEvent) {}

func (l ListenerAdapter) OnGuildMessageReceived(GuildMessageReceivedEvent) {
	println("fuck")
}

func (l ListenerAdapter) OnEvent(event interface{}) {
	println("OnEvent: %v", event)
	switch v := event.(type) {
	case GenericEvent:
		l.OnGenericEvent(v)
	case GuildAvailableEvent:
		l.OnGuildAvailable(v)
	case GuildUnavailableEvent:
		l.OnGuildUnavailable(v)
	case GuildJoinEvent:
		l.OnGuildJoin(v)
	case GuildLeaveEvent:
		l.OnGuildLeave(v)
	case GuildMessageReceivedEvent:
		l.OnGuildMessageReceived(v)
	default:
		log.Errorf("unexpected event received: %#v", event)
	}
}