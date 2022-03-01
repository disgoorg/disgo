package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
)

// NewAutocompleteInteractionCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewAutocompleteInteractionCollector(disgo *core.Bot, filter core.AutocompleteInteractionFilter) (<-chan *core.AutocompleteInteraction, func()) {
	ch := make(chan *core.AutocompleteInteraction)

	col := &AutocompleteInteractionCollector{
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

// AutocompleteInteractionCollector used to collect core.AutocompleteInteraction(s) from a core.Message using a ButtonFilter function
type AutocompleteInteractionCollector struct {
	Filter core.AutocompleteInteractionFilter
	Chan   chan<- *core.AutocompleteInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *AutocompleteInteractionCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.AutocompleteInteractionEvent); ok {
		if !c.Filter(event.AutocompleteInteraction) {
			return
		}
		c.Chan <- event.AutocompleteInteraction
	}
}
