package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var defaultErrorHandler ErrorHandler = func(event *InteractionEvent, err error) {
	event.Client().Logger().Error("error handling interaction", slog.Any("err", err))
}

// New returns a new Router.
func New() *Mux {
	return &Mux{}
}

func newRouter(pattern string, middlewares []Middleware, routes []Route) *Mux {
	return &Mux{
		pattern:     pattern,
		middlewares: middlewares,
		routes:      routes,
	}
}

// Mux is a basic Router implementation.
type Mux struct {
	pattern         string
	middlewares     []Middleware
	routes          []Route
	notFoundHandler NotFoundHandler
	errorHandler    ErrorHandler
	defaultContext  func() context.Context
}

// OnEvent is called when a new event is received.
func (r *Mux) OnEvent(event bot.Event) {
	e, ok := event.(*events.InteractionCreate)
	if !ok {
		return
	}

	var path string
	switch i := e.Interaction.(type) {
	case discord.ApplicationCommandInteraction:
		if sci, ok := i.Data.(discord.SlashCommandInteractionData); ok {
			path = sci.CommandPath()
		} else {
			path = "/" + i.Data.CommandName()
		}
	case discord.AutocompleteInteraction:
		path = i.Data.CommandPath()
	case discord.ComponentInteraction:
		path = i.Data.CustomID()
	case discord.ModalSubmitInteraction:
		path = i.Data.CustomID
	}

	var ctx context.Context
	if r.defaultContext != nil {
		ctx = r.defaultContext()
	} else {
		ctx = context.Background()
	}

	ie := &InteractionEvent{
		InteractionCreate: e,
		Ctx:               ctx,
		Vars:              make(map[string]string),
	}
	if err := r.Handle(path, ie); err != nil {
		if r.errorHandler != nil {
			r.errorHandler(ie, err)
			return
		}
		defaultErrorHandler(ie, err)
	}
}

// Match returns true if the given path matches the Route.
func (r *Mux) Match(path string, t discord.InteractionType, t2 int) bool {
	if r.pattern != "" {
		parts := splitPath(path)
		patternParts := splitPath(r.pattern)

		for i, part := range patternParts {
			path = strings.TrimPrefix(path, "/"+parts[i])
			if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
				continue
			}
			if len(parts) <= i || part != parts[i] {
				return false
			}
		}
	}

	for _, matcher := range r.routes {
		if matcher.Match(path, t, t2) {
			return true
		}
	}
	return false
}

// Handle handles the given interaction event.
func (r *Mux) Handle(path string, event *InteractionEvent) error {
	handlerChain := Handler(func(event *InteractionEvent) error {
		path = parseVariables(path, r.pattern, event.Vars)

		t := event.Type()
		var t2 int
		switch i := event.Interaction.(type) {
		case discord.ApplicationCommandInteraction:
			t2 = int(i.Data.Type())
		case discord.ComponentInteraction:
			t2 = int(i.Data.Type())
		}

		for _, route := range r.routes {
			if route.Match(path, t, t2) {
				return route.Handle(path, event)
			}
		}
		if r.notFoundHandler != nil {
			return r.notFoundHandler(event)
		}
		return nil
	})

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handlerChain = r.middlewares[i](handlerChain)
	}

	return handlerChain(event)
}

// Use adds the given middlewares to the current Router.
func (r *Mux) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// With returns a new Router with the given middlewares.
func (r *Mux) With(middlewares ...Middleware) Router {
	return newRouter("", middlewares, nil)
}

// Group creates a new Router and adds it to the current Router.
func (r *Mux) Group(fn func(router Router)) {
	router := New()
	fn(router)
	r.handle(router)
}

// Route creates a new sub-router with the given pattern and adds it to the current Router.
func (r *Mux) Route(pattern string, fn func(r Router)) Router {
	checkPattern(pattern)
	router := newRouter(pattern, nil, nil)
	fn(router)
	r.handle(router)
	return router
}

// Mount mounts the given router with the given pattern to the current Router.
func (r *Mux) Mount(pattern string, router Router) {
	if pattern == "" {
		r.handle(router)
		return
	}
	r.handle(newRouter(pattern, nil, []Route{router}))
}

