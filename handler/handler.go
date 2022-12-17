package handler

import (
	"sync"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type (
	CommandHandler      func(client bot.Client, event *CommandEvent) error
	AutocompleteHandler func(client bot.Client, event *AutocompleteEvent) error
	ComponentHandler    func(client bot.Client, event *ComponentEvent) error
	ModalHandler        func(client bot.Client, event *ModalEvent) error
)

type Handler interface {
	bot.EventListener

	HandleCommand(commandPath string, handler CommandHandler)
	HandleAutocomplete(commandPath string, handler AutocompleteHandler)
	HandleComponent(componentID string, handler ComponentHandler)
	HandleModal(modalID string, handler ModalHandler)
}

func New() Handler {
	return &handlerImpl{
		commandHandlers: make(map[string]CommandHandler),
	}
}

type handlerImpl struct {
	commandHandlers   map[string]CommandHandler
	commandHandlersMu sync.Mutex

	autocompleteHandlers map[string]AutocompleteHandler
	autocompleteMu       sync.Mutex

	componentHandlers   map[string]ComponentHandler
	componentHandlersMu sync.Mutex

	modalHandlers   map[string]ModalHandler
	modalHandlersMu sync.Mutex
}

func (h *handlerImpl) OnEvent(event bot.Event) {
	var err error
	switch e := event.(type) {
	case *events.ApplicationCommandInteractionCreate:
		path := "/" + e.Data.CommandName()
		if d, ok := e.Data.(discord.SlashCommandInteractionData); ok {
			path = d.CommandPath()
		}

		h.commandHandlersMu.Lock()
		handler, ok := h.commandHandlers[path]
		h.commandHandlersMu.Unlock()
		if ok {
			err = handler(e.Client(), &CommandEvent{ApplicationCommandInteractionCreate: e})
		}

	case *events.AutocompleteInteractionCreate:
		h.commandHandlersMu.Lock()
		handler, ok := h.autocompleteHandlers[e.Data.CommandName]
		h.commandHandlersMu.Unlock()
		if ok {
			err = handler(e.Client(), &AutocompleteEvent{AutocompleteInteractionCreate: e})
		}

	case *events.ComponentInteractionCreate:
		h.commandHandlersMu.Lock()
		handler, ok := h.componentHandlers[e.Data.CustomID()]
		h.commandHandlersMu.Unlock()
		if ok {
			err = handler(e.Client(), &ComponentEvent{ComponentInteractionCreate: e})
		}

	case *events.ModalSubmitInteractionCreate:
		h.commandHandlersMu.Lock()
		handler, ok := h.modalHandlers[e.Data.CustomID]
		h.commandHandlersMu.Unlock()
		if ok {
			err = handler(e.Client(), &ModalEvent{ModalSubmitInteractionCreate: e})
		}
	}

	if err != nil {
		event.Client().Logger().Error("error while handling event: ", err)
	}
}

func (h *handlerImpl) HandleCommand(commandPath string, handler CommandHandler) {
	h.commandHandlersMu.Lock()
	defer h.commandHandlersMu.Unlock()
	if handler == nil {
		delete(h.commandHandlers, commandPath)
		return
	}
	h.commandHandlers[commandPath] = handler
}

func (h *handlerImpl) HandleAutocomplete(commandPath string, handler AutocompleteHandler) {
	h.autocompleteMu.Lock()
	defer h.autocompleteMu.Unlock()
	if handler == nil {
		delete(h.autocompleteHandlers, commandPath)
		return
	}
	h.autocompleteHandlers[commandPath] = handler
}

func (h *handlerImpl) HandleComponent(componentID string, handler ComponentHandler) {
	h.componentHandlersMu.Lock()
	defer h.componentHandlersMu.Unlock()
	if handler == nil {
		delete(h.componentHandlers, componentID)
		return
	}
	h.componentHandlers[componentID] = handler
}

func (h *handlerImpl) HandleModal(modalID string, handler ModalHandler) {
	h.modalHandlersMu.Lock()
	defer h.modalHandlersMu.Unlock()
	if handler == nil {
		delete(h.modalHandlers, modalID)
		return
	}
	h.modalHandlers[modalID] = handler
}
