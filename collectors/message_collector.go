package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewMessageCollector gives you a channel to receive on and a function to close the collector
func NewMessageCollector(disgo *core.Bot, filter core.MessageFilter) (<-chan *core.Message, func()) {
	ch := make(chan *core.Message)

	collector := &MessageCollector{
		Filter: filter,
		Chan:   ch,
	}
	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(collector)
	}
	collector.Close = cls
	disgo.EventManager.AddEventListeners(collector)

	return ch, cls
}

// MessageCollector collects Message(s) using a MessageFilter function
type MessageCollector struct {
	Filter core.MessageFilter
	Chan   chan<- *core.Message
	Close  func()
}

// OnEvent used to get events for the MessageCollector
func (c *MessageCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.MessageCreateEvent); ok {
		if !c.Filter(event.Message) {
			return
		}
		c.Chan <- event.Message
	}
}
