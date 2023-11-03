package rpc

type internalHandler struct {
	handler Handler
	errChan chan error
}

type Handler interface {
	Handle(data MessageData)
}

func NewHandler[T MessageData](handler func(data T)) Handler {
	return &defaultHandler[T]{
		handler: handler,
	}
}

type defaultHandler[T MessageData] struct {
	handler func(data T)
}

func (h *defaultHandler[T]) Handle(data MessageData) {
	if d, ok := data.(T); ok {
		h.handler(d)
	}
}
