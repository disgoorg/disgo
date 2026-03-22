// Package handler is a provides a command handler in the style of an HTTP router.
// It is inspired by the [github.com/go-chi/chi] package.
//
// The handler allows mapping command names (e.g. "/ping" or "/info/commands")
// and component custom IDs (e.g. "/button1" or "/menu/close/{id}") to handlers.
// Paths can contain variables (e.g. "/user/{id}") but need to be a full path segment.
// The variables can be accessed via the Vars field in the event.
//
// Middlewares can be used to intercept and short-circuit the interactions.
// One thing to be aware of is that middlewares are executed while the tree is being traversed.
// This means you can't access all path variables in a middleware, only the ones that have been matched so far in the path.
//
// The handler also supports sub-routers, which can be used to group handlers and middlewares.
// e.g.:
//
// r := handler.New()
//
//	r.Route("/info", func(r handler.Router) {
//		r.Use(someMiddleware)
//		r.SlashCommand("/commands", someCommandHandler)
//	})
//
//	r.Group(func(r handler.Router) {
//		r.Use(someMiddleware2)
//		r.SlashCommand("/options", someCommandHandler2)
//	})
//
// To register a handler, you can use the following methods:
// - [Router.Interaction]: for any interaction type handlers
//
// - [Router.Command]: for application command handlers
// - [Router.SlashCommand]: for slash command handlers
// - [Router.UserCommand]: for user command handlers
// - [Router.MessageCommand]: for message command handlers
// - [Router.EntryPointCommand]: for entry point command handlers
// - [Router.Autocomplete]: for autocomplete command handlers
//
// - [Router.Component]: for component handlers
// - [Router.ButtonComponent]: for button component handlers
// - [Router.SelectMenuComponent]: for select menu component handlers
//
// - [Router.Modal]: for modal handlers
//
// To register a middleware, you can use the following methods:
// - [Router.Use]: to add a middleware to the current router
// - [Router.With]: to create a new router with the given middlewares
//
// To create a sub-router, you can use the following methods:
// - [Router.Group]: to create a new router and add it to the current router
// - [Router.Route]: to create a new sub-router with the given pattern and add it to the current router
// - [Router.Mount]: to mount the given router with the given pattern to the current router
package handler
