package voice

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"slices"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
	botgateway "github.com/disgoorg/disgo/gateway"
)

func TestHandleGatewayCloseStartsFreshConnection(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "invalid session",
			err:  &websocket.CloseError{Code: GatewayCloseEventCodeSessionNoLongerValid.Code},
		},
		{
			name: "resume attempts exhausted",
			err:  errors.New("voice endpoint is unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := &gatewayStub{}
			udp := &udpConnStub{}
			var updates []*snowflake.ID
			var stateAtUpdate State
			var removed bool
			var conn *connImpl

			conn = &connImpl{
				config: defaultConnConfig(),
				state: State{
					GuildID:   1,
					UserID:    2,
					ChannelID: 3,
					SessionID: "old-session",
					Token:     "old-token",
					Endpoint:  "old.discord.media",
				},
				gateway: gateway,
				udp:     udp,
				voiceStateUpdateFunc: func(_ context.Context, _ snowflake.ID, channelID *snowflake.ID, _ bool, _ bool) error {
					updates = append(updates, channelID)
					if channelID != nil {
						stateAtUpdate = conn.state
						conn.openedFunc()
					}
					return nil
				},
				removeConnFunc: func() {
					removed = true
				},
			}

			conn.handleGatewayClose(gateway, tt.err)

			if len(updates) != 1 || updates[0] == nil || *updates[0] != 3 {
				t.Fatalf("expected one voice state update to channel 3, got %v", updates)
			}
			if stateAtUpdate.SessionID != "" || stateAtUpdate.Token != "" || stateAtUpdate.Endpoint != "" {
				t.Fatalf("expected cached voice credentials to be cleared, got %+v", stateAtUpdate)
			}
			if gateway.closeCalls != 0 {
				t.Errorf("expected reopened gateway to remain open, got %d close calls", gateway.closeCalls)
			}
			if udp.closeCalls != 1 {
				t.Errorf("expected stale UDP connection to be closed once, got %d close calls", udp.closeCalls)
			}
			if removed {
				t.Error("expected connection to remain registered")
			}
		})
	}
}

func TestVoiceGatewayWaitsForFreshStateAndServerUpdates(t *testing.T) {
	opened := make(chan State, 1)
	conn := &connImpl{
		config: defaultConnConfig(),
		state: State{
			GuildID:   1,
			UserID:    2,
			ChannelID: 3,
		},
		gateway: &gatewayStub{opened: opened},
	}
	endpoint := "new.discord.media"

	conn.HandleVoiceServerUpdate(botgateway.EventVoiceServerUpdate{
		GuildID:  1,
		Token:    "new-token",
		Endpoint: &endpoint,
	})

	select {
	case state := <-opened:
		t.Fatalf("gateway opened before fresh voice state update: %+v", state)
	case <-time.After(20 * time.Millisecond):
	}

	channelID := snowflake.ID(3)
	conn.HandleVoiceStateUpdate(botgateway.EventVoiceStateUpdate{VoiceState: discord.VoiceState{
		GuildID:   1,
		UserID:    2,
		ChannelID: &channelID,
		SessionID: "new-session",
	}})

	select {
	case state := <-opened:
		if state.SessionID != "new-session" || state.Token != "new-token" || state.Endpoint != endpoint {
			t.Fatalf("gateway opened with incomplete fresh state: %+v", state)
		}
	case <-time.After(time.Second):
		t.Fatal("gateway did not open after both fresh updates arrived")
	}
}

type gatewayStub struct {
	closeCalls int
	opened     chan State
}

func (*gatewayStub) SSRC() uint32 { return 0 }

func (g *gatewayStub) Open(_ context.Context, state State) error {
	if g.opened != nil {
		g.opened <- state
	}
	return nil
}

func (g *gatewayStub) Close() { g.closeCalls++ }

func (*gatewayStub) CloseWithCode(int, string) {}

func (*gatewayStub) Status() Status { return StatusDisconnected }

