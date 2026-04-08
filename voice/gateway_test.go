package voice

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"
)

func newTestGateway() *gatewayImpl {
	cfg := defaultGatewayConfig()
	cfg.AutoReconnect = true
	gw := &gatewayImpl{config: cfg}
	gw.openFunc = gw.open // default; tests override this
	return gw
}

func testState() State {
	channelID := snowflake.ID(1)
	return State{
		GuildID:   snowflake.ID(1),
		UserID:    snowflake.ID(2),
		ChannelID: &channelID,
		SessionID: "test-session",
		Token:     "test-token",
		Endpoint:  "localhost:1234",
	}
}

func TestDoReconnect_FallsBackToIdentifyAfter4006(t *testing.T) {
	gw := newTestGateway()

	// Pre-set resume state to simulate a previously established session.
	gw.ssrc = 12345
	gw.seq = 10

	var calls atomic.Int32

	gw.openFunc = func(ctx context.Context, state State) error {
		n := calls.Add(1)
		switch n {
		case 1:
			// First call: simulate what open() does when Discord rejects
			// a Resume with 4006. open() clears ssrc/seq via Close(), then
			// returns the close error.
			gw.ssrc = 0
			gw.seq = 0
			return &websocket.CloseError{
				Code: GatewayCloseEventCodeSessionNoLongerValid.Code,
				Text: "Session is no longer valid.",
			}
		case 2:
			// Second call: the retry. Verify resume state was cleared so
			// a real open() would send Identify. Simulate success.
			if gw.ssrc != 0 || gw.seq != 0 {
				t.Errorf("expected ssrc=0 seq=0 on retry, got ssrc=%d seq=%d", gw.ssrc, gw.seq)
			}
			return nil
		default:
			t.Fatalf("unexpected call %d to openFunc", n)
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := gw.doReconnect(ctx, testState())
	if err != nil {
		t.Fatalf("doReconnect returned error: %v", err)
	}
	if c := calls.Load(); c != 2 {
		t.Errorf("expected 2 calls to openFunc, got %d", c)
	}
}

func TestDoReconnect_StillFailsOnTrulyNonReconnectableCodes(t *testing.T) {
	testCases := []struct {
		name string
		code int
	}{
		{"4003 Not Authenticated", 4003},
		{"4004 Authentication Failed", 4004},
		{"4009 Session Timeout", 4009},
		{"4011 Server Not Found", 4011},
		{"4014 Disconnected", 4014},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gw := newTestGateway()

			var calls atomic.Int32

			gw.openFunc = func(ctx context.Context, state State) error {
				calls.Add(1)
				return &websocket.CloseError{
					Code: tc.code,
					Text: fmt.Sprintf("close code %d", tc.code),
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := gw.doReconnect(ctx, testState())
			if err == nil {
				t.Fatalf("expected doReconnect to return error for close code %d, got nil", tc.code)
			}
			if c := calls.Load(); c != 1 {
				t.Errorf("expected 1 call to openFunc for code %d, got %d", tc.code, c)
			}
		})
	}
}
