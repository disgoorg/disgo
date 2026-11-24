package bot

import (
	"context"
	"testing"
	"time"
)

type testEvent struct {
	client *Client
}

func (e *testEvent) Client() *Client {
	return e.client
}

func (e *testEvent) SequenceNumber() int {
	return 1
}

func TestWaitForEvent_ContextCancelNoDeadlock(t *testing.T) {
	client := &Client{
		EventManager: NewEventManager(nil),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		WaitForEvent(client, ctx, func(e *testEvent) bool {
			return true
		}, func(e *testEvent) {}, nil)
	}()

	// Cancel context to trigger cleanup
	cancel()

	// Dispatch an event concurrently while cancellation is in-flight
	client.EventManager.DispatchEvent(&testEvent{client: client})

	select {
	case <-done:
		// Success: no deadlock occurred
	case <-time.After(2 * time.Second):
		t.Fatal("WaitForEvent deadlocked on context cancellation")
	}
}

func TestEventManager_ReentrantListenerMutationNoDeadlock(t *testing.T) {
	em := NewEventManager(nil)

	var selfListener EventListener
	executed := false

	selfListener = NewListenerFunc(func(e *testEvent) {
		executed = true
		// Re-entrantly remove self during callback execution
		em.RemoveEventListeners(selfListener)
	})

	em.AddEventListeners(selfListener)

	done := make(chan struct{})
	go func() {
		defer close(done)
		em.DispatchEvent(&testEvent{})
	}()

	select {
	case <-done:
		if !executed {
			t.Fatal("expected listener to execute")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DispatchEvent deadlocked on re-entrant listener mutation")
	}
}