func (*gatewayStub) Send(context.Context, Opcode, GatewayMessageData) error { return nil }

func (*gatewayStub) Latency() time.Duration { return 0 }

type udpConnStub struct {
	closeCalls   int
	opened       bool
	openErr      error
	packet       *Packet
	secretKeySet bool
	secretKeyErr error
}

func (*udpConnStub) LocalAddr() net.Addr { return nil }

func (*udpConnStub) RemoteAddr() net.Addr { return nil }

func (u *udpConnStub) SetSecretKey(EncryptionMode, []byte) error {
	if u.secretKeyErr != nil {
		return u.secretKeyErr
	}
	u.secretKeySet = true
	return nil
}

func (*udpConnStub) SetDeadline(time.Time) error { return nil }

func (*udpConnStub) SetReadDeadline(time.Time) error { return nil }

func (*udpConnStub) SetWriteDeadline(time.Time) error { return nil }

func (u *udpConnStub) Open(context.Context, string, int, uint32) (string, int, error) {
	if u.openErr != nil {
		u.opened = false
		return "", 0, u.openErr
	}
	u.opened = true
	return "", 0, nil
}

// usable reports whether the socket is dialed and the encrypter built.
func (u *udpConnStub) usable() bool { return u.opened && u.secretKeySet }

func (u *udpConnStub) Close() error {
	u.closeCalls++
	return nil
}

func (*udpConnStub) Read([]byte) (int, error) { return 0, nil }

func (u *udpConnStub) ReadPacket() (*Packet, error) { return u.packet, nil }

func (*udpConnStub) Write([]byte) (int, error) { return 0, nil }

// SetOpusFrameProvider and SetOpusFrameReceiver may run before Conn.Open. The
// UDP socket is dialed on the gateway's ready event and its encrypter built
// from the session description, so a loop started by the setter has neither.
func TestAudioLoopsWaitForUDPConn(t *testing.T) {
	conn, udp, rec := newAudioConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	assertLoopsOpened(t, "before any gateway message", rec, 0)

	conn.handleMessage(conn.gateway, OpcodeHello, 0, GatewayMessageDataHello{})
	assertLoopsOpened(t, "after hello", rec, 0)

	gatewayReady(conn)
	if !udp.opened {
		t.Fatal("ready did not open the udp conn")
	}
	assertLoopsOpened(t, "after ready, before the secret key", rec, 0)

	sessionDescription(conn)
	assertLoopsOpened(t, "after session description", rec, 1)
}

// The session description alone does not mean the socket exists.
func TestAudioLoopsWaitForReady(t *testing.T) {
	conn, _, rec := newAudioConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	sessionDescription(conn)

	assertLoopsOpened(t, "after a session description with no ready", rec, 0)
}

// A failed dial leaves the conn nil, so ready arriving is not enough either.
// Both orderings matter: the key may already be installed when the dial fails.
func TestAudioLoopsWaitForSuccessfulDial(t *testing.T) {
	tests := []struct {
		name     string
		keyFirst bool
	}{
		{name: "ready first"},
		{name: "session description first", keyFirst: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, udp, rec := newAudioConn()
			udp.openErr = errors.New("dial failed")

			conn.SetOpusFrameProvider(nopProvider{})
			conn.SetOpusFrameReceiver(nopReceiver{})
			if tt.keyFirst {
				sessionDescription(conn)
				gatewayReady(conn)
			} else {
				gatewayReady(conn)
				sessionDescription(conn)
			}

			assertLoopsOpened(t, "after a failed dial", rec, 0)
		})
	}
}

// A failed SetSecretKey leaves the encrypter nil.
func TestAudioLoopsWaitForSecretKey(t *testing.T) {
	conn, udp, rec := newAudioConn()
	udp.secretKeyErr = errors.New("unknown encryption mode")

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)

	assertLoopsOpened(t, "after a failed secret key", rec, 0)
}

