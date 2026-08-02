package voice

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/disgoorg/godave"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
	botgateway "github.com/disgoorg/disgo/gateway"
)

type testGateway struct {
	mu          sync.Mutex
	status      Status
	statusCalls chan struct{}
	opened      chan State
	sent        chan GatewayMessageData
}

func (g *testGateway) SSRC() uint32 { return 0 }

func (g *testGateway) Open(ctx context.Context, state State) error {
	g.mu.Lock()
	g.status = StatusReady
	g.mu.Unlock()
	select {
	case g.opened <- state:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *testGateway) Close() {
	g.mu.Lock()
	g.status = StatusDisconnected
	g.mu.Unlock()
}

func (g *testGateway) CloseWithCode(int, string) { g.Close() }

func (g *testGateway) Status() Status {
	g.mu.Lock()
	defer g.mu.Unlock()
	select {
	case g.statusCalls <- struct{}{}:
	default:
	}
	return g.status
}

func (g *testGateway) Send(_ context.Context, _ Opcode, data GatewayMessageData) error {
	if g.sent != nil {
		g.sent <- data
	}
	return nil
}

func (*testGateway) Latency() time.Duration { return 0 }

type testDaveSession struct {
	godave.Session
	ready  bool
	closed bool
}

func (s *testDaveSession) Ready() bool {
	return s.ready
}

func (s *testDaveSession) Close() error {
	s.ready = false
	s.closed = true
	return nil
}

func newTestConn(gateway Gateway) *connImpl {
	return NewConn(1, 2, func(context.Context, snowflake.ID, *snowflake.ID, bool, bool) error {
		return nil
	}, func() {}, WithConnGatewayCreateFunc(func(godave.Session, EventHandlerFunc, CloseHandlerFunc, ...GatewayConfigOpt) Gateway {
		return gateway
	})).(*connImpl)
}

func TestHandleMessageReleasesSSRCMutexBeforeCallback(t *testing.T) {
	done := make(chan struct{})
	conn := &connImpl{}
	conn.config.EventHandlerFunc = func(Gateway, Opcode, int, GatewayMessageData) {
		conn.UserIDBySSRC(1)
		close(done)
	}

	go conn.handleMessage(nil, OpcodeClientDisconnect, 0, GatewayMessageDataClientDisconnect{})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("event callback blocked on ssrcsMu")
	}
}

func TestGatewayOpenWaitsForBothVoiceUpdates(t *testing.T) {
	tests := []struct {
		name        string
		serverFirst bool
	}{
		{name: "state first"},
		{name: "server first", serverFirst: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gateway := &testGateway{status: StatusUnconnected, statusCalls: make(chan struct{}, 1), opened: make(chan State, 1)}
			conn := newTestConn(gateway)
			channelID := snowflake.ID(3)
			endpoint := "voice.example"
			stateUpdate := botgateway.EventVoiceStateUpdate{VoiceState: discord.VoiceState{
				GuildID: 1, UserID: 2, ChannelID: &channelID, SessionID: "session",
			}}
			serverUpdate := botgateway.EventVoiceServerUpdate{GuildID: 1, Endpoint: &endpoint, Token: "token"}

			if test.serverFirst {
				conn.HandleVoiceServerUpdate(serverUpdate)
			} else {
				conn.HandleVoiceStateUpdate(stateUpdate)
			}
			<-gateway.statusCalls
			select {
			case <-gateway.opened:
				t.Fatal("gateway opened before both voice updates")
			default:
			}

			if test.serverFirst {
				conn.HandleVoiceStateUpdate(stateUpdate)
			} else {
				conn.HandleVoiceServerUpdate(serverUpdate)
			}

			select {
			case state := <-gateway.opened:
				if state.ChannelID != channelID || state.SessionID != "session" || state.Token != "token" {
					t.Fatalf("opened with incomplete state: %+v", state)
				}
			case <-time.After(time.Second):
				t.Fatal("gateway did not open after both voice updates")
			}
		})
	}
}

func TestClosePreventsQueuedGatewayOpen(t *testing.T) {
	gateway := &testGateway{status: StatusUnconnected, statusCalls: make(chan struct{}, 1), opened: make(chan State, 1)}
	conn := newTestConn(gateway)
	conn.state.ChannelID = 3
	conn.state.SessionID = "session"
	conn.state.Token = "token"
	conn.state.Endpoint = "voice.example"
	conn.voiceStateReceived = true
	conn.voiceServerReceived = true

	conn.gatewayOpenMu.Lock()
	conn.tryOpenGateway()
	conn.Close(context.Background())
	conn.gatewayOpenMu.Unlock()

	select {
	case <-gateway.statusCalls:
	case <-time.After(time.Second):
		t.Fatal("queued gateway open did not run")
	}
	select {
	case <-gateway.opened:
		t.Fatal("gateway opened after the connection was closed")
	default:
	}

	if err := conn.Open(context.Background(), 3, false, false); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("expected net.ErrClosed, got %v", err)
	}
}

