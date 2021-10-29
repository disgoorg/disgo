package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewUserCommandCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewUserCommandCollector(disgo *core.Bot, filter core.UserCommandInteractionFilter) (<-chan *core.UserCommandInteraction, func()) {
	ch := make(chan *core.UserCommandInteraction)

	col := &UserCommandCollector{
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

// UserCommandCollector used to collect core.UserCommandInteraction(s) from a core.Message using a ButtonFilter function
type UserCommandCollector struct {
	Filter core.UserCommandInteractionFilter
	Chan   chan<- *core.UserCommandInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *UserCommandCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.UserCommandEvent); ok {
		if !c.Filter(event.UserCommandInteraction) {
			return
		}
		c.Chan <- event.UserCommandInteraction
	}
}
