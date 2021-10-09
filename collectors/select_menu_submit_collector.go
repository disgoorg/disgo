package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewSelectMenuSubmitCollector gives you a channel to receive on and a function to close the collector
func NewSelectMenuSubmitCollector(disgo *core.Bot, filter core.SelectMenuInteractionFilter) (<-chan *core.SelectMenuInteraction, func()) {
	ch := make(chan *core.SelectMenuInteraction)

	collector := &SelectMenuSubmitCollector{
		Filter:    filter,
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

// SelectMenuSubmitCollector used to collect core.SelectMenuInteraction(s) from a core.Message using a core.SelectMenuInteractionFilter function
type SelectMenuSubmitCollector struct {
	Filter    core.SelectMenuInteractionFilter
	Chan   chan<- *core.SelectMenuInteraction
	Close     func()
}

// OnEvent used to get events for the SelectMenuSubmitCollector
func (c *SelectMenuSubmitCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.SelectMenuSubmitEvent); ok {
		if !c.Filter(event.SelectMenuInteraction) {
			return
		}
		c.Chan <- event.SelectMenuInteraction
	}
}
