package gateway

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Closing the gateway must stop a reconnect loop sitting in its backoff,
// otherwise Close is ignored until the next connect attempt fails.
func TestCloseCancelsPendingReconnect(t *testing.T) {
	g := New("token", nil).(*gatewayImpl)

	ctx, cancel := context.WithCancel(context.Background())
	g.reconnectMu.Lock()
	g.reconnectCtx = ctx
	g.reconnectCancel = cancel
	g.reconnectMu.Unlock()

	closeCtx, closeCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer closeCancel()
	g.Close(closeCtx)

	g.reconnectMu.Lock()
	reconnectCtx := g.reconnectCtx
	g.reconnectMu.Unlock()

	if reconnectCtx != nil {
		select {
		case <-reconnectCtx.Done():
		case <-time.After(2 * time.Second):
			t.Fatal("Close did not cancel the pending reconnect context")
		}
	}
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("expected reconnect context to be canceled, got %v", ctx.Err())
	}
}
