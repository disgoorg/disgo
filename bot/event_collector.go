package bot

import (
	"context"
	"sync"
)

// WaitForEvent waits for an event passing the filterFunc and then calls the actionFunc. You can cancel this function with the passed context.Context and the cancelFunc gets called then.
func WaitForEvent[E Event](client Client, ctx context.Context, filterFunc func(e E) bool, actionFunc func(e E), cancelFunc func()) {
	ch, cancel := NewEventCollector(client, filterFunc)

	select {
	case <-ctx.Done():
		cancel()
		if cancelFunc != nil {
			cancelFunc()
		}
	case e := <-ch:
		cancel()
		if actionFunc != nil {
			actionFunc(e)
		}
	}
}

// NewEventCollector returns a channel in which the events of type T gets sent which pass the passed filter and a function which can be used to stop the event collector.
// The close function needs to be called to stop the event collector.
func NewEventCollector[E Event](client Client, filterFunc func(e E) bool) (<-chan E, func()) {
	ch := make(chan E)
	var once sync.Once

	coll := &collector[E]{
		FilterFunc: filterFunc,
		Chan:       ch,
	}
	client.EventManager().AddEventListeners(coll)

	return ch, func() {
		once.Do(func() {
			client.EventManager().RemoveEventListeners(coll)
			close(ch)
		})
	}
}

type collector[E Event] struct {
	FilterFunc func(e E) bool
	Chan       chan<- E
}

func (c *collector[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		if !c.FilterFunc(event) {
			return
		}
		c.Chan <- event
	}
}
