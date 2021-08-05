package util

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo api.Disgo, filter MessageFilter) (chan *api.Message, func()) {
	ch := make(chan *api.Message)
	col := &MessageCollector{
		Filter:  filter,
		Channel: ch,
	}
	disgo.EventManager().AddEventListeners(col)

	return ch, func() {
		close(ch)
		disgo.EventManager().RemoveEventListener(col)
	}
}

// MessageFilter used to filter api.Message(s) in a MessageCollector
type MessageFilter func(msg *api.Message) bool

// MessageCollector collects api.Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter  MessageFilter
	Channel chan *api.Message
}

// OnEvent used to get events for the MessageCollector
func (c *MessageCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.MessageCreateEvent); ok {
		if !c.Filter(event.Message) {
			return
		}

		c.Channel <- event.Message
	}
}
