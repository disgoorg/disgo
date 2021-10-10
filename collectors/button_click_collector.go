package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewButtonClickCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewButtonClickCollector(disgo *core.Bot, filter core.ButtonInteractionFilter) (<-chan *core.ButtonInteraction, func()) {
	ch := make(chan *core.ButtonInteraction)

	col := &ButtonClickCollector{
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

// ButtonClickCollector used to collect core.ButtonInteraction(s) from a core.Message using a ButtonFilter function
type ButtonClickCollector struct {
	Filter core.ButtonInteractionFilter
	Chan   chan<- *core.ButtonInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *ButtonClickCollector) OnEvent(e interface{}) {
	if event, ok := e.(*events.ButtonClickEvent); ok {
		if !c.Filter(event.ButtonInteraction) {
			return
		}
		c.Chan <- event.ButtonInteraction
	}
}
