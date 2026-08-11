package voice

import (
	"testing"
	"time"
)

func TestReconnectDelay(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 0, want: 0},
		{attempt: 1, want: time.Second},
		{attempt: 2, want: 2 * time.Second},
		{attempt: 3, want: 4 * time.Second},
		{attempt: 4, want: 8 * time.Second},
		{attempt: 5, want: maximumConnectDelay},
		{attempt: 100, want: maximumConnectDelay},
	}

	for _, tt := range tests {
		if got := reconnectDelay(tt.attempt); got != tt.want {
			t.Errorf("reconnectDelay(%d) = %v, want %v", tt.attempt, got, tt.want)
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
