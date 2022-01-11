package core

var _ Collector[Event] = (*defaultCollector[Event])(nil)

type CollectorFilter[T Event] func(T) bool

type Collector[T Event] interface {
	EventListener
}

// NewCollector gives you a channel to receive on and a function to close the collector
func NewCollector[T Event](disgo Bot, filter CollectorFilter[T]) (<-chan T, func()) {
	ch := make(chan T)

	collector := &defaultCollector[T]{
		Filter: filter,
		Chan:   ch,
	}
	cls := func() {
		close(ch)
		disgo.EventManager().RemoveEventListeners(collector)
	}
	collector.Close = cls
	disgo.EventManager().AddEventListeners(collector)

	return ch, cls
}

type defaultCollector[T Event] struct {
	Filter CollectorFilter[T]
	Chan   chan<- T
	Close  func()
}

func (c *defaultCollector[T]) OnEvent(e Event) {
	if event, ok := e.(T); ok {
		if !c.Filter(event) {
			return
		}
		c.Chan <- event
	}
}