// Repeated ready and session-description events must not restart running loops.
func TestDuplicateGatewayEventsDoNotRestartAudioLoops(t *testing.T) {
	conn, _, rec := newAudioConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)
	gatewayReady(conn)
	sessionDescription(conn)

	assertLoopsOpened(t, "after duplicate gateway events", rec, 1)
}

// A loop registered after the conn is usable must start at once.
func TestAudioLoopsStartWhenConnIsReady(t *testing.T) {
	conn, _, rec := newAudioConn()

	gatewayReady(conn)
	sessionDescription(conn)
	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})

	assertLoopsOpened(t, "registered after the conn was usable", rec, 1)
}

// Playback bots register only a provider and recorders only a receiver.
func TestAudioLoopStartsWithoutItsSibling(t *testing.T) {
	t.Run("provider only", func(t *testing.T) {
		conn, _, rec := newAudioConn()

		conn.SetOpusFrameProvider(nopProvider{})
		gatewayReady(conn)
		sessionDescription(conn)

		if rec.senderOpens != 1 || rec.receiverOpens != 0 {
			t.Errorf("open calls: sender=%d receiver=%d, want 1 and 0", rec.senderOpens, rec.receiverOpens)
		}
	})

	t.Run("receiver only", func(t *testing.T) {
		conn, _, rec := newAudioConn()

		conn.SetOpusFrameReceiver(nopReceiver{})
		gatewayReady(conn)
		sessionDescription(conn)

		if rec.receiverOpens != 1 || rec.senderOpens != 0 {
			t.Errorf("open calls: sender=%d receiver=%d, want 0 and 1", rec.senderOpens, rec.receiverOpens)
		}
	})
}

// Setters called after leaving the channel must not start loops.
func TestAudioLoopsDoNotStartAfterChannelLeave(t *testing.T) {
	conn, _, rec := newAudioConn()

	gatewayReady(conn)
	sessionDescription(conn)
	leaveChannel(conn)
	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})

	assertLoopsOpened(t, "after leaving the channel", rec, 0)
}

// Leaving drops the loops, so rejoining must not revive them.
func TestRejoinDoesNotReopenDroppedAudioLoops(t *testing.T) {
	conn, _, rec := newAudioConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)
	leaveChannel(conn)
	gatewayReady(conn)
	sessionDescription(conn)

	assertLoopsOpened(t, "after leaving and handshaking again", rec, 1)
}

// A gateway close re-handshakes on the same conn without tearing down the
// loops, so the re-handshake must not open them a second time either.
func TestAudioLoopsSurviveGatewayClose(t *testing.T) {
	conn, _, rec := newAudioConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)
	conn.handleGatewayClose(conn.gateway, &websocket.CloseError{Code: GatewayCloseEventCodeSessionNoLongerValid.Code})

	if conn.state.SessionID != "" {
		t.Fatal("the gateway close returned before starting a fresh connection")
	}
	if rec.senderCloses != 0 || rec.receiverCloses != 0 {
		t.Errorf("close calls after a gateway close: sender=%d receiver=%d, want 0 and 0", rec.senderCloses, rec.receiverCloses)
	}

	gatewayReady(conn)
	sessionDescription(conn)
	assertLoopsOpened(t, "after a gateway close and re-handshake", rec, 1)
}

// Re-arming a precondition must not reopen a loop that is still registered:
// Open would overwrite its cancel func and leave the first goroutine running.
func TestReconnectDoesNotReopenRunningAudioLoops(t *testing.T) {
	t.Run("failed then successful re-dial", func(t *testing.T) {
		conn, udp, rec := newAudioConn()

		conn.SetOpusFrameProvider(nopProvider{})
		conn.SetOpusFrameReceiver(nopReceiver{})
		gatewayReady(conn)
		sessionDescription(conn)

		udp.openErr = errors.New("dial failed")
		gatewayReady(conn)
		udp.openErr = nil
		gatewayReady(conn)

		assertLoopsOpened(t, "after a failed then successful re-dial", rec, 1)
	})

	t.Run("close then handshake", func(t *testing.T) {
		conn, _, rec := newAudioConn()

		conn.SetOpusFrameProvider(nopProvider{})
		conn.SetOpusFrameReceiver(nopReceiver{})
		gatewayReady(conn)
		sessionDescription(conn)

		conn.Close(context.Background())
		gatewayReady(conn)
		sessionDescription(conn)

		assertLoopsOpened(t, "after closing then handshaking again", rec, 1)
	})
}

