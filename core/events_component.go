package core

// ButtonClickEvent indicates that an core.Button was clicked
type ButtonClickEvent struct {
	*GenericEvent
	*ButtonInteraction
}

func (e *ButtonClickEvent) Bot() *Bot {
	return e.bot
}

// SelectMenuSubmitEvent indicates that an core.SelectMenu was submitted
type SelectMenuSubmitEvent struct {
	*GenericEvent
	*SelectMenuInteraction
}

func (e *SelectMenuSubmitEvent) Bot() *Bot {
	return e.bot
}
