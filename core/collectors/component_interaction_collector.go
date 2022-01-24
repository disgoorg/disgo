package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewComponentInteractionCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewComponentInteractionCollector(disgo *core.Bot, filter core.ComponentInteractionFilter) (<-chan *core.ComponentInteraction, func()) {
	ch := make(chan *core.ComponentInteraction)

	col := &ComponentInteractionCollector{
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

// ComponentInteractionCollector used to collect core.ComponentInteractionComponentInteraction(s) from a core.Message using a ButtonFilter function
type ComponentInteractionCollector struct {
	Filter core.ComponentInteractionFilter
	Chan   chan<- *core.ComponentInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *ComponentInteractionCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.ComponentInteractionEvent); ok {
		if !c.Filter(event.ComponentInteraction) {
			return
		}
		c.Chan <- event.ComponentInteraction
	}
}