// A dropped loop must be closed whether or not it ever started, and the
// replacement must be the instance that ends up running.
func TestDroppedAudioLoopIsClosed(t *testing.T) {
	tests := []struct {
		name  string
		ready bool
		seq   []string
	}{
		{
			name:  "started",
			ready: true,
			seq:   []string{"sender.open", "receiver.open", "sender.close", "sender.open", "receiver.close", "receiver.open"},
		},
		{
			name: "unstarted",
			seq:  []string{"sender.close", "receiver.close"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, _, rec := newAudioConn()
			if tt.ready {
				gatewayReady(conn)
				sessionDescription(conn)
			}

			conn.SetOpusFrameProvider(nopProvider{})
			conn.SetOpusFrameReceiver(nopReceiver{})
			conn.SetOpusFrameProvider(nopProvider{})
			conn.SetOpusFrameReceiver(nopReceiver{})

			if rec.senderCloses != 1 || rec.receiverCloses != 1 {
				t.Errorf("close calls: sender=%d receiver=%d, want 1 and 1", rec.senderCloses, rec.receiverCloses)
			}
			if !slices.Equal(rec.seq, tt.seq) {
				t.Errorf("call sequence = %v, want %v", rec.seq, tt.seq)
			}
			if rec.openedAfterClose {
				t.Error("a loop was opened after it had been closed")
			}
		})
	}
}

// A straggling client disconnect can arrive after the loops are gone.
func TestClientDisconnectAfterLeaveIsSafe(t *testing.T) {
	conn, _, _ := newAudioConn()

	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)
	leaveChannel(conn)

	conn.handleMessage(conn.gateway, OpcodeClientDisconnect, 0, GatewayMessageDataClientDisconnect{UserID: 9})
}

// The real loops, not the stubs: once the conn is usable the sender must
// actually run, and closing it must actually stop it.
func TestRealAudioLoopsRunOnceConnIsUsable(t *testing.T) {
	udp := &udpConnStub{}
	conn := newRealConn()
	conn.udp = udp
	conn.state.ChannelID = 3

	pulled := make(chan struct{}, 64)
	conn.SetOpusFrameProvider(signalProvider{pulled: pulled})
	gatewayReady(conn)
	sessionDescription(conn)

	select {
	case <-pulled:
	case <-time.After(2 * time.Second):
		t.Fatal("the sender never pulled a frame after the conn became usable")
	}

	leaveChannel(conn)
	for len(pulled) > 0 {
		<-pulled
	}

	select {
	case <-pulled:
		t.Fatal("the sender pulled a frame after it was closed")
	case <-time.After(200 * time.Millisecond):
	}
}

// Close can land before the loop goroutine has run, so the cancel func must be
// in place before Open spawns it.
func TestClosingRealAudioLoopImmediatelyStopsIt(t *testing.T) {
	conn := newRealConn()
	conn.udp = &udpConnStub{}
	conn.state.ChannelID = 3

	pulled := make(chan struct{}, 64)
	conn.SetOpusFrameProvider(signalProvider{pulled: pulled})
	gatewayReady(conn)
	sessionDescription(conn)
	leaveChannel(conn)

	time.Sleep(100 * time.Millisecond)
	for len(pulled) > 0 {
		<-pulled
	}

	select {
	case <-pulled:
		t.Fatal("the sender kept pulling frames after an immediate close")
	case <-time.After(200 * time.Millisecond):
	}
}

