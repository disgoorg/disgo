package voice

import (
	"context"
	"log/slog"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// fakeAudioSender is a minimal AudioSender for testing.
type fakeAudioSender struct {
	closed atomic.Bool
}

func (f *fakeAudioSender) Open()  {}
func (f *fakeAudioSender) Close() { f.closed.Store(true) }

// fakeAudioReceiver is a minimal AudioReceiver for testing.
type fakeAudioReceiver struct {
	closed atomic.Bool
}

func (f *fakeAudioReceiver) Open()                      {}
func (f *fakeAudioReceiver) CleanupUser(_ snowflake.ID) {}
func (f *fakeAudioReceiver) Close()                     { f.closed.Store(true) }

// fakeOpusFrameProvider is a minimal OpusFrameProvider for testing.
type fakeOpusFrameProvider struct{}

func (f *fakeOpusFrameProvider) ProvideOpusFrame() ([]byte, error) { return nil, nil }
func (f *fakeOpusFrameProvider) Close()                            {}

// fakeOpusFrameReceiver is a minimal OpusFrameReceiver for testing.
type fakeOpusFrameReceiver struct{}

func (f *fakeOpusFrameReceiver) ReceiveOpusFrame(_ snowflake.ID, _ *Packet) error { return nil }
func (f *fakeOpusFrameReceiver) CleanupUser(_ snowflake.ID)                       {}
func (f *fakeOpusFrameReceiver) Close()                                           {}

// fakeUDPConn is a minimal UDPConn stub that satisfies the interface for testing.
type fakeUDPConn struct{}

func (f *fakeUDPConn) LocalAddr() net.Addr                           { return nil }
func (f *fakeUDPConn) RemoteAddr() net.Addr                          { return nil }
func (f *fakeUDPConn) SetSecretKey(_ EncryptionMode, _ []byte) error { return nil }
func (f *fakeUDPConn) SetDeadline(_ time.Time) error                 { return nil }
func (f *fakeUDPConn) SetReadDeadline(_ time.Time) error             { return nil }
func (f *fakeUDPConn) SetWriteDeadline(_ time.Time) error            { return nil }
func (f *fakeUDPConn) Open(_ context.Context, _ string, _ int, _ uint32) (string, int, error) {
	return "127.0.0.1", 0, nil
}
func (f *fakeUDPConn) Close() error                 { return nil }
func (f *fakeUDPConn) Read(_ []byte) (int, error)   { return 0, nil }
func (f *fakeUDPConn) ReadPacket() (*Packet, error) { return nil, nil }
func (f *fakeUDPConn) Write(_ []byte) (int, error)  { return 0, nil }

func TestHandleMessage_SessionDescription_RestartsAudioPipeline(t *testing.T) {
	var (
		senderCreateCount   atomic.Int32
		receiverCreateCount atomic.Int32
	)

	// Track each sender/receiver instance created.
	var lastSender *fakeAudioSender
	var lastReceiver *fakeAudioReceiver

	senderCreateFunc := func(logger *slog.Logger, provider OpusFrameProvider, conn Conn) AudioSender {
		senderCreateCount.Add(1)
		s := &fakeAudioSender{}
		lastSender = s
		return s
	}

	receiverCreateFunc := func(logger *slog.Logger, receiver OpusFrameReceiver, conn Conn) AudioReceiver {
		receiverCreateCount.Add(1)
		r := &fakeAudioReceiver{}
		lastReceiver = r
		return r
	}

	conn := &connImpl{
		config: connConfig{
			Logger:                  slog.Default(),
			AudioSenderCreateFunc:   senderCreateFunc,
			AudioReceiverCreateFunc: receiverCreateFunc,
		},
		openedChan: make(chan struct{}, 1),
		udp:        &fakeUDPConn{},
	}

	// Simulate user setting up audio pipeline before any reconnect.
	conn.SetOpusFrameProvider(&fakeOpusFrameProvider{})
	conn.SetOpusFrameReceiver(&fakeOpusFrameReceiver{})

	if c := senderCreateCount.Load(); c != 1 {
		t.Fatalf("expected 1 sender create after SetOpusFrameProvider, got %d", c)
	}
	if c := receiverCreateCount.Load(); c != 1 {
		t.Fatalf("expected 1 receiver create after SetOpusFrameReceiver, got %d", c)
	}

	firstSender := lastSender
	firstReceiver := lastReceiver

	// Simulate a SessionDescription arriving after reconnect (Identify path).
	// This should restart the audio sender and receiver.
	conn.handleMessage(nil, OpcodeSessionDescription, 0, GatewayMessageDataSessionDescription{
		Mode:      EncryptionModeAEADXChaCha20Poly1305RTPSize,
		SecretKey: make([]byte, 32),
	})

	// Verify the old sender/receiver were closed.
	if !firstSender.closed.Load() {
		t.Error("expected old audio sender to be closed after SessionDescription")
	}
	if !firstReceiver.closed.Load() {
		t.Error("expected old audio receiver to be closed after SessionDescription")
	}

	// Verify new instances were created.
	if c := senderCreateCount.Load(); c != 2 {
		t.Errorf("expected 2 total sender creates (initial + restart), got %d", c)
	}
	if c := receiverCreateCount.Load(); c != 2 {
		t.Errorf("expected 2 total receiver creates (initial + restart), got %d", c)
	}

	// Drain the openedChan signal.
	select {
	case <-conn.openedChan:
	default:
		t.Error("expected openedChan to be signaled after SessionDescription")
	}
}

func TestHandleMessage_SessionDescription_NoRestartWithoutPriorSetup(t *testing.T) {
	var (
		senderCreateCount   atomic.Int32
		receiverCreateCount atomic.Int32
	)

	senderCreateFunc := func(logger *slog.Logger, provider OpusFrameProvider, conn Conn) AudioSender {
		senderCreateCount.Add(1)
		return &fakeAudioSender{}
	}

	receiverCreateFunc := func(logger *slog.Logger, receiver OpusFrameReceiver, conn Conn) AudioReceiver {
		receiverCreateCount.Add(1)
		return &fakeAudioReceiver{}
	}

	conn := &connImpl{
		config: connConfig{
			Logger:                  slog.Default(),
			AudioSenderCreateFunc:   senderCreateFunc,
			AudioReceiverCreateFunc: receiverCreateFunc,
		},
		openedChan: make(chan struct{}, 1),
		udp:        &fakeUDPConn{},
	}

	// Do NOT call SetOpusFrameProvider/SetOpusFrameReceiver.
	// SessionDescription should not create any audio sender/receiver.
	conn.handleMessage(nil, OpcodeSessionDescription, 0, GatewayMessageDataSessionDescription{
		Mode:      EncryptionModeAEADXChaCha20Poly1305RTPSize,
		SecretKey: make([]byte, 32),
	})

	if c := senderCreateCount.Load(); c != 0 {
		t.Errorf("expected 0 sender creates without prior SetOpusFrameProvider, got %d", c)
	}
	if c := receiverCreateCount.Load(); c != 0 {
		t.Errorf("expected 0 receiver creates without prior SetOpusFrameReceiver, got %d", c)
	}
}
