package voice

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrInvalidPacket    = errors.New("invalid packet")
)

type (
	UDPConnCreateFunc func(ip string, port int, ssrc uint32, opts ...UDPConnConfigOpt) *UDPConn
)

var _ io.Reader = (*UDPConn)(nil)

func NewUDPConn(ip string, port int, ssrc uint32, opts ...UDPConnConfigOpt) *UDPConn {
	config := DefaultUDPConnConfig()
	config.Apply(opts)

	return &UDPConn{
		config:        *config,
		ip:            ip,
		port:          port,
		ssrc:          ssrc,
		receiveBuffer: make([]byte, 1400),
	}
}

type UDPConn struct {
	config UDPConnConfig

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

func (c *UDPConn) SetSecretKey(secretKey [32]byte) {
	c.secretKey = secretKey
}

func (c *UDPConn) Open(ctx context.Context) (string, int, error) {
	fmt.Printf("Opening UDP connection to: %s:%d\n", c.ip, c.port)
	conn, err := c.config.Dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return "", 0, err
	}
	c.conn = conn

	sb := make([]byte, 70)
	binary.BigEndian.PutUint32(sb, c.ssrc)
	if _, err = conn.Write(sb); err != nil {
		return "", 0, err
	}

	rb := make([]byte, 70)
	if _, err = conn.Read(rb); err != nil {
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

func (c *UDPConn) Write(b []byte) (int, error) {
	//fmt.Printf("Opus: %v\n\n", b)
	binary.BigEndian.PutUint16(c.packet[2:4], c.sequence)
	c.sequence++

	binary.BigEndian.PutUint32(c.packet[4:8], c.timestamp)
	c.timestamp += 960

	// Copy the first 12 bytes from the packet into the nonce.
	copy(c.nonce[:12], c.packet[:])

	if _, err := c.conn.Write(secretbox.Seal(c.packet[:12], b, &c.nonce, &c.secretKey)); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *UDPConn) Read(p []byte) (n int, err error) {
	_, reader := c.ReadUser()
	return reader.Read(p)
}

const packetHeaderSize = 12

func (c *UDPConn) ReadUser() (uint32, io.Reader) {
	for {
		i, err := c.conn.Read(c.receiveBuffer)
		if err != nil {
			return 0, errorReader{err: err}
		}

		if i < 12 || (c.receiveBuffer[0] != 0x80 && c.receiveBuffer[0] != 0x90) {
			continue
		}

		copy(c.receiveNonce[:], c.receiveBuffer[0:packetHeaderSize])

		opus, ok := secretbox.Open(c.receiveOpus[:0], c.receiveBuffer[packetHeaderSize:i], &c.receiveNonce, &c.secretKey)
		if !ok {
			return 0, errorReader{err: ErrDecryptionFailed}
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
		return ssrc, bytes.NewReader(opus)
	}
}

func (c *UDPConn) Close() {
	_ = c.conn.Close()
}

type errorReader struct {
	err error
}

func (f errorReader) Read(_ []byte) (int, error) {
	return 0, f.err
}
