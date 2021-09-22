package core

// SlashCommandEvent indicates that a slash core.ApplicationCommand was run
type SlashCommandEvent struct {
	*GenericEvent
	*SlashCommandInteraction
}

func (e *SlashCommandEvent) Bot() *Bot {
	return e.bot
}

type ApplicationCommandAutocompleteEvent struct {
	*GenericEvent
	*ApplicationCommandAutocompleteInteraction
}

func (e *ApplicationCommandAutocompleteEvent) Bot() *Bot {
	return e.bot
}

type UserCommandEvent struct {
	*GenericEvent
	*UserCommandInteraction
}

func (e *UserCommandEvent) Bot() *Bot {
	return e.bot
}

type MessageCommandEvent struct {
	*GenericEvent
	*MessageCommandInteraction
}

func (e *MessageCommandEvent) Bot() *Bot {
	return e.bot
}
