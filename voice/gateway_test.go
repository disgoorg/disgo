package voice

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSendDoesNotBlockStatusWhileWaitingForConnection(t *testing.T) {
	gateway := &gatewayImpl{status: StatusReady}
	gateway.connMu.Lock()

	sendStarted := make(chan struct{})
	sendDone := make(chan error, 1)
	go func() {
		close(sendStarted)
		sendDone <- gateway.Send(context.Background(), OpcodeHeartbeat, GatewayMessageDataHeartbeat{})
	}()
	<-sendStarted

	// Give Send time to reach the connection lock held by this test.
	time.Sleep(10 * time.Millisecond)

	statusDone := make(chan Status, 1)
	go func() {
		statusDone <- gateway.Status()
	}()

	statusBlocked := false
	select {
	case status := <-statusDone:
		if status != StatusReady {
			t.Errorf("unexpected gateway status: %v", status)
		}
	case <-time.After(100 * time.Millisecond):
		statusBlocked = true
	}

	gateway.connMu.Unlock()

	if err := <-sendDone; !errors.Is(err, ErrGatewayNotConnected) {
		t.Errorf("expected ErrGatewayNotConnected, got %v", err)
	}
	if statusBlocked {
		<-statusDone
		t.Fatal("Status was blocked by Send while Send was waiting for the connection lock")
	}
}
