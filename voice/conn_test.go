package voice

import (
	"context"
	"errors"
	"net"
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
	closeCalls int
}

func (*udpConnStub) LocalAddr() net.Addr { return nil }

func (*udpConnStub) RemoteAddr() net.Addr { return nil }

func (*udpConnStub) SetSecretKey(EncryptionMode, []byte) error { return nil }

func (*udpConnStub) SetDeadline(time.Time) error { return nil }

func (*udpConnStub) SetReadDeadline(time.Time) error { return nil }

func (*udpConnStub) SetWriteDeadline(time.Time) error { return nil }

func (*udpConnStub) Open(context.Context, string, int, uint32) (string, int, error) {
	return "", 0, nil
}

func (u *udpConnStub) Close() error {
	u.closeCalls++
	return nil
}

func (*udpConnStub) Read([]byte) (int, error) { return 0, nil }

func (*udpConnStub) ReadPacket() (*Packet, error) { return nil, nil }

func (*udpConnStub) Write([]byte) (int, error) { return 0, nil }
