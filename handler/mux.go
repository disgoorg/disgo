package handler

import (
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var defaultErrorHandler = func(e *events.InteractionCreate, err error) {
	e.Client().Logger().Errorf("error handling interaction: %v\n", err)
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

	if err := r.Handle(path, make(map[string]string), e); err != nil {
		if r.errorHandler != nil {
			r.errorHandler(e, err)
			return
		}
		defaultErrorHandler(e, err)
	}
}

// Match returns true if the given path matches the Route.
func (r *Mux) Match(path string, t discord.InteractionType) bool {
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
		if matcher.Match(path, t) {
			return true
		}
	}
	return false
}

// Handle handles the given interaction event.
func (r *Mux) Handle(path string, variables map[string]string, e *events.InteractionCreate) error {
	handlerChain := func(event *events.InteractionCreate) error {
		path = parseVariables(path, r.pattern, variables)

		for _, route := range r.routes {
			if route.Match(path, e.Type()) {
				return route.Handle(path, variables, e)
			}
		}
		if r.notFoundHandler != nil {
			return r.notFoundHandler(e)
		}
		return nil
	}

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handlerChain = r.middlewares[i](handlerChain)
	}

	return handlerChain(e)
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

// Command registers the given CommandHandler to the current Router.
func (r *Mux) Command(pattern string, h CommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[CommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
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
