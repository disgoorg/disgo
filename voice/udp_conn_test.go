package voice

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/disgoorg/godave"
	"github.com/disgoorg/snowflake/v2"
)

// passthroughEncrypter stands in for Discord's transport encryption: it returns
// the payload unchanged so a test exercises RTP framing alone.
type passthroughEncrypter struct{}

func (passthroughEncrypter) Encrypt(_ [RTPHeaderSize]byte, data []byte) []byte { return data }
func (passthroughEncrypter) Decrypt(rtpHeaderSize int, packet []byte) ([]byte, error) {
	return packet[rtpHeaderSize:], nil
}

// scriptedConn yields a fixed set of datagrams and then reports EOF.
type scriptedConn struct{ packets [][]byte }

func (c *scriptedConn) Read(b []byte) (int, error) {
	if len(c.packets) == 0 {
		return 0, io.EOF
	}
	p := c.packets[0]
	c.packets = c.packets[1:]
	return copy(b, p), nil
}
func (c *scriptedConn) Write(b []byte) (int, error)      { return len(b), nil }
func (c *scriptedConn) Close() error                     { return nil }
func (c *scriptedConn) LocalAddr() net.Addr              { return nil }
func (c *scriptedConn) RemoteAddr() net.Addr             { return nil }
func (c *scriptedConn) SetDeadline(time.Time) error      { return nil }
func (c *scriptedConn) SetReadDeadline(time.Time) error  { return nil }
func (c *scriptedConn) SetWriteDeadline(time.Time) error { return nil }

// rtpPacket builds a minimal RTP audio packet. padding bytes are appended to the
// payload and the final octet carries their count, as RFC 3550 section 5.1
// requires; the P bit is set only when padding is present.
func rtpPacket(payload []byte, padding int) []byte {
	first := byte(0x80) // version 2
	body := append([]byte(nil), payload...)
	if padding > 0 {
		first |= rtpPaddingBit
		for i := 1; i < padding; i++ {
			body = append(body, byte(padding))
		}
		body = append(body, byte(padding))
	}
	h := make([]byte, RTPHeaderSize)
	h[0] = first
	h[1] = RTPPayloadType
	binary.BigEndian.PutUint16(h[2:4], 1)
	binary.BigEndian.PutUint32(h[4:8], 960)
	binary.BigEndian.PutUint32(h[8:12], 42)
	return append(h, body...)
}

func newTestConn(packets ...[]byte) *udpConnImpl {
	return &udpConnImpl{
		config:        udpConnConfig{Logger: slog.New(slog.DiscardHandler)},
		conn:          &scriptedConn{packets: packets},
		encrypter:     passthroughEncrypter{},
		daveSession:   godave.NewNoopSession(slog.New(slog.DiscardHandler), "1", nil),
		ssrcLookup:    func(uint32) snowflake.ID { return snowflake.ID(1) },
		receiveBuffer: make([]byte, 1400),
		decryptBuffer: make([]byte, 1400),
	}
}

// Discord does set the RTP padding bit on audio packets. Leaving the padding
// attached corrupts the payload handed to the DAVE decryptor, which then rejects
// the frame, so this is the behaviour the fix exists to restore.
func TestReadPacketStripsRTPPadding(t *testing.T) {
	opus := []byte{0x78, 0x01, 0x02, 0x03, 0x04, 0x05}

	for _, tt := range []struct {
		name    string
		padding int
	}{
		{"no padding", 0},
		{"one byte", 1},
		{"small", 4},
		{"large", 32},
	} {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestConn(rtpPacket(opus, tt.padding))
			p, err := u.ReadPacket()
			if err != nil {
				t.Fatalf("ReadPacket: %v", err)
			}
			if !bytes.Equal(p.Opus, opus) {
				t.Errorf("padding not removed:\n got %v (%d bytes)\nwant %v (%d bytes)",
					p.Opus, len(p.Opus), opus, len(opus))
			}
		})
	}
}

// A packet consisting only of padding carries no media. It advances the sequence
// number without advancing the timestamp and must not reach the decryptor.
func TestReadPacketSkipsPaddingOnlyPacket(t *testing.T) {
	u := newTestConn(rtpPacket(nil, 16))
	if _, err := u.ReadPacket(); err == nil {
		t.Fatal("expected the padding-only packet to be skipped, but a packet was returned")
	}
}

// The padding bit must not be confused with the CSRC count, which is what a mask
// of 0x04 selects.
func TestRTPFirstOctetMasks(t *testing.T) {
	for _, tt := range []struct {
		name string
		mask byte
		want byte
	}{
		{"padding is bit index 2", rtpPaddingBit, 0x20},
		{"extension is bit index 3", rtpExtensionBit, 0x10},
		{"csrc count is bit index 4 to 7", rtpCSRCCountMask, 0x0F},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mask != tt.want {
				t.Errorf("mask = %#02x, want %#02x", tt.mask, tt.want)
			}
		})
	}
	if rtpPaddingBit&rtpCSRCCountMask != 0 {
		t.Error("padding bit overlaps the CSRC count nibble")
	}
}

func TestStripRTPPadding(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []byte
		want []byte
	}{
		{"empty", []byte{}, []byte{}},
		{"count of one removes only itself", []byte{1, 2, 1}, []byte{1, 2}},
		{"count of three", []byte{9, 3, 3, 3}, []byte{9}},
		{"whole payload is padding", []byte{2, 2}, []byte{}},
		{"count of zero is not padding", []byte{5, 0}, []byte{5, 0}},
		{"count beyond payload is left alone", []byte{1, 200}, []byte{1, 200}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripRTPPadding(tt.in); !bytes.Equal(got, tt.want) {
				t.Errorf("stripRTPPadding(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
