package rpc

type internalHandler struct {
	handler CommandHandler
	errChan chan error
}

type CommandHandler interface {
	Handle(data MessageData)
}

func CmdHandler[T MessageData](handler func(data T)) CommandHandler {
	return &defaultCommandHandler[T]{
		handler: handler,
	}
}

type defaultCommandHandler[T MessageData] struct {
	handler func(data T)
}

func (h *defaultCommandHandler[T]) Handle(data MessageData) {
	if d, ok := data.(T); ok {
		h.handler(d)
	}
}
