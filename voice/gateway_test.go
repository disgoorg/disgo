package voice

import (
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
