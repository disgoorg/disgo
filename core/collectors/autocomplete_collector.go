package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewAutocompleteCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewAutocompleteCollector(disgo *core.Bot, filter core.AutocompleteInteractionFilter) (<-chan *core.AutocompleteInteraction, func()) {
	ch := make(chan *core.AutocompleteInteraction)

	col := &AutocompleteCollector{
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

// AutocompleteCollector used to collect core.AutocompleteInteraction(s) from a core.Message using a ButtonFilter function
type AutocompleteCollector struct {
	Filter core.AutocompleteInteractionFilter
	Chan   chan<- *core.AutocompleteInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *AutocompleteCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.AutocompleteEvent); ok {
		if !c.Filter(event.AutocompleteInteraction) {
			return
		}
		c.Chan <- event.AutocompleteInteraction
	}
}
