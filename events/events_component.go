package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// ButtonClickEvent indicates that an core.Button was clicked
type ButtonClickEvent struct {
	*GenericEvent
	*core.ButtonInteraction
}

func (e *ButtonClickEvent) Bot() *core.Bot {
	return e.bot
}

// SelectMenuSubmitEvent indicates that an core.SelectMenu was submitted
type SelectMenuSubmitEvent struct {
	*GenericEvent
	*core.SelectMenuInteraction
}

func (e *SelectMenuSubmitEvent) Bot() *core.Bot {
	return e.bot
}
