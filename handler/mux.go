package handler

import (
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func New() Router {
	return &mux{}
}

func newRouter(pattern string, middlewares []Middleware, routes []Route) *mux {
	return &mux{
		pattern:     pattern,
		middlewares: middlewares,
		routes:      routes,
	}
}

type mux struct {
	pattern         string
	middlewares     []Middleware
	routes          []Route
	notFoundHandler NotFoundHandler
}

func (r *mux) OnEvent(event bot.Event) {
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
		event.Client().Logger().Errorf("error handling interaction: %v\n", err)
	}
}

func (r *mux) Match(path string, t discord.InteractionType) bool {
	if r.pattern != "" {
		parts := splitPath(path)
		patternParts := splitPath(r.pattern)

		for i, part := range patternParts {
			path = strings.TrimPrefix(path, "/"+parts[i])
			if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
				continue
			}
			if part != parts[i] {
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

func (r *mux) Handle(path string, variables map[string]string, e *events.InteractionCreate) error {
	path = parseVariables(path, r.pattern, variables)
	middlewares := func(event *events.InteractionCreate) {}
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		middlewares = r.middlewares[i](middlewares)
	}
	middlewares(e)

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

func (r *mux) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *mux) With(middlewares ...Middleware) Router {
	return newRouter("", middlewares, nil)
}

func (r *mux) Group(fn func(router Router)) {
	router := New()
	fn(router)
	r.handle(router)
}

func (r *mux) Route(pattern string, fn func(r Router)) Router {
	checkPattern(pattern)
	router := newRouter(pattern, nil, nil)
	fn(router)
	r.handle(router)
	return router
}

func (r *mux) Mount(pattern string, router Router) {
	if pattern == "" {
		r.handle(router)
		return
	}
	r.handle(newRouter(pattern, nil, []Route{router}))
}

func (r *mux) handle(route Route) {
	r.routes = append(r.routes, route)
}

func (r *mux) Command(pattern string, h CommandHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[CommandHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeApplicationCommand,
	})
}

func (r *mux) Autocomplete(pattern string, h AutocompleteHandler) {
	checkPattern(pattern)
	r.handle(&handlerHolder[AutocompleteHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeAutocomplete,
	})
}

func (r *mux) Component(pattern string, h ComponentHandler) {
	checkPatternEmpty(pattern)
	r.handle(&handlerHolder[ComponentHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeComponent,
	})
}

func (r *mux) Modal(pattern string, h ModalHandler) {
	checkPatternEmpty(pattern)
	r.handle(&handlerHolder[ModalHandler]{
		pattern: pattern,
		handler: h,
		t:       discord.InteractionTypeModalSubmit,
	})
}

func (r *mux) NotFound(h NotFoundHandler) {
	r.notFoundHandler = h
}

func checkPatternEmpty(pattern string) {
	if pattern == "" {
		panic("pattern must not be empty")
	}
}

func checkPattern(pattern string) {
	checkPatternEmpty(pattern)
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
