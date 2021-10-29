package collectors

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/events"
)

// NewSlashCommandCollector gives you a channel to receive on and a function to close the collector
//goland:noinspection GoUnusedExportedFunction
func NewSlashCommandCollector(disgo *core.Bot, filter core.SlashCommandInteractionFilter) (<-chan *core.SlashCommandInteraction, func()) {
	ch := make(chan *core.SlashCommandInteraction)

	col := &SlashCommandCollector{
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

// SlashCommandCollector used to collect core.SlashCommandInteraction(s) from a core.Message using a ButtonFilter function
type SlashCommandCollector struct {
	Filter core.SlashCommandInteractionFilter
	Chan   chan<- *core.SlashCommandInteraction
	Close  func()
}

// OnEvent used to get events for the ButtonCollector
func (c *SlashCommandCollector) OnEvent(e core.Event) {
	if event, ok := e.(*events.SlashCommandEvent); ok {
		if !c.Filter(event.SlashCommandInteraction) {
			return
		}
		c.Chan <- event.SlashCommandInteraction
	}
}
