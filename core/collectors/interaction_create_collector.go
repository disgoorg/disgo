package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewInteractionCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewInteractionCollector(disgo *core.Bot, filter core.InteractionFilter) (<-chan core.Interaction, func()) {
	ch := make(chan core.Interaction)

	col := &InteractionCollector{
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

// InteractionCollector used to collect core.InteractionInteraction(s) from a core.Message using a ButtonFilter function
type InteractionCollector struct {
	Filter core.InteractionFilter
	Chan   chan<- core.Interaction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *InteractionCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.InteractionEvent); ok {
		if !c.Filter(event.Interaction) {
			return
		}
		c.Chan <- event.Interaction
	}
}
