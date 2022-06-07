package voice

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrInvalidPacket    = errors.New("invalid packet")
)

type (
	UDPCreateFunc func(ip string, port int, ssrc uint32, opts ...UDPConfigOpt) *UDP
)

var _ io.Reader = (*UDP)(nil)

func NewUDP(ip string, port int, ssrc uint32, opts ...UDPConfigOpt) *UDP {
	config := DefaultUDPConfig()
	config.Apply(opts)

	return &UDP{
		config:        *config,
		ip:            ip,
		port:          port,
		ssrc:          ssrc,
		receiveBuffer: make([]byte, 1400),
	}
}

type UDP struct {
	config UDPConfig

	ip   string
	port int

	conn   net.Conn
	connMu sync.Mutex
	ssrc   uint32

	packet    [12]byte
	secretKey [32]byte

	sequence  uint16
	timestamp uint32
	nonce     [24]byte

	receiveNonce  [24]byte
	receiveBuffer []byte
	receiveOpus   []byte
}

func (c *UDP) SetSecretKey(secretKey [32]byte) {
	c.secretKey = secretKey
}

func (c *UDP) Open(ctx context.Context) (string, int, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	fmt.Printf("Opening UDP connection to: %s:%d\n", c.ip, c.port)
	var err error
	c.conn, err = c.config.Dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return "", 0, err
	}

	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, c.ssrc)
	if _, err = c.conn.Write(sb); err != nil {
		return "", 0, err
	}

	rb := make([]byte, 70)
	if _, err = c.conn.Read(rb); err != nil {
		return "", 0, err
	}

	address := rb[4:68]
	port := binary.BigEndian.Uint16(rb[68:70])

	c.packet = [12]byte{
		0: 0x80, // Version + Flags
		1: 0x78, // Payload Type
		// [2:4] // Sequence
		// [4:8] // Timestamp
	}

	binary.BigEndian.PutUint32(c.packet[8:12], c.ssrc) // SSRC

	return strings.Replace(string(address), "\x00", "", -1), int(port), nil
}

func (c *UDP) Write(b []byte) (int, error) {
	//fmt.Printf("Opus: %v\n\n", b)
	binary.BigEndian.PutUint16(c.packet[2:4], c.sequence)
	c.sequence++

	binary.BigEndian.PutUint32(c.packet[4:8], c.timestamp)
	c.timestamp += 960

	// Copy the first 12 bytes from the packet into the nonce.
	copy(c.nonce[:12], c.packet[:])

	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()
	if _, err := conn.Write(secretbox.Seal(c.packet[:12], b, &c.nonce, &c.secretKey)); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *UDP) Read(p []byte) (n int, err error) {
	packet, err := c.ReadPacket()
	if err != nil {
		return 0, err
	}
	return copy(p, packet.Opus), nil
}

const packetHeaderSize = 12

type Packet struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	Opus      []byte
}

func (c *UDP) SetDeadline(t time.Time) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	return c.conn.SetDeadline(t)
}

func (c *UDP) ReadPacket() (*Packet, error) {
	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()

	for {
		i, err := conn.Read(c.receiveBuffer)
		if err != nil {
			return nil, err
		}

		if i < 12 || (c.receiveBuffer[0] != 0x80 && c.receiveBuffer[0] != 0x90) {
			continue
		}

		copy(c.receiveNonce[:], c.receiveBuffer[0:packetHeaderSize])

		opus, ok := secretbox.Open(c.receiveOpus[:0], c.receiveBuffer[packetHeaderSize:i], &c.receiveNonce, &c.secretKey)
		if !ok {
			return nil, ErrDecryptionFailed
		}

		isExtension := c.receiveBuffer[0]&0x10 == 0x10
		isMarker := c.receiveBuffer[1]&0x80 != 0x0

		if isExtension && !isMarker {
			extLen := binary.BigEndian.Uint16(opus[2:4])
			shift := 4 + 4*int(extLen)

			if len(opus) > shift {
				opus = opus[shift:]
			}
		}
		ssrc := binary.BigEndian.Uint32(c.receiveBuffer[8:12])
		return &Packet{
			SSRC:      ssrc,
			Sequence:  binary.BigEndian.Uint16(c.receiveBuffer[2:4]),
			Timestamp: binary.BigEndian.Uint32(c.receiveBuffer[4:8]),
			Opus:      opus,
		}, nil
	}
}

func (c *UDP) Close() {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	_ = c.conn.Close()
}
