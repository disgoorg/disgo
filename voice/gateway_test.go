package voice

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNextReconnectDelay(t *testing.T) {
	tests := []struct {
		delay time.Duration
		want  time.Duration
	}{
		{delay: 0, want: time.Second},
		{delay: time.Second, want: 2 * time.Second},
		{delay: 2 * time.Second, want: 4 * time.Second},
		{delay: 4 * time.Second, want: 8 * time.Second},
		{delay: 8 * time.Second, want: maximumConnectDelay},
		{delay: maximumConnectDelay, want: maximumConnectDelay},
	}

	for _, tt := range tests {
		if got := nextReconnectDelay(tt.delay); got != tt.want {
			t.Errorf("nextReconnectDelay(%v) = %v, want %v", tt.delay, got, tt.want)
		}
	}
}

func TestNonReconnectableCloseCodesDoNotStartNewConnection(t *testing.T) {
	for _, closeCode := range []GatewayCloseEventCode{
		GatewayCloseEventCodeDisconnected,
		GatewayCloseEventCodeRateLimited,
		GatewayCloseEventCodeCallTerminated,
	} {
		if closeCode.Reconnect || closeCode.NewConnection {
			t.Errorf("close code %d must not reconnect: %+v", closeCode.Code, closeCode)
		}
	}
}

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
