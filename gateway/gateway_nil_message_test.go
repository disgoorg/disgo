package gateway

import (
	"testing"
	"time"
)

type nilMessageTransport struct{}

func (nilMessageTransport) ReceiveMessage() (*Message, error) { return nil, nil }
func (nilMessageTransport) WriteMessage(Message) error        { return nil }
func (nilMessageTransport) WriteClose(int, string) error      { return nil }
func (nilMessageTransport) Close() error                      { return nil }

// A nil message before the gateway is ready means a handshake message (HELLO/READY)
// failed to decode. Waiting for it would never complete and would leave the shard
// half-connected, so listen must fail the open instead of looping.
func TestListenFailsOpenOnUndecodableHandshake(t *testing.T) {
	g := New("token", nil).(*gatewayImpl)
	g.conn = nilMessageTransport{}
	g.status = StatusWaitingForReady

	errCh := make(chan error, 2)
	returned := make(chan struct{})
	go func() {
		g.listen(nilMessageTransport{}, func(err error) { errCh <- err })
		close(returned)
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected listen to fail with an error on undecodable handshake message")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("listen did not call ready after undecodable handshake message")
	}
	select {
	case <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("listen did not return after undecodable handshake message")
	}
}

// While not ready, a nil message must abort the listen loop instead of continuing.
func TestListenReturnsOnNilMessageBeforeReady(t *testing.T) {
	g := New("token", nil).(*gatewayImpl)
	g.status = StatusWaitingForHello

	returned := make(chan struct{})
	go func() {
		g.listen(nilMessageTransport{}, func(error) {})
		close(returned)
	}()

	select {
	case <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("listen kept looping on nil messages before being ready")
	}
}