func (r *Mux) handle(route Route) {
	r.routes = append(r.routes, route)
}

// Interaction registers the given InteractionHandler to the current Router.
// This is a shortcut for Command, Autocomplete, Component and Modal.
func (r *Mux) Interaction(pattern string, h InteractionHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[InteractionHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionType(0),
	})
}

// Command registers the given CommandHandler to the current Router.
func (r *Mux) Command(pattern string, h CommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[CommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
	})
}

// SlashCommand registers the given SlashCommandHandler to the current Router.
func (r *Mux) SlashCommand(pattern string, h SlashCommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[SlashCommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
		t2:      []int{int(discord.ApplicationCommandTypeSlash)},
	})
}

// UserCommand registers the given UserCommandHandler to the current Router.
func (r *Mux) UserCommand(pattern string, h UserCommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[UserCommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
		t2:      []int{int(discord.ApplicationCommandTypeUser)},
	})
}

// MessageCommand registers the given MessageCommandHandler to the current Router.
func (r *Mux) MessageCommand(pattern string, h MessageCommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[MessageCommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
		t2:      []int{int(discord.ApplicationCommandTypeMessage)},
	})
}

// EntryPointCommand registers the given EntryPointCommandHandler to the current Router.
func (r *Mux) EntryPointCommand(pattern string, h EntryPointCommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[EntryPointCommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
		t2:      []int{int(discord.ApplicationCommandTypePrimaryEntryPoint)},
	})
}

// Autocomplete registers the given AutocompleteHandler to the current Router.
func (r *Mux) Autocomplete(pattern string, h AutocompleteHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[AutocompleteHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeAutocomplete,
	})
}

// Component registers the given ComponentHandler to the current Router.
func (r *Mux) Component(pattern string, h ComponentHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[ComponentHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeComponent,
	})
}

// ButtonComponent registers the given ButtonComponentHandler to the current Router.
func (r *Mux) ButtonComponent(pattern string, h ButtonComponentHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[ButtonComponentHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeComponent,
		t2:      []int{int(discord.ComponentTypeButton)},
	})
}

// SelectMenuComponent registers the given SelectMenuComponentHandler to the current Router.
func (r *Mux) SelectMenuComponent(pattern string, h SelectMenuComponentHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[SelectMenuComponentHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeComponent,
		t2: []int{
			int(discord.ComponentTypeStringSelectMenu),
			int(discord.ComponentTypeUserSelectMenu),
			int(discord.ComponentTypeRoleSelectMenu),
			int(discord.ComponentTypeMentionableSelectMenu),
			int(discord.ComponentTypeChannelSelectMenu),
		},
	})
}

// Modal registers the given ModalHandler to the current Router.
func (r *Mux) Modal(pattern string, h ModalHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[ModalHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeModalSubmit,
	})
}

// NotFound sets the NotFoundHandler for this router.
// This handler only works for the root router and will be ignored for sub routers.
func (r *Mux) NotFound(h NotFoundHandler) {
	r.notFoundHandler = h
}

// Error sets the ErrorHandler for this router.
// This handler only works for the root router and will be ignored for sub routers.
func (r *Mux) Error(h ErrorHandler) {
	r.errorHandler = h
}

// DefaultContext sets the default context for this router.
// This context will be used for all interaction events.
func (r *Mux) DefaultContext(ctx func() context.Context) {
	r.defaultContext = ctx
}

func checkPattern(pattern string) {
	if len(pattern) == 0 {
		panic("pattern must not be empty")
	}
	if pattern[0] != '/' {
		panic("pattern must start with /")
	}
}

func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
}

func parseVariables(path string, pattern string, variables map[string]string) string {
	if pattern == "" {
		return path
	}
	parts := splitPath(path)
	patternParts := splitPath(pattern)

	for i := range patternParts {
		path = strings.TrimPrefix(path, "/"+parts[i])
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			variables[patternParts[i][1:len(patternParts[i])-1]] = parts[i]
		}
	}
	return path
}
