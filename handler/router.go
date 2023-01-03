package handler

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type (
	CommandHandler      func(e *CommandEvent) error
	AutocompleteHandler func(e *AutocompleteEvent) error
	ComponentHandler    func(e *ComponentEvent) error
	ModalHandler        func(e *ModalEvent) error
)

var (
	_ Route = (*mux)(nil)
	_ Route = (*handlerHolder[CommandHandler])(nil)
)

type Route interface {
	// Match returns true if the given path matches the Route.
	Match(path string, t discord.InteractionType) bool

	// Handle handles the given interaction event.
	Handle(path string, variables map[string]string, e *events.InteractionCreate) error
}

type Router interface {
	bot.EventListener
	Route

	// Use adds the given middlewares to the current Router
	Use(middlewares ...Middleware)

	// With returns a new Router with the given middlewares
	With(middlewares ...Middleware) Router

	// Group creates a new Router and adds it to the current Router.
	Group(fn func(r Router))

	// Route creates a new sub-router with the given pattern and adds it to the current Router.
	Route(pattern string, fn func(r Router)) Router

	// Mount mounts the given router with the given pattern to the current Router.
	Mount(pattern string, r Router)

	// HandleCommand registers the given CommandHandler to the current Router.
	HandleCommand(pattern string, h CommandHandler)

	// HandleAutocomplete registers the given AutocompleteHandler to the current Router.
	HandleAutocomplete(pattern string, h AutocompleteHandler)

	// HandleComponent registers the given ComponentHandler to the current Router.
	HandleComponent(pattern string, h ComponentHandler)

	// HandleModal registers the given ModalHandler to the current Router.
	HandleModal(pattern string, h ModalHandler)
}
