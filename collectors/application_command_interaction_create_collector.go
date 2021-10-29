package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewApplicationCommandInteractionCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewApplicationCommandInteractionCollector(disgo *core.Bot, filter core.ApplicationCommandInteractionFilter) (<-chan core.ApplicationCommandInteraction, func()) {
	ch := make(chan core.ApplicationCommandInteraction)

	col := &ApplicationCommandInteractionCollector{
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

// ApplicationCommandInteractionCollector used to collect core.ApplicationCommandInteractionApplicationCommandInteraction(s) from a core.Message using a ButtonFilter function
type ApplicationCommandInteractionCollector struct {
	Filter core.ApplicationCommandInteractionFilter
	Chan   chan<- core.ApplicationCommandInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *ApplicationCommandInteractionCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.ApplicationCommandInteractionCreateEvent); ok {
		if !c.Filter(event.ApplicationCommandInteraction) {
			return
		}
		c.Chan <- event.ApplicationCommandInteraction
	}
}