// The receiver's cancel func must stop its loop, the same as the sender's.
func TestClosingRealAudioReceiverStopsIt(t *testing.T) {
	conn := newRealConn()
	conn.udp = &udpConnStub{packet: &Packet{SSRC: 1}}
	conn.state.ChannelID = 3

	received := make(chan struct{}, 64)
	conn.SetOpusFrameReceiver(signalReceiver{received: received})
	gatewayReady(conn)
	sessionDescription(conn)

	select {
	case <-received:
	case <-time.After(2 * time.Second):
		t.Fatal("the receiver never handled a packet after the conn became usable")
	}

	leaveChannel(conn)
	time.Sleep(100 * time.Millisecond)
	for len(received) > 0 {
		<-received
	}

	select {
	case <-received:
		t.Fatal("the receiver kept handling packets after it was closed")
	case <-time.After(200 * time.Millisecond):
	}
}

// A client disconnect has to reach the receiver so it can drop that user.
func TestClientDisconnectCleansUpUser(t *testing.T) {
	conn, _, rec := newAudioConn()

	conn.SetOpusFrameReceiver(nopReceiver{})
	gatewayReady(conn)
	sessionDescription(conn)
	conn.handleMessage(conn.gateway, OpcodeClientDisconnect, 0, GatewayMessageDataClientDisconnect{UserID: 9})

	if !slices.Equal(rec.cleanedUp, []snowflake.ID{9}) {
		t.Errorf("cleaned up = %v, want [9]", rec.cleanedUp)
	}
}

// Closing a loop that never started must not call a nil cancel func.
func TestClosingUnstartedRealAudioLoopIsSafe(t *testing.T) {
	conn := newRealConn()

	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
	conn.SetOpusFrameProvider(nopProvider{})
	conn.SetOpusFrameReceiver(nopReceiver{})
}

// Relies on the default DAVE session reporting ready; if that changes this
// test goes vacuous.
func TestRealAudioLoopsDoNotRunBeforeOpen(t *testing.T) {
	conn := newRealConn()

	pulled := make(chan struct{}, 1)
	conn.SetOpusFrameProvider(signalProvider{pulled: pulled})
	conn.SetOpusFrameReceiver(nopReceiver{})

	select {
	case <-pulled:
		t.Fatal("sender pulled a frame before the udp conn was open")
	case <-time.After(100 * time.Millisecond):
	}
}

func newRealConn() *connImpl {
	return NewConn(1, 2,
		func(context.Context, snowflake.ID, *snowflake.ID, bool, bool) error { return nil },
		func() {},
		WithConnLogger(slog.New(slog.DiscardHandler)),
		WithConnDaveSessionLogger(slog.New(slog.DiscardHandler)),
	).(*connImpl)
}

func newAudioConn() (*connImpl, *udpConnStub, *audioRecord) {
	udp := &udpConnStub{}
	rec := &audioRecord{udp: udp}

	var conn *connImpl
	conn = NewConn(1, 2,
		func(_ context.Context, _ snowflake.ID, channelID *snowflake.ID, _ bool, _ bool) error {
			if channelID != nil && conn.openedFunc != nil {
				conn.openedFunc()
			}
			return nil
		},
		func() {},
		WithConnLogger(slog.New(slog.DiscardHandler)),
		WithConnDaveSessionLogger(slog.New(slog.DiscardHandler)),
		WithConnAudioSenderCreateFunc(func(*slog.Logger, OpusFrameProvider, Conn) AudioSender {
			return &audioSenderStub{rec: rec}
		}),
		WithConnAudioReceiverCreateFunc(func(*slog.Logger, OpusFrameReceiver, Conn) AudioReceiver {
			return &audioReceiverStub{rec: rec}
		}),
	).(*connImpl)
	conn.udp = udp
	conn.state.ChannelID = 3

	return conn, udp, rec
}

