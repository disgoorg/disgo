package handlers

import "github.com/DisgoOrg/disgo/core"

func init() {
	for _, handler := range EventHandlers {
		core.HTTPEventHandlers[handler.EventType()] = handler
	}
}

var EventHandlers = []core.HTTPEventHandler{
	&InteractionCreateWebhookHandler{},
}
