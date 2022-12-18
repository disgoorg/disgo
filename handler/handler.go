package handler

import (
	"sync"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	CommandDelimiter  = "/"
	CustomIDDelimiter = ":"
)

type (
	CommandHandler      func(client bot.Client, event *CommandEvent) error
	AutocompleteHandler func(client bot.Client, event *AutocompleteEvent) error
	ComponentHandler    func(client bot.Client, event *ComponentEvent) error
	ModalHandler        func(client bot.Client, event *ModalEvent) error
)

type Handler interface {
	bot.EventListener

	HandleCommand(commandPath string, handler CommandHandler) func()
	HandleAutocomplete(commandPath string, handler AutocompleteHandler) func()
	HandleComponent(customID string, handler ComponentHandler) func()
	HandleModal(customID string, handler ModalHandler) func()
}

func New() Handler {
	return &handlerImpl{}
}

type handlerImpl struct {
	commandHandlers   []*handlerHolder[CommandHandler]
	commandHandlersMu sync.Mutex

	autocompleteHandlers []*handlerHolder[AutocompleteHandler]
	autocompleteMu       sync.Mutex

	componentHandlers   []*handlerHolder[ComponentHandler]
	componentHandlersMu sync.Mutex

	modalHandlers   []*handlerHolder[ModalHandler]
	modalHandlersMu sync.Mutex
}

func (h *handlerImpl) OnEvent(event bot.Event) {
	var err error
	switch e := event.(type) {
	case *events.ApplicationCommandInteractionCreate:
		path := e.Data.CommandName()
		if d, ok := e.Data.(discord.SlashCommandInteractionData); ok {
			path = d.CommandPath()
		}

		handler, values, ok := findHandler(h.commandHandlers, path, CommandDelimiter)
		if !ok {
			return
		}
		err = handler(e.Client(), &CommandEvent{
			ApplicationCommandInteractionCreate: e,
			Variables:                           values,
		})

	case *events.AutocompleteInteractionCreate:
		handler, values, ok := findHandler(h.autocompleteHandlers, e.Data.CommandName, CommandDelimiter)
		if !ok {
			return
		}
		err = handler(e.Client(), &AutocompleteEvent{
			AutocompleteInteractionCreate: e,
			Variables:                     values,
		})

	case *events.ComponentInteractionCreate:
		handler, values, ok := findHandler(h.componentHandlers, e.Data.CustomID(), CustomIDDelimiter)
		if !ok {
			return
		}
		err = handler(e.Client(), &ComponentEvent{
			ComponentInteractionCreate: e,
			Variables:                  values,
		})

	case *events.ModalSubmitInteractionCreate:
		handler, values, ok := findHandler(h.modalHandlers, e.Data.CustomID, CustomIDDelimiter)
		if !ok {
			return
		}
		err = handler(e.Client(), &ModalEvent{
			ModalSubmitInteractionCreate: e,
			Variables:                    values,
		})
	}

	if err != nil {
		event.Client().Logger().Error("error while handling event: ", err)
	}
}

func (h *handlerImpl) HandleCommand(commandPath string, handler CommandHandler) func() {
	h.commandHandlersMu.Lock()
	defer h.commandHandlersMu.Unlock()

	holder := &handlerHolder[CommandHandler]{
		path:    commandPath,
		handler: handler,
	}

	h.commandHandlers = append(h.commandHandlers, holder)

	var once sync.Once
	return func() {
		once.Do(func() {
			h.modalHandlersMu.Lock()
			defer h.modalHandlersMu.Unlock()

			for i := range h.commandHandlers {
				if h.commandHandlers[i] == holder {
					h.commandHandlers = append(h.commandHandlers[:i], h.commandHandlers[i+1:]...)
					return
				}
			}
		})
	}
}

func (h *handlerImpl) HandleAutocomplete(commandPath string, handler AutocompleteHandler) func() {
	h.autocompleteMu.Lock()
	defer h.autocompleteMu.Unlock()

	holder := &handlerHolder[AutocompleteHandler]{
		path:    commandPath,
		handler: handler,
	}

	h.autocompleteHandlers = append(h.autocompleteHandlers, holder)

	var once sync.Once
	return func() {
		once.Do(func() {
			h.modalHandlersMu.Lock()
			defer h.modalHandlersMu.Unlock()

			for i := range h.autocompleteHandlers {
				if h.autocompleteHandlers[i] == holder {
					h.autocompleteHandlers = append(h.autocompleteHandlers[:i], h.autocompleteHandlers[i+1:]...)
					return
				}
			}
		})
	}
}

func (h *handlerImpl) HandleComponent(componentID string, handler ComponentHandler) func() {
	h.componentHandlersMu.Lock()
	defer h.componentHandlersMu.Unlock()

	holder := &handlerHolder[ComponentHandler]{
		path:    componentID,
		handler: handler,
	}

	h.componentHandlers = append(h.componentHandlers, holder)

	var once sync.Once
	return func() {
		once.Do(func() {
			h.modalHandlersMu.Lock()
			defer h.modalHandlersMu.Unlock()

			for i := range h.componentHandlers {
				if h.componentHandlers[i] == holder {
					h.componentHandlers = append(h.componentHandlers[:i], h.componentHandlers[i+1:]...)
					return
				}
			}
		})
	}
}

func (h *handlerImpl) HandleModal(modalID string, handler ModalHandler) func() {
	h.modalHandlersMu.Lock()
	defer h.modalHandlersMu.Unlock()

	holder := &handlerHolder[ModalHandler]{
		path:    modalID,
		handler: handler,
	}

	h.modalHandlers = append(h.modalHandlers, holder)

	var once sync.Once
	return func() {
		once.Do(func() {
			h.modalHandlersMu.Lock()
			defer h.modalHandlersMu.Unlock()

			for i := range h.modalHandlers {
				if h.modalHandlers[i] == holder {
					h.modalHandlers = append(h.modalHandlers[:i], h.modalHandlers[i+1:]...)
					return
				}
			}
		})
	}
}
