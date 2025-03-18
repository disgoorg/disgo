package handler

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

type (
	// InteractionHandler is a function that handles all types of interactions.
	InteractionHandler func(e *InteractionEvent) error

	// CommandHandler is a function that handles all types of application command interactions.
	CommandHandler func(e *CommandEvent) error
	// SlashCommandHandler is a function that handles slash command interactions.
	SlashCommandHandler func(data discord.SlashCommandInteractionData, e *CommandEvent) error
	// UserCommandHandler is a function that handles user command interactions.
	UserCommandHandler func(data discord.UserCommandInteractionData, e *CommandEvent) error
	// MessageCommandHandler is a function that handles message command interactions.
	MessageCommandHandler func(data discord.MessageCommandInteractionData, e *CommandEvent) error
	// EntryPointCommandHandler is a function that handles entry point command interactions.
	EntryPointCommandHandler func(data discord.EntryPointCommandInteractionData, e *CommandEvent) error

	// AutocompleteHandler is a function that handles autocomplete interactions.
	AutocompleteHandler func(e *AutocompleteEvent) error

	// ComponentHandler is a function that handles all types of component interactions.
	ComponentHandler func(e *ComponentEvent) error
	// ButtonComponentHandler is a function that handles button component interactions.
	ButtonComponentHandler func(data discord.ButtonInteractionData, e *ComponentEvent) error
	// SelectMenuComponentHandler is a function that handles select menu component interactions.
	SelectMenuComponentHandler func(data discord.SelectMenuInteractionData, e *ComponentEvent) error

	// ModalHandler is a function that handles modals.
	ModalHandler func(e *ModalEvent) error

	// NotFoundHandler is a function that is called when no route was found.
	NotFoundHandler func(e *InteractionEvent) error

	// ErrorHandler is a function that is called when an error occurs during handling an interaction.
	ErrorHandler func(e *InteractionEvent, err error)
)

var (
	_ Route = (*Mux)(nil)
	_ Route = (*handlerHolder[CommandHandler])(nil)
	_ Route = (*handlerHolder[AutocompleteHandler])(nil)
	_ Route = (*handlerHolder[ComponentHandler])(nil)
	_ Route = (*handlerHolder[ModalHandler])(nil)
)

// Route is a basic interface for a route in a Router.
type Route interface {
	// Match returns true if the given path matches the Route.
	Match(path string, t discord.InteractionType, t2 int) bool

	// Handle handles the given interaction event.
	Handle(path string, e *InteractionEvent) error
}

// Router provides with the core routing functionality.
// It is used to register handlers and middlewares and sub-routers.
type Router interface {
	bot.EventListener
	Route

	// Use adds the given middlewares to the current Router.
	Use(middlewares ...Middleware)

	// With returns a new Router with the given middlewares.
	With(middlewares ...Middleware) Router

	// Group creates a new Router and adds it to the current Router.
	Group(fn func(r Router))

	// Route creates a new sub-router with the given pattern and adds it to the current Router.
	Route(pattern string, fn func(r Router)) Router

	// Mount mounts the given router with the given pattern to the current Router.
	Mount(pattern string, r Router)

	// Interaction registers the given InteractionHandler to the current Router.
	Interaction(pattern string, h InteractionHandler)

	// Command registers the given CommandHandler to the current Router.
	Command(pattern string, h CommandHandler)

	// SlashCommand registers the given SlashCommandHandler to the current Router.
	SlashCommand(pattern string, h SlashCommandHandler)

	// UserCommand registers the given UserCommandHandler to the current Router.
	UserCommand(pattern string, h UserCommandHandler)

	// MessageCommand registers the given MessageCommandHandler to the current Router.
	MessageCommand(pattern string, h MessageCommandHandler)

	// EntryPointCommand registers the given EntryPointCommandHandler to the current Router.
	EntryPointCommand(pattern string, h EntryPointCommandHandler)

	// Autocomplete registers the given AutocompleteHandler to the current Router.
	Autocomplete(pattern string, h AutocompleteHandler)

	// Component registers the given ComponentHandler to the current Router.
	Component(pattern string, h ComponentHandler)

	// ButtonComponent registers the given ButtonComponentHandler to the current Router.
	ButtonComponent(pattern string, h ButtonComponentHandler)

	// SelectMenuComponent registers the given SelectMenuComponentHandler to the current Router.
	SelectMenuComponent(pattern string, h SelectMenuComponentHandler)

	// Modal registers the given ModalHandler to the current Router.
	Modal(pattern string, h ModalHandler)
}
