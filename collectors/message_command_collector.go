package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewMessageCommandCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewMessageCommandCollector(disgo *core.Bot, filter core.MessageCommandInteractionFilter) (<-chan *core.MessageCommandInteraction, func()) {
	ch := make(chan *core.MessageCommandInteraction)

	col := &MessageCommandCollector{
		Filter: filter,
		Chan:   ch,
	}
	cls := func() {
		close(ch)
		disgo.EventManager.RemoveEventListeners(col)
	}
	col.Close = cls
	disgo.EventManager.AddEventListeners(col)

	return ch, cls
}

// MessageCommandCollector used to collect core.MessageCommandInteraction(s) from a core.Message using a ButtonFilter function
type MessageCommandCollector struct {
	Filter core.MessageCommandInteractionFilter
	Chan   chan<- *core.MessageCommandInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *MessageCommandCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.MessageCommandEvent); ok {
		if !c.Filter(event.MessageCommandInteraction) {
			return
		}
		c.Chan <- event.MessageCommandInteraction
	}
}
