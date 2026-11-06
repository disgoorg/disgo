package gateway

import (
	"context"
	"testing"
)

type lockCheckTransport struct{}

func (lockCheckTransport) ReceiveMessage() (*Message, error) { return nil, nil }
func (lockCheckTransport) WriteMessage(Message) error        { return nil }
func (lockCheckTransport) WriteClose(int, string) error      { return nil }
func (lockCheckTransport) Close() error                      { return nil }

// lockCheckRateLimiter runs inside sendInternal, while connMu is held.
type lockCheckRateLimiter struct {
	g            *gatewayImpl
	statusLocked bool
}

func (r *lockCheckRateLimiter) Close(context.Context) {}
func (r *lockCheckRateLimiter) Reset()                {}
func (r *lockCheckRateLimiter) Unlock()               {}
func (r *lockCheckRateLimiter) Wait(context.Context, RateLimiterCommandType) error {
	if r.g.statusMu.TryLock() {
		r.g.statusMu.Unlock()
	} else {
		r.statusLocked = true
	}
	return nil
}

// CloseWithCode and open lock connMu before statusMu. Send must not hold
// statusMu while it waits for connMu, or it deadlocks against them.
func TestSendDoesNotHoldStatusMuWithConnMu(t *testing.T) {
	rl := &lockCheckRateLimiter{}
	g := New("token", nil, WithRateLimiter(rl)).(*gatewayImpl)
	rl.g = g
	g.status = StatusReady
	g.conn = lockCheckTransport{}

	if err := g.Send(context.Background(), OpcodePresenceUpdate, MessageDataPresenceUpdate{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl.statusLocked {
		t.Fatal("Send held statusMu while holding connMu")
	}
}