func gatewayReady(conn *connImpl) {
	conn.handleMessage(conn.gateway, OpcodeReady, 0, GatewayMessageDataReady{
		SSRC:  1,
		IP:    "127.0.0.1",
		Port:  1234,
		Modes: AllEncryptionModes,
	})
}

func sessionDescription(conn *connImpl) {
	conn.handleMessage(conn.gateway, OpcodeSessionDescription, 0, GatewayMessageDataSessionDescription{
		Mode:      EncryptionModeAEADXChaCha20Poly1305RTPSize,
		SecretKey: make([]byte, 32),
	})
}

func leaveChannel(conn *connImpl) {
	conn.HandleVoiceStateUpdate(botgateway.EventVoiceStateUpdate{VoiceState: discord.VoiceState{
		GuildID: 1,
		UserID:  2,
	}})
}

func assertLoopsOpened(t *testing.T, when string, rec *audioRecord, want int) {
	t.Helper()

	if rec.senderOpens != want || rec.receiverOpens != want {
		t.Errorf("open calls %s: sender=%d receiver=%d, want %d and %d", when, rec.senderOpens, rec.receiverOpens, want, want)
	}
	if rec.openedUnusable {
		t.Errorf("a loop was opened while the udp conn was unusable %s", when)
	}
	if rec.openedAfterClose {
		t.Errorf("a loop was opened after it had been closed %s", when)
	}
}

// audioRecord accumulates across every loop instance the create funcs mint, so
// a test can tell a replacement apart from a reopened instance.
type audioRecord struct {
	udp              *udpConnStub
	senderOpens      int
	senderCloses     int
	receiverOpens    int
	receiverCloses   int
	openedUnusable   bool
	openedAfterClose bool
	cleanedUp        []snowflake.ID
	seq              []string
}

func (r *audioRecord) open(name string, closed bool) {
	r.seq = append(r.seq, name+".open")
	if !r.udp.usable() {
		r.openedUnusable = true
	}
	if closed {
		r.openedAfterClose = true
	}
}

type audioSenderStub struct {
	rec    *audioRecord
	closed bool
}

func (s *audioSenderStub) Open() {
	s.rec.senderOpens++
	s.rec.open("sender", s.closed)
}

func (s *audioSenderStub) Close() {
	s.closed = true
	s.rec.senderCloses++
	s.rec.seq = append(s.rec.seq, "sender.close")
}

type audioReceiverStub struct {
	rec    *audioRecord
	closed bool
}

func (s *audioReceiverStub) Open() {
	s.rec.receiverOpens++
	s.rec.open("receiver", s.closed)
}

func (s *audioReceiverStub) Close() {
	s.closed = true
	s.rec.receiverCloses++
	s.rec.seq = append(s.rec.seq, "receiver.close")
}

func (s *audioReceiverStub) CleanupUser(userID snowflake.ID) {
	s.rec.cleanedUp = append(s.rec.cleanedUp, userID)
}

type nopProvider struct{}

func (nopProvider) ProvideOpusFrame() ([]byte, error) { return nil, nil }
func (nopProvider) Close()                            {}

// signalProvider reports that the sender loop pulled a frame.
type signalProvider struct {
	pulled chan struct{}
}

func (p signalProvider) ProvideOpusFrame() ([]byte, error) {
	select {
	case p.pulled <- struct{}{}:
	default:
	}
	return nil, nil
}
func (signalProvider) Close() {}

// signalReceiver reports that the receiver loop handled a packet.
type signalReceiver struct {
	received chan struct{}
}

func (r signalReceiver) ReceiveOpusFrame(snowflake.ID, *Packet) error {
	select {
	case r.received <- struct{}{}:
	default:
	}
	return nil
}
func (signalReceiver) CleanupUser(snowflake.ID) {}
func (signalReceiver) Close()                   {}

type nopReceiver struct{}

func (nopReceiver) ReceiveOpusFrame(snowflake.ID, *Packet) error { return nil }
func (nopReceiver) CleanupUser(snowflake.ID)                     {}
func (nopReceiver) Close()                                       {}