func TestVoiceStateChangesRecreateVoiceTransports(t *testing.T) {
	var (
		daveSessions    []*testDaveSession
		gatewaySessions []godave.Session
		udpSessions     []godave.Session
		gateways        []*testGateway
	)

	conn := NewConn(1, 2, func(context.Context, snowflake.ID, *snowflake.ID, bool, bool) error {
		return nil
	}, func() {},
		WithConnDaveSessionCreateFunc(func(*slog.Logger, godave.UserID, godave.Callbacks) godave.Session {
			session := &testDaveSession{ready: len(daveSessions) == 0}
			daveSessions = append(daveSessions, session)
			return session
		}),
		WithConnGatewayCreateFunc(func(session godave.Session, _ EventHandlerFunc, _ CloseHandlerFunc, _ ...GatewayConfigOpt) Gateway {
			gatewaySessions = append(gatewaySessions, session)
			gateway := &testGateway{
				status:      StatusUnconnected,
				statusCalls: make(chan struct{}, 1),
				opened:      make(chan State, 1),
				sent:        make(chan GatewayMessageData, 1),
			}
			gateways = append(gateways, gateway)
			return gateway
		}),
		WithUDPConnCreateFunc(func(session godave.Session, lookup SsrcLookupFunc, opts ...UDPConnConfigOpt) UDPConn {
			udpSessions = append(udpSessions, session)
			return NewUDPConn(session, lookup, opts...)
		}),
	).(*connImpl)
	conn.state.ChannelID = 3
	conn.state.SessionID = "session"
	if err := conn.SetSpeaking(context.Background(), SpeakingFlagMicrophone); err != nil {
		t.Fatalf("failed to set initial speaking state: %v", err)
	}
	<-gateways[0].sent

	channelID := snowflake.ID(4)
	conn.HandleVoiceStateUpdate(botgateway.EventVoiceStateUpdate{VoiceState: discord.VoiceState{
		GuildID: 1, UserID: 2, ChannelID: &channelID, SessionID: "session",
	}})

	select {
	case <-gateways[1].statusCalls:
	case <-time.After(time.Second):
		t.Fatal("queued gateway open decision did not run")
	}

	if len(daveSessions) != 2 || len(gatewaySessions) != 2 || len(udpSessions) != 2 {
		t.Fatalf("expected two transport generations, got DAVE=%d gateway=%d UDP=%d", len(daveSessions), len(gatewaySessions), len(udpSessions))
	}
	if !daveSessions[0].closed {
		t.Fatal("old DAVE session was not closed")
	}
	if daveSessions[1].closed || daveSessions[1].Ready() {
		t.Fatal("new DAVE session should be open and waiting for its handshake")
	}
	if conn.DAVE() != daveSessions[1] || gatewaySessions[1] != daveSessions[1] || udpSessions[1] != daveSessions[1] {
		t.Fatal("new voice transports do not share the fresh DAVE session")
	}

	conn.handleMessage(gateways[1], OpcodeSessionDescription, 0, GatewayMessageDataSessionDescription{
		Mode:      EncryptionModeAEADAES256GCMRTPSize,
		SecretKey: make([]byte, 32),
	})
	select {
	case data := <-gateways[1].sent:
		speaking, ok := data.(GatewayMessageDataSpeaking)
		if !ok || speaking.Speaking != SpeakingFlagMicrophone {
			t.Fatalf("expected microphone speaking state on the new gateway, got %#v", data)
		}
	case <-time.After(time.Second):
		t.Fatal("speaking state was not restored on the new gateway")
	}

	conn.HandleVoiceStateUpdate(botgateway.EventVoiceStateUpdate{VoiceState: discord.VoiceState{
		GuildID: 1, UserID: 2, ChannelID: nil, SessionID: "session",
	}})
	if len(daveSessions) != 3 || len(gatewaySessions) != 3 || len(udpSessions) != 3 {
		t.Fatalf("expected disconnect to create a third transport generation, got DAVE=%d gateway=%d UDP=%d", len(daveSessions), len(gatewaySessions), len(udpSessions))
	}
	if !daveSessions[1].closed {
		t.Fatal("disconnected DAVE session was not closed")
	}
	if conn.DAVE() != daveSessions[2] || gatewaySessions[2] != daveSessions[2] || udpSessions[2] != daveSessions[2] {
		t.Fatal("disconnect did not publish fresh voice transports")
	}
	if conn.speakingSet {
		t.Fatal("disconnect did not clear the stale speaking state")
	}

	conn.handleGatewayClose(gateways[1], &websocket.CloseError{
		Code: GatewayCloseEventCodeDisconnected.Code,
		Text: GatewayCloseEventCodeDisconnected.Description,
	})
	if conn.closed {
		t.Fatal("stale disconnected gateway closed the fresh connection")
	}
}

func TestStateDrivenGatewayClosesWaitForVoiceUpdates(t *testing.T) {
	stateUpdates := make(chan *snowflake.ID, 1)
	removed := make(chan struct{}, 1)
	gateway := &testGateway{status: StatusReady, statusCalls: make(chan struct{}, 1), opened: make(chan State, 1)}
	conn := NewConn(1, 2, func(_ context.Context, _ snowflake.ID, channelID *snowflake.ID, _ bool, _ bool) error {
		stateUpdates <- channelID
		return nil
	}, func() {
		removed <- struct{}{}
	}, WithConnGatewayCreateFunc(func(godave.Session, EventHandlerFunc, CloseHandlerFunc, ...GatewayConfigOpt) Gateway {
		return gateway
	})).(*connImpl)
	conn.state.ChannelID = 3
	conn.state.SessionID = "session"
	conn.state.Token = "token"
	conn.state.Endpoint = "voice.example"

	for _, closeCode := range []GatewayCloseEventCode{GatewayCloseEventCodeDisconnected, GatewayCloseEventCodeCallTerminated} {
		conn.handleGatewayClose(gateway, &websocket.CloseError{
			Code: closeCode.Code,
			Text: closeCode.Description,
		})
	}

	select {
	case channelID := <-stateUpdates:
		t.Fatalf("state-driven close sent an unsolicited voice state update for channel %v", channelID)
	default:
	}
	select {
	case <-removed:
		t.Fatal("state-driven close removed the connection before the voice updates arrived")
	default:
	}
	if conn.closed || conn.ChannelID() == nil || *conn.ChannelID() != 3 {
		t.Fatal("state-driven close changed the connection while waiting for voice updates")
	}
}
